package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var ListStyle = def{
	name: "list",
	reqs: style.SurfaceRequirements{MinWidth: 0, MinHeight: 0, PreferredWidth: 0, PreferredHeight: 0, Capability: style.MonoSlow, MinRows: 2, MinCharsPerLine: 8},
	p:    Params{Layout: layoutList},
}
