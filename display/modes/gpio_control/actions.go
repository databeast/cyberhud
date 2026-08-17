package gpio_control

import (
	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
	"github.com/databeast/cyberhud/runtime/action"
)

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

// Handler implements action.Handler for GPIO control mode.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return HandleAction(act, cursor, itemCount)
}

// HandleAction processes user input for GPIO control.
func HandleAction(a action.Action, cursor, itemCount int) action.Result {
	if itemCount == 0 {
		return action.Result{}
	}
	pins := source.CurrentSnapshot()
	switch a {
	case action.Primary, action.Secondary:
		if cursor >= 0 && cursor < len(pins) {
			if pins[cursor].Mode == gpiomgr.ModeOutput {
				_ = source.Toggle(pins[cursor])
				return action.Result{Dirty: true}
			}
		}
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
