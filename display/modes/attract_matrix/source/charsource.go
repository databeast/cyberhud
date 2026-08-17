package source

import (
	"math/rand"
	"time"

	"github.com/databeast/cyberhud/display/widgets/marquee"
)

// Compile-time interface check.
var _ marquee.CharSource = (*RandomSource)(nil)

// characterPool contains the 40 runes used for matrix rain columns:
// digits 0-9 (48–57), uppercase A-Z (65–90), lowercase a-d (97–100).
var characterPool = []rune{
	// Digits 0-9
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	// Uppercase A-Z
	65, 66, 67, 68, 69, 70, 71, 72, 73, 74,
	75, 76, 77, 78, 79, 80, 81, 82, 83, 84,
	85, 86, 87, 88, 89, 90,
	// Lowercase a-d
	97, 98, 99, 100,
}

// defaultMutationInterval is the time a cached character remains stable
// before being replaced with a new random rune on the next CharAt call.
const defaultMutationInterval = 100 * time.Millisecond

// RandomSource is a marquee.CharSource that returns pseudorandom runes from
// the 40-character pool. Characters are cached per cell index and mutate after
// the configured mutation interval elapses.
//
// Performance note: mutation timing uses a frame-level timestamp set via
// SetFrameTime() rather than per-cell time.Now() calls, avoiding syscall
// overhead on embedded hardware.
type RandomSource struct {
	rng              *rand.Rand
	cells            map[int]rune
	lastMutation     map[int]int64 // UnixNano timestamp of last mutation per cell
	mutationInterval int64         // nanoseconds
}

// frameNow holds the current frame's timestamp in UnixNano.
// Updated once per BuildView call via setFrameTime(), eliminating per-cell
// time.Now() syscalls during rendering.
var frameNow int64

// setFrameTime records the current wall-clock time for use by all CharAt
// calls within this frame. Call once at the start of BuildView.
func SetFrameTime() {
	frameNow = time.Now().UnixNano()
}

// NewRandomSource creates a RandomSource for a given column. If seed == 0,
// non-deterministic seeding is used (time-based). Otherwise the seed is
// combined with colIdx via XOR to ensure different columns produce different
// sequences.
func NewRandomSource(seed int64, colIdx int, mutationInterval time.Duration) *RandomSource {
	var src rand.Source
	if seed == 0 {
		src = rand.NewSource(time.Now().UnixNano())
	} else {
		src = rand.NewSource(seed ^ int64(colIdx))
	}
	return &RandomSource{
		rng:              rand.New(src),
		cells:            make(map[int]rune),
		lastMutation:     make(map[int]int64),
		mutationInterval: mutationInterval.Nanoseconds(),
	}
}

// CharAt returns the character for the given logical cell index. If the cell
// has a cached rune that has not yet expired (frame time minus last mutation is
// less than the mutation interval), the cached value is returned. Otherwise a
// new random rune is selected from characterPool, cached, and returned.
//
// This uses the frame-level timestamp (frameNow) rather than calling time.Now()
// per cell, making it safe for high-cell-count panels on embedded hardware.
func (rs *RandomSource) CharAt(index int) rune {
	if r, ok := rs.cells[index]; ok {
		if frameNow-rs.lastMutation[index] < rs.mutationInterval {
			return r
		}
	}
	r := characterPool[rs.rng.Intn(len(characterPool))]
	rs.cells[index] = r
	rs.lastMutation[index] = frameNow
	return r
}
