package menu

import (
	"github.com/databeast/cyberhud/display/region/modehints"
)

// fitnessNotesPostApply generates fitness notes for the applied style
// against the current panel hints when the "style" key changes.
var fitnessNotesPostApply = menuRegistry.FitnessPostApply(modehints.Current, func() string { return GetPolicy().Style })
