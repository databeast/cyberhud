package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 80×160 panels (portrait, w≥64 and w<128).

// ── MonoSlow ──

var MonoSlow80x160Style = def{
	name: "mono-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
}

var MonoSlow80x160MinimalStyle = def{
	name: "mono-slow-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast80x160Style = def{
	name: "mono-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoFast},
}

var MonoFast80x160MinimalStyle = def{
	name: "mono-fast-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow80x160Style = def{
	name: "grayscale-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow80x160MinimalStyle = def{
	name: "grayscale-slow-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast80x160Style = def{
	name: "grayscale-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
}

var GrayscaleFast80x160MinimalStyle = def{
	name: "grayscale-fast-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow80x160Style = def{
	name: "color-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
}

var ColorSlow80x160MinimalStyle = def{
	name: "color-slow-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast80x160Style = def{
	name: "color-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
}

var ColorFast80x160MinimalStyle = def{
	name: "color-fast-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
