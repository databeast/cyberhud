package styles

import "github.com/databeast/cyberhud/display/style"

// Systemd style declarations for all supported panels.
//
// Each entry is a hand-tweakable declaration over the core layouts in
// core.go: adjust Params to tune a display/panel combination, or attach a
// bespoke BuildFn for fully custom rendering.

var ColorStyle128x128 = def{
	name: "color-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle128x160 = def{
	name: "color-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle135x240 = def{
	name: "color-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle160x128 = def{
	name: "color-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle160x80 = def{
	name: "color-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle240x135 = def{
	name: "color-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle240x240 = def{
	name: "color-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle240x320 = def{
	name: "color-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle320x240 = def{
	name: "color-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle320x480 = def{
	name: "color-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle480x320 = def{
	name: "color-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle480x800 = def{
	name: "color-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorStyle800x480 = def{
	name: "color-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorFast},
	p:    Params{Gradient: true},
}

var ColorFast128x32Style = def{
	name: "color-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast128x64Style = def{
	name: "color-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast32x128Style = def{
	name: "color-fast-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast64x128Style = def{
	name: "color-fast-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast104x212Style = def{
	name: "color-fast-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast122x250Style = def{
	name: "color-fast-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast128x296Style = def{
	name: "color-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast176x264Style = def{
	name: "color-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast200x200Style = def{
	name: "color-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast212x104Style = def{
	name: "color-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast250x122Style = def{
	name: "color-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast264x176Style = def{
	name: "color-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast296x128Style = def{
	name: "color-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast300x400Style = def{
	name: "color-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorFast400x300Style = def{
	name: "color-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorFast},
	p:    Params{ColorText: true},
}

var ColorSlow128x32Style = def{
	name: "color-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow128x64Style = def{
	name: "color-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow128x128Style = def{
	name: "color-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow32x128Style = def{
	name: "color-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow64x128Style = def{
	name: "color-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow104x212Style = def{
	name: "color-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow122x250Style = def{
	name: "color-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow128x296Style = def{
	name: "color-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow176x264Style = def{
	name: "color-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow200x200Style = def{
	name: "color-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow212x104Style = def{
	name: "color-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow250x122Style = def{
	name: "color-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow264x176Style = def{
	name: "color-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow296x128Style = def{
	name: "color-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow300x400Style = def{
	name: "color-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow400x300Style = def{
	name: "color-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow480x800Style = def{
	name: "color-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow800x480Style = def{
	name: "color-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow160x80Style = def{
	name: "color-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow128x160Style = def{
	name: "color-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow160x128Style = def{
	name: "color-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow240x135Style = def{
	name: "color-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow135x240Style = def{
	name: "color-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow240x240Style = def{
	name: "color-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow240x320Style = def{
	name: "color-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow320x240Style = def{
	name: "color-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow320x480Style = def{
	name: "color-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var ColorSlow480x320Style = def{
	name: "color-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{ColorText: true},
}

var EinkStyle104x212 = def{
	name: "eink-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle122x250 = def{
	name: "eink-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle128x296 = def{
	name: "eink-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle176x264 = def{
	name: "eink-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle200x200 = def{
	name: "eink-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle212x104 = def{
	name: "eink-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle250x122 = def{
	name: "eink-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle264x176 = def{
	name: "eink-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle296x128 = def{
	name: "eink-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle300x400 = def{
	name: "eink-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle400x300 = def{
	name: "eink-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle480x800 = def{
	name: "eink-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoSlow},
	p:    Params{EinkSuppressed: true, ContentMaxRows: true, Static: true},
}

var EinkStyle800x480 = def{
	name: "eink-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
	p:    Params{BuildFn: einkPoster800x480},
}

var GrayscaleFast128x128Style = def{
	name: "grayscale-fast-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast128x160Style = def{
	name: "grayscale-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast135x240Style = def{
	name: "grayscale-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast160x128Style = def{
	name: "grayscale-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast160x80Style = def{
	name: "grayscale-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true, Static: true},
}

