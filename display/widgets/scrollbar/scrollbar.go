package scrollbar

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/widgets"
)

// Config holds all parameters needed to render a vertical scrollbar.
type Config struct {
	TotalItems   int             // Total number of items in the list
	VisibleItems int             // Number of items visible at once
	ScrollOffset int             // Index of the first visible item
	Bounds       image.Rectangle // Pixel region for rendering
	Foreground   color.RGBA      // Thumb color (zero value defaults to white)
	Background   color.RGBA      // Track color (zero value defaults to black)
}

// Render produces a vertical scrollbar image based on the given configuration.
// Returns nil if Bounds are too small or TotalItems < 1.
func Render(cfg Config) *widgets.Sprite {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()

	// Too small to render.
	if w < 1 || h < 1 {
		return nil
	}

	// No items to represent.
	if cfg.TotalItems < 1 {
		return nil
	}

	fg, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// When all items are visible, fill the entire track with foreground.
	if cfg.TotalItems <= cfg.VisibleItems {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.SetRGBA(x, y, fg)
			}
		}
		return &widgets.Sprite{
			Image:    img,
			Position: cfg.Bounds.Min,
			Label:    "scrollbar",
		}
	}

	// Clamp offset to valid range.
	offset := cfg.ScrollOffset
	if offset < 0 {
		offset = 0
	}
	if offset >= cfg.TotalItems {
		offset = cfg.TotalItems - 1
	}

	// Thumb height: proportional to visible/total ratio, minimum 1px.
	thumbHeight := (h * cfg.VisibleItems) / cfg.TotalItems
	if thumbHeight < 1 {
		thumbHeight = 1
	}

	// Thumb position: proportional to offset, clamped so thumb stays within bounds.
	thumbTop := (h * offset) / cfg.TotalItems
	if thumbTop > h-thumbHeight {
		thumbTop = h - thumbHeight
	}

	// Render: track with background, thumb with foreground
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if y >= thumbTop && y < thumbTop+thumbHeight {
				img.SetRGBA(x, y, fg)
			} else {
				img.SetRGBA(x, y, bg)
			}
		}
	}

	return &widgets.Sprite{
		Image:    img,
		Position: cfg.Bounds.Min,
		Label:    "scrollbar",
	}
}
