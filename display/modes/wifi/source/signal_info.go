package source

import (
	"fmt"
	"os"

	nlwifi "github.com/mdlayher/wifi"
)

func (g *Gatherer) signalInfo(iface string, bssSignal int32, stationInfo func() ([]*nlwifi.StationInfo, error)) (signalDBm, linkSpeedMbps int, unavailable bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.useBSSSignal {
		stations, err := stationInfo()
		switch {
		case err != nil:
			fmt.Fprintf(os.Stderr, "wifi: failed to get station info for %s: %v; using BSS signal for this mode session\n", iface, err)
			g.useBSSSignal = true
		case len(stations) == 0 || stations[0] == nil:
			fmt.Fprintf(os.Stderr, "wifi: no station info for %s; using BSS signal for this mode session\n", iface)
			g.useBSSSignal = true
		default:
			return stations[0].Signal, stations[0].TransmitBitrate / 1_000_000, false
		}
	}

	// The library represents a missing signal attribute as zero, indistinguishable
	// from 0 dBm, so only negative values are usable.
	if bssSignal >= 0 {
		if !g.missingSignalLogged {
			fmt.Fprintf(os.Stderr, "wifi: BSS signal unavailable for %s\n", iface)
			g.missingSignalLogged = true
		}
		return -100, 0, true
	}
	g.missingSignalLogged = false
	return int(bssSignal / 100), 0, false
}
