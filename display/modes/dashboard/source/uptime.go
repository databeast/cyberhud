package source

import (
	"image/color"

	sharedcolor "github.com/databeast/cyberhud/display/style/color"
)

// FormatUptime renders a duration as a human-readable string, max 12 characters.
// Negative durations are treated as 0. Uses day/hour/minute granularity (never
// seconds) to match the minute-granularity RenderCacheKey.
//
// Format rules:
//   - < 1 hour:  "{m}m"       (e.g., "0m", "45m")
//   - < 1 day:   "{h}h {m}m"  (e.g., "2h 30m")
//   - ≥ 1 day:   "{d}d {h}h"  (e.g., "3d 5h", "365d 0h")

// resolveUptimeAccent returns the uptime indicator color based on the accent name.
// When accent is "none", returns opaque white. Otherwise delegates to sharedcolor.ResolveAccent.
func ResolveUptimeAccent(accent string) color.RGBA {
	if accent == "none" {
		return color.RGBA{255, 255, 255, 255}
	}
	return sharedcolor.ResolveAccent(accent)
}
