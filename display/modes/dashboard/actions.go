package dashboard

import "github.com/databeast/cyberhud/runtime/action"

// Handler implements action.Handler for the dashboard mode.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return HandleAction(act)
}

// HandleAction processes an action for the dashboard mode.
func HandleAction(a action.Action) action.Result {
	switch a {
	case action.Primary:
		return action.Result{Navigate: "menu"}
	case action.Secondary:
		return action.Result{Refresh: true, Dirty: true}
	}
	return action.Result{}
}
