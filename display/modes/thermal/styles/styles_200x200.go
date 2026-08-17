package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 200x200 panels (square e-ink).
// Square: width == height, width ≥ 128 → overview, detail, graph, minimal (slow only, no timegraph).

var GrayscaleSlow200x200Style = def{
	name: "grayscale-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow200x200OverviewStyle = def{
	name: "grayscale-slow-200x200-overview",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow200x200DetailStyle = def{
	name: "grayscale-slow-200x200-detail",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow200x200GraphStyle = def{
	name: "grayscale-slow-200x200-graph",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow200x200MinimalStyle = def{
	name: "grayscale-slow-200x200-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow200x200Style = def{
	name: "color-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
}

var ColorSlow200x200OverviewStyle = def{
	name: "color-slow-200x200-overview",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow200x200DetailStyle = def{
	name: "color-slow-200x200-detail",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow200x200GraphStyle = def{
	name: "color-slow-200x200-graph",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow200x200MinimalStyle = def{
	name: "color-slow-200x200-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
