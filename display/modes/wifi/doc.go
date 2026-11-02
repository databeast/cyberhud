// Package wifi implements a display mode for the cyberHUD system that renders
// real-time wireless network status to TFT panels. It reads interface state from
// the Linux kernel via NL80211 (netlink) and presents connection details through
// a visually rich layout featuring a WiFi icon sprite, color-coded signal bars, a
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
//     build-tag-constrained functions (linux vs !linux) returning a WifiData
//     struct, enabling graceful degradation on non-Linux systems.
//   - Widget sprite compositing: signal bar sprite (4 vertical bars with quality-
//     mapped colors), WiFi icon sprite (scaled 8×8 glyph), and horizontal progress
//     bar (link quality) conditionally appended to ViewData based on connection state.
//
// The primary optimization target is the Waveshare 1.44" LCD HAT (128×128 ST7735S),
// though the style registry supports every resolution/capability pairing via the
// same dispatch mechanism used by the clock and dashboard modes.
//
// # User Experience
//
// Promise (verbatim from catalog registration in instance.go init):
//
//	Name:        "WiFi"
//	Description: "Real-time wireless network status with signal bars, quality
//	             bar, and connection details."
//
// Class: informational
//
// Invariants (must hold — a usability defect if violated):
//   - When connected, the SSID is the most prominent row and reads without being
//     clipped (truncated with an ellipsis when it genuinely cannot fit).
//   - Signal strength is always conveyed: the signal-bars sprite in "bars" mode,
//     or a numeric label ("NN%" / "-NNdBm") in percentage/dbm mode.
//   - The signal-bars sprite stands alone in "bars" mode; no redundant
//     "Signal: N/4" text duplicates it.
//   - Frequency and channel, when both shown, occupy a single consolidated row
//     ("Ch 36 (5.2GHz)") rather than two separate rows.
//   - The show_frequency, show_channel, and show_interface toggles actually
//     add/remove their rows across every fast-engine style, not only the few
//     bespoke panels.
//   - Disconnected renders "No Network" (warning color) and Unavailable renders
//     "WiFi N/A" (dimmed), each centered with the state-tinted WiFi icon.
//   - Rows drop from the bottom (details first, SSID/signal last) when the panel
//     is too short; nothing overflows the panel edges.
//
// Intentional (do NOT report as defects):
//   - On very short panels only the SSID and signal row survive; detail rows are
//     deliberately dropped rather than crammed.
//   - Detail rows render in a dimmed color; only SSID and the signal readout are
//     at full brightness/accent.
//   - Link speed and IP address are always shown when present (not policy-gated),
//     matching the bespoke panel layouts.
//
// Not verifiable from a single still frame:
//   - Live refresh cadence and the re-render skipping driven by RenderCacheKey.
//
// Rendering engines (both signal-first; the panel's refresh class picks one):
//   - Fast panels (color/grayscale/mono OLED, TFT) get buildSignalForwardWifi:
//     an anti-aliased 270° arc gauge hero with a plain-language quality word in
//     its center, SSID + icon in the arc's open-bottom wedge, IP footer.
//   - Slow / e-paper panels (MonoSlow and any zero-value-Params style, including
//     color-slow) get buildSlowWifi: the SAME signal-first design rendered with
//     e-paper-safe widgets — a hard-edged SEGMENTED signal bar instead of the
//     arc (crisp under 1-bit quantization, no shimmer), the quality word as the
//     readout, a quiet SSID + state-tinted icon, and the IP footer. Marked
//     Static for no-refresh displays. Both engines composite absolute-positioned
//     sprites only (no renderer text rows), so the two-plane alignment bug cannot
//     occur, and both drop bands from the bottom (detail, then IP) when a short
//     panel runs out of height, centering the surviving block.
//
// Per-resolution notes:
//   - Narrow portrait panels (e.g. 80x160, 128x296) have little horizontal room;
//     the SSID and radio line are the first candidates for size-fitting there
//     (the SSID auto-fits to the largest font that fits rather than truncating).
//   - Narrow-TALL panels (width-starved but tall, e.g. 32x128, aspect height ≥
//     2× width) use a VERTICAL signal bar climbing the long axis as the hero,
//     instead of a width-starved horizontal bar that would waste the height.
//     Both engines do this: the fast path via signalForwardCompactVertical
//     (Linear vertical), the slow/e-paper path via slowConnectedVertical
//     (Segmented vertical, crisp under 1-bit).
//   - 128x128, 128x64 (mono), 240x240, and 320x240 use bespoke build functions;
//     all other fast resolutions share buildSignalForwardWifi in
//     styles/layout_signal.go; all slow resolutions share buildSlowWifi in
//     styles/layouts.go.
package wifi
