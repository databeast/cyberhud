package driver

import (
	"fmt"
	"strings"
)

// Orientation is a user-facing named display orientation.
// Each panel definition maps these names to hardware-specific MADCTL and offset
// values so users never deal with raw register bytes.
type Orientation string

const (
	// OrientationNormal is the panel's default orientation as designed by the
	// manufacturer (connector at bottom/left).
	OrientationNormal Orientation = "normal"

	// OrientationFlip rotates 180° — for upside-down mounting.
	OrientationFlip Orientation = "flip"

	// OrientationCW rotates 90° clockwise.
	OrientationCW Orientation = "cw"

	// OrientationCCW rotates 90° counter-clockwise.
	OrientationCCW Orientation = "ccw"
)

// OrientationConfig holds the hardware register values for a specific named
// orientation. Panel definitions declare one of these per supported orientation.
type OrientationConfig struct {
	MADCTL    byte
	XOffset   int
	YOffset   int
	Width     int  // override display width (0 = no change)
	Height    int  // override display height (0 = no change)
	Rotation  int  // software rotation (0, 90, 180, 270) applied by FlushPath
	MirrorX   bool // horizontal mirror applied by FlushPath before DrawImage
	Rotate180 bool // hardware 180° rotation for drivers that support it (e.g. SSD1305)
}

// ParseOrientation parses a user-provided orientation string (case-insensitive).
// Returns (OrientationNormal, error) for unrecognized values.
func ParseOrientation(s string) (Orientation, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "normal", "":
		return OrientationNormal, nil
	case "flip", "180":
		return OrientationFlip, nil
	case "cw", "90", "clockwise":
		return OrientationCW, nil
	case "ccw", "270", "counter-clockwise", "counterclockwise":
		return OrientationCCW, nil
	default:
		return OrientationNormal, fmt.Errorf("unrecognized orientation %q (valid: normal, flip, cw, ccw)", s)
	}
}

// ValidOrientations returns the list of valid orientation string values for
// documentation and error messages.
func ValidOrientations() []string {
	return []string{"normal", "flip", "cw", "ccw"}
}
