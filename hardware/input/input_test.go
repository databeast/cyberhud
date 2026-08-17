package input_test

import (
	"testing"

	"github.com/databeast/cyberhud/hardware/input"
)

func TestKeyString(t *testing.T) {
	cases := []struct {
		key  input.Key
		want string
	}{
		{input.Key1, "KEY1"},
		{input.Key2, "KEY2"},
		{input.Key3, "KEY3"},
		{input.JoyUp, "JOY_UP"},
		{input.JoyDown, "JOY_DOWN"},
		{input.JoyLeft, "JOY_LEFT"},
		{input.JoyRight, "JOY_RIGHT"},
		{input.JoyPressed, "JOY_PRESS"},
	}
	for _, tc := range cases {
		if got := tc.key.String(); got != tc.want {
			t.Errorf("Key(%d).String() = %q, want %q", tc.key, got, tc.want)
		}
	}
}

func TestHandlerNew_nilPins(t *testing.T) {
	// Nil pins in Config must be silently ignored.
	h, err := input.New(input.Config{}, 0)
	if err != nil {
		t.Fatalf("unexpected error with all-nil config: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Handler")
	}
	// Events channel must be available.
	ch := h.Events()
	if ch == nil {
		t.Fatal("Events() must not return nil")
	}
	// Start and Stop must not panic.
	h.Start()
	h.Stop()
}
