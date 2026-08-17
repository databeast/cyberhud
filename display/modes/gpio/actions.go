package gpio

import "github.com/databeast/cyberhud/runtime/action"

// instanceHandler wraps the instance to implement action.Handler with internal
// scroll state management.
type instanceHandler struct {
	inst *instance
}

func (h instanceHandler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	// Use the instance's internal cursor for accurate state.
	result := HandleAction(act, h.inst.cursor, itemCount)
	h.inst.cursor += result.CursorDelta
	return result
}

// Handler implements action.Handler for the GPIO list mode.
// It processes logical UI input (Up, Down, Primary) for cursor navigation
// and mode switching.
//
// Framework pattern demonstrated: Action handler — a Handler struct implementing
// action.Handler for logical UI input processing.
type Handler struct{}

// HandleAction delegates to the package-level HandleAction function.
func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return HandleAction(act, cursor, itemCount)
}

// HandleAction processes an action for the GPIO list mode.
// cursor is the current selection index; itemCount is len(BuildItems(...)).
//
// Framework pattern demonstrated: Action handler — logical UI input processing
// with cursor delta and navigation results.
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
