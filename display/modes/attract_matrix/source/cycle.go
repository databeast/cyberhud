package source

import "time"

// cyclePeriod is the full duration of one rain intensity cycle.
// The effect starts sparse/slow, ramps to full density/speed at the midpoint,
// then tapers back down.
const cyclePeriod = 45 * time.Second

// cycleStart records when the current cycle began.
var cycleStart time.Time

// cycleProgress returns the current parabolic intensity factor in [0, 1].
// The curve is an inverted parabola: 0 at cycle start, 1 at midpoint, 0 at end.
// Formula: 1 - (2t - 1)^2 where t is normalized progress [0, 1].
func CycleProgress() float64 {
	if cycleStart.IsZero() {
		cycleStart = time.Now()
	}

	elapsed := time.Since(cycleStart)
	t := float64(elapsed%cyclePeriod) / float64(cyclePeriod) // [0, 1)

	// Inverted parabola: peaks at t=0.5, zero at t=0 and t=1.
	x := 2*t - 1 // [-1, 1]
	return 1 - x*x
}

// cycleMinDensity is the density at the start/end of the cycle (sparse).
const cycleMinDensity = 0.1

// cycleMinSpeedFactor scales the policy speed range at cycle start/end.
const cycleMinSpeedFactor = 0.3

// cycleMaxSpeedFactor scales speed at peak density (150% of base).
const cycleMaxSpeedFactor = 1.5

// cycleMinTrailFraction is the fraction of the policy TrailLength used at cycle edges.
const cycleMinTrailFraction = 0.25

// applyCycle modifies the effective policy values based on the current
// position in the rain intensity cycle. Density and trail length are interpolated
// between their minimum (cycle edges) and the policy-configured maximum (cycle peak).
//
// Density is quantized to 0.05 steps and trail length to integer steps to avoid
// rebuilding the strip cache every frame.
func ApplyCycle(p Policy) Policy {
	intensity := CycleProgress() // 0 at edges, 1 at peak

	// Density: interpolate from cycleMinDensity to policy Density,
	// quantized to avoid constant strip cache rebuilds.
	rawDensity := cycleMinDensity + (p.Density-cycleMinDensity)*intensity
	p.Density = quantize(rawDensity, 0.05)
	if p.Density < 0.1 {
		p.Density = 0.1
	}

	// Trail length: interpolate from short (25% of policy) to full policy value.
	// Quantized to int naturally. Minimum 4 enforced by normalizePolicy.
	minTrail := float64(p.TrailLength) * cycleMinTrailFraction
	p.TrailLength = int(minTrail + (float64(p.TrailLength)-minTrail)*intensity)
	if p.TrailLength < 4 {
		p.TrailLength = 4
	}

	return p
}

// cycleSpeedFactor returns the current speed multiplier based on cycle position.
// Ranges from cycleMinSpeedFactor (at cycle edges) to cycleMaxSpeedFactor (at peak).
func CycleSpeedFactor() float64 {
	intensity := CycleProgress()
	return cycleMinSpeedFactor + (cycleMaxSpeedFactor-cycleMinSpeedFactor)*intensity
}

// quantize rounds v to the nearest step.
func quantize(v, step float64) float64 {
	return float64(int(v/step+0.5)) * step
}
