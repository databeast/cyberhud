package styles

import "github.com/databeast/cyberhud/display/style"

// GPIO style declarations for this panel group.
//
// These are hand-tweakable declarations over the shared layouts in core.go.

var ActivityStyle = def{
	name: "activity",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 16, PreferredWidth: 0, PreferredHeight: 0, Capability: style.MonoFast},
	p:    Params{Layout: layoutActivity},
}
