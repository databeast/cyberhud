package image

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("image", newInstance)
}

// instance implements displaymodes.ModeInstance for the image mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string                    { return "image" }
func (i *instance) Activate()                     {} // image has no background work
func (i *instance) Deactivate()                   {} // image has no background work
func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	hints, ok := modehints.Current()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}
	view := BuildView(hints)
	state := style.ViewData{
		Items:   view.Items,
		Sprites: view.Sprites,
	}

	return state
}

func (i *instance) RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(Signature())
}
