package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 135x240 panels (portrait).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var MonoSlow135x240Style = def{
	name: "mono-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
}

var MonoSlow135x240ThermometerStyle = def{
	name: "mono-slow-135x240-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var MonoSlow135x240SparkStyle = def{
	name: "mono-slow-135x240-spark",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var MonoSlow135x240HeatmapStyle = def{
	name: "mono-slow-135x240-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var MonoSlow135x240LEDsStyle = def{
	name: "mono-slow-135x240-leds",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var MonoSlow135x240AvgThermoStyle = def{
	name: "mono-slow-135x240-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var MonoSlow135x240MinimalStyle = def{
	name: "mono-slow-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow135x240Style = def{
	name: "grayscale-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow135x240ThermometerStyle = def{
	name: "grayscale-slow-135x240-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow135x240SparkStyle = def{
	name: "grayscale-slow-135x240-spark",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow135x240HeatmapStyle = def{
	name: "grayscale-slow-135x240-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow135x240LEDsStyle = def{
	name: "grayscale-slow-135x240-leds",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow135x240AvgThermoStyle = def{
	name: "grayscale-slow-135x240-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow135x240MinimalStyle = def{
	name: "grayscale-slow-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast135x240Style = def{
	name: "grayscale-fast-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
}

var GrayscaleFast135x240ThermometerStyle = def{
	name: "grayscale-fast-135x240-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleFast135x240SparkStyle = def{
	name: "grayscale-fast-135x240-spark",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleFast135x240HeatmapStyle = def{
	name: "grayscale-fast-135x240-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleFast135x240LEDsStyle = def{
	name: "grayscale-fast-135x240-leds",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleFast135x240AvgThermoStyle = def{
	name: "grayscale-fast-135x240-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleFast135x240MinimalStyle = def{
	name: "grayscale-fast-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow135x240Style = def{
	name: "color-slow-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
}

var ColorSlow135x240ThermometerStyle = def{
	name: "color-slow-135x240-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow135x240SparkStyle = def{
	name: "color-slow-135x240-spark",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow135x240HeatmapStyle = def{
	name: "color-slow-135x240-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow135x240LEDsStyle = def{
	name: "color-slow-135x240-leds",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow135x240AvgThermoStyle = def{
	name: "color-slow-135x240-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow135x240MinimalStyle = def{
	name: "color-slow-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast135x240Style = def{
	name: "color-135x240",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
}

var ColorFast135x240ThermometerStyle = def{
	name: "color-135x240-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorFast135x240SparkStyle = def{
	name: "color-135x240-spark",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorFast135x240HeatmapStyle = def{
	name: "color-135x240-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorFast135x240LEDsStyle = def{
	name: "color-135x240-leds",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorFast135x240AvgThermoStyle = def{
	name: "color-135x240-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorFast135x240MinimalStyle = def{
	name: "color-135x240-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 135, MinHeight: 240, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
