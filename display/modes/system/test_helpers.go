package system

import (
	"github.com/databeast/cyberhud/display/modes/system/sampler"
	"github.com/databeast/cyberhud/display/modes/system/source"
)

func NormalizePolicy(p Policy) Policy { return normalizePolicy(p) }
func RegisteredStyleNames() []string  { return registeredStyleNames() }
func SetSamplerForTest(s sampler.Sampler) { source.SetSampler(s) }

// fakeSampler is a deterministic sampler for snapshot tests: it returns a fixed
// set of per-core utilizations and top processes so the cores/top styles render
// reproducible, representative content instead of "(no data)" on hosts whose
// real sampler returns nil (e.g. non-Linux dev machines).
type fakeSampler struct {
	cores     []float64
	processes []sampler.ProcessEntry
}

func (f fakeSampler) CPUSample(prev sampler.CPUState) ([]float64, sampler.CPUState) {
	return f.cores, prev
}

func (f fakeSampler) TopProcesses() []sampler.ProcessEntry { return f.processes }

// NewFakeSampler returns a deterministic Sampler seeded with representative CPU
// utilization and top-process data for snapshot rendering.
func NewFakeSampler() sampler.Sampler {
	return fakeSampler{
		cores: []float64{0.62, 0.18, 0.91, 0.07, 0.44, 0.33, 0.78, 0.12},
		processes: []sampler.ProcessEntry{
			{Name: "chromium", PID: 1421, CPUPerc: 87},
			{Name: "cyberhudd", PID: 902, CPUPerc: 41},
			{Name: "node", PID: 3310, CPUPerc: 29},
			{Name: "postgres", PID: 588, CPUPerc: 18},
			{Name: "systemd-journald", PID: 212, CPUPerc: 9},
			{Name: "sshd", PID: 744, CPUPerc: 4},
			{Name: "pipewire", PID: 1105, CPUPerc: 2},
		},
	}
}
