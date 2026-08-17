// Package wifi implements a display mode for the cyberHUD system that renders
// real-time wireless network status to TFT panels. It reads interface state from
// the Linux kernel (/proc/net/wireless) and presents connection details through a
// visually rich layout featuring a WiFi icon sprite, color-coded signal bars, a
// link quality progress bar, and detail rows for SSID, frequency, channel, IP
// address, and interface name.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with typed option metadata (style, show_border, fgcolor, signal_display),
//     summaries, defaults, and allowed-value lists.
//   - Command handling: CLI verb processing with per-key validation, atomic
//     policy mutation, and formatted query responses via cmdutil.CmdHandler.
//   - Policy definition and normalization: a 7-field Policy struct with full
//     normalization (whitespace trimming, lowercasing, invalid-value replacement)
//     ensuring the mode never enters an undefined state.
//   - Data-signature change detection: RenderCacheKey returns a fingerprint
//     incorporating connection state, signal metrics, and all policy fields so
//     the display runtime skips unnecessary re-renders on battery-powered hardware.
//   - Interface-based style dispatch: a style.NewRegistry[WifiSnapshot] holding
//     resolution-specific styles (color-128x128, color-240x240, color-320x240,
//     grayscale-fast-128x128, grayscale-fast-240x240) with the first registered style serving
//     as the registry default.
//   - Panel hints and fitness evaluation: thread-safe panel hints storage with
//     a PostApply hook that generates fitness notes via style.FitnessNotes when
//     the active style changes.
//   - Action handler: a Handler struct implementing action.Handler for logical
//     UI input processing (Left/Right cycle styles, Secondary forces refresh,
//     Primary navigates to menu).
//   - Platform-aware state gathering: WiFi state reading is isolated behind
//     build-tag-constrained functions (linux vs !linux) returning a WifiSnapshot
//     struct, enabling graceful degradation on non-Linux systems.
//   - Widget sprite compositing: signal bar sprite (4 vertical bars with quality-
//     mapped colors), WiFi icon sprite (scaled 8×8 glyph), and horizontal progress
//     bar (link quality) conditionally appended to ViewData based on connection state.
//
// The primary optimization target is the Waveshare 1.44" LCD HAT (128×128 ST7735S),
// though the style registry supports multiple resolutions via the same dispatch
// mechanism used by the clock and dashboard modes.
package wifi
