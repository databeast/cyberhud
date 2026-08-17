//go:build !linux

package sampler

// CPUSample on non-Linux platforms returns nil to signal that /proc/stat is unavailable.
func (ProcSampler) CPUSample(prev CPUState) ([]float64, CPUState) {
	return nil, CPUState{}
}

// TopProcesses on non-Linux platforms returns nil to signal that /proc is unavailable.
func (ProcSampler) TopProcesses() []ProcessEntry {
	return nil
}
