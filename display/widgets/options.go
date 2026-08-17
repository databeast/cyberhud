package widgets

// OptionSet holds common optional parameters for widget constructors.
// When no options are applied, all fields retain their zero values,
// which represent sensible defaults (no caching, no label override).
type OptionSet struct {
	CachingEnabled bool
	LabelOverride  string
}

// Option is a functional option that configures an OptionSet.
// Widget constructors accept variadic ...Option parameters after
// the required Config argument.
type Option func(*OptionSet)

// WithCaching returns an Option that enables render caching on the
// widget. When applied, the widget constructor wraps the render path
// in a cachedRenderer (via NewRenderCache) that skips re-rendering
// when the configuration signature has not changed.
func WithCaching() Option {
	return func(o *OptionSet) { o.CachingEnabled = true }
}

// WithLabel returns an Option that overrides the widget's default
// sprite label with the provided string. This affects the Label field
// on Sprites produced by the widget.
func WithLabel(label string) Option {
	return func(o *OptionSet) { o.LabelOverride = label }
}

// ApplyOptions applies the given options to a fresh OptionSet and
// returns it. This is a convenience for widget constructors.
func ApplyOptions(opts ...Option) OptionSet {
	var set OptionSet
	for _, opt := range opts {
		opt(&set)
	}
	return set
}
