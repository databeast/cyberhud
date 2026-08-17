package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 200×200 panels (square).

// ── MonoSlow ──

var MonoSlow200x200Style = def{
	name: "mono-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoSlow},
}

var MonoSlow200x200MinimalStyle = def{
	name: "mono-slow-200x200-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast200x200Style = def{
	name: "mono-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoFast},
}

var MonoFast200x200MinimalStyle = def{
	name: "mono-fast-200x200-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow200x200Style = def{
	name: "grayscale-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow200x200MinimalStyle = def{
	name: "grayscale-slow-200x200-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast200x200Style = def{
	name: "grayscale-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleFast},
}

var GrayscaleFast200x200MinimalStyle = def{
	name: "grayscale-fast-200x200-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow200x200Style = def{
	name: "color-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
}

var ColorSlow200x200MinimalStyle = def{
	name: "color-slow-200x200-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast200x200Style = def{
	name: "color-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorFast},
}

var ColorFast200x200MinimalStyle = def{
	name: "color-fast-200x200-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
