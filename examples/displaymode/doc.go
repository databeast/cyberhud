// Package displaymode is a canonical skeleton template for building new display
// modes in the cyberHUD framework. It demonstrates every required integration
// point without implementing any real rendering logic, giving developers a
// minimal, compilable starting point they can copy and extend.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with typed option metadata, summaries, defaults, and allowed-value lists.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a thread-safe Policy struct with
//     field-level normalization and automatic fallback to safe defaults when
//     invalid values are supplied.
//   - Data-signature change detection: RenderCacheKey() returns a deterministic
//     change-detection string so the display runtime can skip unnecessary
//     re-renders between frames.
//   - Style interface implementation: 24 concrete styles implementing
//     style.Style[Snapshot] covering all known resolutions from
//     displayresolutions.md, including 7 degraded monochrome-fallback variants
//     for color-capable panels.
//   - Panel hints and fitness evaluation: panel hints storage with thread-safe
//     access, and fitness notes generation via the PostApply hook to report
//     style-panel compatibility at runtime.
//   - Action handler: Handler struct implementing action.Handler for logical UI
//     input processing (directional navigation and primary/secondary actions).
//
// Developers building new display modes should copy this package, replace the
// stub implementations with real logic, and adjust the style set to match the
// resolutions their hardware targets.
package displaymode
