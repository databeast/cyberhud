package menu

import "github.com/databeast/cyberhud/runtime/action"

// Handler implements action.Handler for the menu mode.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return HandleAction(act, cursor, itemCount)
}

// HandleAction processes an action for the menu mode.
// cursor is the current selection index; itemCount is len(Items()).
func HandleAction(a action.Action, cursor, itemCount int) action.Result {
	switch a {
	case action.Up:
		if cursor > 0 {
			return action.Result{CursorDelta: -1, Dirty: true}
		}
	case action.Down:
		if cursor < itemCount-1 {
			return action.Result{CursorDelta: +1, Dirty: true}
		}
	case action.Primary, action.Secondary:
		if cursor >= 0 && cursor < len(Destinations) {
			return action.Result{Navigate: Destinations[cursor]}
		}
	}
	return action.Result{}
}
