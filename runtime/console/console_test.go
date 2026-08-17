package console_test

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/databeast/cyberhud/display/coordinator"
	"github.com/databeast/cyberhud/display/modes/stemma/source"
	tickermode "github.com/databeast/cyberhud/display/modes/ticker/source"
	gpiomgr "github.com/databeast/cyberhud/hardware/gpio"
	"github.com/databeast/cyberhud/runtime/console"

	// Blank imports to trigger mode init() self-registration.
	_ "github.com/databeast/cyberhud/display/modes/clock"
	_ "github.com/databeast/cyberhud/display/modes/dashboard"
	_ "github.com/databeast/cyberhud/display/modes/gpio"
	_ "github.com/databeast/cyberhud/display/modes/gpio_control"
	_ "github.com/databeast/cyberhud/display/modes/image"
	_ "github.com/databeast/cyberhud/display/modes/menu"
	_ "github.com/databeast/cyberhud/display/modes/ticker"
)

func newTestServer(t *testing.T) (*console.Server, string) {
	t.Helper()
	sockPath := filepath.Join(t.TempDir(), "cyberhudd_test.sock")
	if runtime.GOOS == "windows" {
		sockPath = filepath.Join(t.TempDir(), "c.sock")
	}
	scanner := source.New(nil, 0)
	gm := gpiomgr.New()
	modes := coordinator.NewState(
		coordinator.Region{Index: 0, Name: "main", Controller: "st7789", Modes: []string{"menu", "dashboard"}, Default: "menu"},
		coordinator.Region{Index: 1, Name: "aux", Controller: "st7735s", Modes: []string{"stemma", "gpio"}, Default: "stemma"},
	)
	srv := console.New(sockPath, func() *source.Scanner { return scanner }, gm, func() string { return "OK pin report\n  display_region=test" }, modes, nil, nil, "")
	if err := srv.Start(); err != nil {
		t.Fatalf("server start: %v", err)
	}
	t.Cleanup(func() {
		srv.Stop()
		_ = os.Remove(sockPath)
	})
	return srv, sockPath
}

func TestConsoleServer_Status(t *testing.T) {
	_, sockPath := newTestServer(t)

	// Give the server a moment to start.
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)

	// Read greeting.
	if !sc.Scan() {
		t.Fatal("expected greeting")
	}
	greeting := sc.Text()
	if !strings.HasPrefix(greeting, "OK ") {
		t.Fatalf("expected OK greeting, got %q", greeting)
	}

	// Send status command.
	if _, err := conn.Write([]byte("status\n")); err != nil {
		t.Fatalf("write: %v", err)
	}
	if !sc.Scan() {
		t.Fatal("expected status response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK ") {
		t.Fatalf("expected OK response, got %q", resp)
	}
}

func TestConsoleServer_ListStemma(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("stemma status\n"))
	sc.Scan()
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK") {
		t.Fatalf("expected OK, got %q", resp)
	}
}

func TestConsoleServer_ListGPIO(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("gpio status\n"))
	sc.Scan()
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK") {
		t.Fatalf("expected OK, got %q", resp)
	}
}

func TestConsoleServer_UnknownVerb(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("flibble\n"))
	sc.Scan()
	resp := sc.Text()
	if !strings.HasPrefix(resp, "ERR") {
		t.Fatalf("expected ERR, got %q", resp)
	}
}

func TestConsoleServer_Pins(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("gpio pins\n"))
	sc.Scan()
	resp := sc.Text()
	if !strings.Contains(resp, "OK pin report") {
		t.Fatalf("expected pin report, got %q", resp)
	}
}

func TestConsoleServer_GPIOSetUnknownPin(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("gpio set 999 1\n"))
	sc.Scan()
	resp := sc.Text()
	if !strings.HasPrefix(resp, "ERR") {
		t.Fatalf("expected ERR for unknown pin, got %q", resp)
	}
}

func TestSendCommand(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	resp, err := console.SendCommand(sockPath, "status")
	if err != nil {
		t.Fatalf("SendCommand: %v", err)
	}
	if !strings.HasPrefix(resp, "OK") {
		t.Fatalf("expected OK response, got %q", resp)
	}
}

