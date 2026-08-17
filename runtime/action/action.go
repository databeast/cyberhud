// Package action defines hardware-neutral logical UI actions.
// Both the ui package and individual displaymode packages can import this
// without creating circular dependencies.
package action

// Action is a logical UI input independent of physical key wiring.
type Action uint8

const (
	None      Action = iota
	Up               // scroll / cursor up
	Down             // scroll / cursor down
	Primary          // select / enter / confirm
	Secondary        // back / refresh / alternate action
	Left             // joystick left / previous
	Right            // joystick right / next
)

// Handler is implemented by every display mode that responds to input.
// cursor is the current list selection; itemCount is the total number of
// items visible in the current mode (unused by non-list modes).
type Handler interface {
	HandleAction(act Action, cursor, itemCount int) Result
}

// Result is the response a display mode returns after handling an action.
// The Region applies the result: moving the cursor, triggering a redraw, or
// navigating to a different mode.
type Result struct {
	Navigate    string // mode ID to activate ("menu", "dashboard", …) or "" to stay
	CursorDelta int    // relative cursor movement (-1 / 0 / +1)
	Dirty       bool   // request a redraw
	Refresh     bool   // force a data refresh without navigation
}

func (a Action) String() string {
	switch a {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Primary:
		return "Primary"
	case Secondary:
		return "Secondary"
	case Left:
		return "Left"
	case Right:
		return "Right"
	default:
		return "None"
	}
}
