package attract_starfield

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("attract_starfield", newInstance)
}

// instance implements displaymodes.ModeInstance for the starfield mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "attract_starfield" }

func (i *instance) Activate() {} // starfield has no background work (animation is frame-driven)

func (i *instance) Deactivate() {} // starfield has no background work

func (i *instance) ActionHandler() action.Handler { return nil }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := getPanelHints(); ok {
		return BuildView(hints)
	}
	return style.ViewData{Items: []string{"error"}}
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}