func TestConsoleServer_DisplayStatusAndSet(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display status\n"))
	sc.Scan()
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK") {
		t.Fatalf("expected OK, got %q", resp)
	}
	if !sc.Scan() {
		t.Fatal("expected first region status line")
	}
	if got := sc.Text(); !strings.Contains(got, `region=main.0 name="main" controller=st7789 mode=menu modes=menu,dashboard`) {
		t.Fatalf("unexpected region status line: %q", got)
	}
	_ = conn.Close()

	conn, err = net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	sc = bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display set 1 gpio\n"))
	sc.Scan()
	resp = sc.Text()
	if !strings.Contains(resp, "mode=gpio") {
		t.Fatalf("expected mode switch response, got %q", resp)
	}
}

func TestConsoleServer_DisplayModes(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display modes\n"))
	if !sc.Scan() {
		t.Fatal("expected OK line")
	}
	if got := sc.Text(); got != "OK" {
		t.Fatalf("unexpected first line: %q", got)
	}
	if !sc.Scan() {
		t.Fatal("expected region header line")
	}
	if got := sc.Text(); !strings.Contains(got, `region=main.0 name="main" controller=st7789 current=menu`) {
		t.Fatalf("unexpected region header: %q", got)
	}
	if !sc.Scan() {
		t.Fatal("expected mode description line")
	}
	if got := sc.Text(); !strings.Contains(got, `mode=menu title="Menu"`) || strings.Contains(got, `scope=`) {
		t.Fatalf("unexpected mode detail: %q", got)
	}
}

func TestConsoleServer_DisplayNextPrev(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display next 1\n"))
	sc.Scan()
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK region=aux.0 mode=gpio") {
		t.Fatalf("unexpected next response: %q", resp)
	}

	conn.Write([]byte("display prev 1\n"))
	sc.Scan()
	resp = sc.Text()
	if !strings.HasPrefix(resp, "OK region=aux.0 mode=stemma") {
		t.Fatalf("unexpected prev response: %q", resp)
	}
}

func TestConsoleServer_DisplaySetUnknownRegion(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display set 9 gpio\n"))
	sc.Scan()
	resp := sc.Text()
	if !strings.HasPrefix(resp, `ERR unknown region "9"`) {
		t.Fatalf("unexpected response: %q", resp)
	}
}

func TestConsoleServer_DisplaySetUnsupportedMode(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display set 1 nope\n"))
	sc.Scan()
	resp := sc.Text()
	if !strings.HasPrefix(resp, `ERR unsupported mode "nope" for region`) {
		t.Fatalf("unexpected response: %q", resp)
	}
}

func TestConsoleServer_TickerSet(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display ticker set alpha|beta|gamma\n"))
	if !sc.Scan() {
		t.Fatal("expected ticker response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK ticker lines=3") {
		t.Fatalf("unexpected ticker response: %q", resp)
	}
	got := tickermode.Snapshot()
	if len(got) != 3 || got[0] != "alpha" || got[1] != "beta" || got[2] != "gamma" {
		t.Fatalf("unexpected ticker snapshot: %v", got)
	}
}

func TestConsoleServer_TickerPolicy(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display ticker policy\n"))
	if !sc.Scan() {
		t.Fatal("expected ticker policy response")
	}
	if got := sc.Text(); !strings.HasPrefix(got, "OK ticker policy") {
		t.Fatalf("unexpected ticker policy response: %q", got)
	}

	conn.Write([]byte("display ticker policy line_mode=clip direction=none auto_scroll_ms=1500\n"))
	if !sc.Scan() {
		t.Fatal("expected ticker policy set response")
	}
	resp := sc.Text()
	if !strings.Contains(resp, "line_mode=clip") || !strings.Contains(resp, "direction=none") || !strings.Contains(resp, "auto_scroll_ms=1500") {
		t.Fatalf("unexpected ticker policy set response: %q", resp)
	}
}

func TestConsoleServer_ImageClear(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display image clear\n"))
	if !sc.Scan() {
		t.Fatal("expected image clear response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK image cleared") {
		t.Fatalf("unexpected image clear response: %q", resp)
	}
}

func TestConsoleServer_ImagePolicy(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display image policy\n"))
	if !sc.Scan() {
		t.Fatal("expected image policy response")
	}
	if got := sc.Text(); !strings.HasPrefix(got, "OK image policy") {
		t.Fatalf("unexpected image policy response: %q", got)
	}

	conn.Write([]byte("display image policy fit=truncate\n"))
	if !sc.Scan() {
		t.Fatal("expected image policy set response")
	}
	resp := sc.Text()
	if !strings.Contains(resp, "fit=truncate") {
		t.Fatalf("unexpected image policy set response: %q", resp)
	}
}

func TestConsoleServer_HelpModes(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("help modes\n"))
	_ = conn.SetReadDeadline(time.Now().Add(150 * time.Millisecond))
	if !sc.Scan() {
		t.Fatal("expected help modes response")
	}
	if got := sc.Text(); got != "OK" {
		t.Fatalf("unexpected first help line: %q", got)
	}
	var body strings.Builder
	for sc.Scan() {
		body.WriteString(sc.Text())
		body.WriteByte('\n')
	}
	out := body.String()
	if !strings.Contains(out, "verb=ticker") || !strings.Contains(out, "verb=image") {
		t.Fatalf("help modes missing expected commands (ticker=%v image=%v)\n%s",
			strings.Contains(out, "verb=ticker"), strings.Contains(out, "verb=image"), out)
	}
}

func newTestServerWithConfig(t *testing.T, configFn func() string) (*console.Server, string) {
	t.Helper()
	sockPath := filepath.Join(t.TempDir(), "cyberhudd_test.sock")
	if runtime.GOOS == "windows" {
		sockPath = filepath.Join(t.TempDir(), "c.sock")
	}
	scanner := source.New(nil, 0)
	gm := gpiomgr.New()
	modes := coordinator.NewState(
		coordinator.Region{Index: 0, Name: "main", Controller: "st7789", Modes: []string{"menu"}, Default: "menu"},
	)
	srv := console.New(sockPath, func() *source.Scanner { return scanner }, gm, func() string { return "OK" }, modes, configFn, nil, "")
	if err := srv.Start(); err != nil {
		t.Fatalf("server start: %v", err)
	}
	t.Cleanup(func() {
		srv.Stop()
		_ = os.Remove(sockPath)
	})
	return srv, sockPath
}

func TestConsoleServer_ConfigDumpNilClosure(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("config dump\n"))
	if !sc.Scan() {
		t.Fatal("expected config dump response")
	}
	resp := sc.Text()
	if resp != "ERR config snapshot unavailable" {
		t.Fatalf("expected ERR config snapshot unavailable, got %q", resp)
	}
}

