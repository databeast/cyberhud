package console_test

import (
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/coordinator"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	"github.com/databeast/cyberhud/display/region"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
	"github.com/databeast/cyberhud/runtime/console"

	// image mode registers an enum option "fit" with an Allowed list.
	_ "github.com/databeast/cyberhud/display/modes/image"
)

// newOptionsTestServer starts a console server bound to a region whose allowed
// modes include "image" (which advertises an enum option with an Allowed list),
// so "display modes" has option metadata with allowed values to emit.
func newOptionsTestServer(t *testing.T) string {
	t.Helper()
	sockPath := filepath.Join(t.TempDir(), "cyberhudd_opt.sock")
	if runtime.GOOS == "windows" {
		sockPath = filepath.Join(t.TempDir(), "o.sock")
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 240, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}
	rm := region.NewRegionManager(vd)
	if err := rm.Allocate(region.RegionSpec{
		Name: "main", Controller: "st7789",
		Bounds: image.Rect(0, 0, 240, 240),
		Modes:  []string{"image"}, DefaultMode: "image",
	}); err != nil {
		t.Fatalf("Allocate: %v", err)
	}

	modes := coordinator.NewState()
	modes.Bind(rm)

	scanner := source.New(nil, 0)
	gm := gpiomgr.New()
	srv := console.New(sockPath, func() *source.Scanner { return scanner }, gm, func() string { return "OK" }, modes, nil, nil, "")
	if err := srv.Start(); err != nil {
		t.Fatalf("server start: %v", err)
	}
	t.Cleanup(func() {
		srv.Stop()
		_ = os.Remove(sockPath)
	})
	return sockPath
}

func TestConsoleServer_DisplayModesEmitsAllowed(t *testing.T) {
	sockPath := newOptionsTestServer(t)
	time.Sleep(20 * time.Millisecond)

	lines := sendAndDrain(t, sockPath, "display modes")
	joined := strings.Join(lines, "\n")

	// The image mode's "fit" option is an enum (truncate/scale/stretch).
	var optLine string
	for _, l := range lines {
		if strings.Contains(l, "option=fit") {
			optLine = strings.TrimSpace(l)
			break
		}
	}
	if optLine == "" {
		t.Fatalf("no option=fit line in display modes output:\n%s", joined)
	}
	for _, want := range []string{`type=string`, `allowed="`, "truncate", "scale", "stretch"} {
		if !strings.Contains(optLine, want) {
			t.Errorf("option=fit line missing %q: %q", want, optLine)
		}
	}

	// The mode line for image should carry its UI category.
	var modeLine string
	for _, l := range lines {
		if strings.Contains(l, "mode=image ") || strings.HasSuffix(strings.TrimSpace(l), "mode=image") {
			modeLine = strings.TrimSpace(l)
			break
		}
	}
	if modeLine == "" {
		t.Fatalf("no mode=image line in display modes output:\n%s", joined)
	}
	if !strings.Contains(modeLine, `category="Data"`) {
		t.Errorf("mode=image line missing category: %q", modeLine)
	}
}
