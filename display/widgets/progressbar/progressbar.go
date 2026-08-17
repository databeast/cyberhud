package progressbar

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/widgets"
)

// Render produces a progress bar image based on the given configuration.
// Returns nil if Bounds are too small for meaningful rendering, the Style is
// unrecognized, or Ring/Arc bounds are smaller than 3×3.
func Render(cfg Config) *widgets.Sprite {
	// Validate and clamp all config fields in place.
	validate(&cfg)

	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()

	// Too small to render.
	if w < 1 || h < 1 {
		return nil
	}

	// Pie/Ring/Arc charts need at least 3×3 pixels for a meaningful arc.
	switch cfg.Style {
	case Pie, Ring, Arc:
		if w < 3 || h < 3 {
			return nil
		}
	}

	// Resolve orientation into geometry parameters for downstream renderers.
	geom := resolveOrientation(cfg)

	var img *image.RGBA
	var label string

	switch cfg.Style {
	case Linear:
		img = renderLinear(cfg, geom)
		if cfg.RoundedCaps && geom.MinorAxis >= 2 {
			applyRoundedCaps(img, cfg)
		}
		label = "progressbar/linear"
	case Pie:
		img = renderPie(cfg, geom)
		label = "progressbar/pie"
	case Segmented:
		img = renderSegmented(cfg, geom)
		label = "progressbar/segmented"
	case Ring:
		img = renderRing(cfg, geom)
		label = "progressbar/ring"
	case Arc:
		img = renderArc(cfg, geom)
		label = "progressbar/arc"
	default:
		// Unrecognized Style value → nil
		return nil
	}

	// Post-processing: apply border after shape rendering and caps, before animation.
	applyBorder(img, cfg)

	// Apply animation overlay after border.
	applyAnimation(img, cfg)

	// Threshold markers drawn last (on top of fill, track, border, animation).
	drawMarkers(img, cfg)

	return &widgets.Sprite{
		Image:    img,
		Position: cfg.Bounds.Min,
		Label:    label,
	}
}

// renderHorizontal draws a left-to-right filled bar.
// Columns [0, floor(width*value)) are foreground; remaining columns are background.
func renderHorizontal(bounds image.Rectangle, value float64, fg, bg color.RGBA) *image.RGBA {
	w := bounds.Dx()
	h := bounds.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	fillCols := int(float64(w) * value)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x < fillCols {
				img.SetRGBA(x, y, fg)
			} else {
				img.SetRGBA(x, y, bg)
			}
		}
	}

	return img
}

// renderVertical draws a bottom-to-top filled bar.
// Rows [height - floor(height*value), height) are foreground; remaining rows are background.
func renderVertical(bounds image.Rectangle, value float64, fg, bg color.RGBA) *image.RGBA {
	w := bounds.Dx()
	h := bounds.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	fillRows := int(float64(h) * value)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if y >= h-fillRows {
				img.SetRGBA(x, y, fg)
			} else {
				img.SetRGBA(x, y, bg)
			}
		}
	}

	return img
}
