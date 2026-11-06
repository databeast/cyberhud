package console_test

import (
	"bufio"
	"errors"
	"image"
	"io"
	"net"
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

	_ "github.com/databeast/cyberhud/display/modes/clock"
)

// newGeometryTestServer starts a console server backed by a coordinator.State
// bound to a live RegionManager with two allocated regions, so that the
// "display geometry" command has real bounds to report. It returns the socket path.
func newGeometryTestServer(t *testing.T) string {
	t.Helper()
	sockPath := filepath.Join(t.TempDir(), "cyberhudd_geo.sock")
	if runtime.GOOS == "windows" {
		sockPath = filepath.Join(t.TempDir(), "g.sock")
	}

	vd, err := region.NewVirtualDisplay(image.Rect(0, 0, 320, 240))
	if err != nil {
		t.Fatalf("NewVirtualDisplay: %v", err)
	}
	rm := region.NewRegionManager(vd)
	specs := []region.RegionSpec{
		{Name: "main", Controller: "st7789", Bounds: image.Rect(0, 0, 240, 240), Modes: []string{"clock"}, DefaultMode: "clock"},
		{Name: "aux", Controller: "st7735s", Bounds: image.Rect(240, 0, 320, 120), Modes: []string{"clock"}, DefaultMode: "clock"},
	}
	for _, spec := range specs {
		if err := rm.Allocate(spec); err != nil {
			t.Fatalf("Allocate %q: %v", spec.Name, err)
		}
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

// sendAndDrain sends a single command followed by quit, then reads all response
// lines until "OK bye" or EOF. This mirrors cyberhudctl's single-command path and
// is the correct way to consume a multi-line OK-block response (where the first
// line is a bare "OK" and the data follows on indented lines).
func sendAndDrain(t *testing.T, sockPath, cmd string) []string {
	t.Helper()
	conn, err := net.DialTimeout("unix", sockPath, 2*time.Second)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

	r := bufio.NewReader(conn)
	if _, err := r.ReadString('\n'); err != nil { // greeting
		t.Fatalf("read greeting: %v", err)
	}
	if _, err := io.WriteString(conn, cmd+"\n"); err != nil {
		t.Fatalf("write cmd: %v", err)
	}
	if _, err := io.WriteString(conn, "quit\n"); err != nil {
		t.Fatalf("write quit: %v", err)
	}

	var lines []string
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatalf("read: %v", err)
		}
		line = strings.TrimRight(line, "\r\n")
		if strings.TrimSpace(line) == "OK bye" {
			break
		}
		lines = append(lines, line)
	}
	return lines
}

func TestConsoleServer_DisplayGeometry(t *testing.T) {
	sockPath := newGeometryTestServer(t)
	time.Sleep(20 * time.Millisecond)

	lines := sendAndDrain(t, sockPath, "display geometry")

	// Expect: "OK", "  virtual w=320 h=240", "  region=main.0 ...", "  region=aux.0 ..."
	if len(lines) != 4 {
		t.Fatalf("expected 4 response lines, got %d:\n%s", len(lines), strings.Join(lines, "\n"))
	}
	if strings.TrimSpace(lines[0]) != "OK" {
		t.Errorf("line 0: expected bare OK, got %q", lines[0])
	}
	if strings.TrimSpace(lines[1]) != "virtual w=320 h=240" {
		t.Errorf("line 1: expected virtual size, got %q", lines[1])
	}

	region0 := strings.TrimSpace(lines[2])
	for _, want := range []string{`region=main.0`, `index=0`, `name="main"`, `x=0`, `y=0`, `w=240`, `h=240`} {
		if !strings.Contains(region0, want) {
			t.Errorf("region line 0 %q: missing %q", region0, want)
		}
	}

	region1 := strings.TrimSpace(lines[3])
	for _, want := range []string{`region=aux.0`, `index=1`, `name="aux"`, `x=240`, `y=0`, `w=80`, `h=120`} {
		if !strings.Contains(region1, want) {
			t.Errorf("region line 1 %q: missing %q", region1, want)
		}
	}
}

func TestConsoleServer_DisplayGeometryUsage(t *testing.T) {
	sockPath := newGeometryTestServer(t)
	time.Sleep(20 * time.Millisecond)

	lines := sendAndDrain(t, sockPath, "display geometry extra")
	if len(lines) == 0 {
		t.Fatal("expected a usage error response")
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[0]), "ERR usage: display geometry") {
		t.Fatalf("expected usage error, got %q", lines[0])
	}
}
