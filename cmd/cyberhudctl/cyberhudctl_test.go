package main

import (
	"bufio"
	"bytes"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// --- From: freeze_test.go ---
// (freeze-related tests removed: freeze is now a daemon-side protocol command)

// --- From: main_test.go ---

func TestProtocolCommand(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "status", args: []string{"status"}, want: "status"},
		{name: "stemma status", args: []string{"stemma", "status"}, want: "stemma status"},
		{name: "gpio status", args: []string{"gpio", "status"}, want: "gpio status"},
		{name: "gpio pins", args: []string{"gpio", "pins"}, want: "gpio pins"},
		{name: "gpio in", args: []string{"gpio", "in", "17"}, want: "gpio in 17"},
		{name: "gpio set", args: []string{"gpio", "set", "4", "1"}, want: "gpio set 4 1"},
		{name: "display status", args: []string{"display", "status"}, want: "display status"},
		{name: "display list", args: []string{"display", "list"}, want: "display list"},
		{name: "display modes", args: []string{"display", "modes"}, want: "display modes"},
		{name: "display set", args: []string{"display", "set", "1", "GPIO"}, want: "display set 1 gpio"},
		{name: "display next", args: []string{"display", "next", "0"}, want: "display next 0"},
		{name: "display prev", args: []string{"display", "prev", "2"}, want: "display prev 2"},
		{name: "ticker set", args: []string{"display", "ticker", "set", "hello", "world"}, want: "display ticker set hello world"},
		{name: "ticker policy query", args: []string{"display", "ticker", "policy"}, want: "display ticker policy"},
		{name: "ticker policy set", args: []string{"display", "ticker", "policy", "line_mode=clip", "direction=none", "auto_scroll_ms=0"}, want: "display ticker policy line_mode=clip direction=none auto_scroll_ms=0"},
		{name: "image clear", args: []string{"display", "image", "clear"}, want: "display image clear"},
		{name: "image policy query", args: []string{"display", "image", "policy"}, want: "display image policy"},
		{name: "image policy set", args: []string{"display", "image", "policy", "fit=truncate"}, want: "display image policy fit=truncate"},
		{name: "mode sub-command bad key accepted", args: []string{"display", "ticker", "policy", "foo=bar"}, want: "display ticker policy foo=bar"},
		{name: "unknown verb rejected", args: []string{"weather", "set", "city=Berlin"}, wantErr: true},
		{name: "help passthrough", args: []string{"help", "modes"}, want: "help modes"},
		{name: "display regions", args: []string{"display", "regions"}, want: "display regions"},
		{name: "display policy by mode", args: []string{"display", "policy", "clock"}, want: "display policy clock"},
		{name: "display policy by region", args: []string{"display", "policy", "main.0"}, want: "display policy main.0"},
		{name: "freeze as protocol command", args: []string{"freeze"}, want: "freeze"},
		{name: "freeze policy", args: []string{"freeze", "policy"}, want: "freeze policy"},
		{name: "policy dump", args: []string{"policy", "dump"}, want: "policy dump"},
		{name: "region passthrough", args: []string{"region", "main.0"}, want: "region main.0"},
		{name: "bad level", args: []string{"gpio", "set", "4", "9"}, wantErr: true},
		{name: "missing", args: nil, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := protocolCommand(tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRun_StatusOK(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		if line == "status" {
			return []string{"OK stemma_devices=1 gpio_pins=3"}
		}
		if line == "pins" {
			return []string{"OK pin report", "  display_region=test"}
		}
		return []string{"ERR unknown"}
	})
	defer stop()

	var out, errOut bytes.Buffer
	exit := run([]string{"-socket", sock, "status"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run() exit = %d, errOut=%q", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "OK stemma_devices=1") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestRun_PinsOK(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		if line == "gpio pins" {
			return []string{"OK pin report", "  display_region=waveshare-1.3hat"}
		}
		return []string{"ERR unknown"}
	})
	defer stop()

	var out, errOut bytes.Buffer
	exit := run([]string{"-socket", sock, "gpio", "pins"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run() exit = %d, errOut=%q", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "OK pin report") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestRun_CommandError(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		if line == "gpio set 4 1" {
			return []string{"ERR GPIO4 set failed"}
		}
		return []string{"ERR unknown"}
	})
	defer stop()

	var out, errOut bytes.Buffer
	exit := run([]string{"-socket", sock, "gpio", "set", "4", "1"}, &out, &errOut)
	if exit == 0 {
		t.Fatalf("expected non-zero exit, got 0")
	}
	if !strings.Contains(out.String(), "ERR GPIO4 set failed") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestRun_HelpModesPretty(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		if line == "help modes" {
			return []string{
				"OK",
				"  verb=image usage=\"image <set|policy|clear> ...\" summary=\"Manage externally supplied images and fit policy.\" mode_title=\"Image\" scope=any",
				"    option=fit type=string default=\"scale\" summary=\"How the image should fit the panel bounds.\"",
				"  verb=ticker usage=\"ticker <set|get|policy> ...\" summary=\"Manage externally supplied ticker text and policy.\" mode_title=\"Ticker\" scope=any",
				"    option=line_mode type=string default=\"truncate\" summary=\"How long lines are constrained to panel width.\"",
			}
		}
		return []string{"ERR unknown"}
	})
	defer stop()

	var out, errOut bytes.Buffer
	exit := run([]string{"-socket", sock, "help", "modes"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run() exit = %d, errOut=%q out=%q", exit, errOut.String(), out.String())
	}
	got := out.String()
	if !strings.Contains(got, "Mode Commands:") {
		t.Fatalf("expected pretty header, got %q", got)
	}
	if !strings.Contains(got, "- image") || !strings.Contains(got, "- ticker") {
		t.Fatalf("expected image and ticker entries, got %q", got)
	}
	if strings.Contains(got, "verb=image") {
		t.Fatalf("expected pretty output, got raw key-value line: %q", got)
	}
}

