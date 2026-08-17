package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/widgets/icons"
)

// Signal bar sprite dimensions.
const (
	barCount   = 4
	barWidth   = 4
	barSpacing = 1
	// Total width: 4 bars × 4px + 3 gaps × 1px = 19px.
	signalBarsWidth  = barCount*barWidth + (barCount-1)*barSpacing
	signalBarsHeight = 16
)

// barHeights defines the height of each bar (1-indexed: bar 1 is shortest, bar 4 is tallest).
var barHeights = [barCount]int{4, 8, 12, 16}

// darkGray is the fill color for inactive (unfilled) signal bars.
var darkGray = color.RGBA{R: 40, G: 40, B: 40, A: 255}

// renderSignalBars produces a 19×16 pixel RGBA image with 4 vertical bars.
// Bars with 1-based index ≤ barLevel are filled with qualityColor;
// bars with index > barLevel are filled with dark gray (40, 40, 40).
// barLevel 0 means all bars are dark gray.
func renderSignalBars(barLevel int, qualityColor color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, signalBarsWidth, signalBarsHeight))

	for i := 0; i < barCount; i++ {
		// Determine fill color for this bar (1-based index).
		var fill color.RGBA
		if i+1 <= barLevel {
			fill = qualityColor
		} else {
			fill = darkGray
		}

		// X origin for this bar: each bar is barWidth wide with barSpacing gap.
		x0 := i * (barWidth + barSpacing)
		x1 := x0 + barWidth

		// Bar is bottom-aligned within signalBarsHeight, with height barHeights[i].
		h := barHeights[i]
		y0 := signalBarsHeight - h
		y1 := signalBarsHeight

		// Fill the bar rectangle.
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				img.SetRGBA(x, y, fill)
			}
		}
	}

	return img
}

// WiFi icon sprite constants.
const (
	wifiIconSrc = 8  // source icon size (8×8)
	wifiIconDst = 16 // scaled icon size (16×16)
)

// Tint colors for WiFi icon based on connection state.
var (
	disconnectedTint = color.RGBA{R: 255, G: 60, B: 60, A: 255}
	unavailableTint  = color.RGBA{R: 80, G: 80, B: 80, A: 255}
)

// renderWifiIcon produces a 16×16 pixel RGBA image of the WiFi icon,
// scaled from the 8×8 icon glyph using nearest-neighbor interpolation.
// The icon is tinted based on connection state:
//   - Connected → accentColor
//   - Disconnected → red (255, 60, 60)
//   - Unavailable → dimmed gray (80, 80, 80)
//
// Returns nil if the "wifi" icon is not found in the icons registry.
func renderWifiIcon(connState source.ConnectionState, accentColor color.RGBA) *image.RGBA {
	srcImg, ok := icons.Get("wifi")
	if !ok {
		return nil
	}

	// Verify source icon bounds are within expected size.
	srcBounds := srcImg.Bounds()
	if srcBounds.Dx() > wifiIconSrc || srcBounds.Dy() > wifiIconSrc {
		return nil
	}

	// Determine tint color based on connection state.
	var tint color.RGBA
	switch connState {
	case source.Connected:
		tint = accentColor
	case source.Disconnected:
		tint = disconnectedTint
	default: // Unavailable
		tint = unavailableTint
	}

	// Create the scaled output image (16×16).
	dst := image.NewRGBA(image.Rect(0, 0, wifiIconDst, wifiIconDst))

	// Scale using nearest-neighbor: each source pixel maps to a 2×2 block.
	scaleX := wifiIconDst / srcBounds.Dx()
	scaleY := wifiIconDst / srcBounds.Dy()

	for dy := 0; dy < wifiIconDst; dy++ {
		for dx := 0; dx < wifiIconDst; dx++ {
			// Map destination pixel to source pixel.
			sx := srcBounds.Min.X + dx/scaleX
			sy := srcBounds.Min.Y + dy/scaleY

			// Read the source pixel alpha.
			_, _, _, a := srcImg.At(sx, sy).RGBA()
			if a > 0 {
				// Opaque source pixel → apply tint color.
				dst.SetRGBA(dx, dy, tint)
			}
			// Transparent source pixels remain zero (transparent black).
		}
	}

	return dst
}

// RenderSignalBars exposes the signal bar renderer for package-level tests.
func RenderSignalBars(barLevel int, qualityColor color.RGBA) *image.RGBA {
	return renderSignalBars(barLevel, qualityColor)
}