func TestConsoleServer_ConfigDumpSuccess(t *testing.T) {
	wantJSON := `{"socket":"/run/test.sock","scan":"2s"}`
	_, sockPath := newTestServerWithConfig(t, func() string { return wantJSON })
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("config dump\n"))

	// First line should be "OK"
	if !sc.Scan() {
		t.Fatal("expected OK line")
	}
	if got := sc.Text(); got != "OK" {
		t.Fatalf("expected first line to be %q, got %q", "OK", got)
	}

	// Second line should contain the JSON
	if !sc.Scan() {
		t.Fatal("expected JSON line")
	}
	got := sc.Text()
	if got != wantJSON {
		t.Fatalf("expected JSON %q, got %q", wantJSON, got)
	}

	// Validate the JSON is parseable
	if !json.Valid([]byte(got)) {
		t.Fatalf("response is not valid JSON: %q", got)
	}
}

func TestConsoleServer_ConfigDumpUsage(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	// Send "config" with no subcommand
	conn.Write([]byte("config\n"))
	if !sc.Scan() {
		t.Fatal("expected config usage response")
	}
	resp := sc.Text()
	if resp != "ERR usage: config dump" {
		t.Fatalf("expected usage error, got %q", resp)
	}
}

func TestConsoleServer_ConfigDumpDoesNotMutateState(t *testing.T) {
	callCount := 0
	configFn := func() string {
		callCount++
		return `{"count":1}`
	}
	_, sockPath := newTestServerWithConfig(t, configFn)
	time.Sleep(20 * time.Millisecond)

	// Send config dump twice and verify server state is unchanged
	for i := 0; i < 2; i++ {
		conn, err := net.Dial("unix", sockPath)
		if err != nil {
			t.Fatalf("dial: %v", err)
		}
		sc := bufio.NewScanner(conn)
		sc.Scan() // greeting

		conn.Write([]byte("config dump\n"))
		sc.Scan() // OK line
		sc.Scan() // JSON line
		conn.Close()
	}

	// The closure was called exactly twice (once per request), no side effects
	if callCount != 2 {
		t.Fatalf("expected configFn to be called 2 times, got %d", callCount)
	}

	// Verify other commands still work normally after config dump
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("status\n"))
	if !sc.Scan() {
		t.Fatal("expected status response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK") {
		t.Fatalf("expected OK status after config dump, got %q", resp)
	}
}

// mockPolicyStore implements console.PolicyStoreAccess for testing.
type mockPolicyStore struct {
	policies map[string]json.RawMessage
}

