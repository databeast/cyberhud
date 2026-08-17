package gradient

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/widgets"
)

// Compile-time interface assertions.
var (
	_ widgets.Renderable   = (*gradientWidget)(nil)
	_ widgets.Described    = (*gradientWidget)(nil)
	_ widgets.Configurable = (*gradientWidget)(nil)
)

// gradientWidget is the internal struct implementing Renderable, Described,
// and Configurable interfaces for the gradient widget.
type gradientWidget struct {
	cfg   Config
	opts  widgets.OptionSet
	cache widgets.RenderCache[Config, widgets.Sprite]
}

// New creates a gradient widget instance satisfying widgets.Renderable.
// The returned value also implements widgets.Described and widgets.Configurable.
// When WithCaching is provided, the render path is wrapped in a render cache.
// When WithLabel is provided, the sprite label is overridden.
func New(cfg Config, opts ...widgets.Option) widgets.Renderable {
	optSet := widgets.ApplyOptions(opts...)

	w := &gradientWidget{
		cfg:  cfg,
		opts: optSet,
	}

	if optSet.CachingEnabled {
		w.cache = widgets.NewRenderCache[Config, widgets.Sprite](w.renderSprite, sign)
	}

	return w
}

// RenderFrame executes the gradient render logic and returns a positioned
// Sprite, or nil if the widget cannot render (invalid bounds/config).
func (w *gradientWidget) RenderFrame() *widgets.Sprite {
	if w.cache != nil {
		return w.cache.Render(w.cfg)
	}
	return w.renderSprite(w.cfg)
}

// renderSprite performs the actual rendering, producing a *widgets.Sprite.
func (w *gradientWidget) renderSprite(cfg Config) *widgets.Sprite {
	sprite := Render(cfg)
	if sprite == nil {
		return nil
	}

	// Apply label override if WithLabel was used.
	if w.opts.LabelOverride != "" {
		sprite.Label = w.opts.LabelOverride
	}

	return sprite
}

// Describe returns the widget's Descriptor for registry and suppression evaluation.
// Gradient is purely graphical and not eink-safe, so Capabilities is empty.
func (w *gradientWidget) Describe() widgets.Descriptor {
	return widgets.Descriptor{
		Name:         "gradient",
		MinWidth:     1,
		MinHeight:    1,
		Capabilities: []string{},
	}
}

// Configure updates the widget's configuration for the next RenderFrame call.
// Accepts a Config value; panics on type mismatch (programming error).
func (w *gradientWidget) Configure(cfg interface{}) {
	w.cfg = cfg.(Config)
}

// init registers the gradient widget type in the global widget registry.
func init() {
	widgets.Register("gradient", func() widgets.Described {
		return &gradientWidget{
			cfg: Config{
				Style:  Linear,
				Bounds: image.Rect(0, 0, 10, 10),
				Stops: []ColorStop{
					{Position: 0.0, Color: color.RGBA{A: 255}},
					{Position: 1.0, Color: color.RGBA{R: 255, G: 255, B: 255, A: 255}},
				},
			},
		}
	})
}
