package source

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func advanceTyper(ts *TyperState, b *LineBuffer, elapsed time.Duration, speed, rows, columns int) {
	lines, breaks, partial, first := b.SnapshotForTyper()
	ts.Advance(elapsed, speed, lines, breaks, partial, first, rows, columns)
}

func wantTyperFrame(t *testing.T, ts *TyperState, rows int, want ...string) {
	t.Helper()
	if got := ts.Frame(rows); !reflect.DeepEqual(got, want) {
		t.Fatalf("frame = %q, want %q", got, want)
	}
}

func TestTyperRatePartialAndUTF8(t *testing.T) {
	b := NewLineBuffer(5)
	ts := NewTyperState()
	b.Ingest([]byte{'A', 0xe2, 0x82})
	for n := 0; n < 9; n++ {
		advanceTyper(ts, b, 100*time.Millisecond, 1, 3, 10)
	}
	wantTyperFrame(t, ts, 3)
	advanceTyper(ts, b, 100*time.Millisecond, 1, 3, 10)
	wantTyperFrame(t, ts, 3, "A")
	advanceTyper(ts, b, 10*time.Second, 1, 3, 10)
	wantTyperFrame(t, ts, 3, "A")
	b.Ingest([]byte{0xac, 'B'})
	advanceTyper(ts, b, 500*time.Millisecond, 1, 3, 10)
	wantTyperFrame(t, ts, 3, "A")
	advanceTyper(ts, b, 500*time.Millisecond, 1, 3, 10)
	wantTyperFrame(t, ts, 3, "A€")
	advanceTyper(ts, b, time.Second, 1, 3, 10)
	wantTyperFrame(t, ts, 3, "A€B")
}

func TestTyperWrapsAndScrolls(t *testing.T) {
	b := NewLineBuffer(5)
	ts := NewTyperState()
	b.Ingest([]byte("abcd\nef"))
	advanceTyper(ts, b, 4*time.Second, 1, 2, 3)
	wantTyperFrame(t, ts, 2, "d", "")
	advanceTyper(ts, b, 2*time.Second, 1, 2, 3)
	wantTyperFrame(t, ts, 2, "d", "ef")
	b.Ingest([]byte("\nnext"))
	advanceTyper(ts, b, 4*time.Second, 1, 2, 3)
	wantTyperFrame(t, ts, 2, "nex", "t")
}

func TestTyperResumesAfterEviction(t *testing.T) {
	b := NewLineBuffer(2)
	ts := NewTyperState()
	b.Ingest([]byte("abc\n"))
	advanceTyper(ts, b, time.Second, 1, 4, 20)
	wantTyperFrame(t, ts, 4, "a")
	b.Ingest([]byte("two\nthree\nfour"))
	advanceTyper(ts, b, 12*time.Second, 1, 4, 20)
	wantTyperFrame(t, ts, 4, "a", "two", "three", "four")
}

func TestTyperLineCompletionDoesNotRepeatPartial(t *testing.T) {
	b := NewLineBuffer(2)
	ts := NewTyperState()
	b.Ingest([]byte("abc"))
	advanceTyper(ts, b, 3*time.Second, 1, 3, 10)
	b.Ingest([]byte("\nxy"))
	advanceTyper(ts, b, 2*time.Second, 1, 3, 10)
	wantTyperFrame(t, ts, 3, "abc", "xy")
}

func TestTyperHandlesLongLineAndUTF8AtBufferBoundary(t *testing.T) {
	b := NewLineBuffer(3)
	ts := NewTyperState()
	text := strings.Repeat("x", MaxLineBytes-1) + "€z"
	b.Ingest([]byte(text))
	advanceTyper(ts, b, 5*time.Second, 1000, 3, 5000)
	wantTyperFrame(t, ts, 3, text)
	b.Ingest([]byte("\nnext"))
	advanceTyper(ts, b, time.Second, 1000, 3, 5000)
	wantTyperFrame(t, ts, 3, text, "next")
}
