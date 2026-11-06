package storage

import (
	"github.com/databeast/cyberhud/display/region/modehints"
)

// fitnessNotesPostApply generates fitness notes for the applied style against
// the current panel hints when the "style" key changes. It legitimately uses
// the global store: it runs when a console command changes the style, where
// there is no ModeInstance in scope to ask.
var fitnessNotesPostApply = storageRegistry.FitnessPostApply(modehints.Current, func() string { return GetPolicy().Style })
