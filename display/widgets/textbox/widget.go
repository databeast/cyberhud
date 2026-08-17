package textbox

import (
	"image"

	"github.com/databeast/cyberhud/display/widgets"
)

// Compile-time interface assertions.
var (
	_ widgets.Renderable   = (*textboxWidget)(nil)
	_ widgets.Described    = (*textboxWidget)(nil)
	_ widgets.Configurable = (*textboxWidget)(nil)
)

// textboxWidget is the internal struct implementing Renderable, Described,
// and Configurable interfaces for the textbox widget.
type textboxWidget struct {
	cfg   Config
	opts  widgets.OptionSet
	cache widgets.RenderCache[Config, widgets.Sprite]
}

// New creates a textbox widget instance satisfying widgets.Renderable.
// The returned value also implements widgets.Described and widgets.Configurable.
// When WithCaching is provided, the render path is wrapped in a render cache.
// When WithLabel is provided, the sprite label is overridden.
func New(cfg Config, opts ...widgets.Option) widgets.Renderable {
	optSet := widgets.ApplyOptions(opts...)

	w := &textboxWidget{
		cfg:  cfg,
		opts: optSet,
	}

	if optSet.CachingEnabled {
		w.cache = widgets.NewRenderCache[Config, widgets.Sprite](w.renderSprite, sign)
	}

	return w
}

// RenderFrame executes the textbox render logic and returns a positioned
// Sprite, or nil if the widget cannot render (invalid bounds).
func (w *textboxWidget) RenderFrame() *widgets.Sprite {
	if w.cache != nil {
		return w.cache.Render(w.cfg)
	}
	return w.renderSprite(w.cfg)
}

// renderSprite performs the actual rendering, producing a *widgets.Sprite.
func (w *textboxWidget) renderSprite(cfg Config) *widgets.Sprite {
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
// Textbox is eink-safe since it just renders text.
func (w *textboxWidget) Describe() widgets.Descriptor {
	return widgets.Descriptor{
		Name:         "textbox",
		MinWidth:     1,
		MinHeight:    1,
		Capabilities: []string{"eink-safe"},
	}
}

// Configure updates the widget's configuration for the next RenderFrame call.
// Accepts a Config value; panics on type mismatch (programming error).
func (w *textboxWidget) Configure(cfg interface{}) {
	w.cfg = cfg.(Config)
}

// init registers the textbox widget type in the global widget registry.
func init() {
	widgets.Register("textbox", func() widgets.Described {
		return &textboxWidget{
			cfg: Config{
				Bounds: image.Rect(0, 0, 10, 10),
			},
		}
	})
}
