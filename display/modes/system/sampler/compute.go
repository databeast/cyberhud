package sampler

// ComputeUtilization calculates the CPU utilization between two counter
// snapshots as (nonIdleDelta / totalDelta), clamped to [0.0, 1.0].
// Returns 0.0 when totalDelta is zero.
func ComputeUtilization(prev, curr CoreCounters) float64 {
	prevTotal := Total(prev)
	currTotal := Total(curr)

	totalDelta := currTotal - prevTotal
	if totalDelta == 0 {
		return 0.0
	}

	idleDelta := (curr.Idle + curr.IOWait) - (prev.Idle + prev.IOWait)
	nonIdleDelta := totalDelta - idleDelta

	util := float64(nonIdleDelta) / float64(totalDelta)

	// Clamp to [0.0, 1.0]
	if util < 0.0 {
		return 0.0
	}
	if util > 1.0 {
		return 1.0
	}
	return util
}

// Total returns the sum of all fields in a CoreCounters.
func Total(c CoreCounters) uint64 {
	return c.User + c.Nice + c.System + c.Idle + c.IOWait + c.IRQ + c.SoftIRQ + c.Steal
}

// ComputeSample computes per-core utilization from previous and current CPU states.
// Returns nil if prev is nil, empty, or has a different core count than curr.
// This is the platform-independent computation logic used by CPUSample.
func ComputeSample(prev, curr CPUState) []float64 {
	if prev.PerCore == nil || len(prev.PerCore) != len(curr.PerCore) {
		return nil
	}

	utils := make([]float64, len(curr.PerCore))
	for i := range curr.PerCore {
		utils[i] = ComputeUtilization(prev.PerCore[i], curr.PerCore[i])
	}
	return utils
}

// SortAndCapProcesses takes a raw slice of ProcessEntry items, sorts them by
// descending CPUPerc, and returns at most 20 entries. This is the
// platform-independent logic used by TopProcesses.
func SortAndCapProcesses(entries []ProcessEntry) []ProcessEntry {
	if len(entries) == 0 {
		return entries
	}

	// Sort by descending CPU usage
	sorted := make([]ProcessEntry, len(entries))
	copy(sorted, entries)

	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].CPUPerc > sorted[i].CPUPerc {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Cap at 20 entries
	if len(sorted) > 20 {
		sorted = sorted[:20]
	}

	return sorted
}
