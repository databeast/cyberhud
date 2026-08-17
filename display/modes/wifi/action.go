package wifi

import "github.com/databeast/cyberhud/runtime/action"

// Handler implements action.Handler for the WiFi mode.
//
// Framework pattern demonstrated: command handling — input action dispatch with
// policy mutation under write lock and dirty signaling for immediate redraw.
type Handler struct{}

// HandleAction processes joystick and button inputs for the WiFi mode.
// Primary navigates to menu; Secondary forces a refresh; Left/Right cycle styles.
//
// Framework pattern demonstrated: command handling — action handler dispatching
// to policy-mutating helpers with dirty signaling for immediate redraw.
func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	switch act {
	case action.Primary:
		return action.Result{Navigate: "menu"}
	case action.Secondary:
		return action.Result{Refresh: true, Dirty: true}
	case action.Left:
		return cycleStyle(-1)
	case action.Right:
		return cycleStyle(+1)
	default:
		return action.Result{}
	}
}

// cycleStyle advances the current WiFi style by delta positions in the registry,
// wrapping around using modulo arithmetic. Returns Dirty=true when a change occurs.
func cycleStyle(delta int) action.Result {
	policyState.Lock()
	defer policyState.Unlock()

	hints, _ := getPanelHints()
	next := wifiRegistry.Cycle(policyState.policy.Style, delta, hints)
	if next == nil {
		return action.Result{}
	}
	policyState.policy.Style = next.Name()

	return action.Result{Dirty: true}
}
