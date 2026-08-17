package attract_starfield

import (
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// getPanelHints returns the centrally stored panel hints (see modehints).
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }
