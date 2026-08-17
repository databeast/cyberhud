package attract_particles

import (
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_particles/source"
	"github.com/databeast/cyberhud/display/style"
)

func SnapshotStyles() []style.Style[source.Snapshot, source.Policy] {
	return particlesRegistry.Enumerate()
}

func ResetSnapshotState() {
	frameCounter = 0
	lastTick = time.Time{}
	source.ResetParticles()
}
