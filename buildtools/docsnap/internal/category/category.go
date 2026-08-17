// Package category provides shared constants and prefix-matching logic for
// display category classification of snapshot image filenames.
package category

import "strings"

// Category represents a display category name.
type Category string

const (
	Color     Category = "Color"
	EInk      Category = "E-Ink"
	Grayscale Category = "Grayscale"
	Mono      Category = "Mono"
)

// PrefixMapping associates a filename prefix with a display category.
type PrefixMapping struct {
	Prefix   string
	Category Category
}

// Prefixes maps filename prefixes to categories, ordered longest-first
// to ensure correct longest-prefix matching.
var Prefixes = []PrefixMapping{
	{Prefix: "color-slow-", Category: Color},
	{Prefix: "color-fast-", Category: Color},
	{Prefix: "grayscale-fast-", Category: Grayscale},
	{Prefix: "grayscale-slow-", Category: Grayscale},
	{Prefix: "mono-slow-", Category: Mono},
	{Prefix: "mono-fast-", Category: Mono},
	{Prefix: "color-", Category: Color},
	{Prefix: "mono-", Category: Mono},
	{Prefix: "eink-", Category: EInk},
}

// Match returns the category for a filename using longest-prefix matching.
// Returns ("", false) if no prefix matches.
func Match(filename string) (Category, bool) {
	var bestCategory Category
	bestLen := 0

	for _, pm := range Prefixes {
		if strings.HasPrefix(filename, pm.Prefix) && len(pm.Prefix) > bestLen {
			bestCategory = pm.Category
			bestLen = len(pm.Prefix)
		}
	}

	if bestLen == 0 {
		return "", false
	}
	return bestCategory, true
}
