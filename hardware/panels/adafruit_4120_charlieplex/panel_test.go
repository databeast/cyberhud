package adafruit_4120_charlieplex_test

import (
	"strings"
	"testing"
	"unicode"

	_ "github.com/databeast/cyberhud/hardware/driver/is31fl3731"
	_ "github.com/databeast/cyberhud/hardware/panels/adafruit_4120_charlieplex"

	"github.com/databeast/cyberhud/hardware/panels"
	"pgregory.net/rapid"
)

func TestPanelDefinitionMetadata(t *testing.T) {
	t.Run("retrievable via case-insensitive Get", func(t *testing.T) {
		variants := []string{
			"adafruit-4120-charlieplex",
			"Adafruit-4120-Charlieplex",
			"ADAFRUIT-4120-CHARLIEPLEX",
			"  adafruit-4120-charlieplex  ",
		}
		for _, name := range variants {
			def, err := panels.Get(name)
			if err != nil {
				t.Fatalf("Get(%q) returned error: %v", name, err)
			}
			if !strings.EqualFold(def.Name, "adafruit-4120-charlieplex") {
				t.Errorf("Get(%q) returned Name=%q, want %q", name, def.Name, "adafruit-4120-charlieplex")
			}
		}
	})

	t.Run("controller and dimensions", func(t *testing.T) {
		def, err := panels.Get("adafruit-4120-charlieplex")
		if err != nil {
			t.Fatalf("Get returned error: %v", err)
		}
		if def.Controller != "is31fl3731" {
			t.Errorf("Controller = %q, want %q", def.Controller, "is31fl3731")
		}
		if def.Config.Width != 16 {
			t.Errorf("Config.Width = %d, want 16", def.Config.Width)
		}
		if def.Config.Height != 8 {
			t.Errorf("Config.Height = %d, want 8", def.Config.Height)
		}
		if def.Config.Layout != "charlie-bonnet" {
			t.Errorf("Config.Layout = %q, want %q", def.Config.Layout, "charlie-bonnet")
		}
	})

	t.Run("Monochrome is false after normalization", func(t *testing.T) {
		def, err := panels.Get("adafruit-4120-charlieplex")
		if err != nil {
			t.Fatalf("Get returned error: %v", err)
		}
		if def.Monochrome != false {
			t.Errorf("Monochrome = %v, want false", def.Monochrome)
		}
	})

	t.Run("all pin fields are empty strings", func(t *testing.T) {
		def, err := panels.Get("adafruit-4120-charlieplex")
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

	t.Run("Inputs.Any() is false", func(t *testing.T) {
		def, err := panels.Get("adafruit-4120-charlieplex")
		if err != nil {
			t.Fatalf("Get returned error: %v", err)
		}
		if def.Inputs.Any() != false {
			t.Errorf("Inputs.Any() = true, want false")
		}
	})
}

func TestProperty_CaseInsensitivePanelRetrieval(t *testing.T) {

	const baseName = "adafruit-4120-charlieplex"

	rapid.Check(t, func(t *rapid.T) {
		// Build a random case-variant of the base name
		var sb strings.Builder
		for _, ch := range baseName {
			if unicode.IsLetter(ch) {
				if rapid.Bool().Draw(t, "upper") {
					sb.WriteRune(unicode.ToUpper(ch))
				} else {
					sb.WriteRune(unicode.ToLower(ch))
				}
			} else {
				sb.WriteRune(ch)
			}
		}
		variant := sb.String()

		// Optionally add leading whitespace
		leadingSpaces := rapid.IntRange(0, 3).Draw(t, "leadingSpaces")
		variant = strings.Repeat(" ", leadingSpaces) + variant

		// Optionally add trailing whitespace
		trailingSpaces := rapid.IntRange(0, 3).Draw(t, "trailingSpaces")
		variant = variant + strings.Repeat(" ", trailingSpaces)

		def, err := panels.Get(variant)
		if err != nil {
			t.Fatalf("panels.Get(%q) returned error: %v", variant, err)
		}
		if def.Name != baseName {
			t.Fatalf("panels.Get(%q).Name = %q, want %q", variant, def.Name, baseName)
		}
	})
}
