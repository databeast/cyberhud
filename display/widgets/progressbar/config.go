package progressbar

import (
	"image"
	"image/color"
	"time"

	"github.com/databeast/cyberhud/display/widgets/gradient"
)

// Style represents the visual variant of a progress bar.
type Style int

const (
	Linear    Style = iota // Rectangular bar (orientation determines fill direction)
	Pie                    // Solid circular fill (gradient ignored)
	Segmented              // Chunked bar (orientation determines direction)
	Ring                   // Hollow donut (orientation determines start angle)
	Arc                    // Partial circular gauge (orientation determines fill direction)
)

// Orientation controls the primary fill direction for all styles.
type Orientation int

const (
	OrientHorizontal Orientation = iota // Default; left-to-right / clockwise from 12-o'clock
	OrientVertical                      // Bottom-to-top / clockwise from 9-o'clock
)

// Animation enumerates supported animation effects.
type Animation int

const (
	NoAnimation     Animation = iota
	Pulse                     // Sinusoidal brightness modulation
	Shimmer                   // Translucent highlight sweep
	MarchingStripes           // Diagonal stripe scroll
)

// FillDir encodes the resolved fill direction.
type FillDir int

const (
	FillLeftToRight FillDir = iota
	FillBottomToTop
	FillClockwise // Circular shapes; actual start angle in StartAngle field
)

// RenderGeometry is returned by resolveOrientation and encapsulates the
// resolved spatial parameters for the current style+orientation combination.
type RenderGeometry struct {
	PrimaryAxis   int     // Length in pixels along the fill direction
	MinorAxis     int     // Cross-axis dimension in pixels
	FillDirection FillDir // Horizontal vs vertical for linear; start angle for circular
	StartAngle    float64 // Radians; relevant for Ring/Arc/Pie only; 0 for linear shapes
}

// ThresholdMarker defines a visual indicator drawn at a specific value position.
type ThresholdMarker struct {
	Value float64    // Position along bar axis [0.0, 1.0] (clamped)
	Color color.RGBA // Marker line color (zero → opaque white)
}

// GradientFill wraps a slice of gradient color stops for the fill region.
type GradientFill struct {
	Stops []gradient.ColorStop // 2–64 stops; fewer than 2 → solid fallback
}

// AnimationConfig holds animation effect parameters.
type AnimationConfig struct {
	Type   Animation     // Which animation effect
	Period time.Duration // For Pulse: cycle period [100ms, 10000ms], default 1000ms
	Speed  int           // For Shimmer/MarchingStripes: pixels/sec [10, 500]
}

// Config holds all parameters needed to render a progress bar.
type Config struct {
	// Shape
	Style       Style
	Orientation Orientation // Fill direction; zero value = OrientHorizontal (default)

	// Value
	Value float64 // [0.0, 1.0]; clamped; NaN/Inf → 0.0

	// Geometry
	Bounds image.Rectangle

	// Colors
	Foreground color.RGBA
	Background color.RGBA

	// Gradient (optional; nil or <2 stops → solid fill)
	Gradient *GradientFill

	// Segmented options
	SegmentCount int // 0 = auto-compute: floor(length / (4+gap))
	SegmentGap   int // Gap pixels [1, 4]; default 1

	// Ring/Arc options
	Thickness  int     // Arc/ring thickness in px (0 → max(2, 15% radius))
	SweepAngle float64 // Arc only: degrees [90, 350]; default 270

	// Caps and Border
	RoundedCaps bool
	BorderWidth int        // [0, 16]; 0 = no border
	BorderWall  int        // Additional wall thickness; 0 = default border only
	BorderColor color.RGBA // Zero → no border drawn

	// Threshold markers (max 8; excess ignored)
	Markers []ThresholdMarker

	// Animation
	Animation AnimationConfig

	// Internal animation state (advanced by Tick)
	animElapsed time.Duration
}
