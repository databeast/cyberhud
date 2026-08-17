package image

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/action"
)

// Font configuration for the image mode.
// Previously stored in modeFontConfig["image"] = {Family: "spleen", Tier: TierNormal}.
const imageFontFamily = "spleen"

var imageFontTier = tiercatalog.TierNormal

func init() {
	displaymodes.RegisterFactory("image", newInstance)
}

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }

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
	hints, ok := getPanelHints()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}
	view := BuildView(hints)
	state := style.ViewData{
		Items:   view.Items,
		Static:  view.Static,
		Sprites: view.Sprites,
	}

	return state
}

func (i *instance) RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(Signature())
}
