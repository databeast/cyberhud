# Clock

The clock mode displays the current time, date, and weekday, rendering a live-updating clock face with configurable format, timezone, and visual style. It provides a clear, at-a-glance time readout suitable for any panel, with multiple visual styles and layout options.

## Quick Start

```sh
cyberhudctl display set 0 clock
```

## How It Works

The clock mode reads the system time (or a configured timezone) and renders a formatted time display on the panel surface. The display updates every second (when seconds or blink-colon are enabled) or every minute (when both are disabled), keeping the clock face current without unnecessary redraws.

The mode renders up to three rows of information — time, date, and weekday — with automatic font selection and layout optimization for the available panel space. Rows that don't fit are automatically omitted (weekday first, then date), ensuring the time is always visible.

## Styles

The daemon automatically selects a style based on your panel's resolution and capability. You don't normally need to set this manually. Use `cyberhudctl display clock` to see what style is active.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| style | string | Resolution-specific visual layout (auto-selected based on panel) | (auto) | Resolution-specific names (e.g., color-240x240, mono-128x64) |
| show_seconds | bool | Whether to display seconds | true | true, false |
| time_format | string | 12-hour or 24-hour time | 24h | 24h, 12h |
| date_format | string | Date layout or hidden | YYYY-MM-DD | YYYY-MM-DD, DD-MM-YYYY, MM-DD-YYYY, none |
| timezone | string | IANA timezone or system local | local | Any valid IANA timezone (e.g., America/New_York, Europe/London, UTC) or "local" |
| show_weekday | bool | Show the weekday name row | true | true, false |
| blink_colon | bool | Animate colon separator on/off each second | false | true, false |
| fgcolor | string | Foreground color for time text on color panels | cyan | cyan, green, amber, red, white, none |
| show_led | bool | Show LED seconds indicator when seconds digits are hidden | true | true, false |
| seconds_bar | string | Progress bar style showing seconds within current minute | none | none, horizontal, pie |
| show_daybar | bool | Show sparkline bar indicating day progress | false | true, false |
| show_border | bool | Show rounded decorative border frame around the panel edge | false | true, false |
| border_color | string | Border frame color on color panels, or auto to inherit the active fgcolor | auto | cyan, green, emerald, amber, red, white, none, auto |

Configure options via the CLI:

```sh
cyberhudctl display clock <key>=<value> [<key>=<value> ...]
```

## CLI Examples

Set 12-hour time with green foreground:

```sh
cyberhudctl display clock time_format=12h fgcolor=green
```

Show a European date format with no seconds:

```sh
cyberhudctl display clock date_format=DD-MM-YYYY show_seconds=false
```

Display a different timezone:

```sh
cyberhudctl display clock timezone=America/New_York
```

Time-only display — hide date and weekday for maximum digit size:

```sh
cyberhudctl display clock date_format=none show_weekday=false
```

Enable the blinking colon animation:

```sh
cyberhudctl display clock blink_colon=true
```

Set a green foreground color with 12-hour time:

```sh
cyberhudctl display clock fgcolor=green time_format=12h
```

Disable the LED indicator and enable the blinking colon instead:

```sh
cyberhudctl display clock show_led=false blink_colon=true
```

Show a horizontal seconds progress bar with no seconds digits:

```sh
cyberhudctl display clock seconds_bar=horizontal show_seconds=false
```

Enable the day-progress sparkline:

```sh
cyberhudctl display clock show_daybar=true
```

Query all current settings:

```sh
cyberhudctl display clock
```

This returns all key=value pairs reflecting the active configuration.

## Time Format

The `time_format` option controls how hours are displayed:

- **24h** — hours 00–23 with leading zero, no AM/PM indicator (e.g., `14:30:05`)
- **12h** — hours 1–12 without leading zero, with AM/PM suffix (e.g., `2:30:05 PM`)

The AM/PM indicator is included in the font-sizing calculation, so the selected font will always fit the full time string.

## Date Format

The `date_format` option controls the date row:

- **YYYY-MM-DD** — ISO format (e.g., `2025-06-21`)
- **DD-MM-YYYY** — day first (e.g., `21-06-2025`)
- **MM-DD-YYYY** — month first (e.g., `06-21-2025`)
- **none** — hides the date row entirely

## Timezone

Set `timezone` to any valid IANA timezone identifier (e.g., `Europe/London`, `Asia/Tokyo`, `Pacific/Auckland`) to display time for that location. The default `local` uses your system's timezone. Invalid timezone strings silently fall back to local.

## Weekday Display

The weekday row (e.g., "Monday") is shown when `show_weekday` is `true` and the panel has enough vertical space to fit a third text row below time and date. On very small panels (e.g., 128×32), the weekday is automatically omitted regardless of this setting because there isn't room.

## Blinking Colon

When `blink_colon` is enabled, the colon separators in the time string alternate between visible (`:`) on even seconds and hidden (space) on odd seconds. This provides a visual activity indicator. The blink evaluates every second even when `show_seconds` is false.

