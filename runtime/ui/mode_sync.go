package ui

import (
	"log"
	"strings"

	displaymodes "github.com/databeast/cyberhud/display/modes"
	systemdmode "github.com/databeast/cyberhud/display/modes/systemd/source"
	"github.com/databeast/cyberhud/display/region"
)

// syncMode synchronizes the Region's active mode with coordinator.State.
// On each render tick, it compares the coordinator's desired mode with the Region's
// current mode. If they differ, it calls Region.SetMode to construct a fresh
// instance (cursor=0, topRow=0 inherently from fresh construction).
//
// It also handles the systemd boot-mode transition: when the current mode is
// "systemd" and multi-user.target has been reached, it transitions to "dashboard".
//
// The passive normalization (menu → dashboard for panels without input) is applied
// before comparison.
func (rr *RegionRenderer) syncMode(r *region.Region) {
	currentMode := r.CurrentMode()

	// Boot-mode transition: systemd → dashboard once multi-user is reached.
	if currentMode == "systemd" && systemdmode.ReachedMultiUser() {
		if err := r.SetMode("dashboard"); err != nil {
			log.Printf("syncMode: boot transition failed for region %q: %v", r.Name(), err)
		}
		// Also update coordinator.State to reflect the transition.
		if rr.modeState != nil {
			_, _ = rr.modeState.Set(0, "dashboard")
		}
		return
	}

	// If no coordinator.State is configured, nothing to sync.
	if rr.modeState == nil {
		return
	}

	// Read the desired mode from coordinator.State by region name.
	// Each region's name matches its corresponding panel name in coordinator.State
	// (e.g., "main", "left-aux", "right-aux" for multi-screen panels).
	// For single-screen panels where the region is named "default", fall back
	// to panel index 0.
	desiredMode := rr.modeState.CurrentModeByName(r.Name())
	if desiredMode == "" {
		// Fallback for single-screen panels: region name is "default" but
		// coordinator.State uses the panel product name. Use index 0.
		desiredMode = rr.modeState.CurrentMode(0)
	}
	desiredMode = rr.normalizeMode(desiredMode)

	if desiredMode == "" || desiredMode == currentMode {
		return
	}

	// Mode differs — call Region.SetMode to construct a fresh instance.
	if err := r.SetMode(desiredMode); err != nil {
		log.Printf("syncMode: mode switch to %q failed for region %q: %v", desiredMode, r.Name(), err)
	}
}

// normalizeMode applies mode normalization rules:
// - Lowercases and trims whitespace
// - Converts "menu" to "dashboard" for passive (no-input) panels
// - Returns "" for unknown modes
func (rr *RegionRenderer) normalizeMode(mode string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "menu" && rr.inputMapper == nil {
		return "dashboard"
	}
	if !displaymodes.IsKnownInstance(mode) {
		return ""
	}
	return mode
}
