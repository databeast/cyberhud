package main

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// stubDaemon is a minimal console-protocol server over a Unix socket. It mimics
// the real daemon's greeting + verbose "OK" multi-line responses so we can test
// the real Client.Send read strategy and the full fetch→render pipeline.
func stubDaemon(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	sock := filepath.Join(dir, "stub.sock")
	if runtime.GOOS == "windows" {
		sock = filepath.Join(dir, "s.sock")
	}
	ln, err := net.Listen("unix", sock)
	if err != nil {
		t.Fatalf("listen unix: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })

	responses := map[string]string{
		"display status": "OK\n" +
			`  region=main.0 name="main" controller=st7789 mode=clock modes=clock,dashboard`,
		"display modes": "OK\n" +
			`  region=main.0 name="main" controller=st7789 current=clock` + "\n" +
			`    mode=clock title="Clock" scope=any summary="Current time, date, and weekday."` + "\n" +
			`      option=time_format type=string default="24h" summary="12-hour or 24-hour time display." allowed="24h,12h"` + "\n" +
			`      option=show_border type=bool default="false" summary="Show rounded border frame."` + "\n" +
			`    mode=dashboard title="Dashboard" summary="A dashboard."`,
		"display geometry": "OK\n" +
			"  virtual w=240 h=240\n" +
			`  region=main.0 index=0 name="main" x=0 y=0 w=240 h=240`,
		"display policy main.0": `OK mode=clock {"time_format":"12h","show_border":true}`,
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_, _ = io.WriteString(c, "OK cyberhud daemon ready\n")
				sc := bufio.NewScanner(c)
				for sc.Scan() {
					line := strings.TrimSpace(sc.Text())
					if line == "quit" || line == "exit" {
						_, _ = io.WriteString(c, "OK bye\n")
						return
					}
					if resp, ok := responses[line]; ok {
						_, _ = io.WriteString(c, resp+"\n")
					} else {
						_, _ = io.WriteString(c, "ERR unknown verb\n")
					}
				}
			}(conn)
		}
	}()
	return sock
}

// TestEndToEndFetchAndRender exercises the real socket Client against a stub
// daemon and renders the page through the real Server — the Task 7 smoke test.
func TestEndToEndFetchAndRender(t *testing.T) {
	sock := stubDaemon(t)
	time.Sleep(20 * time.Millisecond)

	client := NewClient(sock, 2*time.Second)

	// Read model assembled from the real multi-line socket responses.
	st, err := client.FetchState()
	if err != nil {
		t.Fatalf("FetchState: %v", err)
	}
	if len(st.Regions) != 1 || st.Regions[0].ID != "main.0" || st.Regions[0].Current != "clock" {
		t.Fatalf("unexpected regions: %+v", st.Regions)
	}
	if len(st.Regions[0].Modes) != 2 || st.Regions[0].Modes[1].Title != "Dashboard" {
		t.Fatalf("modes not enriched: %+v", st.Regions[0].Modes)
	}
	// Option metadata (incl. allowed) parsed from display modes.
	clk := st.Regions[0].Modes[0]
	if len(clk.Options) != 2 || clk.Options[0].Key != "time_format" || !clk.Options[0].isEnum() {
		t.Fatalf("clock options not parsed: %+v", clk.Options)
	}
	// Current policy values fetched from display policy and joined in EditFields
	// on the active mode's ModeRef.
	fields := clk.EditFields()
	if len(fields) != 2 || fields[0].Key != "time_format" || fields[0].Value != "12h" {
		t.Fatalf("edit fields / policy not joined: %+v", fields)
	}
	if st.Geometry.VirtualW != 240 || st.Geometry.VirtualH != 240 {
		t.Fatalf("geometry: %+v", st.Geometry)
	}

	// Full page render through the real HTTP handler.
	srv := NewServer(client)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET / = %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"main", "clock", "Dashboard", "<svg", `class="vd-region"`,
		`class="car-track"`, `action="/config"`, `<option value="12h" selected>`} {
		if !strings.Contains(body, want) {
			t.Errorf("rendered page missing %q", want)
		}
	}
}

// TestEndToEndControlCommand verifies a control verb round-trips over the real
// socket and the handler redirects on success.
func TestEndToEndControlCommand(t *testing.T) {
	sock := stubDaemon(t)
	time.Sleep(20 * time.Millisecond)

	client := NewClient(sock, 2*time.Second)
	// The stub returns ERR for unknown verbs; "display set ..." is unknown to it,
	// so this confirms the ERR path surfaces a notice (the stub has no set verb).
	srv := NewServer(client)

	req := httptest.NewRequest(http.MethodPost, "/region/main.0/next", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	// "display next main.0" is unknown to the stub → ERR → redirect with notice.
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rec.Code)
	}
	if loc := rec.Header().Get("Location"); !strings.HasPrefix(loc, "/?notice=") {
		t.Errorf("expected notice redirect, got %q", loc)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
