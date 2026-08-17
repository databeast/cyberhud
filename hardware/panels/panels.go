package panels

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/hardware/driver"
)

// InputPins defines optional GPIO input pin assignments for a panel product.
type InputPins struct {
	Key1       string
	Key2       string
	Key3       string
	JoyUp      string
	JoyDown    string
	JoyLeft    string
	JoyRight   string
	JoyPressed string
}

// Any reports whether any input pin is assigned.
func (p InputPins) Any() bool {
	return p.Key1 != "" || p.Key2 != "" || p.Key3 != "" ||
		p.JoyUp != "" || p.JoyDown != "" || p.JoyLeft != "" ||
		p.JoyRight != "" || p.JoyPressed != ""
}

// Screen describes one display in a potentially multi-display panel product.
type Screen struct {
	Index         int
	Name          string
	Controller    string
	SPI           string
	ExcludedModes []string
	DefaultMode   string
	InputEnabled  bool

	Config driver.DriverConfig

	DCPin        string
	RSTPin       string
	BusyPin      string
	BLPin        string
	I2CBus       string // e.g. "/dev/i2c-1"; empty = system default
	Orientations map[driver.Orientation]driver.OrientationConfig

	// XPosition is the explicit horizontal position of this screen within the
	// virtual display coordinate space. When set (>=0 on any screen), buildPositions
	// uses declared positions instead of accumulating left-to-right by index.
	// -1 means "auto" (use left-to-right accumulation).
	XPosition int

	// YPosition is the explicit vertical position of this screen within the
	// virtual display coordinate space. Defaults to 0 (top-aligned).
	YPosition int

	// Rotation declares the clockwise rotation (0, 90, 180, 270) that FlushPath
	// applies in software before sending pixels to hardware. Use when the logical
	// display orientation differs from the hardware's native pixel scan order.
	// The VD region uses logical dimensions; FlushPath rotates to hardware dimensions.
	Rotation int

	// MirrorX applies a horizontal mirror (left-right flip) in software before
	// sending pixels to hardware. Used when hardware MADCTL orientation introduces
	// an unwanted horizontal mirror that cannot be corrected via MADCTL alone.
	MirrorX bool

	// PPI is the optional per-screen pixels-per-inch override. When greater than
	// zero, this value takes precedence over the parent Definition's PPI for this
	// screen's TextHints. Zero means "use parent PPI."
	PPI float64
}

// Definition describes one board-level panel product.
type Definition struct {
	Name          string
	Description   string
	Controller    string
	Monochrome    bool
	ExcludedModes []string
	DefaultMode   string
	InputEnabled  bool
	Virtual       []Screen
	Layout        *region.RegionLayout // optional; nil means use default generation
	Orientations  map[driver.Orientation]driver.OrientationConfig

	Config driver.DriverConfig

	DCPin   string
	RSTPin  string
	BusyPin string
	BLPin   string
	I2CBus  string // e.g. "/dev/i2c-1"; empty = system default
	Inputs  InputPins

	// PPI is the optional pixels-per-inch for this panel product. Zero means
	// "undeclared." Negative values are clamped to zero during registration.
	PPI float64
}

var (
	defsMu sync.RWMutex
	defs   = map[string]Definition{}
)

func normalize(def Definition) Definition {
	def.Name = strings.ToLower(strings.TrimSpace(def.Name))
	if meta, ok := driver.Get(def.Controller); ok {
		def.Monochrome = meta.Monochrome
	}
	if def.PPI < 0 {
		def.PPI = 0
	}
	return def
}

// Register publishes one panel product definition.
func Register(def Definition) {
	def = normalize(def)
	if def.Name == "" {
		return
	}
	defsMu.Lock()
	defer defsMu.Unlock()
	defs[def.Name] = def
}

// Get returns one built-in panel definition by name.
func Get(name string) (Definition, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	defsMu.RLock()
	def, ok := defs[name]
	defsMu.RUnlock()
	if !ok {
		return Definition{}, fmt.Errorf("unknown panel %q (available: %s)", name, Supported())
	}
	return normalize(def), nil
}

// All returns a copy of all registered panel definitions keyed by name.
func All() map[string]Definition {
	defsMu.RLock()
	defer defsMu.RUnlock()
	out := make(map[string]Definition, len(defs))
	for name, def := range defs {
		out[name] = normalize(def)
	}
	return out
}

// Names returns supported panel names sorted ascending.
func Names() []string {
	all := All()
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Supported returns all panel names as CSV.
func Supported() string {
	return strings.Join(Names(), ", ")
}

// CalcPPI computes pixels-per-inch from resolution and physical diagonal size.
// Returns zero when any input is zero or negative.
func CalcPPI(pixelWidth, pixelHeight int, diagonalInches float64) float64 {
	if pixelWidth <= 0 || pixelHeight <= 0 || diagonalInches <= 0 {
		return 0
	}
	diagPixels := math.Sqrt(float64(pixelWidth*pixelWidth + pixelHeight*pixelHeight))
	return diagPixels / diagonalInches
}
