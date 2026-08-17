package led

import (
	"image"
	"image/color"
	"time"

	"github.com/databeast/cyberhud/display/widgets/gradient"
)

// Shape represents the LED body geometry.
type Shape int

const (
	Circle        Shape = iota // Default; inscribed circle
	Square                     // Full body-area square
	Diamond                    // 45° rotated square (vertices at midpoints)
	RoundedSquare              // Square with 25% corner radius
)

// State represents the LED illumination state.
type State int

const (
	On      State = iota // Fully illuminated
	Off                  // Dimmed outline only
	Warning              // Warning color at full brightness
)

// Animation represents the type of animation effect applied to the LED.
type Animation int

const (
	NoAnimation Animation = iota // No animation effect
	Pulse                        // Sinusoidal brightness modulation
	Blink                        // On/Off toggle at period intervals
	Fade                         // Linear ramp 0→1→0
)

// ShineStyle represents the type of specular highlight drawn on the LED body.
type ShineStyle int

const (
	ShineNone     ShineStyle = iota // No highlight
	ShineDot                        // Small white circle
	ShineCrescent                   // Arc-shaped highlight
)

// Orientation represents the layout direction for LED groups.
type Orientation int

const (
	Horizontal Orientation = iota // Default; left-to-right
	Vertical                      // Top-to-bottom
)

// AnimationConfig holds parameters for LED animation effects.
type AnimationConfig struct {
	Type          Animation     // Which animation effect
	Period        time.Duration // Cycle period [100ms, 5000ms]; default varies by type
	MinBrightness float64       // Pulse only: floor brightness [0.0, 0.99]; default 0.3
}

// GradientConfig holds parameters for gradient fill on the LED body.
type GradientConfig struct {
	Stops []gradient.ColorStop // 2–16 stops; fewer than 2 → solid fallback
}

// GroupEntry defines per-entry overrides for an LED within a group.
// Zero-value fields inherit from the group-level Config.
type GroupEntry struct {
	State        State      // Per-entry state (zero = use group default)
	Foreground   color.RGBA // Per-entry color (zero = use group default)
	WarningColor color.RGBA // Per-entry warning color (zero = use group default)
	Shape        Shape      // Zero = use group shape
	BorderWidth  int        // Zero = use group border
	BorderColor  color.RGBA // Zero = use group border color
	GlowEnabled  bool       // Explicit override flag needed (bool zero = false)
	GlowRadius   int        // Zero = use group glow radius
}

// Config holds all parameters needed to render an LED indicator.
//
// Sentinel values:
//   - Brightness = -1.0: use discrete State (On/Off/Warning)
//   - ShineOpacity = 0: default to 255 (fully opaque)
//   - GlowRadius = 0: use default (30% of body radius) when GlowEnabled is true
type Config struct {
	// Shape
	Shape Shape // LED body shape; zero = Circle (default)

	// State
	State      State   // Discrete state: On, Off, Warning
	Brightness float64 // Continuous brightness; -1.0 = use discrete State

	// Geometry
	Diameter int             // Outer LED dimension in pixels (3–64)
	Bounds   image.Rectangle // Only Bounds.Min used for position

	// Colors
	Foreground   color.RGBA // LED fill color (zero → opaque green)
	Background   color.RGBA // Off-state body interior fill (zero → opaque black); does NOT fill entire output
	WarningColor color.RGBA // Warning state color (zero → amber)

	// Gradient (optional)
	Gradient *GradientConfig // nil or <2 stops → solid fill

	// Glow
	GlowEnabled bool // Whether glow is rendered
	GlowRadius  int  // Pixels beyond body edge; 0 = 30% of body radius

	// Border
	BorderWidth int        // 1–4 pixels; 0 = no border
	BorderColor color.RGBA // Zero → 50% gray when border width > 0

	// Shine
	ShineStyle   ShineStyle // Dot, Crescent, or None
	ShineOpacity uint8      // Alpha for shine pixels; 0 = 255 (fully opaque)

	// Animation
	Animation AnimationConfig

	// Group (optional; nil or empty = single LED)
	Group       []GroupEntry // 2–32 entries for group mode
	Orientation Orientation  // Group layout direction

	// Spacing (group mode)
	Spacing int // Pixels between LEDs in group; default 2, [0, 8]

	// Internal animation state (advanced by Tick)
	animElapsed time.Duration
}
