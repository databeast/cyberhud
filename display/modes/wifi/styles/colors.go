package styles

import (
	"image/color"

	sharedcolor "github.com/databeast/cyberhud/display/style/color"
)

// resolveFGColor returns the RGBA for a named foreground color.
// Unlike the shared ResolveAccent (which defaults to white), the WiFi mode
// defaults unknown names to green — matching the Policy default FGColor.
//
// Mapping:
//
//	"green" → {0, 200, 0, 255}
//	"amber" → {255, 191, 0, 255}
//	"red"   → {255, 0, 0, 255}
//	"white" → {255, 255, 255, 255}
//	"cyan"  → {0, 255, 255, 255}
//	"none"  → {255, 255, 255, 255} (same as white)
//	unknown → {0, 200, 0, 255} (green)
func resolveFGColor(name string) color.RGBA {
	if name == "none" {
		return color.RGBA{255, 255, 255, 255}
	}
	c := sharedcolor.Lookup(name)
	// sharedcolor.Lookup returns white for unrecognized names.
	// WiFi mode defaults to green instead, so detect the fallback case.
	if c == (color.RGBA{255, 255, 255, 255}) && name != "white" {
		return color.RGBA{0, 200, 0, 255} // green default
	}
	return c
}
