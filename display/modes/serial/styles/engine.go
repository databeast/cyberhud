package styles

// The serial instrument-panel engine.
//
// Serial's one look is a console instrument panel: a connection LED, a port/baud
// header, a connection-state line, an optional throughput sparkline + buffer
// gauge, and the ANSI-colored output below. buildFastPanel renders the live
// version (animated LED, throughput sparkline) for refresh-capable panels;
// buildSlowPanel renders an e-paper-safe static version (no animated widgets)
// for slow panels. Both answer "where does this go?" for an already-settled
// snapshot — no clock, no per-frame state.
//
// Both degrade by panel size: the chrome bands (sparkline, buffer gauge) appear
// only when there is vertical room, and on very short panels the panel collapses
// to a compact header + status + data, the way wifi drops its detail/IP bands.
// Policy.Raw suppresses the chrome entirely on either engine, leaving the output
// bytes alone — the one behavioral variant, honoured identically at every size.

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
)

const serialFontFamily = "spleen"

// sparklineMinHeight is the panel content height below which the throughput
// sparkline and buffer gauge are dropped — the chrome would crowd the data.
const sparklineMinHeight = 96

// buildFastPanel renders the live instrument panel for refresh-capable panels:
// an animated connection LED, the header/status chrome, a throughput sparkline
// and buffer gauge when the panel is tall enough, and the ANSI-colored output.
// Non-static so the LED and sparkline refresh.
func buildFastPanel(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildPanel(snap, pol, ctx, d, panelOpts{color: d.p.Color, live: true, static: false})
}

// buildSlowPanel renders the e-paper-safe static panel: the same instrument
// layout, but with a static LED (no blink), no throughput sparkline (animated
// chrome shimmers under 1-bit refresh), and Static output.
func buildSlowPanel(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildPanel(snap, pol, ctx, d, panelOpts{color: d.p.Color, live: false, static: true})
}

// panelOpts carries the engine-variant choices into the shared builder.
type panelOpts struct {
	color  bool // engage the accent palette
	live   bool // animated LED + throughput sparkline (fast panels only)
	static bool // emit a static frame (slow panels)
}

