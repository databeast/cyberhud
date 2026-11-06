// Package gpio implements the GPIO pin monitoring display mode for cyberHUD.
// It renders real-time GPIO pin states across 49 registered styles spanning
// four display categories (color, mono, e-ink, grayscale-fast) plus five
// resolution-independent styles that work on any panel size.
//
// Architecture:
//
//  1. Registry-based dispatch: a style.NewRegistry[GpioSnapshot] holding 44
//     resolution-specific styles plus 5 general-purpose styles, with Lookup/Default
//     fallback replacing the former switch-case dispatch in BuildView.
//  2. LayoutBridge: each resolution style constructs a layout.NewLayoutBridge
//     to compute the content area, border inset, and available pixel dimensions
//     for consistent sprite placement across all panel geometries.
//  3. Compositor: widgets.NewCompositor assembles LED indicators, text labels,
//     and sparkline sprites into a unified sprite list, applying per-widget
//     suppression rules based on panel capabilities.
//  4. SuppressOnEink: e-ink styles configure the Compositor with
//     widgets.SuppressOnEink to exclude non-eink-safe widgets (sparklines)
//     from rendered output, preserving static-friendly display on e-paper.
//  5. RenderCacheKey: deterministic change-detection fingerprint encoding both
//     policy fields and pin states, enabling the display runtime to skip
//     re-renders when nothing has changed.
//  6. TierCatalog: tier catalog metrics (via hints.Catalog) provide layout
//     dimensions (GlyphAdvance, RowHeight) for adaptive rendering, with the
//     registry wrapper handling font resolution through tier declarations.
//  7. PersistentWidgets: package-level LED and sparkline widget instances
//     reconfigured each frame rather than reallocated, reducing GC pressure
//     on memory-constrained hardware.
//  8. FGColor: configurable foreground color theming via sharedcolor.ResolveAccent
//     and sharedcolor.Dim, applied to HIGH/LOW pin state foreground colors
//     across color and grayscale-fast display categories.
//
// Developers building new display modes can use this package alongside the
// clock mode as a structural template for framework integration.
//
// # User Experience
//
// Promise (what the mode delivers):
//
//	Title:   "GPIO"
//	Pitch:   "Your headless machine's GPIO, alive at a glance — every active bus
//	          and driven pin, nothing dead, and the moment one changes you see it."
//
// Note: the catalog Summary in cmd.go still reads "GPIO pin state list on the
// main display or compact GPIO counters on secondary displays." That predates
// this redesign and should be updated to match the liveness framing; the drift
// is a known follow-up, not the mode's actual behavior.
//
// Class: informational
//
// Context — why this mode exists: the deployment target is a HEADLESS machine
// whose only video feedback is this small panel. Nobody is SSH'd in; the screen
// is the single window into what the GPIO header is actually doing. So the mode
// answers one question at arm's length: "is the hardware doing what I deployed
// it to do?" Steadiness reads as "fine"; a change you notice is the event that
// makes you walk over.
//
// What it shows — the LIVE FUNCTIONS of the header, one row each:
//   - A STANDALONE digital pin (output or input) appears individually: its
//     identity (number) + mode (OUT/IN) + level (HI/LO). For these lone lines
//     the level is the whole point.
//   - A BUS / peripheral (SPI0, I2C1, UART0, PWM0, …) spanning several pins
//     COLLAPSES to a single row showing the bus name + "up". Its member pins'
//     individual levels are noise; its presence is the signal. The grouping is
//     derived live from the driver's reported pin function (periph Func()), not
//     a hardcoded table.
//   - UNCONFIGURED/dead pins produce NO row. "Nothing" is not worth screen space
//     on a headless box's only display.
//
// Row order (most actionable first): standalone outputs, then standalone inputs,
// then buses.
//
// Invariants (must hold — a usability defect if violated):
//   - Each row is legible: identity and state both readable, text not clipped at
//     edges, rows not overlapping or desynced from their indicator dot.
//   - A bus appears exactly once (collapsed), never as its individual member
//     pins. A standalone pin appears individually, never merged.
//   - Dead pins are absent. A screen full of idle/unconfigured rows is a defect.
//   - Rows fill the available height rather than clumping at the top.
//   - Indicator dots render as solid marks, not tofu.
//
// Intentional (do NOT report as defects):
//   - A calm, steady screen is the desired resting state — "everything's fine".
//     There is deliberately no animation, flash, or motion; the screen simply
//     reflects current truth every refresh and the viewer's memory supplies the
//     delta.
//   - On e-ink the indicator dots are dropped (not eink-safe); the text row
//     carries the full identity+state, so this loses nothing.
//   - Color-coding (accent for live/HIGH, dimmed for LOW) is policy- and
//     capability-driven; a mono panel showing no color is expected.
//   - On a tiny strip the rows get fewer and tighter, but stay per-function —
//     never collapsed into an aggregate count.
//
// Not verifiable from a single still frame:
//   - "Change as signal" is inherently cross-frame: a still cannot show that a
//     row looks different than it did a moment ago.
//
// Per-resolution notes:
//   - Roomy panels: full rows (dot + "17 OUT HI" / "SPI0 up"), distributed down
//     the panel.
//   - Narrow/portrait and tiny panels: same per-function rows, truncated text
//     and fewer rows; identity is preserved, never replaced by a tally.
package gpio
