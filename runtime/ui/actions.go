package ui

import (
	"github.com/databeast/cyberhud/hardware/input"
	"github.com/databeast/cyberhud/runtime/action"
)

// BuildInputMapper creates a hardware-to-action adapter using explicit pin
// roles detected in the selected panel.
func BuildInputMapper(hasKey1, hasKey2, hasKey3, hasJoyPress bool) InputMapper {
	return func(ev input.Event) Action {
		switch ev.Key {
		case input.Key2:
			if hasKey2 {
				return action.Up
			}
			return action.None
		case input.Key3:
			if hasKey3 {
				return action.Down
			}
			return action.None
		case input.JoyUp:
			return action.Up
		case input.JoyDown:
			return action.Down
		case input.JoyLeft:
			return action.Left
		case input.JoyRight:
			return action.Right
		case input.Key1:
			if hasKey1 {
				return action.Primary
			}
			return action.None
		case input.JoyPressed:
			if hasJoyPress {
				return action.Secondary
			}
			return action.None
		default:
			return action.None
		}
	}
}
