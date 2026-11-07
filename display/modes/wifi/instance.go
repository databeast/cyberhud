package wifi

import (
	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("wifi", newInstance)
}

// instance implements displaymodes.ModeInstance for the wifi mode.
type instance struct {
	gatherer source.Gatherer
}

func newInstance() displaymodes.ModeInstance {
	return &instance{}
}

func (i *instance) ID() string                    { return "wifi" }
func (i *instance) Activate()                     { i.gatherer.Reset() }
func (i *instance) Deactivate()                   {} // wifi has no background work
func (i *instance) ActionHandler() action.Handler { return Handler{} }

func (i *instance) BuildView() style.ViewData {
	vd := buildView(i.gatherer.GatherWifiState())

	return vd
}

func (i *instance) RenderCacheKey() uint32 {
	return renderCacheKey(i.gatherer.GatherWifiState())
}

func init() {
	catalog.Register(catalog.Definition{
		ID:       "wifi",
		Title:    "WiFi",
		Category: "System",
		Summary:  "Real-time wireless network status with signal bars, quality bar, and connection details.",
		Order:    30,
		Options:  append(source.Policy{}.Options(), catalog.OptionDefinition{Key: "style", Type: "string", Summary: "Visual layout variant for WiFi display.", Default: "", Allowed: wifiRegistry.RegisteredNames()}),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "wifi",
		Summary: "Query or set WiFi display options.",
		Usage:   "wifi [style=<name>] [show_border=<true|false>] [fgcolor=<cyan|green|amber|red|white|none>] [signal_display=<bars|percentage|dbm>] [show_frequency=<true|false>] [show_interface=<true|false>] [show_channel=<true|false>]",
		Handle:  HandleCommand,
	})
}
