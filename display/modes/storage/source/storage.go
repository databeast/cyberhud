package source

import (
	"log"
	"sort"
	"sync"
	"time"
)

// Sampler abstracts filesystem-usage collection, enabling mock implementations
// for testing without touching the real system.
type Sampler interface {
	// Sample returns the current usage of all monitored filesystems. The showAll
	// flag requests inclusion of pseudo/virtual filesystems (tmpfs, overlay, ...).
	Sample(showAll bool) (StorageSnapshot, error)
}

// NewSampler returns the platform-appropriate Sampler. On Linux it reads
// /proc/mounts and statfs(2); on other platforms it returns an empty snapshot.
func NewSampler() Sampler { return newSampler() }

// sortMounts orders mounts most-used-first (descending Percent), breaking ties by
// mountpoint so ordering is deterministic. The root filesystem is always promoted
// to the front since it is the one operators care about most.
func sortMounts(mounts []MountUsage) {
	sort.SliceStable(mounts, func(i, j int) bool {
		if mounts[i].Mountpoint == "/" {
			return true
		}
		if mounts[j].Mountpoint == "/" {
			return false
		}
		pi, pj := mounts[i].Percent(), mounts[j].Percent()
		if pi != pj {
			return pi > pj
		}
		return mounts[i].Mountpoint < mounts[j].Mountpoint
	})
}

// ── Background sampling loop ────────────────────────────────────────────────

// loopState holds the lifecycle state for the background sampling goroutine.
var loopState struct {
	sync.Mutex
	stopCh  chan struct{}
	sampler Sampler
}

// Activate starts the background sampling loop using the default platform sampler.
// It is idempotent: calling Activate when the loop is already running is a no-op.
// refreshMS and showAll are re-read each tick so policy changes apply live.
func Activate(refreshMS func() int, showAll func() bool) {
	ActivateWith(NewSampler(), refreshMS, showAll)
}

// ActivateWith starts the background sampling loop using the provided Sampler.
// It is idempotent.
func ActivateWith(sampler Sampler, refreshMS func() int, showAll func() bool) {
	loopState.Lock()
	defer loopState.Unlock()

	if loopState.stopCh != nil {
		return // already running
	}

	// Sample once synchronously so data is available on the activating tick.
	doSample(sampler, readShowAll(showAll))

	loopState.sampler = sampler
	loopState.stopCh = make(chan struct{})
	go samplingLoop(sampler, loopState.stopCh, refreshMS, showAll)
}

// Deactivate stops the background sampling loop. It is idempotent.
func Deactivate() {
	loopState.Lock()
	ch := loopState.stopCh
	loopState.stopCh = nil
	loopState.sampler = nil
	loopState.Unlock()

	if ch != nil {
		close(ch)
	}
}

// SampleActive samples synchronously using the active sampler, updating the
// stored snapshot. Returns (snapshot, true) on success. Used to fill the first
// render tick when the loop has not yet produced data.
func SampleActive(showAll bool) (StorageSnapshot, bool) {
	loopState.Lock()
	sampler := loopState.sampler
	loopState.Unlock()
	if sampler == nil {
		return StorageSnapshot{}, false
	}
	snap, err := sampler.Sample(showAll)
	if err != nil || len(snap.Mounts) == 0 {
		return StorageSnapshot{}, false
	}
	UpdateSnapshot(snap)
	return snap, true
}

func samplingLoop(sampler Sampler, stopCh <-chan struct{}, refreshMSFunc func() int, showAllFunc func() bool) {
	refreshMS := ClampRefreshMS(readRefresh(refreshMSFunc))
	ticker := time.NewTicker(time.Duration(refreshMS) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			doSample(sampler, readShowAll(showAllFunc))

			newRefreshMS := ClampRefreshMS(readRefresh(refreshMSFunc))
			if newRefreshMS != refreshMS {
				refreshMS = newRefreshMS
				ticker.Reset(time.Duration(refreshMS) * time.Millisecond)
			}
		}
	}
}

func readRefresh(f func() int) int {
	if f == nil {
		return DefaultRefreshMS
	}
	return f()
}

func readShowAll(f func() bool) bool {
	if f == nil {
		return false
	}
	return f()
}

func doSample(sampler Sampler, showAll bool) {
	snap, err := sampler.Sample(showAll)
	if err != nil {
		log.Printf("storage: Sample() error: %v", err)
		return
	}
	UpdateSnapshot(snap)
}

// ── Snapshot state ──────────────────────────────────────────────────────────

var snapshotState = struct {
	sync.RWMutex
	snap StorageSnapshot
}{}

// CurrentSnapshot returns the most recent storage snapshot (thread-safe).
func CurrentSnapshot() StorageSnapshot {
	snapshotState.RLock()
	defer snapshotState.RUnlock()
	return snapshotState.snap
}

// UpdateSnapshot stores a new storage snapshot (thread-safe).
func UpdateSnapshot(snap StorageSnapshot) {
	snapshotState.Lock()
	defer snapshotState.Unlock()
	snapshotState.snap = snap
}

// ResetSnapshotStateForTest clears the stored snapshot.
func ResetSnapshotStateForTest() { UpdateSnapshot(StorageSnapshot{}) }
