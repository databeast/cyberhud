package schedule

import (
	"github.com/databeast/cyberhud/runtime/action"
)

// Handler implements action.Handler for the schedule mode.
//
// Schedule is a status view, not a navigable list: the only input it responds to
// is Primary (return to the dashboard). It never moves a cursor, so BuildView
// leaves Cursor at -1 and the view renders as a centered block.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	switch act {
	case action.Primary:
		return action.Result{Navigate: "dashboard"}
	default:
		return action.Result{}
	}
}
