package source

import "strings"

// Version is set at build time via:
//
//	go build -ldflags="-X github.com/databeast/cyberhud/display/modes/dashboard.Version=v1.2.3"
var Version string

// getVersion returns the build version string, falling back to "dev" when
// the variable is empty or whitespace-only (unset ldflags).
func GetVersion() string {
	v := strings.TrimSpace(Version)
	if v == "" {
		return "dev"
	}
	if len(v) > 64 {
		return v[:64]
	}
	return v
}
