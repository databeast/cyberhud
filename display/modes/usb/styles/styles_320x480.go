package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 320×480 panels (portrait, w≥240).

// ── MonoSlow ──

var MonoSlow320x480Style = def{
	name: "mono-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoSlow},
}

var MonoSlow320x480MinimalStyle = def{
	name: "mono-slow-320x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast320x480Style = def{
	name: "mono-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoFast},
}

var MonoFast320x480MinimalStyle = def{
	name: "mono-fast-320x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow320x480Style = def{
	name: "grayscale-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow320x480MinimalStyle = def{
	name: "grayscale-slow-320x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast320x480Style = def{
	name: "grayscale-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleFast},
}

var GrayscaleFast320x480MinimalStyle = def{
	name: "grayscale-fast-320x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow320x480Style = def{
	name: "color-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorSlow},
}

var ColorSlow320x480MinimalStyle = def{
	name: "color-slow-320x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast320x480Style = def{
	name: "color-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorFast},
}

var ColorFast320x480MinimalStyle = def{
	name: "color-fast-320x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
