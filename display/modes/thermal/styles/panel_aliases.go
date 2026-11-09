package styles

import "github.com/databeast/cyberhud/display/style"

var (
	MonoSlow122x250OverviewStyle = def{
		name: "mono-slow-122x250-overview",
		reqs: style.SurfaceRequirements{MinWidth: 122, MinHeight: 250, Capability: style.MonoSlow},
		p:    Params{BuildFn: buildOverviewStyle},
	}
	MonoSlow200x200OverviewStyle = def{
		name: "mono-slow-200x200-overview",
		reqs: style.SurfaceRequirements{MinWidth: 200, MinHeight: 200, Capability: style.MonoSlow},
		p:    Params{BuildFn: buildOverviewStyle},
	}
	MonoSlow250x122OverviewStyle = def{
		name: "mono-slow-250x122-overview",
		reqs: style.SurfaceRequirements{MinWidth: 250, MinHeight: 122, Capability: style.MonoSlow},
		p:    Params{BuildFn: buildOverviewStyle},
	}
	MonoFast168x144OverviewStyle = def{
		name: "mono-fast-168x144-overview",
		reqs: style.SurfaceRequirements{MinWidth: 168, MinHeight: 144, Capability: style.MonoFast},
		p:    Params{BuildFn: buildOverviewStyle},
	}
	MonoFast400x240OverviewStyle = def{
		name: "mono-fast-400x240-overview",
		reqs: style.SurfaceRequirements{MinWidth: 400, MinHeight: 240, Capability: style.MonoFast},
		p:    Params{BuildFn: buildOverviewStyle},
	}
	GrayscaleFast128x64OverviewStyle = def{
		name: "grayscale-fast-128x64-overview",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 64, Capability: style.GrayscaleFast},
		p:    Params{BuildFn: buildOverviewStyle},
	}
	ColorFast128x160OverviewStyle = def{
		name: "color-fast-128x160-overview",
		reqs: style.SurfaceRequirements{MinWidth: 128, MinHeight: 160, Capability: style.ColorFast},
		p:    Params{BuildFn: buildOverviewStyle},
	}
	ColorFast240x320OverviewStyle = def{
		name: "color-fast-240x320-overview",
		reqs: style.SurfaceRequirements{MinWidth: 240, MinHeight: 320, Capability: style.ColorFast},
		p:    Params{BuildFn: buildOverviewStyle},
	}
)
