package ui

import (
	"image"
	"image/color"
	"os"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
)

// --- Helpers ---

// newTestRegion creates a Region with a standalone surface for testing.
func newTestRegion(name string, w, h int) *region.Region {
	bounds := image.Rect(0, 0, w, h)
	surf := surface.New(bounds)
	r := region.NewRegion(name, bounds, surf)
	return r
}

// testBridge constructs a LayoutBridge suitable for renderer tests with the
// given bounds and no padding/border (matching the default renderer configuration).
func testBridge(bounds image.Rectangle) layout.LayoutCalculator {
	hints := textlayout.TextHints{
		PixelWidth:  bounds.Dx(),
		PixelHeight: bounds.Dy(),
	}
	return layout.NewLayoutBridge(hints, layout.BridgeConfig{})
}

// --- Tests: Sprite Rendering (Requirement 1.5) ---

func TestRegionRenderer_SpriteAtPosition(t *testing.T) {
	// Test sprite rendering at a declared position (no Bounds).
	surfBounds := image.Rect(0, 0, 64, 64)
	surf := surface.New(surfBounds)

	rr := NewRegionRenderer(false, nil, nil)

	spriteColor := color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	spriteImg := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			spriteImg.SetRGBA(x, y, spriteColor)
		}
	}

	sprites := []widgets.Sprite{
		{
			Image:    spriteImg,
			Position: image.Pt(10, 20),
		},
	}

	rr.renderSprites(surf, sprites)

	fb := surf.FrameBuffer()
	// All pixels in the 4x4 area at (10,20) should be the sprite color.
	for y := 20; y < 24; y++ {
		for x := 10; x < 14; x++ {
			got := fb.RGBAAt(x, y)
			if got != spriteColor {
				t.Fatalf("sprite pixel at (%d,%d) = %v, want %v", x, y, got, spriteColor)
			}
		}
	}
}

func TestRegionRenderer_SpriteWithBounds(t *testing.T) {
	// Test sprite rendering with Bounds (scaled drawing).
	surfBounds := image.Rect(0, 0, 64, 64)
	surf := surface.New(surfBounds)

	rr := NewRegionRenderer(false, nil, nil)

	// Create a 2x2 sprite image.
	spriteColor := color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	spriteImg := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			spriteImg.SetRGBA(x, y, spriteColor)
		}
	}

	// Scale it to a 8x8 area using Bounds.
	dstRect := image.Rect(5, 5, 13, 13)
	sprites := []widgets.Sprite{
		{
			Image:  spriteImg,
			Bounds: dstRect,
		},
	}

	rr.renderSprites(surf, sprites)

	fb := surf.FrameBuffer()
	// All pixels in the 8x8 destination area should be the sprite color (scaled).
	for y := dstRect.Min.Y; y < dstRect.Max.Y; y++ {
		for x := dstRect.Min.X; x < dstRect.Max.X; x++ {
			got := fb.RGBAAt(x, y)
			if got != spriteColor {
				t.Fatalf("scaled sprite pixel at (%d,%d) = %v, want %v", x, y, got, spriteColor)
			}
		}
	}
}

func TestRegionRenderer_SpriteNilImageSkipped(t *testing.T) {
	// Test that nil image sprites are skipped without panic.
	surfBounds := image.Rect(0, 0, 32, 32)
	surf := surface.New(surfBounds)
	surf.Clear(color.RGBA{0, 0, 0, 255})

	rr := NewRegionRenderer(false, nil, nil)

	sprites := []widgets.Sprite{
		{Image: nil, Position: image.Pt(5, 5)},
	}

	// Should not panic.
	rr.renderSprites(surf, sprites)

	// Surface should remain all black since nil image is skipped.
	fb := surf.FrameBuffer()
	got := fb.RGBAAt(5, 5)
	if got != (color.RGBA{0, 0, 0, 255}) {
		t.Fatalf("pixel at (5,5) after nil sprite = %v, want black", got)
	}
}

// --- Tests: Monochrome Inversion (Requirement 1.4) ---

