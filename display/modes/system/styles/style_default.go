package styles

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/widgets"
)

var DefaultStyle = def{
	name: "default",
	reqs: style.SurfaceRequirements{
		MinRows:         3,
		MinCharsPerLine: 10,
	},
	p: Params{BuildFn: defaultBuild},
}

func defaultBuild(snapshot source.SystemSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	_ = pol
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	items := []string{
		"Host: " + snapshot.Hostname,
		"OS: " + snapshot.OSArch,
		"Uptime: " + snapshot.Uptime,
	}
	if len(snapshot.IPs) == 0 {
		items = append(items, "IP: (none)")
	} else {
		for _, ip := range snapshot.IPs {
			items = append(items, "IP: "+ip)
		}
	}

	// Render status icon at top-right corner of content area.
	var sprites []widgets.Sprite
	if snapshot.GetIcon != nil {
		allDataOK := snapshot.Hostname != "unknown" && snapshot.Uptime != "n/a" && len(snapshot.IPs) > 0
		iconName := "check"
		if !allDataOK {
			iconName = "error"
		}
		if img, ok := snapshot.GetIcon(iconName); ok && img != nil {
			iconSize := 8
			spriteX := bridge.TopRightAnchorX(iconSize)
			_, oy := bridge.ContentOrigin()
			sprites = append(sprites, widgets.Sprite{
				Image:    img,
				Position: image.Point{X: spriteX, Y: oy},
				Label:    "status-" + iconName,
			})
		}
	}

	return style.ViewData{Items: items, Sprites: sprites}
}
