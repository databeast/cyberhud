package attract_geometric

import (
	"math/rand"
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_geometric/source"
	"github.com/databeast/cyberhud/display/modes/attract_geometric/styles"
	"github.com/databeast/cyberhud/display/style"
)

func SnapshotStyles() []style.Style[source.GeometricFrame, source.Policy] {
	return styles.Registry().Enumerate()
}

func ResetSnapshotState() {
	animState.Lock()
	animState.clusters = nil
	animState.pendingClusters = nil
	animState.stashedClusters = nil
	animState.fragmentState = source.FragmentState{}
	animState.perfState = source.PerfState{}
	animState.time = 0
	animState.frameCount = 0
	animState.lastTick = time.Time{}
	animState.rng = nil
	animState.scaleFactor = 0
	animState.panelW = 0
	animState.panelH = 0
	animState.initialized = false
	animState.initFrames = 0
	animState.Unlock()

	slowRefreshState.Lock()
	slowRefreshState.img = nil
	slowRefreshState.fingerprint = ""
	slowRefreshState.panelW = 0
	slowRefreshState.panelH = 0
	slowRefreshState.Unlock()

	clockNow = time.Now
	newRNG = func() *rand.Rand { return rand.New(rand.NewSource(time.Now().UnixNano())) }
	SetPolicy(DefaultPolicy())
}

func SetSnapshotDeterminism(seed int64, start time.Time, step time.Duration) {
	newRNG = func() *rand.Rand { return rand.New(rand.NewSource(seed)) }
	current := start
	clockNow = func() time.Time {
		t := current
		current = current.Add(step)
		return t
	}
}
