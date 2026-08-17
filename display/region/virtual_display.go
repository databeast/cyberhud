package region

import (
	"errors"
	"fmt"
	"image"
)

// VirtualDisplay is the VirtualDisplay architectural component: the unified RGBA
// framebuffer backing all Physical Screens on a single Panel. Regions obtain their
// rendering surfaces as zero-copy sub-images of this framebuffer, sharing the
// underlying pixel memory so that writes to a Region's surface are immediately
// visible in the VirtualDisplay without any copy step.
type VirtualDisplay struct {
	fb     *image.RGBA
	bounds image.Rectangle
}

// NewVirtualDisplay constructs a VirtualDisplay from a bounding rectangle.
// The bounds parameter specifies the desired pixel dimensions; it is normalized to
// origin (0,0) internally. Returns the new VirtualDisplay or an error if bounds has
// zero or negative area.
func NewVirtualDisplay(bounds image.Rectangle) (*VirtualDisplay, error) {
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		return nil, fmt.Errorf("virtual display: bounds must have positive width and height, got %dx%d", bounds.Dx(), bounds.Dy())
	}
	// Normalize to origin (0,0).
	normalized := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
	return &VirtualDisplay{
		fb:     image.NewRGBA(normalized),
		bounds: normalized,
	}, nil
}

// NewVirtualDisplayFromScreens constructs a VirtualDisplay from Physical Screen
// positions. The screens parameter provides the set of hardware screen rectangles
// whose union determines the framebuffer size. The resulting bounds are the minimum
// bounding rectangle encompassing all screen positions with origin at (0,0).
// Returns the new VirtualDisplay or an error if screens is empty or any screen has
// invalid dimensions.
func NewVirtualDisplayFromScreens(screens []ScreenPosition) (*VirtualDisplay, error) {
	if len(screens) == 0 {
		return nil, errors.New("virtual display: at least one physical screen is required")
	}

	// Validate each screen has positive dimensions.
	for _, s := range screens {
		w := s.Bounds.Dx()
		h := s.Bounds.Dy()
		if w <= 0 || h <= 0 {
			return nil, fmt.Errorf("virtual display: screen %q has invalid dimensions %dx%d", s.Name, w, h)
		}
	}

	// Compute minimum bounding rectangle encompassing all screens.
	var maxX, maxY int
	for _, s := range screens {
		if s.Bounds.Max.X > maxX {
			maxX = s.Bounds.Max.X
		}
		if s.Bounds.Max.Y > maxY {
			maxY = s.Bounds.Max.Y
		}
	}

	bounds := image.Rect(0, 0, maxX, maxY)
	return &VirtualDisplay{
		fb:     image.NewRGBA(bounds),
		bounds: bounds,
	}, nil
}

// Bounds returns the VirtualDisplay's bounding rectangle as an image.Rectangle
// with origin at (0,0). Used by the RegionManager to validate that allocated
// Region bounds fit within the framebuffer.
func (vd *VirtualDisplay) Bounds() image.Rectangle {
	return vd.bounds
}

// FrameBuffer returns the underlying *image.RGBA pixel buffer. The RegionManager
// uses this to create zero-copy sub-image surfaces for each Region, and the
// FlushPath reads from it to extract per-screen rectangles for hardware output.
func (vd *VirtualDisplay) FrameBuffer() *image.RGBA {
	return vd.fb
}

// SubImage returns the pixel rectangle at rect as a zero-origin *image.RGBA
// backed by the same underlying memory (zero-copy via image.RGBA.SubImage). The
// rect parameter specifies the sub-region in VirtualDisplay coordinates to extract.
// The returned image shares pixel memory with the VirtualDisplay framebuffer, so
// writes to it are immediately reflected in the parent buffer without copying.
func (vd *VirtualDisplay) SubImage(rect image.Rectangle) *image.RGBA {
	// Intersect with bounds to avoid out-of-range access.
	rect = rect.Intersect(vd.bounds)
	sub := vd.fb.SubImage(rect).(*image.RGBA)

	// Create a zero-origin view that shares the same pixel memory.
	// The sub-image from image.RGBA.SubImage already shares memory via Pix/Stride,
	// but its Rect still refers to the parent coordinate system.
	// We adjust the Rect to be zero-origin while keeping the same Pix slice and Stride.
	return &image.RGBA{
		Pix:    sub.Pix,
		Stride: sub.Stride,
		Rect:   image.Rect(0, 0, rect.Dx(), rect.Dy()),
	}
}
