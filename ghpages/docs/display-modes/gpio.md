# GPIO

The GPIO mode displays the current state of all configured GPIO pins, rendering a live-updating read-only view with configurable visual styles, color coding, and adaptive font selection. It provides a clear, at-a-glance pin state summary suitable for any panel, with multiple display styles ranging from compact grids to detailed per-pin breakdowns.

## Quick Start

```sh
cyberhudctl display set 0 gpio
```

## How It Works

The GPIO mode reads the state of each configured GPIO pin from the hardware GPIO manager and renders a summary on the panel surface, with each pin entry showing the BCM GPIO number, pin mode (input or output), and current logic level (high or low). The display redraws on the standard 1-second render-loop tick and only flushes to the panel when the combined pin-state fingerprint changes, preventing unnecessary redraws when all pins are stable.

Data is sourced from the kernel GPIO subsystem via the hardware GPIO manager — only pins explicitly declared in your configuration file appear in this view. When no GPIO pins are configured or the GPIO subsystem is unavailable, the mode displays a single "GPIO unavailable" message row. Color coding differentiates states visually — green indicates a high level, grey indicates low — making it easy to spot active pins at a glance. On monochrome panels the color coding is suppressed and level is shown by text only.

The mode provides five style variants that produce substantially different visual outputs: **list** renders one pin per row in a scrollable text list with LED indicators beside output pins; **icons** shows a compact grid of LED circles arranged in row-major order with no text labels; **detail** adds user-configured labels to each row for pin identification; **dashboard** displays aggregate counters (total, output, high) in a header followed by grouped LED grids; and **activity** renders a per-pin sparkline graph of recent toggle history for output pins.

!!! note
    The pins shown depend on your GPIO pin configuration. Only pins declared in your config file appear in this view. See [Pin Assignments](../reference/pin-assignments.md) for details on configuring pins.

## Styles

The `style` option controls the visual layout and information density:

- **list** (default) — a scrollable text list showing one pin per row with BCM number, mode, and level. Each row is color-coded by pin level. LED indicators appear beside output pins. Supports adaptive font selection.
- **icons** — a compact grid of LED indicators arranged in row-major order. Each pin is represented as a filled circle (high) or outlined circle (low). No text labels — best for monitoring many pins at once on smaller panels.
- **detail** — an expanded view showing each pin as a row with an LED indicator, BCM number, mode, level, and an optional user-configured label. Use this when you need to identify pins by purpose.
- **dashboard** — a summary view with a header row showing aggregate counts (total pins, output count, high count), followed by LED grids grouped by mode — output pins first, then input pins below. Non-scrolling.
- **activity** — a sparkline graph view showing recent toggle history for each output pin. Each output pin gets a line graph of the last 32 state-change events, making it easy to visualize GPIO activity over time.

Each style automatically selects the best font for the panel dimensions, so text remains readable regardless of resolution.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| `style` | string | Visual layout for pin state display controlling spatial arrangement of GPIO data | list | list, icons, detail, dashboard, activity |
| `color` | bool | Whether to use color-coded status indicators for pin high/low states on color panels | true | true, false |
| `font` | string | Font selection for text rendering, auto selects the best fit for panel dimensions | auto | any valid string (auto or a registered font ID) |
| `fgcolor` | string | Foreground color for GPIO display elements on color-capable panels | cyan | cyan, green, amber, red, white, none |

Configure options via the CLI:

```sh
cyberhudctl display gpio <key>=<value> [<key>=<value> ...]
```

## CLI Examples

Set the display style to icons:

```sh
cyberhudctl display gpio style=icons
```

Switch to the detail view with a specific font:

```sh
cyberhudctl display gpio style=detail font=ascii-5x7
```

Use the dashboard summary:

```sh
cyberhudctl display gpio style=dashboard
```

Watch activity sparklines:

```sh
cyberhudctl display gpio style=activity
```

Disable color coding (useful for monochrome panels):

```sh
cyberhudctl display gpio color=false
```

Combine multiple options in one command:

```sh
cyberhudctl display gpio style=list color=true font=auto
```

Query all current settings:

```sh
cyberhudctl display gpio
```

This returns all key=value pairs reflecting the active configuration.

## Font Selection

The `font` option controls which bitmap font is used for text rendering:

- **auto** (default) — the mode picks the largest font that fits the panel height while maintaining at least 4 visible rows of pins. This adapts automatically to different panel resolutions.
- **Specific font ID** — any registered font identifier (e.g., `ascii-4x6`, `ascii-5x7`). If the ID isn't recognized, the mode silently falls back to automatic selection.

Font selection affects the list, detail, dashboard, and activity styles. The icons style uses LED indicators exclusively and is not affected by font choice.

## Detail Style Labels

The detail style supports user-configured labels for each pin. Labels are stored in the pin-labels policy map, keyed by BCM pin number. When a label is configured, it appears at the end of the pin's row (e.g., "4 OUT HI relay-1"). When no label is set, the label segment is omitted.

Labels are automatically truncated if they exceed the available character width after the fixed columns (pin number, mode, level).

## Activity Tracking

The activity style maintains a sliding window of the last 32 render ticks per output pin. Each tick records whether the pin toggled (1.0) or stayed the same (0.0) since the previous tick. This produces a sparkline graph that reveals patterns — rapid toggling shows as a dense waveform, while idle pins show a flat baseline.

Only output pins are tracked. Input pins and pins with no recorded history display a flat zero-line.

## Layout Behavior

The GPIO mode optimizes rendering automatically:

- **Row omission** — if the pin list exceeds the panel height, it is truncated to the maximum visible rows. The list and detail styles are scrollable; icons and dashboard are not.
- **Adaptive fonts** — the best available font for the panel dimensions is selected on each render, responding to configuration changes.
- **Graceful degradation** — if the panel is too small to render any content (e.g., width of zero pixels), the mode falls back to text-only output without sprites.

## Panel Compatibility

The GPIO mode is non-interactive (read-only) and works on all panels regardless of input controls or resolution. It does not require buttons or joystick input — it only observes configured GPIO states. On monochrome panels, pin state indicators use on/off pixel patterns instead of color-coded values. On slow-refresh panels, the display updates only when a pin state changes rather than on a continuous timer.

| Capability | Description | Behavior |
|---|---|---|
| MonoFast | Fast-refresh monochrome OLED (128×32, 128×64, 128×128) | Fully supported — pin states shown as text labels with on/off indicators in white |
| MonoSlow | Slow-refresh monochrome e-ink (122×250, 176×264, 200×200, 296×128, 400×300, 480×800, 800×480) | Fully supported — static pin state view, refreshes only on state change |
| GrayscaleFast | Fast-refresh grayscale (160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — grayscale intensity indicators for pin states |
| GrayscaleSlow | Slow-refresh grayscale e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static grayscale pin state view, refreshes on state change |
| ColorFast | Fast-refresh color TFT (128×128, 160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — color-coded pin states (green=high, red=low) with real-time updates |
| ColorSlow | Slow-refresh color e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static color pin state view, refreshes on state change |

For interactive pin control (toggling output pins), see [GPIO Control](gpio_control.md).

!!! tip
    Use `style=icons` or `style=dashboard` on small panels for maximum density. Use `style=detail` on larger panels when you need labeled pin identification.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
- [GPIO Control](gpio_control.md) — interactive mode for toggling output pins
- [Pin Assignments](../reference/pin-assignments.md) — configure which pins appear in this view
- [Dashboard](dashboard.md) — consolidated status view including GPIO state summary


