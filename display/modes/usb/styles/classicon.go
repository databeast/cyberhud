package styles

import font "github.com/databeast/cyberhud/display/surface/fonts"

// ClassIcon returns the Material Icons rune for a given device class string.
// Unrecognized or empty class strings return font.IconUsb as fallback.
func ClassIcon(class string) rune {
	switch class {
	case "Storage":
		return font.IconStorage
	case "HID":
		return font.IconKeyboard
	case "Audio":
		return font.IconHeadset
	case "Network":
		return font.IconRouter
	case "Printer":
		return font.IconPrint
	case "Video":
		return font.IconVideocam
	case "Wireless":
		return font.IconWifi
	case "Hub":
		return font.IconHub
	case "Device":
		return font.IconUsb
	default:
		return font.IconUsb
	}
}

// KnownDeviceClasses returns the set of recognized device class strings.
// Useful for test generators.
func KnownDeviceClasses() []string {
	return []string{"Storage", "HID", "Audio", "Network", "Printer", "Video", "Wireless", "Hub", "Device"}
}
