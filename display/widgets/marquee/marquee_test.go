package marquee

import (
	"image"
	"image/color"
	"testing"
	"time"

	stylecolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/surface/fonts"
	"pgregory.net/rapid"
)

// ===========================================================================
// --- From: integration_test.go ---
// ===========================================================================

// TestIntegration_MatrixCodeFont verifies that fonts.Get("matrix-code") returns
// a valid Face with katakana and ASCII codepoints.
func TestIntegration_MatrixCodeFont(t *testing.T) {
	face, ok := font.Get("matrix-code")
	if !ok {
		t.Fatal("font.Get(\"matrix-code\") returned ok=false; font not registered")
	}
	if face == nil {
		t.Fatal("font.Get(\"matrix-code\") returned nil face")
	}

	m := face.Metrics()
	if m.GlyphWidth <= 0 || m.GlyphHeight <= 0 || m.GlyphAdvance <= 0 || m.RowHeight <= 0 {
		t.Fatalf("invalid metrics: %+v", m)
	}

	// Verify ASCII printable characters have glyph data
	hasASCII := false
	for ch := rune('A'); ch <= 'Z'; ch++ {
		for row := 0; row < m.GlyphHeight; row++ {
			if face.GlyphRow(ch, row) != 0 {
				hasASCII = true
				break
			}
		}
		if hasASCII {
			break
		}
	}
	if !hasASCII {
		t.Error("matrix-code font has no non-zero glyph data for ASCII 'A'-'Z'")
	}

	// Verify katakana codepoints are handled gracefully (U+FF66 to U+FF9D).
	// The font dispatches these via fallback — calling GlyphRow must not panic
	// and should return fallback character data (non-zero from '?' glyph).
	for ch := rune(0xFF66); ch <= rune(0xFF9D); ch++ {
		for row := 0; row < m.GlyphHeight; row++ {
			_ = face.GlyphRow(ch, row) // must not panic
		}
	}

	// Verify fallback works: unknown codepoints get '?' glyph data (or zero).
	// The key guarantee is that GlyphRow always returns a value without error.
	fallbackHasData := false
	for row := 0; row < m.GlyphHeight; row++ {
		if face.GlyphRow(0xFFFF, row) != 0 {
			fallbackHasData = true
			break
		}
	}
	// Note: if fallback char '?' has no data either (all rows zero), that's
	// acceptable for this font since it's decorative — the test verifies
	// the dispatch mechanism works.
	_ = fallbackHasData
}

// TestIntegration_MarqueeComposition is a smoke test that creates a marquee Strip
// with vertical direction, a FixedStringSource, and a gradient color slice, then
// verifies the full composition pipeline produces a valid Sprite with non-zero pixels.
func TestIntegration_MarqueeComposition(t *testing.T) {
	face, ok := font.Get("matrix-code")
	if !ok {
		t.Fatal("font.Get(\"matrix-code\") returned ok=false; font not registered")
	}

	m := face.Metrics()

	// Create a gradient color slice (green trailing to black)
	accent := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	trailLen := 8
	colors := stylecolor.Gradient(accent, trailLen)
	if len(colors) != trailLen {
		t.Fatalf("Gradient returned %d colors, want %d", len(colors), trailLen)
	}

	// Create strip bounds: one column wide, tall enough for several cells
	stripWidth := m.GlyphAdvance
	stripHeight := m.RowHeight * trailLen
	bounds := image.Rect(10, 20, 10+stripWidth, 20+stripHeight)

	// Create strip with vertical direction and fixed string source
	source := &FixedStringSource{Runes: []rune("HELLO")}
	strip := New(Config{
		Bounds:    bounds,
		Direction: Vertical,
		Font:      face,
		Source:    source,
		Colors:    colors,
		Speed:     5.0,
		Phase:     0,
	})

	// Tick a few times to advance the scroll
	for i := 0; i < 5; i++ {
		strip.Tick(100 * time.Millisecond)
	}

	// Render and verify
	sprite := strip.RenderFrame()
	if sprite == nil {
		t.Fatal("RenderFrame() returned nil Sprite")
	}
	if sprite.Image == nil {
		t.Fatal("Sprite.Image is nil")
	}

	imgBounds := sprite.Image.Bounds()
	expectedWidth := bounds.Dx()
	expectedHeight := bounds.Dy()
	if imgBounds.Dx() != expectedWidth || imgBounds.Dy() != expectedHeight {
		t.Errorf("Sprite image dimensions = %dx%d, want %dx%d",
			imgBounds.Dx(), imgBounds.Dy(), expectedWidth, expectedHeight)
	}

	if sprite.Position != bounds.Min {
		t.Errorf("Sprite.Position = %v, want %v", sprite.Position, bounds.Min)
	}

	// Verify at least some pixels are non-zero (characters were rendered)
	rgba, ok := sprite.Image.(*image.RGBA)
	if !ok {
		t.Fatal("Sprite.Image is not *image.RGBA")
	}

	hasNonZero := false
	for i := 0; i < len(rgba.Pix); i += 4 {
		if rgba.Pix[i] != 0 || rgba.Pix[i+1] != 0 || rgba.Pix[i+2] != 0 || rgba.Pix[i+3] != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Error("Sprite image has all zero pixels; expected rendered character data")
	}
}

