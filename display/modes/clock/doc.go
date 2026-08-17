// Package clock is the definitive reference implementation for the cyberHUD
// display-mode framework. It demonstrates the full range of framework capabilities
// while delivering a polished clock display on every supported hardware panel.
//
// # Architecture: Centralized Policy / Delegated Layout
//
// The clock mode uses a two-layer architecture that maximizes shared code
// across resolution-specific styles:
//
//   - Policy logic (what/when): buildClock owns all widget decisions — which
//     widgets to include, their theme, color tint, and suppression rules.
//     It reads Policy fields (ShowBorder, BorderColor, FGColor, ShowLED,
//     SecondsBar, etc.) and fully configures each widget before dispatching
//     to a style. No style function ever inspects policy fields to decide
//     whether a widget should appear or what color it should be.
//
//   - Layout logic (where/how big): Per-style Build functions handle only
//     positioning and sizing. They set Bounds on pre-configured widgets to
//     match their target panel dimensions. They never check policy fields or
//     make widget-inclusion decisions.
//
// This separation ensures that adding a new style requires zero
// widget-decision code. A new style only needs to provide layout logic for
// its target dimensions — all policy logic is inherited from buildClock.
//
// # Concrete Example: Border Frame
//
// The border frame widget illustrates this pattern clearly:
//
//  1. buildClock evaluates Policy.ShowBorder and resolves the border color
//     (via BorderColor/FGColor fallback logic). It configures the border
//     widget with theme "rounded" and the resolved ColorTint.
//  2. The style's Build function receives the pre-configured widget and sets
//     only Bounds to match its rendered content region (e.g., 240×240).
//  3. buildClock adds the widget to the Compositor via AddIf — the style
//     never touches the Compositor or checks ShowBorder.
//
// A new 480×800 style added tomorrow would inherit the full border feature
// by sizing Bounds to image.Rect(0, 0, 480, 800) — no conditional logic,
// no color resolution, no Compositor wiring required.
//
// # Guidance for New Display Modes and Styles
//
// Developers building new display modes or adding styles should follow this
// pattern to achieve maximum code sharing and minimal per-style complexity:
//
//  1. Place all widget-inclusion logic in the shared base function (the
//     equivalent of buildClock). Read policy fields, resolve colors,
//     configure widgets with theme and tint here.
//  2. Keep style Build functions purely mechanical: compute text positions,
//     vertical centering, and content bounds for the target resolution, then
//     set Bounds on any pre-configured widgets.
//  3. Use Compositor.AddIf in the shared base — never in styles. This keeps
//     the conditional-inclusion decision centralized and testable in one
//     place.
//  4. New styles inherit all widget features automatically. The only
//     per-style code is layout arithmetic for the new target dimensions.
//
// # Framework Integration Points
//
//   - Catalog registration: mode and command definitions registered via init()
//     with typed option metadata, summaries, defaults, and allowed-value lists.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a Policy struct with full
//     normalization (whitespace trimming, lowercasing, invalid-value replacement)
//     ensuring the mode never enters an undefined state.
//   - Data-signature change detection: RenderCacheKey returns a precision-adaptive
//     fingerprint (second or minute granularity) so the display runtime skips
//     unnecessary re-renders on battery-powered hardware.
//   - Interface-based style dispatch: a style.NewRegistry[Snapshot] holding
//     resolution-specific styles with font selection via tier catalog for each
//     target panel.
//   - Panel hints and fitness evaluation: thread-safe panel hints storage with
//     a PostApply hook that generates fitness notes via style.FitnessNotes when
//     the active style changes.
//   - Action handler: a Handler struct implementing action.Handler for logical
//     UI input processing (Left/Right cycle styles, Primary toggles border).
//   - ViewData conversion: helper functions mapping between mode-internal ViewData
//     and shared style.ViewData with an empty-items guard and per-row font ID
//     assignment.
//   - Widget sprite compositing: border frame, LED seconds indicator, progress
//     bar (horizontal and pie), and sparkline daybar sprites conditionally
//     appended to ViewData via Compositor based on policy flags and panel
//     dimension thresholds.
package clock
