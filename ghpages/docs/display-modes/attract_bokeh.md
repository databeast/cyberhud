# Attract Bokeh

The bokeh mode renders soft out-of-focus light circles that drift gently across the panel with varying sizes, opacities, and positions. The effect mimics bokeh photography — large translucent circles with radial gradient edges that overlap and blend to produce a dreamy, ambient aesthetic.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| `speed` | float | Drift speed multiplier controlling how quickly bokeh circles move across the panel surface | 1.0 | any valid float (0.1–10.0 recommended) |
| `density` | float | Circle count scaling factor from sparse (0.1) to maximum coverage (1.0) on the display | 0.5 | any valid float (0.1–1.0) |
| `size_variance` | float | Radius spread from uniform circles (0.0) to maximum size variation (1.0) in the scene | 0.5 | any valid float (0.0–1.0) |
| `saturation` | float | Color saturation level from desaturated pastel (0.0) to vivid fully-saturated hues (1.0) | 0.7 | any valid float (0.0–1.0) |

## Policy Fields

| Field | Type | Range | Default | Description |
|-------|------|-------|---------|-------------|
| Speed | float64 | [0.1, 10.0] | 1.0 | Drift speed multiplier |
| Density | float64 | [0.1, 1.0] | 0.5 | Circle count scaling |
| SizeVariance | float64 | [0.0, 1.0] | 0.5 | Radius spread (0=uniform, 1=max spread) |
| Saturation | float64 | [0.0, 1.0] | 0.7 | Color saturation |

## CLI Examples

Query the current policy values:

```sh
cyberhudctl display policy attract_bokeh
```

Set a faster drift speed with higher circle density:

```sh
cyberhudctl display config 0 speed=2.5 density=0.8
```

## Panel Compatibility

- **Color TFT** — full-color radial gradient circles with saturation control and smooth blending between overlapping circles.
- **Monochrome OLED** — circles use luminance-only values without color; gradients render as grayscale intensity falloff from center to edge.
- **E-Ink** — a static frozen bokeh scatter frame is produced (see below).

### E-Ink Behavior

On e-ink panels, the bokeh mode renders a static frozen bokeh scatter frame suitable for slow-refresh displays. Rather than animating continuous drift, a single arrangement of circles is composed and held indefinitely, avoiding the ghosting and flicker that partial refreshes would introduce. The Speed policy field has no visible effect on e-ink panels since no animation occurs.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Attract Modes](attract-modes.md) — overview of all attract/screensaver modes

<!-- snapshot-gallery:start -->
## Snapshots

### Color

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/attract_bokeh/color-8x16_0001.png" alt="color-8x16 8x16" style="max-width:320px;width:100%;">
  <figcaption>8x16</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-16x8_0001.png" alt="color-16x8 16x8" style="max-width:320px;width:100%;">
  <figcaption>16x8</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-80x160_0001.png" alt="color-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-104x212_0001.png" alt="color-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-122x250_0001.png" alt="color-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-128x128_0001.png" alt="color-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-128x160_0001.png" alt="color-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-128x296_0001.png" alt="color-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-135x240_0001.png" alt="color-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-160x80_0001.png" alt="color-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-160x128_0001.png" alt="color-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-176x264_0001.png" alt="color-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-200x200_0001.png" alt="color-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-212x104_0001.png" alt="color-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-240x135_0001.png" alt="color-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-240x240_0001.png" alt="color-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-240x320_0001.png" alt="color-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-250x122_0001.png" alt="color-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-264x176_0001.png" alt="color-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-296x128_0001.png" alt="color-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-300x400_0001.png" alt="color-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-320x240_0001.png" alt="color-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-320x480_0001.png" alt="color-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-400x300_0001.png" alt="color-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-480x320_0001.png" alt="color-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-480x800_0001.png" alt="color-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-480x800_0001.png" alt="color-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-800x480_0001.png" alt="color-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/color-slow-800x480_0001.png" alt="color-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Grayscale

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-8x16_0001.png" alt="grayscale-fast-8x16 8x16" style="max-width:320px;width:100%;">
  <figcaption>8x16</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-16x8_0001.png" alt="grayscale-fast-16x8 16x8" style="max-width:320px;width:100%;">
  <figcaption>16x8</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-80x160_0001.png" alt="grayscale-fast-80x160 80x160" style="max-width:320px;width:100%;">
  <figcaption>80x160</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-104x212_0001.png" alt="grayscale-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-122x250_0001.png" alt="grayscale-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-128x128_0001.png" alt="grayscale-fast-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-128x160_0001.png" alt="grayscale-fast-128x160 128x160" style="max-width:320px;width:100%;">
  <figcaption>128x160</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-128x296_0001.png" alt="grayscale-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-135x240_0001.png" alt="grayscale-fast-135x240 135x240" style="max-width:320px;width:100%;">
  <figcaption>135x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-160x80_0001.png" alt="grayscale-fast-160x80 160x80" style="max-width:320px;width:100%;">
  <figcaption>160x80</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-160x128_0001.png" alt="grayscale-fast-160x128 160x128" style="max-width:320px;width:100%;">
  <figcaption>160x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-176x264_0001.png" alt="grayscale-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-200x200_0001.png" alt="grayscale-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-212x104_0001.png" alt="grayscale-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-240x135_0001.png" alt="grayscale-fast-240x135 240x135" style="max-width:320px;width:100%;">
  <figcaption>240x135</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-240x240_0001.png" alt="grayscale-fast-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-240x320_0001.png" alt="grayscale-fast-240x320 240x320" style="max-width:320px;width:100%;">
  <figcaption>240x320</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-250x122_0001.png" alt="grayscale-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-264x176_0001.png" alt="grayscale-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-296x128_0001.png" alt="grayscale-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-300x400_0001.png" alt="grayscale-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-320x240_0001.png" alt="grayscale-fast-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-320x480_0001.png" alt="grayscale-fast-320x480 320x480" style="max-width:320px;width:100%;">
  <figcaption>320x480</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-400x300_0001.png" alt="grayscale-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-480x320_0001.png" alt="grayscale-fast-480x320 480x320" style="max-width:320px;width:100%;">
  <figcaption>480x320</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-480x800_0001.png" alt="grayscale-fast-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-480x800_0001.png" alt="grayscale-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-fast-800x480_0001.png" alt="grayscale-fast-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/grayscale-slow-800x480_0001.png" alt="grayscale-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

