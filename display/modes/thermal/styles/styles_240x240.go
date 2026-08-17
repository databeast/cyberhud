package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 240x240 panels (square).
// Square: width == height, width ≥ 128 → overview, detail, graph, minimal + timegraph for fast.

var MonoSlow240x240Style = def{
	name: "mono-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoSlow240x240OverviewStyle = def{
	name: "mono-slow-240x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var MonoSlow240x240DetailStyle = def{
	name: "mono-slow-240x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var MonoSlow240x240GraphStyle = def{
	name: "mono-slow-240x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var MonoSlow240x240MinimalStyle = def{
	name: "mono-slow-240x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow240x240Style = def{
	name: "grayscale-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x240OverviewStyle = def{
	name: "grayscale-slow-240x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow240x240DetailStyle = def{
	name: "grayscale-slow-240x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow240x240GraphStyle = def{
	name: "grayscale-slow-240x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow240x240MinimalStyle = def{
	name: "grayscale-slow-240x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast240x240Style = def{
	name: "grayscale-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
}

var GrayscaleFast240x240OverviewStyle = def{
	name: "grayscale-fast-240x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleFast240x240TimegraphStyle = def{
	name: "grayscale-fast-240x240-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var GrayscaleFast240x240DetailStyle = def{
	name: "grayscale-fast-240x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleFast240x240GraphStyle = def{
	name: "grayscale-fast-240x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleFast240x240MinimalStyle = def{
	name: "grayscale-fast-240x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow240x240Style = def{
	name: "color-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorSlow},
}

var ColorSlow240x240OverviewStyle = def{
	name: "color-slow-240x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow240x240DetailStyle = def{
	name: "color-slow-240x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow240x240GraphStyle = def{
	name: "color-slow-240x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow240x240MinimalStyle = def{
	name: "color-slow-240x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast240x240Style = def{
	name: "color-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
}

var ColorFast240x240OverviewStyle = def{
	name: "color-240x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorFast240x240TimegraphStyle = def{
	name: "color-240x240-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var ColorFast240x240DetailStyle = def{
	name: "color-240x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorFast240x240GraphStyle = def{
	name: "color-240x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorFast240x240MinimalStyle = def{
	name: "color-240x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
