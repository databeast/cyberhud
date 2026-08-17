package ui

import (
	"log"

	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/hardware/input"
	"github.com/databeast/cyberhud/runtime/action"
)

// RegionInputDispatcher implements region.InputDispatcher. It routes input
// events to the focused Region's active ModeInstance via its ActionHandler.
type RegionInputDispatcher struct {
	inputMapper InputMapper
}

// NewRegionInputDispatcher creates a RegionInputDispatcher with the given
// input mapper. The renderer parameter is accepted for backward compatibility
// but is no longer used — action dispatch goes directly through Region.Instance().
func NewRegionInputDispatcher(renderer *RegionRenderer, inputMapper InputMapper) *RegionInputDispatcher {
	return &RegionInputDispatcher{
		inputMapper: inputMapper,
	}
}

// Dispatch implements region.InputDispatcher. It routes a single input event
// to the specified Region's active ModeInstance:
//  1. Filters out non-press events (release) without dispatching.
//  2. Maps the raw event to a logical action via InputMapper.
//  3. Discards no-action results without invoking the handler.
//  4. If Instance() is nil, discards the action without error.
//  5. If ActionHandler() returns nil, discards the action without error.
//  6. Calls HandleAction on the handler and processes the result:
//     - Navigate: calls region.SetMode with the returned mode ID.
//     - CursorDelta: the instance updates its own internal scroll state.
//     - Dirty: the region is implicitly marked for redraw on next render cycle.
func (d *RegionInputDispatcher) Dispatch(r *region.Region, ev input.Event) {
	// Only press events trigger actions; release (and any future hold) are ignored.
	if ev.Type != input.Press {
		return
	}

	// Map raw input to a logical action.
	act := action.None
	if d.inputMapper != nil {
		act = d.inputMapper(ev)
	}

	// Discard no-action results.
	if act == action.None {
		return
	}

	// Get the active ModeInstance from the Region.
	inst := r.Instance()
	if inst == nil {
		return
	}

	// Get the action handler from the instance.
	handler := inst.ActionHandler()
	if handler == nil {
		return
	}

	// Invoke the handler. The instance owns cursor and item count internally.
	// During the legacy adapter transition, cursor is always 0.
	// Migrated instances will access their own internal scroll state through
	// their handler implementation, making these parameters informational only.
	res := handler.HandleAction(act, 0, 0)

	// Process the result.
	if res.Navigate != "" {
		// Navigate: switch the region to the target mode.
		if err := r.SetMode(res.Navigate); err != nil {
			log.Printf("dispatch: navigate to %q failed for region %q: %v", res.Navigate, r.Name(), err)
		}
		return
	}

	// CursorDelta: the instance updates its own internal scroll state.
	// No external action needed — the instance's handler has already adjusted
	// its internal cursor, and BuildView on the next render will reflect it.
	// (For legacy adapters that don't own cursor state, this is a no-op.)

	// Dirty: the region will be redrawn on the next render cycle.
	// The render loop renders all regions every tick, so no explicit marking
	// is needed — the next frame will pick up any state changes.
}

// Ensure RegionInputDispatcher satisfies region.InputDispatcher at compile time.
var _ region.InputDispatcher = (*RegionInputDispatcher)(nil)
