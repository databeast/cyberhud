package zmq

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/modes/testsnapshot"
	"github.com/databeast/cyberhud/display/modes/zmq/content"
	"github.com/databeast/cyberhud/display/style"
)

var snapshotOutputDir = filepath.Join("snapshots")

func categoryFromStyle(name string, reqs style.SurfaceRequirements) testsnapshot.DisplayCategory {
	switch {
	case strings.HasPrefix(name, "color-"):
		return testsnapshot.CategoryColor
	case strings.HasPrefix(name, "grayscale-"):
		return testsnapshot.CategoryGrayscale
	case strings.HasPrefix(name, "eink-"), strings.HasPrefix(name, "mono-slow-"):
		return testsnapshot.CategoryEink
	case strings.HasPrefix(name, "mono-"):
		return testsnapshot.CategoryMono
	}
	switch reqs.Capability {
	case style.GrayscaleSlow, style.GrayscaleFast:
		return testsnapshot.CategoryGrayscale
	case style.ColorSlow, style.ColorFast:
		return testsnapshot.CategoryColor
	case style.MonoFast:
		return testsnapshot.CategoryMono
	default:
		return testsnapshot.CategoryColor
	}
}

func dimensionsFromRequirements(reqs style.SurfaceRequirements) (int, int) {
	width, height := reqs.MinWidth, reqs.MinHeight
	if width == 0 {
		width = reqs.PreferredWidth
	}
	if height == 0 {
		height = reqs.PreferredHeight
	}
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 240
	}
	return width, height
}

func seedSnapshotMessages() {
	content.ClearMessagesForTest()
	for _, msg := range []string{
		`{"temp":42,"status":"ok"}`,
		"alerts: link up",
		"metrics: cpu=17 mem=42",
	} {
		content.PushMessageForTest(msg)
	}
}

// TestZMQPNGSnapshots enumerates all registered ZMQ styles and renders one
// deterministic PNG for each through the production snapshot pipeline.
func TestZMQPNGSnapshots(t *testing.T) {
	styles := zmqRegistry.Enumerate()
	if len(styles) == 0 {
		t.Fatal("zmq registry contains zero styles")
	}

	if err := os.RemoveAll(snapshotOutputDir); err != nil {
		t.Fatalf("failed to clean snapshot output directory: %v", err)
	}
	if err := os.MkdirAll(snapshotOutputDir, 0o755); err != nil {
		t.Fatalf("failed to create snapshot output directory: %v", err)
	}

	for _, s := range styles {
		s := s
		t.Run(s.Name(), func(t *testing.T) {
			reqs := s.Requirements()
			width, height := dimensionsFromRequirements(reqs)
			policy := content.DefaultPolicy()
			policy.Style = s.Name()
			policy.Font = "auto"
			policy.MaxLines = 24

			pngPath := testsnapshot.RenderSnapshot(t,
				testsnapshot.WithMode("zmq"),
				testsnapshot.WithDimensions(width, height),
				testsnapshot.WithDisplayCategory(categoryFromStyle(s.Name(), reqs)),
				testsnapshot.WithOutputDir(snapshotOutputDir),
				testsnapshot.WithBasename(s.Name()),
				testsnapshot.WithReset(ResetTestState),
				testsnapshot.WithPreRender(func() {
					content.SetPolicy(policy)
					seedSnapshotMessages()
				}),
			)

			testsnapshot.VerifyAll(t, pngPath, width, height)
		})
	}
}
