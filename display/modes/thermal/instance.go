package thermal

import (
	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/action"
)

// Font configuration for the thermal mode.
// Previously stored in modeFontConfig["thermal"] = {Family: "spleen", Tier: TierNormal}.
const thermalFontFamily = "spleen"

var thermalFontTier = tiercatalog.TierNormal

func init() {
	displaymodes.RegisterFactory("thermal", newInstance)
}

// instance implements displaymodes.ModeInstance for the thermal mode.
// It manages the thermal sampler lifecycle via Activate/Deactivate.
type instance struct {
	displaymodes.PanelHints
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string { return "thermal" }

// Activate starts the thermal sampling background goroutine.
func (i *instance) Activate() {
	source.Activate(func() int { return GetPolicy().RefreshMS })
}

// Deactivate stops the thermal sampling background goroutine.
func (i *instance) Deactivate() {
	source.Deactivate()
}

func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	// Hints come from this instance's embedded PanelHints, injected by the
	// hosting Region before Activate. Falls back to the process-wide store
	// (getPanelHints) for instances constructed outside a Region (tests, tooling),
	// then to zero-valued hints if neither source is available.
	hints, ok := i.Hints()
	if !ok {
		hints, ok = getPanelHints()
		if !ok {
			hints = textlayout.TextHints{}
		}
	}
	return buildView(hints)
}

func (i *instance) RenderCacheKey() uint32 {
	snap := source.CurrentSnapshot()
	hints, _ := i.Hints()
	return RenderCacheKey(snap, 0, 0, 0, hints.PixelWidth, hints.PixelHeight)
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "thermal",
		Title:   "Thermal",
		Summary: "Dedicated thermal monitoring with per-zone temperatures, history graphs, and threshold alerts.",
		Order:   65,
		Options: append(source.Policy{}.Options(), catalog.OptionDefinition{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: allowedStyleNames()}),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "thermal",
		Summary: "Query or set thermal display options.",
		Usage:   "thermal [style=<color-320x240-overview|mono-slow-128x64|...>] [font=<auto|font-id>] [refresh_ms=<ms>] [warn_threshold=<°C>] [crit_threshold=<°C>] [show_border=<true|false>] [unit=<C|F>] [fgcolor=<thermal|cyan|green|amber|red|white|none>] [show_led=<true|false>] [show_refresh_bar=<true|false>]",
		Handle:  HandleCommand,
	})
}
