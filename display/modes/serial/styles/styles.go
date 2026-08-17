package styles

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/scrollbar"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
)

var (
	colorConnected    color.Color = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	colorDisconnected color.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
)

func buildItemsFromSnapshot(snap source.Snapshot) []string {
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

func buildColors(snap source.Snapshot, items []string) []color.Color {
	colors := make([]color.Color, len(items))
	for i, item := range items {
		if strings.HasPrefix(item, "State: connected") {
			colors[i] = colorConnected
		} else if strings.HasPrefix(item, "State: disconnected") {
			colors[i] = colorDisconnected
		}
	}
	return colors
}

// buildDefault renders header row (port, baud, state) with LED for connection
// status, followed by scrollable serial output lines using textlabel widgets.
// Uses the Compositor pattern to replace manual render→nil-check→append.
func buildDefault(snap source.Snapshot, p source.Policy, hints textlayout.TextHints, face font.Face) style.ViewData {
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	ox, oy := bridge.ContentOrigin()

	maxVisible := computeMaxVisible(hints, p)

	// Keep fallback Items/Colors for text rendering.
	items := buildItemsFromSnapshot(snap)
	colors := buildColors(snap, items)

	// Truncate to maxVisible rows.
	if len(items) > maxVisible {
		items = items[len(items)-maxVisible:]
		colors = colors[len(colors)-maxVisible:]
	}

	// Layout metrics from hints (already updated via WithFont when face available).
	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	glyphHeight := hints.GlyphHeight
	if glyphHeight <= 0 {
		glyphHeight = textlayout.GlyphHeight
	}
	glyphAdvance := hints.GlyphAdvance
	if glyphAdvance <= 0 {
		glyphAdvance = textlayout.GlyphAdvance
	}

	// Construct SuppressionContext from panel dimensions.
	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)

	xOffset := 0

	// LED for connection status.
	var ledState led.State
	var ledColor color.RGBA
	if snap.Connected {
		ledState = led.On
		ledColor = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	} else {
		ledState = led.Off
		ledColor = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	}

	ledWidget := led.New(led.Config{
		State:      ledState,
		Brightness: -1.0,
		Diameter:   glyphHeight,
		Bounds:     image.Rect(ox, oy, ox+glyphHeight, oy+glyphHeight),
		Foreground: ledColor,
	})
	comp.Add(ledWidget)
	// If LED rendered, offset subsequent content.
	if glyphHeight >= 3 {
		xOffset = glyphHeight + 1
	}

	// Header text: port/baud on first row.
	headerText := buildHeaderText(snap)
	headerBounds := image.Rect(ox+xOffset, oy, ox+hints.PixelWidth, oy+rowHeight)
	comp.Add(textlabel.New(textlabel.Config{
		Text:       headerText,
		Bounds:     headerBounds,
		Font:       face,
		Foreground: color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	}))

	// State row (second header row).
	stateText, stateFg := buildStateText(snap)
	stateBounds := image.Rect(ox+xOffset, oy+rowHeight, ox+hints.PixelWidth, oy+rowHeight*2)
	comp.Add(textlabel.New(textlabel.Config{
		Text:       stateText,
		Bounds:     stateBounds,
		Font:       face,
		Foreground: stateFg,
	}))

	// Serial output lines starting after header rows (row index 2).
	lineStartY := oy + rowHeight*2
	visibleLines := snap.Lines
	if len(visibleLines) == 0 {
		// Placeholder.
		placeholderBounds := image.Rect(ox, lineStartY, ox+hints.PixelWidth, lineStartY+rowHeight)
		comp.Add(textlabel.New(textlabel.Config{
			Text:       "(no data yet)",
			Bounds:     placeholderBounds,
			Font:       face,
			Foreground: color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		}))
	} else {
		// Truncate lines to fit visible area (account for 2 header rows).
		maxLineRows := maxVisible - 2
		if maxLineRows < 0 {
			maxLineRows = 0
		}
		if len(visibleLines) > maxLineRows {
			visibleLines = visibleLines[len(visibleLines)-maxLineRows:]
		}

		for i, rawLine := range visibleLines {
			y := lineStartY + i*rowHeight
			if y+rowHeight > hints.PixelHeight && hints.PixelHeight > 0 {
				break
			}

			// Parse ANSI colors from the line.
			text, segments := source.ParseLine(rawLine)

			if len(segments) <= 1 {
				// Single color (or no ANSI) — render as one textlabel.
				var fg color.RGBA
				if len(segments) == 1 && segments[0].Color != nil {
					fg = colorToRGBA(segments[0].Color)
				} else {
					fg = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
				}
				lineBounds := image.Rect(ox, y, ox+hints.PixelWidth, y+rowHeight)
				comp.Add(textlabel.New(textlabel.Config{
					Text:       text,
					Bounds:     lineBounds,
					Font:       face,
					Foreground: fg,
				}))
			} else {
				// Multiple color segments — render each segment separately.
				xCursor := ox
				for _, seg := range segments {
					segText := text[seg.Start : seg.Start+seg.Length]
					segWidth := len(segText) * glyphAdvance

					var fg color.RGBA
					if seg.Color != nil {
						fg = colorToRGBA(seg.Color)
					} else {
						fg = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
					}

					segBounds := image.Rect(xCursor, y, xCursor+segWidth, y+rowHeight)
					comp.Add(textlabel.New(textlabel.Config{
						Text:       segText,
						Bounds:     segBounds,
						Font:       face,
						Foreground: fg,
					}))
					xCursor += segWidth
				}
			}
		}
	}

	sprites := comp.Sprites()
	if len(sprites) > 0 {
		// Sprite path is authoritative — suppress Items to prevent double-rendering.
		return style.ViewData{
			Sprites: sprites,
			Static:  false,
		}
	}
	// Fallback: no sprites rendered (e.g. nil face) — use items-based rendering.
	return style.ViewData{
		Items:  items,
		Colors: colors,
		Static: false,
	}
}
func buildHeaderText(snap source.Snapshot) string {
	if snap.Port != "" {
		return fmt.Sprintf("Port: %s @%d", snap.Port, snap.Baud)
	} else if snap.AutoSelect {
		return fmt.Sprintf("Port: auto-select @%d", snap.Baud)
	}
	return fmt.Sprintf("Port: (manual) @%d", snap.Baud)
}

