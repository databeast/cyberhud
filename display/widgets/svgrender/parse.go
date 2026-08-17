package svgrender

import (
	"fmt"
	"strings"

	"github.com/srwiley/oksvg"
)

// parse parses an SVG string and returns an *oksvg.SvgIcon scaled to fit
// within the given width and height, preserving aspect ratio. It recovers
// from any panics raised by oksvg, returning nil and an error instead.
func parse(svg string, w, h int) (icon *oksvg.SvgIcon, err error) {
	defer func() {
		if r := recover(); r != nil {
			icon = nil
			err = fmt.Errorf("svgrender: parse panic: %v", r)
		}
	}()

	icon, err = oksvg.ReadIconStream(strings.NewReader(svg))
	if err != nil {
		return nil, fmt.Errorf("svgrender: parse error: %w", err)
	}

	icon.SetTarget(0, 0, float64(w), float64(h))
	return icon, nil
}
