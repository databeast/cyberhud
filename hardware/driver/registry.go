package driver

import (
	"sort"
	"strings"
	"sync"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"

	"github.com/databeast/cyberhud/display/surface/textlayout"
)

// OptionDefinition describes one published driver configuration option.
type OptionDefinition struct {
	Key     string
	Type    string
	Summary string
	Default string
}

// ColorDepth is the native per-pixel color depth a driver's DrawImage actually
// drives on the physical panel. It is the hardware ground truth, independent of
// the refresh-speed dimension that Capability also encodes. Every driver MUST
// declare its ColorDepth; the capability-conformance test asserts that the
// capability the driver advertises at runtime (via TextHints) agrees with the
// depth declared here, so a panel can never silently understate itself (e.g. a
// color OLED defaulting to mono because it forgot to implement TextHints).
type ColorDepth int

const (
	// DepthUnset is the zero value. It is invalid: a driver that leaves ColorDepth
	// unset fails the conformance test rather than silently defaulting to mono.
	DepthUnset ColorDepth = iota
	// DepthMono1 is 1-bit black/white (thresholded). The panel shows two ink states.
	DepthMono1
	// DepthGray is multi-level luminance with no hue (grayscale e-paper, PWM LED
	// matrices). The panel shows a brightness ramp.
	DepthGray
	// DepthColor is full RGB or any multi-hue output (RGB TFT/OLED, tri-color e-ink).
	// Tri-color (black/white/red) counts as color because it carries hue the mono
	// and gray depths cannot.
	DepthColor
)

// String renders the depth as a short label.
func (d ColorDepth) String() string {
	switch d {
	case DepthMono1:
		return "mono1"
	case DepthGray:
		return "gray"
	case DepthColor:
		return "color"
	default:
		return "unset"
	}
}

// DepthForCapability maps a textlayout capability value to the color depth it
// implies. This is the single authority relating the two axes; the conformance
// test uses it to check that a driver's advertised capability is consistent with
// its declared ColorDepth. The refresh-speed half of the capability (slow/fast)
// is intentionally ignored here — depth is the only axis ColorDepth constrains.
func DepthForCapability(capability int) ColorDepth {
	switch capability {
	case textlayout.CapMonoSlow, textlayout.CapMonoFast:
		return DepthMono1
	case textlayout.CapGrayscaleSlow, textlayout.CapGrayscaleFast:
		return DepthGray
	case textlayout.CapColorSlow, textlayout.CapColorFast:
		return DepthColor
	default:
		return DepthUnset
	}
}

// DriverConfig is a driver-agnostic config used by runtime and panel definitions.
type DriverConfig struct {
	Width            int
	Height           int
	SPIHz            physic.Frequency
	MADCTL           byte
	XOffset          int
	YOffset          int
	ColOffset        int
	BusyHigh         bool
	BusyTimeout      time.Duration
	FullRefreshEvery int
	I2CAddr          uint16
	Layout           string // named pixel-address layout for non-linear matrices (driver-specific)
	Rotate180        bool
	InvertColors     bool // when true, driver sends INVON instead of INVOFF during init
}

// Factory creates a configured draw target from an SPI port.
type Factory func(port spi.Port, dc, rst, bl gpio.PinOut, busy gpio.PinIn, cfg DriverConfig) (DrawTarget, error)

// I2CFactory creates a configured draw target from an I2C bus.
type I2CFactory func(bus i2c.Bus, cfg DriverConfig) (DrawTarget, error)

// Definition describes a registered driver and its factory.
type Definition struct {
	NewSPI     Factory    // SPI factory (may be nil for I2C-only drivers)
	NewI2C     I2CFactory // I2C factory (may be nil for SPI-only drivers)
	ID         string
	Title      string
	Summary    string
	OptionDefs []OptionDefinition
	// ColorDepth is the native per-pixel depth this driver's DrawImage drives.
	// Required: a driver that leaves it DepthUnset fails the conformance test.
	// It is the hardware truth against which the advertised Capability is checked.
	ColorDepth ColorDepth
	// Capability is the single declared capability this driver advertises to style
	// selection (ordered tier: mono/gray/color × slow/fast). Required: a driver
	// that leaves it unset (CapMonoSlow zero value) while declaring a richer
	// ColorDepth fails the conformance test.
	Capability int
	Monochrome bool
	IsEPaper   bool // true for e-ink/e-paper controllers (skip clear-on-exit)
}

var (
	defsMu sync.RWMutex
	defs   = map[string]Definition{}
)

// Register publishes a driver definition and its constructor.
func Register(def Definition) {
	def.ID = strings.ToLower(strings.TrimSpace(def.ID))
	if def.ID == "" || (def.NewSPI == nil && def.NewI2C == nil) {
		return
	}
	defsMu.Lock()
	defer defsMu.Unlock()
	defs[def.ID] = def
}

// Get returns one driver definition by ID.
func Get(id string) (Definition, bool) {
	id = strings.ToLower(strings.TrimSpace(id))
	defsMu.RLock()
	defer defsMu.RUnlock()
	def, ok := defs[id]
	return def, ok
}

// Names returns all registered driver IDs sorted ascending.
func Names() []string {
	defsMu.RLock()
	defer defsMu.RUnlock()
	ids := make([]string, 0, len(defs))
	for id := range defs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
