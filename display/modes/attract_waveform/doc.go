// Package attract_waveform implements the "attract_waveform" display mode for
// the cyberHUD framework, rendering animated oscilloscope-style waveform traces
// across the panel. Traces morph between sine, sawtooth, and harmonic wave
// shapes with configurable amplitude, speed, and phosphor persistence trail.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with ID "attract_waveform", title, summary, and command verb for runtime
//     policy queries and mutations.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a mutex-protected Policy struct with
//     speed, density, amplitude, traces, and persistence fields, fully
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
//   - E-ink fallback: static decorative waveform frame for slow-refresh panels.
//
// The waveform mode is purely sprite-based — no text, no fonts. All rendering
// produces image.RGBA sprites with waveforms drawn pixel-by-pixel.
package attract_waveform
