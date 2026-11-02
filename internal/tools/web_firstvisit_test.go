package tools

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// web_firstvisit_test.go — AC row 5 (second clause): "first-visit prompt
// appears once, remembered." D30 layer 5, SPEC-06 §6.5.

func Test22FirstVisitAppearsOncePerDomain(t *testing.T) {
	srv := newCountingHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "hello world")
	}))
	rs := newRecShared(t, nil)
	url := "http://page.test:" + portOf(srv.URL) + "/x"

	r1 := execCall(t, rs.fetchTool(), map[string]any{"url": url})
	if r1.IsError {
		t.Fatalf("first fetch: %s", r1.Text)
	}
	if len(rs.asks) != 1 || rs.asks[0] != "web.fetch|page.test" {
		t.Fatalf("asks = %v, want exactly one web.fetch|page.test", rs.asks)
	}

	r2 := execCall(t, rs.fetchTool(), map[string]any{"url": url})
	if r2.IsError {
		t.Fatalf("second fetch: %s", r2.Text)
	}
	if len(rs.asks) != 1 {
		t.Errorf("second visit re-asked: %v", rs.asks)
	}

	// another domain asks again (per-domain, not global)
	r3 := execCall(t, rs.fetchTool(), map[string]any{"url": "http://other.test:" + portOf(srv.URL) + "/y"})
	if r3.IsError {
		t.Fatalf("other-domain fetch: %s", r3.Text)
	}
	if len(rs.asks) != 2 || rs.asks[1] != "web.fetch|other.test" {
		t.Errorf("asks = %v, want 2nd entry for other.test", rs.asks)
	}
	// and the shared memory spans tools (fetch then open of the same
	// domain = one prompt total, D30⑤'s wording is per DOMAIN)
	rs.askRemember = true
	r4 := execCall(t, rs.openTool(), map[string]any{"url": "http://page.test/x"})
	if r4.IsError {
		t.Fatalf("web.open known domain: %s", r4.Text)
	}
	if len(rs.asks) != 2 {
		t.Errorf("web.open re-asked a remembered domain: %v", rs.asks)
	}
}

func Test22FirstVisitRememberFalseKeepsAsking(t *testing.T) {
	srv := newCountingHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "hi")
	}))
	rs := newRecShared(t, nil)
	rs.askRemember = false // user answered "allow once, do not remember"
	url := "http://page.test:" + portOf(srv.URL) + "/x"
	for i := 1; i <= 3; i++ {
		res := execCall(t, rs.fetchTool(), map[string]any{"url": url})
		if res.IsError {
			t.Fatalf("visit %d: %s", i, res.Text)
		}
		if len(rs.asks) != i {
			t.Fatalf("asks after visit %d = %v, want %d prompts", i, rs.asks, i)
		}
	}
	// and NOT written to the store
	if rs.w.firstVisit.mem.known("page.test") {
		t.Errorf("non-remembered domain was persisted anyway")
	}
}

func Test22FirstVisitRefusalStopsTheCall(t *testing.T) {
	srv := newCountingHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "must never serve a refused domain")
	}))
	rs := newRecShared(t, nil)
	rs.askAllow = false
	res := execCall(t, rs.fetchTool(), map[string]any{"url": "http://page.test:" + portOf(srv.URL) + "/x"})
	if !res.IsError || !strings.Contains(res.Text, "首次访问") {
		t.Errorf("refusal must answer a tool-level rejection: %s", res.Text)
	}
	if n := srv.hits.count(); n != 0 {
		t.Errorf("refused first-visit still reached the server %d times", n)
	}
	// refusal does NOT remember: a later allow is still asked fresh
	rs.askAllow = true
	res2 := execCall(t, rs.fetchTool(), map[string]any{"url": "http://page.test:" + portOf(srv.URL) + "/x"})
	if res2.IsError {
		t.Errorf("after a new allow the fetch should work: %s", res2.Text)
	}
	if len(rs.asks) != 2 {
		t.Errorf("asks = %v, want refusal + allow", rs.asks)
	}
}

