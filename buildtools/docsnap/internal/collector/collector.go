// Package collector implements the core logic for the Image_Collector tool.
// It walks source mode directories, finds PNG snapshots, and copies them
// into the documentation source tree.
package collector

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Summary reports what the collector did.
type Summary struct {
	ModesProcessed int
	ModesSkipped   int
	FilesCopied    int
}

// CollectAll walks srcRoot for mode subdirectories, copies .png files from
// each <mode>/snapshots/ into dstRoot/<mode>/.
// Returns a summary of modes processed and files copied.
func CollectAll(srcRoot, dstRoot string) (Summary, error) {
	var s Summary

	// Read immediate subdirectories of srcRoot (these are mode dirs).
	entries, err := os.ReadDir(srcRoot)
	if err != nil {
		return s, fmt.Errorf("reading source root %q: %w", srcRoot, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		mode := entry.Name()
		snapshotDir := filepath.Join(srcRoot, mode, "snapshots")

		// Check if snapshots/ exists for this mode.
		info, err := os.Stat(snapshotDir)
		if err != nil {
			if os.IsNotExist(err) {
				// Requirement 1.3: skip without error.
				s.ModesSkipped++
				continue
			}
			return s, fmt.Errorf("checking snapshot dir for mode %q: %w", mode, err)
		}
		if !info.IsDir() {
			s.ModesSkipped++
			continue
		}

		// List files and filter to only lowercase .png.
		pngFiles, err := listPNGFiles(snapshotDir)
		if err != nil {
			return s, fmt.Errorf("listing snapshots for mode %q: %w", mode, err)
		}

		// Requirement 1.4: if no .png files, skip without error.
		if len(pngFiles) == 0 {
			s.ModesSkipped++
			continue
		}

		// Mode has .png files — process it.
		dstDir := filepath.Join(dstRoot, mode)

		// Requirement 1.5: clean destination dir before copying (remove stale files).
		if err := os.RemoveAll(dstDir); err != nil {
			return s, fmt.Errorf("cleaning destination dir %q: %w", dstDir, err)
		}

		// Requirement 1.6: create destination directory (including intermediates).
		if err := os.MkdirAll(dstDir, 0o755); err != nil {
			return s, fmt.Errorf("creating destination dir %q: %w", dstDir, err)
		}

		// Copy each .png file byte-for-byte.
		for _, name := range pngFiles {
			srcPath := filepath.Join(snapshotDir, name)
			dstPath := filepath.Join(dstDir, name)
			if err := copyFile(srcPath, dstPath); err != nil {
				return s, fmt.Errorf("copying %q to %q: %w", srcPath, dstPath, err)
			}
			s.FilesCopied++
		}

		s.ModesProcessed++
	}

	return s, nil
}

// listPNGFiles returns filenames in dir that have the exact lowercase ".png" extension
// or ".md" extension (gallery fragments).
func listPNGFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var pngs []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		// Requirement 1.7 / 5.3: only exact lowercase ".png" extension,
		// plus ".md" for gallery fragment files.
		if ext == ".png" || ext == ".md" {
			pngs = append(pngs, e.Name())
		}
	}
	return pngs, nil
}

// copyFile copies src to dst byte-for-byte, preserving content exactly.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}
