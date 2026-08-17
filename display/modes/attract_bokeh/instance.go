package attract_bokeh

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("attract_bokeh", newInstance)
}

// instance implements displaymodes.ModeInstance for the bokeh attract mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "attract_bokeh" }

func (i *instance) Activate() {} // animation is frame-driven, no background work

func (i *instance) Deactivate() {} // no background work to stop

func (i *instance) ActionHandler() action.Handler { return nil }

func (i *instance) BuildView() style.ViewData {
	return BuildView()
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}
