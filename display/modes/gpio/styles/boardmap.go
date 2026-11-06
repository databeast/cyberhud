package styles

import (
	"image"
	"image/color"

	"github.com/databeast/cyberhud/display/modes/gpio/source"
	"github.com/databeast/cyberhud/display/style"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/widgets"
	"github.com/databeast/cyberhud/display/widgets/led"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
)

// boardMapMinPadPitch is the smallest per-column pitch (px) at which a 20-column
// pin strip still reads as distinct pads rather than a smear. On a ~230 PPI 4"
// panel this is roughly a 3mm cell — the floor below which the map stops being
// legible and the caller should fall back to the text list instead.
const boardMapMinPadPitch = 18

// boardMapFits reports whether a panel is roomy enough to draw the 40-pin strip
// legibly: 20 columns across and 2 rows down, each cell at least the minimum
// pitch. Narrow/short panels fail this and get the compact list.
func boardMapFits(availW, availH int) bool {
	if availW <= 0 || availH <= 0 {
		return false
	}
	colPitch := availW / 20
	rowPitch := availH / 2
	return colPitch >= boardMapMinPadPitch && rowPitch >= boardMapMinPadPitch
}

// busHue assigns each bus family a distinct accent so adjacent bus clusters read
// as colored regions of the board. Standalone driven pins use the generic GPIO
// active color; everything else is resolved by the caller.
func busHue(bus string) color.RGBA {
	switch bus {
	case "SPI0", "SPI1":
		return sharedcolor.Lookup("cyan")
	case "I2C1", "I2C0":
		return sharedcolor.Lookup("amber")
	case "UART0", "UART1":
		return sharedcolor.Lookup("emerald")
	default: // PWM, PCM, GPCLK, and any other recognized bus
		return sharedcolor.Lookup("white")
	}
}

// buildBoardMap renders the physical 40-pin connector as a label-free live map:
// two rows of 20 pads at true board positions, lit by what each pin is doing.
// Position carries identity (no text on pads survives a 4" panel), so the
// answer to "where is SPI0 up?" is simply the cyan cluster glowing at pins
// 19/21/23/24/26. This is the roomy-color view; smaller panels fall back to the
// live-functions list.
//
// Pad vocabulary (all visual, no glyphs):
//   - bus pin        → filled, bus hue, glow (live by definition)
//   - output HIGH    → filled, accent, glow
//   - output LOW     → filled, dimmed accent (present but quiet)
//   - input          → hollow (border only) so "sensing" differs from "driving"
//   - unconfigured   → tiny dim ghost dot (the socket is there, doing nothing)
//   - power          → muted square (orientation landmark)
//   - ground         → muted diamond (orientation landmark)
func buildBoardMap(snap GpioSnapshot, ctx style.StyleContext) style.ViewData {
	pol := snap.Policy
	bridge := ctx.Layout(0)
	availW := bridge.AvailableContentWidth()
	availH := bridge.AvailableContentHeight()
	if availW == 0 || availH == 0 {
		return style.ViewData{Static: true}
	}

	originX, originY := bridge.ContentOrigin()
	pads, _ := source.BoardPads(snap.Pins)

	// 20 columns (physical pins advance 2 per column: 1/2, 3/4, …), 2 rows.
	// Odd physical pins (1,3,5,…) are the top row; even (2,4,6,…) the bottom.
	const cols = 20
	colPitch := availW / cols
	rowPitch := availH / 2

	// Pad diameter: fill most of the smaller cell dimension, leaving a gap so
	// adjacent pads read as separate. Cap so pads don't become absurd blobs on
	// a very large panel.
	cell := colPitch
	if rowPitch < cell {
		cell = rowPitch
	}
	diameter := cell * 3 / 4
	if diameter < 4 {
		diameter = 4
	}
	if diameter > 40 {
		diameter = 40
	}
	ghostDiameter := diameter / 2
	if ghostDiameter < 3 {
		ghostDiameter = 3
	}

	// Center the two-row block vertically within the available height.
	blockH := 2 * rowPitch
	offsetY := (availH - blockH) / 2
	if offsetY < 0 {
		offsetY = 0
	}

	accent := resolveFGColor(pol.FGColor)
	mutedCol := sharedcolor.GPIOPalette.Inactive // grey for power/ground landmarks
	ghostCol := sharedcolor.Dim(mutedCol)

	suppCtx := widgets.SuppressionContext{AvailableWidth: availW, AvailableHeight: availH}
	comp := widgets.NewCompositor(suppCtx)

	for _, pad := range pads {
		// Physical pin p sits at column (p-1)/2; odd pins top row, even bottom.
		col := (pad.Phys - 1) / 2
		row := 0
		if pad.Phys%2 == 0 {
			row = 1
		}
		cx := originX + col*colPitch + (colPitch-diameter)/2
		cy := originY + offsetY + row*rowPitch + (rowPitch-diameter)/2

		cfg := led.Config{
			Brightness: -1.0,
			Diameter:   diameter,
			Bounds:     image.Rectangle{Min: image.Pt(cx, cy)},
		}

		switch pad.Class {
		case source.ClassPower:
			cfg.Shape = led.Square
			cfg.State = led.Off
			cfg.Foreground = mutedCol
			cfg.BorderWidth = 1
			cfg.BorderColor = mutedCol

		case source.ClassGround:
			cfg.Shape = led.Diamond
			cfg.State = led.Off
			cfg.Foreground = mutedCol
			cfg.BorderWidth = 1
			cfg.BorderColor = mutedCol

		default: // ClassGPIO
			if !pad.Configured && pad.Bus == "" {
				// Unconfigured: ghost dot — present but doing nothing.
				cfg.Diameter = ghostDiameter
				cfg.Bounds = image.Rectangle{Min: image.Pt(
					originX+col*colPitch+(colPitch-ghostDiameter)/2,
					originY+offsetY+row*rowPitch+(rowPitch-ghostDiameter)/2,
				)}
				cfg.State = led.Off
				cfg.Foreground = ghostCol
				break
			}

			// Determine the pad's hue.
			hue := accent
			if pad.Bus != "" {
				hue = busHue(pad.Bus)
			}
			if !pol.Color {
				hue = accent // mono: renderer drops color anyway; keep bright
			}

			switch {
			case pad.Bus != "":
				// Live bus member: filled, bus hue, glowing cluster.
				cfg.Shape = led.Circle
				cfg.State = led.On
				cfg.Foreground = hue
				cfg.GlowEnabled = pol.Color
			case pad.Mode == gpiomgr.ModeInput:
				// Input: hollow ring so "sensing" reads differently from "driving".
				cfg.Shape = led.Circle
				cfg.State = led.Off
				cfg.Foreground = hue
				cfg.BorderWidth = 2
				cfg.BorderColor = hue
			case pad.Mode == gpiomgr.ModeOutput && pad.Level:
				// Driven HIGH: filled, bright, glowing — your program acting.
				cfg.Shape = led.Circle
				cfg.State = led.On
				cfg.Foreground = hue
				cfg.GlowEnabled = pol.Color
			default:
				// Driven LOW (or configured-but-low): present but quiet.
				cfg.Shape = led.Circle
				cfg.State = led.Off
				cfg.Foreground = sharedcolor.Dim(hue)
			}
		}

		comp.Add(led.New(cfg))
	}

	return style.ViewData{
		Sprites: comp.Sprites(),
		Static:  true,
	}
}
