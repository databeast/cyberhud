# Attract Modes

Attract modes are non-interactive ambient visual effects designed to activate automatically during idle periods. When no user interaction occurs for a configurable duration (the panel's idle timeout), the display catalog transitions to one of the attract modes, producing a continuously animated screensaver-like experience. Once input resumes — a button press, joystick movement, or external command — the panel returns to its previously active display mode.

Each attract mode runs independently of system data feeds and requires no network connectivity or sensor input. They exist purely to keep the panel visually active during downtime, prevent static burn-in on OLED panels, and provide an ambient aesthetic when the device is unattended. The idle timeout that triggers attract mode activation is configured via the `IdleTimeout` policy field on the display catalog.

## Summary

| Mode ID | Title | Description |
|---------|-------|-------------|
| [attract_bokeh](attract_bokeh.md) | Bokeh | Soft out-of-focus light circles drifting with radial gradients |
| [attract_geometric](attract_geometric.md) | Geometric | Pulsing and rotating geometric patterns with neon edges |
| [attract_hacking](attract_hacking.md) | Hacking | Dramatic cyberpunk intrusion sequence with faux terminal output |
| [attract_matrix](attract_matrix.md) | Matrix | Cascading character rain effect inspired by the Matrix |
| [attract_particles](attract_particles.md) | Particles | Drifting firefly-like particles with independent motion and color cycling |
| [attract_plasma](attract_plasma.md) | Plasma | Lava-lamp-style plasma blobs with morphing sine patterns and gradient palette |
| [attract_shapes](attract_shapes.md) | Shapes | Pulsing and rotating geometric polygons in dynamic arrangement |
| [attract_starfield](attract_starfield.md) | Starfield | Perspective depth starfield with stars emanating from central vanishing point |
| [attract_waveform](attract_waveform.md) | Waveform | Oscilloscope-style waveform traces with morphing wave shapes |
