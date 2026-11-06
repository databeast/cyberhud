// Package dashboard implements the "dashboard" display mode: a glanceable
// system-liveness panel. It rolls several host facts (hostname, uptime, IP,
// WiFi, build version) together with an aggregate health verdict into one
// at-a-glance "is this box okay?" view, distinct from the dedicated
// per-subsystem modes (thermal, system, storage).
//
// The mode owns no background work: content is read on demand in
// source.BuildDashboardContent and rendered statically. Per-panel styles live
// in styles/ and are selected by the shared StyleRegistry on capability +
// surface fit.
//
// # User Experience
//
// Promise (verbatim from catalog registration in mode.go init):
//
//	Name:        "Dashboard"
//	Description: "Glanceable liveness panel showing hostname, uptime heartbeat,
//	             IP, WiFi, version, and panel name."
//
// Class: informational
//
// Invariants (must hold — a usability defect if violated):
//   - The panel's identity is legible first: on any panel large enough for one
//     row, the hostname is shown and not truncated into meaninglessness.
//   - Whatever rows are shown fit the panel: no text clipped at the left or
//     right edge, no row drawn past the bottom.
//   - The content block is deliberately placed (centered or top-anchored by
//     design), not floating at an arbitrary offset.
//   - On panels too small for all six facts, the fields shown are the
//     highest-priority ones (identity + health), not simply the first five in
//     struct order with the rest dropped off the bottom.
//   - The health verdict is represented when space permits — it is the fact
//     that distinguishes this mode from a plain host-info readout.
//
// Intentional (do NOT report as defects):
//   - Lower fields (version, WiFi) are dropped first on space-constrained
//     panels; losing them is by design, not a bug.
//   - Mono and grayscale panels render without color; the absence of color is
//     a capability of the panel, not a defect of the style.
//   - Abbreviated field labels / shortened values on narrow panels are
//     deliberate, not truncation bugs.
//
// Not verifiable from a single still frame:
//   - The "uptime heartbeat" cadence — a still frame shows the uptime value,
//     not its liveness pulse.
//
// Per-resolution notes (especially extreme aspect ratios / narrow portrait):
//   - Very small panels (a dimension < 64px, e.g. 128x32, 32x128): only one or
//     two rows fit. Show identity + health, not a clipped five-row stack.
//   - Narrow portrait (e.g. 80x160, 135x240): width is the binding constraint;
//     values must be abbreviated or wrapped, not centered-then-clipped.
//   - Landscape short (e.g. 160x80, 212x104): height is the binding constraint;
//     prioritize rows by importance.
//
// Known promise drift (finding, not yet resolved):
//   - The catalog Description promises "panel name", but DashboardContent
//     carries no panel-name field; it carries Health instead. Either the
//     Description or the content model is stale. Treat a missing panel name as
//     expected until this is reconciled.
package dashboard
