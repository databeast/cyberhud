package widgets

import "time"

// Renderable is the primary widget contract. Any type satisfying Renderable
// can be used with the Compositor and placement helpers.
// RenderFrame executes the widget's render logic for the current frame and
// returns a positioned Sprite, or nil if the widget cannot render (e.g.,
// invalid bounds).
type Renderable interface {
	RenderFrame() *Sprite
}

// Described provides introspectable metadata for registry and suppression.
// Widgets implementing Described can participate in the widget registry and
// be evaluated by suppression rules.
type Described interface {
	Describe() Descriptor
}

// Configurable allows runtime parameter updates between frames.
// The controlling display code calls Configure to change a widget's
// parameters without recreating the instance.
type Configurable interface {
	Configure(cfg interface{})
}

// Animated marks a widget as having time-dependent internal state.
// The controlling display code calls Tick once per frame with the elapsed
// duration since the previous frame, advancing the widget's animation state.
type Animated interface {
	Tick(elapsed time.Duration)
}

// Descriptor holds metadata about a widget type for introspection,
// registry, and suppression rule evaluation.
type Descriptor struct {
	Name         string   // Unique widget type name (e.g., "led", "progressbar").
	MinWidth     int      // Minimum pixel width for valid rendering.
	MinHeight    int      // Minimum pixel height for valid rendering.
	Capabilities []string // Feature tags (e.g., "eink-safe", "animated").
}
