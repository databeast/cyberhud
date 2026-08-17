package source

// DefaultPolicy returns the default WiFi policy with all fields initialized.
func DefaultPolicy() Policy {
	return Policy{
		Style:         "",
		FGColor:       "green",
		SignalDisplay: "bars",
		ShowFrequency: true,
		ShowInterface: true,
		ShowChannel:   true,
	}
}
