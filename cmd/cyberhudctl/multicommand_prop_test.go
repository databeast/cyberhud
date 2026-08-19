package main

import (
	"testing"

	"pgregory.net/rapid"
)

// For any list of command strings where no individual command contains a semicolon
// character, joining them with standalone ';' separators and then splitting by
// semicolons should produce the original list of commands.
func TestProperty_SemicolonSplitRoundTrip(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		numCmds := rapid.IntRange(1, 8).Draw(rt, "numCmds")
		cmds := make([][]string, numCmds)
		for i := 0; i < numCmds; i++ {
			numArgs := rapid.IntRange(1, 5).Draw(rt, "numArgs")
			args := make([]string, numArgs)
			for j := 0; j < numArgs; j++ {
				length := rapid.IntRange(1, 20).Draw(rt, "argLen")
				b := make([]byte, length)
				for k := range b {
					idx := rapid.IntRange(0, 61).Draw(rt, "argChar")
					chars := "abcdefghijklmnopqrstuvwxyz0123456789_-=./ABCDEFGHIJKLMNOPQRSTUVWXYZ"
					b[k] = chars[idx]
				}
				args[j] = string(b)
			}
			cmds[i] = args
		}

		var flat []string
		for i, cmd := range cmds {
			flat = append(flat, cmd...)
			if i < numCmds-1 {
				flat = append(flat, ";")
			}
		}

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
