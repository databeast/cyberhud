package tiercatalog

import (
	"fmt"

	"github.com/databeast/cyberhud/display/surface/fonts"
)

// buildCandidatePool computes the advance budget, wraps all registered fonts
// as bitmapCandidates, and filters to those whose GlyphAdvance fits within
// the budget. Returns the filtered candidates, the computed advance budget,
// and an error if no fonts are registered or none qualify.
func buildCandidatePool(p Params) ([]Candidate, int, error) {
	// Normalize MinChars: default to 10 when <= 0.
	minChars := p.MinChars
	if minChars <= 0 {
		minChars = 10
	}

	// Compute advance budget via integer division.
	advanceBudget := p.PixelWidth / minChars

	// Retrieve all registered fonts.
	allFonts := font.List()
	if len(allFonts) == 0 {
		return nil, 0, fmt.Errorf("tiercatalog: no fonts registered")
	}

	// Wrap each font as a bitmapCandidate and filter by advance budget.
	var candidates []Candidate
	for _, f := range allFonts {
		bc := bitmapCandidate{face: f}
		if bc.MetricsAt(0).GlyphAdvance <= advanceBudget {
			candidates = append(candidates, bc)
		}
	}

	// If no candidates qualify, find the smallest advance for diagnostic info.
	if len(candidates) == 0 {
		smallestAdvance := allFonts[0].Metrics().GlyphAdvance
		for _, f := range allFonts[1:] {
			if adv := f.Metrics().GlyphAdvance; adv < smallestAdvance {
				smallestAdvance = adv
			}
		}
		return nil, advanceBudget, fmt.Errorf(
			"tiercatalog: no font fits region %dx%d with MinChars=%d (smallest advance=%d)",
			p.PixelWidth, p.PixelHeight, minChars, smallestAdvance,
		)
	}

	return candidates, advanceBudget, nil
}
