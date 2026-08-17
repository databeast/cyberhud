package wifi

import (
	"fmt"

	"github.com/databeast/cyberhud/display/modes/wifi/source"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

func BuildView() style.ViewData {
	p := normalizePolicy(GetPolicy())
	data := source.GatherWifiState()

	hints, ok := getPanelHints()
	if !ok {
		return style.ViewData{Items: []string{"error"}}
	}

	s, reason := style.ResolveStyle(wifiRegistry, hints, "wifi", p.Style)
	ctx := style.NewStyleContext(hints)
	vd := s.Build(data, p, ctx)
	vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
	vd.Static = false
	vd.Cursor = -1
	if vd.Sprites == nil {
		vd.Sprites = []widgets.Sprite{}
	}
	return vd
}

func RenderCacheKey() uint32 {
	snap := source.GatherWifiState()
	p := GetPolicy()
	key := fmt.Sprintf("w|%d|%s|%d|%d|%s|%d|%d|%s|%s",
		int(snap.ConnectionState), snap.SSID, snap.SignalStrength, snap.LinkQuality,
		snap.IPAddress, snap.Channel, snap.LinkSpeed, snap.InterfaceName, p.Fingerprint())
	return region.CalcRegionCacheKey(key)
}

// BuildViewWithHints renders WiFi view data with explicit panel hints.
func BuildViewWithHints(hints textlayout.TextHints) style.ViewData {
	p := normalizePolicy(GetPolicy())
	data := source.GatherWifiState()
	s, reason := style.ResolveStyle(wifiRegistry, hints, "wifi", p.Style)
	ctx := style.NewStyleContext(hints)
	vd := s.Build(data, p, ctx)
	vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
	vd.Static = false
	vd.Cursor = -1
	if vd.Sprites == nil {
		vd.Sprites = []widgets.Sprite{}
	}
	return vd
}
