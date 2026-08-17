package led

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/widgets"
)

// Render produces an LED indicator sprite from the given configuration.
// Returns nil when the configuration is irrecoverably invalid (e.g. Diameter < 3).
func Render(cfg Config) *widgets.Sprite {
	// Step 1: Validate and clamp fields; returns false if unrecoverable.
	if !validate(&cfg) {
		return nil
	}

	// Step 2: Apply default colors before any rendering logic.
	resolveColors(&cfg)

	// Step 3: Resolve effective brightness from discrete state or continuous value.
	effectiveBrightness := resolveBrightness(cfg)

	// Step 4: Apply animation curves to modulate effective brightness.
	effectiveBrightness = resolveAnimation(cfg, effectiveBrightness)

	// Step 5: Dispatch group vs single.
	if cfg.Group != nil && len(cfg.Group) > 0 {
		return renderGroup(cfg)
	}

	// Step 6: Single LED rendering.
	return renderSingle(cfg, effectiveBrightness)
}

// renderSingle renders a single LED indicator and returns a positioned Sprite.
func renderSingle(cfg Config, effectiveBrightness float64) *widgets.Sprite {
	// --- Compute output image dimensions ---
	glowRadius := 0
	if cfg.GlowEnabled {
		glowRadius = effectiveGlowRadius(cfg)
	}

	outputSize := cfg.Diameter + 2*glowRadius
	img := image.NewRGBA(image.Rect(0, 0, outputSize, outputSize))

	// --- Compute body rect (inset for border and offset for glow) ---
	bodyInset := cfg.BorderWidth + glowRadius
	bodyRect := image.Rect(bodyInset, bodyInset, outputSize-bodyInset, outputSize-bodyInset)

	// --- Determine fill color and off state ---
	isOff := effectiveBrightness == 0.0
	fillColor := cfg.Foreground
	if cfg.Brightness == -1.0 && cfg.State == Warning {
		fillColor = cfg.WarningColor
	}

	// --- Layer order: glow → border → body → shine ---

	// Glow pass
	if cfg.GlowEnabled && effectiveBrightness > 0.0 {
		applyGlow(img, cfg, glowRadius, effectiveBrightness)
	}

	// Border pass
	if cfg.BorderWidth > 0 {
		applyBorder(img, cfg, glowRadius)
	}

	// Body fill — use gradient fill when applicable, otherwise dispatch to shape renderer.
	if shouldApplyGradient(cfg, isOff) {
		applyGradientFill(img, bodyRect, cfg, effectiveBrightness)
	} else {
		dispatchShapeRenderer(img, bodyRect, fillColor, cfg.Background, effectiveBrightness, isOff, cfg.Shape)
	}

	// Shine pass (last layer)
	if cfg.ShineStyle != ShineNone && effectiveBrightness > 0.0 && cfg.Diameter >= 5 {
		applyShine(img, cfg, bodyRect, effectiveBrightness)
	}

	// --- Assign label ---
	label := assignLabel(cfg, effectiveBrightness)

	return &widgets.Sprite{
		Image:    img,
		Position: cfg.Bounds.Min,
		Label:    label,
	}
}

// dispatchShapeRenderer calls the appropriate shape body renderer based on cfg.Shape.
func dispatchShapeRenderer(img *image.RGBA, bodyRect image.Rectangle, fillColor, bgColor color.RGBA, brightness float64, isOff bool, shape Shape) {
	switch shape {
	case Square:
		renderSquareBody(img, bodyRect, fillColor, bgColor, brightness, isOff)
	case Diamond:
		renderDiamondBody(img, bodyRect, fillColor, bgColor, brightness, isOff)
	case RoundedSquare:
		renderRoundedSquareBody(img, bodyRect, fillColor, bgColor, brightness, isOff)
	default: // Circle (or any unrecognized, which validate already mapped to Circle)
		renderCircleBody(img, bodyRect, fillColor, bgColor, brightness, isOff)
	}
}

// assignLabel determines the Sprite label based on brightness and state.
//
// Logic:
//   - effectiveBrightness > 0.5 → "led/on"
//   - effectiveBrightness == 0.0 → "led/off"
//   - effectiveBrightness ∈ (0.0, 0.5] → "led/warning"
//   - Special cases when Brightness == -1.0 (discrete state mode):
//     State == On → "led/on"; State == Off → "led/off"; State == Warning → "led/warning"
func assignLabel(cfg Config, effectiveBrightness float64) string {
	// When using discrete state (sentinel -1.0), the state directly determines label.
	if cfg.Brightness == -1.0 {
		switch cfg.State {
		case On:
			return "led/on"
		case Off:
			return "led/off"
		case Warning:
			return "led/warning"
		}
	}

	// Continuous brightness mode.
	if effectiveBrightness > 0.5 {
		return "led/on"
	}
	if effectiveBrightness == 0.0 {
		return "led/off"
	}
	// effectiveBrightness ∈ (0.0, 0.5]
	return "led/warning"
}

// effectiveGlowRadius computes the actual glow radius in pixels.
// When GlowRadius is 0, defaults to 30% of Body_Radius (clamped to [1, 32]).
// Otherwise uses the configured value (already clamped by validate to [0, 32]).
func effectiveGlowRadius(cfg Config) int {
	if cfg.GlowRadius == 0 {
		bodyRadius := cfg.Diameter / 2
		r := int(float64(bodyRadius) * 0.3)
		if r < 1 {
			r = 1
		}
		if r > 32 {
			r = 32
		}
		return r
	}
	return cfg.GlowRadius
}
