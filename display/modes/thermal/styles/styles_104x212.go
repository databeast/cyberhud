package styles

import (
	"github.com/databeast/cyberhud/display/style"
)

// Thermal style declarations for 104x212 panels (e-ink portrait).
// Portrait: height > width, height ≥ 128, width ≥ 64 → portrait variants + minimal.

var GrayscaleSlow104x212Style = def{
	name: "grayscale-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
}

var GrayscaleSlow104x212ThermometerStyle = def{
	name: "grayscale-slow-104x212-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var GrayscaleSlow104x212SparkStyle = def{
	name: "grayscale-slow-104x212-spark",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var GrayscaleSlow104x212HeatmapStyle = def{
	name: "grayscale-slow-104x212-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var GrayscaleSlow104x212LEDsStyle = def{
	name: "grayscale-slow-104x212-leds",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var GrayscaleSlow104x212AvgThermoStyle = def{
	name: "grayscale-slow-104x212-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var GrayscaleSlow104x212MinimalStyle = def{
	name: "grayscale-slow-104x212-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.GrayscaleSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}

var ColorSlow104x212Style = def{
	name: "color-slow-104x212",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
}

var ColorSlow104x212ThermometerStyle = def{
	name: "color-slow-104x212-thermometer",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitThermometerStyle},
}

var ColorSlow104x212SparkStyle = def{
	name: "color-slow-104x212-spark",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitSparkStyle},
}

var ColorSlow104x212HeatmapStyle = def{
	name: "color-slow-104x212-heatmap",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitHeatmapStyle},
}

var ColorSlow104x212LEDsStyle = def{
	name: "color-slow-104x212-leds",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitLEDsStyle},
}

var ColorSlow104x212AvgThermoStyle = def{
	name: "color-slow-104x212-avg-thermo",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildPortraitAvgThermoStyle},
}

var ColorSlow104x212MinimalStyle = def{
	name: "color-slow-104x212-minimal",
	reqs: style.SurfaceRequirements{MinWidth: 104, MinHeight: 212, Capability: style.ColorSlow},
	p:    Params{BuildFn: buildMinimalStyle},
}
