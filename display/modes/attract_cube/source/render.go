// Package source renders the cube attract mode.
package source

import (
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
)

type vec3 struct{ x, y, z float64 }
type point struct{ x, y float64 }

const cubeScale = 0.24

var vertices = [8]vec3{
	{-1, -1, -1}, {1, -1, -1}, {-1, 1, -1}, {1, 1, -1},
	{-1, -1, 1}, {1, -1, 1}, {-1, 1, 1}, {1, 1, 1},
}

var faces = [6][4]int{
	{0, 1, 3, 2}, {4, 6, 7, 5},
	{0, 4, 5, 1}, {2, 3, 7, 6},
	{0, 2, 6, 4}, {1, 5, 7, 3},
}

var edges = [12][2]int{
	{0, 1}, {0, 2}, {0, 4}, {1, 3}, {1, 5}, {2, 3},
	{2, 6}, {3, 7}, {4, 5}, {4, 6}, {5, 7}, {6, 7},
}

// RenderFrame renders an opaque, antialiased color frame at the given time in seconds.
// Width and height must be positive, and seconds must be finite.
func RenderFrame(width, height int, seconds float64, base color.RGBA) *image.RGBA {
	return renderFrame(width, height, seconds, base, false)
}

// RenderGrayscaleFrame renders a neutral, high-contrast frame independent of the color policy.
func RenderGrayscaleFrame(width, height int, seconds float64) *image.RGBA {
	return renderFrame(width, height, seconds, color.RGBA{R: 255, G: 255, B: 255, A: 255}, true)
}

// RenderMonoFrame renders a one-bit wireframe without fills or antialiasing.
func RenderMonoFrame(width, height int, seconds float64) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{A: 255})
		}
	}
	_, projected := project(width, height, seconds)
	for _, edge := range edges {
		a, b := projected[edge[0]], projected[edge[1]]
		drawMonoLine(img, int(math.Round(a.x)), int(math.Round(a.y)), int(math.Round(b.x)), int(math.Round(b.y)))
	}
	return img
}

func drawMonoLine(img *image.RGBA, x0, y0, x1, y1 int) {
	dx, dy := absInt(x1-x0), absInt(y1-y0)
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy
	for {
		img.SetRGBA(x0, y0, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		if x0 == x1 && y0 == y1 {
			break
		}
		e := 2 * err
		if e > -dy {
			err -= dy
			x0 += sx
		}
		if e < dx {
			err += dx
			y0 += sy
		}
	}
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func renderFrame(width, height int, seconds float64, base color.RGBA, grayscale bool) *image.RGBA {
	const supersample = 2
	dc := gg.NewContext(width*supersample, height*supersample)
	dc.Scale(supersample, supersample)
	if grayscale {
		dc.SetRGB(0.02, 0.02, 0.02)
	} else {
		dc.SetRGB(0.008, 0.025, 0.039)
	}
	dc.Clear()

	rotated, projected := project(width, height, seconds)
	size := math.Min(float64(width), float64(height)) * cubeScale

	// Broad, low-opacity strokes form a soft halo beneath the translucent cube.
	for _, layer := range []struct{ width, alpha float64 }{
		{14, 0.012}, {9, 0.018}, {5, 0.035}, {2.5, 0.06},
	} {
		r, g, b := tinted(base, 0.16, 0.86, 1)
		dc.SetRGBA(r, g, b, layer.alpha)
		dc.SetLineWidth(math.Max(0.5, size*layer.width/100))
		strokeEdges(dc, projected)
	}

	type faceDepth struct {
		index int
		depth float64
	}
	ordered := make([]faceDepth, len(faces))
	for i, face := range faces {
		for _, vertex := range face {
			ordered[i].depth += rotated[vertex].z / 4
		}
		ordered[i].index = i
	}
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].depth < ordered[j].depth })
	for _, face := range ordered {
		indices := faces[face.index]
		dc.MoveTo(projected[indices[0]].x, projected[indices[0]].y)
		for _, i := range indices[1:] {
			dc.LineTo(projected[i].x, projected[i].y)
		}
		dc.ClosePath()
		// Face-dependent opacity keeps overlapping surfaces legible.
		r, g, b := tinted(base, 0.05, 0.75, 0.92)
		dc.SetRGBA(r, g, b, 0.105+0.025*float64(face.index%3))
		dc.Fill()
	}

	r, g, b := tinted(base, 0.24, 0.91, 1)
	dc.SetRGBA(r, g, b, 0.8)
	dc.SetLineWidth(math.Max(0.7, size*0.012))
	strokeEdges(dc, projected)
	r, g, b = tinted(base, 0.78, 1, 1)
	dc.SetRGBA(r, g, b, 0.7)
	dc.SetLineWidth(math.Max(0.3, size*0.004))
	strokeEdges(dc, projected)

	result := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(result, result.Bounds(), dc.Image(), dc.Image().Bounds(), draw.Src, nil)
	return result
}

func project(width, height int, seconds float64) (rotated [8]vec3, projected [8]point) {
	size := math.Min(float64(width), float64(height)) * cubeScale
	for i, v := range vertices {
		rotated[i] = rotate(v, seconds)
		perspective := 5.0 / (5.0 - rotated[i].z)
		projected[i] = point{
			float64(width)/2 + rotated[i].x*size*perspective,
			float64(height)/2 - rotated[i].y*size*perspective,
		}
	}
	return
}

func strokeEdges(dc *gg.Context, points [8]point) {
	for _, edge := range edges {
		dc.MoveTo(points[edge[0]].x, points[edge[0]].y)
		dc.LineTo(points[edge[1]].x, points[edge[1]].y)
	}
	dc.Stroke()
}

func rotate(v vec3, seconds float64) vec3 {
	// Start with the body diagonal vertical, so the cube rests on a corner.
	v = vec3{
		(v.x - v.z) / math.Sqrt2,
		(v.x + v.y + v.z) / math.Sqrt(3),
		(v.x - 2*v.y + v.z) / math.Sqrt(6),
	}
	sx, cx := math.Sincos(seconds * 0.21)
	sy, cy := math.Sincos(seconds * 0.37)
	sz, cz := math.Sincos(seconds * 0.13)
	v = vec3{v.x, v.y*cx - v.z*sx, v.y*sx + v.z*cx}
	v = vec3{v.x*cy + v.z*sy, v.y, -v.x*sy + v.z*cy}
	return vec3{v.x*cz - v.y*sz, v.x*sz + v.y*cz, v.z}
}
