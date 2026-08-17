package source

import "time"

// Constants ported verbatim from the TypeScript source
// (website/src/animations/geometric/).
const (
	SizeMin = 20
	SizeMax = 120

	PeakOpacityMin = 0.05
	PeakOpacityMax = 0.6

	CycleDurationMin = 3.0
	CycleDurationMax = 10.0

	ClusterSquaresMin = 8
	ClusterSquaresMax = 20

	FadeInDurationMin = 1.0
	FadeInDurationMax = 3.0

	MinPhaseOffsetDiff = 0.5

	MaxActiveFragments = 3
	MinSpawnInterval   = 3.0

	ClusterProximityRadius = 100.0

	MaxTickElapsed     = 80 * time.Millisecond
	DeferredInitFrames = 3

	ActiveOpacityThreshold  = 0.1
	MaxReplacementsPerFrame = 1
	MinGlowSpreadPx         = 4

	WindowSize             = 10
	ReductionThresholdMs   = 33.0
	RestorationThresholdMs = 25.0
	MinSquareCount         = 4

	ClampedFrameThresholdMs = 100.0
	SpeedMultiplier         = 2.0
)

// AspectPool is the weighted pool of aspect ratios for rectangle generation.
// Square (1:1) has 6× weight; non-square ratios appear once each.
var AspectPool = []float64{
	1.0, 1.0, 1.0, 1.0, 1.0, 1.0,
	4.0 / 3.0, 3.0 / 4.0, 16.0 / 9.0, 9.0 / 16.0, 21.0 / 9.0, 9.0 / 21.0,
}
