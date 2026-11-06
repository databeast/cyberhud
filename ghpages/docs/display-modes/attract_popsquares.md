# Attract Pop Squares

The pop-squares mode renders a full-bleed grid of flat color cells that rest at a calm base color while a sparse scatter of cells gently pulses brighter and eases back. Inspired by the XScreenSaver "popsquares" hack, it evokes a supercomputer idling: banks of status cells, something always quietly ticking over somewhere, no urgency and no chaos. The base field is a slow drifting gradient so neighboring cells stay related — one cohesive system rather than random confetti — and each cell's pulse is driven by a per-cell hash with an irrational-period cycle, so the twinkle never falls into a visible repeating pattern.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| `speed` | float | Animation rate controlling the pop-and-fade cadence | 1.0 | any valid float (0.1–10.0) |
| `grid_scale` | float | Cell size multiplier; smaller values pack more, finer squares | 1.0 | any valid float (0.5–4.0) |
| `pop_rate` | float | Fraction of cells gently active at once; low values read as idle | 0.35 | any valid float (0.1–1.0) |
| `warmth` | float | Palette temperature: 0 cool/idle (blue-teal), 1 warm/busy (amber) | 0.15 | any valid float (0.0–1.0) |

## Policy Fields

| Field | Type | Range | Default | Description |
|-------|------|-------|---------|-------------|
| Speed | float64 | [0.1, 10.0] | 1.0 | Pop-and-fade cadence |
| GridScale | float64 | [0.5, 4.0] | 1.0 | Cell size multiplier |
| PopRate | float64 | [0.1, 1.0] | 0.35 | Fraction of cells active at once |
| Warmth | float64 | [0.0, 1.0] | 0.15 | Palette temperature (cool→warm) |

## CLI Examples

Query current pop-squares policy settings:

```sh
cyberhudctl display attract_popsquares
```

Pack the grid tighter and warm up the palette slightly:

```sh
cyberhudctl display config 0 grid_scale=1.5 warmth=0.3
```

## Panel Compatibility

- **Color TFT** — full RGB rendering with a cool nominal palette and bright pulse highlights. This is the ideal panel type for the pop-squares effect.
- **Monochrome OLED** — cells are rendered in grayscale, mapping base color and pulse to white levels. The idle twinkle remains visible without color variation.
- **E-Ink** — renders a single static calm frame due to slow refresh characteristics. See E-Ink Behavior below.

### E-Ink Behavior

On e-ink panels, the pop-squares mode renders a single frozen frame — the calm base field with a handful of cells caught mid-pulse — and holds it. The RenderCacheKey remains stable after the initial frame, preventing unnecessary refreshes that would cause flicker or ghosting on e-ink displays.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Attract Modes](attract-modes.md) — overview of all attract/screensaver modes

<!-- snapshot-gallery:start -->
<!-- snapshot-gallery:end -->
