package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// fakeFetcher is a test double for the daemon client.
type fakeFetcher struct {
	state    *State
	stateErr error
	// sent records control commands issued via Send.
	sent []string
	// reply maps a command to a canned response; missing commands return OK.
	reply map[string][]string
}

func (f *fakeFetcher) FetchState() (*State, error) {
	return f.state, f.stateErr
}

func (f *fakeFetcher) Send(cmd string) ([]string, error) {
	f.sent = append(f.sent, cmd)
	if r, ok := f.reply[cmd]; ok {
		return r, nil
	}
	return []string{"OK ok"}, nil
}

func twoRegionState() *State {
	return &State{
		Regions: []Region{
			{Index: 0, ID: "main.0", Name: "main", Controller: "st7789", Current: "clock",
				Modes: []ModeRef{
					{ID: "clock", Title: "Clock", Category: "Clock",
						Options: []Option{
							{Key: "show_seconds", Type: "bool", Default: "true", Summary: "Include seconds in the time string.", Allowed: []string{"true", "false"}},
							{Key: "time_format", Type: "string", Default: "24h", Summary: "12-hour or 24-hour time display.", Allowed: []string{"24h", "12h"}},
							{Key: "date_format", Type: "string", Default: "YYYY-MM-DD", Summary: "Date layout or none to hide date row.", Allowed: []string{"YYYY-MM-DD", "DD-MM-YYYY", "MM-DD-YYYY", "none"}},
							{Key: "timezone", Type: "string", Default: "local", Summary: "IANA timezone identifier or local."},
							{Key: "fgcolor", Type: "string", Default: "cyan", Summary: "Foreground color for time text on color panels.", Allowed: []string{"cyan", "green", "amber", "red", "white", "none"}},
							{Key: "show_border", Type: "bool", Default: "false", Summary: "Show rounded border frame."},
						},
						Policy: map[string]string{
							"show_seconds": "true", "time_format": "24h", "date_format": "DD-MM-YYYY",
							"timezone": "Europe/Berlin", "fgcolor": "amber", "show_border": "true",
						}},
					{ID: "dashboard", Title: "Dashboard", Category: "Meta"},
				}},
			{Index: 1, ID: "aux.0", Name: "aux", Controller: "st7735s", Current: "stemma",
				Modes: []ModeRef{{ID: "stemma", Title: "STEMMA", Category: "Data"}}},
		},
		Geometry: Geometry{
			VirtualW: 320, VirtualH: 240,
			Regions: []RegionGeom{
				{Index: 0, Name: "main", X: 0, Y: 0, W: 240, H: 240},
				{Index: 1, Name: "aux", X: 240, Y: 0, W: 80, H: 120},
			},
		},
	}
}

func TestIndexRendersRegionsAndModes(t *testing.T) {
	srv := NewServer(&fakeFetcher{state: twoRegionState()})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"main", "aux", "clock", "dashboard", "stemma", "st7789"} {
		if !strings.Contains(body, want) {
			t.Errorf("index body missing %q", want)
		}
	}
	// Each region renders a carousel of compact selector chips.
	if !strings.Contains(body, `class="car-track"`) {
		t.Errorf("missing carousel track")
	}
	if !strings.Contains(body, `class="mode-chip`) {
		t.Errorf("missing compact mode chips")
	}
	// Chips carry the panel + mode wiring for JS show/hide.
	if !strings.Contains(body, `data-panel="cfg-main"`) || !strings.Contains(body, `data-mode="clock"`) {
		t.Errorf("chip missing data-panel/data-mode wiring; body:\n%s", body)
	}
	// The non-active mode (dashboard) offers a switch form in the config card;
	// the active mode (clock) shows an active badge instead.
	if !strings.Contains(body, `action="/region/main.0/set"`) {
		t.Errorf("missing switch form for a non-active mode on main.0")
	}
	if !strings.Contains(body, `mode-active-badge`) {
		t.Errorf("active mode should render an active badge")
	}
}

func TestSingleConfigCardPerRegion(t *testing.T) {
	srv := NewServer(&fakeFetcher{state: twoRegionState()})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	body := rec.Body.String()

	// One config card per region (keyed by name).
	if !strings.Contains(body, `id="cfg-main"`) || !strings.Contains(body, `id="cfg-aux"`) {
		t.Errorf("expected one config card per region")
	}
	// Within main's card, every mode has a config section, but only the active
	// mode (clock) is visible; others are hidden by default (JS reveals them).
	// clock section visible:
	if !strings.Contains(body, `<section class="mode-config" data-mode="clock">`) {
		t.Errorf("active mode (clock) config section should be visible (not hidden)")
	}
	// dashboard section present but hidden:
	if !strings.Contains(body, `data-mode="dashboard" hidden>`) {
		t.Errorf("inactive mode (dashboard) config section should be hidden by default; body:\n%s", body)
	}
}

