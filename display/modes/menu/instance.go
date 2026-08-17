package menu

import (
	"image"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/action"
)

// Font configuration for the menu mode.
const menuFontFamily = "spleen"

var menuFontTier = tiercatalog.TierNormal

func init() {
	displaymodes.RegisterFactory("menu", newInstance)
}

// instance implements displaymodes.ModeInstance for the menu mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string                    { return "menu" }
func (i *instance) Activate()                     {}
func (i *instance) Deactivate()                   {}
func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := getPanelHints(); ok {
		return BuildView(hints, nil)
	}
	return BuildView(textlayout.DefaultTextHints(image.Rect(0, 0, 240, 240)), nil)
}

func (i *instance) RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(GetPolicy().Fingerprint())
}
