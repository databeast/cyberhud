package doccheck

import "strings"

// panelGroups maps panel registry names that share a single documentation page
// to their common base name. Panels sharing the same hardware (e.g., SPI/I2C
// variants of the same board) are documented on one combined page.
var panelGroups = map[string]string{
	"waveshare-2.23-oled-hat-spi": "waveshare-2-23-oled-hat",
	"waveshare-2.23-oled-hat-i2c": "waveshare-2-23-oled-hat",
	"waveshare-1.44":              "waveshare-1-44-lcd-hat",
}

// utilityModeExclusions are mode IDs exempt from the user-facing mode page
// requirement. They are documented under the Utility Modes section instead.
var utilityModeExclusions = map[string]bool{
	"testfonts":    true,
	"testicons":    true,
	"testwidgets":  true,
	"snapshottest": true,
	"testpattern":  true,
}

// stubDirectoryExclusions are display/modes/ subdirectories that contain no
// catalog.Register() call and should not trigger missing-documentation errors.
var stubDirectoryExclusions = map[string]bool{
	"schedule":     true,
	"storage":      true,
	"servicemon":   true,
	"frameclock":   true,
	"testsnapshot": true,
	"tests":        true,
}

// panelGroupBase returns the documentation base name for a panel. For grouped
// panels (SPI/I2C variants sharing one page), it returns the pre-slugified
// group base from panelGroups. For all other panels, it returns the raw panel
// name unchanged.
func panelGroupBase(panelName string) string {
	if base, ok := panelGroups[panelName]; ok {
		return base
	}
	return panelName
}

// panelDocFilename converts a panel registry name to its expected documentation
// markdown filename. Dots become hyphens, spaces become hyphens.
// e.g., "waveshare-2.23-oled-hat-spi" → "waveshare-2-23-oled-hat.md"
//
//	(grouped with -i2c variant on same page)
//
// e.g., "adafruit-2.13-ssd1680" → "adafruit-2-13-ssd1680.md"
func panelDocFilename(panelName string) string {
	base := panelGroupBase(panelName)
	slug := strings.ReplaceAll(base, ".", "-")
	return slug + ".md"
}
