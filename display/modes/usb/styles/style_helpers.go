package styles

import (
	"image"
)

// IconConsumer is a local interface for styles that accept an icon getter.
// BuildView checks for this interface via a type assertion and injects the
// getter before calling Build.
type IconConsumer interface {
	SetIconGetter(fn func(name string) (image.Image, bool))
}

// Compile-time interface compliance checks.
var (
	_ IconConsumer = (*def)(nil)
)

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
