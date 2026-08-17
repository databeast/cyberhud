package source

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TripPoint represents a kernel-defined temperature threshold for a thermal zone.
type TripPoint struct {
	Type  string  // e.g., "critical", "hot", "passive"
	TempC float64 // temperature in degrees Celsius
}

// ZoneReading represents a single thermal zone's current state.
type ZoneReading struct {
	ZoneID     int
	Label      string
	TempC      float64
	TripPoints []TripPoint
}

// zoneHistory is a ring buffer storing the last 64 temperature samples for a zone.
type zoneHistory struct {
	samples [64]float64
	head    int // next write position
	count   int // number of populated slots (0..64)
}

// ZoneHistory exposes the thermal history ring buffer for package-level tests.
type ZoneHistory = zoneHistory

// Append adds a temperature sample to the ring buffer.
func (h *zoneHistory) Append(temp float64) {
	h.samples[h.head] = temp
	h.head = (h.head + 1) % 64
	if h.count < 64 {
		h.count++
	}
}

// History returns 64 values oldest-first, zero-padded for unpopulated slots.
func (h *zoneHistory) History() []float64 {
	result := make([]float64, 64)
	if h.count == 0 {
		return result
	}
	// Start reading from the oldest entry
	start := (h.head - h.count + 64) % 64
	// Leading zeros for unpopulated slots
	offset := 64 - h.count
	for i := 0; i < h.count; i++ {
		idx := (start + i) % 64
		result[offset+i] = h.samples[idx]
	}
	return result
}

// historyState holds per-zone sliding window histories, keyed by zone ID.
var historyState = struct {
	sync.RWMutex
	zones map[int]*zoneHistory
}{
	zones: make(map[int]*zoneHistory),
}

// RecordHistory appends a temperature sample for the given zone ID.
// Called by the sampler/refresh loop after each successful sample.
func RecordHistory(zoneID int, tempC float64) {
	historyState.Lock()
	defer historyState.Unlock()
	h, ok := historyState.zones[zoneID]
	if !ok {
		h = &zoneHistory{}
		historyState.zones[zoneID] = h
	}
	h.Append(tempC)
}

// getHistory returns the 64-sample history for a zone (oldest first, zero-padded).
func GetHistory(zoneID int) []float64 {
	historyState.RLock()
	defer historyState.RUnlock()
	h, ok := historyState.zones[zoneID]
	if !ok {
		return make([]float64, 64)
	}

	return h.History()
}

// ResetHistoryStateForTest clears all zone histories.
func ResetHistoryStateForTest() {
	historyState.Lock()
	defer historyState.Unlock()
	historyState.zones = make(map[int]*zoneHistory)
}

// loopState holds the lifecycle state for the background sampling goroutine.
var loopState struct {
	sync.Mutex
	stopCh  chan struct{}
	sampler Sampler
}

// Activate starts the background sampling loop using the default LinuxSampler.
// It is idempotent: calling Activate() when the loop is already running is a no-op.
func Activate(refreshMS func() int) {
	log.Printf("thermal: Activate() called")
	ActivateWith(NewLinuxSampler(), refreshMS)
}

// ActivateWith starts the background sampling loop using the provided Sampler.
// It is idempotent: calling ActivateWith() when the loop is already running is a no-op.
func ActivateWith(sampler Sampler, refreshMS func() int) {
	loopState.Lock()
	defer loopState.Unlock()

	// Already running â€” idempotent no-op.
	if loopState.stopCh != nil {
		log.Printf("thermal: ActivateWith() called but already running, skipping")
		return
	}

	log.Printf("thermal: starting sampling loop")

	// Perform the first sample synchronously so data is available immediately
	// on the same render tick that triggers activation.
	doSample(sampler)

	loopState.sampler = sampler
	loopState.stopCh = make(chan struct{})
	go samplingLoop(sampler, loopState.stopCh, refreshMS)
}

// Deactivate stops the background sampling loop.
// It is idempotent: calling Deactivate() when the loop is not running is a no-op.
func Deactivate() {
	loopState.Lock()
	ch := loopState.stopCh
	loopState.stopCh = nil
	loopState.sampler = nil
	loopState.Unlock()

	if ch != nil {
		close(ch)
	}
}

