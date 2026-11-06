package source

import (
	"math"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

// TyperState tracks a character cursor independently of the pager's scroll and page strategies.
type TyperState struct {
	nextLine uint64
	bytePos  int
	credit   int64
	rows     []string
	revision atomic.Uint64
}

func NewTyperState() *TyperState {
	return &TyperState{}
}

// Advance reveals complete UTF-8 characters from the buffer at the given rate.
func (ts *TyperState) Advance(delta time.Duration, speed int, lines []string, lineBreaks []bool, partial string, firstLine uint64, rows, columns int) {
	if rows <= 0 || columns <= 0 || speed <= 0 {
		return
	}
	if ts.nextLine < firstLine {
		if ts.bytePos > 0 {
			ts.newRow(rows)
		}
		ts.nextLine = firstLine
		ts.bytePos = 0
	}
	if delta > 0 {
		if int64(delta) > (math.MaxInt64-ts.credit)/int64(speed) {
			ts.credit = math.MaxInt64
		} else {
			ts.credit += int64(delta) * int64(speed)
		}
	}

	for {
		offset := ts.nextLine - firstLine
		if offset > uint64(len(lines)) {
			ts.credit = 0
			break
		}
		complete := offset < uint64(len(lines))
		line := partial
		if complete {
			line = lines[int(offset)]
		}
		if ts.bytePos >= len(line) {
			if !complete {
				ts.credit = 0
				break
			}
			if lineBreaks[int(offset)] {
				ts.newRow(rows)
			}
			ts.nextLine++
			ts.bytePos = 0
			continue
		}
		remaining := line[ts.bytePos:]
		if !complete && !utf8.FullRuneInString(remaining) {
			ts.credit = 0
			break
		}
		if ts.credit < int64(time.Second) {
			break
		}
		r, size := utf8.DecodeRuneInString(remaining)
		if len(ts.rows) == 0 {
			ts.rows = append(ts.rows, "")
		}
		if utf8.RuneCountInString(ts.rows[len(ts.rows)-1]) >= columns {
			ts.newRow(rows)
		}
		ts.rows[len(ts.rows)-1] += string(r)
		ts.bytePos += size
		ts.credit -= int64(time.Second)
		ts.revision.Add(1)
	}
}

func (ts *TyperState) newRow(limit int) {
	if len(ts.rows) == 0 {
		ts.rows = append(ts.rows, "")
	}
	ts.rows = append(ts.rows, "")
	if len(ts.rows) > limit {
		ts.rows = append([]string(nil), ts.rows[len(ts.rows)-limit:]...)
	}
	ts.revision.Add(1)
}

func (ts *TyperState) Frame(rows int) []string {
	if rows <= 0 || len(ts.rows) == 0 {
		return nil
	}
	start := len(ts.rows) - rows
	if start < 0 {
		start = 0
	}
	return append([]string(nil), ts.rows[start:]...)
}

func (ts *TyperState) Revision() uint64 {
	return ts.revision.Load()
}
