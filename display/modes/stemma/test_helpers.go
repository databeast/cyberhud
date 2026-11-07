package stemma

import (
	"image"
	"image/color"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/modes/stemma/styles"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

type Policy = source.Policy

var ColorPresent = styles.ColorPresent
var ColorAbsent = styles.ColorAbsent

func BuildItems(devs []*source.Device) []string { return styles.BuildItems(devs) }
func BuildColors(devs []*source.Device, present, absent color.Color) []color.Color {
	return styles.BuildColors(devs, present, absent)
}
func BuildSprites(devs []*source.Device, rowHeight int, getIcon func(string) (image.Image, bool)) []widgets.Sprite {
	bridge := layout.NewLayoutBridge(textlayout.TextHints{}, layout.BridgeConfig{})
	return styles.BuildSprites(devs, bridge, layout.NewMetrics(1, rowHeight, ""), getIcon)
}

func StemmaRegistryEnumerate() []style.Style[source.StemmaSnapshot, source.Policy] {
	return stemmaRegistry.Enumerate()
}

func StemmaRegistryDefaultName() string {
	s, _ := style.ResolveStyle(stemmaRegistry, textlayout.TextHints{}, "stemma", "")
	return s.Name()
}

func NewInstanceForTest() displaymodes.ModeInstance { return newInstance() }
