package content

import "sync"

// buffer is an internal ring buffer holding received messages.
type buffer struct {
	mu       sync.RWMutex
	lines    []string
	seq      uint64
	maxLines int
}

// Buffer exposes the internal ring buffer behavior to integration tests without exposing its fields.
type Buffer = buffer

func newBuffer(maxLines int) *buffer {
	if maxLines < 1 {
		maxLines = 1
	}
	return &buffer{maxLines: maxLines}
}

// NewBuffer creates a bounded message buffer.
func NewBuffer(maxLines int) *Buffer { return newBuffer(maxLines) }

// MaxLines returns the configured buffer capacity.
func (b *buffer) MaxLines() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.maxLines
}

func (b *buffer) Push(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = append(b.lines, msg)
	if len(b.lines) > b.maxLines {
		excess := len(b.lines) - b.maxLines
		b.lines = append([]string(nil), b.lines[excess:]...)
	}
	b.seq++
}

func (b *buffer) Snapshot() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if len(b.lines) == 0 {
		return nil
	}
	out := make([]string, len(b.lines))
	copy(out, b.lines)
	return out
}

func (b *buffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = nil
	b.seq++
}

func (b *buffer) SetMaxLines(maxLines int) {
	if maxLines < 1 {
		maxLines = 1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.maxLines = maxLines
	if len(b.lines) > b.maxLines {
		excess := len(b.lines) - b.maxLines
		b.lines = append([]string(nil), b.lines[excess:]...)
		b.seq++
	}
}

func (b *buffer) Seq() uint64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.seq
}

func (b *buffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.lines)
}
