package source

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/widgets"
)

func RenderBokehFrame(x, y int) {

}

// buildEinkView produces a static bokeh scatter frame for e-ink panels.
// The frame is frozen (Static=true) and the cache key is stable across
// calls unless the policy changes.
func RenderEInkBokehPicture(panelWidth, panelHeight int, p Policy) style.ViewData {
	img := renderEinkFrame(panelWidth, panelHeight, p)

	return style.ViewData{
		Static: true,
		Sprites: []widgets.Sprite{
			{
				Image:    img,
				Position: image.Point{},
				Bounds:   image.Rect(0, 0, panelWidth, panelHeight),
				Label:    "bokeh-eink",
			},
		},
		StyleReport: style.StyleReport{Name: "eink-bokeh", Reason: "fitness"},
	}
}

// represents a single bokeh light circle in the animation.
type circle struct {
	x       float64 // center x position
	y       float64 // center y position
	radius  float64 // circle radius in pixels
	dx      float64 // horizontal velocity (pixels/sec at speed=1)
	dy      float64 // vertical velocity (pixels/sec at speed=1)
	hue     float64 // hue angle [0, 360)
	opacity float64 // peak opacity [0.15, 0.6]
}

// AdvanceCircles moves all circles by the given elapsed time scaled by speed.
// Circles that fully exit bounds re-enter on the opposite edge.
func AdvanceCircles(elapsed time.Duration, speed float64, panelWidth, panelHeight int) {
	dt := elapsed.Seconds() * speed
	for i := range circles {
		circles[i].x += circles[i].dx * dt
		circles[i].y += circles[i].dy * dt

		// Check if circle fully exited bounds and wrap to opposite edge.
		r := circles[i].radius
		w := float64(panelWidth)
		h := float64(panelHeight)

		if circles[i].x-r > w {
			circles[i].x = -r
			circles[i].y = rand.Float64() * h
		} else if circles[i].x+r < 0 {
			circles[i].x = w + r
			circles[i].y = rand.Float64() * h
		}

		if circles[i].y-r > h {
			circles[i].y = -r
			circles[i].x = rand.Float64() * w
		} else if circles[i].y+r < 0 {
			circles[i].y = h + r
			circles[i].x = rand.Float64() * w
		}
	}
}

// RenderCircles produces an image.RGBA with all circles drawn as radial gradients.
func RenderCircles(panelWidth, panelHeight int, p Policy, mono bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))

	for _, c := range circles {
		drawCircle(img, c, p, mono)
	}

	return img
}

// drawCircle renders a single circle with a radial gradient (soft edges).
func drawCircle(img *image.RGBA, c circle, p Policy, mono bool) {
	bounds := img.Bounds()
	r := c.radius
	if r < 1 {
		r = 1
	}

	// Bounding box for the circle.
	minX := int(math.Floor(c.x - r))
	maxX := int(math.Ceil(c.x + r))
	minY := int(math.Floor(c.y - r))
	maxY := int(math.Ceil(c.y + r))

	// Clamp to image bounds.
	if minX < bounds.Min.X {
		minX = bounds.Min.X
	}
	if maxX > bounds.Max.X {
		maxX = bounds.Max.X
	}
	if minY < bounds.Min.Y {
		minY = bounds.Min.Y
	}
	if maxY > bounds.Max.Y {
		maxY = bounds.Max.Y
	}

	for py := minY; py < maxY; py++ {
		for px := minX; px < maxX; px++ {
			dx := float64(px) + 0.5 - c.x
			dy := float64(py) + 0.5 - c.y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist >= r {
				continue
			}

			// Radial gradient: full opacity at center, zero at edge.
			alpha := (1.0 - dist/r) * c.opacity

			var col color.RGBA
			if mono {
				// Monochrome: luminance only.
				lum := uint8(alpha * 255.0)
				col = color.RGBA{R: lum, G: lum, B: lum, A: lum}
			} else {
				// Color: use hue with saturation from policy.
				cr, cg, cb := hslToRGB(c.hue, p.Saturation, 0.6)
				col = color.RGBA{
					R: uint8(float64(cr) * alpha),
					G: uint8(float64(cg) * alpha),
					B: uint8(float64(cb) * alpha),
					A: uint8(alpha * 255.0),
				}
			}

			// Alpha-blend over existing pixel.
			existing := img.RGBAAt(px, py)
			blended := alphaBlend(existing, col)
			img.SetRGBA(px, py, blended)
		}
	}
}

