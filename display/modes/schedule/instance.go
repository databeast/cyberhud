package schedule

import (
	"time"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/schedule/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("schedule", newInstance)
}

// instance implements displaymodes.ModeInstance for the schedule mode. It owns
// the calendar-fetcher lifecycle via Activate/Deactivate and reads its Region's
// geometry from the embedded PanelHints.
type instance struct {
	displaymodes.PanelHints
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "schedule" }

// Activate starts the background calendar fetcher.
func (i *instance) Activate() {
	source.Activate()
}

// Deactivate stops the background calendar fetcher.
func (i *instance) Deactivate() {
	source.Deactivate()
}

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := i.Hints(); ok {
		return BuildView(time.Now(), hints)
	}
	return style.ViewData{Items: []string{"error"}, Cursor: -1}
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey(time.Now())
}
