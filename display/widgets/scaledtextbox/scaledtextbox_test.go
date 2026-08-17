package scaledtextbox

import (
	"image"
	"image/color"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/scale"
	"github.com/databeast/cyberhud/display/widgets/textbox"
	"pgregory.net/rapid"
)

// --- From: scaledtextbox_prop_test.go ---

// **Feature: textbox-widget, Property 16: ScaledTextBox output matches reference scaling**
//
// For any valid ScaledTextBox Config with Border=false, the output Image SHALL be
// pixel-identical to calling scale.NearestNeighbor(textbox.Render(innerConfig).Image,
// TargetSize.X, TargetSize.Y) where innerConfig uses LogicalSize as Bounds.
//

func TestPropertyScaledTextBoxMatchesReferenceScaling(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()

		// Generate LogicalSize (20-100 × 10-50).
		logicalX := rapid.IntRange(20, 100).Draw(t, "logicalX")
		logicalY := rapid.IntRange(10, 50).Draw(t, "logicalY")

		// Generate TargetSize (10-80 × 5-40).
		targetX := rapid.IntRange(10, 80).Draw(t, "targetX")
		targetY := rapid.IntRange(5, 40).Draw(t, "targetY")

		// Generate text content (1-10 printable chars).
		text := rapid.StringMatching(`[A-Za-z0-9]{1,10}`).Draw(t, "text")

		// Generate foreground color (non-zero to avoid default logic differences).
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "r")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "g")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "b")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "a")),
		}

		// Generate alignment and overflow.
		alignment := textbox.Alignment(rapid.IntRange(0, 2).Draw(t, "alignment"))
		valign := textbox.VAlign(rapid.IntRange(0, 2).Draw(t, "valign"))
		overflow := textbox.Overflow(rapid.IntRange(0, 2).Draw(t, "overflow"))

		// Generate small padding values that still leave positive effective area.
		maxPadX := (logicalX - 1) / 2
		if maxPadX > 5 {
			maxPadX = 5
		}
		if maxPadX < 0 {
			maxPadX = 0
		}
		maxPadY := (logicalY - 1) / 2
		if maxPadY > 5 {
			maxPadY = 5
		}
		if maxPadY < 0 {
			maxPadY = 0
		}
		padX := rapid.IntRange(0, maxPadX).Draw(t, "padX")
		padY := rapid.IntRange(0, maxPadY).Draw(t, "padY")

		cfg := Config{
			LogicalSize: image.Point{X: logicalX, Y: logicalY},
			TargetSize:  image.Point{X: targetX, Y: targetY},
			Text:        text,
			Font:        face,
			Alignment:   alignment,
			VAlign:      valign,
			Overflow:    overflow,
			Foreground:  fg,
			PadX:        padX,
			PadY:        padY,
			Border:      false, // Property 16 requires Border=false.
		}

		result := Render(cfg)

		// Construct reference: render inner textbox then scale.
		innerConfig := textbox.Config{
			Bounds:     image.Rect(0, 0, logicalX, logicalY),
			Text:       text,
			Font:       face,
			Alignment:  alignment,
			VAlign:     valign,
			Overflow:   overflow,
			Foreground: fg,
			PadX:       padX,
			PadY:       padY,
			Border:     false,
		}
		innerResult := textbox.Render(innerConfig)

		if innerResult == nil {
			// If inner textbox returns nil, ScaledTextBox should also return nil.
			if result != nil {
				t.Fatalf("expected nil result when inner textbox returns nil")
			}
			return
		}

		if result == nil {
			t.Fatalf("expected non-nil result when inner textbox is valid")
		}

		reference := scale.NearestNeighbor(innerResult.Image, targetX, targetY)
		if reference == nil {
			t.Fatalf("reference scale returned nil unexpectedly")
		}

		// Compare pixel by pixel.
		for y := 0; y < targetY; y++ {
			for x := 0; x < targetX; x++ {
				rr, rg, rb, ra := result.Image.At(x, y).RGBA()
				er, eg, eb, ea := reference.At(x, y).RGBA()
				if rr != er || rg != eg || rb != eb || ra != ea {
					t.Fatalf("pixel mismatch at (%d, %d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
						x, y, rr, rg, rb, ra, er, eg, eb, ea)
				}
			}
		}
	})
}

