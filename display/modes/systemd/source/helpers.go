package source

import "fmt"

// BootStatusLabel returns a single-line summary of the boot state.
// Returns "Boot Complete" when finished, otherwise the currently-loading
// target name (falling back to the desired target if loading is empty).
func BootStatusLabel(snap Snapshot) string {
	if snap.BootComplete {
		return "Boot Complete"
	}
	if snap.Loading != "" {
		return snap.Loading
	}
	return snap.Desired
}

// BuildDefaultItems returns the multi-row boot status text list for reuse
// by styles that display line-by-line boot progress (e-ink, mono, default).
// Format: "Booting...." header, loaded/loading detail, optional desired line.
func BuildDefaultItems(snap Snapshot) []string {
	if snap.BootComplete {
		return []string{"Boot Complete"}
	}
	items := []string{"Booting...."}
	if snap.Loaded != "" && snap.Loaded != snap.Loading {
		items = append(items, fmt.Sprintf("Loaded: %s : Loading: %s", snap.Loaded, snap.Loading))
	} else {
		items = append(items, fmt.Sprintf("Loading: %s", snap.Loading))
	}
	if snap.Desired != "" && snap.Desired != snap.Loading {
		items = append(items, "Desired: "+snap.Desired)
	}
	return items
}
