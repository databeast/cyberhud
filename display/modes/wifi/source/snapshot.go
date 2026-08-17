package source

import (
	"math"
)

// ConnectionState represents the WiFi adapter state.
type ConnectionState int

const (
	// Connected indicates the wireless interface is associated and has link.
	Connected ConnectionState = iota
	// Disconnected indicates the wireless interface exists but has no association.
	Disconnected
	// Unavailable indicates no wireless interface was detected or the platform
	// does not support WiFi state reading (e.g., non-Linux OS).
	Unavailable
)

// WifiData holds all WiFi state consumed by style Build methods.
// It is produced by gatherWifiState and passed to the active style's Build method.
type WifiData struct {
	SSID            string  // network name, max 32 chars
	SignalStrength  int     // dBm, -100 to 0
	LinkQuality     int     // 0-100 percentage
	IPAddress       string  // IPv4 address string
	Frequency       float64 // GHz
	Channel         int     // derived from Frequency via frequencyToChannel
	InterfaceName   string  // e.g. "wlan0"
	LinkSpeed       int     // Mbps, 0 if unavailable
	ConnectionState ConnectionState
	StatusMessage   string // "No Network" or "WiFi N/A" when disconnected/unavailable
}

// frequencyToChannel maps a WiFi frequency in GHz to the corresponding channel number.
// Returns 0 for frequencies that do not match any known 2.4 GHz or 5 GHz channel.
//
// Mapping rules:
//   - 2.412–2.472 GHz → channels 1–13: ch = (freq_MHz - 2407) / 5
//   - 2.484 GHz → channel 14 (special case, Japan)
//   - 5.180–5.885 GHz → channels 36–177: ch = (freq_MHz - 5000) / 5
//   - Other → 0 (unknown)
func frequencyToChannel(freqGHz float64) int {
	// Convert GHz to MHz for integer arithmetic.
	freqMHz := int(math.Round(freqGHz * 1000))

	// 2.4 GHz band: channels 1–13 (2412–2472 MHz)
	if freqMHz >= 2412 && freqMHz <= 2472 {
		return (freqMHz - 2407) / 5
	}

	// 2.4 GHz band: channel 14 (2484 MHz, special case)
	if freqMHz == 2484 {
		return 14
	}

	// 5 GHz band: channels 36–177 (5180–5885 MHz)
	if freqMHz >= 5180 && freqMHz <= 5885 {
		return (freqMHz - 5000) / 5
	}

	return 0
}

// FrequencyToChannel maps a WiFi frequency in GHz to a channel number.
func FrequencyToChannel(freqGHz float64) int { return frequencyToChannel(freqGHz) }
