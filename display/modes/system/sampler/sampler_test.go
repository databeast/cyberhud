package sampler

import (
	"math"
	"testing"

	"pgregory.net/rapid"
)

// --- From: sampler_prop_test.go ---

// TestPropertyCPUComputation verifies that CPU utilization computation is correct.
//
// Property 8: CPU utilization computation correctness
//
// For any two valid CoreCounters snapshots prev and curr (where curr total ≥ prev total),
// the computed utilization SHALL equal clamp((nonIdleDelta) / (totalDelta), 0.0, 1.0)
// where totalDelta = sum(curr) - sum(prev) and nonIdleDelta = totalDelta -
// (curr.Idle + curr.IOWait - prev.Idle - prev.IOWait). When totalDelta = 0,
// utilization SHALL be 0.0. The output slice length SHALL equal the number of cores
// in the input.

func TestPropertyCPUComputation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate arbitrary prev CoreCounters with constrained uint64 values.
		// Use a moderate range to avoid overflow issues in totalDelta computation.
		prev := drawCoreCounters(t, "prev")

		// Generate curr such that each field >= corresponding prev field,
		// ensuring curr total >= prev total.
		curr := drawCoreCountersAtLeast(t, "curr", prev)

		// Compute expected utilization using the specification formula.
		prevTotal := Total(prev)
		currTotal := Total(curr)
		totalDelta := currTotal - prevTotal

		var expected float64
		if totalDelta == 0 {
			expected = 0.0
		} else {
			idleDelta := (curr.Idle + curr.IOWait) - (prev.Idle + prev.IOWait)
			nonIdleDelta := totalDelta - idleDelta
			expected = float64(nonIdleDelta) / float64(totalDelta)
			// Clamp to [0.0, 1.0]
			if expected < 0.0 {
				expected = 0.0
			}
			if expected > 1.0 {
				expected = 1.0
			}
		}

		// Compute actual utilization using the implementation.
		actual := ComputeUtilization(prev, curr)

		// Verify the result matches the expected formula within floating-point tolerance.
		if math.Abs(actual-expected) > 1e-12 {
			t.Fatalf("utilization mismatch: got %v, want %v (prev=%+v, curr=%+v, totalDelta=%d)",
				actual, expected, prev, curr, totalDelta)
		}

		// Verify clamping: result must be in [0.0, 1.0].
		if actual < 0.0 || actual > 1.0 {
			t.Fatalf("utilization out of range [0.0, 1.0]: got %v (prev=%+v, curr=%+v)",
				actual, prev, curr)
		}
	})
}

// TestPropertyCPUComputationSliceLength verifies that the output slice length
// equals the number of cores in the input when computing utilization for
// multiple cores.

func TestPropertyCPUComputationSliceLength(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random number of cores (1 to 32).
		numCores := rapid.IntRange(1, 32).Draw(t, "numCores")

		prevCores := make([]CoreCounters, numCores)
		currCores := make([]CoreCounters, numCores)

		for i := 0; i < numCores; i++ {
			prevCores[i] = drawCoreCounters(t, "prev")
			currCores[i] = drawCoreCountersAtLeast(t, "curr", prevCores[i])
		}

		// Compute utilization for all cores.
		utils := make([]float64, len(currCores))
		for i := range currCores {
			utils[i] = ComputeUtilization(prevCores[i], currCores[i])
		}

		// Verify output slice length equals number of cores.
		if len(utils) != numCores {
			t.Fatalf("slice length mismatch: got %d, want %d", len(utils), numCores)
		}
	})
}

// TestPropertyCPUComputationZeroDelta verifies that when totalDelta is zero,
// the utilization is 0.0.

func TestPropertyCPUComputationZeroDelta(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate arbitrary counters and use the same for both prev and curr.
		cc := drawCoreCounters(t, "cc")

		actual := ComputeUtilization(cc, cc)

		if actual != 0.0 {
			t.Fatalf("expected 0.0 when totalDelta is 0, got %v (counters=%+v)", actual, cc)
		}
	})
}

