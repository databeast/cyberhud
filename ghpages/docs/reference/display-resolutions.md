# Display Resolutions

CyberHUD supports a range of display panels spanning monochrome OLEDs, color LCDs, and e-ink displays. This reference lists all supported display resolutions with their characteristics.

## Landscape (Native) Display Targets

| Resolution | Aspect Ratio | Color Depth | Pixel Format | Refresh Model |
|-----------|--------------|-------------|--------------|---------------|
| 128×32 | 4:1 | 1-bit | Mono packed | Deferred |
| 128×64 | 2:1 | 1-bit | Mono packed | Deferred |
| 128×128 | 1:1 | 1-bit | Mono packed | Deferred |
| 128×128 | 1:1 | RGB565 | RGB565 | Continuous |
| 160×80 | 2:1 | RGB565 | RGB565 | Continuous |
| 160×128 | 5:4 | RGB565 | RGB565 | Continuous |
| 240×135 | 16:9 | RGB565 | RGB565 | Continuous |
| 240×240 | 1:1 | RGB565 | RGB565 | Continuous |
| 320×240 | 4:3 | RGB565 | RGB565 | Continuous |
| 480×320 | 3:2 | RGB565 / RGB666 | RGB565 / RGB666 | Continuous |
| 800×480 | 5:3 | RGB888 / RGB666 | RGB888 / RGB666 | Continuous |
| 122×250 | variable | 1-bit | Mono packed | Deferred |
| 176×264 | variable | 1-bit | Mono packed | Deferred |
| 200×200 | 1:1 | 1-bit | Mono packed | Deferred |
| 212×104 | ~2:1 | 1-bit | Mono packed | Deferred |
| 296×128 | ~2.3:1 | 1-bit / 2-bit | Mono packed | Deferred |
| 400×300 | 4:3 | 1-bit / 2-bit | Mono packed | Deferred |
| 480×800 | 5:3 | 1-bit / 2-bit | Mono packed | Deferred |
| 800×480 | 5:3 | 1-bit / 2-bit | Mono packed | Deferred |

!!! note "Square displays"
    Square displays (128×128, 240×240, 200×200) are rotation-invariant — they produce identical width×height dimensions when rotated 90° or 270° and do not require separate portrait layout designs.

## Rotated (Portrait) Display Targets

When a panel is configured with rotation (via `rotate` or `orientation` config fields), the effective resolution swaps width and height. These are the resulting portrait targets:

| Resolution | Aspect Ratio | Color Depth | Pixel Format | Refresh Model | Source |
|-----------|--------------|-------------|--------------|---------------|--------|
| 32×128 | 1:4 | 1-bit | Mono packed | Deferred | 128×32 |
| 64×128 | 1:2 | 1-bit | Mono packed | Deferred | 128×64 |
| 80×160 | 1:2 | RGB565 | RGB565 | Continuous | 160×80 |
| 128×160 | 4:5 | RGB565 | RGB565 | Continuous | 160×128 |
| 104×212 | ~1:2 | 1-bit | Mono packed | Deferred | 212×104 |
| 250×122 | variable | 1-bit | Mono packed | Deferred | 122×250 |
| 135×240 | 9:16 | RGB565 | RGB565 | Continuous | 240×135 |
| 128×296 | ~1:2.3 | 1-bit / 2-bit | Mono packed | Deferred | 296×128 |
| 264×176 | variable | 1-bit | Mono packed | Deferred | 176×264 |
| 240×320 | 3:4 | RGB565 | RGB565 | Continuous | 320×240 |
| 300×400 | 3:4 | 1-bit / 2-bit | Mono packed | Deferred | 400×300 |
| 320×480 | 2:3 | RGB565 / RGB666 | RGB565 / RGB666 | Continuous | 480×320 |
| 480×800 | 3:5 | RGB888 / RGB666 | RGB888 / RGB666 | Continuous | 800×480 |

## Column Descriptions

| Column | Description |
|--------|-------------|
| **Resolution** | Width × Height in pixels |
| **Aspect Ratio** | Display aspect ratio (marked "variable" for non-standard e-ink panels) |
| **Color Depth** | Bits per pixel — 1-bit (monochrome), 2-bit (grayscale e-ink), RGB565 (16-bit color), RGB666 (18-bit), or RGB888 (24-bit) |
| **Pixel Format** | Internal pixel buffer format used by the rendering pipeline |
| **Refresh Model** | **Continuous** for LCD/OLED (frame-by-frame updates) or **Deferred** for e-ink (full-panel refresh) |
| **Source** | (Portrait table only) The native landscape resolution this target is derived from |
