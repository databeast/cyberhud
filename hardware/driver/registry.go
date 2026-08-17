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
	ID          string
	Title       string
	Summary     string
	Monochrome  bool
	IsEPaper    bool // true for e-ink/e-paper controllers (skip clear-on-exit)
	OptionDefs  []OptionDefinition
	DefaultText textlayout.TextHints
	NewSPI      Factory    // SPI factory (may be nil for I2C-only drivers)
	NewI2C      I2CFactory // I2C factory (may be nil for SPI-only drivers)
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
