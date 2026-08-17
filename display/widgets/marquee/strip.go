package marquee

import (
	"image"
	"image/color"
	"time"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets"
)

// Compile-time interface checks.
var (
	_ widgets.Renderable = (*Strip)(nil)
	_ widgets.Animated   = (*Strip)(nil)
)

// Strip is the marquee widget instance. It satisfies widgets.Renderable and widgets.Animated.
// It produces a continuously scrolling column (or row) of characters that a compositor
// can place into a full-frame image.
type Strip struct {
	cfg    Config
	offset float64 // Current scroll offset in cell units (advanced by Tick)
}

// New creates a marquee strip with the given configuration.
// It stores the config, initializes the internal scroll offset to cfg.Phase,
// and falls back to fonts.Default() if Font is nil.
func New(cfg Config) *Strip {
	face := cfg.Font
	if face == nil {
		face = font.Default()
	}
	return &Strip{
		cfg: Config{
			Bounds:    cfg.Bounds,
			Direction: cfg.Direction,
			Font:      face,
			Source:    cfg.Source,
			Colors:    cfg.Colors,
			Speed:     cfg.Speed,
			Phase:     cfg.Phase,
		},
		offset: cfg.Phase,
	}
}

// Source returns the CharSource configured for this strip.
func (s *Strip) Source() CharSource {
	return s.cfg.Source
}

// Speed returns the configured scroll speed in cells per second.
func (s *Strip) Speed() float64 {
	return s.cfg.Speed
}

// Offset returns the current scroll offset in cell units.
func (s *Strip) Offset() float64 {
	return s.offset
}

// Tick advances the scroll position by elapsed time.
// The offset grows by Speed * elapsed.Seconds(), making the strip move.
// For vertical strips, the offset wraps when the entire trail has exited
// the visible area, creating a continuous rain cycle.
func (s *Strip) Tick(elapsed time.Duration) {
	s.offset += s.cfg.Speed * elapsed.Seconds()

	// Wrap for vertical strips: once lead + full trail has exited the bottom,
	// reset to start a new drop from the top.
	if s.cfg.Direction == Vertical {
		metrics := s.cfg.Font.Metrics()
		if metrics.RowHeight > 0 {
			numVisible := s.cfg.Bounds.Dy() / metrics.RowHeight
			trailLen := len(s.cfg.Colors)
			wrapPoint := float64(numVisible + trailLen)
			if s.offset >= wrapPoint {
				s.offset -= wrapPoint
			}
		}
	}
}

// RenderFrame produces a Sprite for the current strip state.
// Returns nil if Source is nil or Bounds has zero/negative dimensions.
//
// Smooth scrolling: the fractional part of the scroll offset is converted to
// a pixel shift, so characters glide between cell positions rather than jumping.
// An extra cell is rendered at the leading edge to cover the partial entry.
func (s *Strip) RenderFrame() *widgets.Sprite {
	cfg := s.cfg
	if cfg.Source == nil {
		return nil
	}

	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()
	if w <= 0 || h <= 0 {
		return nil
	}

	metrics := cfg.Font.Metrics()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	numCells := s.visibleCells(metrics)
	frontIndex := int(s.offset)

	// Sub-pixel offset: fractional part of offset converted to pixels.
	// This is how far the leading cell has scrolled past its grid position.
	frac := s.offset - float64(frontIndex)
	var pixelShift int
	if cfg.Direction == Horizontal {
		pixelShift = int(frac * float64(metrics.GlyphAdvance))
	} else {
		pixelShift = int(frac * float64(metrics.RowHeight))
	}

	// Render one extra cell to cover the partial character entering the visible area.
	cellsToRender := numCells + 1
	if len(cfg.Colors) > 0 && cellsToRender > len(cfg.Colors) {
		cellsToRender = len(cfg.Colors)
	}

	for i := 0; i < cellsToRender; i++ {
		var ch rune
		if cfg.Direction == Horizontal {
			// Horizontal: text scrolls left (enters from right, exits left).
			ch = cfg.Source.CharAt(frontIndex + i)
		} else {
			// Vertical (rain): lead (i=0) is the lowest visible cell,
			// trail extends upward from it.
			ch = cfg.Source.CharAt(frontIndex - i)
		}
		fg := s.colorForCell(i)
		var x, y int

		if cfg.Direction == Horizontal {
			x, y = s.cellPosition(i, metrics)
			x -= pixelShift
		} else {
			// Vertical: position each cell relative to the current frontIndex.
			// Lead (i=0) at row frontIndex, trail (i=1..) at rows above.
			// Clamp to panel bounds — cells above y=0 are simply not drawn.
			row := frontIndex - i
			y = row*metrics.RowHeight + pixelShift
			x = 0

			// Skip cells that are above the visible area.
			if y < -metrics.RowHeight {
				continue
			}
			// Skip cells that are below the visible area.
			if y >= h {
				continue
			}
		}

		font.DrawGlyph(img, cfg.Font, ch, x, y, fg)
	}

	return &widgets.Sprite{
		Image:    img,
		Position: cfg.Bounds.Min,
		Label:    "marquee",
	}
}

// visibleCells calculates how many cells fit in the visible area based on direction.
func (s *Strip) visibleCells(metrics font.Metrics) int {
	if s.cfg.Direction == Horizontal {
		if metrics.GlyphAdvance <= 0 {
			return 0
		}
		return s.cfg.Bounds.Dx() / metrics.GlyphAdvance
	}
	// Vertical
	if metrics.RowHeight <= 0 {
		return 0
	}
	return s.cfg.Bounds.Dy() / metrics.RowHeight
}

// colorForCell returns the color for visible cell at index i.
// If Colors is empty, defaults to opaque white.
func (s *Strip) colorForCell(i int) color.RGBA {
	if len(s.cfg.Colors) == 0 {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	if i < len(s.cfg.Colors) {
		return s.cfg.Colors[i]
	}
	return color.RGBA{R: 255, G: 255, B: 255, A: 255}
}

// cellPosition computes the pixel position for visible cell at index i.
func (s *Strip) cellPosition(i int, metrics font.Metrics) (x, y int) {
	if s.cfg.Direction == Horizontal {
		return i * metrics.GlyphAdvance, 0
	}
	// Vertical: cell 0 at the top, increasing i moves down.
	return 0, i * metrics.RowHeight
}
