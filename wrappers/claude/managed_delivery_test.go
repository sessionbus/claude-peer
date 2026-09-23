// SPDX-License-Identifier: MIT
package claude

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"path/filepath"
	"sync"
	"testing"

	kit "github.com/antst/sessionbus/bus/sdk/go"
	"github.com/antst/sessionbus/bus/sdk/go/protocol"
	"github.com/sessionbus/peer-common/testsocket"
)

// Managed lane delivery since mandatory wake: Deliver refuses with NotRunning
// before any native write, and a retained delivery wakes a fresh run through
// the native query path. See docs/migration/MANAGED-DELIVERY-REGRESSIONS.md.

// frameRecorder is a non-blocking native stdin. A regression that writes
// fails as an extra recorded frame instead of blocking on a pipe.
type frameRecorder struct {
	mu      sync.Mutex
	buffer  bytes.Buffer
	written chan struct{}
}

func (r *frameRecorder) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, err := r.buffer.Write(p)
	r.written <- struct{}{}
	return n, err
}
func (r *frameRecorder) Close() error { return nil }
func (r *frameRecorder) frames(t *testing.T) []map[string]json.RawMessage {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()
	var frames []map[string]json.RawMessage
	decoder := json.NewDecoder(bytes.NewReader(r.buffer.Bytes()))
	for decoder.More() {
		var frame map[string]json.RawMessage
		if err := decoder.Decode(&frame); err != nil {
			t.Fatal(err)
		}
		frames = append(frames, frame)
	}
	return frames
}

func recordedStream(t *testing.T) (*stream, *frameRecorder, *json.Encoder) {
	t.Helper()
	recorder := &frameRecorder{written: make(chan struct{}, 64)}
	reader, outgoing := io.Pipe()
	s := newStream(recorder, reader, nil)
	t.Cleanup(func() { s.stop(io.EOF); _ = outgoing.Close() })
	return s, recorder, json.NewEncoder(outgoing)
}

func managedDelivery(body string) kit.DeliveryRequest {
	return kit.DeliveryRequest{MessageID: "m", Body: body, From: kit.DeliverySource{SessionID: "sender@local", Product: "claude-peer", Groups: []string{"g"}}}
}

func frameContent(t *testing.T, frame map[string]json.RawMessage) string {
	t.Helper()
	var message struct{ Content string }
	if err := json.Unmarshal(frame["message"], &message); err != nil {
		t.Fatal(err)
	}
	return message.Content
}

func TestManagedDeliveryNeverWritesIdleOrDuringActiveRun(t *testing.T) {
	for _, state := range []string{"idle", "active-before-replay", "active-after-replay"} {
		t.Run(state, func(t *testing.T) {
			s, recorder, native := recordedStream(t)
			p := New(t.TempDir())
			p.opened, p.stream, p.identity = true, s, "native-id"
			done := make(chan runResult, 1)
			admitted := make(chan struct{}, 1)
			var id string
			replay := func() {
				if err := native.Encode(map[string]any{"type": "user", "session_id": "native-id", "uuid": id, "isReplay": true}); err != nil {
					t.Fatal(err)
				}
				<-admitted
			}
			if state != "idle" {
				go func() {
					v, err := s.run(context.Background(), "native-id", "original input", func() { admitted <- struct{}{} })
					done <- runResult{v, err}
				}()
				<-recorder.written
				id = rawString(t, recorder.frames(t)[0]["uuid"])
				if state == "active-after-replay" {
					replay()
				}
			}
			receipt, err := p.Deliver(context.Background(), managedDelivery("must wake a managed run"), nil)
			var refusal *kit.ProtocolError
			if !errors.As(err, &refusal) || refusal.Code != protocol.NotRunning || receipt != (kit.DeliveryReceipt{}) {
				t.Fatalf("delivery = %+v, %v", receipt, err)
			}
			want := 0
			if state != "idle" {
				want = 1
			}
			if frames := recorder.frames(t); len(frames) != want {
				t.Fatalf("delivery wrote native input: %d frames", len(frames))
			}
			if state == "idle" {
				return
			}
			if state == "active-before-replay" {
				replay()
			}
			if err := native.Encode(map[string]any{"type": "result", "session_id": "native-id", "subtype": "success", "result": "original result", "user_message_uuids": []string{id}}); err != nil {
				t.Fatal(err)
			}
			got := <-done
			if got.err != nil || got.value.Result != "original result" || got.value.Outcome != "completed" {
				t.Fatalf("active run changed: %+v", got)
			}
			frames := recorder.frames(t)
			if len(frames) != 1 || frameContent(t, frames[0]) != "original input" || string(frames[0]["shouldQuery"]) != "true" {
				t.Fatalf("native input after the active run: %+v", frames)
			}
		})
	}
}

