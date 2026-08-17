package styles

import "github.com/databeast/cyberhud/display/style"

var CompactStyle = def{
	name: "compact",
	reqs: style.SurfaceRequirements{
		MinWidth:        0,
		MinHeight:       0,
		PreferredWidth:  0,
		PreferredHeight: 0,
		Capability:      style.MonoSlow,
		MinRows:         1,
		MinCharsPerLine: 10,
	},
	p: Params{Layout: layoutSummary},
}
