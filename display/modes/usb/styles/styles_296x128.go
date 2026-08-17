package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 296×128 panels (landscape, w≥240).

// ── MonoSlow ──

var MonoSlow296x128Style = def{
	name: "mono-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow296x128MinimalStyle = def{
	name: "mono-slow-296x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast296x128Style = def{
	name: "mono-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast296x128MinimalStyle = def{
	name: "mono-fast-296x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow296x128Style = def{
	name: "grayscale-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow296x128MinimalStyle = def{
	name: "grayscale-slow-296x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast296x128Style = def{
	name: "grayscale-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleFast},
}

var GrayscaleFast296x128MinimalStyle = def{
	name: "grayscale-fast-296x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow296x128Style = def{
	name: "color-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow296x128MinimalStyle = def{
	name: "color-slow-296x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast296x128Style = def{
	name: "color-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorFast},
}

var ColorFast296x128MinimalStyle = def{
	name: "color-fast-296x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
