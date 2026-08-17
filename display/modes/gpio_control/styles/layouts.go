package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/gpio_control/source"
	"github.com/databeast/cyberhud/display/style"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

// BuildColors returns a Colors slice with one entry per Items row.
// Output-mode pins get ColorOutput, input-mode pins get ColorInput.
// If pins is empty, returns a single-entry Colors slice for the placeholder.
func BuildColors(pins []gpiomgr.PinState) []color.Color {
	if len(pins) == 0 {
		return []color.Color{ColorInput}
	}
	colors := make([]color.Color, len(pins))
	for i, p := range pins {
		if p.Mode == gpiomgr.ModeOutput {
			colors[i] = ColorOutput
		} else {
			colors[i] = ColorInput
		}
	}
	return colors
}

// BuildSprites returns one Sprite per output pin, positioned at the left
// edge of the corresponding row. Uses the Compositor pattern with led.New()
// to produce state indicators: led.On for HIGH, led.Off for LOW.
// Input pins have no Sprite. If pins is empty, returns nil.
func BuildSprites(pins []gpiomgr.PinState, hints textlayout.TextHints) []widgets.Sprite {
	if len(pins) == 0 {
		return nil
	}

	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}

	// LED diameter matches row height for alignment with text.
	diameter := led.DiameterForRow(rowHeight)

	ctx := widgets.SuppressionContext{
		AvailableWidth:  diameter,
		AvailableHeight: rowHeight * len(pins),
	}
	comp := widgets.NewCompositor(ctx)

	for i, p := range pins {
		state := led.Off
		if p.Level {
			state = led.On
		}

		comp.AddIf(p.Mode == gpiomgr.ModeOutput, led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(0, i*rowHeight)},
			Foreground: ColorOutput,
		}, widgets.WithLabel(fmt.Sprintf("pin%d-state", p.Number))))
	}
	if len(comp.Sprites()) == 0 {
		return nil
	}
	return comp.Sprites()
}

// BuildGridView renders pins as an LED grid with borderframe cursor highlight.
// Each cell is 2*diameter wide and 2*diameter tall. The cursor cell gets a
// borderframe highlight. Scroll offset is applied when pins exceed visible cells.
// Uses the Compositor pattern for widget orchestration.
func BuildGridView(pins []gpiomgr.PinState, hints textlayout.TextHints, cursor int) []widgets.Sprite {
	if len(pins) == 0 {
		return nil
	}

	diameter := hints.GlyphHeight
	if diameter <= 0 {
		diameter = textlayout.GlyphHeight
	}

	cellSize := 2 * diameter
	if cellSize <= 0 {
		return nil
	}

	columns := hints.PixelWidth / cellSize
	if columns < 1 {
		columns = 1
	}

	// Compute visible rows and scroll offset.
	visibleRows := hints.PixelHeight / cellSize
	if visibleRows < 1 {
		visibleRows = 1
	}
	visibleCells := columns * visibleRows

	// Compute scroll offset to keep cursor in view.
	scrollOffset := 0
	if len(pins) > visibleCells && cursor >= visibleCells {
		// Scroll so that the cursor row is the last visible row.
		cursorRow := cursor / columns
		topRow := cursorRow - visibleRows + 1
		if topRow < 0 {
			topRow = 0
		}
		scrollOffset = topRow * columns
	}

	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)

	// Render visible pins.
	endIndex := scrollOffset + visibleCells
	if endIndex > len(pins) {
		endIndex = len(pins)
	}

	for i := scrollOffset; i < endIndex; i++ {
		p := pins[i]
		localIndex := i - scrollOffset
		col := localIndex % columns
		row := localIndex / columns

		x := col * cellSize
		y := row * cellSize

		state := led.Off
		if p.Level {
			state = led.On
		}

		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(x, y)},
			Foreground: ColorOutput,
		}, widgets.WithLabel(fmt.Sprintf("pin%d-grid", p.Number))))
	}

	// Render borderframe cursor highlight.
	if cursor >= scrollOffset && cursor < endIndex {
		localIndex := cursor - scrollOffset
		cursorCol := localIndex % columns
		cursorRow := localIndex / columns

		cellBounds := image.Rect(
			cursorCol*cellSize,
			cursorRow*cellSize,
			cursorCol*cellSize+cellSize,
			cursorRow*cellSize+cellSize,
		)

		comp.Add(borderframe.New(borderframe.Config{Bounds: cellBounds}))
	}

	return comp.Sprites()
}

