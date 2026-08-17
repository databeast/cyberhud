package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/databeast/cyberhud/display/catalog"
	"github.com/databeast/cyberhud/display/coordinator"
	displaymodes "github.com/databeast/cyberhud/display/modes"
	// Mode self-registration via init().
	_ "github.com/databeast/cyberhud/display/modes/attract_bokeh"
	_ "github.com/databeast/cyberhud/display/modes/attract_geometric"
	_ "github.com/databeast/cyberhud/display/modes/attract_hacking"
	_ "github.com/databeast/cyberhud/display/modes/attract_matrix"
	_ "github.com/databeast/cyberhud/display/modes/attract_particles"
	_ "github.com/databeast/cyberhud/display/modes/attract_plasma"
	_ "github.com/databeast/cyberhud/display/modes/attract_shapes"
	_ "github.com/databeast/cyberhud/display/modes/attract_starfield"
	_ "github.com/databeast/cyberhud/display/modes/attract_waveform"
	_ "github.com/databeast/cyberhud/display/modes/clock"
	_ "github.com/databeast/cyberhud/display/modes/cycle"
	_ "github.com/databeast/cyberhud/display/modes/dashboard"
	_ "github.com/databeast/cyberhud/display/modes/gpio"
	_ "github.com/databeast/cyberhud/display/modes/gpio_control"
	_ "github.com/databeast/cyberhud/display/modes/image"
	_ "github.com/databeast/cyberhud/display/modes/menu"
	_ "github.com/databeast/cyberhud/display/modes/pager"
	_ "github.com/databeast/cyberhud/display/modes/serial"
	_ "github.com/databeast/cyberhud/display/modes/stemma"
	_ "github.com/databeast/cyberhud/display/modes/system"
	_ "github.com/databeast/cyberhud/display/modes/systemd"
	_ "github.com/databeast/cyberhud/display/modes/testfonts"
	_ "github.com/databeast/cyberhud/display/modes/testicons"
	_ "github.com/databeast/cyberhud/display/modes/testpattern"
	_ "github.com/databeast/cyberhud/display/modes/testwidgets"
	_ "github.com/databeast/cyberhud/display/modes/thermal"
	_ "github.com/databeast/cyberhud/display/modes/ticker"
	_ "github.com/databeast/cyberhud/display/modes/usb"
	_ "github.com/databeast/cyberhud/display/modes/wifi"
	_ "github.com/databeast/cyberhud/display/modes/zmq"
	"github.com/databeast/cyberhud/hardware/panels"
)

// expectedModeIDs lists all 28 mode IDs that must be registered after init().
var expectedModeIDs = []string{
	"menu",
	"dashboard",
	"stemma",
	"gpio",
	"system",
	"systemd",
	"clock",
	"ticker",
	"image",
	"usb",
	"serial",
	"thermal",
	"gpio-control",
	"testpattern",
	"testfonts",
	"testicons",
	"testwidgets",
	"cycle",
	"attract_matrix",
	"zmq",
	"wifi",
	"attract_waveform",
	"attract_particles",
	"attract_plasma",
	"attract_starfield",
	"attract_bokeh",
	"attract_shapes",
	"attract_geometric",
	"attract_hacking",
}

// validateModeRegistrations asserts that all 27 expected mode IDs are registered
// in the mode registry. It must be called after all init() functions have fired
// (i.e., from main or any function called from main). If any mode IDs are missing,
// it terminates the process with a descriptive error.
func validateModeRegistrations() {
	var missing []string
	for _, id := range expectedModeIDs {
		if !displaymodes.IsKnownInstance(id) {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		log.Fatalf("startup validation failed: %d mode(s) not registered: %s",
			len(missing), fmt.Sprintf("%v", missing))
	}
}

func configureDisplayModes(state *coordinator.State, profile panels.Definition, inputEnabled bool) {
	if state == nil {
		return
	}
	boot := catalog.SystemdBootPolicy{}
	policy := catalog.StandardPolicy{}
	profile.InputEnabled = inputEnabled
	for i := range profile.Virtual {
		profile.Virtual[i].InputEnabled = inputEnabled
	}
	displays := panels.Displays(profile, policy)
	for i := range displays {
		if mode := boot.BootMode(displays[i].Index); mode != "" {
			displays[i].Modes = prependMode(mode, displays[i].Modes)
			displays[i].Default = mode
		}
	}
	state.ResetWithFallback(displays, "dashboard")
}

func prependMode(mode string, modes []string) []string {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		out := make([]string, len(modes))
		copy(out, modes)
		return out
	}
	out := make([]string, 0, len(modes)+1)
	out = append(out, mode)
	for _, m := range modes {
		norm := strings.ToLower(strings.TrimSpace(m))
		if norm == "" || norm == mode {
			continue
		}
		out = append(out, norm)
	}
	return out
}
