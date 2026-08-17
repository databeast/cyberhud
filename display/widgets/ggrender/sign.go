package ggrender

import (
	"encoding/binary"
	"hash/fnv"
)

// Sign produces a deterministic uint64 hash of a Config for render cache memoization.
// It uses FNV-1a 64-bit hashing over all Config fields in declaration order.
func Sign(cfg Config) uint64 {
	h := fnv.New64a()

	// Hash Bounds (Min.X, Min.Y, Max.X, Max.Y)
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Min.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Min.Y))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Max.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Max.Y))

	// Hash Label
	h.Write([]byte(cfg.Label))

	// Hash Color (R, G, B, A)
	h.Write([]byte{cfg.Color.R, cfg.Color.G, cfg.Color.B, cfg.Color.A})

	return h.Sum64()
}
