# STEMMA

Lists detected STEMMA QT / QWIIC devices on the I2C bus. Shows device addresses, bus assignments, and presence status with periodic re-scanning.

## Quick Start

```sh
cyberhudctl display set 0 stemma
```

## How It Works

The stemma mode renders a list of I2C devices detected on configured STEMMA QT / QWIIC buses, showing each device as a row with a presence icon (green check for present, red error for absent), bus name, hexadecimal address, and a human-readable device name from a built-in lookup table of common Adafruit and SparkFun peripherals. The rows are arranged vertically in list order, with the number of visible rows adapting to the panel's available height.

The display updates on a fixed interval driven by the background scanner goroutine, which rescans all configured I2C buses periodically (interval set at daemon startup via scanner configuration). Each scan queries every known address on each bus and updates the presence flag for each device. A change in any device's presence state triggers a cache-key change, causing the display to redraw.

The data source is the Linux I2C subsystem accessed via direct bus reads on the configured bus paths. When no I2C devices are detected, the mode displays a placeholder message. When the scanner is nil (no I2C hardware available or not yet initialized), the mode gracefully returns a valid view with zero data items. Unrecognized device addresses are labeled "unknown" in the device name column.

The `style` option controls the visual presentation:

- **list** (default) — displays each detected device as a separate row with a presence icon, bus assignment, hex address, and device name. Rows are color-coded green (present) or red (absent). The number of visible rows adapts to the panel's available vertical space.
- **compact** — shows a single summary line with the count of present devices out of total detected (e.g., "3/5 present"). Useful for narrow panels or when a quick glance is all you need.

!!! tip
    Use this mode to verify hardware connections before switching to a mode that depends on a specific I2C peripheral.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| style | string | Visual layout or explicit profile name | list | list, compact, and the registered profile matrix |

Configure options via the CLI:

```sh
cyberhudctl display stemma style=<value>
```

### Style Variants

- **list** — displays each detected device as a separate row with a presence icon, bus assignment, hex address, and device name. Rows are color-coded green (present) or red (absent). The number of visible rows adapts to the panel's available vertical space.
- **compact** — shows a single summary line with the count of present devices out of total detected (e.g. `3/5 present`). Useful for narrow panels or when a quick glance is all you need.
- **profile variants** — the registry also exposes explicit panel profiles across mono, grayscale, and color families so the mode resolves like the other display modes instead of relying on one generic layout.

## CLI Examples

Activate the stemma mode on the main region:

```sh
cyberhudctl display set 0 stemma
```

Query the current stemma settings:

```sh
cyberhudctl display stemma
```

Switch to compact summary layout:

```sh
cyberhudctl display stemma style=compact
```

Switch back to the full device list:

```sh
cyberhudctl display stemma style=list
```

## Panel Compatibility

The stemma mode is non-interactive and works on all panels. It does not require input controls — no buttons or joystick pins are needed. On monochrome panels, device labels and sensor values render in the native foreground color without color differentiation. The mode adapts its layout density based on panel dimensions to show as many detected I2C devices as will fit.

