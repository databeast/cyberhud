package styles

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/databeast/cyberhud/display/modes/gpio/source"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

// Hardcoded LED state colors for the 128×128 two-column layout.
// These encode direction×level as distinct colors, independent of policy accent.
var (
	colorOutputHigh = color.RGBA{R: 255, G: 0, B: 0, A: 255}   // Red: OUTPUT + HIGH
	colorOutputLow  = color.RGBA{R: 255, G: 165, B: 0, A: 255} // Orange: OUTPUT + LOW
	colorInputHigh  = color.RGBA{R: 0, G: 255, B: 0, A: 255}   // Green: INPUT + HIGH
	colorInputLow   = color.RGBA{R: 0, G: 0, B: 255, A: 255}   // Blue: INPUT + LOW
)

// ledColorForState maps a pin's mode and level to the LED widget state and foreground color.
// OUTPUT+HIGH → On, Red; OUTPUT+LOW → Off, Orange; INPUT+HIGH → On, Green; INPUT+LOW → Off, Blue.
func ledColorForState(pin gpiomgr.PinState) (led.State, color.RGBA) {
	if pin.Mode == gpiomgr.ModeOutput {
		if pin.Level {
			return led.On, colorOutputHigh
		}
		return led.Off, colorOutputLow
	}
	if pin.Level {
		return led.On, colorInputHigh
	}
	return led.Off, colorInputLow
}

// buildAccentColors returns a color slice using the accent-resolved colors.
// HIGH pins use the resolved accent color; LOW pins use the dimmed accent.
// Returns nil when there are no pins or when colorEnabled is false.
func buildAccentColors(pins []gpiomgr.PinState, pol Policy) []color.Color {
	if !pol.Color {
		return nil
	}
	if len(pins) == 0 {
		return nil
	}
	accent := resolveFGColor(pol.FGColor)
	dimmed := dimFGColor(pol.FGColor)
	colors := make([]color.Color, len(pins))
	for i, p := range pins {
		if p.Level {
			colors[i] = accent
		} else {
			colors[i] = dimmed
		}
	}
	return colors
}

// buildTextLabelSprites renders each pin row as a TextLabel sprite for list style.
// Returns nil if any textlabel widget returns nil (triggering fallback to Items).
// Returns an empty non-nil slice if pins is empty.
// Uses the Compositor pattern for widget orchestration.
func buildTextLabelSprites(pins []gpiomgr.PinState, hints textlayout.TextHints, face font.Face, colors []color.Color) []widgets.Sprite {
	if len(pins) == 0 {
		return []widgets.Sprite{}
	}

	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	pixelWidth := hints.PixelWidth

	items := BuildItemsTruncated(pins, textlayout.MaxCharsPerRow(hints, 0))

	ctx := widgets.SuppressionContext{
		AvailableWidth:  pixelWidth,
		AvailableHeight: rowHeight * len(pins),
	}
	comp := widgets.NewCompositor(ctx)

	for i := range pins {
		rowText := items[i]

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

		w := textlabel.New(textlabel.Config{
			Text:       rowText,
			Bounds:     image.Rect(0, i*rowHeight, pixelWidth, (i+1)*rowHeight),
			Font:       face,
			Alignment:  textlabel.Left,
			Foreground: fg,
		})
		comp.Add(w)
	}

	// All-or-nothing semantic: if any row failed to render (compositor discarded it),
	// abandon the entire batch and return nil for fallback.
	if len(comp.Sprites()) != len(pins) {
		return nil
	}
	return comp.Sprites()
}

