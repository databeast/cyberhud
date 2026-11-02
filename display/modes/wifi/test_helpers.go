package wifi

import (
	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
)

var WifiRegistryExported = wifiRegistry
var AllowedFGColors = source.AllowedFGColors
var AllowedSignalDisplay = source.AllowedSignalDisplay

func RegisteredStyleNames() []string { return registeredStyleNames() }
func BuildViewForTest(_ interface{}, hints textlayout.TextHints) style.ViewData {
	return BuildViewWithHints(hints)
}

func SetTestWifiState(data source.WifiData) { source.SetTestWifiState(data) }
func ResetTestWifiState()                   { source.ResetTestWifiState() }
