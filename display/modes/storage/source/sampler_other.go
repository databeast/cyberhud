//go:build !linux

package source

import "time"

// otherSampler is the non-Linux stub. statfs/proc-mounts are unavailable, so it
// reports no filesystems and the mode renders a "no storage data" placeholder.
type otherSampler struct{}

func newSampler() Sampler { return otherSampler{} }

func (otherSampler) Sample(showAll bool) (StorageSnapshot, error) {
	_ = showAll
	return StorageSnapshot{Timestamp: time.Now()}, nil
}