var GrayscaleFast240x135Style = def{
	name: "grayscale-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast240x240Style = def{
	name: "grayscale-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast240x320Style = def{
	name: "grayscale-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast320x240Style = def{
	name: "grayscale-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast320x480Style = def{
	name: "grayscale-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast480x320Style = def{
	name: "grayscale-fast-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast480x800Style = def{
	name: "grayscale-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast800x480Style = def{
	name: "grayscale-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleFast},
	p:    Params{SingleLabel: true},
}

var GrayscaleFast128x32Style = def{
	name: "grayscale-fast-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleFast},
}

var GrayscaleFast128x64Style = def{
	name: "grayscale-fast-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleFast},
}

var GrayscaleFast32x128Style = def{
	name: "grayscale-fast-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.GrayscaleFast},
}

var GrayscaleFast64x128Style = def{
	name: "grayscale-fast-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleFast},
}

var GrayscaleFast104x212Style = def{
	name: "grayscale-fast-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleFast},
}

var GrayscaleFast122x250Style = def{
	name: "grayscale-fast-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleFast},
}

var GrayscaleFast128x296Style = def{
	name: "grayscale-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleFast},
}

var GrayscaleFast176x264Style = def{
	name: "grayscale-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleFast},
}

var GrayscaleFast200x200Style = def{
	name: "grayscale-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleFast},
}

var GrayscaleFast212x104Style = def{
	name: "grayscale-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleFast},
}

var GrayscaleFast250x122Style = def{
	name: "grayscale-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleFast},
}

var GrayscaleFast264x176Style = def{
	name: "grayscale-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleFast},
}

var GrayscaleFast296x128Style = def{
	name: "grayscale-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleFast},
}

var GrayscaleFast300x400Style = def{
	name: "grayscale-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleFast},
}

var GrayscaleFast400x300Style = def{
	name: "grayscale-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleFast},
}

var GrayscaleSlow128x32Style = def{
	name: "grayscale-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x64Style = def{
	name: "grayscale-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x128Style = def{
	name: "grayscale-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow32x128Style = def{
	name: "grayscale-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow64x128Style = def{
	name: "grayscale-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow104x212Style = def{
	name: "grayscale-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow122x250Style = def{
	name: "grayscale-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x296Style = def{
	name: "grayscale-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow176x264Style = def{
	name: "grayscale-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow200x200Style = def{
	name: "grayscale-slow-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow212x104Style = def{
	name: "grayscale-slow-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow250x122Style = def{
	name: "grayscale-slow-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow264x176Style = def{
	name: "grayscale-slow-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow296x128Style = def{
	name: "grayscale-slow-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow300x400Style = def{
	name: "grayscale-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow400x300Style = def{
	name: "grayscale-slow-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow480x800Style = def{
	name: "grayscale-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow800x480Style = def{
	name: "grayscale-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow160x80Style = def{
	name: "grayscale-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x160Style = def{
	name: "grayscale-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow160x128Style = def{
	name: "grayscale-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x135Style = def{
	name: "grayscale-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow135x240Style = def{
	name: "grayscale-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x240Style = def{
	name: "grayscale-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x320Style = def{
	name: "grayscale-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow320x240Style = def{
	name: "grayscale-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow320x480Style = def{
	name: "grayscale-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow480x320Style = def{
	name: "grayscale-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.GrayscaleSlow},
}

var MonoStyle128x128 = def{
	name: "mono-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoFast},
}

var MonoStyle128x32 = def{
	name: "mono-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoFast},
	p:    Params{Summary: true},
}

var MonoStyle128x64 = def{
	name: "mono-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoFast},
}

var MonoStyle32x128 = def{
	name: "mono-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoFast},
}

