package attract_bokeh

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

func init() {
	region.RegisterTickRate("attract_bokeh", &tickProvider{})
}

// tickProvider implements region.TickRateProvider for the bokeh attract mode.
// It returns a fixed 33ms interval (~30fps) for smooth circle drift animation.
type tickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps).
func (tickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}

// maxTickElapsed caps the per-frame advancement to prevent visual jumps.
const maxTickElapsed = 80 * time.Millisecond
