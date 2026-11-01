package styles

import (
	"image"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
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

// itemsJoined joins the rendered item rows for substring assertions.
func itemsJoined(vd style.ViewData) string {
	return strings.Join(vd.Items, "\n")
}

// fastColorDef mirrors the registration of a color-fast style large enough to
// show every detail row (800x600 fits the full stack).
func fastColorDef() def {
	return def{
		name: "color-fast-800x600",
		reqs: style.SurfaceRequirements{MinWidth: 800, MinHeight: 600, Capability: style.ColorFast},
		p:    Params{Fast: true, Color: true},
	}
}

func buildFast(t *testing.T, pol source.Policy) style.ViewData {
	t.Helper()
	hints := textlayout.DefaultTextHints(image.Rect(0, 0, 800, 600))
	ctx := style.NewStyleContext(hints)
	return buildFastWifi(connectedSnapshot(), pol, ctx, fastColorDef())
}

// The shared fast engine previously ignored show_frequency/show_channel/
// show_interface. These tests lock in that the toggles now take effect.
func TestFastWifi_HidesInterfaceWhenDisabled(t *testing.T) {
	pol := source.DefaultPolicy()
	pol.ShowInterface = false

	vd := buildFast(t, pol)
	if strings.Contains(itemsJoined(vd), "wlan0") {
		t.Errorf("interface row shown despite show_interface=false: %q", vd.Items)
	}
}

func TestFastWifi_ShowsInterfaceWhenEnabled(t *testing.T) {
	pol := source.DefaultPolicy()
	pol.ShowInterface = true

	vd := buildFast(t, pol)
	if !strings.Contains(itemsJoined(vd), "wlan0") {
		t.Errorf("interface row missing despite show_interface=true: %q", vd.Items)
	}
}

func TestFastWifi_HidesRadioLineWhenBothDisabled(t *testing.T) {
	pol := source.DefaultPolicy()
	pol.ShowFrequency = false
	pol.ShowChannel = false

	vd := buildFast(t, pol)
	joined := itemsJoined(vd)
	if strings.Contains(joined, "GHz") || strings.Contains(joined, "Ch ") {
		t.Errorf("radio line shown despite freq+channel disabled: %q", vd.Items)
	}
}

func TestFastWifi_ConsolidatesFreqAndChannelOnOneRow(t *testing.T) {
	pol := source.DefaultPolicy() // both enabled
	vd := buildFast(t, pol)

	var radioRows int
	for _, item := range vd.Items {
		if strings.Contains(item, "GHz") || strings.Contains(item, "Ch ") {
			radioRows++
		}
	}
	if radioRows != 1 {
		t.Errorf("expected exactly one consolidated radio row, got %d: %q", radioRows, vd.Items)
	}
	if !strings.Contains(itemsJoined(vd), "Ch 36") || !strings.Contains(itemsJoined(vd), "GHz") {
		t.Errorf("consolidated radio row should include channel and frequency: %q", vd.Items)
	}
}

// In bars mode the signal sprite stands alone — the redundant "Signal: N/4"
// text must be gone. In dbm/percentage mode the numeric label is kept.
func TestFastWifi_DropsRedundantSignalLabelInBarsMode(t *testing.T) {
	pol := source.DefaultPolicy()
	pol.SignalDisplay = "bars"

	vd := buildFast(t, pol)
	if strings.Contains(itemsJoined(vd), "Signal:") {
		t.Errorf("redundant signal label shown in bars mode: %q", vd.Items)
	}
}

func TestFastWifi_KeepsSignalLabelInDbmMode(t *testing.T) {
	pol := source.DefaultPolicy()
	pol.SignalDisplay = "dbm"

	vd := buildFast(t, pol)
	if !strings.Contains(itemsJoined(vd), "dBm") {
		t.Errorf("dbm signal label missing in dbm mode: %q", vd.Items)
	}
}
