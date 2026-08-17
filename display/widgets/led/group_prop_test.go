package led

import (
	"image"
	"image/color"
	"testing"

	"pgregory.net/rapid"
)

// TestProperty21_GroupOutputDimensions verifies that for any valid group Config
// with N entries (2 ≤ N ≤ 32) and Diameter ≥ 3, the output Sprite dimensions SHALL be:
//   - Horizontal: width = N × Diameter + (N − 1) × spacing, height = Diameter
//   - Vertical: width = Diameter, height = N × Diameter + (N − 1) × spacing

func TestProperty21_GroupOutputDimensions(t *testing.T) {
	t.Run("Horizontal", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			n := rapid.IntRange(2, 8).Draw(t, "numEntries")
			diameter := rapid.IntRange(3, 32).Draw(t, "diameter")
			spacing := rapid.IntRange(0, 8).Draw(t, "spacing")

			entries := make([]GroupEntry, n)
			for i := range entries {
				entries[i] = GroupEntry{
					State:      State(rapid.IntRange(0, 1).Draw(t, "entryState")),
					Foreground: color.RGBA{R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")), G: 100, B: 50, A: 255},
				}
			}

			cfg := Config{
				Shape:       Circle,
				State:       On,
				Brightness:  -1.0,
				Diameter:    diameter,
				Bounds:      image.Rect(0, 0, diameter, diameter),
				Foreground:  color.RGBA{R: 0, G: 200, B: 0, A: 255},
				Group:       entries,
				Orientation: Horizontal,
				Spacing:     spacing,
			}

			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil sprite for valid group config")
			}

			expectedWidth := n*diameter + (n-1)*spacing
			expectedHeight := diameter

			imgBounds := result.Image.Bounds()
			actualWidth := imgBounds.Dx()
			actualHeight := imgBounds.Dy()

			if actualWidth != expectedWidth {
				t.Fatalf("horizontal width mismatch: got %d, want %d (n=%d, diameter=%d, spacing=%d)",
					actualWidth, expectedWidth, n, diameter, spacing)
			}
			if actualHeight != expectedHeight {
				t.Fatalf("horizontal height mismatch: got %d, want %d (n=%d, diameter=%d)",
					actualHeight, expectedHeight, n, diameter)
			}
		})
	})

	t.Run("Vertical", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			n := rapid.IntRange(2, 8).Draw(t, "numEntries")
			diameter := rapid.IntRange(3, 32).Draw(t, "diameter")
			spacing := rapid.IntRange(0, 8).Draw(t, "spacing")

			entries := make([]GroupEntry, n)
			for i := range entries {
				entries[i] = GroupEntry{
					State:      State(rapid.IntRange(0, 1).Draw(t, "entryState")),
					Foreground: color.RGBA{R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")), G: 100, B: 50, A: 255},
				}
			}

			cfg := Config{
				Shape:       Circle,
				State:       On,
				Brightness:  -1.0,
				Diameter:    diameter,
				Bounds:      image.Rect(0, 0, diameter, diameter),
				Foreground:  color.RGBA{R: 0, G: 200, B: 0, A: 255},
				Group:       entries,
				Orientation: Vertical,
				Spacing:     spacing,
			}

			result := Render(cfg)
			if result == nil {
				t.Fatal("expected non-nil sprite for valid group config")
			}

			expectedWidth := diameter
			expectedHeight := n*diameter + (n-1)*spacing

			imgBounds := result.Image.Bounds()
			actualWidth := imgBounds.Dx()
			actualHeight := imgBounds.Dy()

			if actualWidth != expectedWidth {
				t.Fatalf("vertical width mismatch: got %d, want %d (n=%d, diameter=%d)",
					actualWidth, expectedWidth, n, diameter)
			}
			if actualHeight != expectedHeight {
				t.Fatalf("vertical height mismatch: got %d, want %d (n=%d, diameter=%d, spacing=%d)",
					actualHeight, expectedHeight, n, diameter, spacing)
			}
		})
	})
}

// TestProperty22_SingleEntryGroupRendersAsStandardSingleLED verifies that for any
// group Config with exactly one entry, the rendered output SHALL be pixel-identical
// to a single LED rendered with the group-level settings merged with the entry's overrides.

