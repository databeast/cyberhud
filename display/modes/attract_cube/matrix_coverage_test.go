package attract_cube

import (
	"testing"

	geometric "github.com/databeast/cyberhud/display/modes/attract_geometric/styles"
)

func TestSurfaceMatrixCoversExistingAttractModes(t *testing.T) {
	for _, s := range geometric.Registry().Enumerate() {
		if cubeStyles.Lookup(s.Name()) == nil {
			t.Errorf("missing cube surface style %s", s.Name())
		}
	}
}
