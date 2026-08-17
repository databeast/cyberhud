package attract_starfield

import (
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_starfield/source"
	"github.com/databeast/cyberhud/display/style"
)

func SnapshotStyles() []style.Style[source.Snapshot, source.Policy] {
	return starfieldRegistry.Enumerate()
}

func ResetSnapshotState() {
	frameCounter = 0
	lastTick = time.Time{}
	source.ResetStarfield()
}
