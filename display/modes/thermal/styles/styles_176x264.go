package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 176x264 panels (e-ink portrait).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var GrayscaleSlow176x264Style = def{
	name: "grayscale-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow176x264ThermometerStyle = def{
	name: "grayscale-slow-176x264-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow176x264SparkStyle = def{
	name: "grayscale-slow-176x264-spark",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow176x264HeatmapStyle = def{
	name: "grayscale-slow-176x264-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow176x264LEDsStyle = def{
	name: "grayscale-slow-176x264-leds",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow176x264AvgThermoStyle = def{
	name: "grayscale-slow-176x264-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow176x264MinimalStyle = def{
	name: "grayscale-slow-176x264-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow176x264Style = def{
	name: "color-slow-176x264",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
}

var ColorSlow176x264ThermometerStyle = def{
	name: "color-slow-176x264-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow176x264SparkStyle = def{
	name: "color-slow-176x264-spark",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow176x264HeatmapStyle = def{
	name: "color-slow-176x264-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow176x264LEDsStyle = def{
	name: "color-slow-176x264-leds",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow176x264AvgThermoStyle = def{
	name: "color-slow-176x264-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow176x264MinimalStyle = def{
	name: "color-slow-176x264-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 176, MinHeight: 264, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
