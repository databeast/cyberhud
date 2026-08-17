package tests_test

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/modes/usb"
	"github.com/databeast/cyberhud/display/modes/usb/source"
)

// snapshotOutputDir is the persistent directory where usb snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("..", "snapshots")

// defaultTestWidth is the standard test width used for usb styles
// (which have unconstrained dimensions). 128px is a common OLED panel width.
const defaultTestWidth = 128

// defaultTestHeight is the standard test height used for usb styles.
const defaultTestHeight = 128

// testUSBSnapshot returns a deterministic USB Snapshot for snapshot rendering.
func testUSBSnapshot() usb.Snapshot {
	return usb.Snapshot{
		Sequence:        7,
		Connected:       true,
		HasLast:         true,
		LastConnectedAt: time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
		Device: source.DeviceInfo{
			Key:          "usb1-1.2",
			VendorID:     "2341",
			ProductID:    "0043",
			Manufacturer: "Arduino",
			Product:      "Mega 2560",
			Serial:       "557353132383519011C2",
			BusNum:       "001",
			DevNum:       "004",
			DeviceClass:  "CDC",
		},
	}
}

// stubIconGetter provides a minimal icon getter for USB device status icons.
// It returns a small 8x8 colored image for any requested icon name.
func stubIconGetter(name string) (image.Image, bool) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 200, B: 0, A: 255})
		}
	}
	return img, true
}

// TestUSBPNGSnapshots enumerates all registered usb styles and generates
// a snapshot PNG for each through the full production pipeline via the
// snapshottest framework.
//
// USB styles have unconstrained dimensions (MinWidth=0, MinHeight=0), so we
// use fixed standard test dimensions. The WithIconGetter option provides USB
// device status icons, and WithPreRender seeds the monitor state with a
// deterministic USB device snapshot.
//

func TestUSBPNGSnapshots(t *testing.T) {
	styles := usb.UsbRegistryEnumerate()
	if len(styles) == 0 {
		t.Fatal("usbRegistry contains zero styles")
	}

	// Ensure the output directory exists and is clean.
	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	// Prepare deterministic USB device state.
	snap := testUSBSnapshot()

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			// Build a policy that targets the current style.
			p := usb.DefaultPolicy()
			p.Style = s.Name()

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithDimensions(defaultTestWidth, defaultTestHeight),
				testsnapshot.WithMode("usb"),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithDisplayCategory(testsnapshot.CategoryMono),
				testsnapshot.WithIconGetter(stubIconGetter),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithPreRender(func() {
					usb.SetTestMonitorState(snap, p)
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, defaultTestWidth, defaultTestHeight)
		})
	}
}
