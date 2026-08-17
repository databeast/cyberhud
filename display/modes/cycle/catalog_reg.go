package cycle

import "github.com/databeast/cyberhud/display/catalog"

func init() {
	catalog.Register(catalog.Definition{
		ID:      "cycle",
		Title:   "Cycle",
		Summary: "Auto-cycles through configured display modes on one or more regions at a configurable interval.",
		Order:   90,
		Options: []catalog.OptionDefinition{
			{Key: "interval", Type: "duration", Summary: "Time between automatic mode switches.", Default: "30s"},
			{Key: "modes", Type: "string-list", Summary: "Ordered list of mode IDs to cycle through.", Default: ""},
			{Key: "regions", Type: "int-list", Summary: "Region indices to cycle on.", Default: ""},
		},
	})

	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "cycle",
		Summary: "Query or set cycle mode configuration.",
		Usage:   "cycle [interval=<duration>] [modes=<id,id,...>] [regions=<0,1,...>]",
		Handle:  HandleCommand,
	})

	catalog.RegisterSnapshotter("cycle", cycleSnapshotter{})
}
