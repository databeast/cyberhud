package styles

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
)

// connectedSnapshot returns a realistic connected WifiData for layout tests.
func connectedSnapshot() source.WifiData {
	return source.WifiData{
		SSID:            "CyberNet-5G",
		SignalStrength:  -45,
		LinkQuality:     78,
		IPAddress:       "192.168.1.42",
		Frequency:       5.18,
		Channel:         36,
		InterfaceName:   "wlan0",
		LinkSpeed:       867,
		ConnectionState: source.Connected,
	}
}

// The signal-forward layout composites text into sprites, so the policy-toggle
// behavior is verified against the string-assembly helpers that decide what
// text exists, rather than by scanning rendered rows.

func TestDetailLine_EmptyWhenAllTogglesOff(t *testing.T) {
	pol := source.DefaultPolicy() // detail toggles default OFF now
	if got := signalForwardDetailLine(connectedSnapshot(), pol); got != "" {
		t.Errorf("detail line = %q, want empty when all toggles off (glance-first default)", got)
	}
}

func TestDetailLine_ShowsInterfaceWhenEnabled(t *testing.T) {
	pol := source.DefaultPolicy()
	pol.ShowInterface = true
	if got := signalForwardDetailLine(connectedSnapshot(), pol); !strings.Contains(got, "wlan0") {
		t.Errorf("detail line = %q, want it to include interface when show_interface=true", got)
	}
}

func TestDetailLine_ConsolidatesFreqAndChannel(t *testing.T) {
	pol := source.DefaultPolicy()
	pol.ShowFrequency = true
	pol.ShowChannel = true
	got := signalForwardDetailLine(connectedSnapshot(), pol)
	if !strings.Contains(got, "Ch 36") || !strings.Contains(got, "GHz") {
		t.Errorf("detail line = %q, want a single consolidated 'Ch N (G GHz)' fragment", got)
	}
	// Channel and frequency must share one fragment, not appear as two separate ones.
	if strings.Count(got, "GHz") != 1 {
		t.Errorf("detail line = %q, want exactly one frequency mention", got)
	}
}

func TestCombinedRadioLine_HonorsToggles(t *testing.T) {
	snap := connectedSnapshot()
	if got := combinedRadioLine(snap, source.Policy{ShowChannel: true, ShowFrequency: true}); !strings.Contains(got, "Ch 36") {
		t.Errorf("both on: got %q", got)
	}
	if got := combinedRadioLine(snap, source.Policy{ShowChannel: true}); got != "Ch 36" {
		t.Errorf("channel only: got %q, want %q", got, "Ch 36")
	}
	if got := combinedRadioLine(snap, source.Policy{}); got != "" {
		t.Errorf("both off: got %q, want empty", got)
	}
}

// In bars mode the gauge center shows a plain-language quality word that names
// the arc (not a number, not jargon). In dbm/percentage the exact figure shows.
func TestReadout_BarsModeShowsQualityWord(t *testing.T) {
	got := signalForwardReadout(connectedSnapshot(), source.Policy{SignalDisplay: "bars"}, 4)
	if got != "STRONG" {
		t.Errorf("bars readout at level 4 = %q, want %q", got, "STRONG")
	}
	if strings.Contains(got, "dBm") || strings.Contains(got, "/") {
		t.Errorf("bars readout = %q, should be a word, not a number/fraction", got)
	}
}

func TestReadout_QualityWordPerLevel(t *testing.T) {
	want := map[int]string{4: "STRONG", 3: "GOOD", 2: "FAIR", 1: "WEAK", 0: "NONE"}
	for level, w := range want {
		if got := signalQualityWord(level); got != w {
			t.Errorf("signalQualityWord(%d) = %q, want %q", level, got, w)
		}
	}
}

func TestReadout_DbmMode(t *testing.T) {
	got := signalForwardReadout(connectedSnapshot(), source.Policy{SignalDisplay: "dbm"}, 4)
	if !strings.Contains(got, "dBm") {
		t.Errorf("dbm readout = %q, want a dBm value", got)
	}
}

func TestReadout_PercentageMode(t *testing.T) {
	got := signalForwardReadout(connectedSnapshot(), source.Policy{SignalDisplay: "percentage"}, 4)
	if !strings.Contains(got, "%") {
		t.Errorf("percentage readout = %q, want a %% value", got)
	}
}
