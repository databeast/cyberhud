package progressbar

import (
	"image"
	"time"

	"github.com/databeast/cyberhud/display/widgets"
)

// progressbarWidget is the internal struct implementing Renderable, Described,
// Configurable, and Animated for the progressbar widget.
type progressbarWidget struct {
	cfg   Config
	opts  widgets.OptionSet
	cache widgets.RenderCache[Config, widgets.Sprite]
}

// New creates a progressbar widget instance satisfying widgets.Renderable.
// The returned value also implements widgets.Described, widgets.Configurable,
// and widgets.Animated.
// When WithCaching is provided, the render path is wrapped in a render cache.
// When WithLabel is provided, the sprite label is overridden.
func New(cfg Config, opts ...widgets.Option) widgets.Renderable {
	optSet := widgets.ApplyOptions(opts...)

	w := &progressbarWidget{
		cfg:  cfg,
		opts: optSet,
	}

	if optSet.CachingEnabled {
		w.cache = widgets.NewRenderCache[Config, widgets.Sprite](w.renderSprite, sign)
	}

	return w
}

// RenderFrame executes the progressbar render logic and returns a positioned
// Sprite, or nil if the widget cannot render (invalid bounds).
func (w *progressbarWidget) RenderFrame() *widgets.Sprite {
	if w.cache != nil {
		return w.cache.Render(w.cfg)
	}
	return w.renderSprite(w.cfg)
}

// renderSprite performs the actual rendering, producing a *widgets.Sprite.
func (w *progressbarWidget) renderSprite(cfg Config) *widgets.Sprite {
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

// Tick advances animation state. Implements widgets.Animated.
func (w *progressbarWidget) Tick(elapsed time.Duration) {
	if w.cfg.Animation.Type == NoAnimation {
		return
	}
	w.cfg.animElapsed += elapsed
}

// Describe returns the widget's Descriptor for registry and suppression evaluation.
func (w *progressbarWidget) Describe() widgets.Descriptor {
	return widgets.Descriptor{
		Name:         "progressbar",
		MinWidth:     1,
		MinHeight:    1,
		Capabilities: []string{},
	}
}

// Configure updates the widget's configuration for the next RenderFrame call.
// Accepts a Config value; panics on type mismatch (programming error).
func (w *progressbarWidget) Configure(cfg interface{}) {
	w.cfg = cfg.(Config)
}

// init registers the progressbar widget type in the global widget registry.
func init() {
	widgets.Register("progressbar", func() widgets.Described {
		return &progressbarWidget{
			cfg: Config{
				Bounds: image.Rect(0, 0, 10, 10),
			},
		}
	})
}
