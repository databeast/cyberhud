package attract_popsquares

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

func init() {
	region.RegisterTickRate("attract_popsquares", &popsquaresTickProvider{})
}

// maxTickElapsed caps the per-frame advancement to prevent visual jumps when
// frames are occasionally delayed (e.g., GC pause or I/O stall).
const maxTickElapsed = 80 * time.Millisecond

// popsquaresTickProvider implements region.TickRateProvider for the popsquares
// mode. The pop-and-fade animation is gentle, but it still wants a smooth frame
// rate so the cell fades read as continuous rather than stepping. 33ms (~30fps)
// matches the other attract modes.
type popsquaresTickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps).
func (popsquaresTickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}
