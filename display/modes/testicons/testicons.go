package testicons

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/icons"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	catalog.Register(catalog.Definition{
		ID:      "testicons",
		Title:   "Test Icons",
		Summary: "Cycles through all registered icons showing each enlarged with its name.",
		Order:   92,
	})
}

// Cycle interval: each icon is displayed for 8 seconds.
const displayDuration = 1000 // milliseconds

// Rendering colors.
var (
	colorBackground = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	colorLabel      = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
)

// Handler implements action.Handler for the testicons mode (no-op, non-interactive).
type Handler struct{}

// HandleAction returns a zero-value Result for all actions since the testicons
// mode is non-interactive.
func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return action.Result{}
}

// ViewData holds the rendered icon showcase output.
type ViewData struct {
	Title   string
	Hint    string
	Static  bool
	Image   *image.RGBA
	Sprites []widgets.Sprite
}

// frameRenderable wraps a pre-rendered image as a widgets.Renderable so it can
// be passed through the Compositor pattern.
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

// BuildView renders the current icon showcase for the given panel hints and time.
func BuildView(hints textlayout.TextHints, now time.Time) ViewData {
	names := icons.Names()

	if len(names) == 0 {
		img := image.NewRGBA(image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight))
		draw.Draw(img, img.Bounds(), &image.Uniform{colorBackground}, image.Point{}, draw.Src)

		ctx := widgets.SuppressionContext{
			AvailableWidth:  hints.PixelWidth,
			AvailableHeight: hints.PixelHeight,
		}
		comp := widgets.NewCompositor(ctx)
		comp.Add(&frameRenderable{img: img, label: "testicons"})

		return ViewData{Title: "TEST ICONS", Hint: "", Static: false, Image: img, Sprites: comp.Sprites()}
	}

	index := int(now.UnixMilli()/displayDuration) % len(names)
	iconImg, _ := icons.Get(names[index])

	// Determine reference font for label rendering.
	refFace := resolveRefFace(hints)
	refMetrics := refFace.Metrics()
	labelHeight := refMetrics.RowHeight
	if labelHeight < 8 {
		labelHeight = 8
	}

	// Scale icon to fill available height minus label area using nearest-neighbor.
	srcBounds := iconImg.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	availableHeight := hints.PixelHeight - labelHeight
	if availableHeight < 1 {
		availableHeight = 1
	}

	scale := availableHeight / srcHeight
	if scale < 1 {
		scale = 1
	}
	// Also constrain by width to maintain 1:1 aspect ratio.
	widthScale := hints.PixelWidth / srcWidth
	if widthScale < scale {
		scale = widthScale
	}
	if scale < 1 {
		scale = 1
	}

	scaledWidth := srcWidth * scale
	scaledHeight := srcHeight * scale

	// Create full panel-sized image with black background.
	img := image.NewRGBA(image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{colorBackground}, image.Point{}, draw.Src)

	// Draw scaled icon centered horizontally at top.
	offsetX := (hints.PixelWidth - scaledWidth) / 2
	for sy := 0; sy < srcHeight; sy++ {
		for sx := 0; sx < srcWidth; sx++ {
			c := iconImg.At(srcBounds.Min.X+sx, srcBounds.Min.Y+sy)
			r, g, b, a := c.RGBA()
			if a == 0 {
				continue
			}
			pixel := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
			// Fill a scale×scale block in the destination.
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					destX := offsetX + sx*scale + dx
					destY := sy*scale + dy
					if destX >= 0 && destX < hints.PixelWidth && destY >= 0 && destY < hints.PixelHeight {
						img.SetRGBA(destX, destY, pixel)
					}
				}
			}
		}
	}

	// Draw label text below the scaled icon.
	labelY := scaledHeight + 1
	labelText := names[index]
	if refMetrics.GlyphAdvance > 0 {
		maxChars := hints.PixelWidth / refMetrics.GlyphAdvance
		if len(labelText) > maxChars {
			labelText = labelText[:maxChars]
		}
	}

	x := 2
	for _, ch := range labelText {
		renderGlyph(img, refFace, ch, x, labelY, colorLabel)
		x += refMetrics.GlyphAdvance
	}

	// Use Compositor to collect the full-frame sprite.
	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)
	comp.Add(&frameRenderable{img: img, label: "testicons"})

	return ViewData{Title: "TEST ICONS", Hint: "", Static: false, Image: img, Sprites: comp.Sprites()}
}

// RenderCacheKey returns a string that changes when the displayed icon changes.
func RenderCacheKey(hints textlayout.TextHints, now time.Time) string {
	names := icons.Names()
	if len(names) == 0 {
		return "testicons:empty"
	}
	index := int(now.UnixMilli()/displayDuration) % len(names)
	return fmt.Sprintf("testicons:%s", names[index])
}

// resolveRefFace returns a catalog-validated face for use as the label/reference font.
// Uses TierSmall via hints.Face as the default tier face.
func resolveRefFace(hints textlayout.TextHints) font.Face {
	if hints.Catalog.PixelWidth() <= 0 {
		return fallbackFace{}
	}
	face := hints.Face("spleen", tiercatalog.TierSmall)
	if face == nil {
		return fallbackFace{}
	}
	return face
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
type fallbackFace struct{}

func (fallbackFace) ID() string                { return "fallback" }
func (fallbackFace) Metrics() font.Metrics     { return font.Metrics{} }
func (fallbackFace) GlyphRow(rune, int) uint32 { return 0 }
