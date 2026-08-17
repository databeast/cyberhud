package catalog

// BootPolicy determines the startup mode override for a panel.
type BootPolicy interface {
	// BootMode returns the mode to inject at boot, or "" for no override.
	BootMode(panelIndex int) string
}

// SystemdBootPolicy injects "systemd" on panel 0 during boot.
type SystemdBootPolicy struct{}

func (SystemdBootPolicy) BootMode(panelIndex int) string {
	if panelIndex == 0 {
		return "systemd"
	}
	return ""
}
