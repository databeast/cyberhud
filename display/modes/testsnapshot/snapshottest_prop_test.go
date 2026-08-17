package testsnapshot

import (
	"os"
	"testing"

	"github.com/databeast/cyberhud/hardware/panels/pngpanel"
	"pgregory.net/rapid"

	// Blank imports to trigger mode init() self-registration.
	_ "github.com/databeast/cyberhud/display/modes/attract_matrix"
	_ "github.com/databeast/cyberhud/display/modes/clock"
	_ "github.com/databeast/cyberhud/display/modes/cycle"
	_ "github.com/databeast/cyberhud/display/modes/dashboard"
	_ "github.com/databeast/cyberhud/display/modes/gpio"
	_ "github.com/databeast/cyberhud/display/modes/gpio_control"
	_ "github.com/databeast/cyberhud/display/modes/image"
	_ "github.com/databeast/cyberhud/display/modes/menu"
	_ "github.com/databeast/cyberhud/display/modes/serial"
	_ "github.com/databeast/cyberhud/display/modes/stemma"
	_ "github.com/databeast/cyberhud/display/modes/system"
	_ "github.com/databeast/cyberhud/display/modes/systemd"
	_ "github.com/databeast/cyberhud/display/modes/testfonts"
	_ "github.com/databeast/cyberhud/display/modes/testpattern"
	_ "github.com/databeast/cyberhud/display/modes/thermal"
	_ "github.com/databeast/cyberhud/display/modes/ticker"
	_ "github.com/databeast/cyberhud/display/modes/usb"
	_ "github.com/databeast/cyberhud/display/modes/wifi"
	_ "github.com/databeast/cyberhud/display/modes/zmq"
)

// dimPair holds a width/height pair from a registered style.
type dimPair struct {
	w, h int
}

// registeredStyleDimensions returns the set of panel dimensions used by actual
// registered styles across all display modes. Extracted from the style registries
// (clock, dashboard, gpio, etc.). This is Option 3 from the design doc - tests
// the exact dimensions the framework encounters in production.
func registeredStyleDimensions() []dimPair {
	return []dimPair{
		// Monochrome OLED
		{128, 32}, {128, 64}, {128, 128}, {64, 128},
		// Color TFT (landscape and portrait)
		{160, 80}, {160, 128}, {240, 135}, {240, 240}, {320, 240},
		{480, 320}, {800, 480}, {80, 160}, {128, 160},
		{135, 240}, {240, 320}, {320, 480}, {480, 800},
		// E-ink
		{122, 250}, {176, 264}, {200, 200}, {212, 104}, {296, 128},
		{400, 300}, {104, 212}, {250, 122},
		{128, 296}, {264, 176}, {300, 400},
	}
}

// knownModes returns the list of registered display mode IDs that work
// without external dependencies (no GPIO manager, no STEMMA scanner, etc.).
// Menu is excluded because passive mode normalizes it to dashboard.
func knownModes() []string {
	return []string{
		"clock",
		"dashboard",
		"system",
		"systemd",
		"ticker",
		"testpattern",
		"testfonts",
		"cycle",
		"attract_matrix",
		"wifi",
		"zmq",
	}
}

// knownColorModes returns the supported PNGPanel color modes for property testing.
func knownColorModes() []pngpanel.ColorMode {
	return []pngpanel.ColorMode{
		pngpanel.ColorModeFullColor,
		pngpanel.ColorModeGrayscale,
		pngpanel.ColorModeMonochrome,
	}
}

// For any valid panel dimensions where tiercatalog.Build succeeds (drawn from
// registered style dimensions), any registered display mode, and any supported
// color mode, calling RenderSnapshot SHALL produce a decodable PNG file whose
// image dimensions equal the specified width and height.

