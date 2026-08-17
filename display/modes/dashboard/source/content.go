package source

// DashboardSnapshot captures all state needed by dashboard style Build methods.
// The dashboard mode renders a glanceable liveness panel.
type DashboardContent struct {
	Hostname  string // System hostname
	Uptime    string // System uptime duration
	IPAddress string // First non-loopback IPv4 address (or "(no network)")
	WifiSSID  string // Connected SSID ≤32 chars, or "(no wifi)"
	Version   string // Build version ≤64 chars, or "dev"
}

// buildSnapshot resolves current system state into a DashboardSnapshot.
// If a test override is set (via SetTestDashboardContent), returns that
// instead of querying live system state.
func BuildDashboardContent() DashboardContent {
	if data, ok := getTestOverride(); ok {
		return data
	}

	ip := getPrimaryIP()
	if ip == "" {
		ip = "(no network)"
	}

	return DashboardContent{
		Hostname:  getHostname(),
		Uptime:    getUptime(),
		IPAddress: ip,
		WifiSSID:  GetWifiSSID(),
		Version:   GetVersion(),
	}
}
