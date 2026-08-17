package styles

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/thermal/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/sparkline"
)

// buildPortraitSparkStyle renders sparklines rotated 90° so the data axis runs along
// the tall dimension of portrait side panels.
// It is the shared BuildFn used by per-resolution styles that want the spark layout.
func buildPortraitSparkStyle(snapshot source.ThermalSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	if len(snapshot.Zones) == 0 {
		return style.ViewData{Items: []string{"no thermal data"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 5})

	ox, oy := bridge.ContentOrigin()
	width := bridge.AvailableContentWidth()
	height := bridge.AvailableContentHeight()

	// Zone truncation: each sparkline needs at least 8px of width after rotation.
	maxZones := width / 8
	if maxZones < 1 {
		maxZones = 1
	}

	zones := snapshot.Zones
	if len(zones) > maxZones {
		zones = zones[:maxZones]
	}

	// Equal width allocation per zone (after rotation).
	perZoneWidth := width / len(zones)
	if perZoneWidth < 1 {
		perZoneWidth = 1
	}

	vd := style.ViewData{
		Items:  []string{"thermal"},
		Colors: []color.Color{color.RGBA{R: 255, G: 255, B: 255, A: 255}},
	}

	for i, z := range zones {
		ec := effectiveCritical(z, float64(pol.CritThreshold))
		sev := severity(z.TempC, float64(pol.WarnThreshold), ec)
		sevColor := severityColorRGBA(sev)

		// Get zone's 64-sample history and normalize.
		history := source.GetHistory(z.ZoneID)
		sparkData := normalizeSparkline(history, ec)

		// Render the sparkline in landscape orientation:
		// width = height (data axis, long dimension)
		// height = perZoneWidth (amplitude axis, will become the horizontal width after rotation)
		sparkResult := sparkline.Render(sparkline.Config{
			Data:       sparkData,
			Style:      sparkline.Line,
			Bounds:     image.Rect(0, 0, height, perZoneWidth),
			Foreground: sevColor,
		})
		if sparkResult == nil {
			continue
		}

		// Apply heat gradient to the sparkline fill area.
		applyHeatGradient(sparkResult.Image, perZoneWidth)

		// Rotate 90° CW: the data axis (horizontal) becomes vertical (top-to-bottom).
		rotated := rotateCW90(sparkResult.Image)

		// Position: side-by-side horizontally, anchored to content origin.
		xPos := ox + i*perZoneWidth

		vd.Sprites = append(vd.Sprites, widgets.Sprite{
			Image:    rotated,
			Position: image.Pt(xPos, oy),
			Label:    fmt.Sprintf("portrait-spark-zone-%d", i),
		})
	}

	return vd
}

// applyHeatGradient replaces filled (non-background) pixels in a sparkline image
// with a vertical heat gradient. Pixels near the bottom of the fill area (low amplitude)
// are blue/green, transitioning through yellow to red at the top (high amplitude).
// The height parameter is the total amplitude height of the sparkline.
func applyHeatGradient(img image.Image, height int) {
	rgba, ok := img.(*image.RGBA)
	if !ok || height < 1 {
		return
	}

	bounds := rgba.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Background color is the sparkline's bg (black by default).
	bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	for y := 0; y < h; y++ {
		// t = 0.0 at bottom (low temp), 1.0 at top (high temp).
		// In image coords, y=0 is top, y=h-1 is bottom.
		t := 1.0 - float64(y)/float64(h-1)
		if h == 1 {
			t = 1.0
		}
		heatColor := heatGradientColor(t)

		for x := 0; x < w; x++ {
			off := (y-bounds.Min.Y)*rgba.Stride + (x-bounds.Min.X)*4
			px := color.RGBA{
				R: rgba.Pix[off],
				G: rgba.Pix[off+1],
				B: rgba.Pix[off+2],
				A: rgba.Pix[off+3],
			}
			// Replace any non-background pixel with the heat color.
			if px != bg {
				rgba.Pix[off] = heatColor.R
				rgba.Pix[off+1] = heatColor.G
				rgba.Pix[off+2] = heatColor.B
				rgba.Pix[off+3] = heatColor.A
			}
		}
	}
}

// heatGradientColor maps a normalized value t ∈ [0,1] to a heat color.
// 0.0 = cool (dark blue), 0.33 = green, 0.66 = yellow, 1.0 = red.
func heatGradientColor(t float64) color.RGBA {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	var r, g, b uint8
	switch {
	case t < 0.33:
		// Dark blue → green
		f := t / 0.33
		r = 0
		g = uint8(255 * f)
		b = uint8(80 * (1 - f))
	case t < 0.66:
		// Green → yellow
		f := (t - 0.33) / 0.33
		r = uint8(255 * f)
		g = 255
		b = 0
	default:
		// Yellow → red
		f := (t - 0.66) / 0.34
		r = 255
		g = uint8(255 * (1 - f))
		b = 0
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// rotateCW90 rotates an image 90° clockwise.
// A source image of size W×H becomes H×W.
func rotateCW90(src image.Image) image.Image {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, srcH, srcW))

	// Fast path for *image.RGBA sources (common case from sparkline.Render).
	if srcRGBA, ok := src.(*image.RGBA); ok {
		srcPix := srcRGBA.Pix
		srcStride := srcRGBA.Stride
		dstPix := dst.Pix
		dstStride := dst.Stride

		// 90° CW: dst(sy, srcW-1-sx) = src(sx, sy)
		for sy := 0; sy < srcH; sy++ {
			for sx := 0; sx < srcW; sx++ {
				dx := sy
				dy := srcW - 1 - sx
				srcOff := (sy-srcBounds.Min.Y)*srcStride + (sx-srcBounds.Min.X)*4
				dstOff := dy*dstStride + dx*4
				dstPix[dstOff] = srcPix[srcOff]
				dstPix[dstOff+1] = srcPix[srcOff+1]
				dstPix[dstOff+2] = srcPix[srcOff+2]
				dstPix[dstOff+3] = srcPix[srcOff+3]
			}
		}
		return dst
	}

	// Slow path for other image types.
	for sy := srcBounds.Min.Y; sy < srcBounds.Max.Y; sy++ {
		for sx := srcBounds.Min.X; sx < srcBounds.Max.X; sx++ {
			dx := sy - srcBounds.Min.Y
			dy := srcW - 1 - (sx - srcBounds.Min.X)
			dst.Set(dx, dy, src.At(sx, sy))
		}
	}
	return dst
}
