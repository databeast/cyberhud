package scrollbar

import (
	"image"

	"github.com/databeast/cyberhud/display/widgets"
)

// Compile-time interface assertions.
var (
	_ widgets.Renderable   = (*scrollbarWidget)(nil)
	_ widgets.Described    = (*scrollbarWidget)(nil)
	_ widgets.Configurable = (*scrollbarWidget)(nil)
)

// scrollbarWidget is the internal struct implementing Renderable, Described,
// and Configurable for the scrollbar widget.
type scrollbarWidget struct {
	cfg   Config
	opts  widgets.OptionSet
	cache widgets.RenderCache[Config, widgets.Sprite]
}

// New creates a scrollbar widget instance satisfying widgets.Renderable.
// The returned value also implements widgets.Described and widgets.Configurable.
// When WithCaching is provided, the render path is wrapped in a render cache.
// When WithLabel is provided, the sprite label is overridden.
func New(cfg Config, opts ...widgets.Option) widgets.Renderable {
	optSet := widgets.ApplyOptions(opts...)

	w := &scrollbarWidget{
		cfg:  cfg,
		opts: optSet,
	}

	if optSet.CachingEnabled {
		w.cache = widgets.NewRenderCache[Config, widgets.Sprite](w.renderSprite, sign)
	}

	return w
}

// RenderFrame executes the scrollbar render logic and returns a positioned
// Sprite, or nil if the widget cannot render (invalid bounds or TotalItems < 1).
func (w *scrollbarWidget) RenderFrame() *widgets.Sprite {
	if w.cache != nil {
		return w.cache.Render(w.cfg)
	}
	return w.renderSprite(w.cfg)
}

// renderSprite performs the actual rendering, producing a *widgets.Sprite.
func (w *scrollbarWidget) renderSprite(cfg Config) *widgets.Sprite {
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
// The scrollbar is purely graphical with no special capabilities (not eink-safe).
func (w *scrollbarWidget) Describe() widgets.Descriptor {
	return widgets.Descriptor{
		Name:         "scrollbar",
		MinWidth:     1,
		MinHeight:    1,
		Capabilities: []string{},
	}
}

// Configure updates the widget's configuration for the next RenderFrame call.
// Accepts a Config value; panics on type mismatch (programming error).
func (w *scrollbarWidget) Configure(cfg interface{}) {
	w.cfg = cfg.(Config)
}

// init registers the scrollbar widget type in the global widget registry.
func init() {
	widgets.Register("scrollbar", func() widgets.Described {
		return &scrollbarWidget{
			cfg: Config{
				Bounds: image.Rect(0, 0, 4, 20),
			},
		}
	})
}