func TestProperty_RenderRoundTripProducesValidPNG(t *testing.T) {
	dims := registeredStyleDimensions()
	modes := knownModes()
	colorModes := knownColorModes()

	rapid.Check(t, func(rt *rapid.T) {
		// Draw from registered style dimensions.
		d := rapid.SampledFrom(dims).Draw(rt, "dimensions")

		// Draw a safe mode.
		mode := rapid.SampledFrom(modes).Draw(rt, "mode")

		// Draw a color mode.
		cm := rapid.SampledFrom(colorModes).Draw(rt, "colorMode")

		// Create a temp directory for this iteration's output.
		outputDir, err := os.MkdirTemp("", "prop1-render-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(outputDir)

		// Render through the full production pipeline.
		pngPath := RenderSnapshot(t,
			WithDimensions(d.w, d.h),
			WithMode(mode),
			WithOutputDir(outputDir),
			WithColorMode(cm),
		)

		// Verify the PNG is decodable and has correct dimensions.
		VerifyPNG(t, pngPath)
		VerifyDimensions(t, pngPath, d.w, d.h)
	})
}

// For any PNG file produced by RenderSnapshot, comparing that file against itself
// via VerifyGolden SHALL pass (no pixel differences detected).

func TestProperty_GoldenSelfComparisonIsIdentity(t *testing.T) {
	modes := knownModes()
	colorModes := knownColorModes()

	rapid.Check(t, func(rt *rapid.T) {
		// Pick a random mode from the known-good set.
		modeID := rapid.SampledFrom(modes).Draw(rt, "mode")

		// Pick a random color mode.
		colorMode := rapid.SampledFrom(colorModes).Draw(rt, "colorMode")

		// Use known-good dimensions that work with tiercatalog.
		// 240x320 is a standard panel size supported by all modes.
		width := 240
		height := 320

		// Create a temp directory for this iteration's output.
		outputDir, err := os.MkdirTemp("", "golden-self-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(outputDir)

		// Render the snapshot through the full production pipeline.
		pngPath := RenderSnapshot(t,
			WithDimensions(width, height),
			WithMode(modeID),
			WithColorMode(colorMode),
			WithOutputDir(outputDir),
		)

		// THE PROPERTY: comparing the rendered PNG against itself must pass.
		// VerifyGolden does pixel-by-pixel comparison; self-comparison must
		// produce zero differences.
		VerifyGolden(t, pngPath, pngPath)
	})
}

// For any valid configuration (where tiercatalog succeeds) and any frame count N >= 1,
// rendering N consecutive frames SHALL produce a decodable PNG with dimensions matching
// the configured panel dimensions.
func TestProperty_FrameCountDoesNotAffectOutputValidity(t *testing.T) {
	knownGoodDims := []dimPair{
		{240, 320},
		{320, 240},
		{128, 160},
		{160, 128},
	}

	// Modes known to work without extra dependencies.
	knownGoodModes := []string{"clock", "dashboard"}

	rapid.Check(t, func(rt *rapid.T) {
		// Generate random frame count in [1, 10].
		frameCount := rapid.IntRange(1, 10).Draw(rt, "frameCount")

		// Pick a random valid dimension pair.
		dims := rapid.SampledFrom(knownGoodDims).Draw(rt, "dims")

		// Pick a random mode from known-good set.
		mode := rapid.SampledFrom(knownGoodModes).Draw(rt, "mode")

		// Create a temp directory for this iteration's output.
		outputDir, err := os.MkdirTemp("", "prop5-frames-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(outputDir)

		// Render with the generated frame count.
		pngPath := RenderSnapshot(t,
			WithDimensions(dims.w, dims.h),
			WithMode(mode),
			WithOutputDir(outputDir),
			WithFrameCount(frameCount),
			WithDisplayCategory(CategoryColor),
		)

		// Verify the output is a decodable PNG with correct dimensions.
		VerifyPNG(t, pngPath)
		VerifyDimensions(t, pngPath, dims.w, dims.h)
	})
}

// For any valid panel dimensions and any DisplayCategory value, the framework SHALL
// configure the PNGPanel's color encoding and RegionRenderer's monochrome flag according
// to the category mapping table.
func TestProperty_DisplayCategoryDerivesColorMode(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate random valid dimensions (within the valid range for the framework).
		width := rapid.IntRange(60, 4096).Draw(rt, "width")
		height := rapid.IntRange(60, 4096).Draw(rt, "height")

		// Generate a random DisplayCategory value (0-3).
		categoryInt := rapid.IntRange(0, 3).Draw(rt, "category")
		category := DisplayCategory(categoryInt)

		// Construct a snapshotConfig with the display category set.
		cfg := &snapshotConfig{
			width:           width,
			height:          height,
			displayCategory: &category,
		}

		// Call resolveColorMode to derive colorMode and monochrome from the category.
		resolveColorMode(cfg)

		// Verify against the expected mapping table.
		var wantColorMode pngpanel.ColorMode
		var wantMonochrome bool

		switch category {
		case CategoryColor:
			wantColorMode = pngpanel.ColorModeFullColor
			wantMonochrome = false
		case CategoryMono:
			wantColorMode = pngpanel.ColorModeMonochrome
			wantMonochrome = true
		case CategoryEink:
			wantColorMode = pngpanel.ColorModeMonochrome
			wantMonochrome = true
		case CategoryGrayscale:
			wantColorMode = pngpanel.ColorModeGrayscale
			wantMonochrome = false
		default:
			// Out-of-range categories default to FullColor with monochrome=false.
			wantColorMode = pngpanel.ColorModeFullColor
			wantMonochrome = false
		}

		if cfg.colorMode != wantColorMode {
			t.Fatalf("category %d at %dx%d: colorMode = %d, want %d",
				categoryInt, width, height, cfg.colorMode, wantColorMode)
		}
		if cfg.monochrome != wantMonochrome {
			t.Fatalf("category %d at %dx%d: monochrome = %v, want %v",
				categoryInt, width, height, cfg.monochrome, wantMonochrome)
		}
	})
}

// For any DisplayCategory integer value outside the range [0, 3], calling
// resolveColorMode SHALL set colorMode = ColorModeFullColor and monochrome = false.
func TestProperty_OutOfRangeDisplayCategoryDefaultsToFullColor(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a DisplayCategory value outside [0, 3].
		// Use a custom generator that picks either negative values or values >= 4.
		categoryInt := rapid.OneOf(
			rapid.IntRange(-1000, -1),
			rapid.IntRange(4, 1000),
		).Draw(rt, "outOfRangeCategory")
		category := DisplayCategory(categoryInt)

		// Construct a snapshotConfig with the out-of-range display category.
		cfg := &snapshotConfig{
			width:           240,
			height:          320,
			displayCategory: &category,
		}

		// Call resolveColorMode — should fall through to default case.
		resolveColorMode(cfg)

		// THE PROPERTY: out-of-range categories must default to FullColor + monochrome=false.
		if cfg.colorMode != pngpanel.ColorModeFullColor {
			t.Fatalf("category %d: colorMode = %d, want ColorModeFullColor (%d)",
				categoryInt, cfg.colorMode, pngpanel.ColorModeFullColor)
		}
		if cfg.monochrome != false {
			t.Fatalf("category %d: monochrome = true, want false", categoryInt)
		}
	})
}

