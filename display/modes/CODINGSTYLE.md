# Coding Style for Display Modes

This file is the single source of truth for display-mode style conventions.

## Core rule

Use **hand-crafted, explicit styles**. Prefer one file per resolution/capability variant, with per-profile tweaks in code rather than generated output.

## Idiomatic pattern

Most production modes follow this shape:

```go
type Params struct {
	BuildFn func(Snapshot, Policy, StyleContext, def) ViewData
}

type def struct {
	name string
	reqs style.SurfaceRequirements
	p    Params
}

var NameStyle = def{
	name: "kebab-case-name",
	reqs: style.SurfaceRequirements{
		MinWidth:   WIDTH,
		MinHeight:  HEIGHT,
		Capability: style.MonoFast, // or MonoSlow, GrayscaleSlow, GrayscaleFast, ColorSlow, ColorFast
	},
	p: Params{BuildFn: buildName},
}
```

## File layout

- `core.go`: `Params`, `def`, interface plumbing
- `styles.go` or mode registry file: `style.NewRegistry[...]`
- `style_*.go`: one file per explicit style
- `style_helpers.go`: shared helpers only

## Render function shape

```go
func buildName(snap Snapshot, p Policy, ctx style.StyleContext, _ def) style.ViewData {
	hints := ctx.Hints()
	bridge := layout.NewLayoutBridge(hints, layout.BridgeConfig{PaddingPct: 0})
	if bridge.AvailableContentWidth() <= 0 || bridge.AvailableContentHeight() <= 0 {
		return style.ViewData{Items: []string{"(too small)"}, Static: true}
	}

	entry := ctx.Entry(tiercatalog.TierNormal)
	// Compute rows, truncate text, align offsets, fit vertically.
	return style.ViewData{Items: items, Tiers: tiers, LineOffsets: offsets, OffsetY: offsetY}
}
```

## Registry rule

Register styles in capability order. The first entry is the default.

```go
var registry = style.NewRegistry[Snapshot, Policy](
	styles.DefaultStyle,
	styles.MonoSlow128x64Style,
	styles.MonoFast128x64Style,
	// ...
)
```

## Hardware aliases

Register panel defaults in `hardware/panels/*/panel.go`:

```go
style.RegisterAliases("panel-name", map[string]string{
	"clock":   "color-240x240",
	"serial":  "mono-fast-128x64",
	"gpio":    "color-320x240",
})
```

## Capability guidance

- **MonoSlow**: e-paper, set `Static: true`
- **MonoFast**: OLED, no `Static`
- **GrayscaleSlow/Fast**: same def/Params pattern, use grayscale-capable layouts
- **ColorSlow**: preserve color output, e-paper-safe layouting
- **ColorFast**: widgets, richer color, multi-tier layouts when needed

## Exceptions

- **Attract modes** use parameterized structs when one code path must cover many resolutions.
- Use a closure for registry construction only when the registry needs extra setup.

## Validation

- Build the package
- Run the mode tests
- Verify snapshots
- Keep per-style names stable; aliases depend on them

## Rule of thumb

If you need a new display-mode style document, update this file instead.
