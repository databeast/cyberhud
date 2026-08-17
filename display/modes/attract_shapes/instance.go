package attract_shapes

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("attract_shapes", newInstance)
}

// instance implements displaymodes.ModeInstance for the attract_shapes mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "attract_shapes" }

func (i *instance) Activate() {} // shapes is frame-driven, no background work

func (i *instance) Deactivate() {} // shapes is frame-driven, no background work

func (i *instance) ActionHandler() action.Handler { return nil }

func (i *instance) BuildView() style.ViewData {
	return BuildView()
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}
