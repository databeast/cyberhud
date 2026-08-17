// Package attract_plasma implements the "attract_plasma" display mode for the
// cyberHUD framework, rendering a lava-lamp-style plasma blob effect on any
// supported hardware panel.
//
// The mode computes a full-panel plasma pattern by evaluating a sum of
// overlapping sine-based color functions at each pixel position. The pattern
// morphs continuously over time, producing smooth organic blob shapes that
// shift and merge. A 256-step gradient palette cycles through hues at a
// configurable rate, while the spatial frequency of the pattern is controlled
// by the BlobScale policy field.
//
// Architecture:
//
//   - Catalog registration: mode definition registered via init() with ID
//     "attract_plasma", title, summary, and an "attract_plasma" command verb
//     for runtime policy queries and mutations.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a mutex-protected Policy struct with
//     speed, density, cycle rate, and blob scale fields, fully normalized on
//     every write to prevent undefined states.
//   - Change detection: RenderCacheKey returns a monotonically increasing frame
//     counter combined with the policy fingerprint, ensuring every frame is
//     unique for continuous animation.
//   - Full-panel sprite rendering: a single image.RGBA sprite covers the entire
//     panel, with each pixel computed from the plasma function. No font or text
//     rendering is involved.
//   - Widget sprite compositing: widgets.Compositor wraps the single plasma
//     sprite into the final frame output.
//   - Interface-based style dispatch: a style.NewRegistry[Snapshot] holding
//     per-resolution styles (color TFT, monochrome OLED, e-ink) with panel-
//     appropriate rendering paths.
//   - Panel hints and fitness evaluation: thread-safe panel hints storage with
//     a PostApply hook for fitness note generation.
//   - Action handler: returns nil (non-interactive ambient display mode).
//   - E-ink fallback: renders a single static plasma frame at time=0 with
//     stable RenderCacheKey for slow-refresh panels.
package attract_plasma