// **Feature: textbox-widget, Property 17: ScaledTextBox 1:1 scaling is identity**
//
// For any ScaledTextBox Config where LogicalSize equals TargetSize and Border=false,
// the output Image SHALL be pixel-identical to a plain TextBox Render with Bounds
// set to that size.
//

func TestPropertyScaledTextBox1to1Identity(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()

		// Generate size used for both LogicalSize and TargetSize.
		sizeX := rapid.IntRange(10, 80).Draw(t, "sizeX")
		sizeY := rapid.IntRange(8, 40).Draw(t, "sizeY")

		// Generate text content.
		text := rapid.StringMatching(`[A-Za-z0-9]{1,10}`).Draw(t, "text")

		// Generate foreground color (non-zero).
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "r")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "g")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "b")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "a")),
		}

		alignment := textbox.Alignment(rapid.IntRange(0, 2).Draw(t, "alignment"))
		valign := textbox.VAlign(rapid.IntRange(0, 2).Draw(t, "valign"))
		overflow := textbox.Overflow(rapid.IntRange(0, 2).Draw(t, "overflow"))

		// Small padding.
		maxPadX := (sizeX - 1) / 2
		if maxPadX > 3 {
			maxPadX = 3
		}
		if maxPadX < 0 {
			maxPadX = 0
		}
		maxPadY := (sizeY - 1) / 2
		if maxPadY > 3 {
			maxPadY = 3
		}
		if maxPadY < 0 {
			maxPadY = 0
		}
		padX := rapid.IntRange(0, maxPadX).Draw(t, "padX")
		padY := rapid.IntRange(0, maxPadY).Draw(t, "padY")

		// ScaledTextBox with 1:1 scaling.
		scaledCfg := Config{
			LogicalSize: image.Point{X: sizeX, Y: sizeY},
			TargetSize:  image.Point{X: sizeX, Y: sizeY},
			Text:        text,
			Font:        face,
			Alignment:   alignment,
			VAlign:      valign,
			Overflow:    overflow,
			Foreground:  fg,
			PadX:        padX,
			PadY:        padY,
			Border:      false,
		}

		scaledResult := Render(scaledCfg)

		// Plain TextBox with same size as Bounds.
		plainCfg := textbox.Config{
			Bounds:     image.Rect(0, 0, sizeX, sizeY),
			Text:       text,
			Font:       face,
			Alignment:  alignment,
			VAlign:     valign,
			Overflow:   overflow,
			Foreground: fg,
			PadX:       padX,
			PadY:       padY,
			Border:     false,
		}

		plainResult := textbox.Render(plainCfg)

		if plainResult == nil {
			// If plain textbox returns nil, scaled should too.
			if scaledResult != nil {
				t.Fatalf("expected nil result when plain textbox returns nil")
			}
			return
		}

		if scaledResult == nil {
			t.Fatalf("expected non-nil result when plain textbox is valid")
		}

		// Compare pixel by pixel.
		for y := 0; y < sizeY; y++ {
			for x := 0; x < sizeX; x++ {
				rr, rg, rb, ra := scaledResult.Image.At(x, y).RGBA()
				er, eg, eb, ea := plainResult.Image.At(x, y).RGBA()
				if rr != er || rg != eg || rb != eb || ra != ea {
					t.Fatalf("pixel mismatch at (%d, %d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
						x, y, rr, rg, rb, ra, er, eg, eb, ea)
				}
			}
		}
	})
}

// **Feature: textbox-widget, Property 18: ScaledTextBox invalid dimensions produce nil**
//
// For any ScaledTextBox Config where LogicalSize has zero or negative X or Y,
// OR TargetSize has zero or negative X or Y, Render SHALL return nil.
//

