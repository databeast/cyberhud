package ui

import (
	"testing"

	"github.com/databeast/cyberhud/hardware/input"
	"github.com/databeast/cyberhud/runtime/action"
)

func TestBuildInputMapper_JoyLeftMapsToActionLeft(t *testing.T) {
	mapper := BuildInputMapper(true, true, true, true)
	ev := input.Event{Key: input.JoyLeft, Type: input.Press}
	got := mapper(ev)
	if got != action.Left {
		t.Errorf("JoyLeft mapped to %v, want action.Left", got)
	}
}

func TestBuildInputMapper_JoyRightMapsToActionRight(t *testing.T) {
	mapper := BuildInputMapper(true, true, true, true)
	ev := input.Event{Key: input.JoyRight, Type: input.Press}
	got := mapper(ev)
	if got != action.Right {
		t.Errorf("JoyRight mapped to %v, want action.Right", got)
	}
}

func TestBuildInputMapper_ExistingMappingsUnchanged(t *testing.T) {
	mapper := BuildInputMapper(true, true, true, true)

	tests := []struct {
		name string
		key  input.Key
		want action.Action
	}{
		{"JoyUp maps to Up", input.JoyUp, action.Up},
		{"JoyDown maps to Down", input.JoyDown, action.Down},
		{"Key1 maps to Primary", input.Key1, action.Primary},
		{"JoyPressed maps to Secondary", input.JoyPressed, action.Secondary},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ev := input.Event{Key: tc.key, Type: input.Press}
			got := mapper(ev)
			if got != tc.want {
				t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestBuildInputMapper_Key3MapsToDown(t *testing.T) {
	mapper := BuildInputMapper(true, true, true, true)
	ev := input.Event{Key: input.Key3, Type: input.Press}
	got := mapper(ev)
	if got != action.Down {
		t.Errorf("Key3 mapped to %v, want action.Down", got)
	}
}
