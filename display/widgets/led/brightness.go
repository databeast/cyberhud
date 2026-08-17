package led

import (
	"math"
	"time"
)

// resolveBrightness converts the Config's Brightness and State fields into a
// single effective brightness value in [0.0, 1.0].
//
// When Brightness == -1.0 (sentinel), the discrete State field determines the
// effective brightness:
//   - On      → 1.0
//   - Off     → 0.0
//   - Warning → 1.0
//
// When Brightness is any other value, it is used directly (with a defensive
// clamp for NaN/Inf and out-of-range values that should already be caught by
// validate).
//
// This function does NOT apply animation curves — see resolveAnimation (task 8.1).
func resolveBrightness(cfg Config) float64 {
	if cfg.Brightness == -1.0 {
		switch cfg.State {
		case On:
			return 1.0
		case Off:
			return 0.0
		case Warning:
			return 1.0
		default:
			// Unrecognized state treated as Off (defensive; validate should
			// have already normalized this).
			return 0.0
		}
	}

	// Defensive checks: NaN or Inf → 0.0
	if math.IsNaN(cfg.Brightness) || math.IsInf(cfg.Brightness, 0) {
		return 0.0
	}

	// Clamp to [0.0, 1.0]
	if cfg.Brightness < 0.0 {
		return 0.0
	}
	if cfg.Brightness > 1.0 {
		return 1.0
	}

	return cfg.Brightness
}

// resolveAnimation applies animation curves to the base brightness value.
// It returns the effective brightness after applying the configured animation
// at the current elapsed time.
//
// When animation is disabled (NoAnimation, unrecognized type, or Period ≤ 0),
// baseBrightness is returned unchanged.
//
// Animation curves:
//   - Pulse: sinusoidal modulation between MinBrightness and 1.0
//   - Blink: alternates between 1.0 and 0.0 with equal half-periods
//   - Fade: linear ramp 0.0 → 1.0 → 0.0 over the configured period
func resolveAnimation(cfg Config, baseBrightness float64) float64 {
	// Animation disabled: NoAnimation, unrecognized type, or non-positive period.
	if cfg.Animation.Type == NoAnimation || cfg.Animation.Period <= 0 {
		return baseBrightness
	}

	// Only recognized animation types proceed.
	switch cfg.Animation.Type {
	case Pulse, Blink, Fade:
		// Valid — continue below.
	default:
		// Unrecognized animation type: treat as disabled.
		return baseBrightness
	}

	period := cfg.Animation.Period
	elapsed := cfg.animElapsed

	// Compute phase as fraction of period in [0.0, 1.0).
	// Use modulo to wrap elapsed time into a single cycle.
	periodNs := period.Nanoseconds()
	if periodNs <= 0 {
		return baseBrightness
	}
	elapsedNs := elapsed.Nanoseconds()
	// Handle negative elapsed gracefully (shouldn't happen, but defensive).
	if elapsedNs < 0 {
		elapsedNs = 0
	}
	cycleNs := elapsedNs % periodNs
	phase := float64(cycleNs) / float64(periodNs)

	switch cfg.Animation.Type {
	case Pulse:
		// Sinusoidal modulation between MinBrightness and 1.0.
		// At phase=0: cos(0)=1, (1-cos(0))/2 = 0, brightness = minBrightness
		// At phase=0.5: cos(π)=-1, (1-cos(π))/2 = 1, brightness = 1.0
		minB := cfg.Animation.MinBrightness
		brightness := minB + (1.0-minB)*(1.0-math.Cos(2.0*math.Pi*phase))/2.0
		return brightness

	case Blink:
		// Alternate on/off with equal half-periods.
		// At phase=0: on (brightness = 1.0).
		// phase < 0.5: on; phase >= 0.5: off.
		halfPeriodNs := periodNs / 2
		if cycleNs < halfPeriodNs {
			return 1.0
		}
		return 0.0

	case Fade:
		// Linear ramp: 0→1 in first half, 1→0 in second half.
		// At phase=0: brightness = 0.0
		// At phase=0.5: brightness = 1.0
		if phase < 0.5 {
			return phase * 2.0
		}
		return (1.0 - phase) * 2.0
	}

	return baseBrightness
}

// defaultAnimationPeriod returns the default period for each animation type.
// This is used during validation when Period is not explicitly set.
func defaultAnimationPeriod(animType Animation) time.Duration {
	switch animType {
	case Pulse:
		return 1000 * time.Millisecond
	case Blink:
		return 500 * time.Millisecond
	case Fade:
		return 2000 * time.Millisecond
	default:
		return 0
	}
}
