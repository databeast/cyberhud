package borderframe

import (
	"math"
	"time"

	"github.com/databeast/cyberhud/display/widgets"
)

// Compile-time assertion: borderframeWidget satisfies widgets.Animated.
var _ widgets.Animated = (*borderframeWidget)(nil)

// animationState holds the running state for all frame-driven animations.
type animationState struct {
	// Pulse animation.
	pulsePhase float64 // [0.0, 1.0)

	// Scan line animation.
	scanPosition float64 // Current pixel offset along clockwise perimeter.

	// Corner flash animation.
	flashElapsed time.Duration // Time within current flash cycle.

	// Segment reveal animation.
	revealedTiles int           // Number of tiles currently visible.
	revealAccum   time.Duration // Accumulated time for next tile reveal.
	revealDone    bool          // All tiles have been revealed.
}

// hasAnimation returns true when at least one animation field in cfg is active.
func hasAnimation(cfg Config) bool {
	return cfg.PulseCycle > 0 ||
		cfg.ScanSpeed > 0 ||
		cfg.CornerFlash ||
		cfg.SegmentReveal
}

// Tick advances all active animations atomically based on elapsed duration.
// It is a no-op when no animation Config fields are active.
func (w *borderframeWidget) Tick(elapsed time.Duration) {
	if !hasAnimation(w.cfg) {
		return
	}

	// Pulse animation.
	if w.cfg.PulseCycle > 0 {
		w.tickPulse(elapsed)
	}

	// Scan line animation.
	if w.cfg.ScanSpeed > 0 {
		w.tickScanLine(elapsed)
	}

	// Corner flash animation.
	if w.cfg.CornerFlash && w.cfg.CornerAccent {
		w.tickCornerFlash(elapsed)
	}

	// Segment reveal animation.
	if w.cfg.SegmentReveal && w.cfg.RevealSpeed > 0 && !w.anim.revealDone {
		w.tickSegmentReveal(elapsed)
	}
}

// ---------------------------------------------------------------------------
// Pulse animation logic
// ---------------------------------------------------------------------------

// tickPulse advances pulse phase by (elapsed / PulseCycle), wrapping at 1.0.
func (w *borderframeWidget) tickPulse(elapsed time.Duration) {
	advance := float64(elapsed) / float64(w.cfg.PulseCycle)
	w.anim.pulsePhase += advance
	// Wrap at 1.0.
	w.anim.pulsePhase -= math.Floor(w.anim.pulsePhase)
}

// pulseIntensity returns the current glow intensity factor in [0.3, 1.0].
// Uses sinusoidal mapping: intensity = 0.3 + 0.7 * (0.5 + 0.5*cos(2π*phase))
// This gives: 100% at phase 0.0, 30% at phase 0.5, back to 100% at phase 1.0.
// When PulseCycle is zero, returns 1.0 (constant full intensity).
func (w *borderframeWidget) pulseIntensity() float64 {
	if w.cfg.PulseCycle <= 0 {
		return 1.0
	}
	return 0.3 + 0.7*(0.5+0.5*math.Cos(2*math.Pi*w.anim.pulsePhase))
}

// ---------------------------------------------------------------------------
// Scan line animation logic
// ---------------------------------------------------------------------------

// tickScanLine advances scan position by (ScanSpeed * elapsed.Seconds()) pixels
// along the clockwise perimeter, wrapping seamlessly at perimeter end.
func (w *borderframeWidget) tickScanLine(elapsed time.Duration) {
	// Determine effective perimeter in pixels.
	perimeter := w.perimeterPixels()

	// No-op conditions: ShowBorder explicitly false, or perimeter < 32.
	if w.cfg.ShowBorder != nil && !*w.cfg.ShowBorder {
		return
	}
	if perimeter < 32 {
		return
	}

	// Clamp ScanSpeed to [1, 1000].
	speed := clampInt(w.cfg.ScanSpeed, 1, 1000)

	// Effective ScanLength.
	scanLen := w.effectiveScanLength(perimeter)

	// If ScanLength >= perimeter, the scan line covers everything — no advancement.
	if scanLen >= perimeter {
		return
	}

	// Advance position.
	advance := float64(speed) * elapsed.Seconds()
	w.anim.scanPosition += advance

	// Wrap seamlessly at perimeter end (modulo total perimeter).
	perimF := float64(perimeter)
	if w.anim.scanPosition >= perimF {
		w.anim.scanPosition = math.Mod(w.anim.scanPosition, perimF)
	}
}

