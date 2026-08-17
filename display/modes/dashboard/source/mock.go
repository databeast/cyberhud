package source

// MockDashboardContent returns a deterministic DashboardContent representing a
// realistic liveness panel state. Used by snapshot tests for reproducible
// output that does not depend on the hostname, IP, uptime, or WiFi state of
// the machine running the test (e.g. a GitHub Actions runner).
//
// Fixed values (no randomness):
//
//	Hostname: "cyberdeck-01", Uptime: "5d 12h", IP: "192.168.1.42",
//	WifiSSID: "CyberNet-5G", Version: "v1.0.0"
func MockDashboardContent() DashboardContent {
	return DashboardContent{
		Hostname:  "cyberdeck-01",
		Uptime:    "5d 12h",
		IPAddress: "192.168.1.42",
		WifiSSID:  "CyberNet-5G",
		Version:   "v1.0.0",
	}
}
