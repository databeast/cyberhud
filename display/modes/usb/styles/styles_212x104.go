package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 212×104 panels (landscape).

// ── MonoSlow ──

var MonoSlow212x104Style = def{
	name: "mono-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoSlow},
}

var MonoSlow212x104MinimalStyle = def{
	name: "mono-slow-212x104-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast212x104Style = def{
	name: "mono-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoFast},
}

var MonoFast212x104MinimalStyle = def{
	name: "mono-fast-212x104-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow212x104Style = def{
	name: "grayscale-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow212x104MinimalStyle = def{
	name: "grayscale-slow-212x104-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast212x104Style = def{
	name: "grayscale-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleFast},
}

var GrayscaleFast212x104MinimalStyle = def{
	name: "grayscale-fast-212x104-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow212x104Style = def{
	name: "color-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
}

var ColorSlow212x104MinimalStyle = def{
	name: "color-slow-212x104-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast212x104Style = def{
	name: "color-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorFast},
}

var ColorFast212x104MinimalStyle = def{
	name: "color-fast-212x104-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
