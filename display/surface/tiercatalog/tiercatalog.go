// Package tiercatalog computes which font sizes correspond to each tier for a
// given region's pixel dimensions. This is the region's declaration of "what
// fits here."
package tiercatalog

import (
	"strings"

	"github.com/databeast/cyberhud/display/surface/fonts"
)

// Tier is a semantic size name.
type Tier string

const (
	TierSmall    Tier = "small"
	TierNormal   Tier = "normal"
	TierLarge    Tier = "large"
	TierHuge     Tier = "huge"
	TierColossal Tier = "colossal"
	TierFull     Tier = "full"
	TierFullsize Tier = "fullsize" // Backward-compat alias; resolves to TierFull in Catalog.Get
)

// DefaultTargetsMM defines the default physical height targets in millimeters.
// TierFull has no mm target — it always selects the largest qualifying font.
var DefaultTargetsMM = map[Tier]float64{
	TierSmall:    3.0,
	TierNormal:   5.0,
	TierLarge:    8.0,
	TierHuge:     12.0,
	TierColossal: 18.0,
}

// DefaultTargetsPx defines fallback pixel height targets when PPI is unknown.
// TierFull has no pixel target — it always selects the largest qualifying font.
var DefaultTargetsPx = map[Tier]int{
	TierSmall:    8,
	TierNormal:   14,
	TierLarge:    20,
	TierHuge:     28,
	TierColossal: 40,
}

// tierOrder defines the canonical ascending order of tiers for monotonicity enforcement.
var tierOrder = []Tier{TierSmall, TierNormal, TierLarge, TierHuge, TierColossal, TierFull}

// Entry holds the target font metrics for a given tier.
type Entry struct {
	GlyphWidth   int
	GlyphHeight  int
	GlyphAdvance int
	RowHeight    int
	FontID       string // Identifies the selected font face; empty when unset.
}

// defaultMinChars is the MinChars applied when a caller passes zero or a
// negative value. Ten characters is the legibility floor the display system
// assumes: a region that cannot show ten monospace glyphs across is treated as
// too narrow for text at that tier.
const defaultMinChars = 10

// Catalog maps tiers to their target font metrics for a specific region.
type Catalog struct {
	entries  map[Tier]Entry
	width    int // region pixel width
	height   int // region pixel height
	minChars int // minimum characters constraint actually achieved during build

	// requestedMinChars records what the caller asked for, which differs from
	// minChars only when relaxation kicked in. Kept for diagnostics so a
	// degraded region can be reported honestly rather than silently.
	requestedMinChars int
	relaxed           bool
}

// Get returns the entry for the given tier. Returns zero Entry and false if
// the tier has no assignment.
//
// Prefer [Catalog.Entry] in styles and renderers: it cannot fail, and its
// fallbacks resolve to real registered fonts instead of the fabricated metrics
// that callers used to substitute by hand. Get remains for callers to whom the
// absence of a tier is itself meaningful.
func (c Catalog) Get(tier Tier) (Entry, bool) {
	e, ok := c.entries[NormalizeTier(tier)]
	return e, ok
}

// Tiers returns all defined tiers in ascending size order.
func (c Catalog) Tiers() []Tier {
	return tierOrder
}

// PixelWidth returns the region width used to build this catalog.
func (c Catalog) PixelWidth() int {
	return c.width
}

// PixelHeight returns the region height used to build this catalog.
func (c Catalog) PixelHeight() int {
	return c.height
}

// MinChars returns the minimum characters constraint used to build this catalog.
func (c Catalog) MinChars() int {
	return c.minChars
}

// Params controls catalog construction.
type Params struct {
	PixelWidth  int     // Region pixel width
	PixelHeight int     // Region pixel height
	MinChars    int     // Minimum characters that must fit horizontally (default: 10)
	PPI         float64 // Optional: panel PPI for glyph height scaling. Zero = no scaling.

	// AllowRelaxedMinChars opts into best-effort construction: when no registered
	// font is narrow enough to fit MinChars characters across PixelWidth, Build
	// lowers MinChars to whatever the narrowest registered font does achieve
	// instead of returning an error.
	//
	// Why this is opt-in rather than the default: the strict behaviour — Build
	// fails when the legibility floor cannot be met — is a deliberate contract
	// with dedicated tests in this package, and callers performing capability
	// checks legitimately want to be told "text does not fit here."
	//
	// Why the display pipeline opts in: the alternative is worse. When Build fails,
	// region.buildCatalogIfNeeded logs and leaves hints.Catalog zero, and every
	// style then substitutes textlayout's 6x10 default metrics, which belong to no
	// registered font. The style lays out for a font the renderer will never draw
	// with, and text lands in the wrong place at the wrong size. Rendering real
	// small text is strictly better than rendering imaginary text.
	//
	// When relaxation occurs, [Catalog.MinChars] reports the achieved value and
	// [Catalog.Relaxed] returns true, so the width-safety invariant
	// GlyphAdvance*MinChars <= PixelWidth still holds against the catalog's own
	// reported constraint.
	AllowRelaxedMinChars bool

	// TierTargetsMM provides custom mm targets per tier (overrides DefaultTargetsMM).
	// A nil map means "use defaults."
	TierTargetsMM map[Tier]float64

	// TierTargetsPx provides custom pixel targets per tier when PPI is zero (overrides DefaultTargetsPx).
	// A nil map means "use defaults."
	TierTargetsPx map[Tier]int
}