func TestPropertyScaledTextBoxInvalidDimensionsNil(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()

		// Generate text content.
		text := rapid.StringMatching(`[A-Za-z0-9]{1,10}`).Draw(t, "text")

		fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}

		// Choose which dimension to make invalid.
		// 0 = LogicalSize.X invalid
		// 1 = LogicalSize.Y invalid
		// 2 = TargetSize.X invalid
		// 3 = TargetSize.Y invalid
		invalidChoice := rapid.IntRange(0, 3).Draw(t, "invalidChoice")

		// Default valid values.
		logicalX := rapid.IntRange(10, 50).Draw(t, "logicalX")
		logicalY := rapid.IntRange(10, 50).Draw(t, "logicalY")
		targetX := rapid.IntRange(10, 50).Draw(t, "targetX")
		targetY := rapid.IntRange(10, 50).Draw(t, "targetY")

		// Make one dimension invalid (zero or negative).
		invalidVal := rapid.IntRange(-10, 0).Draw(t, "invalidVal")
		switch invalidChoice {
		case 0:
			logicalX = invalidVal
		case 1:
			logicalY = invalidVal
		case 2:
			targetX = invalidVal
		case 3:
			targetY = invalidVal
		}

		cfg := Config{
			LogicalSize: image.Point{X: logicalX, Y: logicalY},
			TargetSize:  image.Point{X: targetX, Y: targetY},
			Text:        text,
			Font:        face,
			Foreground:  fg,
		}

		result := Render(cfg)
		if result != nil {
			t.Fatalf("expected nil result for invalid dimensions (logicalSize=%v, targetSize=%v, invalidChoice=%d, invalidVal=%d)",
				cfg.LogicalSize, cfg.TargetSize, invalidChoice, invalidVal)
		}
	})
}

// **Feature: textbox-widget, Property 19: ScaledTextBox output dimensions match target**
//
// For any valid ScaledTextBox Config, the output Image SHALL have dimensions
// exactly equal to TargetSize.X × TargetSize.Y.
//

func TestPropertyScaledTextBoxOutputDimensionsMatchTarget(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()

		// Generate LogicalSize (20-100 × 10-50).
		logicalX := rapid.IntRange(20, 100).Draw(t, "logicalX")
		logicalY := rapid.IntRange(10, 50).Draw(t, "logicalY")

		// Generate TargetSize (10-80 × 5-40).
		targetX := rapid.IntRange(10, 80).Draw(t, "targetX")
		targetY := rapid.IntRange(5, 40).Draw(t, "targetY")

		// Generate text content (may be empty for transparent output).
		text := rapid.StringMatching(`[A-Za-z0-9]{0,10}`).Draw(t, "text")

		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "r")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "g")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "b")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "a")),
		}

		// Generate small padding that still leaves positive effective area.
		maxPadX := (logicalX - 1) / 2
		if maxPadX > 5 {
			maxPadX = 5
		}
		if maxPadX < 0 {
			maxPadX = 0
		}
		maxPadY := (logicalY - 1) / 2
		if maxPadY > 5 {
			maxPadY = 5
		}
		if maxPadY < 0 {
			maxPadY = 0
		}
		padX := rapid.IntRange(0, maxPadX).Draw(t, "padX")
		padY := rapid.IntRange(0, maxPadY).Draw(t, "padY")

		border := rapid.Bool().Draw(t, "border")

		// Skip if border is true and TargetSize is too small.
		if border && (targetX <= 2 || targetY <= 2) {
			return
		}

		cfg := Config{
			LogicalSize: image.Point{X: logicalX, Y: logicalY},
			TargetSize:  image.Point{X: targetX, Y: targetY},
			Text:        text,
			Font:        face,
			Foreground:  fg,
			PadX:        padX,
			PadY:        padY,
			Border:      border,
		}

		result := Render(cfg)
		if result == nil {
			// Possible if padding overflows logical size. Skip these cases.
			return
		}

		// Verify image dimensions match TargetSize exactly.
		imgBounds := result.Image.Bounds()
		gotWidth := imgBounds.Dx()
		gotHeight := imgBounds.Dy()

		if gotWidth != targetX {
			t.Fatalf("image width mismatch: got %d, want %d (logicalSize=%v, targetSize=%v)",
				gotWidth, targetX, cfg.LogicalSize, cfg.TargetSize)
		}
		if gotHeight != targetY {
			t.Fatalf("image height mismatch: got %d, want %d (logicalSize=%v, targetSize=%v)",
				gotHeight, targetY, cfg.LogicalSize, cfg.TargetSize)
		}
	})
}

// **Feature: textbox-widget, Property 22: ScaledTextBox border is 1px at TargetSize**
//
// For any valid ScaledTextBox Config with Border=true, the output Image SHALL have
// a 1px border drawn at the TargetSize edges. Every pixel on the outermost row and
// column of the output SHALL have RGBA values matching the effective Foreground
// color, and the border SHALL be exactly 1px wide.
//

