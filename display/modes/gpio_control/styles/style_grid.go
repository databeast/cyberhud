package styles

import (
	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	"github.com/databeast/cyberhud/display/style"
)

var GridStyle = def{
	name: "grid",
	reqs: style.SurfaceRequirements{
		MinWidth:        64,
		MinHeight:       64,
		PreferredWidth:  0,
		PreferredHeight: 0,
		Capability:      style.ColorSlow,
	},
	p: Params{BuildFn: gridBuild},
}

func gridBuild(snapshot source.Data, _ source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()

	pins := snapshot.Pins

	items := source.BuildItems(pins)
	colors := BuildColors(pins)

	if len(pins) == 0 {
		return style.ViewData{
			Items:  items,
			Colors: colors,
		}
	}

	gridSprites := BuildGridView(pins, hints, snapshot.Cursor)

	return style.ViewData{
		Items:   items,
		Colors:  colors,
		Sprites: gridSprites,
		Static:  true,
	}
}
