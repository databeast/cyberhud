package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// System style declarations for 296x128 panels (landscape).
// Variant set: skeleton (adaptiveBuild) + default, compact, cores, top.

var ColorSlow296x128Style = def{
	name: "color-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow296x128DefaultStyle = def{
	name: "color-slow-296x128-default",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildDefaultStyle},
}

var ColorSlow296x128CompactStyle = def{
	name: "color-slow-296x128-compact",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCompactStyle},
}

var ColorSlow296x128CoresStyle = def{
	name: "color-slow-296x128-cores",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildCoresStyle},
}

var ColorSlow296x128TopStyle = def{
	name: "color-slow-296x128-top",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildTopStyle},
}

