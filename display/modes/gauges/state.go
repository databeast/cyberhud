package gauges

import (
	"sync"

	"github.com/databeast/cyberhud/display/modes/gauges/source"
)

var state = struct {
	sync.RWMutex
	policy  source.Policy
	snap    source.GaugeSet
	version uint64
}{
	policy: source.DefaultPolicy(),
	snap:   source.GaugeSet{},
}

func SetPolicy(p source.Policy) {
	state.Lock()
	defer state.Unlock()
	state.policy = source.NormalizePolicy(p)
	state.version++
}

func sourcePolicySnapshot() source.Policy {
	state.RLock()
	defer state.RUnlock()
	return state.policy
}

func sourceSnapshot() source.GaugeSet {
	state.RLock()
	defer state.RUnlock()
	return source.CloneGaugeSet(state.snap)
}

func sourceVersion() uint64 {
	state.RLock()
	defer state.RUnlock()
	return state.version
}

func SetSnapshot(s source.GaugeSet) {
	state.Lock()
	defer state.Unlock()
	state.snap = source.CloneGaugeSet(s)
	state.version++
}
