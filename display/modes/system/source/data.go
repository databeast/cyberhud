package source

import (
	"image"

	"github.com/databeast/cyberhud/display/modes/system/sampler"
)

// SystemSnapshot captures all state fields needed by system mode styles for rendering.
type SystemSnapshot struct {
	Hostname  string
	OSArch    string
	Uptime    string
	IPs       []string
	CPUSample []float64
	Processes []sampler.ProcessEntry
	GetIcon   func(name string) (image.Image, bool)
}
