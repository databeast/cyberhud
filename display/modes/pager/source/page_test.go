package source

import (
	"reflect"
	"testing"
	"time"
)

func tickPage(ps *PageState, buf *LineBuffer, delta time.Duration, rows, maxWait int) {
	lines, first := buf.SnapshotWithPosition()
	ps.Tick(delta, lines, first, rows, 0, 0, maxWait)
}

func finishPageTransition(t *testing.T, ps *PageState, buf *LineBuffer, rows int) {
	t.Helper()
	if ps.Phase() != PhaseFadeOut {
		t.Fatalf("phase = %v, want fade out", ps.Phase())
	}
	tickPage(ps, buf, time.Millisecond, rows, 5)
	tickPage(ps, buf, time.Millisecond, rows, 5)
	if ps.Phase() != PhaseIdle {
		t.Fatalf("phase = %v, want idle", ps.Phase())
	}
}

func expectPage(t *testing.T, ps *PageState, want ...string) {
	t.Helper()
	if got := ps.CurrentPage(); !reflect.DeepEqual(got, want) {
		t.Fatalf("page = %q, want %q", got, want)
	}
}

func TestPageAdvancesOnlyForUnreadLines(t *testing.T) {
	buf := NewLineBuffer(10)
	ps := NewPageState(2, 1)
	buf.Ingest([]byte("one\ntwo\nthree\nfour\n"))

	tickPage(ps, buf, 3*time.Second, 2, 5)
	finishPageTransition(t, ps, buf, 2)
	expectPage(t, ps, "one", "two")

	tickPage(ps, buf, 3*time.Second, 2, 5)
	finishPageTransition(t, ps, buf, 2)
	expectPage(t, ps, "three", "four")

	tickPage(ps, buf, 30*time.Second, 2, 5)
	if ps.Phase() != PhaseIdle {
		t.Fatalf("phase = %v with no unread lines, want idle", ps.Phase())
	}
	expectPage(t, ps, "three", "four")

	// Identical text at a new position is still new input.
	buf.Ingest([]byte("three\nfour\n"))
	tickPage(ps, buf, 3*time.Second, 2, 5)
	finishPageTransition(t, ps, buf, 2)
	expectPage(t, ps, "three", "four")
}

func TestPageWaitsForPartialUnreadContent(t *testing.T) {
	buf := NewLineBuffer(10)
	ps := NewPageState(2, 1)
	buf.Ingest([]byte("one\ntwo\n"))
	tickPage(ps, buf, 3*time.Second, 2, 5)
	finishPageTransition(t, ps, buf, 2)

	// Idle time must not count towards the next partial page's timeout.
	tickPage(ps, buf, 30*time.Second, 2, 5)
	buf.Ingest([]byte("three\n"))
	tickPage(ps, buf, 4*time.Second, 2, 5)
	if ps.Phase() != PhaseIdle {
		t.Fatalf("phase = %v before partial timeout, want idle", ps.Phase())
	}
	expectPage(t, ps, "one", "two")

	tickPage(ps, buf, time.Second, 2, 5)
	finishPageTransition(t, ps, buf, 2)
	expectPage(t, ps, "three", "")
	tickPage(ps, buf, 10*time.Second, 2, 5)
	if ps.Phase() != PhaseIdle {
		t.Fatalf("phase = %v after partial page, want idle", ps.Phase())
	}
}

func TestPageSkipsEvictedUnreadLines(t *testing.T) {
	buf := NewLineBuffer(3)
	ps := NewPageState(2, 1)
	buf.Ingest([]byte("one\ntwo\n"))
	tickPage(ps, buf, 3*time.Second, 2, 5)

	// The selected page remains stable even if the buffer changes during a fade.
	buf.Ingest([]byte("three\nfour\nfive\nsix\n"))
	finishPageTransition(t, ps, buf, 2)
	expectPage(t, ps, "one", "two")

	tickPage(ps, buf, 3*time.Second, 2, 5)
	finishPageTransition(t, ps, buf, 2)
	expectPage(t, ps, "four", "five")
	tickPage(ps, buf, 5*time.Second, 2, 5)
	finishPageTransition(t, ps, buf, 2)
	expectPage(t, ps, "six", "")
}

func TestBufferSnapshotPositionTracksTrimmingAndPartialLines(t *testing.T) {
	buf := NewLineBuffer(3)
	buf.Ingest([]byte("one\ntwo\nthree\npartial"))
	lines, first := buf.SnapshotWithPosition()
	if first != 0 || !reflect.DeepEqual(lines, []string{"one", "two", "three"}) {
		t.Fatalf("snapshot = %q at %d", lines, first)
	}

	buf.Ingest([]byte("\nfive\n"))
	lines, first = buf.SnapshotWithPosition()
	if first != 2 || !reflect.DeepEqual(lines, []string{"three", "partial", "five"}) {
		t.Fatalf("snapshot = %q at %d", lines, first)
	}
	buf.SetMaxLines(1)
	lines, first = buf.SnapshotWithPosition()
	if first != 4 || !reflect.DeepEqual(lines, []string{"five"}) {
		t.Fatalf("trimmed snapshot = %q at %d", lines, first)
	}
}
