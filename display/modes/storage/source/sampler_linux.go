//go:build linux

package source

import (
	"bufio"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// linuxSampler reads mounted filesystems from /proc/mounts and queries each with
// statfs(2). It is the production Sampler on Linux.
type linuxSampler struct {
	mountsPath string // defaults to /proc/mounts; overridable for tests
}

func newSampler() Sampler { return &linuxSampler{mountsPath: "/proc/mounts"} }

// pseudoFstypes are filesystem types that do not represent real block storage.
// They are hidden unless show_all is set.
var pseudoFstypes = map[string]bool{
	"proc": true, "sysfs": true, "devtmpfs": true, "devpts": true,
	"tmpfs": true, "cgroup": true, "cgroup2": true, "pstore": true,
	"bpf": true, "tracefs": true, "debugfs": true, "securityfs": true,
	"mqueue": true, "hugetlbfs": true, "overlay": true, "squashfs": true,
	"autofs": true, "fusectl": true, "configfs": true, "ramfs": true,
	"binfmt_misc": true, "efivarfs": true, "nsfs": true,
}

func (s *linuxSampler) Sample(showAll bool) (StorageSnapshot, error) {
	f, err := os.Open(s.mountsPath)
	if err != nil {
		return StorageSnapshot{Timestamp: time.Now()}, err
	}
	defer f.Close()

	seen := make(map[string]bool)
	var mounts []MountUsage

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}
		device := unescapeMount(fields[0])
		mountpoint := unescapeMount(fields[1])
		fstype := fields[2]

		if !showAll && pseudoFstypes[fstype] {
			continue
		}
		// De-duplicate by mountpoint (bind mounts can repeat entries).
		if seen[mountpoint] {
			continue
		}

		usage, ok := statfsUsage(mountpoint)
		if !ok || usage.TotalBytes == 0 {
			continue
		}
		usage.Mountpoint = mountpoint
		usage.Device = device
		usage.Fstype = fstype

		seen[mountpoint] = true
		mounts = append(mounts, usage)
	}
	if err := scanner.Err(); err != nil {
		return StorageSnapshot{Timestamp: time.Now()}, err
	}

	sortMounts(mounts)
	return StorageSnapshot{Mounts: mounts, Timestamp: time.Now()}, nil
}

// statfsUsage queries statfs(2) for a mountpoint and returns byte totals.
// Used/free are computed against the blocks available to unprivileged users,
// matching what `df` reports.
func statfsUsage(mountpoint string) (MountUsage, bool) {
	var st unix.Statfs_t
	if err := unix.Statfs(mountpoint, &st); err != nil {
		return MountUsage{}, false
	}
	bsize := uint64(st.Bsize)
	if bsize == 0 {
		return MountUsage{}, false
	}
	total := st.Blocks * bsize
	avail := st.Bavail * bsize
	// Used = total minus free-to-root; reserved blocks count as used so the bar
	// matches df's "Use%".
	freeToRoot := st.Bfree * bsize
	used := total - freeToRoot
	return MountUsage{
		TotalBytes: total,
		UsedBytes:  used,
		FreeBytes:  avail,
	}, true
}

// unescapeMount decodes the octal escapes used in /proc/mounts for spaces (\040),
// tabs (\011), newlines (\012) and backslashes (\134).
func unescapeMount(s string) string {
	if !strings.ContainsRune(s, '\\') {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+3 < len(s) {
			// Attempt a 3-digit octal escape (e.g. \040 for space).
			if v, ok := parseOctal(s[i+1 : i+4]); ok {
				b.WriteByte(v)
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func parseOctal(s string) (byte, bool) {
	if len(s) != 3 {
		return 0, false
	}
	v := 0
	for _, c := range s {
		if c < '0' || c > '7' {
			return 0, false
		}
		v = v*8 + int(c-'0')
	}
	if v > 255 {
		return 0, false
	}
	return byte(v), true
}
