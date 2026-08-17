package widgets

// cacheEntry holds a single cached render result.
type cacheEntry[R any] struct {
	sig    uint64
	result *R
	valid  bool
}

// RenderCache is the public interface for the caching subsystem. Widget
// packages store a RenderCache and call Render instead of their raw render
// function. The underlying implementation (cachedRenderer) is unexported.
type RenderCache[C any, R any] interface {
	Render(C) *R
}

// NewRenderCache creates a render cache wrapping a render function and a
// signature hash function. This is the only public entry point to the caching
// subsystem — the underlying cachedRenderer type is an internal detail.
func NewRenderCache[C any, R any](render func(C) *R, sign func(C) uint64) RenderCache[C, R] {
	return newCachedRenderer(render, sign)
}

// cachedRenderer provides opt-in memoization for widget Render functions.
// It caches the two most recent render results using a 2-slot ring buffer,
// preventing thrashing when a widget alternates between two Config signatures
// (e.g., an LED blinking On/Off every frame).
//
// This type is unexported — widgets opt into caching via the WithCaching
// functional option, and the constructor is accessed via NewRenderCache.
//
// Type parameters:
//
//	C — the widget's Config type
//	R — the widget's Result type (e.g., progressbar.Result)
type cachedRenderer[C any, R any] struct {
	render  func(C) *R // underlying stateless render function
	sign    func(C) uint64
	entries [2]cacheEntry[R]
	next    int // index of the slot to evict on the next miss (alternates 0→1→0→1)
}

// newCachedRenderer wraps a Render function with two-entry caching.
// The sign function hashes a Config into a uint64; if the signature matches
// either cached slot, the cached *R is returned without calling render.
func newCachedRenderer[C any, R any](render func(C) *R, sign func(C) uint64) *cachedRenderer[C, R] {
	return &cachedRenderer[C, R]{
		render: render,
		sign:   sign,
	}
}

// Render returns the cached result when the Config signature matches either
// cache slot, or calls the underlying render function, stores the result in
// the next eviction slot, and advances the ring index.
func (c *cachedRenderer[C, R]) Render(cfg C) *R {
	sig := c.sign(cfg)

	// Check both cache slots for a hit.
	for i := range c.entries {
		if c.entries[i].valid && c.entries[i].sig == sig {
			return c.entries[i].result
		}
	}

	// Cache miss — render and store in the next eviction slot.
	result := c.render(cfg)
	c.entries[c.next] = cacheEntry[R]{sig: sig, result: result, valid: true}
	c.next = 1 - c.next // alternate: 0→1→0→1
	return result
}
