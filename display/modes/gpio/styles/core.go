package styles

// Core layout assembly for GPIO styles.
//
// Every per-resolution style in this package is declared as a def value in its
// styles_WxH.go file: a name, surface requirements, and a Params block selecting
// one of the shared GPIO layouts. Set Params.BuildFn for a fully bespoke layout.

import (
	"fmt"
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/style/layout"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/surface/tiercatalog"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/borderframe"
	"github.com/databeast/cyberhud/display/widgets/led"
	"github.com/databeast/cyberhud/display/widgets/textlabel"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

type layoutKind uint8

const (
	layoutNotImplemented layoutKind = iota
	layoutGrayscaleFast
	layoutEink
	layoutList
	layoutIcons
	layoutDetail
	layoutDashboard
	layoutActivity
	layoutColorFastRows
	layoutMonoRows
)

// Params holds the hand-tweakable knobs for a GPIO style declaration.
// Layout selects a real, existing GPIO rendering behavior; the zero value keeps
// the legacy skeleton "NOT IMPLEMENTED" output.
type Params struct {
	Layout  layoutKind
	BuildFn func(data GpioSnapshot, pol Policy, ctx style.StyleContext, d def) style.ViewData
}

type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var _ style.Style[GpioSnapshot, Policy] = def{}

func (d def) Name() string { return d.name }

func (d def) Requirements() style.SurfaceRequirements { return d.reqs }

func (d def) Supports(hints textlayout.TextHints) style.Fitness {
	return style.EvaluateFitness(d.reqs, hints)
}

func (d def) Build(data GpioSnapshot, pol Policy, ctx style.StyleContext) style.ViewData {
	if d.p.BuildFn != nil {
		return d.p.BuildFn(data, pol, ctx, d)
	}
	switch d.p.Layout {
	case layoutGrayscaleFast:
		return grayscaleFastBuild(data, pol, ctx, d)
	case layoutEink:
		return einkBuild(data, pol, ctx, d)
	case layoutList:
		return listBuild(data, pol, ctx, d)
	case layoutIcons:
		return iconsBuild(data, pol, ctx, d)
	case layoutDetail:
		return detailBuild(data, pol, ctx, d)
	case layoutDashboard:
		return dashboardBuild(data, pol, ctx, d)
	case layoutActivity:
		return activityBuild(data, pol, ctx, d)
	case layoutColorFastRows:
		return colorFastRowsBuild(data, pol, ctx, d)
	case layoutMonoRows:
		return monoRowsBuild(data, pol, ctx, d)
	default:
		return notImplementedBuild(data, pol, ctx, d)
	}
}

func notImplementedBuild(_ GpioSnapshot, _ Policy, _ style.StyleContext, _ def) style.ViewData {
	return style.ViewData{Items: []string{"NOT IMPLEMENTED"}}
}

func grayscaleFastBuild(snap GpioSnapshot, _ Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildGrayscaleFastStyle(snap, ctx, d.reqs)
}

func einkBuild(snap GpioSnapshot, _ Policy, ctx style.StyleContext, d def) style.ViewData {
	return buildEinkStyle(snap, ctx, d.reqs)
}

func listBuild(snap GpioSnapshot, _ Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildListInternal(snap, ctx)
}

func iconsBuild(snap GpioSnapshot, _ Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildIconsInternal(snap, ctx)
}

func detailBuild(snap GpioSnapshot, _ Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildDetailInternal(snap, ctx)
}

func dashboardBuild(snap GpioSnapshot, _ Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildDashboardInternal(snap, ctx)
}

func activityBuild(snap GpioSnapshot, _ Policy, ctx style.StyleContext, _ def) style.ViewData {
	return buildActivityInternal(snap, ctx)
}

func maxPinLabelLen(pins []gpiomgr.PinState) int {
	maxLen := 6 // minimum "## X X"
	for _, p := range pins {
		l := len(fmt.Sprintf("%d %s %s", p.Number, p.Mode.String(), pinLevelStr(p)))
		if l > maxLen {
			maxLen = l
		}
	}
	return maxLen
}

func pinLevelStr(p gpiomgr.PinState) string {
	if p.Level {
		return "HI"
	}
	return "LO"
}

func pinRowLabel(p gpiomgr.PinState) string {
	return fmt.Sprintf("%-2d %s %s", p.Number, p.Mode.String(), pinLevelStr(p))
}

func configuredPins(pins []gpiomgr.PinState) []gpiomgr.PinState {
	out := make([]gpiomgr.PinState, 0, len(pins))
	for _, p := range pins {
		if p.Mode != gpiomgr.ModeUnknown {
			out = append(out, p)
		}
	}
	return out
}

func noConfiguredItems(totalPins, maxChars int) []string {
	items := []string{
		"GPIO idle",
		fmt.Sprintf("Configured 0/%d", totalPins),
	}
	if maxChars > 0 {
		for i, item := range items {
			items[i] = textlayout.Truncate(item, maxChars)
		}
	}
	return items
}

func buildEinkStyle(snap GpioSnapshot, sctx style.StyleContext, reqs style.SurfaceRequirements) style.ViewData {
	// 1. Construct own LayoutBridge from hints.
	hints := sctx.Hints()
	_ = hints
	bridge := sctx.Layout(0)

	// 2. Check for zero content area.
	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{Static: true}
	}

	// 3. Use tier catalog metrics for layout calculations.
	entry, ok := sctx.FontCatalog().Get(tiercatalog.TierNormal)
	if !ok {
		entry = tiercatalog.Entry{
			GlyphWidth:   bridge.GlyphAdvance(),
			GlyphHeight:  bridge.RowHeight(),
			GlyphAdvance: bridge.GlyphAdvance(),
			RowHeight:    bridge.RowHeight(),
		}
	}
	rowHeight := entry.RowHeight
	if rowHeight <= 0 {
		rowHeight = 10 // default row height
	}
	maxRows := bridge.AvailableContentHeight() / rowHeight

	// 4. Create Compositor with SuppressOnEink rule.
	suppCtx := widgets.SuppressionContext{
		IsEink:          true,
		AvailableWidth:  bridge.AvailableContentWidth(),
		AvailableHeight: bridge.AvailableContentHeight(),
	}
	comp := widgets.NewCompositor(suppCtx, widgets.SuppressOnEink())

	// 5. Build text rows for pin states.
	visiblePins := configuredPins(snap.Pins)
	if len(visiblePins) == 0 {
		items := noConfiguredItems(len(snap.Pins), maxCharsFromBridge(bridge))
		offsets := make([]int, len(items))
		for i, item := range items {
			offsets[i] = bridge.CenterXWith(len(item), entry.GlyphAdvance)
		}
		tiers := make([]tiercatalog.Tier, len(items))
		for i := range tiers {
			tiers[i] = tiercatalog.TierNormal
		}
		return style.ViewData{
			Items:       items,
			Tiers:       tiers,
			LineOffsets: offsets,
			Static:      true,
		}
	}
	if maxRows > 0 && len(visiblePins) > maxRows {
		visiblePins = visiblePins[:maxRows]
	}

	items := make([]string, len(visiblePins))
	accent := resolveFGColor(snap.Policy.FGColor)

	originX, originY := bridge.ContentOrigin()

	for i, p := range visiblePins {
		items[i] = p.String()

		// Add LED indicator via compositor (will be suppressed if not eink-safe).
		state := led.Off
		if p.Level {
			state = led.On
		}
		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   led.DiameterForRow(rowHeight),
			Bounds:     image.Rectangle{Min: image.Pt(originX, originY+i*rowHeight)},
			Foreground: accent,
		}))
	}

	// 6. Compute centering for LineOffsets.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = bridge.CenterXWith(len(item), entry.GlyphAdvance)
	}

	// 7. Build tiers slice (all TierNormal).
	tiers := make([]tiercatalog.Tier, len(items))
	for i := range tiers {
		tiers[i] = tiercatalog.TierNormal
	}

	// 8. Return ViewData with Static=true, nil Colors.
	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		Colors:      nil,
		Static:      true,
		Sprites:     comp.Sprites(),
	}
}

