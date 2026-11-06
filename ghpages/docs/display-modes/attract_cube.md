# Attract Cube

A translucent cube tumbles in 3D with illuminated edges and a subtle glow.

## Options

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `base_color` | string | `cyan` | Named accent (`cyan`, `green`, `emerald`, `amber`, `red`, `white`) or RGB hex (`#RRGGBB`). |

```sh
cyberhudctl display set 0 attract_cube base_color=amber
cyberhudctl display config main.0 'base_color=#ff6600'
```

Quote hex values in shells where `#` starts a comment. All cube faces, edges,
and glow derive from the base color. Very dark colors may be hard to see.

## Panel Compatibility

Fast-refresh color and grayscale displays show the animated cube; slow-refresh
color and grayscale panels show a static initial frame. Grayscale uses a
neutral high-contrast palette independent of `base_color`. Monochrome displays
use a one-bit wireframe; fast-refresh panels animate and slow-refresh panels
show a static initial frame. Monochrome also ignores `base_color`.

The mode registers surface styles for the standard landscape and portrait
display resolutions, including tiny LED matrices. Other region dimensions
render at their actual size using the panel's own capability.