func TestIndexRendersCategoryJumps(t *testing.T) {
	srv := NewServer(&fakeFetcher{state: twoRegionState()})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	body := rec.Body.String()

	// main has two categories (Clock, Meta) so category-jump chips appear.
	if !strings.Contains(body, `class="cat-jump"`) {
		t.Errorf("expected category-jump nav for a multi-category region")
	}
	if !strings.Contains(body, `data-cat="Clock"`) || !strings.Contains(body, `data-cat="Meta"`) {
		t.Errorf("expected Clock and Meta category chips; body:\n%s", body)
	}
}

func TestIndexRendersConfigEditor(t *testing.T) {
	srv := NewServer(&fakeFetcher{state: twoRegionState()})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	body := rec.Body.String()

	// Config forms target the global /config endpoint and carry the mode ID.
	if !strings.Contains(body, `action="/config"`) {
		t.Errorf("missing config form")
	}
	if !strings.Contains(body, `<input type="hidden" name="mode" value="clock">`) {
		t.Errorf("config form should carry the target mode id")
	}
	// Enum option "date_format" rendered as a select with the live value selected.
	if !strings.Contains(body, `name="opt.date_format"`) {
		t.Errorf("missing date_format select")
	}
	if !strings.Contains(body, `<option value="DD-MM-YYYY" selected>`) {
		t.Errorf("date_format select should preselect live value 'DD-MM-YYYY'; body:\n%s", body)
	}
	// Free-text option (timezone has no Allowed list) rendered as a text input.
	if !strings.Contains(body, `name="opt.timezone" value="Europe/Berlin"`) {
		t.Errorf("timezone should render as a text input with its live value")
	}
	// Bool option rendered as a checkbox with its hidden marker, checked (true).
	if !strings.Contains(body, `name="bool.show_border"`) {
		t.Errorf("missing bool marker for show_border")
	}
	if !strings.Contains(body, `name="opt.show_border" value="true" checked`) {
		t.Errorf("show_border checkbox should be checked")
	}
}

func TestIndexDaemonErrorBanner(t *testing.T) {
	srv := NewServer(&fakeFetcher{stateErr: &daemonError{msg: "no active mode"}})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "daemon unreachable") || !strings.Contains(body, "no active mode") {
		t.Errorf("expected error banner with daemon message, got:\n%s", body)
	}
}

func TestLayoutSvgRenderedToScale(t *testing.T) {
	srv := NewServer(&fakeFetcher{state: twoRegionState()})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	body := rec.Body.String()

	if !strings.Contains(body, "<svg") {
		t.Fatalf("expected an SVG layout, got none")
	}
	// Two region rects (plus the outline rect) expected.
	if n := strings.Count(body, `class="vd-region"`); n != 2 {
		t.Errorf("expected 2 region rects, got %d", n)
	}
	// Region links target the region anchors for the mode picker.
	if !strings.Contains(body, `href="#region-main"`) {
		t.Errorf("expected layout link to #region-main")
	}
}

func TestSetModeSuccessRedirects(t *testing.T) {
	ff := &fakeFetcher{state: twoRegionState()}
	srv := NewServer(ff)

	form := url.Values{"mode": {"dashboard"}}
	req := httptest.NewRequest(http.MethodPost, "/region/main.0/set", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/" {
		t.Errorf("redirect to %q, want /", loc)
	}
	if len(ff.sent) != 1 || ff.sent[0] != "display set main.0 dashboard" {
		t.Errorf("sent commands = %v, want [display set main.0 dashboard]", ff.sent)
	}
}

func TestSetModeErrorSurfacesNotice(t *testing.T) {
	ff := &fakeFetcher{
		state: twoRegionState(),
		reply: map[string][]string{
			"display set main.0 bogus": {`ERR unsupported mode "bogus" for region main.0`},
		},
	}
	srv := NewServer(ff)

	form := url.Values{"mode": {"bogus"}}
	req := httptest.NewRequest(http.MethodPost, "/region/main.0/set", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "/?notice=") {
		t.Fatalf("expected redirect with notice, got %q", loc)
	}
	if !strings.Contains(loc, url.QueryEscape("unsupported mode")) {
		t.Errorf("notice missing error text: %q", loc)
	}
}

func TestNextMode(t *testing.T) {
	ff := &fakeFetcher{state: twoRegionState()}
	srv := NewServer(ff)

	req := httptest.NewRequest(http.MethodPost, "/region/aux.0/next", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rec.Code)
	}
	if len(ff.sent) != 1 || ff.sent[0] != "display next aux.0" {
		t.Errorf("sent = %v, want [display next aux.0]", ff.sent)
	}
}

func TestFreezeButtonRendered(t *testing.T) {
	srv := NewServer(&fakeFetcher{state: twoRegionState()})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	body := rec.Body.String()
	if !strings.Contains(body, `action="/freeze"`) {
		t.Errorf("expected a freeze form in the page")
	}
}

func TestFreezeButtonHiddenWhenDaemonDown(t *testing.T) {
	srv := NewServer(&fakeFetcher{stateErr: &daemonError{msg: "unreachable"}})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if strings.Contains(rec.Body.String(), `action="/freeze"`) {
		t.Errorf("freeze form should be hidden when there is no state")
	}
}

