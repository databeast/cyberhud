package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 135×240 panels (portrait, 128≤w<240).

// ── MonoSlow ──

var MonoSlow135x240Style = def{
	name: "mono-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoSlow135x240MinimalStyle = def{
	name: "mono-slow-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast135x240Style = def{
	name: "mono-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoFast},
}

var MonoFast135x240MinimalStyle = def{
	name: "mono-fast-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow135x240Style = def{
	name: "grayscale-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow135x240MinimalStyle = def{
	name: "grayscale-slow-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast135x240Style = def{
	name: "grayscale-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
}

var GrayscaleFast135x240MinimalStyle = def{
	name: "grayscale-fast-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow135x240Style = def{
	name: "color-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
}

var ColorSlow135x240MinimalStyle = def{
	name: "color-slow-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast135x240Style = def{
	name: "color-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
}

var ColorFast135x240MinimalStyle = def{
	name: "color-fast-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
