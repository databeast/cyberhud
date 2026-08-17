package led

import (
	"time"

	"github.com/databeast/cyberhud/display/widgets"
)

// Compile-time interface assertions.
var (
	_ widgets.Renderable   = (*ledWidget)(nil)
	_ widgets.Described    = (*ledWidget)(nil)
	_ widgets.Configurable = (*ledWidget)(nil)
	_ widgets.Animated     = (*ledWidget)(nil)
)

// ledWidget is the internal struct implementing Renderable, Described,
// Configurable, and Animated interfaces for the LED widget.
type ledWidget struct {
	cfg   Config
	opts  widgets.OptionSet
	cache widgets.RenderCache[Config, widgets.Sprite]
}

// New creates an LED widget instance satisfying widgets.Renderable.
// The returned value also implements Described, Configurable, and Animated.
// When WithCaching is provided, the render path is wrapped in a 2-slot ring
// buffer cache that checks sign(cfg) before re-rendering.
func New(cfg Config, opts ...widgets.Option) widgets.Renderable {
	set := widgets.ApplyOptions(opts...)

	w := &ledWidget{
		cfg:  cfg,
		opts: set,
	}

	if set.CachingEnabled {
		w.cache = widgets.NewRenderCache[Config, widgets.Sprite](w.renderSprite, sign)
	}

	return w
}

// RenderFrame produces a Sprite for the current frame, or nil if the
// LED cannot render (diameter < 3).
func (w *ledWidget) RenderFrame() *widgets.Sprite {
	if w.cache != nil {
		return w.cache.Render(w.cfg)
	}
	return w.renderSprite(w.cfg)
}

// renderSprite is the internal render function that converts a Config into a Sprite.
func (w *ledWidget) renderSprite(cfg Config) *widgets.Sprite {
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
func (w *ledWidget) Describe() widgets.Descriptor {
	return widgets.Descriptor{
		Name:         "led",
		MinWidth:     3,
		MinHeight:    3,
		Capabilities: []string{"animated"},
	}
}

// Configure updates the widget's parameters between frames.
// Accepts a Config value; panics on type mismatch (programming error).
func (w *ledWidget) Configure(cfg interface{}) {
	w.cfg = cfg.(Config)
}

// Tick advances the LED's animation state. When the configured animation type
// is not NoAnimation, elapsed time is accumulated into cfg.animElapsed which
// drives Pulse, Blink, and Fade curves. When NoAnimation is set, Tick is a no-op.
func (w *ledWidget) Tick(elapsed time.Duration) {
	if w.cfg.Animation.Type == NoAnimation {
		return
	}
	w.cfg.animElapsed += elapsed
}

// init registers the LED widget type in the global widget registry.
func init() {
	widgets.Register("led", func() widgets.Described {
		return &ledWidget{
			cfg: Config{
				State:      Off,
				Diameter:   10,
				Brightness: -1.0,
			},
		}
	})
}
