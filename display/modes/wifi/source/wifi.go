package source

// DefaultPolicy returns the default WiFi policy with all fields initialized.
//
// The detail toggles default to false: the WiFi mode is an ambient, glanceable
// display whose hero is the live signal strength, not a diagnostic readout. The
// channel/frequency/interface fields are tinkerer detail, off by default; the
// user turns them on when debugging. (IP address is not gated by a toggle — on a
// headless box it is the one operationally critical field and is always shown.)
func DefaultPolicy() Policy {
	return Policy{
		Style:         "",
		FGColor:       "green",
		SignalDisplay: "bars",
		ShowFrequency: false,
		ShowInterface: false,
		ShowChannel:   false,
	}
}
