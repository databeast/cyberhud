package region

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"log"
	"sort"
)

// FlushPath is the FlushPath architectural component: it extracts Physical
// Screen rectangles from the [VirtualDisplay] framebuffer and pushes each
// screen's pixel data to hardware via the screen's DrawTarget. One FlushPath
// exists per Panel and is called once per frame by the [RenderLoop].
type FlushPath struct {
	vd      *VirtualDisplay
	screens []ScreenPosition
	buffers []*image.RGBA // pre-allocated per-screen buffers (logical dimensions), reused each flush
	hwBufs  []*image.RGBA // pre-allocated per-screen buffers (hardware dimensions) for rotated screens; nil if no rotation
}

// NewFlushPath creates a [FlushPath] for the given VirtualDisplay and screens.
// The vd parameter is the VirtualDisplay whose framebuffer will be read during
// Flush. The screens parameter provides the set of Physical Screen positions
// that define extraction rectangles and DrawTargets. Screens are sorted by
// ascending Index internally so [FlushPath.Flush] always processes them in a
// deterministic order. Per-screen RGBA buffers are pre-allocated to avoid
// per-frame allocations during Flush.
//
// Returns the initialized FlushPath ready for use by the RenderLoop.
func NewFlushPath(vd *VirtualDisplay, screens []ScreenPosition) *FlushPath {
	// Make a sorted copy so we don't mutate the caller's slice.
	sorted := make([]ScreenPosition, len(screens))
	copy(sorted, screens)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Index < sorted[j].Index
	})

	// Pre-allocate a destination buffer for each screen (logical dimensions).
	buffers := make([]*image.RGBA, len(sorted))
	hwBufs := make([]*image.RGBA, len(sorted))
	needsHWBufs := false
	for i, screen := range sorted {
		screenBounds := screen.Bounds.Intersect(vd.bounds)
		if screenBounds.Empty() {
			continue
		}
		buffers[i] = image.NewRGBA(image.Rect(0, 0, screenBounds.Dx(), screenBounds.Dy()))

		// If screen is rotated, pre-allocate a hardware-dimension buffer for the
		// rotated output.
		if screen.Rotation == 90 || screen.Rotation == 270 {
			// Swap width/height for the hardware buffer.
			hwBufs[i] = image.NewRGBA(image.Rect(0, 0, screenBounds.Dy(), screenBounds.Dx()))
			needsHWBufs = true
		} else if screen.Rotation == 180 {
			// Same dimensions, different pixel order.
			hwBufs[i] = image.NewRGBA(image.Rect(0, 0, screenBounds.Dx(), screenBounds.Dy()))
			needsHWBufs = true
		}
	}

	if !needsHWBufs {
		hwBufs = nil
	}

	return &FlushPath{
		vd:      vd,
		screens: sorted,
		buffers: buffers,
		hwBufs:  hwBufs,
	}
}

// Flush extracts each screen's sub-rectangle from the VirtualDisplay framebuffer,
// fills uncovered pixels with opaque black, translates to zero-origin, and calls
// DrawImage on the screen's DrawTarget. Processing continues on error so that a
// single failing screen does not block the others. All errors encountered are
// combined via [errors.Join] and returned as a single error (nil when all screens
// succeed).
//
// Pre-allocated buffers are reused each frame to avoid per-flush allocations.
func (fp *FlushPath) Flush() error {
	var errs []error

	// Iterate each Physical Screen in index order. Each screen corresponds to a
	// distinct hardware display that needs its own sub-rectangle extracted from the
	// VirtualDisplay. Processing per-screen (rather than per-region) is necessary
	// because hardware targets accept exactly one complete image per flush — if we
	// skipped screens or merged them, hardware displays would show stale or
	// incomplete frames.
	for i, screen := range fp.screens {
		// Skip screens with nil Target silently.
		if screen.Target == nil {
			continue
		}

		screenBounds := screen.Bounds.Intersect(fp.vd.bounds)
		if screenBounds.Empty() {
			continue
		}

		dst := fp.buffers[i]
		if dst == nil {
			continue
		}

		// Fill the reusable buffer with opaque black before compositing. This
		// guarantees pixels not covered by any Region render as opaque black
		// rather than transparent or garbage from a prior frame. Without this
		// step, uncovered areas would show stale pixel data from the previous
		// flush cycle because buffers are reused across frames.
		fillOpaqueBlack(dst)

		// Composite the VirtualDisplay sub-rect onto the black background using
		// draw.Over. Regions write with full alpha so their pixels fully overwrite
		// the black fill. draw.Over (not draw.Src) is required because transparent
		// pixels in the VD must not replace the black background — if draw.Src were
		// used, uncovered pixels would appear transparent on hardware that interprets
		// alpha, causing visual artifacts.
		draw.Draw(dst, dst.Bounds(), fp.vd.fb, screenBounds.Min, draw.Over)

		// Push the composited image to the hardware display. This must be the final
		// step per screen — doing it before compositing is complete would send a
		// partially-rendered frame to the display.
		//
		// If the screen has a rotation, apply pixel rotation from the logical buffer
		// into the hardware-dimensions buffer before sending.
		var output *image.RGBA
		if screen.Rotation != 0 && fp.hwBufs != nil && fp.hwBufs[i] != nil {
			rotateBuf(dst, fp.hwBufs[i], screen.Rotation)
			output = fp.hwBufs[i]
		} else {
			output = dst
		}

		// If the screen has a horizontal mirror, flip pixels left-to-right in-place.
		if screen.MirrorX {
			mirrorXBuf(output)
		}

		if err := screen.Target.DrawImage(output); err != nil {
			log.Printf("flush: screen %q (index %d): DrawImage error: %v", screen.Name, screen.Index, err)
			errs = append(errs, fmt.Errorf("screen %q (index %d): %w", screen.Name, screen.Index, err))
		}
	}

	return errors.Join(errs...)
}

