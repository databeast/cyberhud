package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 250x122 panels (e-ink landscape).
// Landscape: width > height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal (slow only, no timegraph).

var GrayscaleSlow250x122Style = def{
	name: "grayscale-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow250x122OverviewStyle = def{
	name: "grayscale-slow-250x122-overview",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow250x122DetailStyle = def{
	name: "grayscale-slow-250x122-detail",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow250x122GraphStyle = def{
	name: "grayscale-slow-250x122-graph",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow250x122MinimalStyle = def{
	name: "grayscale-slow-250x122-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow250x122Style = def{
	name: "color-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
}

var ColorSlow250x122OverviewStyle = def{
	name: "color-slow-250x122-overview",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow250x122DetailStyle = def{
	name: "color-slow-250x122-detail",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow250x122GraphStyle = def{
	name: "color-slow-250x122-graph",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow250x122MinimalStyle = def{
	name: "color-slow-250x122-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
