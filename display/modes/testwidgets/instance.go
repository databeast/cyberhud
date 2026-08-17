package testwidgets

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

// Font configuration for the testwidgets mode.
const (
	testwidgetsFontFamily = "spleen"
)

var testwidgetsFontTier = tiercatalog.TierSmall
var nowFunc = time.Now

func init() {
	displaymodes.RegisterFactory("testwidgets", newInstance)
}

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }

// instance implements displaymodes.ModeInstance for the testwidgets mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "testwidgets" }

func (i *instance) Activate() {} // testwidgets has no background work

func (i *instance) Deactivate() {} // testwidgets has no background work

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	hints, ok := getPanelHints()
	if !ok {
		return style.ViewData{
			Items:        []string{"error"},
			VisibleCount: 1,
			StyleReport:  style.StyleReport{Name: "testwidgets", Reason: "builtin"},
		}
	}
	view := BuildView(hints, nowFunc())
	state := style.ViewData{
		Static:       view.Static,
		Sprites:      view.Sprites,
		VisibleCount: view.VisibleCount,
		StyleReport: style.StyleReport{
			Name:   "testwidgets",
			Reason: "builtin",
		},
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