// drawCoreCounters generates a random CoreCounters with field values in a
// moderate range to keep sums within safe uint64 bounds.
func drawCoreCounters(t *rapid.T, label string) CoreCounters {
	// Use range [0, 10^12] — large enough to be realistic but sums won't overflow uint64.
	const maxVal = 1_000_000_000_000
	return CoreCounters{
		User:    rapid.Uint64Range(0, maxVal).Draw(t, label+"_user"),
		Nice:    rapid.Uint64Range(0, maxVal).Draw(t, label+"_nice"),
		System:  rapid.Uint64Range(0, maxVal).Draw(t, label+"_system"),
		Idle:    rapid.Uint64Range(0, maxVal).Draw(t, label+"_idle"),
		IOWait:  rapid.Uint64Range(0, maxVal).Draw(t, label+"_iowait"),
		IRQ:     rapid.Uint64Range(0, maxVal).Draw(t, label+"_irq"),
		SoftIRQ: rapid.Uint64Range(0, maxVal).Draw(t, label+"_softirq"),
		Steal:   rapid.Uint64Range(0, maxVal).Draw(t, label+"_steal"),
	}
}

// drawCoreCountersAtLeast generates a CoreCounters where each field is >= the
// corresponding field in min, ensuring curr total >= prev total.
func drawCoreCountersAtLeast(t *rapid.T, label string, min CoreCounters) CoreCounters {
	const maxDelta = 1_000_000_000
	return CoreCounters{
		User:    min.User + rapid.Uint64Range(0, maxDelta).Draw(t, label+"_user"),
		Nice:    min.Nice + rapid.Uint64Range(0, maxDelta).Draw(t, label+"_nice"),
		System:  min.System + rapid.Uint64Range(0, maxDelta).Draw(t, label+"_system"),
		Idle:    min.Idle + rapid.Uint64Range(0, maxDelta).Draw(t, label+"_idle"),
		IOWait:  min.IOWait + rapid.Uint64Range(0, maxDelta).Draw(t, label+"_iowait"),
		IRQ:     min.IRQ + rapid.Uint64Range(0, maxDelta).Draw(t, label+"_irq"),
		SoftIRQ: min.SoftIRQ + rapid.Uint64Range(0, maxDelta).Draw(t, label+"_softirq"),
		Steal:   min.Steal + rapid.Uint64Range(0, maxDelta).Draw(t, label+"_steal"),
	}
}

// TestPropertyProcessSamplerInvariants verifies process sampler output invariants.
//
// Property 9: Process sampler output invariants
//
// For any valid /proc data yielding K processes (K ≥ 0), the sampler output
// SHALL have length min(K, 20), SHALL be sorted by descending CPUPerc, and
// each entry SHALL have a non-empty Name, PID > 0, and CPUPerc ≥ 0.

func TestPropertyProcessSamplerInvariants(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random number of processes (0 to 50 to exercise the cap).
		numProcesses := rapid.IntRange(0, 50).Draw(t, "numProcesses")

		// Generate a slice of ProcessEntry with valid fields.
		entries := make([]ProcessEntry, numProcesses)
		for i := 0; i < numProcesses; i++ {
			// Name: non-empty string (1 to 20 chars from alphanumeric).
			name := rapid.StringMatching(`[a-zA-Z][a-zA-Z0-9_\-]{0,19}`).Draw(t, "name")
			// PID: positive integer.
			pid := rapid.IntRange(1, 1000000).Draw(t, "pid")
			// CPUPerc: non-negative integer (0 to 400 for multi-core).
			cpuPerc := rapid.IntRange(0, 400).Draw(t, "cpuPerc")

			entries[i] = ProcessEntry{
				Name:    name,
				PID:     pid,
				CPUPerc: cpuPerc,
			}
		}

		// Apply the sort-and-cap logic.
		result := SortAndCapProcesses(entries)

		// Invariant 1: Output length ≤ 20 and equals min(K, 20).
		expectedLen := numProcesses
		if expectedLen > 20 {
			expectedLen = 20
		}
		if len(result) != expectedLen {
			t.Fatalf("output length: got %d, want %d (input had %d entries)",
				len(result), expectedLen, numProcesses)
		}

		// Invariant 2: Sorted by descending CPUPerc.
		for i := 0; i < len(result)-1; i++ {
			if result[i].CPUPerc < result[i+1].CPUPerc {
				t.Fatalf("not sorted by descending CPUPerc at index %d: %d < %d",
					i, result[i].CPUPerc, result[i+1].CPUPerc)
			}
		}

		// Invariant 3: Each entry has non-empty Name.
		for i, entry := range result {
			if entry.Name == "" {
				t.Fatalf("entry %d has empty Name", i)
			}
		}

		// Invariant 4: Each entry has PID > 0.
		for i, entry := range result {
			if entry.PID <= 0 {
				t.Fatalf("entry %d has PID <= 0: %d", i, entry.PID)
			}
		}

		// Invariant 5: Each entry has CPUPerc ≥ 0.
		for i, entry := range result {
			if entry.CPUPerc < 0 {
				t.Fatalf("entry %d has CPUPerc < 0: %d", i, entry.CPUPerc)
			}
		}
	})
}

