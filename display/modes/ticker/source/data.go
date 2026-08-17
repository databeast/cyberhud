package source

import (
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// Compile-time interface compliance checks.

// TickerSnapshot captures the current ticker state for style Build methods.
// It provides all fields needed by both PlainStyle and BorderedStyle to render.
type TickerSnapshot struct {
	Directives   []LineDirective      // Current feed directives.
	Policy       Policy               // Current ticker policy.
	ScrollOffset int                  // Vertical scroll offset (line index).
	StripSprites []widgets.Sprite     // Pre-rendered sprites from active strips; nil if inactive.
	Hints        textlayout.TextHints // TextHints for catalog-validated font resolution.
}

// tickerRegistry is declared in styles.go with the full multi-capability
// registration. This comment marks the former location for traceability.

// allEmptyItems returns true if items is nil, empty, or contains only empty strings.
func AllEmptyItems(items []string) bool {
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
