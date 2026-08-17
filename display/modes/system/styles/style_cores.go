package styles

import (
	"fmt"
	"image"
	"strconv"

	"github.com/databeast/cyberhud/display/modes/system/source"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/progressbar"
)

// Compile-time check that progressbar.New returns a Renderable (used by Compositor).
var _ = progressbar.New(progressbar.Config{})

var CoresStyle = def{
	name: "cores",
	reqs: style.SurfaceRequirements{
		MinWidth:  32,
		MinHeight: 16,
	},
	p: Params{BuildFn: coresBuild},
}

func coresBuild(snapshot source.SystemSnapshot, pol source.Policy, ctx style.StyleContext, d def) style.ViewData {
	_ = pol
	sample := snapshot.CPUSample
	if sample == nil {
		return style.ViewData{Items: []string{"(no data)"}}
	}

	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})

	coreCount := len(sample)
	visibleRows := bridge.MaxVisibleRows()

	// Handle too-small case: fewer than 2 visible rows AND core count exceeds visible rows.
	if visibleRows < 2 && coreCount > visibleRows {
		return style.ViewData{Items: []string{"(too small)"}}
	}

	// Determine how many cores to render: min(coreCount, visibleRows).
	count := coreCount
	if visibleRows < count {
		count = visibleRows
	}

	// Compute layout dimensions.
	rowHeight := bridge.RowHeight()

	// Label width: reserve space for core index labels (e.g., "0 ", "1 ").
	glyphAdvance := bridge.GlyphAdvance()
	labelChars := len(strconv.Itoa(coreCount-1)) + 1
	labelWidth := labelChars * glyphAdvance

	pixelWidth := bridge.AvailableContentWidth()
	ox, oy := bridge.ContentOrigin()

	// Construct SuppressionContext from available dimensions.
	suppCtx := widgets.SuppressionContext{
		AvailableWidth:  pixelWidth,
		AvailableHeight: bridge.AvailableContentHeight(),
	}

	// Create Compositor for per-core progress bar sprites.
	comp := widgets.NewCompositor(suppCtx)

	var items []string
	for i := 0; i < count; i++ {
		bounds := image.Rect(ox+labelWidth, oy+i*rowHeight, ox+pixelWidth, oy+(i+1)*rowHeight)

		label := strconv.Itoa(i)
		items = append(items, label)

		// Add progress bar widget for this core via Compositor.
		comp.Add(progressbar.New(progressbar.Config{
			Style:  progressbar.Linear,
			Value:  sample[i],
			Bounds: bounds,
		}))
	}

	sprites := comp.Sprites()

	// If all sprite renders failed (compositor produced no sprites), fall back to text.
	if len(sprites) == 0 {
		textItems := make([]string, 0, count)
		for i := 0; i < count; i++ {
			pct := int(sample[i] * 100)
			textItems = append(textItems, fmt.Sprintf("%d: %d%%", i, pct))
		}
		return style.ViewData{Items: textItems}
	}

	return style.ViewData{
		Items:   items,
		Sprites: sprites,
	}
}
