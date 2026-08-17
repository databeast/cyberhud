package attract_plasma

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

func init() {
	region.RegisterTickRate("attract_plasma", &plasmaTickProvider{})
}

// maxTickElapsed caps the per-frame advancement to prevent visual jumps when
// frames are occasionally delayed (e.g., GC pause or I/O stall).
const maxTickElapsed = 80 * time.Millisecond

// plasmaTickProvider implements region.TickRateProvider for the plasma mode.
// It returns a fixed 33ms interval (~30fps) since the plasma effect requires
// continuous animation at a consistent frame rate for smooth morphing.
type plasmaTickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps). The plasma effect is
// continuously animated and needs rapid refresh to produce smooth color
// transitions and blob morphing.
func (plasmaTickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}
