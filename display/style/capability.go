package style

import "fmt"

// Capability encodes a panel's hardware capability level as a single ordered value.
// The ordering is: MonoSlow < MonoFast < GrayscaleSlow < GrayscaleFast < ColorSlow < ColorFast
type Capability int

const (
	MonoSlow      Capability = iota // 1-bit, slow refresh (e-paper mono)
	MonoFast                        // 1-bit, fast refresh (OLED mono)
	GrayscaleSlow                   // Multi-level luminance, slow refresh (grayscale e-paper)
	GrayscaleFast                   // Multi-level luminance, fast refresh (grayscale LED matrix)
	ColorSlow                       // RGB, slow refresh (color e-paper)
	ColorFast                       // RGB, fast refresh (color TFT)
)

// String returns the kebab-case representation of the Capability.
func (c Capability) String() string {
	switch c {
	case MonoSlow:
		return "mono-slow"
	case MonoFast:
		return "mono-fast"
	case GrayscaleSlow:
		return "grayscale-slow"
	case GrayscaleFast:
		return "grayscale-fast"
	case ColorSlow:
		return "color-slow"
	case ColorFast:
		return "color-fast"
	default:
		return fmt.Sprintf("capability(%d)", int(c))
	}
}

// ParseCapability parses a kebab-case capability string into its Capability constant.
// Returns (MonoSlow, error) for unrecognized values.
func ParseCapability(s string) (Capability, error) {
	switch s {
	case "mono-slow":
		return MonoSlow, nil
	case "mono-fast":
		return MonoFast, nil
	case "grayscale-slow":
		return GrayscaleSlow, nil
	case "grayscale-fast":
		return GrayscaleFast, nil
	case "color-slow":
		return ColorSlow, nil
	case "color-fast":
		return ColorFast, nil
	default:
		return MonoSlow, fmt.Errorf("unrecognized capability: %q", s)
	}
}
