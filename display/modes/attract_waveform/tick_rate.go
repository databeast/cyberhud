package attract_waveform

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

func init() {
	region.RegisterTickRate("attract_waveform", &tickProvider{})
}

// tickProvider implements region.TickRateProvider for the attract_waveform mode.
// It returns a fixed 33ms interval (~30fps) for smooth waveform animation.
type tickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps). The waveform effect is
// continuously animated and needs rapid refresh for smooth scrolling.
func (tickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}

// maxTickElapsed caps the per-frame advancement to prevent visual jumps when
// frames are occasionally delayed (e.g., GC pause or I/O stall).
const maxTickElapsed = 80 * time.Millisecond