func TestManagedDeliveryUnavailableOrCancelledNeverWrites(t *testing.T) {
	receipt, err := New(t.TempDir()).Deliver(context.Background(), managedDelivery("x"), nil)
	if err != nil || receipt.Disposition != "rejected" || receipt.Reason != "native_unavailable" {
		t.Fatalf("unopened delivery = %+v, %v", receipt, err)
	}
	s, recorder, _ := recordedStream(t)
	p := New(t.TempDir())
	p.opened, p.stream, p.identity = true, s, "native-id"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if receipt, err := p.Deliver(ctx, managedDelivery("x"), nil); !errors.Is(err, context.Canceled) || receipt != (kit.DeliveryReceipt{}) {
		t.Fatalf("cancelled delivery = %+v, %v", receipt, err)
	}
	p.closing = true
	if receipt, err := p.Deliver(context.Background(), managedDelivery("x"), nil); err != nil || receipt.Reason != "native_unavailable" {
		t.Fatalf("closing delivery = %+v, %v", receipt, err)
	}
	if frames := recorder.frames(t); len(frames) != 0 {
		t.Fatalf("refused delivery wrote %d frames", len(frames))
	}
}

// The wake path keeps the removed append path's attempted-write rule: native
// loss after the write is uncertain, never a refusal or an admission.
func TestWakeNativeLossAfterWriteIsUncertain(t *testing.T) {
	for _, loss := range []string{"stream-stop", "native-eof"} {
		t.Run(loss, func(t *testing.T) {
			recorder := &frameRecorder{written: make(chan struct{}, 64)}
			reader, outgoing := io.Pipe()
			s := newStream(recorder, reader, nil)
			t.Cleanup(func() { s.stop(io.EOF) })
			type answer struct {
				receipt kit.DeliveryReceipt
				err     error
			}
			reported := make(chan answer, 1)
			done := make(chan error, 1)
			go func() {
				_, err := s.execute(context.Background(), "native-id", "message", func() {}, func(r kit.DeliveryReceipt, err error) error {
					reported <- answer{r, err}
					return nil
				})
				done <- err
			}()
			<-recorder.written
			if loss == "stream-stop" {
				s.stop(io.EOF)
			} else {
				_ = outgoing.Close()
			}
			got := <-reported
			var protocolError *kit.ProtocolError
			if !errors.As(got.err, &protocolError) || protocolError.Code != -32603 || string(protocolError.Data) != `"uncertain_native_admission"` || got.receipt.Disposition != "" {
				t.Fatalf("native loss after write = %+v", got)
			}
			if err := <-done; err == nil {
				t.Fatal("wake run succeeded without native admission")
			}
		})
	}
}

func TestWakeAdmissionRequiresExactUUIDAndSession(t *testing.T) {
	f := streamFixture(t)
	reported := make(chan kit.DeliveryReceipt, 1)
	done := make(chan runResult, 1)
	go func() {
		v, err := f.s.execute(context.Background(), "native-id", "message", func() {}, func(r kit.DeliveryReceipt, err error) error {
			if err != nil {
				return err
			}
			reported <- r
			return nil
		})
		done <- runResult{v, err}
	}()
	id := rawString(t, f.next(t)["uuid"])
	f.send(t, map[string]any{"type": "user", "session_id": "native-id", "uuid": "other-uuid", "isReplay": true})
	f.send(t, map[string]any{"type": "user", "session_id": "other-id", "uuid": id, "isReplay": true})
	f.barrier(t)
	select {
	case r := <-reported:
		t.Fatalf("uncorrelated replay admitted wake delivery: %+v", r)
	default:
	}
	f.send(t, map[string]any{"type": "user", "session_id": "native-id", "uuid": id, "isReplay": true})
	if r := <-reported; r.Disposition != "injected" {
		t.Fatal(r)
	}
	f.send(t, map[string]any{"type": "result", "session_id": "native-id", "subtype": "success", "result": "wake result", "user_message_uuids": []string{id}})
	if r := <-done; r.err != nil || r.value.Result != "wake result" {
		t.Fatalf("%+v", r)
	}
}

