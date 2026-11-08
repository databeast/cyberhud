package source

import "time"

// Event is a single calendar entry reduced to the fields the display cares about.
type Event struct {
	Summary  string
	Location string
	Start    time.Time
	End      time.Time
	AllDay   bool
}

// Active reports whether the event spans the instant t (start inclusive, end
// exclusive). An event with a zero End is treated as a point in time and is
// never "active" beyond its start minute.
func (e Event) Active(t time.Time) bool {
	if e.Start.IsZero() {
		return false
	}
	if t.Before(e.Start) {
		return false
	}
	if e.End.IsZero() {
		return false
	}
	return t.Before(e.End)
}

// Snapshot is a read-only view of the live schedule state. It carries the full
// parsed, time-sorted event set so the render path can pick the event active for
// the current instant on every tick, independent of the fetch cadence.
type Snapshot struct {
	Sequence   uint64    // Bumped on every state mutation; drives render change detection.
	Loaded     bool      // True once a calendar has been successfully parsed at least once.
	Source     string    // The source the snapshot was built from.
	LastError  string    // Most recent fetch/parse error, empty when healthy.
	FetchedAt  time.Time // When the calendar was last successfully fetched.
	Events     []Event   // All future-or-current events, sorted by start ascending.
	CalendarName string  // X-WR-CALNAME / NAME if the calendar declared one.

	// Now is the instant the view is being rendered for, stamped by the render
	// path. The active-event decision moves with the clock even when the fetched
	// event set does not, so a style resolves "what is on now" against this field
	// rather than reading a clock itself.
	Now time.Time
}

// ActiveEvent returns the event active at instant t, or (zero, false) when none
// is. When events overlap, the one that started most recently wins, which is the
// intuitive "what am I in right now" answer for back-to-back or nested entries.
func (s Snapshot) ActiveEvent(t time.Time) (Event, bool) {
	var best Event
	found := false
	for _, e := range s.Events {
		if !e.Active(t) {
			continue
		}
		if !found || e.Start.After(best.Start) {
			best = e
			found = true
		}
	}
	return best, found
}

// NextEvent returns the earliest event that starts strictly after instant t, or
// (zero, false) when none remain. Events are start-sorted, so the first match is
// the soonest.
func (s Snapshot) NextEvent(t time.Time) (Event, bool) {
	for _, e := range s.Events {
		if e.Start.After(t) {
			return e, true
		}
	}
	return Event{}, false
}
