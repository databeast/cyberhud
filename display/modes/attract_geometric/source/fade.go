package source

import (
	"math"
)

// ComputeFadeOpacity returns the sinusoidal fade opacity for a square at a given time.
// Formula: peakOpacity * (0.5 + 0.5*sin(2π*(time+phaseOffset)/cycleDuration - π/2))
// Result is clamped to [0, peakOpacity] to guard against floating-point edge cases.
func ComputeFadeOpacity(time, phaseOffset, cycleDuration, peakOpacity float64) float64 {
	raw := peakOpacity * (0.5 + 0.5*math.Sin(2*math.Pi*(time+phaseOffset)/cycleDuration-math.Pi/2))
	if raw < 0 {
		return 0
	}
	if raw > peakOpacity {
		return peakOpacity
	}
	return raw
}

// ComputeGlowMultiplier returns a distance-based glow intensity multiplier.
// Formula: clamp(1.0 - 0.4*(dist/boundingRadius), 0.6, 1.0)
// Returns 1.0 when boundingRadius ≤ 0.
func ComputeGlowMultiplier(dist, boundingRadius float64) float64 {
	if boundingRadius <= 0 {
		return 1.0
	}
	m := 1.0 - 0.4*(dist/boundingRadius)
	if m < 0.6 {
		return 0.6
	}
	if m > 1.0 {
		return 1.0
	}
	return m
}

// ComputeGlowColor returns a brighter variant of the input color for glow rendering.
// Same H and S, with L = min(baseL + 20, 100).
func ComputeGlowColor(c HSLColor) HSLColor {
	l := c.L + 20
	if l > 100 {
		l = 100
	}
	return HSLColor{H: c.H, S: c.S, L: l}
}
