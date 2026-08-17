package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 300×400 panels (portrait, w≥240).

// ── MonoSlow ──

var MonoSlow300x400Style = def{
	name: "mono-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoSlow},
}

var MonoSlow300x400MinimalStyle = def{
	name: "mono-slow-300x400-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast300x400Style = def{
	name: "mono-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoFast},
}

var MonoFast300x400MinimalStyle = def{
	name: "mono-fast-300x400-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow300x400Style = def{
	name: "grayscale-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow300x400MinimalStyle = def{
	name: "grayscale-slow-300x400-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast300x400Style = def{
	name: "grayscale-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleFast},
}

var GrayscaleFast300x400MinimalStyle = def{
	name: "grayscale-fast-300x400-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow300x400Style = def{
	name: "color-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
}

var ColorSlow300x400MinimalStyle = def{
	name: "color-slow-300x400-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast300x400Style = def{
	name: "color-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorFast},
}

var ColorFast300x400MinimalStyle = def{
	name: "color-fast-300x400-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
