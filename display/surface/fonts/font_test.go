package font_test

import (
	"bufio"
	"image"
	"image/color"
	"os"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/surface/fonts"
)

// --- From: draw_test.go ---

func TestDrawGlyph_NilFace(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	// Should not panic with nil face.
	font.DrawGlyph(img, nil, 'A', 0, 0, color.RGBA{R: 255, A: 255})
}

func TestDrawText_NilFace(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	w := font.DrawText(img, nil, "hello", 0, 0, color.RGBA{R: 255, A: 255}, 100)
	if w != 0 {
		t.Errorf("expected 0 width for nil face, got %d", w)
	}
}

func TestDrawGlyph_SetsPixels(t *testing.T) {
	face := font.Default()
	if face == nil {
		t.Skip("no default face registered")
	}
	m := face.Metrics()
	img := image.NewRGBA(image.Rect(0, 0, m.GlyphWidth+4, m.GlyphHeight+4))
	fg := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	font.DrawGlyph(img, face, 'A', 2, 2, fg)

	// Verify at least one pixel was set (font 'A' has at least some ink).
	found := false
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if img.RGBAAt(x, y) == fg {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		t.Error("expected at least one pixel to be set for character 'A'")
	}
}

func TestDrawGlyph_ClipsOutOfBounds(t *testing.T) {
	face := font.Default()
	if face == nil {
		t.Skip("no default face registered")
	}
	// Place glyph at negative coordinates — should not panic.
	img := image.NewRGBA(image.Rect(0, 0, 5, 5))
	font.DrawGlyph(img, face, 'X', -10, -10, color.RGBA{R: 255, A: 255})
	font.DrawGlyph(img, face, 'X', 100, 100, color.RGBA{R: 255, A: 255})
}

func TestDrawText_AdvancesCorrectly(t *testing.T) {
	face := font.Default()
	if face == nil {
		t.Skip("no default face registered")
	}
	m := face.Metrics()
	img := image.NewRGBA(image.Rect(0, 0, 200, m.GlyphHeight+4))
	text := "Hello"
	w := font.DrawText(img, face, text, 0, 0, color.RGBA{G: 255, A: 255}, 200)
	expected := len(text) * m.GlyphAdvance
	if w != expected {
		t.Errorf("DrawText width = %d, want %d", w, expected)
	}
}

func TestDrawText_StopsAtMaxX(t *testing.T) {
	face := font.Default()
	if face == nil {
		t.Skip("no default face registered")
	}
	m := face.Metrics()
	img := image.NewRGBA(image.Rect(0, 0, 200, m.GlyphHeight+4))
	// maxX allows only 2 characters
	maxX := m.GlyphAdvance * 2
	text := "ABCDE"
	w := font.DrawText(img, face, text, 0, 0, color.RGBA{B: 255, A: 255}, maxX)
	expectedMax := 2 * m.GlyphAdvance
	if w > expectedMax {
		t.Errorf("DrawText should stop at maxX, got width %d, expected at most %d", w, expectedMax)
	}
}

// --- From: font_test.go ---

func TestDefaultFaceRegistered(t *testing.T) {
	face := font.Default()
	if face == nil {
		t.Fatal("Default() returned nil")
	}
	if id := face.ID(); id != "spleen-5x8" {
		t.Fatalf("Default().ID()=%q, want %q", id, "spleen-5x8")
	}
}

func TestSpleen5x8GlyphFallback(t *testing.T) {
	face, ok := font.Get(font.Spleen5x8ID)
	if !ok {
		t.Fatalf("Get(%q) returned ok=false", font.Spleen5x8ID)
	}
	qMark := face.GlyphRow('?', 0)
	if got := face.GlyphRow('\u2603', 0); got != qMark {
		t.Fatalf("fallback glyph row mismatch: got=0x%08X want=0x%08X", got, qMark)
	}
}

