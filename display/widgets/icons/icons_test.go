package icons

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"pgregory.net/rapid"
)

// --- From: borders_test.go ---

// For all border tile names in the set of 15 tiles, retrieving the tile from
// the registry SHALL return a non-nil *image.Alpha with bounds exactly equal
// to an 8×8 rectangle.

func TestProperty_BorderTiles8x8Alpha(t *testing.T) {
	borderTileNames := []string{
		"border/h", "border/v",
		"border/corner-tl", "border/corner-tr", "border/corner-bl", "border/corner-br",
		"border/tee-l", "border/tee-r", "border/tee-t", "border/tee-b",
		"border/cross",
		"border/round-tl", "border/round-tr", "border/round-bl", "border/round-br",
	}

	rapid.Check(t, func(t *rapid.T) {
		name := rapid.SampledFrom(borderTileNames).Draw(t, "tileName")

		img, ok := Get(name)
		if !ok {
			t.Fatalf("Get(%q) returned false, expected border tile to be registered", name)
		}
		if img == nil {
			t.Fatalf("Get(%q) returned nil image, expected non-nil", name)
		}

		alpha, ok := img.(*image.Alpha)
		if !ok {
			t.Fatalf("Get(%q) returned %T, expected *image.Alpha", name, img)
		}

		expectedBounds := image.Rect(0, 0, 8, 8)
		if alpha.Bounds() != expectedBounds {
			t.Fatalf("Get(%q) bounds = %v, expected %v", name, alpha.Bounds(), expectedBounds)
		}
	})
}

// For all border tile names excluding rounded corners, every opaque pixel in the
// tile SHALL lie on row 3 or row 4 (horizontal band) and/or column 3 or column 4
// (vertical band).

func TestProperty_BorderTilesCenterBands(t *testing.T) {
	nonRoundedNames := []string{
		"border/h",
		"border/v",
		"border/corner-tl",
		"border/corner-tr",
		"border/corner-bl",
		"border/corner-br",
		"border/tee-l",
		"border/tee-r",
		"border/tee-t",
		"border/tee-b",
		"border/cross",
	}

	rapid.Check(t, func(t *rapid.T) {
		name := rapid.SampledFrom(nonRoundedNames).Draw(t, "tileName")

		img, ok := Get(name)
		if !ok {
			t.Fatalf("Get(%q) returned false, expected tile to be registered", name)
		}

		alpha, ok := img.(*image.Alpha)
		if !ok {
			t.Fatalf("Get(%q) returned non-Alpha image type %T", name, img)
		}

		bounds := alpha.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				a := alpha.AlphaAt(x, y).A
				if a > 0 {
					inHorizontalBand := y == 3 || y == 4
					inVerticalBand := x == 3 || x == 4
					if !inHorizontalBand && !inVerticalBand {
						t.Fatalf("tile %q has opaque pixel at (%d, %d) which is outside both center bands (rows 3-4, cols 3-4)", name, x, y)
					}
				}
			}
		}
	})
}

// For any border tile that has a line exiting a given edge, the boundary pixels at
// that edge SHALL be opaque: at column 0 or 7 the pixels on rows 3 and 4 must be
// opaque (for left/right exits), and at row 0 or 7 the pixels on columns 3 and 4
// must be opaque (for top/bottom exits).

// edgeExits declares which edges each border tile has a line exiting from.
type edgeExits struct {
	Left, Right, Top, Bottom bool
}

var borderEdgeMap = map[string]edgeExits{
	"border/h":         {Left: true, Right: true},
	"border/v":         {Top: true, Bottom: true},
	"border/corner-tl": {Right: true, Bottom: true},
	"border/corner-tr": {Left: true, Bottom: true},
	"border/corner-bl": {Right: true, Top: true},
	"border/corner-br": {Left: true, Top: true},
	"border/tee-l":     {Left: true, Top: true, Bottom: true},
	"border/tee-r":     {Right: true, Top: true, Bottom: true},
	"border/tee-t":     {Left: true, Right: true, Top: true},
	"border/tee-b":     {Left: true, Right: true, Bottom: true},
	"border/cross":     {Left: true, Right: true, Top: true, Bottom: true},
	"border/round-tl":  {Right: true, Bottom: true},
	"border/round-tr":  {Left: true, Bottom: true},
	"border/round-bl":  {Right: true, Top: true},
	"border/round-br":  {Left: true, Top: true},
}

