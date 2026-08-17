package attract_hacking

import (
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_hacking/source"
	"github.com/databeast/cyberhud/display/style"
)

func SnapshotStyles() []style.Style[source.Data, source.Policy] {
	return hackingRegistry.Enumerate()
}

func ResetSnapshotState() {
	frameCounter = 0
	lastTick = time.Time{}
	phase = 0
}
