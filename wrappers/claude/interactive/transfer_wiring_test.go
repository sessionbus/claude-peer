// SPDX-License-Identifier: MIT
package interactive

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	kit "github.com/antst/sessionbus/bus/sdk/go"
	"github.com/antst/sessionbus/bus/sdk/go/protocol"
)

// Reuse the production-codec owner fixture. Only the daemon and native delivery
// callback are simulated; paths and prepared handles use the SDK filesystem.
func transferOwner(t *testing.T, cwd string) (*Owner, *wire, <-chan *wire) {
	t.Helper()
	o, wires := testOwner(t)
	o.nativeSocket = "/native.sock"
	o.registryRead = func() (startupIdentity, error) {
		return startupIdentity{SessionID: "native-id", Socket: "/native.sock", Cwd: cwd}, nil
	}
	done := report(t, o, "Stop", "native-id", "name")
	w := <-wires
	w.reply(t, w.next(t), map[string]any{"capabilities": []string{protocol.TransferCapability}})
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	return o, w, wires
}
func transferDescriptor(id, state string) map[string]any {
	return map[string]any{"id": id, "sender": "native-id", "receiver": "peer@local", "name": "file", "created_by": "native-id", "state": state, "bytes": 0}
}
func ownerOffer(t *testing.T, o *Owner, w *wire, source, id string) {
	t.Helper()
	raw, _ := json.Marshal(map[string]string{"target": "peer@local", "source": source})
	done := make(chan error, 1)
	go func() { _, err := o.Action(context.Background(), "transfer.offer", raw); done <- err }()
	f := w.next(t)
	var offer map[string]any
	if f.Method != "transfer.offer" || json.Unmarshal(f.Params, &offer) != nil || len(offer) != 2 || offer["target"] != "peer@local" || offer["name"] != filepath.Base(source) {
		t.Fatalf("offer leaked/changed local binding: %s %s", f.Method, f.Params)
	}
	w.reply(t, f, transferDescriptor(id, "offered"))
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
func transferRequest(t *testing.T, w *wire, method string, params any) protocol.Frame {
	t.Helper()
	w.request(t, method, params)
	f := w.next(t)
	if f.ID != w.nextID || f.Error != nil {
		t.Fatalf("%s reply: %+v", method, f)
	}
	if _, err := protocol.DecodeResult(method, f.Result); err != nil {
		t.Fatalf("%s invalid result: %s %v", method, f.Result, err)
	}
	return f
}
func TestInteractiveTransferOfferPreparesTheNativeDirectory(t *testing.T) {
	cwd := t.TempDir()
	body := []byte("the native session file")
	if err := os.WriteFile(filepath.Join(cwd, "notes.txt"), body, 0600); err != nil {
		t.Fatal(err)
	}
	o, w, _ := transferOwner(t, cwd)
	ownerOffer(t, o, w, "notes.txt", "source-one")
	f := transferRequest(t, w, "transfer.prepare", map[string]any{"id": "source-one", "role": "source"})
	var prepared protocol.TransferPrepared
	if protocol.UnmarshalResult("transfer.prepare", f.Result, &prepared) != nil || prepared.Size == nil || int(*prepared.Size) != len(body) {
		t.Fatalf("wrong native source: %s", f.Result)
	}
	o.End() // joins the SDK source before this owner is disposed.
}
func TestInteractiveTransferUnknownOrStaleCwdDoesNotBorrowTheHelperDirectory(t *testing.T) {
	o, w, _ := transferOwner(t, "/known/native")
	for _, tc := range []struct {
		name string
		row  startupIdentity
		err  error
	}{
		{name: "missing"},
		{name: "relative", row: startupIdentity{SessionID: "native-id", Socket: "/native.sock", Cwd: "relative"}},
		{name: "foreign session", row: startupIdentity{SessionID: "other", Socket: "/native.sock", Cwd: "/foreign"}},
		{name: "foreign socket", row: startupIdentity{SessionID: "native-id", Socket: "/other.sock", Cwd: "/foreign"}},
		{name: "partial read", err: errors.New("incomplete registry")},
		{name: "ended parent", err: errStartupParentEnded},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o.mu.Lock()
			o.registryRead = func() (startupIdentity, error) { return tc.row, tc.err }
			o.mu.Unlock()
			done := make(chan error, 1)
			go func() {
				_, err := o.Action(context.Background(), "transfer.offer", json.RawMessage(`{"target":"peer@local","source":"relative.txt"}`))
				done <- err
			}()
			select {
			case f := <-w.frames:
				t.Fatalf("unbound local path reached wire: %s %s", f.Method, f.Params)
			case err := <-done:
				if err == nil {
					t.Fatal("unknown/stale native cwd accepted")
				}
			}
		})
	}
	// A missing registry does not invent a scope restriction: absolute paths work.
	o.mu.Lock()
	o.registryRead = nil
	o.mu.Unlock()
	absolute := filepath.Join(t.TempDir(), "absolute.txt")
	if err := os.WriteFile(absolute, []byte("absolute"), 0600); err != nil {
		t.Fatal(err)
	}
	ownerOffer(t, o, w, absolute, "absolute-one")
	f := transferRequest(t, w, "transfer.prepare", map[string]any{"id": "absolute-one", "role": "source"})
	var prepared protocol.TransferPrepared
	if protocol.UnmarshalResult("transfer.prepare", f.Result, &prepared) != nil || prepared.Size == nil || *prepared.Size != 8 {
		t.Fatalf("absolute source: %s", f.Result)
	}
}
func TestInteractiveTransferFinalProjectionLeavesTheReaderAvailable(t *testing.T) {
	cwd := t.TempDir()
	o, w, _ := transferOwner(t, cwd)
	deliveries := make(chan kit.DeliveryRequest, 4)
	o.deliver = func(r Recipient, d kit.DeliveryRequest) (kit.DeliveryReceipt, error) {
		if r.SessionID != "native-id" {
			t.Errorf("wrong native owner: %+v", r)
		}
		deliveries <- d
		return kit.DeliveryReceipt{Disposition: "written"}, nil
	}
	body := []byte("received bytes")
	digest := sha256.Sum256(body)
	sum := hex.EncodeToString(digest[:])
	peer := kit.DeliverySource{SessionID: "peer@local", Product: "fixture", Groups: []string{}}
	transferRequest(t, w, "transfer.notice", map[string]any{"kind": "offer", "id": "received", "name": "out.bin", "from": peer, "size": len(body)})
	done := make(chan error, 1)
	go func() {
		_, err := o.Action(context.Background(), "transfer.accept", json.RawMessage(`{"id":"received","destination":"out.bin"}`))
		done <- err
	}()
	accepted := w.next(t)
	if accepted.Method != "transfer.accept" || string(accepted.Params) != `{"id":"received"}` {
		t.Fatalf("accept leaked path: %s %s", accepted.Method, accepted.Params)
	}
	descriptor := transferDescriptor("received", "preparing")
	descriptor["sender"] = "peer@local"
	descriptor["receiver"] = "native-id"
	descriptor["created_by"] = "peer@local"
	w.reply(t, accepted, descriptor)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	for _, step := range []struct {
		method string
		params map[string]any
	}{
		{"transfer.prepare", map[string]any{"id": "received", "role": "destination", "size": len(body)}},
		{"transfer.start", map[string]any{"id": "received", "role": "destination", "size": len(body)}},
		{"transfer.chunk", map[string]any{"id": "received", "offset": 0, "data": base64.StdEncoding.EncodeToString(body)}},
		{"transfer.finish", map[string]any{"id": "received", "bytes": len(body), "sha256": sum}},
	} {
		transferRequest(t, w, step.method, step.params)
	}
	path := filepath.Join(cwd, "out.bin")
	if got, err := os.ReadFile(path); err != nil || string(got) != string(body) {
		t.Fatalf("native destination: %q %v", got, err)
	}
	terminal := map[string]any{"kind": "terminal", "id": "received", "outcome": "completed", "bytes": len(body), "sha256": sum}
	raw, _ := json.Marshal(terminal)
	// A broken on-reader projection fails this finite write, not a longer timeout.
	if err := w.fd.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatal(err)
	}
	w.request(t, "message.deliver", kit.DeliveryRequest{MessageID: "transfer-notice-received", From: peer, Body: string(raw)})
	terminalID := w.nextID
	w.request(t, "message.deliver", kit.DeliveryRequest{MessageID: "later", From: peer, Body: "later"})
	laterID := w.nextID
	first := w.next(t)
	if first.ID != laterID || first.Error != nil {
		t.Fatalf("reader blocked on terminal projection: %+v", first)
	}
	w.request(t, "transfer.notice", terminal)
	machineID := w.nextID
	replies := map[int64]protocol.Frame{}
	for len(replies) < 2 {
		f := w.next(t)
		replies[f.ID] = f
	}
	for _, id := range []int64{terminalID, machineID} {
		if f, ok := replies[id]; !ok || f.Error != nil {
			t.Fatalf("terminal reply %d: %+v", id, f)
		}
	}
	firstDelivery, lastDelivery := <-deliveries, <-deliveries
	if firstDelivery.MessageID != "later" || lastDelivery.MessageID != "transfer-notice-received" || !strings.Contains(lastDelivery.Body, `"destination":"`+path+`"`) {
		t.Fatalf("native notice projection/order: %+v %+v", firstDelivery, lastDelivery)
	}
}
func TestInteractiveTransferEndJoinsUnpublishedDestination(t *testing.T) {
	cwd := t.TempDir()
	o, w, _ := transferOwner(t, cwd)
	peer := kit.DeliverySource{SessionID: "peer@local", Product: "fixture", Groups: []string{}}
	transferRequest(t, w, "transfer.notice", map[string]any{"kind": "offer", "id": "pending", "name": "out.bin", "from": peer})
	done := make(chan error, 1)
	go func() {
		_, err := o.Action(context.Background(), "transfer.accept", json.RawMessage(`{"id":"pending","destination":"out.bin"}`))
		done <- err
	}()
	request := w.next(t)
	descriptor := transferDescriptor("pending", "preparing")
	descriptor["sender"] = "peer@local"
	descriptor["receiver"] = "native-id"
	descriptor["created_by"] = "peer@local"
	w.reply(t, request, descriptor)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	transferRequest(t, w, "transfer.prepare", map[string]any{"id": "pending", "role": "destination", "size": 1})
	files, err := os.ReadDir(cwd)
	if err != nil || len(files) != 1 || !strings.HasPrefix(files[0].Name(), ".sessionbus-transfer-") {
		t.Fatalf("expected owned pending temp: %v %v", files, err)
	}
	o.End()
	files, err = os.ReadDir(cwd)
	if err != nil || len(files) != 0 {
		t.Fatalf("owner returned before transfer cleanup: %v %v", files, err)
	}
}

