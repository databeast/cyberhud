package source

import (
	"bytes"
	"sync"
)

// maxLineBytes is the maximum number of bytes allowed in a single line
// without a newline terminator. Lines exceeding this limit are truncated
// at this boundary and the truncation point is treated as a line boundary.
const maxLineBytes = 4096

const MaxLineBytes = maxLineBytes

// lineBuffer is a thread-safe ring buffer that stores ingested text lines.
// It splits incoming byte data on newline boundaries, retains trailing partial
// bytes, and discards the oldest lines when the configured capacity is exceeded.
type LineBuffer struct {
	mu       sync.RWMutex
	lines    []string
	partial  []byte // incomplete line (no trailing newline yet)
	seq      uint64 // monotonic mutation counter
	maxLines int    // configurable capacity [1..1000]
}

// newLineBuffer creates a lineBuffer with the given capacity.
// The capacity is clamped to [1, 1000].
func NewLineBuffer(maxLines int) *LineBuffer {
	maxLines = clampMaxLines(maxLines)
	return &LineBuffer{
		lines:    make([]string, 0, maxLines),
		maxLines: maxLines,
	}
}

// Ingest processes raw bytes from the data source. It prepends any existing
// partial bytes, splits on '\n' boundaries, stores complete lines in the ring
// buffer, and retains trailing bytes without '\n' as partial. If a contiguous
// run exceeds maxLineBytes without a newline, it is truncated at that limit
// and treated as a line boundary.
func (b *LineBuffer) Ingest(data []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Prepend any leftover partial bytes from the previous call.
	var buf []byte
	if len(b.partial) > 0 {
		buf = make([]byte, 0, len(b.partial)+len(data))
		buf = append(buf, b.partial...)
		buf = append(buf, data...)
		b.partial = nil
	} else {
		buf = data
	}

	var newLines []string

	for len(buf) > 0 {
		idx := bytes.IndexByte(buf, '\n')
		if idx >= 0 {
			// Found a newline — everything up to it is a complete line.
			newLines = append(newLines, string(buf[:idx]))
			buf = buf[idx+1:]
		} else if len(buf) > maxLineBytes {
			// No newline found and exceeds max line length — truncate.
			newLines = append(newLines, string(buf[:maxLineBytes]))
			buf = buf[maxLineBytes:]
		} else {
			// No newline and within limit — retain as partial.
			b.partial = make([]byte, len(buf))
			copy(b.partial, buf)
			buf = nil
		}
	}

	// Append complete lines to the ring buffer.
	if len(newLines) > 0 {
		b.lines = append(b.lines, newLines...)
		// Discard oldest lines if exceeding capacity.
		if len(b.lines) > b.maxLines {
			discard := len(b.lines) - b.maxLines
			// Allow GC of discarded strings.
			for i := 0; i < discard; i++ {
				b.lines[i] = ""
			}
			b.lines = b.lines[discard:]
		}
	}

	// Increment sequence counter if we produced new lines or modified partial.
	b.seq++
}

// Snapshot returns a copy of the current buffered lines. The returned slice
// is safe for the caller to use without holding any lock.
func (b *LineBuffer) Snapshot() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	out := make([]string, len(b.lines))
	copy(out, b.lines)
	return out
}

// Len returns the current number of complete lines in the buffer.
func (b *LineBuffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.lines)
}

// Seq returns the monotonic mutation counter. It is incremented on each
// Ingest call.
func (b *LineBuffer) Seq() uint64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.seq
}

// SetMaxLines updates the buffer capacity. The value is clamped to [1, 1000].
// If the new capacity is smaller than the current line count, the oldest lines
// are discarded.
func (b *LineBuffer) SetMaxLines(n int) {
	n = clampMaxLines(n)

	b.mu.Lock()
	defer b.mu.Unlock()

	b.maxLines = n
	if len(b.lines) > n {
		discard := len(b.lines) - n
		for i := 0; i < discard; i++ {
			b.lines[i] = ""
		}
		b.lines = b.lines[discard:]
	}
}

// clampMaxLines enforces the valid range [1, 1000] for buffer capacity.
func clampMaxLines(n int) int {
	if n < 1 {
		return 1
	}
	if n > 1000 {
		return 1000
	}
	return n
}
