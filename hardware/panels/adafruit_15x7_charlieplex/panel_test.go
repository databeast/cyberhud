package adafruit_15x7_charlieplex_test

import (
	"strings"
	"testing"

	_ "github.com/databeast/cyberhud/hardware/driver/is31fl3731"

	"github.com/databeast/cyberhud/hardware/panels"
)

func TestPanelDefinitionMetadata(t *testing.T) {
	t.Run("retrievable via case-insensitive Get", func(t *testing.T) {
		variants := []string{
			"adafruit-15x7-charlieplex",
			"Adafruit-15x7-Charlieplex",
			"ADAFRUIT-15X7-CHARLIEPLEX",
			"  adafruit-15x7-charlieplex  ",
		}
		for _, name := range variants {
			def, err := panels.Get(name)
			if err != nil {
				t.Fatalf("Get(%q) returned error: %v", name, err)
			}
			if !strings.EqualFold(def.Name, "adafruit-15x7-charlieplex") {
				t.Errorf("Get(%q) returned Name=%q, want %q", name, def.Name, "adafruit-15x7-charlieplex")
			}
		}
	})

	t.Run("controller and dimensions", func(t *testing.T) {
		def, err := panels.Get("adafruit-15x7-charlieplex")
		if err != nil {
			t.Fatalf("Get returned error: %v", err)
		}
		if def.Controller != "is31fl3731" {
			t.Errorf("Controller = %q, want %q", def.Controller, "is31fl3731")
		}
		if def.Config.Width != 15 {
			t.Errorf("Config.Width = %d, want 15", def.Config.Width)
		}
		if def.Config.Height != 7 {
			t.Errorf("Config.Height = %d, want 7", def.Config.Height)
		}
	})

	t.Run("Monochrome is false after normalization", func(t *testing.T) {
		def, err := panels.Get("adafruit-15x7-charlieplex")
		if err != nil {
			t.Fatalf("Get returned error: %v", err)
		}
		if def.Monochrome != false {
			t.Errorf("Monochrome = %v, want false", def.Monochrome)
		}
	})

	t.Run("all pin fields are empty strings", func(t *testing.T) {
		def, err := panels.Get("adafruit-15x7-charlieplex")
		if err != nil {
			t.Fatalf("Get returned error: %v", err)
		}
		if def.DCPin != "" {
			t.Errorf("DCPin = %q, want empty string", def.DCPin)
		}
		if def.RSTPin != "" {
			t.Errorf("RSTPin = %q, want empty string", def.RSTPin)
		}
		if def.BusyPin != "" {
			t.Errorf("BusyPin = %q, want empty string", def.BusyPin)
		}
		if def.BLPin != "" {
			t.Errorf("BLPin = %q, want empty string", def.BLPin)
		}
	})
}
