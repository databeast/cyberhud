package source

import (
	"image"
	"image/color"
	"math"
)

// RenderFrame renders the geometric frame for the given clusters and fragments.
func RenderFrame(clusters []ClusterConfig, fragments []ActiveFragment, time float64, panelW, panelH int, scaleFactor, glowIntensity float64) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, panelW, panelH))

	type rectData struct {
		absX, absY float64
		sq         SquareConfig
		opacity    float64
		cluster    *ClusterConfig
	}

	var allRects []rectData
	var glowCandidates []rectData

	for ci := range clusters {
		cl := &clusters[ci]
		elapsed := time - cl.SpawnTime
		if elapsed <= 0 {
			continue
		}
		fadeInFactor := elapsed / cl.FadeInDuration
		if fadeInFactor > 1.0 {
			fadeInFactor = 1.0
		}
		var glowCandidate *rectData
		for si := range cl.Squares {
			sq := cl.Squares[si]
			absX := (cl.CenterXPct/100.0)*float64(panelW) + sq.OffsetX
			absY := (cl.CenterYPct/100.0)*float64(panelH) + sq.OffsetY
			dist := math.Sqrt(sq.OffsetX*sq.OffsetX + sq.OffsetY*sq.OffsetY)
			fadeOp := ComputeFadeOpacity(time, sq.PhaseOffset, sq.CycleDuration, sq.PeakOpacity)
			glowMult := ComputeGlowMultiplier(dist, cl.BoundingRadius)
			opacity := fadeOp * glowMult * fadeInFactor
			opacity = math.Min(opacity, centralZoneOpacityCap(absX, panelW))
			rd := rectData{absX: absX, absY: absY, sq: sq, opacity: opacity, cluster: cl}
			allRects = append(allRects, rd)
			if opacity > 0.05 {
				if glowCandidate == nil || sq.Size >= glowCandidate.sq.Size {
					rdCopy := rd
					glowCandidate = &rdCopy
				}
			}
		}
		if glowCandidate != nil {
			glowCandidates = append(glowCandidates, *glowCandidate)
		}
	}

	for _, gc := range glowCandidates {
		renderGlow(img, gc.absX, gc.absY, gc.sq, gc.opacity, glowIntensity)
	}
	for _, rd := range allRects {
		if rd.opacity <= 0 {
			continue
		}
		renderSquare(img, rd.absX, rd.absY, rd.sq, rd.opacity)
	}
	for _, f := range fragments {
		fOpacity := ComputeFragmentOpacity(f, time)
		if fOpacity <= 0 {
			continue
		}
		fOpacity = math.Min(fOpacity, centralZoneOpacityCap(f.X, panelW))
		if fOpacity <= 0 {
			continue
		}
		renderFragmentText(img, f, fOpacity)
	}

	return img
}

func alphaBlend(dst, src color.RGBA) color.RGBA {
	if src.A == 0 {
		return dst
	}
	if src.A == 255 {
		return src
	}

	srcA := uint32(src.A)
	invA := 255 - srcA

	outR := (uint32(src.R)*255 + uint32(dst.R)*invA) / 255
	outG := (uint32(src.G)*255 + uint32(dst.G)*invA) / 255
	outB := (uint32(src.B)*255 + uint32(dst.B)*invA) / 255
	outA := srcA + uint32(dst.A)*invA/255

	return color.RGBA{R: uint8(outR), G: uint8(outG), B: uint8(outB), A: uint8(outA)}
}

