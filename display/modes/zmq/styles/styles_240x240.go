package styles

import (
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
)

var ColorMedium240x240Style = def{
	name: "color-medium-240x240",
	reqs: style.SurfaceRequirements{
		MinWidth:  240,
		MinHeight: 240,
	},
	p: Params{BuildFn: colorMedium240x240Build},
}

func colorMedium240x240Build(snapshot content.ZMQData, pol content.Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	items := snapshot.Lines
	if len(items) == 0 {
		items = []string{"Waiting for messages..."}
	}

	maxVisible := bridge.MaxVisibleRows()

	var cursor, topRow int
	if maxVisible > 0 && len(items) > maxVisible {
		cursor, topRow = clampList(cursor, topRow, maxVisible, len(items))
	}

	return style.ViewData{
		Items:  items,
		Cursor: cursor,
		TopRow: topRow,
		Static: false,
	}
}

func clampList(cursor, topRow, maxVisible, itemCount int) (int, int) {
	if itemCount <= 0 {
		return 0, 0
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= itemCount {
		cursor = itemCount - 1
	}
	if maxVisible <= 0 || itemCount <= maxVisible {
		return cursor, 0
	}
	if topRow < 0 {
		topRow = 0
	}
	if topRow > itemCount-maxVisible {
		topRow = itemCount - maxVisible
	}
	if cursor < topRow {
		topRow = cursor
	}
	if cursor >= topRow+maxVisible {
		topRow = cursor - maxVisible + 1
	}
	return cursor, topRow
}
