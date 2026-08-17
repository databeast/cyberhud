package testfonts

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sort"
	"strings"
	"time"

	"github.com/databeast/cyberhud/display/catalog"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/surface/tierselect"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	catalog.Register(catalog.Definition{
		ID:      "testfonts",
		Title:   "Test Fonts",
		Summary: "Cycles through all registered fonts showing sample text.",
		Order:   91,
	})
}

// Sample text constants.
const (
	sampleUpper     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sampleLower     = "abcdefghijklmnopqrstuvwxyz"
	sampleDigits    = "0123456789"
	sampleSymbols   = "!@#$%^&*()-+=[]{}|;:',.<>/?"
	colorDuration   = 2 * time.Second           // time per color phase
	numColors       = 4                         // white, blue, green, red
	displayDuration = colorDuration * numColors // total time per font (8s)
)

// Rendering colors.
var (
	colorBackground = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	colorHeader     = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}

	// Sample text cycles through these colors every 2 seconds.
	sampleColors = [numColors]color.RGBA{
		{0xFF, 0xFF, 0xFF, 0xFF}, // white
		{0x00, 0x80, 0xFF, 0xFF}, // blue
		{0x00, 0xFF, 0x00, 0xFF}, // green
		{0xFF, 0x00, 0x00, 0xFF}, // red
	}
)

// Handler implements action.Handler for the testfonts mode (no-op, non-interactive).
type Handler struct{}

// HandleAction returns a zero-value Result for all actions since the testfonts
// mode is non-interactive.
func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return action.Result{}
}

// ViewData holds the rendered font showcase output.
type ViewData struct {
	Title     string
	Hint      string
	Static    bool
	Image     *image.RGBA
	Sprites   []widgets.Sprite
	FontIndex int
}

// pageLayout describes the vertical structure of a single font showcase page.
type pageLayout struct {
	headerY      int
	headerHeight int
	gapHeight    int
	bodyY        int
	bodyRows     int
}

// tierRow represents a catalog-resolved tier and the sample text drawn for it.
type tierRow struct {
	tier   tiercatalog.Tier
	face   font.Face
	label  string
	sample string
}

// sortedFaces returns catalog-validated faces sorted by ascending GlyphHeight,
// with ties broken by ascending font ID for deterministic ordering.
// Uses hints.AllFaces() to enumerate only catalog-validated faces.
func sortedFaces(hints textlayout.TextHints) []font.Face {
	if hints.Catalog.PixelWidth() > 0 {
		faces := hints.AllFaces("spleen")
		if len(faces) > 0 {
			sort.Slice(faces, func(i, j int) bool {
				hi := faces[i].Metrics().GlyphHeight
				hj := faces[j].Metrics().GlyphHeight
				if hi != hj {
					return hi < hj
				}
				return faces[i].ID() < faces[j].ID()
			})
			return faces
		}
	}

	all := font.List()
	if len(all) == 0 {
		return nil
	}
	faces := make([]font.Face, 0, len(all))
	for _, face := range all {
		if face == nil || face.Metrics().GlyphHeight <= 0 {
			continue
		}
		faces = append(faces, face)
	}
	if len(faces) == 0 {
		return nil
	}
	sort.Slice(faces, func(i, j int) bool {
		hi := faces[i].Metrics().GlyphHeight
		hj := faces[j].Metrics().GlyphHeight
		if hi != hj {
			return hi < hj
		}
		return faces[i].ID() < faces[j].ID()
	})
	return faces
}

func tierRows(hints textlayout.TextHints) []tierRow {
	catalog, ok := resolvedCatalog(hints)
	if ok {
		hints.Catalog = catalog
		rows := make([]tierRow, 0, len(catalog.Tiers()))
		for _, tier := range catalog.Tiers() {
			face := hints.Face("spleen", tier)
			if face == nil {
				face, _ = tierselect.TrySelect(catalog, tierselect.Request{Family: "spleen", Tier: tier})
			}
			if face == nil {
				continue
			}
			rows = append(rows, tierRow{
				tier:   tier,
				face:   face,
				label:  strings.ToLower(string(tier)),
				sample: sampleForTier(face),
			})
		}
		if len(rows) > 0 {
			return rows
		}
	}

	faces := sortedFaces(hints)
	if len(faces) > 0 {
		rows := make([]tierRow, 0, len(faces))
		for i, face := range faces {
			if face == nil || face.Metrics().GlyphAdvance <= 0 {
				continue
			}
			label := "small"
			switch {
			case i == 0:
				label = "small"
			case i == 1:
				label = "normal"
			case i == 2:
				label = "large"
			default:
				label = fmt.Sprintf("tier%d", i+1)
			}
			rows = append(rows, tierRow{
				tier:   tiercatalog.TierSmall,
				face:   face,
				label:  label,
				sample: sampleForTier(face),
			})
		}
		if len(rows) > 0 {
			return rows
		}
	}

	defaultFace := font.Default()
	if defaultFace != nil {
		return []tierRow{{
			tier:   tiercatalog.TierSmall,
			face:   defaultFace,
			label:  "small",
			sample: sampleForTier(defaultFace),
		}}
	}
	return nil
}