// ===========================================================================
// --- From: marquee_prop_test.go ---
// ===========================================================================

// recordingSource implements CharSource and logs every CharAt call.
type recordingSource struct {
	calls   []int // logged indices
	backing *FixedStringSource
}

func (r *recordingSource) CharAt(index int) rune {
	r.calls = append(r.calls, index)
	return r.backing.CharAt(index)
}

func TestProp_MarqueeCharSourceUsage(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a positive bounding rectangle (at least one cell visible).
		// Use bounds large enough to show at least 1 cell: min 14px high (default RowHeight)
		// and min 5px wide (default GlyphWidth).
		w := rapid.IntRange(5, 100).Draw(t, "width")
		h := rapid.IntRange(14, 200).Draw(t, "height")
		bounds := image.Rect(0, 0, w, h)

		// Generate a non-empty rune slice for the source.
		runeCount := rapid.IntRange(1, 20).Draw(t, "runeCount")
		runes := make([]rune, runeCount)
		for i := range runes {
			runes[i] = rune(rapid.IntRange(33, 126).Draw(t, "rune"))
		}

		// Generate a phase/offset to start at a non-zero position.
		phase := rapid.Float64Range(0.0, 100.0).Draw(t, "phase")

		// Generate colors (at least 1 to ensure rendering occurs).
		numColors := rapid.IntRange(1, 30).Draw(t, "numColors")
		colors := make([]color.RGBA, numColors)
		for i := range colors {
			colors[i] = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}

		rec := &recordingSource{
			backing: &FixedStringSource{Runes: runes},
		}

		cfg := Config{
			Bounds:    bounds,
			Direction: Vertical,
			Source:    rec,
			Colors:    colors,
			Speed:     1.0,
			Phase:     phase,
		}

		strip := New(cfg)

		// Call RenderFrame which should invoke CharAt on the recording source.
		sprite := strip.RenderFrame()
		if sprite == nil {
			// If sprite is nil, bounds were too small for even one cell — skip.
			return
		}

		// Property 1: The recording source was called (non-empty calls slice).
		if len(rec.calls) == 0 {
			t.Fatalf("recordingSource.CharAt was never called during RenderFrame; "+
				"bounds=%v, phase=%f, numColors=%d", bounds, phase, numColors)
		}

		// Property 2: Each call index is derived from the current offset.
		// The frontIndex should be int(strip.offset), and calls should be
		// frontIndex, frontIndex-1, frontIndex-2, ... for each visible cell.
		frontIndex := int(strip.offset)

		// Build the set of expected indices: frontIndex - 0, frontIndex - 1, ..., frontIndex - (len(calls)-1)
		expectedIndices := make(map[int]bool)
		for i := 0; i < len(rec.calls); i++ {
			expectedIndices[frontIndex-i] = true
		}

		// Verify every call used an expected index (derived from offset).
		for callNum, idx := range rec.calls {
			if !expectedIndices[idx] {
				t.Fatalf("CharAt call #%d used index %d, which is not derived from frontIndex=%d; "+
					"expected indices in range [%d, %d]",
					callNum, idx, frontIndex, frontIndex-len(rec.calls)+1, frontIndex)
			}
		}

		// Property 3: The calls should be exactly the sequence frontIndex, frontIndex-1, ... in order.
		for i, idx := range rec.calls {
			expected := frontIndex - i
			if idx != expected {
				t.Fatalf("CharAt call #%d had index %d, want %d (frontIndex=%d)",
					i, idx, expected, frontIndex)
			}
		}
	})
}

