package usb

import (
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }

// fitnessNotesPostApply generates fitness notes for the applied style
// against the current panel hints when the "style" key changes.
var fitnessNotesPostApply = usbRegistry.FitnessPostApply(modehints.Current, func() string { return PolicySnapshot().Style })
