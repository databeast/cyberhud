package tests_test

import (
	"image"
	"testing"

	"github.com/databeast/cyberhud/display/modes/wifi"
	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

func TestBuildViewWithHints_UsesMonoFast128x64OnWaveshare13OLEDHat(t *testing.T) {
	t.Cleanup(func() {
		wifi.SetPolicy(source.DefaultPolicy())
	})
	wifi.SetPolicy(source.DefaultPolicy())

	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 128, 64))
	hints.PanelProduct = "waveshare-1.3-oled-hat"
	hints.Capability = textlayout.CapMonoFast

	vd := wifi.BuildViewWithHints(hints)
	if vd.StyleReport.Name != "mono-fast-128x64" {
		t.Fatalf("StyleReport.Name = %q, want %q", vd.StyleReport.Name, "mono-fast-128x64")
	}
	if vd.Cursor != -1 {
		t.Fatalf("Cursor = %d, want -1 for non-navigable WiFi view", vd.Cursor)
	}
	if len(vd.Items) == 0 {
		t.Fatal("expected at least one rendered row")
	}
}
