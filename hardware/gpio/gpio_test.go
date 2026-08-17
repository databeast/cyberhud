package gpio_test

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/hardware/gpio"
)

func TestManagerSnapshot_noHardware(t *testing.T) {
	// On a non-RPi host periph will register no pins.  Manager must still
	// initialise without panicking.
	m := gpio.New()
	if m == nil {
		t.Fatal("expected non-nil Manager")
	}
	snap := m.Snapshot()
	// Pins may be empty on a non-RPi host; the slice must at least be non-nil.
	if snap == nil {
		t.Fatal("Snapshot must not return nil")
	}
}

func TestManagerSetOutput_unknownPin(t *testing.T) {
	m := gpio.New()
	err := m.SetOutput(999, false)
	if err == nil {
		t.Fatal("expected error for unknown pin")
	}
}

func TestManagerSetOutput_nilCapability(t *testing.T) {
	// On a non-RPi host all managed pins have out == nil.
	// SetOutput must return a descriptive error with the pin number.
	m := gpio.New()
	err := m.SetOutput(4, false)
	if err == nil {
		t.Skip("pin 4 has output capability on this host")
	}
	if !strings.Contains(err.Error(), "4") {
		t.Errorf("error %q must contain pin number", err)
	}
	if !strings.Contains(err.Error(), "output") {
		t.Errorf("error %q must mention unsupported capability", err)
	}
}

func TestManagerSetInput_unknownPin(t *testing.T) {
	m := gpio.New()
	err := m.SetInput(999, 0)
	if err == nil {
		t.Fatal("expected error for unknown pin")
	}
}

func TestManagerSetInput_nilCapability(t *testing.T) {
	// On a non-RPi host all managed pins have in == nil.
	// SetInput must return a descriptive error with the pin number.
	m := gpio.New()
	err := m.SetInput(4, 0)
	if err == nil {
		t.Skip("pin 4 has input capability on this host")
	}
	if !strings.Contains(err.Error(), "4") {
		t.Errorf("error %q must contain pin number", err)
	}
	if !strings.Contains(err.Error(), "input") {
		t.Errorf("error %q must mention unsupported capability", err)
	}
}

func TestManagerRead_unknownPin(t *testing.T) {
	m := gpio.New()
	_, err := m.Read(999)
	if err == nil {
		t.Fatal("expected error for unknown pin")
	}
}

func TestPinStateString(t *testing.T) {
	p := gpio.PinState{Number: 4, Name: "GPIO4", Mode: gpio.ModeOutput, Level: true}
	s := p.String()
	if s == "" {
		t.Fatal("PinState.String() must not be empty")
	}
}

func TestPinModeString(t *testing.T) {
	cases := []struct {
		mode gpio.PinMode
		want string
	}{
		{gpio.ModeInput, "IN"},
		{gpio.ModeOutput, "OUT"},
		{gpio.ModeAlt, "ALT"},
		{gpio.ModeUnknown, "???"},
	}
	for _, tc := range cases {
		if got := tc.mode.String(); got != tc.want {
			t.Errorf("PinMode(%d).String() = %q, want %q", tc.mode, got, tc.want)
		}
	}
}

func TestPinStateStringUnknownMode(t *testing.T) {
	p := gpio.PinState{Number: 4, Name: "GPIO4", Mode: gpio.ModeUnknown, Level: false}
	s := p.String()
	if !strings.Contains(s, "--") {
		t.Fatalf("PinState.String()=%q, want '--' marker for unknown mode", s)
	}
	if strings.Contains(s, " LO") || strings.Contains(s, " HI") {
		t.Fatalf("PinState.String()=%q, unknown mode should not imply HI/LO", s)
	}
}
