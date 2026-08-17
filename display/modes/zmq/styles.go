package zmq

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/modes/zmq/styles"
	"github.com/databeast/cyberhud/display/style"
)

// zmqRegistry is the per-mode StyleRegistry for the ZMQ display mode.
// The first style registered becomes the default.
var zmqRegistry = style.NewRegistry[content.ZMQData, content.Policy](styles.ColorMedium240x240Style, styles.MonoSlow800x480Style)

func registeredStyleNames() []string {
	styles := zmqRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}
