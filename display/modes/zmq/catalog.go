package zmq

import (
	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/zmq/content"
)

func init() {
	catalog.Register(catalog.Definition{
		ID:      "zmq",
		Title:   "ZMQ Monitor",
		Summary: "Live ZeroMQ message stream with JSON field filtering.",
		Order:   37,
		Options: append(content.Policy{}.Options(), catalog.OptionDefinition{
			Key: "style", Type: "string", Summary: "Visual presentation style.",
			Default: "", Allowed: registeredStyleNames(),
		}),
	})
	catalog.RegisterCommand(catalog.CommandDefinition{
		Verb:    "zmq",
		Summary: "Query or set ZMQ display mode options.",
		Usage:   "zmq [key=value ...] | zmq clear",
		Handle:  HandleCommand,
	})
}
