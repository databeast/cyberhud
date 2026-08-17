package styles

import (
	"github.com/databeast/cyberhud/display/modes/menu/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
)

var PlainStyle = def{
	name: "plain",
	reqs: style.SurfaceRequirements{
		MinWidth:        0,
		MinHeight:       0,
		PreferredWidth:  0,
		PreferredHeight: 0,
		Capability:      style.MonoSlow,
		MinRows:         2,
		MinCharsPerLine: 6,
	},
	p: Params{BuildFn: plainBuild},
}

func plainBuild(snapshot source.MenuSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	items := snapshot.Items
	if len(items) == 0 {
		items = []string{"(no menu items)"}
	}

	return style.ViewData{
		Items:  items,
		Static: true,
	}
}
