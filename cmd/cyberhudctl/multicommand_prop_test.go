package main

import (
	"strings"
	"testing"

	"github.com/databeast/cyberhud/display/regionid"
	"pgregory.net/rapid"
)

// For any active region context with ID R and any scoped command (mode M,
// config K=V, next, prev, status), ResolveCommand should produce a fully-qualified
// protocol command containing R's canonical string representation and the original
// arguments.
//
// Specifically:
// - mode X → display set <R> X
// - config k=v → display config <R> k=v
// - next → display next <R>
// - prev → display prev <R>
// - status → display policy <R>

// validSurfaceFirstChars contains all characters valid as the first character of a surface name.
const validSurfaceFirstChars = "abcdefghijklmnopqrstuvwxyz"

// validSurfaceTailChars contains all characters valid in the tail of a surface name.
const validSurfaceTailChars = "abcdefghijklmnopqrstuvwxyz0123456789-"

// genValidSurface generates a random valid surface name matching [a-z][a-z0-9-]*.
func genValidSurface(rt *rapid.T) string {
	firstIdx := rapid.IntRange(0, len(validSurfaceFirstChars)-1).Draw(rt, "firstIdx")
	firstChar := validSurfaceFirstChars[firstIdx]
	tailLen := rapid.IntRange(0, 12).Draw(rt, "tailLen")
	tail := make([]byte, tailLen)
	for i := range tail {
		charIdx := rapid.IntRange(0, len(validSurfaceTailChars)-1).Draw(rt, "tailCharIdx")
		tail[i] = validSurfaceTailChars[charIdx]
	}
	return string(firstChar) + string(tail)
}

// genRegionID generates a random valid region ID.
func genRegionID(rt *rapid.T) regionid.ID {
	surface := genValidSurface(rt)
	index := rapid.IntRange(0, 100).Draw(rt, "index")
	return regionid.ID{Surface: surface, Index: index}
}

// identChars are characters valid in mode names and config key/value strings.
const identChars = "abcdefghijklmnopqrstuvwxyz0123456789_"

// genIdent generates a random identifier (mode name, config key, or config value).
func genIdent(rt *rapid.T, label string) string {
	length := rapid.IntRange(1, 16).Draw(rt, label+"Len")
	b := make([]byte, length)
	for i := range b {
		idx := rapid.IntRange(0, len(identChars)-1).Draw(rt, label+"Char")
		b[i] = identChars[idx]
	}
	return string(b)
}

func TestProperty_ScopedCommandExpansion(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a random region ID as context.
		rid := genRegionID(rt)
		regionStr := rid.String()
		ctx := &RegionContext{Active: &rid}

		// Choose a random scoped command type.
		cmdType := rapid.IntRange(0, 4).Draw(rt, "cmdType")

		switch cmdType {
		case 0: // mode X → display set <R> X
			modeName := genIdent(rt, "mode")
			result, err := ctx.ResolveCommand([]string{"mode", modeName})
			if err != nil {
				t.Fatalf("ResolveCommand([mode %s]) error: %v", modeName, err)
			}
			expected := "display set " + regionStr + " " + modeName
			if result != expected {
				t.Fatalf("mode expansion: got %q, want %q", result, expected)
			}

		case 1: // config k=v → display config <R> k=v
			numKV := rapid.IntRange(1, 5).Draw(rt, "numKV")
			args := []string{"config"}
			kvParts := make([]string, numKV)
			for i := 0; i < numKV; i++ {
				key := genIdent(rt, "cfgKey")
				val := genIdent(rt, "cfgVal")
				kv := key + "=" + val
				args = append(args, kv)
				kvParts[i] = kv
			}
			result, err := ctx.ResolveCommand(args)
			if err != nil {
				t.Fatalf("ResolveCommand(%v) error: %v", args, err)
			}
			expected := "display config " + regionStr + " " + strings.Join(kvParts, " ")
			if result != expected {
				t.Fatalf("config expansion: got %q, want %q", result, expected)
			}

		case 2: // next → display next <R>
			result, err := ctx.ResolveCommand([]string{"next"})
			if err != nil {
				t.Fatalf("ResolveCommand([next]) error: %v", err)
			}
			expected := "display next " + regionStr
			if result != expected {
				t.Fatalf("next expansion: got %q, want %q", result, expected)
			}

		case 3: // prev → display prev <R>
			result, err := ctx.ResolveCommand([]string{"prev"})
			if err != nil {
				t.Fatalf("ResolveCommand([prev]) error: %v", err)
			}
			expected := "display prev " + regionStr
			if result != expected {
				t.Fatalf("prev expansion: got %q, want %q", result, expected)
			}

		case 4: // status → display policy <R>
			result, err := ctx.ResolveCommand([]string{"status"})
			if err != nil {
				t.Fatalf("ResolveCommand([status]) error: %v", err)
			}
			expected := "display policy " + regionStr
			if result != expected {
				t.Fatalf("status expansion: got %q, want %q", result, expected)
			}
		}
	})
}

