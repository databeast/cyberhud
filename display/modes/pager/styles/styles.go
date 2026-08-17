package styles

import "github.com/databeast/cyberhud/display/style"

// Pager style declarations for supported fast and slow surfaces.
//
// Most entries share the core layout in core.go. Slow-refresh declarations use
// the zero-value Params; fast-refresh declarations set Params.Fast to retain
// smooth scrolling via snapshot OffsetY.

// Mono Slow declarations.
var (
	PagerMonoSlow122x250Style = def{
		name: "mono-slow-122x250",
		reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow176x264Style = def{
		name: "mono-slow-176x264",
		reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow200x200Style = def{
		name: "mono-slow-200x200",
		reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow212x104Style = def{
		name: "mono-slow-212x104",
		reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow296x128Style = def{
		name: "mono-slow-296x128",
		reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow400x300Style = def{
		name: "mono-slow-400x300",
		reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow480x800Style = def{
		name: "mono-slow-480x800",
		reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow800x480Style = def{
		name: "mono-slow-800x480",
		reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow104x212Style = def{
		name: "mono-slow-104x212",
		reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow250x122Style = def{
		name: "mono-slow-250x122",
		reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow128x296Style = def{
		name: "mono-slow-128x296",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow264x176Style = def{
		name: "mono-slow-264x176",
		reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerMonoSlow300x400Style = def{
		name: "mono-slow-300x400",
		reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoSlow, MinRows: 1, MinCharsPerLine: 1},
	}
)

// Mono Fast declarations.
var (
	PagerMonoFast128x32Style = def{
		name: "mono-128x32",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerMonoFast128x64Style = def{
		name: "mono-128x64",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerMonoFast128x128Style = def{
		name: "mono-128x128",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerMonoFast32x128Style = def{
		name: "mono-32x128",
		reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerMonoFast64x128Style = def{
		name: "mono-64x128",
		reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
)

// Grayscale Slow declarations.
var (
	PagerGrayscaleSlow122x250Style = def{
		name: "grayscale-slow-122x250",
		reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow176x264Style = def{
		name: "grayscale-slow-176x264",
		reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow200x200Style = def{
		name: "grayscale-slow-200x200",
		reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow212x104Style = def{
		name: "grayscale-slow-212x104",
		reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow296x128Style = def{
		name: "grayscale-slow-296x128",
		reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow400x300Style = def{
		name: "grayscale-slow-400x300",
		reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow480x800Style = def{
		name: "grayscale-slow-480x800",
		reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow800x480Style = def{
		name: "grayscale-slow-800x480",
		reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow104x212Style = def{
		name: "grayscale-slow-104x212",
		reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow250x122Style = def{
		name: "grayscale-slow-250x122",
		reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow128x296Style = def{
		name: "grayscale-slow-128x296",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow264x176Style = def{
		name: "grayscale-slow-264x176",
		reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerGrayscaleSlow300x400Style = def{
		name: "grayscale-slow-300x400",
		reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow, MinRows: 1, MinCharsPerLine: 1},
	}
)

// Grayscale Fast declarations.
var (
	PagerGrayscaleFast160x80Style = def{
		name: "grayscale-fast-160x80",
		reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast160x128Style = def{
		name: "grayscale-fast-160x128",
		reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast240x135Style = def{
		name: "grayscale-fast-240x135",
		reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast240x240Style = def{
		name: "grayscale-fast-240x240",
		reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast320x240Style = def{
		name: "grayscale-fast-320x240",
		reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast480x320Style = def{
		name: "grayscale-fast-480x320",
		reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast800x480Style = def{
		name: "grayscale-fast-800x480",
		reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast80x160Style = def{
		name: "grayscale-fast-80x160",
		reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast128x160Style = def{
		name: "grayscale-fast-128x160",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast135x240Style = def{
		name: "grayscale-fast-135x240",
		reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast240x320Style = def{
		name: "grayscale-fast-240x320",
		reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast320x480Style = def{
		name: "grayscale-fast-320x480",
		reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast480x800Style = def{
		name: "grayscale-fast-480x800",
		reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerGrayscaleFast128x128Style = def{
		name: "grayscale-fast-128x128",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
)

// Color Slow declarations.
var (
	PagerColorSlow122x250Style = def{
		name: "color-slow-122x250",
		reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow176x264Style = def{
		name: "color-slow-176x264",
		reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow200x200Style = def{
		name: "color-slow-200x200",
		reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow212x104Style = def{
		name: "color-slow-212x104",
		reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow296x128Style = def{
		name: "color-slow-296x128",
		reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow400x300Style = def{
		name: "color-slow-400x300",
		reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow480x800Style = def{
		name: "color-slow-480x800",
		reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow800x480Style = def{
		name: "color-slow-800x480",
		reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow104x212Style = def{
		name: "color-slow-104x212",
		reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow250x122Style = def{
		name: "color-slow-250x122",
		reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow128x296Style = def{
		name: "color-slow-128x296",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow264x176Style = def{
		name: "color-slow-264x176",
		reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
	PagerColorSlow300x400Style = def{
		name: "color-slow-300x400",
		reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow, MinRows: 1, MinCharsPerLine: 1},
	}
)

// Color Fast declarations.
var (
	PagerColorFast160x80Style = def{
		name: "color-160x80",
		reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast160x128Style = def{
		name: "color-160x128",
		reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast240x135Style = def{
		name: "color-240x135",
		reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast240x240Style = def{
		name: "color-240x240",
		reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast320x240Style = def{
		name: "color-320x240",
		reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast480x320Style = def{
		name: "color-480x320",
		reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast800x480Style = def{
		name: "color-800x480",
		reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast80x160Style = def{
		name: "color-80x160",
		reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast128x160Style = def{
		name: "color-128x160",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast135x240Style = def{
		name: "color-135x240",
		reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast240x320Style = def{
		name: "color-240x320",
		reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast320x480Style = def{
		name: "color-320x480",
		reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast480x800Style = def{
		name: "color-480x800",
		reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
	PagerColorFast128x128Style = def{
		name: "color-128x128",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast, MinRows: 1, MinCharsPerLine: 1},
		p:    Params{Fast: true},
	}
)
