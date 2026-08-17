// Package systemd implements the systemd boot progress display mode for cyberHUD.
// It renders real-time boot target transition status across 44 registered styles
// spanning four display categories (color, mono, e-ink, grayscale-fast), with
// resolution-specific rendering optimized for each panel geometry.
//
// Architecture:
//
//  1. Registry-based dispatch: a style.NewRegistry[Snapshot] holding 44
//     resolution-specific styles, with Lookup/Default fallback replacing the
//     former switch-case dispatch in BuildView.
//  2. LayoutBridge: each resolution style constructs a layout.NewLayoutBridge
//     to compute the content area, border inset, and available pixel dimensions
//     for consistent sprite placement across all panel geometries.
//  3. Compositor: widgets.NewCompositor assembles gradient backgrounds, text
//     labels, and status sprites into a unified sprite list, applying per-widget
//     suppression rules based on panel capabilities.
//  4. SuppressOnEink: e-ink styles configure the Compositor with
//     widgets.SuppressOnEink to exclude animated widgets (gradients, progress
//     bars) from rendered output, preserving static-friendly display on e-paper.
//  5. RenderCacheKey: deterministic change-detection fingerprint encoding policy
//     fields (style, show_border, color_accent) and boot state (targets, boot
//     completion), enabling the display runtime to skip re-renders when nothing
//     has changed.
//  6. BestFontFor: style.BestFontFor selects the optimal font for each panel's
//     pixel dimensions and row/column requirements, with fonts.Default fallback
//     when no registered font satisfies the constraints.
//  7. GradientWidget: gradient.New produces a full-panel linear gradient sprite
//     used as a background fill behind boot status text, transitioning from
//     black to the resolved accent color at a position proportional to boot
//     fraction.
//  8. ColorAccent: configurable accent theming via sharedcolor.ResolveAccent,
//     applied to gradient stop colors and text foreground across color and
//     grayscale-fast display categories, with automatic green override on boot
//     completion.
//
// Developers building new display modes can use this package alongside the
// clock and GPIO modes as a structural template for framework integration.
package systemd
