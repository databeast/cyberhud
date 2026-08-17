package source

import (
	"fmt"
	"image"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/databeast/cyberhud/display/modes/system/sampler"
)

// BuildItems returns system information rows for the display panel.
func BuildItems() []string {
	hostname := "unknown"
	if h, err := os.Hostname(); err == nil && strings.TrimSpace(h) != "" {
		hostname = h
	}

	items := []string{
		"Host: " + hostname,
		"OS: " + runtime.GOOS + "/" + runtime.GOARCH,
		"Uptime: " + Uptime(),
	}

	ips := IPAddresses()
	if len(ips) == 0 {
		items = append(items, "IP: (none)")
	} else {
		items = append(items, "IP: "+ips[0])
		for _, ip := range ips[1:] {
			items = append(items, "IP: "+ip)
		}
	}
	return items
}

// getHostname returns the system hostname.
func getHostname() string {
	if h, err := os.Hostname(); err == nil && strings.TrimSpace(h) != "" {
		return h
	}
	return "unknown"
}

// IPAddresses returns the IPv4 addresses of all non-loopback UP interfaces.
func IPAddresses() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	addrs := make([]string, 0, 4)
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		ifAddrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range ifAddrs {
			ip, _, err := net.ParseCIDR(a.String())
			if err != nil || ip == nil {
				continue
			}
			if v4 := ip.To4(); v4 != nil {
				addrs = append(addrs, v4.String())
			}
		}
	}
	return addrs
}

// Uptime reads the system uptime from /proc/uptime on Linux.
// Returns "n/a" on non-Linux systems or read errors.
func Uptime() string {
	b, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "n/a"
	}
	fields := strings.Fields(string(b))
	if len(fields) < 1 {
		return "n/a"
	}
	secs, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "n/a"
	}
	return FormatUptime(time.Duration(secs * float64(time.Second)))
}

// FormatUptime formats a duration as a human-readable uptime string.
func FormatUptime(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int64(d.Seconds())
	days := total / 86400
	hours := (total % 86400) / 3600
	mins := (total % 3600) / 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}

// BuildSnapshot assembles the model snapshot consumed by system styles.
func BuildSnapshot(cpuSample []float64, processes []sampler.ProcessEntry, getIcon func(name string) (image.Image, bool)) SystemSnapshot {
	return SystemSnapshot{
		Hostname:  getHostname(),
		OSArch:    runtime.GOOS + "/" + runtime.GOARCH,
		Uptime:    Uptime(),
		IPs:       IPAddresses(),
		CPUSample: cpuSample,
		Processes: processes,
		GetIcon:   getIcon,
	}
}

// ConsumeCPUSample returns and clears the CPU sample primed by PrimeCPUSample.
func ConsumeCPUSample() []float64 {
	sample := lastCPUSample
	lastCPUSample = nil
	return sample
}

// ConsumeTopProcesses returns and clears the process list primed by PrimeTopProcesses.
func ConsumeTopProcesses() []sampler.ProcessEntry {
	processes := lastProcesses
	lastProcesses = nil
	return processes
}

// SampleCPU samples CPU data without storing it for later consumption.
func SampleCPU() []float64 {
	if activeSampler == nil {
		return nil
	}
	var newState sampler.CPUState
	sample, newState := activeSampler.CPUSample(cpuState)
	cpuState = newState
	return sample
}

// PrimeCPUSample samples CPU data and stores it for the subsequent BuildView call.
func PrimeCPUSample() []float64 {
	lastCPUSample = SampleCPU()
	return lastCPUSample
}

// TopProcesses returns current top-process sampler data.
func TopProcesses() []sampler.ProcessEntry {
	if activeSampler == nil {
		return nil
	}
	return activeSampler.TopProcesses()
}

// PrimeTopProcesses samples top-process data and stores it for the subsequent BuildView call.
func PrimeTopProcesses() []sampler.ProcessEntry {
	lastProcesses = TopProcesses()
	return lastProcesses
}

// activeSampler holds the current sampler instance for CPU/process data.
var activeSampler sampler.Sampler = sampler.ProcSampler{}

// cpuState holds the previous CPU counters for delta computation.
var cpuState sampler.CPUState

// lastCPUSample holds the most recent CPU sample (set by RenderCacheKey).
var lastCPUSample []float64

// lastProcesses holds the most recent process list (set by RenderCacheKey).
var lastProcesses []sampler.ProcessEntry

// SetSampler replaces the active sampler (for testing with mocks).
func SetSampler(s sampler.Sampler) {
	activeSampler = s
}