// perimeterPixels returns the total border perimeter in pixels.
// Perimeter = 2*(width + height) - 4*tileSize (avoid double-counting corners).
// For a frame of cols×rows tiles: perimeter in tiles = 2*(cols+rows-2), in pixels = that * tileSize.
func (w *borderframeWidget) perimeterPixels() int {
	width := w.cfg.Bounds.Dx()
	height := w.cfg.Bounds.Dy()
	if width < 16 || height < 16 {
		return 0
	}
	cols := width / tileSize
	rows := height / tileSize
	tiles := 2 * (cols + rows - 2)
	return tiles * tileSize
}

// effectiveScanLength returns the clamped ScanLength value.
// Default 16 if zero-value; clamp to [1, perimeter].
func (w *borderframeWidget) effectiveScanLength(perimeter int) int {
	length := w.cfg.ScanLength
	if length == 0 {
		length = 16
	}
	return clampInt(length, 1, perimeter)
}

// ---------------------------------------------------------------------------
// Corner flash animation logic
// ---------------------------------------------------------------------------

// tickCornerFlash advances flashElapsed by elapsed duration, wrapping at FlashInterval.
// The flash cycle has two phases:
//   - Peak phase: alpha = 1.0 for FlashDuration (default 150ms, clamped [50ms, 1000ms])
//   - Base phase: alpha = 0.4 for the remainder of FlashInterval (default 2000ms, clamped [200ms, 10000ms])
func (w *borderframeWidget) tickCornerFlash(elapsed time.Duration) {
	interval := w.clampedFlashInterval()

	w.anim.flashElapsed += elapsed
	if w.anim.flashElapsed >= interval {
		w.anim.flashElapsed = w.anim.flashElapsed % interval
	}
}

// cornerFlashAlpha returns the current alpha factor for corner accent tiles.
//   - If CornerFlash is false or CornerAccent is false, returns 1.0 (no modulation).
//   - If flashElapsed < clampedFlashDuration: returns 1.0 (peak).
//   - Otherwise: returns 0.4 (base).
func (w *borderframeWidget) cornerFlashAlpha() float64 {
	if !w.cfg.CornerFlash || !w.cfg.CornerAccent {
		return 1.0
	}

	duration := w.clampedFlashDuration()
	if w.anim.flashElapsed < duration {
		return 1.0 // Peak phase.
	}
	return 0.4 // Base phase.
}

// clampedFlashDuration returns the effective FlashDuration clamped to [50ms, 1000ms],
// defaulting to 150ms when the configured value is zero.
func (w *borderframeWidget) clampedFlashDuration() time.Duration {
	d := w.cfg.FlashDuration
	if d == 0 {
		d = 150 * time.Millisecond
	}
	return clampDuration(d, 50*time.Millisecond, 1000*time.Millisecond)
}

// clampedFlashInterval returns the effective FlashInterval clamped to [200ms, 10000ms],
// defaulting to 2000ms when the configured value is zero.
func (w *borderframeWidget) clampedFlashInterval() time.Duration {
	d := w.cfg.FlashInterval
	if d == 0 {
		d = 2000 * time.Millisecond
	}
	return clampDuration(d, 200*time.Millisecond, 10000*time.Millisecond)
}

// ---------------------------------------------------------------------------
// Segment reveal animation logic
// ---------------------------------------------------------------------------