func SampleActive() (ThermalSnapshot, bool) {
	loopState.Lock()
	sampler := loopState.sampler
	loopState.Unlock()
	if sampler == nil {
		return ThermalSnapshot{}, false
	}
	snap, err := sampler.Sample()
	if err != nil || len(snap.Zones) == 0 {
		return ThermalSnapshot{}, false
	}
	UpdateSnapshot(snap)
	for _, z := range snap.Zones {
		RecordHistory(z.ZoneID, z.TempC)
	}
	return snap, true
}

// samplingLoop is the background goroutine that periodically samples thermal data.
// It exits when stopCh is closed. On each tick it calls sampler.Sample(), and on
// success updates the snapshot and records history for each zone.
// It re-reads GetPolicy().RefreshMS each tick to pick up policy changes dynamically.
func samplingLoop(sampler Sampler, stopCh <-chan struct{}, refreshMSFunc func() int) {
	refreshMS := clampRefreshMS(refreshMSFunc)
	ticker := time.NewTicker(time.Duration(refreshMS) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			doSample(sampler)

			// Re-read policy to pick up RefreshMS changes without restarting the loop.
			newRefreshMS := clampRefreshMS(refreshMSFunc)
			if newRefreshMS != refreshMS {
				refreshMS = newRefreshMS
				ticker.Reset(time.Duration(refreshMS) * time.Millisecond)
			}

		}
	}
}

func clampRefreshMS(refreshMSFunc func() int) int {
	ms := DefaultPolicy().RefreshMS
	if refreshMSFunc != nil {
		ms = refreshMSFunc()
	}
	if ms <= 0 {
		ms = 2000
	}
	if ms < 500 {
		ms = 500
	}
	if ms > 120000 {
		ms = 120000
	}
	return ms
}

// doSample calls the sampler and, on success, updates the snapshot state and records
// history for each zone. On error it does nothing (preserving existing snapshot state).
func doSample(sampler Sampler) {
	snap, err := sampler.Sample()
	if err != nil {
		log.Printf("thermal: Sample() error: %v", err)
		return
	}
	UpdateSnapshot(snap)
	for _, z := range snap.Zones {
		RecordHistory(z.ZoneID, z.TempC)
	}
}

// snapshotState holds the latest sampled thermal snapshot, protected by a mutex.
var snapshotState = struct {
	sync.RWMutex
	snap ThermalSnapshot
}{}

// CurrentSnapshot returns the most recent thermal snapshot (thread-safe).
// If no sample has been taken yet, returns an empty ThermalSnapshot.
func CurrentSnapshot() ThermalSnapshot {
	snapshotState.RLock()
	defer snapshotState.RUnlock()
	return snapshotState.snap
}

// UpdateSnapshot stores a new thermal snapshot (thread-safe).
// Called by the refresh loop when a new sample is collected.
func UpdateSnapshot(snap ThermalSnapshot) {
	snapshotState.Lock()
	defer snapshotState.Unlock()
	snapshotState.snap = snap
}

// Sampler abstracts thermal data collection, enabling mock implementations
// for testing without access to /sys/class/thermal/.
type Sampler interface {
	Sample() (ThermalSnapshot, error)
}

// LinuxSampler implements Sampler by reading from the Linux sysfs thermal
// interface under /sys/class/thermal/thermal_zone*.
type LinuxSampler struct {
	baseDir string // defaults to /sys/class/thermal
}

// NewLinuxSampler returns a LinuxSampler reading from the default sysfs path.
func NewLinuxSampler() *LinuxSampler {
	return &LinuxSampler{baseDir: "/sys/class/thermal"}
}

// NewLinuxSamplerAt returns a LinuxSampler reading from the specified base
// directory. This is useful for testing with a mock sysfs tree.
func NewLinuxSamplerAt(baseDir string) *LinuxSampler {
	return &LinuxSampler{baseDir: baseDir}
}