// buildStateText returns the state text and foreground color for the state row.
func buildStateText(snap source.Snapshot) (string, color.RGBA) {
	if snap.Connected {
		return "State: connected", color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	}
	return "State: disconnected", color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
}

// colorToRGBA converts a color.Color to color.RGBA.
func colorToRGBA(c color.Color) color.RGBA {
	if c == nil {
		return color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	}
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

// buildRaw renders only serial buffer lines with no header, no widgets, no sprites.
// ANSI escape sequences are stripped and the first segment's color is used per row.
func buildRaw(snap source.Snapshot, p source.Policy, hints textlayout.TextHints, face font.Face) style.ViewData {
	_ = layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	maxVisible := computeMaxVisible(hints, p)

	var items []string
	var colors []color.Color

	if len(snap.Lines) == 0 {
		items = []string{"(no data yet)"}
		colors = []color.Color{nil}
	} else {
		items = make([]string, 0, len(snap.Lines))
		colors = make([]color.Color, 0, len(snap.Lines))
		for _, raw := range snap.Lines {
			text, segments := source.ParseLine(raw)
			items = append(items, text)
			// Use the first segment's color for the entire row.
			var rowColor color.Color
			if len(segments) > 0 {
				rowColor = segments[0].Color
			}
			colors = append(colors, rowColor)
		}
	}

	// Truncate to maxVisible rows.
	if len(items) > maxVisible {
		items = items[len(items)-maxVisible:]
		colors = colors[len(colors)-maxVisible:]
	}

	return style.ViewData{
		Items:   items,
		Colors:  colors,
		Sprites: nil,
		Static:  false,
	}
}

// buildDashboard renders a bordered frame with LED, port/baud label, sparkline,
// progressbar, and textlabel-based serial output. Static=true (non-scrollable).
// Uses the Compositor pattern to replace manual render→nil-check→append.
func buildDashboard(snap source.Snapshot, p source.Policy, hints textlayout.TextHints, face font.Face) style.ViewData {
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	ox, oy := bridge.ContentOrigin()

	maxVisible := computeMaxVisible(hints, p)

	// Keep fallback Items/Colors for text rendering.
	items := buildItemsFromSnapshot(snap)
	colors := buildColors(snap, items)

	// Truncate to maxVisible rows.
	if len(items) > maxVisible {
		items = items[len(items)-maxVisible:]
		colors = colors[len(colors)-maxVisible:]
	}

	// Layout metrics from hints (already updated via WithFont when face available).
	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	glyphHeight := hints.GlyphHeight
	if glyphHeight <= 0 {
		glyphHeight = textlayout.GlyphHeight
	}
	glyphAdvance := hints.GlyphAdvance
	if glyphAdvance <= 0 {
		glyphAdvance = textlayout.GlyphAdvance
	}

	// Construct SuppressionContext from panel dimensions.
	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)

	currentY := oy

	// --- Row 0: LED (connection) + port/baud textlabel ---
	xOffset := 0

	// LED for connection status.
	var ledState led.State
	var ledColor color.RGBA
	if snap.Connected {
		ledState = led.On
		ledColor = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	} else {
		ledState = led.Off
		ledColor = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	}

	ledWidget := led.New(led.Config{
		State:      ledState,
		Brightness: -1.0,
		Diameter:   glyphHeight,
		Bounds:     image.Rect(ox, currentY, ox+glyphHeight, currentY+glyphHeight),
		Foreground: ledColor,
	})
	comp.Add(ledWidget)
	if glyphHeight >= 3 {
		xOffset = glyphHeight + 1
	}

	// Port/baud text using textlabel.
	headerText := buildHeaderText(snap)
	headerBounds := image.Rect(ox+xOffset, currentY, ox+hints.PixelWidth, currentY+rowHeight)
	comp.Add(textlabel.New(textlabel.Config{
		Text:       headerText,
		Bounds:     headerBounds,
		Font:       face,
		Foreground: color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	}))
	currentY += rowHeight

	// --- Row 1: Sparkline (throughput graph, full width, height ~10px) ---
	sparklineHeight := 10
	if rowHeight > sparklineHeight {
		sparklineHeight = rowHeight
	}

	// Normalize throughput: maxBPS = Baud / 10 (theoretical max bytes/sec for 8N1).
	maxBPS := p.Baud / 10
	if maxBPS <= 0 {
		maxBPS = 1
	}
	throughputData := make([]float64, 32)
	for i, v := range snap.Throughput {
		ratio := float64(v) / float64(maxBPS)
		if ratio > 1.0 {
			ratio = 1.0
		}
		if ratio < 0.0 {
			ratio = 0.0
		}
		throughputData[i] = ratio
	}

	sparklineBounds := image.Rect(ox, currentY, ox+hints.PixelWidth, currentY+sparklineHeight)
	comp.Add(sparkline.New(sparkline.Config{
		Data:       throughputData,
		Style:      sparkline.Line,
		Bounds:     sparklineBounds,
		Foreground: color.RGBA{R: 0, G: 200, B: 200, A: 255},
	}))
	currentY += sparklineHeight

	// --- Row 2: Progressbar (buffer fill, full width, height ~4px) ---
	progressHeight := 4
	fillRatio := float64(len(snap.Lines)) / float64(p.MaxLines)
	if fillRatio > 1.0 {
		fillRatio = 1.0
	}
	if fillRatio < 0.0 {
		fillRatio = 0.0
	}

	progressBounds := image.Rect(ox, currentY, ox+hints.PixelWidth, currentY+progressHeight)
	comp.Add(progressbar.New(progressbar.Config{
		Style:      progressbar.Linear,
		Value:      fillRatio,
		Bounds:     progressBounds,
		Foreground: color.RGBA{R: 0, G: 180, B: 0, A: 255},
	}))
	currentY += progressHeight

	// --- Remaining rows: Serial output lines via textlabel with ANSI colors ---
	lineStartY := currentY
	scrollbarWidth := 4 // Reserve space for scrollbar on right edge.
	contentWidth := hints.PixelWidth - scrollbarWidth
	if contentWidth < 1 {
		contentWidth = hints.PixelWidth
	}

	// Compute how many serial lines can fit in the remaining space.
	remainingHeight := hints.PixelHeight - (lineStartY - oy)
	if remainingHeight < 0 {
		remainingHeight = 0
	}
	maxLineRows := 0
	if rowHeight > 0 {
		maxLineRows = remainingHeight / rowHeight
	}
	if maxLineRows <= 0 {
		maxLineRows = maxVisible
	}

	visibleLines := snap.Lines
	if len(visibleLines) == 0 {
		// Placeholder.
		placeholderBounds := image.Rect(ox, lineStartY, ox+contentWidth, lineStartY+rowHeight)
		comp.Add(textlabel.New(textlabel.Config{
			Text:       "(no data yet)",
			Bounds:     placeholderBounds,
			Font:       face,
			Foreground: color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		}))
	} else {
		// Apply scroll offset and truncation.
		if len(visibleLines) > maxLineRows {
			// Use scroll offset to determine which lines to show.
			startIdx := len(visibleLines) - maxLineRows
			if snap.ScrollOffset > 0 && snap.ScrollOffset < startIdx {
				startIdx = snap.ScrollOffset
			}
			endIdx := startIdx + maxLineRows
			if endIdx > len(visibleLines) {
				endIdx = len(visibleLines)
			}
			visibleLines = visibleLines[startIdx:endIdx]
		}

		for i, rawLine := range visibleLines {
			y := lineStartY + i*rowHeight
			if hints.PixelHeight > 0 && y+rowHeight > hints.PixelHeight {
				break
			}

			// Parse ANSI colors from the line.
			text, segments := source.ParseLine(rawLine)

			if len(segments) <= 1 {
				// Single color (or no ANSI) — render as one textlabel.
				var fg color.RGBA
				if len(segments) == 1 && segments[0].Color != nil {
					fg = colorToRGBA(segments[0].Color)
				} else {
					fg = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
				}
				lineBounds := image.Rect(ox, y, ox+contentWidth, y+rowHeight)
				comp.Add(textlabel.New(textlabel.Config{
					Text:       text,
					Bounds:     lineBounds,
					Font:       face,
					Foreground: fg,
				}))
			} else {
				// Multiple color segments — render each segment separately.
				xCursor := ox
				for _, seg := range segments {
					segText := text[seg.Start : seg.Start+seg.Length]
					segWidth := len(segText) * glyphAdvance

					var fg color.RGBA
					if seg.Color != nil {
						fg = colorToRGBA(seg.Color)
					} else {
						fg = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
					}

					segBounds := image.Rect(xCursor, y, xCursor+segWidth, y+rowHeight)
					comp.Add(textlabel.New(textlabel.Config{
						Text:       segText,
						Bounds:     segBounds,
						Font:       face,
						Foreground: fg,
					}))
					xCursor += segWidth
				}
			}
		}
	}

	// --- Right edge: Scrollbar (if buffer > visible rows AND scroll offset != tail) ---
	tailOffset := len(snap.Lines) - maxLineRows
	if tailOffset < 0 {
		tailOffset = 0
	}
	showScrollbar := len(snap.Lines) > maxLineRows && snap.ScrollOffset != tailOffset
	comp.AddIf(showScrollbar, scrollbar.New(scrollbar.Config{
		TotalItems:   len(snap.Lines),
		VisibleItems: maxLineRows,
		ScrollOffset: snap.ScrollOffset,
		Bounds: image.Rect(
			ox+hints.PixelWidth-scrollbarWidth, lineStartY,
			ox+hints.PixelWidth, hints.PixelHeight,
		),
		Foreground: color.RGBA{R: 128, G: 128, B: 128, A: 255},
	}))

	sprites := comp.Sprites()
	if len(sprites) > 0 {
		return style.ViewData{
			Sprites: sprites,
			Static:  true,
		}
	}
	return style.ViewData{
		Items:  items,
		Colors: colors,
		Static: true,
	}
}

