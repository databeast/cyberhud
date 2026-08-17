# Test Widgets

Cycles through all registered UI widget types showing a demo preview of each. Use this mode to verify that widget rendering is correct on your panel and to preview the visual style of each widget component available in CyberHUD.

## Quick Start

```sh
cyberhudctl display set 0 testwidgets
```

## How It Works

The test widgets mode showcases every registered widget type one at a time, cycling through them on a fixed 8-second interval. Each widget is rendered with sample data so you can see its visual appearance and verify correct rendering on your panel.

The following widget types are demonstrated in sequence:

- **borderframe** — decorative border container
- **gradient** — color/shade gradient fill
- **LED** — simulated LED indicator
- **progressbar** — horizontal progress indicator
- **scaledtextbox** — auto-scaling text container
- **scrollbar** — vertical scroll position indicator
- **sparkline** — miniature line chart
- **textbox** — standard text container
- **textlabel** — single-line text label

After displaying the last widget, the cycle restarts from the beginning.

## Options

This mode has no configurable options.

## Panel Compatibility

Test widgets is non-interactive and ignores all input actions. Works on all panels. Useful for verifying widget rendering across different display resolutions and color modes.

## Related Pages

- [Display Modes](index.md) — overview of all available modes
- [Test Pattern](testpattern.md) — fixed diagnostic pattern for validating screen configuration
<!-- snapshot-gallery:start -->
## Snapshots

### Color

<figure>
  <img src="../img/testwidgets/color-240x240_0001.png" alt="color-240x240 240x240" style="max-width:320px;width:100%;">
  <figcaption>240x240</figcaption>
</figure>

<figure>
  <img src="../img/testwidgets/color-320x240_0001.png" alt="color-320x240 320x240" style="max-width:320px;width:100%;">
  <figcaption>320x240</figcaption>
</figure>

### Mono

<figure>
  <img src="../img/testwidgets/mono-128x64_0001.png" alt="mono-128x64 128x64" style="max-width:320px;width:100%;">
  <figcaption>128x64</figcaption>
</figure>

<!-- snapshot-gallery:end -->
