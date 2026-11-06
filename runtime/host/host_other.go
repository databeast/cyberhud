//go:build !linux

package host

import (
	"net"
	"os"
	"strings"
	"time"
)

// Available reports whether host reads are backed by a real procfs/sysfs on
// this platform. False off Linux: hostname and IP still work via the stdlib,
// but /proc and /sys metrics are unavailable and return zero values.
func Available() bool { return false }

// Hostname returns the system hostname, or "unknown" on error. Works on all
// platforms via the standard library.
func Hostname() string {
	if h, err := os.Hostname(); err == nil && strings.TrimSpace(h) != "" {
		return h
	}
	return "unknown"
}

// OSArch returns the "goos/goarch" identity string.
func OSArch() string { return goosArch() }

// Uptime returns 0 off Linux (no portable uptime source wired here).
func Uptime() time.Duration { return 0 }

// IPv4s returns the IPv4 addresses of all non-loopback, up interfaces. Works
// on all platforms via the standard library.
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

// Load1 returns 0 off Linux.
func Load1() float64 { return 0 }

// Mem returns a zero reading off Linux.
func Mem() Memory { return Memory{} }

// HottestZone returns Present=false off Linux.
func HottestZone() Thermal { return Thermal{} }

// Net reports only interface liveness off Linux (no kernel counters available).
func Net() Network { return Network{AnyUp: len(IPv4s()) > 0} }

// ReportWith computes health off Linux using only what is available: network
// liveness. CPU/mem/thermal contribute nothing (SeverityOK) since they cannot
// be read, so the verdict reflects reachability alone.
func ReportWith(t Thresholds) Health {
	netw := Net()
	h := Health{Overall: SeverityOK}
	sev := classifyNet(netw)
	h.Overall = sev
	if sev > SeverityOK {
		detail := "down"
		h.Degraded = append(h.Degraded, Contributor{Name: "NET", Detail: detail, Severity: sev})
	}
	return h
}
