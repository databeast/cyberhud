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

## Compositing depth - the text plane is the datum, sprites carry `Z`

The renderer draws the text/cell plane first, at depth 0, then draws sprites in
ascending `widgets.Sprite.Z` order around it. The default `Z: 0` therefore draws a
sprite *on top of* the text - correct for the common case (icons, bars, labels
decorating a row). A sprite with `Z < 0` sinks *beneath* the text.

A full-panel backdrop (gradient, fill, frame fill) is **`Z: -1`**. Leaving it at the
default buries the text under an opaque sprite and the mode renders blank - every
mechanical check still passes, so nothing warns you. If your text has vanished and a
sprite covers the panel, the lever is the backdrop's `Z`, not the text's position.

99% of styles never set `Z`, because their elements do not overlap. Only reach for it
when a sprite and the text (or another sprite) occupy the same pixels and order
matters. Within one depth, declaration order breaks ties.

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

## A single centred label: bind the font, or center it wrong

A bare `Items` row (text, no `Tiers`, no `FontIDs`) opts out of the font contract.
The renderer's `resolveTextFitFonts` fallback then re-picks the *largest font that
fits the width* - usually a bigger face than the `ctx.Entry`/`ctx.Metrics` advance you
measured `CenterX` against. You center for a narrow face; the renderer draws a wide
one; the label leans right and clips (the classic "BOOT COMP" overflow).

Two correct fixes, strongest last:

- **Emit the tier you measured:** set `ViewData.Tiers = []tiercatalog.Tier{tier}`
  alongside the row. `resolveTierFonts` then binds the row to that face and the fit
  fallback early-returns. Minimal, and it is what the `ctx.Metrics` doc prescribes.
- **Rasterize the label as a sprite** with `textlabel.Render` /
  `textlabel.New`, using the face from `ctx.Face(family, tier)`, placed with
  `bridge.ContentOrigin()` + `AvailableContentWidth/Height()` and
  `textlabel.Center`. The sprite owns its glyphs, so measurement and drawing are the
  same face by construction and the renderer never re-resolves. This is the idiomatic
  path when the label also sits over a backdrop (serial's engine, systemd's gradient
  and single-label builders): the backdrop is `Z: -1`, the label sprite is the default
  `Z: 0`, and it reads on top. See `display/modes/serial/styles/engine.go`.

Either way: a style must never measure with one font and let the renderer draw with
another.

## Rendered text is ASCII-only - non-ASCII glyphs render as tofu

The bitmap fonts (spleen and friends) carry ASCII and little else. Any non-ASCII
character you put in an `Items` row or a `textlabel` renders as the font's
missing-glyph box - visible on-panel as a stray `?`. A time window written
`"11:29–12:29"` with an en-dash shows up as `"11:29?12:29"`; the same trap caught a
degree symbol (`52.3?C`) and a `…` ellipsis elsewhere.

Use ASCII equivalents in any string a style emits for drawing:

- `–` / `—` (en/em dash) and `→` (arrow) → `-` / `->`
- `·` (middot) → `-` or `|`
- `…` (ellipsis) → `...` (and budget three columns for it, not one)
- `°`, `×`, curly quotes, etc. → ASCII or drop them

Icons are the exception and go through a *different* path: the material-icons font
is a private-use-area glyph set addressed by codepoint via the icon widgets, not by
typing a Unicode character into a text row. Decorative symbols belong there, not in
`Items`.

## Fit text to the width - never truncate it

A row clipped to `"Design re..."` is a visible failure, not information: the glance
it exists for cannot use it. So a style never truncates a line to force a fit. It
either makes the line fit, or it does not show the line.

The primitive is `layout.Metrics.TextWidth(s)` - the exact horizontal pixels a
string needs in that font, by the real draw rule ((n-1)·advance + last glyph's
ink width, not advance·n). `Metrics.Fits(s, widthPx)` is the predicate. Both are
pure and unit-tested in `display/style/layout` - a fit decision needs no render.

The two structural responses when a line is too wide, in order:

1. **Size the row down.** `ctx.FitTier(text, widthPx, ceiling)` returns the
   largest tier (no bigger than `ceiling`) whose font renders the whole string
   within the width. The text drives the font, capped by the tier the panel
   height suggests. Measure and emit that tier so layout and drawing agree.
2. **Drop the row.** When `FitTier` returns `ok == false`, the string will not fit
   even at the smallest tier - omit the row rather than clip it. `FitRows` already
   drops rows that do not fit *vertically*; this is the horizontal twin.

The one exception is a mandatory headline (e.g. an event title): it must always
show, so when nothing fits it falls back to the smallest tier and hard-clips at
the panel edge - a clean cut, never an ellipsis apology. Decide per row which are
mandatory and which are sheddable; secondary/preview rows are sheddable.

`MaxChars` (width / advance) is a coarse character budget, fine for a quick bound;
it is not the fit test. For "does this exact string fit?", use `TextWidth`/`Fits`.

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
- Measured tier == emitted tier (or the label is a face-bound sprite). No bare
  centred `Items` row relying on the renderer's font-fit fallback.
- Any full-panel backdrop sprite is `Z: -1` so the text reads on top.
- No non-ASCII characters in any string a style emits for drawing (dashes,
  arrows, middots, ellipses, degree signs) - the bitmap fonts render them as
  tofu. Use ASCII; icons go through the icon widgets, not text rows.
- No truncation-to-fit. Every emitted row fits its width whole (size it down via
  `ctx.FitTier`, or drop it when `ok == false`). A `...`-clipped row is a bug, not
  a layout. Mandatory headlines may hard-clip at the edge as the sole exception.
- `PaddingPct` set whenever padding is non-zero.
- Run `go build ./...` and the mode's tests (snapshot tests catch positional drift).
- To find usages and blast radius, use graphify (`$G query`, `$G affected`) before grepping.
