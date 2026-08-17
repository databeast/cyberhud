package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 128x128 panels (square).
// Square: width == height, width ≥ 128 → overview, detail, graph, minimal + timegraph for fast.
// Note: Mono128x128Style (MonoFast) is declared in styles.go with explicit BuildFn (polished detail).

var MonoSlow128x128Style = def{
	name: "mono-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow128x128OverviewStyle = def{
	name: "mono-slow-128x128-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var MonoSlow128x128DetailStyle = def{
	name: "mono-slow-128x128-detail",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var MonoSlow128x128GraphStyle = def{
	name: "mono-slow-128x128-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var MonoSlow128x128MinimalStyle = def{
	name: "mono-slow-128x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

// MonoFast 128x128 variants (skeleton is in styles.go as polished detail)
var MonoFast128x128OverviewStyle = def{
	name: "mono-128x128-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var MonoFast128x128TimegraphStyle = def{
	name: "mono-128x128-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var MonoFast128x128GraphStyle = def{
	name: "mono-128x128-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var MonoFast128x128MinimalStyle = def{
	name: "mono-128x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow128x128Style = def{
	name: "grayscale-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x128OverviewStyle = def{
	name: "grayscale-slow-128x128-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow128x128DetailStyle = def{
	name: "grayscale-slow-128x128-detail",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow128x128GraphStyle = def{
	name: "grayscale-slow-128x128-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow128x128MinimalStyle = def{
	name: "grayscale-slow-128x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast128x128Style = def{
	name: "grayscale-fast-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
}

var GrayscaleFast128x128OverviewStyle = def{
	name: "grayscale-fast-128x128-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleFast128x128TimegraphStyle = def{
	name: "grayscale-fast-128x128-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var GrayscaleFast128x128DetailStyle = def{
	name: "grayscale-fast-128x128-detail",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleFast128x128GraphStyle = def{
	name: "grayscale-fast-128x128-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleFast128x128MinimalStyle = def{
	name: "grayscale-fast-128x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow128x128Style = def{
	name: "color-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow128x128OverviewStyle = def{
	name: "color-slow-128x128-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow128x128DetailStyle = def{
	name: "color-slow-128x128-detail",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow128x128GraphStyle = def{
	name: "color-slow-128x128-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow128x128MinimalStyle = def{
	name: "color-slow-128x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast128x128Style = def{
	name: "color-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
}

var ColorFast128x128OverviewStyle = def{
	name: "color-128x128-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorFast128x128TimegraphStyle = def{
	name: "color-128x128-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var ColorFast128x128DetailStyle = def{
	name: "color-128x128-detail",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorFast128x128GraphStyle = def{
	name: "color-128x128-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorFast128x128MinimalStyle = def{
	name: "color-128x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
