// Package usb implements a display mode for the cyberHUD system that identifies
// the most recently connected USB device, for bench work where you plug a device
// in and want the panel to tell you immediately what it is. It reads device
// identity from the Linux sysfs tree (/sys/bus/usb/devices), prefers instant
// netlink-driven change notification with a polling fallback, and presents the
// device's name, VID:PID, bus/device numbers, serial, and connection status.
//
// Architecture:
//
//   - Catalog registration: mode and command definitions registered via init()
//     with typed option metadata (poll_ms, hold_unplugged_ms, hide_root_hubs,
//     style), summaries, defaults, and allowed-value lists. The style option's
//     allowed values are the live registry style-name list.
//   - Source layer (source/): sysfs scanning (detector.go) that ranks devices by
//     directory mtime so the newest connection sorts first, device-class
//     classification from bDeviceClass / bInterfaceClass into human labels
//     (Hub/Storage/HID/Audio/…), root-hub filtering, an instant netlink monitor
//     (watcher.go) with a polling fallback, and the Snapshot/Transition state
//     machine (data.go) that tracks the current bench device and its connected
//     flag.
//   - Monitor state: a process-global monitor (usb.go) that debounces scans to
//     the poll interval, drains instant-notification events to force a rescan,
//     and applies the hold-unplugged policy so a just-removed device can linger
//     on screen for a configured window.
//   - Policy definition and normalization: a Policy struct (poll_ms,
//     hold_unplugged_ms, hide_root_hubs, style) with normalization that floors
//     the poll interval, clamps negative holds, and resolves the style name
//     against the registry (empty means auto-detect via fitness).
//   - Data-signature change detection: Signature()/RenderCacheKey return a
//     fingerprint incorporating the snapshot sequence, connected state, device
//     key, and all policy fields so the display runtime skips re-renders on
//     battery-powered hardware.
//   - Interface-based style dispatch: a style.NewRegistry[Snapshot, Policy]
//     holding one skeleton style plus one "-minimal" variant per
//     resolution × capability (mono/grayscale/color × slow/fast). An init()
//     guard enforces that every registered style name matches the expected
//     <capability><WxH>[-variant] pattern.
//   - Adaptive layout (styles/core.go): skeleton styles share adaptiveBuild,
//     which routes by panel pixel size to buildFull (w≥240 && h≥240),
//     buildCompact (w≥128), or buildMinimalLayout (w<128). The "-minimal"
//     variants pin buildMinimal regardless of size.
//   - Panel hints and fitness evaluation: per-instance panel hints injected by
//     the hosting Region (instance.go) with a modehints fallback for tests and
//     tooling, plus a PostApply hook that generates fitness notes when the
//     active style changes.
//   - Action handler: a Handler implementing action.Handler (Primary navigates
//     to the dashboard).
//
// # User Experience
//
// Promise (verbatim from catalog registration in usb.go init):
//
//	Name:        "USB Bench"
//	Description: "Highlights the most recently connected USB device for bench
//	             identification."
//
// Class: informational
//
// The screen's job is INSTANT DEVICE-TYPE RECOGNITION. On a bench you plug
// something in and want the panel to announce what it is before you have read a
// word — so on panels with room the device CLASS is the hero: a large,
// class-colored glyph (disk / keyboard / headset / router …) is the focal point,
// with the device name as the title and the identity details as a quiet
// supporting cast. "Legible and tidy" is the floor, not the goal; the goal is
// that a glance tells you "storage" or "audio" from across the bench.
//
// Invariants (must hold — a usability defect if violated):
//   - On a large panel (w≥240 && h≥240) the class glyph is the visual hero: big,
//     centered, and obviously the dominant element — not a small icon tucked
//     beside text. A keyboard must look like a keyboard at a glance.
//   - On color panels the class carries a distinct color (storage green, HID
//     cyan, audio orange, …) on the glyph and the accent bar, so device types
//     are told apart by color as well as shape.
//   - The device's human name (e.g. "Mega 2560") is the title, prominent and
//     read without a tofu/ellipsis box — fitted to a smaller font rather than
//     cut when it does not fit.
//   - The hero group (glyph + name + accent + details) is vertically centered as
//     one unit; it does not clump at the top of a tall panel leaving dead space.
//   - The identity details remain available and legible: VID:PID, Bus/Dev (when
//     known), serial (when present), and connected/unplugged status. They read
//     as secondary to the hero, never overflow the panel edge.
//   - When no device has been seen yet, the panel says so in plain language
//     ("Waiting for USB...") rather than rendering blank or stale.
//   - Connection status is unambiguous: a connected device and an unplugged
//     (held) device are visibly distinguishable, not identical frames.
//   - Nothing load-bearing renders in a color that vanishes against the panel
//     background — on 1-bit/grayscale panels the glyph, name, accent rule and
//     details are all full white (class color and dimmed detail text are a
//     color-panel enhancement only).
//
// Intentional (do NOT report as defects):
//   - The hero glyph is an INTEGER-upscaled bitmap, so it reads as crisp-but-
//     chunky pixel art rather than a smooth icon. The blocky look is the
//     aesthetic, not a scaling artifact.
//   - On color panels the detail rows render dimmed (lower brightness than the
//     title) for hierarchy; on mono they stay full white by necessity.
//   - On mid-size panels (128–239px wide) the screen shows a compact badge: the
//     class glyph above the device name, with the VID:PID as a quiet sub-line —
//     no full detail stack. On short landscape strips without vertical room for
//     a glyph, it falls back to name + VID:PID only.
//   - On narrow panels (<128px wide, and the "-minimal" variants) the screen
//     shows the class glyph above the device name — these panels are narrow but
//     tall, so the glyph fills the vertical space and still conveys device type.
//     A panel too short for the glyph shows the name alone. This is the minimal
//     contract, not an overflow failure.
//
// Not verifiable from a single still frame:
//   - That the panel updates "immediately" on a plug/unplug event — the instant
//     netlink monitor, the polling fallback cadence, and the hold-unplugged
//     linger window are all time-behaviors a still frame cannot show.
//   - That the most RECENTLY connected device is the one shown (mtime ranking);
//     a still frame shows one device, not the selection among several.
//
// Per-resolution notes (especially extreme aspect ratios / narrow portrait):
//   - The hero (buildHero) scales the 24px class glyph by the largest integer
//     factor fitting ~45% of the panel's smaller dimension, clamped to [1,3]; on
//     small-but-full panels the factor drops so the glyph still fits.
//   - Very wide/short panels (e.g. 296x128, 320x240) have limited vertical room,
//     so the bottom detail rows are the first dropped when the hero stack does
//     not fully fit.
//   - Compact and minimal layouts (below the hero threshold) render a compact
//     badge — class glyph stacked above the name (plus VID:PID sub-line on
//     compact) — vertically centered via buildCompactBadge, which drops the
//     glyph when a panel is too short to seat it. The narrow-tall portrait
//     strips (32x128, 80x160, 104x212, 122x250) are the main beneficiaries:
//     the glyph fills what was otherwise a mostly-empty column.
package usb
