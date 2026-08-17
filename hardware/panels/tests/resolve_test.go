package tests_test

import (
	"testing"

	"github.com/databeast/cyberhud/hardware/panels"
	_ "github.com/databeast/cyberhud/hardware/panels/all"
)

func TestResolveOverrides(t *testing.T) {
	p, err := panels.Resolve("waveshare-1.3hat", panels.Overrides{Width: 240, Height: 320, MADCTL: "0x28", XOffset: 2, YOffset: 4, DCPin: "GPIO9", RSTPin: "GPIO11", BLPin: "none"})
	if err != nil {
		t.Fatalf("Resolve() error=%v", err)
	}
	if p.Config.Width != 240 || p.Config.Height != 320 {
		t.Fatalf("unexpected size: %+v", p.Config)
	}
	if p.Config.MADCTL != 0x28 {
		t.Fatalf("unexpected MADCTL: 0x%02X", p.Config.MADCTL)
	}
	if p.Config.XOffset != 2 || p.Config.YOffset != 4 {
		t.Fatalf("unexpected offsets: %+v", p.Config)
	}
	if p.DCPin != "GPIO9" || p.RSTPin != "GPIO11" {
		t.Fatalf("unexpected pin overrides: dc=%s rst=%s", p.DCPin, p.RSTPin)
	}
	if p.BLPin != "" {
		t.Fatalf("expected BLPin disabled, got %q", p.BLPin)
	}
}
