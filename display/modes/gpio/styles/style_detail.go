package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var DetailStyle = def{
	name: "detail",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 32, PreferredWidth: 0, PreferredHeight: 0, Capability: style.MonoSlow, MinRows: 2, MinCharsPerLine: 12},
	p:    Params{Layout: layoutDetail},
}
