package serial

import (
	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/runtime/action"
)

// Handler implements action.Handler for the serial monitor.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return HandleAction(act, cursor, itemCount)
}

// HandleAction supports scrolling through serial history and returning home.
func HandleAction(a action.Action, cursor, itemCount int) action.Result {
	if itemCount == 0 {
		return action.Result{}
	}
	switch a {
	case action.Primary:
		return action.Result{Navigate: "dashboard"}
	case action.Secondary:
		source.Clear()
		return action.Result{Dirty: true}
	case action.Up:
		if cursor > 0 {
			return action.Result{CursorDelta: -1, Dirty: true}
		}
	case action.Down:
		if cursor < itemCount-1 {
			return action.Result{CursorDelta: +1, Dirty: true}
		}
	}
	return action.Result{}
}
