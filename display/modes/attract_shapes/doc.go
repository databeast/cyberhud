// Package attract_shapes implements the "attract_shapes" display mode for
// the cyberHUD framework, rendering pulsing and rotating geometric shapes
// (regular polygons) across the panel for a dynamic abstract art effect.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with ID "attract_shapes", title, summary, and command verb for runtime
//     policy queries and mutations.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a mutex-protected Policy struct with
//     speed, density, shape_count, pulse_rate, and complexity fields, fully
//     normalized on every write to prevent undefined states.
//   - Change detection: RenderCacheKey returns a monotonically increasing frame
//     counter combined with the policy fingerprint, ensuring every frame is
//     unique for continuous animation.
//   - Widget sprite compositing: widgets.Compositor assembles per-frame sprites
//     into the final rendered output.
//   - Interface-based style dispatch: a style.NewRegistry[Snapshot] holding
//     per-resolution styles (color TFT, monochrome OLED, e-ink) covering all
//     six capability tiers.
//   - Panel hints: thread-safe panel hints storage for fitness evaluation.
//   - E-ink fallback: static decorative shapes frame for slow-refresh panels.
//
// The shapes mode is purely sprite-based — no text, no fonts. All rendering
// produces image.RGBA sprites with regular polygons drawn using vertex computation.
package attract_shapes
