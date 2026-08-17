package sparkline

import (
	"encoding/binary"
	"hash/fnv"
	"math"
)

// sign produces a deterministic uint64 hash of a Config for render cache memoization.
// It uses FNV-1a 64-bit hashing over all Config fields in declaration order.
func sign(cfg Config) uint64 {
	h := fnv.New64a()

	// Hash Data slice (length + each float64 element)
	binary.Write(h, binary.LittleEndian, int32(len(cfg.Data)))
	for _, v := range cfg.Data {
		binary.Write(h, binary.LittleEndian, math.Float64bits(v))
	}

	// Hash Style
	binary.Write(h, binary.LittleEndian, int32(cfg.Style))

	// Hash Bounds (Min.X, Min.Y, Max.X, Max.Y)
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Min.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Min.Y))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Max.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Max.Y))

	// Hash Foreground (R, G, B, A)
	h.Write([]byte{cfg.Foreground.R, cfg.Foreground.G, cfg.Foreground.B, cfg.Foreground.A})

	// Hash Background (R, G, B, A)
	h.Write([]byte{cfg.Background.R, cfg.Background.G, cfg.Background.B, cfg.Background.A})

	return h.Sum64()
}