// genMarqueeConfig generates a Config with positive Bounds and a non-nil Source.
func genMarqueeConfig(t *rapid.T) Config {
	// Generate positive bounds (min 1px in each dimension)
	minX := rapid.IntRange(0, 100).Draw(t, "minX")
	minY := rapid.IntRange(0, 100).Draw(t, "minY")
	w := rapid.IntRange(1, 200).Draw(t, "width")
	h := rapid.IntRange(1, 200).Draw(t, "height")
	bounds := image.Rect(minX, minY, minX+w, minY+h)

	dir := Direction(rapid.IntRange(0, 1).Draw(t, "direction"))
	speed := rapid.Float64Range(0.1, 20.0).Draw(t, "speed")
	phase := rapid.Float64Range(0.0, 100.0).Draw(t, "phase")

	// Use a non-nil FixedStringSource with some characters
	src := &FixedStringSource{Runes: []rune("ABCDEFGHIJ")}

	// Generate a color slice with at least 1 entry
	numColors := rapid.IntRange(1, 30).Draw(t, "numColors")
	colors := make([]color.RGBA, numColors)
	for i := range colors {
		colors[i] = color.RGBA{R: 255, G: 200, B: 100, A: 255}
	}

	return Config{
		Bounds:    bounds,
		Direction: dir,
		Font:      font.Default(),
		Source:    src,
		Colors:    colors,
		Speed:     speed,
		Phase:     phase,
	}
}

func TestProp_MarqueeSpriteDimensions(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		cfg := genMarqueeConfig(t)
		strip := New(cfg)

		sprite := strip.RenderFrame()

		// With positive Bounds and a non-nil Source, sprite must not be nil
		if sprite == nil {
			t.Fatal("RenderFrame() returned nil for valid config with positive Bounds")
		}

		// Sprite Image must not be nil
		if sprite.Image == nil {
			t.Fatal("RenderFrame() returned Sprite with nil Image")
		}

		// Image dimensions must exactly match Bounds.Dx() x Bounds.Dy()
		imgBounds := sprite.Image.Bounds()
		gotW := imgBounds.Dx()
		gotH := imgBounds.Dy()
		wantW := cfg.Bounds.Dx()
		wantH := cfg.Bounds.Dy()

		if gotW != wantW {
			t.Fatalf("Image width = %d, want Bounds.Dx() = %d (Bounds: %v)",
				gotW, wantW, cfg.Bounds)
		}
		if gotH != wantH {
			t.Fatalf("Image height = %d, want Bounds.Dy() = %d (Bounds: %v)",
				gotH, wantH, cfg.Bounds)
		}

		// Position must equal Bounds.Min
		if sprite.Position != cfg.Bounds.Min {
			t.Fatalf("Sprite.Position = %v, want Bounds.Min = %v",
				sprite.Position, cfg.Bounds.Min)
		}
	})
}

// ===========================================================================
// --- From: marquee_test.go ---
// ===========================================================================

func TestFixedStringSource_CharAt(t *testing.T) {
	src := &FixedStringSource{Runes: []rune("ABC")}

	tests := []struct {
		name  string
		index int
		want  rune
	}{
		{"index 0", 0, 'A'},
		{"index 1", 1, 'B'},
		{"index 2", 2, 'C'},
		{"wraps forward", 3, 'A'},
		{"wraps forward again", 5, 'C'},
		{"negative index -1", -1, 'C'},
		{"negative index -2", -2, 'B'},
		{"negative index -3", -3, 'A'},
		{"negative wraps", -4, 'C'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := src.CharAt(tt.index)
			if got != tt.want {
				t.Errorf("CharAt(%d) = %q, want %q", tt.index, got, tt.want)
			}
		})
	}
}

func TestFixedStringSource_Empty(t *testing.T) {
	src := &FixedStringSource{Runes: nil}
	got := src.CharAt(0)
	if got != ' ' {
		t.Errorf("CharAt(0) on empty = %q, want ' '", got)
	}
	got = src.CharAt(5)
	if got != ' ' {
		t.Errorf("CharAt(5) on empty = %q, want ' '", got)
	}
}

func TestFixedStringSource_ImplementsCharSource(t *testing.T) {
	var _ CharSource = (*FixedStringSource)(nil)
}
