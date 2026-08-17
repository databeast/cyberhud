package attract_hacking

import (
	"github.com/databeast/cyberhud/display/region/modehints"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// getPanelHints returns the centrally stored panel hints for the current region.
func getPanelHints() (textlayout.TextHints, bool) { return modehints.Current() }
