// Package gpio manages the general-purpose I/O pins exposed on the Adafruit
// Cyberdeck HAT 40-pin GPIO breakout header.
//
// It tracks the current mode (input/output) and state (high/low) of each
// configurable pin and exposes a snapshot suitable for display on the
// Waveshare LCD.
package gpio

import (
	"fmt"
	"strings"
	"sync"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/pin"
)

// PinMode describes how a GPIO pin is currently configured.
type PinMode uint8

const (
	ModeUnknown PinMode = iota
	ModeInput
	ModeOutput
	ModeAlt // ALT function (SPI, I2C, UART, PWM, …)
)

func (m PinMode) String() string {
	switch m {
	case ModeInput:
		return "IN"
	case ModeOutput:
		return "OUT"
	case ModeAlt:
		return "ALT"
	default:
		return "???"
	}
}

// PinState holds the last-known state for a single GPIO pin.
type PinState struct {
	Number int
	Name   string
	Mode   PinMode
	Level  gpio.Level

	// Func is the pin's current hardware function as reported by the driver
	// (periph pin.Func), e.g. "SPI0_MOSI", "I2C1_SDA", "UART0_TX", "In", "Out",
	// or "" when unknown/unavailable (non-RPi host). Unlike Mode — which
	// collapses every alternate function to the single ModeAlt bucket — Func
	// preserves WHICH peripheral an ALT pin is serving, so consumers can group
	// pins belonging to the same bus (SPI0, I2C1, …) into one logical entity.
	// This field is additive: existing consumers that read only Mode/Level are
	// unaffected.
	Func string
}

// BusGroup returns the bus/peripheral family a pin belongs to, derived from
// Func, or "" when the pin is not serving a recognized multi-pin bus.
//
// The family is the Func prefix up to the first underscore for known bus
// peripherals (SPI, I2C, UART, PWM, PCM, GPCLK, SMI). For example "SPI0_MOSI"
// and "SPI0_SCLK" both return "SPI0", so a consumer can collapse every pin of
// one bus into a single row. Plain digital functions ("In", "Out") and
// unrecognized values return "" — those are standalone pins, displayed
// individually.
func (p PinState) BusGroup() string {
	return busFamily(p.Func)
}

// busFamily extracts the bus/peripheral group name from a periph function
// string. It returns the portion before the first underscore when that portion
// names a known multi-pin bus, else "".
func busFamily(fn string) string {
	if fn == "" {
		return ""
	}
	i := strings.IndexByte(fn, '_')
	if i <= 0 {
		return "" // no underscore → plain function like "In"/"Out"/"Low"/"High"
	}
	prefix := fn[:i]
	// Recognize known bus/peripheral families. The trailing digits (bus index)
	// are part of the group name so SPI0 and SPI1 stay distinct.
	for _, fam := range []string{"SPI", "I2C", "UART", "PWM", "PCM", "GPCLK", "SMI", "DPI"} {
		if strings.HasPrefix(prefix, fam) {
			return prefix
		}
	}
	return ""
}

func (p PinState) String() string {
	if p.Mode == ModeUnknown {
		return fmt.Sprintf("%-2d --", p.Number)
	}
	lvl := "LO"
	if p.Level {
		lvl = "HI"
	}
	return fmt.Sprintf("%-2d %-3s %s", p.Number, p.Mode.String(), lvl)
}

// Manager tracks the user-accessible GPIO pins on the Cyberdeck HAT breakout.
// Pins used by the display HAT itself are excluded.
type Manager struct {
	mu   sync.RWMutex
	pins map[int]*managedPin
}

type managedPin struct {
	state   PinState
	out     gpio.PinOut
	in      gpio.PinIn
	pinfunc pin.PinFunc // non-nil when the driver reports pin function
}

// readFunc returns the pin's current function string via the PinFunc interface,
// or "" when unavailable (nil pin, or a driver that does not implement it).
func (mp *managedPin) readFunc() string {
	if mp.pinfunc == nil {
		return ""
	}
	return string(mp.pinfunc.Func())
}