// BuildDetailView renders the "detail" style: one LED indicator + text row per pin.
// Each row shows: LED (On for HIGH, Off for LOW) | "{BCM} {mode} {HI|LO} {label}".
// Labels are sourced from the labels map keyed by BCM pin number; if empty or absent,
// the label segment is omitted. Visible rows are truncated to MaxVisibleRows(hints, 0),
// and label text is truncated if it exceeds remaining character width.
//
// Framework pattern demonstrated: Compositor — per-row LED + TextLabel widget
// compositing with truncation and label integration.
func BuildDetailView(pins []gpiomgr.PinState, hints textlayout.TextHints, labels map[int]string, face font.Face) []widgets.Sprite {
	maxRows := textlayout.MaxVisibleRows(hints, 0)
	maxChars := textlayout.MaxCharsPerRow(hints, 0)

	visiblePins := pins
	if maxRows > 0 && len(visiblePins) > maxRows {
		visiblePins = visiblePins[:maxRows]
	}

	rowHeight := hints.RowHeight
	if rowHeight < 1 {
		rowHeight = textlayout.RowHeight
	}
	diameter := led.DiameterForRow(rowHeight)
	pixelWidth := hints.PixelWidth

	ctx := widgets.SuppressionContext{
		AvailableWidth:  pixelWidth,
		AvailableHeight: rowHeight * len(visiblePins),
	}
	comp := widgets.NewCompositor(ctx)

	for i, p := range visiblePins {
		// --- LED indicator ---
		state := led.Off
		fg := ColorLow
		if p.Level {
			state = led.On
			fg = ColorHigh
		}
		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(0, i*rowHeight)},
			Foreground: fg,
		}))

		// --- Text row ---
		// Build the text: "{BCM} {mode} {HI|LO} {label}"
		lvl := "LO"
		if p.Level {
			lvl = "HI"
		}
		modeStr := p.Mode.String()
		bcmStr := fmt.Sprintf("%d", p.Number)

		// Fixed portion: "{BCM} {mode} {HI|LO}"
		fixedText := bcmStr + " " + modeStr + " " + lvl

		// Compute how many chars the LED occupies in character units.
		// The LED occupies `diameter` pixels; the text area starts at X=diameter.
		// So the text area width in pixels is pixelWidth - diameter.
		textPixelWidth := pixelWidth - diameter
		advance := hints.GlyphAdvance
		if advance <= 0 {
			advance = textlayout.GlyphAdvance
		}
		textMaxChars := 0
		if textPixelWidth > 0 && advance > 0 {
			textMaxChars = textPixelWidth / advance
		}

		// Also respect the overall maxChars (computed from full pixelWidth) but
		// subtract the chars consumed by the LED width. We use textMaxChars as
		// the effective character budget for the text portion.
		_ = maxChars // textMaxChars is more accurate for the offset text area

		// Check if there's a label for this pin.
		label := ""
		if labels != nil {
			label = labels[p.Number]
		}

		var rowText string
		if label != "" {
			rowText = fixedText + " " + label
		} else {
			rowText = fixedText
		}

		// Truncate if rowText exceeds textMaxChars.
		if textMaxChars > 0 && len(rowText) > textMaxChars {
			// We want to truncate the label portion, not the fixed portion.
			// If the fixed text alone fits, truncate the label.
			if label != "" && len(fixedText)+1 < textMaxChars {
				// Available chars for label = textMaxChars - len(fixedText) - 1 (for the space)
				availLabel := textMaxChars - len(fixedText) - 1
				if availLabel > 0 && len(label) > availLabel {
					label = label[:availLabel]
				} else if availLabel <= 0 {
					label = ""
				}
				if label != "" {
					rowText = fixedText + " " + label
				} else {
					rowText = fixedText
				}
			} else {
				// Fixed text alone exceeds budget — truncate the whole thing.
				rowText = rowText[:textMaxChars]
			}
		}

		// Render text label at offset (diameter, i*rowHeight) with width = pixelWidth - diameter.
		comp.Add(textlabel.New(textlabel.Config{
			Text:       rowText,
			Bounds:     image.Rect(diameter, i*rowHeight, pixelWidth, (i+1)*rowHeight),
			Font:       face,
			Alignment:  textlabel.Left,
			Foreground: pinColor(p),
		}))
	}

	return comp.Sprites()
}

// pinColor returns ColorHigh for HIGH pins, ColorLow for LOW pins.
func pinColor(p gpiomgr.PinState) color.RGBA {
	if p.Level {
		return ColorHigh
	}
	return ColorLow
}

