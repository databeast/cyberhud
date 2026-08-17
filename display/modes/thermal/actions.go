package thermal

import "github.com/databeast/cyberhud/runtime/action"

// Handler implements action.Handler for the thermal display mode.
// Left/Right cycle the display style, Primary toggles the border, and Secondary navigates to the menu.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	switch act {
	case action.Left:
		return cycleStyle(-1)
	case action.Right:
		return cycleStyle(+1)
	case action.Primary:

	case action.Secondary:
		return action.Result{Navigate: "menu"}
	}
	return action.Result{}
}

func cycleStyle(delta int) action.Result {
	policyState.Lock()
	defer policyState.Unlock()

	hints, _ := getPanelHints()
	next := thermalRegistry.Cycle(policyState.policy.Style, delta, hints)
	if next == nil {
		return action.Result{}
	}
	policyState.policy.Style = next.Name()
	return action.Result{Dirty: true}
}
