// Package attract_geometric implements the "attract_geometric" display mode for
// the cyberHUD framework, rendering clusters of semi-transparent rotating
// rectangles with sinusoidal fade cycles, distance-based glow effects, and
// drifting pseudocode text fragments.
//
// The effect is a 1:1 port of the website's geometric background animation
// (website/src/animations/geometric/) into the CyberHUD display system. All
// mathematical formulas, constants, algorithms, draw order, and visual behavior
// are preserved verbatim from the TypeScript source. Three adaptations are made:
// proportional scaling for panels smaller than 240px shortest dimension,
// CyberHUD integration via the standard catalog/policy/instance pattern, and
// image.RGBA rendering replacing the Canvas 2D API.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with ID "attract_geometric", title, summary, and an "attract_geometric"
//     command verb for runtime policy queries and mutations.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a mutex-protected Policy struct with
//     Speed, Density, GlowIntensity, and FragmentRate fields, fully normalized
//     on every write to prevent undefined states.
//   - Change detection: RenderCacheKey returns a monotonically increasing frame
//     counter combined with the policy fingerprint, ensuring every animated
//     frame produces a unique key.
//   - Cluster generation: grid-distributed clusters of 8–20 rectangles with
//     exponential size decay, anchor-point grid positioning, and computed
//     bounding radii for glow distance attenuation.
//   - Sinusoidal fade cycles: each rectangle oscillates in opacity via a
//     per-rectangle phase offset and cycle duration, producing organic shimmer.
//   - Distance-based glow: the largest rectangle per cluster receives a glow
//     layer with elevated lightness and blur spread, attenuated by distance
//     from cluster center.
//   - Central zone opacity cap: elements within the middle 60% of panel width
//     are capped at 0.4 opacity to keep text-overlay areas readable.
//   - Pseudocode text fragments: drifting monospace text snippets spawned near
//     active clusters with fade-in/hold/fade-out lifecycle.
//   - Adaptive performance scaling: a rolling frame-time window triggers
//     cluster reduction or restoration to maintain target frame rate.
//   - Widget sprite compositing: widgets.Compositor assembles per-frame sprites
//     from all active rectangles, glow layers, and fragments into the final
//     frame output.
//   - Interface-based style dispatch: a style.NewRegistry holding per-resolution
//     styles (color TFT, monochrome OLED, slow-refresh) with panel-appropriate
//     rendering.
//   - Panel hints and fitness evaluation: thread-safe panel hints storage.
//   - Action handler: nil (non-interactive ambient display mode).
//
// On slow-refresh panels, a single static frame is produced at time=0 with all
// clusters at full fade-in and no fragment text.
package attract_geometric
