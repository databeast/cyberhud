package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 264×176 panels (landscape, w≥240).

// ── MonoSlow ──

var MonoSlow264x176Style = def{
	name: "mono-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoSlow},
}

var MonoSlow264x176MinimalStyle = def{
	name: "mono-slow-264x176-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast264x176Style = def{
	name: "mono-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoFast},
}

var MonoFast264x176MinimalStyle = def{
	name: "mono-fast-264x176-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow264x176Style = def{
	name: "grayscale-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow264x176MinimalStyle = def{
	name: "grayscale-slow-264x176-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast264x176Style = def{
	name: "grayscale-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleFast},
}

var GrayscaleFast264x176MinimalStyle = def{
	name: "grayscale-fast-264x176-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow264x176Style = def{
	name: "color-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
}

var ColorSlow264x176MinimalStyle = def{
	name: "color-slow-264x176-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast264x176Style = def{
	name: "color-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorFast},
}

var ColorFast264x176MinimalStyle = def{
	name: "color-fast-264x176-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
