package styles

import (
	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
)

var CompactStyle = def{
	name: "compact",
	reqs: style.SurfaceRequirements{
		MinRows:         1,
		MinCharsPerLine: 10,
	},
	p: Params{BuildFn: compactBuild},
}

func compactBuild(snapshot source.SystemSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	_ = pol
	hints := ctx.Hints()
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	ip := "(none)"
	if len(snapshot.IPs) > 0 {
		ip = snapshot.IPs[0]
	}

	return style.ViewData{Items: []string{snapshot.Hostname + ": " + ip}}
}
