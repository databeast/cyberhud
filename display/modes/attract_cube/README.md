# Cube attract mode

Configure `base_color` with a named accent (`cyan`, `green`, `emerald`, `amber`,
`red`, or `white`) or a six-digit RGB color such as `#ff6600`. The default is
`cyan`. Quote hex values in shells where `#` introduces a comment.

```sh
cyberhudctl display set 0 attract_cube base_color=amber
cyberhudctl display config main.0 'base_color=#ff6600'
```

The cube animates on fast-refresh color and grayscale displays; slow-refresh
displays show a static initial frame. Grayscale uses a neutral high-contrast
palette, while monochrome uses a solid one-bit wireframe; neither depends on
`base_color`. Monochrome fast-refresh panels animate and slow-refresh panels
show a static initial frame. To render snapshots through the existing pipeline:

```sh
go test ./display/modes/attract_cube -run TestCubePNGSnapshots -count=1
```

The generated PNGs and gallery fragment are written to
`display/modes/attract_cube/snapshots/`. They include every registered
color, grayscale, and monochrome surface style in both refresh categories,
plus several animation moments and alternate colors.
