package tests_test

import (
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/hardware/panels"
	_ "github.com/databeast/cyberhud/hardware/panels/all"

	_ "github.com/databeast/cyberhud/display/modes/attract_matrix"
	// Blank imports to trigger mode init() self-registration.
	_ "github.com/databeast/cyberhud/display/modes/clock"
	_ "github.com/databeast/cyberhud/display/modes/cycle"
	_ "github.com/databeast/cyberhud/display/modes/dashboard"
	_ "github.com/databeast/cyberhud/display/modes/gpio"
	_ "github.com/databeast/cyberhud/display/modes/gpio_control"
	_ "github.com/databeast/cyberhud/display/modes/image"
	_ "github.com/databeast/cyberhud/display/modes/menu"
	_ "github.com/databeast/cyberhud/display/modes/serial"
	_ "github.com/databeast/cyberhud/display/modes/stemma"
	_ "github.com/databeast/cyberhud/display/modes/system"
	_ "github.com/databeast/cyberhud/display/modes/systemd"
	_ "github.com/databeast/cyberhud/display/modes/testfonts"
	_ "github.com/databeast/cyberhud/display/modes/testpattern"
	_ "github.com/databeast/cyberhud/display/modes/thermal"
	_ "github.com/databeast/cyberhud/display/modes/ticker"
	_ "github.com/databeast/cyberhud/display/modes/usb"
	_ "github.com/databeast/cyberhud/display/modes/wifi"
	_ "github.com/databeast/cyberhud/display/modes/zmq"
)

func hasMode(modes []string, want string) bool {
	for _, mode := range modes {
		if mode == want {
			return true
		}
	}
	return false
}

func TestDisplaysSingleDefaultModeWithInput(t *testing.T) {
	p, err := panels.Get("waveshare-2.2")
	if err != nil {
		t.Fatalf("Get() error=%v", err)
	}
	p.DefaultMode = ""
	p.ExcludedModes = nil
	p.InputEnabled = true
	displays := panels.Displays(p, catalog.StandardPolicy{})
	if len(displays) != 1 {
		t.Fatalf("Displays len=%d, want 1", len(displays))
	}
	if displays[0].Default != "menu" {
		t.Fatalf("Displays[0].Default=%q, want menu", displays[0].Default)
	}
	if len(displays[0].Modes) == 0 {
		t.Fatal("Displays[0].Modes should be populated")
	}
	if !hasMode(displays[0].Modes, "menu") {
		t.Fatalf("expected dynamic mode list to include menu, got %v", displays[0].Modes)
	}
}

func TestDisplaysVirtualDefaults(t *testing.T) {
	p, err := panels.Get("waveshare-triple-screen")
	if err != nil {
		t.Fatalf("Get() error=%v", err)
	}
	for i := range p.Virtual {
		p.Virtual[i].ExcludedModes = nil
		p.Virtual[i].DefaultMode = ""
		p.Virtual[i].InputEnabled = false
	}
	displays := panels.Displays(p, catalog.StandardPolicy{})
	if len(displays) != 3 {
		t.Fatalf("Displays len=%d, want 3", len(displays))
	}
	for i := range displays {
		if displays[i].Default != "dashboard" {
			t.Fatalf("displays[%d].Default=%q, want dashboard", i, displays[i].Default)
		}
		if len(displays[i].Modes) == 0 {
			t.Fatalf("displays[%d].Modes should be populated", i)
		}
	}
}

func TestDisplaysExcludedModes(t *testing.T) {
	p, err := panels.Get("waveshare-2.2")
	if err != nil {
		t.Fatalf("Get() error=%v", err)
	}
	p.ExcludedModes = []string{"menu", "ticker"}
	p.InputEnabled = true
	displays := panels.Displays(p, catalog.StandardPolicy{})
	if len(displays) != 1 {
		t.Fatalf("Displays len=%d, want 1", len(displays))
	}
	if hasMode(displays[0].Modes, "menu") || hasMode(displays[0].Modes, "ticker") {
		t.Fatalf("excluded modes leaked into mode list: %v", displays[0].Modes)
	}
	if displays[0].Default == "menu" {
		t.Fatalf("default mode should fall back when menu excluded, got %q", displays[0].Default)
	}
}

func TestDisplaysAdafruitEInkIncludesClockMode(t *testing.T) {
	p, err := panels.Get("adafruit-2.13-ssd1680")
	if err != nil {
		t.Fatalf("Get() error=%v", err)
	}
	p.InputEnabled = true
	displays := panels.Displays(p, catalog.StandardPolicy{})
	if len(displays) != 1 {
		t.Fatalf("Displays len=%d, want 1", len(displays))
	}
	if !hasMode(displays[0].Modes, "clock") {
		t.Fatalf("expected clock mode for %q, got %v", p.Name, displays[0].Modes)
	}
}
