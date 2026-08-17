package dashboard

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

// instance implements displaymodes.ModeInstance for the dashboard mode.
type instance struct {
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string                    { return "dashboard" }
func (i *instance) Activate()                     {} // dashboard has no background work
func (i *instance) Deactivate()                   {} // dashboard has no background work
func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := getPanelHints(); ok {
		return BuildView(hints)
	}
	return style.ViewData{Items: []string{"error"}}
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}
