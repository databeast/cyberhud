package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 64x128 panels (portrait).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var MonoSlow64x128Style = def{
	name: "mono-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
}

var MonoSlow64x128ThermometerStyle = def{
	name: "mono-slow-64x128-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var MonoSlow64x128SparkStyle = def{
	name: "mono-slow-64x128-spark",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var MonoSlow64x128HeatmapStyle = def{
	name: "mono-slow-64x128-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var MonoSlow64x128LEDsStyle = def{
	name: "mono-slow-64x128-leds",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var MonoSlow64x128AvgThermoStyle = def{
	name: "mono-slow-64x128-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var MonoSlow64x128MinimalStyle = def{
	name: "mono-slow-64x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var MonoFast64x128Style = def{
	name: "mono-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
}

var MonoFast64x128ThermometerStyle = def{
	name: "mono-64x128-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var MonoFast64x128SparkStyle = def{
	name: "mono-64x128-spark",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var MonoFast64x128HeatmapStyle = def{
	name: "mono-64x128-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var MonoFast64x128LEDsStyle = def{
	name: "mono-64x128-leds",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var MonoFast64x128AvgThermoStyle = def{
	name: "mono-64x128-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var MonoFast64x128MinimalStyle = def{
	name: "mono-64x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.MonoFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow64x128Style = def{
	name: "grayscale-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow64x128ThermometerStyle = def{
	name: "grayscale-slow-64x128-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow64x128SparkStyle = def{
	name: "grayscale-slow-64x128-spark",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow64x128HeatmapStyle = def{
	name: "grayscale-slow-64x128-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow64x128LEDsStyle = def{
	name: "grayscale-slow-64x128-leds",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow64x128AvgThermoStyle = def{
	name: "grayscale-slow-64x128-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow64x128MinimalStyle = def{
	name: "grayscale-slow-64x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow64x128Style = def{
	name: "color-slow-64x128",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
}

var ColorSlow64x128ThermometerStyle = def{
	name: "color-slow-64x128-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow64x128SparkStyle = def{
	name: "color-slow-64x128-spark",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow64x128HeatmapStyle = def{
	name: "color-slow-64x128-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow64x128LEDsStyle = def{
	name: "color-slow-64x128-leds",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow64x128AvgThermoStyle = def{
	name: "color-slow-64x128-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow64x128MinimalStyle = def{
	name: "color-slow-64x128-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 64, MinHeight: 128, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
