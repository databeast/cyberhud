//go:build !linux

package source

// gatherWifiState on non-Linux platforms returns a WifiData indicating that
// WiFi state reading is unavailable.
func GatherWifiState() WifiData {
	if data, ok := getTestOverride(); ok {
		return data
	}
	return WifiData{
		ConnectionState: Unavailable,
		StatusMessage:   "WiFi N/A",
	}
}
