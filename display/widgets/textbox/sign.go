package textbox

import (
	"encoding/binary"
	"hash/fnv"
)

// sign produces a deterministic uint64 hash of a Config for render cache memoization.
// It uses FNV-1a 64-bit hashing over all Config fields in declaration order.
func sign(cfg Config) uint64 {
	h := fnv.New64a()

	// Hash Text
	h.Write([]byte(cfg.Text))

	// Hash Bounds (Min.X, Min.Y, Max.X, Max.Y)
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Min.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Min.Y))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Max.X))
	binary.Write(h, binary.LittleEndian, int32(cfg.Bounds.Max.Y))

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
