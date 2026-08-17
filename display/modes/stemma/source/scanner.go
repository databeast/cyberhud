// Package stemma scans the I2C bus(es) exposed by the Adafruit Cyberdeck HAT
// (STEMMA QT / QWIIC connectors) and maintains a live inventory of attached
// devices.
//
// The Adafruit Cyberdeck HAT routes the Raspberry Pi's primary I2C bus (bus 1,
// GPIO2/SDA and GPIO3/SCL) through two JST SH 4-pin STEMMA QT / QWIIC
// connectors.  A second I2C bus (bus 3, routed through GPIO header pins) may
// also be available depending on the operating system configuration.
package source

import (
	"fmt"
	"sync"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

// Device describes a detected I2C device on a STEMMA QT / QWIIC connector.
type Device struct {
	Bus     string // bus identifier, e.g. "/dev/i2c-1"
	Addr    uint16 // 7-bit I2C address
	Name    string // human-readable name (from knownDevices), or "unknown"
	Present bool   // true if device was detected on last scan
}

// Scanner polls one or more I2C buses for STEMMA QT / QWIIC devices and
// maintains an up-to-date device list.
type Scanner struct {
	busNames []string
	interval time.Duration

	mu      sync.RWMutex
	devices map[string]*Device // key = "bus:addr"

	stopCh   chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
}

// New creates a Scanner that will monitor the given I2C bus names.
// If busNames is empty the default bus ("/dev/i2c-1") is used.
// interval controls how often the buses are re-scanned.
func New(busNames []string, interval time.Duration) *Scanner {
	if len(busNames) == 0 {
		busNames = []string{"/dev/i2c-1"}
	}
	if interval <= 0 {
		interval = 2 * time.Second
	}
	return &Scanner{
		busNames: busNames,
		interval: interval,
		devices:  make(map[string]*Device),
		stopCh:   make(chan struct{}),
	}
}

// Start begins background scanning.
func (s *Scanner) Start() {
	s.wg.Add(1)
	go s.run()
}

// Stop stops background scanning.
func (s *Scanner) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopCh)
	})
	s.wg.Wait()
}

// Devices returns a snapshot of all currently known devices.
func (s *Scanner) Devices() []*Device {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Device, 0, len(s.devices))
	for _, d := range s.devices {
		cp := *d
		out = append(out, &cp)
	}
	return out
}

// PresentDevices returns only the devices that were present on the last scan.
func (s *Scanner) PresentDevices() []*Device {
	all := s.Devices()
	present := all[:0]
	for _, d := range all {
		if d.Present {
			present = append(present, d)
		}
	}
	return present
}

func (s *Scanner) run() {
	defer s.wg.Done()

	// Run an immediate scan before the first tick.
	s.scanAll()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.scanAll()
		}
	}
}

func (s *Scanner) scanAll() {
	for _, busName := range s.busNames {
		s.scanBus(busName)
	}
}

func (s *Scanner) scanBus(busName string) {
	bus, err := i2creg.Open(busName)
	if err != nil {
		// Bus may not be present on all systems; skip silently.
		return
	}
	defer bus.Close()

	// Mark all existing entries for this bus as absent; we'll update below.
	s.mu.Lock()
	for _, d := range s.devices {
		if d.Bus == busName {
			d.Present = false
		}
	}
	s.mu.Unlock()

	for addr := uint16(0x03); addr <= 0x77; addr++ {
		if s.probeAddr(bus, busName, addr) {
			s.markPresent(busName, addr)
		}
	}
}

// probeAddr attempts a zero-length read to addr; returns true if ACKed.
func (s *Scanner) probeAddr(bus i2c.Bus, busName string, addr uint16) bool {
	dev := &i2c.Dev{Bus: bus, Addr: addr}
	err := dev.Tx(nil, make([]byte, 1))
	return err == nil
}

func (s *Scanner) markPresent(busName string, addr uint16) {
	key := fmt.Sprintf("%s:0x%02X", busName, addr)
	s.mu.Lock()
	defer s.mu.Unlock()
	if d, ok := s.devices[key]; ok {
		d.Present = true
	} else {
		s.devices[key] = &Device{
			Bus:     busName,
			Addr:    addr,
			Name:    knownDeviceName(addr),
			Present: true,
		}
	}
}