// Worker boundary for managed active delivery (installed CLW922D shape): the
// active run is untouched, Claude refuses before writing, and the daemon's
// later wake of the retained delivery is a distinct native query-path run.
func TestActualWorkerActiveDeliveryDefersToDistinctWakeRun(t *testing.T) {
	listener, err := net.Listen("unix", filepath.Join(testsocket.Directory(t), "bus.sock"))
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	t.Setenv("SESSIONBUS_SOCKET", listener.Addr().String())
	t.Setenv("SESSIONBUS_LAUNCH_TOKEN", "active-delivery-test")
	t.Setenv("SESSIONBUS_LOCAL_KEY", "")
	s, recorder, native := recordedStream(t)
	base := New(t.TempDir())
	base.ctx, base.cancel = context.WithCancel(context.Background())
	base.stream = s
	worker := kit.NewWorker(&workerStreamProduct{base})
	base.SetCaller(worker.Caller())
	base.SetShutdown(worker.Shutdown)
	served := make(chan error, 1)
	go func() { served <- worker.Serve(context.Background()) }()
	c, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close(); base.cancel(); <-served })
	wire := bufio.NewReader(c)
	next := func() protocol.Frame {
		t.Helper()
		line, err := wire.ReadBytes('\n')
		if err != nil {
			t.Fatal(err)
		}
		frame, err := protocol.DecodeFrame(line[:len(line)-1])
		if err != nil {
			t.Fatal(err)
		}
		return frame
	}
	send := func(b []byte, err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		if _, err = c.Write(b); err != nil {
			t.Fatal(err)
		}
	}
	nativeSend := func(v any) {
		t.Helper()
		if err := native.Encode(v); err != nil {
			t.Fatal(err)
		}
	}
	hello := next()
	send(protocol.ResultBytes(hello.ID, "session.hello", struct{}{}))
	send(protocol.RequestBytes(1, "session.open", kit.OpenRequest{Name: "parent/child@local", Groups: []string{}, Policy: &kit.LanePolicy{IdleMessage: "run"}}))
	_ = next()

	send(protocol.RequestBytes(2, "turn.execute", protocol.ExecuteRequest{SessionID: "native-id@local", RunID: "g/1", Input: "original input"}))
	if started := next(); started.ID != 2 || started.Error != nil {
		t.Fatal(started)
	}
	<-recorder.written
	original := rawString(t, recorder.frames(t)[0]["uuid"])
	nativeSend(map[string]any{"type": "user", "session_id": "native-id", "uuid": original, "isReplay": true})

	delivery := managedDelivery("must wake a managed run")
	send(protocol.RequestBytes(3, "message.deliver", delivery))
	if refused := next(); refused.ID != 3 || refused.Error == nil || refused.Error.Code != protocol.NotRunning {
		t.Fatal(refused)
	}
	if frames := recorder.frames(t); len(frames) != 1 {
		t.Fatalf("active delivery wrote native input: %d frames", len(frames))
	}
	nativeSend(map[string]any{"type": "result", "session_id": "native-id", "subtype": "success", "result": "original result", "user_message_uuids": []string{original}})
	ready := next()
	if ready.Method != "turn.ready" {
		t.Fatal(ready)
	}
	var first protocol.TurnReady
	if json.Unmarshal(ready.Params, &first) != nil || first.RunID != "g/1" || first.Outcome != "completed" {
		t.Fatal(ready)
	}
	send(protocol.ResultBytes(ready.ID, "turn.ready", struct{}{}))

	// The daemon, not this wrapper, retains the refused delivery and seeds g/2.
	delivery.RunID = "g/2"
	send(protocol.RequestBytes(4, "message.deliver", delivery))
	<-recorder.written
	frames := recorder.frames(t)
	if len(frames) != 2 || string(frames[1]["shouldQuery"]) != "true" {
		t.Fatalf("wake frames = %+v", frames)
	}
	var envelope struct {
		From    kit.DeliverySource `json:"from"`
		Message string             `json:"message"`
	}
	if json.Unmarshal([]byte(frameContent(t, frames[1])), &envelope) != nil || envelope.From.SessionID != "sender@local" || envelope.Message != delivery.Body {
		t.Fatal("wake lost the delivery source or body")
	}
	wake := rawString(t, frames[1]["uuid"])
	if wake == original {
		t.Fatal("wake reused the original native turn")
	}
	nativeSend(map[string]any{"type": "user", "session_id": "native-id", "uuid": wake, "isReplay": true})
	receipt := next()
	var admitted kit.DeliveryReceipt
	if receipt.ID != 4 || protocol.UnmarshalResult("message.deliver", receipt.Result, &admitted) != nil || admitted.Disposition != "injected" {
		t.Fatal(receipt)
	}
	nativeSend(map[string]any{"type": "result", "session_id": "native-id", "subtype": "success", "result": "wake result", "user_message_uuids": []string{wake}})
	ready = next()
	var second protocol.TurnReady
	if ready.Method != "turn.ready" || json.Unmarshal(ready.Params, &second) != nil || second.RunID != "g/2" || second.Outcome != "completed" {
		t.Fatal(ready)
	}
	send(protocol.ResultBytes(ready.ID, "turn.ready", struct{}{}))
}
