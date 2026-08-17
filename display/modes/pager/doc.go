// Package pager implements the "pager" display mode for the cyberHUD
// framework, tailing a configured data source (file, named pipe, or Unix
// domain socket) and rendering its text content on a display surface.
//
// The mode adapts its presentation strategy to the display hardware:
//
//   - Fast displays (OLED, TFT): smooth pixel-by-pixel upward scrolling at
//     approximately 30 fps, with velocity adaptation when buffered lines exceed
//     visible capacity.
//   - Slow displays (e-ink): full-page transitions using fade-out/fade-in
//     animation timed to a reading cadence derived from visible row count and a
//     configurable per-line reading time.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with ID "pager", title, summary, and a "pager" command verb for runtime
//     policy queries and mutations.
//   - Policy definition and normalization: a mutex-protected Policy struct with
//     Source, ScrollSpeed, MaxLines, ScanMS, FadeOutMS, FadeInMS, LineTimeMS,
//     MaxWaitS, Font, and Style fields.
//   - Data ingestion: a tailReader goroutine tails the configured data source,
//     splitting on newlines into a thread-safe ring buffer.
//   - Rendering strategies: Smooth_Scroll for fast displays and Page_Transition
//     for slow displays, selected automatically based on TextHints capabilities.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query/mutation responses.
//   - Change detection: RenderCacheKey incorporates buffer sequence, scroll
//     offset (or page phase), and policy fingerprint to enable skip-redraw
//     optimization.
package pager
