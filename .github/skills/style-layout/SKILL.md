---
name: style-layout
description: The idiomatic way to lay out a cyberhud display-mode style - bridge construction, font metrics, row solving, and text-cell to pixel placement. Use when writing or editing a Build function under display/modes/*/styles/, positioning sprites or widgets beside text rows, or touching coordinate math (row/column vs pixels, OffsetY, ContentOrigin).
---

# Style Layout

A style owns layout; the renderer owns drawing. `style.ViewData` (display/style/style.go)
is the whole interface between them. Follow this recipe instead of hand-rolling
`ox + i*rowHeight` arithmetic. The reference implementation is
`display/modes/clock/styles/core.go` (`adaptiveBuild`).

## Two coordinate systems - never mix them silently

| Frame | Type | Origin | Used for |
|---|---|---|---|
| TTY block | `layout.TTYX` / `layout.TTYY` | the centred, padded text block | `ViewData.OffsetY`, `LineOffsets`, `CenterX`, `CenterBlockY` |
| World (panel pixels) | plain `int` / `image.Point` | panel top-left | `widgets.Sprite.Position`, widget `Bounds` |

TTY values are distinct types so they cannot be assigned into a world field.
Cross the boundary deliberately with `bridge.TTYToPixel(x, y)`.

## The recipe

1. **Bridge:** `bridge := ctx.Layout(pct)`. Keep `pct` in one named constant
   (see `clockPaddingPct`) and reuse it for step 5. Do not call
   `layout.NewLayoutBridge` with a different padding for the same view.
2. **Metrics:** measure with `ctx.Metrics(tier)`, or `ctx.DefaultMetrics()` for bare
   rows with no tier. Never use raw `hints.RowHeight` / `hints.GlyphAdvance` or the
   `textlayout` defaults - they belong to no real font. Emit the same tier you measured.
3. **Solve rows once:** `result := bridge.LayoutRows([]layout.RowInput{{TextLen, Font}})`.
   Do not compute `(availH - blockH) / 2` or per-row centring by hand;
   `FitRows`/`CenterBlockY`/`CenterX` already do it. Uniform single-font blocks may use
   `CenterBlockY` directly.
4. **Place row-tied sprites with cell helpers:**
   - `bridge.CellToPixel(row, col, m)` - world pixel of a character cell's top-left.
   - `bridge.CellRect(row, col, m)` - the cell's world rectangle.
   - `bridge.PixelToCell(p, m)` - inverse; `ok` is false outside the content area.
   These are relative to the content origin and do NOT include `OffsetY`. When the
   block is centred, add it: `bridge.CellToPixel(i, col, m).Add(image.Pt(0, int(vd.OffsetY)))`.
   Prefer these over `ContentOrigin()` plus manual multiplication.
5. **Return the whole solution in `ViewData`:** `Items`, `Tiers`, `LineOffsets`,
   `OffsetY`, `RowHeights`, `Spacing`, `VisibleCount`, and `PaddingPct` (the same
   value passed to `ctx.Layout`). Omitting a field makes the renderer re-derive it and
   disagree with `OffsetY`. Call `style.GuardEmptyItems(&vd)` where the clock does.
6. **Sprite builders take the bridge:** helpers that emit row-aligned sprites accept
   `(bridge layout.LayoutCalculator, m layout.Metrics)`, not `hints` or a bare
   `rowHeight`. Examples: `gpio_control/styles.BuildSprites`, `stemma/styles.BuildSprites`.

## When NOT to use the cell helpers

Pixel-space composition is legitimate and stays as is: centred icons or gauges, pixel
bands, tile grids, and rows with a stretched pitch (for example `gpio` row pitch that
grows to fill the panel). Use `bridge.ContentOrigin()` and `AvailableContentWidth()` /
`AvailableContentHeight()` as the base there. The cell helpers are for rows laid out
at a font's own advance and row height.

## Checklist before finishing

- No `ox + n*rowHeight`, `oy + i*rowHeight`, or `col*glyphAdvance` left that a cell
  helper could express.
- No TTY value assigned to a world field without `TTYToPixel`.
- Measured tier == emitted tier.
- `PaddingPct` set whenever padding is non-zero.
- Run `go build ./...` and the mode's tests (snapshot tests catch positional drift).
- To find usages and blast radius, use graphify (`$G query`, `$G affected`) before grepping.