## Border Frame

When `show_border` is enabled, an 8-pixel decorative tile border is drawn around the panel edge. Content is inset by 8 pixels on each side to avoid overlap. The border is identical across all styles and requires a panel of at least 16×16 pixels — smaller panels skip the border automatically.

## Foreground Color

The `fgcolor` option controls the text color on color-capable panels. The time row is rendered in the full foreground color, while date and weekday rows use a dimmed variant (each RGB channel halved) for visual hierarchy.

Available foreground colors:

- **cyan** (default) — bright cyan (0, 255, 255)
- **green** — soft green (0, 200, 0)
- **amber** — warm amber (255, 191, 0)
- **red** — vivid red (255, 0, 0)
- **white** — plain white (255, 255, 255)
- **none** — renders all text in white, same as "white"

On monochrome panels the foreground color setting is ignored and all text uses the panel's native foreground color.

## LED Seconds Indicator

The `show_led` option enables a small 6×6 pixel LED dot in the top-right corner of the display. It blinks on/off each second, providing a subtle activity indicator that confirms the clock is updating — useful when seconds digits are hidden.

The LED is automatically suppressed when:

- `show_seconds` is `true` (seconds digits already show updates)
- The LED's bounding box would overlap a text row

The LED defaults to on (`true`), so it appears automatically whenever seconds are hidden.

## Seconds Progress Bar

The `seconds_bar` option renders a visual progress indicator showing how far through the current minute the clock has advanced:

- **none** (default) — no bar shown
- **horizontal** — a thin 4-pixel-tall bar spanning the full content width at the bottom of the display, filling left-to-right as seconds tick from 0 to 59
- **pie** — a 16×16 pixel pie-chart widget in the bottom-right corner, sweeping clockwise

The progress bar is suppressed when the effective display height is less than 48 pixels. On color panels the bar uses the active foreground color; on monochrome panels it renders in white.

When the seconds bar is active, the available content height is reduced by the bar's height (4px for horizontal, 16px for pie) so text rows don't overlap it.

## Day Progress Sparkline

The `show_daybar` option renders a horizontal sparkline at the bottom of the display showing how far through the current 24-hour day you are. It computes progress as `(hour × 60 + minute) / 1440`.

The sparkline is 8 pixels tall and spans the full effective width, positioned 12 pixels from the bottom panel edge. On color panels it uses the foreground color at 30% brightness for a subtle background element.

The daybar requires a panel height of at least 128 pixels — it is suppressed on shorter displays.

## Layout Behavior

The clock mode optimizes vertical space automatically:

- **Row omission** — if content doesn't fit the panel height, rows are dropped from the bottom (weekday first, then date). The time row is always preserved.
- **Vertical centering** — the content block is centered vertically within the available space.
- **Horizontal centering** — each row is individually centered horizontally using its font metrics.
- **Adaptive fonts** — the largest available font that fits the panel width is selected automatically, recalculated on each render to respond to configuration changes.

## Panel Compatibility

The clock mode is non-interactive and works on all panels regardless of input controls. It does not require buttons, joystick, or any specific resolution — the selected style adapts its rendering to fit the available display area. On monochrome panels, color information is discarded and all text renders in the native foreground color. On slow-refresh panels, the blink_colon feature is automatically suppressed to avoid excessive redraws.

