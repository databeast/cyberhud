package attract_waveform

import (
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/modes/attract_waveform/source"
)

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

// GetPolicy returns the current waveform policy (thread-safe read under RWMutex).
func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the waveform policy after normalization (thread-safe write under Mutex).
func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// normalizePolicy ensures policy fields contain valid values.
func normalizePolicy(p source.Policy) source.Policy {
	if p.Speed < 0.1 {
		p.Speed = 0.1
	}
	if p.Speed > 10.0 {
		p.Speed = 10.0
	}
	if p.Density < 0.1 {
		p.Density = 0.1
	}
	if p.Density > 1.0 {
		p.Density = 1.0
	}
	if p.Amplitude < 0.1 {
		p.Amplitude = 0.1
	}
	if p.Amplitude > 1.0 {
		p.Amplitude = 1.0
	}
	if p.Traces < 1 {
		p.Traces = 1
	}
	if p.Traces > 8 {
		p.Traces = 8
	}
	if p.Persistence < 0.1 {
		p.Persistence = 0.1
	}
	if p.Persistence > 1.0 {
		p.Persistence = 1.0
	}

	switch strings.ToLower(strings.TrimSpace(p.Direction)) {
	case "horizontal", "h":
		p.Direction = "horizontal"
	case "vertical", "v":
		p.Direction = "vertical"
	case "auto", "":
		p.Direction = "auto"
	default:
		p.Direction = "auto"
	}

	return p
}
