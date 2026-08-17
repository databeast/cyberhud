package scaledtextbox

import (
	"encoding/binary"
	"hash/fnv"
)

// sign produces a deterministic uint64 hash of a Config for render cache memoization.
// It uses FNV-1a 64-bit hashing over all Config fields in declaration order.
func sign(cfg Config) uint64 {
	h := fnv.New64a()

	// Hash LogicalSize (X, Y)
	binary.Write(h, binary.LittleEndian, int32(cfg.LogicalSize.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.LogicalSize.Y))

	// Hash TargetSize (X, Y)
	binary.Write(h, binary.LittleEndian, int32(cfg.TargetSize.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.TargetSize.Y))

	// Hash Position (X, Y)
	binary.Write(h, binary.LittleEndian, int32(cfg.Position.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.Position.Y))

	// Hash Text
	h.Write([]byte(cfg.Text))

	// Hash Font (using ID string for deterministic hashing)
	if cfg.Font != nil {
		h.Write([]byte(cfg.Font.ID()))
	}

	// Hash Alignment
	binary.Write(h, binary.LittleEndian, int32(cfg.Alignment))

	// Hash VAlign
	binary.Write(h, binary.LittleEndian, int32(cfg.VAlign))

	// Hash Overflow
	binary.Write(h, binary.LittleEndian, int32(cfg.Overflow))

	// Hash Foreground (R, G, B, A)
	h.Write([]byte{cfg.Foreground.R, cfg.Foreground.G, cfg.Foreground.B, cfg.Foreground.A})

	// Hash LineSpacing
	binary.Write(h, binary.LittleEndian, int32(cfg.LineSpacing))

	// Hash PadX
	binary.Write(h, binary.LittleEndian, int32(cfg.PadX))

	// Hash PadY
	binary.Write(h, binary.LittleEndian, int32(cfg.PadY))

	// Hash Border
	if cfg.Border {
		h.Write([]byte{1})
	} else {
		h.Write([]byte{0})
	}

	// Hash Label
	h.Write([]byte(cfg.Label))

	// Hash FontOverrides length
	binary.Write(h, binary.LittleEndian, int32(len(cfg.FontOverrides)))

	return h.Sum64()
}
