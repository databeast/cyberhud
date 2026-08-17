package system

import "github.com/databeast/cyberhud/runtime/action"

// Handler implements action.Handler for the system info mode.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return HandleAction(act)
}

// HandleAction processes an action for the system info mode.
// Any Primary or Secondary action returns to the main menu.
func HandleAction(a action.Action) action.Result {
	switch a {
	case action.Primary, action.Secondary:
		return action.Result{Navigate: "menu"}
	}
	return action.Result{}
}