// For any valid panel dimensions (where tiercatalog succeeds) and any padding
// percentage in [0, 50], the output PNG dimensions SHALL equal the configured
// panel dimensions (padding affects internal layout only, not the output image size).
func TestProperty_PaddingDoesNotAffectOutputPNGDimensions(t *testing.T) {
	// Use large dimensions so that even high padding values leave enough content
	// area for tiercatalog. At padding P%, the content width = W - 2*(P*W/100).
	// For tiercatalog to succeed, the content area needs at least ~60px wide
	// (smallest font advance=6 × MinChars=10). We use dimensions >= 240 so
	// that padding up to ~37% still leaves adequate content area, and we cap
	// the generated padding per-iteration based on the selected dimensions to
	// ensure tiercatalog success.
	knownGoodDims := []dimPair{
		{240, 320},
		{320, 240},
		{240, 240},
		{480, 320},
		{320, 480},
		{800, 480},
		{480, 800},
	}

	modes := []string{"clock", "dashboard", "system", "ticker"}

	rapid.Check(t, func(rt *rapid.T) {
		// Pick a random valid dimension pair.
		dims := rapid.SampledFrom(knownGoodDims).Draw(rt, "dims")

		// Compute the maximum padding that still leaves at least 70px in both
		// dimensions for tiercatalog (with some safety margin).
		// Content width = W - 2*(P*W/100), solve for P: P < (W - 70)*100/(2*W)
		minDim := dims.w
		if dims.h < minDim {
			minDim = dims.h
		}
		maxPad := (minDim - 70) * 100 / (2 * minDim)
		if maxPad > 50 {
			maxPad = 50
		}
		if maxPad < 0 {
			maxPad = 0
		}

		// Generate random padding in [0, maxPad].
		padding := rapid.IntRange(0, maxPad).Draw(rt, "padding")

		// Pick a random mode from known-good set.
		mode := rapid.SampledFrom(modes).Draw(rt, "mode")

		// Create a temp directory for this iteration's output.
		outputDir := t.TempDir()

		// Render with the generated padding value.
		pngPath := RenderSnapshot(t,
			WithDimensions(dims.w, dims.h),
			WithMode(mode),
			WithOutputDir(outputDir),
			WithPadding(padding),
			WithDisplayCategory(CategoryColor),
		)

		// THE PROPERTY: output PNG dimensions must equal configured panel dimensions
		// regardless of padding value. Padding affects internal layout only.
		VerifyPNG(t, pngPath)
		VerifyDimensions(t, pngPath, dims.w, dims.h)
	})
}
