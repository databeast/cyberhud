package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 320x240 panels (landscape).
// Landscape: width > height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal + timegraph for fast.
// Note: ColorFast 320x240 polished variants (overview, timegraph) are declared in styles.go with explicit BuildFn.

var MonoSlow320x240Style = def{
	name: "mono-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoSlow320x240OverviewStyle = def{
	name: "mono-slow-320x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var MonoSlow320x240DetailStyle = def{
	name: "mono-slow-320x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var MonoSlow320x240GraphStyle = def{
	name: "mono-slow-320x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var MonoSlow320x240MinimalStyle = def{
	name: "mono-slow-320x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow320x240Style = def{
	name: "grayscale-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow320x240OverviewStyle = def{
	name: "grayscale-slow-320x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow320x240DetailStyle = def{
	name: "grayscale-slow-320x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow320x240GraphStyle = def{
	name: "grayscale-slow-320x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow320x240MinimalStyle = def{
	name: "grayscale-slow-320x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast320x240Style = def{
	name: "grayscale-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
}

var GrayscaleFast320x240OverviewStyle = def{
	name: "grayscale-fast-320x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleFast320x240TimegraphStyle = def{
	name: "grayscale-fast-320x240-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var GrayscaleFast320x240DetailStyle = def{
	name: "grayscale-fast-320x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleFast320x240GraphStyle = def{
	name: "grayscale-fast-320x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleFast320x240MinimalStyle = def{
	name: "grayscale-fast-320x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow320x240Style = def{
	name: "color-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorSlow},
}

var ColorSlow320x240OverviewStyle = def{
	name: "color-slow-320x240-overview",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow320x240DetailStyle = def{
	name: "color-slow-320x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow320x240GraphStyle = def{
	name: "color-slow-320x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow320x240MinimalStyle = def{
	name: "color-slow-320x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast320x240Style = def{
	name: "color-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
}

// ColorFast 320x240 polished variants (overview, timegraph) are in styles.go.
// Add detail, graph, minimal here since they're not in the polished set.

var ColorFast320x240DetailStyle = def{
	name: "color-320x240-detail",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorFast320x240GraphStyle = def{
	name: "color-320x240-graph",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorFast320x240MinimalStyle = def{
	name: "color-320x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
