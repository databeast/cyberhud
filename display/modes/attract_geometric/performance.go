package attract_geometric

import "github.com/databeast/cyberhud/display/modes/attract_geometric/source"

// PerformanceAction represents the type of scaling action to take.
type PerformanceAction int

const (
	PerfNone PerformanceAction = iota
	PerfReduce
	PerfRestore
)

// PerformanceDecision describes a performance scaling action and the new target count.
type PerformanceDecision struct {
	Action   PerformanceAction
	NewCount int
}

// evaluatePerformance processes a frame time and returns any scaling decision.
// It maintains a rolling window of frame times and triggers reduction when
// performance drops (avg > 33ms) or restoration when it recovers (avg < 25ms).
func evaluatePerformance(state *source.PerfState, frameTimeMs float64) PerformanceDecision {
	// Skip clamped frames entirely — don't add to window or evaluate.
	if frameTimeMs > source.ClampedFrameThresholdMs {
		return PerformanceDecision{Action: PerfNone}
	}

	// Add to rolling window.
	state.FrameTimes = append(state.FrameTimes, frameTimeMs)
	if len(state.FrameTimes) > source.WindowSize {
		state.FrameTimes = state.FrameTimes[1:]
	}

	// No decision until window is full.
	if len(state.FrameTimes) < source.WindowSize {
		return PerformanceDecision{Action: PerfNone}
	}

	// Compute average frame time.
	var sum float64
	for _, ft := range state.FrameTimes {
		sum += ft
	}
	avg := sum / float64(len(state.FrameTimes))

	// Reduction check (priority over restoration).
	if avg > source.ReductionThresholdMs && state.CurrentSquareCount > source.MinSquareCount {
		newCount := state.CurrentSquareCount / 2
		if newCount < source.MinSquareCount {
			newCount = source.MinSquareCount
		}
		state.HasReduced = true
		state.CurrentSquareCount = newCount
		return PerformanceDecision{Action: PerfReduce, NewCount: newCount}
	}

	// Restoration check.
	if avg < source.RestorationThresholdMs && state.HasReduced && state.CurrentSquareCount < state.OriginalSquareCount {
		removed := state.OriginalSquareCount - state.CurrentSquareCount
		restoreAmount := removed / 4
		if restoreAmount < 1 {
			restoreAmount = 1
		}
		newCount := state.CurrentSquareCount + restoreAmount
		if newCount > state.OriginalSquareCount {
			newCount = state.OriginalSquareCount
		}
		state.CurrentSquareCount = newCount
		if newCount == state.OriginalSquareCount {
			state.HasReduced = false
		}
		return PerformanceDecision{Action: PerfRestore, NewCount: newCount}
	}

	return PerformanceDecision{Action: PerfNone}
}