func Test22FirstVisitUnwiredFailsClosed(t *testing.T) {
	// A composition that forgot the prompt seam: WebDeps{} through the REAL
	// constructor, so this exercises production wiring shape, not a hand-
	// built struct. No network is attempted either (the refusal precedes
	// the guard).
	entries := BuiltinWebEntries(WebDeps{})
	var ft Tool
	for _, e := range entries {
		if e.Tool.Name() == "web.fetch" {
			ft = e.Tool
		}
	}
	if ft == nil {
		t.Fatal("no web.fetch in BuiltinWebEntries")
	}
	res, err := ft.Execute(context.Background(),
		mustJSON(t, map[string]any{"url": "https://example.com/x"}), nil)
	if err != nil {
		t.Fatalf("Go fault: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "首次访问") || !strings.Contains(res.Text, "fail-closed") {
		t.Errorf("unwired prompt channel must refuse fail-closed, got: %s", res.Text)
	}
}

func Test22FirstVisitIPsAreNotDomains(t *testing.T) {
	srv := newCountingHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "ip-served")
	}))
	port := portOf(srv.URL)
	rs := newRecShared(t, nil)
	// stand a public-shape IP in front of the local server BELOW the judge
	rs.w.guard = localGuard(t, func(o *guardOptions) {
		o.rawDial = func(ctx context.Context, network, addr string) (net.Conn, error) {
			d := &net.Dialer{Timeout: 2 * time.Second}
			return d.DialContext(ctx, network, "127.0.0.1:"+port)
		}
	})
	res := execCall(t, rs.fetchTool(), map[string]any{"url": "http://8.8.8.8:" + port + "/x"})
	if res.IsError {
		t.Fatalf("public IP fetch: %s", res.Text)
	}
	if len(rs.asks) != 0 {
		t.Errorf("IP literal asked a domain prompt: %v", rs.asks)
	}
}

func Test22FirstVisitMemoryPersistsAcrossProcesses(t *testing.T) {
	store := filepath.Join(t.TempDir(), "domains.json")
	srv := newCountingHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "persist me")
	}))
	url := "http://page.test:" + portOf(srv.URL) + "/x"

	rs := newRecShared(t, func(d *WebDeps) { d.DomainStorePath = store })
	if res := execCall(t, rs.fetchTool(), map[string]any{"url": url}); res.IsError {
		t.Fatalf("first-process fetch: %s", res.Text)
	}
	if len(rs.asks) != 1 {
		t.Fatalf("asks = %v", rs.asks)
	}
	b, err := os.ReadFile(store)
	if err != nil {
		t.Fatalf("store file: %v", err)
	}
	if !strings.Contains(string(b), "page.test") {
		t.Errorf("store content missing domain: %s", b)
	}

	// a FRESH shared over the same path = the next process launch: no re-prompt
	rs2 := newRecShared(t, func(d *WebDeps) { d.DomainStorePath = store })
	if res := execCall(t, rs2.fetchTool(), map[string]any{"url": url}); res.IsError {
		t.Fatalf("second-process fetch: %s", res.Text)
	}
	if len(rs2.asks) != 0 {
		t.Errorf("second process re-asked a remembered domain: %v", rs2.asks)
	}
}

func Test22FirstVisitStoreCorruptionStillPrompts(t *testing.T) {
	// A corrupt store must NEVER silently ALLOW (memory only suppresses
	// prompts; losing it means more prompts, the fail-closed direction).
	store := filepath.Join(t.TempDir(), "domains.json")
	if err := os.WriteFile(store, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	srv := newCountingHTTP(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "ok")
	}))
	rs := newRecShared(t, func(d *WebDeps) { d.DomainStorePath = store })
	url := "http://page.test:" + portOf(srv.URL) + "/x"
	if res := execCall(t, rs.fetchTool(), map[string]any{"url": url}); res.IsError {
		t.Fatalf("fetch with corrupt store: %s", res.Text)
	}
	if len(rs.asks) != 1 {
		t.Errorf("corrupt store must still ask (fail toward prompts), got %v", rs.asks)
	}
	// and a later remember rewrites the file cleanly
	b, _ := os.ReadFile(store)
	if !strings.Contains(string(b), "page.test") {
		t.Errorf("remember did not repair the store: %s", b)
	}
}
