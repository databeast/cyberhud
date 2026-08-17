package usb

import "fmt"

// BuildItems returns rows for the USB bench screen.
func BuildItems() []string {
	snap := SnapshotNow()
	if !snap.HasLast {
		return []string{
			"Waiting for USB device...",
			"Plug device into bench",
		}
	}
	lines := []string{displayName(snap.Device)}
	lines = append(lines, fmt.Sprintf("VID:PID %s:%s", snap.Device.VendorID, snap.Device.ProductID))
	if snap.Device.BusNum != "" || snap.Device.DevNum != "" {
		lines = append(lines, fmt.Sprintf("Bus %s Dev %s", safeValue(snap.Device.BusNum, "?"), safeValue(snap.Device.DevNum, "?")))
	}
	if snap.Device.Serial != "" {
		lines = append(lines, "SN "+snap.Device.Serial)
	}
	if snap.Connected {
		lines = append(lines, "Status connected")
	} else {
		lines = append(lines, "Status unplugged")
	}
	return lines
}