func hslToRGBA(c HSLColor, alpha float64) color.RGBA {
	h := c.H / 360.0
	s := c.S / 100.0
	l := c.L / 100.0

	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hueToRGB(p, q, h+1.0/3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3.0)
	}

	if alpha < 0 {
		alpha = 0
	}
	if alpha > 1 {
		alpha = 1
	}

	return color.RGBA{
		R: uint8(math.Round(r * 255)),
		G: uint8(math.Round(g * 255)),
		B: uint8(math.Round(b * 255)),
		A: uint8(math.Round(alpha * 255)),
	}
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func renderSquare(img *image.RGBA, absX, absY float64, sq SquareConfig, opacity float64) {
	if opacity <= 0 {
		return
	}

	w := sq.Size * sq.Aspect
	h := sq.Size
	halfW := w / 2
	halfH := h / 2
	rad := sq.Rotation * math.Pi / 180
	cosR := math.Cos(rad)
	sinR := math.Sin(rad)

	corners := [4][2]float64{
		{absX + (-halfW)*cosR - (-halfH)*sinR, absY + (-halfW)*sinR + (-halfH)*cosR},
		{absX + (halfW)*cosR - (-halfH)*sinR, absY + (halfW)*sinR + (-halfH)*cosR},
		{absX + (halfW)*cosR - (halfH)*sinR, absY + (halfW)*sinR + (halfH)*cosR},
		{absX + (-halfW)*cosR - (halfH)*sinR, absY + (-halfW)*sinR + (halfH)*cosR},
	}

	c := hslToRGBA(sq.Color, opacity)
	for i := 0; i < 4; i++ {
		j := (i + 1) % 4
		drawLine(img, corners[i][0], corners[i][1], corners[j][0], corners[j][1], c, 2)
	}
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 float64, c color.RGBA, strokeWidth int) {
	dx := x1 - x0
	dy := y1 - y0
	length := math.Sqrt(dx*dx + dy*dy)
	if length < 0.5 {
		return
	}

	steps := int(math.Ceil(length))
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		px := int(math.Round(x0 + t*dx))
		py := int(math.Round(y0 + t*dy))
		half := strokeWidth / 2
		for ox := -half; ox < strokeWidth-half; ox++ {
			for oy := -half; oy < strokeWidth-half; oy++ {
				setBlendedPixel(img, px+ox, py+oy, c)
			}
		}
	}
}

func setBlendedPixel(img *image.RGBA, x, y int, c color.RGBA) {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return
	}
	off := img.PixOffset(x, y)
	dst := color.RGBA{R: img.Pix[off], G: img.Pix[off+1], B: img.Pix[off+2], A: img.Pix[off+3]}
	blended := alphaBlend(dst, c)
	img.Pix[off] = blended.R
	img.Pix[off+1] = blended.G
	img.Pix[off+2] = blended.B
	img.Pix[off+3] = blended.A
}

func renderGlow(img *image.RGBA, absX, absY float64, sq SquareConfig, opacity, glowIntensity float64) {
	if glowIntensity <= 0 || opacity <= 0 {
		return
	}
	glowColor := hslToRGBA(ComputeGlowColor(sq.Color), opacity*0.4*glowIntensity)
	halfW := sq.Size * sq.Aspect / 2
	halfH := sq.Size / 2
	minX := int(math.Floor(absX - halfW - 4))
	maxX := int(math.Ceil(absX + halfW + 4))
	minY := int(math.Floor(absY - halfH - 4))
	maxY := int(math.Ceil(absY + halfH + 4))
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			dx := math.Abs(float64(x)-absX) / (halfW + 4)
			dy := math.Abs(float64(y)-absY) / (halfH + 4)
			dist := math.Max(dx, dy)
			if dist > 1 {
				continue
			}
			a := uint8(float64(glowColor.A) * (1 - dist))
			setBlendedPixel(img, x, y, color.RGBA{R: glowColor.R, G: glowColor.G, B: glowColor.B, A: a})
		}
	}
}

func renderFragmentText(img *image.RGBA, f ActiveFragment, opacity float64) {
	if opacity <= 0 {
		return
	}
	w := int(math.Max(8, float64(len(f.Text))*f.FontSize*0.55))
	h := int(math.Max(8, f.FontSize*0.9))
	x0 := int(math.Round(f.X))
	y0 := int(math.Round(f.Y))
	col := hslToRGBA(f.Color, opacity)
	for y := y0; y < y0+h; y++ {
		for x := x0; x < x0+w; x++ {
			if x == x0 || x == x0+w-1 || y == y0 || y == y0+h-1 {
				setBlendedPixel(img, x, y, col)
			}
		}
	}
}