// userPins lists the GPIO numbers available on the Cyberdeck HAT breakout that
// are not reserved for the Waveshare LCD HAT, SPI, or I2C.
var userPins = []int{
	// I2C (available via STEMMA QT)
	2, 3,
	// UART
	14, 15,
	// PWM
	12, 13,
	// General purpose (not used by LCD HAT or SPI)
	4, 17, 22, 23, 25, 27,
}

// outputOf returns the PinOut interface of p, or nil if p does not implement PinOut.
func outputOf(p gpio.PinIO) gpio.PinOut {
	out, _ := p.(gpio.PinOut)
	return out
}

// inputOf returns the PinIn interface of p, or nil if p does not implement PinIn.
func inputOf(p gpio.PinIO) gpio.PinIn {
	in, _ := p.(gpio.PinIn)
	return in
}

// New creates a Manager and registers all user-accessible pins.
func New() *Manager {
	m := &Manager{pins: make(map[int]*managedPin)}
	for _, num := range userPins {
		name := fmt.Sprintf("GPIO%d", num)
		p := gpioreg.ByName(name)
		entry := &managedPin{
			state: PinState{Number: num, Name: name, Mode: ModeUnknown},
		}
		if p != nil {
			entry.out = outputOf(p)
			entry.in = inputOf(p)
			if pf, ok := p.(pin.PinFunc); ok {
				entry.pinfunc = pf
			}
		}
		// Capture the pin's function as the system has it configured RIGHT NOW.
		// A pin already serving a bus (SPI/I2C/UART/…) at startup is live work
		// the display must surface, even though nothing in this process set it;
		// reflect that by marking it ModeAlt so it is not treated as idle.
		entry.state.Func = entry.readFunc()
		if busFamily(entry.state.Func) != "" {
			entry.state.Mode = ModeAlt
		}
		m.pins[num] = entry
	}
	return m
}

// Snapshot returns the current state of all managed pins.
func (m *Manager) Snapshot() []PinState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]PinState, 0, len(m.pins))
	for _, p := range m.pins {
		out = append(out, p.state)
	}
	// Sort by pin number for deterministic ordering.
	sortPinStates(out)
	return out
}

// SetOutput configures pinNum as an output and drives it to level.
func (m *Manager) SetOutput(pinNum int, level gpio.Level) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.pins[pinNum]
	if !ok {
		return fmt.Errorf("gpio: pin %d is not managed", pinNum)
	}
	if p.out == nil {
		return fmt.Errorf("gpio: pin %d does not support output", pinNum)
	}
	if err := p.out.Out(level); err != nil {
		return err
	}
	p.state.Mode = ModeOutput
	p.state.Level = level
	p.state.Func = p.readFunc()
	return nil
}

// SetInput configures pinNum as an input with the given pull direction and
// performs an immediate read to update the stored level.
func (m *Manager) SetInput(pinNum int, pull gpio.Pull) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.pins[pinNum]
	if !ok {
		return fmt.Errorf("gpio: pin %d is not managed", pinNum)
	}
	if p.in == nil {
		return fmt.Errorf("gpio: pin %d does not support input", pinNum)
	}
	if err := p.in.In(pull, gpio.NoEdge); err != nil {
		return err
	}
	p.state.Mode = ModeInput
	p.state.Level = p.in.Read()
	p.state.Func = p.readFunc()
	return nil
}

// Read refreshes the level of an already-configured input pin.
func (m *Manager) Read(pinNum int) (gpio.Level, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.pins[pinNum]
	if !ok {
		return gpio.Low, fmt.Errorf("gpio: pin %d is not managed", pinNum)
	}
	if p.in == nil || p.state.Mode != ModeInput {
		return gpio.Low, fmt.Errorf("gpio: pin %d is not configured as input", pinNum)
	}
	lvl := p.in.Read()
	p.state.Level = lvl
	return lvl, nil
}

// RefreshAll reads the level of every input pin.
func (m *Manager) RefreshAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, p := range m.pins {
		if p.in != nil && p.state.Mode == ModeInput {
			p.state.Level = p.in.Read()
		}
	}
}

// sortPinStates sorts in place by pin number (simple insertion sort).
func sortPinStates(s []PinState) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j].Number < s[j-1].Number; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
