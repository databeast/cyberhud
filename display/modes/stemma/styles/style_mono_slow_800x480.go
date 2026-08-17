package styles

import "github.com/databeast/cyberhud/display/style"

var MonoSlow800x480Style = def{
	name: "mono-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p: Params{
		Layout:       layoutPoster,
		HeaderText:   "STEMMA QT",
		RowFormatter: defaultDeviceRow,
		UseBorder:    true,
		BorderTheme:  "sharp",
	},
}
