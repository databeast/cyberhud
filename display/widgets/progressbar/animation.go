package progressbar

import (
	"image"
	"math"
	"time"

	"github.com/databeast/cyberhud/display/widgets"
)

// applyAnimation is a post-processing dispatcher called after the main render
// pass produces the image. It modifies the image in-place based on the
// configured animation type.
//
// When cfg.Animation.Type == NoAnimation, this is a no-op.
// When cfg.Animation.Type == Pulse, sinusoidal brightness modulation is applied.
// When cfg.Animation.Type == Shimmer, a translucent highlight band sweeps the fill.
// When cfg.Animation.Type == MarchingStripes, diagonal stripes scroll across the fill.
func applyAnimation(img *image.RGBA, cfg Config) {
	switch cfg.Animation.Type {
	case NoAnimation:
		return
	case Pulse:
		applyPulse(img, cfg)
	case Shimmer:
		applyShimmer(img, cfg)
	case MarchingStripes:
		applyMarchingStripes(img, cfg)
	}
}

// isFillPixel returns true if the pixel at (x, y) is a "fill" pixel — i.e.,
// not transparent and not background. This is the same classification used by
// applyPulse and all animation overlays.
func isFillPixel(img *image.RGBA, x, y int, bg [4]uint8) bool {
	c := img.RGBAAt(x, y)
	if c.A == 0 {
		return false
	}
	if c.R == bg[0] && c.G == bg[1] && c.B == bg[2] && c.A == bg[3] {
		return false
	}
	return true
}

// blendWhite50 applies a 50% white blend to the pixel at (x, y).
// R = (R + 255) / 2, G = (G + 255) / 2, B = (B + 255) / 2.
func blendWhite50(img *image.RGBA, x, y int) {
	c := img.RGBAAt(x, y)
	c.R = uint8((uint16(c.R) + 255) / 2)
	c.G = uint8((uint16(c.G) + 255) / 2)
	c.B = uint8((uint16(c.B) + 255) / 2)
	img.SetRGBA(x, y, c)
}

// applyShimmer applies a translucent highlight band sweeping across the fill.
//
// For Linear/Segmented bars: the band sweeps along the primary axis.
// Band width = 20% of the fill extent (pixels). Position advances at
// cfg.Animation.Speed pixels per second, wrapping end-to-start.
//
// For Ring/Arc bars: the band sweeps along the angular extent of the fill.
// Band width = 20% of the filled arc angle.
//
// Only fill pixels (non-transparent, non-background) within the band are
// highlighted with 50% white blend.
func applyShimmer(img *image.RGBA, cfg Config) {
	speed := cfg.Animation.Speed
	if speed <= 0 {
		return // disabled
	}

	_, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)
	bgArr := [4]uint8{bg.R, bg.G, bg.B, bg.A}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	geom := resolveOrientation(cfg)

	switch cfg.Style {
	case Ring, Arc:
		applyShimmerCircular(img, cfg, geom, bgArr)
	default:
		// Linear, Segmented, Pie (though Pie is unlikely to use shimmer)
		applyShimmerLinear(img, cfg, geom, bgArr, w, h, speed)
	}
}

// applyShimmerLinear handles shimmer for Linear and Segmented bars.
func applyShimmerLinear(img *image.RGBA, cfg Config, geom RenderGeometry, bgArr [4]uint8, w, h, speed int) {
	primaryAxis := geom.PrimaryAxis
	if primaryAxis <= 0 {
		return
	}

	// Fill extent: how many pixels along primary axis are filled.
	fillExtent := int(float64(primaryAxis) * cfg.Value)
	if fillExtent <= 0 {
		return
	}

	// Band width = 20% of fill extent, minimum 1 pixel.
	bandWidth := fillExtent * 20 / 100
	if bandWidth < 1 {
		bandWidth = 1
	}

	// Compute offset from elapsed time.
	offset := int(float64(speed) * float64(cfg.animElapsed) / float64(time.Second))
	bandStart := offset % primaryAxis

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !isFillPixel(img, x, y, bgArr) {
				continue
			}

			// Determine position along primary axis.
			var pos int
			if geom.FillDirection == FillBottomToTop {
				// Vertical: primary axis is Y from bottom.
				pos = (h - 1) - y
			} else {
				// Horizontal: primary axis is X from left.
				pos = x
			}

			// Check if this position is within the shimmer band (wrapping).
			inBand := ((pos - bandStart + primaryAxis) % primaryAxis) < bandWidth
			if inBand {
				blendWhite50(img, x, y)
			}
		}
	}
}