func TestProperty_BorderTilesEdgeConnectivity(t *testing.T) {
	edgeMapKeys := make([]string, 0, len(borderEdgeMap))
	for k := range borderEdgeMap {
		edgeMapKeys = append(edgeMapKeys, k)
	}

	rapid.Check(t, func(t *rapid.T) {
		name := rapid.SampledFrom(edgeMapKeys).Draw(t, "tileName")
		exits := borderEdgeMap[name]

		img, ok := Get(name)
		if !ok {
			t.Fatalf("Get(%q) returned false, expected tile to be registered", name)
		}

		alpha, ok := img.(*image.Alpha)
		if !ok {
			t.Fatalf("Get(%q) returned non-Alpha image type %T", name, img)
		}

		// Left exit: pixels at (0,3) and (0,4) must be opaque
		if exits.Left {
			if alpha.AlphaAt(0, 3).A != 0xFF {
				t.Fatalf("tile %q declares left exit but pixel (0,3) alpha = 0x%02X, expected 0xFF", name, alpha.AlphaAt(0, 3).A)
			}
			if alpha.AlphaAt(0, 4).A != 0xFF {
				t.Fatalf("tile %q declares left exit but pixel (0,4) alpha = 0x%02X, expected 0xFF", name, alpha.AlphaAt(0, 4).A)
			}
		}

		// Right exit: pixels at (7,3) and (7,4) must be opaque
		if exits.Right {
			if alpha.AlphaAt(7, 3).A != 0xFF {
				t.Fatalf("tile %q declares right exit but pixel (7,3) alpha = 0x%02X, expected 0xFF", name, alpha.AlphaAt(7, 3).A)
			}
			if alpha.AlphaAt(7, 4).A != 0xFF {
				t.Fatalf("tile %q declares right exit but pixel (7,4) alpha = 0x%02X, expected 0xFF", name, alpha.AlphaAt(7, 4).A)
			}
		}

		// Top exit: pixels at (3,0) and (4,0) must be opaque
		if exits.Top {
			if alpha.AlphaAt(3, 0).A != 0xFF {
				t.Fatalf("tile %q declares top exit but pixel (3,0) alpha = 0x%02X, expected 0xFF", name, alpha.AlphaAt(3, 0).A)
			}
			if alpha.AlphaAt(4, 0).A != 0xFF {
				t.Fatalf("tile %q declares top exit but pixel (4,0) alpha = 0x%02X, expected 0xFF", name, alpha.AlphaAt(4, 0).A)
			}
		}

		// Bottom exit: pixels at (3,7) and (4,7) must be opaque
		if exits.Bottom {
			if alpha.AlphaAt(3, 7).A != 0xFF {
				t.Fatalf("tile %q declares bottom exit but pixel (3,7) alpha = 0x%02X, expected 0xFF", name, alpha.AlphaAt(3, 7).A)
			}
			if alpha.AlphaAt(4, 7).A != 0xFF {
				t.Fatalf("tile %q declares bottom exit but pixel (4,7) alpha = 0x%02X, expected 0xFF", name, alpha.AlphaAt(4, 7).A)
			}
		}
	})
}

// TestProperty_BorderTilesRoundedDifference verifies that each rounded corner tile
// differs from its corresponding sharp corner tile in at least one pixel position.
//
// For any rounded corner tile ("border/round-tl", "border/round-tr", "border/round-bl",
// "border/round-br"), its pixel data SHALL differ from the corresponding sharp corner
// tile ("border/corner-tl", "border/corner-tr", "border/corner-bl", "border/corner-br")
// in at least one pixel position.

