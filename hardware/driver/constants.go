package driver

// ST7789 / ST7735S MADCTL flags exported for panel definitions and driver configs.
const (
	MadctlMY  byte = 0x80
	MadctlMX  byte = 0x40
	MadctlMV  byte = 0x20
	MadctlML  byte = 0x10
	MadctlRGB byte = 0x00
	MadctlBGR byte = 0x08
)
