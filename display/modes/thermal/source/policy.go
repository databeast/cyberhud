package source

import (
	"fmt"

	"github.com/databeast/cyberhud/display/catalog"
)

// Policy captures all runtime-configurable parameters for the thermal mode.
type Policy struct {
	Style          string // "overview", "detail", "graph", "minimal"
	Font           string // "auto" or a registered font ID
	RefreshMS      int    // sampling interval in milliseconds [500, 120000]
	WarnThreshold  int    // warning temperature in °C [0, ∞)
	CritThreshold  int    // critical temperature in °C [0, ∞), must be > WarnThreshold
	Unit           string // "C" or "F"
	FGColor        string // "thermal", "cyan", "green", "amber", "red", "white", "none"
	ShowLED        bool   // LED activity indicator
	ShowRefreshBar bool   // refresh progress bar
	ShowBorder     bool   // border frame around panel
}

// AllowedUnits lists valid temperature display units.
var AllowedUnits = []string{"C", "F"}

// AllowedFGColors lists valid foreground color values for the thermal mode.
var AllowedFGColors = []string{"thermal", "cyan", "green", "amber", "red", "white", "none"}

func (p Policy) Options() []catalog.OptionDefinition {
	return []catalog.OptionDefinition{
		{Key: "font", Type: "string", Summary: "Font selection (auto or a registered font ID).", Default: "auto"},
		{Key: "refresh_ms", Type: "int", Summary: "Sampling interval in milliseconds.", Default: "2000"},
		{Key: "warn_threshold", Type: "int", Summary: "Warning temperature threshold in °C.", Default: "70"},
		{Key: "crit_threshold", Type: "int", Summary: "Critical temperature threshold in °C.", Default: "90"},
		{Key: "unit", Type: "string", Summary: "Temperature display unit.", Default: "C", Allowed: AllowedUnits},
		{Key: "fgcolor", Type: "string", Summary: "Foreground color palette.", Default: "thermal", Allowed: AllowedFGColors},
		{Key: "show_led", Type: "bool", Summary: "Show LED activity indicator.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "show_refresh_bar", Type: "bool", Summary: "Show refresh progress bar.", Default: "true", Allowed: []string{"true", "false"}},
		{Key: "show_border", Type: "bool", Summary: "Show border frame.", Default: "false", Allowed: []string{"true", "false"}},
	}
}

// DefaultPolicy returns the default thermal policy.
func DefaultPolicy() Policy {
	return Policy{
		Style:          "",
		Font:           "auto",
		RefreshMS:      2000,
		WarnThreshold:  70,
		CritThreshold:  90,
		Unit:           "C",
		FGColor:        "thermal",
		ShowLED:        true,
		ShowRefreshBar: true,
		ShowBorder:     false,
	}
}

func (p Policy) Fingerprint() string {
	return fmt.Sprintf("%s|%s|%d|%d|%d|%s|%s|%v|%v|%v",
		p.Style, p.Font, p.RefreshMS, p.WarnThreshold, p.CritThreshold,
		p.Unit, p.FGColor, p.ShowLED, p.ShowRefreshBar, p.ShowBorder)
}

func (p Policy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"style":            p.Style,
		"font":             p.Font,
		"refresh_ms":       p.RefreshMS,
		"warn_threshold":   p.WarnThreshold,
		"crit_threshold":   p.CritThreshold,
		"unit":             p.Unit,
		"fgcolor":          p.FGColor,
		"show_led":         p.ShowLED,
		"show_refresh_bar": p.ShowRefreshBar,
		"show_border":      p.ShowBorder,
	}
}
