package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 160×80 panels (landscape, 128≤w<240).

// ── MonoSlow ──

var MonoSlow160x80Style = def{
	name: "mono-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
}

var MonoSlow160x80MinimalStyle = def{
	name: "mono-slow-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast160x80Style = def{
	name: "mono-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoFast},
}

var MonoFast160x80MinimalStyle = def{
	name: "mono-fast-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow160x80Style = def{
	name: "grayscale-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow160x80MinimalStyle = def{
	name: "grayscale-slow-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast160x80Style = def{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
}

var GrayscaleFast160x80MinimalStyle = def{
	name: "grayscale-fast-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow160x80Style = def{
	name: "color-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
}

var ColorSlow160x80MinimalStyle = def{
	name: "color-slow-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast160x80Style = def{
	name: "color-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
}

var ColorFast160x80MinimalStyle = def{
	name: "color-fast-160x80-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
