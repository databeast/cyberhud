package cyberhudd

import (
	"testing"

	"github.com/databeast/cyberhud/display/modes/serial"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"
)

func TestMatrixSerialAliasesFit(t *testing.T) {
	for _, name := range []string{"adafruit-13x9-is31fl3741", "adafruit-15x7-charlieplex", "adafruit-4120-charlieplex"} {
		t.Run(name, func(t *testing.T) {
			panel, err := panels.Get(name)
			if err != nil {
				t.Fatal(err)
			}
			d, ok := driver.Get(panel.Controller)
			if !ok {
				t.Fatalf("unknown controller %q", panel.Controller)
			}
			hints := textlayout.TextHints{
				PanelProduct: name, PixelWidth: panel.Config.Width,
				PixelHeight: panel.Config.Height, Capability: d.Capability,
			}
			selected, reason := style.ResolveStyle(serial.Registry(), hints, "serial", "")
			if reason != "alias" || selected.Supports(hints) == style.Unsupported {
				t.Fatalf("serial style %q (%s) does not support panel", selected.Name(), reason)
			}
			vd := selected.Build(serial.Snapshot{Lines: []string{"test"}}, serial.DefaultPolicy(), style.NewStyleContext(hints))
			if len(vd.Sprites) != 1 || vd.Sprites[0].Label != "serial character" {
				t.Fatalf("serial style did not render a character: %#v", vd.Sprites)
			}
		})
	}
}
