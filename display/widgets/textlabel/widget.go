package textlabel

import (
	"image"

	"github.com/databeast/cyberhud/display/widgets"
)

// Compile-time interface assertions.
var (
	_ widgets.Renderable   = (*textlabelWidget)(nil)
	_ widgets.Described    = (*textlabelWidget)(nil)
	_ widgets.Configurable = (*textlabelWidget)(nil)
)

// textlabelWidget is the internal struct implementing Renderable, Described,
// and Configurable for the textlabel widget.
type textlabelWidget struct {
	cfg   Config
	opts  widgets.OptionSet
	cache widgets.RenderCache[Config, widgets.Sprite]
}

// New creates a textlabel widget instance satisfying widgets.Renderable.
// The returned value also implements widgets.Described and widgets.Configurable.
// When WithCaching is provided, the render path is wrapped in a render cache.
// When WithLabel is provided, the sprite label is overridden.
func New(cfg Config, opts ...widgets.Option) widgets.Renderable {
	optSet := widgets.ApplyOptions(opts...)

	w := &textlabelWidget{
		cfg:  cfg,
		opts: optSet,
	}

	if optSet.CachingEnabled {
		w.cache = widgets.NewRenderCache[Config, widgets.Sprite](w.renderSprite, sign)
	}

	return w
}

// RenderFrame executes the textlabel render logic and returns a positioned
// Sprite, or nil if the widget cannot render (invalid bounds).
func (w *textlabelWidget) RenderFrame() *widgets.Sprite {
	if w.cache != nil {
		return w.cache.Render(w.cfg)
	}
	return w.renderSprite(w.cfg)
}

// renderSprite performs the actual rendering, producing a *widgets.Sprite.
func (w *textlabelWidget) renderSprite(cfg Config) *widgets.Sprite {
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
func (w *textlabelWidget) Describe() widgets.Descriptor {
	return widgets.Descriptor{
		Name:         "textlabel",
		MinWidth:     1,
		MinHeight:    1,
		Capabilities: []string{"eink-safe"},
	}
}

// Configure updates the widget's configuration for the next RenderFrame call.
// Accepts a Config value; panics on type mismatch (programming error).
func (w *textlabelWidget) Configure(cfg interface{}) {
	w.cfg = cfg.(Config)
}

// init registers the textlabel widget type in the global widget registry.
func init() {
	widgets.Register("textlabel", func() widgets.Described {
		return &textlabelWidget{
			cfg: Config{
				Bounds: image.Rect(0, 0, 10, 10),
			},
		}
	})
}
