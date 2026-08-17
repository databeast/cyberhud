package region

import (
	"sync"
	"time"
)

// TickRateResolver is the RenderLoop-layer interface that determines the render
// tick interval for a Region based on its current display mode and mode policy.
// Implementations are typically backed by a registry of [TickRateProvider] values.
type TickRateResolver interface {
	// TickInterval returns the appropriate tick duration for the mode identified
	// by modeID. The returned duration controls how often the Region is rendered
	// when per-region deadline scheduling is active.
	TickInterval(modeID string) time.Duration
}

// TickRateProvider is the RenderLoop-layer interface implemented by display modes
// that need a custom tick interval. Modes register themselves via [RegisterTickRate]
// during package init(); the [DefaultTickRateResolver] queries the registry at
// runtime instead of maintaining a hardcoded switch.
type TickRateProvider interface {
	// PreferredTickInterval returns the mode's desired tick duration.
	// The resolver clamps the result to [1ms, 10000ms] before use.
	PreferredTickInterval() time.Duration
}

const (
	// defaultTickInterval is the fallback interval for modes without timing policies.
	defaultTickInterval = 1000 * time.Millisecond

	// minTickInterval is the minimum allowed tick interval.
	minTickInterval = 1 * time.Millisecond

	// maxTickInterval is the maximum allowed tick interval.
	maxTickInterval = 10000 * time.Millisecond
)

var (
	providersMu sync.RWMutex
	providers   = map[string]TickRateProvider{}
)

// RegisterTickRate registers a [TickRateProvider] for the given mode ID in the
// global tick rate registry. Called from mode package init() functions so that
// [DefaultTickRateResolver] can look up preferred intervals at runtime.
//
// The modeID parameter identifies the display mode (e.g. "attract_bokeh").
// The p parameter is the provider whose PreferredTickInterval will be queried.
// If modeID is empty or p is nil, the call is silently ignored.
func RegisterTickRate(modeID string, p TickRateProvider) {
	if modeID == "" || p == nil {
		return
	}
	providersMu.Lock()
	defer providersMu.Unlock()
	providers[modeID] = p
}

// DefaultTickRateResolver is the RenderLoop-layer default implementation of
// [TickRateResolver]. It queries the global tick rate provider registry to derive
// tick intervals. For registered modes, it returns the provider's preferred
// interval clamped to [1ms, 10000ms]. For all other modes, it returns the
// default 1000ms interval.
type DefaultTickRateResolver struct{}

// TickInterval returns the appropriate tick duration for the display mode
// identified by modeID. If a [TickRateProvider] is registered for modeID and
// returns a value greater than zero, that value is clamped to [1ms, 10000ms]
// and returned. Otherwise the default interval of 1000ms is returned.
func (d *DefaultTickRateResolver) TickInterval(modeID string) time.Duration {
	providersMu.RLock()
	p, ok := providers[modeID]
	providersMu.RUnlock()

	if ok {
		return clampInterval(p.PreferredTickInterval())
	}
	return defaultTickInterval
}

// clampInterval restricts a duration to the [minTickInterval, maxTickInterval] range.
func clampInterval(d time.Duration) time.Duration {
	if d < minTickInterval {
		return minTickInterval
	}
	if d > maxTickInterval {
		return maxTickInterval
	}
	return d
}
