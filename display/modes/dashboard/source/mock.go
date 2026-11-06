package source

import "github.com/databeast/cyberhud/runtime/host"

// MockDashboardContent returns a deterministic DashboardContent representing a
// realistic liveness panel state. Used by snapshot tests for reproducible
// output that does not depend on the hostname, IP, uptime, or WiFi state of
// the machine running the test (e.g. a GitHub Actions runner).
//
// The Health verdict is deliberately WARN (not OK) so snapshots exercise the
// health row and its severity colorization — an all-OK fixture would never
// render anything but green and would hide a broken warn/crit path.
func MockDashboardContent() DashboardContent {
	return DashboardContent{
		Hostname:  "cyberdeck-01",
		Uptime:    "5d 12h",
		IPAddress: "192.168.1.42",
		WifiSSID:  "CyberNet-5G",
		Version:   "v1.0.0",
		Health: host.Health{
			Overall:  host.SeverityWarn,
			Degraded: []host.Contributor{{Name: "TEMP", Detail: "78°C", Severity: host.SeverityWarn}},
		},
	}
}