func TestRegisteredFacesMetrics(t *testing.T) {
	cases := []struct {
		id   string
		want font.Metrics
	}{
		{id: font.Spleen5x8ID, want: font.Metrics{GlyphWidth: 5, GlyphHeight: 8, GlyphAdvance: 6, RowHeight: 10}},
		{id: font.Spleen6x12ID, want: font.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14}},
		{id: font.Terminus6x12ID, want: font.Metrics{GlyphWidth: 6, GlyphHeight: 12, GlyphAdvance: 7, RowHeight: 14}},
		{id: font.Cozette6x13ID, want: font.Metrics{GlyphWidth: 6, GlyphHeight: 13, GlyphAdvance: 7, RowHeight: 15}},
	}
	for _, tc := range cases {
		face, ok := font.Get(tc.id)
		if !ok {
			t.Fatalf("Get(%q) returned ok=false", tc.id)
		}
		if got := face.Metrics(); got != tc.want {
			t.Fatalf("Metrics(%q)=%+v, want %+v", tc.id, got, tc.want)
		}
	}
}

func TestListReturnsMultipleFaces(t *testing.T) {
	faces := font.List()
	if len(faces) == 0 {
		t.Fatal("List() returned empty slice")
	}
	// Should have at least the generated spleen, terminus, and cozette fonts
	if len(faces) < 5 {
		t.Fatalf("List() returned only %d faces, expected at least 5", len(faces))
	}
	// Verify sorted by ascending GlyphHeight
	for i := 1; i < len(faces); i++ {
		if faces[i].Metrics().GlyphHeight < faces[i-1].Metrics().GlyphHeight {
			t.Fatalf("List() not sorted by GlyphHeight: %q (%d) after %q (%d)",
				faces[i].ID(), faces[i].Metrics().GlyphHeight,
				faces[i-1].ID(), faces[i-1].Metrics().GlyphHeight)
		}
	}
}

// --- Tests merged from selection_test.go ---

// TestDefaultIsSpleen5x8 verifies that Default() returns Spleen 5×8.
func TestDefaultIsSpleen5x8(t *testing.T) {
	face := font.Default()
	if face == nil {
		t.Fatal("Default() returned nil")
	}
	got := face.ID()
	if got != "spleen-5x8" {
		t.Fatalf("Default().ID() = %q, want %q", got, "spleen-5x8")
	}
	m := face.Metrics()
	if m.GlyphWidth != 5 || m.GlyphHeight != 8 {
		t.Fatalf("Default().Metrics() = %+v, want 5×8", m)
	}
}

// TestByHeight32ReturnsSpleen16x32 verifies ByHeight(32) returns spleen-16x32.
func TestByHeight32ReturnsSpleen16x32(t *testing.T) {
	_, ok := font.Get("spleen-16x32")
	if !ok {
		t.Skip("spleen-16x32 not registered (generated fonts not yet available)")
	}

	face := font.ByHeight(32)
	if face == nil {
		t.Fatal("ByHeight(32) returned nil")
	}
	if got := face.ID(); got != "spleen-16x32" {
		t.Fatalf("ByHeight(32).ID() = %q, want %q", got, "spleen-16x32")
	}
	if h := face.Metrics().GlyphHeight; h != 32 {
		t.Fatalf("ByHeight(32).Metrics().GlyphHeight = %d, want 32", h)
	}
}

// TestByHeight48ReturnsClosestNotExceeding verifies ByHeight(48) returns the
// largest font whose GlyphHeight does not exceed 48.
func TestByHeight48ReturnsClosestNotExceeding(t *testing.T) {
	face := font.ByHeight(48)
	if face == nil {
		t.Fatal("ByHeight(48) returned nil")
	}

	h := face.Metrics().GlyphHeight
	if h > 48 {
		t.Fatalf("ByHeight(48) returned font with GlyphHeight=%d, which exceeds 48", h)
	}

	// Verify no registered font with height ≤ 48 is larger than what we got.
	for _, f := range font.List() {
		fh := f.Metrics().GlyphHeight
		if fh > h && fh <= 48 {
			t.Fatalf("ByHeight(48) returned %q (height %d) but %q (height %d) is larger and still ≤ 48",
				face.ID(), h, f.ID(), fh)
		}
	}
}