func TestPropertyScaledTextBoxBorder1pxAtTargetSize(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		face := font.Default()

		// Generate LogicalSize (20-100 × 10-50).
		logicalX := rapid.IntRange(20, 100).Draw(t, "logicalX")
		logicalY := rapid.IntRange(10, 50).Draw(t, "logicalY")

		// Generate TargetSize (must be > 2 for border to be valid).
		targetX := rapid.IntRange(4, 80).Draw(t, "targetX")
		targetY := rapid.IntRange(4, 40).Draw(t, "targetY")

		// Generate text (may be empty - border should still be drawn).
		text := rapid.StringMatching(`[A-Za-z0-9]{0,10}`).Draw(t, "text")

		// Generate foreground color. Non-zero ensures we can verify border color.
		fg := color.RGBA{
			R: uint8(rapid.IntRange(1, 255).Draw(t, "r")),
			G: uint8(rapid.IntRange(0, 255).Draw(t, "g")),
			B: uint8(rapid.IntRange(0, 255).Draw(t, "b")),
			A: uint8(rapid.IntRange(1, 255).Draw(t, "a")),
		}

		// Small padding that leaves positive effective area at logical size.
		maxPadX := (logicalX - 1) / 2
		if maxPadX > 3 {
			maxPadX = 3
		}
		if maxPadX < 0 {
			maxPadX = 0
		}
		maxPadY := (logicalY - 1) / 2
		if maxPadY > 3 {
			maxPadY = 3
		}
		if maxPadY < 0 {
			maxPadY = 0
		}
		padX := rapid.IntRange(0, maxPadX).Draw(t, "padX")
		padY := rapid.IntRange(0, maxPadY).Draw(t, "padY")

		cfg := Config{
			LogicalSize: image.Point{X: logicalX, Y: logicalY},
			TargetSize:  image.Point{X: targetX, Y: targetY},
			Text:        text,
			Font:        face,
			Foreground:  fg,
			PadX:        padX,
			PadY:        padY,
			Border:      true,
		}

		result := Render(cfg)
		if result == nil {
			// Possible if padding overflows logical size. Skip these cases.
			return
		}

		// Determine effective foreground color (same logic as implementation).
		effectiveFg := fg
		if effectiveFg == (color.RGBA{}) {
			effectiveFg = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}

		// Convert to RGBA 16-bit for comparison (how image.At().RGBA() returns values).
		wantR := uint32(effectiveFg.R) * 0x101
		wantG := uint32(effectiveFg.G) * 0x101
		wantB := uint32(effectiveFg.B) * 0x101
		wantA := uint32(effectiveFg.A) * 0x101

		// Check top row (y=0).
		for x := 0; x < targetX; x++ {
			r, g, b, a := result.Image.At(x, 0).RGBA()
			if r != wantR || g != wantG || b != wantB || a != wantA {
				t.Fatalf("top row pixel (%d, 0): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					x, r, g, b, a, wantR, wantG, wantB, wantA)
			}
		}

		// Check bottom row (y=targetY-1).
		for x := 0; x < targetX; x++ {
			r, g, b, a := result.Image.At(x, targetY-1).RGBA()
			if r != wantR || g != wantG || b != wantB || a != wantA {
				t.Fatalf("bottom row pixel (%d, %d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					x, targetY-1, r, g, b, a, wantR, wantG, wantB, wantA)
			}
		}

		// Check left column (x=0).
		for y := 0; y < targetY; y++ {
			r, g, b, a := result.Image.At(0, y).RGBA()
			if r != wantR || g != wantG || b != wantB || a != wantA {
				t.Fatalf("left col pixel (0, %d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					y, r, g, b, a, wantR, wantG, wantB, wantA)
			}
		}

		// Check right column (x=targetX-1).
		for y := 0; y < targetY; y++ {
			r, g, b, a := result.Image.At(targetX-1, y).RGBA()
			if r != wantR || g != wantG || b != wantB || a != wantA {
				t.Fatalf("right col pixel (%d, %d): got RGBA(%d,%d,%d,%d), want RGBA(%d,%d,%d,%d)",
					targetX-1, y, r, g, b, a, wantR, wantG, wantB, wantA)
			}
		}

		// Verify border is exactly 1px wide: the second row/column inward should
		// NOT all be the border color (unless the text or scaling happens to produce
		// that color at those positions). We verify by checking that at least some
		// interior pixels (row 1, col 1) differ from the border color when text is empty.
		// A stronger check: for empty text, interior should be transparent.
		if text == "" && targetX > 2 && targetY > 2 {
			// Interior pixels (not on border) should be fully transparent.
			for y := 1; y < targetY-1; y++ {
				for x := 1; x < targetX-1; x++ {
					_, _, _, a := result.Image.At(x, y).RGBA()
					if a != 0 {
						t.Fatalf("border width > 1px: interior pixel (%d, %d) is not transparent for empty text (alpha=%d)",
							x, y, a)
					}
				}
			}
		}
	})
}

// --- From: scaledtextbox_test.go ---

// Test2xLogicalTo1xTarget verifies that rendering at 2× logical size and
// scaling down to 1× target produces output with correct TargetSize dimensions.

func Test2xLogicalTo1xTarget(t *testing.T) {
	result := Render(Config{
		LogicalSize: image.Point{X: 60, Y: 20},
		TargetSize:  image.Point{X: 30, Y: 10},
		Text:        "Hi",
		Font:        font.Default(),
	})
	if result == nil {
		t.Fatal("expected non-nil result for 2× logical → 1× target")
	}
	bounds := result.Image.Bounds()
	if bounds.Dx() != 30 || bounds.Dy() != 10 {
		t.Fatalf("expected 30×10 image, got %d×%d", bounds.Dx(), bounds.Dy())
	}
}

// TestFractionalScaling verifies that a non-integer scaling ratio (3:2)
// produces output with correct TargetSize dimensions.

func TestFractionalScaling(t *testing.T) {
	result := Render(Config{
		LogicalSize: image.Point{X: 30, Y: 20},
		TargetSize:  image.Point{X: 20, Y: 14},
		Text:        "Hi",
		Font:        font.Default(),
	})
	if result == nil {
		t.Fatal("expected non-nil result for fractional scaling (3:2 ratio)")
	}
	bounds := result.Image.Bounds()
	if bounds.Dx() != 20 || bounds.Dy() != 14 {
		t.Fatalf("expected 20×14 image, got %d×%d", bounds.Dx(), bounds.Dy())
	}
}

// TestUpscaleCase verifies that when LogicalSize is smaller than TargetSize
// (upscale), output has correct TargetSize dimensions.

func TestUpscaleCase(t *testing.T) {
	result := Render(Config{
		LogicalSize: image.Point{X: 20, Y: 10},
		TargetSize:  image.Point{X: 40, Y: 20},
		Text:        "Hi",
		Font:        font.Default(),
	})
	if result == nil {
		t.Fatal("expected non-nil result for upscale case")
	}
	bounds := result.Image.Bounds()
	if bounds.Dx() != 40 || bounds.Dy() != 20 {
		t.Fatalf("expected 40×20 image, got %d×%d", bounds.Dx(), bounds.Dy())
	}
}

// TestEmptyTextTransparent verifies that empty text at scaled resolution
// produces a fully transparent image at TargetSize.

func TestEmptyTextTransparent(t *testing.T) {
	result := Render(Config{
		LogicalSize: image.Point{X: 40, Y: 20},
		TargetSize:  image.Point{X: 20, Y: 10},
		Text:        "",
		Font:        font.Default(),
	})
	if result == nil {
		t.Fatal("expected non-nil result for empty text")
	}
	bounds := result.Image.Bounds()
	if bounds.Dx() != 20 || bounds.Dy() != 10 {
		t.Fatalf("expected 20×10 image, got %d×%d", bounds.Dx(), bounds.Dy())
	}

	// All pixels should be fully transparent.
	for y := 0; y < 10; y++ {
		for x := 0; x < 20; x++ {
			_, _, _, a := result.Image.At(x, y).RGBA()
			if a != 0 {
				t.Fatalf("expected transparent pixel at (%d,%d), got alpha=%d", x, y, a)
			}
		}
	}
}

// TestPositionPassedThrough verifies that the Position field from Config
// is passed through unchanged to the Result.

func TestPositionPassedThrough(t *testing.T) {
	result := Render(Config{
		LogicalSize: image.Point{X: 40, Y: 20},
		TargetSize:  image.Point{X: 20, Y: 10},
		Position:    image.Point{X: 15, Y: 25},
		Text:        "Hi",
		Font:        font.Default(),
	})
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Position != (image.Point{X: 15, Y: 25}) {
		t.Fatalf("expected position (15,25), got %v", result.Position)
	}
}

// TestLabelDefaultsToScaledtextbox verifies that when no Label is provided,
// the Result Label defaults to "scaledtextbox".

func TestLabelDefaultsToScaledtextbox(t *testing.T) {
	result := Render(Config{
		LogicalSize: image.Point{X: 40, Y: 20},
		TargetSize:  image.Point{X: 20, Y: 10},
		Text:        "Hi",
		Font:        font.Default(),
	})
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Label != "scaledtextbox" {
		t.Fatalf("expected label \"scaledtextbox\", got %q", result.Label)
	}
}

// TestBorderAt1px verifies that with 2× downscale and Border=true, the border
// is drawn at 1px at TargetSize (not 2px from scaling a logical border).

func TestBorderAt1px(t *testing.T) {
	fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	result := Render(Config{
		LogicalSize: image.Point{X: 60, Y: 20},
		TargetSize:  image.Point{X: 30, Y: 10},
		Text:        "Hi",
		Font:        font.Default(),
		Border:      true,
		Foreground:  fg,
	})
	if result == nil {
		t.Fatal("expected non-nil result with border")
	}

	img := result.Image.(*image.RGBA)
	width := 30
	height := 10

	// Outer edge pixels should be the foreground color (border drawn at TargetSize).
	for x := 0; x < width; x++ {
		if img.RGBAAt(x, 0) != fg {
			t.Fatalf("top border pixel at (%d,0) should be foreground", x)
		}
		if img.RGBAAt(x, height-1) != fg {
			t.Fatalf("bottom border pixel at (%d,%d) should be foreground", x, height-1)
		}
	}
	for y := 0; y < height; y++ {
		if img.RGBAAt(0, y) != fg {
			t.Fatalf("left border pixel at (0,%d) should be foreground", y)
		}
		if img.RGBAAt(width-1, y) != fg {
			t.Fatalf("right border pixel at (%d,%d) should be foreground", width-1, y)
		}
	}

	// Verify border is exactly 1px: pixels at (1,1) through (28,8) are NOT
	// all border color. At least some interior pixels should differ from fg
	// (either text or transparent — but not uniform border).
	allBorder := true
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			if img.RGBAAt(x, y) != fg {
				allBorder = false
				break
			}
		}
		if !allBorder {
			break
		}
	}
	if allBorder {
		t.Fatal("interior pixels at (1,1)-(28,8) should NOT all be border color; border should be 1px only")
	}
}

