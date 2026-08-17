package system

import (
	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets/icons"
	"github.com/databeast/cyberhud/runtime/action"
)

// Font configuration for the system mode.
// Previously stored in modeFontConfig["system"] = {Family: "spleen", Tier: TierNormal}.
const systemFontFamily = "spleen"

var systemFontTier = tiercatalog.TierNormal

func init() {
	displaymodes.RegisterFactory("system", newInstance)
}

// instance implements displaymodes.ModeInstance for the system mode.
type instance struct{}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string                    { return "system" }
func (i *instance) Activate()                     {} // system has no background work
func (i *instance) Deactivate()                   {} // system has no background work
func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	hints, ok := getPanelHints()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}
	view := BuildView(hints, icons.Get)
	state := style.ViewData{
		Items:       view.Items,
		Static:      view.Static,
		Sprites:     view.Sprites,
		StyleReport: view.StyleReport,
	}

	return state
}

func (i *instance) RenderCacheKey() uint32 {
	return RenderCacheKey()
}

func init() {
	catalog.Register(catalog.Definition{
		ID:      "system",
		Title:   "System",
		Summary: "Host, uptime, and IP address information.",
		Order:   50,
		Options: []catalog.OptionDefinition{
			{Key: "style", Type: "string", Summary: "Visual presentation style.", Default: "", Allowed: registeredStyleNames()},
			{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
		},
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "system",
		Summary: "Query or set system display options.",
		Usage:   "system [style=<default|compact|cores|top>] [font=<auto|font-id>]",
		Handle:  HandleCommand,
	})
}
