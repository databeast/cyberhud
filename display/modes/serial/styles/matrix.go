package styles

import (
	"fmt"
	"image"
	"image/color"
	"unicode"

	"github.com/databeast/cyberhud/display/modes/serial/source"
	"github.com/databeast/cyberhud/display/style"
	font "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets"
	"golang.org/x/image/draw"
)

func matrixStyle(width, height int, capability style.Capability) def {
	return def{
		name: fmt.Sprintf("%s-%dx%d", capability, width, height),
		reqs: style.SurfaceRequirements{MinWidth: width, MinHeight: height, Capability: capability},
		p:    Params{BuildFn: buildMatrix},
	}
}

var (
	ColorFast13x9Style     = matrixStyle(13, 9, style.ColorFast)
	GrayscaleFast15x7Style = matrixStyle(15, 7, style.GrayscaleFast)
	GrayscaleFast16x8Style = matrixStyle(16, 8, style.GrayscaleFast)
)

// A matrix has room for only one glyph; show the latest printable character.
func buildMatrix(snap source.Snapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	vd := style.ViewData{Cursor: -1}
	bridge := ctx.Layout(0)
	width, height := bridge.AvailableContentWidth(), bridge.AvailableContentHeight()
	if width <= 0 || height <= 0 {
		return vd
	}

	ch, ok := lastMatrixCharacter(snap.Lines)
	if !ok {
		if pol.Raw {
			return vd
		}
		fg := colorWhite
		if !snap.Connected {
			fg = color.RGBA{R: 128, G: 128, B: 128, A: 255}
			if snap.LastError != "" {
				fg = color.RGBA{R: 64, G: 64, B: 64, A: 255}
			}
		}
		if d.reqs.Capability == style.ColorFast {
			switch {
			case snap.Connected:
				fg = colorSerialOK
			case snap.LastError != "":
				fg = colorSerialErr
			default:
				fg = colorSerialWarn
			}
		}
		pixel := image.NewRGBA(image.Rect(0, 0, 1, 1))
		pixel.SetRGBA(0, 0, fg)
		ox, oy := bridge.ContentOrigin()
		vd.Sprites = []widgets.Sprite{{Image: pixel, Position: image.Pt(ox+width/2, oy+height/2), Label: "serial status"}}
		return vd
	}

	face := font.Default()
	metrics := face.Metrics()
	glyph := image.NewRGBA(image.Rect(0, 0, metrics.GlyphWidth, metrics.GlyphHeight))
	font.DrawGlyph(glyph, face, ch, 0, 0, colorWhite)
	bounds := image.Rectangle{}
	for y := 0; y < metrics.GlyphHeight; y++ {
		for x := 0; x < metrics.GlyphWidth; x++ {
			if glyph.RGBAAt(x, y).A != 0 {
				pixel := image.Rect(x, y, x+1, y+1)
				if bounds.Empty() {
					bounds = pixel
				} else {
					bounds = bounds.Union(pixel)
				}
			}
		}
	}
	if bounds.Empty() {
		return vd
	}
	glyphWidth, glyphHeight := min(bounds.Dx(), width), min(bounds.Dy(), height)
	imageOut := image.NewRGBA(image.Rect(0, 0, glyphWidth, glyphHeight))
	draw.NearestNeighbor.Scale(imageOut, imageOut.Bounds(), glyph, bounds, draw.Src, nil)
	ox, oy := bridge.ContentOrigin()
	vd.Sprites = []widgets.Sprite{{
		Image: imageOut, Position: image.Pt(ox+(width-glyphWidth)/2, oy+(height-glyphHeight)/2),
		Label: "serial character",
	}}
	return vd
}

func lastMatrixCharacter(lines []string) (rune, bool) {
	for i := len(lines) - 1; i >= 0; i-- {
		text, _ := source.ParseLine(lines[i])
		runes := []rune(text)
		for j := len(runes) - 1; j >= 0; j-- {
			ch := runes[j]
			if unicode.IsSpace(ch) || unicode.IsControl(ch) {
				continue
			}
			if ch < 33 || ch > 126 {
				ch = '?'
			}
			return ch, true
		}
	}
	return 0, false
}
