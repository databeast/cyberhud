package testwidgets

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
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/gradient"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
	"github.com/databeast/cyberhud/display/widgets/scaledtextbox"
	"github.com/databeast/cyberhud/display/widgets/scrollbar"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
	"github.com/databeast/cyberhud/display/widgets/textbox"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
	"github.com/databeast/cyberhud/runtime/action"
)

func init() {
	catalog.Register(catalog.Definition{
		ID:      "testwidgets",
		Title:   "Test Widgets",
		Summary: "Cycles through all registered widget types showing a demo preview of each.",
		Order:   93,
	})
}

// Cycle interval: each widget is displayed for 8 seconds.
const displayDuration = 8000 // milliseconds

// Rendering colors.
var (
	colorBackground = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	colorLabel      = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
)

// Handler implements action.Handler for the testwidgets mode (no-op, non-interactive).
type Handler struct{}

// HandleAction returns a zero-value Result for all actions since the testwidgets
// mode is non-interactive.
func (Handler) HandleAction(act action.Action, cursor, itemCount int) action.Result {
	return action.Result{}
}

// ViewData holds the rendered widget showcase output.
type ViewData struct {
	Title        string
	Hint         string
	Static       bool
	VisibleCount int
	Image        *image.RGBA
	Sprites      []widgets.Sprite
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

// BuildView renders the current widget showcase for the given panel hints and time.
func BuildView(hints textlayout.TextHints, now time.Time) ViewData {
	names := showcaseNames()
	if len(names) == 0 {
		return buildPlaceholderView(hints, "No widgets available")
	}

	index := int(now.UnixMilli()/displayDuration) % len(names)
	labelText, widgetSprite := buildVisibleWidget(names, index, hints, now)
	if widgetSprite == nil || widgetSprite.Image == nil {
		return buildPlaceholderView(hints, "No visible widgets")
	}

	// Determine reference font for label rendering.
	refFace := resolveRefFace(hints)
	refMetrics := refFace.Metrics()
	labelHeight := refMetrics.RowHeight
	if labelHeight < 8 {
		labelHeight = 8
	}

	// Create full panel-sized image with black background.
	img := image.NewRGBA(image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{colorBackground}, image.Point{}, draw.Src)

	// Draw label text at the top (above widget preview area).
	if refMetrics.GlyphAdvance > 0 {
		maxChars := hints.PixelWidth / refMetrics.GlyphAdvance
		if len(labelText) > maxChars {
			labelText = labelText[:maxChars]
		}
	}

	x := 2
	for _, ch := range labelText {
		renderGlyph(img, refFace, ch, x, 1, colorLabel)
		x += refMetrics.GlyphAdvance
	}

	// If the widget produced a sprite, draw it into the panel image below the label.
	widgetY := labelHeight + 1
	srcBounds := widgetSprite.Image.Bounds()
	dstRect := image.Rect(0, widgetY, srcBounds.Dx(), widgetY+srcBounds.Dy())
	// Clip to panel bounds.
	if dstRect.Max.X > hints.PixelWidth {
		dstRect.Max.X = hints.PixelWidth
	}
	if dstRect.Max.Y > hints.PixelHeight {
		dstRect.Max.Y = hints.PixelHeight
	}
	draw.Draw(img, dstRect, widgetSprite.Image, srcBounds.Min, draw.Over)

	// Use Compositor to collect the full-frame sprite.
	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)
	comp.Add(&frameRenderable{img: img, label: "testwidgets"})

	return ViewData{Title: "TEST WIDGETS", Hint: "", Static: false, VisibleCount: 1, Image: img, Sprites: comp.Sprites()}
}

func buildVisibleWidget(names []string, start int, hints textlayout.TextHints, now time.Time) (string, *widgets.Sprite) {
	if len(names) == 0 {
		return "", nil
	}
	for offset := 0; offset < len(names); offset++ {
		name := names[(start+offset)%len(names)]
		sprite := renderShowcaseWidget(name, hints, now)
		if sprite == nil || sprite.Image == nil {
			continue
		}
		return name, sprite
	}
	return "", nil
}

func showcaseNames() []string {
	return []string{
		"borderframe",
		"gradient",
		"led",
		"progressbar",
		"scaledtextbox",
		"scrollbar",
		"sparkline",
		"textbox",
		"textlabel",
	}
}