// --- From: sampler_test.go ---

// TestComputeUtilization_MultiCore verifies correct utilization calculation
// for a 4-core system with different usage patterns.
func TestComputeUtilization_MultiCore(t *testing.T) {
	tests := []struct {
		name     string
		prev     CoreCounters
		curr     CoreCounters
		wantUtil float64
	}{
		{
			name: "core0_50_percent_utilization",
			prev: CoreCounters{User: 100, Nice: 0, System: 50, Idle: 200, IOWait: 50, IRQ: 0, SoftIRQ: 0, Steal: 0},
			curr: CoreCounters{User: 200, Nice: 0, System: 100, Idle: 400, IOWait: 100, IRQ: 0, SoftIRQ: 0, Steal: 0},
			// totalDelta = 800-400 = 400, idleDelta = (400+100)-(200+50) = 250, nonIdle = 400-250 = 150
			// util = 150/400 = 0.375
			wantUtil: 0.375,
		},
		{
			name: "core1_high_utilization",
			prev: CoreCounters{User: 1000, Nice: 10, System: 500, Idle: 100, IOWait: 10, IRQ: 5, SoftIRQ: 5, Steal: 0},
			curr: CoreCounters{User: 1500, Nice: 20, System: 700, Idle: 120, IOWait: 15, IRQ: 10, SoftIRQ: 10, Steal: 0},
			// prevTotal = 1630, currTotal = 2375, totalDelta = 745
			// idleDelta = (120+15)-(100+10) = 25, nonIdle = 745-25 = 720
			// util = 720/745 ≈ 0.9664
			wantUtil: 720.0 / 745.0,
		},
		{
			name: "core2_mostly_idle",
			prev: CoreCounters{User: 10, Nice: 0, System: 5, Idle: 900, IOWait: 80, IRQ: 0, SoftIRQ: 0, Steal: 5},
			curr: CoreCounters{User: 15, Nice: 0, System: 8, Idle: 1800, IOWait: 160, IRQ: 0, SoftIRQ: 0, Steal: 7},
			// prevTotal = 1000, currTotal = 1990, totalDelta = 990
			// idleDelta = (1800+160)-(900+80) = 980, nonIdle = 990-980 = 10
			// util = 10/990 ≈ 0.0101
			wantUtil: 10.0 / 990.0,
		},
		{
			name: "core3_with_steal_time",
			prev: CoreCounters{User: 100, Nice: 5, System: 50, Idle: 200, IOWait: 20, IRQ: 5, SoftIRQ: 5, Steal: 15},
			curr: CoreCounters{User: 200, Nice: 10, System: 100, Idle: 300, IOWait: 30, IRQ: 10, SoftIRQ: 10, Steal: 40},
			// prevTotal = 400, currTotal = 700, totalDelta = 300
			// idleDelta = (300+30)-(200+20) = 110, nonIdle = 300-110 = 190
			// util = 190/300 ≈ 0.6333
			wantUtil: 190.0 / 300.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ComputeUtilization(tc.prev, tc.curr)
			if math.Abs(got-tc.wantUtil) > 1e-9 {
				t.Errorf("ComputeUtilization() = %v, want %v", got, tc.wantUtil)
			}
		})
	}
}

// TestComputeUtilization_SingleCore verifies computation works for a single core.
func TestComputeUtilization_SingleCore(t *testing.T) {
	prev := CoreCounters{User: 500, Nice: 10, System: 200, Idle: 800, IOWait: 50, IRQ: 10, SoftIRQ: 5, Steal: 0}
	curr := CoreCounters{User: 600, Nice: 15, System: 250, Idle: 900, IOWait: 60, IRQ: 15, SoftIRQ: 10, Steal: 0}
	// prevTotal = 1575, currTotal = 1850, totalDelta = 275
	// idleDelta = (900+60)-(800+50) = 110, nonIdle = 275-110 = 165
	// util = 165/275 = 0.6
	got := ComputeUtilization(prev, curr)
	want := 165.0 / 275.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("ComputeUtilization() = %v, want %v", got, want)
	}
}

