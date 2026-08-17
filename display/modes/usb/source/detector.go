package source

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DeviceInfo holds a compact USB identity for bench display.
type DeviceInfo struct {
	Key          string
	VendorID     string
	ProductID    string
	Manufacturer string
	Product      string
	Serial       string
	BusNum       string
	DevNum       string
	DeviceClass  string // Human-readable label ("Storage", "HID", "Hub", etc.)
}

func (d DeviceInfo) identity() string {
	return strings.ToLower(strings.TrimSpace(d.VendorID)) + ":" + strings.ToLower(strings.TrimSpace(d.ProductID)) + ":" + strings.TrimSpace(d.Serial)
}

func ScanDevices(root string, policy Policy) ([]DeviceInfo, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	type candidate struct {
		info    DeviceInfo
		modTime time.Time
	}
	candidates := make([]candidate, 0, len(entries))
	seen := make(map[string]bool)
	for _, ent := range entries {
		name := strings.TrimSpace(ent.Name())
		if name == "" || strings.Contains(name, ":") {
			continue
		}
		devDir := filepath.Join(root, name)
		// Deduplicate symlinks that resolve to the same directory.
		resolved, resolveErr := filepath.EvalSymlinks(devDir)
		if resolveErr != nil {
			resolved = devDir
		}
		if seen[resolved] {
			continue
		}
		seen[resolved] = true
		vendor := readTrimmed(filepath.Join(devDir, "idVendor"))
		productID := readTrimmed(filepath.Join(devDir, "idProduct"))
		if vendor == "" || productID == "" {
			continue
		}
		product := readTrimmed(filepath.Join(devDir, "product"))
		deviceClass := readTrimmed(filepath.Join(devDir, "bDeviceClass"))
		if policy.HideRootHubs && isRootHub(product, deviceClass) {
			continue
		}
		fi, statErr := os.Stat(devDir)
		modTime := time.Time{}
		if statErr == nil {
			modTime = fi.ModTime()
		}
		candidates = append(candidates, candidate{
			info: DeviceInfo{
				Key:          name,
				VendorID:     strings.ToUpper(vendor),
				ProductID:    strings.ToUpper(productID),
				Manufacturer: readTrimmed(filepath.Join(devDir, "manufacturer")),
				Product:      product,
				Serial:       readTrimmed(filepath.Join(devDir, "serial")),
				BusNum:       readTrimmed(filepath.Join(devDir, "busnum")),
				DevNum:       readTrimmed(filepath.Join(devDir, "devnum")),
				DeviceClass:  classifyDevice(devDir),
			},
			modTime: modTime,
		})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].modTime.Equal(candidates[j].modTime) {
			return candidates[i].info.Key < candidates[j].info.Key
		}
		return candidates[i].modTime.After(candidates[j].modTime)
	})
	out := make([]DeviceInfo, 0, len(candidates))
	for _, c := range candidates {
		out = append(out, c.info)
	}
	return out, nil
}

func Transition(prev map[string]DeviceInfo, curr []DeviceInfo, last Snapshot, now time.Time) (map[string]DeviceInfo, Snapshot, bool) {
	next := make(map[string]DeviceInfo, len(curr))
	for _, dev := range curr {
		next[dev.Key] = dev
	}
	snapshot := last
	changed := false

	for _, dev := range curr {
		prevDev, ok := prev[dev.Key]
		if !ok || prevDev.identity() != dev.identity() {
			snapshot.HasLast = true
			snapshot.Connected = true
			snapshot.Device = dev
			snapshot.LastConnectedAt = now
			changed = true
			break
		}
	}

	if snapshot.HasLast {
		if active, ok := next[snapshot.Device.Key]; ok && active.identity() == snapshot.Device.identity() {
			snapshot.Connected = true
		} else if snapshot.Connected {
			snapshot.Connected = false
			changed = true
		}
	}

	return next, snapshot, changed
}

func readTrimmed(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func isRootHub(product, classCode string) bool {
	if strings.Contains(strings.ToLower(strings.TrimSpace(product)), "root hub") {
		return true
	}
	classCode = strings.ToLower(strings.TrimSpace(classCode))
	classCode = strings.TrimPrefix(classCode, "0x")
	return classCode == "09"
}

// classCodeMap maps lowercase two-char hex class codes to human-readable labels.
var classCodeMap = map[string]string{
	"09": "Hub",
	"08": "Storage",
	"03": "HID",
	"01": "Audio",
	"02": "Network",
	"07": "Printer",
	"0e": "Video",
	"e0": "Wireless",
}

// classifyDevice reads bDeviceClass (and optionally bInterfaceClass) from the
// given sysfs device directory and returns a human-readable device class label.
func classifyDevice(devDir string) string {
	raw := readTrimmed(filepath.Join(devDir, "bDeviceClass"))
	code := normalizeClassCode(raw)
	if code == "" {
		return "Device"
	}

	if code == "00" {
		// Class defined at interface level — read from lowest-numbered interface subdir.
		code = readInterfaceClass(devDir)
		if code == "" {
			return "Device"
		}
	}

	if label, ok := classCodeMap[code]; ok {
		return label
	}
	return "Device"
}

// normalizeClassCode normalizes a raw sysfs class code value to a lowercase
// two-character hex string. It handles formats like "0x09", "09", "9", etc.
func normalizeClassCode(raw string) string {
	if raw == "" {
		return ""
	}
	raw = strings.ToLower(strings.TrimSpace(raw))
	raw = strings.TrimPrefix(raw, "0x")
	if raw == "" {
		return ""
	}
	// Parse as hex to validate and normalize.
	val, err := strconv.ParseUint(raw, 16, 8)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%02x", val)
}

// readInterfaceClass finds the lowest-numbered interface subdirectory under
// devDir and reads its bInterfaceClass attribute.
func readInterfaceClass(devDir string) string {
	entries, err := os.ReadDir(devDir)
	if err != nil {
		return ""
	}

	// Interface subdirectories have names like "1-2:1.0", "1-2:1.1", etc.
	// The interface number is the part after the last dot.
	type ifaceEntry struct {
		path  string
		ifNum int
	}

	var ifaces []ifaceEntry
	for _, ent := range entries {
		name := ent.Name()
		if !strings.Contains(name, ":") {
			continue
		}
		// Extract interface number from the last segment after "."
		dotIdx := strings.LastIndex(name, ".")
		if dotIdx < 0 || dotIdx >= len(name)-1 {
			continue
		}
		numStr := name[dotIdx+1:]
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}
		ifaces = append(ifaces, ifaceEntry{
			path:  filepath.Join(devDir, name),
			ifNum: num,
		})
	}

	if len(ifaces) == 0 {
		return ""
	}

	// Sort by interface number to find the lowest.
	sort.Slice(ifaces, func(i, j int) bool {
		return ifaces[i].ifNum < ifaces[j].ifNum
	})

	raw := readTrimmed(filepath.Join(ifaces[0].path, "bInterfaceClass"))
	return normalizeClassCode(raw)
}