// buildCompact renders a single-row status bar (LED + port + baud) followed by
// serial output lines using the tier-catalog-resolved face to maximize visible
// line count. Uses the Compositor pattern.
func buildCompact(snap source.Snapshot, p source.Policy, hints textlayout.TextHints, face font.Face) style.ViewData {
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	ox, oy := bridge.ContentOrigin()

	maxVisible := computeMaxVisible(hints, p)

	// Keep fallback Items/Colors for text rendering.
	items := buildItemsFromSnapshot(snap)
	colors := buildColors(snap, items)

	// Truncate to maxVisible rows.
	if len(items) > maxVisible {
		items = items[len(items)-maxVisible:]
		colors = colors[len(colors)-maxVisible:]
	}

	// Layout metrics from hints (already updated via WithFont when face available).
	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	glyphHeight := hints.GlyphHeight
	if glyphHeight <= 0 {
		glyphHeight = textlayout.GlyphHeight
	}
	glyphAdvance := hints.GlyphAdvance
	if glyphAdvance <= 0 {
		glyphAdvance = textlayout.GlyphAdvance
	}

	// Construct SuppressionContext from panel dimensions.
	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)

	xOffset := 0

	// LED for connection status in status bar row.
	var ledState led.State
	var ledColor color.RGBA
	if snap.Connected {
		ledState = led.On
		ledColor = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	} else {
		ledState = led.Off
		ledColor = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	}

	ledWidget := led.New(led.Config{
		State:      ledState,
		Brightness: -1.0,
		Diameter:   glyphHeight,
		Bounds:     image.Rect(ox, oy, ox+glyphHeight, oy+glyphHeight),
		Foreground: ledColor,
	})
	comp.Add(ledWidget)
	if glyphHeight >= 3 {
		xOffset = glyphHeight + 1
	}

	// Port@baud status text as textlabel in same row.
	statusText := fmt.Sprintf("%s@%d", portDisplay(snap), snap.Baud)
	statusBounds := image.Rect(ox+xOffset, oy, ox+hints.PixelWidth, oy+rowHeight)
	comp.Add(textlabel.New(textlabel.Config{
		Text:       statusText,
		Bounds:     statusBounds,
		Font:       face,
		Foreground: color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
	}))

	// Serial output lines starting after the status bar row (row 1 onward).
	lineStartY := oy + rowHeight
	visibleLines := snap.Lines
	if len(visibleLines) == 0 {
		// Placeholder.
		placeholderBounds := image.Rect(ox, lineStartY, ox+hints.PixelWidth, lineStartY+rowHeight)
		comp.Add(textlabel.New(textlabel.Config{
			Text:       "(no data yet)",
			Bounds:     placeholderBounds,
			Font:       face,
			Foreground: color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		}))
	} else {
		// Maximize visible serial line rows (account for 1 status bar row).
		maxLineRows := maxVisible - 1
		if maxLineRows < 0 {
			maxLineRows = 0
		}
		if len(visibleLines) > maxLineRows {
			visibleLines = visibleLines[len(visibleLines)-maxLineRows:]
		}

		for i, rawLine := range visibleLines {
			y := lineStartY + i*rowHeight
			if y+rowHeight > hints.PixelHeight && hints.PixelHeight > 0 {
				break
			}

			// Parse ANSI colors from the line.
			text, segments := source.ParseLine(rawLine)

			if len(segments) <= 1 {
				// Single color (or no ANSI) — render as one textlabel.
				var fg color.RGBA
				if len(segments) == 1 && segments[0].Color != nil {
					fg = colorToRGBA(segments[0].Color)
				} else {
					fg = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
				}
				lineBounds := image.Rect(ox, y, ox+hints.PixelWidth, y+rowHeight)
				comp.Add(textlabel.New(textlabel.Config{
					Text:       text,
					Bounds:     lineBounds,
					Font:       face,
					Foreground: fg,
				}))
			} else {
				// Multiple color segments — render each segment separately.
				xCursor := ox
				for _, seg := range segments {
					segText := text[seg.Start : seg.Start+seg.Length]
					segWidth := len(segText) * glyphAdvance

					var fg color.RGBA
					if seg.Color != nil {
						fg = colorToRGBA(seg.Color)
					} else {
						fg = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
					}

					segBounds := image.Rect(xCursor, y, xCursor+segWidth, y+rowHeight)
					comp.Add(textlabel.New(textlabel.Config{
						Text:       segText,
						Bounds:     segBounds,
						Font:       face,
						Foreground: fg,
					}))
					xCursor += segWidth
				}
			}
		}
	}

	sprites := comp.Sprites()
	if len(sprites) > 0 {
		return style.ViewData{
			Sprites: sprites,
			Static:  false,
		}
	}
	return style.ViewData{
		Items:  items,
		Colors: colors,
		Static: false,
	}
}

