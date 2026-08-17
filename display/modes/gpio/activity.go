package gpio

import (
	"image"
	"sync"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

const activityWindowSize = 32

// activityState stores per-pin toggle history for the "activity" style.
type activityState struct {
	sync.RWMutex
	history map[int][]float64 // BCM pin number → sliding window of up to 32 entries
	prev    map[int]bool      // BCM pin number → previous level for change detection
}

// activity is the package-level singleton holding activity tracking state.
var activity = &activityState{
	history: make(map[int][]float64),
	prev:    make(map[int]bool),
}

// RecordActivity compares previous and current pin states and records transitions.
// For each pin, if the Level changed between prev and curr, record 1.0; otherwise 0.0.
// The sliding window holds at most 32 entries; when full, the oldest is evicted.
//
// Framework pattern demonstrated: PersistentWidgets — maintains per-pin toggle
// history that feeds sparkline visualization across frames.
func RecordActivity(prev, curr []gpiomgr.PinState) {
	activity.Lock()
	defer activity.Unlock()

	// Build a lookup for previous snapshot levels.
	prevLevels := make(map[int]bool, len(prev))
	for _, p := range prev {
		prevLevels[p.Number] = bool(p.Level)
	}

	for _, c := range curr {
		prevLevel, hasPrev := prevLevels[c.Number]
		currLevel := bool(c.Level)

		var value float64
		if hasPrev && prevLevel != currLevel {
			value = 1.0
		}

		// Also check against stored prev map for pins not in prev slice.
		if !hasPrev {
			if storedPrev, ok := activity.prev[c.Number]; ok && storedPrev != currLevel {
				value = 1.0
			}
		}

		hist := activity.history[c.Number]
		if len(hist) >= activityWindowSize {
			// Shift left: drop the oldest entry.
			hist = hist[1:]
		}
		hist = append(hist, value)
		activity.history[c.Number] = hist

		// Update stored previous level.
		activity.prev[c.Number] = currLevel
	}
}

// GetActivityHistory returns a copy of the activity history for a given pin.
// If no history exists, returns nil. The caller receives a snapshot that
// won't be mutated by concurrent RecordActivity calls.
func GetActivityHistory(pin int) []float64 {
	activity.RLock()
	defer activity.RUnlock()
	h, ok := activity.history[pin]
	if !ok {
		return nil
	}
	cp := make([]float64, len(h))
	copy(cp, h)
	return cp
}

// ResetActivity clears all activity tracking state.
// Exported for use in tests to ensure a clean starting state.
func ResetActivity() {
	activity.Lock()
	defer activity.Unlock()
	activity.history = make(map[int][]float64)
	activity.prev = make(map[int]bool)
}

// BuildActivityView constructs a ViewData for the "activity" style.
// It renders one sparkline per visible output pin showing toggle history.
// Pins with no history get 32 zero-valued data points.
// If there are no output pins, no sparklines are rendered.
//
// Framework pattern demonstrated: Compositor — per-pin sparkline widget compositing
// with PersistentWidgets for sliding-window history visualization.
func BuildActivityView(pins []gpiomgr.PinState, hints textlayout.TextHints, face font.Face) style.ViewData {
	maxRows := textlayout.MaxVisibleRows(hints, 0)

	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	pixelWidth := hints.PixelWidth

	// Filter to output pins only.
	var outputPins []gpiomgr.PinState
	for _, p := range pins {
		if p.Mode == gpiomgr.ModeOutput {
			outputPins = append(outputPins, p)
		}
	}

	// Limit to maxRows visible pins.
	if maxRows > 0 && len(outputPins) > maxRows {
		outputPins = outputPins[:maxRows]
	}

	// If no output pins, return empty ViewData.
	if len(outputPins) == 0 {
		return style.ViewData{
			Static: false,
		}
	}

	ctx := widgets.SuppressionContext{
		AvailableWidth:  pixelWidth,
		AvailableHeight: rowHeight * len(outputPins),
	}
	comp := widgets.NewCompositor(ctx)

	activity.RLock()
	defer activity.RUnlock()

	for i, p := range outputPins {
		// Get history for this pin, or default to 32 zeros.
		data := activity.history[p.Number]
		if len(data) == 0 {
			data = make([]float64, activityWindowSize)
		} else {
			// Pad to 32 entries if history is shorter.
			if len(data) < activityWindowSize {
				padded := make([]float64, activityWindowSize)
				// Place existing data at the end (most recent).
				copy(padded[activityWindowSize-len(data):], data)
				data = padded
			} else {
				// Make a copy to avoid holding the lock over rendering.
				cp := make([]float64, activityWindowSize)
				copy(cp, data[len(data)-activityWindowSize:])
				data = cp
			}
		}

		bounds := image.Rect(0, i*rowHeight, pixelWidth, (i+1)*rowHeight)
		comp.Add(sparkline.New(sparkline.Config{
			Data:       data,
			Style:      sparkline.Line,
			Bounds:     bounds,
			Foreground: ColorHigh,
		}))
	}

	return style.ViewData{
		Sprites: comp.Sprites(),
		Static:  false,
	}
}
