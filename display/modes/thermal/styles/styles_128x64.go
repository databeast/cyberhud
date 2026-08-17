package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 128x64 panels (landscape).
// Landscape: width ≥ height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal + timegraph for fast.
// Note: MonoSlow128x64Style is declared in styles.go with explicit BuildFn (polished).
// The MonoFast skeleton is declared here.

var MonoFast128x64Style = def{
	name: "mono-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
}

var MonoFast128x64OverviewStyle = def{
	name: "mono-128x64-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var MonoFast128x64TimegraphStyle = def{
	name: "mono-128x64-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var MonoFast128x64DetailStyle = def{
	name: "mono-128x64-detail",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var MonoFast128x64GraphStyle = def{
	name: "mono-128x64-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var MonoFast128x64MinimalStyle = def{
	name: "mono-128x64-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow128x64Style = def{
	name: "grayscale-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x64OverviewStyle = def{
	name: "grayscale-slow-128x64-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow128x64DetailStyle = def{
	name: "grayscale-slow-128x64-detail",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow128x64GraphStyle = def{
	name: "grayscale-slow-128x64-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow128x64MinimalStyle = def{
	name: "grayscale-slow-128x64-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow128x64Style = def{
	name: "color-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
}

var ColorSlow128x64OverviewStyle = def{
	name: "color-slow-128x64-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow128x64DetailStyle = def{
	name: "color-slow-128x64-detail",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow128x64GraphStyle = def{
	name: "color-slow-128x64-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow128x64MinimalStyle = def{
	name: "color-slow-128x64-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

// MonoSlow 128x64 variants (skeleton is in styles.go as polished)
var MonoSlow128x64OverviewStyle = def{
	name: "mono-slow-128x64-overview",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var MonoSlow128x64DetailStyle = def{
	name: "mono-slow-128x64-detail",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var MonoSlow128x64GraphStyle = def{
	name: "mono-slow-128x64-graph",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var MonoSlow128x64MinimalStyle = def{
	name: "mono-slow-128x64-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
