package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/widgets/icons"
)

// wifiIconDst is the default inline WiFi icon size (16×16), used as a lower bound
// when the signal-forward layout sizes its icons.
const wifiIconDst = 16

// scaledWifiIcon produces a size×size RGBA image of the WiFi icon, scaled from
// the 8×8 glyph with nearest-neighbor and tinted with the given color. The size
// and tint are explicit so the signal-forward layout can draw a large hero icon
// for the status screens or a small inline icon beside the SSID. Returns nil if
// the icon is missing from the registry or size is non-positive.
func scaledWifiIcon(tint color.RGBA, size int) *image.RGBA {
	if size <= 0 {
		return nil
	}
	srcImg, ok := icons.Get("wifi")
	if !ok {
		return nil
	}
	srcBounds := srcImg.Bounds()
	sw, sh := srcBounds.Dx(), srcBounds.Dy()
	if sw <= 0 || sh <= 0 {
		return nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	for dy := 0; dy < size; dy++ {
		for dx := 0; dx < size; dx++ {
			sx := srcBounds.Min.X + dx*sw/size
			sy := srcBounds.Min.Y + dy*sh/size
			if _, _, _, a := srcImg.At(sx, sy).RGBA(); a > 0 {
				dst.SetRGBA(dx, dy, tint)
			}
		}
	}
	return dst
}
