package pager

import (
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/modes/pager/source"
)

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current pager policy (thread-safe read under RWMutex).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the pager policy after normalization (thread-safe write).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// policyFingerprint returns a deterministic string representation of all policy
// fields, used for change detection in RenderCacheKey.
func policyFingerprint(p Policy) string {
	return p.Fingerprint()
}

type Policy = source.Policy

func DefaultPolicy() source.Policy {
	return source.DefaultPolicy()
}

func normalizePolicy(p source.Policy) source.Policy {
	p = source.NormalizePolicy(p)
	p.Source = strings.TrimSpace(p.Source)
	p.Font = strings.TrimSpace(p.Font)
	p.Style = strings.ToLower(strings.TrimSpace(p.Style))
	if p.Style != "" && pagerRegistry.Lookup(p.Style) == nil {
		p.Style = ""
	}
	return p
}

// pagerSnapshotter implements catalog.PolicySnapshotter for the pager mode.
type pagerSnapshotter struct{}

// SnapshotPolicy returns the current pager policy as a JSON-serializable map
// using snake_case keys matching the protocol wire format.
func (pagerSnapshotter) SnapshotPolicy() map[string]interface{} {
	p := GetPolicy()
	return map[string]interface{}{
		"source":       p.Source,
		"scroll_speed": p.ScrollSpeed,
		"max_lines":    p.MaxLines,
		"scan_ms":      p.ScanMS,
		"font":         p.Font,
		"style":        p.Style,
		"fade_out_ms":  p.FadeOutMS,
		"fade_in_ms":   p.FadeInMS,
		"line_time_ms": p.LineTimeMS,
		"max_wait_s":   p.MaxWaitS,
	}
}

// RestorePolicy applies policy values from a JSON map, running them through
// the same normalization as SetPolicy. Keys use the protocol wire format.
func (pagerSnapshotter) RestorePolicy(data map[string]interface{}) error {
	p := source.DefaultPolicy()

	if v, ok := data["source"]; ok {
		if s, ok := v.(string); ok {
			p.Source = s
		}
	}
	if v, ok := data["scroll_speed"]; ok {
		if n, ok := toInt(v); ok {
			p.ScrollSpeed = n
		}
	}
	if v, ok := data["max_lines"]; ok {
		if n, ok := toInt(v); ok {
			p.MaxLines = n
		}
	}
	if v, ok := data["scan_ms"]; ok {
		if n, ok := toInt(v); ok {
			p.ScanMS = n
		}
	}
	if v, ok := data["font"]; ok {
		if s, ok := v.(string); ok {
			p.Font = s
		}
	}
	if v, ok := data["style"]; ok {
		if s, ok := v.(string); ok {
			p.Style = s
		}
	}
	if v, ok := data["fade_out_ms"]; ok {
		if n, ok := toInt(v); ok {
			p.FadeOutMS = n
		}
	}
	if v, ok := data["fade_in_ms"]; ok {
		if n, ok := toInt(v); ok {
			p.FadeInMS = n
		}
	}
	if v, ok := data["line_time_ms"]; ok {
		if n, ok := toInt(v); ok {
			p.LineTimeMS = n
		}
	}
	if v, ok := data["max_wait_s"]; ok {
		if n, ok := toInt(v); ok {
			p.MaxWaitS = n
		}
	}

	SetPolicy(p)
	return nil
}

// toInt extracts an int from an interface value, handling both float64
// (JSON numbers decode as float64) and direct int types.
func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}

func init() {
	catalog.RegisterSnapshotter("pager", pagerSnapshotter{})
}
