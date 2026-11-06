// Package storage implements the "storage" display mode: a glanceable readout of
// filesystem usage across the host's mounted volumes.
//
// # User Experience
//
// Storage answers one question at a glance — "how full is this machine?" Each
// monitored filesystem is a row: its mountpoint, a compact used/total summary
// with a usage percentage, and a progress bar whose fill and color encode disk
// pressure. The color follows a green→amber→red heat ramp (on color-capable
// panels), so a nearly-full volume reads as alarming red while a roomy one stays
// calm green. Mounts are ordered most-used-first, with the root filesystem
// always promoted, so the most urgent volume is on top and survives truncation
// on short panels.
//
// On space-constrained panels the mode collapses to a single-filesystem focus
// readout of the most-used volume — percentage, free space, and a full-width
// bar — rather than cramming illegible rows.
//
// # Class
//
// Informational / periodic-poll. A background goroutine samples filesystem usage
// on an interval (default 5s, configurable via refresh_ms). Storage changes
// slowly, so the default cadence is relaxed relative to CPU or thermal modes.
//
// # Data Source
//
// On Linux the sampler enumerates /proc/mounts and queries each filesystem with
// statfs(2), reporting used/total/free to match df's accounting (reserved blocks
// count as used). Pseudo and virtual filesystems (tmpfs, overlay, cgroup, ...)
// are hidden by default and shown with show_all=true. On non-Linux platforms the
// sampler reports no filesystems and the mode renders a placeholder.
//
// # Styles
//
// Styles follow the per-resolution + capability-class registry pattern shared by
// thermal and system. Each resolution has a skeleton (adaptive) entry plus two
// explicit variants: "bars" (multi-mount rows) and "focus" (single filesystem).
// A bare resolution requirement means larger panels in the same capability class
// still resolve via the fitness scorer.
package storage
