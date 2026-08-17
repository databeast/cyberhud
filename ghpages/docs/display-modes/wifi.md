# WiFi

The wifi display mode shows real-time wireless network connection status including signal strength, SSID, IP address, channel, frequency, and link speed. It reads interface state via NL80211 (netlink) on Linux and presents connection details through a visually rich layout with a WiFi icon sprite, color-coded signal bars, and a link quality progress bar.

## Quick Start

```sh
cyberhudctl display set 0 wifi
```

## How It Works

The WiFi mode queries the system's wireless interface via NL80211 netlink on each render cycle and displays connection details as a multi-row layout with graphical indicators. The display renders a WiFi icon sprite, signal strength bars (color-coded by quality level), a horizontal link quality progress bar, and detail rows showing SSID, frequency band, channel number, IP address, and interface name. Each render produces a data-signature fingerprint incorporating connection state, signal metrics, and all policy fields so the display runtime skips unnecessary re-renders when nothing has changed.

When no wireless interface is detected the mode displays "WiFi N/A", and when the interface exists but is not associated it shows "No Network" — both states render gracefully without errors. The signal strength is represented in one of three formats controlled by the `signal_display` option: vertical bars (0–4 levels mapped from dBm thresholds), a percentage (linear interpolation from -100 dBm to -30 dBm), or raw dBm value. On color panels, signal bars use quality-mapped colors: the accent foreground color for strong signals (3–4 bars), amber for moderate (2 bars), and red for weak (0–1 bars).

The mode supports multiple resolution-specific styles organized by panel capability class (MonoSlow, MonoFast, GrayscaleSlow, GrayscaleFast, ColorSlow, ColorFast) with automatic style selection based on panel hints. Users can cycle styles via Left/Right input actions, force a refresh with the Secondary button, or navigate to the menu with Primary.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| `fgcolor` | string | Foreground color for WiFi display elements including signal bars, icon, and text on color-capable panels | green | cyan, green, amber, red, white, none |
| `signal_display` | string | Signal strength display format controlling how wireless signal quality is shown visually on the panel | bars | bars, percentage, dbm |
| `show_frequency` | bool | Show the WiFi frequency band (2.4GHz or 5GHz) in the connection details area of the display | true | true, false |
| `show_interface` | bool | Show the wireless interface name (e.g. wlan0) in the connection details area on the display | true | true, false |
| `show_channel` | bool | Show the WiFi channel number derived from the current frequency in the details area | true | true, false |
| `style` | string | Visual layout variant for WiFi display controlling spatial arrangement and detail level per panel resolution | | Resolution-specific style names |

## CLI Examples

Activate the WiFi mode on the main display region:

```sh
cyberhudctl display set 0 wifi
```

Switch the signal display format to show dBm values instead of bars:

```sh
cyberhudctl display wifi signal_display=dbm
```

Set a cyan foreground color and hide the interface name:

```sh
cyberhudctl display wifi fgcolor=cyan show_interface=false
```

Show signal as a percentage and hide the channel number:

```sh
cyberhudctl display wifi signal_display=percentage show_channel=false
```

Query all current WiFi mode settings:

```sh
cyberhudctl display wifi
```

## Panel Compatibility

The WiFi mode is fully supported across all six panel capability classes through resolution-specific style variants that adapt the layout to each display type. All 17 resolution families are available for every capability tier.

### Resolution Families

| Resolution | Portrait Variant | Styles |
|-----------|-----------------|--------|
| 128×32 | 32×128 | 12 |
| 128×64 | — | 6 |
| 128×128 | — | 6 |
| 160×80 | 80×160 | 12 |
| 160×128 | 128×160 | 12 |
| 200×200 | — | 6 |
| 212×104 | — | 6 |
| 240×135 | 135×240 | 12 |
| 240×240 | — | 6 |
| 250×122 | — | 6 |
| 264×176 | — | 6 |
| 296×128 | 128×296 | 12 |
| 320×240 | 240×320 | 12 |
| 400×300 | — | 6 |
| 480×320 | 320×480 | 12 |
| 800×480 | 480×800 | 12 |
| 800×600 | — | 6 |

### Capability Tiers

