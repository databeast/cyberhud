// Package schedule implements the "schedule" display mode: a status view that
// shows the calendar event active right now, read from an iCalendar (.ics)
// source that may be a local file path or an http(s) URL.
//
// # Shape
//
// The mode follows the single-intent pattern of clock/serial/wifi:
//
//   - source/ owns all live data: the Policy and Snapshot types, the background
//     fetch goroutine (Activate/Deactivate), and the iCalendar fetch+parse step.
//     The parsed, start-sorted event set lives on the Snapshot; the "what is on
//     now" decision is made at render time against the current instant so it
//     tracks the clock between fetches.
//   - styles/ owns placement: one adaptive layout (styles.adaptiveBuild) that
//     sizes the active event's title to the panel, shows its time window beneath,
//     and optionally previews the next event. Per-resolution declarations exist
//     only so the registry can resolve a capability/size cell for a panel.
//   - the mode package wires the two together: BuildView dispatches through the
//     style registry; policy.go persists config via a PolicySnapshotter; cmd.go
//     registers the "schedule" catalog definition and console verb; instance.go
//     registers the factory and drives the fetcher lifecycle.
//
// # Views
//
// When an event is active its title shows in the panel's largest fitting tier,
// accented (amber on color panels) to read as "live now", with its time window
// below. With nothing active the mode previews the next upcoming event and how
// soon it starts. Unconfigured, loading, and error states each state themselves
// plainly in a restrained, unaccented treatment.
//
// # Limitation
//
// Recurring events (RRULE) are not yet expanded: only the single DTSTART/DTEND
// an event literally declares is considered. A recurring entry therefore appears
// only for its first occurrence. Expanding recurrences is deferred to a later
// iteration.
package schedule
