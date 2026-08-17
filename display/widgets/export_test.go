package widgets

// export_test.go re-exports unexported symbols for use by external test files
// (package widgets_test). This is a standard Go testing pattern.

// NewCachedRendererForTest exposes newCachedRenderer for external test packages.
// It returns the concrete *cachedRenderer pointer (not the RenderCache interface)
// so that tests can verify pointer identity on cache hits.
func NewCachedRendererForTest[C any, R any](render func(C) *R, sign func(C) uint64) *cachedRenderer[C, R] {
	return newCachedRenderer(render, sign)
}