var MonoStyle64x128 = def{
	name: "mono-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast104x212Style = def{
	name: "mono-fast-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.MonoFast},
}

var MonoFast122x250Style = def{
	name: "mono-fast-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoFast},
}

var MonoFast128x296Style = def{
	name: "mono-fast-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.MonoFast},
}

var MonoFast176x264Style = def{
	name: "mono-fast-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.MonoFast},
}

var MonoFast200x200Style = def{
	name: "mono-fast-200x200",
	reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoFast},
}

var MonoFast212x104Style = def{
	name: "mono-fast-212x104",
	reqs: style.SurfaceRequirements{MinWidth: 212, MinHeight: 104, Capability: style.MonoFast},
}

var MonoFast250x122Style = def{
	name: "mono-fast-250x122",
	reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoFast},
}

var MonoFast264x176Style = def{
	name: "mono-fast-264x176",
	reqs: style.SurfaceRequirements{MinWidth: 264, MinHeight: 176, Capability: style.MonoFast},
}

var MonoFast296x128Style = def{
	name: "mono-fast-296x128",
	reqs: style.SurfaceRequirements{MinWidth: 296, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast300x400Style = def{
	name: "mono-fast-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.MonoFast},
}

var MonoFast400x300Style = def{
	name: "mono-fast-400x300",
	reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 300, Capability: style.MonoFast},
}

var MonoFast160x80Style = def{
	name: "mono-fast-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoFast},
}

var MonoFast128x160Style = def{
	name: "mono-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoFast},
}

var MonoFast160x128Style = def{
	name: "mono-fast-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast240x135Style = def{
	name: "mono-fast-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoFast},
}

var MonoFast135x240Style = def{
	name: "mono-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoFast},
}

var MonoFast240x240Style = def{
	name: "mono-fast-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoFast},
}

var MonoFast240x320Style = def{
	name: "mono-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoFast},
}

var MonoFast320x240Style = def{
	name: "mono-fast-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoFast},
}

var MonoFast320x480Style = def{
	name: "mono-fast-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoFast},
}

var MonoFast480x320Style = def{
	name: "mono-fast-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoFast},
}

var MonoFast480x800Style = def{
	name: "mono-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoFast},
}

var MonoFast800x480Style = def{
	name: "mono-fast-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoFast},
}

var MonoSlow128x32Style = def{
	name: "mono-slow-128x32",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 32, Capability: style.MonoSlow},
}

var MonoSlow128x64Style = def{
	name: "mono-slow-128x64",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.MonoSlow},
}

var MonoSlow32x128Style = def{
	name: "mono-slow-32x128",
	reqs: style.SurfaceRequirements{MinWidth: 32, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow64x128Style = def{
	name: "mono-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow128x128Style = def{
	name: "mono-slow-128x128",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow160x80Style = def{
	name: "mono-slow-160x80",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 80, Capability: style.MonoSlow},
}

var MonoSlow128x160Style = def{
	name: "mono-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
}

var MonoSlow160x128Style = def{
	name: "mono-slow-160x128",
	reqs: style.SurfaceRequirements{MinWidth: 160, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow240x135Style = def{
	name: "mono-slow-240x135",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 135, Capability: style.MonoSlow},
}

var MonoSlow135x240Style = def{
	name: "mono-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoSlow240x240Style = def{
	name: "mono-slow-240x240",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoSlow240x320Style = def{
	name: "mono-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoSlow320x240Style = def{
	name: "mono-slow-320x240",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoSlow320x480Style = def{
	name: "mono-slow-320x480",
	reqs: style.SurfaceRequirements{MinWidth: 320, MinHeight: 480, Capability: style.MonoSlow},
}

var MonoSlow480x320Style = def{
	name: "mono-slow-480x320",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoSlow480x800Style = def{
	name: "mono-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.MonoSlow},
}

var MonoSlow800x480Style = def{
	name: "mono-slow-800x480",
	reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 480, Capability: style.MonoSlow},
}
