package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 400x300 panels (landscape).
// Landscape: width > height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal + timegraph for fast.
// Note: GrayscaleFast400x300Style is declared in styles.go with explicit BuildFn (polished overview).

var GrayscaleSlow400x300Style = def{
	name: "grayscale-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow400x300OverviewStyle = def{
	name: "grayscale-slow-400x300-overview",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var GrayscaleSlow400x300DetailStyle = def{
	name: "grayscale-slow-400x300-detail",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow400x300GraphStyle = def{
	name: "grayscale-slow-400x300-graph",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow400x300MinimalStyle = def{
	name: "grayscale-slow-400x300-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

// GrayscaleFast 400x300 variants (skeleton is in styles.go as polished overview)
var GrayscaleFast400x300TimegraphStyle = def{
	name: "grayscale-fast-400x300-timegraph",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildTimegraphStyle},
}

var GrayscaleFast400x300DetailStyle = def{
	name: "grayscale-fast-400x300-detail",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleFast400x300GraphStyle = def{
	name: "grayscale-fast-400x300-graph",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleFast400x300MinimalStyle = def{
	name: "grayscale-fast-400x300-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow400x300Style = def{
	name: "color-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
}

var ColorSlow400x300OverviewStyle = def{
	name: "color-slow-400x300-overview",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow400x300DetailStyle = def{
	name: "color-slow-400x300-detail",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow400x300GraphStyle = def{
	name: "color-slow-400x300-graph",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow400x300MinimalStyle = def{
	name: "color-slow-400x300-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
