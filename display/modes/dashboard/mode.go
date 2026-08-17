package dashboard

import (
	"github.com/databeast/cyberhud/display/catalog"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/dashboard/source"
)

func init() {
	catalog.Register(catalog.Definition{
		ID:      "dashboard",
		Title:   "Dashboard",
		Summary: "Glanceable liveness panel showing hostname, uptime heartbeat, IP, WiFi, version, and panel name.",
		Order:   20,
		Options: append(source.Policy{}.Options(), styleOptions()...),
	})

	displaymodes.RegisterFactory("dashboard", newInstance)

}

func styleOptions() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{
			Key:     "style",
			Type:    "string",
			Summary: "Visual presentation style.",
			Default: "",
			Allowed: registeredStyleNames(),
		},
	}
}
