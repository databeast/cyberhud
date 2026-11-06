// Package host provides stateless reads of the machine this service runs on:
// uptime, hostname, IP addresses, CPU load, memory pressure, thermal, and
// network error counters. Every function returns a current value with no
// retained state and no background goroutine — modes that need sampling over
// time build their own sampler on top of these reads.
//
// The Linux implementation (host_linux.go) is backed by github.com/prometheus/procfs
// rather than hand-rolled /proc and /sys parsing. Non-Linux builds (host_other.go)
// return zero values and Available()==false so callers degrade gracefully.
package host

import "runtime"

// Severity classifies a health contributor from best to worst. The zero value
// is SeverityOK, so an unpopulated contributor reads as healthy.
type Severity int

const (
	SeverityOK Severity = iota
	SeverityWarn
	SeverityCritical
)

// String renders a severity as a short status word.
func (s Severity) String() string {
	switch s {
	case SeverityWarn:
		return "WARN"
	case SeverityCritical:
		return "CRIT"
	default:
		return "OK"
	}
}

// Memory is a point-in-time memory reading in bytes.
type Memory struct {
	TotalBytes     uint64
	AvailableBytes uint64
}

// UsedFraction returns used memory as a fraction [0,1] of total. Zero when
// total is unknown.
func (m Memory) UsedFraction() float64 {
	if m.TotalBytes == 0 {
		return 0
	}
	used := m.TotalBytes - m.AvailableBytes
	return float64(used) / float64(m.TotalBytes)
}

// Thermal is the hottest thermal zone found, in degrees Celsius. Present is
// false when the host exposes no thermal zones (common on non-Pi hardware).
type Thermal struct {
	Present bool
	Label   string
	TempC   float64
}

// Network summarizes interface liveness and the kernel's own error/drop
// counters across non-loopback interfaces — packet loss as the OS already
// counts it, not a probe we run.
type Network struct {
	AnyUp     bool   // at least one non-loopback interface is up with an IPv4
	RxErrors  uint64 // cumulative receive errors across up interfaces
	RxDropped uint64 // cumulative receive drops
	TxErrors  uint64 // cumulative transmit errors
	TxDropped uint64 // cumulative transmit drops
	RxBytes   uint64 // cumulative received bytes (denominator for a loss rate)
	TxBytes   uint64 // cumulative transmitted bytes
}

// Contributor is one named input to the overall health roll-up, carrying its
// own severity and a short human-readable detail for the UI to call out.
type Contributor struct {
	Name     string // short label, e.g. "CPU", "MEM", "TEMP", "NET"
	Detail   string // e.g. "3.8", "92%", "81°C", "12 drop"
	Severity Severity
}

// Health is the aggregate liveness verdict plus every contributor that is
// dragging it below OK. Overall is the worst contributor severity. Degraded
// lists only the warn/crit contributors, so a healthy box has an empty slice
// and the UI shows nothing to call out.
//
// Build one with [Report] (default thresholds) or [ReportWith].
type Health struct {
	Overall  Severity
	Degraded []Contributor
}

// StatusWord renders the overall verdict for a headline: OK / DEGRADED / CRITICAL.
func (h Health) StatusWord() string {
	switch h.Overall {
	case SeverityCritical:
		return "CRITICAL"
	case SeverityWarn:
		return "DEGRADED"
	default:
		return "OK"
	}
}

// Report computes the aggregate health verdict using DefaultThresholds.
func Report() Health { return ReportWith(DefaultThresholds()) }

// Thresholds parameterize the health classification. The zero value is not
// useful; use DefaultThresholds.
type Thresholds struct {
	// LoadWarnPerCore / LoadCritPerCore are 1-minute loadavg divided by core
	// count: warn at fully saturated, crit at heavily oversubscribed.
	LoadWarnPerCore float64
	LoadCritPerCore float64
	// MemWarn / MemCrit are used-memory fractions [0,1].
	MemWarn float64
	MemCrit float64
	// TempWarnC / TempCritC are hottest-zone temperatures in Celsius. Used only
	// when a zone exposes no "critical" trip point of its own.
	TempWarnC float64
	TempCritC float64
}

// DefaultThresholds are sensible defaults for a small single-board host.
func DefaultThresholds() Thresholds {
	return Thresholds{
		LoadWarnPerCore: 1.0,
		LoadCritPerCore: 2.0,
		MemWarn:         0.85,
		MemCrit:         0.95,
		TempWarnC:       75,
		TempCritC:       85,
	}
}

// classifyLoad reduces a 1-minute load average to a severity against the core
// count. cores <= 0 is treated as 1.
func classifyLoad(load1 float64, cores int, t Thresholds) Severity {
	if cores < 1 {
		cores = 1
	}
	perCore := load1 / float64(cores)
	switch {
	case perCore >= t.LoadCritPerCore:
		return SeverityCritical
	case perCore >= t.LoadWarnPerCore:
		return SeverityWarn
	default:
		return SeverityOK
	}
}

// classifyMem reduces a memory reading to a severity.
func classifyMem(m Memory, t Thresholds) Severity {
	f := m.UsedFraction()
	switch {
	case f >= t.MemCrit:
		return SeverityCritical
	case f >= t.MemWarn:
		return SeverityWarn
	default:
		return SeverityOK
	}
}

// classifyTemp reduces the hottest zone to a severity against the thresholds.
func classifyTemp(th Thermal, t Thresholds) Severity {
	if !th.Present {
		return SeverityOK
	}
	switch {
	case th.TempC >= t.TempCritC:
		return SeverityCritical
	case th.TempC >= t.TempWarnC:
		return SeverityWarn
	default:
		return SeverityOK
	}
}

// classifyNet reduces network state to a severity: critical when nothing is
// up, warn when the kernel is counting errors/drops, otherwise ok.
func classifyNet(n Network) Severity {
	if !n.AnyUp {
		return SeverityCritical
	}
	if n.RxErrors+n.RxDropped+n.TxErrors+n.TxDropped > 0 {
		return SeverityWarn
	}
	return SeverityOK
}

// worst returns the higher of two severities.
func worst(a, b Severity) Severity {
	if b > a {
		return b
	}
	return a
}

// goosArch and numCPU are tiny runtime wrappers kept here (platform-neutral) so
// both the Linux and stub implementations share one definition.

func goosArch() string { return runtime.GOOS + "/" + runtime.GOARCH }

func numCPU() int {
	if n := runtime.NumCPU(); n > 0 {
		return n
	}
	return 1
}
