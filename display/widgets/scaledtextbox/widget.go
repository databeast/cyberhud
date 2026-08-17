package scaledtextbox

import (
	"github.com/databeast/cyberhud/display/widgets"
)

// scaledTextBoxWidget is the internal struct implementing Renderable, Described,
// and Configurable interfaces for the ScaledTextBox widget.
type scaledTextBoxWidget struct {
	cfg   Config
	opts  widgets.OptionSet
	cache widgets.RenderCache[Config, widgets.Sprite]
}

// New creates a ScaledTextBox widget instance satisfying widgets.Renderable.
// The returned value also implements Described and Configurable.
func New(cfg Config, opts ...widgets.Option) widgets.Renderable {
	set := widgets.ApplyOptions(opts...)

	w := &scaledTextBoxWidget{
		cfg:  cfg,
		opts: set,
	}

	if set.CachingEnabled {
		w.cache = widgets.NewRenderCache[Config, widgets.Sprite](w.renderSprite, sign)
	}

	return w
}

// RenderFrame produces a Sprite for the current frame, or nil if the
// ScaledTextBox cannot render (invalid dimensions).
func (w *scaledTextBoxWidget) RenderFrame() *widgets.Sprite {
	if w.cache != nil {
		return w.cache.Render(w.cfg)
	}
	return w.renderSprite(w.cfg)
}

// renderSprite is the internal render function that converts a Config into a Sprite.
func (w *scaledTextBoxWidget) renderSprite(cfg Config) *widgets.Sprite {
	sprite := Render(cfg)
	if sprite == nil {
		return nil
	}

	if w.opts.LabelOverride != "" {
		sprite.Label = w.opts.LabelOverride
	}

	return sprite
}

// Describe returns the widget's metadata for registry and suppression.
func (w *scaledTextBoxWidget) Describe() widgets.Descriptor {
	return widgets.Descriptor{
		Name:         "scaledtextbox",
		MinWidth:     1,
		MinHeight:    1,
		Capabilities: []string{"eink-safe"},
	}
}

// Configure updates the widget's parameters between frames.
// Accepts a Config value; panics on type mismatch (programming error).
func (w *scaledTextBoxWidget) Configure(cfg interface{}) {
	w.cfg = cfg.(Config)
}

// init registers the ScaledTextBox widget type in the global widget registry.
func init() {
	widgets.Register("scaledtextbox", func() widgets.Described {
		w := &scaledTextBoxWidget{
			cfg: Config{},
		}
		return w
	})
}