// BuildDashboardView constructs a ViewData for the "dashboard" style.
// It renders a header with pin counts, then output and input LED grids.
//
// Framework pattern demonstrated: Compositor — multi-section widget compositing
// with header TextLabel and partitioned LED grids.
func BuildDashboardView(pins []gpiomgr.PinState, hints textlayout.TextHints, face font.Face) style.ViewData {
	total := len(pins)
	var outputs, high int
	var outputPins, inputPins []gpiomgr.PinState
	for _, p := range pins {
		if p.Mode == gpiomgr.ModeOutput {
			outputs++
			outputPins = append(outputPins, p)
		} else {
			inputPins = append(inputPins, p)
		}
		if p.Level {
			high++
		}
	}

	headerText := fmt.Sprintf("Pins:%d Out:%d Hi:%d", total, outputs, high)

	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	pixelWidth := hints.PixelWidth
	pixelHeight := hints.PixelHeight

	ctx := widgets.SuppressionContext{
		AvailableWidth:  pixelWidth,
		AvailableHeight: pixelHeight,
	}
	comp := widgets.NewCompositor(ctx)

	// Render the header row as a TextLabel at position (0,0).
	comp.Add(textlabel.New(textlabel.Config{
		Text:      headerText,
		Bounds:    image.Rect(0, 0, pixelWidth, rowHeight),
		Font:      face,
		Alignment: textlabel.Left,
	}))

	// Cell size for LED grids: max(GlyphHeight, 8).
	cellSize := hints.GlyphHeight
	if cellSize < 8 {
		cellSize = 8
	}

	// Columns for LED grid.
	cols := pixelWidth / cellSize
	if cols < 1 {
		cols = 1
	}

	// Render output pins LED grid starting at Y = rowHeight.
	outputStartY := rowHeight
	for i, p := range outputPins {
		row := i / cols
		col := i % cols
		x := col * cellSize
		y := outputStartY + row*cellSize

		state := led.Off
		if p.Level {
			state = led.On
		}
		fg := ColorLow
		if p.Level {
			fg = ColorHigh
		}
		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   cellSize,
			Bounds:     image.Rectangle{Min: image.Pt(x, y)},
			Foreground: fg,
		}))
	}

	// Compute rows used by output grid.
	outputRows := 0
	if len(outputPins) > 0 {
		outputRows = (len(outputPins) + cols - 1) / cols
	}

	// Input grid starts one cell below output grid.
	inputStartY := outputStartY + outputRows*cellSize + cellSize

	// Check if there's remaining height for input grid.
	remainingHeight := pixelHeight - (outputStartY + outputRows*cellSize)
	if remainingHeight >= cellSize && len(inputPins) > 0 {
		for i, p := range inputPins {
			row := i / cols
			col := i % cols
			x := col * cellSize
			y := inputStartY + row*cellSize

			state := led.Off
			if p.Level {
				state = led.On
			}
			fg := ColorLow
			if p.Level {
				fg = ColorHigh
			}
			comp.Add(led.New(led.Config{
				State:      state,
				Brightness: -1.0,
				Diameter:   cellSize,
				Bounds:     image.Rectangle{Min: image.Pt(x, y)},
				Foreground: fg,
			}))
		}
	}

	return style.ViewData{
		Sprites: comp.Sprites(),
		Static:  true,
	}
}