| Capability | Description | Behavior |
|------------|-------------|----------|
| MonoSlow | Monochrome e-paper displays with slow refresh | Fully supported. Static text-only layout with border frame optimized for infrequent refresh. Available at all 17 resolutions. |
| MonoFast | Monochrome OLED displays with fast refresh | Fully supported. Renders WiFi icon, signal bars, and text rows in white. Available at all 17 resolutions. |
| GrayscaleSlow | Grayscale e-paper with slow refresh | Fully supported. Static layout with grayscale shading optimized for slow refresh panels. Available at all 17 resolutions. |
| GrayscaleFast | Grayscale LED matrix or fast grayscale displays | Fully supported. Uses grayscale shading for signal indicators with icon and progress bar. Available at all 17 resolutions. |
| ColorSlow | Color e-paper with slow refresh | Fully supported. Color layout with accent colors adapted for infrequent refresh. Available at all 17 resolutions. |
| ColorFast | Color TFT displays with fast refresh | Fully supported. Full color rendering with quality-mapped signal bar colors and configurable foreground accent. Available at all 17 resolutions. |

The mode does not require physical input controls — joystick Left/Right cycle styles and Primary navigates to menu, but these are optional conveniences. WiFi status displays correctly on headless panels without any input hardware.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
- [Cycle](cycle.md) — auto-cycles through modes including wifi
- [Dashboard](dashboard.md) — overview mode that also displays WiFi SSID

<!-- snapshot-gallery:start -->
## Snapshots

### Color

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/wifi/color-slow-32x128_0001.png" alt="color-slow-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-32x128_0001.png" alt="color-fast-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-80x160_0001.png" alt="color-fast-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-80x160_0001.png" alt="color-slow-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-128x32_0001.png" alt="color-slow-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-128x32_0001.png" alt="color-fast-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-128x64_0001.png" alt="color-fast-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-128x64_0001.png" alt="color-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-128x128_0001.png" alt="color-slow-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-128x128_0001.png" alt="color-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-128x160_0001.png" alt="color-slow-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-128x160_0001.png" alt="color-fast-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-128x296_0001.png" alt="color-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-128x296_0001.png" alt="color-fast-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-135x240_0001.png" alt="color-slow-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-135x240_0001.png" alt="color-fast-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-160x80_0001.png" alt="color-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-160x80_0001.png" alt="color-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-160x128_0001.png" alt="color-slow-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-160x128_0001.png" alt="color-fast-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-200x200_0001.png" alt="color-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-200x200_0001.png" alt="color-fast-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-212x104_0001.png" alt="color-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-212x104_0001.png" alt="color-fast-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-240x135_0001.png" alt="color-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-240x135_0001.png" alt="color-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-240x240_0001.png" alt="color-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-240x240_0001.png" alt="color-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-240x320_0001.png" alt="color-fast-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-240x320_0001.png" alt="color-slow-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-250x122_0001.png" alt="color-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-250x122_0001.png" alt="color-fast-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-264x176_0001.png" alt="color-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-264x176_0001.png" alt="color-fast-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-296x128_0001.png" alt="color-fast-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-296x128_0001.png" alt="color-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-320x240_0001.png" alt="color-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-320x240_0001.png" alt="color-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-320x480_0001.png" alt="color-fast-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-320x480_0001.png" alt="color-slow-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-400x300_0001.png" alt="color-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-400x300_0001.png" alt="color-fast-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-480x320_0001.png" alt="color-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-480x320_0001.png" alt="color-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-480x800_0001.png" alt="color-fast-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-480x800_0001.png" alt="color-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-800x480_0001.png" alt="color-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-800x480_0001.png" alt="color-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-fast-800x600_0001.png" alt="color-fast-800x600 800x600" style="max-width:320px;width:100%;">
  <figcaption>800x600</figcaption>
</figure>

<figure>
  <img src="../img/wifi/color-slow-800x600_0001.png" alt="color-slow-800x600 800x600" style="max-width:320px;width:100%;">
  <figcaption>800x600</figcaption>
</figure>

</div>

