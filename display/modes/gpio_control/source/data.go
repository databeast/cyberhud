package source

import (
	"fmt"

	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

// Data captures the state needed by all gpio_control styles for rendering.
type Data struct {
	Pins   []gpiomgr.PinState
	Cursor int
	TopRow int
}

// BuildItems returns compact display rows for GPIO control mode.
// Format: "NN OUT HI", "NN OUT LO", "NN IN" — terse, no redundant prefixes.
func BuildItems(pins []gpiomgr.PinState) []string {
	if len(pins) == 0 {
		return []string{"(no pins)"}
	}
	items := make([]string, len(pins))
	for i, p := range pins {
		switch p.Mode {
		case gpiomgr.ModeOutput:
			if p.Level {
				items[i] = fmt.Sprintf("%-2d OUT HI", p.Number)
			} else {
				items[i] = fmt.Sprintf("%-2d OUT LO", p.Number)
			}
		case gpiomgr.ModeInput:
			items[i] = fmt.Sprintf("%-2d IN", p.Number)
		default:
			items[i] = fmt.Sprintf("%-2d --", p.Number)
		}
	}
	return items
}
