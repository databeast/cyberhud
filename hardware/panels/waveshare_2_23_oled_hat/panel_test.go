package waveshare_2_23_oled_hat_test

import (
	"testing"

	"github.com/databeast/cyberhud/hardware/panels"

	// Side-effect imports to trigger init() registration.
	_ "github.com/databeast/cyberhud/hardware/driver/ssd1305"
	_ "github.com/databeast/cyberhud/hardware/panels/waveshare_2_23_oled_hat"
)

func TestI2CPanelDefinition(t *testing.T) {
	def, err := panels.Get("waveshare-2.23-oled-hat-i2c")
	if err != nil {
		t.Fatalf("panels.Get(i2c) error: %v", err)
	}

	if def.Name != "waveshare-2.23-oled-hat-i2c" {
		t.Errorf("Name = %q, want %q", def.Name, "waveshare-2.23-oled-hat-i2c")
	}
	if def.Controller != "ssd1305" {
		t.Errorf("Controller = %q, want %q", def.Controller, "ssd1305")
	}
	if !def.Monochrome {
		t.Errorf("Monochrome = false, want true")
	}
	if def.Config.Width != 128 {
		t.Errorf("Config.Width = %d, want 128", def.Config.Width)
	}
	if def.Config.Height != 32 {
		t.Errorf("Config.Height = %d, want 32", def.Config.Height)
	}
	if def.Config.I2CAddr != 0x3C {
		t.Errorf("Config.I2CAddr = 0x%02X, want 0x3C", def.Config.I2CAddr)
	}
	if def.I2CBus != "/dev/i2c-1" {
		t.Errorf("I2CBus = %q, want %q", def.I2CBus, "/dev/i2c-1")
	}
	if def.DCPin != "" {
		t.Errorf("DCPin = %q, want empty", def.DCPin)
	}
	if def.RSTPin != "" {
		t.Errorf("RSTPin = %q, want empty", def.RSTPin)
	}
	if def.BLPin != "" {
		t.Errorf("BLPin = %q, want empty", def.BLPin)
	}
	if def.BusyPin != "" {
		t.Errorf("BusyPin = %q, want empty", def.BusyPin)
	}
}

func TestSPIPanelDefinition(t *testing.T) {
	def, err := panels.Get("waveshare-2.23-oled-hat-spi")
	if err != nil {
		t.Fatalf("panels.Get(spi) error: %v", err)
	}

	if def.Name != "waveshare-2.23-oled-hat-spi" {
		t.Errorf("Name = %q, want %q", def.Name, "waveshare-2.23-oled-hat-spi")
	}
	if def.Controller != "ssd1305" {
		t.Errorf("Controller = %q, want %q", def.Controller, "ssd1305")
	}
	if !def.Monochrome {
		t.Errorf("Monochrome = false, want true")
	}
	if def.Config.Width != 128 {
		t.Errorf("Config.Width = %d, want 128", def.Config.Width)
	}
	if def.Config.Height != 32 {
		t.Errorf("Config.Height = %d, want 32", def.Config.Height)
	}
	if def.DCPin != panels.GPIO24 {
		t.Errorf("DCPin = %q, want %q", def.DCPin, panels.GPIO24)
	}
	if def.RSTPin != panels.GPIO25 {
		t.Errorf("RSTPin = %q, want %q", def.RSTPin, panels.GPIO25)
	}
	if def.BLPin != panels.GPIO18 {
		t.Errorf("BLPin = %q, want %q", def.BLPin, panels.GPIO18)
	}
}

func TestPanelNamesIncludesBothVariants(t *testing.T) {
	names := panels.Names()

	found := map[string]bool{
		"waveshare-2.23-oled-hat-i2c": false,
		"waveshare-2.23-oled-hat-spi": false,
	}

	for _, n := range names {
		if _, ok := found[n]; ok {
			found[n] = true
		}
	}

	for name, present := range found {
		if !present {
			t.Errorf("panels.Names() missing %q", name)
		}
	}
}
