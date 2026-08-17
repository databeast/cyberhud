package testpattern

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/databeast/cyberhud/display/catalog"
	fontpkg "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/runtime/action"
)

// Rendering constants.
const (
	cornerMarkerSize       = 4
	minDimensionForCorners = 8
	swatchCount            = 5
	monoSwatchCount        = 2
	swatchHeightMultiplier = 2
	guideIntensityMax      = 63
)

// Color palette.
var (
	colorBorder     = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF} // white
	colorBackground = color.RGBA{0x00, 0x00, 0x00, 0xFF} // black
	colorSwatchRed  = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	colorSwatchGrn  = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	colorSwatchBlu  = color.RGBA{0x00, 0x00, 0xFF, 0xFF}
	colorSwatchWht  = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	colorSwatchBlk  = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	colorGuide      = color.RGBA{0x3F, 0x3F, 0x3F, 0xFF} // 25% gray
	colorGlyph      = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF} // white
	colorHintText   = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF} // white
)

func init() {
	catalog.Register(catalog.Definition{
		ID:      "testpattern",
		Title:   "Test Pattern",
		Summary: "Fixed diagnostic pattern for screen validation.",
		Order:   90,
	})
}

// Handler implements action.Handler for the test pattern mode (no-op).
type Handler struct{}

// HandleAction returns a zero-value Result for all actions since the test
// pattern mode is non-interactive.
func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return action.Result{}
}

// ViewData holds the rendered test pattern output.
type ViewData struct {
	Title   string
	Hint    string
	Image   *image.RGBA
	Sprites []widgets.Sprite
	Static  bool
}

// RenderCacheKey returns a deterministic string derived from the TextHints fields.
// Same hints produce the same key; different hints produce different keys.
func RenderCacheKey(hints textlayout.TextHints) string {
	return fmt.Sprintf(
		"testpattern:PW=%d,PH=%d,GW=%d,GH=%d,GA=%d,RH=%d",
		hints.PixelWidth, hints.PixelHeight,
		hints.GlyphWidth, hints.GlyphHeight,
		hints.GlyphAdvance, hints.RowHeight,
	)
}

// frameRenderable wraps a pre-rendered image as a widgets.Renderable so it can
// be passed through the Compositor pattern. This replaces manual sprite slice
// construction with the standard Add/Sprites flow.
type frameRenderable struct {
	img   *image.RGBA
	label string
}

func (f *frameRenderable) RenderFrame() *widgets.Sprite {
	if f.img == nil {
		return nil
	}
	return &widgets.Sprite{Image: f.img, Position: image.Point{}, Label: f.label}
}

// BuildView renders the complete test pattern for the given TextHints.
func BuildView(hints textlayout.TextHints, monochrome bool) ViewData {
	if hints.PixelWidth <= 0 || hints.PixelHeight <= 0 {
		return ViewData{Title: "TEST PATTERN", Hint: "", Image: image.NewRGBA(image.Rect(0, 0, 0, 0)), Static: true}
	}

	img := image.NewRGBA(image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight))

	// Fill background black.
	draw.Draw(img, img.Bounds(), &image.Uniform{colorBackground}, image.Point{}, draw.Src)

	// Render sub-elements in order (guide grid behind, then overlay elements).
	drawGuideGrid(img, hints)
	drawBorder(img, hints)
	drawCornerMarkers(img, hints)
	drawColorSwatches(img, hints, monochrome)
	drawGlyphGrid(img, hints)
	drawHintSummary(img, hints)

	// Use Compositor to collect the full-frame sprite. SuppressionContext is
	// constructed from the panel dimensions; no suppression rules are needed
	// for the test pattern since it is a single full-frame diagnostic image.
	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)
	comp.Add(&frameRenderable{img: img, label: "testpattern"})

	return ViewData{Title: "TEST PATTERN", Hint: "", Image: img, Sprites: comp.Sprites(), Static: true}
}

// --- Rendering sub-function stubs ---

