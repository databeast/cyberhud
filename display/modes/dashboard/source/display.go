package source

import "github.com/databeast/cyberhud/display/surface/textlayout"

// getPanelName derives the active panel name from TextHints.
// It combines PanelProduct and ScreenName with the following logic:
//   - Both set → "product/screen" truncated to 64 chars
//   - Only product set → "product" truncated to 64 chars
//   - Only screen set → "(unknown)/screen" truncated to 64 chars
//   - Neither set → "(unknown)"
func getPanelName(hints textlayout.TextHints) string {
	product := hints.PanelProduct
	screen := hints.ScreenName

	var name string
	switch {
	case product != "" && screen != "":
		name = product + "/" + screen
	case product != "":
		name = product
	case screen != "":
		name = "(unknown)/" + screen
	default:
		return "(unknown)"
	}
	if len(name) > 64 {
		name = name[:64]
	}
	return name
}

func GetPanelName(hints textlayout.TextHints) string { return getPanelName(hints) }
