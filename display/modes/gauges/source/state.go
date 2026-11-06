package source

import (
	"sync"
)

var state = struct {
	sync.RWMutex
	policy  Policy
	snap    GaugeSet
	version uint64
}{
	policy: DefaultPolicy(),
	snap:   GaugeSet{},
}

func SetPolicy(p Policy) {
	state.Lock()
	defer state.Unlock()
	state.policy = NormalizePolicy(p)
	state.version++
}

func PolicySnapshot() Policy {
	state.RLock()
	defer state.RUnlock()
	return state.policy
}

func SetSnapshot(s GaugeSet) {
	state.Lock()
	defer state.Unlock()
	state.snap = CloneGaugeSet(s)
	state.version++
}

func Snapshot() GaugeSet {
	state.RLock()
	defer state.RUnlock()
	return CloneGaugeSet(state.snap)
}

func Version() uint64 {
	state.RLock()
	defer state.RUnlock()
	return state.version
}

func SetPayload(payload string) error {
	snap, err := ParsePayload(payload, PolicySnapshot())
	if err != nil {
		return err
	}
	SetSnapshot(snap)
	return nil
}

func SerializeSnapshotNow() (string, error) {
	return SerializeSnapshot(Snapshot())
}
