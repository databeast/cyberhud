package clock

import "github.com/databeast/cyberhud/runtime/action"

// HandleAction processes joystick and button inputs for the clock mode.
// Left/Right cycle through styles; Primary toggles the border.
//
// Framework pattern demonstrated: command handling — action handler dispatching
// to policy-mutating helpers with dirty signaling for immediate redraw.
func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	switch act {
	case action.Left:
		return cycleStyle(-1)
	case action.Right:
		return cycleStyle(+1)
	default:
		return action.Result{}
	}
}

// cycleStyle advances the current style by delta positions in the registry,
// wrapping around using modulo arithmetic. Returns Dirty=true when a change
// occurs. Uses the clockRegistry's Cycle method for dispatch.
// Demonstrates the framework pattern: policy mutation under write lock with dirty signaling.
func cycleStyle(delta int) action.Result {
	policyState.Lock()
	defer policyState.Unlock()

	hints, _ := getPanelHints()
	next := clockRegistry.Cycle(policyState.policy.Style, delta, hints)
	if next == nil {
		return action.Result{}
	}
	policyState.policy.Style = next.Name()

	return action.Result{Dirty: true}
}