func longestPinString(pins []gpiomgr.PinState) int {
	maxLen := 0
	for _, p := range pins {
		if l := len(p.String()); l > maxLen {
			maxLen = l
		}
	}
	return maxLen
}

func monoSummaryBuild(snap GpioSnapshot, pol Policy, ctx style.StyleContext, d def) style.ViewData {
	// 1. Construct own LayoutBridge from hints.
	hints := ctx.Hints()
	_ = hints
	bridge := ctx.Layout(0)

	// 2. Check for zero content area.
	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{}
	}

	// 3. Compute summary: "Pins:{total} Hi:{high}"
	total := len(snap.Pins)
	high := 0
	for _, p := range snap.Pins {
		if p.Level {
			high++
		}
	}
	summary := fmt.Sprintf("Pins:%d Hi:%d", total, high)
	items := []string{summary}

	// 4. Use tier catalog metrics for layout calculations.
	entry, ok := ctx.FontCatalog().Get(tiercatalog.TierNormal)
	if !ok {
		entry = tiercatalog.Entry{
			GlyphWidth:   bridge.GlyphAdvance(),
			GlyphHeight:  bridge.RowHeight(),
			GlyphAdvance: bridge.GlyphAdvance(),
			RowHeight:    bridge.RowHeight(),
		}
	}

	// 5. Compute per-row centering offsets.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = bridge.CenterXWith(len(item), entry.GlyphAdvance)
	}

	// 6. Build tiers slice (all TierNormal).
	tiers := make([]tiercatalog.Tier, len(items))
	for i := range tiers {
		tiers[i] = tiercatalog.TierNormal
	}

	// 7. Return ViewData with nil Colors and Static=true.
	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		Colors:      nil,
		Static:      true,
	}
}

