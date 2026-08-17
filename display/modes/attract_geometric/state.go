package attract_geometric

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_geometric/source"
)

// animState holds the package-level animation state, protected by a mutex.
// This is accessed by the render loop goroutine (BuildView) and the command
// handler goroutine (policy changes), following the existing attract mode pattern.
var animState = struct {
	sync.Mutex
	clusters        []source.ClusterConfig
	pendingClusters []source.ClusterConfig
	stashedClusters []source.ClusterConfig
	fragmentState   source.FragmentState
	perfState       source.PerfState
	time            float64
	frameCount      uint64
	lastTick        time.Time
	rng             *rand.Rand
	scaleFactor     float64
	panelW, panelH  int
	initialized     bool
	initFrames      int
}{}

var (
	clockNow = time.Now
	newRNG   = func() *rand.Rand { return rand.New(rand.NewSource(time.Now().UnixNano())) }
)

// initAnimState initializes the animation state for the given panel dimensions.
func initAnimState(panelW, panelH int) {
	animState.rng = newRNG()
	animState.scaleFactor = computeScaleFactor(panelW, panelH)
	animState.panelW = panelW
	animState.panelH = panelH

	policy := tunePolicyForUltraLowRes(GetPolicy(), panelW, panelH)
	baseCount := int(math.Round(float64(10) * policy.Density))
	allClusters := source.InitializeClusters(panelW, panelH, baseCount, animState.rng)

	animState.pendingClusters = allClusters
	animState.clusters = nil
	animState.stashedClusters = nil
	animState.fragmentState = source.FragmentState{}
	animState.perfState = source.PerfState{
		CurrentSquareCount:  len(allClusters),
		OriginalSquareCount: len(allClusters),
	}
	animState.time = 0
	animState.frameCount = 0
	animState.lastTick = time.Time{}
	animState.initialized = true
	animState.initFrames = 0
}

// tickFrame advances the animation state by one frame.
func tickFrame() {
	now := clockNow()

	// Compute elapsed, capped at 80ms.
	// On first call (no previous tick), use 16ms default.
	var elapsed time.Duration
	if animState.lastTick.IsZero() {
		elapsed = 16 * time.Millisecond
	} else {
		elapsed = now.Sub(animState.lastTick)
		if elapsed > source.MaxTickElapsed {
			elapsed = source.MaxTickElapsed
		}
	}
	animState.lastTick = now

	policy := tunePolicyForUltraLowRes(GetPolicy(), animState.panelW, animState.panelH)

	// Advance animation time: elapsed * 2.0 * Speed.
	animState.time += elapsed.Seconds() * source.SpeedMultiplier * policy.Speed
	animState.frameCount++

	// Deferred initialization (first 3 frames).
	if animState.initFrames < source.DeferredInitFrames && len(animState.pendingClusters) > 0 {
		animState.initFrames++
		remaining := source.DeferredInitFrames - animState.initFrames + 1
		count := int(math.Ceil(float64(len(animState.pendingClusters)) / float64(remaining)))
		if count > len(animState.pendingClusters) {
			count = len(animState.pendingClusters)
		}
		animState.clusters = append(animState.clusters, animState.pendingClusters[:count]...)
		animState.pendingClusters = animState.pendingClusters[count:]
	}

	// Determine active clusters (opacity >= 0.1 threshold).
	var activeClusters []source.ClusterConfig
	for _, cl := range animState.clusters {
		elapsedSinceSpawn := animState.time - cl.SpawnTime
		if elapsedSinceSpawn <= 0 {
			continue
		}
		fadeInFactor := math.Min(1.0, elapsedSinceSpawn/cl.FadeInDuration)
		active := false
		for _, sq := range cl.Squares {
			baseOp := source.ComputeFadeOpacity(animState.time, sq.PhaseOffset, sq.CycleDuration, sq.PeakOpacity)
			if baseOp*fadeInFactor >= source.ActiveOpacityThreshold {
				active = true
				break
			}
		}
		if active {
			activeClusters = append(activeClusters, cl)
		}
	}

	// Rectangle lifecycle: replace dead squares (max 1 per frame).
	replacements := 0
	for ci := range animState.clusters {
		if replacements >= source.MaxReplacementsPerFrame {
			break
		}
		cl := &animState.clusters[ci]
		elapsedSinceSpawn := animState.time - cl.SpawnTime
		if elapsedSinceSpawn <= 0 {
			continue
		}
		fadeInFactor := elapsedSinceSpawn / cl.FadeInDuration
		if fadeInFactor < 1.0 {
			continue
		}

		// Iterate from last to first.
		for si := len(cl.Squares) - 1; si >= 0 && replacements < source.MaxReplacementsPerFrame; si-- {
			sq := cl.Squares[si]
			age := animState.time / sq.CycleDuration
			if age > 1.0 {
				fadeOp := source.ComputeFadeOpacity(animState.time, sq.PhaseOffset, sq.CycleDuration, sq.PeakOpacity)
				if fadeOp < 0.01 {
					newSq := spawnReplacement(cl, si, animState.rng, animState.scaleFactor, animState.panelW, animState.panelH)
					cl.Squares[si] = newSq
					replacements++
				}
			}
		}
	}

	// Fragment lifecycle: remove expired fragments.
	alive := animState.fragmentState.ActiveFragments[:0]
	for _, f := range animState.fragmentState.ActiveFragments {
		if source.ComputeFragmentOpacity(f, animState.time) >= 0 {
			alive = append(alive, f)
		}
	}
	animState.fragmentState.ActiveFragments = alive

	// Spawn fragment if conditions are met.
	if shouldRenderFragments(animState.panelW, animState.panelH) && len(activeClusters) > 0 {
		source.SpawnFragment(&animState.fragmentState, animState.time, activeClusters, animState.panelW, animState.panelH, animState.rng, animState.scaleFactor)
	}

	// Performance evaluation.
	frameTimeMs := elapsed.Seconds() * 1000
	decision := evaluatePerformance(&animState.perfState, frameTimeMs)
	switch decision.Action {
	case PerfReduce:
		if decision.NewCount < len(animState.clusters) {
			stashed := make([]source.ClusterConfig, len(animState.clusters)-decision.NewCount)
			copy(stashed, animState.clusters[decision.NewCount:])
			animState.stashedClusters = append(animState.stashedClusters, stashed...)
			animState.clusters = animState.clusters[:decision.NewCount]
		}
	case PerfRestore:
		toRestore := decision.NewCount - len(animState.clusters)
		if toRestore > len(animState.stashedClusters) {
			toRestore = len(animState.stashedClusters)
		}
		if toRestore > 0 {
			animState.clusters = append(animState.clusters, animState.stashedClusters[:toRestore]...)
			animState.stashedClusters = animState.stashedClusters[toRestore:]
		}
	}
}

