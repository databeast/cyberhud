package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 80x160 panels (portrait).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var MonoSlow80x160Style = def{
	name: "mono-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
}

var MonoSlow80x160ThermometerStyle = def{
	name: "mono-slow-80x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var MonoSlow80x160SparkStyle = def{
	name: "mono-slow-80x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var MonoSlow80x160HeatmapStyle = def{
	name: "mono-slow-80x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var MonoSlow80x160LEDsStyle = def{
	name: "mono-slow-80x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var MonoSlow80x160AvgThermoStyle = def{
	name: "mono-slow-80x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var MonoSlow80x160MinimalStyle = def{
	name: "mono-slow-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow80x160Style = def{
	name: "grayscale-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow80x160ThermometerStyle = def{
	name: "grayscale-slow-80x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow80x160SparkStyle = def{
	name: "grayscale-slow-80x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow80x160HeatmapStyle = def{
	name: "grayscale-slow-80x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow80x160LEDsStyle = def{
	name: "grayscale-slow-80x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow80x160AvgThermoStyle = def{
	name: "grayscale-slow-80x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow80x160MinimalStyle = def{
	name: "grayscale-slow-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast80x160Style = def{
	name: "grayscale-fast-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
}

var GrayscaleFast80x160ThermometerStyle = def{
	name: "grayscale-fast-80x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleFast80x160SparkStyle = def{
	name: "grayscale-fast-80x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleFast80x160HeatmapStyle = def{
	name: "grayscale-fast-80x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleFast80x160LEDsStyle = def{
	name: "grayscale-fast-80x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleFast80x160AvgThermoStyle = def{
	name: "grayscale-fast-80x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleFast80x160MinimalStyle = def{
	name: "grayscale-fast-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow80x160Style = def{
	name: "color-slow-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
}

var ColorSlow80x160ThermometerStyle = def{
	name: "color-slow-80x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow80x160SparkStyle = def{
	name: "color-slow-80x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow80x160HeatmapStyle = def{
	name: "color-slow-80x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow80x160LEDsStyle = def{
	name: "color-slow-80x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow80x160AvgThermoStyle = def{
	name: "color-slow-80x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow80x160MinimalStyle = def{
	name: "color-slow-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast80x160Style = def{
	name: "color-80x160",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
}

var ColorFast80x160ThermometerStyle = def{
	name: "color-80x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorFast80x160SparkStyle = def{
	name: "color-80x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorFast80x160HeatmapStyle = def{
	name: "color-80x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorFast80x160LEDsStyle = def{
	name: "color-80x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorFast80x160AvgThermoStyle = def{
	name: "color-80x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorFast80x160MinimalStyle = def{
	name: "color-80x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 80, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