func TestProperty_BorderTilesRoundedDifference(t *testing.T) {
	roundedToSharp := map[string]string{
		"border/round-tl": "border/corner-tl",
		"border/round-tr": "border/corner-tr",
		"border/round-bl": "border/corner-bl",
		"border/round-br": "border/corner-br",
	}

	roundedNames := []string{
		"border/round-tl",
		"border/round-tr",
		"border/round-bl",
		"border/round-br",
	}

	rapid.Check(t, func(t *rapid.T) {
		roundedName := rapid.SampledFrom(roundedNames).Draw(t, "roundedTile")
		sharpName := roundedToSharp[roundedName]

		roundedImg, ok := Get(roundedName)
		if !ok {
			t.Fatalf("Get(%q) returned false, expected tile to be registered", roundedName)
		}
		sharpImg, ok := Get(sharpName)
		if !ok {
			t.Fatalf("Get(%q) returned false, expected tile to be registered", sharpName)
		}

		roundedAlpha, ok := roundedImg.(*image.Alpha)
		if !ok {
			t.Fatalf("Get(%q) returned %T, expected *image.Alpha", roundedName, roundedImg)
		}
		sharpAlpha, ok := sharpImg.(*image.Alpha)
		if !ok {
			t.Fatalf("Get(%q) returned %T, expected *image.Alpha", sharpName, sharpImg)
		}

		// Compare pixel-by-pixel; at least one pixel must differ
		bounds := roundedAlpha.Bounds()
		differs := false
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if roundedAlpha.AlphaAt(x, y).A != sharpAlpha.AlphaAt(x, y).A {
					differs = true
					break
				}
			}
			if differs {
				break
			}
		}

		if !differs {
			t.Fatalf("rounded tile %q has identical pixel data to sharp tile %q; expected at least one pixel to differ", roundedName, sharpName)
		}
	})
}

// --- From: icons_test.go ---

// For any non-empty normalized name string and for any non-nil image.Image,
// calling Register(name, img) followed by Get(name) SHALL return the same
// image.Image value and true.

// registryRoundTripInput holds a randomly generated test input for the round-trip property.
type registryRoundTripInput struct {
	Name string // Non-empty name (after normalization)
	W, H int    // Image dimensions (1..16)
	Seed int64  // For reproducible image generation
}

// Generate implements quick.Generator for property-based testing.
func (registryRoundTripInput) Generate(r *rand.Rand, size int) reflect.Value {
	// Generate a non-empty name with random lowercase characters + unique prefix.
	nameLen := 1 + r.Intn(20)
	chars := make([]byte, nameLen)
	for i := range chars {
		chars[i] = byte('a' + r.Intn(26))
	}
	// Prefix with a unique test ID to avoid collisions with other tests.
	name := "pbt8_" + string(chars)

	return reflect.ValueOf(registryRoundTripInput{
		Name: name,
		W:    1 + r.Intn(16),
		H:    1 + r.Intn(16),
		Seed: r.Int63(),
	})
}

func TestProperty8_IconRegistryRoundTrip(t *testing.T) {
	cfg := &quick.Config{MaxCount: 200}

	prop := func(input registryRoundTripInput) bool {
		// Generate a non-nil image with the given dimensions.
		rng := rand.New(rand.NewSource(input.Seed))
		img := generateRandomImage(rng, input.W, input.H)

		// Register the icon.
		Register(input.Name, img)

		// Get it back.
		got, ok := Get(input.Name)

		// Must be found.
		if !ok {
			t.Logf("Get(%q) returned false after Register", input.Name)
			return false
		}

		// Must be the same image value (pointer equality).
		if got != img {
			t.Logf("Get(%q) returned different image: got %p, want %p", input.Name, got, img)
			return false
		}

		return true
	}

	if err := quick.Check(prop, cfg); err != nil {
		t.Errorf("Property 8 failed: %v", err)
	}
}