// applyShimmerCircular handles shimmer for Ring and Arc bars.
// The band sweeps along the angular extent of the fill.
// Band width = 20% of the filled arc angle.
func applyShimmerCircular(img *image.RGBA, cfg Config, geom RenderGeometry, bgArr [4]uint8) {
	speed := cfg.Animation.Speed
	if speed <= 0 {
		return
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	// Determine the total angular sweep of the fill.
	var sweepRad float64
	var arcStart float64

	if cfg.Style == Ring {
		sweepRad = 2 * math.Pi
		arcStart = geom.StartAngle
	} else {
		// Arc
		sweepRad = cfg.SweepAngle * math.Pi / 180.0
		arcStart = geom.StartAngle - sweepRad/2.0
	}

	// Filled arc angle.
	fillAngle := cfg.Value * sweepRad
	if fillAngle <= 0 {
		return
	}

	// Band angular width = 20% of filled arc angle.
	bandAngle := fillAngle * 0.20
	if bandAngle <= 0 {
		return
	}

	// Use a virtual "linear" position concept: the band sweeps at speed px/sec
	// mapped to the angular extent. We use the circumference at the midpoint
	// radius as the reference for mapping px/sec to radians/sec.
	outerR := math.Min(float64(w), float64(h)) / 2.0
	innerR := outerR - float64(cfg.Thickness)
	if innerR < 0 {
		innerR = 0
	}
	midR := (outerR + innerR) / 2.0
	if midR <= 0 {
		midR = 1
	}

	// Map speed (px/sec) to angular speed (rad/sec).
	angularSpeed := float64(speed) / midR

	// Current angular offset of the band start.
	elapsed := float64(cfg.animElapsed) / float64(time.Second)
	bandOffset := math.Mod(angularSpeed*elapsed, fillAngle)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !isFillPixel(img, x, y, bgArr) {
				continue
			}

			// Compute angle for this pixel.
			dx := float64(x) + 0.5 - cx
			dy := float64(y) + 0.5 - cy

			// Angle clockwise from 12-o'clock.
			angle := math.Atan2(dx, -dy)

			// Compute relative angle from arc start.
			var relAngle float64
			if cfg.Style == Ring {
				// For ring: compute relative to start offset.
				var startOffset float64
				if cfg.Orientation == OrientVertical {
					startOffset = 3.0 * math.Pi / 2.0
				}
				// Normalize angle to [0, 2π).
				if angle < 0 {
					angle += 2 * math.Pi
				}
				relAngle = angle - startOffset
				if relAngle < 0 {
					relAngle += 2 * math.Pi
				}
			} else {
				// Arc: relative to arcStart.
				relAngle = angle - arcStart
				relAngle = math.Mod(relAngle, 2*math.Pi)
				if relAngle < 0 {
					relAngle += 2 * math.Pi
				}
			}

			// Only apply to fill pixels within the filled angular extent.
			if relAngle >= fillAngle {
				continue
			}

			// Check if within the shimmer band (wrapping within fill extent).
			posInFill := relAngle - bandOffset
			posInFill = math.Mod(posInFill, fillAngle)
			if posInFill < 0 {
				posInFill += fillAngle
			}

			if posInFill < bandAngle {
				blendWhite50(img, x, y)
			}
		}
	}
}

// applyMarchingStripes applies diagonal stripe scrolling across fill pixels.
//
// For Linear/Segmented: 45° diagonal stripes (4px wide, 4px gap) scroll at
// configured speed. Pattern: ((x - y + offset) mod 8) < 4 → stripe pixel.
//
// For Ring/Arc: MarchingStripes is NOT supported — falls back to Shimmer.
//
// Only fill pixels (non-transparent, non-background) receive the stripe overlay.
func applyMarchingStripes(img *image.RGBA, cfg Config) {
	// Ring/Arc: fall back to Shimmer.
	if cfg.Style == Ring || cfg.Style == Arc {
		applyShimmer(img, cfg)
		return
	}

	speed := cfg.Animation.Speed
	if speed <= 0 {
		return // disabled
	}

	_, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)
	bgArr := [4]uint8{bg.R, bg.G, bg.B, bg.A}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	// Compute offset from elapsed time.
	offset := int(float64(speed) * float64(cfg.animElapsed) / float64(time.Second))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !isFillPixel(img, x, y, bgArr) {
				continue
			}

			// 45° diagonal stripe pattern: ((x - y + offset) mod 8) < 4
			stripe := ((x-y+offset)%8 + 8) % 8 // ensure positive modulo
			if stripe < 4 {
				blendWhite50(img, x, y)
			}
		}
	}
}

// applyPulse applies sinusoidal brightness modulation to all non-background,
// non-transparent pixels in the image.
//
// The brightness factor oscillates between 0.30 and 1.00:
//
//	phase = sin(2π × elapsed / period)
//	factor = 0.65 + 0.35 × phase   → range [0.30, 1.00]
//
// Pixels are classified as "fill" if they are neither transparent (alpha == 0)
// nor equal to the resolved background color. Only fill pixels are modulated.
// This correctly handles:
//   - Linear style: all fill region pixels are modulated
//   - Ring/Arc styles: all filled arc pixels are modulated (track is background)
//   - Segmented style: filled cells are modulated; gaps and unfilled cells are background
func applyPulse(img *image.RGBA, cfg Config) {
	period := cfg.Animation.Period
	if period <= 0 {
		return // disabled
	}

	elapsed := cfg.animElapsed
	phase := math.Sin(2.0 * math.Pi * float64(elapsed) / float64(period))
	factor := 0.65 + 0.35*phase // [0.30, 1.00]

	_, bg := widgets.ResolveColors(cfg.Foreground, cfg.Background)

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.RGBAAt(x, y)
			// Skip transparent pixels (e.g., from rounded caps masking).
			if c.A == 0 {
				continue
			}
			// Skip background/track pixels.
			if c == bg {
				continue
			}
			// Modulate RGB channels by factor.
			c.R = uint8(math.Round(float64(c.R) * factor))
			c.G = uint8(math.Round(float64(c.G) * factor))
			c.B = uint8(math.Round(float64(c.B) * factor))
			img.SetRGBA(x, y, c)
		}
	}
}