// Sample reads all thermal zones from sysfs and returns a ThermalSnapshot.
// Zones with unreadable or non-numeric temp files are skipped gracefully.
func (s *LinuxSampler) Sample() (ThermalSnapshot, error) {
	pattern := filepath.Join(s.baseDir, "thermal_zone*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return ThermalSnapshot{Timestamp: time.Now()}, err
	}

	var zones []ZoneReading
	for _, zonePath := range matches {
		reading, ok := s.readZone(zonePath)
		if !ok {
			continue
		}
		zones = append(zones, reading)
	}

	// Sort by ZoneID for deterministic ordering.
	sort.Slice(zones, func(i, j int) bool {
		return zones[i].ZoneID < zones[j].ZoneID
	})

	return ThermalSnapshot{
		Zones:     zones,
		Timestamp: time.Now(),
	}, nil
}

// readZone reads a single thermal zone directory and returns a ZoneReading.
// Returns false if the zone's temp file is unreadable or non-numeric.
func (s *LinuxSampler) readZone(zonePath string) (ZoneReading, bool) {
	// Read temperature (required).
	tempPath := filepath.Join(zonePath, "temp")
	tempData, err := os.ReadFile(tempPath)
	if err != nil {
		return ZoneReading{}, false
	}
	millideg, err := strconv.ParseInt(strings.TrimSpace(string(tempData)), 10, 64)
	if err != nil {
		return ZoneReading{}, false
	}
	tempC := float64(millideg) / 1000.0

	// Parse zone ID from directory name (e.g., "thermal_zone0" â†’ 0).
	dirName := filepath.Base(zonePath)
	zoneID := parseZoneID(dirName)

	// Read type label (optional, fallback to dirname).
	label := s.readLabel(zonePath, dirName)

	// Read trip points (optional, up to 20).
	tripPoints := s.readTripPoints(zonePath)

	return ZoneReading{
		ZoneID:     zoneID,
		Label:      label,
		TempC:      tempC,
		TripPoints: tripPoints,
	}, true
}

// readLabel reads the "type" file from a zone directory. If unreadable or
// empty, it falls back to the directory basename.
func (s *LinuxSampler) readLabel(zonePath, dirName string) string {
	typePath := filepath.Join(zonePath, "type")
	data, err := os.ReadFile(typePath)
	if err != nil {
		return dirName
	}
	label := strings.TrimSpace(string(data))
	if label == "" {
		return dirName
	}
	return label
}

// readTripPoints reads up to 20 trip points from a zone directory by globbing
// trip_point_*_temp files and reading matching trip_point_*_type files.
func (s *LinuxSampler) readTripPoints(zonePath string) []TripPoint {
	pattern := filepath.Join(zonePath, "trip_point_*_temp")
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return nil
	}

	// Sort matches for deterministic ordering.
	sort.Strings(matches)

	var tripPoints []TripPoint
	for _, tempFile := range matches {
		if len(tripPoints) >= 20 {
			break
		}

		// Read trip point temperature.
		data, err := os.ReadFile(tempFile)
		if err != nil {
			continue
		}
		millideg, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		if err != nil {
			continue
		}
		tripTempC := float64(millideg) / 1000.0

		// Derive the type file path from the temp file path.
		// e.g., trip_point_0_temp â†’ trip_point_0_type
		typeFile := strings.TrimSuffix(tempFile, "_temp") + "_type"
		tripType := readTrimmed(typeFile)

		tripPoints = append(tripPoints, TripPoint{
			Type:  tripType,
			TempC: tripTempC,
		})
	}

	return tripPoints
}

// parseZoneID extracts the numeric zone ID from a directory name like
// "thermal_zone0". Returns 0 if parsing fails.
func parseZoneID(dirName string) int {
	const prefix = "thermal_zone"
	if !strings.HasPrefix(dirName, prefix) {
		return 0
	}
	id, err := strconv.Atoi(dirName[len(prefix):])
	if err != nil {
		return 0
	}
	return id
}

// readTrimmed reads a file and returns its content trimmed of whitespace.
// Returns an empty string if the file is unreadable.
func readTrimmed(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
