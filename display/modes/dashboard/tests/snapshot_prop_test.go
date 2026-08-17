package tests_test

import (
	"testing"

	"github.com/databeast/cyberhud/display/modes/dashboard"
	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"pgregory.net/rapid"
)

// For any input string to the WiFi SSID, Version, and Active Panel Name
// population functions, the resulting DashboardSnapshot fields SHALL satisfy:
// len(WifiSSID) ≤ 32, len(Version) ≤ 64, len(ActivePanelName) ≤ 64, and
// Version != "".

func TestProperty_FieldLengthConstraints(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// --- WifiSSID constraint: len ≤ 32 ---
		// On non-Linux (Windows), getWifiSSID always returns "(no wifi)"
		// which trivially satisfies ≤ 32. Verify the invariant holds.
		ssid := dashboard.GetWifiSSIDExported()
		if len(ssid) > 32 {
			t.Fatalf("WifiSSID length %d exceeds 32: %q", len(ssid), ssid)
		}

		// --- Version constraint: len ≤ 64 and non-empty ---
		// Inject a random string (0–200 chars) into dashboard.Version,
		// then call GetVersionExported and verify constraints.
		// Include whitespace chars to test the trimming logic.
		randomVersion := rapid.StringMatching(`[a-zA-Z0-9 \t\n]{0,200}`).Draw(t, "version")
		source.Version = randomVersion
		version := dashboard.GetVersionExported()
		if len(version) > 64 {
			t.Fatalf("Version length %d exceeds 64 for input %q: got %q", len(version), randomVersion, version)
		}
		if version == "" {
			t.Fatalf("Version must not be empty for input %q: got empty string", randomVersion)
		}

		// --- ActivePanelName constraint: len ≤ 64 ---
		// Create TextHints with random PanelProduct and ScreenName (0–200 chars)
		// then call GetPanelNameExported and verify length.
		panelProduct := rapid.StringMatching(`[a-zA-Z0-9\-_]{0,200}`).Draw(t, "panelProduct")
		screenName := rapid.StringMatching(`[a-zA-Z0-9\-_]{0,200}`).Draw(t, "screenName")
		hints := textlayout.TextHints{
			PanelProduct: panelProduct,
			ScreenName:   screenName,
		}
		panelName := dashboard.GetPanelNameExported(hints)
		if len(panelName) > 64 {
			t.Fatalf("ActivePanelName length %d exceeds 64 for product=%q screen=%q: got %q",
				len(panelName), panelProduct, screenName, panelName)
		}
	})
}
