package gpio

import (
	"image"
	"image/color"
	"strings"
	"sync"

	"github.com/databeast/cyberhud/display/modes/gpio/source"
	"github.com/databeast/cyberhud/display/modes/gpio/styles"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/led"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

type Policy = source.Policy

func DefaultPolicy() Policy { return source.DefaultPolicy() }

// gpioRegistry is the per-mode StyleRegistry for the gpio display mode.
var gpioRegistry = style.NewRegistry[source.GpioSnapshot, source.Policy](
	// Legacy (5) — first entry is the default returned by gpioRegistry.Default().
	styles.ListStyle,
	styles.IconsStyle,
	styles.DetailStyle,
	styles.DashboardStyle,
	styles.ActivityStyle,
	// MonoSlow
	styles.EinkStyle104x212,
	styles.EinkStyle122x250,
	styles.EinkStyle128x296,
	styles.EinkStyle176x264,
	styles.EinkStyle200x200,
	styles.EinkStyle212x104,
	styles.EinkStyle250x122,
	styles.EinkStyle264x176,
	styles.EinkStyle296x128,
	styles.EinkStyle300x400,
	styles.EinkStyle400x300,
	styles.EinkStyle480x800,
	styles.EinkStyle800x480,
	styles.MonoSlow128x32Style,
	styles.MonoSlow128x64Style,
	styles.MonoSlow128x128Style,
	styles.MonoSlow32x128Style,
	styles.MonoSlow64x128Style,
	styles.MonoSlow160x80Style,
	styles.MonoSlow128x160Style,
	styles.MonoSlow160x128Style,
	styles.MonoSlow135x240Style,
	styles.MonoSlow240x135Style,
	styles.MonoSlow240x240Style,
	styles.MonoSlow240x320Style,
	styles.MonoSlow320x240Style,
	styles.MonoSlow320x480Style,
	styles.MonoSlow480x320Style,
	// MonoFast
	styles.MonoStyle128x32,
	styles.MonoStyle128x64,
	styles.MonoStyle128x128,
	styles.MonoStyle32x128,
	styles.MonoStyle64x128,
	styles.MonoFast160x80Style,
	styles.MonoFast128x160Style,
	styles.MonoFast160x128Style,
	styles.MonoFast135x240Style,
	styles.MonoFast240x135Style,
	styles.MonoFast240x240Style,
	styles.MonoFast240x320Style,
	styles.MonoFast320x240Style,
	styles.MonoFast320x480Style,
	styles.MonoFast480x320Style,
	styles.MonoFast480x800Style,
	styles.MonoFast800x480Style,
	styles.MonoFast104x212Style,
	styles.MonoFast122x250Style,
	styles.MonoFast128x296Style,
	styles.MonoFast176x264Style,
	styles.MonoFast200x200Style,
	styles.MonoFast212x104Style,
	styles.MonoFast250x122Style,
	styles.MonoFast264x176Style,
	styles.MonoFast296x128Style,
	styles.MonoFast300x400Style,
	styles.MonoFast400x300Style,
	// GrayscaleSlow
	styles.GrayscaleSlow128x32Style,
	styles.GrayscaleSlow128x64Style,
	styles.GrayscaleSlow128x128Style,
	styles.GrayscaleSlow32x128Style,
	styles.GrayscaleSlow64x128Style,
	styles.GrayscaleSlow160x80Style,
	styles.GrayscaleSlow128x160Style,
	styles.GrayscaleSlow160x128Style,
	styles.GrayscaleSlow135x240Style,
	styles.GrayscaleSlow240x135Style,
	styles.GrayscaleSlow240x240Style,
	styles.GrayscaleSlow240x320Style,
	styles.GrayscaleSlow320x240Style,
	styles.GrayscaleSlow320x480Style,
	styles.GrayscaleSlow480x320Style,
	styles.GrayscaleSlow480x800Style,
	styles.GrayscaleSlow800x480Style,
	styles.GrayscaleSlow104x212Style,
	styles.GrayscaleSlow122x250Style,
	styles.GrayscaleSlow128x296Style,
	styles.GrayscaleSlow176x264Style,
	styles.GrayscaleSlow200x200Style,
	styles.GrayscaleSlow212x104Style,
	styles.GrayscaleSlow250x122Style,
	styles.GrayscaleSlow264x176Style,
	styles.GrayscaleSlow296x128Style,
	styles.GrayscaleSlow300x400Style,
	styles.GrayscaleSlow400x300Style,
	// GrayscaleFast
	styles.GrayscaleFast128x32Style,
	styles.GrayscaleFast128x64Style,
	styles.GrayscaleFast128x128Style,
	styles.GrayscaleFast32x128Style,
	styles.GrayscaleFast64x128Style,
	styles.GrayscaleFast160x80Style,
	styles.GrayscaleFast128x160Style,
	styles.GrayscaleFast160x128Style,
	styles.GrayscaleFast135x240Style,
	styles.GrayscaleFast240x135Style,
	styles.GrayscaleFast240x240Style,
	styles.GrayscaleFast240x320Style,
	styles.GrayscaleFast320x240Style,
	styles.GrayscaleFast320x480Style,
	styles.GrayscaleFast480x320Style,
	styles.GrayscaleFast480x800Style,
	styles.GrayscaleFast800x480Style,
	styles.GrayscaleFast104x212Style,
	styles.GrayscaleFast122x250Style,
	styles.GrayscaleFast128x296Style,
	styles.GrayscaleFast176x264Style,
	styles.GrayscaleFast200x200Style,
	styles.GrayscaleFast212x104Style,
	styles.GrayscaleFast250x122Style,
	styles.GrayscaleFast264x176Style,
	styles.GrayscaleFast296x128Style,
	styles.GrayscaleFast300x400Style,
	styles.GrayscaleFast400x300Style,
	// ColorSlow
	styles.ColorSlow128x32Style,
	styles.ColorSlow128x64Style,
	styles.ColorSlow128x128Style,
	styles.ColorSlow32x128Style,
	styles.ColorSlow64x128Style,
	styles.ColorSlow160x80Style,
	styles.ColorSlow128x160Style,
	styles.ColorSlow160x128Style,
	styles.ColorSlow135x240Style,
	styles.ColorSlow240x135Style,
	styles.ColorSlow240x240Style,
	styles.ColorSlow240x320Style,
	styles.ColorSlow320x240Style,
	styles.ColorSlow320x480Style,
	styles.ColorSlow480x320Style,
	styles.ColorSlow480x800Style,
	styles.ColorSlow800x480Style,
	styles.ColorSlow104x212Style,
	styles.ColorSlow122x250Style,
	styles.ColorSlow128x296Style,
	styles.ColorSlow176x264Style,
	styles.ColorSlow200x200Style,
	styles.ColorSlow212x104Style,
	styles.ColorSlow250x122Style,
	styles.ColorSlow264x176Style,
	styles.ColorSlow296x128Style,
	styles.ColorSlow300x400Style,
	styles.ColorSlow400x300Style,
	// ColorFast
	styles.ColorStyle160x80,
	styles.ColorStyle128x128,
	styles.ColorStyle128x160,
	styles.ColorStyle160x128,
	styles.ColorStyle240x135,
	styles.ColorStyle135x240,
	styles.ColorStyle240x240,
	styles.ColorStyle240x320,
	styles.ColorStyle320x240,
	styles.ColorStyle320x480,
	styles.ColorStyle480x320,
	styles.ColorStyle480x800,
	styles.ColorStyle800x480,
	styles.ColorFast128x32Style,
	styles.ColorFast128x64Style,
	styles.ColorFast32x128Style,
	styles.ColorFast64x128Style,
	styles.ColorFast104x212Style,
	styles.ColorFast122x250Style,
	styles.ColorFast128x296Style,
	styles.ColorFast176x264Style,
	styles.ColorFast200x200Style,
	styles.ColorFast212x104Style,
	styles.ColorFast250x122Style,
	styles.ColorFast264x176Style,
	styles.ColorFast296x128Style,
	styles.ColorFast300x400Style,
	styles.ColorFast400x300Style,
)

// SetGPIOManager registers the GPIO manager used by gpio instances.
func SetGPIOManager(mgr interface{ Snapshot() []gpiomgr.PinState }) {
	source.SetGPIOManager(mgr)
}

// registeredStyleNames returns the list of style names from the registry.
// Used by catalog registration and cmdHandler for allowed-value validation.
func registeredStyleNames() []string {
	styles := gpioRegistry.Enumerate()
	names := make([]string, len(styles))
	for i, s := range styles {
		names[i] = s.Name()
	}
	return names
}

// Allowed style values for the GPIO mode.
const (
	StyleList      = "list"
	StyleIcons     = "icons"
	StyleDetail    = "detail"
	StyleDashboard = "dashboard"
	StyleActivity  = "activity"
)

// Colors for pin level indication, sourced from the shared color package.
var (
	ColorHigh = sharedcolor.GPIOPalette.Active   // green for HIGH
	ColorLow  = sharedcolor.GPIOPalette.Inactive // grey for LOW
)

// allowedFGColors lists all valid fgcolor policy values.
var allowedFGColors = []string{"cyan", "green", "amber", "red", "white", "none"}

// AllowedFGColorValues returns a copy of the valid fgcolor values.
// Exported for use in property-based tests.
//
// Framework pattern demonstrated: FGColor — exposes the canonical foreground
// color palette for test generators and validation.
func AllowedFGColorValues() []string {
	out := make([]string, len(allowedFGColors))
	copy(out, allowedFGColors)
	return out
}

// policyState holds the mutex-protected active policy.
var policyState = struct {
	sync.RWMutex
	policy Policy
}{
	policy: DefaultPolicy(),
}

// GetPolicy returns the current GPIO policy (thread-safe read).
func GetPolicy() Policy {
	policyState.RLock()
	defer policyState.RUnlock()
	return policyState.policy
}

// SetPolicy updates the GPIO policy after normalization (thread-safe write).
// Invalid fields are reset to defaults by normalizePolicy.
func SetPolicy(p Policy) {
	policyState.Lock()
	defer policyState.Unlock()
	policyState.policy = normalizePolicy(p)
}

// isAllowedAccent checks if value is present in the allowedFGColors list.
func isAllowedAccent(value string) bool {
	for _, a := range allowedFGColors {
		if value == a {
			return true
		}
	}
	return false
}

// fontValidator rejects whitespace-only values; accepts "auto" and any non-empty trimmed string.
func fontValidator(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return "must be \"auto\" or a non-empty font ID"
	}
	return ""
}

