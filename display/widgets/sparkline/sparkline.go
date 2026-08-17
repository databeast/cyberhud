package sparkline

import (
	"image"
	"image/color"
	"math"

	"github.com/databeast/cyberhud/display/widgets"
)

// Style represents the visual variant of a sparkline.
type Style int

const (
	Line Style = iota
	Bar
)

// Config holds all parameters for sparkline rendering.
type Config struct {
	Data       []float64       // Normalized data points [0.0, 1.0]
	Style      Style           // Line or Bar
	Bounds     image.Rectangle // Pixel region for rendering
	Foreground color.RGBA      // Line/bar color (zero value → opaque white)
	Background color.RGBA      // Fill color (zero value → opaque black)
}

func init() {
	widgets.Register("sparkline", func() widgets.Described {
		return &sparklineWidget{}
	})
}

// sparklineWidget is the persistent instance satisfying Renderable, Described, and Configurable.
type sparklineWidget struct {
	cfg   Config
	opts  widgets.OptionSet
	cache widgets.RenderCache[Config, widgets.Sprite]
}

// New creates a sparkline Widget_Instance satisfying Renderable, Described, and Configurable.
// It accepts optional functional options such as WithCaching and WithLabel.
func New(cfg Config, opts ...widgets.Option) widgets.Renderable {
	optSet := widgets.ApplyOptions(opts...)
	w := &sparklineWidget{
		cfg:  cfg,
		opts: optSet,
	}
	if optSet.CachingEnabled {
		w.cache = widgets.NewRenderCache[Config, widgets.Sprite](w.renderSprite, sign)
	}
	return w
}

// RenderFrame implements the widgets.Renderable interface.
func (w *sparklineWidget) RenderFrame() *widgets.Sprite {
	if w.cache != nil {
		return w.cache.Render(w.cfg)
	}
	return w.renderSprite(w.cfg)
}

// renderSprite performs the actual rendering, producing a *widgets.Sprite.
func (w *sparklineWidget) renderSprite(cfg Config) *widgets.Sprite {
	sprite := Render(cfg)
	if sprite == nil {
		return nil
	}
	if w.opts.LabelOverride != "" {
		sprite.Label = w.opts.LabelOverride
	}
	return sprite
}

// Describe implements the widgets.Described interface.
func (w *sparklineWidget) Describe() widgets.Descriptor {
	return widgets.Descriptor{
		Name:         "sparkline",
		MinWidth:     1,
		MinHeight:    1,
		Capabilities: []string{},
	}
}

// Configure implements the widgets.Configurable interface.
func (w *sparklineWidget) Configure(cfg interface{}) {
	if c, ok := cfg.(Config); ok {
		w.cfg = c
	}
}

// Render produces a sparkline graph bitmap based on the given configuration.
// Returns nil if Bounds are too small for meaningful rendering.
func Render(cfg Config) *widgets.Sprite {
	w := cfg.Bounds.Dx()
	h := cfg.Bounds.Dy()

	if w < 1 || h < 1 {
		return nil
	}

	fg, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)

	// Clamp and truncate data
	data := clampData(cfg.Data)
	if len(data) > w {
		data = data[len(data)-w:]
	}

	var img *image.RGBA
	var label string

	switch cfg.Style {
	case Bar:
		img = renderBar(w, h, data, fg, bg)
		label = "sparkline/bar"
	default:
		img = renderLine(w, h, data, fg, bg)
		label = "sparkline/line"
	}

	return &widgets.Sprite{
		Image:    img,
		Position: cfg.Bounds.Min,
		Label:    label,
	}
}

// clampData creates a new slice with each value clamped to [0.0, 1.0].
// NaN values are treated as 0.0.
func clampData(data []float64) []float64 {
	result := make([]float64, len(data))
	for i, v := range data {
		if math.IsNaN(v) || v < 0.0 {
			result[i] = 0.0
		} else if v > 1.0 {
			result[i] = 1.0
		} else {
			result[i] = v
		}
	}
	return result
}

// renderBar draws vertical bars packed left-to-right.
// Each bar is floor(w / count) pixels wide (minimum 1).
// Bar height = floor(h * value) pixels from the bottom.
func renderBar(w, h int, data []float64, fg, bg color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Empty data: fill entire image with background
	if len(data) == 0 {
		fillImage(img, w, h, bg)
		return img
	}

	count := len(data)
	barWidth := w / count
	if barWidth < 1 {
		barWidth = 1
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Determine which bar this column belongs to
			barIndex := x / barWidth
			if barIndex >= count {
				// Columns past the last bar get background
				img.SetRGBA(x, y, bg)
				continue
			}

			barHeight := int(float64(h) * data[barIndex])
			if y >= h-barHeight {
				img.SetRGBA(x, y, fg)
			} else {
				img.SetRGBA(x, y, bg)
			}
		}
	}

	return img
}

// renderLine draws a connected line graph with filled area below.
// For count=1: single column at x=0, fill from bottom.
// For count>1: linear interpolation between data points.
func renderLine(w, h int, data []float64, fg, bg color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Empty data: fill entire image with background
	if len(data) == 0 {
		fillImage(img, w, h, bg)
		return img
	}

	count := len(data)

	if count == 1 {
		// Single data point: single column at x=0
		fillHeight := int(float64(h) * data[0])
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if x == 0 && y >= h-fillHeight {
					img.SetRGBA(x, y, fg)
				} else {
					img.SetRGBA(x, y, bg)
				}
			}
		}
		return img
	}

	// Multiple data points: compute x position for each point
	// Point i maps to x_i = floor(i * (w-1) / (count-1))
	xPositions := make([]int, count)
	for i := 0; i < count; i++ {
		xPositions[i] = i * (w - 1) / (count - 1)
	}

	// For each pixel column, compute interpolated value
	for x := 0; x < w; x++ {
		// Find which segment this x falls into
		interpValue := interpolateAtX(x, xPositions, data, count)
		fillHeight := int(float64(h) * interpValue)

		for y := 0; y < h; y++ {
			if y >= h-fillHeight {
				img.SetRGBA(x, y, fg)
			} else {
				img.SetRGBA(x, y, bg)
			}
		}
	}

	return img
}

// interpolateAtX computes the linearly interpolated data value at pixel column x.
func interpolateAtX(x int, xPositions []int, data []float64, count int) float64 {
	// Find the segment: the two data points that bracket this x
	// Start from the last segment and work backwards for efficiency isn't needed;
	// simple linear scan is fine for small counts.
	for i := 0; i < count-1; i++ {
		x0 := xPositions[i]
		x1 := xPositions[i+1]
		if x >= x0 && x <= x1 {
			if x0 == x1 {
				// Points map to the same column
				return data[i]
			}
			// Linear interpolation
			t := float64(x-x0) / float64(x1-x0)
			return data[i] + t*(data[i+1]-data[i])
		}
	}
	// Should not reach here for valid x in [0, w-1], but fallback to last point
	return data[count-1]
}

// fillImage fills the entire image with the given color.
func fillImage(img *image.RGBA, w, h int, c color.RGBA) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
}