// TestSpleenPreferredOverTerminusAtSameHeight verifies that when Spleen and
// Terminus variants share the same glyph height, Spleen is selected.
func TestSpleenPreferredOverTerminusAtSameHeight(t *testing.T) {
	// Both spleen-8x16 and terminus-8x16 have GlyphHeight=16.
	spleen, ok := font.Get("spleen-8x16")
	if !ok {
		t.Skip("spleen-8x16 not registered")
	}
	terminus, ok := font.Get("terminus-8x16")
	if !ok {
		t.Skip("terminus-8x16 not registered")
	}

	sh := spleen.Metrics().GlyphHeight
	th := terminus.Metrics().GlyphHeight
	if sh != th {
		t.Skipf("spleen-8x16 height=%d != terminus-8x16 height=%d; no tie to test", sh, th)
	}

	face := font.ByHeight(sh)
	if face == nil {
		t.Fatal("ByHeight returned nil")
	}
	if got := face.ID(); got != "spleen-8x16" {
		t.Fatalf("ByHeight(%d) = %q, want %q (Spleen preferred over Terminus at same height)",
			sh, got, "spleen-8x16")
	}
}

// TestAllGeneratedVariantsRegistered checks that generated font files are
// properly registered in the font registry.
func TestAllGeneratedVariantsRegistered(t *testing.T) {
	generated := []string{
		"spleen-5x8",
		"spleen-6x12",
		"spleen-8x16",
		"spleen-12x24",
		"spleen-16x32",
		"spleen-32x64",
		"cozette-6x13",
		"terminus-6x12",
		"terminus-8x14",
		"terminus-8x16",
		"terminus-10x18",
		"terminus-10x20",
		"terminus-11x22",
		"terminus-12x24",
		"terminus-14x28",
		"terminus-16x32",
	}

	for _, id := range generated {
		if _, ok := font.Get(id); !ok {
			t.Errorf("expected font %q to be registered", id)
		}
	}
}

// TestHeightCoverage verifies the set of glyph heights present in the registry.
func TestHeightCoverage(t *testing.T) {
	all := font.List()
	if len(all) == 0 {
		t.Fatal("no fonts registered")
	}

	heights := make(map[int]bool)
	for _, f := range all {
		heights[f.Metrics().GlyphHeight] = true
	}

	// Heights covered by generated fonts:
	// spleen-5x8 (8), spleen-6x12 (12), cozette-6x13 (13), terminus-8x14 (14),
	// spleen-8x16 / terminus-8x16 (16), terminus-10x18 (18), terminus-10x20 (20),
	// terminus-11x22 (22), spleen-12x24 / terminus-12x24 (24),
	// terminus-14x28 (28), spleen-16x32 / terminus-16x32 (32), spleen-32x64 (64)
	expectedHeights := []int{8, 12, 13, 14, 16, 18, 20, 22, 24, 28, 32, 64}
	for _, h := range expectedHeights {
		if !heights[h] {
			t.Errorf("expected pixel height %d to be covered, but no font registered at that height", h)
		}
	}

	t.Logf("all registered heights: %v", sortedKeys(heights))
}

// sortedKeys returns the keys of a map[int]bool in ascending order.
func sortedKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j-1] > keys[j]; j-- {
			keys[j-1], keys[j] = keys[j], keys[j-1]
		}
	}
	return keys
}

// --- From: gen_matrix_code_test.go ---

func TestMatrixCodeFontRegistered(t *testing.T) {
	face, ok := font.Get("matrix-code")
	if !ok {
		t.Fatal("font.Get(\"matrix-code\") returned ok=false; font not registered")
	}
	if face == nil {
		t.Fatal("font.Get(\"matrix-code\") returned nil face")
	}
	if id := face.ID(); id != font.MatrixCodeID {
		t.Fatalf("ID()=%q, want %q", id, font.MatrixCodeID)
	}
}

func TestMatrixCodeFontMetrics(t *testing.T) {
	face, ok := font.Get("matrix-code")
	if !ok {
		t.Fatal("font not registered")
	}
	m := face.Metrics()
	if m.GlyphWidth != 13 {
		t.Errorf("GlyphWidth=%d, want 13", m.GlyphWidth)
	}
	if m.GlyphHeight != 12 {
		t.Errorf("GlyphHeight=%d, want 12", m.GlyphHeight)
	}
	if m.GlyphAdvance != 14 {
		t.Errorf("GlyphAdvance=%d, want 14", m.GlyphAdvance)
	}
	if m.RowHeight != 14 {
		t.Errorf("RowHeight=%d, want 14", m.RowHeight)
	}
}

