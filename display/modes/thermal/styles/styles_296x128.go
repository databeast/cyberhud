package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 296x128 panels (e-ink landscape).
// Note: GrayscaleSlow296x128Style is declared in styles.go with explicit BuildFn (polished).
// Landscape: width > height, width ≥ 128, height ≥ 64 → overview, detail, graph, minimal (slow only, no timegraph).

// GrayscaleSlow 296x128 variants (skeleton is in styles.go as polished)
var GrayscaleSlow296x128DetailStyle = def{
	name: "grayscale-slow-296x128-detail",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var GrayscaleSlow296x128GraphStyle = def{
	name: "grayscale-slow-296x128-graph",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var GrayscaleSlow296x128MinimalStyle = def{
	name: "grayscale-slow-296x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow296x128Style = def{
	name: "color-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow296x128OverviewStyle = def{
	name: "color-slow-296x128-overview",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildOverviewStyle},
}

var ColorSlow296x128DetailStyle = def{
	name: "color-slow-296x128-detail",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDetailStyle},
}

var ColorSlow296x128GraphStyle = def{
	name: "color-slow-296x128-graph",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildGraphStyle},
}

var ColorSlow296x128MinimalStyle = def{
	name: "color-slow-296x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
