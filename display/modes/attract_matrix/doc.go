// Package matrix implements the "matrix" display mode for the cyberHUD
// framework, rendering the iconic "digital rain" effect from The Matrix (1999)
// on any supported hardware panel.
//
// The mode tiles vertical columns of falling characters across the panel width,
// each scrolling at a randomized speed with a green-fade trail that decays from
// bright white-green through dark green to black. Characters are drawn from the
// matrix-code bitmap font (fonts.MatrixCodeID), which maps film-accurate
// katakana-like glyphs to standard ASCII code point positions.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with ID "matrix", title, summary, and a "matrix" command verb for runtime
//     policy queries and mutations.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a mutex-protected Policy struct with
//     speed range, trail length, density, and background toggle fields, fully
//     normalized on every write to prevent undefined states.
//   - Change detection: RenderCacheKey returns a monotonically
//     increasing frame counter combined with the policy fingerprint, ensuring
//     every frame is unique for continuous animation.
//   - Vertical marquee strips: each Rain_Column is a marquee.Strip configured
//     with vertical direction, a RandomSource character provider, and a
//     per-column color gradient array for the trail fade effect. Rain text
//     uses the matrix family, which includes a compact 10x10 variant for 128x64
//     panels and the original matrix-code face for larger layouts.
//   - Radial gradient background: optional background layer rendered via
//     gradient.Render with radial style, composited beneath the rain columns.
//   - Widget sprite compositing: widgets.Compositor assembles per-frame sprites
//     from all active strips and the optional background gradient into the
//     final frame output.
//   - Interface-based style dispatch: a style.NewRegistry[MatrixSnapshot]
//     holding per-resolution styles (color TFT, monochrome OLED, e-ink) with
//     panel-appropriate trail lengths and color palettes.
//   - Panel hints and fitness evaluation: thread-safe panel hints storage with
//     a PostApply hook that generates fitness notes via style.FitnessNotes when
//     the active style changes.
//   - Action handler: a Handler struct implementing action.Handler that returns
//     a no-op result (non-list mode with no interactive navigation).
//
// Developers building new display modes can study this package alongside the
// clock reference implementation for guidance on framework integration patterns
// including strip persistence between frames, elapsed-time tracking, and
// density-based column activation.
package attract_matrix
