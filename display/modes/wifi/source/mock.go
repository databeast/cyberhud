package source

// MockWifiData returns a deterministic WifiData representing a realistic
// connected WiFi state. Used by snapshot tests for reproducible output.
//
// Fixed values (no randomness):
//
//	SSID: "CyberNet-5G", Signal: -45 dBm (bar level 4),
//	LinkQuality: 78, IP: "192.168.1.42", Frequency: 5.18 GHz (ch 36),
//	LinkSpeed: 867 Mbps, Interface: "wlan0"
func MockWifiData() WifiData {
	return WifiData{
		SSID:            "CyberNet-5G",
		SignalStrength:  -45,
		LinkQuality:     78,
		IPAddress:       "192.168.1.42",
		Frequency:       5.18,
		Channel:         frequencyToChannel(5.18), // 36
		InterfaceName:   "wlan0",
		LinkSpeed:       867,
		ConnectionState: Connected,
	}
}
