package tiercatalog

// enforceMonotonicity walks tiers in ascending order and promotes any tier
// whose GlyphHeight is less than its predecessor's. This guarantees
// non-decreasing heights across the tier sequence and allows propagation
// (multiple tiers may collapse to the same entry).
func enforceMonotonicity(entries map[Tier]Entry) {
	for i := 1; i < len(tierOrder); i++ {
		prev := entries[tierOrder[i-1]]
		curr := entries[tierOrder[i]]
		if curr.GlyphHeight < prev.GlyphHeight {
			entries[tierOrder[i]] = prev
		}
	}
}