// buildPanel is the shared instrument-panel assembly for both engine variants.
func buildPanel(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, _ def, opts panelOpts) style.ViewData {
	hints := ctx.Hints()
	face := ctx.Face(serialFontFamily, tiercatalog.TierNormal)
	hints = textlayout.WithFont(hints, face)

	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	cw := bridge.AvailableContentWidth()
	chh := bridge.AvailableContentHeight()
	if cw <= 0 || chh <= 0 {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	// Raw mode: just the output bytes, no chrome. Same at every panel size.
	if pol.Raw {
		return rawView(snap, ctx, bridge, opts.static)
	}

	entry := ctx.Entry(tiercatalog.TierNormal)
	rowHeight := entry.RowHeight
	glyphAdvance := entry.GlyphAdvance
	glyphHeight := entry.GlyphHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	if glyphAdvance <= 0 {
		glyphAdvance = textlayout.GlyphAdvance
	}
	if glyphHeight <= 0 {
		glyphHeight = rowHeight
	}

	ox, oy := bridge.ContentOrigin()
	maxChars := cw / glyphAdvance

	comp := widgets.NewCompositor(widgets.SuppressionContext{
		IsEink:          opts.static,
		AvailableWidth:  cw,
		AvailableHeight: chh,
	})

	y := oy

	// --- Header row: connection LED + port/baud ---
	ledSize := glyphHeight
	xOffset := 0
	if ledSize >= 3 {
		ledState := led.Off
		ledColor := ledRed
		if snap.Connected {
			ledState = led.On
			ledColor = ledGreen
		}
		comp.Add(led.New(led.Config{
			State:      ledState,
			Brightness: -1.0,
			Diameter:   ledSize,
			Bounds:     image.Rect(ox, y, ox+ledSize, y+ledSize),
			Foreground: ledColor,
		}))
		xOffset = ledSize + 1
	}
	header := compactHeader(snap)
	if maxChars > 0 {
		header = textlayout.Truncate(header, maxCharsAfter(cw, xOffset, glyphAdvance))
	}
	comp.Add(textlabel.New(textlabel.Config{
		Text:       header,
		Bounds:     image.Rect(ox+xOffset, y, ox+cw, y+rowHeight),
		Font:       face,
		Foreground: headerColor(opts.color),
	}))
	y += rowHeight

	// --- State row ---
	stateStr, stateCol := statusText(snap)
	if maxChars > 0 {
		stateStr = textlayout.Truncate(stateStr, maxChars)
	}
	comp.Add(textlabel.New(textlabel.Config{
		Text:       stateStr,
		Bounds:     image.Rect(ox, y, ox+cw, y+rowHeight),
		Font:       face,
		Foreground: stateColor(stateCol, opts.color),
	}))
	y += rowHeight

	// --- Throughput sparkline + buffer gauge (live, tall panels only) ---
	if opts.live && chh >= sparklineMinHeight {
		sparkH := rowHeight
		comp.Add(sparkline.New(sparkline.Config{
			Data:       normalizeThroughput(snap, pol),
			Style:      sparkline.Line,
			Bounds:     image.Rect(ox, y, ox+cw, y+sparkH),
			Foreground: color.RGBA{R: 0, G: 200, B: 200, A: 255},
		}))
		y += sparkH

		barH := 4
		fill := bufferFill(snap, pol)
		comp.Add(progressbar.New(progressbar.Config{
			Style:      progressbar.Linear,
			Value:      fill,
			Bounds:     image.Rect(ox, y, ox+cw, y+barH),
			Foreground: color.RGBA{R: 0, G: 180, B: 0, A: 255},
		}))
		y += barH
	}

	// --- Output lines, ANSI-colored, filling the remainder ---
	lineBudget := 0
	if rowHeight > 0 {
		lineBudget = (oy + chh - y) / rowHeight
	}
	placeDataLines(comp, snap, ox, y, cw, rowHeight, glyphAdvance, face, lineBudget, maxChars, opts.color)

	return style.ViewData{
		Sprites: comp.Sprites(),
		Cursor:  -1,
	}
}

// placeDataLines composites the tail of the output buffer as ANSI-colored text
// labels, one row per line, segment-by-segment where a line carries multiple
// colors. Shows a placeholder when the buffer is empty.
func placeDataLines(comp *widgets.Compositor, snap source.Snapshot, ox, startY, cw, rowHeight, glyphAdvance int, face font.Face, budget, maxChars int, useColor bool) {
	if budget <= 0 {
		return
	}
	if len(snap.Lines) == 0 {
		comp.Add(textlabel.New(textlabel.Config{
			Text:       "(no data yet)",
			Bounds:     image.Rect(ox, startY, ox+cw, startY+rowHeight),
			Font:       face,
			Foreground: colorWhite,
		}))
		return
	}

	lines := snap.Lines
	if len(lines) > budget {
		lines = lines[len(lines)-budget:]
	}
	for i, raw := range lines {
		y := startY + i*rowHeight
		text, segments := source.ParseLine(raw)
		if !useColor || len(segments) <= 1 {
			fg := colorSerialData
			if useColor && len(segments) == 1 && segments[0].Color != nil {
				fg = colorToRGBA(segments[0].Color)
			} else if !useColor {
				fg = colorWhite
			}
			if maxChars > 0 {
				text = textlayout.Truncate(text, maxChars)
			}
			comp.Add(textlabel.New(textlabel.Config{
				Text:       text,
				Bounds:     image.Rect(ox, y, ox+cw, y+rowHeight),
				Font:       face,
				Foreground: fg,
			}))
			continue
		}
		// Multi-segment colored line: place each segment at its own x cursor.
		xCursor := ox
		for _, seg := range segments {
			if seg.Start+seg.Length > len(text) {
				continue
			}
			segText := text[seg.Start : seg.Start+seg.Length]
			segWidth := len(segText) * glyphAdvance
			fg := colorSerialData
			if seg.Color != nil {
				fg = colorToRGBA(seg.Color)
			}
			comp.Add(textlabel.New(textlabel.Config{
				Text:       segText,
				Bounds:     image.Rect(xCursor, y, xCursor+segWidth, y+rowHeight),
				Font:       face,
				Foreground: fg,
			}))
			xCursor += segWidth
		}
	}
}

// rawView renders output bytes only — no LED, header, or status — on the row
// plane, honouring Policy.Raw. The one chrome-free variant.
func rawView(snap source.Snapshot, ctx style.StyleContext, bridge layout.LayoutCalculator, static bool) style.ViewData {
	entry := ctx.Entry(tiercatalog.TierNormal)
	rowHeight := entry.RowHeight
	glyphAdvance := entry.GlyphAdvance
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	cw := bridge.AvailableContentWidth()
	ox, _ := bridge.ContentOrigin()
	maxChars := 0
	if glyphAdvance > 0 {
		maxChars = cw / glyphAdvance
	}
	budget := bridge.MaxVisibleRows(layout.NewMetrics(glyphAdvance, rowHeight, entry.FontID))

	var items []string
	var colors []color.Color
	if len(snap.Lines) == 0 {
		items = []string{"(no data yet)"}
		colors = []color.Color{nil}
	} else {
		lines := snap.Lines
		if budget > 0 && len(lines) > budget {
			lines = lines[len(lines)-budget:]
		}
		for _, raw := range lines {
			text, segments := source.ParseLine(raw)
			if maxChars > 0 {
				text = textlayout.Truncate(text, maxChars)
			}
			items = append(items, text)
			var rowColor color.Color
			if len(segments) > 0 {
				rowColor = segments[0].Color
			}
			colors = append(colors, rowColor)
		}
	}

	offsets := make([]style.TTYX, len(items))
	for i := range offsets {
		offsets[i] = style.TTYX(ox)
	}
	return style.ViewData{
		Items:       items,
		Colors:      colors,
		LineOffsets: offsets,
		Cursor:      -1,
	}
}

// normalizeThroughput scales the snapshot's bytes-per-second history to [0,1]
// against the theoretical max for the configured baud (baud/10 for 8N1), so the
// sparkline reads as "fraction of line capacity in use."
func normalizeThroughput(snap source.Snapshot, pol source.Policy) []float64 {
	maxBPS := pol.Baud / 10
	if maxBPS <= 0 {
		maxBPS = 1
	}
	out := make([]float64, len(snap.Throughput))
	for i, v := range snap.Throughput {
		r := float64(v) / float64(maxBPS)
		if r > 1 {
			r = 1
		}
		if r < 0 {
			r = 0
		}
		out[i] = r
	}
	return out
}

// bufferFill returns the buffer occupancy [0,1] for the capacity gauge.
func bufferFill(snap source.Snapshot, pol source.Policy) float64 {
	if pol.MaxLines <= 0 {
		return 0
	}
	f := float64(len(snap.Lines)) / float64(pol.MaxLines)
	if f > 1 {
		f = 1
	}
	if f < 0 {
		f = 0
	}
	return f
}

// headerColor returns the header text color, accent on color panels and white
// elsewhere.
func headerColor(useColor bool) color.RGBA {
	if useColor {
		return colorSerialHeader
	}
	return colorWhite
}

// stateColor returns the state-row color, the traffic-light tint on color
// panels and white elsewhere.
func stateColor(tint color.RGBA, useColor bool) color.RGBA {
	if useColor {
		return tint
	}
	return colorWhite
}

// maxCharsAfter returns how many glyphs fit in the content width to the right
// of an x inset.
func maxCharsAfter(contentWidth, xOffset, glyphAdvance int) int {
	if glyphAdvance <= 0 {
		return 0
	}
	avail := contentWidth - xOffset
	if avail < 0 {
		return 0
	}
	return avail / glyphAdvance
}
