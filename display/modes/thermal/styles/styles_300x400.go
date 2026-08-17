package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 300x400 panels (e-ink portrait).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var GrayscaleSlow300x400Style = def{
	name: "grayscale-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow300x400ThermometerStyle = def{
	name: "grayscale-slow-300x400-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow300x400SparkStyle = def{
	name: "grayscale-slow-300x400-spark",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow300x400HeatmapStyle = def{
	name: "grayscale-slow-300x400-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow300x400LEDsStyle = def{
	name: "grayscale-slow-300x400-leds",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow300x400AvgThermoStyle = def{
	name: "grayscale-slow-300x400-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow300x400MinimalStyle = def{
	name: "grayscale-slow-300x400-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow300x400Style = def{
	name: "color-slow-300x400",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
}

var ColorSlow300x400ThermometerStyle = def{
	name: "color-slow-300x400-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow300x400SparkStyle = def{
	name: "color-slow-300x400-spark",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow300x400HeatmapStyle = def{
	name: "color-slow-300x400-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow300x400LEDsStyle = def{
	name: "color-slow-300x400-leds",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow300x400AvgThermoStyle = def{
	name: "color-slow-300x400-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow300x400MinimalStyle = def{
	name: "color-slow-300x400-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 300, MinHeight: 400, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
