package zmq

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/runtime/action"
)

// Handler implements the action.Handler interface for the ZMQ display mode.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return HandleAction(act, cursor, itemCount)
}

// HandleAction supports scrolling through ZMQ message history and returning home.
func HandleAction(a action.Action, cursor, itemCount int) action.Result {
	if itemCount == 0 {
		return action.Result{}
	}
	switch a {
	case action.Primary:
		return action.Result{Navigate: "dashboard"}
	case action.Secondary:
		content.Clear()
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
