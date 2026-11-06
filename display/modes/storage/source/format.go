package source

import "fmt"

// FormatBytes renders a byte count as a short human-readable string according to
// the unit policy:
//
//	"si"   → 1000-based with KB/MB/GB/TB suffixes
//	"iec"  → 1024-based with KiB/MiB/GiB/TiB suffixes
//	"auto" → 1024-based with KB/MB/GB/TB suffixes (the common, compact default)
//
// Values are shown with one decimal place for sub-100 magnitudes and no decimals
// above, keeping each value narrow for small panels (e.g. "931G", "12.4G").
func FormatBytes(n uint64, unit string) string {
	var base float64
	var suffixes []string
	switch unit {
	case "si":
		base = 1000
		suffixes = []string{"B", "KB", "MB", "GB", "TB", "PB"}
	case "iec":
		base = 1024
		suffixes = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	default: // "auto"
		base = 1024
		suffixes = []string{"B", "K", "M", "G", "T", "P"}
	}

	v := float64(n)
	i := 0
	for v >= base && i < len(suffixes)-1 {
		v /= base
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%d%s", n, suffixes[0])
	}
	if v < 10 {
		return fmt.Sprintf("%.1f%s", v, suffixes[i])
	}
	if v < 100 {
		return fmt.Sprintf("%.0f%s", v, suffixes[i])
	}
	return fmt.Sprintf("%.0f%s", v, suffixes[i])
}