// For any list of command strings where no individual command contains a semicolon
// character, joining them with standalone ';' separators and then splitting by
// semicolons should produce the original list of commands.
//
// SplitCommands(Join(cmds, ";")) == cmds

// argChars are characters valid in command arguments (excludes semicolons).
const argChars = "abcdefghijklmnopqrstuvwxyz0123456789_-=./ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// genArg generates a random non-empty argument string without semicolons.
func genArg(rt *rapid.T, label string) string {
	length := rapid.IntRange(1, 20).Draw(rt, label+"Len")
	b := make([]byte, length)
	for i := range b {
		idx := rapid.IntRange(0, len(argChars)-1).Draw(rt, label+"Char")
		b[i] = argChars[idx]
	}
	return string(b)
}

func TestProperty_SemicolonSplitRoundTrip(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a random list of commands (each is a list of non-semicolon args).
		numCmds := rapid.IntRange(1, 8).Draw(rt, "numCmds")
		cmds := make([][]string, numCmds)
		for i := 0; i < numCmds; i++ {
			numArgs := rapid.IntRange(1, 5).Draw(rt, "numArgs")
			args := make([]string, numArgs)
			for j := 0; j < numArgs; j++ {
				args[j] = genArg(rt, "arg")
			}
			cmds[i] = args
		}

		// Join with standalone ";" separators into a flat arg list.
		var flat []string
		for i, cmd := range cmds {
			flat = append(flat, cmd...)
			if i < numCmds-1 {
				flat = append(flat, ";")
			}
		}

		// Split by semicolons and verify round-trip.
		got := SplitCommands(flat)
		if len(got) != len(cmds) {
			t.Fatalf("round-trip length mismatch: got %d commands, want %d\nflat: %v\ngot: %v\nwant: %v",
				len(got), len(cmds), flat, got, cmds)
		}
		for i := range cmds {
			if len(got[i]) != len(cmds[i]) {
				t.Fatalf("command %d length mismatch: got %v, want %v", i, got[i], cmds[i])
			}
			for j := range cmds[i] {
				if got[i][j] != cmds[i][j] {
					t.Fatalf("command %d arg %d mismatch: got %q, want %q", i, j, got[i][j], cmds[i][j])
				}
			}
		}
	})
}

// For any command sequence where command at index K is `region <R>` and all commands
// at indices > K are scoped commands (mode, config, next, prev, status), every
// resolved command after index K should target region R.

func TestProperty_RegionContextPersistence(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a random region ID.
		rid := genRegionID(rt)
		regionStr := rid.String()

		// Build a command sequence: start with `region <R>`, followed by random scoped commands.
		ctx := &RegionContext{}

		// First command: set the region context.
		_, err := ctx.ResolveCommand([]string{"region", regionStr})
		if err != nil {
			t.Fatalf("region command error: %v", err)
		}

		// Generate N subsequent scoped commands, all should target R.
		numCmds := rapid.IntRange(1, 10).Draw(rt, "numCmds")
		for i := 0; i < numCmds; i++ {
			cmdType := rapid.IntRange(0, 4).Draw(rt, "cmdType")
			var args []string
			switch cmdType {
			case 0:
				modeName := genIdent(rt, "mode")
				args = []string{"mode", modeName}
			case 1:
				key := genIdent(rt, "cfgKey")
				val := genIdent(rt, "cfgVal")
				args = []string{"config", key + "=" + val}
			case 2:
				args = []string{"next"}
			case 3:
				args = []string{"prev"}
			case 4:
				args = []string{"status"}
			}

			result, err := ctx.ResolveCommand(args)
			if err != nil {
				t.Fatalf("command %d (%v) error: %v", i, args, err)
			}

			// Every resolved command should contain the region string.
			if !strings.Contains(result, regionStr) {
				t.Fatalf("command %d (%v) resolved to %q which does not contain region %q",
					i, args, result, regionStr)
			}
		}
	})
}
