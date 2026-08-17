package source

import "image/color"

// buildColorArray constructs the per-cell color gradient for a rain column.
// The returned slice has length trailLength+1, where index 0 is the lead cell
// and subsequent indices represent the fading trail behind it.
//
// In normal (color) mode:
//   - Index 0: bright white-green lead {180, 255, 180, 255}
//   - Indices 1..trailLength/2: linear interpolation from bright green
//     {0, 255, 70, 255} to dark green {0, 100, 0, 255}
//   - Indices trailLength/2+1..trailLength: linear interpolation from dark
//     green {0, 100, 0, 255} to black {0, 0, 0, 255}
//
// In mono mode:
//   - Index 0: white {255, 255, 255, 255}
//   - Indices 1..trailLength: linear interpolation from white to black with
//     equal R, G, B channels
//
// If trailLength is 0, the returned slice contains only the lead color.
func buildColorArray(trailLength int, mono bool) []color.RGBA {
	if trailLength <= 0 {
		if mono {
			return []color.RGBA{{255, 255, 255, 255}}
		}
		return []color.RGBA{{180, 255, 180, 255}}
	}

	colors := make([]color.RGBA, trailLength+1)

	if mono {
		// Mono mode: lead is white, trail linearly fades to black.
		colors[0] = color.RGBA{255, 255, 255, 255}
		// Interpolate from white {255,255,255} to black {0,0,0} across indices 1..trailLength.
		// segLen = trailLength so that at i=trailLength the formula yields exactly 0 (black).
		segLen := trailLength
		for i := 1; i <= trailLength; i++ {
			v := uint8(lerpInt(255, 0, i, segLen))
			colors[i] = color.RGBA{v, v, v, 255}
		}
		return colors
	}

	// Normal (color) mode.
	colors[0] = color.RGBA{180, 255, 180, 255}

	// First half: indices 1..trailLength/2
	// Interpolate from bright green {0, 255, 70, 255} to dark green {0, 100, 0, 255}.
	firstHalfEnd := trailLength / 2
	segLen := firstHalfEnd // at i=firstHalfEnd, t=1 → reaches dark green
	for i := 1; i <= firstHalfEnd; i++ {
		r := uint8(lerpInt(0, 0, i, segLen))
		g := uint8(lerpInt(255, 100, i, segLen))
		b := uint8(lerpInt(70, 0, i, segLen))
		colors[i] = color.RGBA{r, g, b, 255}
	}

	// Second half: indices trailLength/2+1..trailLength
	// Interpolate from dark green {0, 100, 0, 255} to black {0, 0, 0, 255}.
	secondHalfStart := firstHalfEnd + 1
	numSecondHalf := trailLength - firstHalfEnd
	segLen = numSecondHalf // at i=numSecondHalf, t=1 → reaches black
	for i := secondHalfStart; i <= trailLength; i++ {
		idx := i - secondHalfStart + 1
		r := uint8(lerpInt(0, 0, idx, segLen))
		g := uint8(lerpInt(100, 0, idx, segLen))
		b := uint8(lerpInt(0, 0, idx, segLen))
		colors[i] = color.RGBA{r, g, b, 255}
	}

	return colors
}

// lerpInt performs integer linear interpolation from a to b at step i out of
// segLen total steps. It uses integer arithmetic with truncation toward zero:
//
//	int(float64(a) + float64(b-a)*float64(i)/float64(segLen))
func lerpInt(a, b, i, segLen int) int {
	if segLen <= 0 {
		return a
	}
	return int(float64(a) + float64(b-a)*float64(i)/float64(segLen))
}
