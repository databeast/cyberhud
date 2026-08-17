package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 480x320 panels (landscape large).
// Landscape: width > height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal + timegraph for fast.

var MonoSlow480x320Style = def{
	name: "mono-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoSlow480x320OverviewStyle = def{
	name: "mono-slow-480x320-overview",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var MonoSlow480x320DetailStyle = def{
	name: "mono-slow-480x320-detail",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var MonoSlow480x320GraphStyle = def{
	name: "mono-slow-480x320-graph",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var MonoSlow480x320MinimalStyle = def{
	name: "mono-slow-480x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow480x320Style = def{
	name: "grayscale-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow480x320OverviewStyle = def{
	name: "grayscale-slow-480x320-overview",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow480x320DetailStyle = def{
	name: "grayscale-slow-480x320-detail",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow480x320GraphStyle = def{
	name: "grayscale-slow-480x320-graph",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow480x320MinimalStyle = def{
	name: "grayscale-slow-480x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast480x320Style = def{
	name: "grayscale-fast-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
}

var GrayscaleFast480x320OverviewStyle = def{
	name: "grayscale-fast-480x320-overview",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleFast480x320TimegraphStyle = def{
	name: "grayscale-fast-480x320-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var GrayscaleFast480x320DetailStyle = def{
	name: "grayscale-fast-480x320-detail",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleFast480x320GraphStyle = def{
	name: "grayscale-fast-480x320-graph",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleFast480x320MinimalStyle = def{
	name: "grayscale-fast-480x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow480x320Style = def{
	name: "color-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorSlow},
}

var ColorSlow480x320OverviewStyle = def{
	name: "color-slow-480x320-overview",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow480x320DetailStyle = def{
	name: "color-slow-480x320-detail",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow480x320GraphStyle = def{
	name: "color-slow-480x320-graph",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow480x320MinimalStyle = def{
	name: "color-slow-480x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast480x320Style = def{
	name: "color-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
}

var ColorFast480x320OverviewStyle = def{
	name: "color-480x320-overview",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorFast480x320TimegraphStyle = def{
	name: "color-480x320-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var ColorFast480x320DetailStyle = def{
	name: "color-480x320-detail",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorFast480x320GraphStyle = def{
	name: "color-480x320-graph",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorFast480x320MinimalStyle = def{
	name: "color-480x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
