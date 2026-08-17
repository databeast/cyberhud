package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var DashboardStyle = def{
	name: "dashboard",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 64, PreferredWidth: 0, PreferredHeight: 0, Capability: style.ColorSlow},
	p:    Params{Layout: layoutDashboard},
}
