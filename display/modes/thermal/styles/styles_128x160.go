package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 128x160 panels (portrait).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var MonoSlow128x160Style = def{
	name: "mono-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
}

var MonoSlow128x160ThermometerStyle = def{
	name: "mono-slow-128x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var MonoSlow128x160SparkStyle = def{
	name: "mono-slow-128x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var MonoSlow128x160HeatmapStyle = def{
	name: "mono-slow-128x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var MonoSlow128x160LEDsStyle = def{
	name: "mono-slow-128x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var MonoSlow128x160AvgThermoStyle = def{
	name: "mono-slow-128x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var MonoSlow128x160MinimalStyle = def{
	name: "mono-slow-128x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow128x160Style = def{
	name: "grayscale-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x160ThermometerStyle = def{
	name: "grayscale-slow-128x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow128x160SparkStyle = def{
	name: "grayscale-slow-128x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow128x160HeatmapStyle = def{
	name: "grayscale-slow-128x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow128x160LEDsStyle = def{
	name: "grayscale-slow-128x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow128x160AvgThermoStyle = def{
	name: "grayscale-slow-128x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow128x160MinimalStyle = def{
	name: "grayscale-slow-128x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast128x160Style = def{
	name: "grayscale-fast-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
}

var GrayscaleFast128x160ThermometerStyle = def{
	name: "grayscale-fast-128x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleFast128x160SparkStyle = def{
	name: "grayscale-fast-128x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleFast128x160HeatmapStyle = def{
	name: "grayscale-fast-128x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleFast128x160LEDsStyle = def{
	name: "grayscale-fast-128x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleFast128x160AvgThermoStyle = def{
	name: "grayscale-fast-128x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleFast128x160MinimalStyle = def{
	name: "grayscale-fast-128x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow128x160Style = def{
	name: "color-slow-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
}

var ColorSlow128x160ThermometerStyle = def{
	name: "color-slow-128x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow128x160SparkStyle = def{
	name: "color-slow-128x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow128x160HeatmapStyle = def{
	name: "color-slow-128x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow128x160LEDsStyle = def{
	name: "color-slow-128x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow128x160AvgThermoStyle = def{
	name: "color-slow-128x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow128x160MinimalStyle = def{
	name: "color-slow-128x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast128x160Style = def{
	name: "color-128x160",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
}

var ColorFast128x160ThermometerStyle = def{
	name: "color-128x160-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorFast128x160SparkStyle = def{
	name: "color-128x160-spark",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorFast128x160HeatmapStyle = def{
	name: "color-128x160-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorFast128x160LEDsStyle = def{
	name: "color-128x160-leds",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorFast128x160AvgThermoStyle = def{
	name: "color-128x160-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorFast128x160MinimalStyle = def{
	name: "color-128x160-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
