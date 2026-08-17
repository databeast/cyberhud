# Image

Displays a static image on the panel. Useful for splash screens, logos, or custom artwork rendered to the display framebuffer. Content is pushed via the CLI from an external image file.

## Quick Start

```sh
cyberhudctl display set 0 image
```

## How It Works

The image mode renders a single externally-supplied PNG or JPEG image as a full-panel sprite, covering the entire display surface according to the configured fit policy. The display is static and event-driven — it only redraws when a new image is pushed via the CLI (`cyberhudctl display image set`) or the fit/style policy is changed, with no periodic polling or refresh.

The image data source is entirely external: you provide an image file path or base64-encoded data through the console command, and the mode decodes it into memory. When no image has been set (or after a `clear` command), the mode shows a single text row "(no image)" as a placeholder. The decoded image persists in memory until explicitly replaced or cleared.

The **default** style renders the image without decoration; the **bordered** style adds an 8-pixel decorative tile border frame around the panel edge and insets the image content by 8 pixels on each side to avoid overlap. The border is suppressed on panels smaller than 16×16 pixels. Within the content area, scaling is controlled by the `fit` option: **scale** (default) uniformly scales the image to the largest size that preserves aspect ratio, **stretch** fills the entire viewport ignoring proportions, and **truncate** draws at native resolution and clips any overflow.

!!! tip
    Use `fit=scale` (the default) for most images to ensure the full content is visible. Switch to `fit=stretch` only when filling the entire panel area matters more than preserving proportions.

## Options

| Key | Type | Description | Default | Allowed Values |
|-----|------|-------------|---------|----------------|
| fit | string | How the image is adjusted to panel dimensions | scale | truncate, scale, stretch |
| style | string | Border treatment around the image | default | default, bordered |

Configure options via the CLI:

```sh
cyberhudctl display image fit=<value> style=<value>
```

## CLI Examples

Activate the image mode on the main region:

```sh
cyberhudctl display set 0 image
```

Scale an image to fit the panel while preserving aspect ratio:

```sh
cyberhudctl display image fit=scale
```

Stretch an image to fill the entire panel area:

```sh
cyberhudctl display image fit=stretch
```

Add a decorative border frame around the image:

```sh
cyberhudctl display image style=bordered
```

Combine fit and style options in one command:

```sh
cyberhudctl display image fit=scale style=bordered
```

Query current image settings:

```sh
cyberhudctl display image
```

## Panel Compatibility

Image mode is non-interactive and works on all panels. Image rendering adapts to the panel's pixel dimensions and color depth — on monochrome panels, images are dithered to 1-bit; on grayscale panels, images are converted to luminance values. Larger panels show more detail, while smaller panels benefit from the `scale` fit policy to avoid clipping.

| Capability | Description | Behavior |
|---|---|---|
| MonoFast | Fast-refresh monochrome OLED (128×32, 128×64, 128×128) | Fully supported — image dithered to 1-bit monochrome, scaled to panel dimensions |
| MonoSlow | Slow-refresh monochrome e-ink (122×250, 176×264, 200×200, 296×128, 400×300, 480×800, 800×480) | Fully supported — 1-bit dithered image, single static frame rendered on load |
| GrayscaleFast | Fast-refresh grayscale (160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — image converted to grayscale with luminance mapping |
| GrayscaleSlow | Slow-refresh grayscale e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — grayscale image rendered as single static frame |
| ColorFast | Fast-refresh color TFT (128×128, 160×80, 240×135, 240×240, 320×240, 480×320, 800×480) | Fully supported — full color image rendering with native color depth |
| ColorSlow | Slow-refresh color e-ink (122×250, 176×264, 200×200, 400×300, 800×480) | Fully supported — full color image rendered as single static frame |

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Getting Started: CLI Usage](../getting-started/cli.md) — introduction to `cyberhudctl` commands
- [Configuration](../configuration/index.md) — persistent configuration options
<!-- snapshot-gallery:start -->
## Snapshots

### Color

<figure>
  <img src="../img/image/color-240x240-scale-bordered_0001.png" alt="color-240x240-scale-bordered 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/image/color-320x240-stretch_0001.png" alt="color-320x240-stretch 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

### Mono

<figure>
  <img src="../img/image/mono-128x64-truncate_0001.png" alt="mono-128x64-truncate 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<!-- snapshot-gallery:end -->
