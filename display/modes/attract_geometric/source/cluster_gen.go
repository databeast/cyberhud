package source

import (
	"math"
	"math/rand"
	"sort"
)

// randRange returns a random float64 in [min, max).
func randRange(rng *rand.Rand, min, max float64) float64 {
	return min + rng.Float64()*(max-min)
}

// computePhaseOffset generates a phase offset that differs from all constrained
// values by more than MinPhaseOffsetDiff seconds. It uses rejection sampling
// with up to 50 attempts, then falls back to a deterministic scan at 0.1s steps.
// If even the deterministic scan fails (very dense), a random value is returned.
func computePhaseOffset(rng *rand.Rand, cycleDuration float64, constrained []float64) float64 {
	if len(constrained) == 0 {
		return randRange(rng, 0, cycleDuration)
	}

	// Rejection sampling: up to 50 attempts.
	for attempt := 0; attempt < 50; attempt++ {
		candidate := randRange(rng, 0, cycleDuration)
		if phaseOffsetValid(candidate, constrained) {
			return candidate
		}
	}

	// Deterministic fallback: scan from 0 at 0.1 steps.
	for t := 0.0; t < cycleDuration; t += 0.1 {
		if phaseOffsetValid(t, constrained) {
			return t
		}
	}

	// Very dense case: just use a random value.
	return randRange(rng, 0, cycleDuration)
}

// phaseOffsetValid returns true if candidate differs from all constrained values
// by more than MinPhaseOffsetDiff.
func phaseOffsetValid(candidate float64, constrained []float64) bool {
	for _, c := range constrained {
		diff := candidate - c
		if diff < 0 {
			diff = -diff
		}
		if diff <= MinPhaseOffsetDiff {
			return false
		}
	}
	return true
}

// GenerateSquare produces a single SquareConfig with randomized parameters.
// constrainedPhaseOffsets, when non-empty, triggers rejection sampling to ensure
// the new square's phase offset is sufficiently different from existing ones.
func GenerateSquare(rng *rand.Rand, constrainedPhaseOffsets []float64) SquareConfig {
	// Color in green spectrum.
	color := HSLColor{
		H: randRange(rng, 100, 140),
		S: randRange(rng, 70, 100),
		L: randRange(rng, 20, 55),
	}

	size := randRange(rng, SizeMin, SizeMax)

	// 90% chance rotation=0, 10% chance rotation=45.
	var rotation float64
	if rng.Float64() < 0.1 {
		rotation = 45
	}

	cycleDuration := randRange(rng, CycleDurationMin, CycleDurationMax)
	peakOpacity := randRange(rng, PeakOpacityMin, PeakOpacityMax)
	aspect := AspectPool[rng.Intn(len(AspectPool))]

	// Phase offset with rejection sampling when constrained.
	phaseOffset := computePhaseOffset(rng, cycleDuration, constrainedPhaseOffsets)

	return SquareConfig{
		OffsetX:       0,
		OffsetY:       0,
		Size:          size,
		Aspect:        aspect,
		Rotation:      rotation,
		Color:         color,
		PhaseOffset:   phaseOffset,
		CycleDuration: cycleDuration,
		PeakOpacity:   peakOpacity,
	}
}

// generateCluster produces a cluster at the given center position with spawn time.
func generateCluster(rng *rand.Rand, cxPct, cyPct, spawn float64) ClusterConfig {
	count := ClusterSquaresMin + rng.Intn(ClusterSquaresMax-ClusterSquaresMin+1)

	// Generate squares.
	squares := make([]SquareConfig, count)
	for i := range squares {
		squares[i] = GenerateSquare(rng, nil)
	}

	// Override sizes with exponential decay.
	for i := range squares {
		var t float64
		if count > 1 {
			t = float64(i) / float64(count-1)
		}
		squares[i].Size = SizeMin + (SizeMax-SizeMin)*math.Pow(1-t, 3.5)
	}

	// Override phaseOffsets.
	for i := range squares {
		squares[i].PhaseOffset = rng.Float64() * squares[i].CycleDuration
	}

	// Sort largest-to-smallest + position on grid.
	positionSquaresOnGrid(rng, squares)

	// Compute bounding radius.
	var maxDist float64
	for _, sq := range squares {
		dist := math.Sqrt(sq.OffsetX*sq.OffsetX + sq.OffsetY*sq.OffsetY)
		if dist > maxDist {
			maxDist = dist
		}
	}

	return ClusterConfig{
		CenterXPct:     cxPct,
		CenterYPct:     cyPct,
		Squares:        squares,
		BoundingRadius: maxDist,
		SpawnTime:      spawn,
		FadeInDuration: randRange(rng, FadeInDurationMin, FadeInDurationMax),
	}
}

