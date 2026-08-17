package attract_geometric

import "github.com/databeast/cyberhud/display/modes/attract_geometric/source"

// computeScaleFactor returns the proportional scale factor for the given panel
// dimensions. Panels with a shortest dimension >= 240 return 1.0; smaller panels
// return a proportional fraction.
func computeScaleFactor(panelWidth, panelHeight int) float64 {
	shortest := panelWidth
	if panelHeight < shortest {
		shortest = panelHeight
	}
	sf := float64(shortest) / 240.0
	if sf > 1.0 {
		sf = 1.0
	}
	return sf
}

// scaledSizeMin returns the scaled SIZE_MIN value. When the scaled SIZE_MAX
// would fall below 20, the fallback clamp of 4 is applied.
func scaledSizeMin(sf float64) float64 {
	if source.SizeMax*sf < 20 {
		return 4
	}
	return source.SizeMin * sf
}

// scaledSizeMax returns the scaled SIZE_MAX value. When the scaled SIZE_MAX
// would fall below 20, the fallback clamp of 20 is applied.
func scaledSizeMax(sf float64) float64 {
	if source.SizeMax*sf < 20 {
		return 20
	}
	return source.SizeMax * sf
}

// scaledProximityRadius returns the scaled CLUSTER_PROXIMITY_RADIUS value.
func scaledProximityRadius(sf float64) float64 {
	return source.ClusterProximityRadius * sf
}

// scaledFontSizeMin returns the scaled minimum font size for fragments.
func scaledFontSizeMin(sf float64) float64 {
	return 10 * sf
}

// scaledFontSizeMax returns the scaled maximum font size for fragments.
func scaledFontSizeMax(sf float64) float64 {
	return 16 * sf
}

// shouldRenderFragments returns false when the panel's shortest dimension is
// below 80 pixels, indicating that fragment text should be skipped entirely.
func shouldRenderFragments(panelWidth, panelHeight int) bool {
	shortest := panelWidth
	if panelHeight < shortest {
		shortest = panelHeight
	}
	return shortest >= 80
}
