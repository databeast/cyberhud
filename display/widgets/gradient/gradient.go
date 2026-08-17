package gradient

import (
	"encoding/binary"
	"hash/fnv"
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/databeast/cyberhud/display/widgets"
)

// Style represents the gradient rendering mode.
type Style int

const (
	Linear Style = iota
	Radial
)

// ColorStop defines a color at a normalized position along the gradient.
type ColorStop struct {
	Position float64    // Normalized position in [0.0, 1.0] (clamped before use)
	Color    color.RGBA // RGBA color at this position
}

// Config holds all parameters needed to render a gradient.
type Config struct {
	Style  Style           // Linear or Radial
	Angle  float64         // Direction in degrees for Linear (clockwise from 12-o'clock)
	Bounds image.Rectangle // Pixel region for rendering
	Stops  []ColorStop     // Color stops (minimum 2, maximum 64)
}

// Render produces a gradient image based on the given configuration.
// Returns nil for invalid configs.
func Render(cfg Config) *widgets.Sprite {
	if !validate(cfg) {
		return nil
	}

	stops := normalizeStops(cfg.Stops)
	angle := normalizeAngle(cfg.Angle)

	var img *image.RGBA
	var label string

	switch cfg.Style {
	case Linear:
		img = renderLinear(cfg.Bounds, angle, stops)
		label = "gradient/linear"
	case Radial:
		img = renderRadial(cfg.Bounds, stops)
		label = "gradient/radial"
	default:
		return nil
	}

	return &widgets.Sprite{
		Image:    img,
		Position: cfg.Bounds.Min,
		Label:    label,
	}
}

// sign produces a stable uint64 hash of a Config for render cache memoization.
// It uses FNV-1a 64-bit hashing over Style, Angle, Bounds, and all color stops.
func sign(cfg Config) uint64 {
	h := fnv.New64a()
	var buf [8]byte

	// 1. Style as 1 byte (uint8 cast)
	buf[0] = byte(cfg.Style)
	h.Write(buf[:1])

	// 2. Angle as 8 bytes via math.Float64bits
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(cfg.Angle))
	h.Write(buf[:])

	// 3. Bounds: Min.X, Min.Y, Max.X, Max.Y each as int64 (8 bytes each)
	binary.LittleEndian.PutUint64(buf[:], uint64(int64(cfg.Bounds.Min.X)))
	h.Write(buf[:])
	binary.LittleEndian.PutUint64(buf[:], uint64(int64(cfg.Bounds.Min.Y)))
	h.Write(buf[:])
	binary.LittleEndian.PutUint64(buf[:], uint64(int64(cfg.Bounds.Max.X)))
	h.Write(buf[:])
	binary.LittleEndian.PutUint64(buf[:], uint64(int64(cfg.Bounds.Max.Y)))
	h.Write(buf[:])

	// 4. Each stop: Position bits (8 bytes) + R, G, B, A (4 bytes)
	for _, stop := range cfg.Stops {
		binary.LittleEndian.PutUint64(buf[:], math.Float64bits(stop.Position))
		h.Write(buf[:])
		buf[0] = stop.Color.R
		buf[1] = stop.Color.G
		buf[2] = stop.Color.B
		buf[3] = stop.Color.A
		h.Write(buf[:4])
	}

	return h.Sum64()
}

// validate returns false for invalid configurations.
func validate(cfg Config) bool {
	if cfg.Bounds.Dx() < 1 || cfg.Bounds.Dy() < 1 {
		return false
	}

	if len(cfg.Stops) < 2 {
		return false
	}

	if cfg.Style != Linear && cfg.Style != Radial {
		return false
	}

	for _, stop := range cfg.Stops {
		if math.IsNaN(stop.Position) || math.IsInf(stop.Position, 0) {
			return false
		}
	}

	if math.IsNaN(cfg.Angle) || math.IsInf(cfg.Angle, 0) {
		return false
	}

	return true
}

// normalizeStops clamps stop positions to [0.0, 1.0], stable-sorts by position,
// and truncates to the first 64 stops if there are more.
func normalizeStops(stops []ColorStop) []ColorStop {
	// Make a copy to avoid mutating the input slice
	out := make([]ColorStop, len(stops))
	copy(out, stops)

	// Clamp each stop position to [0.0, 1.0].
	for i := range out {
		if out[i].Position < 0.0 {
			out[i].Position = 0.0
		} else if out[i].Position > 1.0 {
			out[i].Position = 1.0
		}
	}

	// Stable-sort by position, preserving original order for same-position stops.
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Position < out[j].Position
	})

	// Cap at 64 stops maximum.
	if len(out) > 64 {
		out = out[:64]
	}

	return out
}

// normalizeAngle maps any finite float64 angle to the range [0.0, 360.0).
func normalizeAngle(angle float64) float64 {
	return math.Mod(math.Mod(angle, 360)+360, 360)
}

