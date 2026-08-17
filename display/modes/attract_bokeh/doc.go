// Package attract_bokeh implements the "attract_bokeh" display mode for the
// cyberHUD framework, rendering soft out-of-focus light circles that drift
// gently across the panel with varying sizes, opacities, and positions.
//
// The effect mimics bokeh photography — large translucent circles with radial
// gradient edges (full opacity at center, zero at edge) that overlap and blend
// to produce a dreamy, ambient aesthetic. Circles vary in radius between 2% and
// 15% of the panel's shortest dimension based on the SizeVariance policy field.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with ID "attract_bokeh", title, summary, and an "attract_bokeh" command
//     verb for runtime policy queries and mutations.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a mutex-protected Policy struct with
//     Speed, Density, SizeVariance, and Saturation fields, fully normalized on
//     every write to prevent undefined states.
//   - Change detection: RenderCacheKey returns a monotonically increasing frame
//     counter combined with the policy fingerprint, ensuring every animated
//     frame produces a unique key.
//   - Radial gradient circles: each circle is rendered with a soft edge gradient
//     from full opacity at center to zero at edge, with peak opacity between
//     0.15 and 0.6 per circle.
//   - Widget sprite compositing: widgets.Compositor assembles per-frame sprites
//     from all active circles into the final frame output.
//   - Interface-based style dispatch: a style.NewRegistry[BokehSnapshot] holding
//     per-resolution styles (color TFT, monochrome OLED, e-ink) with
//     panel-appropriate rendering.
//   - Panel hints and fitness evaluation: thread-safe panel hints storage.
//   - Action handler: nil (non-interactive ambient display mode).
//
// On monochrome panels, circles use luminance-only values without color.
// On e-ink panels, a static frozen bokeh scatter frame is produced.
package attract_bokeh
