//go:build linux

package source

import (
	"fmt"
	"os"

	nlwifi "github.com/mdlayher/wifi"
)

// getWifiSSID queries the connected WiFi network name via NL80211 (netlink).
// It opens an NL80211 client, finds the first station-mode interface, retrieves
// the BSS, and returns the SSID truncated to 32 characters.
// On any error or when no BSS is available, it returns "(no wifi)".
func GetWifiSSID() string {
	c, err := nlwifi.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "dashboard: failed to open NL80211 client: %v\n", err)
		return "(no wifi)"
	}
	defer c.Close()

	interfaces, err := c.Interfaces()
	if err != nil {
		fmt.Fprintf(os.Stderr, "dashboard: failed to enumerate interfaces: %v\n", err)
		return "(no wifi)"
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
		return "(no wifi)"
	}

	// Retrieve the BSS (Basic Service Set) for the associated network.
	bss, err := c.BSS(ifi)
	if err != nil || bss == nil {
		return "(no wifi)"
	}

	ssid := bss.SSID
	if len(ssid) > 32 {
		ssid = ssid[:32]
	}
	if ssid == "" {
		return "(no wifi)"
	}
	return ssid
}
