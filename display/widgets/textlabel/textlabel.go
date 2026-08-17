package textlabel

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets"
)

// Alignment specifies horizontal text positioning within bounds.
type Alignment int

const (
	Left Alignment = iota // Default
	Center
	Right
)

// Config holds all parameters for text label rendering.
type Config struct {
	Text       string          // Text content to render
	Bounds     image.Rectangle // Pixel region for rendering
	Font       font.Face       // Bitmap font face for glyph data (nil → fonts.Default())
	Alignment  Alignment       // Horizontal alignment (default: Left)
	Foreground color.RGBA      // Text color (zero value → opaque white)
}

// Render produces a text label bitmap. Returns nil for invalid bounds.
func Render(cfg Config) *widgets.Sprite {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()

	if w < 1 || h < 1 {
		return nil
	}

	// Default font fallback.
	face := cfg.Font
	if face == nil {
		face = font.Default()
	}

	// Default foreground to opaque white when zero-value color provided.
	fg, _ := widgets.ResolveColors(cfg.Foreground, color.RGBA{})

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Empty text → transparent image with correct dimensions.
	if len(cfg.Text) == 0 {
		return &widgets.Sprite{
			Image:    img,
			Position: cfg.Bounds.Min,
			Label:    "textlabel",
		}
	}

	metrics := face.Metrics()
	glyphWidth := metrics.GlyphWidth
	glyphHeight := metrics.GlyphHeight
	glyphAdvance := metrics.GlyphAdvance

	// Compute the total pixel width of the rendered text.
	// Each character takes GlyphAdvance pixels except the last which takes GlyphWidth.
	charCount := len(cfg.Text)
	textPixelWidth := charCount * glyphAdvance
	if charCount > 0 {
		textPixelWidth = (charCount-1)*glyphAdvance + glyphWidth
	}

	// Vertical centering: top offset = floor((height - glyphHeight) / 2).
	topOffset := (h - glyphHeight) / 2

	// Horizontal alignment positioning.
	var xStart int
	switch cfg.Alignment {
	case Left:
		xStart = 0
	case Center:
		xStart = (w - textPixelWidth) / 2
		if xStart < 0 {
			xStart = 0
		}
	case Right:
		xStart = w - textPixelWidth
		if xStart < 0 {
			xStart = 0
		}
	}

	// Render each character.
	xCursor := xStart
	for _, ch := range cfg.Text {
		for row := 0; row < glyphHeight; row++ {
			bits := face.GlyphRow(ch, row)
			for col := 0; col < glyphWidth; col++ {
				x := xCursor + col
				// Right-edge pixel clipping.
				if x >= w {
					break
				}
				if x < 0 {
					continue
				}
				// Bit 31-col being set means pixel at that column is "on".
				if bits&(1<<uint(31-col)) != 0 {
					y := topOffset + row
					if y >= 0 && y < h {
						img.SetRGBA(x, y, fg)
					}
				}
			}
		}
		xCursor += glyphAdvance
		// If cursor is already beyond bounds, stop rendering more characters.
		if xCursor >= w {
			break
		}
	}

	return &widgets.Sprite{
		Image:    img,
		Position: cfg.Bounds.Min,
		Label:    "textlabel",
	}
}
