package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 480x800 panels (portrait extra large).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var GrayscaleSlow480x800Style = def{
	name: "grayscale-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow480x800ThermometerStyle = def{
	name: "grayscale-slow-480x800-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow480x800SparkStyle = def{
	name: "grayscale-slow-480x800-spark",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow480x800HeatmapStyle = def{
	name: "grayscale-slow-480x800-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow480x800LEDsStyle = def{
	name: "grayscale-slow-480x800-leds",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow480x800AvgThermoStyle = def{
	name: "grayscale-slow-480x800-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow480x800MinimalStyle = def{
	name: "grayscale-slow-480x800-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast480x800Style = def{
	name: "grayscale-fast-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
}

var GrayscaleFast480x800ThermometerStyle = def{
	name: "grayscale-fast-480x800-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleFast480x800SparkStyle = def{
	name: "grayscale-fast-480x800-spark",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleFast480x800HeatmapStyle = def{
	name: "grayscale-fast-480x800-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleFast480x800LEDsStyle = def{
	name: "grayscale-fast-480x800-leds",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleFast480x800AvgThermoStyle = def{
	name: "grayscale-fast-480x800-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleFast480x800MinimalStyle = def{
	name: "grayscale-fast-480x800-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow480x800Style = def{
	name: "color-slow-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
}

var ColorSlow480x800ThermometerStyle = def{
	name: "color-slow-480x800-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow480x800SparkStyle = def{
	name: "color-slow-480x800-spark",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow480x800HeatmapStyle = def{
	name: "color-slow-480x800-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow480x800LEDsStyle = def{
	name: "color-slow-480x800-leds",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow480x800AvgThermoStyle = def{
	name: "color-slow-480x800-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow480x800MinimalStyle = def{
	name: "color-slow-480x800-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast480x800Style = def{
	name: "color-480x800",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
}

var ColorFast480x800ThermometerStyle = def{
	name: "color-480x800-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorFast480x800SparkStyle = def{
	name: "color-480x800-spark",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorFast480x800HeatmapStyle = def{
	name: "color-480x800-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorFast480x800LEDsStyle = def{
	name: "color-480x800-leds",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorFast480x800AvgThermoStyle = def{
	name: "color-480x800-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorFast480x800MinimalStyle = def{
	name: "color-480x800-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 480, MinHeight: 800, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
