package source

import (
	"sync"

	"github.com/databeast/cyberhud/display/style"
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

// renderSmoothScroll produces a ViewData representing the visible window of
// text lines given the current scroll state, buffer snapshot, and layout.
//
// The function selects which lines are visible based on offsetPx and rowHeight,
// and sets OffsetY to the negative sub-line pixel offset so the renderer
// positions text with sub-pixel vertical scrolling.
//
// This is a standalone rendering function — it does NOT modify BuildView
// directly. Task 13.1 handles BuildView wiring.
func RenderSmoothScroll(scroll *ScrollState, lines []string, layout Layout) style.ViewData {
	if layout.RowHeight <= 0 || layout.VisibleRows <= 0 || len(lines) == 0 {
		return style.ViewData{
			Items:  []string{},
			Static: true,
		}
	}

	scroll.mu.Lock()
	offsetPx := scroll.offsetPx
	scroll.mu.Unlock()

	// Determine how many full lines have been scrolled past.
	linesScrolled := 0
	if layout.RowHeight > 0 {
		linesScrolled = offsetPx / layout.RowHeight
	}

	// Sub-line pixel offset within the top visible line.
	subOffset := offsetPx % layout.RowHeight

	// We need visibleRows + 1 lines to account for the partially visible
	// top line being scrolled out and the partially visible bottom line
	// being scrolled in.
	needLines := layout.VisibleRows + 1

	// Calculate the start index into the buffer.
	startIdx := linesScrolled
	if startIdx >= len(lines) {
		// We've scrolled past all content — show empty.
		return style.ViewData{
			Items:  []string{},
			Static: true,
		}
	}

	// Calculate the end index.
	endIdx := startIdx + needLines
	if endIdx > len(lines) {
		endIdx = len(lines)
	}

	// Extract the visible slice.
	visible := make([]string, endIdx-startIdx)
	copy(visible, lines[startIdx:endIdx])

	// Truncate lines to VisibleColumns if layout provides it.
	if layout.VisibleColumns > 0 {
		for i, line := range visible {
			if len(line) > layout.VisibleColumns {
				visible[i] = line[:layout.VisibleColumns]
			}
		}
	}

	return style.ViewData{
		Items:   visible,
		OffsetY: -subOffset, // negative offset scrolls content upward
		Static:  false,
	}
}
