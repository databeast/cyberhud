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
package gpio