func TestFormatHelpModesResponseFallback(t *testing.T) {
	if got := formatHelpModesResponse("OK stemma_devices=1"); got != "" {
		t.Fatalf("formatHelpModesResponse fallback expected empty, got %q", got)
	}
}

func TestParseKeyValueLineQuoted(t *testing.T) {
	kv := parseKeyValueLine(`verb=image usage="image <set|policy|clear> ..." scope=any`)
	if kv["verb"] != "image" {
		t.Fatalf("verb=%q", kv["verb"])
	}
	if kv["usage"] != "image <set|policy|clear> ..." {
		t.Fatalf("usage=%q", kv["usage"])
	}
	if kv["scope"] != "any" {
		t.Fatalf("scope=%q", kv["scope"])
	}
}

func startFakeServer(t *testing.T, sock string, dispatch func(string) []string) func() {
	t.Helper()
	_ = os.Remove(sock)
	ln, err := net.Listen("unix", sock)
	if err != nil {
		t.Fatalf("listen unix: %v", err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_, _ = c.Write([]byte("OK cyberhud daemon ready\n"))
				sc := bufio.NewScanner(c)
				for sc.Scan() {
					line := strings.TrimSpace(sc.Text())
					if line == "" {
						continue
					}
					for _, resp := range dispatch(line) {
						_, _ = c.Write([]byte(resp + "\n"))
					}
					if line == "quit" || line == "exit" {
						_, _ = c.Write([]byte("OK bye\n"))
						return
					}
				}
			}(conn)
		}
	}()
	return func() {
		_ = ln.Close()
		<-done
	}
}

func testSocketPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	name := "console.sock"
	if runtime.GOOS == "windows" {
		// unix sockets in tests require a short path on Windows.
		name = "c.sock"
	}
	return filepath.Join(dir, name)
}

func TestRun_StyleSetWithFitnessNotes(t *testing.T) {
	// Simulate a daemon response that includes fitness note lines after a style-set.
	// The notes should be printed to stdout and NOT cause a non-zero exit code.
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		if strings.HasPrefix(line, "display config") {
			return []string{
				"OK thermal style=graph font=auto refresh_ms=2000 warn_threshold=70 crit_threshold=90 show_border=false unit=C color_accent=thermal show_led=false show_refresh_bar=false",
				"note: style requires minimum width 200 but panel is 128",
				"note: style prefers 240×240 but panel is 128×64",
			}
		}
		return []string{"ERR unknown"}
	})
	defer stop()

	var out, errOut bytes.Buffer
	exit := run([]string{"-socket", sock, "display", "config", "0", "style=graph"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run() exit = %d; notes should not cause non-zero exit; errOut=%q", exit, errOut.String())
	}

	got := out.String()
	if !strings.Contains(got, "OK thermal") {
		t.Fatalf("expected OK response in output, got %q", got)
	}
	if !strings.Contains(got, "note: style requires minimum width 200 but panel is 128") {
		t.Fatalf("expected first fitness note in output, got %q", got)
	}
	if !strings.Contains(got, "note: style prefers 240×240 but panel is 128×64") {
		t.Fatalf("expected second fitness note in output, got %q", got)
	}
}

func TestRun_StyleSetNoNotes(t *testing.T) {
	// When fitness is Full/Optimal, no note lines are included in the response.
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		if strings.HasPrefix(line, "display config") {
			return []string{
				"OK clock style=digital show_seconds=true time_format=24h",
			}
		}
		return []string{"ERR unknown"}
	})
	defer stop()

	var out, errOut bytes.Buffer
	exit := run([]string{"-socket", sock, "display", "config", "0", "style=digital"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run() exit = %d, errOut=%q", exit, errOut.String())
	}

	got := out.String()
	if !strings.Contains(got, "OK clock") {
		t.Fatalf("expected OK response, got %q", got)
	}
	if strings.Contains(got, "note:") {
		t.Fatalf("expected no fitness notes when Full/Optimal, got %q", got)
	}
}

