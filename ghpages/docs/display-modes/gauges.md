# Gauges

The gauges display mode renders user-supplied numeric data as progress bars, rings, arcs, and dial-style widgets. It is designed for personal dashboards and project-specific status displays where you want to provide your own values instead of reading a fixed system source.

## Quick Start

```sh
cyberhudctl display set 0 gauges
cyberhudctl display gauges set 42
```

## How It Works

The gauges mode renders one or more progress-bar widgets arranged in an auto-computed grid (rows × columns) sized to fill the panel surface, with each tile containing a filled indicator and an optional text label above it. Data is supplied externally via the `cyberhudctl display gauges set` command — you push a single number, a JSON object, or a JSON array of gauge objects — and the display updates event-driven on the next render tick (1 second) after data arrives.

Each gauge value is normalized to a 0–1 percentage within its declared min/max range and rendered using the configured shape: **linear** (a horizontal progress bar filling left to right), **ring** (a circular outline that fills clockwise), **arc** (a semicircular sweep), or **pie** (a filled wedge). When `shape` is set to **auto** (the default), the mode selects `linear` unless the payload overrides it per gauge.

When no data has been pushed yet (empty gauge set), the mode displays a placeholder "gauges idle" text row rather than rendering empty widgets. The grid layout auto-sizes based on panel pixel dimensions — columns and rows are computed to maximize tile area — but you can override both explicitly. On color panels, gauge fills use the configured accent color; on monochrome panels, fills render in white.

The mode is well suited to:

- a home server panel showing CPU, RAM, disk, and network usage
- a 3D printer display showing print progress, nozzle temperature, and chamber heat
- a solar or battery monitor showing charge level, inverter load, and voltage
- a plant or aquarium dashboard showing moisture, water level, and pump duty cycle

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| `style` | string | Visual style name for gauge rendering or empty for automatic fitness-based selection | gauges-default | any valid string (style names or empty for auto) |
| `shape` | string | Default gauge shape used when the payload does not override it per individual gauge | auto | auto, linear, ring, arc, pie |
| `show_labels` | bool | Render labels above each gauge tile showing the gauge name text for identification | true | true, false |
| `label_tier` | string | Font tier used for gauge labels controlling text size relative to gauge tile height | normal | auto, small, normal, large, fullsize |
| `accent` | string | Named accent color for gauge fill indicators providing visual theming on color panels | cyan | amber, cyan, emerald, green, red, white, none |
| `default_min` | float | Default minimum value used when a gauge payload omits the min field boundary specification | 0 | any valid float |
| `default_max` | float | Default maximum value used when a gauge payload omits the max field boundary specification | 100 | any valid float |
| `columns` | int | Explicit gauge grid columns for layout, or 0 for automatic column count based on panel width | 0 | any valid int (0 for auto) |
| `rows` | int | Explicit gauge grid rows for layout, or 0 for automatic row count based on panel height | 0 | any valid int (0 for auto) |
| `tile_gap_px` | int | Gap between gauge tiles in pixels controlling visual spacing between adjacent gauge widgets | 1 | any valid int |
| `padding_pct` | int | Layout inset percentage applied to the panel reducing the drawable area from panel edges | 0 | any valid int (0–50) |

## Policy Fields

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `style` | string | gauges-default | Visual style name, or empty for automatic selection |
| `shape` | string | auto | Default gauge shape (`auto`, `linear`, `ring`, `arc`, `pie`) |
| `show_labels` | bool | true | Render a label above each gauge tile |
| `label_tier` | string | normal | Font tier used for labels |
| `accent` | string | cyan | Accent color for fills |
| `default_min` | float | 0 | Default minimum when a payload omits `min` |
| `default_max` | float | 100 | Default maximum when a payload omits `max` |
| `columns` | int | 0 | Explicit grid columns, or 0 for auto |
| `rows` | int | 0 | Explicit grid rows, or 0 for auto |
| `tile_gap_px` | int | 1 | Pixel gap between gauge tiles |
| `padding_pct` | int | 0 | Layout inset percentage |

## CLI Examples

Show a single gauge from a plain number:

```sh
cyberhudctl display gauges set 72
```

Push a home server status payload:

```sh
cyberhudctl display gauges set '{"gauges":[{"label":"CPU","value":42,"min":0,"max":100},{"label":"RAM","value":73,"min":0,"max":100},{"label":"Disk","value":81,"min":0,"max":100}]}'
```

Render a battery monitor with rings instead of bars:

```sh
cyberhudctl display gauges policy shape=ring show_labels=true accent=green
cyberhudctl display gauges set '{"gauges":[{"label":"Battery","value":64,"min":0,"max":100},{"label":"Load","value":18,"min":0,"max":100}]}'
```

Use a compact 2-column layout for a small TFT panel:

```sh
cyberhudctl display gauges policy columns=2 tile_gap_px=2 padding_pct=2
```

Hide labels for a very small panel:

```sh
cyberhudctl display gauges policy show_labels=false
```

Drive a 3D printer dashboard from your own script:

```sh
cyberhudctl display gauges set '{"gauges":[{"label":"Progress","value":18,"min":0,"max":100},{"label":"Nozzle","value":204,"min":0,"max":260},{"label":"Bed","value":58,"min":0,"max":110}]}'
```

## Payload Examples

Single-value payload:

```json
42
```

Object payload:

```json
{"label":"CPU","value":42,"min":0,"max":100}
```

Multi-gauge payload:

```json
{
  "gauges": [
    {"label":"CPU","value":42,"min":0,"max":100},
    {"label":"RAM","value":73,"min":0,"max":100},
    {"label":"Disk","value":81,"min":0,"max":100}
  ]
}
```

## Practical Notes

- Values outside the configured range are clamped before rendering.
- Empty or malformed payloads fall back to a compact placeholder instead of crashing.
- The first registered style is the default; specific panels can still resolve a better fit automatically.

## Panel Compatibility

The gauges mode is non-interactive and renders on all six panel capability classes without requiring buttons or any minimum resolution. The mode uses a multi-style registry that automatically selects a style variant tuned for each panel type.

| Capability | Description | Behavior |
|---|---|---|
| MonoFast | Fast-refresh monochrome OLED | Fully supported — normal label tier, standard tile gap rendering |
| MonoSlow | Slow-refresh monochrome e-ink | Fully supported — labels use small tier to conserve vertical space on slow-refresh displays |
| GrayscaleFast | Fast-refresh grayscale | Fully supported — normal label tier, full grayscale fill rendering |
| GrayscaleSlow | Slow-refresh grayscale e-ink | Fully supported — labels use small tier, static output avoids unnecessary redraws |
| ColorFast | Fast-refresh color TFT | Fully supported — accent color applied to gauge fills, normal label tier |
| ColorSlow | Slow-refresh color e-ink | Fully supported — accent color applied, labels use small tier for compact layout |

On slow-refresh panels (MonoSlow, GrayscaleSlow, ColorSlow) the label font tier defaults to small to keep gauge tiles compact and reduce vertical overflow. The accent color option only produces visible color differences on color-capable panels (ColorFast, ColorSlow); monochrome and grayscale panels render fills in their native foreground color.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
