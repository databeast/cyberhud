package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 240×320 panels (portrait, w≥240).

// ── MonoSlow ──

var MonoSlow240x320Style = def{
	name: "mono-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoSlow240x320MinimalStyle = def{
	name: "mono-slow-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast240x320Style = def{
	name: "mono-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoFast},
}

var MonoFast240x320MinimalStyle = def{
	name: "mono-fast-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow240x320Style = def{
	name: "grayscale-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x320MinimalStyle = def{
	name: "grayscale-slow-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast240x320Style = def{
	name: "grayscale-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
}

var GrayscaleFast240x320MinimalStyle = def{
	name: "grayscale-fast-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow240x320Style = def{
	name: "color-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
}

var ColorSlow240x320MinimalStyle = def{
	name: "color-slow-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast240x320Style = def{
	name: "color-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
}

var ColorFast240x320MinimalStyle = def{
	name: "color-fast-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
