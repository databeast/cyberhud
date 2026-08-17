package displaymode

import (
	"regexp"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// genTextHints generates random TextHints values for property testing.
func genTextHints(t *rapid.T) textlayout.TextHints {
	return textlayout.TextHints{
		PixelWidth:         rapid.IntRange(0, 2048).Draw(t, "pixelWidth"),
		PixelHeight:        rapid.IntRange(0, 2048).Draw(t, "pixelHeight"),
		GlyphWidth:         rapid.IntRange(1, 20).Draw(t, "glyphWidth"),
		GlyphHeight:        rapid.IntRange(1, 20).Draw(t, "glyphHeight"),
		GlyphAdvance:       rapid.IntRange(1, 24).Draw(t, "glyphAdvance"),
		RowHeight:          rapid.IntRange(1, 32).Draw(t, "rowHeight"),
		PreferEventRefresh: rapid.Bool().Draw(t, "preferEventRefresh"),
	}
}

// =============================================================================
//
// For any style in the template registry, Name() returns a lowercase ASCII
// string between 1 and 32 characters matching the pattern ^[a-z0-9-]+$ with
// the correct prefix (mono-, color-, eink-, or grayscale-fast-).

// =============================================================================

func TestPropertyStyleNameValidity(t *testing.T) {
	namePattern := regexp.MustCompile(`^[a-z0-9-]+$`)

	rapid.Check(t, func(t *rapid.T) {
		styles := templateRegistry.Enumerate()
		idx := rapid.IntRange(0, len(styles)-1).Draw(t, "styleIndex")
		s := styles[idx]

		name := s.Name()

		// Length: 1–32 characters.
		if len(name) < 1 || len(name) > 32 {
			t.Fatalf("style name %q has length %d; want 1–32", name, len(name))
		}

		// Pattern: lowercase ASCII with hyphens and digits only.
		if !namePattern.MatchString(name) {
			t.Fatalf("style name %q does not match pattern ^[a-z0-9-]+$", name)
		}

		// Prefix: must be one of the valid prefixes.
		validPrefixes := []string{"mono-", "color-", "eink-", "grayscale-fast-"}
		hasValidPrefix := false
		for _, prefix := range validPrefixes {
			if strings.HasPrefix(name, prefix) {
				hasValidPrefix = true
				break
			}
		}
		if !hasValidPrefix {
			t.Fatalf("style name %q does not start with a valid prefix (mono-, color-, eink-, grayscale-fast-)", name)
		}
	})
}

// =============================================================================
//
// For any style in the template registry, Requirements() returns a
// SurfaceRequirements where MinWidth/MinHeight match the style's documented
// target resolution, and Capability matches the style's declared capability needs.

// =============================================================================

func TestPropertySurfaceRequirementsCorrectness(t *testing.T) {
	type expectedReqs struct {
		width, height int
		capability    style.Capability
	}

	// Lookup table of expected values per style name.
	expected := map[string]expectedReqs{
		// Primary mono styles
		"mono-128x32": {128, 32, style.MonoFast},
		"mono-128x64": {128, 64, style.MonoFast},
		// Primary color styles
		"color-160x80":  {160, 80, style.ColorFast},
		"color-160x128": {160, 128, style.ColorFast},
		"color-240x135": {240, 135, style.ColorFast},
		"color-240x240": {240, 240, style.ColorFast},
		"color-320x240": {320, 240, style.ColorFast},
		"color-480x320": {480, 320, style.ColorFast},
		"color-800x480": {800, 480, style.ColorFast},
		// Primary e-ink styles
		"eink-122x250": {122, 250, style.MonoSlow},
		"eink-176x264": {176, 264, style.MonoSlow},
		"eink-200x200": {200, 200, style.MonoSlow},
		"eink-212x104": {212, 104, style.MonoSlow},
		"eink-296x128": {296, 128, style.MonoSlow},
		"eink-400x300": {400, 300, style.MonoSlow},
		"eink-480x800": {480, 800, style.MonoSlow},
		"eink-800x480": {800, 480, style.MonoSlow},
		// Grayscale-fast styles (grayscale-fast capability)
		"grayscale-fast-160x80":  {160, 80, style.GrayscaleFast},
		"grayscale-fast-160x128": {160, 128, style.GrayscaleFast},
		"grayscale-fast-240x135": {240, 135, style.GrayscaleFast},
		"grayscale-fast-240x240": {240, 240, style.GrayscaleFast},
		"grayscale-fast-320x240": {320, 240, style.GrayscaleFast},
		"grayscale-fast-480x320": {480, 320, style.GrayscaleFast},
		"grayscale-fast-800x480": {800, 480, style.GrayscaleFast},
	}

	rapid.Check(t, func(t *rapid.T) {
		styles := templateRegistry.Enumerate()
		idx := rapid.IntRange(0, len(styles)-1).Draw(t, "styleIndex")
		s := styles[idx]

		name := s.Name()
		reqs := s.Requirements()

		exp, ok := expected[name]
		if !ok {
			t.Fatalf("style %q not found in expected requirements table", name)
		}

		if reqs.MinWidth != exp.width {
			t.Fatalf("style %q: MinWidth=%d, want %d", name, reqs.MinWidth, exp.width)
		}
		if reqs.MinHeight != exp.height {
			t.Fatalf("style %q: MinHeight=%d, want %d", name, reqs.MinHeight, exp.height)
		}
		if reqs.Capability != exp.capability {
			t.Fatalf("style %q: Capability=%v, want %v", name, reqs.Capability, exp.capability)
		}
	})
}

// =============================================================================
//
// For any style in the template registry and for any valid TextHints value,
// style.Supports(hints) returns the same Fitness value as
// style.EvaluateFitness(style.Requirements(), hints).

// =============================================================================

func TestPropertySupportsDelegatesToEvaluateFitness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		styles := templateRegistry.Enumerate()
		idx := rapid.IntRange(0, len(styles)-1).Draw(t, "styleIndex")
		s := styles[idx]

		hints := genTextHints(t)

		got := s.Supports(hints)
		want := style.EvaluateFitness(s.Requirements(), hints)

		if got != want {
			t.Fatalf("style %q: Supports(hints)=%v, EvaluateFitness(Requirements(), hints)=%v",
				s.Name(), got, want)
		}
	})
}

