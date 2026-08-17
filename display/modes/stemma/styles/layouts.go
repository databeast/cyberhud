package styles

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/databeast/cyberhud/display/modes/stemma/source"
	sharedcolor "github.com/databeast/cyberhud/display/style/color"
	"github.com/databeast/cyberhud/display/surface/textlayout"
	"github.com/databeast/cyberhud/display/widgets"
)

// iconWidthInGlyphs returns the number of glyph cells consumed by an icon,
// computed as ceil(iconPixelWidth / glyphAdvance). Returns 0 if no icon is available
// or glyphAdvance is zero.
func iconWidthInGlyphs(glyphAdvance int) int {
	if glyphAdvance <= 0 {
		return 0
	}
	// Standard icon width is 8 pixels.
	const iconPixelWidth = 8
	return int(math.Ceil(float64(iconPixelWidth) / float64(glyphAdvance)))
}

func summaryLine(present, total int) string {
	return fmt.Sprintf("%d/%d present", present, total)
}

func deviceBusLabel(d *source.Device) string {
	bus := strings.TrimSpace(d.Bus)
	if bus == "" {
		return "(no bus)"
	}
	return bus
}

func defaultDeviceRow(d *source.Device, nameWidth int) string {
	return fmt.Sprintf("%s %s 0x%02X %s", devicePresenceLabel(d), deviceBusLabel(d), d.Addr, textlayout.Truncate(d.Name, nameWidth))
}

func formatDeviceRow(d *source.Device, nameWidth int, formatter func(*source.Device, int) string) string {
	if formatter != nil {
		return formatter(d, nameWidth)
	}
	return defaultDeviceRow(d, nameWidth)
}

func devicePresenceLabel(d *source.Device) string {
	if d.Present {
		return "✓"
	}
	return "✗"
}

// BuildItems returns STEMMA mode row text.
func BuildItems(devs []*source.Device) []string {
	if len(devs) == 0 {
		return []string{"(no devices found)"}
	}
	items := make([]string, len(devs))
	for i, d := range devs {
		items[i] = defaultDeviceRow(d, 18)
	}
	return items
}

// BuildItemsTruncated returns STEMMA mode rows truncated to the given max character width,
// accounting for the icon width in glyph cells at column 0.
func BuildItemsTruncated(devs []*source.Device, maxChars, iconGlyphCells int) []string {
	if len(devs) == 0 {
		return []string{"(no devices found)"}
	}
	nameWidth := maxChars - iconGlyphCells
	if nameWidth < 1 {
		nameWidth = 1
	}
	items := make([]string, len(devs))
	for i, d := range devs {
		items[i] = textlayout.Truncate(defaultDeviceRow(d, 18), nameWidth)
	}
	return items
}

// BuildColors returns row colours matching STEMMA device presence.
// Returns nil when no devices are provided.
func BuildColors(devs []*source.Device, present, absent color.Color) []color.Color {
	if len(devs) == 0 {
		return nil
	}
	states := make([]bool, len(devs))
	for i, d := range devs {
		states[i] = d.Present
	}
	return sharedcolor.BuildSlice(states, sharedcolor.NewBinaryPalette(toRGBA(present), toRGBA(absent)), len(devs) > 0)
}

// BuildSprites returns a Sprite per device using check (present) or error (absent) icons.
// Each sprite is positioned at column 0 of the device's row.
// The getIcon parameter resolves icon names to images.
func BuildSprites(devs []*source.Device, rowHeight int, getIcon func(name string) (image.Image, bool)) []widgets.Sprite {
	if len(devs) == 0 {
		return nil
	}
	var checkImg, errorImg image.Image
	if getIcon != nil {
		checkImg, _ = getIcon("check")
		errorImg, _ = getIcon("error")
	}

	var sprites []widgets.Sprite
	for i, d := range devs {
		var img image.Image
		var iconName string
		if d.Present {
			img = checkImg
			iconName = "check"
		} else {
			img = errorImg
			iconName = "error"
		}
		if img == nil {
			continue
		}
		sprites = append(sprites, widgets.Sprite{
			Image:    img,
			Position: image.Pt(0, i*rowHeight),
			Label:    fmt.Sprintf("stemma-0x%02X-status", d.Addr),
		})
		_ = iconName // icon name preserved in Label for debug purposes
	}
	return sprites
}

// deviceIconRenderable wraps a single device's icon rendering as a widgets.Renderable.
// It resolves the appropriate icon (check or error) and positions it using
// ContentOrigin as the base coordinate.
type deviceIconRenderable struct {
	dev      *source.Device
	rowIndex int
	rowH     int
	getIcon  func(name string) (image.Image, bool)
	originX  int
	originY  int
}

// RenderFrame produces a Sprite for this device's status icon, or nil if the
// icon cannot be resolved.
func (d *deviceIconRenderable) RenderFrame() *widgets.Sprite {
	if d.getIcon == nil {
		return nil
	}

	iconName := "error"
	if d.dev.Present {
		iconName = "check"
	}

	img, ok := d.getIcon(iconName)
	if !ok || img == nil {
		return nil
	}

	return &widgets.Sprite{
		Image:    img,
		Position: image.Pt(d.originX, d.originY+d.rowIndex*d.rowH),
		Label:    fmt.Sprintf("stemma-0x%02X-status", d.dev.Addr),
	}
}

// StemmaPalette is the standard present/absent color pair for STEMMA devices.
var StemmaPalette = sharedcolor.BinaryPalette{
	Active:   color.RGBA{0x00, 0xCC, 0x44, 0xFF}, // green
	Inactive: color.RGBA{0xCC, 0x00, 0x00, 0xFF}, // red
}

var (
	ColorPresent color.Color = StemmaPalette.Active
	ColorAbsent  color.Color = StemmaPalette.Inactive
)

func toRGBA(c color.Color) color.RGBA {
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}