func drawBorder(img *image.RGBA, hints textlayout.TextHints) {
	w := hints.PixelWidth
	h := hints.PixelHeight
	if w <= 0 || h <= 0 {
		return
	}

	// Top row (y=0).
	for x := 0; x < w; x++ {
		img.SetRGBA(x, 0, colorBorder)
	}
	// Bottom row (y=h-1).
	for x := 0; x < w; x++ {
		img.SetRGBA(x, h-1, colorBorder)
	}
	// Left column (x=0).
	for y := 0; y < h; y++ {
		img.SetRGBA(0, y, colorBorder)
	}
	// Right column (x=w-1).
	for y := 0; y < h; y++ {
		img.SetRGBA(w-1, y, colorBorder)
	}
}

func drawCornerMarkers(img *image.RGBA, hints textlayout.TextHints) {
	w := hints.PixelWidth
	h := hints.PixelHeight
	if w < minDimensionForCorners || h < minDimensionForCorners {
		return
	}

	// Top-left corner (0,0) to (3,3).
	for y := 0; y < cornerMarkerSize; y++ {
		for x := 0; x < cornerMarkerSize; x++ {
			img.SetRGBA(x, y, colorBorder)
		}
	}
	// Top-right corner (W-4,0) to (W-1,3).
	for y := 0; y < cornerMarkerSize; y++ {
		for x := w - cornerMarkerSize; x < w; x++ {
			img.SetRGBA(x, y, colorBorder)
		}
	}
	// Bottom-left corner (0,H-4) to (3,H-1).
	for y := h - cornerMarkerSize; y < h; y++ {
		for x := 0; x < cornerMarkerSize; x++ {
			img.SetRGBA(x, y, colorBorder)
		}
	}
	// Bottom-right corner (W-4,H-4) to (W-1,H-1).
	for y := h - cornerMarkerSize; y < h; y++ {
		for x := w - cornerMarkerSize; x < w; x++ {
			img.SetRGBA(x, y, colorBorder)
		}
	}
}

