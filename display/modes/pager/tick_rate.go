package pager

import "time"

// tickProvider implements region.TickRateProvider for the pager mode.
// It returns a fixed 33ms interval (~30fps) for smooth scroll animation.
// Registration into the instance lifecycle happens during Activate() (task 6.1).
type tickProvider struct{}

// PreferredTickInterval returns 33ms (~30fps).
func (tickProvider) PreferredTickInterval() time.Duration {
	return 33 * time.Millisecond
}