// boolStr returns "true" or "false" for a bool value.
func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// sortInts sorts a slice of ints in ascending order (insertion sort for small slices).
func sortInts(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j-1] > a[j]; j-- {
			a[j-1], a[j] = a[j], a[j-1]
		}
	}
}

// BuildItems returns GPIO mode row text for the primary list screen.
// Each entry is the full string representation of a PinState.
func BuildItems(pins []gpiomgr.PinState) []string {
	items := make([]string, len(pins))
	for i, p := range pins {
		items[i] = p.String()
	}
	return items
}

// BuildItemsTruncated returns GPIO mode rows truncated to the given max character width.
func BuildItemsTruncated(pins []gpiomgr.PinState, maxChars int) []string {
	items := make([]string, len(pins))
	for i, p := range pins {
		s := p.String()
		if maxChars > 0 && len(s) > maxChars {
			s = s[:maxChars]
		}
		items[i] = s
	}
	return items
}

// BuildColors returns a Colors slice mapping pin levels to accent-derived colors.
// HIGH pins receive the active palette color; LOW pins receive the inactive color.
// If colorEnabled is false, returns nil.
//
// Framework pattern demonstrated: FGColor integration — pin-level color
// mapping using the shared color palette.
func BuildColors(pins []gpiomgr.PinState, colorEnabled bool) []color.Color {
	states := make([]bool, len(pins))
	for i, p := range pins {
		states[i] = bool(p.Level)
	}
	return sharedcolor.BuildSlice(states, sharedcolor.GPIOPalette, colorEnabled)
}

