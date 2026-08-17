package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// USB style declarations for 128x296 panels (portrait e-ink).
// Width ≥ 128 → skeleton + minimal variant for all capability levels.

// ── MonoSlow ──

var MonoSlow128x296Style = def{
	name: "mono-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoSlow},
}

var MonoSlow128x296MinimalStyle = def{
	name: "mono-slow-128x296-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast128x296Style = def{
	name: "mono-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoFast},
}

var MonoFast128x296MinimalStyle = def{
	name: "mono-fast-128x296-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow128x296Style = def{
	name: "grayscale-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x296MinimalStyle = def{
	name: "grayscale-slow-128x296-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast128x296Style = def{
	name: "grayscale-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleFast},
}

var GrayscaleFast128x296MinimalStyle = def{
	name: "grayscale-fast-128x296-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow128x296Style = def{
	name: "color-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
}

var ColorSlow128x296MinimalStyle = def{
	name: "color-slow-128x296-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast128x296Style = def{
	name: "color-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorFast},
}

var ColorFast128x296MinimalStyle = def{
	name: "color-fast-128x296-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
