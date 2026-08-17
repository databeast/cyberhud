package testfonts

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

// Font configuration for the testfonts mode.
// Previously stored in modeFontConfig["testfonts"] = {Family: "spleen", Tier: TierSmall}.
const (
	testfontsFontFamily = "spleen"
)

var testfontsFontTier = tiercatalog.TierSmall
var nowFunc = time.Now

func init() {
	displaymodes.RegisterFactory("testfonts", newInstance)
}

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }

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