// =============================================================================
//
// For any style in the template registry and for any Snapshot and TextHints
// values, Build().Items contains at least one non-empty string.

// =============================================================================

func TestPropertyBuildProducesNonEmptyItems(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		styles := templateRegistry.Enumerate()
		idx := rapid.IntRange(0, len(styles)-1).Draw(t, "styleIndex")
		s := styles[idx]

		hints := genTextHints(t)
		ctx := style.NewStyleContext(hints)
		result := s.Build(Snapshot{}, Policy{}, ctx)

		if len(result.Items) == 0 {
			t.Fatalf("style %q: Build() returned empty Items slice", s.Name())
		}

		hasNonEmpty := false
		for _, item := range result.Items {
			if item != "" {
				hasNonEmpty = true
				break
			}
		}
		if !hasNonEmpty {
			t.Fatalf("style %q: Build() Items has no non-empty strings: %v", s.Name(), result.Items)
		}
	})
}

// =============================================================================
//
// For each color-capable resolution (160×80, 160×128, 240×135, 240×240,
// 320×240, 480×320, 800×480), the registry contains both a primary color style
// (Capability=ColorFast) and a grayscale-fast style (Capability=GrayscaleFast)
// with identical MinWidth and MinHeight values for that resolution.

// =============================================================================

func TestPropertyGrayscaleFastStylePairingInvariant(t *testing.T) {
	// Color-capable resolutions that must have both primary and grayscale-fast variants.
	type resolution struct {
		width, height int
	}
	colorResolutions := []resolution{
		{160, 80},
		{160, 128},
		{240, 135},
		{240, 240},
		{320, 240},
		{480, 320},
		{800, 480},
	}

	styles := templateRegistry.Enumerate()

	for _, res := range colorResolutions {
		var foundPrimary, foundGrayscaleFast bool

		for _, s := range styles {
			reqs := s.Requirements()
			if reqs.MinWidth != res.width || reqs.MinHeight != res.height {
				continue
			}

			name := s.Name()
			if strings.HasPrefix(name, "color-") && reqs.Capability == style.ColorFast {
				foundPrimary = true
			}
			if strings.HasPrefix(name, "grayscale-fast-") && reqs.Capability == style.GrayscaleFast {
				foundGrayscaleFast = true
			}
		}

		if !foundPrimary {
			t.Fatalf("resolution %dx%d: missing primary color style (Capability=ColorFast)",
				res.width, res.height)
		}
		if !foundGrayscaleFast {
			t.Fatalf("resolution %dx%d: missing grayscale-fast style (Capability=GrayscaleFast)",
				res.width, res.height)
		}
	}
}
