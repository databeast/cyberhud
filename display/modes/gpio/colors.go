package gpio

import (
	"image/color"

	sharedcolor "github.com/databeast/cyberhud/display/style/color"
)

// resolveFGColor returns the primary RGBA for a named foreground color.
// When accent is "none" (or unrecognized), resolves to opaque white {255, 255, 255, 255}.
func resolveFGColor(accent string) color.RGBA {
	return sharedcolor.ResolveAccent(accent)
}

// dimFGColor returns the dimmed variant of the resolved foreground color.
// Used for LOW-state pins in list/detail styles to produce a halved-intensity foreground.
func dimFGColor(accent string) color.RGBA {
	return sharedcolor.Dim(resolveFGColor(accent))
}

// borderInset returns the pixel inset for border rendering.
// After the layout-padding-refactor, border is decorative only — always returns 0.
func borderInset(p Policy) int {
	return 0
}