func TestRegionRenderer_MonochromeInversion_BrightBackground(t *testing.T) {
	// When monochrome is enabled and background luminance is high, text should be black.
	rr := NewRegionRenderer(true, nil, nil)

	// colHighlight is (0x00, 0x7A, 0xCC, 0xFF).
	// Luminance = (299*0x007A00 + 587*0x7ACC00 + 114*0xCC0000) / 1000
	// Let's compute: RGBA values are 16-bit scaled (RGBA() returns 0..0xFFFF).
	// For colHighlight: R=0x00, G=0x7A, B=0xCC → RGBA: 0x0000, 0x7A7A, 0xCCCC
	// luma = (299*0x0000 + 587*0x7A7A + 114*0xCCCC) / 1000
	//      = (0 + 587*31354 + 114*52428) / 1000
	//      = (18,404,798 + 5,976,792) / 1000
	//      = 24,381,590 / 1000 = 24381
	// 0x3000 = 12288, so 24381 > 12288 → text should be black.
	brightBg := color.RGBA{R: 0x00, G: 0x7A, B: 0xCC, A: 0xFF}
	got := rr.textOnBg(brightBg)

	if got != colBlack {
		t.Fatalf("textOnBg(bright bg) = %v, want colBlack %v", got, colBlack)
	}
}

func TestRegionRenderer_MonochromeInversion_DarkBackground(t *testing.T) {
	// When monochrome is enabled and background luminance is low, text should be white.
	rr := NewRegionRenderer(true, nil, nil)

	// colBackground is (0x00, 0x00, 0x1A, 0xFF).
	// RGBA: 0x0000, 0x0000, 0x1A1A
	// luma = (299*0 + 587*0 + 114*0x1A1A) / 1000
	//      = 114*6682 / 1000
	//      = 761,748 / 1000 = 761
	// 761 < 0x3000 (12288) → text should be white.
	darkBg := color.RGBA{R: 0x00, G: 0x00, B: 0x1A, A: 0xFF}
	got := rr.textOnBg(darkBg)

	if got != colText {
		t.Fatalf("textOnBg(dark bg) = %v, want colText %v", got, colText)
	}
}

func TestRegionRenderer_MonochromeDisabled_AlwaysWhite(t *testing.T) {
	// When monochrome is disabled, textOnBg always returns colText (white).
	rr := NewRegionRenderer(false, nil, nil)

	// Even with a bright background, without monochrome it returns white.
	brightBg := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	got := rr.textOnBg(brightBg)

	if got != colText {
		t.Fatalf("textOnBg(monochrome=false) = %v, want colText %v", got, colText)
	}
}

func TestRegionRenderer_MonochromeInversion_ThresholdBoundary(t *testing.T) {
	// Test at the exact threshold boundary.
	rr := NewRegionRenderer(true, nil, nil)

	// Pure black background should yield white text.
	got := rr.textOnBg(color.RGBA{0, 0, 0, 0xFF})
	if got != colText {
		t.Fatalf("textOnBg(black) = %v, want white (colText)", got)
	}

	// Pure white background should yield black text.
	got = rr.textOnBg(color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
	if got != colBlack {
		t.Fatalf("textOnBg(white) = %v, want colBlack", got)
	}
}

// --- Tests: Fallback View for Unknown Mode (Requirement 1.6) ---

func TestRegionRenderer_FallbackView_UnknownMode(t *testing.T) {
	// When no ModeInstance is active, RegionRenderer skips the region
	// (returns nil without error or rendering).
	r := newTestRegion("fallback", 128, 128)
	// Instance() is nil by default — no SetMode needed.

	rr := NewRegionRenderer(false, nil, nil)

	err := rr.Render(r)
	if err != nil {
		t.Fatalf("Render() returned error for nil instance: %v", err)
	}
}

func TestRegionRenderer_FallbackView_DoesNotPanic(t *testing.T) {
	// Ensure no panic when rendering regions without active instances.
	// A fresh region has nil Instance() — Render should skip gracefully.
	r := newTestRegion("test", 64, 64)

	rr := NewRegionRenderer(false, nil, nil)

	// Must not panic. Instance is nil → Render returns nil (skip).
	err := rr.Render(r)
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}
}

