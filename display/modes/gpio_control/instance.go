package gpio_control

import (
	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	displaymodes.RegisterFactory("gpio-control", newInstance)
}

// instance implements displaymodes.ModeInstance for the GPIO Control mode.
type instance struct {
	hints  textlayout.TextHints
	mgr    interface{ Snapshot() []gpiomgr.PinState }
	cursor int
	topRow int
}

func newInstance() displaymodes.ModeInstance {
	return &instance{
		hints: initialHints(),
		mgr:   source.GPIOControlManager(), // self-sourced from package singleton
	}
}

func (i *instance) ID() string { return "gpio-control" }

func (i *instance) Activate() {
	if i.mgr != nil {
		source.SetGPIOControlManager(i.mgr)
	}
}

func (i *instance) Deactivate() {
	source.SetGPIOControlManager(nil)
}

func (i *instance) ActionHandler() action.Handler { return instanceHandler{i} }

func (i *instance) BuildView() style.ViewData {
	pins := i.pins()
	// Update the shared snapshot so the action handler can read it for Toggle.
	source.SetSnapshot(pins)
	vd := BuildView(pins, i.hints, i.cursor)
	// Apply scroll state clamping.
	maxVisible := textlayout.MaxVisibleRows(i.hints, 0)
	i.cursor, i.topRow = style.ClampList(i.cursor, i.topRow, maxVisible, len(vd.Items))
	vd.Cursor = i.cursor
	vd.TopRow = i.topRow

	return vd
}

// pins returns the current GPIO pin snapshot, or nil if the manager is unavailable.
func (i *instance) pins() []gpiomgr.PinState {
	if i.mgr == nil {
		return nil
	}
	return i.mgr.Snapshot()
}

func initialHints() textlayout.TextHints {
	if h, ok := modehints.Current(); ok {
		return h
	}
	return textlayout.TextHints{}
}
