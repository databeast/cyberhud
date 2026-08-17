package svgrender

import (
	"image"
	"image/color"
	"time"

	"github.com/databeast/cyberhud/display/widgets"
	"github.com/srwiley/rasterx"
)

// Sprite is a type alias for widgets.Sprite.
type Sprite = widgets.Sprite

// Config holds parameters for SVG-based widget rendering.
type Config struct {
	Bounds     image.Rectangle // Pixel region for rendering
	Label      string          // Sprite identification (max 128 chars)
	Color      color.RGBA      // Tint color (zero = no tint)
	SVG        string          // SVG markup for static rendering
	Frames     []Frame         // Frame sequence for animation (takes precedence over SVG)
	FrameIndex int             // Current frame index (set by Animator or caller)
}

// Frame represents a single animation keyframe with SVG content and duration.
type Frame struct {
	SVG      string        // SVG markup for this keyframe
	Duration time.Duration // How long this frame is displayed
}

// Render produces an SVG-rendered sprite from the given Config.
// Returns nil for invalid bounds, empty SVG source, or unparseable SVG content.
func Render(cfg Config) *Sprite {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()

	// Resolution guard: reject bounds below minimum threshold.
	if w < MinBoundsWidth || h < MinBoundsHeight {
		return nil
	}

	// Determine SVG source.
	var svgSource string
	if len(cfg.Frames) > 0 {
		// Frames take precedence over SVG field.
		idx := cfg.FrameIndex
		if idx < 0 {
			idx = 0
		}
		if idx >= len(cfg.Frames) {
			idx = len(cfg.Frames) - 1
		}

		// Try the selected frame first, then iterate forward to find valid SVG.
		svgSource = ""
		for i := 0; i < len(cfg.Frames); i++ {
			candidate := (idx + i) % len(cfg.Frames)
			icon, err := parse(cfg.Frames[candidate].SVG, w, h)
			if err == nil && icon != nil {
				svgSource = cfg.Frames[candidate].SVG
				break
			}
		}
		if svgSource == "" {
			return nil
		}
	} else {
		// Use the SVG field directly.
		if cfg.SVG == "" {
			return nil
		}
		svgSource = cfg.SVG
	}

	// Parse SVG source.
	icon, err := parse(svgSource, w, h)
	if err != nil || icon == nil {
		return nil
	}

	// Create canvas and rasterize.
	canvas := NewCanvas(w, h)
	if canvas == nil {
		return nil
	}

	// Rasterize using the rasterx scanline pipeline.
	scanner := rasterx.NewScannerGV(w, h, canvas.img, canvas.img.Bounds())
	dasher := rasterx.NewDasher(w, h, scanner)
	icon.Draw(dasher, 1.0)

	// Apply tint if Color is non-zero.
	if cfg.Color != (color.RGBA{}) {
		applyTint(canvas.img, cfg.Color)
	}

	// Produce sprite with position and truncated label.
	return canvas.ToResult(cfg.Bounds.Min, cfg.Label)
}