// TestComputeUtilization_ZeroDelta verifies that zero total delta returns 0.0.
func TestComputeUtilization_ZeroDelta(t *testing.T) {
	counters := CoreCounters{User: 100, Nice: 10, System: 50, Idle: 200, IOWait: 20, IRQ: 5, SoftIRQ: 5, Steal: 10}
	got := ComputeUtilization(counters, counters)
	if got != 0.0 {
		t.Errorf("ComputeUtilization() with identical counters = %v, want 0.0", got)
	}
}

// TestComputeUtilization_FullUtilization verifies that 100% CPU returns 1.0.
func TestComputeUtilization_FullUtilization(t *testing.T) {
	prev := CoreCounters{User: 100, Nice: 0, System: 0, Idle: 500, IOWait: 0, IRQ: 0, SoftIRQ: 0, Steal: 0}
	curr := CoreCounters{User: 200, Nice: 0, System: 0, Idle: 500, IOWait: 0, IRQ: 0, SoftIRQ: 0, Steal: 0}
	// prevTotal = 600, currTotal = 700, totalDelta = 100
	// idleDelta = (500+0)-(500+0) = 0, nonIdle = 100-0 = 100
	// util = 100/100 = 1.0
	got := ComputeUtilization(prev, curr)
	if got != 1.0 {
		t.Errorf("ComputeUtilization() full utilization = %v, want 1.0", got)
	}
}

// TestComputeUtilization_Clamping verifies utilization is clamped to [0.0, 1.0].
func TestComputeUtilization_Clamping(t *testing.T) {
	t.Run("clamped_to_1.0", func(t *testing.T) {
		// A scenario where idle decreases (counter wrap or anomaly) could produce > 1.0
		// without clamping. We simulate by having idleDelta be negative (curr idle < prev idle).
		prev := CoreCounters{User: 0, Nice: 0, System: 0, Idle: 100, IOWait: 100, IRQ: 0, SoftIRQ: 0, Steal: 0}
		curr := CoreCounters{User: 100, Nice: 0, System: 100, Idle: 50, IOWait: 50, IRQ: 0, SoftIRQ: 0, Steal: 0}
		// prevTotal = 200, currTotal = 300, totalDelta = 100
		// idleDelta = (50+50)-(100+100) = -100, nonIdle = 100-(-100) = 200
		// util = 200/100 = 2.0 → clamped to 1.0
		got := ComputeUtilization(prev, curr)
		if got != 1.0 {
			t.Errorf("ComputeUtilization() should be clamped to 1.0, got %v", got)
		}
	})
}

// TestComputeSample_MultiCore verifies per-core utilization computation with 4 cores.
func TestComputeSample_MultiCore(t *testing.T) {
	prev := CPUState{
		PerCore: []CoreCounters{
			{User: 100, Nice: 0, System: 50, Idle: 200, IOWait: 50, IRQ: 0, SoftIRQ: 0, Steal: 0},
			{User: 200, Nice: 0, System: 100, Idle: 100, IOWait: 0, IRQ: 0, SoftIRQ: 0, Steal: 0},
			{User: 50, Nice: 0, System: 25, Idle: 400, IOWait: 25, IRQ: 0, SoftIRQ: 0, Steal: 0},
			{User: 300, Nice: 0, System: 150, Idle: 50, IOWait: 0, IRQ: 0, SoftIRQ: 0, Steal: 0},
		},
	}
	curr := CPUState{
		PerCore: []CoreCounters{
			{User: 200, Nice: 0, System: 100, Idle: 400, IOWait: 100, IRQ: 0, SoftIRQ: 0, Steal: 0},
			{User: 400, Nice: 0, System: 200, Idle: 100, IOWait: 0, IRQ: 0, SoftIRQ: 0, Steal: 0},
			{User: 60, Nice: 0, System: 30, Idle: 800, IOWait: 110, IRQ: 0, SoftIRQ: 0, Steal: 0},
			{User: 600, Nice: 0, System: 300, Idle: 100, IOWait: 0, IRQ: 0, SoftIRQ: 0, Steal: 0},
		},
	}

	utils := ComputeSample(prev, curr)
	if utils == nil {
		t.Fatal("ComputeSample returned nil, expected 4-element slice")
	}
	if len(utils) != 4 {
		t.Fatalf("ComputeSample returned %d values, want 4", len(utils))
	}

	// Verify each core is in [0.0, 1.0]
	for i, u := range utils {
		if u < 0.0 || u > 1.0 {
			t.Errorf("core %d utilization %v out of range [0.0, 1.0]", i, u)
		}
	}

	// Verify core 1 (all non-idle): totalDelta = 700-400 = 300, idleDelta = (100+0)-(100+0) = 0
	// nonIdle = 300, util = 300/300 = 1.0
	if math.Abs(utils[1]-1.0) > 1e-9 {
		t.Errorf("core 1 utilization = %v, want 1.0", utils[1])
	}
}

