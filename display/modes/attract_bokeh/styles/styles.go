package styles

import "github.com/databeast/cyberhud/display/style"
import "github.com/databeast/cyberhud/display/modes/attract_bokeh/source"

// Compile-time interface compliance checks.
var (
	_ style.Style[source.BokehFrame, source.Policy] = BokehColorFastStyle{}
	_ style.Style[source.BokehFrame, source.Policy] = BokehMonoFastStyle{}
	_ style.Style[source.BokehFrame, source.Policy] = BokehMonoSlowStyle{}
	_ style.Style[source.BokehFrame, source.Policy] = BokehGrayscaleFastStyle{}
	_ style.Style[source.BokehFrame, source.Policy] = BokehGrayscaleSlowStyle{}
	_ style.Style[source.BokehFrame, source.Policy] = BokehColorSlowStyle{}
)
