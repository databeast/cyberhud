package scale

import (
	"image"
	"image/color"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

// For any source image.Image with positive dimensions and for any destination
// image.Rectangle with positive width and height that intersects the surface bounds,
// calling NearestNeighbor SHALL produce a result where each destination pixel
// (dx, dy) contains the color of the source pixel at
// (srcMinX + dx * srcW / dstW, srcMinY + dy * srcH / dstH) (integer division, nearest-neighbor mapping).

// scaleInput holds a randomly generated test input for the nearest-neighbor scaling property test.
type scaleInput struct {
	SrcW, SrcH int   // Source image dimensions (1..32)
	DstW, DstH int   // Destination dimensions (1..64)
	Seed       int64 // RNG seed for reproducible image generation
}

// Generate implements quick.Generator for property-based testing.
func (scaleInput) Generate(r *rand.Rand, size int) reflect.Value {
	input := scaleInput{
		SrcW: 1 + r.Intn(32),
		SrcH: 1 + r.Intn(32),
		DstW: 1 + r.Intn(64),
		DstH: 1 + r.Intn(64),
		Seed: r.Int63(),
	}
	return reflect.ValueOf(input)
}

func TestPropertyNearestNeighborCorrectness(t *testing.T) {
	config := &quick.Config{MaxCount: 200}

	prop := func(input scaleInput) bool {
		// Generate a random RGBA source image.
		rng := rand.New(rand.NewSource(input.Seed))
		src := image.NewRGBA(image.Rect(0, 0, input.SrcW, input.SrcH))
		for y := 0; y < input.SrcH; y++ {
			for x := 0; x < input.SrcW; x++ {
				src.SetRGBA(x, y, color.RGBA{
					R: uint8(rng.Intn(256)),
					G: uint8(rng.Intn(256)),
					B: uint8(rng.Intn(256)),
					A: uint8(rng.Intn(256)),
				})
			}
		}

		// Call NearestNeighbor directly.
		result := NearestNeighbor(src, input.DstW, input.DstH)

		// Verify the result dimensions.
		if result.Bounds().Dx() != input.DstW || result.Bounds().Dy() != input.DstH {
			t.Errorf("result dimensions mismatch: got %dx%d, want %dx%d",
				result.Bounds().Dx(), result.Bounds().Dy(), input.DstW, input.DstH)
			return false
		}

		// Verify each pixel matches the expected nearest-neighbor mapping.
		srcBounds := src.Bounds()
		srcMinX := srcBounds.Min.X
		srcMinY := srcBounds.Min.Y
		srcW := srcBounds.Dx()
		srcH := srcBounds.Dy()

		for dy := 0; dy < input.DstH; dy++ {
			for dx := 0; dx < input.DstW; dx++ {
				// Compute expected source coordinate using integer division.
				expectedSrcX := srcMinX + dx*srcW/input.DstW
				expectedSrcY := srcMinY + dy*srcH/input.DstH

				// Get the expected color from the source at that coordinate.
				expectedColor := src.RGBAAt(expectedSrcX, expectedSrcY)

				// Get the actual color from the scaled result.
				gotColor := result.RGBAAt(dx, dy)

				if gotColor != expectedColor {
					t.Errorf("pixel (%d,%d): got %v, want %v (mapped from src(%d,%d), srcSize=%dx%d, dstSize=%dx%d)",
						dx, dy, gotColor, expectedColor,
						expectedSrcX, expectedSrcY,
						input.SrcW, input.SrcH, input.DstW, input.DstH)
					return false
				}
			}
		}

		return true
	}

	if err := quick.Check(prop, config); err != nil {
		t.Errorf("Property 3 failed: %v", err)
	}
}

func TestNearestNeighborNilSource(t *testing.T) {
	result := NearestNeighbor(nil, 10, 10)
	if result != nil {
		t.Error("expected nil result for nil source")
	}
}

func TestNearestNeighborZeroDimensions(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 4))

	tests := []struct {
		name string
		w, h int
	}{
		{"zero width", 0, 10},
		{"zero height", 10, 0},
		{"negative width", -1, 10},
		{"negative height", 10, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NearestNeighbor(src, tt.w, tt.h)
			if result != nil {
				t.Errorf("expected nil for %s, got non-nil", tt.name)
			}
		})
	}
}

func TestNearestNeighborIdentity(t *testing.T) {
	// Scaling to the same size should produce identical pixels.
	src := image.NewRGBA(image.Rect(0, 0, 4, 4))
	src.SetRGBA(0, 0, color.RGBA{R: 255, A: 255})
	src.SetRGBA(3, 3, color.RGBA{B: 255, A: 255})

	result := NearestNeighbor(src, 4, 4)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			got := result.RGBAAt(x, y)
			want := src.RGBAAt(x, y)
			if got != want {
				t.Errorf("pixel (%d,%d): got %v, want %v", x, y, got, want)
			}
		}
	}
}
