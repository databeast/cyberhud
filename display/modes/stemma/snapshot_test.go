package stemma_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/databeast/cyberhud/display/modes/stemma"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/modes/stemma/tests"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
)

// snapshotOutputDir is the persistent directory where stemma snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// defaultTestWidth is the standard test width used for stemma styles
// (which have unconstrained dimensions). 128px is a common OLED panel width.
const defaultTestWidth = 128

// defaultTestHeight is the standard test height used for stemma styles.
const defaultTestHeight = 128

// TestStemmaPNGSnapshots enumerates all registered stemma styles and generates
// a snapshot PNG for each through the full production pipeline via the
// snapshottest framework.
//
// The snapshots feed a fake scanner into the production pipeline so each style
// renders realistic bus/address/state content instead of the empty placeholder.
// A fixed 128x128 test frame keeps the output deterministic while still
// exercising the profile matrix.
//

func TestStemmaPNGSnapshots(t *testing.T) {
	styles := tests.StemmaRegistryEnumerate()
	if len(styles) == 0 {
		t.Fatal("stemmaRegistry contains zero styles")
	}
	if len(styles) < 108 {
		t.Fatalf("expected 108+ stemma styles, got %d", len(styles))
	}

	// Ensure the output directory exists and is clean.
	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			// Build a policy that targets the current style.
			p := stemma.DefaultPolicy()
			p.Style = s.Name()

			scanner := source.NewTestScanner(
				&source.Device{Bus: "/dev/i2c-1", Addr: 0x3C, Name: "SSD1306 OLED display (128x64)", Present: true},
				&source.Device{Bus: "/dev/i2c-1", Addr: 0x76, Name: "BME280/BMP280 env sensor", Present: true},
				&source.Device{Bus: "/dev/i2c-3", Addr: 0x44, Name: "SHT30/SHT31 humidity sensor", Present: false},
				&source.Device{Bus: "/dev/i2c-3", Addr: 0x29, Name: "TSL2591 light sensor", Present: true},
				&source.Device{Bus: "/dev/i2c-6", Addr: 0x68, Name: "MPU6050 motion sensor", Present: true},
			)

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithDimensions(defaultTestWidth, defaultTestHeight),
				testsnapshot.WithMode("stemma"),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithDisplayCategory(testsnapshot.CategoryMono),
				testsnapshot.WithScanner(scanner),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithPreRender(func() {
					stemma.SetPolicy(p)
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, defaultTestWidth, defaultTestHeight)
		})
	}
}
