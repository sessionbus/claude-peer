// SPDX-License-Identifier: MIT
package claude

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	kit "github.com/antst/sessionbus/bus/sdk/go"
)

// The existing worker owns the transfer manager; this control isolates the
// product's cwd binding before that caller. Full lane I/O acceptance remains a separate live check.
func TestLaneTransferRequiresTheOpenedNativeBase(t *testing.T) {
	p := New(t.TempDir())
	p.ctx, p.cancel = context.WithCancel(context.Background())
	t.Cleanup(p.cancel)
	p.opened = true
	p.stream = &stream{}
	p.SetCaller(kit.NewCaller(func(context.Context, string, any) (json.RawMessage, error) {
		t.Fatal("unexpected generic call")
		return nil, nil
	}))
	b := &laneBackend{owner: p, ctx: p.ctx}
	raw := json.RawMessage(`{"target":"peer@local","source":"notes.txt"}`)
	if _, err := b.Action(context.Background(), "transfer.offer", raw); err == nil || !strings.Contains(err.Error(), "absolute path") {
		t.Fatalf("unknown native base: %v", err)
	}
	p.cwd = "/native/lane"
	if _, err := b.Action(context.Background(), "transfer.offer", raw); !errors.Is(err, kit.ErrTransferUnsupported) {
		t.Fatalf("bound action did not reach SDK: %v", err)
	}
	if _, err := b.Action(context.Background(), "transfer.offer", json.RawMessage(`{"target":"peer@local","Source":"notes.txt"}`)); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("case-fold path bypass: %v", err)
	}
	p.closing = true
	if _, err := b.Action(context.Background(), "transfer.offer", raw); err == nil {
		t.Fatal("closed lane accepted transfer")
	}
}