func TestFreezeSuccess(t *testing.T) {
	ff := &fakeFetcher{
		state: twoRegionState(),
		reply: map[string][]string{"freeze policy": {"OK"}},
	}
	srv := NewServer(ff)

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/freeze", nil))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rec.Code)
	}
	if len(ff.sent) != 1 || ff.sent[0] != "freeze policy" {
		t.Fatalf("sent = %v, want [freeze policy]", ff.sent)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "/?notice=") || !strings.Contains(loc, url.QueryEscape("saved to disk")) {
		t.Errorf("expected success notice, got %q", loc)
	}
}

func TestFreezeErrorSurfacesNotice(t *testing.T) {
	ff := &fakeFetcher{
		state: twoRegionState(),
		reply: map[string][]string{"freeze policy": {"ERR config path not configured"}},
	}
	srv := NewServer(ff)

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/freeze", nil))

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "/?notice=") || !strings.Contains(loc, url.QueryEscape("freeze failed")) {
		t.Errorf("expected freeze-failed notice, got %q", loc)
	}
	if !strings.Contains(loc, url.QueryEscape("config path not configured")) {
		t.Errorf("notice should include the daemon error text, got %q", loc)
	}
}

func TestConfigApplySendsKeyValues(t *testing.T) {
	ff := &fakeFetcher{state: twoRegionState()}
	srv := NewServer(ff)

	// Config is mode-scoped: the form carries the target mode id and applies via
	// the mode-verb form "display <mode> key=value" (writes global policy,
	// whether or not the mode is active on any region).
	form := url.Values{
		"mode":             {"clock"},
		"opt.time_format":  {"12h"},
		"bool.show_border": {"1"},
		"opt.show_border":  {"true"},
	}
	req := httptest.NewRequest(http.MethodPost, "/config", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rec.Code)
	}
	if len(ff.sent) != 1 {
		t.Fatalf("want 1 command, got %v", ff.sent)
	}
	cmd := ff.sent[0]
	if !strings.HasPrefix(cmd, "display clock ") {
		t.Fatalf("unexpected command prefix: %q", cmd)
	}
	// Order of map iteration is non-deterministic; assert both pairs present.
	if !strings.Contains(cmd, "time_format=12h") {
		t.Errorf("missing time_format=12h in %q", cmd)
	}
	if !strings.Contains(cmd, "show_border=true") {
		t.Errorf("missing show_border=true in %q", cmd)
	}
}

func TestConfigApplyToInactiveMode(t *testing.T) {
	// Configuring a mode that is NOT the region's active mode still works: the
	// command targets the mode by id, not the region. (dashboard is not active.)
	ff := &fakeFetcher{state: twoRegionState()}
	srv := NewServer(ff)

	form := url.Values{"mode": {"dashboard"}, "opt.color_accent": {"green"}}
	req := httptest.NewRequest(http.MethodPost, "/config", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if len(ff.sent) != 1 || ff.sent[0] != "display dashboard color_accent=green" {
		t.Fatalf("sent = %v, want [display dashboard color_accent=green]", ff.sent)
	}
}

func TestConfigApplyUncheckedBoolSendsFalse(t *testing.T) {
	ff := &fakeFetcher{state: twoRegionState()}
	srv := NewServer(ff)

	// Unchecked checkbox: marker present, but no opt.show_border value.
	form := url.Values{
		"mode":             {"clock"},
		"bool.show_border": {"1"},
		"opt.time_format":  {"24h"},
	}
	req := httptest.NewRequest(http.MethodPost, "/config", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if len(ff.sent) != 1 {
		t.Fatalf("want 1 command, got %v", ff.sent)
	}
	if !strings.Contains(ff.sent[0], "show_border=false") {
		t.Errorf("unchecked bool should send false, got %q", ff.sent[0])
	}
}

func TestConfigApplyMissingModeIsRejected(t *testing.T) {
	ff := &fakeFetcher{state: twoRegionState()}
	srv := NewServer(ff)

	form := url.Values{"opt.time_format": {"12h"}} // no mode field
	req := httptest.NewRequest(http.MethodPost, "/config", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if len(ff.sent) != 0 {
		t.Errorf("no command should be sent without a mode, got %v", ff.sent)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, url.QueryEscape("no mode")) {
		t.Errorf("expected a no-mode notice, got %q", loc)
	}
}

func TestConfigApplyErrorSurfacesNotice(t *testing.T) {
	ff := &fakeFetcher{
		state: twoRegionState(),
		reply: map[string][]string{}, // default OK; override below
	}
	// Make the clock config command return an ERR.
	ff.reply["display clock time_format=bogus"] = []string{`ERR time_format: must be one of [12h, 24h]`}
	srv := NewServer(ff)

	form := url.Values{"mode": {"clock"}, "opt.time_format": {"bogus"}}
	req := httptest.NewRequest(http.MethodPost, "/config", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "/?notice=") || !strings.Contains(loc, url.QueryEscape("must be one of")) {
		t.Errorf("expected validation error notice, got %q", loc)
	}
}
