package attract_particles

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

func init() {
	region.RegisterTickRate("attract_particles", &tickProvider{})
}

// tickProvider implements region.TickRateProvider for the particles mode.
type tickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps) for smooth particle animation.
func (tickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}

// Package-level frame state.
var (
	frameCounter uint64
	lastTick     time.Time
)

// maxTickElapsed caps the per-frame advancement to prevent visual jumps.
const maxTickElapsed = 80 * time.Millisecond
