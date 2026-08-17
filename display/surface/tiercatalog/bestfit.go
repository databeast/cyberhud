package tiercatalog

// bestFit selects the candidate closest to targetPx by GlyphHeight.
// Tie-breaking order:
//  1. Prefer smaller GlyphHeight (round down)
//  2. Prefer higher family priority (Spleen=3 > Terminus=2 > Cozette=1 > other=0)
//  3. Prefer lexicographically smaller FontID
//
// The caller guarantees candidates is non-empty.
func bestFit(candidates []Candidate, targetPx int) Candidate {
	best := candidates[0]
	for _, cand := range candidates[1:] {
		distBest := abs(best.MetricsAt(targetPx).GlyphHeight - targetPx)
		distCand := abs(cand.MetricsAt(targetPx).GlyphHeight - targetPx)

		switch {
		case distCand < distBest:
			best = cand
		case distCand == distBest:
			candH := cand.MetricsAt(targetPx).GlyphHeight
			bestH := best.MetricsAt(targetPx).GlyphHeight
			if candH < bestH {
				// Prefer smaller GlyphHeight (round down).
				best = cand
			} else if candH == bestH {
				candPri := familyPriority(cand.ID())
				bestPri := familyPriority(best.ID())
				if candPri > bestPri {
					// Prefer higher family priority.
					best = cand
				} else if candPri == bestPri {
					if cand.ID() < best.ID() {
						// Prefer lexicographically smaller FontID.
						best = cand
					}
				}
			}
		}
	}
	return best
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
