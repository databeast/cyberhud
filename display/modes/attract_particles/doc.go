// Package attract_particles implements the "attract_particles" display mode for
// the cyberHUD framework, rendering drifting firefly-like particles with
// independent motion, color cycling, and edge-wrapping behavior.
//
// The mode produces a calm ambient glow effect by filling the panel with small
// filled circles that move in independent directions at varying speeds. Each
// particle's heading is perturbed every frame by a configurable drift randomness
// factor, and its hue cycles over time at a rate proportional to the glow
// intensity policy field. When a particle exits one panel edge it re-enters at
// a random position on the opposite edge, retaining its current color phase.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with ID "attract_particles", title, summary, and an "attract_particles"
//     command verb for runtime policy queries and mutations.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a mutex-protected source.Policy with
//     speed, density, drift, and glow fields, fully normalized on every write to
//     prevent undefined states.
//   - Change detection: RenderCacheKey returns a monotonically increasing frame
//     counter combined with the policy fingerprint, ensuring every frame is
//     unique for continuous animation.
//   - Particle system: each particle has an independent position, direction,
//     speed factor, and color phase. Particle count scales with panel area and
//     density: clamp(area/512 × density, 4, 200).
//   - Widget sprite compositing: widgets.Compositor assembles per-frame sprites
//     from the rendered particle image into the final frame output.
//   - Interface-based style dispatch: a style.NewRegistry[source.Snapshot,
//     source.Policy] holding per-resolution styles (color TFT, monochrome OLED,
//     e-ink) with panel-appropriate rendering.
//   - Panel hints and fitness evaluation: thread-safe panel hints storage.
//   - E-ink fallback: static particle scatter frame for slow-refresh panels.
package attract_particles
