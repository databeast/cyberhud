package progressbar

import (
	"math"
	"time"
)

// validate clamps all Config fields in place per the design's clamping table.
// It mutates the provided Config pointer. After validate returns, all fields
// are guaranteed to be within their valid ranges.
func validate(cfg *Config) {
	// Value: NaN/Inf → 0.0, then clamp [0, 1]
	if math.IsNaN(cfg.Value) || math.IsInf(cfg.Value, 0) {
		cfg.Value = 0.0
	}
	if cfg.Value < 0.0 {
		cfg.Value = 0.0
	}
	if cfg.Value > 1.0 {
		cfg.Value = 1.0
	}

	// Orientation: unrecognized values → OrientHorizontal
	if cfg.Orientation != OrientHorizontal && cfg.Orientation != OrientVertical {
		cfg.Orientation = OrientHorizontal
	}

	// SegmentGap: clamp [1, 4]
	if cfg.SegmentGap < 1 {
		cfg.SegmentGap = 1
	}
	if cfg.SegmentGap > 4 {
		cfg.SegmentGap = 4
	}

	// SweepAngle: default 270 when 0, then clamp [90, 350]
	if cfg.SweepAngle == 0 {
		cfg.SweepAngle = 270
	}
	if cfg.SweepAngle < 90 {
		cfg.SweepAngle = 90
	}
	if cfg.SweepAngle > 350 {
		cfg.SweepAngle = 350
	}

	// BorderWidth: clamp [0, 16]
	if cfg.BorderWidth < 0 {
		cfg.BorderWidth = 0
	}
	if cfg.BorderWidth > 16 {
		cfg.BorderWidth = 16
	}

	// BorderWall: clamp [0, 16] and is applied in addition to BorderWidth.
	if cfg.BorderWall < 0 {
		cfg.BorderWall = 0
	}
	if cfg.BorderWall > 16 {
		cfg.BorderWall = 16
	}

	// Thickness: auto-compute when 0
	if cfg.Thickness == 0 {
		dx := cfg.Bounds.Dx()
		dy := cfg.Bounds.Dy()
		minDim := dx
		if dy < minDim {
			minDim = dy
		}
		radius := minDim / 2
		computed := int(0.15 * float64(radius))
		if computed < 2 {
			computed = 2
		}
		cfg.Thickness = computed
	}

	// Markers: truncate to first 8
	if len(cfg.Markers) > 8 {
		cfg.Markers = cfg.Markers[:8]
	}

	// Gradient.Stops: truncate to first 64
	if cfg.Gradient != nil && len(cfg.Gradient.Stops) > 64 {
		cfg.Gradient.Stops = cfg.Gradient.Stops[:64]
	}

	// Animation.Period: ≤0 → keep as-is (treated as disabled by animation renderer);
	// when positive, clamp [100ms, 10s]
	if cfg.Animation.Period > 0 {
		if cfg.Animation.Period < 100*time.Millisecond {
			cfg.Animation.Period = 100 * time.Millisecond
		}
		if cfg.Animation.Period > 10*time.Second {
			cfg.Animation.Period = 10 * time.Second
		}
	}

	// Animation.Speed:
	//   0 → type-specific default (60 for Shimmer, 30 for MarchingStripes)
	//   negative → disabled (leave as-is; animation renderer treats negative as disabled)
	//   positive → clamp [10, 500]
	if cfg.Animation.Speed == 0 {
		switch cfg.Animation.Type {
		case Shimmer:
			cfg.Animation.Speed = 60
		case MarchingStripes:
			cfg.Animation.Speed = 30
		}
	} else if cfg.Animation.Speed > 0 {
		if cfg.Animation.Speed < 10 {
			cfg.Animation.Speed = 10
		}
		if cfg.Animation.Speed > 500 {
			cfg.Animation.Speed = 500
		}
	}
}

// resolveOrientation is a pure function that takes a Config by value and returns
// a RenderGeometry struct containing the resolved spatial parameters for the
// current style+orientation combination.
func resolveOrientation(cfg Config) RenderGeometry {
	dx := cfg.Bounds.Dx()
	dy := cfg.Bounds.Dy()

	switch cfg.Style {
	case Linear, Segmented:
		if cfg.Orientation == OrientVertical {
			return RenderGeometry{
				PrimaryAxis:   dy,
				MinorAxis:     dx,
				FillDirection: FillBottomToTop,
				StartAngle:    0,
			}
		}
		return RenderGeometry{
			PrimaryAxis:   dx,
			MinorAxis:     dy,
			FillDirection: FillLeftToRight,
			StartAngle:    0,
		}

	case Ring, Pie:
		minDim := dx
		if dy < minDim {
			minDim = dy
		}
		startAngle := -math.Pi / 2.0 // 12-o'clock for horizontal
		if cfg.Orientation == OrientVertical {
			startAngle = math.Pi // 9-o'clock for vertical
		}
		return RenderGeometry{
			PrimaryAxis:   minDim,
			MinorAxis:     minDim,
			FillDirection: FillClockwise,
			StartAngle:    startAngle,
		}

	case Arc:
		minDim := dx
		if dy < minDim {
			minDim = dy
		}
		startAngle := -math.Pi / 2.0 // horizontal default
		if cfg.Orientation == OrientVertical {
			startAngle = math.Pi // vertical: bottom endpoint
		}
		return RenderGeometry{
			PrimaryAxis:   minDim,
			MinorAxis:     minDim,
			FillDirection: FillClockwise,
			StartAngle:    startAngle,
		}

	default:
		// Unrecognized style — return zero geometry (Render will return nil)
		return RenderGeometry{}
	}
}
