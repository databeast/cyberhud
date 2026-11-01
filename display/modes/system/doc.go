// Package system implements the "system" display mode: an at-a-glance monitor
// of the host the cyberHUD is attached to. It surfaces identity and health
// (hostname, OS/arch, uptime, network addresses) and live load (per-core CPU
// utilization and the top CPU-consuming processes) across four sub-styles —
// default, compact, cores, and top — selected via the "system" command verb.
//
// Data comes from a swappable sampler.Sampler (the /proc-backed ProcSampler on
// Linux); identity fields come from the source package. Styles are registered
// in systemRegistry and dispatched by fitness through the shared style
// framework; each style's Build method owns only layout, reading a
// source.SystemSnapshot assembled per frame.
//
// # User Experience
//
// Promise (verbatim from the catalog registration in instance.go):
//   Title:   "System"
//   Summary: "Host, uptime, and IP address information."
//
// PROMISE GAP (flagged for the developer, not invented intent): the registered
// summary advertises only host/uptime/IP — it does NOT mention per-core CPU
// utilization or top processes, yet the cores and top sub-styles deliver
// exactly those. The menu under-sells the mode: a user would not know the
// htop-style views exist. The summary should likely be widened to name CPU and
// process monitoring so the pitch matches what the mode actually does. Until it
// is, the code-blind reviewer is handed the narrow summary above, and the
// cores/top views legitimately exceed what was promised rather than fail it.
//
// Class: informational. Function first — a user picks this to KNOW the state of
// their machine at a glance. Appeal is a bonus, but this mode has a specific
// aesthetic birthright it should claim: it is a little glass cockpit for a
// computer, and the good version feels like a server-room status panel, not a
// text dump. Appeal here is earned by making data legible AT A GLANCE (shape,
// fill, and color carrying meaning before a single word is read), not by
// decoration.
//
// The four sub-styles are not four copies of the same screen; each answers a
// different question, and each has a distinct "at its best":
//
//   - default — "Who and where is this box, and is it OK?" The identity card.
//     At its best: a clean stack — a prominent hostname, then OS/arch, uptime,
//     and each network address — anchored by a genuine HEALTH indicator the
//     eye lands on first. The status glyph should read unmistakably as a verdict
//     (a calm check when all is well, an alarmed mark when something is missing
//     or wrong), not as a stray character. Uptime should read as a badge of
//     reliability ("up 14d 3h"), the thing a sysadmin is quietly proud of.
//
//   - compact — "The one line I'd put on a tiny always-on strip." The name tag.
//     At its best: the single most identifying fact, dense and dignified, sized
//     to own a short status-bar panel (think a 1-2 row strip, not a full
//     portrait screen). Hostname plus primary address is the floor; if room
//     allows, a health dot earns its place. Compact should feel deliberately
//     minimal — a label you'd stick on the front of a rack unit — NOT a full
//     screen left 90% empty. Emptiness is the failure mode to design against.
//
//   - cores — "How hard is every core working, right now?" The load wall.
//     At its best: one bar PER CORE, filling the panel as a tidy stack or grid,
//     every core always shown (idle cores draw an empty track, never vanish),
//     each bar's FILL and COLOR carrying the load — cool green at idle climbing
//     through amber to an urgent red near saturation — so a glance across the
//     wall reads the machine's effort instantly, like a mixing desk's meters.
//     The percentage is the caption; the color is the headline. A flat
//     monochrome bar that ignores load, or a panel showing only the first few
//     cores, squanders the whole point of the view.
//
//   - top — "What's eating my CPU?" The leaderboard. At its best: a ranked
//     descending list of the heaviest processes, name left and percentage right,
//     the worst offender unmistakably on top, with the heaviest rows accented
//     (the same cool-to-hot logic as cores) so the troublemaker draws the eye.
//     Type rhythm should stay even — one row shrinking to fit breaks the
//     leaderboard feel. Names truncate gracefully with the percentage always
//     preserved; the number is the point.
//
// Invariants (must hold — a usability defect if violated):
//   - Text is legible with correct glyphs; no tofu, overlap, or clipping of the
//     values that carry meaning (a temperature, a percentage, an address).
//   - cores shows ALL cores reported by the sampler, each with a visible bar or
//     empty track; it never silently drops cores.
//   - top lists processes in descending CPU order with the percentage always
//     visible, even when the name is truncated.
//   - The default status indicator reflects real health and is recognizable as a
//     status verdict, not ambiguous clutter.
//   - A style fills its panel proportionate to the information it carries; large
//     dead space on an informational screen reads as broken or half-loaded.
//
// Intentional (do NOT report as defects):
//   - compact deliberately shows very little; judged against a short strip, its
//     sparseness is the design (but see the emptiness caveat above when it is
//     rendered on a tall panel it was not meant to fill).
//   - A truncated process or address name is acceptable when the panel is
//     genuinely too narrow, provided the numeric value survives.
//
// Not verifiable from a single still frame:
//   - Whether per-core bars and top-process rows UPDATE live and smoothly.
//   - Whether uptime, CPU, and process values are current vs stale.
//   - Whether cores shown as idle are truly idle at that instant or simply not
//     yet sampled.
//
// Per-resolution notes:
//   - These styles are content layouts, not per-resolution variants; the same
//     style must behave on a short status strip and a tall portrait panel.
//   - Tall panels expose the sparse-readout failure mode (compact, partially
//     drawn cores); wide/short panels expose row-count and truncation limits in
//     top and cores. Both extremes are worth a look.
//   - Where the panel supports color, cores and top should use it to encode load
//     (green→amber→red); on mono/grayscale panels the same meaning falls back to
//     fill level and ordering alone.
package system
