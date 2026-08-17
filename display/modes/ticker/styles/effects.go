package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/ticker/source"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
)

// isColorFast returns true when the panel supports full-color fast refresh.
func isColorFast(hints textlayout.TextHints) bool {
	return hints.Capability == textlayout.CapColorFast
}

// buildBorderSprites creates a borderframe sprite at the panel bounds,
// optionally tinted with the accent color on ColorFast panels.
func buildBorderSprites(hints textlayout.TextHints, effective source.Policy) []widgets.Sprite {
	suppCtx := widgets.SuppressionContext{
		AvailableWidth:  hints.PixelWidth,
		AvailableHeight: hints.PixelHeight,
	}
	comp := widgets.NewCompositor(suppCtx)

	cfg := borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
	}
	if isColorFast(hints) {
		cfg.ColorTint = sharedcolor.ResolveAccent(effective.Accent)
	}

	comp.AddIf(true, borderframe.New(cfg))
	return comp.Sprites()
}

// buildGlowBackground creates a panel-bounds *image.RGBA filled with a linear
// gradient from the dimmed accent color at the top to black at the bottom.
func buildGlowBackground(width, height int, accent color.RGBA) widgets.Sprite {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	dimmed := sharedcolor.Dim(accent)

	for y := 0; y < height; y++ {
		// Linear interpolation: dimmed at y=0, black at y=height-1.
		t := float64(y) / float64(height)
		r := uint8(float64(dimmed.R) * (1.0 - t))
		g := uint8(float64(dimmed.G) * (1.0 - t))
		b := uint8(float64(dimmed.B) * (1.0 - t))

		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return widgets.Sprite{
		Image:    img,
		Position: image.Point{X: 0, Y: 0},
		Bounds:   image.Rect(0, 0, width, height),
		Label:    "color-glow-bg",
	}
}

// buildGlowSprites creates one glow sprite per visible text line.
// Each glow sprite is an *image.RGBA rectangle extending at least 1px
// beyond the text glyph bounds in all directions, filled with the accent
// color at reduced opacity (alpha ~100).
func buildGlowSprites(formatted []source.FormattedLine, hints textlayout.TextHints, offsetY int, accent color.RGBA) []widgets.Sprite {
	if len(formatted) == 0 {
		return nil
	}

	sprites := make([]widgets.Sprite, 0, len(formatted))
	currentY := offsetY

	// Glow color: accent at reduced alpha for soft glow effect.
	const glowAlpha = 100
	glowColor := color.RGBA{R: accent.R, G: accent.G, B: accent.B, A: glowAlpha}

	// Extension: glow extends 2px beyond glyph bounds on all sides.
	const glowExtend = 2

	for i, fl := range formatted {
		if fl.Text == "" {
			// Still advance Y for empty lines.
			face := source.ResolveFace(hints, "spleen", fl.Tier)
			currentY += face.Metrics().RowHeight
			continue
		}

		face := source.ResolveFace(hints, "spleen", fl.Tier)
		metrics := face.Metrics()
		rowHeight := metrics.RowHeight
		textWidth := len([]rune(fl.Text)) * metrics.GlyphAdvance

		// Glow bounds: extend beyond text bounds by glowExtend pixels.
		glowX := 0 - glowExtend
		if glowX < 0 {
			glowX = 0
		}
		glowY := currentY - glowExtend
		if glowY < 0 {
			glowY = 0
		}
		glowW := textWidth + 2*glowExtend
		glowH := rowHeight + 2*glowExtend

		// Clamp to panel bounds.
		if glowX+glowW > hints.PixelWidth {
			glowW = hints.PixelWidth - glowX
		}
		if glowY+glowH > hints.PixelHeight {
			glowH = hints.PixelHeight - glowY
		}

		if glowW > 0 && glowH > 0 {
			img := image.NewRGBA(image.Rect(0, 0, glowW, glowH))
			for py := 0; py < glowH; py++ {
				for px := 0; px < glowW; px++ {
					img.SetRGBA(px, py, glowColor)
				}
			}

			sprites = append(sprites, widgets.Sprite{
				Image:    img,
				Position: image.Point{X: glowX, Y: glowY},
				Label:    fmt.Sprintf("color-glow-line-%d", i),
			})
		}

		currentY += rowHeight
	}

	return sprites
}
