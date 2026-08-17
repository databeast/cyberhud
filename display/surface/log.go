package surface

import (
	"image"
	"image/color"
)

// DrawCall records a single draw operation performed on a Surface.
type DrawCall struct {
	Type  string          // "text", "rect", "clear"
	X     int             // X coordinate (for text)
	Y     int             // Y coordinate (for text)
	Text  string          // text content (for text draws)
	Rect  image.Rectangle // rectangle (for rect draws)
	Color color.Color     // color used in the draw operation
}

// NewWithLog creates a Surface that records all draw operations in a log.
// The returned surface behaves identically to a normal Surface but also stores
// DrawCall entries that can be retrieved via DrawLog().
func NewWithLog(bounds image.Rectangle) *Surface {
	s := New(bounds)
	s.logging = true
	s.drawLog = nil
	return s
}

// DrawLog returns the recorded draw calls. Returns nil if logging is not enabled.
func (s *Surface) DrawLog() []DrawCall {
	if !s.logging {
		return nil
	}
	return s.drawLog
}

// ClearLog resets the draw log, keeping logging enabled.
func (s *Surface) ClearLog() {
	s.drawLog = nil
}