func TestSendMultipleCommands_AllSucceed(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		switch line {
		case "display set main.0 attract_matrix":
			return []string{"OK region=main.0 mode=attract_matrix"}
		case "display config main.0 density=0.8":
			return []string{"OK attract_matrix density=0.8"}
		case "quit":
			return nil
		default:
			return []string{"ERR unknown command"}
		}
	})
	defer stop()

	cmds := []string{
		"display set main.0 attract_matrix",
		"display config main.0 density=0.8",
	}
	responses, err := sendMultipleCommands(sock, cmds, 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if !strings.Contains(responses[0], "OK region=main.0 mode=attract_matrix") {
		t.Fatalf("unexpected first response: %q", responses[0])
	}
	if !strings.Contains(responses[1], "OK attract_matrix density=0.8") {
		t.Fatalf("unexpected second response: %q", responses[1])
	}
}

func TestSendMultipleCommands_StopsOnError(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		switch line {
		case "display set main.0 attract_matrix":
			return []string{"OK region=main.0 mode=attract_matrix"}
		case "display set main.0 nonexistent":
			return []string{"ERR unknown mode \"nonexistent\""}
		case "display config main.0 density=0.8":
			return []string{"OK attract_matrix density=0.8"}
		case "quit":
			return nil
		default:
			return []string{"ERR unknown command"}
		}
	})
	defer stop()

	cmds := []string{
		"display set main.0 attract_matrix",
		"display set main.0 nonexistent",
		"display config main.0 density=0.8", // should not execute
	}
	responses, err := sendMultipleCommands(sock, cmds, 2*time.Second)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	// Should have 2 responses: the first success and the error
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses (success + error), got %d: %v", len(responses), responses)
	}
	if !strings.Contains(responses[0], "OK") {
		t.Fatalf("expected first response to be OK, got %q", responses[0])
	}
	if !strings.Contains(responses[1], "ERR") {
		t.Fatalf("expected second response to be ERR, got %q", responses[1])
	}
}

func TestSendMultipleCommands_EmptyCmds(t *testing.T) {
	_, err := sendMultipleCommands("/nonexistent.sock", nil, 2*time.Second)
	if err == nil {
		t.Fatalf("expected error for empty commands, got nil")
	}
	if !strings.Contains(err.Error(), "no commands") {
		t.Fatalf("expected 'no commands' error, got %q", err.Error())
	}
}

func TestSendMultipleCommands_MultiLineResponse(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		switch line {
		case "display regions":
			// Multi-line response ending with OK
			return []string{
				"surface=main index=0 mode=attract_matrix modes=attract_matrix,clock,ticker",
				"surface=left-aux index=0 mode=clock modes=clock,ticker",
				"OK 2 regions",
			}
		case "quit":
			return nil
		default:
			return []string{"ERR unknown command"}
		}
	})
	defer stop()

	cmds := []string{"display regions"}
	responses, err := sendMultipleCommands(sock, cmds, 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(responses))
	}
	if !strings.Contains(responses[0], "surface=main") {
		t.Fatalf("expected multi-line response with surface data, got %q", responses[0])
	}
	if !strings.Contains(responses[0], "OK 2 regions") {
		t.Fatalf("expected OK terminator in response, got %q", responses[0])
	}
}

func TestRun_DisplayRegions(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		if line == "display regions" {
			return []string{
				"OK",
				`  region=main.0 name="main" controller=st7789 mode=clock modes=clock,dashboard`,
				`  region=left-aux.0 name="left-aux" controller=st7735s mode=stemma modes=stemma,gpio`,
			}
		}
		return []string{"ERR unknown"}
	})
	defer stop()

	var out, errOut bytes.Buffer
	exit := run([]string{"-socket", sock, "display", "regions"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run() exit = %d, errOut=%q", exit, errOut.String())
	}
	got := out.String()
	if !strings.Contains(got, "main.0") {
		t.Fatalf("expected region main.0 in output, got %q", got)
	}
	if !strings.Contains(got, "left-aux.0") {
		t.Fatalf("expected region left-aux.0 in output, got %q", got)
	}
}

func TestRun_PolicyDump(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		if line == "policy dump" {
			return []string{`OK {"attract_matrix":{"density":0.8},"clock":{}}`}
		}
		return []string{"ERR unknown"}
	})
	defer stop()

	var out, errOut bytes.Buffer
	exit := run([]string{"-socket", sock, "raw", "policy", "dump"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run() exit = %d, errOut=%q", exit, errOut.String())
	}
	got := out.String()
	if !strings.Contains(got, "attract_matrix") {
		t.Fatalf("expected attract_matrix in policy dump output, got %q", got)
	}
}

func TestRun_FreezePolicyProtocol(t *testing.T) {
	sock := testSocketPath(t)
	stop := startFakeServer(t, sock, func(line string) []string {
		if line == "freeze policy" {
			return []string{"OK"}
		}
		return []string{"ERR unknown"}
	})
	defer stop()

	var out, errOut bytes.Buffer
	exit := run([]string{"-socket", sock, "raw", "freeze", "policy"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("run() exit = %d, errOut=%q", exit, errOut.String())
	}
	got := out.String()
	if !strings.Contains(got, "OK") {
		t.Fatalf("expected OK in freeze policy output, got %q", got)
	}
}
