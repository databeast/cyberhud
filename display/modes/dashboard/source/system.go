package source

import (
	"fmt"
	"time"

	"github.com/databeast/cyberhud/runtime/host"
)

// GetUptime returns the formatted system uptime string, composing the shared
// host read with the dashboard's formatter.
func GetUptime() string { return FormatUptime(host.Uptime()) }

// FormatUptime renders a duration as a compact "Nd Nh" / "Nh Nm" / "Nm" string.
// Host facts (hostname, uptime, IP) now come from runtime/host; this formatter
// is the one dashboard-specific presentation concern that remains here.
func FormatUptime(d time.Duration) string {
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