| Capability | Description | Behavior |
|---|---|---|
| MonoFast | Fast-refresh monochrome OLED (128×32, 128×64, 128×128) | Fully supported — I2C device list with addresses and values in white text |
| MonoSlow | Slow-refresh monochrome e-ink (122×250, 176×264, 200×200, 296×128, 400×300, 480×800, 800×480) | Fully supported — static device list, refreshes on scan interval |
| GrayscaleFast | Fast-refresh grayscale (160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — grayscale text rendering of device data |
| GrayscaleSlow | Slow-refresh grayscale e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static grayscale device list, refreshes on scan interval |
| ColorFast | Fast-refresh color TFT (128×128, 160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — color-coded device status with periodic scan updates |
| ColorSlow | Slow-refresh color e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — static color device list, refreshes on scan interval |

On panels with buttons, the primary button navigates back to the dashboard.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
- [Dashboard](dashboard.md) — consolidated status overview including STEMMA device summary
<!-- snapshot-gallery:start -->
## Snapshots

### Color

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/stemma/color-slow-32x128_0001.png" alt="color-slow-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-64x128_0001.png" alt="color-slow-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-80x160_0001.png" alt="color-slow-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-80x160_0001.png" alt="color-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-104x212_0001.png" alt="color-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-122x250_0001.png" alt="color-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-128x32_0001.png" alt="color-slow-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-128x64_0001.png" alt="color-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-128x128_0001.png" alt="color-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-128x128_0001.png" alt="color-slow-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-128x160_0001.png" alt="color-slow-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-128x160_0001.png" alt="color-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-128x296_0001.png" alt="color-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-135x240_0001.png" alt="color-slow-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-160x80_0001.png" alt="color-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-160x80_0001.png" alt="color-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-160x128_0001.png" alt="color-slow-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-160x128_0001.png" alt="color-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-176x264_0001.png" alt="color-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-200x200_0001.png" alt="color-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-212x104_0001.png" alt="color-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-240x135_0001.png" alt="color-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-240x240_0001.png" alt="color-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-240x320_0001.png" alt="color-slow-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-250x122_0001.png" alt="color-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-264x176_0001.png" alt="color-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-296x128_0001.png" alt="color-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-300x400_0001.png" alt="color-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-320x240_0001.png" alt="color-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-320x480_0001.png" alt="color-slow-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-400x300_0001.png" alt="color-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-480x320_0001.png" alt="color-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-480x800_0001.png" alt="color-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/stemma/color-slow-800x480_0001.png" alt="color-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Grayscale

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/stemma/grayscale-slow-32x128_0001.png" alt="grayscale-slow-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-64x128_0001.png" alt="grayscale-slow-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-80x160_0001.png" alt="grayscale-fast-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-80x160_0001.png" alt="grayscale-slow-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-104x212_0001.png" alt="grayscale-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-122x250_0001.png" alt="grayscale-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-128x32_0001.png" alt="grayscale-slow-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-128x64_0001.png" alt="grayscale-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-128x128_0001.png" alt="grayscale-slow-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-128x128_0001.png" alt="grayscale-fast-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-128x160_0001.png" alt="grayscale-slow-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-128x160_0001.png" alt="grayscale-fast-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-128x296_0001.png" alt="grayscale-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-135x240_0001.png" alt="grayscale-slow-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-135x240_0001.png" alt="grayscale-fast-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-160x80_0001.png" alt="grayscale-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-160x80_0001.png" alt="grayscale-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-160x128_0001.png" alt="grayscale-fast-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-160x128_0001.png" alt="grayscale-slow-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-176x264_0001.png" alt="grayscale-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-200x200_0001.png" alt="grayscale-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-212x104_0001.png" alt="grayscale-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-240x135_0001.png" alt="grayscale-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-240x135_0001.png" alt="grayscale-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-240x240_0001.png" alt="grayscale-fast-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-240x240_0001.png" alt="grayscale-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-240x320_0001.png" alt="grayscale-fast-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-240x320_0001.png" alt="grayscale-slow-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-250x122_0001.png" alt="grayscale-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-264x176_0001.png" alt="grayscale-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-296x128_0001.png" alt="grayscale-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-300x400_0001.png" alt="grayscale-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-320x240_0001.png" alt="grayscale-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-320x240_0001.png" alt="grayscale-fast-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-320x480_0001.png" alt="grayscale-slow-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-320x480_0001.png" alt="grayscale-fast-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-400x300_0001.png" alt="grayscale-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-480x320_0001.png" alt="grayscale-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-480x320_0001.png" alt="grayscale-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-480x800_0001.png" alt="grayscale-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-480x800_0001.png" alt="grayscale-fast-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-fast-800x480_0001.png" alt="grayscale-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/stemma/grayscale-slow-800x480_0001.png" alt="grayscale-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Mono

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/stemma/mono-32x128_0001.png" alt="mono-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-32x128_0001.png" alt="mono-slow-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-64x128_0001.png" alt="mono-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-64x128_0001.png" alt="mono-slow-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-80x160_0001.png" alt="mono-slow-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-128x32_0001.png" alt="mono-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-128x32_0001.png" alt="mono-slow-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-128x64_0001.png" alt="mono-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-128x64_0001.png" alt="mono-slow-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-128x128_0001.png" alt="mono-slow-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-128x128_0001.png" alt="mono-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-128x160_0001.png" alt="mono-slow-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-135x240_0001.png" alt="mono-slow-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-160x80_0001.png" alt="mono-slow-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-160x128_0001.png" alt="mono-slow-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-240x135_0001.png" alt="mono-slow-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-240x240_0001.png" alt="mono-slow-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-240x320_0001.png" alt="mono-slow-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-320x240_0001.png" alt="mono-slow-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-320x480_0001.png" alt="mono-slow-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-480x320_0001.png" alt="mono-slow-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/stemma/mono-slow-800x480_0001.png" alt="mono-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

<!-- snapshot-gallery:end -->


