package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// USB style declarations for 128x32 panels (landscape narrow strip).
// Height < 64 → skeleton + minimal variant only.

// ── MonoSlow ──

var MonoSlow128x32Style = def{
	name: "mono-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
}

var MonoSlow128x32MinimalStyle = def{
	name: "mono-slow-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast128x32Style = def{
	name: "mono-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
}

var MonoFast128x32MinimalStyle = def{
	name: "mono-fast-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow128x32Style = def{
	name: "grayscale-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x32MinimalStyle = def{
	name: "grayscale-slow-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast128x32Style = def{
	name: "grayscale-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleFast},
}

var GrayscaleFast128x32MinimalStyle = def{
	name: "grayscale-fast-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow128x32Style = def{
	name: "color-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
}

var ColorSlow128x32MinimalStyle = def{
	name: "color-slow-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast128x32Style = def{
	name: "color-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorFast},
}

var ColorFast128x32MinimalStyle = def{
	name: "color-fast-128x32-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
