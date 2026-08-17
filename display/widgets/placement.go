package widgets

import (
	"image"

	"github.com/databeast/cyberhud/display/style/layout"
)

// PlaceTopRight returns the bounding rectangle for an element anchored at the
// top-right of the content area. Returns image.Rectangle{} (zero value) when
// the element dimensions exceed the available content area.
func PlaceTopRight(bridge layout.LayoutCalculator, elementWidth, elementHeight int) image.Rectangle {
	w := bridge.AvailableContentWidth()
	h := bridge.AvailableContentHeight()
	if elementWidth > w || elementHeight > h {
		return image.Rectangle{}
	}

	originX, originY := bridge.ContentOrigin()

	return image.Rectangle{
		Min: image.Point{X: originX + w - elementWidth, Y: originY},
		Max: image.Point{X: originX + w, Y: originY + elementHeight},
	}
}

// PlaceBottom returns the bounding rectangle for a full-width element anchored
// at the bottom of the content area. Returns image.Rectangle{} (zero value) when
// the element height exceeds the available content height.
func PlaceBottom(bridge layout.LayoutCalculator, elementHeight int) image.Rectangle {
	w := bridge.AvailableContentWidth()
	h := bridge.AvailableContentHeight()
	if elementHeight > h {
		return image.Rectangle{}
	}

	originX, originY := bridge.ContentOrigin()

	return image.Rectangle{
		Min: image.Point{X: originX, Y: originY + h - elementHeight},
		Max: image.Point{X: originX + w, Y: originY + h},
	}
}

// PlaceBottomRight returns the bounding rectangle for an element anchored at
// the bottom-right corner of the content area. Returns image.Rectangle{} (zero
// value) when the element dimensions exceed the available content area.
func PlaceBottomRight(bridge layout.LayoutCalculator, elementWidth, elementHeight int) image.Rectangle {
	w := bridge.AvailableContentWidth()
	h := bridge.AvailableContentHeight()
	if elementWidth > w || elementHeight > h {
		return image.Rectangle{}
	}

	originX, originY := bridge.ContentOrigin()

	return image.Rectangle{
		Min: image.Point{X: originX + w - elementWidth, Y: originY + h - elementHeight},
		Max: image.Point{X: originX + w, Y: originY + h},
	}
}
