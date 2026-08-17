package source

import "fmt"

// NewTestScanner constructs an in-memory scanner populated with the provided
// devices. It is intended for snapshot tests that need realistic STEMMA data
// without opening a real I2C bus.
func NewTestScanner(devs ...*Device) *Scanner {
	s := &Scanner{
		devices: make(map[string]*Device, len(devs)),
		stopCh:  make(chan struct{}),
	}
	for _, d := range devs {
		if d == nil {
			continue
		}
		cp := *d
		if cp.Bus == "" {
			cp.Bus = "(no bus)"
		}
		s.devices[fmt.Sprintf("%s:0x%02X", cp.Bus, cp.Addr)] = &cp
	}
	return s
}
