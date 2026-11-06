package source

import (
	"sync"
)

// scrollState tracks the smooth-scroll position and velocity for fast displays.
// It is advanced each frame (33ms tick) and produces the visible text window.
// All fields are protected by mu since BuildView (render goroutine) and
// SetBaseSpeed (command handler goroutine) access state concurrently.
type ScrollState struct {
	mu        sync.Mutex
	offsetPx  int  // current pixel offset within the top visible line
	velocity  int  // current pixels-per-second (may be accelerated)
	baseSpeed int  // configured scroll_speed (pixels/sec)
	paused    bool // true when buffer is empty
}

// newScrollState creates a scrollState with the given base speed (pixels/sec).
func NewScrollState(baseSpeed int) *ScrollState {
	return &ScrollState{
		baseSpeed: baseSpeed,
		velocity:  baseSpeed,
	}
}

// Advance moves the scroll offset forward by the amount dictated by the current
// velocity over the given frame delta (in seconds). When the offset exceeds
// rowHeight, it wraps around — the caller must handle line consumption externally.
//
// frameDeltaSec is typically 0.033 (33ms at ~30fps).
func (s *ScrollState) Advance(frameDeltaSec float64, rowHeight int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.paused || s.velocity <= 0 || rowHeight <= 0 {
		return
	}
	pixels := int(float64(s.velocity) * frameDeltaSec)
	if pixels < 1 {
		pixels = 1
	}
	s.offsetPx += pixels
}

// AdaptVelocity adjusts scroll velocity based on the relationship between
// buffered lines and visible rows:
//   - bufferedLines > visibleRows: velocity = baseSpeed × min(4.0, bufferedLines/visibleRows)
//   - bufferedLines == 0: paused, velocity = 0
//   - bufferedLines ≤ visibleRows (and > 0): velocity = baseSpeed (restore normal speed)
func (s *ScrollState) AdaptVelocity(bufferedLines, visibleRows int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if bufferedLines == 0 {
		s.paused = true
		s.velocity = 0
		return
	}

	s.paused = false

	if visibleRows <= 0 {
		s.velocity = s.baseSpeed
		return
	}

	if bufferedLines > visibleRows {
		multiplier := float64(bufferedLines) / float64(visibleRows)
		if multiplier > 4.0 {
			multiplier = 4.0
		}
		s.velocity = int(float64(s.baseSpeed) * multiplier)
	} else {
		s.velocity = s.baseSpeed
	}
}

// SetBaseSpeed updates the base scroll speed without altering the current
// pixel offset. The new speed applies from the next AdaptVelocity call.
func (s *ScrollState) SetBaseSpeed(speed int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.baseSpeed = speed
}

func (s *ScrollState) OffsetPx() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.offsetPx
}

// ScrollFrame resolves the visible window of text lines for the current scroll
// state, returning the already-truncated lines and the sub-line pixel offset.
//
// It selects which lines are visible from offsetPx and rowHeight, truncates each
// to the column budget, and returns the negative sub-line offset so the frame
// scrolls content upward. This is pure frame resolution — it produces data for a
// style to place, not a ViewData. Returns (nil, 0) when there is nothing to show.
func ScrollFrame(scroll *ScrollState, lines []string, layout Layout) (visible []string, offsetY int) {
	if layout.RowHeight <= 0 || layout.VisibleRows <= 0 || len(lines) == 0 {
		return nil, 0
	}

	scroll.mu.Lock()
	offsetPx := scroll.offsetPx
	scroll.mu.Unlock()

	// Full lines scrolled past, and the sub-line offset within the top line.
	linesScrolled := offsetPx / layout.RowHeight
	subOffset := offsetPx % layout.RowHeight

	// One extra line beyond the visible count covers the partially-scrolled top
	// and bottom rows.
	needLines := layout.VisibleRows + 1

	startIdx := linesScrolled
	if startIdx >= len(lines) {
		// Scrolled past all content.
		return nil, 0
	}
	endIdx := startIdx + needLines
	if endIdx > len(lines) {
		endIdx = len(lines)
	}

	visible = make([]string, endIdx-startIdx)
	copy(visible, lines[startIdx:endIdx])

	if layout.VisibleColumns > 0 {
		for i, line := range visible {
			visible[i] = truncateCols(line, layout.VisibleColumns)
		}
	}

	return visible, -subOffset
}

// truncateCols trims a line to at most cols runes, counting runes so multi-byte
// content is not split mid-codepoint.
func truncateCols(s string, cols int) string {
	if cols <= 0 {
		return ""
	}
	n := 0
	for i := range s {
		if n == cols {
			return s[:i]
		}
		n++
	}
	return s
}
