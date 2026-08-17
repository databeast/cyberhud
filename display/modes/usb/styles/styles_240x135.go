package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 240×135 panels (landscape, w≥240).

// ── MonoSlow ──

var MonoSlow240x135Style = def{
	name: "mono-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
}

var MonoSlow240x135MinimalStyle = def{
	name: "mono-slow-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast240x135Style = def{
	name: "mono-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoFast},
}

var MonoFast240x135MinimalStyle = def{
	name: "mono-fast-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow240x135Style = def{
	name: "grayscale-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x135MinimalStyle = def{
	name: "grayscale-slow-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast240x135Style = def{
	name: "grayscale-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
}

var GrayscaleFast240x135MinimalStyle = def{
	name: "grayscale-fast-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow240x135Style = def{
	name: "color-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
}

var ColorSlow240x135MinimalStyle = def{
	name: "color-slow-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast240x135Style = def{
	name: "color-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
}

var ColorFast240x135MinimalStyle = def{
	name: "color-fast-240x135-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