func TestInteractiveTransferRegistryLookupCannotReviveEndedResident(t *testing.T) {
	o, w, wires := transferOwner(t, t.TempDir())
	entered, release := make(chan struct{}), make(chan struct{})
	o.mu.Lock()
	o.registryRead = func() (startupIdentity, error) {
		close(entered)
		<-release
		return startupIdentity{SessionID: "native-id", Socket: "/native.sock", Cwd: "/old/native"}, nil
	}
	o.mu.Unlock()
	done := make(chan error, 1)
	go func() {
		_, err := o.Action(context.Background(), "transfer.offer", json.RawMessage(`{"target":"peer@local","source":"a.txt"}`))
		done <- err
	}()
	<-entered
	report(t, o, "SessionEnd", "native-id", "")
	close(release)
	if err := <-done; err == nil {
		t.Fatal("ended resident accepted held registry result")
	}
	for f := range w.frames {
		if f.Method == "transfer.offer" {
			t.Fatal("stale transfer escaped after SessionEnd")
		}
	}
	o.mu.Lock()
	o.registryRead = nil
	o.mu.Unlock()
	// Ending the captured identity must not disable a later native session.
	observed := report(t, o, "Stop", "replacement", "replacement")
	next := <-wires
	next.reply(t, next.next(t), map[string]any{"capabilities": []string{protocol.TransferCapability}})
	if err := <-observed; err != nil {
		t.Fatal(err)
	}
	listed := make(chan error, 1)
	go func() { _, err := o.Action(context.Background(), "list", json.RawMessage(`{}`)); listed <- err }()
	f := next.next(t)
	if f.Method != "session.list" {
		t.Fatal(f.Method)
	}
	next.reply(t, f, map[string]any{"sessions": []any{}})
	if err := <-listed; err != nil {
		t.Fatal(err)
	}
}
