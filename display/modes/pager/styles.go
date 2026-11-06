package pager

import (
	"github.com/databeast/cyberhud/display/modes/pager/source"
	"github.com/databeast/cyberhud/display/modes/pager/styles"
	"github.com/databeast/cyberhud/display/style"
)

// pagerRegistry is the per-mode StyleRegistry for the pager display mode. The
// style set is generated programmatically in the styles subpackage (a dimension
// table per capability class run through shared constructors) and handed here
// via styles.All, in capability order — MonoSlow → MonoFast → GrayscaleSlow →
// GrayscaleFast → ColorSlow → ColorFast, which also defines fitness tie-breaking.
var pagerRegistry = style.NewRegistry[source.PagerSnapshot, source.Policy](styles.All()...)
