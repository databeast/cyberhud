# Attract Plasma

The plasma mode renders a lava-lamp-style plasma blob effect by evaluating overlapping sine-based color functions at each pixel position. The pattern morphs continuously over time, producing smooth organic blob shapes that shift and merge through a cycling 256-step gradient palette.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| `speed` | float | Morph speed multiplier controlling how quickly the plasma blobs shift and evolve over time | 1.0 | any valid float (0.1–10.0 recommended) |
| `density` | float | Density parameter kept for uniformity across attract modes, does not visibly affect plasma output | 0.5 | any valid float (0.1–1.0) |
| `cycle_rate` | float | Color palette cycles per second controlling how rapidly the gradient hues rotate through the spectrum | 1.0 | any valid float (0.1–5.0 recommended) |
| `blob_scale` | float | Spatial frequency multiplier for blob size, higher values produce smaller more numerous blobs | 1.0 | any valid float (0.5–4.0) |

## Policy Fields

| Field | Type | Range | Default | Description |
|-------|------|-------|---------|-------------|
| Speed | float64 | [0.1, 10.0] | 1.0 | Morph speed multiplier |
| Density | float64 | [0.1, 1.0] | 0.5 | Reserved for uniformity |
| CycleRate | float64 | [0.1, 5.0] | 1.0 | Color palette cycles/second |
| BlobScale | float64 | [0.5, 4.0] | 1.0 | Spatial frequency multiplier |

## CLI Examples

Query current plasma policy settings:

```sh
cyberhudctl display attract_plasma
```

Set a faster morph speed with tighter blob scaling:

```sh
cyberhudctl display config 0 speed=3.0 blob_scale=2.0
```

## Panel Compatibility

- **Color TFT** — full RGB plasma rendering with smooth gradient transitions across the 256-step palette. This is the ideal panel type for the plasma effect.
- **Monochrome OLED** — the plasma pattern is rendered in grayscale, mapping palette intensity to white levels. Blob shapes remain visible but without color variation.
- **E-Ink** — renders a single static frame due to slow refresh characteristics. See E-Ink Behavior below.

### E-Ink Behavior

On e-ink panels, the plasma mode renders a single static plasma frame at time=0 for slow-refresh panels. The RenderCacheKey remains stable after the initial frame, preventing unnecessary refreshes that would cause flicker or ghosting on e-ink displays.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Attract Modes](attract-modes.md) — overview of all attract/screensaver modes

<!-- snapshot-gallery:start -->
## Snapshots

### Color

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/attract_plasma/color-8x16_0001.png" alt="color-8x16 8x16" style="max-width:320px;width:100%;">
  <figcaption>8x16</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-16x8_0001.png" alt="color-16x8 16x8" style="max-width:320px;width:100%;">
  <figcaption>16x8</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-80x160_0001.png" alt="color-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-104x212_0001.png" alt="color-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-122x250_0001.png" alt="color-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-128x128_0001.png" alt="color-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-128x160_0001.png" alt="color-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-128x296_0001.png" alt="color-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-135x240_0001.png" alt="color-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-160x80_0001.png" alt="color-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-160x128_0001.png" alt="color-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-176x264_0001.png" alt="color-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-200x200_0001.png" alt="color-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-212x104_0001.png" alt="color-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-240x135_0001.png" alt="color-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-240x240_0001.png" alt="color-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-240x320_0001.png" alt="color-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-250x122_0001.png" alt="color-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-264x176_0001.png" alt="color-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-296x128_0001.png" alt="color-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-300x400_0001.png" alt="color-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-320x240_0001.png" alt="color-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-320x480_0001.png" alt="color-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-400x300_0001.png" alt="color-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-480x320_0001.png" alt="color-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-480x800_0001.png" alt="color-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-480x800_0001.png" alt="color-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-800x480_0001.png" alt="color-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/color-slow-800x480_0001.png" alt="color-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Grayscale

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/attract_plasma/grayscale-fast-8x16_0001.png" alt="grayscale-fast-8x16 8x16" style="max-width:320px;width:100%;">
  <figcaption>8x16</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-16x8_0001.png" alt="grayscale-fast-16x8 16x8" style="max-width:320px;width:100%;">
  <figcaption>16x8</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-80x160_0001.png" alt="grayscale-fast-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-104x212_0001.png" alt="grayscale-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-122x250_0001.png" alt="grayscale-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-128x128_0001.png" alt="grayscale-fast-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-128x160_0001.png" alt="grayscale-fast-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-128x296_0001.png" alt="grayscale-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-135x240_0001.png" alt="grayscale-fast-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-160x80_0001.png" alt="grayscale-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-160x128_0001.png" alt="grayscale-fast-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-176x264_0001.png" alt="grayscale-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-200x200_0001.png" alt="grayscale-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-212x104_0001.png" alt="grayscale-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-240x135_0001.png" alt="grayscale-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-240x240_0001.png" alt="grayscale-fast-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-240x320_0001.png" alt="grayscale-fast-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-250x122_0001.png" alt="grayscale-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-264x176_0001.png" alt="grayscale-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-296x128_0001.png" alt="grayscale-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-300x400_0001.png" alt="grayscale-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-320x240_0001.png" alt="grayscale-fast-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-320x480_0001.png" alt="grayscale-fast-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-400x300_0001.png" alt="grayscale-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-480x320_0001.png" alt="grayscale-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-480x800_0001.png" alt="grayscale-fast-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-480x800_0001.png" alt="grayscale-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-fast-800x480_0001.png" alt="grayscale-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/grayscale-slow-800x480_0001.png" alt="grayscale-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Mono

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/attract_plasma/mono-8x16_0001.png" alt="mono-8x16 8x16" style="max-width:320px;width:100%;">
  <figcaption>8x16</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-16x8_0001.png" alt="mono-16x8 16x8" style="max-width:320px;width:100%;">
  <figcaption>16x8</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-32x128_0001.png" alt="mono-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-64x128_0001.png" alt="mono-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-104x212_0001.png" alt="mono-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-122x250_0001.png" alt="mono-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-128x32_0001.png" alt="mono-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-128x64_0001.png" alt="mono-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-128x128_0001.png" alt="mono-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-128x296_0001.png" alt="mono-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-176x264_0001.png" alt="mono-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-200x200_0001.png" alt="mono-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-212x104_0001.png" alt="mono-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-250x122_0001.png" alt="mono-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-264x176_0001.png" alt="mono-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-296x128_0001.png" alt="mono-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-300x400_0001.png" alt="mono-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-400x300_0001.png" alt="mono-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-480x800_0001.png" alt="mono-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_plasma/mono-slow-800x480_0001.png" alt="mono-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

<!-- snapshot-gallery:end -->




