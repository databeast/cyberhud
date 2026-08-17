package attract_particles

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("attract_particles", newInstance)
}

// instance implements displaymodes.ModeInstance for the particles mode.
type instance struct{}

// Compile-time interface compliance check.
var _ displaymodes.ModeInstance = (*instance)(nil)

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "attract_particles" }

func (i *instance) Activate() {} // particles mode is frame-driven, no background work

func (i *instance) Deactivate() {} // particles mode is frame-driven, no background work

func (i *instance) ActionHandler() action.Handler { return nil }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := getPanelHints(); ok {
		return buildView(hints)
	}
	return style.ViewData{Items: []string{"error"}}
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}
