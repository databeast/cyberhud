package source

import gpiomgr "github.com/databeast/cyberhud/hardware/gpio"

// GpioSnapshot captures the current GPIO pin state for style Build methods.
type GpioSnapshot struct {
	Pins   []gpiomgr.PinState
	Policy Policy
}
