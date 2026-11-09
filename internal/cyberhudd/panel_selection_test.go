package cyberhudd

import (
	"context"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPanelSelection(t *testing.T) {
	if os.Getenv("CYBERHUD_TEST_PANEL_SELECTION") == "1" {
		flag.CommandLine = flag.NewFlagSet("cyberhudd", flag.ExitOnError)
		os.Args = []string{"cyberhudd"}
		if args := os.Getenv("CYBERHUD_TEST_ARGS"); args != "" {
			os.Args = append(os.Args, strings.Split(args, "\x1f")...)
		}
		Main()
		return
	}

	dir := t.TempDir()
	config := func(name, contents string) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	missing := config("missing.json", `{"display":{}}`)
	blank := config("blank.json", `{"display":{"panel":"  "}}`)
	valid := config("valid.json", `{"display":{"panel":"waveshare-2.2"}}`)
	invalid := config("invalid.json", `{"display":{"panel":"not-a-panel"}}`)
	headless := config("headless.json", `{"display":{"disabled":true}}`)

	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"missing", nil, "specify a display panel"},
		{"missing config panel", []string{"-config", missing}, "specify a display panel"},
		{"blank config panel", []string{"-config", blank}, "specify a display panel"},
		{"blank flag overrides config", []string{"-config", valid, "-panel", " "}, "specify a display panel"},
		{"headless flag", []string{"-nodisplay"}, "specify a display panel"},
		{"headless config", []string{"-config", headless}, "specify a display panel"},
		{"config panel", []string{"-config", valid, "-display-madctl", "invalid"}, `parsing "invalid"`},
		{"flag panel", []string{"-panel", "waveshare-1.3hat", "-display-madctl", "invalid"}, `parsing "invalid"`},
		{"flag overrides config", []string{"-config", invalid, "-panel", "waveshare-1.3hat", "-display-madctl", "invalid"}, `parsing "invalid"`},
		{"invalid flag overrides config", []string{"-config", valid, "-panel", "not-a-panel"}, `unknown panel "not-a-panel"`},
		{"unknown panel", []string{"-panel", "not-a-panel"}, `unknown panel "not-a-panel"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestPanelSelection$")
			cmd.Env = append(os.Environ(), "CYBERHUD_TEST_PANEL_SELECTION=1", "CYBERHUD_TEST_ARGS="+strings.Join(tc.args, "\x1f"))
			output, err := cmd.CombinedOutput()
			if ctx.Err() != nil || err == nil || !strings.Contains(string(output), tc.want) {
				t.Fatalf("error = %v, output = %q; want failure containing %q", err, output, tc.want)
			}
		})
	}
}
