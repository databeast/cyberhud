package styles

import (
	"fmt"
	"strings"
	"time"

	"github.com/databeast/cyberhud/display/modes/schedule/source"
)

// resolveView turns the raw snapshot into the three display rows for the instant
// snap.Now, applying the policy's time format and upcoming-preview toggle.
//
// This is pure formatting: event selection lives on the Snapshot
// (ActiveEvent/NextEvent), and the instant arrives on the snapshot. The function
// holds no clock and no I/O, so it is trivially testable and identical across
// every panel size.
func resolveView(snap source.Snapshot, pol source.Policy) View {
	now := snap.Now
	if now.IsZero() {
		now = time.Now()
	}

	// Not configured / not yet loaded / errored: say so plainly, no accent.
	if snap.Source == "" {
		return View{Title: "No calendar set", Subtitle: "set schedule source"}
	}
	if snap.LastError != "" && !snap.Loaded {
		return View{Title: "Calendar error", Subtitle: shortError(snap.LastError)}
	}
	if !snap.Loaded {
		return View{Title: "Loading..."}
	}

	if active, ok := snap.ActiveEvent(now); ok {
		v := View{
			Title:  active.Summary,
			Accent: true,
		}
		v.Subtitle = eventWindow(active, pol)
		if active.Location != "" {
			// Prefer the time window; fall back to location when there is no
			// window worth showing (should be rare for an active event).
			if v.Subtitle == "" {
				v.Subtitle = active.Location
			}
		}
		if pol.ShowUpcoming {
			if next, ok := snap.NextEvent(now); ok {
				v.Upcoming = "next: " + next.Summary + " " + relativeStart(next, now)
			}
		}
		return v
	}

	// Nothing active: preview the next event if there is one.
	if next, ok := snap.NextEvent(now); ok {
		return View{
			Title:    next.Summary,
			Subtitle: "in " + humanizeUntil(next.Start.Sub(now)) + " - " + eventWindow(next, pol),
		}
	}

	return View{Title: "Nothing scheduled"}
}

// eventWindow formats an event's start–end time window per the policy's time
// format. All-day events read as a date. A zero end collapses to just the start.
func eventWindow(e source.Event, pol source.Policy) string {
	if e.Start.IsZero() {
		return ""
	}
	if e.AllDay {
		return "all day - " + e.Start.Format("Mon Jan 2")
	}
	start := formatClock(e.Start, pol.TimeFormat)
	if e.End.IsZero() {
		return start
	}
	// Same-day events show only the end time; multi-day include the end date.
	// Separators are ASCII: the bitmap fonts carry no en-dash/arrow glyph, so a
	// Unicode "–"/"→" renders as tofu ("11:29?12:29").
	if sameDay(e.Start, e.End) {
		return start + "-" + formatClock(e.End, pol.TimeFormat)
	}
	return start + " -> " + e.End.Format("Jan 2 ") + formatClock(e.End, pol.TimeFormat)
}

// relativeStart renders a short parenthetical for how soon an upcoming event
// begins, e.g. "(in 20m)".
func relativeStart(e source.Event, now time.Time) string {
	d := e.Start.Sub(now)
	if d < 0 {
		return ""
	}
	return "(in " + humanizeUntil(d) + ")"
}

// humanizeUntil renders a duration as a compact, human "2h 05m" / "45m" / "30s".
func humanizeUntil(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		h := int(d / time.Hour)
		m := int((d % time.Hour) / time.Minute)
		return fmt.Sprintf("%dh %02dm", h, m)
	}
	days := int(d / (24 * time.Hour))
	h := int((d % (24 * time.Hour)) / time.Hour)
	return fmt.Sprintf("%dd %dh", days, h)
}

// formatClock renders a time per the "24h"/"12h" policy format.
func formatClock(t time.Time, format string) string {
	if format == "12h" {
		return strings.TrimSpace(t.Format("3:04PM"))
	}
	return t.Format("15:04")
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// shortError trims a verbose error chain to its leading clause for a single row.
func shortError(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, ":"); i > 0 && i < 24 {
		// Keep the operation prefix ("fetch", "open", "parse") plus a little.
		return s
	}
	const max = 40
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}
