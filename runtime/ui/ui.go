// Package ui provides both an interactive menu system for button-equipped
// ST7789 displays and a passive status dashboard for display-only setups.
//
// Default screen layout (240×240):
//
//	┌──────────────────────────┐  y=0
//	│  ▶ CYBERHUD  v1          │  title bar   (0–19)
//	├──────────────────────────┤  y=20
//	│  item 0                  │  \
//	│  item 1  ◀               │   > scrollable list (20–219)
//	│  …                       │  /
//	├──────────────────────────┤  y=220
//	│  [K1] Back [K2]▲ [K3]▼  │  hint bar    (220–239)
//	└──────────────────────────┘  y=240
package ui

import (
	"image/color"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/hardware/input"
	"github.com/databeast/cyberhud/runtime/action"
)

// Colors used throughout the UI.
var (
	colBackground = color.RGBA{R: 0x00, G: 0x00, B: 0x1A, A: 0xFF} // dark navy
	colTitleBar   = color.RGBA{R: 0x00, G: 0x4E, B: 0x8A, A: 0xFF} // cobalt blue
	colHintBar    = color.RGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF} // dark grey
	colHighlight  = color.RGBA{R: 0x00, G: 0x7A, B: 0xCC, A: 0xFF} // bright blue
	colText       = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // white
	colSubText    = color.RGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF} // grey
	colPresent    = color.RGBA{R: 0x00, G: 0xFF, B: 0x88, A: 0xFF} // green
	colAbsent     = color.RGBA{R: 0xFF, G: 0x44, B: 0x44, A: 0xFF} // red
	colBlack      = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // pure black (monochrome text on bright bg)
)

const (
	// Default bar heights for large panels (≥80px tall).
	defaultTitleBarHeight = 20
	defaultHintBarHeight  = 20
	// Compact bar heights for small OLED panels (<80px tall).
	compactTitleBarHeight = 10
	compactHintBarHeight  = 10

	defaultRowHeight = 10
)

// Action is a logical UI input independent of physical key wiring.
// Re-exported from internal/action for API consistency.
type Action = action.Action

// Logical action constants (see internal/action for definitions).
const (
	ActionNone      = action.None
	ActionUp        = action.Up
	ActionDown      = action.Down
	ActionPrimary   = action.Primary
	ActionSecondary = action.Secondary
)

// InputMapper converts a raw input event into a logical action.
type InputMapper func(input.Event) Action

func defaultInputMapper(ev input.Event) Action {
	switch ev.Key {
	case input.Key2, input.JoyUp:
		return action.Up
	case input.Key3, input.JoyDown:
		return action.Down
	case input.Key1:
		return action.Primary
	case input.JoyPressed:
		return action.Secondary
	default:
		return action.None
	}
}

func truncateHint(s string, max int) string {
	return textlayout.Truncate(s, max)
}
