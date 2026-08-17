package pngpanel

import (
	"os"
	"testing"

	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// Property 2: Preservation - Non-Grayscale TextHints Unchanged
// For all ColorModeFullColor panels with valid dimensions:
// Capability == CapColorFast, all scroll flags true, ticker direction vertical.

func TestProperty_Preservation_FullColorCapability(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 4096).Draw(t, "width")
		height := rapid.IntRange(1, 4096).Draw(t, "height")

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(ColorModeFullColor),
			WithOutputDir(os.TempDir()),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		hints := panel.TextHints()

		// Capability must be CapColorFast (5) for full-color panels.
		if hints.Capability != textlayout.CapColorFast {
			t.Fatalf("FullColor Capability = %d, want %d (CapColorFast)", hints.Capability, textlayout.CapColorFast)
		}

		// All scroll flags must be true.
		if !hints.SupportsVerticalScroll {
			t.Fatal("FullColor: SupportsVerticalScroll should be true")
		}
		if !hints.SupportsHorizontalScroll {
			t.Fatal("FullColor: SupportsHorizontalScroll should be true")
		}
		if !hints.SupportsAutoScroll {
			t.Fatal("FullColor: SupportsAutoScroll should be true")
		}

		// PreferEventRefresh must be false.
		if hints.PreferEventRefresh {
			t.Fatal("FullColor: PreferEventRefresh should be false")
		}

		// DefaultTickerDirection must be vertical.
		if hints.DefaultTickerDirection != textlayout.TickerDirectionVertical {
			t.Fatalf("FullColor: DefaultTickerDirection = %q, want %q", hints.DefaultTickerDirection, textlayout.TickerDirectionVertical)
		}
	})
}

// Property 2: Preservation - Non-Grayscale TextHints Unchanged
// For all ColorModeMonochrome panels with valid dimensions:
// Capability == CapMonoSlow, PreferEventRefresh == true, ticker direction none.

func TestProperty_Preservation_MonochromeCapability(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 4096).Draw(t, "width")
		height := rapid.IntRange(1, 4096).Draw(t, "height")

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(ColorModeMonochrome),
			WithOutputDir(os.TempDir()),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		hints := panel.TextHints()

		// Capability must be CapMonoSlow (0) for monochrome panels.
		if hints.Capability != textlayout.CapMonoSlow {
			t.Fatalf("Monochrome Capability = %d, want %d (CapMonoSlow)", hints.Capability, textlayout.CapMonoSlow)
		}

		// PreferEventRefresh must be true.
		if !hints.PreferEventRefresh {
			t.Fatal("Monochrome: PreferEventRefresh should be true")
		}

		// DefaultTickerDirection must be none.
		if hints.DefaultTickerDirection != textlayout.TickerDirectionNone {
			t.Fatalf("Monochrome: DefaultTickerDirection = %q, want %q", hints.DefaultTickerDirection, textlayout.TickerDirectionNone)
		}

		// Scroll flags must all be false for monochrome.
		if hints.SupportsVerticalScroll {
			t.Fatal("Monochrome: SupportsVerticalScroll should be false")
		}
		if hints.SupportsHorizontalScroll {
			t.Fatal("Monochrome: SupportsHorizontalScroll should be false")
		}
		if hints.SupportsAutoScroll {
			t.Fatal("Monochrome: SupportsAutoScroll should be false")
		}
	})
}

// Property 2: Preservation - Non-Grayscale TextHints Unchanged
// For all modes: glyph metric fields remain unset at the panel layer and
// DefaultLineMode is always LineModeTruncate.

func TestProperty_Preservation_GlyphMetrics(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 4096).Draw(t, "width")
		height := rapid.IntRange(1, 4096).Draw(t, "height")
		mode := rapid.SampledFrom([]ColorMode{ColorModeFullColor, ColorModeMonochrome}).Draw(t, "colorMode")

		panel, err := New(
			WithDimensions(width, height),
			WithColorMode(mode),
			WithOutputDir(os.TempDir()),
		)
		if err != nil {
			t.Fatalf("unexpected construction error: %v", err)
		}

		hints := panel.TextHints()

		// Glyph metrics must be absent; Region selects the baseline font later.
		if hints.GlyphWidth != 0 {
			t.Fatalf("GlyphWidth = %d, want 0", hints.GlyphWidth)
		}
		if hints.GlyphHeight != 0 {
			t.Fatalf("GlyphHeight = %d, want 0", hints.GlyphHeight)
		}
		if hints.GlyphAdvance != 0 {
			t.Fatalf("GlyphAdvance = %d, want 0", hints.GlyphAdvance)
		}
		if hints.RowHeight != 0 {
			t.Fatalf("RowHeight = %d, want 0", hints.RowHeight)
		}

		// DefaultLineMode must always be truncate.
		if hints.DefaultLineMode != textlayout.LineModeTruncate {
			t.Fatalf("DefaultLineMode = %q, want %q", hints.DefaultLineMode, textlayout.LineModeTruncate)
		}

		// PixelWidth and PixelHeight must match configured dimensions.
		if hints.PixelWidth != width {
			t.Fatalf("PixelWidth = %d, want %d", hints.PixelWidth, width)
		}
		if hints.PixelHeight != height {
			t.Fatalf("PixelHeight = %d, want %d", hints.PixelHeight, height)
		}
	})
}