// generateRandomImage creates a random RGBA image with the given dimensions.
func generateRandomImage(rng *rand.Rand, w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(rng.Intn(256)),
				G: uint8(rng.Intn(256)),
				B: uint8(rng.Intn(256)),
				A: uint8(rng.Intn(256)),
			})
		}
	}
	return img
}

// For any name string and for any two distinct non-nil images img1 and img2,
// calling Register(name, img1) then Register(name, img2) then Get(name) SHALL return img2 and true.
func TestProperty9_IconRegistryOverwriteSemantics(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}

	iteration := 0
	err := quick.Check(func(seed uint16) bool {
		iteration++
		// Generate a unique name per iteration to avoid collisions with other tests
		name := fmt.Sprintf("overwrite_test_%d_%d", iteration, seed)

		// Create two distinct non-nil images with different pixel content
		img1 := image.NewRGBA(image.Rect(0, 0, 4, 4))
		img1.Set(0, 0, color.RGBA{R: uint8(seed % 256), G: 0, B: 0, A: 255})

		img2 := image.NewRGBA(image.Rect(0, 0, 4, 4))
		img2.Set(0, 0, color.RGBA{R: 0, G: uint8(seed % 256), B: 0, A: 255})

		// Register first image
		Register(name, img1)

		// Overwrite with second image
		Register(name, img2)

		// Get should return img2 and true
		got, ok := Get(name)
		if !ok {
			t.Logf("iteration %d: Get(%q) returned false, expected true", iteration, name)
			return false
		}
		if got != img2 {
			t.Logf("iteration %d: Get(%q) returned img1 instead of img2", iteration, name)
			return false
		}

		return true
	}, config)

	if err != nil {
		t.Errorf("Property 9 failed: %v", err)
	}
}

// For any name string that has never been passed to Register, calling Get(name)
// SHALL return nil and false.
func TestProperty10_IconRegistryUnregisteredNameReturnsNil(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}

	iteration := 0
	err := quick.Check(func(seed uint32) bool {
		iteration++
		// Use a prefix guaranteed to not match any registered icon,
		// combined with a unique iteration counter and seed for uniqueness.
		name := fmt.Sprintf("unregistered_prop10_%d_%d", iteration, seed)

		img, ok := Get(name)
		if img != nil {
			t.Logf("iteration %d: Get(%q) returned non-nil image, expected nil", iteration, name)
			return false
		}
		if ok {
			t.Logf("iteration %d: Get(%q) returned true, expected false", iteration, name)
			return false
		}
		return true
	}, config)

	if err != nil {
		t.Errorf("Property 10 failed: %v", err)
	}
}

// --- Unit Tests: Icon Registry Edge Cases and Built-in Icons ---

// TestRegisterEmptyName verifies that Register with "" and whitespace-only names
// doesn't store anything in the registry.

func TestRegisterEmptyName(t *testing.T) {
	cases := []struct {
		name string
		desc string
	}{
		{"", "empty string"},
		{"   ", "whitespace only (spaces)"},
		{"\t\n", "whitespace only (tab and newline)"},
	}

	img := image.NewRGBA(image.Rect(0, 0, 4, 4))

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			Register(tc.name, img)
			// After normalization these all become "", which should not be retrievable.
			got, ok := Get(tc.name)
			if ok {
				t.Errorf("Get(%q) returned true, expected false for empty/whitespace name", tc.name)
			}
			if got != nil {
				t.Errorf("Get(%q) returned non-nil image, expected nil for empty/whitespace name", tc.name)
			}
		})
	}
}

// TestRegisterNilImage verifies that Register with a valid name but nil image
// doesn't store anything in the registry.

func TestRegisterNilImage(t *testing.T) {
	name := "nil_image_test_unique_key"
	Register(name, nil)

	got, ok := Get(name)
	if ok {
		t.Errorf("Get(%q) returned true after registering nil image, expected false", name)
	}
	if got != nil {
		t.Errorf("Get(%q) returned non-nil image after registering nil image", name)
	}
}

// TestBuiltinIconsExist verifies that all built-in icons are registered during
// package initialization with correct dimensions (between 4x4 and 16x16 inclusive).