// tickSegmentReveal accumulates time and reveals tiles at RevealSpeed tiles/sec.
// The reveal starts from the configured origin corner (RevealOrigin, default TopLeft)
// and traverses the perimeter clockwise. Sets revealDone when all tiles are revealed.
// When RevealSpeed ≤ 0, SegmentReveal is treated as disabled (all tiles visible).
func (w *borderframeWidget) tickSegmentReveal(elapsed time.Duration) {
	// RevealSpeed must be in valid range [1, 240]; the Tick guard already ensures > 0.
	speed := clampInt(w.cfg.RevealSpeed, 1, 240)

	totalTiles := w.perimeterTileCount()
	if totalTiles <= 0 {
		w.anim.revealDone = true
		return
	}

	// Accumulate elapsed time.
	w.anim.revealAccum += elapsed

	// Calculate how many tiles should be revealed based on accumulated time.
	interval := time.Second / time.Duration(speed)
	for w.anim.revealAccum >= interval && w.anim.revealedTiles < totalTiles {
		w.anim.revealedTiles++
		w.anim.revealAccum -= interval
	}

	// Mark done when all tiles revealed — hold fully-revealed state until Configure restarts.
	if w.anim.revealedTiles >= totalTiles {
		w.anim.revealedTiles = totalTiles
		w.anim.revealDone = true
	}
}

// revealedTileCount returns the number of currently visible perimeter tiles.
// If SegmentReveal is disabled (RevealSpeed ≤ 0) or reveal is complete,
// returns the total perimeter tile count (all tiles visible).
func (w *borderframeWidget) revealedTileCount() int {
	total := w.perimeterTileCount()
	if total <= 0 {
		return 0
	}

	// Disabled: all tiles visible immediately.
	if !w.cfg.SegmentReveal || w.cfg.RevealSpeed <= 0 {
		return total
	}

	// Reveal complete: all tiles visible.
	if w.anim.revealDone {
		return total
	}

	return w.anim.revealedTiles
}

// revealStartIndex returns the perimeter tile index where reveal begins for
// the given origin corner. The perimeter is traversed clockwise with indices:
//
//	Top:    [0, cols)              — cols tiles (TL at 0, TR at cols-1)
//	Right:  [cols, cols+rows-2)   — rows-2 tiles (excluding corners)
//	Bottom: [cols+rows-2, 2*cols+rows-2) — cols tiles (BR at cols+rows-2, BL at 2*cols+rows-3)
//	Left:   [2*cols+rows-2, 2*(cols+rows-2)) — rows-2 tiles (excluding corners)
//
// Origin corner mapping:
//   - TopLeft:     index 0
//   - TopRight:    index cols-1
//   - BottomRight: index cols+rows-3
//   - BottomLeft:  index 2*cols+rows-3
func revealStartIndex(cols, rows int, origin Corner) int {
	if cols < 2 || rows < 2 {
		return 0
	}

	switch origin {
	case TopLeft:
		return 0
	case TopRight:
		return cols - 1
	case BottomRight:
		// Top edge (cols) + right edge (rows-2) - 1 = cols + rows - 3.
		return cols + rows - 3
	case BottomLeft:
		// Top edge (cols) + right edge (rows-2) + bottom-right-to-BL (cols-1) = 2*cols + rows - 3.
		return 2*cols + rows - 3
	default:
		return 0
	}
}

// perimeterTileCount returns the total number of perimeter tiles.
func (w *borderframeWidget) perimeterTileCount() int {
	width := w.cfg.Bounds.Dx()
	height := w.cfg.Bounds.Dy()
	if width < 16 || height < 16 {
		return 0
	}
	cols := width / tileSize
	rows := height / tileSize
	return 2 * (cols + rows - 2)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// clampInt clamps v to [min, max].
func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// clampDuration clamps d to [min, max].
func clampDuration(d, min, max time.Duration) time.Duration {
	if d < min {
		return min
	}
	if d > max {
		return max
	}
	return d
}
