package attract_hacking

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/databeast/cyberhud/display/modes/attract_hacking/source"
	"github.com/databeast/cyberhud/display/modes/attract_hacking/styles"
	"github.com/databeast/cyberhud/display/region"
	"github.com/databeast/cyberhud/display/style"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

var (
	frameCounter uint64
	lastTick     time.Time
	phase        float64
)

var hackingRegistry = style.NewRegistry[source.Data, source.Policy](
	styles.MonoSlowStyle{Width: 122, Height: 250}, styles.MonoSlowStyle{Width: 176, Height: 264}, styles.MonoSlowStyle{Width: 200, Height: 200}, styles.MonoSlowStyle{Width: 212, Height: 104}, styles.MonoSlowStyle{Width: 296, Height: 128}, styles.MonoSlowStyle{Width: 400, Height: 300}, styles.MonoSlowStyle{Width: 480, Height: 800}, styles.MonoSlowStyle{Width: 800, Height: 480}, styles.MonoSlowStyle{Width: 104, Height: 212}, styles.MonoSlowStyle{Width: 250, Height: 122}, styles.MonoSlowStyle{Width: 128, Height: 296}, styles.MonoSlowStyle{Width: 264, Height: 176}, styles.MonoSlowStyle{Width: 300, Height: 400},
	styles.MonoFastStyle{Width: 16, Height: 8}, styles.MonoFastStyle{Width: 8, Height: 16}, styles.MonoFastStyle{Width: 128, Height: 32}, styles.MonoFastStyle{Width: 128, Height: 64}, styles.MonoFastStyle{Width: 128, Height: 128}, styles.MonoFastStyle{Width: 32, Height: 128}, styles.MonoFastStyle{Width: 64, Height: 128},
	styles.GrayscaleSlowStyle{Width: 122, Height: 250}, styles.GrayscaleSlowStyle{Width: 176, Height: 264}, styles.GrayscaleSlowStyle{Width: 200, Height: 200}, styles.GrayscaleSlowStyle{Width: 212, Height: 104}, styles.GrayscaleSlowStyle{Width: 296, Height: 128}, styles.GrayscaleSlowStyle{Width: 400, Height: 300}, styles.GrayscaleSlowStyle{Width: 480, Height: 800}, styles.GrayscaleSlowStyle{Width: 800, Height: 480}, styles.GrayscaleSlowStyle{Width: 104, Height: 212}, styles.GrayscaleSlowStyle{Width: 250, Height: 122}, styles.GrayscaleSlowStyle{Width: 128, Height: 296}, styles.GrayscaleSlowStyle{Width: 264, Height: 176}, styles.GrayscaleSlowStyle{Width: 300, Height: 400},
	styles.GrayscaleFastStyle{Width: 16, Height: 8}, styles.GrayscaleFastStyle{Width: 8, Height: 16}, styles.GrayscaleFastStyle{Width: 160, Height: 80}, styles.GrayscaleFastStyle{Width: 160, Height: 128}, styles.GrayscaleFastStyle{Width: 240, Height: 135}, styles.GrayscaleFastStyle{Width: 240, Height: 240}, styles.GrayscaleFastStyle{Width: 320, Height: 240}, styles.GrayscaleFastStyle{Width: 480, Height: 320}, styles.GrayscaleFastStyle{Width: 800, Height: 480}, styles.GrayscaleFastStyle{Width: 80, Height: 160}, styles.GrayscaleFastStyle{Width: 128, Height: 160}, styles.GrayscaleFastStyle{Width: 135, Height: 240}, styles.GrayscaleFastStyle{Width: 240, Height: 320}, styles.GrayscaleFastStyle{Width: 320, Height: 480}, styles.GrayscaleFastStyle{Width: 480, Height: 800}, styles.GrayscaleFastStyle{Width: 128, Height: 128},
	styles.ColorSlowStyle{Width: 122, Height: 250}, styles.ColorSlowStyle{Width: 176, Height: 264}, styles.ColorSlowStyle{Width: 200, Height: 200}, styles.ColorSlowStyle{Width: 212, Height: 104}, styles.ColorSlowStyle{Width: 296, Height: 128}, styles.ColorSlowStyle{Width: 400, Height: 300}, styles.ColorSlowStyle{Width: 480, Height: 800}, styles.ColorSlowStyle{Width: 800, Height: 480}, styles.ColorSlowStyle{Width: 104, Height: 212}, styles.ColorSlowStyle{Width: 250, Height: 122}, styles.ColorSlowStyle{Width: 128, Height: 296}, styles.ColorSlowStyle{Width: 264, Height: 176}, styles.ColorSlowStyle{Width: 300, Height: 400},
	styles.ColorFastStyle{Width: 16, Height: 8}, styles.ColorFastStyle{Width: 8, Height: 16}, styles.ColorFastStyle{Width: 160, Height: 80}, styles.ColorFastStyle{Width: 160, Height: 128}, styles.ColorFastStyle{Width: 240, Height: 135}, styles.ColorFastStyle{Width: 240, Height: 240}, styles.ColorFastStyle{Width: 320, Height: 240}, styles.ColorFastStyle{Width: 480, Height: 320}, styles.ColorFastStyle{Width: 800, Height: 480}, styles.ColorFastStyle{Width: 80, Height: 160}, styles.ColorFastStyle{Width: 128, Height: 160}, styles.ColorFastStyle{Width: 135, Height: 240}, styles.ColorFastStyle{Width: 240, Height: 320}, styles.ColorFastStyle{Width: 320, Height: 480}, styles.ColorFastStyle{Width: 480, Height: 800}, styles.ColorFastStyle{Width: 128, Height: 128},
)

