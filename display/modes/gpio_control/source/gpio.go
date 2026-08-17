package source

import (
	"fmt"
	"sync"

	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
	"periph.io/x/conn/v3/gpio"
)

// gpioControlMgr is the package-level GPIO manager singleton used by the
// instance to obtain pin state snapshots. Set via SetGPIOControlManager
// before any BuildView calls. If nil, BuildView returns valid ViewData with
// zero data items.
type Snapshotter interface {
	Snapshot() []gpiomgr.PinState
}

var (
	state = struct {
		sync.RWMutex
		pins []gpiomgr.PinState
		mgr  *gpiomgr.Manager
	}{
		pins: []gpiomgr.PinState{},
	}
	gpioControlMgr Snapshotter
)

// SetGPIOControlManager registers the GPIO manager for the gpio-control mode
// instance to self-source pin data. Call from the daemon wiring (cmd/cyberhudd)
// before any mode activation.
func SetGPIOControlManager(mgr Snapshotter) {
	gpioControlMgr = mgr
}

// GPIOControlManager returns the configured package-level GPIO snapshotter.
func GPIOControlManager() Snapshotter {
	return gpioControlMgr
}

// CurrentSnapshot returns the most recent cached pin snapshot.
func CurrentSnapshot() []gpiomgr.PinState {
	state.RLock()
	defer state.RUnlock()
	out := make([]gpiomgr.PinState, len(state.pins))
	copy(out, state.pins)
	return out
}

// SetSnapshot updates the internal pin cache.
func SetSnapshot(pins []gpiomgr.PinState) {
	state.Lock()
	state.pins = append([]gpiomgr.PinState(nil), pins...)
	state.Unlock()
}

// SetManager registers the GPIO manager for control operations.
func SetManager(mgr *gpiomgr.Manager) {
	state.Lock()
	state.mgr = mgr
	state.Unlock()
}

// Toggle attempts to flip the state of a GPIO pin.
func Toggle(pin gpiomgr.PinState) error {
	state.RLock()
	mgr := state.mgr
	state.RUnlock()
	if mgr == nil {
		return fmt.Errorf("GPIO manager unavailable")
	}
	nextLevel := gpio.Low
	if !pin.Level {
		nextLevel = gpio.High
	}
	return mgr.SetOutput(pin.Number, nextLevel)
}