func renderShowcaseWidget(name string, hints textlayout.TextHints, now time.Time) *widgets.Sprite {
	panelW := hints.PixelWidth
	panelH := hints.PixelHeight
	if panelW <= 0 || panelH <= 0 {
		return nil
	}

	face := resolveRefFace(hints)
	bg := colorBackground
	fg := colorLabel
	phase := float64(now.UnixMilli()%displayDuration) / float64(displayDuration)

	switch name {
	case "borderframe":
		r := borderframe.New(borderframe.Config{
			Bounds:     image.Rect(0, 0, panelW, panelH),
			TileSet:    "border",
			Background: bg,
			ColorTint:  fg,
		})
		return renderWidget(r)
	case "gradient":
		r := gradient.New(gradient.Config{
			Style:  gradient.Radial,
			Bounds: image.Rect(0, 0, panelW, panelH),
			Stops: []gradient.ColorStop{
				{Position: 0.0, Color: color.RGBA{0x00, 0x66, 0x22, 0xFF}},
				{Position: 0.6, Color: color.RGBA{0x00, 0x22, 0x11, 0xFF}},
				{Position: 1.0, Color: color.RGBA{0x00, 0x00, 0x00, 0xFF}},
			},
		})
		return renderWidget(r)
	case "led":
		r := led.New(led.Config{
			Shape:       led.Circle,
			State:       led.On,
			Brightness:  -1.0,
			Diameter:    18,
			Bounds:      image.Rect(0, 0, 18, 18),
			Foreground:  color.RGBA{0x22, 0xFF, 0x66, 0xFF},
			GlowEnabled: true,
			GlowRadius:  4,
			ShineStyle:  led.ShineCrescent,
		})
		return renderWidget(r)
	case "progressbar":
		r := progressbar.New(progressbar.Config{
			Style:        progressbar.Segmented,
			Orientation:  progressbar.OrientHorizontal,
			Value:        phase,
			Bounds:       image.Rect(0, 0, panelW, 14),
			Foreground:   color.RGBA{0x33, 0xCC, 0xFF, 0xFF},
			Background:   color.RGBA{0x00, 0x11, 0x22, 0xFF},
			SegmentCount: 16,
			SegmentGap:   1,
		})
		return renderWidget(r)
	case "scaledtextbox":
		r := scaledtextbox.New(scaledtextbox.Config{
			LogicalSize:   image.Point{X: 64, Y: 16},
			TargetSize:    image.Point{X: panelW, Y: maxInt(16, panelH/3)},
			Position:      image.Point{},
			Text:          "TEST WIDGETS",
			Font:          face,
			Alignment:     textbox.Center,
			VAlign:        textbox.Middle,
			Overflow:      textbox.Truncate,
			Foreground:    fg,
			LineSpacing:   0,
			PadX:          1,
			PadY:          1,
			Border:        true,
			FontOverrides: nil,
			Label:         "scaledtextbox",
		})
		return renderWidget(r)
	case "scrollbar":
		r := scrollbar.New(scrollbar.Config{
			TotalItems:   12,
			VisibleItems: 4,
			ScrollOffset: int(now.UnixMilli()/300) % 12,
			Bounds:       image.Rect(panelW-4, 0, panelW, panelH),
			Foreground:   color.RGBA{0xFF, 0xCC, 0x33, 0xFF},
			Background:   color.RGBA{0x20, 0x20, 0x20, 0xFF},
		})
		return renderWidget(r)
	case "sparkline":
		r := sparkline.New(sparkline.Config{
			Data: []float64{
				0.12, 0.28, 0.48, 0.22, 0.68, 0.94, 0.55, 0.34,
				0.62, 0.15, 0.8, 0.45, 0.9, 0.25, 0.7, 0.52,
			},
			Style:      sparkline.Line,
			Bounds:     image.Rect(0, panelH/2, panelW-4, panelH),
			Foreground: color.RGBA{0xFF, 0x55, 0xAA, 0xFF},
			Background: color.RGBA{0x00, 0x00, 0x11, 0xFF},
		})
		return renderWidget(r)
	case "textbox":
		r := textbox.New(textbox.Config{
			Bounds:     image.Rect(0, 0, panelW, maxInt(18, panelH/2)),
			Text:       "TEST WIDGETS",
			Font:       face,
			Alignment:  textbox.Center,
			VAlign:     textbox.Middle,
			Overflow:   textbox.Truncate,
			Foreground: fg,
			Border:     true,
			Label:      "textbox",
		})
		return renderWidget(r)
	case "textlabel":
		r := textlabel.New(textlabel.Config{
			Text:       "WIDGET LABEL",
			Bounds:     image.Rect(0, 0, panelW, maxInt(12, panelH/4)),
			Font:       face,
			Alignment:  textlabel.Center,
			Foreground: fg,
		})
		return renderWidget(r)
	default:
		return nil
	}
}

func renderWidget(r widgets.Renderable) *widgets.Sprite {
	if r == nil {
		return nil
	}
	return r.RenderFrame()
}

func buildPlaceholderView(hints textlayout.TextHints, labelText string) ViewData {
	img := image.NewRGBA(image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{colorBackground}, image.Point{}, draw.Src)

	refFace := resolveRefFace(hints)
	refMetrics := refFace.Metrics()
	if refMetrics.GlyphAdvance > 0 {
		x := 2
		y := 2
		for _, ch := range labelText {
			renderGlyph(img, refFace, ch, x, y, colorLabel)
			x += refMetrics.GlyphAdvance
		}
	}

	ctx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(ctx)
	comp.Add(&frameRenderable{img: img, label: "testwidgets"})

	return ViewData{Title: "TEST WIDGETS", Hint: "", Static: false, VisibleCount: 1, Image: img, Sprites: comp.Sprites()}
}

// RenderCacheKey returns a string that changes when the displayed widget changes.
func RenderCacheKey(hints textlayout.TextHints, now time.Time) string {
	names := showcaseNames()
	if len(names) == 0 {
		return "testwidgets:empty"
	}
	index := int(now.UnixMilli()/displayDuration) % len(names)
	return fmt.Sprintf("testwidgets:%s", names[index])
}

// resolveRefFace returns a catalog-validated face for use as the label/reference font.
// Uses TierSmall via hints.Face as the default tier face, falling back to the
// global default font when the catalog is not yet populated or the tier lookup fails.
func resolveRefFace(hints textlayout.TextHints) font.Face {
	if hints.Catalog.PixelWidth() > 0 {
		if face := hints.Face("spleen", tiercatalog.TierSmall); face != nil {
			return face
		}
	}
	return font.Default()
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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
