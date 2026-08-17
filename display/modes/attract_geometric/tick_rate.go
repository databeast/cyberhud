package attract_geometric

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

func init() {
	region.RegisterTickRate("attract_geometric", &geometricTickProvider{})
}

// geometricTickProvider implements region.TickRateProvider for the geometric attract mode.
// It returns a fixed 33ms interval (~30fps) for smooth animation.
type geometricTickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps).
func (geometricTickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}
