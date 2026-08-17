package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 400×300 panels (landscape, w≥240).

// ── MonoSlow ──

var MonoSlow400x300Style = def{
	name: "mono-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoSlow},
}

var MonoSlow400x300MinimalStyle = def{
	name: "mono-slow-400x300-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast400x300Style = def{
	name: "mono-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoFast},
}

var MonoFast400x300MinimalStyle = def{
	name: "mono-fast-400x300-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow400x300Style = def{
	name: "grayscale-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow400x300MinimalStyle = def{
	name: "grayscale-slow-400x300-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast400x300Style = def{
	name: "grayscale-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
}

var GrayscaleFast400x300MinimalStyle = def{
	name: "grayscale-fast-400x300-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow400x300Style = def{
	name: "color-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
}

var ColorSlow400x300MinimalStyle = def{
	name: "color-slow-400x300-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast400x300Style = def{
	name: "color-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorFast},
}

var ColorFast400x300MinimalStyle = def{
	name: "color-fast-400x300-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