// buildGrayscaleFastStyle is the shared Build implementation for all grayscale-fast GPIO styles.
// It uses an area threshold of 16384 to determine label visibility:
//   - Area < 16384: LED indicators only (no text labels)
//   - Area ≥ 16384: LED indicators with truncated text labels
//
// Capability: GrayscaleFast for all grayscale-fast styles.
func buildGrayscaleFastStyle(snap GpioSnapshot, sctx style.StyleContext, reqs style.SurfaceRequirements) style.ViewData {
	p := snap.Policy

	// 1. Construct own LayoutBridge from hints.
	hints := sctx.Hints()
	_ = hints
	bridge := sctx.Layout(0)

	// 2. Check for zero content area.
	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{}
	}

	// 3. Use tier catalog metrics for layout calculations.
	entry, ok := sctx.FontCatalog().Get(tiercatalog.TierNormal)
	if !ok {
		entry = tiercatalog.Entry{
			GlyphWidth:   bridge.GlyphAdvance(),
			GlyphHeight:  bridge.RowHeight(),
			GlyphAdvance: bridge.GlyphAdvance(),
			RowHeight:    bridge.RowHeight(),
		}
	}

	// 4. Resolve accent colors.
	accent := resolveFGColor(p.FGColor)
	dimmed := dimFGColor(p.FGColor)

	// 5. Determine label visibility via area threshold.
	area := reqs.MinWidth * reqs.MinHeight
	showLabels := area >= 16384

	// 6. Build LED sprites per pin using Compositor.
	// Distribute rows evenly across available height so pins fill the panel,
	// but never go below the font's minimum row height for legibility.
	contentH := bridge.AvailableContentHeight()
	numPins := len(snap.Pins)
	if numPins < 1 {
		numPins = 1
	}
	minRowHeight := entry.GlyphHeight + 2
	if minRowHeight < 8 {
		minRowHeight = 8
	}
	rowHeight := contentH / numPins
	if rowHeight < minRowHeight {
		rowHeight = minRowHeight
	}

	// LED indicators sized to match row height for alignment with text.
	diameter := led.DiameterForRow(rowHeight)
	suppCtx := widgets.SuppressionContext{
		AvailableWidth:  bridge.AvailableContentWidth(),
		AvailableHeight: bridge.AvailableContentHeight(),
	}
	comp := widgets.NewCompositor(suppCtx)

	originX, originY := bridge.ContentOrigin()

	visiblePins := snap.Pins
	maxRows := bridge.AvailableContentHeight() / rowHeight
	// Ensure last sprite bottom ((maxRows-1)*rowHeight + diameter) fits within content area.
	if diameter > rowHeight && maxRows > 0 {
		maxBySprite := (bridge.AvailableContentHeight()-diameter)/rowHeight + 1
		if maxBySprite < 0 {
			maxBySprite = 0
		}
		if maxBySprite < maxRows {
			maxRows = maxBySprite
		}
	}
	if maxRows > 0 && len(visiblePins) > maxRows {
		visiblePins = visiblePins[:maxRows]
	}

	var items []string
	var colors []color.Color

	maxChars := 0
	if bridge.GlyphAdvance() > 0 {
		maxChars = bridge.AvailableContentWidth() / bridge.GlyphAdvance()
	}

	for i, pin := range visiblePins {
		// LED sprite for each pin, vertically centered within the row.
		state := led.Off
		fg := dimmed
		if pin.Level {
			state = led.On
			fg = accent
		}
		ledYOffset := (rowHeight - diameter) / 2
		if ledYOffset < 0 {
			ledYOffset = 0
		}
		ledY := originY + i*rowHeight + ledYOffset
		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(originX, ledY)},
			Foreground: fg,
		}))

		// Build truncated text items if labels are visible.
		if showLabels {
			label := grayscaleFastPinLabel(pin, maxChars)
			items = append(items, label)
			if pin.Level {
				colors = append(colors, accent)
			} else {
				colors = append(colors, dimmed)
			}
		}
	}

	// 7. Compute LineOffsets: position text right of the LED indicator.
	var offsets []int
	if len(items) > 0 {
		ledTextGap := diameter + 4 // LED diameter + 4px gap
		offsets = make([]int, len(items))
		for i := range offsets {
			offsets[i] = originX + ledTextGap
		}
	}

	// Handle color-disabled policy.
	if !p.Color {
		colors = nil
	}

	return style.ViewData{
		Items:       items,
		LineOffsets: offsets,
		Colors:      colors,
		Sprites:     comp.Sprites(),
	}
}

// grayscaleFastPinLabel formats a pin label for grayscale-fast display, truncating to maxChars.
func grayscaleFastPinLabel(pin gpiomgr.PinState, maxChars int) string {
	label := pinRowLabel(pin)
	if maxChars > 0 && len(label) > maxChars {
		label = label[:maxChars]
	}
	return label
}

// grayscaleFastMaxLabelLen returns the length of the longest grayscale-fast pin label string.
func grayscaleFastMaxLabelLen(pins []gpiomgr.PinState) int {
	maxLen := 6 // minimum "## X X"
	for _, p := range pins {
		l := len(fmt.Sprintf("%d %s", p.Number, pinLevelStr(p)))
		if l > maxLen {
			maxLen = l
		}
	}
	return maxLen
}

// Shared GPIO style helpers.
const activityWindowSize = 32

var (
	ColorHigh = sharedcolor.GPIOPalette.Active
	ColorLow  = sharedcolor.GPIOPalette.Inactive
)