func resolvedCatalog(hints textlayout.TextHints) (tiercatalog.Catalog, bool) {
	if hints.Catalog.PixelWidth() > 0 && hints.Catalog.PixelHeight() > 0 {
		return hints.Catalog, true
	}
	if hints.PixelWidth <= 0 || hints.PixelHeight <= 0 {
		return tiercatalog.Catalog{}, false
	}
	cat, err := tiercatalog.Build(tiercatalog.Params{
		PixelWidth:  hints.PixelWidth,
		PixelHeight: hints.PixelHeight,
		MinChars:    10,
		PPI:         hints.PPI,
	})
	if err != nil {
		return tiercatalog.Catalog{}, false
	}
	return cat, true
}

func sampleForTier(face font.Face) string {
	if face == nil {
		return "Aa"
	}
	m := face.Metrics()
	if m.GlyphAdvance <= 0 || m.GlyphWidth <= 0 {
		return "Aa"
	}
	if m.GlyphHeight >= 20 {
		return "ABC"
	}
	if m.GlyphHeight >= 12 {
		return "abc"
	}
	return "Aa"
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

// BuildView renders the current font showcase page for the given panel hints and time.
func BuildView(hints textlayout.TextHints, now time.Time) ViewData {
	rows := tierRows(hints)
	if len(rows) == 0 {
		img := image.NewRGBA(image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight))
		draw.Draw(img, img.Bounds(), &image.Uniform{colorBackground}, image.Point{}, draw.Src)

		ctx := widgets.SuppressionContext{
			AvailableWidth:  hints.PixelWidth,
			AvailableHeight: hints.PixelHeight,
		}
		comp := widgets.NewCompositor(ctx)
		comp.Add(&frameRenderable{img: img, label: "testfonts"})

		return ViewData{Title: "TEST FONTS", Hint: "", Static: false, Image: img, Sprites: comp.Sprites(), FontIndex: 0}
	}

	index := int(now.UnixMilli()/displayDuration.Milliseconds()) % len(rows)
	row := rows[index]
	colorIndex := int(now.UnixMilli()/colorDuration.Milliseconds()) % numColors
	sampleColor := sampleColors[colorIndex]

	img := image.NewRGBA(image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{colorBackground}, image.Point{}, draw.Src)

	layout := computeLayout(hints, row.face)
	renderHeader(img, hints, row.face, layout)
	renderTierRows(img, hints, rows, row, layout, sampleColor)

	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)
	comp.Add(&frameRenderable{img: img, label: "testfonts"})

	return ViewData{Title: "TEST FONTS", Hint: "", Static: false, Image: img, Sprites: comp.Sprites(), FontIndex: index}
}

// RenderCacheKey returns a string that changes when the displayed font or color changes.
func RenderCacheKey(hints textlayout.TextHints, now time.Time) string {
	rows := tierRows(hints)
	if len(rows) == 0 {
		return "testfonts:empty"
	}
	index := int(now.UnixMilli()/displayDuration.Milliseconds()) % len(rows)
	colorIndex := int(now.UnixMilli()/colorDuration.Milliseconds()) % numColors
	return fmt.Sprintf("testfonts:%s:%d:%d", rows[index].face.ID(), index, colorIndex)
}

// computeLayout determines vertical positioning of header and body.
func computeLayout(hints textlayout.TextHints, showcaseFace font.Face) pageLayout {
	refFace := resolveRefFace(hints)
	refMetrics := refFace.Metrics()
	showcaseMetrics := showcaseFace.Metrics()

	headerY := 2
	headerHeight := refMetrics.RowHeight
	if headerHeight <= 0 {
		headerHeight = 8
	}
	gapHeight := 2
	if hints.PixelHeight <= 64 {
		gapHeight = 1
	}
	bodyY := headerY + headerHeight + gapHeight

	availableHeight := hints.PixelHeight - bodyY
	bodyRows := 0
	if showcaseMetrics.RowHeight > 0 && availableHeight > 0 {
		bodyRows = availableHeight / showcaseMetrics.RowHeight
	}
	if hints.PixelHeight <= 64 && bodyRows > 3 {
		bodyRows = 3
	}
	if bodyRows < 1 && availableHeight > 0 && showcaseMetrics.RowHeight > 0 {
		bodyRows = 1
	}

	return pageLayout{
		headerY:      headerY,
		headerHeight: headerHeight,
		gapHeight:    gapHeight,
		bodyY:        bodyY,
		bodyRows:     bodyRows,
	}
}

