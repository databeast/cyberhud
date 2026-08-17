package serial

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// BuildView returns the serial mode view data with adaptive font selection,
// style dispatch, and full style.ViewData population.
func BuildView(hints textlayout.TextHints) style.ViewData {
	if hints.PixelWidth == 0 || hints.PixelHeight == 0 {
		if current, ok := getPanelHints(); ok {
			hints = current
		}
	}
	snap := source.SnapshotNow()
	p := normalizePolicy(PolicySnapshot())

	// Registry-based dispatch: configured → alias → fitness.
	s, reason := style.ResolveStyle(serialRegistry, hints, "serial", p.Style)

	// Construct StyleContext for the style boundary.
	ctx := style.NewStyleContext(hints)
	vd := s.Build(snap, p, ctx)

	// Report style resolution metadata to the registry layer.
	vd.StyleReport = style.StyleReport{
		Name:   s.Name(),
		Reason: reason,
	}

	return vd
}

// Color constants for connection status.
var (
	colorConnected    color.Color = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	colorDisconnected color.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
)

// BuildItems returns text rows for the serial monitor display.
func BuildItems() []string {
	return BuildView(textlayout.TextHints{}).Items
}

// buildItemsFromSnapshot builds text items from a snapshot (used by style builders).
func buildItemsFromSnapshot(snap Snapshot) []string {
	if snap.Port == "" && len(snap.Lines) == 0 && snap.LastError == "" && !snap.Connected {
		return []string{"Waiting for serial port...", "Use `serial policy` or `serial set`"}
	}

	items := []string{}
	if snap.Port != "" {
		items = append(items, fmt.Sprintf("Port: %s @%d", snap.Port, snap.Baud))
	} else if snap.AutoSelect {
		items = append(items, fmt.Sprintf("Port: auto-select @%d", snap.Baud))
	} else {
		items = append(items, fmt.Sprintf("Port: (manual) @%d", snap.Baud))
	}

	if snap.Connected {
		items = append(items, "State: connected")
	} else {
		items = append(items, "State: disconnected")
	}
	if snap.LastError != "" {
		items = append(items, "Error: "+snap.LastError)
	}
	if len(snap.Lines) == 0 {
		items = append(items, "(no data yet)")
	} else {
		items = append(items, snap.Lines...)
	}
	return items
}

// buildColors creates a Colors slice for the default style, highlighting
// the connection status row with green (connected) or red (disconnected).
func buildColors(snap Snapshot, items []string) []color.Color {
	colors := make([]color.Color, len(items))

	// Find the status row index and assign color.
	for i, item := range items {
		if strings.HasPrefix(item, "State: connected") {
			colors[i] = colorConnected
		} else if strings.HasPrefix(item, "State: disconnected") {
			colors[i] = colorDisconnected
		}
		// All other rows (port info, error, serial output) use nil (default foreground).
	}
	return colors
}

// Signature returns a stable change token for UI refresh.
func Signature() uint32 {
	return region.CalcRegionCacheKey(signatureFromSnapshot(source.SnapshotNow(), PolicySnapshot()))
}

func signatureFromSnapshot(snap Snapshot, p Policy) string {
	conn := "0"
	if snap.Connected {
		conn = "1"
	}
	p = normalizePolicy(p)
	return fmt.Sprintf("serial:%d:%s:%s:%d:%t:%s:%d", snap.Sequence, conn, snap.Port, snap.Baud, snap.AutoSelect, snap.LastError, len(snap.Lines))
}