func resolveFGColor(accent string) color.RGBA {
	return sharedcolor.ResolveAccent(accent)
}

func dimFGColor(accent string) color.RGBA {
	return sharedcolor.Dim(resolveFGColor(accent))
}

func BuildItemsTruncated(pins []gpiomgr.PinState, maxChars int) []string {
	return source.BuildItemsTruncated(pins, maxChars)
}

// BuildIconGrid returns sprites arranged as a grid of LED indicators for icons style.
func BuildIconGrid(pins []gpiomgr.PinState, pixelWidth, pixelHeight, glyphWidth, glyphHeight int) []widgets.Sprite {
	cellW := glyphWidth
	if cellW < 8 {
		cellW = 8
	}
	cellH := glyphHeight
	if cellH < 8 {
		cellH = 8
	}
	cols := 1
	if pixelWidth > 0 && cellW > 0 {
		cols = pixelWidth / cellW
	}
	if cols < 1 {
		cols = 1
	}
	diameter := cellW
	if cellH < diameter {
		diameter = cellH
	}
	ctx := widgets.SuppressionContext{AvailableWidth: pixelWidth, AvailableHeight: pixelHeight}
	comp := widgets.NewCompositor(ctx)
	for i, p := range pins {
		row, col := i/cols, i%cols
		x, y := col*cellW, row*cellH
		if pixelHeight > 0 && y+cellH > pixelHeight {
			break
		}
		state := led.Off
		fg := ColorLow
		if p.Level {
			state = led.On
			fg = ColorHigh
		}
		comp.Add(led.New(led.Config{State: state, Brightness: -1.0, Diameter: diameter, Bounds: image.Rectangle{Min: image.Pt(x, y)}, Foreground: fg}))
	}
	return comp.Sprites()
}

type activityState struct {
	sync.RWMutex
	history map[int][]float64
	prev    map[int]bool
}

var activity = &activityState{history: make(map[int][]float64), prev: make(map[int]bool)}

func RecordActivity(prev, curr []gpiomgr.PinState) {
	activity.Lock()
	defer activity.Unlock()
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
		if !hasPrev {
			if storedPrev, ok := activity.prev[c.Number]; ok && storedPrev != currLevel {
				value = 1.0
			}
		}
		hist := activity.history[c.Number]
		if len(hist) >= activityWindowSize {
			hist = hist[1:]
		}
		hist = append(hist, value)
		activity.history[c.Number] = hist
		activity.prev[c.Number] = currLevel
	}
}

func BuildActivityView(pins []gpiomgr.PinState, hints textlayout.TextHints, face font.Face) style.ViewData {
	maxRows := textlayout.MaxVisibleRows(hints, 0)
	rowHeight := hints.RowHeight
	if rowHeight <= 0 {
		rowHeight = textlayout.RowHeight
	}
	pixelWidth := hints.PixelWidth
	var outputPins []gpiomgr.PinState
	for _, p := range pins {
		if p.Mode == gpiomgr.ModeOutput {
			outputPins = append(outputPins, p)
		}
	}
	if maxRows > 0 && len(outputPins) > maxRows {
		outputPins = outputPins[:maxRows]
	}
	if len(outputPins) == 0 {
		return style.ViewData{Static: false}
	}
	ctx := widgets.SuppressionContext{AvailableWidth: pixelWidth, AvailableHeight: rowHeight * len(outputPins)}
	comp := widgets.NewCompositor(ctx)
	activity.RLock()
	defer activity.RUnlock()
	for i, p := range outputPins {
		data := activity.history[p.Number]
		if len(data) == 0 {
			data = make([]float64, activityWindowSize)
		} else if len(data) < activityWindowSize {
			padded := make([]float64, activityWindowSize)
			copy(padded[activityWindowSize-len(data):], data)
			data = padded
		} else {
			cp := make([]float64, activityWindowSize)
			copy(cp, data[len(data)-activityWindowSize:])
			data = cp
		}
		bounds := image.Rect(0, i*rowHeight, pixelWidth, (i+1)*rowHeight)
		comp.Add(sparkline.New(sparkline.Config{Data: data, Style: sparkline.Line, Bounds: bounds, Foreground: ColorHigh}))
	}
	_ = face
	return style.ViewData{Sprites: comp.Sprites(), Static: false}
}
