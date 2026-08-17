package source

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"sort"
	"time"

	fonts "github.com/databeast/cyberhud/display/surface/fonts"
	"github.com/databeast/cyberhud/display/widgets/marquee"
)

// computeColumnCount returns the number of full-width columns that fit within
// panelWidth given a glyph advance. Returns 0 if glyphAdvance <= 0 or
// panelWidth < glyphAdvance.
func computeColumnCount(panelWidth, glyphAdvance int) int {
	if glyphAdvance <= 0 || panelWidth < glyphAdvance {
		return 0
	}
	return panelWidth / glyphAdvance
}

// computeVisibleCells returns the number of vertical cells that fit within
// panelHeight given a row height. Returns 0 if rowHeight <= 0.
func computeVisibleCells(panelHeight, rowHeight int) int {
	if rowHeight <= 0 {
		return 0
	}
	return panelHeight / rowHeight
}

// computeActiveColumns returns a sorted slice of column indices that should be
// active given the total column count and density. The selection uses a random
// permutation seeded by totalColumns so the same density produces the same
// columns across frames. If totalColumns == 1, always returns [0].
func computeActiveColumns(totalColumns int, density float64) []int {
	if totalColumns <= 0 {
		return nil
	}
	if totalColumns == 1 {
		return []int{0}
	}

	activeCount := int(math.Floor(float64(totalColumns) * density))
	if activeCount < 1 {
		activeCount = 1
	}
	if activeCount >= totalColumns {
		// All columns active — return full range.
		cols := make([]int, totalColumns)
		for i := range cols {
			cols[i] = i
		}
		return cols
	}

	// Seeded permutation for deterministic column selection.
	rng := rand.New(rand.NewSource(int64(totalColumns)))
	perm := rng.Perm(totalColumns)

	cols := make([]int, activeCount)
	for i := 0; i < activeCount; i++ {
		cols[i] = perm[i]
	}
	sort.Ints(cols)
	return cols
}

// Package-level strip cache. Strips persist between frames. When density
// changes, columns are added/removed incrementally: existing strips keep
// their current scroll position, new strips start at phase 0 (entering from
// the top of the screen).
var (
	stripCache            []*marquee.Strip
	stripColumnMap        map[int]*marquee.Strip // colIdx → existing strip
	lastPolicyFingerprint string
	lastCacheKey          string
	firstBuildDone        bool // tracks whether the initial staggered build has happened
)

// buildCacheKey constructs a composite key from the layout-affecting parameters.
// Only density (which changes active column count) triggers a rebuild — speed
// changes are applied via Tick without rebuilding strips.
func buildCacheKey(p Policy, panelWidth, panelHeight, glyphAdvance, rowHeight int, mono bool) string {
	return fmt.Sprintf("%v|%d|%d|%d|%d|%d|%v",
		p.Density, p.TrailLength, panelWidth, panelHeight, glyphAdvance, rowHeight, mono)
}

// rebuildStrips creates or updates the strip set for the current active columns.
// On first call, strips get random phase offsets for visual staggering.
// On subsequent calls (density changes), existing strips are preserved and new
// strips start at phase 0 so they scroll in from the top naturally.
// The face parameter is the catalog-validated font face resolved through hints.Face.
func rebuildStrips(p Policy, panelWidth, panelHeight, glyphAdvance, rowHeight int, mono bool, seed int64, face fonts.Face) []*marquee.Strip {
	columnCount := computeColumnCount(panelWidth, glyphAdvance)
	visibleCells := computeVisibleCells(panelHeight, rowHeight)

	if visibleCells == 0 {
		stripCache = nil
		stripColumnMap = nil
		return nil
	}

	activeColumns := computeActiveColumns(columnCount, p.Density)

	// Clamp trail length.
	clampedTrailLength := p.TrailLength
	if clampedTrailLength < 4 {
		clampedTrailLength = 4
	}
	if clampedTrailLength > visibleCells {
		clampedTrailLength = visibleCells
	}

	colors := buildColorArray(clampedTrailLength, mono)

	// Font face is passed in from the hints.Face (catalog-validated).

	// PRNG for speed (and initial phase on first build).
	var rng *rand.Rand
	if seed == 0 {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		rng = rand.New(rand.NewSource(seed))
	}

	// Build the new set, preserving existing strips by column index.
	newMap := make(map[int]*marquee.Strip, len(activeColumns))
	strips := make([]*marquee.Strip, 0, len(activeColumns))

	for _, colIdx := range activeColumns {
		// Reuse existing strip if this column was already active.
		if existing, ok := stripColumnMap[colIdx]; ok {
			newMap[colIdx] = existing
			strips = append(strips, existing)
			continue
		}

		// New column: random speed, but phase depends on whether this is
		// the initial build (staggered) or a mid-cycle addition (start from top).
		speed := p.MinSpeed + rng.Float64()*(p.MaxSpeed-p.MinSpeed)
		var phase float64
		if !firstBuildDone {
			// Initial startup: stagger columns so they don't all start together.
			phase = float64(rng.Intn(visibleCells*2 + 1))
		}
		// else: phase = 0, strip enters from top of screen.

		cfg := marquee.Config{
			Bounds:    image.Rect(colIdx*glyphAdvance, 0, (colIdx+1)*glyphAdvance, panelHeight),
			Direction: marquee.Vertical,
			Font:      face,
			Source:    NewRandomSource(seed, colIdx, defaultMutationInterval),
			Colors:    colors,
			Speed:     speed,
			Phase:     phase,
		}
		strip := marquee.New(cfg)
		newMap[colIdx] = strip
		strips = append(strips, strip)
	}

	firstBuildDone = true
	stripColumnMap = newMap
	stripCache = strips
	return strips
}

// getOrRebuildStrips returns the cached strips if the cache key is unchanged,
// otherwise incrementally updates the strip set.
// The face parameter is the catalog-validated font face resolved through hints.Face.
func GetOrRebuildStrips(p Policy, panelWidth, panelHeight, glyphAdvance, rowHeight int, mono bool, seed int64, face fonts.Face) []*marquee.Strip {
	key := buildCacheKey(p, panelWidth, panelHeight, glyphAdvance, rowHeight, mono)
	if key == lastCacheKey && stripCache != nil {
		return stripCache
	}
	lastCacheKey = key
	lastPolicyFingerprint = PolicyFingerprint(p)
	return rebuildStrips(p, panelWidth, panelHeight, glyphAdvance, rowHeight, mono, seed, face)
}
