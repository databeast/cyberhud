# Storage

The storage mode shows how full the host's filesystems are, at a glance. Each mounted volume is a row with its mountpoint, a compact used/total summary, a usage percentage, and a progress bar whose fill and color encode disk pressure — green when there's room, amber as it fills, red when a volume is nearly full. It works on any panel size, collapsing to a single-volume focus readout on the smallest displays.

## Quick Start

```sh
cyberhudctl display set 0 storage
```

## How It Works

On Linux the storage mode enumerates the mounted filesystems listed in `/proc/mounts` and queries each one with `statfs(2)`, reporting used, total, and free bytes. Used and total match the accounting that `df` reports: reserved (root-only) blocks count as used, and free space is what is available to an unprivileged user. Sampling happens on a background goroutine at a configurable interval (default: every 5 seconds). Because disk usage changes slowly, the default cadence is relaxed compared with CPU or thermal monitoring.

Pseudo and virtual filesystems — `tmpfs`, `overlay`, `cgroup`, `proc`, `sysfs`, and similar — are hidden by default so the view shows real block storage. Set `show_all=true` to include them.

Mounts are ordered most-used-first, with the root filesystem always promoted to the top. On a short panel where not every row fits, this keeps the fullest (most urgent) volumes visible and drops the least-used ones. A green→amber→red heat ramp colors each row and its bar on color-capable panels, so a nearly-full disk reads as alarming red while a roomy one stays green. On monochrome and grayscale panels, color cannot carry that meaning, so the bar's fill level conveys the magnitude instead.

On non-Linux hosts the mode reports no filesystems and renders a placeholder.

## Styles

Each style targets a specific panel form factor (resolution and display capability) and comes in two variants, plus an adaptive skeleton the daemon auto-selects when no style is configured:

- **`<capability>-<WxH>` (skeleton)** — adaptive layout that picks bars or focus based on panel size. This is what fitness scoring selects for an unconfigured panel.
- **`<capability>-<WxH>-bars`** — the multi-mount view: one row per filesystem with a usage-colored progress bar filling the trailing width.
- **`<capability>-<WxH>-focus`** — a single-filesystem readout of the most-used volume (mountpoint, percentage, free space, full-width bar), intended for space-constrained panels.

Style names follow the pattern `<capability-prefix>-<width>x<height>[-<variant>]`, for example `color-fast-240x240-bars` or `mono-fast-128x64-focus`. When no explicit style is configured, the daemon auto-selects the best style for the detected panel using fitness scoring: each registered style declares the resolution and capability it supports, and the highest-scoring style for the current panel wins, with registration-order tie-breaking. A bare resolution requirement means a larger panel in the same capability class still resolves sensibly.

Each style automatically selects the best bitmap font for the available panel space.

<!-- snapshot-gallery:start -->
<!-- snapshot-gallery:end -->

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| `style` | string | Visual presentation style (resolution-qualified name, or empty for auto-selection) | (auto-selected) | any registered style name (e.g. `color-fast-240x240-bars`, `mono-fast-128x64-focus`) |
| `font` | string | Font selection (auto or a registered font ID) | auto | any valid string (auto or a registered font ID) |
| `refresh_ms` | int | Sampling interval in milliseconds controlling how often filesystems are polled | 5000 | 500–120000 |
| `unit` | string | Capacity unit style for byte values | auto | auto, si, iec |
| `show_all` | bool | Include pseudo/virtual filesystems (tmpfs, overlay, ...) | false | true, false |

Configure options via the CLI:

```sh
cyberhudctl display storage <key>=<value> [<key>=<value> ...]
```

## CLI Examples

Activate the storage display mode on the main panel:

```sh
cyberhudctl display set 0 storage
```

Select the multi-mount bars layout for a 240×240 color panel:

```sh
cyberhudctl display storage style=color-fast-240x240-bars
```

Use the single-filesystem focus layout on a 128×64 OLED:

```sh
cyberhudctl display storage style=mono-fast-128x64-focus
```

Show IEC (1024-based) units with `KiB`/`MiB`/`GiB` suffixes:

```sh
cyberhudctl display storage unit=iec
```

Include tmpfs and other pseudo filesystems:

```sh
cyberhudctl display storage show_all=true
```

Set a faster 2-second refresh rate:

```sh
cyberhudctl display storage refresh_ms=2000
```

Query all current settings (zero arguments returns the active configuration):

```sh
cyberhudctl display storage
```

This returns all key=value pairs reflecting the active configuration, for example:

```
OK storage style= font=auto refresh_ms=5000 unit=auto show_all=false
```

## Capacity Units

The `unit` option controls how byte counts are formatted:

| Value | Base | Suffixes | Example |
|-------|------|----------|---------|
| `auto` | 1024 | B, K, M, G, T | `54G`, `931G` |
| `si` | 1000 | B, KB, MB, GB, TB | `58GB` |
| `iec` | 1024 | B, KiB, MiB, GiB, TiB | `54GiB` |

Values below 10 of a given magnitude show one decimal place (`4.0G`); larger values are rounded to whole units to keep each row narrow on small panels.

## Panel Compatibility

The storage mode works on all panels regardless of input controls. It does not require buttons, joystick, or any specific resolution — the selected style adapts its rendering to fit the available display area. On monochrome and grayscale panels, usage is conveyed by bar fill level rather than color. On slow-refresh (e-ink) panels, the view is a static frame that updates at the configured refresh interval rather than animating.

| Capability | Description | Behavior |
|---|---|---|
| MonoFast | Fast-refresh monochrome OLED (128×32, 128×64, 128×128) | Fully supported — mountpoint, usage summary, and bar fill in white, periodic refresh |
| MonoSlow | Slow-refresh monochrome e-ink (122×250, 128×296, 200×200, 240×240, 480×320, 800×480) | Fully supported — static usage frame, updates at configured refresh_ms interval |
| GrayscaleFast | Fast-refresh grayscale (128×128, 160×128, 240×240, 320×240, 480×320, 800×480) | Fully supported — grayscale usage rendering with fill-level bars |
| GrayscaleSlow | Slow-refresh grayscale e-ink (122×250, 128×296, 200×200, 240×240, 400×300, 800×480) | Fully supported — static grayscale usage frame, updates at refresh interval |
| ColorFast | Fast-refresh color TFT (128×128, 135×240, 160×128, 240×240, 240×320, 320×240, 480×320, 800×480) | Fully supported — green→amber→red usage heat ramp with colored bars |
| ColorSlow | Slow-refresh color e-ink (122×250, 128×296, 200×200, 240×240, 400×300, 800×480) | Fully supported — static color usage frame, updates at refresh interval |

On panels with joystick or button inputs, physical controls provide quick access:

- **Left / Right** — cycle through registered styles (wrapping around)
- **Secondary** — navigate to the menu

!!! tip
    Use a higher `refresh_ms` on battery-powered or e-ink setups to reduce sampling frequency; disk usage rarely changes second-to-second.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [System](system.md) — host, uptime, and network information
- [Thermal](thermal.md) — per-zone temperature monitoring
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
