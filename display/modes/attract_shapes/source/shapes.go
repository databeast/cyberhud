package source

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/databeast/cyberhud/display/widgets"
)

// shape represents a single geometric shape in the animation.
type shape struct {
	x             float64 // center x position (fraction of panel width)
	y             float64 // center y position (fraction of panel height)
	baseSize      float64 // base radius in pixels
	rotationAngle float64 // current rotation angle in radians
	currentScale  float64 // current scale factor (oscillates 0.5-1.5)
	hue           float64 // color hue [0, 360)
	sides         int     // number of polygon sides
}

// Package-level animation state.
var (
	shapes            []shape
	shapesInitialized bool
	shapesWidth       int
	shapesHeight      int
	shapesCount       int
	animPhase         float64
)

// initShapes initializes or re-initializes the shape array when panel dimensions
// or shape count changes.
func InitShapes(p Policy, panelWidth, panelHeight int) {
	count := p.ShapeCount
	if count < 1 {
		count = 1
	}
	if count > 50 {
		count = 50
	}

	if shapesInitialized && shapesWidth == panelWidth && shapesHeight == panelHeight && shapesCount == count {
		return
	}

	shapes = make([]shape, count)
	minDim := panelWidth
	if panelHeight < minDim {
		minDim = panelHeight
	}
	baseRadius := float64(minDim) * 0.08 // base shape size is 8% of shortest dimension

	sides := p.Complexity
	if sides < 3 {
		sides = 3
	}
	if sides > 8 {
		sides = 8
	}

	for i := range shapes {
		shapes[i] = shape{
			x:             rand.Float64(),
			y:             rand.Float64(),
			baseSize:      baseRadius * (0.5 + rand.Float64()),
			rotationAngle: rand.Float64() * 2 * math.Pi,
			currentScale:  0.5 + rand.Float64(),
			hue:           rand.Float64() * 360.0,
			sides:         sides,
		}
	}

	shapesInitialized = true
	shapesWidth = panelWidth
	shapesHeight = panelHeight
	shapesCount = count
}

// advanceShapes updates rotation and scale for each shape based on elapsed time.
func AdvanceShapes(elapsedSec float64, p Policy) {
	animPhase += elapsedSec

	// Rotation: Speed degrees per second, converted to radians.
	rotationDelta := elapsedSec * (math.Pi / 180.0) * 90.0 // 90 deg/sec base × speed factor

	for i := range shapes {
		// Advance rotation.
		shapes[i].rotationAngle += rotationDelta

		// Oscillate scale between 0.5 and 1.5 using sine wave at PulseRate.
		// PulseRate is oscillations per second.
		shapes[i].currentScale = 1.0 + 0.5*math.Sin(animPhase*p.PulseRate*2*math.Pi+float64(i)*0.7)
	}
}

// buildSprites generates geometric shape sprites as an image.RGBA sprite.
func BuildSprites(p Policy, panelWidth, panelHeight int, mono bool) []widgets.Sprite {
	if panelWidth <= 0 || panelHeight <= 0 {
		return []widgets.Sprite{{Image: image.NewRGBA(image.Rect(0, 0, 1, 1)), Position: image.Point{}}}
	}

	img := image.NewRGBA(image.Rect(0, 0, panelWidth, panelHeight))

	sides := p.Complexity
	if sides < 3 {
		sides = 3
	}
	if sides > 8 {
		sides = 8
	}

	for _, s := range shapes {
		cx := s.x * float64(panelWidth)
		cy := s.y * float64(panelHeight)
		radius := s.baseSize * s.currentScale

		if mono {
			// Monochrome: single-color outlines only.
			DrawPolygonOutline(img, cx, cy, radius, s.sides, s.rotationAngle, color.RGBA{255, 255, 255, 255})
		} else {
			// Color: filled polygons with color variation.
			fillColor := hueToRGBA(s.hue, 200)
			drawPolygonFilled(img, cx, cy, radius, s.sides, s.rotationAngle, fillColor, panelWidth, panelHeight)
		}
	}

	return []widgets.Sprite{
		{Image: img, Position: image.Point{X: 0, Y: 0}, Label: "shapes-frame"},
	}
}

// drawPolygonOutline draws the outline of a regular polygon.
func DrawPolygonOutline(img *image.RGBA, cx, cy, radius float64, sides int, angle float64, c color.RGBA) {
	bounds := img.Bounds()
	vertices := computeVertices(cx, cy, radius, sides, angle)

	for i := 0; i < len(vertices); i++ {
		v1 := vertices[i]
		v2 := vertices[(i+1)%len(vertices)]
		drawLine(img, v1.X, v1.Y, v2.X, v2.Y, c, bounds)
	}
}

