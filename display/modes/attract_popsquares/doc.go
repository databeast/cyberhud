// Package attract_popsquares implements the "attract_popsquares" display mode
// for the cyberHUD framework: a full-bleed grid of flat colour cells that sit
// at a calm base colour while a sparse scatter of cells gently pulses brighter
// and eases back, evoking a supercomputer idling. Inspired by the XScreenSaver
// "popsquares" hack.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with typed option metadata (speed, grid_scale, pop_rate, warmth).
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a mutex-protected Policy struct,
//     normalized on every write to clamp out-of-range values.
//   - Change detection: RenderCacheKey returns a per-frame key on animated
//     panels and a stable (policy-only) key on e-ink panels.
//   - Widget sprite compositing: the render pass composes a single full-panel
//     grid sprite per frame via widgets.Compositor.
//   - Interface-based style dispatch: a style.NewRegistry[Snapshot, Policy]
//     holding per-resolution parameterized styles across all six capability
//     tiers (the attract-mode exception in CODINGSTYLE.md — one code path,
//     many resolutions).
//   - Panel hints: thread-safe panel hints storage for fitness evaluation.
//   - E-ink fallback: a single static calm frame for slow-refresh panels.
//
// The mode is purely sprite-based — no text, no fonts. All rendering produces
// image.RGBA cells drawn block-by-block.
//
// # User Experience
//
// Promise (verbatim from catalog registration in cmd.go init):
//
//	Name:        "Pop Squares"
//	Description: "A calm grid of color cells that gently pulse and fade — a supercomputer idling."
//
// Class: attract/decorative
//
// The fantasy: a vast machine AT REST. Banks of status cells, something always
// quietly ticking over somewhere, no urgency, no chaos — the ambient pulse of a
// system that is alive and waiting. You could watch it all day and it would
// never alarm you, but it never goes dead either. The craft is RESTRAINT: the
// temptation to make it flashy must be resisted, because flashy breaks the
// "idling" illusion. Calm is the product.
//
// Invariants (must hold — a usability defect if violated):
//   - A full-bleed grid of flat colour cells fills the panel edge to edge
//     (no borders, no gaps, no background showing through).
//   - The field is COHERENT: neighbouring cells are related in colour (a slow
//     drifting gradient), not random confetti — it reads as one system.
//   - Most cells rest at a calm, dim base colour; only a sparse minority are
//     lit brighter at any instant (the idle twinkle, not a light show).
//   - The default palette is cool/nominal (blue-teal-green) — the colour of a
//     machine at rest, not an alarm board.
//   - The twinkle does not fall into a visible repeating pattern: pulses are
//     scattered per-cell, not synchronised rows or marching diagonals.
//
// Intentional (do NOT report as defects):
//   - The calm, low-energy feel. Sparse activity is the point; "it's not very
//     busy" is success, not a defect.
//   - On grayscale/mono panels the hue collapses to brightness — colour
//     differentiation is not promised there.
//   - No text, labels, numbers, or icons — the mode is a bare ambient field.
//   - Partial cells at the right/bottom edge where the grid does not divide the
//     panel evenly.
//
// Not verifiable from a single still frame:
//   - The pulse-and-fade MOTION and the fact that the active cells change over
//     time without repeating. A still shows one instant: a calm field with a
//     few cells currently brighter than their neighbours. The living, breathing
//     quality only exists on the device at ~30fps.
//
// Per-resolution notes:
//   - Cell size scales with the panel's short edge and the grid_scale knob, so
//     the bank of squares stays legible from tiny OLEDs to large TFTs.
//   - Ultra-low-res panels (≤16x8 / ≤8x16) get fewer, slower pops so the grid
//     reads as calm rather than noisy.
//   - Slow-refresh (e-ink/e-paper) panels render a single frozen calm frame
//     with a handful of cells lit, no animation.
package attract_popsquares