// TestComputeSample_FirstSampleReturnsNil verifies that the first sample (nil prev)
// returns nil utilization since no delta can be computed.
func TestComputeSample_FirstSampleReturnsNil(t *testing.T) {
	curr := CPUState{
		PerCore: []CoreCounters{
			{User: 100, Nice: 0, System: 50, Idle: 200, IOWait: 50, IRQ: 0, SoftIRQ: 0, Steal: 0},
		},
	}

	utils := ComputeSample(CPUState{}, curr)
	if utils != nil {
		t.Errorf("ComputeSample with nil prev.PerCore should return nil, got %v", utils)
	}
}

// TestComputeSample_CoreCountMismatchReturnsNil verifies that a change in core
// count between samples returns nil.
func TestComputeSample_CoreCountMismatchReturnsNil(t *testing.T) {
	prev := CPUState{
		PerCore: []CoreCounters{
			{User: 100, Nice: 0, System: 50, Idle: 200, IOWait: 50, IRQ: 0, SoftIRQ: 0, Steal: 0},
			{User: 200, Nice: 0, System: 100, Idle: 100, IOWait: 0, IRQ: 0, SoftIRQ: 0, Steal: 0},
		},
	}
	curr := CPUState{
		PerCore: []CoreCounters{
			{User: 200, Nice: 0, System: 100, Idle: 400, IOWait: 100, IRQ: 0, SoftIRQ: 0, Steal: 0},
		},
	}

	utils := ComputeSample(prev, curr)
	if utils != nil {
		t.Errorf("ComputeSample with mismatched core count should return nil, got %v", utils)
	}
}

// TestTotal verifies the sum of all fields in CoreCounters.
func TestTotal(t *testing.T) {
	c := CoreCounters{User: 10, Nice: 20, System: 30, Idle: 40, IOWait: 50, IRQ: 60, SoftIRQ: 70, Steal: 80}
	got := Total(c)
	want := uint64(10 + 20 + 30 + 40 + 50 + 60 + 70 + 80)
	if got != want {
		t.Errorf("Total() = %v, want %v", got, want)
	}
}

// TestComputeSample_SingleCore verifies computation works for a single core.
func TestComputeSample_SingleCore(t *testing.T) {
	prev := CPUState{
		PerCore: []CoreCounters{
			{User: 500, Nice: 10, System: 200, Idle: 800, IOWait: 50, IRQ: 10, SoftIRQ: 5, Steal: 0},
		},
	}
	curr := CPUState{
		PerCore: []CoreCounters{
			{User: 600, Nice: 15, System: 250, Idle: 900, IOWait: 60, IRQ: 15, SoftIRQ: 10, Steal: 0},
		},
	}

	utils := ComputeSample(prev, curr)
	if utils == nil {
		t.Fatal("ComputeSample returned nil, expected 1-element slice")
	}
	if len(utils) != 1 {
		t.Fatalf("ComputeSample returned %d values, want 1", len(utils))
	}

	// prevTotal = 1575, currTotal = 1850, totalDelta = 275
	// idleDelta = (900+60)-(800+50) = 110, nonIdle = 275-110 = 165
	// util = 165/275 = 0.6
	want := 165.0 / 275.0
	if math.Abs(utils[0]-want) > 1e-9 {
		t.Errorf("core 0 utilization = %v, want %v", utils[0], want)
	}
}
