package progressbar

import (
	"encoding/binary"
	"hash/fnv"
	"math"
)

// sign produces a deterministic uint64 hash of a Config for render cache memoization.
// It uses FNV-1a 64-bit hashing over all Config fields in declaration order.
func sign(cfg Config) uint64 {
	h := fnv.New64a()
	var buf [8]byte

	// Hash Style
	binary.LittleEndian.PutUint32(buf[:4], uint32(cfg.Style))
	h.Write(buf[:4])

	// Hash Orientation
	binary.LittleEndian.PutUint32(buf[:4], uint32(cfg.Orientation))
	h.Write(buf[:4])

	// Hash Value (float64 as uint64 bits)
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(cfg.Value))
	h.Write(buf[:])

	// Hash Bounds (Min.X, Min.Y, Max.X, Max.Y)
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.Bounds.Min.X)))
	h.Write(buf[:4])
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.Bounds.Min.Y)))
	h.Write(buf[:4])
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.Bounds.Max.X)))
	h.Write(buf[:4])
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.Bounds.Max.Y)))
	h.Write(buf[:4])

	// Hash Foreground (R, G, B, A)
	h.Write([]byte{cfg.Foreground.R, cfg.Foreground.G, cfg.Foreground.B, cfg.Foreground.A})

	// Hash Background (R, G, B, A)
	h.Write([]byte{cfg.Background.R, cfg.Background.G, cfg.Background.B, cfg.Background.A})

	// Hash Gradient
	if cfg.Gradient != nil {
		binary.LittleEndian.PutUint32(buf[:4], uint32(len(cfg.Gradient.Stops)))
		h.Write(buf[:4])
		for _, stop := range cfg.Gradient.Stops {
			binary.LittleEndian.PutUint64(buf[:], math.Float64bits(stop.Position))
			h.Write(buf[:])
			h.Write([]byte{stop.Color.R, stop.Color.G, stop.Color.B, stop.Color.A})
		}
	} else {
		binary.LittleEndian.PutUint32(buf[:4], 0)
		h.Write(buf[:4])
	}

	// Hash SegmentCount
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.SegmentCount)))
	h.Write(buf[:4])

	// Hash SegmentGap
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.SegmentGap)))
	h.Write(buf[:4])

	// Hash Thickness
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.Thickness)))
	h.Write(buf[:4])

	// Hash SweepAngle
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(cfg.SweepAngle))
	h.Write(buf[:])

	// Hash RoundedCaps
	if cfg.RoundedCaps {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	h.Write(buf[:1])

	// Hash BorderWidth
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.BorderWidth)))
	h.Write(buf[:4])

	// Hash BorderWall
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.BorderWall)))
	h.Write(buf[:4])

	// Hash BorderColor
	h.Write([]byte{cfg.BorderColor.R, cfg.BorderColor.G, cfg.BorderColor.B, cfg.BorderColor.A})

	// Hash Markers
	binary.LittleEndian.PutUint32(buf[:4], uint32(len(cfg.Markers)))
	h.Write(buf[:4])
	for _, m := range cfg.Markers {
		binary.LittleEndian.PutUint64(buf[:], math.Float64bits(m.Value))
		h.Write(buf[:])
		h.Write([]byte{m.Color.R, m.Color.G, m.Color.B, m.Color.A})
	}

	// Hash Animation config
	binary.LittleEndian.PutUint32(buf[:4], uint32(cfg.Animation.Type))
	h.Write(buf[:4])
	binary.LittleEndian.PutUint64(buf[:], uint64(cfg.Animation.Period))
	h.Write(buf[:])
	binary.LittleEndian.PutUint32(buf[:4], uint32(int32(cfg.Animation.Speed)))
	h.Write(buf[:4])

	// Hash animElapsed (as int64 nanoseconds)
	binary.LittleEndian.PutUint64(buf[:], uint64(cfg.animElapsed.Nanoseconds()))
	h.Write(buf[:])

	return h.Sum64()
}
