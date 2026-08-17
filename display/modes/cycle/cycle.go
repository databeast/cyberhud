// Package cycle implements the "cycle" display mode — a configurable meta-mode
// that auto-rotates through a subset of display modes on one or more regions.
package cycle

import (
	"sync"
	"time"

	"github.com/databeast/cyberhud/display/coordinator"
)

const (
	// ModeID is the catalog identifier for the cycle mode.
	ModeID = "cycle"

	// DefaultInterval is the time between automatic mode switches.
	DefaultInterval = 30 * time.Second

	// MinInterval is the minimum allowed cycling interval to prevent display thrashing.
	MinInterval = 5 * time.Second

	// MaxInterval is the maximum allowed cycling interval.
	MaxInterval = 10 * time.Minute

	// maxRegions is the maximum number of region entries allowed in policy.
	maxRegions = 16
)

// Policy holds configurable parameters for the cycle mode.
type Policy struct {
	Interval time.Duration // time between mode switches (clamped: 5s–10m, default 30s)
	Modes    []string      // ordered list of mode IDs to cycle (empty = all on region)
	Regions  []int         // region indices to cycle (empty = activating region only)
}

// normalizeInterval clamps d to [MinInterval, MaxInterval] and defaults to
// DefaultInterval when d is zero or negative.
func normalizeInterval(d time.Duration) time.Duration {
	if d <= 0 {
		return DefaultInterval
	}
	if d < MinInterval {
		return MinInterval
	}
	if d > MaxInterval {
		return MaxInterval
	}
	return d
}

// dedup removes duplicate strings from s, preserving first-occurrence order.
func dedup(s []string) []string {
	if len(s) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(s))
	out := make([]string, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// dedupInts removes duplicate integers from s, preserving first-occurrence order.
func dedupInts(s []int) []int {
	if len(s) == 0 {
		return nil
	}
	seen := make(map[int]struct{}, len(s))
	out := make([]int, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// policyStore holds the cycle mode policy state with read-write lock protection.
type policyStore struct {
	mu     sync.RWMutex
	policy Policy
}

// policyState holds the current cycle mode policy.
var policyState = &policyStore{
	policy: Policy{Interval: DefaultInterval},
}

// GetPolicy returns a copy of the current cycle mode policy.
func GetPolicy() Policy {
	policyState.mu.RLock()
	defer policyState.mu.RUnlock()
	return policyState.policy
}

// SetPolicy updates the cycle mode policy. It deduplicates Modes and Regions,
// and enforces the 16-entry cap on Regions (truncating to the first 16).
func SetPolicy(p Policy) {
	p.Modes = dedup(p.Modes)
	p.Regions = dedupInts(p.Regions)
	if len(p.Regions) > maxRegions {
		p.Regions = p.Regions[:maxRegions]
	}
	policyState.mu.Lock()
	defer policyState.mu.Unlock()
	policyState.policy = p
}

// globalModeStateMu protects access to globalModeState.
var globalModeStateMu sync.RWMutex

// globalModeState holds the coordinator state reference injected at startup.
var globalModeState *coordinator.State

// WireState injects the coordinator state reference. Called once at startup
// from cmd/cyberhudd.
func WireState(s *coordinator.State) {
	globalModeStateMu.Lock()
	defer globalModeStateMu.Unlock()
	globalModeState = s
}

// getGlobalModeState returns the current coordinator state reference.
func getGlobalModeState() *coordinator.State {
	globalModeStateMu.RLock()
	defer globalModeStateMu.RUnlock()
	return globalModeState
}
