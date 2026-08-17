package attract_shapes

import (
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_shapes/source"
	"github.com/databeast/cyberhud/display/style"
)

type Policy = source.Policy

func SnapshotStyles() []style.Style[source.Snapshot, source.Policy] {
	return shapesRegistry.Enumerate()
}

func ResetSnapshotState() {
	frameCounter = 0
	lastTick = time.Time{}
	source.ResetShapes()
}
