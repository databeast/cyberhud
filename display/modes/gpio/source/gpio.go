package source

import gpiomgr "github.com/databeast/cyberhud/hardware/gpio"

var gpioMgr interface {
	Snapshot() []gpiomgr.PinState
}

// SetGPIOManager registers the GPIO manager for the gpio mode instance to self-source pin data.
func SetGPIOManager(mgr interface{ Snapshot() []gpiomgr.PinState }) {
	gpioMgr = mgr
}

// GPIOManager returns the currently registered GPIO manager.
func GPIOManager() interface{ Snapshot() []gpiomgr.PinState } {
	return gpioMgr
}

// BuildItems returns GPIO mode row text for the primary list screen.
func BuildItems(pins []gpiomgr.PinState) []string {
	items := make([]string, len(pins))
	for i, p := range pins {
		items[i] = p.String()
	}
	return items
}

// BuildItemsTruncated returns GPIO mode rows truncated to maxChars.
func BuildItemsTruncated(pins []gpiomgr.PinState, maxChars int) []string {
	items := make([]string, len(pins))
	for i, p := range pins {
		s := p.String()
		if maxChars > 0 && len(s) > maxChars {
			s = s[:maxChars]
		}
		items[i] = s
	}
	return items
}
