package source

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"testing"
)

func TestRotationPreservesCubeAndStartsOnPoint(t *testing.T) {
	for _, seconds := range []float64{0, 1, 4, 20} {
		for _, edge := range edges {
			a, b := rotate(vertices[edge[0]], seconds), rotate(vertices[edge[1]], seconds)
			length := math.Sqrt((a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y) + (a.z-b.z)*(a.z-b.z))
			if math.Abs(length-2) > 1e-12 {
				t.Fatalf("edge length at %g seconds = %g, want 2", seconds, length)
			}
		}
	}
	bottom := rotate(vertices[0], 0)
	top := rotate(vertices[7], 0)
	if math.Abs(bottom.x) > 1e-12 || math.Abs(bottom.z) > 1e-12 ||
		math.Abs(top.x) > 1e-12 || math.Abs(top.z) > 1e-12 || bottom.y >= 0 || top.y <= 0 {
		t.Fatalf("body diagonal is not vertical: bottom=%+v top=%+v", bottom, top)
	}
}

func TestRenderFrameDeterministicAndAnimated(t *testing.T) {
	for _, size := range []image.Point{{240, 240}, {320, 240}, {480, 320}} {
		cyan := color.RGBA{G: 255, B: 255, A: 255}
		initial := RenderFrame(size.X, size.Y, 0, cyan)
		repeated := RenderFrame(size.X, size.Y, 0, cyan)
		later := RenderFrame(size.X, size.Y, 4, cyan)
		if initial.Bounds() != image.Rect(0, 0, size.X, size.Y) {
			t.Fatalf("unexpected bounds %v", initial.Bounds())
		}

		if !bytes.Equal(initial.Pix, repeated.Pix) {
			t.Fatal("same time produced different frames")
		}
		if bytes.Equal(initial.Pix, later.Pix) {
			t.Fatal("cube did not animate")
		}
		background := initial.RGBAAt(0, 0)
		changed := 0
		for y := 0; y < size.Y; y++ {
			for x := 0; x < size.X; x++ {
				if initial.RGBAAt(x, y) != background {
					changed++
				}
			}

		}
		if changed < size.X*size.Y/100 {
			t.Fatalf("cube covers only %d pixels at %v", changed, size)
		}
	}
}

func TestTumbleStaysWithinFrame(t *testing.T) {
	// Maximize 5*x/(5-z) subject to x*x+z*z <= 3.
	maxProjected := cubeScale * 5 * math.Sqrt(3) / math.Sqrt(25-3)
	if maxProjected+cubeScale*0.14/2 >= 0.5 {
		t.Fatalf("cube or glow can reach the edge of the frame: half-size %g", maxProjected)
	}
}