func (m *mockPolicyStore) Get(modeID string) json.RawMessage {
	return m.policies[modeID]
}

func (m *mockPolicyStore) Set(modeID string, data json.RawMessage) {
	if m.policies == nil {
		m.policies = make(map[string]json.RawMessage)
	}
	m.policies[modeID] = data
}

func (m *mockPolicyStore) All() map[string]json.RawMessage {
	cp := make(map[string]json.RawMessage, len(m.policies))
	for k, v := range m.policies {
		cp[k] = v
	}
	return cp
}

func newTestServerWithPolicyStore(t *testing.T, ps console.PolicyStoreAccess) (*console.Server, string) {
	t.Helper()
	sockPath := filepath.Join(t.TempDir(), "cyberhudd_test.sock")
	if runtime.GOOS == "windows" {
		sockPath = filepath.Join(t.TempDir(), "c.sock")
	}
	scanner := source.New(nil, 0)
	gm := gpiomgr.New()
	modes := coordinator.NewState(
		coordinator.Region{Index: 0, Name: "main", Controller: "st7789", Modes: []string{"menu"}, Default: "menu"},
	)
	srv := console.New(sockPath, func() *source.Scanner { return scanner }, gm, func() string { return "OK" }, modes, nil, ps, "")
	if err := srv.Start(); err != nil {
		t.Fatalf("server start: %v", err)
	}
	t.Cleanup(func() {
		srv.Stop()
		_ = os.Remove(sockPath)
	})
	return srv, sockPath
}

func TestConsoleServer_PolicyDump(t *testing.T) {
	ps := &mockPolicyStore{
		policies: map[string]json.RawMessage{
			"attract_matrix": json.RawMessage(`{"density":0.8,"trail_length":20}`),
			"clock":          json.RawMessage(`{}`),
		},
	}
	_, sockPath := newTestServerWithPolicyStore(t, ps)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("policy dump\n"))
	if !sc.Scan() {
		t.Fatal("expected policy dump response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK ") {
		t.Fatalf("expected OK response, got %q", resp)
	}

	// Extract and parse the JSON
	jsonStr := strings.TrimPrefix(resp, "OK ")
	if !json.Valid([]byte(jsonStr)) {
		t.Fatalf("response is not valid JSON: %q", jsonStr)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal policy dump: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 modes in dump, got %d", len(result))
	}
	if _, ok := result["attract_matrix"]; !ok {
		t.Fatal("missing attract_matrix in policy dump")
	}
	if _, ok := result["clock"]; !ok {
		t.Fatal("missing clock in policy dump")
	}
}

func TestConsoleServer_PolicyDumpNilStore(t *testing.T) {
	// newTestServer creates a server with nil policyStore
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("policy dump\n"))
	if !sc.Scan() {
		t.Fatal("expected policy dump response")
	}
	resp := sc.Text()
	if resp != "ERR policy store not configured" {
		t.Fatalf("expected policy store nil error, got %q", resp)
	}
}