// Build constructs a TierCatalog for the given region dimensions using the
// physical-size-aware best-fit algorithm. It filters the font registry by
// advance budget, resolves pixel targets per tier, selects the closest font
// for each tier, enforces monotonicity, and returns the catalog.
// Returns an error if no fonts are registered or none fit the region.
func Build(p Params) (Catalog, error) {
	// 1. Normalize negative dimensions.
	if p.PixelWidth < 0 {
		p.PixelWidth = 0
	}
	if p.PixelHeight < 0 {
		p.PixelHeight = 0
	}
	if p.PPI < 0 {
		p.PPI = 0
	}

	// 2. Resolve the effective MinChars before filtering.
	//
	// Relaxation is applied here rather than inside buildCandidatePool so that
	// helper keeps its single responsibility (filter by a given budget, report why
	// nothing qualified) and its existing signature, which package tests call
	// directly.
	requestedMinChars := p.MinChars
	if requestedMinChars <= 0 {
		requestedMinChars = defaultMinChars
	}
	effectiveMinChars := requestedMinChars
	if p.AllowRelaxedMinChars {
		effectiveMinChars = relaxMinChars(p.PixelWidth, requestedMinChars)
	}
	p.MinChars = effectiveMinChars

	// 3. Build candidate pool (handles MinChars normalization, error reporting).
	candidates, _, err := buildCandidatePool(p)
	if err != nil {
		return Catalog{}, err
	}

	// 3. Resolve pixel targets.
	pixelTargets := resolvePixelTargets(p)

	// 4. Select best-fit font for each tier (except TierFull).
	entries := make(map[Tier]Entry, len(tierOrder))
	for _, tier := range tierOrder {
		if tier == TierFull {
			continue
		}
		selected := bestFit(candidates, pixelTargets[tier])
		entries[tier] = candidateToEntry(selected, pixelTargets[tier])
	}

	// 5. TierFull: select the largest qualifying font.
	largest := selectLargest(candidates)
	entries[TierFull] = candidateToEntry(largest, 0)

	// 6. Enforce monotonicity.
	enforceMonotonicity(entries)

	return Catalog{
		entries:           entries,
		width:             p.PixelWidth,
		height:            p.PixelHeight,
		minChars:          effectiveMinChars,
		requestedMinChars: requestedMinChars,
		relaxed:           effectiveMinChars != requestedMinChars,
	}, nil
}

// relaxMinChars returns requested when at least one registered font already fits
// that many characters across pixelWidth, and otherwise the largest character
// count the narrowest registered font can achieve.
//
// Returning a value below requested is what lets a 32px-wide region produce a
// real font instead of failing the build. The floor of 1 matters: a region narrower
// than the narrowest glyph still gets a catalog whose entries carry that glyph's
// true metrics, which is the honest answer — one clipped character is a visible,
// diagnosable symptom, whereas fabricated metrics silently corrupt every style's
// layout arithmetic.
func relaxMinChars(pixelWidth, requested int) int {
	if pixelWidth <= 0 {
		return requested
	}

	allFonts := font.List()
	if len(allFonts) == 0 {
		// Let buildCandidatePool report the empty-registry error unchanged.
		return requested
	}

	smallestAdvance := allFonts[0].Metrics().GlyphAdvance
	for _, f := range allFonts[1:] {
		if adv := f.Metrics().GlyphAdvance; adv > 0 && adv < smallestAdvance {
			smallestAdvance = adv
		}
	}
	if smallestAdvance <= 0 {
		return requested
	}

	// Already satisfiable at the requested constraint.
	if pixelWidth/requested >= smallestAdvance {
		return requested
	}

	achievable := pixelWidth / smallestAdvance
	if achievable < 1 {
		achievable = 1
	}
	return achievable
}

// candidateToEntry converts a Candidate's metrics into a tier Entry.
// pixelHeight is passed to MetricsAt so scalable fonts produce metrics at the
// exact requested size. Bitmap fonts ignore the argument.
func candidateToEntry(c Candidate, pixelHeight int) Entry {
	m := c.MetricsAt(pixelHeight)
	return Entry{
		GlyphWidth:   m.GlyphWidth,
		GlyphHeight:  m.GlyphHeight,
		GlyphAdvance: m.GlyphAdvance,
		RowHeight:    m.RowHeight,
		FontID:       c.ID(),
	}
}

// familyPriority returns the tie-breaking priority for a font face based on its ID.
// Spleen has highest priority (3), then Terminus (2), then Cozette (1), others (0).
func familyPriority(id string) int {
	lower := strings.ToLower(id)
	switch {
	case strings.HasPrefix(lower, "spleen-"):
		return 3
	case strings.HasPrefix(lower, "terminus-"):
		return 2
	case strings.HasPrefix(lower, "cozette-"):
		return 1
	default:
		return 0
	}
}
