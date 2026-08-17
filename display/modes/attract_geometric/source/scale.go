package source

// scaledSizeMin returns the scaled SIZE_MIN value. When the scaled SIZE_MAX
// would fall below 20, the fallback clamp of 4 is applied.
func scaledSizeMin(sf float64) float64 {
	if SizeMax*sf < 20 {
		return 4
	}
	return SizeMin * sf
}

// scaledSizeMax returns the scaled SIZE_MAX value. When the scaled SIZE_MAX
// would fall below 20, the fallback clamp of 20 is applied.
func scaledSizeMax(sf float64) float64 {
	if SizeMax*sf < 20 {
		return 20
	}
	return SizeMax * sf
}

// scaledProximityRadius returns the scaled CLUSTER_PROXIMITY_RADIUS value.
func scaledProximityRadius(sf float64) float64 {
	return ClusterProximityRadius * sf
}

// scaledFontSizeMin returns the scaled minimum font size for fragments.
func scaledFontSizeMin(sf float64) float64 {
	return 10 * sf
}

// scaledFontSizeMax returns the scaled maximum font size for fragments.
func scaledFontSizeMax(sf float64) float64 {
	return 16 * sf
}