// fillOpaqueBlack fills an RGBA buffer with opaque black (0,0,0,255) without
// allocating an image.Uniform. This is a hot path called every flush per screen.
func fillOpaqueBlack(img *image.RGBA) {
	pix := img.Pix
	for i := 0; i < len(pix); i += 4 {
		pix[i] = 0
		pix[i+1] = 0
		pix[i+2] = 0
		pix[i+3] = 255
	}
}

// rotateBuf rotates the pixels from src into dst by the given degrees (90, 180, 270 CW).
// src and dst must have compatible dimensions:
//   - 90/270: dst is (srcH × srcW)
//   - 180: dst is (srcW × srcH)
//
// This is a hot path called every flush per rotated screen.
func rotateBuf(src, dst *image.RGBA, degrees int) {
	srcW := src.Bounds().Dx()
	srcH := src.Bounds().Dy()
	srcPix := src.Pix
	srcStride := src.Stride
	dstPix := dst.Pix
	dstStride := dst.Stride

	switch degrees {
	case 90:
		// 90° CW: dst(x,y) = src(y, srcW-1-x)
		// dst dimensions: srcH × srcW
		for sy := 0; sy < srcH; sy++ {
			for sx := 0; sx < srcW; sx++ {
				dx := sy
				dy := srcW - 1 - sx
				srcOff := sy*srcStride + sx*4
				dstOff := dy*dstStride + dx*4
				dstPix[dstOff] = srcPix[srcOff]
				dstPix[dstOff+1] = srcPix[srcOff+1]
				dstPix[dstOff+2] = srcPix[srcOff+2]
				dstPix[dstOff+3] = srcPix[srcOff+3]
			}
		}
	case 180:
		// 180°: dst(x,y) = src(srcW-1-x, srcH-1-y)
		for sy := 0; sy < srcH; sy++ {
			for sx := 0; sx < srcW; sx++ {
				dx := srcW - 1 - sx
				dy := srcH - 1 - sy
				srcOff := sy*srcStride + sx*4
				dstOff := dy*dstStride + dx*4
				dstPix[dstOff] = srcPix[srcOff]
				dstPix[dstOff+1] = srcPix[srcOff+1]
				dstPix[dstOff+2] = srcPix[srcOff+2]
				dstPix[dstOff+3] = srcPix[srcOff+3]
			}
		}
	case 270:
		// 270° CW (= 90° CCW): dst(x,y) = src(srcH-1-y, x)
		// dst dimensions: srcH × srcW
		for sy := 0; sy < srcH; sy++ {
			for sx := 0; sx < srcW; sx++ {
				dx := srcH - 1 - sy
				dy := sx
				srcOff := sy*srcStride + sx*4
				dstOff := dy*dstStride + dx*4
				dstPix[dstOff] = srcPix[srcOff]
				dstPix[dstOff+1] = srcPix[srcOff+1]
				dstPix[dstOff+2] = srcPix[srcOff+2]
				dstPix[dstOff+3] = srcPix[srcOff+3]
			}
		}
	}
}

// mirrorXBuf flips an RGBA image horizontally (left-right) in-place.
// Each row's pixels are reversed without allocating a new buffer.
func mirrorXBuf(img *image.RGBA) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	stride := img.Stride
	pix := img.Pix

	for y := 0; y < h; y++ {
		rowStart := y * stride
		for x := 0; x < w/2; x++ {
			l := rowStart + x*4
			r := rowStart + (w-1-x)*4
			// Swap 4 bytes (RGBA)
			pix[l], pix[r] = pix[r], pix[l]
			pix[l+1], pix[r+1] = pix[r+1], pix[l+1]
			pix[l+2], pix[r+2] = pix[r+2], pix[l+2]
			pix[l+3], pix[r+3] = pix[r+3], pix[l+3]
		}
	}
}
