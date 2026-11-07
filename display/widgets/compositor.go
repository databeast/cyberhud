package widgets

// Compositor accumulates rendered Sprites from widgets, applying optional
// suppression rules before invoking RenderFrame. It is created per BuildView
// call and holds no inter-frame state.
type Compositor struct {
	sprites []Sprite
	rules   []SuppressionRule
	ctx     SuppressionContext
}

// NewCompositor creates a Compositor with the given suppression context and
// optional suppression rules. Rules are evaluated in order; if any rule
// returns true for a widget's Descriptor, that widget is skipped entirely
// without invoking RenderFrame.
func NewCompositor(ctx SuppressionContext, rules ...SuppressionRule) *Compositor {
	return &Compositor{
		ctx:   ctx,
		rules: rules,
	}
}

// Add invokes RenderFrame on r. If r also implements Described and any
// configured suppression rule matches, the widget is skipped entirely
// (RenderFrame is NOT called). Otherwise, RenderFrame is called and
// non-nil results are appended to the internal sprite slice.
// Nil results from RenderFrame are discarded.
func (c *Compositor) Add(r Renderable) {
	if c.suppressed(r) {
		return
	}
	sprite := r.RenderFrame()
	if sprite != nil {
		c.sprites = append(c.sprites, *sprite)
	}
}

// AddIf is a conditional Add. When condition is false, r is not invoked
// at all and no sprite is added. When condition is true, behaves
// identically to Add.
func (c *Compositor) AddIf(condition bool, r Renderable) {
	if !condition {
		return
	}
	c.Add(r)
}

// Sprites returns the accumulated sprite slice in insertion order.
//
// Insertion order is z-order only within a single compositing depth: the region
// renderer orders sprites across depths by Sprite.Z (ascending, around the text
// plane) and preserves this insertion order to break ties within equal Z. A
// backdrop added first but left at the default Z still composites above a later
// sprite given a positive Z — the depth wins, insertion order only settles ties.
// The insertion order defines z-order: first added renders below,
// last added renders on top.
func (c *Compositor) Sprites() []Sprite {
	return c.sprites
}

// suppressed checks whether r implements Described and, if so, evaluates
// all configured suppression rules. Returns true if any rule suppresses
// the widget (short-circuit OR).
func (c *Compositor) suppressed(r Renderable) bool {
	if len(c.rules) == 0 {
		return false
	}
	d, ok := r.(Described)
	if !ok {
		return false
	}
	desc := d.Describe()
	for _, rule := range c.rules {
		if rule(desc, c.ctx) {
			return true
		}
	}
	return false
}
