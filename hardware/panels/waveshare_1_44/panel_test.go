package waveshare_1_44_test

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/hardware/driver"
	"github.com/databeast/cyberhud/hardware/panels"

	// Side-effect imports to trigger init() registration.
	_ "github.com/databeast/cyberhud/hardware/driver/st7735s"
	_ "github.com/databeast/cyberhud/hardware/panels/waveshare_1_44"

	// Blank imports to trigger mode init() self-registration.
	_ "github.com/databeast/cyberhud/display/modes/attract_matrix"
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

func TestWaveshare144NoExcludedModes(t *testing.T) {
	def, err := panels.Get("waveshare-1.44")
	if err != nil {
		t.Fatalf("panels.Get(\"waveshare-1.44\") error: %v", err)
	}

	if len(def.ExcludedModes) != 0 {
		t.Errorf("ExcludedModes = %v, want nil or empty", def.ExcludedModes)
	}

	// Verify Displays returns all registered catalog modes.
	def.InputEnabled = true
	displays := panels.Displays(def, catalog.StandardPolicy{})
	if len(displays) != 1 {
		t.Fatalf("Displays() returned %d panels, want 1", len(displays))
	}

	allModes := panels.AllModes()
	if len(allModes) == 0 {
		t.Fatal("AllModes() returned empty; expected registered modes")
	}

	got := displays[0].Modes
	if len(got) != len(allModes) {
		t.Errorf("Displays()[0].Modes has %d modes, want %d (all registered modes)\ngot:  %v\nwant: %v",
			len(got), len(allModes), got, allModes)
	}

	// Verify every registered mode is present in the output.
	modeSet := make(map[string]struct{}, len(got))
	for _, m := range got {
		modeSet[m] = struct{}{}
	}
	for _, want := range allModes {
		if _, ok := modeSet[want]; !ok {
			t.Errorf("mode %q missing from Displays output", want)
		}
	}
}

func TestPanelRegistration(t *testing.T) {
	_, err := panels.Get("waveshare-1.44")
	if err != nil {
		t.Fatalf("panels.Get(\"waveshare-1.44\") error: %v", err)
	}

	names := panels.Names()
	found := false
	for _, n := range names {
		if n == "waveshare-1.44" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("panels.Names() does not include \"waveshare-1.44\", got %v", names)
	}
}

func TestPanelDefinitionFields(t *testing.T) {
	def, err := panels.Get("waveshare-1.44")
	if err != nil {
		t.Fatalf("panels.Get error: %v", err)
	}

	if def.Controller != "st7735s" {
		t.Errorf("Controller = %q, want %q", def.Controller, "st7735s")
	}
	if def.Config.Width != 128 {
		t.Errorf("Config.Width = %d, want 128", def.Config.Width)
	}
	if def.Config.Height != 128 {
		t.Errorf("Config.Height = %d, want 128", def.Config.Height)
	}
	if def.Config.XOffset != 2 {
		t.Errorf("Config.XOffset = %d, want 2", def.Config.XOffset)
	}
	if def.Config.YOffset != 1 {
		t.Errorf("Config.YOffset = %d, want 1", def.Config.YOffset)
	}

	wantMADCTL := driver.MadctlMY | driver.MadctlMV | driver.MadctlBGR
	if def.Config.MADCTL != wantMADCTL {
		t.Errorf("Config.MADCTL = 0x%02X, want 0x%02X", def.Config.MADCTL, wantMADCTL)
	}

	if def.DCPin != panels.GPIO25 {
		t.Errorf("DCPin = %q, want %q", def.DCPin, panels.GPIO25)
	}
	if def.RSTPin != panels.GPIO27 {
		t.Errorf("RSTPin = %q, want %q", def.RSTPin, panels.GPIO27)
	}
	if def.BLPin != panels.GPIO24 {
		t.Errorf("BLPin = %q, want %q", def.BLPin, panels.GPIO24)
	}
	if def.Monochrome {
		t.Errorf("Monochrome = true, want false")
	}
}

func TestInputPinAssignments(t *testing.T) {
	def, err := panels.Get("waveshare-1.44")
	if err != nil {
		t.Fatalf("panels.Get error: %v", err)
	}

	checks := []struct {
		name string
		got  string
		want string
	}{
		{"Key1", def.Inputs.Key1, panels.GPIO21},
		{"Key2", def.Inputs.Key2, panels.GPIO20},
		{"Key3", def.Inputs.Key3, panels.GPIO16},
		{"JoyUp", def.Inputs.JoyUp, panels.GPIO6},
		{"JoyDown", def.Inputs.JoyDown, panels.GPIO19},
		{"JoyLeft", def.Inputs.JoyLeft, panels.GPIO5},
		{"JoyRight", def.Inputs.JoyRight, panels.GPIO26},
		{"JoyPressed", def.Inputs.JoyPressed, panels.GPIO13},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("Inputs.%s = %q, want %q", c.name, c.got, c.want)
		}
	}

	if !def.Inputs.Any() {
		t.Errorf("Inputs.Any() = false, want true")
	}
}

func TestPinConflictReporting(t *testing.T) {
	def, err := panels.Get("waveshare-1.44")
	if err != nil {
		t.Fatalf("panels.Get error: %v", err)
	}

	notices := panels.PinNotices(def)
	if len(notices) != 1 {
		t.Fatalf("PinNotices() returned %d notices, want 1; got %v", len(notices), notices)
	}

	notice := notices[0]
	if !strings.Contains(notice, "GPIO13") {
		t.Errorf("notice does not contain \"GPIO13\": %q", notice)
	}
	if !strings.Contains(notice, "unavailable") {
		t.Errorf("notice does not contain \"unavailable\": %q", notice)
	}
	if !strings.Contains(notice, "input_press") {
		t.Errorf("notice does not contain \"input_press\": %q", notice)
	}

	report := panels.BuildPinReport(def)

	// GPIO13 connector should show conflict with input_press.
	if !strings.Contains(report, "3-pin connector GPIO13") {
		t.Errorf("BuildPinReport missing \"3-pin connector GPIO13\"")
	}
	if !strings.Contains(report, "conflict") {
		t.Errorf("BuildPinReport missing \"conflict\" status")
	}
	if !strings.Contains(report, "input_press") {
		t.Errorf("BuildPinReport missing \"input_press\" conflict reference")
	}

	// GPIO18 connector should be free.
	if !strings.Contains(report, "3-pin connector GPIO18") {
		t.Errorf("BuildPinReport missing \"3-pin connector GPIO18\"")
	}

	// Find the GPIO18 line and verify it says "free".
	for _, line := range strings.Split(report, "\n") {
		if strings.Contains(line, "3-pin connector GPIO18") {
			if !strings.Contains(line, "free") {
				t.Errorf("GPIO18 connector line does not contain \"free\": %q", line)
			}
			break
		}
	}
}
