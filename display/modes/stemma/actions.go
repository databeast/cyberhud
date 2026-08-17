package stemma

import "github.com/databeast/cyberhud/runtime/action"

// Handler implements action.Handler for the STEMMA list mode.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return HandleAction(act, cursor, itemCount)
}

// HandleAction processes an action for the STEMMA list mode.
// cursor is the current selection index; itemCount is len(BuildItems(...)).
func HandleAction(a action.Action, cursor, itemCount int) action.Result {
	switch a {
	case action.Primary:
		return action.Result{Navigate: "dashboard"}
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
