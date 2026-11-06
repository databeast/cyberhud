package tests

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/catalog"

	// Register the three modes migrated onto the uniform cmdutil.CmdHandler
	// policy path. Their init() funcs register both the command and the
	// snapshotter required for PolicyStore capture.
	_ "github.com/databeast/cyberhud/display/modes/gauges"
	_ "github.com/databeast/cyberhud/display/modes/image"
	_ "github.com/databeast/cyberhud/display/modes/ticker"
)

// TestMigratedModesCapturePolicyOnApply verifies that applying a key=value
// policy edit via each migrated mode's command handler drives the centralized
// capture chain (cmdutil.OnApplied → catalog.OnPolicyChange), which is what
// makes `freeze policy` persist the change. Before the migration these three
// modes hand-rolled their policy apply and never fired OnApplied, so freeze
// silently dropped their settings.
func TestMigratedModesCapturePolicyOnApply(t *testing.T) {
	cases := []struct {
		verb      string
		args      []string // key=value edit forwarded to Handle
		wantKey   string
		wantValue string // expected value in the captured snapshot (JSON string form)
	}{
		{verb: "image", args: []string{"fit=stretch"}, wantKey: "fit", wantValue: "stretch"},
		{verb: "ticker", args: []string{"accent=green"}, wantKey: "accent", wantValue: "green"},
		{verb: "gauges", args: []string{"shape=ring"}, wantKey: "shape", wantValue: "ring"},
	}

	for _, tc := range cases {
		t.Run(tc.verb, func(t *testing.T) {
			def, ok := catalog.Command(tc.verb)
			if !ok {
				t.Fatalf("no command registered for %q", tc.verb)
			}
			// A snapshotter must exist or the capture bridge no-ops. This is the
			// piece gauges was missing before the migration.
			if catalog.Snapshotter(tc.verb) == nil {
				t.Fatalf("mode %q has no registered snapshotter", tc.verb)
			}

			// Install a capture hook mirroring the daemon's PolicyStore wiring.
			var gotMode string
			var gotData json.RawMessage
			prev := catalog.OnPolicyChange
			t.Cleanup(func() { catalog.OnPolicyChange = prev })
			catalog.OnPolicyChange = func(modeID string, data json.RawMessage) {
				gotMode = modeID
				gotData = data
			}

			resp := def.Handle(tc.args)
			if strings.HasPrefix(resp, "ERR") {
				t.Fatalf("apply %v returned error: %s", tc.args, resp)
			}

			if gotMode != tc.verb {
				t.Fatalf("OnPolicyChange not fired for %q (got %q)", tc.verb, gotMode)
			}

			var snap map[string]any
			if err := json.Unmarshal(gotData, &snap); err != nil {
				t.Fatalf("captured snapshot is not valid JSON: %v", err)
			}
			got, ok := snap[tc.wantKey]
			if !ok {
				t.Fatalf("captured snapshot missing key %q: %v", tc.wantKey, snap)
			}
			if gotStr := valueToString(got); gotStr != tc.wantValue {
				t.Errorf("captured %s = %q, want %q", tc.wantKey, gotStr, tc.wantValue)
			}
		})
	}
}

// TestMigratedModesRejectUnknownKey verifies the uniform CmdHandler path rejects
// unknown option keys (proving edits really go through CmdHandler, not a
// hand-rolled switch).
func TestMigratedModesRejectUnknownKey(t *testing.T) {
	for _, verb := range []string{"image", "ticker", "gauges"} {
		t.Run(verb, func(t *testing.T) {
			def, ok := catalog.Command(verb)
			if !ok {
				t.Fatalf("no command registered for %q", verb)
			}
			resp := def.Handle([]string{"definitely_not_a_key=1"})
			if !strings.HasPrefix(resp, "ERR") {
				t.Fatalf("expected ERR for unknown key, got %q", resp)
			}
		})
	}
}

func valueToString(v any) string {
	switch n := v.(type) {
	case string:
		return n
	case bool:
		if n {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
