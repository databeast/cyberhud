package styles

import "github.com/databeast/cyberhud/display/style"

// USB style declarations for 250×122 panels (landscape, w≥240).

// ── MonoSlow ──

var MonoSlow250x122Style = def{
	name: "mono-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoSlow},
}

var MonoSlow250x122MinimalStyle = def{
	name: "mono-slow-250x122-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── MonoFast ──

var MonoFast250x122Style = def{
	name: "mono-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoFast},
}

var MonoFast250x122MinimalStyle = def{
	name: "mono-fast-250x122-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleSlow ──

var GrayscaleSlow250x122Style = def{
	name: "grayscale-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow250x122MinimalStyle = def{
	name: "grayscale-slow-250x122-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── GrayscaleFast ──

var GrayscaleFast250x122Style = def{
	name: "grayscale-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleFast},
}

var GrayscaleFast250x122MinimalStyle = def{
	name: "grayscale-fast-250x122-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorSlow ──

var ColorSlow250x122Style = def{
	name: "color-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
}

var ColorSlow250x122MinimalStyle = def{
	name: "color-slow-250x122-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimal},
}

// ── ColorFast ──

var ColorFast250x122Style = def{
	name: "color-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorFast},
}

var ColorFast250x122MinimalStyle = def{
	name: "color-fast-250x122-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimal},
}
