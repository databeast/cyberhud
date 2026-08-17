package led

import (
	"image/color"
	"math"
	"time"
)

// validate clamps all Config fields to their valid ranges and returns false
// when the configuration is irrecoverably invalid (Render should return nil).
// When it returns false, the caller should return nil immediately.
func validate(cfg *Config) bool {
	// --- Nil conditions (return false) ---

	// Diameter < 3 → invalid
	if cfg.Diameter < 3 {
		return false
	}

	// Group non-nil with 0 entries → invalid
	if cfg.Group != nil && len(cfg.Group) == 0 {
		return false
	}

	// --- Enum fallbacks ---

	// Unrecognized Shape → Circle
	if cfg.Shape < Circle || cfg.Shape > RoundedSquare {
		cfg.Shape = Circle
	}

	// Unrecognized State → Off
	if cfg.State < On || cfg.State > Warning {
		cfg.State = Off
	}

	// Unrecognized Animation type → NoAnimation
	if cfg.Animation.Type < NoAnimation || cfg.Animation.Type > Fade {
		cfg.Animation.Type = NoAnimation
	}

	// Unrecognized Orientation → Horizontal
	if cfg.Orientation < Horizontal || cfg.Orientation > Vertical {
		cfg.Orientation = Horizontal
	}

	// --- Brightness clamping ---
	// NaN/Inf → 0.0
	if math.IsNaN(cfg.Brightness) || math.IsInf(cfg.Brightness, 0) {
		cfg.Brightness = 0.0
	} else if cfg.Brightness < -1.0 {
		// < -1.0 → 0.0
		cfg.Brightness = 0.0
	} else if cfg.Brightness > -1.0 && cfg.Brightness < 0.0 {
		// in (-1.0, 0.0) → 0.0
		cfg.Brightness = 0.0
	} else if cfg.Brightness > 1.0 {
		// > 1.0 → 1.0
		cfg.Brightness = 1.0
	}
	// -1.0 and [0.0, 1.0] pass through unchanged

	// --- GlowRadius clamping to [0, 32] ---
	if cfg.GlowRadius < 0 {
		cfg.GlowRadius = 0
	} else if cfg.GlowRadius > 32 {
		cfg.GlowRadius = 32
	}

	// --- BorderWidth clamping ---
	// First clamp to [0, 4]
	if cfg.BorderWidth < 0 {
		cfg.BorderWidth = 0
	} else if cfg.BorderWidth > 4 {
		cfg.BorderWidth = 4
	}
	// Then clamp to floor(Body_Radius / 3)
	// Body_Radius = (Diameter - 2*BorderWidth) / 2, but at this stage we use
	// Diameter/2 as the body radius before border is applied for the cap calculation.
	bodyRadius := cfg.Diameter / 2
	maxBorder := bodyRadius / 3
	if cfg.BorderWidth > maxBorder {
		cfg.BorderWidth = maxBorder
	}

	// --- Animation.Period clamping ---
	if cfg.Animation.Type != NoAnimation {
		if cfg.Animation.Period <= 0 {
			// ≤0 → animation disabled
			cfg.Animation.Type = NoAnimation
		} else {
			if cfg.Animation.Period < 100*time.Millisecond {
				cfg.Animation.Period = 100 * time.Millisecond
			} else if cfg.Animation.Period > 5000*time.Millisecond {
				cfg.Animation.Period = 5000 * time.Millisecond
			}
		}
	}

	// --- Animation.MinBrightness clamping to [0.0, 0.99] ---
	if cfg.Animation.MinBrightness < 0.0 {
		cfg.Animation.MinBrightness = 0.0
	} else if cfg.Animation.MinBrightness > 0.99 {
		cfg.Animation.MinBrightness = 0.99
	}

	// --- Spacing clamping to [0, 8] ---
	if cfg.Spacing < 0 {
		cfg.Spacing = 0
	} else if cfg.Spacing > 8 {
		cfg.Spacing = 8
	}

	// --- Gradient.Stops: discard NaN/Inf positions, truncate to 16 ---
	if cfg.Gradient != nil && len(cfg.Gradient.Stops) > 0 {
		// Discard stops with NaN/Inf positions
		valid := cfg.Gradient.Stops[:0]
		for _, stop := range cfg.Gradient.Stops {
			if !math.IsNaN(stop.Position) && !math.IsInf(stop.Position, 0) {
				valid = append(valid, stop)
			}
		}
		cfg.Gradient.Stops = valid

		// Truncate to first 16
		if len(cfg.Gradient.Stops) > 16 {
			cfg.Gradient.Stops = cfg.Gradient.Stops[:16]
		}
	}

	// --- Group: truncate to first 32; treat 0-entry as nil ---
	if cfg.Group != nil {
		if len(cfg.Group) == 0 {
			cfg.Group = nil
		} else if len(cfg.Group) > 32 {
			cfg.Group = cfg.Group[:32]
		}
	}

	return true
}

// resolveColors applies default colors to zero-value color fields in the Config.
// This must be called before any rendering logic (dimming, gradient, glow, brightness scaling).
func resolveColors(cfg *Config) {
	zero := color.RGBA{}

	// Foreground: zero → opaque green
	if cfg.Foreground == zero {
		cfg.Foreground = color.RGBA{R: 0, G: 200, B: 0, A: 255}
	}

	// Background: zero → opaque black
	if cfg.Background == zero {
		cfg.Background = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}

	// WarningColor: zero → opaque amber
	if cfg.WarningColor == zero {
		cfg.WarningColor = color.RGBA{R: 255, G: 191, B: 0, A: 255}
	}

	// BorderColor: zero + border width > 0 → 50% gray
	if cfg.BorderColor == zero && cfg.BorderWidth > 0 {
		cfg.BorderColor = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	}
}