// resolveRefFace returns a catalog-validated face for use as the header/reference font.
// Uses TierSmall via hints.Face as the default tier face; if the catalog is unavailable,
// it falls back to the built-in default font so the diagnostic remains visible.
func resolveRefFace(hints textlayout.TextHints) font.Face {
	catalog, ok := resolvedCatalog(hints)
	if ok {
		hints.Catalog = catalog
		face := hints.Face("spleen", tiercatalog.TierSmall)
		if face != nil && face.Metrics().GlyphAdvance > 0 {
			return face
		}
	}
	if defaultFace := font.Default(); defaultFace != nil && defaultFace.Metrics().GlyphAdvance > 0 {
		return defaultFace
	}
	return fallbackFace{}
}

// renderHeader draws the font ID label at the top using the reference font.
func renderHeader(img *image.RGBA, hints textlayout.TextHints, face font.Face, layout pageLayout) {
	refFace := resolveRefFace(hints)
	refMetrics := refFace.Metrics()

	label := face.ID()
	if refMetrics.GlyphAdvance > 0 {
		maxChars := hints.PixelWidth / refMetrics.GlyphAdvance
		if len(label) > maxChars {
			label = label[:maxChars]
		}
	}

	x := 2
	for _, ch := range label {
		renderGlyph(img, refFace, ch, x, layout.headerY, colorHeader)
		x += refMetrics.GlyphAdvance
	}
}

// renderSampleBody draws rows of sample characters in the showcased font.
func renderSampleBody(img *image.RGBA, hints textlayout.TextHints, face font.Face, layout pageLayout, fg color.RGBA) {
	m := face.Metrics()
	if m.GlyphAdvance <= 0 || m.RowHeight <= 0 {
		return
	}

	maxChars := hints.PixelWidth / m.GlyphAdvance
	samples := []string{sampleUpper, sampleLower, sampleDigits, sampleSymbols}

	y := layout.bodyY
	for row := 0; row < layout.bodyRows; row++ {
		sampleIdx := row % len(samples)
		line := samples[sampleIdx]

		if len(line) > maxChars {
			line = line[:maxChars]
		}

		x := 0
		for _, ch := range line {
			renderGlyph(img, face, ch, x, y, fg)
			x += m.GlyphAdvance
		}
		y += m.RowHeight
	}
}

func renderTierRows(img *image.RGBA, hints textlayout.TextHints, rows []tierRow, active tierRow, layout pageLayout, fg color.RGBA) {
	if len(rows) == 0 {
		return
	}
	refFace := resolveRefFace(hints)
	refMetrics := refFace.Metrics()
	if refMetrics.GlyphAdvance <= 0 {
		refMetrics.GlyphAdvance = 6
	}
	if refMetrics.RowHeight <= 0 {
		refMetrics.RowHeight = 8
	}
	maxRows := (hints.PixelHeight - layout.bodyY) / maxInt(refMetrics.RowHeight, 8)
	if maxRows < 1 {
		maxRows = 1
	}
	visible := rows
	if len(visible) > maxRows {
		visible = visible[:maxRows]
	}
	for i, row := range visible {
		face := row.face
		if face == nil {
			continue
		}
		m := face.Metrics()
		rowHeight := maxInt(refMetrics.RowHeight, m.RowHeight+1)
		y := layout.bodyY + i*rowHeight
		if y >= hints.PixelHeight {
			break
		}
		labelText := strings.TrimSpace(row.label)
		if len(labelText) > 5 {
			labelText = labelText[:5]
		}
		x := 1
		labelColor := fg
		if row.tier == active.tier {
			labelColor = colorHeader
		}
		for _, ch := range labelText {
			renderGlyph(img, refFace, ch, x, y, labelColor)
			x += refMetrics.GlyphAdvance
		}
		sampleText := row.sample
		if len(sampleText) > 0 {
			maxChars := maxInt(1, (hints.PixelWidth-x-2)/maxInt(m.GlyphAdvance, 1))
			if len(sampleText) > maxChars {
				sampleText = sampleText[:maxChars]
			}
		}
		if sampleText == "" {
			sampleText = "Aa"
		}
		cx := x + 2
		for _, ch := range sampleText {
			colorUse := fg
			if row.tier == active.tier {
				colorUse = colorHeader
			}
			renderGlyph(img, face, ch, cx, y, colorUse)
			cx += m.GlyphAdvance
		}
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// renderGlyph draws a single character at (px, py) using the given font face.
func renderGlyph(img *image.RGBA, face font.Face, ch rune, px, py int, fg color.RGBA) {
	m := face.Metrics()
	bounds := img.Bounds()
	for row := 0; row < m.GlyphHeight; row++ {
		bits := face.GlyphRow(ch, row)
		for col := 0; col < m.GlyphWidth; col++ {
			if bits&(1<<uint(31-col)) != 0 {
				x, y := px+col, py+row
				if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
					img.SetRGBA(x, y, fg)
				}
			}
		}
	}
}

// fallbackFace is a zero-metrics font.Face used when the catalog is unavailable.
// This prevents nil panics in edge cases where BuildView is called without a valid catalog.
type fallbackFace struct{}

func (fallbackFace) ID() string                { return "fallback" }
func (fallbackFace) Metrics() font.Metrics     { return font.Metrics{} }
func (fallbackFace) GlyphRow(rune, int) uint32 { return 0 }
