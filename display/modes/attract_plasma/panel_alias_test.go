package attract_plasma

import (
	"testing"

	"github.com/databeast/cyberhud/display/style"
	_ "github.com/databeast/cyberhud/hardware/panels/waveshare_3_5_lcd_a"
)

func TestPanelStyleAlias(t *testing.T) {
	name, ok := style.Lookup("waveshare-3.5-lcd-a", "attract_plasma")
	if !ok || plasmaRegistry.Lookup(name) == nil {
		t.Fatalf("unresolved panel style alias %q", name)
	}
}
