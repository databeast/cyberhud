package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 104×212 panels (portrait, w≥64 and w<128).

// ── MonoSlow ──

var MonoSlow104x212Style = def{
	name: "mono-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoSlow},
}

var MonoSlow104x212MinimalStyle = def{
	name: "mono-slow-104x212-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast104x212Style = def{
	name: "mono-fast-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoFast},
}

var MonoFast104x212MinimalStyle = def{
	name: "mono-fast-104x212-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow104x212Style = def{
	name: "grayscale-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow104x212MinimalStyle = def{
	name: "grayscale-slow-104x212-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast104x212Style = def{
	name: "grayscale-fast-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleFast},
}

var GrayscaleFast104x212MinimalStyle = def{
	name: "grayscale-fast-104x212-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow104x212Style = def{
	name: "color-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
}

var ColorSlow104x212MinimalStyle = def{
	name: "color-slow-104x212-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast104x212Style = def{
	name: "color-fast-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorFast},
}

var ColorFast104x212MinimalStyle = def{
	name: "color-fast-104x212-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
