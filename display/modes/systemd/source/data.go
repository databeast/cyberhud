package source

import "fmt"

// Snapshot is one boot-progress sample used for rendering and signatures.
type Snapshot struct {
	Desired      string
	Loading      string
	Loaded       string
	BootComplete bool
	// Progress fields for bar/pie styles.
	ActiveTargets int // Number of currently active targets
	TotalTargets  int // Total known targets
}

// BuildItems returns formatted boot-progress lines for panel rendering.
func BuildItems() []string {
	return BuildDefaultItems(PollSnapshot())
}

// Signature returns a stable key for boot-progress change detection.
func Signature(p Policy) string {
	snap := PollSnapshot()
	return fmt.Sprintf("%t|%s|%s|%s|%d/%d|%s|%s",
		snap.BootComplete, snap.Desired, snap.Loading, snap.Loaded,
		snap.ActiveTargets, snap.TotalTargets,
		p.Style, p.ColorAccent)
}
