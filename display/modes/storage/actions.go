package storage

import (
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/runtime/action"
)

// Handler implements action.Handler for the storage display mode.
// Left/Right cycle the display style; Secondary navigates back to the menu.
type Handler struct{}

func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	switch act {
	case action.Left:
		return cycleStyle(-1)
	case action.Right:
		return cycleStyle(+1)
	case action.Secondary:
		return action.Result{Navigate: "menu"}
	}
	return action.Result{}
}

func cycleStyle(delta int) action.Result {
	policyState.Lock()
	defer policyState.Unlock()

	hints, _ := modehints.Current()
	next := storageRegistry.Cycle(policyState.policy.Style, delta, hints)
	if next == nil {
		return action.Result{}
	}
	policyState.policy.Style = next.Name()
	return action.Result{Dirty: true}
}
