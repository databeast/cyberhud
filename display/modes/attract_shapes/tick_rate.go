package attract_shapes

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

func init() {
	region.RegisterTickRate("attract_shapes", &tickProvider{})
}

// tickProvider implements region.TickRateProvider for the attract_shapes mode.
// It returns a fixed 33ms interval (~30fps) for smooth shape animation.
type tickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps). The shapes effect is
// continuously animated and needs rapid refresh for smooth rotation and pulsing.
func (tickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}
