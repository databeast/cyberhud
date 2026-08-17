// Package sampler provides system metrics collection for the htop-style views
// in the system display mode. It defines types for CPU state tracking and
// process information, along with a Sampler interface that abstracts the
// platform-specific data collection from /proc.
package sampler

// CPUState holds the raw counters from the previous /proc/stat read.
// Used to compute deltas on the next sample.
type CPUState struct {
	PerCore []CoreCounters
}

// CoreCounters holds per-core jiffies from /proc/stat.
type CoreCounters struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	IOWait  uint64
	IRQ     uint64
	SoftIRQ uint64
	Steal   uint64
}

// ProcessEntry describes a single process for the Top Consumers view.
type ProcessEntry struct {
	Name    string // Process command name (from /proc/[pid]/stat comm field)
	PID     int    // Process ID
	CPUPerc int    // CPU usage percentage (0-100 per core, can exceed 100 on multi-core)
}

// Sampler provides system metrics for the htop-style views.
type Sampler interface {
	// CPUSample returns per-core utilization as float64 values in [0.0, 1.0]
	// and an updated CPUState for the next call.
	// Returns a nil slice on read failure.
	CPUSample(prev CPUState) ([]float64, CPUState)

	// TopProcesses returns the top processes sorted by descending CPU usage.
	// Returns at most 20 entries. Returns a nil slice on read failure.
	TopProcesses() []ProcessEntry
}

// ProcSampler implements Sampler using the /proc filesystem.
// The actual implementation is in platform-specific build-tagged files.
type ProcSampler struct{}