func TestProperty22_SingleEntryGroupRendersAsStandardSingleLED(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		diameter := rapid.IntRange(3, 32).Draw(t, "diameter")
		shape := Shape(rapid.IntRange(0, 3).Draw(t, "shape"))
		state := State(rapid.IntRange(0, 2).Draw(t, "entryState"))

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
			A: 255,
		}

		groupFg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "groupFgR")),
			G: uint8(rapid.IntRange(1, 255).Draw(t, "groupFgG")),
			B: uint8(rapid.IntRange(1, 255).Draw(t, "groupFgB")),
			A: 255,
		}

		bounds := image.Rect(0, 0, diameter, diameter)

		// Single-entry group config
		groupCfg := Config{
			Shape:      shape,
			State:      On,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: groupFg,
			Group: []GroupEntry{
				{
					State:      state,
					Foreground: fg,
				},
			},
			Orientation: Horizontal,
			Spacing:     2,
		}

		groupResult := Render(groupCfg)
		if groupResult == nil {
			t.Fatal("expected non-nil sprite from single-entry group")
		}

		// Build the equivalent single LED config using buildEntryCfg logic:
		// Entry state is always used directly, entry foreground overrides group if non-zero.
		singleCfg := Config{
			Shape:      shape,
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     bounds,
			Foreground: fg, // Entry foreground overrides group-level
		}

		singleResult := Render(singleCfg)
		if singleResult == nil {
			t.Fatal("expected non-nil sprite from single LED")
		}

		// Compare pixel-for-pixel
		groupImg := groupResult.Image.(*image.RGBA)
		singleImg := singleResult.Image.(*image.RGBA)

		if groupImg.Bounds() != singleImg.Bounds() {
			t.Fatalf("dimension mismatch: group=%v, single=%v",
				groupImg.Bounds(), singleImg.Bounds())
		}

		bounds2 := groupImg.Bounds()
		for y := bounds2.Min.Y; y < bounds2.Max.Y; y++ {
			for x := bounds2.Min.X; x < bounds2.Max.X; x++ {
				gc := groupImg.RGBAAt(x, y)
				sc := singleImg.RGBAAt(x, y)
				if gc != sc {
					t.Fatalf("pixel mismatch at (%d,%d): group=%v, single=%v "+
						"[diameter=%d, shape=%d, state=%d, fg=%v, groupFg=%v]",
						x, y, gc, sc, diameter, shape, state, fg, groupFg)
				}
			}
		}
	})
}

// TestProperty23_GroupEntryIndependence verifies that for any group Config with N
// entries where each entry has different state and foreground color, each LED in the
// output image SHALL be rendered using its own state and color.
//
// We verify this by checking the center pixel of each LED cell: On-state LEDs should
// have their own foreground color (brightness-scaled), while Off-state LEDs should
// have the background color (body interior).

func TestProperty23_GroupEntryIndependence(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(2, 6).Draw(t, "numEntries")
		diameter := rapid.IntRange(5, 20).Draw(t, "diameter")
		spacing := rapid.IntRange(0, 4).Draw(t, "spacing")
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		// Generate entries with distinct foreground colors and alternating states
		entries := make([]GroupEntry, n)
		colors := make([]color.RGBA, n)
		states := make([]State, n)

		for i := range entries {
			// Alternate On/Off so we have both states
			st := On
			if i%2 == 1 {
				st = Off
			}
			states[i] = st

			// Each entry gets a unique, distinguishable foreground color
			fg := color.RGBA{
				R: uint8(50 + i*40),
				G: uint8(30 + i*30),
				B: uint8(20 + i*20),
				A: 255,
			}
			colors[i] = fg
			entries[i] = GroupEntry{
				State:      st,
				Foreground: fg,
			}
		}

		cfg := Config{
			Shape:       Circle,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      image.Rect(0, 0, diameter, diameter),
			Foreground:  color.RGBA{R: 0, G: 200, B: 0, A: 255},
			Background:  color.RGBA{R: 0, G: 0, B: 0, A: 255},
			Group:       entries,
			Orientation: orientation,
			Spacing:     spacing,
		}

		result := Render(cfg)
		if result == nil {
			t.Fatal("expected non-nil sprite for valid group config")
		}

		img := result.Image.(*image.RGBA)

		// Check the center pixel of each LED cell
		for i := 0; i < n; i++ {
			var cellCenterX, cellCenterY int
			if orientation == Horizontal {
				cellCenterX = i*(diameter+spacing) + diameter/2
				cellCenterY = diameter / 2
			} else {
				cellCenterX = diameter / 2
				cellCenterY = i*(diameter+spacing) + diameter/2
			}

			pixel := img.RGBAAt(cellCenterX, cellCenterY)

			if states[i] == On {
				// On state: center pixel should be the entry's foreground color
				expected := colors[i]
				if pixel.R != expected.R || pixel.G != expected.G || pixel.B != expected.B {
					t.Fatalf("entry %d (On): center pixel (%d,%d) = %v, want foreground %v "+
						"[n=%d, diameter=%d, spacing=%d, orientation=%d]",
						i, cellCenterX, cellCenterY, pixel, expected,
						n, diameter, spacing, orientation)
				}
			} else {
				// Off state: center pixel should be background (black)
				bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}
				if pixel != bg {
					t.Fatalf("entry %d (Off): center pixel (%d,%d) = %v, want background %v "+
						"[n=%d, diameter=%d, spacing=%d, orientation=%d]",
						i, cellCenterX, cellCenterY, pixel, bg,
						n, diameter, spacing, orientation)
				}
			}
		}
	})
}