### Grayscale

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/wifi/grayscale-fast-32x128_0001.png" alt="grayscale-fast-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-32x128_0001.png" alt="grayscale-slow-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-80x160_0001.png" alt="grayscale-slow-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-80x160_0001.png" alt="grayscale-fast-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-128x32_0001.png" alt="grayscale-slow-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-128x32_0001.png" alt="grayscale-fast-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-128x64_0001.png" alt="grayscale-fast-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-128x64_0001.png" alt="grayscale-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-128x128_0001.png" alt="grayscale-slow-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-128x128_0001.png" alt="grayscale-fast-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-128x160_0001.png" alt="grayscale-fast-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-128x160_0001.png" alt="grayscale-slow-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-128x296_0001.png" alt="grayscale-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-128x296_0001.png" alt="grayscale-fast-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-135x240_0001.png" alt="grayscale-fast-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-135x240_0001.png" alt="grayscale-slow-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-160x80_0001.png" alt="grayscale-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-160x80_0001.png" alt="grayscale-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-160x128_0001.png" alt="grayscale-slow-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-160x128_0001.png" alt="grayscale-fast-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-200x200_0001.png" alt="grayscale-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-200x200_0001.png" alt="grayscale-fast-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-212x104_0001.png" alt="grayscale-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-212x104_0001.png" alt="grayscale-fast-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-240x135_0001.png" alt="grayscale-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-240x135_0001.png" alt="grayscale-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-240x240_0001.png" alt="grayscale-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-240x240_0001.png" alt="grayscale-fast-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-240x320_0001.png" alt="grayscale-fast-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-240x320_0001.png" alt="grayscale-slow-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-250x122_0001.png" alt="grayscale-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-250x122_0001.png" alt="grayscale-fast-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-264x176_0001.png" alt="grayscale-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-264x176_0001.png" alt="grayscale-fast-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-296x128_0001.png" alt="grayscale-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-296x128_0001.png" alt="grayscale-fast-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-320x240_0001.png" alt="grayscale-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-320x240_0001.png" alt="grayscale-fast-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-320x480_0001.png" alt="grayscale-fast-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-320x480_0001.png" alt="grayscale-slow-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-400x300_0001.png" alt="grayscale-fast-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-400x300_0001.png" alt="grayscale-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-480x320_0001.png" alt="grayscale-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-480x320_0001.png" alt="grayscale-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-480x800_0001.png" alt="grayscale-fast-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-480x800_0001.png" alt="grayscale-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-800x480_0001.png" alt="grayscale-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-800x480_0001.png" alt="grayscale-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-fast-800x600_0001.png" alt="grayscale-fast-800x600 800x600" style="max-width:320px;width:100%;">
  <figcaption>800x600</figcaption>
</figure>

<figure>
  <img src="../img/wifi/grayscale-slow-800x600_0001.png" alt="grayscale-slow-800x600 800x600" style="max-width:320px;width:100%;">
  <figcaption>800x600</figcaption>
</figure>

</div>

### Mono

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/wifi/mono-fast-32x128_0001.png" alt="mono-fast-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-32x128_0001.png" alt="mono-slow-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-80x160_0001.png" alt="mono-slow-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-80x160_0001.png" alt="mono-fast-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-128x32_0001.png" alt="mono-slow-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-128x32_0001.png" alt="mono-fast-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-128x64_0001.png" alt="mono-fast-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-128x64_0001.png" alt="mono-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-128x128_0001.png" alt="mono-slow-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-128x128_0001.png" alt="mono-fast-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-128x160_0001.png" alt="mono-fast-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-128x160_0001.png" alt="mono-slow-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-128x296_0001.png" alt="mono-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-128x296_0001.png" alt="mono-fast-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-135x240_0001.png" alt="mono-fast-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-135x240_0001.png" alt="mono-slow-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-160x80_0001.png" alt="mono-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-160x80_0001.png" alt="mono-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-160x128_0001.png" alt="mono-slow-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-160x128_0001.png" alt="mono-fast-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-200x200_0001.png" alt="mono-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-200x200_0001.png" alt="mono-fast-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-212x104_0001.png" alt="mono-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-212x104_0001.png" alt="mono-fast-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-240x135_0001.png" alt="mono-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-240x135_0001.png" alt="mono-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-240x240_0001.png" alt="mono-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-240x240_0001.png" alt="mono-fast-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-240x320_0001.png" alt="mono-fast-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-240x320_0001.png" alt="mono-slow-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-250x122_0001.png" alt="mono-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-250x122_0001.png" alt="mono-fast-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-264x176_0001.png" alt="mono-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-264x176_0001.png" alt="mono-fast-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-296x128_0001.png" alt="mono-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-296x128_0001.png" alt="mono-fast-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-320x240_0001.png" alt="mono-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-320x240_0001.png" alt="mono-fast-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-320x480_0001.png" alt="mono-fast-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-320x480_0001.png" alt="mono-slow-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-400x300_0001.png" alt="mono-fast-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-400x300_0001.png" alt="mono-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-480x320_0001.png" alt="mono-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-480x320_0001.png" alt="mono-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-480x800_0001.png" alt="mono-fast-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-480x800_0001.png" alt="mono-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-800x480_0001.png" alt="mono-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-800x480_0001.png" alt="mono-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-fast-800x600_0001.png" alt="mono-fast-800x600 800x600" style="max-width:320px;width:100%;">
  <figcaption>800x600</figcaption>
</figure>

<figure>
  <img src="../img/wifi/mono-slow-800x600_0001.png" alt="mono-slow-800x600 800x600" style="max-width:320px;width:100%;">
  <figcaption>800x600</figcaption>
</figure>

</div>

<!-- snapshot-gallery:end -->



