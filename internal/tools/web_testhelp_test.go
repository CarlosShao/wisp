package tools

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync"
	"testing"
	"time"
)

// web_testhelp_test.go — the harness the ticket-22 suites share.
//
// Guard-side layering, stated once so no test can quietly weaken anything:
//
//   - the POLICY (scheme/length/decode/range/redirect) always runs;
//   - what tests inject is the PLUMBING below/beside the policy:
//     lookup (fake DNS), rawDial (a veth to the local listener instead of
//     the real internet — the IP it dials is the one the policy already
//     vetted), and blockedAddr ONLY in the conversion/matrix suites, where
//     the SSRF assertion is not the one being made (a separate suite,
//     web_guard_test.go, runs the frozen predicate with no overrides at
//     all).
//   - newNetGuard with zero options = the production posture. The guard
//     suite uses that shape exclusively.

// hitCount is a listener-side counter: every accepted conn bumps it before
// any bytes move, so "receiver asserts zero hits" is a connect-level claim.
type hitCount struct {
	mu    sync.Mutex
	n     int
	paths []string
}

func (h *hitCount) count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.n
}

func (h *hitCount) record(path string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	h.paths = append(h.paths, path)
}

func (h *hitCount) seen(prefix string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, p := range h.paths {
		if strings.HasPrefix(p, prefix) {
			return true
		}
	}
	return false
}

// countingServer wraps httptest with the hit ledger.
type countingServer struct {
	*httptest.Server
	hits *hitCount
}

func newCountingTLS(t *testing.T, h http.Handler) *countingServer {
	t.Helper()
	hits := &hitCount{}
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.record(r.URL.Path)
		h.ServeHTTP(w, r)
	})
	ts := httptest.NewTLSServer(mux)
	t.Cleanup(ts.Close)
	return &countingServer{Server: ts, hits: hits}
}

func newCountingHTTP(t *testing.T, h http.Handler) *countingServer {
	t.Helper()
	hits := &hitCount{}
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.record(r.URL.Path)
		h.ServeHTTP(w, r)
	})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return &countingServer{Server: ts, hits: hits}
}

// fakeAddr renders the port half of a server URL, so tests can build
// same-port URLs on a fake domain ("http://page.test:PORT/x").
func portOf(u string) string {
	_, p, err := net.SplitHostPort(strings.TrimPrefix(strings.TrimPrefix(u, "https://"), "http://"))
	if err != nil {
		panic(err)
	}
	return p
}

// loopbackLookup maps every DNS name to 127.0.0.1 (used by the conversion /
// matrix / redirect suites where the SERVER is local; the guard's own
// encoded-literal policy still applies before this is consulted).
func loopbackLookup(_ context.Context, _ string) ([]netip.Addr, error) {
	return []netip.Addr{netip.MustParseAddr("127.0.0.1")}, nil
}

// localGuard is the conversion-suite guard: policy shape checks + redirect
// rules intact; the range predicate relaxed to loopback-only-blocking ONLY
// where a test is not making the SSRF claim, plus rawDial straight to the
// resolved numeric address (httptest lives on localhost). See file header:
// the SSRF suite never builds this, it builds newNetGuard(guardOptions{}).
func localGuard(t *testing.T, opts ...func(*guardOptions)) *netGuard {
	t.Helper()
	o := guardOptions{
		lookup:       loopbackLookup,
		dialTimeout:  2 * time.Second,
		fetchTimeout: 5 * time.Second,
		blockedAddr: func(ip netip.Addr) bool {
			return strings.HasPrefix(ip.String(), "192.0.2.") // TEST-NET sentinel, loopback stays reachable
		},
	}
	for _, f := range opts {
		f(&o)
	}
	g := newNetGuard(o)
	// Test cert trust for httptest's TLS server (test-scoped only; the
	// production newWebShared never touches Transport fields here).
	if tr, ok := g.client.Transport.(*http.Transport); ok {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // httptest self-signed CA
	}
	return g
}

// recShared assembles a webShared over a local guard for tool-level suites
// (NOT the guard suite). Everything D30 owns stays wired and on:
// first-visit asks are recorded and (per script) allowed, the Launch hand-off
// is recorded, and the guard runs with its real redirect policy.
type recShared struct {
	w           *webShared
	asks        []string // "tool|domain" per first-visit prompt
	launches    []string // "kind|target" per delivered open
	logs        []string
	askAllow    bool
	askRemember bool
}

func newRecShared(t *testing.T, mut func(*WebDeps)) *recShared {
	t.Helper()
	rs := &recShared{askAllow: true, askRemember: true}
	d := WebDeps{
		FetchTimeout:     5 * time.Second,
		DialTimeout:      2 * time.Second,
		MaxResponseBytes: defaultFetchMaxBytes,
	}
	if mut != nil {
		mut(&d)
	}
	w := &webShared{
		deps:             d,
		fetchTimeout:     5 * time.Second,
		searchTimeout:    5 * time.Second,
		maxResponseBytes: orDefaultInt(d.MaxResponseBytes, defaultFetchMaxBytes),
	}
	w.guard = localGuard(t, func(o *guardOptions) {
		if d.LookupIP != nil {
			o.lookup = d.LookupIP
		}
		if d.MaxRedirects > 0 {
			o.maxRedirects = d.MaxRedirects
		}
		if d.DialTimeout > 0 {
			o.dialTimeout = d.DialTimeout
		}
	})
	rs.w = w
	w.logf = func(format string, args ...any) {
		rs.logs = append(rs.logs, fmt.Sprintf(format, args...))
	}
	w.provider = d.SearchProvider
	if w.provider == nil {
		w.provider = NewAnySearchProvider(d.SearchBaseURL, d.SearchKey, nil)
	}
	w.firstVisit = &firstVisitGate{
		mem: newDomainMemory(d.DomainStorePath),
		ask: func(_ context.Context, tool, domain string) (bool, bool) {
			rs.asks = append(rs.asks, tool+"|"+domain)
			return rs.askAllow, rs.askRemember
		},
		log: w.logf,
	}
	// the Launch seam records instead of touching the OS
	base := d.Launch
	w.deps.Launch = func(ctx context.Context, kind, target string) error {
		rs.launches = append(rs.launches, kind+"|"+target)
		if base != nil {
			return base(ctx, kind, target)
		}
		return nil
	}
	return rs
}

func (r *recShared) fetchTool() webFetch {
	return webFetch{w: r.w}
}

func (r *recShared) searchTool(provider SearchProvider) webSearch {
	if provider != nil {
		r.w.provider = provider
	}
	return webSearch{w: r.w}
}

func (r *recShared) openTool() webOpen     { return webOpen{w: r.w} }
func (r *recShared) fileTool() fileOpen    { return fileOpen{w: r.w} }
func (r *recShared) launchTool() appLaunch { return appLaunch{w: r.w} }

func mustJSON(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// execCall runs one tool Execute with parsed params (unit-level, no bridge).
func execCall(t *testing.T, tool interface {
	Execute(context.Context, json.RawMessage, func(string)) (Result, error)
}, params any,
) Result {
	t.Helper()
	res, err := tool.Execute(context.Background(), mustJSON(t, params), nil)
	if err != nil {
		t.Fatalf("tool returned a Go fault: %v", err)
	}
	return res
}
