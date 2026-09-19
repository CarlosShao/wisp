package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
)

// Server is the mockllm core. All state is mutex-guarded; the server itself
// is stateless per request except control injections, counters and golden
// cursors.
type Server struct {
	goldenDir string
	models    []string

	mu       sync.Mutex
	faults   []fault           // FIFO of injected failing responses
	latency  int               // per-chunk delay, ms
	truncate int               // cut streaming responses after N chunks (0 = off)
	totalReq uint64            // reset by /__control/reset
	perRoute map[string]uint64 // reset by /__control/reset

	goldenCur map[string]int // golden name -> next response section index
}

// fault is one injected failure.
type fault struct {
	status     int
	retryAfter string // raw Retry-After header value
	body       string // response body; empty = synthesized error JSON
}

// NewServer builds the mock with the given golden directory and model list.
func NewServer(goldenDir string, models []string) *Server {
	return &Server{
		goldenDir: goldenDir,
		models:    models,
		goldenCur: map[string]int{},
		perRoute:  map[string]uint64{},
	}
}

// Register wires the routes on mux.
func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/chat/completions", s.counted("chat", s.handleChat))
	mux.HandleFunc("POST /v1/responses", s.counted("responses", s.handleResponses))
	mux.HandleFunc("POST /v1/messages", s.counted("messages", s.handleMessages))
	mux.HandleFunc("GET /v1/models", s.counted("models", s.handleModels))
	mux.HandleFunc("POST /__control/fail_next", s.handleFailNext)
	mux.HandleFunc("POST /__control/latency", s.handleLatency)
	mux.HandleFunc("POST /__control/truncate", s.handleTruncate)
	mux.HandleFunc("POST /__control/reset", s.handleReset)
	mux.HandleFunc("GET /__control/state", s.handleState)
}

// counted bumps the request counters around a handler.
func (s *Server) counted(route string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		s.totalReq++
		s.perRoute[route]++
		s.mu.Unlock()
		h(w, r)
	}
}

// takeFault pops the next injected failure, if any.
func (s *Server) takeFault() (fault, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.faults) == 0 {
		return fault{}, false
	}
	f := s.faults[0]
	s.faults = s.faults[1:]
	return f, true
}

func (s *Server) snapshot() (latency, truncate int, queued int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.latency, s.truncate, len(s.faults)
}

// nextGolden consumes and returns the next response section for name.
func (s *Server) nextGolden(name string, sections []goldenResponse) (goldenResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.goldenCur[name]
	if i >= len(sections) {
		return goldenResponse{}, false
	}
	s.goldenCur[name] = i + 1
	return sections[i], true
}

// resetGolden clears golden cursors and counters (used by /__control/reset).
func (s *Server) resetGolden() {
	s.mu.Lock()
	s.goldenCur = map[string]int{}
	s.totalReq = 0
	s.perRoute = map[string]uint64{}
	s.mu.Unlock()
}

// stats snapshots the counters for /__control/state.
func (s *Server) stats() (total uint64, routes map[string]uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	routes = make(map[string]uint64, len(s.perRoute))
	for k, v := range s.perRoute {
		routes[k] = v
	}
	return s.totalReq, routes
}

// reqID renders a unique id for synthesized responses.
func (s *Server) reqID(prefix string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("%s_mock_%d", prefix, s.totalReq)
}

// goldenName extracts the golden selector: model "golden/<name>" or the
// X-Wisp-Golden header.
func goldenName(model string, header string) string {
	if header != "" {
		return header
	}
	if after, ok := strings.CutPrefix(model, "golden/"); ok {
		return after
	}
	return ""
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// listen binds addr with :0 support (tests ask for a free port).
func listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
