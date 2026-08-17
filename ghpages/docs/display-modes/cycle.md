# Cycle

Auto-cycles through configured display modes on one or more regions at a configurable interval. Cycle is a meta-mode — rather than rendering its own content, it runs a background sequencer that advances other modes on a timer. This lets you create rotating displays that automatically showcase multiple modes without manual intervention.

## Quick Start

```sh
cyberhudctl display set 0 cycle
```

With default settings, cycle advances to the next mode on the activating region every 30 seconds, skipping itself to avoid recursion.

## How It Works

When activated, the cycle mode starts a background goroutine that fires on the configured interval. Each tick:

1. Re-reads the current policy (supporting hot-reload of interval and mode list changes).
2. Determines which regions to cycle — either the configured `regions` list or the region that activated the mode.
3. Advances each region to the next mode, skipping "cycle" itself and any mode not in the configured `modes` list (when specified).

The sequencer uses `coordinator.Next()` to advance modes, which wraps around from the last mode to the first. If a configured mode list is provided, only modes in that list are cycled to — others are skipped during advancement.

The mode renders a static placeholder view ("Cycling modes...") since its primary function is orchestrating other modes rather than producing display content.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| interval | duration | Time between automatic mode switches | 30s | 5s to 10m (any valid Go duration) |
| modes | string-list | Ordered list of mode IDs to cycle through (empty = all available on region) | (empty) | comma-separated mode IDs |
| regions | int-list | Region indices to cycle on (empty = activating region only) | (empty) | comma-separated non-negative integers (max 16) |

Configure options via the CLI:

```sh
cyberhudctl display cycle <key>=<value> [<key>=<value> ...]
```

## CLI Examples

Query all current settings:

```sh
cyberhudctl display cycle
```

Set the cycling interval to 1 minute:

```sh
cyberhudctl display cycle interval=1m
```

Cycle only through clock, dashboard, and system modes:

```sh
cyberhudctl display cycle modes=clock,dashboard,system
```

Cycle across multiple regions (main and left-aux):

```sh
cyberhudctl display cycle regions=0,1
```

Combine interval with a mode list:

```sh
cyberhudctl display cycle interval=45s modes=clock,thermal,wifi
```

Clear the mode filter to cycle all available modes:

```sh
cyberhudctl display cycle modes=
```

Switch to cycle mode on a specific region:

```sh
cyberhudctl display set 0 cycle
```

Switch to cycle mode with inline policy:

```sh
cyberhudctl display set 0 cycle interval=20s modes=clock,system
```

## Interval Normalization

The `interval` value is clamped to the range [5s, 10m]:

- Values below 5 seconds are raised to 5s to prevent display thrashing.
- Values above 10 minutes are capped at 10m.
- Zero or negative values default to 30s.

The interval supports any valid Go duration string (e.g., `10s`, `1m30s`, `2m`).

## Mode List Behavior

When `modes` is empty (the default), the sequencer cycles through all modes available on the region, skipping only "cycle" itself.

When `modes` contains specific mode IDs, only those modes are cycled to. If a mode in the list is not available on the region, it is skipped during advancement. Duplicate entries are deduplicated (first occurrence preserved).

## Region Targeting

When `regions` is empty, cycling applies only to the region that activated the cycle mode.

When `regions` contains specific indices (e.g., `0,1,2`), all listed regions are advanced on each tick. This allows a single cycle instance to orchestrate multiple displays simultaneously. The list is capped at 16 entries and deduplicated.

## Panel Compatibility

Cycle mode works on all panels. It produces a static placeholder display and does not require input controls. Since cycle is a meta-mode, the actual display content comes from the modes being cycled to — the cycle mode itself only renders a brief transition indicator before switching.

| Capability | Description | Behavior |
|---|---|---|
| MonoFast | Fast-refresh monochrome OLED (128×32, 128×64, 128×128) | Fully supported — placeholder text rendered in native foreground, transitions occur at configured interval |
| MonoSlow | Slow-refresh monochrome e-ink (122×250, 176×264, 200×200, 296×128, 400×300, 480×800, 800×480) | Fully supported — static placeholder frame, mode switches occur without animation |
| GrayscaleFast | Fast-refresh grayscale (160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — grayscale placeholder rendering with timed transitions |
| GrayscaleSlow | Slow-refresh grayscale e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static placeholder, transitions occur without animation |
| ColorFast | Fast-refresh color TFT (128×128, 160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — color placeholder text, smooth timed transitions |
| ColorSlow | Slow-refresh color e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static placeholder, transitions occur without animation |

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
- [Clock](clock.md) — commonly included in cycle mode lists
- [Dashboard](dashboard.md) — commonly included in cycle mode lists
- [System](system.md) — commonly included in cycle mode lists
- [Thermal](thermal.md) — commonly included in cycle mode lists
- [WiFi](wifi.md) — commonly included in cycle mode lists
<!-- snapshot-gallery:start -->
## Snapshots

### Mono

<figure>
  <img src="../img/cycle/mono-128x64-placeholder_0001.png" alt="mono-128x64-placeholder 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<!-- snapshot-gallery:end -->
