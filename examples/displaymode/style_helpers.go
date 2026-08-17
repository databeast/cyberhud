package displaymode

import (
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// newDisplaymodeBridge constructs a LayoutBridge for the displaymode template.
// The template mode uses no border, no title bar, and no hint bar.
func newDisplaymodeBridge(hints textlayout.TextHints) layout.LayoutCalculator {
	return layout.NewLayoutBridge(hints, layout.BridgeConfig{})
}

// ViewData is the mode-internal representation of rendering output.
// In the template it contains only Items; real modes add domain-specific fields.
type ViewData struct {
	Items []string
}

// convertViewData converts the mode-internal ViewData to the shared style.ViewData.
// If all items are empty (nil, zero-length, or only empty strings), it substitutes
// a placeholder slice to satisfy the non-empty items contract.
func convertViewData(vd ViewData) style.ViewData {
	items := vd.Items
	if allEmptyItems(items) {
		items = []string{"(template)"}
	}
	return style.ViewData{Items: items}
}

// convertFromStyleViewData converts a shared style.ViewData back to mode-internal ViewData.
// This is the inverse of convertViewData, used by BuildView after registry dispatch.
func convertFromStyleViewData(svd style.ViewData) ViewData {
	return ViewData{Items: svd.Items}
}

// allEmptyItems returns true if items is nil, empty, or contains only empty strings.
func allEmptyItems(items []string) bool {
	if len(items) == 0 {
		return true
	}
	for _, item := range items {
		if item != "" {
			return false
		}
	}
	return true
}

// registeredStyleNames returns the list of style names from the registry in
// registration order. Used by catalog registration and cmdHandler for
// allowed-value validation.
func registeredStyleNames() []string {
	styles := templateRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}
