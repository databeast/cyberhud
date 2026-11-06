package source

import "time"

// MountUsage captures the usage of a single mounted filesystem at a point in time.
type MountUsage struct {
	// Mountpoint is the path the filesystem is mounted at (e.g. "/", "/boot").
	Mountpoint string
	// Device is the backing device or source (e.g. "/dev/sda1", "tmpfs"). May be empty.
	Device string
	// Fstype is the filesystem type (e.g. "ext4", "vfat"). May be empty.
	Fstype string
	// TotalBytes is the total capacity of the filesystem.
	TotalBytes uint64
	// UsedBytes is the space in use (TotalBytes - available-to-unprivileged).
	UsedBytes uint64
	// FreeBytes is the space available to unprivileged users.
	FreeBytes uint64
}

// Percent returns the fraction of the filesystem in use in [0.0, 1.0].
// Returns 0 for a zero-capacity filesystem.
func (m MountUsage) Percent() float64 {
	if m.TotalBytes == 0 {
		return 0
	}
	p := float64(m.UsedBytes) / float64(m.TotalBytes)
	if p < 0 {
		return 0
	}
	if p > 1 {
		return 1
	}
	return p
}

// StorageSnapshot captures the usage of all monitored filesystems at a point in time.
// Mounts are ordered with the most-used filesystem first so space-constrained panels
// surface the most urgent filesystem.
type StorageSnapshot struct {
	Mounts    []MountUsage
	Timestamp time.Time
}