// TestProperty24_GroupTruncationTo32Entries verifies that for any group Config with
// more than 32 entries, the rendered output SHALL be pixel-identical to the same
// Config with only the first 32 entries.

func TestProperty24_GroupTruncationTo32Entries(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate 33–40 entries
		totalEntries := rapid.IntRange(33, 40).Draw(t, "totalEntries")
		diameter := rapid.IntRange(3, 10).Draw(t, "diameter")
		spacing := rapid.IntRange(0, 4).Draw(t, "spacing")
		orientation := Orientation(rapid.IntRange(0, 1).Draw(t, "orientation"))

		allEntries := make([]GroupEntry, totalEntries)
		for i := range allEntries {
			allEntries[i] = GroupEntry{
				State: State(rapid.IntRange(0, 1).Draw(t, "entryState")),
				Foreground: color.RGBA{
					R: uint8(rapid.IntRange(1, 255).Draw(t, "fgR")),
					G: uint8(rapid.IntRange(1, 255).Draw(t, "fgG")),
					B: uint8(rapid.IntRange(1, 255).Draw(t, "fgB")),
					A: 255,
				},
			}
		}

		bounds := image.Rect(0, 0, diameter, diameter)
		fg := color.RGBA{R: 0, G: 200, B: 0, A: 255}

		// Config with all entries (>32)
		fullCfg := Config{
			Shape:       Circle,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      bounds,
			Foreground:  fg,
			Group:       allEntries,
			Orientation: orientation,
			Spacing:     spacing,
		}

		// Config with only first 32 entries
		truncatedEntries := make([]GroupEntry, 32)
		copy(truncatedEntries, allEntries[:32])
		truncCfg := Config{
			Shape:       Circle,
			State:       On,
			Brightness:  -1.0,
			Diameter:    diameter,
			Bounds:      bounds,
			Foreground:  fg,
			Group:       truncatedEntries,
			Orientation: orientation,
			Spacing:     spacing,
		}

		fullResult := Render(fullCfg)
		truncResult := Render(truncCfg)

		if fullResult == nil || truncResult == nil {
			t.Fatal("expected non-nil sprites for valid group configs")
		}

		fullImg := fullResult.Image.(*image.RGBA)
		truncImg := truncResult.Image.(*image.RGBA)

		if fullImg.Bounds() != truncImg.Bounds() {
			t.Fatalf("dimension mismatch: full(%d entries)=%v, truncated(32)=%v",
				totalEntries, fullImg.Bounds(), truncImg.Bounds())
		}

		imgBounds := fullImg.Bounds()
		for y := imgBounds.Min.Y; y < imgBounds.Max.Y; y++ {
			for x := imgBounds.Min.X; x < imgBounds.Max.X; x++ {
				fc := fullImg.RGBAAt(x, y)
				tc := truncImg.RGBAAt(x, y)
				if fc != tc {
					t.Fatalf("pixel mismatch at (%d,%d): full=%v, truncated=%v "+
						"[totalEntries=%d, diameter=%d, spacing=%d, orientation=%d]",
						x, y, fc, tc, totalEntries, diameter, spacing, orientation)
				}
			}
		}
	})
}