// buildFramed renders serial output inside a borderframe widget with 8px inset.
// Uses borderframe.New for a decorative tile border around the panel edge,
// with content inset by 8px on each side. Shows a scrollbar when the buffer
// exceeds visible rows and the scroll offset is not at the tail position.
// Uses the Compositor pattern.
func buildFramed(snap source.Snapshot, p source.Policy, hints textlayout.TextHints, face font.Face) style.ViewData {
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	ox, oy := bridge.ContentOrigin()

	maxVisible := computeMaxVisible(hints, p)

	// Fallback Items/Colors for text rendering.
	items := buildItemsFromSnapshot(snap)
	colors := buildColors(snap, items)

	// Truncate to maxVisible rows.
	if len(items) > maxVisible {
		items = items[len(items)-maxVisible:]
		colors = colors[len(colors)-maxVisible:]
	}

	// Construct SuppressionContext from panel dimensions.
	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)

	// Build border frame via Renderable — returns nil sprite if bounds < 16×16.
	borderWidget := borderframe.New(borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
	})
	comp.Add(borderWidget)

	// Determine content bounds. If border rendered (panel ≥ 16×16), inset by 8px.
	// The border is a panel-covering element; content inset accounts for border tiles.
	var contentBounds image.Rectangle
	if hints.PixelWidth >= 16 && hints.PixelHeight >= 16 {
		contentBounds = image.Rect(ox+8, oy+8, ox+hints.PixelWidth-8, oy+hints.PixelHeight-8)
	} else {
		contentBounds = image.Rect(ox, oy, ox+hints.PixelWidth, oy+hints.PixelHeight)
	}

	// Layout metrics from hints (already updated via WithFont when face available).
	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	glyphAdvance := hints.GlyphAdvance
	if glyphAdvance <= 0 {
		glyphAdvance = textlayout.GlyphAdvance
	}

	// Scrollbar width reserved on right edge of content area.
	scrollbarWidth := 4
	contentWidth := contentBounds.Dx() - scrollbarWidth
	if contentWidth < 1 {
		contentWidth = contentBounds.Dx()
	}

	// Compute how many serial lines fit in the content area.
	contentHeight := contentBounds.Dy()
	maxLineRows := 0
	if rowHeight > 0 && contentHeight > 0 {
		maxLineRows = contentHeight / rowHeight
	}
	if maxLineRows <= 0 {
		maxLineRows = maxVisible
	}

	// Determine visible lines with scroll offset.
	visibleLines := snap.Lines
	if len(visibleLines) == 0 {
		// Placeholder.
		placeholderBounds := image.Rect(
			contentBounds.Min.X, contentBounds.Min.Y,
			contentBounds.Min.X+contentWidth, contentBounds.Min.Y+rowHeight,
		)
		comp.Add(textlabel.New(textlabel.Config{
			Text:       "(no data yet)",
			Bounds:     placeholderBounds,
			Font:       face,
			Foreground: color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		}))
	} else {
		// Apply scroll offset and truncation.
		if len(visibleLines) > maxLineRows {
			startIdx := len(visibleLines) - maxLineRows
			if snap.ScrollOffset > 0 && snap.ScrollOffset < startIdx {
				startIdx = snap.ScrollOffset
			}
			endIdx := startIdx + maxLineRows
			if endIdx > len(visibleLines) {
				endIdx = len(visibleLines)
			}
			visibleLines = visibleLines[startIdx:endIdx]
		}

		for i, rawLine := range visibleLines {
			y := contentBounds.Min.Y + i*rowHeight
			if y+rowHeight > contentBounds.Max.Y {
				break
			}

			// Parse ANSI colors from the line.
			text, segments := source.ParseLine(rawLine)

			if len(segments) <= 1 {
				// Single color (or no ANSI) — render as one textlabel.
				var fg color.RGBA
				if len(segments) == 1 && segments[0].Color != nil {
					fg = colorToRGBA(segments[0].Color)
				} else {
					fg = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
				}
				lineBounds := image.Rect(contentBounds.Min.X, y, contentBounds.Min.X+contentWidth, y+rowHeight)
				comp.Add(textlabel.New(textlabel.Config{
					Text:       text,
					Bounds:     lineBounds,
					Font:       face,
					Foreground: fg,
				}))
			} else {
				// Multiple color segments — render each segment separately.
				xCursor := contentBounds.Min.X
				for _, seg := range segments {
					segText := text[seg.Start : seg.Start+seg.Length]
					segWidth := len(segText) * glyphAdvance

					var fg color.RGBA
					if seg.Color != nil {
						fg = colorToRGBA(seg.Color)
					} else {
						fg = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
					}

					segBounds := image.Rect(xCursor, y, xCursor+segWidth, y+rowHeight)
					comp.Add(textlabel.New(textlabel.Config{
						Text:       segText,
						Bounds:     segBounds,
						Font:       face,
						Foreground: fg,
					}))
					xCursor += segWidth
				}
			}
		}
	}

	// Scrollbar on right edge of content area when buffer > visible rows
	// AND scroll offset is not at the tail position.
	tailOffset := len(snap.Lines) - maxLineRows
	if tailOffset < 0 {
		tailOffset = 0
	}
	showScrollbar := len(snap.Lines) > maxLineRows && snap.ScrollOffset != tailOffset
	comp.AddIf(showScrollbar, scrollbar.New(scrollbar.Config{
		TotalItems:   len(snap.Lines),
		VisibleItems: maxLineRows,
		ScrollOffset: snap.ScrollOffset,
		Bounds: image.Rect(
			contentBounds.Max.X-scrollbarWidth, contentBounds.Min.Y,
			contentBounds.Max.X, contentBounds.Max.Y,
		),
		Foreground: color.RGBA{R: 128, G: 128, B: 128, A: 255},
	}))

	sprites := comp.Sprites()
	if len(sprites) > 0 {
		return style.ViewData{
			Sprites: sprites,
			Static:  false,
		}
	}
	return style.ViewData{
		Items:  items,
		Colors: colors,
		Static: false,
	}
}

// computeMaxVisible determines the maximum number of visible rows based on
// TextHints and policy. Uses textlayout.MaxVisibleRows if RowHeight/PixelHeight
// are positive, otherwise falls back to Policy.MaxLines.
func computeMaxVisible(hints textlayout.TextHints, p source.Policy) int {
	if hints.RowHeight > 0 && hints.PixelHeight > 0 {
		if mv := textlayout.MaxVisibleRows(hints, 0); mv > 0 {
			return mv
		}
	}
	return p.MaxLines
}

// portDisplay returns a formatted port string for compact display.
func portDisplay(snap source.Snapshot) string {
	if snap.Port != "" {
		return snap.Port
	}
	if snap.AutoSelect {
		return "auto"
	}
	return "(none)"
}

// statusColor returns the appropriate color for connection status indicators.
func statusColor(connected bool) color.Color {
	if connected {
		return colorConnected
	}
	return colorDisconnected
}