func TestBuiltinIconsExist(t *testing.T) {
	builtins := []string{"wifi", "bluetooth", "battery", "error", "check"}

	for _, name := range builtins {
		t.Run(name, func(t *testing.T) {
			img, ok := Get(name)
			if !ok {
				t.Fatalf("Get(%q) returned false, expected built-in icon to be registered", name)
			}
			if img == nil {
				t.Fatalf("Get(%q) returned nil image, expected non-nil built-in icon", name)
			}

			bounds := img.Bounds()
			w := bounds.Dx()
			h := bounds.Dy()

			if w < 4 || w > 16 {
				t.Errorf("Get(%q) icon width = %d, expected between 4 and 16 inclusive", name, w)
			}
			if h < 4 || h > 16 {
				t.Errorf("Get(%q) icon height = %d, expected between 4 and 16 inclusive", name, h)
			}
		})
	}
}

// --- Unit Tests: Names() edge cases ---

// TestNames_EmptyRegistry verifies that Names() returns a non-nil empty slice
// when no icons are registered (after Reset).
func TestNames_EmptyRegistry(t *testing.T) {
	Reset()
	defer Reset()

	names := Names()
	if names == nil {
		t.Fatal("Names() returned nil, expected non-nil empty slice")
	}
	if len(names) != 0 {
		t.Fatalf("Names() returned %d elements, expected 0", len(names))
	}
}

// TestNames_SingleEntry verifies that Names() returns a single-element slice
// when exactly one icon is registered.
func TestNames_SingleEntry(t *testing.T) {
	Reset()
	defer Reset()

	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	Register("alpha", img)

	names := Names()
	expected := []string{"alpha"}
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("Names() = %v, expected %v", names, expected)
	}
}

// TestNames_MultiEntry_Sorted verifies that Names() returns all registered names
// in sorted lexicographic order regardless of insertion order.
func TestNames_MultiEntry_Sorted(t *testing.T) {
	Reset()
	defer Reset()

	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	Register("zulu", img)
	Register("alpha", img)
	Register("mike", img)

	names := Names()
	expected := []string{"alpha", "mike", "zulu"}
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("Names() = %v, expected %v", names, expected)
	}
}

// TestNames_DynamicAddition verifies that newly registered icons appear in
// subsequent Names() calls.
func TestNames_DynamicAddition(t *testing.T) {
	Reset()
	defer Reset()

	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	Register("bravo", img)

	t.Run("before addition", func(t *testing.T) {
		names := Names()
		expected := []string{"bravo"}
		if !reflect.DeepEqual(names, expected) {
			t.Fatalf("Names() = %v, expected %v", names, expected)
		}
	})

	Register("alpha", img)

	t.Run("after addition", func(t *testing.T) {
		names := Names()
		expected := []string{"alpha", "bravo"}
		if !reflect.DeepEqual(names, expected) {
			t.Fatalf("Names() = %v, expected %v", names, expected)
		}
	})
}

// TestConcurrentRegisterGet launches multiple goroutines doing Register and Get
// simultaneously to verify the icon registry is safe for concurrent access.
// Run with -race to detect data races.

func TestConcurrentRegisterGet(t *testing.T) {
	const goroutines = 20
	const iterations = 100

	done := make(chan struct{})

	// Launch writers
	for i := 0; i < goroutines/2; i++ {
		go func(id int) {
			defer func() { done <- struct{}{} }()
			for j := 0; j < iterations; j++ {
				name := fmt.Sprintf("concurrent_icon_%d", j%10)
				img := image.NewRGBA(image.Rect(0, 0, 8, 8))
				Register(name, img)
			}
		}(i)
	}

	// Launch readers
	for i := 0; i < goroutines/2; i++ {
		go func(id int) {
			defer func() { done <- struct{}{} }()
			for j := 0; j < iterations; j++ {
				name := fmt.Sprintf("concurrent_icon_%d", j%10)
				Get(name)
			}
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < goroutines; i++ {
		<-done
	}
}
