package console_test

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/coordinator"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
	"github.com/databeast/cyberhud/runtime/console"
	"pgregory.net/rapid"
)

// For any set of configured regions (with surface names, indices, current modes,
// and available modes), the formatted `display regions` response should contain
// every surface name, every region index, every current mode, and every available
// mode from the input configuration.

// propRegionFirstChars contains all characters valid as the first character of a surface name.
const propRegionFirstChars = "abcdefghijklmnopqrstuvwxyz"

// propRegionTailChars contains all characters valid in the tail of a surface name.
const propRegionTailChars = "abcdefghijklmnopqrstuvwxyz0123456789-"

// modePool is a set of mode names that can be assigned to regions in the property test.
// These are fictional mode names that won't collide with each other.
var modePool = []string{
	"alpha", "bravo", "charlie", "delta", "echo",
	"foxtrot", "golf", "hotel", "india", "juliet",
	"kilo", "lima", "mike", "november", "oscar",
}

func TestProperty_RegionListingFormatCompleteness(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate 1-5 regions with unique surface names and random mode lists.
		numRegions := rapid.IntRange(1, 5).Draw(rt, "numRegions")

		type regionInput struct {
			name    string
			index   int
			modes   []string
			current string // the first mode in the list (default)
		}

		regions := make([]regionInput, 0, numRegions)
		coordRegions := make([]coordinator.Region, 0, numRegions)

		for i := 0; i < numRegions; i++ {
			// Generate a unique surface name matching [a-z][a-z0-9-]*
			// Append index suffix to guarantee uniqueness across regions.
			firstIdx := rapid.IntRange(0, len(propRegionFirstChars)-1).Draw(rt, fmt.Sprintf("region%d_firstChar", i))
			firstChar := propRegionFirstChars[firstIdx]

			tailLen := rapid.IntRange(1, 6).Draw(rt, fmt.Sprintf("region%d_tailLen", i))
			tail := make([]byte, tailLen)
			for j := range tail {
				charIdx := rapid.IntRange(0, len(propRegionTailChars)-1).Draw(rt, fmt.Sprintf("region%d_tailChar%d", i, j))
				tail[j] = propRegionTailChars[charIdx]
			}
			// Suffix with region index to ensure uniqueness without retries.
			name := fmt.Sprintf("%c%s%d", firstChar, string(tail), i)

			// Generate 1-4 unique modes for this region from the mode pool.
			numModes := rapid.IntRange(1, 4).Draw(rt, fmt.Sprintf("region%d_numModes", i))
			// Pick a start offset into the mode pool and take sequential modes (wrapping).
			startIdx := rapid.IntRange(0, len(modePool)-1).Draw(rt, fmt.Sprintf("region%d_modeStart", i))
			modes := make([]string, numModes)
			for m := 0; m < numModes; m++ {
				modes[m] = modePool[(startIdx+m)%len(modePool)]
			}

			ri := regionInput{
				name:    name,
				index:   i,
				modes:   modes,
				current: modes[0], // default is first mode
			}
			regions = append(regions, ri)
			coordRegions = append(coordRegions, coordinator.Region{
				Index:   i,
				Name:    name,
				Modes:   modes,
				Default: modes[0],
			})
		}

		// Create coordinator state with the generated regions.
		state := coordinator.NewState(coordRegions...)

		// Create a console server with this state.
		sockPath := filepath.Join(t.TempDir(), "test.sock")
		if runtime.GOOS == "windows" {
			sockPath = filepath.Join(t.TempDir(), "t.sock")
		}
		scanner := source.New(nil, 0)
		gm := gpiomgr.New()
		srv := console.New(sockPath, func() *source.Scanner { return scanner }, gm, nil, state, nil, nil, "")
		if err := srv.Start(); err != nil {
			t.Fatalf("server start: %v", err)
		}
		defer func() {
			srv.Stop()
			_ = os.Remove(sockPath)
		}()

		// Give the server a moment to start.
		time.Sleep(15 * time.Millisecond)

		// Connect and send "display regions".
		conn, err := net.Dial("unix", sockPath)
		if err != nil {
			t.Fatalf("dial: %v", err)
		}
		defer conn.Close()

		sc := bufio.NewScanner(conn)
		sc.Scan() // read greeting

		conn.Write([]byte("display regions\n"))

		// Read the full multi-line response.
		// The response format is:
		//   OK
		//   <name>.0 mode=<current> modes=<m1>,<m2>,...
		//   ...
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		var responseLines []string
		for sc.Scan() {
			line := sc.Text()
			responseLines = append(responseLines, line)
			// The first line is "OK" and subsequent lines are region entries.
			// We know how many regions to expect.
			if len(responseLines) > numRegions {
				break
			}
		}

		if len(responseLines) == 0 {
			t.Fatal("no response received")
		}

		// First line should be "OK"
		if responseLines[0] != "OK" {
			t.Fatalf("expected first line to be 'OK', got %q", responseLines[0])
		}

		// Join all response lines for completeness checks.
		fullResponse := strings.Join(responseLines, "\n")

		// Verify: every surface name appears in the response.
		for _, r := range regions {
			if !strings.Contains(fullResponse, r.name) {
				t.Fatalf("surface name %q not found in response:\n%s", r.name, fullResponse)
			}
		}

		// Verify: every current mode appears in the response.
		for _, r := range regions {
			if !strings.Contains(fullResponse, r.current) {
				t.Fatalf("current mode %q for region %q not found in response:\n%s", r.current, r.name, fullResponse)
			}
		}

		// Verify: every available mode from every region appears in the response.
		for _, r := range regions {
			for _, mode := range r.modes {
				if !strings.Contains(fullResponse, mode) {
					t.Fatalf("available mode %q for region %q not found in response:\n%s", mode, r.name, fullResponse)
				}
			}
		}
	})
}
