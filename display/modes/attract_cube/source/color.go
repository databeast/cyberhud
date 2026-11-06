package source

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"

	accent "github.com/databeast/cyberhud/display/style/color"
)

// ParseBaseColor accepts a shared named accent or a six-digit RGB hex value.
func ParseBaseColor(value string) (string, color.RGBA, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, name := range accent.Names() {
		if value == name {
			return value, accent.Lookup(name), nil
		}
	}
	if len(value) == 7 && value[0] == '#' {
		valid := true
		for _, c := range value[1:] {
			if !(c >= '0' && c <= '9') && !(c >= 'a' && c <= 'f') {
				valid = false
				break
			}
		}
		if valid {
			n, err := strconv.ParseUint(value[1:], 16, 24)
			if err == nil {
				return value, color.RGBA{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n), A: 255}, nil
			}
		}
	}
	return "", color.RGBA{}, fmt.Errorf("base_color must be a named accent (%s) or #RRGGBB", strings.Join(accent.Names(), ", "))
}

func rgbToHSV(r, g, b float64) (h, s, v float64) {
	maximum := math.Max(r, math.Max(g, b))
	minimum := math.Min(r, math.Min(g, b))
	delta := maximum - minimum
	v = maximum
	if maximum == 0 || delta == 0 {
		return 0, 0, v
	}
	s = delta / maximum
	switch maximum {
	case r:
		h = (g - b) / delta
	case g:
		h = 2 + (b-r)/delta
	default:
		h = 4 + (r-g)/delta
	}
	h = math.Mod(h*60+360, 360)
	return
}

func hsvToRGB(h, s, v float64) (r, g, b float64) {
	h = math.Mod(h+360, 360) / 60
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h, 2)-1))
	m := v - c
	switch {
	case h < 1:
		r, g = c, x
	case h < 2:
		r, g = x, c
	case h < 3:
		g, b = c, x
	case h < 4:
		g, b = x, c
	case h < 5:
		r, b = x, c
	default:
		r, b = c, x
	}
	return r + m, g + m, b + m
}

func tinted(base color.RGBA, r, g, b float64) (float64, float64, float64) {
	h, s, v := rgbToHSV(r, g, b)
	baseH, baseS, baseV := rgbToHSV(float64(base.R)/255, float64(base.G)/255, float64(base.B)/255)
	return hsvToRGB(h+baseH-180, s*baseS, v*baseV)
}