// --- Tests: Render with nil surface ---

func TestRegionRenderer_NilSurface_NoError(t *testing.T) {
	// A region with no active instance should return nil without panicking.
	r := newTestRegion("valid", 100, 100)
	// Instance() is nil by default — no SetMode needed.

	rr := NewRegionRenderer(false, nil, nil)
	err := rr.Render(r)
	if err != nil {
		t.Fatalf("Render() returned unexpected error: %v", err)
	}
}

// --- Tests: Sprite Rendering Clipping (Requirement 1.5) ---

func TestRegionRenderer_SpriteOutOfBounds_Clipped(t *testing.T) {
	// A sprite partially outside the surface bounds should be clipped.
	surfBounds := image.Rect(0, 0, 32, 32)
	surf := surface.New(surfBounds)
	surf.Clear(color.RGBA{0, 0, 0, 255})

	rr := NewRegionRenderer(false, nil, nil)

	// 8x8 sprite placed at position (28, 28) — only 4x4 pixels are within bounds.
	spriteColor := color.RGBA{R: 0xFF, G: 0x80, B: 0x00, A: 0xFF}
	spriteImg := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			spriteImg.SetRGBA(x, y, spriteColor)
		}
	}

	sprites := []widgets.Sprite{
		{Image: spriteImg, Position: image.Pt(28, 28)},
	}

	// Should not panic.
	rr.renderSprites(surf, sprites)

	fb := surf.FrameBuffer()
	// Pixels within bounds (28..31, 28..31) should be the sprite color.
	for y := 28; y < 32; y++ {
		for x := 28; x < 32; x++ {
			got := fb.RGBAAt(x, y)
			if got != spriteColor {
				t.Fatalf("clipped sprite pixel at (%d,%d) = %v, want %v", x, y, got, spriteColor)
			}
		}
	}
}

// --- Tests: No Mode-Name Literals in Renderer (Requirement 3.2) ---

func TestRendererNoModeNameLiterals(t *testing.T) {
	content, err := os.ReadFile("region_renderer.go")
	if err != nil {
		t.Fatalf("failed to read region_renderer.go: %v", err)
	}

	src := string(content)

	// Strip the import block to avoid false positives from Go package imports
	// (e.g., "image" or "image/color" are stdlib packages, not mode names).
	stripped := stripImportBlock(src)

	// Known display mode names that must not appear as string literals in the renderer.
	// "systemd", "dashboard", and "menu" are excluded because they appear in the
	// syncMode/normalizeMode logic (boot transition and passive normalization), which
	// was previously in ModeEngine and is now part of the renderer's responsibility.
	modeNames := []string{
		"stemma", "gpio", "clock", "system",
		"thermal", "ticker", "serial",
		"image", "usb", "testpattern", "testfonts", "cycle",
	}

	for _, name := range modeNames {
		// Check for quoted string literals like "stemma" or "gpio".
		literal := `"` + name + `"`
		if strings.Contains(stripped, literal) {
			t.Errorf("region_renderer.go contains mode-name literal %s — renderer must be mode-identity-free (Requirement 3.2)", literal)
		}
	}
}

// stripImportBlock removes the import(...) block from Go source to avoid
// false positives when checking for mode-name string literals.
func stripImportBlock(src string) string {
	const importStart = "import ("
	start := strings.Index(src, importStart)
	if start == -1 {
		return src
	}
	end := strings.Index(src[start:], ")")
	if end == -1 {
		return src
	}
	return src[:start] + src[start+end+1:]
}

// --- Integration Test: Region renderer uses LayoutBridge for content positioning ---
// After layout-padding-refactor: Region no longer has padding. The renderer uses
// LayoutBridge (via testBridge) with PaddingPct=0 for positioning.

