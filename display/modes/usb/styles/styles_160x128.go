package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 160×128 panels (landscape, 128≤w<240).

// ── MonoSlow ──

var MonoSlow160x128Style = def{
	name: "mono-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow160x128MinimalStyle = def{
	name: "mono-slow-160x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast160x128Style = def{
	name: "mono-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast160x128MinimalStyle = def{
	name: "mono-fast-160x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow160x128Style = def{
	name: "grayscale-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow160x128MinimalStyle = def{
	name: "grayscale-slow-160x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast160x128Style = def{
	name: "grayscale-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleFast},
}

var GrayscaleFast160x128MinimalStyle = def{
	name: "grayscale-fast-160x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow160x128Style = def{
	name: "color-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow160x128MinimalStyle = def{
	name: "color-slow-160x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast160x128Style = def{
	name: "color-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorFast},
}

var ColorFast160x128MinimalStyle = def{
	name: "color-fast-160x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
