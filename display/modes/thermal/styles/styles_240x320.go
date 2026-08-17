package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 240x320 panels (portrait).
// Note: ColorFast 240x320 polished variants (thermometer, spark, heatmap, leds, avg-thermo)
// are declared in styles.go with explicit BuildFn.
// Skeleton defaults for each capability class are declared here.
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var MonoSlow240x320Style = def{
	name: "mono-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
}

var MonoSlow240x320ThermometerStyle = def{
	name: "mono-slow-240x320-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var MonoSlow240x320SparkStyle = def{
	name: "mono-slow-240x320-spark",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var MonoSlow240x320HeatmapStyle = def{
	name: "mono-slow-240x320-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var MonoSlow240x320LEDsStyle = def{
	name: "mono-slow-240x320-leds",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var MonoSlow240x320AvgThermoStyle = def{
	name: "mono-slow-240x320-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var MonoSlow240x320MinimalStyle = def{
	name: "mono-slow-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.MonoSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleSlow240x320Style = def{
	name: "grayscale-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow240x320ThermometerStyle = def{
	name: "grayscale-slow-240x320-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow240x320SparkStyle = def{
	name: "grayscale-slow-240x320-spark",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow240x320HeatmapStyle = def{
	name: "grayscale-slow-240x320-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow240x320LEDsStyle = def{
	name: "grayscale-slow-240x320-leds",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow240x320AvgThermoStyle = def{
	name: "grayscale-slow-240x320-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow240x320MinimalStyle = def{
	name: "grayscale-slow-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var GrayscaleFast240x320Style = def{
	name: "grayscale-fast-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
}

var GrayscaleFast240x320ThermometerStyle = def{
	name: "grayscale-fast-240x320-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleFast240x320SparkStyle = def{
	name: "grayscale-fast-240x320-spark",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleFast240x320HeatmapStyle = def{
	name: "grayscale-fast-240x320-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleFast240x320LEDsStyle = def{
	name: "grayscale-fast-240x320-leds",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleFast240x320AvgThermoStyle = def{
	name: "grayscale-fast-240x320-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleFast240x320MinimalStyle = def{
	name: "grayscale-fast-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.GrayscaleFast},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow240x320Style = def{
	name: "color-slow-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
}

var ColorSlow240x320ThermometerStyle = def{
	name: "color-slow-240x320-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow240x320SparkStyle = def{
	name: "color-slow-240x320-spark",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow240x320HeatmapStyle = def{
	name: "color-slow-240x320-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow240x320LEDsStyle = def{
	name: "color-slow-240x320-leds",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow240x320AvgThermoStyle = def{
	name: "color-slow-240x320-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow240x320MinimalStyle = def{
	name: "color-slow-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorFast240x320Style = def{
	name: "color-240x320",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
}

// ColorFast 240x320 polished variants are in styles.go (thermometer, spark, heatmap, leds, avg-thermo).
// Add minimal variant here since it's not in the polished set.

var ColorFast240x320MinimalStyle = def{
	name: "color-240x320-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
	p:    Params{BuildFn: buildMinimalStyle},
}