// drawPolygonFilled draws a filled regular polygon using scanline fill.
func drawPolygonFilled(img *image.RGBA, cx, cy, radius float64, sides int, angle float64, c color.RGBA, panelWidth, panelHeight int) {
	vertices := computeVertices(cx, cy, radius, sides, angle)

	// Find bounding box.
	minY, maxY := vertices[0].Y, vertices[0].Y
	for _, v := range vertices {
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}

	// Scanline fill.
	startY := int(math.Floor(minY))
	endY := int(math.Ceil(maxY))
	if startY < 0 {
		startY = 0
	}
	if endY >= panelHeight {
		endY = panelHeight - 1
	}

	for y := startY; y <= endY; y++ {
		fy := float64(y) + 0.5
		var intersections []float64

		n := len(vertices)
		for i := 0; i < n; i++ {
			v1 := vertices[i]
			v2 := vertices[(i+1)%n]

			if (v1.Y <= fy && v2.Y > fy) || (v2.Y <= fy && v1.Y > fy) {
				t := (fy - v1.Y) / (v2.Y - v1.Y)
				ix := v1.X + t*(v2.X-v1.X)
				intersections = append(intersections, ix)
			}
		}

		// Sort intersections.
		for i := 0; i < len(intersections)-1; i++ {
			for j := i + 1; j < len(intersections); j++ {
				if intersections[j] < intersections[i] {
					intersections[i], intersections[j] = intersections[j], intersections[i]
				}
			}
		}

		// Fill between pairs.
		for i := 0; i+1 < len(intersections); i += 2 {
			startX := int(math.Ceil(intersections[i]))
			endX := int(math.Floor(intersections[i+1]))
			if startX < 0 {
				startX = 0
			}
			if endX >= panelWidth {
				endX = panelWidth - 1
			}
			for x := startX; x <= endX; x++ {
				img.SetRGBA(x, y, c)
			}
		}
	}

	// Draw outline on top for definition.
	outlineColor := color.RGBA{
		R: uint8(min(int(c.R)+40, 255)),
		G: uint8(min(int(c.G)+40, 255)),
		B: uint8(min(int(c.B)+40, 255)),
		A: 255,
	}
	DrawPolygonOutline(img, cx, cy, radius, sides, angle, outlineColor)
}

// vertex represents a 2D point.
type vertex struct {
	X, Y float64
}

// computeVertices computes the vertices of a regular polygon.
func computeVertices(cx, cy, radius float64, sides int, angle float64) []vertex {
	verts := make([]vertex, sides)
	for i := 0; i < sides; i++ {
		theta := angle + float64(i)*2.0*math.Pi/float64(sides)
		verts[i] = vertex{
			X: cx + radius*math.Cos(theta),
			Y: cy + radius*math.Sin(theta),
		}
	}
	return verts
}

// drawLine draws a line between two points using Bresenham's algorithm.
func drawLine(img *image.RGBA, x0, y0, x1, y1 float64, c color.RGBA, bounds image.Rectangle) {
	ix0, iy0 := int(math.Round(x0)), int(math.Round(y0))
	ix1, iy1 := int(math.Round(x1)), int(math.Round(y1))

	dx := abs(ix1 - ix0)
	dy := -abs(iy1 - iy0)
	sx := 1
	if ix0 > ix1 {
		sx = -1
	}
	sy := 1
	if iy0 > iy1 {
		sy = -1
	}
	err := dx + dy

	for {
		if ix0 >= bounds.Min.X && ix0 < bounds.Max.X && iy0 >= bounds.Min.Y && iy0 < bounds.Max.Y {
			img.SetRGBA(ix0, iy0, c)
		}
		if ix0 == ix1 && iy0 == iy1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			ix0 += sx
		}
		if e2 <= dx {
			err += dx
			iy0 += sy
		}
	}
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// hueToRGBA converts a hue value [0, 360) to an RGBA color with the given alpha.
func hueToRGBA(hue float64, alpha uint8) color.RGBA {
	h := math.Mod(hue, 360.0)
	if h < 0 {
		h += 360.0
	}
	c := 1.0
	x := 1.0 - math.Abs(math.Mod(h/60.0, 2.0)-1.0)

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: alpha,
	}
}
