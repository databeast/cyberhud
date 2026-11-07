package styles

import (
	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

var ListStyle = def{
	name: "list",
	reqs: style.SurfaceRequirements{
		MinWidth:        0,
		MinHeight:       0,
		PreferredWidth:  0,
		PreferredHeight: 0,
		Capability:      style.MonoSlow,
		MinRows:         2,
		MinCharsPerLine: 8,
	},
	p: Params{BuildFn: listBuild},
}

func listBuild(snapshot source.Data, _ source.Policy, ctx style.StyleContext, d def) style.ViewData {
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
	m := ctx.DefaultMetrics()
	sprites := BuildSprites(pins, bridge, m)

	// Truncate text to MaxCharsPerRow.
	maxCols := bridge.MaxChars(m)
	if maxCols > 0 {
		for i, item := range items {
			if len(item) > maxCols {
				items[i] = item[:maxCols]
			}
		}
	}

	// Limit visible rows to MaxVisibleRows.
	maxRows := bridge.MaxVisibleRows(m)
	if maxRows > 0 && len(items) > maxRows {
		items = items[:maxRows]
		if len(colors) > maxRows {
			colors = colors[:maxRows]
		}
	}

	// TextLabel rendering: produce text label sprites if a tier-resolved face is available.
	if hints.Catalog.PixelWidth() > 0 {
		spriteFace := hints.Face("terminus", tiercatalog.TierNormal)
		if spriteFace != nil {
			fm := spriteFace.Metrics()
			faceMetrics := layout.NewMetrics(fm.GlyphAdvance, fm.RowHeight, "")
			textLabelSprites := buildControlTextLabelSprites(pins, items, bridge, faceMetrics, spriteFace, colors, textlayout.MaxVisibleRows(textlayout.WithFont(hints, spriteFace), 0))
			if textLabelSprites != nil {
				sprites = append(sprites, textLabelSprites...)
			}
		}
	}

	return style.ViewData{
		Items:   items,
		Colors:  colors,
		Sprites: sprites,
		
	}
}
