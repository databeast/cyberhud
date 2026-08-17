package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 122x250 panels (e-ink portrait).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var GrayscaleSlow122x250Style = def{
	name: "grayscale-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow122x250ThermometerStyle = def{
	name: "grayscale-slow-122x250-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow122x250SparkStyle = def{
	name: "grayscale-slow-122x250-spark",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow122x250HeatmapStyle = def{
	name: "grayscale-slow-122x250-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow122x250LEDsStyle = def{
	name: "grayscale-slow-122x250-leds",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow122x250AvgThermoStyle = def{
	name: "grayscale-slow-122x250-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow122x250MinimalStyle = def{
	name: "grayscale-slow-122x250-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow122x250Style = def{
	name: "color-slow-122x250",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
}

var ColorSlow122x250ThermometerStyle = def{
	name: "color-slow-122x250-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow122x250SparkStyle = def{
	name: "color-slow-122x250-spark",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow122x250HeatmapStyle = def{
	name: "color-slow-122x250-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow122x250LEDsStyle = def{
	name: "color-slow-122x250-leds",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow122x250AvgThermoStyle = def{
	name: "color-slow-122x250-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow122x250MinimalStyle = def{
	name: "color-slow-122x250-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
