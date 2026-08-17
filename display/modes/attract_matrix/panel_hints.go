package attract_matrix

import (
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }

// fitnessNotesPostApply is the PostApply hook for the attract_matrix CmdHandler.
// This mode auto-resolves styles from panel hints; no policy key affects
// style selection, so no fitness notes are generated.
func fitnessNotesPostApply(appliedKeys []string) []string { return nil }
