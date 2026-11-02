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
//
// # User Experience
//
// Promise (verbatim from catalog registration in cmd.go init):
//
//	Name:        "Waveform"
//	Description: "Oscilloscope-style animated waveform traces with phosphor persistence."
//
// Class: attract/decorative
//
// The core fantasy: a window onto an UNKNOWN signal. The observer has no
// legend, no units, no idea what is being measured — only the trace — and the
// trace must carry the whole sense that "something is out there, being
// measured." The single most important property is therefore that it never
// looks twice-the-same: not a repeating tile within a frame, not a recurring
// pattern across frames. A guessable pattern collapses the illusion into "a
// computer drew a wave." Unpredictability is the feature.
//
// Implementation note (why it does not repeat): each trace is the sum of
// several sine partials at mutually-incommensurate (irrational-ratio)
// frequencies, so the composite is quasi-periodic — its true period is
// astronomically long and it never visibly repeats across the panel. A
// slow-wandering macro-envelope plus sparse transient "events" (bursts) give it
// dynamics (surges, lulls, jolts), and a cross-frame phosphor buffer carries
// the fading trail of recent sweeps.
//
// Invariants (must hold — a usability defect if violated):
//   - The frame reads as a live oscilloscope signal: a restless trace sweeping
//     across the panel, not noise, not blank, not a single flat line.
//   - The trace does NOT visibly repeat within the frame — no row of identical
//     humps. One broad, irregular signal spans the sweep.
//   - At least two distinct traces, visually distinguishable (on color panels
//     by punchy phosphor hues — green/cyan/amber/magenta; on mono/grayscale by
//     shape and brightness).
//   - Traces use a meaningful fraction of the panel's amplitude and run the
//     full sweep, so it looks alive rather than a thin line hugging the center.
//   - A phosphor glow/bloom is visible around the trace on fast-refresh panels
//     (the beam reads as lit, not a hard vector stroke).
//
// Intentional (do NOT report as defects):
//   - Overlapping / crossing traces, dense tangles, bursts, traces touching the
//     panel edges. This is a live scope, not a tidy chart.
//   - Asymmetry and irregularity along the sweep — that IS the signal.
//   - On grayscale/mono panels the per-trace color coding collapses to
//     shades/white — color differentiation is not promised there.
//   - Lack of visible text, labels, grid or graticule — the mode is a bare,
//     unexplained window onto the signal by design.
//
// Not verifiable from a single still frame:
//   - The non-repetition ACROSS TIME and the moving phosphor-persistence trail.
//     A still can show one instant is not a repeating tile and that a static
//     glow exists; it cannot show the beam sweeping or the trail fading. The
//     living quality only fully exists on the device at ~30fps.
//
// Per-resolution notes:
//   - Landscape panels (w>h) such as 320x240 sweep horizontally; portrait
//     panels sweep vertically ("auto" direction keys off aspect ratio).
//   - Ultra-low-res panels (≤16x8 / ≤8x16) degrade to a single calmer trace.
//   - Slow-refresh (e-ink/e-paper) panels render a single bold frozen slice of
//     the signal with no animation or glow by design.
package attract_waveform
