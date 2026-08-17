package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 240x135 panels (landscape wide).
// Landscape: width > height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal + timegraph for fast.

var MonoSlow240x135Style = def{
	name: "mono-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
}

var MonoSlow240x135OverviewStyle = def{
	name: "mono-slow-240x135-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var MonoSlow240x135DetailStyle = def{
	name: "mono-slow-240x135-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var MonoSlow240x135GraphStyle = def{
	name: "mono-slow-240x135-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var MonoSlow240x135MinimalStyle = def{
	name: "mono-slow-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow240x135Style = def{
	name: "grayscale-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x135OverviewStyle = def{
	name: "grayscale-slow-240x135-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow240x135DetailStyle = def{
	name: "grayscale-slow-240x135-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow240x135GraphStyle = def{
	name: "grayscale-slow-240x135-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow240x135MinimalStyle = def{
	name: "grayscale-slow-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast240x135Style = def{
	name: "grayscale-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
}

var GrayscaleFast240x135OverviewStyle = def{
	name: "grayscale-fast-240x135-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleFast240x135TimegraphStyle = def{
	name: "grayscale-fast-240x135-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var GrayscaleFast240x135DetailStyle = def{
	name: "grayscale-fast-240x135-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleFast240x135GraphStyle = def{
	name: "grayscale-fast-240x135-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleFast240x135MinimalStyle = def{
	name: "grayscale-fast-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow240x135Style = def{
	name: "color-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
}

var ColorSlow240x135OverviewStyle = def{
	name: "color-slow-240x135-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow240x135DetailStyle = def{
	name: "color-slow-240x135-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow240x135GraphStyle = def{
	name: "color-slow-240x135-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow240x135MinimalStyle = def{
	name: "color-slow-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast240x135Style = def{
	name: "color-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
}

var ColorFast240x135OverviewStyle = def{
	name: "color-240x135-overview",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorFast240x135TimegraphStyle = def{
	name: "color-240x135-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var ColorFast240x135DetailStyle = def{
	name: "color-240x135-detail",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorFast240x135GraphStyle = def{
	name: "color-240x135-graph",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorFast240x135MinimalStyle = def{
	name: "color-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
