package gpio_control

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	gcstyles "github.com/databeast/cyberhud/display/modes/gpio_control/styles"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

var (
	ColorOutput = gcstyles.ColorOutput
	ColorInput  = gcstyles.ColorInput
)

func BuildItems(pins []gpiomgr.PinState) []string { return source.BuildItems(pins) }

func CurrentSnapshot() []gpiomgr.PinState { return source.CurrentSnapshot() }

func SetSnapshot(pins []gpiomgr.PinState) { source.SetSnapshot(pins) }

func BuildColors(pins []gpiomgr.PinState) []color.Color { return gcstyles.BuildColors(pins) }

func BuildSprites(pins []gpiomgr.PinState, hints textlayout.TextHints) []widgets.Sprite {
	return gcstyles.BuildSprites(pins, hints)
}

func Registry() *style.StyleRegistry[source.Data, source.Policy] { return gpioControlRegistry }
