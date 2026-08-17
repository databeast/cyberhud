package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 800×480 panels (landscape, w≥240).

// ── MonoSlow ──

var MonoSlow800x480Style = def{
	name: "mono-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
}

var MonoSlow800x480MinimalStyle = def{
	name: "mono-slow-800x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast800x480Style = def{
	name: "mono-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoFast},
}

var MonoFast800x480MinimalStyle = def{
	name: "mono-fast-800x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow800x480Style = def{
	name: "grayscale-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow800x480MinimalStyle = def{
	name: "grayscale-slow-800x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast800x480Style = def{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
}

var GrayscaleFast800x480MinimalStyle = def{
	name: "grayscale-fast-800x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow800x480Style = def{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
}

var ColorSlow800x480MinimalStyle = def{
	name: "color-slow-800x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast800x480Style = def{
	name: "color-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
}

var ColorFast800x480MinimalStyle = def{
	name: "color-fast-800x480-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
