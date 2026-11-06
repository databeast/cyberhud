package styles

import "github.com/databeast/cyberhud/display/style"

// All returns every pager style in capability order (MonoSlow → MonoFast →
// GrayscaleSlow → GrayscaleFast → ColorSlow → ColorFast). The mode root feeds
// this directly into style.NewRegistry. Registration order defines fitness
// tie-breaking, so this order is the contract.
func All() []style.Style[snapshot, policy] {
	return allStyles()
}
