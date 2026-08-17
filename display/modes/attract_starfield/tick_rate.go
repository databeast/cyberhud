package attract_starfield

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

// Package-level frame state.
var (
	frameCounter uint64
	lastTick     time.Time
)

func init() {
	region.RegisterTickRate("attract_starfield", &starfieldTickProvider{})
}

// starfieldTickProvider implements region.TickRateProvider for the starfield mode.
type starfieldTickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps). The starfield effect requires
// continuous animation at a consistent frame rate for smooth star movement.
func (starfieldTickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}