func drawColorSwatches(img *image.RGBA, hints textlayout.TextHints, monochrome bool) {
	w := hints.PixelWidth
	h := hints.PixelHeight
	rh := hints.RowHeight

	// Position the swatch band below the top corner marker area.
	startY := 1
	hasCorners := w >= minDimensionForCorners && h >= minDimensionForCorners
	if hasCorners {
		startY = cornerMarkerSize
	}

	// Do not draw over the bottom border row.
	maxY := h - 1
	if startY >= maxY {
		return
	}

	// Compute swatch band height: min(2 * RowHeight, available height).
	bandHeight := swatchHeightMultiplier * rh
	if bandHeight <= 0 {
		bandHeight = maxY - startY
	}
	if bandHeight > maxY-startY {
		bandHeight = maxY - startY
	}
	if bandHeight <= 0 {
		return
	}

	var colors []color.RGBA
	var count int
	if monochrome {
		colors = []color.RGBA{colorSwatchWht, colorSwatchBlk}
		count = monoSwatchCount
	} else {
		colors = []color.RGBA{colorSwatchRed, colorSwatchGrn, colorSwatchBlu, colorSwatchWht, colorSwatchBlk}
		count = swatchCount
	}

	swatchWidth := w / count

	for i, c := range colors {
		x0 := i * swatchWidth
		for y := startY; y < startY+bandHeight; y++ {
			for x := x0; x < x0+swatchWidth && x < w; x++ {
				// Skip border pixels.
				if x == 0 || x == w-1 || y == 0 || y == h-1 {
					continue
				}
				// Skip corner marker regions.
				if hasCorners {
					inTopLeft := x < cornerMarkerSize && y < cornerMarkerSize
					inTopRight := x >= w-cornerMarkerSize && y < cornerMarkerSize
					inBottomLeft := x < cornerMarkerSize && y >= h-cornerMarkerSize
					inBottomRight := x >= w-cornerMarkerSize && y >= h-cornerMarkerSize
					if inTopLeft || inTopRight || inBottomLeft || inBottomRight {
						continue
					}
				}
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func drawGlyphGrid(img *image.RGBA, hints textlayout.TextHints) {
	ga := hints.GlyphAdvance
	gw := hints.GlyphWidth
	gh := hints.GlyphHeight
	rh := hints.RowHeight
	w := hints.PixelWidth
	h := hints.PixelHeight

	if ga <= 0 || w < ga {
		return
	}

	maxChars := (w - 2*cornerMarkerSize) / ga
	if maxChars <= 0 {
		return
	}

	face := fontpkg.Default()

	startX := cornerMarkerSize

	// Compute where the swatch band ends to position the grid below it.
	hasCorners := w >= minDimensionForCorners && h >= minDimensionForCorners
	swatchStartY := 1
	if hasCorners {
		swatchStartY = cornerMarkerSize
	}
	swatchBandH := swatchHeightMultiplier * rh
	if rh <= 0 {
		swatchBandH = 0
	}
	available := h - swatchStartY
	if swatchBandH > available {
		swatchBandH = available
	}
	if swatchBandH < 0 {
		swatchBandH = 0
	}
	gridStartY := swatchStartY + swatchBandH

	if gridStartY >= h {
		return
	}

	// Helper to render one character at pixel position (px, py).
	renderGlyph := func(ch rune, px, py int) {
		for row := 0; row < gh && py+row < h; row++ {
			bits := face.GlyphRow(ch, row)
			for col := 0; col < gw && px+col < w; col++ {
				if bits&(1<<uint(31-col)) != 0 {
					img.SetRGBA(px+col, py+row, colorGlyph)
				}
			}
		}
	}

	// Row 1: uppercase A-Z.
	alphaRow := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	count := maxChars
	if count > len(alphaRow) {
		count = len(alphaRow)
	}
	for ci := 0; ci < count; ci++ {
		px := startX + ci*ga
		renderGlyph(rune(alphaRow[ci]), px, gridStartY)
	}

	// Row 2: digits 0-9 (at gridStartY + RowHeight).
	if rh <= 0 {
		return
	}
	digitRowY := gridStartY + rh
	if digitRowY >= h {
		return
	}
	digitRow := "0123456789"
	count = maxChars
	if count > len(digitRow) {
		count = len(digitRow)
	}
	for ci := 0; ci < count; ci++ {
		px := startX + ci*ga
		renderGlyph(rune(digitRow[ci]), px, digitRowY)
	}
}

func drawHintSummary(img *image.RGBA, hints textlayout.TextHints) {
	ga := hints.GlyphAdvance
	gw := hints.GlyphWidth
	gh := hints.GlyphHeight
	rh := hints.RowHeight
	w := hints.PixelWidth
	h := hints.PixelHeight

	if ga <= 0 || rh <= 0 {
		return
	}

	face := fontpkg.Default()

	startX := cornerMarkerSize
	startY := cornerMarkerSize

	labels := []struct {
		name  string
		value int
	}{
		{"PW", hints.PixelWidth},
		{"PH", hints.PixelHeight},
		{"GW", hints.GlyphWidth},
		{"GH", hints.GlyphHeight},
		{"GA", hints.GlyphAdvance},
		{"RH", hints.RowHeight},
	}

	curY := startY
	for _, lbl := range labels {
		text := fmt.Sprintf("%s=%d", lbl.name, lbl.value)

		// Check vertical fit.
		if curY+gh > h {
			break
		}

		// Check horizontal fit.
		textWidth := len(text) * ga
		if startX+textWidth > w {
			// Try to fit: truncate label to available chars.
			maxChars := (w - startX) / ga
			if maxChars <= 0 {
				break
			}
			if maxChars < len(text) {
				text = text[:maxChars]
			}
		}

		// Render each character.
		for ci, ch := range text {
			px := startX + ci*ga
			for row := 0; row < gh && curY+row < h; row++ {
				bits := face.GlyphRow(ch, row)
				for col := 0; col < gw && px+col < w; col++ {
					if bits&(1<<uint(31-col)) != 0 {
						img.SetRGBA(px+col, curY+row, colorHintText)
					}
				}
			}
		}

		curY += rh
	}
}

func drawGuideGrid(img *image.RGBA, hints textlayout.TextHints) {
	w := hints.PixelWidth
	h := hints.PixelHeight
	rh := hints.RowHeight
	ga := hints.GlyphAdvance

	if rh <= 0 || ga <= 0 {
		return
	}

	// Horizontal lines at y = n * RowHeight.
	numH := h / rh
	for n := 0; n < numH; n++ {
		y := n * rh
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, colorGuide)
		}
	}

	// Vertical lines at x = n * GlyphAdvance.
	numV := w / ga
	for n := 0; n < numV; n++ {
		x := n * ga
		for y := 0; y < h; y++ {
			img.SetRGBA(x, y, colorGuide)
		}
	}
}
