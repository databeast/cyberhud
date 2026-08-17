package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 160x80 panels (landscape compact).
// Landscape: width ≥ height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal + timegraph for fast.

var MonoSlow160x80Style = def{
	name: "mono-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
}

var MonoSlow160x80OverviewStyle = def{
	name: "mono-slow-160x80-overview",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var MonoSlow160x80DetailStyle = def{
	name: "mono-slow-160x80-detail",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var MonoSlow160x80GraphStyle = def{
	name: "mono-slow-160x80-graph",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var MonoSlow160x80MinimalStyle = def{
	name: "mono-slow-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow160x80Style = def{
	name: "grayscale-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow160x80OverviewStyle = def{
	name: "grayscale-slow-160x80-overview",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow160x80DetailStyle = def{
	name: "grayscale-slow-160x80-detail",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow160x80GraphStyle = def{
	name: "grayscale-slow-160x80-graph",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow160x80MinimalStyle = def{
	name: "grayscale-slow-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast160x80Style = def{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
}

var GrayscaleFast160x80OverviewStyle = def{
	name: "grayscale-fast-160x80-overview",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleFast160x80TimegraphStyle = def{
	name: "grayscale-fast-160x80-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var GrayscaleFast160x80DetailStyle = def{
	name: "grayscale-fast-160x80-detail",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleFast160x80GraphStyle = def{
	name: "grayscale-fast-160x80-graph",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleFast160x80MinimalStyle = def{
	name: "grayscale-fast-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow160x80Style = def{
	name: "color-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
}

var ColorSlow160x80OverviewStyle = def{
	name: "color-slow-160x80-overview",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow160x80DetailStyle = def{
	name: "color-slow-160x80-detail",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow160x80GraphStyle = def{
	name: "color-slow-160x80-graph",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow160x80MinimalStyle = def{
	name: "color-slow-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast160x80Style = def{
	name: "color-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
}

var ColorFast160x80OverviewStyle = def{
	name: "color-160x80-overview",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorFast160x80TimegraphStyle = def{
	name: "color-160x80-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var ColorFast160x80DetailStyle = def{
	name: "color-160x80-detail",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorFast160x80GraphStyle = def{
	name: "color-160x80-graph",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorFast160x80MinimalStyle = def{
	name: "color-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
