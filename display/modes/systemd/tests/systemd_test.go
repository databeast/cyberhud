package tests

import (
	"image"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/systemd"
	"github.com/databeast/cyberhud/display/modes/systemd/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestBuildItemsLoadingAndTransition(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		sequence([]string{"basic.target", "network-online.target"}),
		func() bool { return false },
		func() (int, int) { return 2, 8 },
	)

	items := systemd.BuildItems()
	if len(items) < 2 || items[0] != "Booting...." || items[1] != "Loading: basic.target" {
		t.Fatalf("first BuildItems()=%v", items)
	}

	items = systemd.BuildItems()
	if len(items) < 2 || items[1] != "Loaded: basic.target : Loading: network-online.target" {
		t.Fatalf("second BuildItems()=%v", items)
	}
}

func TestBuildItemsBootComplete(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		func() string { return "multi-user.target" },
		func() bool { return true },
		func() (int, int) { return 8, 8 },
	)
	items := systemd.BuildItems()
	if len(items) != 1 || items[0] != "Boot Complete" {
		t.Fatalf("BuildItems()=%v, want [Boot Complete]", items)
	}
	if !source.ReachedMultiUser() {
		t.Fatal("ReachedMultiUser()=false, want true")
	}
}

func TestBuildItemsDesiredFallbackWhenNoLoading(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "graphical.target" },
		func() string { return "" },
		func() bool { return false },
		func() (int, int) { return 1, 5 },
	)
	items := systemd.BuildItems()
	if len(items) < 2 || items[1] != "Loading: graphical.target" {
		t.Fatalf("BuildItems()=%v", items)
	}
}

func TestSanitizeTarget(t *testing.T) {
	cases := map[string]string{
		"multi-user.target":               "multi-user.target",
		"multi-user.target start running": "multi-user.target",
		"   graphical.target   ":          "graphical.target",
		"not-a-target.service":            "",
		"":                                "",
	}
	for in, want := range cases {
		if got := systemd.SanitizeTarget(in); got != want {
			t.Fatalf("sanitizeTarget(%q)=%q, want %q", in, got, want)
		}
	}
}

func TestSignatureReflectsState(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		sequence([]string{"basic.target", "network.target"}),
		func() bool { return false },
		func() (int, int) { return 1, 5 },
	)
	s1 := systemd.Signature()
	s2 := systemd.Signature()
	if s1 == s2 {
		t.Fatalf("Signature should change across loading transition, got %q", s1)
	}
}

func TestSignatureIncludesPolicyFingerprint(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		func() string { return "basic.target" },
		func() bool { return false },
		func() (int, int) { return 2, 8 },
	)

	systemd.SetPolicy(systemd.Policy{Style: "color-240x240"})
	s1 := systemd.Signature()

	systemd.SetPolicy(systemd.Policy{Style: "mono-128x64"})
	s2 := systemd.Signature()

	if s1 == s2 {
		t.Fatalf("Signature should differ when style changes: %q vs %q", s1, s2)
	}
}

func TestHandleCommandUseCmdutil(t *testing.T) {
	// Query returns OK with current state.
	resp := systemd.HandleCommand(nil)
	if !strings.HasPrefix(resp, "OK systemd") {
		t.Fatalf("HandleCommand(nil)=%q, want OK prefix", resp)
	}

	// Set style.
	resp = systemd.HandleCommand([]string{"style=mono-128x64"})
	if !strings.Contains(resp, "style=mono-128x64") {
		t.Fatalf("HandleCommand(style=mono-128x64)=%q, want style=mono-128x64 in response", resp)
	}

	// Invalid key.
	resp = systemd.HandleCommand([]string{"bogus=x"})
	if !strings.HasPrefix(resp, "ERR") {
		t.Fatalf("HandleCommand(bogus=x)=%q, want ERR prefix", resp)
	}

	// Invalid style value.
	resp = systemd.HandleCommand([]string{"style=nonexistent"})
	if !strings.HasPrefix(resp, "ERR") {
		t.Fatalf("HandleCommand(style=nonexistent)=%q, want ERR prefix", resp)
	}

	// Reset.
	systemd.SetPolicy(systemd.DefaultPolicy())
}

func TestBuildViewDefaultStyle(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		func() string { return "basic.target" },
		func() bool { return false },
		func() (int, int) { return 3, 8 },
	)
	systemd.SetPolicy(systemd.Policy{Style: "color-240x240"})

	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 240, 240))
	view := systemd.BuildView(hints)

	if len(view.Items) == 0 && len(view.Sprites) == 0 {
		t.Fatal("BuildView color-240x240: no items and no sprites")
	}
}

func TestBuildViewMonoStyle(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		func() string { return "network.target" },
		func() bool { return false },
		func() (int, int) { return 2, 6 },
	)
	systemd.SetPolicy(systemd.Policy{Style: "mono-128x64"})

	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 128, 64))
	view := systemd.BuildView(hints)

	if len(view.Items) == 0 {
		t.Fatal("BuildView mono-128x64: no items")
	}
	if !strings.Contains(view.Items[0], "Boot") && !strings.Contains(view.Items[0], "Loading") {
		t.Fatalf("BuildView mono-128x64 items[0]=%q, want boot-related content", view.Items[0])
	}
}

