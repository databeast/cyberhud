package systemd

import "github.com/databeast/cyberhud/runtime/action"

// Handler implements action.Handler for the boot-progress view.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	switch act {
	case action.Primary, action.Secondary:
		return action.Result{Navigate: "dashboard"}
	}
	return action.Result{}
}