func TestConsoleServer_DisplayPolicyByRegionDotNotation(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	// Query region main.0 — should return active mode "menu" and its policy.
	conn.Write([]byte("display policy main.0\n"))
	if !sc.Scan() {
		t.Fatal("expected display policy response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK mode=menu ") {
		t.Fatalf("expected 'OK mode=menu ...', got %q", resp)
	}
	// The policy portion should be valid JSON.
	jsonPart := strings.TrimPrefix(resp, "OK mode=menu ")
	if !json.Valid([]byte(jsonPart)) {
		t.Fatalf("policy portion is not valid JSON: %q", jsonPart)
	}
}

func TestConsoleServer_DisplayPolicyByRegionBareInt(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	// Query bare integer 0 — resolves to region at coordinator index 0 (main).
	conn.Write([]byte("display policy 0\n"))
	if !sc.Scan() {
		t.Fatal("expected display policy response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK mode=menu ") {
		t.Fatalf("expected 'OK mode=menu ...', got %q", resp)
	}
}

func TestConsoleServer_DisplayPolicyByRegionUnknown(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	// Query a region that doesn't exist.
	conn.Write([]byte("display policy nonexistent.0\n"))
	if !sc.Scan() {
		t.Fatal("expected display policy error response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "ERR unknown mode or region") {
		t.Fatalf("expected error for unknown region, got %q", resp)
	}
	if !strings.Contains(resp, "main.0") {
		t.Fatalf("expected available regions in error, got %q", resp)
	}
}

func TestConsoleServer_DisplayPolicyModeHasPriority(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	// Query by mode name "menu" — should return just the policy (no "mode=..." prefix).
	conn.Write([]byte("display policy menu\n"))
	if !sc.Scan() {
		t.Fatal("expected display policy response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK {") && !strings.HasPrefix(resp, "OK {}") {
		t.Fatalf("expected 'OK {...' (mode-level policy), got %q", resp)
	}
	// Should NOT contain "mode=" prefix since this is a direct mode query.
	if strings.Contains(resp, "mode=") {
		t.Fatalf("mode-level query should not have mode= prefix, got %q", resp)
	}
}

func TestConsoleServer_PolicyUsage(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	// No sub-command
	conn.Write([]byte("policy\n"))
	if !sc.Scan() {
		t.Fatal("expected policy usage response")
	}
	resp := sc.Text()
	if resp != "ERR usage: policy dump" {
		t.Fatalf("expected usage error, got %q", resp)
	}
}

func newTestServerWithPolicyStoreAndConfig(t *testing.T, ps console.PolicyStoreAccess, configPath string) (*console.Server, string) {
	t.Helper()
	sockPath := filepath.Join(t.TempDir(), "cyberhudd_test.sock")
	if runtime.GOOS == "windows" {
		sockPath = filepath.Join(t.TempDir(), "c.sock")
	}
	scanner := source.New(nil, 0)
	gm := gpiomgr.New()
	modes := coordinator.NewState(
		coordinator.Region{Index: 0, Name: "main", Controller: "st7789", Modes: []string{"menu"}, Default: "menu"},
	)
	srv := console.New(sockPath, func() *source.Scanner { return scanner }, gm, func() string { return "OK" }, modes, nil, ps, configPath)
	if err := srv.Start(); err != nil {
		t.Fatalf("server start: %v", err)
	}
	t.Cleanup(func() {
		srv.Stop()
		_ = os.Remove(sockPath)
	})
	return srv, sockPath
}

func TestConsoleServer_FreezePolicySuccess(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "cyberhud.cfg")

	// Write an existing config file with some hardware settings.
	existingConfig := `{"socket":"/run/test.sock","display":{"panel":"waveshare-1.3hat"}}`
	if err := os.WriteFile(configPath, []byte(existingConfig), 0644); err != nil {
		t.Fatalf("write existing config: %v", err)
	}

	ps := &mockPolicyStore{
		policies: map[string]json.RawMessage{
			"attract_matrix": json.RawMessage(`{"density":0.8,"trail_length":20}`),
			"clock":          json.RawMessage(`{}`),
		},
	}
	_, sockPath := newTestServerWithPolicyStoreAndConfig(t, ps, configPath)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("freeze policy\n"))
	if !sc.Scan() {
		t.Fatal("expected freeze policy response")
	}
	resp := sc.Text()
	if resp != "OK" {
		t.Fatalf("expected OK, got %q", resp)
	}

	// Verify the config file was written with policies merged.
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config after freeze: %v", err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}

	// Check non-policy fields preserved.
	if string(result["socket"]) != `"/run/test.sock"` {
		t.Fatalf("socket field not preserved: %s", result["socket"])
	}

	// Check policies field was added.
	if _, ok := result["policies"]; !ok {
		t.Fatal("policies field missing from written config")
	}

	var policies map[string]json.RawMessage
	if err := json.Unmarshal(result["policies"], &policies); err != nil {
		t.Fatalf("unmarshal policies: %v", err)
	}
	if len(policies) != 2 {
		t.Fatalf("expected 2 policies, got %d", len(policies))
	}
	if _, ok := policies["attract_matrix"]; !ok {
		t.Fatal("missing attract_matrix in persisted policies")
	}

	// Verify backup was created.
	bakPath := configPath + ".bak"
	bakData, err := os.ReadFile(bakPath)
	if err != nil {
		t.Fatalf("backup file not created: %v", err)
	}
	if string(bakData) != existingConfig {
		t.Fatalf("backup content mismatch: got %q", string(bakData))
	}
}

func TestConsoleServer_FreezePolicyNoConfigPath(t *testing.T) {
	ps := &mockPolicyStore{
		policies: map[string]json.RawMessage{
			"clock": json.RawMessage(`{}`),
		},
	}
	// Empty config path.
	_, sockPath := newTestServerWithPolicyStoreAndConfig(t, ps, "")
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("freeze policy\n"))
	if !sc.Scan() {
		t.Fatal("expected freeze policy response")
	}
	resp := sc.Text()
	if resp != "ERR config path not configured" {
		t.Fatalf("expected config path error, got %q", resp)
	}
}

func TestConsoleServer_FreezePolicyNilStore(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "cyberhud.cfg")

	// Server with config path but nil policy store.
	sockPath := filepath.Join(t.TempDir(), "cyberhudd_test.sock")
	if runtime.GOOS == "windows" {
		sockPath = filepath.Join(t.TempDir(), "c.sock")
	}
	scanner := source.New(nil, 0)
	gm := gpiomgr.New()
	modes := coordinator.NewState(
		coordinator.Region{Index: 0, Name: "main", Controller: "st7789", Modes: []string{"menu"}, Default: "menu"},
	)
	srv := console.New(sockPath, func() *source.Scanner { return scanner }, gm, func() string { return "OK" }, modes, nil, nil, configPath)
	if err := srv.Start(); err != nil {
		t.Fatalf("server start: %v", err)
	}
	defer srv.Stop()
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("freeze policy\n"))
	if !sc.Scan() {
		t.Fatal("expected freeze policy response")
	}
	resp := sc.Text()
	if resp != "ERR policy store not configured" {
		t.Fatalf("expected policy store nil error, got %q", resp)
	}
}

