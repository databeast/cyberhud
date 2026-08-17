// Package color provides shared color primitives for all display modes.
//
// It exposes a named accent palette, lookup functions, and color manipulation
// utilities. The package depends only on image/color from the standard library
// and imports no other cyberhud module packages.
package color

import (
	"image/color"
	"sort"
)

// palette is the immutable registry of accent colors.
// Concurrency safety is achieved through immutability (no writes after init).
var palette = map[string]color.RGBA{
	"cyan":    {0, 255, 255, 255},
	"green":   {0, 200, 0, 255},
	"emerald": {0, 105, 55, 255},
	"amber":   {255, 191, 0, 255},
	"red":     {255, 0, 0, 255},
	"white":   {255, 255, 255, 255},
}

// Names returns all registered palette entry names as a sorted string slice.
func Names() []string {
	names := make([]string, 0, len(palette))
	for name := range palette {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Lookup returns the RGBA for a registered name, or opaque white for
// unregistered or empty names.
func Lookup(name string) color.RGBA {
	if c, ok := palette[name]; ok {
		return c
	}
	return color.RGBA{255, 255, 255, 255}
}

// ResolveAccent maps a named accent (or "none") to its RGBA value.
// Unrecognized and "none" values resolve to opaque white {255, 255, 255, 255}.
func ResolveAccent(name string) color.RGBA {
	if name == "none" {
		return color.RGBA{255, 255, 255, 255}
	}
	return Lookup(name)
}

// Dim returns a dimmed variant: each RGB channel halved via integer division,
// alpha forced to 255.
func Dim(c color.RGBA) color.RGBA {
	return color.RGBA{
		R: c.R / 2,
		G: c.G / 2,
		B: c.B / 2,
		A: 255,
	}
}

// BinaryPalette maps boolean states to two colors.
type BinaryPalette struct {
	Active   color.RGBA
	Inactive color.RGBA
}

// NewBinaryPalette constructs a BinaryPalette from active and inactive colors.
func NewBinaryPalette(active, inactive color.RGBA) BinaryPalette {
	return BinaryPalette{Active: active, Inactive: inactive}
}

// Select returns Active when state is true, Inactive when false.
func (bp BinaryPalette) Select(state bool) color.RGBA {
	if state {
		return bp.Active
	}
	return bp.Inactive
}

// GPIOPalette is the standard HIGH/LOW color pair for GPIO pins.
var GPIOPalette = BinaryPalette{
	Active:   color.RGBA{0x00, 0xCC, 0x44, 0xFF}, // green
	Inactive: color.RGBA{0x66, 0x66, 0x66, 0xFF}, // grey
}

// BuildSlice maps a slice of boolean states through a BinaryPalette.
// Returns nil when enabled is false. Returns an empty non-nil slice when states
// is empty and enabled is true.
func BuildSlice(states []bool, palette BinaryPalette, enabled bool) []color.Color {
	if !enabled {
		return nil
	}
	result := make([]color.Color, len(states))
	for i, s := range states {
		result[i] = palette.Select(s)
	}
	return result
}
