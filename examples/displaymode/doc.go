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
//   - Style interface implementation: 24 per-resolution concrete styles
//     implementing style.Style[Snapshot] covering all known resolutions,
//     including 7 degraded monochrome-fallback variants for color-capable
//     panels, plus one graphics-example style (see below).
//   - Panel hints and fitness evaluation: panel hints storage with thread-safe
//     access, and fitness notes generation via the PostApply hook to report
//     style-panel compatibility at runtime.
//   - Action handler: Handler struct implementing action.Handler for logical UI
//     input processing (directional navigation and primary/secondary actions).
//
// Developers building new display modes should copy this package, replace the
// stub implementations with real logic, and adjust the style set to match the
// resolutions their hardware targets.
//
// # Text vs graphics: the two layout planes
//
// The 24 per-resolution styles here are text-only (they emit ViewData.Items).
// The moment a mode draws GRAPHICS — a bar, a gauge, an icon — next to text,
// there is a non-obvious rule that, ignored, produces views that render as an
// unreadable clump. A ViewData is drawn through TWO independent coordinate
// systems:
//
//   - The TTY / row text plane (ViewData.Items): the renderer owns vertical
//     placement and flows rows by its own row height. The style does not
//     control a row's pixel Y.
//   - The pixel-anchored sprite plane (ViewData.Sprites): absolute coordinates
//     the style sets exactly.
//
// Do NOT mix them by putting a caption in Items and a bar in Sprites and
// expecting them to line up — the row height the renderer uses and the one the
// style assumes can differ (notably on panels that report no glyph metrics),
// and the two planes drift apart. Commit a text+graphics view entirely to the
// pixel plane: render labels as sprites with the textlabel widget and leave
// Items empty. The GraphicsExampleStyle (style_graphics_example.go) is the
// fully-worked, heavily-commented reference for this, including getting real
// font metrics from the StyleContext and distributing rows to fill the panel.
package displaymode
