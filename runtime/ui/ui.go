// Package ui provides both an interactive menu system for button-equipped
// ST7789 displays and a passive status dashboard for display-only setups.
package ui

import (
	"image/color"

	"github.com/databeast/cyberhud/hardware/input"
	"github.com/databeast/cyberhud/runtime/action"
)

// Colors used throughout the UI.
var (
	colBackground = color.RGBA{R: 0x00, G: 0x00, B: 0x1A, A: 0xFF} // dark navy
	colHighlight  = color.RGBA{R: 0x00, G: 0x7A, B: 0xCC, A: 0xFF} // bright blue
	colText       = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // white
	colBlack      = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // pure black (monochrome text on bright bg)
)

const defaultRowHeight = 10

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