// buildControlTextLabelSprites renders each pin row as a TextLabel sprite for list style.
// Returns nil if any textlabel widget returns nil (triggering fallback to Items).
// Returns an empty non-nil slice if there are no rows to render.
// Uses the Compositor pattern for widget orchestration.
func buildControlTextLabelSprites(pins []gpiomgr.PinState, items []string, hints textlayout.TextHints, face font.Face, colors []color.Color, maxRows int) []widgets.Sprite {
	if len(pins) == 0 {
		return []widgets.Sprite{}
	}

	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	pixelWidth := hints.PixelWidth

	// Use the truncated/limited items list for rendering.
	visibleItems := items
	visiblePins := pins
	if maxRows > 0 && len(visiblePins) > maxRows {
		visiblePins = visiblePins[:maxRows]
	}
	if len(visibleItems) > len(visiblePins) {
		visibleItems = visibleItems[:len(visiblePins)]
	}

	ctx := widgets.SuppressionContext{
		AvailableWidth:  pixelWidth,
		AvailableHeight: rowHeight * len(visiblePins),
	}
	comp := widgets.NewCompositor(ctx)

	for i, p := range visiblePins {
		rowText := ""
		if i < len(visibleItems) {
			rowText = visibleItems[i]
		}

		// Determine foreground color for this row.
		var fg color.RGBA
		if i < len(colors) && colors[i] != nil {
			if rgba, ok := colors[i].(color.RGBA); ok {
				fg = rgba
			} else {
				r, g, b, a := colors[i].RGBA()
				fg = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
			}
		}

		comp.Add(textlabel.New(textlabel.Config{
			Text:       rowText,
			Bounds:     image.Rect(0, i*rowHeight, pixelWidth, (i+1)*rowHeight),
			Font:       face,
			Alignment:  textlabel.Left,
			Foreground: fg,
		}, widgets.WithLabel(fmt.Sprintf("textlabel-pin%d", p.Number))))
	}

	// All-or-nothing semantic: if any row failed to render (compositor discarded it),
	// abandon the entire batch and return nil for fallback.
	if len(comp.Sprites()) != len(visiblePins) {
		return nil
	}
	return comp.Sprites()
}

// sharedListBuild is the default layout for mono and grayscale capability styles
// whose BuildFn is nil. Produces text-based rows with cursor navigation.
// Output is always Static=true.
func sharedListBuild(data source.Data, _ source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()
	bridge := ctx.Layout(0)

	pins := data.Pins
	items := source.BuildItems(pins)
	colors := BuildColors(pins)

	if len(pins) == 0 {
		return style.ViewData{Items: items, Colors: colors, Static: true}
	}

	sprites := BuildSprites(pins, hints)

	// Add cursor indicator to the selected row.
	for i := range items {
		if i == data.Cursor {
			items[i] = ">" + items[i]
		} else {
			items[i] = " " + items[i]
		}
	}

	// Truncate text to available width.
	maxCols := 0
	if bridge.GlyphAdvance() > 0 {
		maxCols = bridge.AvailableContentWidth() / bridge.GlyphAdvance()
	}
	if maxCols > 0 {
		for i, item := range items {
			if len(item) > maxCols {
				items[i] = item[:maxCols]
			}
		}
	}

	// Limit visible rows with scroll offset.
	maxRows := bridge.MaxVisibleRows()
	if maxRows > 0 && len(items) > maxRows {
		topRow := data.TopRow
		if topRow < 0 {
			topRow = 0
		}
		if topRow+maxRows > len(items) {
			topRow = len(items) - maxRows
		}
		items = items[topRow : topRow+maxRows]
		if len(colors) > topRow+maxRows {
			colors = colors[topRow : topRow+maxRows]
		} else if topRow < len(colors) {
			colors = colors[topRow:]
		}
	}

	return style.ViewData{
		Items:   items,
		Colors:  colors,
		Sprites: sprites,
		Cursor:  data.Cursor,
		TopRow:  data.TopRow,
		Static:  true,
	}
}

// sharedGridBuild is the default layout for color capability styles (≥64×64)
// whose BuildFn is nil. Produces LED grid sprites with cursor highlight.
// Output is always Static=true.
func sharedGridBuild(data source.Data, _ source.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()

	pins := data.Pins
	items := source.BuildItems(pins)
	colors := BuildColors(pins)

	if len(pins) == 0 {
		return style.ViewData{Items: items, Colors: colors, Static: true}
	}

	gridSprites := BuildGridView(pins, hints, data.Cursor)

	return style.ViewData{
		Items:   items,
		Colors:  colors,
		Sprites: gridSprites,
		Static:  true,
	}
}

// Colors for output vs input pins. They differ in at least one RGBA channel.
var (
	ColorOutput = color.RGBA{R: 0x00, G: 0xAA, B: 0xFF, A: 0xFF} // Blue for output pins
	ColorInput  = color.RGBA{R: 0xAA, G: 0xAA, B: 0x66, A: 0xFF} // Muted yellow for input pins
)
