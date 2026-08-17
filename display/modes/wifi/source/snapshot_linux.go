//go:build linux

package source

import (
	"fmt"
	"net"
	"os"

	nlwifi "github.com/mdlayher/wifi"
)

// gatherWifiState queries wireless interface status via NL80211 (netlink) and
// returns a fully populated WifiData. On any failure the function logs the
// error and returns a snapshot with ConnectionState set to Unavailable or
// Disconnected as appropriate — it never panics.
func GatherWifiState() WifiData {
	if data, ok := getTestOverride(); ok {
		return data
	}

	c, err := nlwifi.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wifi: failed to open NL80211 client: %v\n", err)
		return WifiData{
			ConnectionState: Unavailable,
			StatusMessage:   "WiFi N/A",
		}
	}
	defer c.Close()

	interfaces, err := c.Interfaces()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wifi: failed to enumerate interfaces: %v\n", err)
		return WifiData{
			ConnectionState: Unavailable,
			StatusMessage:   "WiFi N/A",
		}
	}

	// Find the first station-mode interface.
	var ifi *nlwifi.Interface
	for _, iface := range interfaces {
		if iface.Type == nlwifi.InterfaceTypeStation {
			ifi = iface
			break
		}
	}
	if ifi == nil {
		return WifiData{
			ConnectionState: Unavailable,
			StatusMessage:   "WiFi N/A",
		}
	}

	// Check if we're associated by retrieving the BSS.
	bss, err := c.BSS(ifi)
	if err != nil || bss == nil {
		return WifiData{
			InterfaceName:   ifi.Name,
			ConnectionState: Disconnected,
			StatusMessage:   "No Network",
		}
	}

	// We have an associated BSS — gather station info for signal & bitrate.
	var signalDBm int
	var linkSpeedMbps int
	stations, err := c.StationInfo(ifi)
	if err == nil && len(stations) > 0 {
		sta := stations[0]
		signalDBm = sta.Signal
		// TransmitBitrate is in bits/second; convert to Mbps.
		linkSpeedMbps = sta.TransmitBitrate / 1_000_000
	}

	// Frequency from BSS is in MHz; convert to GHz.
	freqGHz := float64(bss.Frequency) / 1000.0
	channel := frequencyToChannel(freqGHz)
	ipAddr := resolveIPAddress(ifi.Name)

	return WifiData{
		SSID:            truncateSSID(bss.SSID),
		SignalStrength:  clampSignal(signalDBm),
		LinkQuality:     signalToQuality(signalDBm),
		IPAddress:       ipAddr,
		Frequency:       freqGHz,
		Channel:         channel,
		InterfaceName:   ifi.Name,
		LinkSpeed:       linkSpeedMbps,
		ConnectionState: Connected,
	}
}

// signalToQuality converts a dBm value to a link quality percentage (0–100)
// using the same linear mapping as the old /proc/net/wireless approach:
// -110 dBm → 0%, -40 dBm → 100%, linearly interpolated and clamped.
func signalToQuality(dBm int) int {
	if dBm <= -110 {
		return 0
	}
	if dBm >= -40 {
		return 100
	}
	return (dBm + 110) * 100 / 70
}

// resolveIPAddress finds the first IPv4 address on the named wireless interface.
func resolveIPAddress(ifaceName string) string {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wifi: failed to look up interface %s: %v\n", ifaceName, err)
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wifi: failed to get addresses for %s: %v\n", ifaceName, err)
		return ""
	}

	for _, a := range addrs {
		ip, _, err := net.ParseCIDR(a.String())
		if err != nil || ip == nil {
			continue
		}
		if v4 := ip.To4(); v4 != nil {
			return v4.String()
		}
	}

	return ""
}

// truncateSSID ensures the SSID does not exceed 32 characters.
func truncateSSID(ssid string) string {
	if len(ssid) > 32 {
		return ssid[:32]
	}
	return ssid
}

// clampSignal clamps the signal level to the range -100 to 0.
func clampSignal(dBm int) int {
	if dBm < -100 {
		return -100
	}
	if dBm > 0 {
		return 0
	}
	return dBm
}
