package ticker

import (
	"time"

	"github.com/databeast/cyberhud/display/region"
)

func init() {
	region.RegisterTickRate("ticker", &tickerTickProvider{})
}

// tickerTickProvider implements region.TickRateProvider for the ticker mode.
// It reads the current ticker policy's AutoScrollMS value to derive the
// preferred tick interval.
type tickerTickProvider struct{}

// PreferredTickInterval returns the ticker mode's preferred tick interval
// based on the current AutoScrollMS policy value. If AutoScrollMS is <= 0,
// falls back to 1000ms (the default tick interval).
func (t *tickerTickProvider) PreferredTickInterval() time.Duration {
	policy := PolicySnapshot()
	if policy.AutoScrollMS > 0 {
		return time.Duration(policy.AutoScrollMS) * time.Millisecond
	}
	return 1000 * time.Millisecond
}