func BuildView(hints textlayout.TextHints) style.ViewData {
	if !lastTick.IsZero() {
		elapsed := time.Since(lastTick)
		if elapsed > 250*time.Millisecond {
			elapsed = 250 * time.Millisecond
		}
		phase += elapsed.Seconds() * GetPolicy().Speed
	}
	lastTick = time.Now()
	frameCounter++

	p := GetPolicy()
	img := buildTerminalFrame(hints.PixelWidth, hints.PixelHeight, p)

	snap := source.Data{Sprites: []widgets.Sprite{{Image: img, Position: image.Point{0, 0}, Label: "attract_hacking"}}}
	ctx := style.NewStyleContext(hints)
	s, reason := style.ResolveStyle(hackingRegistry, hints, "attract_hacking", "")
	vd := s.Build(snap, p, ctx)
	vd.Static = false
	vd.StyleReport = style.StyleReport{Name: s.Name(), Reason: reason}
	return vd
}

func RenderCacheKey() uint32 {
	return region.CalcRegionCacheKey(fmt.Sprintf("attract_hacking:%d|%s", frameCounter, GetPolicy().Fingerprint()))
}

func buildTerminalFrame(w, h int, p source.Policy) *image.RGBA {
	if w <= 0 || h <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	bg := color.RGBA{5, 10, 18, 255}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, bg)
		}
	}

	for y := 0; y < h; y += 4 {
		for x := 0; x < w; x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.RGBA{22, 36, 48, 255})
			}
		}
	}

	midX := w / 2
	midY := h / 2
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			dx := float64(x-midX) / float64(max(1, w/2))
			dy := float64(y-midY) / float64(max(1, h/2))
			d := math.Sqrt(dx*dx + dy*dy)
			if d < 1.2 {
				img.Set(x, y, color.RGBA{26, 128, 255, 64})
			}
		}
	}

	glow := color.RGBA{56, 240, 255, 255}
	wire := color.RGBA{18, 255, 150, 255}
	warn := color.RGBA{255, 80, 180, 255}

	// Terminal header bar
	for y := 0; y < 18; y++ {
		for x := 0; x < w; x++ {
			if x < w-1 && y > 0 {
				img.Set(x, y, color.RGBA{16, 22, 32, 255})
			}
		}
	}
	labelX := 10
	labelY := 14
	drawText(img, labelX, labelY, "ROOT//KRAKEN@NET", glow)
	drawText(img, w-80, labelY, "RUNNING", wire)

	for i, line := range hackinglines {
		y := 32 + i*10 + int(math.Sin(phase*0.9+float64(i))*4)
		if y > h-22 {
			continue
		}
		col := glow
		if i%3 == 0 {
			col = wire
		}
		if i%5 == 0 {
			col = warn
		}
		if i%2 == 0 {
			drawText(img, 16, y, line, col)
		} else {
			drawText(img, 24+int(math.Sin(phase+float64(i))*6), y, line, col)
		}
	}

	// Center panel / holographic frame.
	frameW, frameH := w/2, h/2
	x0, y0 := midX-frameW/2, midY-frameH/2
	for x := x0; x < x0+frameW; x++ {
		for y := y0; y < y0+frameH; y++ {
			if x == x0 || y == y0 || x == x0+frameW-1 || y == y0+frameH-1 {
				img.Set(x, y, glow)
			}
		}
	}

	for y := y0 + 12; y < y0+frameH-12; y += 8 {
		for x := x0 + 12; x < x0+frameW-12; x += 9 {
			img.Set(x, y, color.RGBA{72, 210, 255, 150})
		}
	}

	// Progress bar
	progress := int((0.5 + 0.5*math.Sin(phase*1.3)) * float64(w-70))
	barY := h - 22
	for x := 12; x < w-12; x++ {
		if x < 12+progress {
			img.Set(x, barY, glow)
			img.Set(x, barY+1, glow)
		} else if x%7 == 0 {
			img.Set(x, barY, color.RGBA{30, 42, 58, 255})
		}
	}
	for i := 0; i < 5; i++ {
		img.Set(12+i*40, barY-6, color.RGBA{18, 255, 150, 220})
		img.Set(12+i*40, barY+5, color.RGBA{18, 255, 150, 220})
	}
	drawText(img, 14, h-12, "[ ACCESS_VECTOR: LIVE ]", wire)

	// Retro glitch offsets
	glitch := int(p.Glitch * 6)
	if glitch > 0 {
		for y := 0; y < h; y += 9 {
			if int(math.Sin(float64(y)+phase))%2 == 0 {
				for x := 0; x < 3; x++ {
					if x+glitch < w {
						img.Set(x+glitch, y, color.RGBA{255, 255, 255, 170})
					}
				}
			}
		}
	}

	return img
}

func drawText(img *image.RGBA, x, y int, s string, c color.Color) {
	if len(s) == 0 {
		return
	}
	d := &font.Drawer{Dst: img, Src: image.NewUniform(c), Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
