package cycle

import (
	"time"

	"github.com/databeast/cyberhud/display/catalog"
)

// cycleSnapshotter implements catalog.PolicySnapshotter for the cycle mode.
type cycleSnapshotter struct{}

// SnapshotPolicy returns the current cycle policy as a JSON-serializable map.
func (cycleSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()

	modes := make([]interface{}, len(p.Modes))
	for i, m := range p.Modes {
		modes[i] = m
	}

	regions := make([]interface{}, len(p.Regions))
	for i, r := range p.Regions {
		regions[i] = r
	}

	return map[string]interface{}{
		"interval": p.Interval.String(),
		"modes":    modes,
		"regions":  regions,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Always returns nil (graceful degradation).
func (cycleSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := Policy{Interval: DefaultInterval}

	if v, ok := data["interval"]; ok {
		if s, ok := v.(string); ok {
			if d, err := time.ParseDuration(s); err == nil {
				p.Interval = normalizeInterval(d)
			}
		}
	}

	if v, ok := data["modes"]; ok {
		if items, ok := v.([]interface{}); ok {
			for _, item := range items {
				if s, ok := item.(string); ok {
					if _, known := catalog.Describe(s); known {
						p.Modes = append(p.Modes, s)
					}
				}
			}
		}
	}

	if v, ok := data["regions"]; ok {
		if items, ok := v.([]interface{}); ok {
			for _, item := range items {
				if f, ok := item.(float64); ok && f >= 0 {
					p.Regions = append(p.Regions, int(f))
				}
			}
		}
	}

	SetPolicy(p)
	return nil
}
