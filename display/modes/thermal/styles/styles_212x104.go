package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 212x104 panels (e-ink landscape).
// Landscape: width > height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal (slow only, no timegraph).

var GrayscaleSlow212x104Style = def{
	name: "grayscale-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow212x104OverviewStyle = def{
	name: "grayscale-slow-212x104-overview",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow212x104DetailStyle = def{
	name: "grayscale-slow-212x104-detail",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow212x104GraphStyle = def{
	name: "grayscale-slow-212x104-graph",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow212x104MinimalStyle = def{
	name: "grayscale-slow-212x104-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow212x104Style = def{
	name: "color-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
}

var ColorSlow212x104OverviewStyle = def{
	name: "color-slow-212x104-overview",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow212x104DetailStyle = def{
	name: "color-slow-212x104-detail",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow212x104GraphStyle = def{
	name: "color-slow-212x104-graph",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow212x104MinimalStyle = def{
	name: "color-slow-212x104-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
