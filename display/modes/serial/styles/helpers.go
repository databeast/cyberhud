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
	"strings"

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

// compactHeader returns the "port @baud" header row, compacted to fit maxChars by
// shedding low-value parts rather than clipping meaningful ones. A non-positive
// maxChars imposes no budget (the full form).
//
// Priority when space is tight: the port basename identifies the device, so the
// "/dev/" prefix is dropped first; the baud is kept because it is the other half
// of the setting. The device basename itself is never clipped — a half-port name
// is useless — so on a panel too narrow even for "ttyUSB0 @115200" the header
// falls back to the basename alone.
func compactHeader(snap source.Snapshot, maxChars int) string {
	port := snap.Port
	if port == "" {
		if snap.AutoSelect {
			port = "auto"
		} else {
			port = "(none)"
		}
	}

	full := fmt.Sprintf("%s @%d", port, snap.Baud)
	if maxChars <= 0 || len([]rune(full)) <= maxChars {
		return full
	}

	// Shed the "/dev/" prefix: the basename is what identifies the device.
	short := fmt.Sprintf("%s @%d", basename(port), snap.Baud)
	if len([]rune(short)) <= maxChars {
		return short
	}

	// Still too wide: keep the basename alone (never a clipped port name).
	return basename(port)
}

// basename returns the final path element of a device path, so "/dev/ttyUSB0"
// becomes "ttyUSB0". Paths without a slash are returned unchanged.
func basename(p string) string {
	if i := strings.LastIndexByte(p, '/'); i >= 0 && i < len(p)-1 {
		return p[i+1:]
	}
	return p
}

// statusText returns the connection-state string and its colour, compacted to
// maxChars. OK and disconnected are already short; the error case sheds its
// detail message (which belongs in the output stream, not a one-line status) down
// to a bare "ERR" rather than clipping the message mid-word. A non-positive
// maxChars imposes no budget.
func statusText(snap source.Snapshot, maxChars int) (string, color.RGBA) {
	switch {
	case snap.Connected:
		return "OK", colorSerialOK
	case snap.LastError != "":
		full := "ERR: " + snap.LastError
		if maxChars <= 0 || len([]rune(full)) <= maxChars {
			return full, colorSerialErr
		}
		return "ERR", colorSerialErr
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
