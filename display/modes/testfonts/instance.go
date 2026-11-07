package testfonts

import (
	"time"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

var nowFunc = time.Now

func init() {
	displaymodes.RegisterFactory("testfonts", newInstance)
}

// instance implements displaymodes.ModeInstance for the testfonts mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "testfonts" }

func (i *instance) Activate() {} // testfonts has no background work

func (i *instance) Deactivate() {} // testfonts has no background work

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	hints, ok := modehints.Current()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}
	view := BuildView(hints, nowFunc())
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
	return region.CalcRegionCacheKey(RenderCacheKey(hints, nowFunc()))
}
