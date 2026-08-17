package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 128x296 panels (e-ink portrait).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var GrayscaleSlow128x296Style = def{
	name: "grayscale-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow128x296ThermometerStyle = def{
	name: "grayscale-slow-128x296-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow128x296SparkStyle = def{
	name: "grayscale-slow-128x296-spark",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow128x296HeatmapStyle = def{
	name: "grayscale-slow-128x296-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow128x296LEDsStyle = def{
	name: "grayscale-slow-128x296-leds",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow128x296AvgThermoStyle = def{
	name: "grayscale-slow-128x296-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow128x296MinimalStyle = def{
	name: "grayscale-slow-128x296-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow128x296Style = def{
	name: "color-slow-128x296",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
}

var ColorSlow128x296ThermometerStyle = def{
	name: "color-slow-128x296-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow128x296SparkStyle = def{
	name: "color-slow-128x296-spark",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow128x296HeatmapStyle = def{
	name: "color-slow-128x296-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow128x296LEDsStyle = def{
	name: "color-slow-128x296-leds",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow128x296AvgThermoStyle = def{
	name: "color-slow-128x296-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow128x296MinimalStyle = def{
	name: "color-slow-128x296-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 296, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