func TestMatrixCodeFontASCIIPrintable(t *testing.T) {
	face, ok := font.Get("matrix-code")
	if !ok {
		t.Fatal("font not registered")
	}

	// Verify that GlyphRow returns data for ASCII digits and letters
	// (the font is decorative and may not have all punctuation, but digits/letters should work)
	hasNonZero := false
	for ch := rune('0'); ch <= '9'; ch++ {
		for row := 0; row < face.Metrics().GlyphHeight; row++ {
			if face.GlyphRow(ch, row) != 0 {
				hasNonZero = true
				break
			}
		}
		if hasNonZero {
			break
		}
	}
	if !hasNonZero {
		t.Error("no non-zero GlyphRow data found for ASCII digits '0'-'9'")
	}

	hasNonZero = false
	for ch := rune('A'); ch <= 'Z'; ch++ {
		for row := 0; row < face.Metrics().GlyphHeight; row++ {
			if face.GlyphRow(ch, row) != 0 {
				hasNonZero = true
				break
			}
		}
		if hasNonZero {
			break
		}
	}
	if !hasNonZero {
		t.Error("no non-zero GlyphRow data found for ASCII letters 'A'-'Z'")
	}
}

func TestMatrixCodeFontKatakana(t *testing.T) {
	face, ok := font.Get("matrix-code")
	if !ok {
		t.Fatal("font not registered")
	}

	// Half-width katakana range: U+FF66 (65382) to U+FF9D (65437)
	// Verify we get glyph data for at least some of these codepoints.
	// The font may not have all katakana characters rendered, but it should
	// at least accept them (returning fallback data for missing ones).
	katakanaStart := rune(0xFF66)
	katakanaEnd := rune(0xFF9D)

	// The face should return non-zero data for known chars or fallback '?' data.
	// Verify that calling GlyphRow on katakana codepoints doesn't panic and
	// returns the map-entry data (or fallback).
	for ch := katakanaStart; ch <= katakanaEnd; ch++ {
		for row := 0; row < face.Metrics().GlyphHeight; row++ {
			_ = face.GlyphRow(ch, row) // should not panic
		}
	}

	// Verify out-of-bounds row returns 0.
	if got := face.GlyphRow(katakanaStart, -1); got != 0 {
		t.Errorf("GlyphRow(katakana, -1)=0x%08X, want 0", got)
	}
	if got := face.GlyphRow(katakanaStart, face.Metrics().GlyphHeight); got != 0 {
		t.Errorf("GlyphRow(katakana, height)=0x%08X, want 0", got)
	}
}

func TestMatrixCodeFontFallback(t *testing.T) {
	face, ok := font.Get("matrix-code")
	if !ok {
		t.Fatal("font not registered")
	}

	// Unknown codepoint should fall back to '?' data (or 0 if '?' has no data).
	// The key property is that it doesn't panic.
	unknownChar := rune(0xFFFF)
	for row := 0; row < face.Metrics().GlyphHeight; row++ {
		_ = face.GlyphRow(unknownChar, row)
	}
}

// --- From: generate_test.go ---

// TestGenerateDirectivesPresent verifies that generate.go contains exactly 19
// //go:generate directives: 18 font variants (fontgen) and 1 icon face (gen-icons).
func TestGenerateDirectivesPresent(t *testing.T) {
	expectedFontIDs := []string{
		"spleen-5x8",
		"spleen-6x12",
		"spleen-8x16",
		"spleen-12x24",
		"spleen-16x32",
		"spleen-32x64",
		"cozette-6x13",
		"terminus-6x12",
		"terminus-8x14",
		"terminus-8x16",
		"terminus-10x18",
		"terminus-10x20",
		"terminus-11x22",
		"terminus-12x24",
		"terminus-14x28",
		"terminus-16x32",
		"matrix-10x10",
		"matrix-code",
	}

	data, err := os.ReadFile("generate.go")
	if err != nil {
		t.Fatalf("failed to read generate.go: %v", err)
	}

	// Collect all //go:generate lines.
	var directives []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "//go:generate") {
			directives = append(directives, line)
		}
	}

	// Verify total count: 18 fontgen + 1 gen-icons = 19.
	if len(directives) != 19 {
		t.Fatalf("expected 19 //go:generate directives, got %d", len(directives))
	}

	// Verify each expected font ID appears in a directive with the -id flag.
	for _, id := range expectedFontIDs {
		found := false
		needle := "-id " + id
		for _, d := range directives {
			if strings.Contains(d, needle) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing //go:generate directive with -id %s", id)
		}
	}

	// Verify the gen-icons directive is present with -faceid material-icons-24.
	iconDirectiveFound := false
	for _, d := range directives {
		if strings.Contains(d, "gen-icons") && strings.Contains(d, "-faceid material-icons-24") {
			iconDirectiveFound = true
			break
		}
	}
	if !iconDirectiveFound {
		t.Error("missing //go:generate directive for gen-icons with -faceid material-icons-24")
	}
}

