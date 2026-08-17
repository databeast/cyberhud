package led

import (
	"encoding/binary"
	"hash/fnv"
	"math"
)

// sign produces a deterministic uint64 hash of a Config for render cache memoization.
// It uses FNV-1a 64-bit hashing over all Config fields in declaration order.
// The function is a pure computation with no side effects and completes in O(N) time
// where N is the number of group entries (max 32), or O(1) for single-LED configs.
func sign(cfg Config) uint64 {
	h := fnv.New64a()
	buf := make([]byte, 8) // reusable buffer for binary encoding

	// 1. Shape (int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.Shape)))
	h.Write(buf[:4])

	// 2. State (int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.State)))
	h.Write(buf[:4])

	// 3. Brightness (float64 via Float64bits → uint64)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(cfg.Brightness))
	h.Write(buf[:8])

	// 4. Diameter (int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.Diameter)))
	h.Write(buf[:4])

	// 5. Bounds (Min.X, Min.Y, Max.X, Max.Y as int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.Bounds.Min.X)))
	h.Write(buf[:4])
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.Bounds.Min.Y)))
	h.Write(buf[:4])
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.Bounds.Max.X)))
	h.Write(buf[:4])
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.Bounds.Max.Y)))
	h.Write(buf[:4])

	// 6. Foreground (4 bytes: R, G, B, A)
	h.Write([]byte{cfg.Foreground.R, cfg.Foreground.G, cfg.Foreground.B, cfg.Foreground.A})

	// 7. Background (4 bytes: R, G, B, A)
	h.Write([]byte{cfg.Background.R, cfg.Background.G, cfg.Background.B, cfg.Background.A})

	// 8. WarningColor (4 bytes: R, G, B, A)
	h.Write([]byte{cfg.WarningColor.R, cfg.WarningColor.G, cfg.WarningColor.B, cfg.WarningColor.A})

	// 9. Gradient (length prefix uint32 + per-stop: position as float64 + color 4 bytes)
	if cfg.Gradient != nil {
		binary.LittleEndian.PutUint32(buf, uint32(len(cfg.Gradient.Stops)))
		h.Write(buf[:4])
		for _, stop := range cfg.Gradient.Stops {
			binary.LittleEndian.PutUint64(buf, math.Float64bits(stop.Position))
			h.Write(buf[:8])
			h.Write([]byte{stop.Color.R, stop.Color.G, stop.Color.B, stop.Color.A})
		}
	} else {
		binary.LittleEndian.PutUint32(buf, 0)
		h.Write(buf[:4])
	}

	// 10. GlowEnabled (1 byte)
	if cfg.GlowEnabled {
		h.Write([]byte{1})
	} else {
		h.Write([]byte{0})
	}

	// 11. GlowRadius (int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.GlowRadius)))
	h.Write(buf[:4])

	// 12. BorderWidth (int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.BorderWidth)))
	h.Write(buf[:4])

	// 13. BorderColor (4 bytes: R, G, B, A)
	h.Write([]byte{cfg.BorderColor.R, cfg.BorderColor.G, cfg.BorderColor.B, cfg.BorderColor.A})

	// 14. ShineStyle (int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.ShineStyle)))
	h.Write(buf[:4])

	// 15. ShineOpacity (1 byte)
	h.Write([]byte{cfg.ShineOpacity})

	// 16. Animation.Type (int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.Animation.Type)))
	h.Write(buf[:4])

	// 17. Animation.Period (int64 nanoseconds)
	binary.LittleEndian.PutUint64(buf, uint64(cfg.Animation.Period.Nanoseconds()))
	h.Write(buf[:8])

	// 18. Animation.MinBrightness (float64)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(cfg.Animation.MinBrightness))
	h.Write(buf[:8])

	// 19. Group (length prefix + per-entry fields)
	binary.LittleEndian.PutUint32(buf, uint32(len(cfg.Group)))
	h.Write(buf[:4])
	for _, entry := range cfg.Group {
		// State (int32)
		binary.LittleEndian.PutUint32(buf, uint32(int32(entry.State)))
		h.Write(buf[:4])
		// Foreground (4 bytes)
		h.Write([]byte{entry.Foreground.R, entry.Foreground.G, entry.Foreground.B, entry.Foreground.A})
		// WarningColor (4 bytes)
		h.Write([]byte{entry.WarningColor.R, entry.WarningColor.G, entry.WarningColor.B, entry.WarningColor.A})
		// Shape (int32)
		binary.LittleEndian.PutUint32(buf, uint32(int32(entry.Shape)))
		h.Write(buf[:4])
		// BorderWidth (int32)
		binary.LittleEndian.PutUint32(buf, uint32(int32(entry.BorderWidth)))
		h.Write(buf[:4])
		// BorderColor (4 bytes)
		h.Write([]byte{entry.BorderColor.R, entry.BorderColor.G, entry.BorderColor.B, entry.BorderColor.A})
		// GlowEnabled (1 byte)
		if entry.GlowEnabled {
			h.Write([]byte{1})
		} else {
			h.Write([]byte{0})
		}
		// GlowRadius (int32)
		binary.LittleEndian.PutUint32(buf, uint32(int32(entry.GlowRadius)))
		h.Write(buf[:4])
	}

	// 20. Orientation (int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.Orientation)))
	h.Write(buf[:4])

	// 21. Spacing (int32)
	binary.LittleEndian.PutUint32(buf, uint32(int32(cfg.Spacing)))
	h.Write(buf[:4])

	// 22. animElapsed (int64 nanoseconds)
	binary.LittleEndian.PutUint64(buf, uint64(cfg.animElapsed.Nanoseconds()))
	h.Write(buf[:8])

	return h.Sum64()
}