func TestConsoleServer_FreezePolicyNewFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "new_config.cfg")

	ps := &mockPolicyStore{
		policies: map[string]json.RawMessage{
			"ticker": json.RawMessage(`{"style":"default"}`),
		},
	}
	_, sockPath := newTestServerWithPolicyStoreAndConfig(t, ps, configPath)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	// Config file doesn't exist yet — should still succeed.
	conn.Write([]byte("freeze policy\n"))
	if !sc.Scan() {
		t.Fatal("expected freeze policy response")
	}
	resp := sc.Text()
	if resp != "OK" {
		t.Fatalf("expected OK, got %q", resp)
	}

	// File should now exist with policies.
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal new config: %v", err)
	}
	if _, ok := result["policies"]; !ok {
		t.Fatal("policies field missing from new config")
	}
}

func TestConsoleServer_FreezeWithoutArgs(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "cyberhud.cfg")
	os.WriteFile(configPath, []byte(`{"socket":"/run/test.sock"}`), 0644)

	ps := &mockPolicyStore{
		policies: map[string]json.RawMessage{
			"clock": json.RawMessage(`{}`),
		},
	}
	_, sockPath := newTestServerWithPolicyStoreAndConfig(t, ps, configPath)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	// "freeze" with no args defaults to freeze policy.
	conn.Write([]byte("freeze\n"))
	if !sc.Scan() {
		t.Fatal("expected freeze response")
	}
	resp := sc.Text()
	if resp != "OK" {
		t.Fatalf("expected OK, got %q", resp)
	}
}

func TestConsoleServer_FreezeInvalidSubcommand(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "cyberhud.cfg")

	ps := &mockPolicyStore{
		policies: map[string]json.RawMessage{},
	}
	_, sockPath := newTestServerWithPolicyStoreAndConfig(t, ps, configPath)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("freeze blah\n"))
	if !sc.Scan() {
		t.Fatal("expected freeze error response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "ERR") {
		t.Fatalf("expected ERR for invalid freeze subcommand, got %q", resp)
	}
}

func TestConsoleServer_DisplayRegions(t *testing.T) {
	_, sockPath := newTestServer(t)
	time.Sleep(20 * time.Millisecond)

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	sc := bufio.NewScanner(conn)
	sc.Scan() // greeting

	conn.Write([]byte("display regions\n"))

	// Read the OK line.
	if !sc.Scan() {
		t.Fatal("expected display regions response")
	}
	resp := sc.Text()
	if !strings.HasPrefix(resp, "OK") {
		t.Fatalf("expected OK response, got %q", resp)
	}

	// Read multi-line body.
	_ = conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	var body strings.Builder
	body.WriteString(resp)
	body.WriteByte('\n')
	for sc.Scan() {
		body.WriteString(sc.Text())
		body.WriteByte('\n')
	}
	full := body.String()

	// Should contain info for both regions in our test setup.
	if !strings.Contains(full, "main.0") {
		t.Fatalf("expected 'main.0' region in display regions output, got %q", full)
	}
	if !strings.Contains(full, "aux.0") {
		t.Fatalf("expected 'aux.0' region in display regions output, got %q", full)
	}
	// Should contain mode information.
	if !strings.Contains(full, "mode=menu") {
		t.Fatalf("expected current mode info in output, got %q", full)
	}
}
