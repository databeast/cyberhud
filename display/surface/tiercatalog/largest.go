package tiercatalog

// selectLargest returns the candidate with the largest GlyphHeight from the
// given pool. Tie-breaking when multiple candidates share the same largest
// GlyphHeight: prefer higher family priority (Spleen > Terminus > Cozette >
// other), then lexicographically smaller FontID.
//
// The caller guarantees that candidates is non-empty.
func selectLargest(candidates []Candidate) Candidate {
	best := candidates[0]
	for _, cand := range candidates[1:] {
		candHeight := cand.MetricsAt(0).GlyphHeight
		bestHeight := best.MetricsAt(0).GlyphHeight
		switch {
		case candHeight > bestHeight:
			best = cand
		case candHeight == bestHeight:
			candPri := familyPriority(cand.ID())
			bestPri := familyPriority(best.ID())
			if candPri > bestPri {
				best = cand
			} else if candPri == bestPri && cand.ID() < best.ID() {
				best = cand
			}
		}
	}
	return best
}