func monoRowsBuild(snap GpioSnapshot, pol Policy, ctx style.StyleContext, d def) style.ViewData {
	// 1. Construct own LayoutBridge from hints.
	hints := ctx.Hints()
	_ = hints
	bridge := ctx.Layout(0)

	// 2. Check for zero content area.
	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{}
	}

	// 3. Determine visible rows and character width.
	maxRows := bridge.MaxVisibleRows()
	maxChars := maxCharsFromBridge(bridge)

	visiblePins := configuredPins(snap.Pins)
	if len(visiblePins) == 0 {
		items := noConfiguredItems(len(snap.Pins), maxChars)
		offsets := make([]int, len(items))
		for i, item := range items {
			offsets[i] = bridge.CenterXWith(len(item), bridge.GlyphAdvance())
		}
		tiers := make([]tiercatalog.Tier, len(items))
		for i := range tiers {
			tiers[i] = tiercatalog.TierNormal
		}
		return style.ViewData{
			Items:       items,
			Tiers:       tiers,
			LineOffsets: offsets,
			Colors:      nil,
			Static:      true,
		}
	}
	if maxRows > 0 && len(visiblePins) > maxRows {
		visiblePins = visiblePins[:maxRows]
	}

	// 4. Build per-pin text items, truncated to maxChars.
	items := BuildItemsTruncated(visiblePins, maxChars)

	// 5. Use tier catalog metrics for layout calculations.
	entry, ok := ctx.FontCatalog().Get(tiercatalog.TierNormal)
	if !ok {
		entry = tiercatalog.Entry{
			GlyphWidth:   bridge.GlyphAdvance(),
			GlyphHeight:  bridge.RowHeight(),
			GlyphAdvance: bridge.GlyphAdvance(),
			RowHeight:    bridge.RowHeight(),
		}
	}

	// 6. Compute per-row centering offsets.
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = bridge.CenterXWith(len(item), entry.GlyphAdvance)
	}

	// 7. Build tiers slice (all TierNormal).
	tiers := make([]tiercatalog.Tier, len(items))
	for i := range tiers {
		tiers[i] = tiercatalog.TierNormal
	}

	// 8. Return ViewData with nil Colors and Static=true.
	return style.ViewData{
		Items:       items,
		Tiers:       tiers,
		LineOffsets: offsets,
		Colors:      nil,
		Static:      true,
	}
}