// TestBorderWithEmptyText verifies that with Border=true and empty text at
// scaled resolution, the border is drawn at edges and interior is transparent.

func TestBorderWithEmptyText(t *testing.T) {
	fg := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	result := Render(Config{
		LogicalSize: image.Point{X: 40, Y: 20},
		TargetSize:  image.Point{X: 20, Y: 10},
		Text:        "",
		Font:        font.Default(),
		Border:      true,
		Foreground:  fg,
	})
	if result == nil {
		t.Fatal("expected non-nil result with border and empty text")
	}

	img := result.Image.(*image.RGBA)
	width := 20
	height := 10

	// Border pixels should be foreground.
	for x := 0; x < width; x++ {
		if img.RGBAAt(x, 0) != fg {
			t.Fatalf("top border pixel at (%d,0) should be foreground", x)
		}
		if img.RGBAAt(x, height-1) != fg {
			t.Fatalf("bottom border pixel at (%d,%d) should be foreground", x, height-1)
		}
	}
	for y := 0; y < height; y++ {
		if img.RGBAAt(0, y) != fg {
			t.Fatalf("left border pixel at (0,%d) should be foreground", y)
		}
		if img.RGBAAt(width-1, y) != fg {
			t.Fatalf("right border pixel at (%d,%d) should be foreground", width-1, y)
		}
	}

	// Interior should be fully transparent.
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			c := img.RGBAAt(x, y)
			if c.A != 0 {
				t.Fatalf("interior pixel at (%d,%d) should be transparent, got %v", x, y, c)
			}
		}
	}
}
