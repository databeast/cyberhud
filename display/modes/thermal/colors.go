package thermal

import "image/color"

// accentColor returns the RGBA value for a named accent color.
// For "thermal" and "none", it returns white as a default since those
// are handled specially by the caller (thermal delegates to severity colors,
// none renders text in white).
func accentColor(name string) color.RGBA {
	switch name {
	case "cyan":
		return color.RGBA{R: 0, G: 255, B: 255, A: 255}
	case "green":
		return color.RGBA{R: 0, G: 200, B: 0, A: 255}
	case "amber":
		return color.RGBA{R: 255, G: 191, B: 0, A: 255}
	case "red":
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	case "white":
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	default:
		// "thermal", "none", or unrecognized → white
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
}

// resolveAccentForLabel returns the color to use for zone labels, column headers,
// and unit symbols based on the accent setting and panel type.
//
// Rules:
//   - Monochrome panel: returns nativeFG (ignore accent entirely)
//   - "thermal": returns the severity color for the zone
//   - Named accent (cyan/green/amber/red/white): returns the accent color
//   - "none": returns white (255, 255, 255, 255)
func resolveAccentForLabel(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	// Monochrome panels ignore accent, use native foreground.
	if !isColor {
		return nativeFG
	}

	switch accent {
	case "thermal":
		return severityColorRGBA(severityLevel)
	case "none":
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	default:
		// Named accent: cyan, green, amber, red, white
		return accentColor(accent)
	}
}

// resolveAccentForTemp returns the color to use for temperature numeric values
// based on the accent setting and panel type.
//
// Rules:
//   - Monochrome panel: returns nativeFG (ignore accent entirely)
//   - "thermal": returns severity color (green/yellow/red)
//   - Named accent: returns severity color (temps always use severity)
//   - "none": returns white (255, 255, 255, 255)
func resolveAccentForTemp(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	// Monochrome panels ignore accent, use native foreground.
	if !isColor {
		return nativeFG
	}

	switch accent {
	case "none":
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	default:
		// "thermal" and all named accents: temperature values use severity colors
		return severityColorRGBA(severityLevel)
	}
}

// resolveAccentForBar returns the color to use for progress bar fills
// based on the accent setting and panel type.
//
// Rules:
//   - Monochrome panel: returns nativeFG
//   - All accents (including "none"): bar fills use severity colors
func resolveAccentForBar(accent string, isColor bool, severityLevel int, nativeFG color.RGBA) color.RGBA {
	// Monochrome panels ignore accent, use native foreground.
	if !isColor {
		return nativeFG
	}

	// Bar fills always use severity colors regardless of accent setting.
	return severityColorRGBA(severityLevel)
}

// severityColorRGBA returns the severity color as color.RGBA.
//
//	0 (Normal)   → green  (0, 255, 0, 255)
//	1 (Warning)  → yellow (255, 255, 0, 255)
//	2 (Critical) → red    (255, 0, 0, 255)
func severityColorRGBA(level int) color.RGBA {
	switch level {
	case 1:
		return color.RGBA{R: 255, G: 255, B: 0, A: 255}
	case 2:
		return color.RGBA{R: 255, G: 0, B: 0, A: 255}
	default:
		return color.RGBA{R: 0, G: 255, B: 0, A: 255}
	}
}
