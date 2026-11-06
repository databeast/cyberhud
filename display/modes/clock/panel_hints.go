package clock

import (
	"github.com/databeast/cyberhud/display/region/modehints"
)

// fitnessNotesPostApply generates fitness notes for the applied style
// against the current panel hints when the "style" key changes.
//
// This legitimately uses the global store: it runs when a console command changes
// the style, where there is no ModeInstance in scope to ask. It is advisory output
// for the operator, not layout input, so a single-panel approximation is acceptable
// here in a way it is not on the render path.
var fitnessNotesPostApply = clockRegistry.FitnessPostApply(modehints.Current, func() string { return GetPolicy().Style })