// spawnReplacement creates a replacement square for a dead rectangle in a cluster.
func spawnReplacement(cl *source.ClusterConfig, removedIdx int, rng *rand.Rand, scaleFactor float64, panelW, panelH int) source.SquareConfig {
	removed := cl.Squares[removedIdx]

	// Replacement size: (15 + rng()*40) * scaleFactor, with fallback clamp.
	newSize := (15 + rng.Float64()*40) * scaleFactor
	if scaledSizeMax(scaleFactor) < 20 {
		newSize = math.Max(4, newSize)
	}

	// Try 10 random anchor-grid positions, reject if within newSize*0.5 of removed.
	var offsetX, offsetY float64
	placed := false
	minDist := newSize * 0.5

	// Build list of remaining squares (excluding the one being replaced).
	remaining := make([]source.SquareConfig, 0, len(cl.Squares)-1)
	for i, sq := range cl.Squares {
		if i != removedIdx {
			remaining = append(remaining, sq)
		}
	}

	for attempt := 0; attempt < 10 && !placed; attempt++ {
		if len(remaining) == 0 {
			break
		}
		parentIdx := rng.Intn(len(remaining))
		parent := remaining[parentIdx]

		step := parent.Size / 4.0
		halfSize := parent.Size / 2.0
		gridX := rng.Intn(5)
		gridY := rng.Intn(5)
		localX := -halfSize + float64(gridX)*step
		localY := -halfSize + float64(gridY)*step

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

		dx := worldX - removed.OffsetX
		dy := worldY - removed.OffsetY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist >= minDist {
			offsetX = worldX
			offsetY = worldY
			placed = true
		}
	}

	if !placed {
		// Fallback: angle + radius from cluster center (0,0).
		angle := rng.Float64() * 2 * math.Pi
		var radius float64
		if scaledSizeMax(scaleFactor) < 20 {
			radius = 20 + rng.Float64()*15
		} else {
			radius = (60 + rng.Float64()*40) * scaleFactor
		}
		offsetX = radius * math.Cos(angle)
		offsetY = radius * math.Sin(angle)
	}

	// Clamp absolute position to panel bounds to prevent offset drift off-screen over time.
	clusterAbsX := (cl.CenterXPct / 100.0) * float64(panelW)
	clusterAbsY := (cl.CenterYPct / 100.0) * float64(panelH)
	absX := clusterAbsX + offsetX
	absY := clusterAbsY + offsetY
	if absX < 0 {
		absX = 0
	} else if absX > float64(panelW) {
		absX = float64(panelW)
	}
	if absY < 0 {
		absY = 0
	} else if absY > float64(panelH) {
		absY = float64(panelH)
	}
	offsetX = absX - clusterAbsX
	offsetY = absY - clusterAbsY

	// Generate a new square with randomized parameters, then override size and position.
	sq := source.GenerateSquare(rng, nil)
	sq.Size = newSize
	sq.OffsetX = offsetX
	sq.OffsetY = offsetY

	return sq
}
