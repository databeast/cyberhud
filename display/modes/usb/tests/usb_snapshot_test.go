package tests_test

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/modes/usb"
	"github.com/databeast/cyberhud/display/modes/usb/source"
)

// snapshotOutputDir is the persistent directory where usb snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("..", "snapshots")

// categoryFromStyleName derives the testsnapshot.DisplayCategory from a usb
// style's name prefix. USB styles follow the capability naming convention:
// mono-slow-*, mono-fast-*, grayscale-slow-*, grayscale-fast-*, color-slow-*,
// color-fast-*.
func categoryFromStyleName(name string) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "grayscale-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	case strings.HasPrefix(name, "eink-"):
		return testsnapshot.CategoryEink
	default:
		return testsnapshot.CategoryMono
	}
}

// cleanDirContents removes every entry inside dir without removing dir itself.
// This avoids the Windows "device or resource busy" failure that os.RemoveAll
// hits when an external process (e.g. an IDE file watcher) holds an open handle
// on the directory, while still leaving the directory empty for a fresh run.
func cleanDirContents(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if err := os.RemoveAll(filepath.Join(dir, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

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
// Each style is rendered at its OWN declared panel size (reqs.MinWidth ×
// reqs.MinHeight) with the display category derived from its name prefix, so
// the full / compact / minimal layout branches and every aspect ratio are
// actually exercised. The WithIconGetter option provides USB device status
// icons, and WithPreRender seeds the monitor state with a deterministic USB
// device snapshot.
func TestUSBPNGSnapshots(t *testing.T) {
	styles := usb.UsbRegistryEnumerate()
	if len(styles) == 0 {
		t.Fatal("usbRegistry contains zero styles")
	}

	// Ensure the output directory exists and is clean. Clear the directory's
	// CONTENTS rather than removing the directory itself: on Windows an IDE file
	// watcher commonly holds an open handle on the dir, which makes RemoveAll on
	// the directory fail with "device or resource busy" even though every file
	// inside can be deleted. Clearing contents preserves the "start clean" intent
	// while staying resilient to that watcher lock.
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}
	if err := cleanDirContents(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}

	// Prepare deterministic USB device state.
	snap := testUSBSnapshot()

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()
			if reqs.MinWidth == 0 || reqs.MinHeight == 0 {
				t.Skip("skipping: style has unconstrained dimensions")
			}

			// Derive display category from the style name prefix.
			category := categoryFromStyleName(s.Name())

			// Build a policy that targets the current style.
			p := usb.DefaultPolicy()
			p.Style = s.Name()

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithDimensions(reqs.MinWidth, reqs.MinHeight),
				testsnapshot.WithMode("usb"),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithDisplayCategory(category),
				testsnapshot.WithIconGetter(stubIconGetter),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithPreRender(func() {
					usb.SetTestMonitorState(snap, p)
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, reqs.MinWidth, reqs.MinHeight)
		})
	}
}
