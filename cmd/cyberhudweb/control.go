package main

import (
	"net/http"
	"net/url"
	"strings"
)

// Control handlers forward mode-switch actions to the daemon as the existing
// socket verbs (display set/next/prev). These are fire-and-confirm: a non-ERR
// reply means success; the verbose OK body is not parsed. On success or failure
// we redirect back to "/" (POST-redirect-GET) so a browser refresh is safe and
// the page re-renders from fresh daemon state. Failures surface via a ?notice=.

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		s.redirectNotice(w, r, "bad form: "+err.Error())
		return
	}
	mode := strings.TrimSpace(r.FormValue("mode"))
	if mode == "" {
		s.redirectNotice(w, r, "no mode selected")
		return
	}
	s.sendControl(w, r, "display set "+id+" "+mode)
}

func (s *Server) handleNext(w http.ResponseWriter, r *http.Request) {
	s.sendControl(w, r, "display next "+r.PathValue("id"))
}

// handleConfig applies option changes to a specific mode by its ID, via the
// mode-verb form "display <mode-id> key=value ...". Because mode policy is
// global per mode (not per region), this writes that mode's policy whether or
// not it is the region's active mode — so the user can pre-configure a mode
// without switching to it. The target mode is carried in the "mode" form field.
//
// Only keys submitted in the form are sent; values are validated by the daemon.
// Bool checkboxes are normalized to true/false. Fire-and-confirm: a non-ERR
// reply is success; the daemon applies and auto-captures the new policy.
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.redirectNotice(w, r, "bad form: "+err.Error())
		return
	}
	mode := strings.TrimSpace(r.FormValue("mode"))
	if mode == "" {
		s.redirectNotice(w, r, "no mode specified for config")
		return
	}

	// Form fields are named "opt.<key>". Checkboxes that are unchecked do not
	// appear in the form, so we rely on a companion hidden "bool.<key>" marker
	// to know a bool option was present and should be sent as false when absent.
	pairs := collectConfigPairs(r)
	if len(pairs) == 0 {
		s.redirectNotice(w, r, "no options to apply")
		return
	}

	cmd := "display " + mode + " " + strings.Join(pairs, " ")
	s.sendControl(w, r, cmd)
}

// collectConfigPairs builds the key=value tokens from the posted form. Text,
// number, and select inputs arrive as "opt.<key>". Bool options arrive as a
// hidden "bool.<key>" marker (always present) plus a checkbox "opt.<key>" that
// is present only when checked.
func collectConfigPairs(r *http.Request) []string {
	var pairs []string
	seen := map[string]bool{}

	// Bool options first, driven by their marker so unchecked boxes send false.
	for key := range r.Form {
		const prefix = "bool."
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		optKey := strings.TrimPrefix(key, prefix)
		val := "false"
		if r.FormValue("opt."+optKey) != "" {
			val = "true"
		}
		pairs = append(pairs, optKey+"="+val)
		seen[optKey] = true
	}

	// Non-bool options.
	for key := range r.Form {
		const prefix = "opt."
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		optKey := strings.TrimPrefix(key, prefix)
		if seen[optKey] {
			continue // already handled as a bool
		}
		val := strings.TrimSpace(r.FormValue(key))
		pairs = append(pairs, optKey+"="+val)
	}
	return pairs
}

func (s *Server) handlePrev(w http.ResponseWriter, r *http.Request) {
	s.sendControl(w, r, "display prev "+r.PathValue("id"))
}

// handleFreeze persists the daemon's current in-memory policy for every mode to
// its config file via the "freeze policy" socket command, so settings survive a
// daemon restart. It is a global action (not per-region). Unlike the other
// controls, success is surfaced with a confirmation notice because the effect
// (a file write) is otherwise invisible in the UI.
func (s *Server) handleFreeze(w http.ResponseWriter, r *http.Request) {
	lines, err := s.fetcher.Send("freeze policy")
	switch {
	case err != nil:
		s.redirectNotice(w, r, err.Error())
	case IsErr(lines):
		s.redirectNotice(w, r, "freeze failed: "+ErrText(lines))
	default:
		s.redirectNotice(w, r, "configuration saved to disk")
	}
}

// sendControl issues a control command and redirects home, surfacing an error
// notice if the daemon rejected the command or the socket failed.
func (s *Server) sendControl(w http.ResponseWriter, r *http.Request, cmd string) {
	lines, err := s.fetcher.Send(cmd)
	switch {
	case err != nil:
		s.redirectNotice(w, r, err.Error())
	case IsErr(lines):
		s.redirectNotice(w, r, ErrText(lines))
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s *Server) redirectNotice(w http.ResponseWriter, r *http.Request, msg string) {
	http.Redirect(w, r, "/?notice="+url.QueryEscape(msg), http.StatusSeeOther)
}
