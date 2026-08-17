package attract_waveform

import (
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_waveform/source"
	"github.com/databeast/cyberhud/display/style"
)

type Policy = source.Policy

func DefaultPolicy() source.Policy { return source.DefaultPolicy() }

func SnapshotStyles() []style.Style[source.Data, source.Policy] {
	return waveformRegistry.Enumerate()
}

func ResetSnapshotState() {
	frameCounter = 0
	lastTick = time.Time{}
	animPhase = 0
}