func TestIntegration_NoPaddingOnRegion_240x320(t *testing.T) {
	// Construct a Region with bounds (0,0)-(240,320).
	// After refactor: no paddingPct, content starts at bridge.ContentOrigin() X=0 with width=240.
	bounds := image.Rect(0, 0, 240, 320)
	surf := surface.NewWithLog(bounds)

	// --- Render items ---
	rr := NewRegionRenderer(false, nil, nil)

	// renderItems now takes the whole ViewData, so that it can honour the layout
	// the style solved (per-row heights, spacing, offsets, row budget) instead of
	// re-deriving geometry. Cursor -1 keeps the original intent: no row highlighted.
	items := []string{"Hello", "World"}
	rr.renderItems(surf, style.ViewData{Items: items, Cursor: -1}, testBridge(bounds))

	// --- Verify draw log ---
	log := surf.DrawLog()
	if len(log) == 0 {
		t.Fatal("draw log is empty after renderItems")
	}

	var textCalls []surface.DrawCall
	var rectCalls []surface.DrawCall
	for _, call := range log {
		switch call.Type {
		case "text":
			textCalls = append(textCalls, call)
		case "rect":
			rectCalls = append(rectCalls, call)
		}
	}

	// Verify text calls have X = 0 (LayoutBridge with PaddingPct=0 → ContentOrigin X=0).
	if len(textCalls) == 0 {
		t.Fatal("no text draw calls recorded")
	}
	for i, call := range textCalls {
		if call.X != 0 {
			t.Fatalf("text call %d: X = %d, want 0", i, call.X)
		}
	}

	// Verify rect calls (row backgrounds) start at X=0 with width=240 (full bounds).
	if len(rectCalls) == 0 {
		t.Fatal("no rect draw calls recorded")
	}
	for i, call := range rectCalls {
		if call.Rect.Min.X != 0 {
			t.Fatalf("rect call %d: Min.X = %d, want 0", i, call.Rect.Min.X)
		}
		gotWidth := call.Rect.Dx()
		if gotWidth != 240 {
			t.Fatalf("rect call %d: width = %d, want 240", i, gotWidth)
		}
	}
}

// --- Tests: resolveTextFitFonts ---

func TestResolveTextFitFonts_NoOpWhenTiersSet(t *testing.T) {
	cat := textlayout.DefaultTextHints(image.Rect(0, 0, 128, 64)).Catalog
	v := style.ViewData{
		Items: []string{"hello"},
		Tiers: []tiercatalog.Tier{tiercatalog.TierNormal},
	}
	resolveTextFitFonts(&v, cat, 128)
	if v.FontIDs != nil {
		t.Fatal("resolveTextFitFonts should be no-op when Tiers already set")
	}
}

func TestResolveTextFitFonts_NoOpWhenFontIDsSet(t *testing.T) {
	cat := textlayout.DefaultTextHints(image.Rect(0, 0, 128, 64)).Catalog
	v := style.ViewData{
		Items:   []string{"hello"},
		FontIDs: []string{"spleen-5x8"},
	}
	resolveTextFitFonts(&v, cat, 128)
	if len(v.FontIDs) != 1 || v.FontIDs[0] != "spleen-5x8" {
		t.Fatal("resolveTextFitFonts should not overwrite existing FontIDs")
	}
}

func TestResolveTextFitFonts_SelectsFittingFont(t *testing.T) {
	// 128x64 panel: "[pager] no source" is 18 chars.
	// The function must pick a font whose advance * 18 <= 128.
	cat := textlayout.DefaultTextHints(image.Rect(0, 0, 128, 64)).Catalog
	if cat.PixelWidth() == 0 {
		t.Skip("no catalog available in test environment")
	}

	msg := "[pager] no source"
	v := style.ViewData{Items: []string{msg}, Static: true}
	resolveTextFitFonts(&v, cat, 128)

	if len(v.FontIDs) != 1 || v.FontIDs[0] == "" {
		t.Fatal("resolveTextFitFonts did not populate FontIDs")
	}

	runeCount := len([]rune(msg))
	for _, tier := range cat.Tiers() {
		e, ok := cat.Get(tier)
		if !ok || e.FontID != v.FontIDs[0] {
			continue
		}
		if e.GlyphAdvance*runeCount > 128 {
			t.Fatalf("selected font %q (advance %d) does not fit %d chars in 128px (%d needed)",
				v.FontIDs[0], e.GlyphAdvance, runeCount, e.GlyphAdvance*runeCount)
		}
		return
	}
}
