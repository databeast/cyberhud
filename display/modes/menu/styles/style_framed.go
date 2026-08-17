package styles

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/menu/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
)

var FramedStyle = def{
	name: "framed",
	reqs: style.SurfaceRequirements{
		MinWidth:        32,
		MinHeight:       32,
		PreferredWidth:  0,
		PreferredHeight: 0,
		Capability:      style.MonoSlow,
		MinRows:         2,
		MinCharsPerLine: 6,
	},
	p: Params{BuildFn: framedBuild},
}

func framedBuild(snapshot source.MenuSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()

	items := snapshot.Items
	if len(items) == 0 {
		items = []string{"(no menu items)"}
	}

	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})

	// Construct SuppressionContext from bridge dimensions.
	suppCtx := widgets.SuppressionContext{
		AvailableWidth:  bridge.AvailableContentWidth(),
		AvailableHeight: bridge.AvailableContentHeight(),
	}
	comp := widgets.NewCompositor(suppCtx)

	// Border covers the full panel — use hints for panel-covering bounds.
	bounds := image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight)
	comp.Add(borderframe.New(borderframe.Config{Bounds: bounds}))

	return style.ViewData{
		Items:   items,
		Sprites: comp.Sprites(),
		Static:  true,
	}
}
