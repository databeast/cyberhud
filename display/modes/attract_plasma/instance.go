package attract_plasma

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

// Compile-time interface compliance check.
var _ displaymodes.ModeInstance = (*instance)(nil)

func init() {
	displaymodes.RegisterFactory("attract_plasma", newInstance)
}

// instance implements displaymodes.ModeInstance for the plasma mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "attract_plasma" }

func (i *instance) Activate() {} // plasma has no background work (animation is frame-driven)

func (i *instance) Deactivate() {} // plasma has no background work

func (i *instance) ActionHandler() action.Handler { return nil }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := getPanelHints(); ok {
		return buildView(hints)
	}
	return style.ViewData{}
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}