// renderLinear allocates an *image.RGBA and fills it with a linear gradient.
// The angle is in degrees, clockwise from 12-o'clock (0° = top-to-bottom, 90° = left-to-right).
// Stops must already be normalized (clamped and sorted).
func renderLinear(bounds image.Rectangle, angle float64, stops []ColorStop) *image.RGBA {
	w := bounds.Dx()
	h := bounds.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Convert angle to radians (clockwise from north/12-o'clock)
	rad := angle * math.Pi / 180.0

	// Direction vector: sin for X (clockwise from north), cos for Y
	dx := math.Sin(rad)
	dy := math.Cos(rad)

	// Compute min and max projections across all four corners (0,0), (w-1,0), (0,h-1), (w-1,h-1)
	// to determine the full range of the gradient axis.
	// Using float64 corner coordinates for the pixel centers.
	fw := float64(w - 1)
	fh := float64(h - 1)

	// Project all four corners
	p00 := 0.0           // (0,0) dot (dx,dy) = 0
	p10 := fw * dx       // (w-1,0)
	p01 := fh * dy       // (0,h-1)
	p11 := fw*dx + fh*dy // (w-1,h-1)

	// Find min and max projection values
	minProj := p00
	maxProj := p00
	for _, p := range [3]float64{p10, p01, p11} {
		if p < minProj {
			minProj = p
		}
		if p > maxProj {
			maxProj = p
		}
	}

	// Compute range; if zero (degenerate case: 1px dimension along axis), all pixels get t=0
	projRange := maxProj - minProj

	// Iterate over all pixels and compute their gradient position
	pix := img.Pix
	stride := img.Stride

	for y := 0; y < h; y++ {
		fy := float64(y)
		rowBase := y * stride
		for x := 0; x < w; x++ {
			fx := float64(x)
			proj := fx*dx + fy*dy

			var t float64
			if projRange > 0 {
				t = (proj - minProj) / projRange
			}

			c := interpolateColor(stops, t)

			off := rowBase + x*4
			pix[off+0] = c.R
			pix[off+1] = c.G
			pix[off+2] = c.B
			pix[off+3] = c.A
		}
	}

	return img
}

// renderRadial produces a radial gradient image where color transitions outward
// from the center in concentric circles. The gradient radius is the inscribed circle
// radius: min(width, height) / 2. Pixels outside this radius get the last stop color.
func renderRadial(bounds image.Rectangle, stops []ColorStop) *image.RGBA {
	w := bounds.Dx()
	h := bounds.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Center via integer division.
	cx := w / 2
	cy := h / 2

	// Inscribed circle radius.
	radius := min(w, h) / 2
	radiusF := float64(radius)

	// Use float center for accurate Euclidean distance
	cxF := float64(cx)
	cyF := float64(cy)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cxF
			dy := float64(y) - cyF
			dist := math.Sqrt(dx*dx + dy*dy)

			// Normalize distance by radius; pixels beyond the radius clamp to 1.0.
			t := dist / radiusF
			if t > 1.0 {
				t = 1.0
			}

			c := interpolateColor(stops, t)
			img.SetRGBA(x, y, c)
		}
	}

	return img
}

// interpolateColor finds the two bracketing stops for the given normalized position t
// and linearly interpolates each RGBA channel independently.
// Assumes stops are already sorted by position (ascending) and clamped to [0,1].
// For t at or beyond a stop with duplicate positions, uses the last stop in original order.
func interpolateColor(stops []ColorStop, t float64) color.RGBA {
	n := len(stops)

	// For t at or before the first stop, use the last stop sharing that position
	// (handles duplicate positions at the start).
	if t <= stops[0].Position {
		last := 0
		for i := 1; i < n; i++ {
			if stops[i].Position == stops[0].Position {
				last = i
			} else {
				break
			}
		}
		return stops[last].Color
	}

	// For t >= last stop position: return last stop color
	if t >= stops[n-1].Position {
		return stops[n-1].Color
	}

	// Find the two bracketing stops for t.
	// We want the last stop with Position <= t as the left bracket,
	// and the first stop with Position > t as the right bracket.
	// This naturally handles duplicates: if t equals a stop position,
	// we advance past all stops at that position to use the last one.
	left := 0
	for i := 1; i < n; i++ {
		if stops[i].Position <= t {
			left = i
		} else {
			break
		}
	}

	right := left + 1

	// If left and right have the same position (shouldn't happen given the loop above,
	// but guard against it), just return the left color
	if stops[left].Position == stops[right].Position {
		// Use last stop at this position — which is `right` since sorted stably
		return stops[right].Color
	}

	// Compute interpolation fraction between the two bracketing stops
	fraction := (t - stops[left].Position) / (stops[right].Position - stops[left].Position)

	// Linearly interpolate each RGBA channel and round to nearest uint8
	r := math.Round(float64(stops[left].Color.R) + (float64(stops[right].Color.R)-float64(stops[left].Color.R))*fraction)
	g := math.Round(float64(stops[left].Color.G) + (float64(stops[right].Color.G)-float64(stops[left].Color.G))*fraction)
	b := math.Round(float64(stops[left].Color.B) + (float64(stops[right].Color.B)-float64(stops[left].Color.B))*fraction)
	a := math.Round(float64(stops[left].Color.A) + (float64(stops[right].Color.A)-float64(stops[left].Color.A))*fraction)

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}
