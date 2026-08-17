package systemd

import (
	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/systemd/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/runtime/action"
)

// Font configuration for the systemd mode.
// Previously stored in modeFontConfig["systemd"] = {Family: "spleen", Tier: TierNormal}.
const systemdFontFamily = "spleen"

var systemdFontTier = tiercatalog.TierNormal

func init() {
	displaymodes.RegisterFactory("systemd", newInstance)
}

// instance implements displaymodes.ModeInstance for the systemd mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string                    { return "systemd" }
func (i *instance) Activate()                     {} // systemd has no background work
func (i *instance) Deactivate()                   {} // systemd has no background work
func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	if hints, ok := getPanelHints(); ok {
		return BuildView(hints)
	}
	return style.ViewData{Items: []string{"error"}}
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "systemd",
		Title:   "Systemd Boot",
		Summary: "Boot progress and target transition diagnostics.",
		Order:   5,
		Options: []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: registeredStyleNames()},
			{Key: "color_accent", Type: "string", Summary: "Boot progress accent color.", Default: "amber", Allowed: source.AllowedAccents},
		},
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "systemd",
		Summary: "Query or set systemd display options.",
		Usage:   "systemd [style=<name>] [show_border=<true|false>] [color_accent=<cyan|green|amber|red|white|none>]",
		Handle:  HandleCommand,
	})
}
