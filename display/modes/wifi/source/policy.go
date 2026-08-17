package source

import (
	"fmt"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
)

// Allowed FGColor values for the WiFi mode.
var AllowedFGColors = []string{"cyan", "green", "amber", "red", "white", "none"}

// Allowed SignalDisplay values for the WiFi mode.
var AllowedSignalDisplay = []string{"bars", "percentage", "dbm"}

// Policy captures all runtime-configurable parameters for the WiFi mode.
type Policy struct {
	Style         string
	FGColor       string
	SignalDisplay string
	ShowFrequency bool
	ShowInterface bool
	ShowChannel   bool
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s|%s|%v|%v|%v", p.Style, p.FGColor, p.SignalDisplay, p.ShowFrequency, p.ShowInterface, p.ShowChannel)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style":          p.Style,
		"fgcolor":        p.FGColor,
		"signal_display": p.SignalDisplay,
		"show_frequency": p.ShowFrequency,
		"show_interface": p.ShowInterface,
		"show_channel":   p.ShowChannel,
	}
}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "fgcolor", Type: "string", Summary: "Foreground color for WiFi display elements.", Default: "green", Allowed: AllowedFGColors},
		{Key: "signal_display", Type: "string", Summary: "Signal strength display format.", Default: "bars", Allowed: AllowedSignalDisplay},
		{Key: "show_frequency", Type: "bool", Summary: "Show WiFi frequency.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "show_interface", Type: "bool", Summary: "Show wireless interface name.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "show_channel", Type: "bool", Summary: "Show WiFi channel.", Default: "true", Allowed: []string{"true", "false"}},
	}
}

func NormalizeText(value string) string { return strings.ToLower(strings.TrimSpace(value)) }
