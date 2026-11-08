package tests

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
)

func fitCtx(t *testing.T, w, h int) style.StyleContext {
	t.Helper()
	cat, err := tiercatalog.Build(tiercatalog.Params{PixelWidth: w, PixelHeight: h})
	if err != nil {
		t.Fatalf("build catalog: %v", err)
	}
	return style.NewStyleContext(textlayout.TextHints{PixelWidth: w, PixelHeight: h, Catalog: cat})
}

// A short string fits at the ceiling tier itself — no sizing down needed.
func TestFitTier_ShortStringFitsAtCeiling(t *testing.T) {
	ctx := fitCtx(t, 480, 800)
	tier, ok := ctx.FitTier("Hi", 424, tiercatalog.TierColossal)
	if !ok {
		t.Fatal("short string should fit")
	}
	if tier != tiercatalog.TierColossal {
		t.Errorf("short string should fit at the ceiling, got %v", tier)
	}
}

// A long string that overflows the ceiling must come back at a smaller tier,
// still fitting the width — the "size down" behavior.
func TestFitTier_SizesDownToFit(t *testing.T) {
	ctx := fitCtx(t, 480, 800)
	const width = 424
	s := "next: Design review (in 1h 30m)"
	tier, ok := ctx.FitTier(s, width, tiercatalog.TierColossal)
	if !ok {
		t.Fatal("string should fit at some tier on a 480px panel")
	}
	// It must actually fit at the returned tier...
	if !ctx.Metrics(tier).Fits(s, width) {
		t.Errorf("returned tier %v does not fit %q in %dpx", tier, s, width)
	}
	// ...and must have sized below the ceiling (colossal cannot hold 31 glyphs here).
	if tier == tiercatalog.TierColossal {
		t.Error("expected a tier smaller than the colossal ceiling")
	}
}

// When even the smallest tier cannot render the string, ok is false — the
// signal to drop the row rather than clip it.
func TestFitTier_ReportsUnfittable(t *testing.T) {
	ctx := fitCtx(t, 128, 64)
	s := strings.Repeat("W", 200) // absurdly long for a 128px panel
	if _, ok := ctx.FitTier(s, 120, tiercatalog.TierNormal); ok {
		t.Error("a 200-char string must not fit a 128px panel at any tier")
	}
}

// The ceiling is honored: a tier larger than the ceiling is never returned even
// if it would also fit.
func TestFitTier_NeverExceedsCeiling(t *testing.T) {
	ctx := fitCtx(t, 480, 800)
	tier, ok := ctx.FitTier("Hi", 424, tiercatalog.TierNormal)
	if !ok {
		t.Fatal("short string should fit")
	}
	if tier != tiercatalog.TierNormal {
		t.Errorf("result %v exceeded the Normal ceiling", tier)
	}
}
