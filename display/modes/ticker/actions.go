package ticker

import "github.com/databeast/cyberhud/runtime/action"

// Handler implements action.Handler for ticker mode.
// Ticker is externally driven, so local button actions are no-ops for now.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return action.Result{}
}