// initializeClusters generates grid-distributed clusters for the given panel.
func InitializeClusters(w, h int, baseCount int, rng *rand.Rand) []ClusterConfig {
	clusterCount := baseCount
	if clusterCount < 5 {
		clusterCount = 5
	}

	cols := int(math.Ceil(math.Sqrt(float64(clusterCount))))
	rows := int(math.Ceil(float64(clusterCount) / float64(cols)))

	cellWidth := 100.0 / float64(cols)
	cellHeight := 100.0 / float64(rows)

	clusters := make([]ClusterConfig, clusterCount)
	for i := 0; i < clusterCount; i++ {
		col := i % cols
		row := i / cols

		// Cell center + jitter (±30% of cell dimensions).
		centerX := cellWidth*(float64(col)+0.5) + (rng.Float64()*2-1)*0.3*cellWidth
		centerY := cellHeight*(float64(row)+0.5) + (rng.Float64()*2-1)*0.3*cellHeight

		// Clamp to [5%, 95%].
		centerX = clamp(centerX, 5, 95)
		centerY = clamp(centerY, 5, 95)

		spawnTime := float64(i) * randRange(rng, 0.1, 0.5)

		clusters[i] = generateCluster(rng, centerX, centerY, spawnTime)
	}

	return clusters
}

// clamp restricts v to [min, max].
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// positionSquaresOnGrid positions rectangles using anchor-point grid snapping.
// It sorts squares descending by size, places the first at origin, then snaps
// each subsequent square to a random anchor point on a random parent's 5×5 grid,
// offset by one of the new square's 4 corners chosen at random.
func positionSquaresOnGrid(rng *rand.Rand, squares []SquareConfig) {
	if len(squares) < 2 {
		return
	}

	// Sort descending by size.
	sort.Slice(squares, func(i, j int) bool {
		return squares[i].Size > squares[j].Size
	})

	// First square at origin.
	squares[0].OffsetX = 0
	squares[0].OffsetY = 0

	for i := 1; i < len(squares); i++ {
		// Select random parent from [0, i).
		parentIdx := rng.Intn(i)
		parent := squares[parentIdx]

		// Compute anchor point on parent's 5×5 grid (4×4 subdivision).
		step := parent.Size / 4.0
		halfSize := parent.Size / 2.0
		gridX := rng.Intn(5)
		gridY := rng.Intn(5)
		localX := -halfSize + float64(gridX)*step
		localY := -halfSize + float64(gridY)*step

		// Transform to world space based on parent rotation.
		var worldX, worldY float64
		if parent.Rotation == 45 {
			cos45 := math.Cos(math.Pi / 4)
			sin45 := math.Sin(math.Pi / 4)
			worldX = parent.OffsetX + localX*cos45 - localY*sin45
			worldY = parent.OffsetY + localX*sin45 + localY*cos45
		} else {
			worldX = parent.OffsetX + localX
			worldY = parent.OffsetY + localY
		}

		// Select random corner of new rectangle.
		halfH := squares[i].Size / 2.0
		corners := [4][2]float64{
			{-halfH, -halfH},
			{halfH, -halfH},
			{halfH, halfH},
			{-halfH, halfH},
		}
		corner := corners[rng.Intn(4)]

		// Position = anchor - corner offset.
		squares[i].OffsetX = worldX - corner[0]
		squares[i].OffsetY = worldY - corner[1]
	}
}
