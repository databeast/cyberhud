package styles

// Shared helpers for the serial instrument-panel engine.
//
// These are pure text/geometry helpers the fast and slow engines both use:
// the header/status strings, the serial accent palette, ANSI colour conversion,
// and the row-selection math. They answer "what text, what colour, how many
// rows" — not "where on screen", which is the engine's job.

import (
	"fmt"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// Serial accent palette. The header is a cool blue, status is traffic-light
// (green OK / amber waiting / red error), and data rows are near-white.
var (
	colorSerialHeader = color.RGBA{R: 0x50, G: 0xC8, B: 0xFF, A: 0xFF}
	colorSerialOK     = color.RGBA{R: 0x40, G: 0xE0, B: 0x80, A: 0xFF}
	colorSerialErr    = color.RGBA{R: 0xFF, G: 0x70, B: 0x70, A: 0xFF}
	colorSerialWarn   = color.RGBA{R: 0xFF, G: 0xC8, B: 0x40, A: 0xFF}
	colorSerialData   = color.RGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF}

	colorWhite = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	ledGreen   = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	ledRed     = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
)

// compactHeader returns a short "port @baud" string for the header row.
func compactHeader(snap source.Snapshot) string {
	port := snap.Port
	if port == "" {
		if snap.AutoSelect {
			port = "auto"
		} else {
			port = "(none)"
		}
	}
	return fmt.Sprintf("%s @%d", port, snap.Baud)
}

// statusText returns the connection-state word and its colour: OK (connected),
// ERR + message (error present), or disconnected.
func statusText(snap source.Snapshot) (string, color.RGBA) {
	switch {
	case snap.Connected:
		return "OK", colorSerialOK
	case snap.LastError != "":
		return "ERR: " + snap.LastError, colorSerialErr
	default:
		return "disconnected", colorSerialWarn
	}
}

// noPort reports whether the monitor has no port at all (no manual port, no
// auto-select, not connected) — the cold-start state that shows a bare prompt
// instead of a header/status pair.
func noPort(snap source.Snapshot) bool {
	return snap.Port == "" && !snap.AutoSelect && !snap.Connected
}

// dataLines returns the tail of the output buffer that fits budget rows, each
// stripped of ANSI escapes and truncated to maxChars. Returns nil when budget
// is non-positive.
func dataLines(snap source.Snapshot, budget, maxChars int) []string {
	if budget <= 0 {
		return nil
	}
	lines := snap.Lines
	if len(lines) > budget {
		lines = lines[len(lines)-budget:]
	}
	out := make([]string, len(lines))
	for i, raw := range lines {
		text, _ := source.ParseLine(raw)
		if maxChars > 0 {
			text = textlayout.Truncate(text, maxChars)
		}
		out[i] = text
	}
	return out
}

// colorToRGBA converts a color.Color to color.RGBA, defaulting to white for nil.
func colorToRGBA(c color.Color) color.RGBA {
	if c == nil {
		return colorWhite
	}
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}
