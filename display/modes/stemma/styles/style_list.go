package styles

import "github.com/databeast/cyberhud/display/style"

var ListStyle = def{
	name: "list",
	reqs: style.SurfaceRequirements{
		MinWidth:        0,
		MinHeight:       0,
		PreferredWidth:  0,
		PreferredHeight: 0,
		Capability:      style.MonoSlow,
		MinRows:         2,
		MinCharsPerLine: 10,
	},
	p: Params{Layout: layoutList, RowFormatter: defaultDeviceRow},
}
