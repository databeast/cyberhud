package source

import "time"

// ClockData is the mode-specific data consumed by style Build methods.
// Contains pre-formatted time, date, and weekday strings plus the active
// Policy fields needed for rendering decisions (widget flags, accent color).
// Now carries the raw time.Time so styles can pass it to appendClockWidgets
// for second-precision widget updates (LED blink, progress bar, sparkline).
type ClockData struct {
	Time    string    // Formatted time string (e.g., "14:30:05" or "2:30 PM")
	Date    string    // Formatted date string (e.g., "2024-01-15") or "" when DateFormat="none"
	Weekday string    // Day name (e.g., "Monday") or "" when ShowWeekday=false
	Now     time.Time // Raw time for widget helpers (LED, progress bar, sparkline)
}
