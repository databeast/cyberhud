package styles

import (
	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	"github.com/databeast/cyberhud/display/style"
)

var CompactStyle = def{
	name: "compact",
	reqs: style.SurfaceRequirements{
		MinWidth:        0,
		MinHeight:       0,
		PreferredWidth:  0,
		PreferredHeight: 0,
		Capability:      style.MonoSlow,
	},
	p: Params{BuildFn: compactBuild},
}

func compactBuild(snapshot source.Data, _ source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()
	bridge := ctx.Layout(0)

	pins := snapshot.Pins

	items := source.BuildItems(pins)
	colors := BuildColors(pins)

	if len(pins) == 0 {
		return style.ViewData{
			Items:  items,
			Colors: colors,
		}
	}

	// Build LED sprites using bridge metrics.
	sprites := BuildSprites(pins, hints)

	// Truncate text to MaxCharsPerRow.
	maxCols := 0
	if bridge.GlyphAdvance() > 0 {
		maxCols = bridge.AvailableContentWidth() / bridge.GlyphAdvance()
	}
	if maxCols > 0 {
		for i, item := range items {
			if len(item) > maxCols {
				items[i] = item[:maxCols]
			}
		}
	}

	// Limit visible rows to MaxVisibleRows.
	maxRows := bridge.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
		if len(colors) > maxRows {
			colors = colors[:maxRows]
		}
	}

	return style.ViewData{
		Items:   items,
		Colors:  colors,
		Sprites: sprites,
		Static:  true,
	}
}
