package testpattern

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("testpattern", newInstance)
}

// instance implements displaymodes.ModeInstance for the testpattern mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "testpattern" }

func (i *instance) Activate() {} // testpattern has no background work

func (i *instance) Deactivate() {} // testpattern has no background work

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	hints, ok := modehints.Current()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}
	view := BuildView(hints, false)
	state := style.ViewData{
		Sprites: view.Sprites,
	}

	return state
}

func (i *instance) RenderCacheKey() uint32 {
	hints, ok := modehints.Current()
	if !ok {
		return 0 // no hints yet; key changes once hints propagate
	}
	return region.CalcRegionCacheKey(RenderCacheKey(hints))
}