| Capability | Description | Behavior |
|---|---|---|
| MonoFast | Fast-refresh monochrome OLED (128×32, 128×64, 128×128) | Fully supported — time text renders in white with adaptive font sizing, colon blink animates smoothly |
| MonoSlow | Slow-refresh monochrome e-ink (122×250, 176×264, 200×200, 296×128, 400×300, 480×800, 800×480) | Fully supported — static time display, blink_colon forced off to prevent flicker, updates once per minute |
| GrayscaleFast | Fast-refresh grayscale (160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — grayscale text rendering with smooth colon blink animation |
| GrayscaleSlow | Slow-refresh grayscale e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static time display, blink_colon forced off, updates once per minute |
| ColorFast | Fast-refresh color TFT (128×128, 160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — full color rendering with accent-colored text and smooth animation |
| ColorSlow | Slow-refresh color e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static color time display, blink_colon forced off, updates once per minute |

!!! tip
    Use `show_seconds=false` and `blink_colon=false` on battery-powered setups to reduce redraw frequency from once per second to once per minute.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
- [Cycle](cycle.md) — auto-cycles through modes including clock
<!-- snapshot-gallery:start -->
## Snapshots

### Color

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/clock/color-slow-32x128_0001.png" alt="color-slow-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-64x128_0001.png" alt="color-slow-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-80x160_0001.png" alt="color-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-80x160_0001.png" alt="color-slow-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-104x212_0001.png" alt="color-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-122x250_0001.png" alt="color-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-128x32_0001.png" alt="color-slow-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-128x64_0001.png" alt="color-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-128x128_0001.png" alt="color-slow-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-128x128_0001.png" alt="color-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-128x160_0001.png" alt="color-slow-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-128x160_0001.png" alt="color-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-128x296_0001.png" alt="color-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-135x240_0001.png" alt="color-slow-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-135x240_0001.png" alt="color-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-160x80_0001.png" alt="color-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-160x80_0001.png" alt="color-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-160x128_0001.png" alt="color-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-160x128_0001.png" alt="color-slow-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-176x264_0001.png" alt="color-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-200x200_0001.png" alt="color-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-212x104_0001.png" alt="color-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-240x135_0001.png" alt="color-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-240x135_0001.png" alt="color-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-240x240_0001.png" alt="color-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-240x240_0001.png" alt="color-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-240x320_0001.png" alt="color-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-240x320_0001.png" alt="color-slow-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-250x122_0001.png" alt="color-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-264x176_0001.png" alt="color-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-296x128_0001.png" alt="color-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-300x400_0001.png" alt="color-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-320x240_0001.png" alt="color-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-320x240_0001.png" alt="color-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-320x480_0001.png" alt="color-slow-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-320x480_0001.png" alt="color-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-400x300_0001.png" alt="color-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-480x320_0001.png" alt="color-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-480x320_0001.png" alt="color-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-480x800_0001.png" alt="color-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-480x800_0001.png" alt="color-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-800x480_0001.png" alt="color-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/clock/color-slow-800x480_0001.png" alt="color-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Grayscale

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/clock/grayscale-slow-32x128_0001.png" alt="grayscale-slow-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-64x128_0001.png" alt="grayscale-slow-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-80x160_0001.png" alt="grayscale-fast-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-80x160_0001.png" alt="grayscale-slow-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-104x212_0001.png" alt="grayscale-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-122x250_0001.png" alt="grayscale-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-128x32_0001.png" alt="grayscale-slow-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-128x64_0001.png" alt="grayscale-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-128x128_0001.png" alt="grayscale-slow-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-128x128_0001.png" alt="grayscale-fast-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-128x160_0001.png" alt="grayscale-slow-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-128x160_0001.png" alt="grayscale-fast-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-128x296_0001.png" alt="grayscale-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-135x240_0001.png" alt="grayscale-slow-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-135x240_0001.png" alt="grayscale-fast-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-160x80_0001.png" alt="grayscale-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-160x80_0001.png" alt="grayscale-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-160x128_0001.png" alt="grayscale-fast-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-160x128_0001.png" alt="grayscale-slow-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-176x264_0001.png" alt="grayscale-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-200x200_0001.png" alt="grayscale-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-212x104_0001.png" alt="grayscale-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-240x135_0001.png" alt="grayscale-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-240x135_0001.png" alt="grayscale-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-240x240_0001.png" alt="grayscale-fast-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-240x240_0001.png" alt="grayscale-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-240x320_0001.png" alt="grayscale-fast-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-240x320_0001.png" alt="grayscale-slow-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-250x122_0001.png" alt="grayscale-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-264x176_0001.png" alt="grayscale-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-296x128_0001.png" alt="grayscale-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-300x400_0001.png" alt="grayscale-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-320x240_0001.png" alt="grayscale-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-320x240_0001.png" alt="grayscale-fast-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-320x480_0001.png" alt="grayscale-slow-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-320x480_0001.png" alt="grayscale-fast-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-400x300_0001.png" alt="grayscale-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-480x320_0001.png" alt="grayscale-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-480x320_0001.png" alt="grayscale-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-480x800_0001.png" alt="grayscale-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-480x800_0001.png" alt="grayscale-fast-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-fast-800x480_0001.png" alt="grayscale-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/clock/grayscale-slow-800x480_0001.png" alt="grayscale-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Mono

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/clock/mono-32x128_0001.png" alt="mono-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-32x128_0001.png" alt="mono-slow-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-64x128_0001.png" alt="mono-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-64x128_0001.png" alt="mono-slow-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-80x160_0001.png" alt="mono-slow-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-128x32_0001.png" alt="mono-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-128x32_0001.png" alt="mono-slow-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-128x64_0001.png" alt="mono-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-128x64_0001.png" alt="mono-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-128x128_0001.png" alt="mono-slow-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-128x128_0001.png" alt="mono-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-128x160_0001.png" alt="mono-slow-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-135x240_0001.png" alt="mono-slow-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-160x80_0001.png" alt="mono-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-160x128_0001.png" alt="mono-slow-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-240x135_0001.png" alt="mono-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-240x240_0001.png" alt="mono-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-240x320_0001.png" alt="mono-slow-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-320x240_0001.png" alt="mono-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-320x480_0001.png" alt="mono-slow-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-480x320_0001.png" alt="mono-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/clock/mono-slow-800x480_0001.png" alt="mono-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

<!-- snapshot-gallery:end -->




