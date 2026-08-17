package clock

import (
	"time"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("clock", newInstance)
}

// instance implements displaymodes.ModeInstance for the clock mode.
//
// The embedded PanelHints holds this instance's own Region hints, injected by
// Region.SetMode before activation. Clock is the reference implementation of that
// pattern: it reads geometry from instance state rather than from the process-wide
// modehints store, so two Regions running clock on different panels lay out
// independently. See displaymodes.PanelHints for the migration guide.
type instance struct {
	displaymodes.PanelHints
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "clock" }

func (i *instance) Activate() {} // clock has no background work

func (i *instance) Deactivate() {} // clock has no background work

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	// Hints come from this instance, not from a global. Falls back to the legacy
	// store only for instances built outside a Region (tooling and tests).
	if hints, ok := i.Hints(); ok {
		return BuildView(time.Now(), hints)
	}
	return style.ViewData{Items: []string{"error"}}
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey(time.Now())
}
