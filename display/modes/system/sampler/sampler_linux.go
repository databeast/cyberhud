//go:build linux

package sampler

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// CPUSample reads /proc/stat, parses per-core CPU counters, computes
// utilization deltas against the previous state, and returns a slice of
// per-core utilization values in [0.0, 1.0].
//
// On the first call (prev.PerCore is nil) or when the core count changes,
// it returns a nil utilization slice along with the current state so that
// the next call can compute a valid delta.
//
// Returns a nil slice on any read failure.
func (ProcSampler) CPUSample(prev CPUState) ([]float64, CPUState) {
	cores, err := readProcStat()
	if err != nil {
		return nil, CPUState{}
	}

	curr := CPUState{PerCore: cores}

	utils := ComputeSample(prev, curr)
	return utils, curr
}

// readProcStat reads /proc/stat and returns CoreCounters for each per-core
// line (cpu0, cpu1, ...). The aggregate "cpu " line is skipped.
func readProcStat() ([]CoreCounters, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cores []CoreCounters
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip the aggregate cpu line (starts with "cpu " followed by space).
		// We only want per-core lines like "cpu0", "cpu1", etc.
		if !strings.HasPrefix(line, "cpu") {
			continue
		}
		// Distinguish "cpu " (aggregate) from "cpu0", "cpu1", etc.
		rest := line[3:]
		if len(rest) == 0 || rest[0] == ' ' {
			// This is the aggregate "cpu " line — skip it.
			continue
		}

		cc, err := parseCPULine(line)
		if err != nil {
			continue
		}
		cores = append(cores, cc)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(cores) == 0 {
		return nil, fmt.Errorf("no per-core cpu lines found in /proc/stat")
	}

	return cores, nil
}

// parseCPULine parses a single cpu line from /proc/stat.
// Expected format: "cpuN user nice system idle iowait irq softirq steal [guest guest_nice]"
func parseCPULine(line string) (CoreCounters, error) {
	fields := strings.Fields(line)
	// Need at least the label + 8 numeric fields.
	if len(fields) < 9 {
		return CoreCounters{}, fmt.Errorf("too few fields in cpu line: %q", line)
	}

	vals := make([]uint64, 8)
	for i := 0; i < 8; i++ {
		v, err := strconv.ParseUint(fields[i+1], 10, 64)
		if err != nil {
			return CoreCounters{}, fmt.Errorf("failed to parse field %d in %q: %w", i, line, err)
		}
		vals[i] = v
	}

	return CoreCounters{
		User:    vals[0],
		Nice:    vals[1],
		System:  vals[2],
		Idle:    vals[3],
		IOWait:  vals[4],
		IRQ:     vals[5],
		SoftIRQ: vals[6],
		Steal:   vals[7],
	}, nil
}

// TopProcesses reads /proc/[pid]/stat entries, computes CPU usage as a
// percentage of system uptime, sorts by descending CPU usage, and returns
// at most 20 entries. Returns nil if /proc cannot be read.
func (ProcSampler) TopProcesses() []ProcessEntry {
	uptime, err := readUptime()
	if err != nil || uptime <= 0 {
		return nil
	}

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}

	numCPUs := runtime.NumCPU()
	if numCPUs < 1 {
		numCPUs = 1
	}

	clkTck := uint64(100) // Standard clock ticks per second on Linux

	var processes []ProcessEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue // Not a numeric PID directory
		}

		name, utime, stime, err := readProcPidStat(pid)
		if err != nil {
			continue
		}

		// CPU% = (process_cpu_time_seconds / (uptime * num_cpus)) * 100
		processTimeSec := float64(utime+stime) / float64(clkTck)
		cpuPerc := int(processTimeSec / (uptime * float64(numCPUs)) * 100)
		if cpuPerc < 0 {
			cpuPerc = 0
		}

		processes = append(processes, ProcessEntry{
			Name:    name,
			PID:     pid,
			CPUPerc: cpuPerc,
		})
	}

	// Sort by descending CPU usage
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPUPerc > processes[j].CPUPerc
	})

	// Cap at 20 entries
	if len(processes) > 20 {
		processes = processes[:20]
	}

	return processes
}

// readUptime reads /proc/uptime and returns the system uptime in seconds.
func readUptime() (float64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0, fmt.Errorf("empty /proc/uptime")
	}
	return strconv.ParseFloat(fields[0], 64)
}

// readProcPidStat reads /proc/[pid]/stat and extracts the comm name,
// utime (field 14), and stime (field 15).
// The comm field is enclosed in parentheses and may contain spaces.
func readProcPidStat(pid int) (name string, utime, stime uint64, err error) {
	path := fmt.Sprintf("/proc/%d/stat", pid)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", 0, 0, err
	}

	content := string(data)

	// The comm field is enclosed in parentheses and may contain spaces or
	// special characters. Find the first '(' and last ')' to extract it.
	openParen := strings.IndexByte(content, '(')
	closeParen := strings.LastIndexByte(content, ')')
	if openParen < 0 || closeParen < 0 || closeParen <= openParen {
		return "", 0, 0, fmt.Errorf("cannot parse comm field in %s", path)
	}

	name = content[openParen+1 : closeParen]

	// Fields after the closing paren (space-separated).
	// Field numbering from /proc/[pid]/stat (1-indexed):
	//   1: pid, 2: comm, 3: state, ..., 14: utime, 15: stime
	// After extracting comm, the remaining fields start at field 3.
	// So utime is at index 11 (14 - 3) and stime is at index 12 (15 - 3)
	// in the zero-indexed remaining fields array.
	remaining := strings.TrimSpace(content[closeParen+1:])
	fields := strings.Fields(remaining)

	// We need at least 13 fields (indices 0..12 correspond to fields 3..15)
	if len(fields) < 13 {
		return "", 0, 0, fmt.Errorf("too few fields after comm in %s", path)
	}

	utime, err = strconv.ParseUint(fields[11], 10, 64)
	if err != nil {
		return "", 0, 0, fmt.Errorf("cannot parse utime in %s: %w", path, err)
	}

	stime, err = strconv.ParseUint(fields[12], 10, 64)
	if err != nil {
		return "", 0, 0, fmt.Errorf("cannot parse stime in %s: %w", path, err)
	}

	return name, utime, stime, nil
}
