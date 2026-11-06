package storage_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/modes/storage"
	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/style"
)

// snapshotOutputDir is the persistent directory where storage snapshot PNGs are written.
var snapshotOutputDir = filepath.Join("snapshots")

// testStorageSnapshot returns a deterministic StorageSnapshot for rendering,
// spanning the full green→amber→red usage range so colored bars are exercised.
func testStorageSnapshot() storage.StorageSnapshot {
	gib := uint64(1024 * 1024 * 1024)
	return storage.StorageSnapshot{
		Mounts: []storage.MountUsage{
			{Mountpoint: "/", Device: "/dev/mmcblk0p2", Fstype: "ext4",
				TotalBytes: 58 * gib, UsedBytes: 54 * gib, FreeBytes: 4 * gib},
			{Mountpoint: "/boot", Device: "/dev/mmcblk0p1", Fstype: "vfat",
				TotalBytes: 256 * gib / 1024, UsedBytes: 160 * gib / 1024, FreeBytes: 96 * gib / 1024},
			{Mountpoint: "/mnt/data", Device: "/dev/sda1", Fstype: "ext4",
				TotalBytes: 931 * gib, UsedBytes: 210 * gib, FreeBytes: 721 * gib},
		},
		Timestamp: time.Date(2024, 1, 15, 14, 30, 5, 0, time.UTC),
	}
}

func categoryFromStyle(name string, reqs style.SurfaceRequirements) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "grayscale-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	}
	switch reqs.Capability {
	case style.GrayscaleSlow, style.GrayscaleFast:
		return testsnapshot.CategoryGrayscale
	case style.ColorSlow, style.ColorFast:
		return testsnapshot.CategoryColor
	case style.MonoFast:
		return testsnapshot.CategoryMono
	default:
		return testsnapshot.CategoryColor
	}
}

func dimensionsFromRequirements(reqs style.SurfaceRequirements) (int, int) {
	width, height := reqs.MinWidth, reqs.MinHeight
	if width == 0 {
		width = reqs.PreferredWidth
	}
	if height == 0 {
		height = reqs.PreferredHeight
	}
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 240
	}
	return width, height
}

// TestStoragePNGSnapshots enumerates all registered storage styles and renders
// one deterministic PNG for each through the production snapshot pipeline.
func TestStoragePNGSnapshots(t *testing.T) {
	styles := storage.StorageRegistryEnumerate()
	if len(styles) == 0 {
		t.Fatal("storage registry contains zero styles")
	}

	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	snap := testStorageSnapshot()

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()
			width, height := dimensionsFromRequirements(reqs)

			p := storage.DefaultPolicy()
			p.Style = s.Name()

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("storage"),
				testsnapshot.WithDimensions(width, height),
				testsnapshot.WithDisplayCategory(categoryFromStyle(s.Name(), reqs)),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(func() {
					storage.SetPolicy(storage.DefaultPolicy())
					storage.ResetSnapshotState()
				}),
				testsnapshot.WithPreRender(func() {
					storage.SetPolicy(p)
					storage.UpdateSnapshot(snap)
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, width, height)
		})
	}
}