### Mono

<div class="grid" style="display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;">

<figure>
  <img src="../img/attract_bokeh/mono-8x16_0001.png" alt="mono-8x16 8x16" style="max-width:320px;width:100%;">
  <figcaption>8x16</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-16x8_0001.png" alt="mono-16x8 16x8" style="max-width:320px;width:100%;">
  <figcaption>16x8</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-32x128_0001.png" alt="mono-32x128 32x128" style="max-width:320px;width:100%;">
  <figcaption>32x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-64x128_0001.png" alt="mono-64x128 64x128" style="max-width:320px;width:100%;">
  <figcaption>64x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-104x212_0001.png" alt="mono-slow-104x212 104x212" style="max-width:320px;width:100%;">
  <figcaption>104x212</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-122x250_0001.png" alt="mono-slow-122x250 122x250" style="max-width:320px;width:100%;">
  <figcaption>122x250</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-128x32_0001.png" alt="mono-128x32 128x32" style="max-width:320px;width:100%;">
  <figcaption>128x32</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-128x64_0001.png" alt="mono-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-128x128_0001.png" alt="mono-128x128 128x128" style="max-width:320px;width:100%;">
  <figcaption>128x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-128x296_0001.png" alt="mono-slow-128x296 128x296" style="max-width:320px;width:100%;">
  <figcaption>128x296</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-176x264_0001.png" alt="mono-slow-176x264 176x264" style="max-width:320px;width:100%;">
  <figcaption>176x264</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-200x200_0001.png" alt="mono-slow-200x200 200x200" style="max-width:320px;width:100%;">
  <figcaption>200x200</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-212x104_0001.png" alt="mono-slow-212x104 212x104" style="max-width:320px;width:100%;">
  <figcaption>212x104</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-250x122_0001.png" alt="mono-slow-250x122 250x122" style="max-width:320px;width:100%;">
  <figcaption>250x122</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-264x176_0001.png" alt="mono-slow-264x176 264x176" style="max-width:320px;width:100%;">
  <figcaption>264x176</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-296x128_0001.png" alt="mono-slow-296x128 296x128" style="max-width:320px;width:100%;">
  <figcaption>296x128</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-300x400_0001.png" alt="mono-slow-300x400 300x400" style="max-width:320px;width:100%;">
  <figcaption>300x400</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-400x300_0001.png" alt="mono-slow-400x300 400x300" style="max-width:320px;width:100%;">
  <figcaption>400x300</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-480x800_0001.png" alt="mono-slow-480x800 480x800" style="max-width:320px;width:100%;">
  <figcaption>480x800</figcaption>
</figure>

<figure>
  <img src="../img/attract_bokeh/mono-slow-800x480_0001.png" alt="mono-slow-800x480 800x480" style="max-width:320px;width:100%;">
  <figcaption>800x480</figcaption>
</figure>

</div>

<!-- snapshot-gallery:end -->




