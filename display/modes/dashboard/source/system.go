package source

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// getUptime returns the system uptime. On Linux it reads /proc/uptime;
// on other platforms it returns 0.
func getUptime() string {
	if runtime.GOOS != "linux" {
		return formatUptime(0)
	}
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return formatUptime(0)
	}
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return formatUptime(0)
	}
	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return formatUptime(0)
	}
	return formatUptime(time.Duration(seconds * float64(time.Second)))
}

func GetUptime() string { return getUptime() }

// getHostname returns the system hostname.
func getHostname() string {
	if h, err := os.Hostname(); err == nil && strings.TrimSpace(h) != "" {
		return h
	}
	return "unknown"
}

// getPrimaryIP returns the first non-loopback IPv4 address.
func getPrimaryIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
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
				return v4.String()
			}
		}
	}
	return ""
}

func formatUptime(d time.Duration) string {
	if d < 0 {
		d = 0
	}

	totalMinutes := int(d.Minutes())
	days := totalMinutes / (60 * 24)
	hours := (totalMinutes / 60) % 24
	minutes := totalMinutes % 60

	switch {
	case days >= 1:
		return fmt.Sprintf("%dd %dh", days, hours)
	case hours >= 1:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	default:
		return fmt.Sprintf("%dm", minutes)
	}
}

func FormatUptime(d time.Duration) string { return formatUptime(d) }
