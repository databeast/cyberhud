package attract_hacking

import (
	"sync"

	"github.com/databeast/cyberhud/display/modes/attract_hacking/source"
)

var policyState = struct {
	sync.RWMutex
	policy source.Policy
}{
	policy: source.DefaultPolicy(),
}

func GetPolicy() source.Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

func SetPolicy(p source.Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

func normalizePolicy(p source.Policy) source.Policy {
	if p.Speed < 0.1 {
		p.Speed = 0.1
	}
	if p.Speed > 3.0 {
		p.Speed = 3.0
	}
	if p.Density < 0.1 {
		p.Density = 0.1
	}
	if p.Density > 1.0 {
		p.Density = 1.0
	}
	if p.Glitch < 0.0 {
		p.Glitch = 0.0
	}
	if p.Glitch > 1.0 {
		p.Glitch = 1.0
	}
	if p.Intensity < 0.1 {
		p.Intensity = 0.1
	}
	if p.Intensity > 1.0 {
		p.Intensity = 1.0
	}
	if p.Pulse < 0.1 {
		p.Pulse = 0.1
	}
	if p.Pulse > 1.5 {
		p.Pulse = 1.5
	}
	return p
}