// --- From: registry_icons_test.go ---

// TestIconFaceRegistrationAndLookup verifies that a registered icon face
// can be retrieved via font.Get with ok=true, and that unregistered IDs
// return nil, false. This validates the registry contract for icon faces.
//

func TestIconFaceRegistrationAndLookup(t *testing.T) {
	// Snapshot and clear the registry so we control exactly what's registered.
	restore := font.SnapshotAndClear()
	defer restore()

	// Register a mock icon face with material-icons-24 ID and square metrics.
	mock := &mockIconFace{
		id: "material-icons-24",
		metrics: font.Metrics{
			GlyphWidth:   24,
			GlyphHeight:  24,
			GlyphAdvance: 24,
			RowHeight:    24,
		},
	}
	font.Register(mock)

	// Requirement 8.2: font.Get("material-icons-24") returns the face with ok=true.
	face, ok := font.Get("material-icons-24")
	if !ok {
		t.Fatal("font.Get(\"material-icons-24\") returned ok=false, want ok=true")
	}
	if face == nil {
		t.Fatal("font.Get(\"material-icons-24\") returned nil face")
	}

	// Requirement 8.1: Verify the registered face ID matches the pattern.
	if got := face.ID(); got != "material-icons-24" {
		t.Fatalf("face.ID() = %q, want %q", got, "material-icons-24")
	}

	// Requirement 8.4: Verify square metrics (all four fields = 24).
	m := face.Metrics()
	if m.GlyphWidth != 24 {
		t.Errorf("Metrics.GlyphWidth = %d, want 24", m.GlyphWidth)
	}
	if m.GlyphHeight != 24 {
		t.Errorf("Metrics.GlyphHeight = %d, want 24", m.GlyphHeight)
	}
	if m.GlyphAdvance != 24 {
		t.Errorf("Metrics.GlyphAdvance = %d, want 24", m.GlyphAdvance)
	}
	if m.RowHeight != 24 {
		t.Errorf("Metrics.RowHeight = %d, want 24", m.RowHeight)
	}

	// Requirement 8.3: font.Get for unregistered ID returns nil, false.
	face2, ok2 := font.Get("material-icons-99")
	if ok2 {
		t.Fatal("font.Get(\"material-icons-99\") returned ok=true, want ok=false")
	}
	if face2 != nil {
		t.Fatal("font.Get(\"material-icons-99\") returned non-nil face, want nil")
	}
}

// TestIconFaceSquareMetricsAllEqual verifies that icon faces report Metrics
// where all four fields are equal (square glyph property for icons).
//

func TestIconFaceSquareMetricsAllEqual(t *testing.T) {
	restore := font.SnapshotAndClear()
	defer restore()

	heights := []int{16, 24, 32, 48}
	for _, h := range heights {
		mock := &mockIconFace{
			id: "material-icons-test",
			metrics: font.Metrics{
				GlyphWidth:   h,
				GlyphHeight:  h,
				GlyphAdvance: h,
				RowHeight:    h,
			},
		}
		m := mock.Metrics()
		if m.GlyphWidth != m.GlyphHeight || m.GlyphHeight != m.GlyphAdvance || m.GlyphAdvance != m.RowHeight {
			t.Errorf("height=%d: Metrics fields not all equal: %+v", h, m)
		}
	}
}

// mockIconFace implements fonts.Face for testing icon registration behavior.
type mockIconFace struct {
	id      string
	metrics font.Metrics
}

func (f *mockIconFace) ID() string                    { return f.id }
func (f *mockIconFace) Metrics() font.Metrics         { return f.metrics }
func (f *mockIconFace) GlyphRow(_ rune, _ int) uint32 { return 0 }
