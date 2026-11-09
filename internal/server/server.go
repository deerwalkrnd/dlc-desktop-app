// Package server exposes the library over HTTP: a small JSON API, video
// streaming and the embedded single-page web app.
package server

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/deerwalkrnd/dlc-desktop-app/internal/library"
)

// Options configures a Server.
type Options struct {
	// UI is the built web app (the directory containing index.html). If it
	// has no index.html, a placeholder page explains how to build it.
	UI fs.FS

	// AllowedHosts lists acceptable Host header values. Requests with any
	// other Host are rejected, which defeats DNS-rebinding attacks against
	// the local server. Empty means allow any host (LAN mode).
	AllowedHosts []string

	Logger  *slog.Logger
	Version string
}

// Server is the application's HTTP handler.
type Server struct {
	lib     *library.Library
	log     *slog.Logger
	version string
	hosts   map[string]bool
	handler http.Handler
}

// New builds the HTTP handler for lib.
func New(lib *library.Library, opts Options) (*Server, error) {
	s := &Server{lib: lib, log: opts.Logger, version: opts.Version}
	if s.log == nil {
		s.log = slog.Default()
	}
	if len(opts.AllowedHosts) > 0 {
		s.hosts = map[string]bool{}
		for _, h := range opts.AllowedHosts {
			s.hosts[strings.ToLower(h)] = true
		}
	}

	spa, err := newSPA(opts.UI)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.healthz)
	mux.HandleFunc("GET /api/classes", s.listClasses)
	mux.HandleFunc("GET /api/classes/{number}", s.getClass)
	mux.HandleFunc("GET /api/classes/{number}/subjects/{slug}", s.getSubject)
	mux.HandleFunc("GET /api/teachers", s.listTeachers)
	mux.HandleFunc("GET /api/library", s.libraryStatus)
	mux.HandleFunc("POST /api/library/rescan", s.rescan)
	mux.HandleFunc("GET /api/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "no such API endpoint")
	})
	mux.HandleFunc("GET /media/{id}", s.media)
	mux.Handle("GET /", spa)

	s.handler = s.recoverer(securityHeaders(s.hostGuard(mux)))
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

// HTTPServer wraps the handler in an *http.Server with sensible timeouts.
// There is deliberately no WriteTimeout: it would cut off video streams.
func (s *Server) HTTPServer() *http.Server {
	return &http.Server{
		Handler:           s,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       2 * time.Minute,
		ErrorLog:          slog.NewLogLogger(s.log.Handler(), slog.LevelDebug),
	}
}

// LoopbackHosts returns the Host header values a browser sends for a server
// listening on a loopback address, or nil if addr is not loopback.
func LoopbackHosts(addr string) []string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil
	}
	if host != "localhost" {
		ip := net.ParseIP(host)
		if ip == nil || !ip.IsLoopback() {
			return nil
		}
	}
	return []string{
		"localhost:" + port,
		"127.0.0.1:" + port,
		"[::1]:" + port,
	}
}

func (s *Server) hostGuard(next http.Handler) http.Handler {
	if s.hosts == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.hosts[strings.ToLower(r.Host)] {
			http.Error(w, "unexpected Host header", http.StatusMisdirectedRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// sameOrigin rejects cross-site state-changing requests. Browsers always
// send Origin on cross-origin POSTs, so its absence means a same-origin or
// non-browser client.
func sameOrigin(r *http.Request) bool {
	if site := r.Header.Get("Sec-Fetch-Site"); site != "" && site != "same-origin" && site != "none" {
		return false
	}
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	u, err := url.Parse(origin)
	return err == nil && strings.EqualFold(u.Host, r.Host)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Cross-Origin-Resource-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				if v == http.ErrAbortHandler {
					panic(v)
				}
				s.log.Error("panic serving request", "path", r.URL.Path, "panic", fmt.Sprint(v), "stack", string(debug.Stack()))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
