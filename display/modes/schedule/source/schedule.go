package source

import (
	"log"
	"sync"
	"time"
)

// state holds the live schedule state behind a mutex. It mirrors the serial
// mode's shape: a policy, the latest parsed data, a sequence counter for change
// detection, and a stop channel for the background fetch goroutine.
var state = struct {
	sync.RWMutex
	policy    Policy
	loaded    bool
	source    string
	lastError string
	fetchedAt time.Time
	events    []Event
	calName   string
	seq       uint64
	stopCh    chan struct{}
	now       func() time.Time
}{
	policy: DefaultPolicy(),
	now:    time.Now,
}

// PolicySnapshot returns the current schedule policy.
func PolicySnapshot() Policy {
	state.RLock()
	defer state.RUnlock()
	return state.policy
}

// SetPolicy updates the schedule policy. When the source changes, the cached
// calendar is dropped so the next fetch repopulates from the new location.
func SetPolicy(policy Policy) {
	state.Lock()
	policy = NormalizePolicy(policy)
	prev := state.policy
	state.policy = policy
	if prev.Source != policy.Source {
		state.loaded = false
		state.events = nil
		state.lastError = ""
		state.calName = ""
	}
	state.seq++
	running := state.stopCh != nil
	state.Unlock()

	// A source change should take effect promptly rather than waiting out the
	// current refresh interval, so kick an immediate fetch when active.
	if running && prev.Source != policy.Source {
		go fetchOnce(policy)
	}
}

// SnapshotNow returns the current schedule snapshot.
func SnapshotNow() Snapshot {
	state.RLock()
	defer state.RUnlock()
	return snapshotLocked()
}

func snapshotLocked() Snapshot {
	events := append([]Event(nil), state.events...)
	return Snapshot{
		Sequence:     state.seq,
		Loaded:       state.loaded,
		Source:       state.source,
		LastError:    state.lastError,
		FetchedAt:    state.fetchedAt,
		Events:       events,
		CalendarName: state.calName,
	}
}

// Activate starts the background calendar fetcher. Safe to call repeatedly;
// only one goroutine runs at a time.
func Activate() {
	state.Lock()
	if state.stopCh != nil {
		state.Unlock()
		return
	}
	state.stopCh = make(chan struct{})
	stop := state.stopCh
	state.Unlock()
	go fetchLoop(stop)
}

// Deactivate stops the background calendar fetcher.
func Deactivate() {
	state.Lock()
	ch := state.stopCh
	state.stopCh = nil
	state.Unlock()
	if ch != nil {
		close(ch)
	}
}

func fetchLoop(stop <-chan struct{}) {
	// Fetch immediately on activation, then on the policy's refresh interval.
	fetchOnce(PolicySnapshot())
	for {
		p := PolicySnapshot()
		delay := time.Duration(p.RefreshSec) * time.Second
		if delay <= 0 {
			delay = DefaultRefreshSec * time.Second
		}
		select {
		case <-stop:
			return
		case <-time.After(delay):
			fetchOnce(PolicySnapshot())
		}
	}
}

// fetchOnce loads the calendar named by the policy and stores the result (or the
// error) into state. An empty source clears the buffer and records a waiting
// message rather than an error.
func fetchOnce(p Policy) {
	if p.Source == "" {
		state.Lock()
		state.source = ""
		state.loaded = false
		state.events = nil
		state.calName = ""
		state.lastError = ""
		state.seq++
		state.Unlock()
		return
	}

	now := clock()
	cal, err := loadCalendar(p.Source)
	if err != nil {
		state.Lock()
		state.source = p.Source
		state.lastError = err.Error()
		state.seq++
		state.Unlock()
		log.Printf("[cyberhudd] schedule mode: %v", err)
		return
	}

	events := extractEvents(cal, now)
	name := calendarName(cal)

	state.Lock()
	state.source = p.Source
	state.loaded = true
	state.lastError = ""
	state.fetchedAt = now
	state.events = events
	state.calName = name
	state.seq++
	state.Unlock()
}

func clock() time.Time {
	state.RLock()
	fn := state.now
	state.RUnlock()
	if fn == nil {
		return time.Now()
	}
	return fn()
}

// SetSnapshotForTest installs a fixed snapshot and policy for tests.
func SetSnapshotForTest(s Snapshot, p Policy) {
	state.Lock()
	defer state.Unlock()
	state.loaded = s.Loaded
	state.source = s.Source
	state.lastError = s.LastError
	state.fetchedAt = s.FetchedAt
	state.events = append([]Event(nil), s.Events...)
	state.calName = s.CalendarName
	state.seq = s.Sequence
	state.policy = NormalizePolicy(p)
}

// ResetForTest returns state to defaults.
func ResetForTest() {
	state.Lock()
	defer state.Unlock()
	state.loaded = false
	state.source = ""
	state.lastError = ""
	state.fetchedAt = time.Time{}
	state.events = nil
	state.calName = ""
	state.seq = 0
	state.policy = DefaultPolicy()
}
