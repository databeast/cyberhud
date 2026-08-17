package gpio

import (
	"image/color"

	"github.com/databeast/cyberhud/display/modes/gpio/source"
	"github.com/databeast/cyberhud/display/style"
)

var AllowedFGColors = source.AllowedFGColors

func ResolveFGColor(accent string) color.RGBA { return resolveFGColor(accent) }

func DimFGColor(accent string) color.RGBA { return dimFGColor(accent) }

func CenterOffsets(items []string, glyphAdvance, width int) []int {
	offsets := make([]int, len(items))
	for i, item := range items {
		textWidth := len([]rune(item)) * glyphAdvance
		if textWidth < width {
			offsets[i] = (width - textWidth) / 2
		}
	}
	return offsets
}

func Registry() *style.StyleRegistry[source.GpioSnapshot, source.Policy] { return gpioRegistry }
