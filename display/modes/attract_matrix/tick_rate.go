package attract_matrix

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

// Package-level frame state.
var (
	frameCounter uint64
	lastTick     time.Time
)

// maxTickElapsed caps the per-frame advancement to prevent visual jumps when
// frames are occasionally delayed (e.g., GC pause or I/O stall). Under normal
// operation the render loop fires at ~33ms intervals (via TickRateProvider),
// so this cap only activates on anomalous frame gaps.
const maxTickElapsed = 80 * time.Millisecond

func init() {
	region.RegisterTickRate("attract_matrix", &matrixTickProvider{})
}

// matrixTickProvider implements region.TickRateProvider for the matrix mode.
// It returns a fixed 33ms interval (~30fps) since the matrix rain effect
// requires continuous animation at a consistent frame rate for smooth scrolling.
type matrixTickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps). The matrix rain effect is
// continuously animated and needs rapid refresh to produce smooth pixel-level
// scrolling of the rain columns.
func (matrixTickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}
