package styles

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

var GaugesDefaultStyle = def{
	name: "gauges-default",
	reqs: style.SurfaceRequirements{
		MinWidth:   1,
		MinHeight:  1,
		Capability: style.MonoFast,
	},
	p: Params{
		Shape:      0,
		LabelTier:  tiercatalog.TierNormal,
		PaddingPct: 0,
		TileGapPx:  1,
		Columns:    0,
	},
}

var GaugesMonoSlowStyle = def{
	name: "gauges-mono-slow",
	reqs: style.SurfaceRequirements{
		MinWidth:   1,
		MinHeight:  1,
		Capability: style.MonoSlow,
	},
	p: Params{
		Shape:      0,
		LabelTier:  tiercatalog.TierSmall,
		PaddingPct: 0,
		TileGapPx:  1,
	},
}

var GaugesMonoFastStyle = def{
	name: "gauges-mono-fast",
	reqs: style.SurfaceRequirements{
		MinWidth:   1,
		MinHeight:  1,
		Capability: style.MonoFast,
	},
	p: Params{
		Shape:      0,
		LabelTier:  tiercatalog.TierNormal,
		PaddingPct: 0,
		TileGapPx:  1,
	},
}

var GaugesGrayscaleSlowStyle = def{
	name: "gauges-grayscale-slow",
	reqs: style.SurfaceRequirements{
		MinWidth:   1,
		MinHeight:  1,
		Capability: style.GrayscaleSlow,
	},
	p: Params{
		Shape:      0,
		LabelTier:  tiercatalog.TierSmall,
		PaddingPct: 0,
		TileGapPx:  1,
	},
}

var GaugesGrayscaleFastStyle = def{
	name: "gauges-grayscale-fast",
	reqs: style.SurfaceRequirements{
		MinWidth:   1,
		MinHeight:  1,
		Capability: style.GrayscaleFast,
	},
	p: Params{
		Shape:      0,
		LabelTier:  tiercatalog.TierNormal,
		PaddingPct: 0,
		TileGapPx:  1,
	},
}

var GaugesColorSlowStyle = def{
	name: "gauges-color-slow",
	reqs: style.SurfaceRequirements{
		MinWidth:   1,
		MinHeight:  1,
		Capability: style.ColorSlow,
	},
	p: Params{
		Shape:      0,
		LabelTier:  tiercatalog.TierSmall,
		PaddingPct: 0,
		TileGapPx:  1,
	},
}

var GaugesColorFastStyle = def{
	name: "gauges-color-fast",
	reqs: style.SurfaceRequirements{
		MinWidth:   1,
		MinHeight:  1,
		Capability: style.ColorFast,
	},
	p: Params{
		Shape:      0,
		LabelTier:  tiercatalog.TierNormal,
		PaddingPct: 0,
		TileGapPx:  1,
	},
}
