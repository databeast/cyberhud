package console_test

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/coordinator"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	gpiopkg "github.com/databeast/cyberhud/hardware/gpio"
	"github.com/databeast/cyberhud/runtime/console"
	"pgregory.net/rapid"
)

// For any region_id that does not match any configured region, and any non-empty set of
// available regions, the error response should contain the canonical string representation
// of every available region.

// validSurfaceFirst contains characters valid as the first character of a surface name.
const validSurfaceFirst = "abcdefghijklmnopqrstuvwxyz"

// validSurfaceTail contains characters valid in the tail of a surface name.
const validSurfaceTail = "abcdefghijklmnopqrstuvwxyz0123456789-"

// genSurfaceName generates a valid surface name matching [a-z][a-z0-9-]*.
func genSurfaceName(t *rapid.T, label string) string {
	firstIdx := rapid.IntRange(0, len(validSurfaceFirst)-1).Draw(t, label+"_first")
	first := validSurfaceFirst[firstIdx]
	tailLen := rapid.IntRange(1, 10).Draw(t, label+"_tailLen")
	tail := make([]byte, tailLen)
	for i := range tail {
		charIdx := rapid.IntRange(0, len(validSurfaceTail)-1).Draw(t, label+"_tailChar")
		tail[i] = validSurfaceTail[charIdx]
	}
	return string(first) + string(tail)
}

// genInvalidRegionID generates a region_id string that will NOT match any of the given
// surface names. It avoids bare integers and generates identifiers with an unknown surface.
func genInvalidRegionID(t *rapid.T, surfaceNames []string) string {
	nameSet := make(map[string]bool, len(surfaceNames))
	for _, n := range surfaceNames {
		nameSet[n] = true
	}

	// Generate a surface name that doesn't exist in the available set.
	// We prefix with "zz" to make collisions extremely unlikely with generated names.
	suffix := rapid.IntRange(100, 99999).Draw(t, "invalidSuffix")
	candidate := fmt.Sprintf("zz%d", suffix)

	// Double-check it's not accidentally in the set
	for nameSet[candidate] {
		suffix = rapid.IntRange(100000, 999999).Draw(t, "invalidSuffixRetry")
		candidate = fmt.Sprintf("zz%d", suffix)
	}

	return candidate + ".0"
}

func TestProperty_InvalidRegionError_ListsAvailableRegions(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate 1–5 unique surface names for the configured regions.
		numRegions := rapid.IntRange(1, 5).Draw(rt, "numRegions")
		surfaceNames := make([]string, 0, numRegions)
		namesSeen := make(map[string]bool)

		for i := 0; i < numRegions; i++ {
			var name string
			for {
				name = genSurfaceName(rt, fmt.Sprintf("surface%d", i))
				if !namesSeen[name] {
					break
				}
			}
			namesSeen[name] = true
			surfaceNames = append(surfaceNames, name)
		}

		// Build coordinator regions. Each region gets at least one mode.
		regions := make([]coordinator.Region, numRegions)
		for i, name := range surfaceNames {
			regions[i] = coordinator.Region{
				Index:   i,
				Name:    name,
				Modes:   []string{"dashboard"},
				Default: "dashboard",
			}
		}

		// Create a server with these regions.
		// Use os.MkdirTemp with a short prefix to avoid exceeding Unix socket
		// path length limits on Windows (test name + rapid iteration is too long).
		tmpDir, err := os.MkdirTemp("", "cre")
		if err != nil {
			t.Fatalf("mkdirtemp: %v", err)
		}
		defer os.RemoveAll(tmpDir)
		sockPath := filepath.Join(tmpDir, "t.sock")
		scanner := source.New(nil, 0)
		gm := gpiopkg.New()
		modes := coordinator.NewState(regions...)
		srv := console.New(sockPath, func() *source.Scanner { return scanner }, gm, func() string { return "OK" }, modes, nil, nil, "")
		if err := srv.Start(); err != nil {
			t.Fatalf("server start: %v", err)
		}
		defer func() {
			srv.Stop()
			_ = os.Remove(sockPath)
		}()
		time.Sleep(15 * time.Millisecond)

		// Generate an invalid region ID that doesn't match any configured surface.
		invalidID := genInvalidRegionID(rt, surfaceNames)

		// Connect and send a display set command with the invalid region.
		conn, err := net.Dial("unix", sockPath)
		if err != nil {
			t.Fatalf("dial: %v", err)
		}
		defer conn.Close()

		sc := bufio.NewScanner(conn)
		sc.Scan() // greeting

		cmd := fmt.Sprintf("display set %s dashboard\n", invalidID)
		conn.Write([]byte(cmd))
		sc.Scan()
		resp := sc.Text()

		// Verify the response starts with ERR and contains "available:"
		if !strings.HasPrefix(resp, "ERR") {
			t.Fatalf("expected ERR response for invalid region %q, got: %q", invalidID, resp)
		}

		// Verify the error response contains the canonical string of every available region.
		for _, name := range surfaceNames {
			canonical := name + ".0"
			if !strings.Contains(resp, canonical) {
				t.Fatalf("error response %q does not contain available region %q", resp, canonical)
			}
		}
	})
}
