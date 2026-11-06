package source

import (
	"github.com/databeast/cyberhud/runtime/host"
)

// DashboardContent captures all state needed by dashboard style Build methods.
// The dashboard mode renders a glanceable liveness panel.
type DashboardContent struct {
	Hostname  string // System hostname
	Uptime    string // System uptime duration
	IPAddress string // First non-loopback IPv4 address (or "(no network)")
	WifiSSID  string // Connected SSID ≤32 chars, or "(no wifi)"
	Version   string // Build version ≤64 chars, or "dev"

	// Health is the aggregate liveness verdict (OK / DEGRADED / CRITICAL) plus
	// the contributors dragging it down, if any. It is what distinguishes the
	// dashboard from the dedicated per-subsystem modes: not thermal's zones or
	// system's CPU detail, but a single "is this box okay" roll-up with the
	// culprit named when it is not.
	Health host.Health
}

// BuildDashboardContent resolves current system state into a DashboardContent.
// Host facts (hostname, uptime, IP) and the health roll-up come from the shared
// runtime/host package; wifi SSID and build version are dashboard-specific. If a
// test override is set (via SetTestDashboardContent), returns that instead.
func BuildDashboardContent() DashboardContent {
	if data, ok := getTestOverride(); ok {
		return data
	}

	ip := host.PrimaryIP()
	if ip == "" {
		ip = "(no network)"
	}

	return DashboardContent{
		Hostname:  host.Hostname(),
		Uptime:    FormatUptime(host.Uptime()),
		IPAddress: ip,
		WifiSSID:  GetWifiSSID(),
		Version:   GetVersion(),
		Health:    host.Report(),
	}
}
