package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 264x176 panels (e-ink landscape).
// Landscape: width > height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal (slow only, no timegraph).

var GrayscaleSlow264x176Style = def{
	name: "grayscale-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow264x176OverviewStyle = def{
	name: "grayscale-slow-264x176-overview",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow264x176DetailStyle = def{
	name: "grayscale-slow-264x176-detail",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow264x176GraphStyle = def{
	name: "grayscale-slow-264x176-graph",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow264x176MinimalStyle = def{
	name: "grayscale-slow-264x176-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow264x176Style = def{
	name: "color-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
}

var ColorSlow264x176OverviewStyle = def{
	name: "color-slow-264x176-overview",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow264x176DetailStyle = def{
	name: "color-slow-264x176-detail",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow264x176GraphStyle = def{
	name: "color-slow-264x176-graph",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow264x176MinimalStyle = def{
	name: "color-slow-264x176-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