// alphaBlend blends src over dst using premultiplied alpha.
func alphaBlend(dst, src color.RGBA) color.RGBA {
	sa := uint32(src.A)
	da := uint32(dst.A)
	outA := sa + da*(255-sa)/255

	if outA == 0 {
		return color.RGBA{}
	}

	outR := (uint32(src.R)*255 + uint32(dst.R)*(255-sa)) / 255
	outG := (uint32(src.G)*255 + uint32(dst.G)*(255-sa)) / 255
	outB := (uint32(src.B)*255 + uint32(dst.B)*(255-sa)) / 255

	return color.RGBA{
		R: uint8(outR),
		G: uint8(outG),
		B: uint8(outB),
		A: uint8(outA),
	}
}

// hslToRGB converts HSL (hue 0-360, saturation 0-1, lightness 0-1) to RGB (0-255).
func hslToRGB(h, s, l float64) (r, g, b uint8) {
	if s == 0 {
		v := uint8(l * 255.0)
		return v, v, v
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	hNorm := h / 360.0
	rF := hueToRGB(p, q, hNorm+1.0/3.0)
	gF := hueToRGB(p, q, hNorm)
	bF := hueToRGB(p, q, hNorm-1.0/3.0)

	return uint8(rF * 255.0), uint8(gF * 255.0), uint8(bF * 255.0)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}

// CircleCount computes the number of circles for the given panel area and density.
// Formula: max(1, floor(area/4096 * density))
func CircleCount(panelWidth, panelHeight int, density float64) int {
	area := panelWidth * panelHeight
	count := int(math.Floor(float64(area) / 4096.0 * density))
	if count < 1 {
		count = 1
	}
	return count
}

// circles holds the current bokeh circle state.
var circles []circle

// circlesInitialized tracks whether circles have been initialized for the current panel.
var circlesInitialized bool
var circlesPanelW, circlesPanelH int

// InitCircles initializes the circle slice for the given panel dimensions and policy.
func InitCircles(panelWidth, panelHeight int, p Policy) {
	count := CircleCount(panelWidth, panelHeight, p.Density)
	shortDim := panelWidth
	if panelHeight < shortDim {
		shortDim = panelHeight
	}

	circles = make([]circle, count)
	for i := range circles {
		circles[i] = randomCircle(panelWidth, panelHeight, shortDim, p)
	}
	circlesInitialized = true
	circlesPanelW = panelWidth
	circlesPanelH = panelHeight
}

// EnsureCircles initializes or refreshes circles for the current panel and policy.
func EnsureCircles(panelWidth, panelHeight int, p Policy) {
	expectedCount := CircleCount(panelWidth, panelHeight, p.Density)
	if !circlesInitialized || circlesPanelW != panelWidth || circlesPanelH != panelHeight || len(circles) != expectedCount {
		InitCircles(panelWidth, panelHeight, p)
	}
}

// ResetCircles clears the animated circle state.
func ResetCircles() {
	circles = nil
	circlesInitialized = false
	circlesPanelW = 0
	circlesPanelH = 0
}

// randomCircle creates a circle with random position, size, velocity, and color.
func randomCircle(panelWidth, panelHeight, shortDim int, p Policy) circle {
	// Radius between 2% and 15% of shortest dimension, varied by SizeVariance.
	minRadius := 0.02 * float64(shortDim)
	maxRadius := 0.15 * float64(shortDim)
	// SizeVariance controls the spread: 0 = all at midpoint, 1 = full range.
	midRadius := (minRadius + maxRadius) / 2.0
	spread := (maxRadius - minRadius) / 2.0
	radius := midRadius + (rand.Float64()*2.0-1.0)*spread*p.SizeVariance

	// Clamp radius to valid range.
	if radius < minRadius {
		radius = minRadius
	}
	if radius > maxRadius {
		radius = maxRadius
	}

	// Random position within panel bounds.
	x := rand.Float64() * float64(panelWidth)
	y := rand.Float64() * float64(panelHeight)

	// Random velocity — base speed of ~20-60 pixels/sec at speed=1.
	angle := rand.Float64() * 2.0 * math.Pi
	speed := 20.0 + rand.Float64()*40.0
	dx := math.Cos(angle) * speed
	dy := math.Sin(angle) * speed

	// Random hue.
	hue := rand.Float64() * 360.0

	// Peak opacity between 0.15 and 0.6.
	opacity := 0.15 + rand.Float64()*0.45

	return circle{
		x:       x,
		y:       y,
		radius:  radius,
		dx:      dx,
		dy:      dy,
		hue:     hue,
		opacity: opacity,
	}
}