// BuildSprites returns sprites for output-mode pins using the Compositor pattern.
// Input-only pins do not receive a sprite.
// The sprite is positioned at (0, row * rowHeight).
//
// Framework pattern demonstrated: Compositor — widget sprite compositing with
// conditional inclusion based on pin mode.
func BuildSprites(pins []gpiomgr.PinState, rowHeight int, glyphWidth, glyphHeight int) []widgets.Sprite {
	diameter := led.DiameterForRow(rowHeight)

	ctx := widgets.SuppressionContext{
		AvailableWidth:  diameter,
		AvailableHeight: rowHeight * len(pins),
	}
	comp := widgets.NewCompositor(ctx)

	for i, p := range pins {
		state := led.Off
		if p.Level {
			state = led.On
		}
		comp.AddIf(p.Mode == gpiomgr.ModeOutput, led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(0, i*rowHeight)},
			Foreground: ColorHigh,
		}))
	}
	if len(comp.Sprites()) == 0 {
		return nil
	}
	return comp.Sprites()
}

// BuildIconGrid returns sprites arranged as a grid of LED indicators for "icons" style.
// Each pin gets an LED sprite (On for HIGH, Off for LOW) arranged in row-major order
// within the given panel bounds.
//
// Framework pattern demonstrated: Compositor — grid-based widget sprite compositing
// with dynamic column computation.
func BuildIconGrid(pins []gpiomgr.PinState, pixelWidth, pixelHeight, glyphWidth, glyphHeight int) []widgets.Sprite {
	// Icon cell size based on glyph metrics (minimum 8x8).
	cellW := glyphWidth
	if cellW < 8 {
		cellW = 8
	}
	cellH := glyphHeight
	if cellH < 8 {
		cellH = 8
	}

	cols := 1
	if pixelWidth > 0 && cellW > 0 {
		cols = pixelWidth / cellW
	}
	if cols < 1 {
		cols = 1
	}

	diameter := cellW
	if cellH < diameter {
		diameter = cellH
	}

	ctx := widgets.SuppressionContext{
		AvailableWidth:  pixelWidth,
		AvailableHeight: pixelHeight,
	}
	comp := widgets.NewCompositor(ctx)

	for i, p := range pins {
		row := i / cols
		col := i % cols
		x := col * cellW
		y := row * cellH

		// Stop if we're past the vertical bounds.
		if pixelHeight > 0 && y+cellH > pixelHeight {
			break
		}

		state := led.Off
		if p.Level {
			state = led.On
		}
		fg := ColorLow
		if p.Level {
			fg = ColorHigh
		}
		comp.Add(led.New(led.Config{
			State:      state,
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(x, y)},
			Foreground: fg,
		}))
	}
	return comp.Sprites()
}

// Summary reports compact counters for small secondary displays.
// Returns total pin count, output pin count, and HIGH pin count.
func Summary(pins []gpiomgr.PinState) (count, outputs, high int) {
	count = len(pins)
	for _, p := range pins {
		if p.Mode == gpiomgr.ModeOutput {
			outputs++
		}
		if p.Level {
			high++
		}
	}
	return count, outputs, high
}
