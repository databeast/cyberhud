package source

import gpiomgr "github.com/databeast/cyberhud/hardware/gpio"

// PinClass is the coarse role of a physical pin on the 40-pin connector, used
// to decide how a pad is drawn on the board-map view.
type PinClass uint8

const (
	ClassGPIO   PinClass = iota // a general-purpose line (may be in/out/alt)
	ClassPower                  // 3V3 or 5V rail
	ClassGround                 // ground
)

// HeaderPin is one physical position on the 40-pin connector: where it sits
// (Phys, 1..40), what it fundamentally is (Class), and for GPIO lines the BCM
// number used to join it to live state.
type HeaderPin struct {
	Phys  int
	Class PinClass
	BCM   int // valid only when Class == ClassGPIO
}

// header40 is the canonical Raspberry Pi 40-pin physical layout, identical
// across every 40-pin model (B+ through Pi 5). Index is physical pin order
// 1..40. Odd physical pins form one row, even pins the other; physical pin 1 is
// the corner. This is a fixed hardware fact, not configuration — but it is
// cross-checked against the live driver's reported BCM numbers at render time
// (see BoardPads), so a wrong entry here surfaces rather than silently
// mislabels a pad.
//
// Power/ground positions and the physical→BCM mapping follow the standard Pi
// pinout (e.g. physical pin 11 = GPIO17, physical pin 13 = GPIO27).
var header40 = []HeaderPin{
	{Phys: 1, Class: ClassPower},          // 3V3
	{Phys: 2, Class: ClassPower},          // 5V
	{Phys: 3, Class: ClassGPIO, BCM: 2},   // SDA1 / I2C1
	{Phys: 4, Class: ClassPower},          // 5V
	{Phys: 5, Class: ClassGPIO, BCM: 3},   // SCL1 / I2C1
	{Phys: 6, Class: ClassGround},         //
	{Phys: 7, Class: ClassGPIO, BCM: 4},   //
	{Phys: 8, Class: ClassGPIO, BCM: 14},  // TXD / UART0
	{Phys: 9, Class: ClassGround},         //
	{Phys: 10, Class: ClassGPIO, BCM: 15}, // RXD / UART0
	{Phys: 11, Class: ClassGPIO, BCM: 17}, //
	{Phys: 12, Class: ClassGPIO, BCM: 18}, // PCM_CLK / PWM0
	{Phys: 13, Class: ClassGPIO, BCM: 27}, //
	{Phys: 14, Class: ClassGround},        //
	{Phys: 15, Class: ClassGPIO, BCM: 22}, //
	{Phys: 16, Class: ClassGPIO, BCM: 23}, //
	{Phys: 17, Class: ClassPower},         // 3V3
	{Phys: 18, Class: ClassGPIO, BCM: 24}, //
	{Phys: 19, Class: ClassGPIO, BCM: 10}, // MOSI / SPI0
	{Phys: 20, Class: ClassGround},        //
	{Phys: 21, Class: ClassGPIO, BCM: 9},  // MISO / SPI0
	{Phys: 22, Class: ClassGPIO, BCM: 25}, //
	{Phys: 23, Class: ClassGPIO, BCM: 11}, // SCLK / SPI0
	{Phys: 24, Class: ClassGPIO, BCM: 8},  // CE0 / SPI0
	{Phys: 25, Class: ClassGround},        //
	{Phys: 26, Class: ClassGPIO, BCM: 7},  // CE1 / SPI0
	{Phys: 27, Class: ClassGPIO, BCM: 0},  // ID_SD
	{Phys: 28, Class: ClassGPIO, BCM: 1},  // ID_SC
	{Phys: 29, Class: ClassGPIO, BCM: 5},  //
	{Phys: 30, Class: ClassGround},        //
	{Phys: 31, Class: ClassGPIO, BCM: 6},  //
	{Phys: 32, Class: ClassGPIO, BCM: 12}, // PWM0
	{Phys: 33, Class: ClassGPIO, BCM: 13}, // PWM1
	{Phys: 34, Class: ClassGround},        //
	{Phys: 35, Class: ClassGPIO, BCM: 19}, //
	{Phys: 36, Class: ClassGPIO, BCM: 16}, //
	{Phys: 37, Class: ClassGPIO, BCM: 26}, //
	{Phys: 38, Class: ClassGPIO, BCM: 20}, //
	{Phys: 39, Class: ClassGround},        //
	{Phys: 40, Class: ClassGPIO, BCM: 21}, //
}

// Header40 returns the canonical physical connector layout (physical order).
func Header40() []HeaderPin {
	out := make([]HeaderPin, len(header40))
	copy(out, header40)
	return out
}

// PadState is one physical position prepared for the board-map render: its
// fixed geometry (Phys, Class) merged with whatever live state the driver
// reports for that position. For non-GPIO positions, only Class is meaningful.
type PadState struct {
	Phys  int
	Class PinClass

	// For ClassGPIO positions:
	Configured bool // true when the driver reports a real mode for this pin
	Mode       gpiomgr.PinMode
	Level      bool   // meaningful when Mode is Input/Output
	Bus        string // bus family (SPI0, I2C1, …) or "" — from BusGroup()
}

// Lit reports whether the pad should render as "active" (bright) vs quiet:
// a driven-HIGH output, a HIGH input, or any pin serving a bus is live.
func (p PadState) Lit() bool {
	if p.Bus != "" {
		return true
	}
	return p.Configured && p.Level
}

// BoardPads overlays a live pin snapshot onto the fixed 40-pin layout, yielding
// one PadState per physical position (1..40) in physical order. GPIO positions
// are matched to the snapshot by BCM number; positions the snapshot does not
// cover render as unconfigured (ghosted) GPIO, and power/ground positions carry
// only their class.
//
// This is also the integrity check for header40: every GPIO PinState in the
// snapshot should find a home at exactly one physical position. A mismatch
// (a BCM in the data with no physical slot) is reported via the second return —
// callers may surface it rather than silently dropping the pin.
func BoardPads(pins []gpiomgr.PinState) (pads []PadState, unmatched []int) {
	byBCM := make(map[int]gpiomgr.PinState, len(pins))
	for _, p := range pins {
		byBCM[p.Number] = p
	}
	matched := make(map[int]bool, len(pins))

	pads = make([]PadState, 0, len(header40))
	for _, hp := range header40 {
		pad := PadState{Phys: hp.Phys, Class: hp.Class}
		if hp.Class == ClassGPIO {
			if st, ok := byBCM[hp.BCM]; ok {
				matched[hp.BCM] = true
				pad.Configured = st.Mode != gpiomgr.ModeUnknown
				pad.Mode = st.Mode
				pad.Level = bool(st.Level)
				pad.Bus = st.BusGroup()
			}
		}
		pads = append(pads, pad)
	}

	for _, p := range pins {
		if p.Mode == gpiomgr.ModeUnknown {
			continue
		}
		if !matched[p.Number] {
			unmatched = append(unmatched, p.Number)
		}
	}
	return pads, unmatched
}
