package testicons

import (
	"time"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/action"
)

// Font configuration for the testicons mode.
const (
	testiconsFontFamily = "spleen"
)

var testiconsFontTier = tiercatalog.TierSmall
var nowFunc = time.Now

func init() {
	displaymodes.RegisterFactory("testicons", newInstance)
}

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }

// instance implements displaymodes.ModeInstance for the testicons mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "testicons" }

func (i *instance) Activate() {} // testicons has no background work

func (i *instance) Deactivate() {} // testicons has no background work

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	hints, ok := getPanelHints()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}
	view := BuildView(hints, nowFunc())
	state := style.ViewData{
		Static:  view.Static,
		Sprites: view.Sprites,
	}

	return state
}

func (i *instance) RenderCacheKey() uint32 {
	hints, ok := getPanelHints()
	if !ok {
		return 0 // no hints yet; key changes once hints propagate
	}
	return region.CalcRegionCacheKey(RenderCacheKey(hints, nowFunc()))
}
