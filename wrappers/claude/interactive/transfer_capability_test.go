// SPDX-License-Identifier: MIT
package interactive

import (
	"encoding/json"
	"testing"
)

func TestInteractiveTransferCapabilityIsNegotiated(t *testing.T) {
	o, wires := testOwner(t)
	done := report(t, o, "Stop", "native-id", "name")
	w := <-wires
	hello := w.next(t)
	var params struct {
		Info map[string]any `json:"info"`
	}
	if json.Unmarshal(hello.Params, &params) != nil || params.Info["sessionbus_transfer"] != float64(1) {
		t.Fatalf("native resident did not negotiate SDK transfers: %s", hello.Params)
	}
	w.reply(t, hello, map[string]any{})
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
