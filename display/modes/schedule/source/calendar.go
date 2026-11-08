package source

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

// fetchTimeout bounds a single HTTP calendar fetch.
const fetchTimeout = 15 * time.Second

// maxCalendarBytes caps how much a remote or local calendar may be before the
// fetcher refuses it, guarding against a pathological source exhausting memory
// on a small device.
const maxCalendarBytes = 8 << 20 // 8 MiB

// httpClient is the shared client for remote calendar fetches.
var httpClient = &http.Client{Timeout: fetchTimeout}

// isURL reports whether source names a remote http(s) calendar rather than a
// local filesystem path.
func isURL(source string) bool {
	s := strings.ToLower(source)
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// loadCalendar reads and parses the calendar named by source, which is either an
// http(s) URL or a local filesystem path.
func loadCalendar(source string) (*ics.Calendar, error) {
	source = strings.TrimSpace(source)
	if source == "" {
		return nil, fmt.Errorf("no calendar source configured")
	}

	var r io.Reader
	if isURL(source) {
		resp, err := httpClient.Get(source)
		if err != nil {
			return nil, fmt.Errorf("fetch: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("fetch: unexpected status %s", resp.Status)
		}
		r = io.LimitReader(resp.Body, maxCalendarBytes)
	} else {
		f, err := os.Open(source)
		if err != nil {
			return nil, fmt.Errorf("open: %w", err)
		}
		defer f.Close()
		r = io.LimitReader(f, maxCalendarBytes)
	}

	cal, err := ics.ParseCalendar(r)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return cal, nil
}

// extractEvents reduces a parsed calendar to the display Event set, keeping only
// events that have not already ended relative to now and sorting them by start
// time ascending.
//
// Recurring events (RRULE) are not expanded in this iteration: only the single
// DTSTART/DTEND the event literally declares is used. This is a known limitation
// noted in the mode docs.
func extractEvents(cal *ics.Calendar, now time.Time) []Event {
	if cal == nil {
		return nil
	}
	raw := cal.Events()
	out := make([]Event, 0, len(raw))
	for _, ve := range raw {
		ev, ok := convertEvent(ve)
		if !ok {
			continue
		}
		// Drop events that have already ended. Point events (zero End) are kept
		// until their start passes by a day so a just-started all-day entry does
		// not vanish immediately.
		if !ev.End.IsZero() {
			if ev.End.Before(now) {
				continue
			}
		} else if ev.Start.Before(now.Add(-24 * time.Hour)) {
			continue
		}
		out = append(out, ev)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Start.Before(out[j].Start)
	})
	return out
}

// convertEvent maps a parsed VEVENT into a display Event, resolving timed and
// all-day start/end forms. Returns false when the event has no usable start.
func convertEvent(ve *ics.VEvent) (Event, bool) {
	var ev Event
	ev.Summary = propValue(ve, ics.ComponentPropertySummary)
	ev.Location = propValue(ve, ics.ComponentPropertyLocation)

	if start, err := ve.GetStartAt(); err == nil {
		ev.Start = start
	} else if start, err := ve.GetAllDayStartAt(); err == nil {
		ev.Start = start
		ev.AllDay = true
	} else {
		return Event{}, false
	}

	if ev.AllDay {
		if end, err := ve.GetAllDayEndAt(); err == nil {
			ev.End = end
		} else {
			// An all-day event with no explicit end spans its single day.
			ev.End = ev.Start.Add(24 * time.Hour)
		}
	} else if end, err := ve.GetEndAt(); err == nil {
		ev.End = end
	}

	if ev.Summary == "" {
		ev.Summary = "(untitled)"
	}
	return ev, true
}

// propValue returns the trimmed value of a component property, or "" when absent.
func propValue(ve *ics.VEvent, prop ics.ComponentProperty) string {
	p := ve.GetProperty(prop)
	if p == nil {
		return ""
	}
	// iCal escapes commas, semicolons and newlines in TEXT values; FromText
	// reverses that so summaries/locations read naturally.
	return strings.TrimSpace(ics.FromText(p.Value))
}

// calendarName returns a display name for the calendar, preferring X-WR-CALNAME
// then NAME, or "" when neither is set.
func calendarName(cal *ics.Calendar) string {
	if cal == nil {
		return ""
	}
	for _, prop := range cal.CalendarProperties {
		switch strings.ToUpper(prop.IANAToken) {
		case "X-WR-CALNAME", "NAME":
			if v := strings.TrimSpace(prop.Value); v != "" {
				return v
			}
		}
	}
	return ""
}
