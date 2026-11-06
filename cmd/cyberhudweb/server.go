package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

// Server is the HTTP front end. It renders the read model from a Fetcher and
// forwards control actions to the daemon.
type Server struct {
	fetcher Fetcher
	tmpl    *template.Template
	mux     *http.ServeMux
}

// NewServer builds the HTTP server around a Fetcher (the live daemon Client in
// production, a fake in tests).
func NewServer(fetcher Fetcher) *Server {
	s := &Server{
		fetcher: fetcher,
		tmpl:    template.Must(template.New("").Funcs(templateFuncs()).ParseFS(templateFS, "templates/*.html")),
		mux:     http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	static, _ := fs.Sub(staticFS, "static")
	s.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	s.mux.HandleFunc("GET /", s.handleIndex)
	s.mux.HandleFunc("POST /region/{id}/set", s.handleSet)
	s.mux.HandleFunc("POST /region/{id}/next", s.handleNext)
	s.mux.HandleFunc("POST /region/{id}/prev", s.handlePrev)
	s.mux.HandleFunc("POST /config", s.handleConfig)
	s.mux.HandleFunc("POST /freeze", s.handleFreeze)
}

// ServeHTTP lets Server satisfy http.Handler (used by httptest).
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }

// ListenAndServe starts the HTTP server on addr.
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

// pageData is the template model for the index page.
type pageData struct {
	State  *State
	Err    string // daemon/connection error banner
	Notice string // transient action result (e.g. switch failure)
	Layout *LayoutView
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := pageData{Notice: r.URL.Query().Get("notice")}

	st, err := s.fetcher.FetchState()
	if err != nil {
		data.Err = err.Error()
	} else {
		data.State = st
		data.Layout = BuildLayout(st, layoutTargetW, layoutTargetH)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		// lower is used to make case-insensitive comparisons in templates.
		"lower": strings.ToLower,
	}
}
