package attract_matrix

import (
	"time"

	"github.com/databeast/cyberhud/display/modes/attract_matrix/source"
	"github.com/databeast/cyberhud/display/style"
)

type Policy = source.Policy
type MatrixSnapshot = source.MatrixSnapshot

func NormalizePolicy(p source.Policy) source.Policy { return normalizePolicy(p) }

func SnapshotPolicyForTest() map[string]interface{} { return matrixSnapshotter{}.SnapshotPolicy() }

func RestorePolicyForTest(data map[string]interface{}) error {
	return matrixSnapshotter{}.RestorePolicy(data)
}

func SnapshotStyles() []style.Style[source.MatrixSnapshot, source.Policy] {
	return matrixRegistry.Enumerate()
}

func ResetSnapshotState() {
	frameCounter = 0
	lastTick = time.Time{}
	cachedPanelTypeSet = false
	source.ResetLayoutCache()
	source.ResetFontCache()
	source.ResetCycleStart()
}

func SetCycleStartForTest(t time.Time) { source.SetCycleStart(t) }
