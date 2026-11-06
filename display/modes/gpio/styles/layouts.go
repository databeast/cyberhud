package styles

import (
	"image/color"

	sharedcolor "github.com/databeast/cyberhud/display/style/color"
)

// Accent color resolution for the GPIO live-functions view.
//
// A live function (a HIGH standalone pin, or an active bus) renders in the
// resolved accent; a quiet one (a LOW pin) renders dimmed. The renderer drops
// color entirely on mono/grayscale panels, so these are only consulted when the
// panel is color-capable and policy enables color.

func resolveFGColor(accent string) color.RGBA {
	return sharedcolor.ResolveAccent(accent)
}

func dimFGColor(accent string) color.RGBA {
	return sharedcolor.Dim(resolveFGColor(accent))
}
