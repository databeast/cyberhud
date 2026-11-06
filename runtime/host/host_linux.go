//go:build linux

package host

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/sysfs"
)

// Available reports whether host reads are backed by a real procfs/sysfs on
// this platform. Always true on Linux.
func Available() bool { return true }

// fs returns the default /proc filesystem handle, or ok=false if /proc is not
// mounted (unusual, but handled so callers degrade rather than panic).
func fs() (procfs.FS, bool) {
	f, err := procfs.NewDefaultFS()
	if err != nil {
		return procfs.FS{}, false
	}
	return f, true
}

// Hostname returns the system hostname, or "unknown" on error.
func Hostname() string {
	if h, err := os.Hostname(); err == nil && strings.TrimSpace(h) != "" {
		return h
	}
	return "unknown"
}

// OSArch returns the "goos/goarch" identity string.
func OSArch() string { return goosArch() }

// Uptime returns the system uptime, derived from the kernel boot time in
// /proc/stat. Returns 0 if unavailable.
func Uptime() time.Duration {
	f, ok := fs()
	if !ok {
		return 0
	}
	st, err := f.Stat()
	if err != nil || st.BootTime == 0 {
		return 0
	}
	boot := time.Unix(int64(st.BootTime), 0)
	d := time.Since(boot)
	if d < 0 {
		return 0
	}
	return d
}

// IPv4s returns the IPv4 addresses of all non-loopback, up interfaces.
func IPv4s() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var out []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ip, _, err := net.ParseCIDR(a.String())
			if err != nil || ip == nil {
				continue
			}
			if v4 := ip.To4(); v4 != nil {
				out = append(out, v4.String())
			}
		}
	}
	return out
}

// PrimaryIP returns the first non-loopback IPv4 address, or "" if none.
func PrimaryIP() string {
	if ips := IPv4s(); len(ips) > 0 {
		return ips[0]
	}
	return ""
}

// Load1 returns the 1-minute load average, or 0 on error.
func Load1() float64 {
	f, ok := fs()
	if !ok {
		return 0
	}
	la, err := f.LoadAvg()
	if err != nil || la == nil {
		return 0
	}
	return la.Load1
}

// Mem returns the current memory reading. TotalBytes is 0 when unavailable.
func Mem() Memory {
	f, ok := fs()
	if !ok {
		return Memory{}
	}
	mi, err := f.Meminfo()
	if err != nil {
		return Memory{}
	}
	var m Memory
	if mi.MemTotalBytes != nil {
		m.TotalBytes = *mi.MemTotalBytes
	} else if mi.MemTotal != nil {
		m.TotalBytes = *mi.MemTotal * 1024
	}
	// Prefer MemAvailable (accounts for reclaimable cache) over MemFree.
	switch {
	case mi.MemAvailableBytes != nil:
		m.AvailableBytes = *mi.MemAvailableBytes
	case mi.MemAvailable != nil:
		m.AvailableBytes = *mi.MemAvailable * 1024
	case mi.MemFreeBytes != nil:
		m.AvailableBytes = *mi.MemFreeBytes
	case mi.MemFree != nil:
		m.AvailableBytes = *mi.MemFree * 1024
	}
	return m
}

// HottestZone returns the hottest thermal zone, or Present=false when the host
// exposes none.
func HottestZone() Thermal {
	sfs, err := sysfs.NewDefaultFS()
	if err != nil {
		return Thermal{}
	}
	zones, err := sfs.ClassThermalZoneStats()
	if err != nil || len(zones) == 0 {
		return Thermal{}
	}
	hottest := Thermal{}
	for _, z := range zones {
		tempC := float64(z.Temp) / 1000.0
		if !hottest.Present || tempC > hottest.TempC {
			hottest = Thermal{Present: true, Label: z.Type, TempC: tempC}
		}
	}
	return hottest
}

// Net summarizes interface liveness and the kernel's cumulative error/drop
// counters across non-loopback interfaces.
func Net() Network {
	n := Network{AnyUp: len(IPv4s()) > 0}

	f, ok := fs()
	if !ok {
		return n
	}
	devs, err := f.NetDev()
	if err != nil {
		return n
	}
	for name, d := range devs {
		if name == "lo" || strings.HasPrefix(name, "lo") {
			continue
		}
		n.RxErrors += d.RxErrors
		n.RxDropped += d.RxDropped
		n.TxErrors += d.TxErrors
		n.TxDropped += d.TxDropped
		n.RxBytes += d.RxBytes
		n.TxBytes += d.TxBytes
	}
	return n
}

// ReportWith computes the aggregate health verdict using the given thresholds.
func ReportWith(t Thresholds) Health {
	cores := numCPU()
	load1 := Load1()
	mem := Mem()
	therm := HottestZone()
	netw := Net()

	type c struct {
		name   string
		detail string
		sev    Severity
	}
	cands := []c{
		{"CPU", fmt.Sprintf("%.1f", load1), classifyLoad(load1, cores, t)},
		{"MEM", fmt.Sprintf("%d%%", int(mem.UsedFraction()*100)), classifyMem(mem, t)},
		{"NET", netDetail(netw), classifyNet(netw)},
	}
	if therm.Present {
		cands = append(cands, c{"TEMP", fmt.Sprintf("%.0f°C", therm.TempC), classifyTemp(therm, t)})
	}

	h := Health{Overall: SeverityOK}
	for _, cc := range cands {
		h.Overall = worst(h.Overall, cc.sev)
		if cc.sev > SeverityOK {
			h.Degraded = append(h.Degraded, Contributor{Name: cc.name, Detail: cc.detail, Severity: cc.sev})
		}
	}
	return h
}

// netDetail renders a short network detail string: "down" when nothing is up,
// a drop/error count when the kernel is counting them, else "ok".
func netDetail(n Network) string {
	if !n.AnyUp {
		return "down"
	}
	bad := n.RxErrors + n.RxDropped + n.TxErrors + n.TxDropped
	if bad > 0 {
		return fmt.Sprintf("%d drop", bad)
	}
	return "ok"
}
