package style

import (
	"fmt"
	"math"

	"github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// FitnessReporter is the minimal interface required by FitnessNotes.
// Any Style[S,P] satisfies this interface, allowing FitnessNotes to work
// with styles of any snapshot type.
type FitnessReporter interface {
	Requirements() SurfaceRequirements
	Supports(hints textlayout.TextHints) Fitness
}

// FitnessNotes generates informational "note:" lines for unmet requirements.
// Returns nil when fitness is Full or Optimal. When Degraded or Unsupported,
// returns one note per unmet SurfaceRequirements field describing the mismatch
// between the style's needs and the panel's actual characteristics.
func FitnessNotes(s FitnessReporter, hints textlayout.TextHints) []string {
	fitness := s.Supports(hints)
	if fitness >= Full {
		return nil
	}

	reqs := s.Requirements()
	var notes []string

	// Check minimum dimensions.
	if reqs.MinWidth > 0 && hints.PixelWidth < reqs.MinWidth {
		notes = append(notes, fmt.Sprintf(
			"note: style requires minimum width %d but panel is %d",
			reqs.MinWidth, hints.PixelWidth))
	}
	if reqs.MinHeight > 0 && hints.PixelHeight < reqs.MinHeight {
		notes = append(notes, fmt.Sprintf(
			"note: style requires minimum height %d but panel is %d",
			reqs.MinHeight, hints.PixelHeight))
	}

	// Check preferred dimensions.
	if reqs.PreferredWidth > 0 && hints.PixelWidth < reqs.PreferredWidth {
		notes = append(notes, fmt.Sprintf(
			"note: style prefers %d\u00d7%d but panel is %d\u00d7%d",
			reqs.PreferredWidth, reqs.PreferredHeight, hints.PixelWidth, hints.PixelHeight))
	} else if reqs.PreferredHeight > 0 && hints.PixelHeight < reqs.PreferredHeight {
		notes = append(notes, fmt.Sprintf(
			"note: style prefers %d\u00d7%d but panel is %d\u00d7%d",
			reqs.PreferredWidth, reqs.PreferredHeight, hints.PixelWidth, hints.PixelHeight))
	}

	// Check capability requirement.
	if reqs.Capability > Capability(hints.Capability) {
		notes = append(notes, fmt.Sprintf("note: style requires %s but panel provides %s", reqs.Capability, Capability(hints.Capability)))
	}

	// Check MinRows text-level constraint.
	if reqs.MinRows > 0 {
		smallestRowHeight := 0
		canSatisfy := false
		for _, f := range font.List() {
			m := f.Metrics()
			if m.RowHeight <= 0 {
				continue
			}
			if smallestRowHeight == 0 || m.RowHeight < smallestRowHeight {
				smallestRowHeight = m.RowHeight
			}
			if reqs.MinRows*m.RowHeight <= hints.PixelHeight {
				canSatisfy = true
			}
		}
		if !canSatisfy && smallestRowHeight > 0 {
			maxRows := hints.PixelHeight / smallestRowHeight
			notes = append(notes, fmt.Sprintf(
				"note: style requires %d rows but panel (%dpx height) fits at most %d with smallest font (RowHeight=%d)",
				reqs.MinRows, hints.PixelHeight, maxRows, smallestRowHeight))
		}
	}

	// Check MinCharsPerLine text-level constraint.
	if reqs.MinCharsPerLine > 0 {
		smallestGlyphAdvance := 0
		canSatisfy := false
		for _, f := range font.List() {
			m := f.Metrics()
			if m.GlyphAdvance <= 0 {
				continue
			}
			if smallestGlyphAdvance == 0 || m.GlyphAdvance < smallestGlyphAdvance {
				smallestGlyphAdvance = m.GlyphAdvance
			}
			if reqs.MinCharsPerLine*m.GlyphAdvance <= hints.PixelWidth {
				canSatisfy = true
			}
		}
		if !canSatisfy && smallestGlyphAdvance > 0 {
			maxChars := hints.PixelWidth / smallestGlyphAdvance
			notes = append(notes, fmt.Sprintf(
				"note: style requires %d chars/line but panel (%dpx width) fits at most %d with narrowest font (GlyphAdvance=%d)",
				reqs.MinCharsPerLine, hints.PixelWidth, maxChars, smallestGlyphAdvance))
		}
	}

	return notes
}

// FitnessPostApply returns a PostApply hook for a mode's command handler.
// When the "style" key is among the applied keys, it looks up the currently
// selected style (via styleName) in the registry and generates fitness notes
// against the panel hints supplied by hints (typically modehints.Current).
//
// Both inputs are injected so this package gains no dependency on where
// hints or policy live; each mode wires its own sources:
//
//	var fitnessNotesPostApply = modeRegistry.FitnessPostApply(
//		modehints.Current, func() string { return GetPolicy().Style })
func (r *StyleRegistry[S, P]) FitnessPostApply(hints func() (textlayout.TextHints, bool), styleName func() string) func(appliedKeys []string) []string {
	return func(appliedKeys []string) []string {
		styleChanged := false
		for _, k := range appliedKeys {
			if k == "style" {
				styleChanged = true
				break
			}
		}
		if !styleChanged {
			return nil
		}

		h, ok := hints()
		if !ok {
			// No panel hints available (e.g., testing or headless mode).
			return nil
		}

		s := r.Lookup(styleName())
		if s == nil {
			return nil
		}

		return FitnessNotes(s, h)
	}
}

// Fitness indicates how well a style can render on a specific panel.
// Values encode both a tier (Unsupported/Degraded/Full/Optimal) and a
// tie-breaking score within each tier. Higher values are always better fits.
// Use the tier constants for threshold comparisons (e.g., fitness >= Full).
type Fitness int

const (
	Unsupported Fitness = 0    // Cannot render at all
	Degraded    Fitness = 1000 // Renders with visual loss
	Full        Fitness = 2000 // Meets all minimum requirements
	Optimal     Fitness = 3000 // Meets preferred dimensions and capability set
)

// EvaluateFitness implements the standard fitness evaluation algorithm.
// All display mode Supports() methods should delegate to this function.
//
// The algorithm follows the design specification:
//   - Zero dimensions → Unsupported
//   - Capability ordering violation (reqs > hints) → Unsupported
//   - Below MinWidth/MinHeight → Unsupported
//   - TextFitness delegation → Unsupported if text constraints fail
//   - Meets all min reqs and preferred → Optimal + bonus
//   - Meets all min reqs but not preferred → Full + bonus
//
// Within a tier, a bonus score (0–999) breaks ties by preferring styles whose
// capability is closest to the panel's actual capability (higher weight) and
// whose resolution is the tightest fit for the panel dimensions (lower weight).
func EvaluateFitness(reqs SurfaceRequirements, hints textlayout.TextHints) Fitness {
	// Zero-dimension panels cannot render anything.
	if hints.PixelWidth == 0 || hints.PixelHeight == 0 {
		return Unsupported
	}

	// Capability ordering check: a style requiring a higher capability level
	// than the panel provides cannot render on that panel.
	if reqs.Capability > Capability(hints.Capability) {
		return Unsupported
	}

	// Check minimum dimension requirements.
	if reqs.MinWidth > 0 && hints.PixelWidth < reqs.MinWidth {
		return Unsupported
	}
	if reqs.MinHeight > 0 && hints.PixelHeight < reqs.MinHeight {
		return Unsupported
	}

	// Text-level constraint check.
	if reqs.MinRows > 0 || reqs.MinCharsPerLine > 0 {
		if TextFitness(reqs, hints) == Unsupported {
			return Unsupported
		}
	}

	// PPI range check (only when panel PPI is known).
	if hints.PPI > 0 {
		if reqs.MinPPI > 0 && hints.PPI < reqs.MinPPI {
			return Unsupported
		}
		if reqs.MaxPPI > 0 && hints.PPI > reqs.MaxPPI {
			return Unsupported
		}
	}

	// At this point, all minimum requirements are met.
	// Determine the base tier.
	meetsPreferred := true
	if reqs.PreferredWidth > 0 && hints.PixelWidth < reqs.PreferredWidth {
		meetsPreferred = false
	}
	if reqs.PreferredHeight > 0 && hints.PixelHeight < reqs.PreferredHeight {
		meetsPreferred = false
	}

	baseTier := Full
	if meetsPreferred {
		baseTier = Optimal
	}

	// ── Tie-breaking bonus (0–999) ───────────────────────────────────────────
	// Three components:
	//   capBonus (0–500): prefer styles whose required capability is closest to
	//                     the panel's actual capability. An exact match scores 500;
	//                     each tier of distance reduces the score proportionally.
	//   resBonus (0–299): prefer styles whose min dimensions are closest to the
	//                     panel dimensions (less wasted area = tighter fit).
	//   ppiBonus (0–200): prefer styles whose PPI range center is closest to the
	//                     panel's actual PPI.

	// Capability proximity bonus.
	// capDistance is 0 when exact match, up to 5 when MonoSlow on a ColorFast panel.
	panelCap := int(hints.Capability)
	styleCap := int(reqs.Capability)
	capDistance := panelCap - styleCap // always >= 0 because we passed the gate above
	// Map distance 0→500, 1→400, 2→300, 3→200, 4→100, 5→0
	capBonus := 0
	if capDistance <= 5 {
		capBonus = (5 - capDistance) * 100
	}

	// Resolution proximity bonus.
	// Prefer styles whose MinWidth×MinHeight is closest to the panel area.
	// A style that exactly matches the panel resolution gets 299; a style that
	// is much smaller (wasted panel space) scores lower.
	resBonus := 0
	panelArea := hints.PixelWidth * hints.PixelHeight
	if panelArea > 0 {
		styleW := reqs.MinWidth
		if styleW <= 0 {
			styleW = hints.PixelWidth
		}
		styleH := reqs.MinHeight
		if styleH <= 0 {
			styleH = hints.PixelHeight
		}
		styleArea := styleW * styleH
		// ratio is styleArea/panelArea clamped to [0, 1].
		// Multiply by 299 to get the bonus.
		if styleArea >= panelArea {
			resBonus = 299
		} else {
			resBonus = (styleArea * 299) / panelArea
		}
	}

	// PPI proximity bonus (0-200): prefer styles whose PPI range center
	// is closest to the panel's actual PPI.
	ppiBonus := 0
	if hints.PPI > 0 && (reqs.MinPPI > 0 || reqs.MaxPPI > 0) {
		center := (reqs.MinPPI + reqs.MaxPPI) / 2
		if reqs.MinPPI == 0 {
			center = reqs.MaxPPI
		}
		if reqs.MaxPPI == 0 {
			center = reqs.MinPPI
		}
		distance := math.Abs(hints.PPI - center)
		// Normalize: 0 distance → 200 bonus, 500+ distance → 0 bonus.
		if distance < 500 {
			ppiBonus = int((500 - distance) * 200 / 500)
		}
	}

	return baseTier + Fitness(capBonus) + Fitness(resBonus) + Fitness(ppiBonus)
}

// TextFitness evaluates whether any registered font can satisfy the text-level
// constraints declared in reqs within the panel dimensions from hints.
// Returns Unsupported if no font can satisfy the constraints, Full otherwise.
// Returns Full (no-op) when both MinRows and MinCharsPerLine are 0.
func TextFitness(reqs SurfaceRequirements, hints textlayout.TextHints) Fitness {
	minRows := reqs.MinRows
	if minRows < 0 {
		minRows = 0
	}
	minChars := reqs.MinCharsPerLine
	if minChars < 0 {
		minChars = 0
	}

	// Both unconstrained — nothing to check.
	if minRows == 0 && minChars == 0 {
		return Full
	}

	for _, f := range font.List() {
		m := f.Metrics()
		if m.RowHeight == 0 || m.GlyphAdvance == 0 {
			continue
		}
		if minRows > 0 && minRows*m.RowHeight > hints.PixelHeight {
			continue
		}
		if minChars > 0 && minChars*m.GlyphAdvance > hints.PixelWidth {
			continue
		}
		// At least one font satisfies both constraints.
		return Full
	}

	return Unsupported
}