func colorFastRowsBuild(snap GpioSnapshot, pol Policy, ctx style.StyleContext, d def) style.ViewData {
	p := snap.Policy

	// 1. Construct own LayoutBridge from hints.
	hints := ctx.Hints()
	_ = hints
	bridge := ctx.Layout(0)

	// 2. Check for zero content area.
	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{}
	}

	// 3. Use tier catalog metrics for layout calculations.
	entry, ok := ctx.FontCatalog().Get(tiercatalog.TierNormal)
	if !ok {
		entry = tiercatalog.Entry{
			GlyphWidth:   bridge.GlyphAdvance(),
			GlyphHeight:  bridge.RowHeight(),
			GlyphAdvance: bridge.GlyphAdvance(),
			RowHeight:    bridge.RowHeight(),
		}
	}

	// 4. Resolve accent colors.
	accent := resolveFGColor(p.FGColor)
	dimmed := dimFGColor(p.FGColor)

	// 6. Build LED sprites per pin using Compositor.
	// Distribute rows evenly across available height so pins fill the panel.
	contentH := bridge.AvailableContentHeight()
	numPins := len(snap.Pins)
	if numPins < 1 {
		numPins = 1
	}
	minRowHeight := entry.GlyphHeight + 2
	if minRowHeight < 8 {
		minRowHeight = 8
	}
	rowHeight := contentH / numPins
	if rowHeight < minRowHeight {
		rowHeight = minRowHeight
	}

	// LED indicators sized to match row height for alignment with text.
	diameter := led.DiameterForRow(rowHeight)
	suppCtx := widgets.SuppressionContext{
		AvailableWidth:  bridge.AvailableContentWidth(),
		AvailableHeight: bridge.AvailableContentHeight(),
	}
	comp := widgets.NewCompositor(suppCtx)

	originX, originY := bridge.ContentOrigin()

	visiblePins := snap.Pins
	maxRows := bridge.AvailableContentHeight() / rowHeight
	// Ensure last sprite bottom ((maxRows-1)*rowHeight + diameter) fits within content area.
	if diameter > rowHeight && maxRows > 0 {
		maxBySprite := (bridge.AvailableContentHeight()-diameter)/rowHeight + 1
		if maxBySprite < 0 {
			maxBySprite = 0
		}
		if maxBySprite < maxRows {
			maxRows = maxBySprite
		}
	}
	if maxRows > 0 && len(visiblePins) > maxRows {
		visiblePins = visiblePins[:maxRows]
	}

	var items []string
	var colors []color.Color

	for i, pin := range visiblePins {
		// LED sprite for each pin, vertically centered within the row.
		state := led.Off
		fg := dimmed
		if pin.Level {
			state = led.On
			fg = accent
		}
		ledYOffset := (rowHeight - diameter) / 2
		if ledYOffset < 0 {
			ledYOffset = 0
		}
		ledY := originY + i*rowHeight + ledYOffset
		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(originX, ledY)},
			Foreground: fg,
		}))

	}

	// 7. Compute LineOffsets: position text right of the LED indicator.
	var offsets []int
	if len(items) > 0 {
		ledTextGap := diameter + 4 // LED diameter + 4px gap
		offsets = make([]int, len(items))
		for i := range offsets {
			offsets[i] = originX + ledTextGap
		}
	}

	// Handle color-disabled policy.
	if !p.Color {
		colors = nil
	}

	return style.ViewData{
		Items:       items,
		LineOffsets: offsets,
		Colors:      colors,
		Sprites:     comp.Sprites(),
	}
}

func colorFast128Build(snap GpioSnapshot, pol Policy, ctx style.StyleContext, d def) style.ViewData {
	p := snap.Policy

	// 1. Construct own LayoutBridge from hints.
	hints := ctx.Hints()
	_ = hints
	bridge := ctx.Layout(0)

	// 2. Return early for zero content area (Requirement 6.3).
	if bridge.AvailableContentWidth() == 0 || bridge.AvailableContentHeight() == 0 {
		return style.ViewData{}
	}

	// 3. Return static with empty sprites for zero-pin input (Requirement 7.1).
	numPins := len(snap.Pins)
	if numPins == 0 {
		return style.ViewData{Static: true}
	}

	// 4. Resolve font face for textlabel rendering.
	entry := ctx.Entry(tiercatalog.TierNormal)
	selectedFont := ctx.Face("spleen", tiercatalog.TierNormal)

	glyphHeight := entry.GlyphHeight

	// 5. Compute two-column layout geometry.
	leftCount := (numPins + 1) / 2 // ceil(numPins/2)
	contentH := bridge.AvailableContentHeight()
	rowHeight := contentH / leftCount
	if rowHeight < 6 {
		rowHeight = 6 // Clamp to LED diameter minimum (Requirement 7.4)
	}

	originX, originY := bridge.ContentOrigin()
	colWidth := bridge.AvailableContentWidth() / 2
	rightColX := originX + colWidth

	// 6. Resolve accent color for textlabel foreground.
	accent := resolveFGColor(p.FGColor)

	// 7. Create Compositor and render pin cells (Requirement 5.1).
	suppCtx := widgets.SuppressionContext{
		AvailableWidth:  bridge.AvailableContentWidth(),
		AvailableHeight: bridge.AvailableContentHeight(),
	}
	comp := widgets.NewCompositor(suppCtx)

	ledDiameter := 6
	ledTextGap := 4

	for i, pin := range snap.Pins {
		// Determine column and row position.
		var cellX int
		var rowIndex int
		if i < leftCount {
			cellX = originX
			rowIndex = i
		} else {
			cellX = rightColX
			rowIndex = i - leftCount
		}

		cellY := originY + rowIndex*rowHeight
		ledY := cellY + (rowHeight-ledDiameter)/2

		// LED widget for this pin (Requirement 5.2).
		state, fg := ledColorForState(pin)
		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   ledDiameter,
			Bounds:     image.Rectangle{Min: image.Pt(cellX, ledY)},
			Foreground: fg,
		}))

		// Textlabel widget for this pin (Requirement 3.3).
		textBoundsH := rowHeight
		if glyphHeight > textBoundsH {
			textBoundsH = glyphHeight // Glyph clamping (Requirement 3.6)
		}
		// Clamp textBoundsH so it does not exceed the content area bottom edge.
		contentBottom := originY + contentH
		if cellY+textBoundsH > contentBottom {
			textBoundsH = contentBottom - cellY
		}
		comp.Add(textlabel.New(textlabel.Config{
			Text:       fmt.Sprintf("%02d", i),
			Bounds:     image.Rect(cellX+ledDiameter+ledTextGap, cellY, cellX+colWidth, cellY+textBoundsH),
			Font:       selectedFont,
			Alignment:  textlabel.Left,
			Foreground: accent,
		}))
	}

	// 8. Return sprites-only ViewData (Requirement 1.1, 1.2).
	return style.ViewData{Static: true, Sprites: comp.Sprites()}
}

