package source

import "image"

// StemmaSnapshot captures the state needed by all stemma styles for rendering.
// Reuses the existing Device type from scanner.go.
type StemmaSnapshot struct {
	Devices []*Device
	GetIcon func(name string) (image.Image, bool)
}
