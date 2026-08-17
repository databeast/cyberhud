package dashboard

import (
	"github.com/databeast/cyberhud/display/modes/dashboard/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

type Policy = source.Policy

func DefaultPolicy() Policy { return source.DefaultPolicy() }

func NormalizePolicy(p Policy) Policy { return normalizePolicy(p) }

func Registry() *style.StyleRegistry[source.DashboardContent, source.Policy] {
	return dashboardRegistry
}

var AllowedAccentsExported = source.AllowedAccents

func GetWifiSSIDExported() string { return source.GetWifiSSID() }

func GetVersionExported() string { return source.GetVersion() }

func GetPanelNameExported(hints textlayout.TextHints) string { return source.GetPanelName(hints) }
