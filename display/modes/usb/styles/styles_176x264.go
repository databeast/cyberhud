package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 176×264 panels (portrait, 128≤w<240).

// ── MonoSlow ──

var MonoSlow176x264Style = def{
	name: "mono-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoSlow},
}

var MonoSlow176x264MinimalStyle = def{
	name: "mono-slow-176x264-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast176x264Style = def{
	name: "mono-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoFast},
}

var MonoFast176x264MinimalStyle = def{
	name: "mono-fast-176x264-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow176x264Style = def{
	name: "grayscale-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow176x264MinimalStyle = def{
	name: "grayscale-slow-176x264-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast176x264Style = def{
	name: "grayscale-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleFast},
}

var GrayscaleFast176x264MinimalStyle = def{
	name: "grayscale-fast-176x264-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow176x264Style = def{
	name: "color-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
}

var ColorSlow176x264MinimalStyle = def{
	name: "color-slow-176x264-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast176x264Style = def{
	name: "color-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorFast},
}

var ColorFast176x264MinimalStyle = def{
	name: "color-fast-176x264-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
