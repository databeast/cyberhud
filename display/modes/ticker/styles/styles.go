package styles

import "github.com/databeast/cyberhud/display/style"

// Ticker style declarations.
//
// Entries here are data declarations over the shared layouts in core.go:
// zero Params selects the fast/shared body, Params{Slow:true} selects the
// slow-refresh body, and Params{BuildFn: ...} attaches a bespoke layout.

var TickerMonoSlow122x250Style = def{
	name: "mono-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow176x264Style = def{
	name: "mono-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow200x200Style = def{
	name: "mono-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow212x104Style = def{
	name: "mono-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow296x128Style = def{
	name: "mono-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow400x300Style = def{
	name: "mono-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow480x800Style = def{
	name: "mono-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow104x212Style = def{
	name: "mono-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow250x122Style = def{
	name: "mono-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow128x296Style = def{
	name: "mono-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow264x176Style = def{
	name: "mono-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}

var TickerMonoSlow300x400Style = def{
	name: "mono-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoSlow},
	p:    Params{Slow: true},
}


var TickerMonoFast128x32Style = def{
	name: "mono-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
}

var TickerMonoFast128x64Style = def{
	name: "mono-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
}

var TickerMonoFast128x128Style = def{
	name: "mono-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
}

var TickerMonoFast32x128Style = def{
	name: "mono-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoFast},
}

var TickerMonoFast64x128Style = def{
	name: "mono-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
}

var TickerGrayscaleSlow122x250Style = def{
	name: "grayscale-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow176x264Style = def{
	name: "grayscale-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow200x200Style = def{
	name: "grayscale-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow212x104Style = def{
	name: "grayscale-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow296x128Style = def{
	name: "grayscale-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow400x300Style = def{
	name: "grayscale-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow480x800Style = def{
	name: "grayscale-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow800x480Style = def{
	name: "grayscale-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow104x212Style = def{
	name: "grayscale-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow250x122Style = def{
	name: "grayscale-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow128x296Style = def{
	name: "grayscale-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow264x176Style = def{
	name: "grayscale-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleSlow300x400Style = def{
	name: "grayscale-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
	p:    Params{Slow: true},
}

var TickerGrayscaleFast160x80Style = def{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast160x128Style = def{
	name: "grayscale-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast240x135Style = def{
	name: "grayscale-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast240x240Style = def{
	name: "grayscale-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast320x240Style = def{
	name: "grayscale-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast480x320Style = def{
	name: "grayscale-fast-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast800x480Style = def{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast80x160Style = def{
	name: "grayscale-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast128x160Style = def{
	name: "grayscale-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast135x240Style = def{
	name: "grayscale-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast240x320Style = def{
	name: "grayscale-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast320x480Style = def{
	name: "grayscale-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast480x800Style = def{
	name: "grayscale-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
}

var TickerGrayscaleFast128x128Style = def{
	name: "grayscale-fast-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
}

var TickerColorSlow122x250Style = def{
	name: "color-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow176x264Style = def{
	name: "color-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow200x200Style = def{
	name: "color-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow212x104Style = def{
	name: "color-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow296x128Style = def{
	name: "color-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow400x300Style = def{
	name: "color-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow480x800Style = def{
	name: "color-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow800x480Style = def{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow104x212Style = def{
	name: "color-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow250x122Style = def{
	name: "color-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow128x296Style = def{
	name: "color-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow264x176Style = def{
	name: "color-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorSlow300x400Style = def{
	name: "color-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{Slow: true},
}

var TickerColorFast128x128Style = def{
	name: "color-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
}

var TickerColorFast160x80Style = def{
	name: "color-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
}

var TickerColorFast160x128Style = def{
	name: "color-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorFast},
}

var TickerColorFast240x135Style = def{
	name: "color-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
}

var TickerColorFast240x240Style = def{
	name: "color-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
}

var TickerColorFast240x320Style = def{
	name: "color-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
}

var TickerColorFast320x240Style = def{
	name: "color-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
}

var TickerColorFast320x480Style = def{
	name: "color-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorFast},
}

var TickerColorFast480x320Style = def{
	name: "color-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
}

var TickerColorFast480x800Style = def{
	name: "color-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
}

var TickerColorFast800x480Style = def{
	name: "color-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
}

var TickerColorFast80x160Style = def{
	name: "color-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
}

var TickerColorFast128x160Style = def{
	name: "color-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
}

var TickerColorFast135x240Style = def{
	name: "color-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
}

var MonoSlow800x480Style = def{
	name: "mono-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: monoSlow800x480Build},
}
