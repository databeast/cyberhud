package source

import (
	"image"
	"image/color"

	fonts "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets"
)

const splashText = "CYBERHUD"

// splashDuration is the total time the splash is visible (centered on peak).
// The splash fades in for half this duration and fades out for the other half.
const splashPeakThreshold = 0.95 // intensity must be >= this for splash to show

// splashColor is a saturated matrix green for the label.
var splashColor = color.RGBA{0, 255, 50, 255}

// renderSplash produces a sprite with the "CYBERHUD" text centered on the panel,
// with the given alpha (0-255) for fade in/out.
// The face parameter is the catalog-validated font resolved through hints.Face.
func RenderSplash(panelWidth, panelHeight int, alpha uint8, face fonts.Face) *widgets.Sprite {
	if face == nil {
		return nil
	}

	m := face.Metrics()
	textWidth := len(splashText) * m.GlyphAdvance
	textHeight := m.RowHeight

	// If the text doesn't fit the panel width, skip splash rendering.
	if textWidth > panelWidth {
		return nil
	}

	// Center the text on the panel.
	x := (panelWidth - textWidth) / 2
	y := (panelHeight - textHeight) / 2

	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))

	// Use alpha channel for proper fade via DrawGlyph's alpha blending.
	fg := color.RGBA{
		R: splashColor.R,
		G: splashColor.G,
		B: splashColor.B,
		A: alpha,
	}

	for i, ch := range splashText {
		fonts.DrawGlyph(img, face, ch, x+i*m.GlyphAdvance, y, fg)
	}

	return &widgets.Sprite{
		Image:    img,
		Position: image.Point{},
		Label:    "splash",
	}
}

// splashAlpha computes the fade alpha for the splash based on cycle intensity.
// Returns 0 if intensity is below the threshold, otherwise fades in/out
// smoothly within the peak region.
func SplashAlpha(intensity float64) uint8 {
	if intensity < splashPeakThreshold {
		return 0
	}
	// Map [threshold, 1.0] to [0, 1] for fade calculation.
	// Fade in from threshold to 1.0, symmetrically fades out on the way down.
	t := (intensity - splashPeakThreshold) / (1.0 - splashPeakThreshold)
	// t goes 0→1 as we approach peak, then 1→0 as we leave.
	// Use it directly as alpha fraction.
	alpha := int(t * 255)
	if alpha > 255 {
		alpha = 255
	}
	if alpha < 0 {
		alpha = 0
	}
	return uint8(alpha)
}