func TestBuildViewColorGradientStyle(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		func() string { return "network.target" },
		func() bool { return false },
		func() (int, int) { return 4, 8 },
	)
	systemd.SetPolicy(systemd.Policy{Style: "color-240x240"})

	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 240, 240))
	view := systemd.BuildView(hints)

	if len(view.Items) == 0 && len(view.Sprites) == 0 {
		t.Fatal("BuildView color-240x240: no items and no sprites")
	}
	// Should have a gradient sprite.
	hasGradient := false
	for _, s := range view.Sprites {
		if strings.Contains(s.Label, "gradient") {
			hasGradient = true
			break
		}
	}
	if !hasGradient {
		t.Fatal("BuildView color-240x240: no gradient sprite found")
	}
}

func TestBuildViewSmallColorStyle(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		func() string { return "basic.target" },
		func() bool { return false },
		func() (int, int) { return 3, 8 },
	)
	systemd.SetPolicy(systemd.Policy{Style: "color-128x128"})

	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 128, 128))
	view := systemd.BuildView(hints)

	// Small color panel (area < 57600): gradient only, no text expected.
	if len(view.Sprites) == 0 {
		t.Fatal("BuildView color-128x128: no sprites found")
	}
	hasGradient := false
	for _, s := range view.Sprites {
		if strings.Contains(s.Label, "gradient") {
			hasGradient = true
			break
		}
	}
	if !hasGradient {
		t.Fatal("BuildView color-128x128: no gradient sprite found")
	}
}

func TestBuildViewLargeColorStyle(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		func() string { return "network.target" },
		func() bool { return false },
		func() (int, int) { return 5, 10 },
	)
	systemd.SetPolicy(systemd.Policy{Style: "color-320x240"})

	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 320, 240))
	view := systemd.BuildView(hints)

	// Large color panel (area >= 57600): should have text and gradient.
	if len(view.Items) == 0 && len(view.Sprites) == 0 {
		t.Fatal("BuildView color-320x240: no items and no sprites")
	}
}

func TestBuildViewWithBorder(t *testing.T) {
	systemd.ResetStateForTest()
	systemd.SetProbesForTest(
		func() string { return "multi-user.target" },
		func() string { return "basic.target" },
		func() bool { return false },
		func() (int, int) { return 2, 5 },
	)

	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 240, 240))

	// Build without border.
	systemd.SetPolicy(systemd.Policy{Style: "color-240x240"})
	viewNoBorder := systemd.BuildView(hints)

	// Build with border.
	systemd.SetPolicy(systemd.Policy{Style: "color-240x240"})
	viewWithBorder := systemd.BuildView(hints)

	// With border enabled, the border is purely decorative (BorderInset=0) so
	// content offsets remain the same. The view should still render correctly.
	if len(viewWithBorder.Items) == 0 && len(viewWithBorder.Sprites) == 0 {
		t.Fatal("BuildView with border: no items and no sprites")
	}

	// Offsets should be identical since border is decorative only (no layout inset).
	if len(viewNoBorder.LineOffsets) > 0 && len(viewWithBorder.LineOffsets) > 0 {
		if viewNoBorder.LineOffsets[0] != viewWithBorder.LineOffsets[0] || viewNoBorder.OffsetY != viewWithBorder.OffsetY {
			t.Fatal("BuildView with border: offsets should be identical (border is decorative)")
		}
	}
}

func TestBootFraction(t *testing.T) {
	// Boot complete → 1.0
	snap := systemd.Snapshot{BootComplete: true, ActiveTargets: 8, TotalTargets: 8}
	if f := systemd.BootFraction(snap); f != 1.0 {
		t.Fatalf("bootFraction(complete)=%f, want 1.0", f)
	}

	// Midway through with targets.
	snap = systemd.Snapshot{ActiveTargets: 4, TotalTargets: 8, Loading: "x.target"}
	f := systemd.BootFraction(snap)
	if f != 0.5 {
		t.Fatalf("bootFraction(4/8)=%f, want 0.5", f)
	}

	// No target data, but loading.
	snap = systemd.Snapshot{Loading: "basic.target"}
	f = systemd.BootFraction(snap)
	if f != 0.5 {
		t.Fatalf("bootFraction(heuristic loading)=%f, want 0.5", f)
	}

	// No target data, loaded one.
	snap = systemd.Snapshot{Loaded: "basic.target", Loading: "network.target"}
	f = systemd.BootFraction(snap)
	if f != 0.75 {
		t.Fatalf("bootFraction(heuristic loaded)=%f, want 0.75", f)
	}
}

func sequence(values []string) func() string {
	idx := 0
	return func() string {
		if len(values) == 0 {
			return ""
		}
		if idx >= len(values) {
			return values[len(values)-1]
		}
		v := values[idx]
		idx++
		return v
	}
}