func einkFramed800Build(snap GpioSnapshot, pol Policy, ctx style.StyleContext, d def) style.ViewData {
	hints := ctx.Hints()

	// Step 1: Guard
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 2})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Static: true}
	}

	// Step 2: Borderframe
	frameCfg := borderframe.Config{
		Bounds: image.Rect(0, 0, hints.PixelWidth, hints.PixelHeight),
		Theme:  "sharp",
	}
	frameSprite := borderframe.Render(frameCfg)
	var sprites []widgets.Sprite
	if frameSprite != nil {
		sprites = append(sprites, *frameSprite)
	}

	// Step 3: Tier resolution
	normalEntry := ctx.Entry(tiercatalog.TierNormal)
	rowHeight := normalEntry.RowHeight
	if rowHeight <= 0 {
		rowHeight = 10
	}

	// Step 4: Row assembly — one row per pin
	visiblePins := snap.Pins

	// Step 5: Adaptive fitting — drop trailing pins if they exceed available height
	maxRows := bridge.AvailableContentHeight() / rowHeight
	if maxRows > 0 && len(visiblePins) > maxRows {
		visiblePins = visiblePins[:maxRows]
	}

	items := make([]string, len(visiblePins))
	for i, p := range visiblePins {
		items[i] = p.String()
	}

	// Step 6: Truncation
	ga := normalEntry.GlyphAdvance
	if ga > 0 {
		maxChars := bridge.AvailableContentWidth() / ga
		if maxChars > 0 {
			for i := range items {
				items[i] = textlayout.Truncate(items[i], maxChars)
			}
		}
	}

	// Step 7: Horizontal centering
	offsets := make([]int, len(items))
	for i, item := range items {
		offsets[i] = bridge.CenterXWith(len([]rune(item)), ga)
	}

	// Step 8: Vertical centering
	blockHeight := len(items) * rowHeight
	offsetY := (bridge.AvailableContentHeight() - blockHeight) / 2
	if offsetY < 0 {
		offsetY = 0
	}

	// Step 9: LED sprites via Compositor with SuppressOnEink (LEDs are suppressed
	// because they lack the "eink-safe" capability).
	originX, originY := bridge.ContentOrigin()
	accent := resolveFGColor(snap.Policy.FGColor)

	suppCtx := widgets.SuppressionContext{
		IsEink:          true,
		AvailableWidth:  bridge.AvailableContentWidth(),
		AvailableHeight: bridge.AvailableContentHeight(),
	}
	comp := widgets.NewCompositor(suppCtx, widgets.SuppressOnEink())

	for i, p := range visiblePins {
		state := led.Off
		if p.Level {
			state = led.On
		}
		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   led.DiameterForRow(rowHeight),
			Bounds:     image.Rectangle{Min: image.Pt(originX, originY+offsetY+i*rowHeight)},
			Foreground: accent,
		}))
	}
	sprites = append(sprites, comp.Sprites()...)

	return style.ViewData{
		Items:       items,
		LineOffsets: offsets,
		Colors:      nil,
		Static:      true,
		OffsetY:     offsetY,
		Sprites:     sprites,
	}
}
