package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 122×250 panels (portrait, w≥64 and w<128).

// ── MonoSlow ──

var MonoSlow122x250Style = def{
	name: "mono-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoSlow},
}

var MonoSlow122x250MinimalStyle = def{
	name: "mono-slow-122x250-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast122x250Style = def{
	name: "mono-fast-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoFast},
}

var MonoFast122x250MinimalStyle = def{
	name: "mono-fast-122x250-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow122x250Style = def{
	name: "grayscale-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow122x250MinimalStyle = def{
	name: "grayscale-slow-122x250-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast122x250Style = def{
	name: "grayscale-fast-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleFast},
}

var GrayscaleFast122x250MinimalStyle = def{
	name: "grayscale-fast-122x250-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow122x250Style = def{
	name: "color-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
}

var ColorSlow122x250MinimalStyle = def{
	name: "color-slow-122x250-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast122x250Style = def{
	name: "color-fast-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorFast},
}

var ColorFast122x250MinimalStyle = def{
	name: "color-fast-122x250-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
