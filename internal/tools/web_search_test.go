package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
)

// web_search_test.go — D34 web.search row + the R4 ruling (ledger line 50:
// implementation path = anysearch API, scrape demoted to fallback). The
// provider wire face follows the R9-registered contract, and the tool's
// call sites go through SearchProvider only — switching providers must not
// touch call sites (ticket constraint), pinned here by running the SAME
// tool over two providers.

type fakeProvider struct {
	mu    sync.Mutex
	calls []SearchQuery
	hits  []SearchResult
	err   error
	// egress is the D30(5) domain the tool asks about.
	egress string
}

func (p *fakeProvider) Name() string { return "fake" }
func (p *fakeProvider) EgressHost() string {
	return p.egress
}

func (p *fakeProvider) Search(_ context.Context, q SearchQuery) ([]SearchResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls = append(p.calls, q)
	if p.err != nil {
		return nil, p.err
	}
	return p.hits, nil
}

func (p *fakeProvider) count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.calls)
}

func Test22SearchToolRendersMarkdown(t *testing.T) {
	fp := &fakeProvider{egress: "api.anysearch.com", hits: []SearchResult{
		{Title: "First [Result]", URL: "https://a.test/1", Snippet: "some   snippet"},
		{Title: "", URL: "https://b.test/2", Snippet: ""},
	}}
	rs := newRecShared(t, nil)
	res := execCall(t, rs.searchTool(fp), map[string]any{"query": "wisp test", "max_results": 3})
	if res.IsError {
		t.Fatalf("search: %s", res.Text)
	}
	if !strings.Contains(res.Text, "[First \\[Result\\]](https://a.test/1)") {
		t.Errorf("link render/escape wrong:\n%s", res.Text)
	}
	if !strings.Contains(res.Text, "https://b.test/2") {
		t.Errorf("titleless hit should fall to url:\n%s", res.Text)
	}
	if fp.count() != 1 || fp.calls[0].MaxResults != 3 {
		t.Errorf("provider calls = %+v", fp.calls)
	}
	// D30(5): the provider EGRESS domain gets the one-time prompt (the
	// result URLs are not contacted by web.search at all).
	if len(rs.asks) != 1 || rs.asks[0] != "web.search|api.anysearch.com" {
		t.Errorf("asks = %v", rs.asks)
	}
}

func Test22SearchRefusalAndEmpty(t *testing.T) {
	rs := newRecShared(t, nil)
	res := execCall(t, rs.searchTool(&fakeProvider{err: fmt.Errorf("quota exceeded")}),
		map[string]any{"query": "x"})
	if !res.IsError || !strings.Contains(res.Text, "quota exceeded") {
		t.Errorf("provider failure should surface as tool error: %s", res.Text)
	}
	res = execCall(t, rs.searchTool(&fakeProvider{egress: "api.test"}), map[string]any{"query": "x"})
	if res.IsError || !strings.Contains(res.Text, "无结果") {
		t.Errorf("empty results are a normal answer: %s", res.Text)
	}
	res = execCall(t, rs.searchTool(&fakeProvider{egress: "api.test"}), map[string]any{"query": "  "})
	if !res.IsError {
		t.Errorf("blank query must be refused")
	}
}

func Test22SearchProviderSwitchDoesNotTouchCallSites(t *testing.T) {
	// same tool value, two providers: the R4 "switching providers must not
	// touch call-sites" constraint as an executable sentence.
	rs := newRecShared(t, nil)
	a := &fakeProvider{egress: "a.test", hits: []SearchResult{{Title: "A", URL: "https://a.test"}}}
	b := &fakeProvider{egress: "b.test", hits: []SearchResult{{Title: "B", URL: "https://b.test"}}}
	tool := rs.searchTool(nil)
	tool.w.provider = a
	r1 := execCall(t, tool, map[string]any{"query": "q"})
	tool.w.provider = b
	r2 := execCall(t, tool, map[string]any{"query": "q"})
	if r1.IsError || r2.IsError || !strings.Contains(r1.Text, "https://a.test") || !strings.Contains(r2.Text, "https://b.test") {
		t.Fatalf("provider swap broken: %s / %s", r1.Text, r2.Text)
	}
	if len(rs.asks) != 2 {
		t.Errorf("each provider egress asked once: %v", rs.asks)
	}
}

func Test22ScrapeProviderIsExplicitlyUnlanded(t *testing.T) {
	// R4 demoted scrape to fallback; selecting it must SAY so, not fail
	// silently (the ticket-25 review shape).
	rs := newRecShared(t, nil)
	res := execCall(t, rs.searchTool(unlandedProvider{name: "scrape"}), map[string]any{"query": "q"})
	if !res.IsError || !strings.Contains(res.Text, "尚未落地") {
		t.Errorf("scrape selection: %s", res.Text)
	}
}

// ---------------------------------------------------------------------------
// the anysearch wire face (R9 contract), over a local TLS stand-in
// ---------------------------------------------------------------------------

type capturedReq struct {
	path  string
	auth  string
	body  map[string]any
	count int
}

func Test22AnySearchWireContract(t *testing.T) {
	var mu sync.Mutex
	var got capturedReq
	srv := newCountingTLS(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		got.path = r.URL.Path
		got.auth = r.Header.Get("Authorization")
		got.count++
		mu.Unlock()
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		mu.Lock()
		got.body = body
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"results":[{"title":"T1","url":"https://r1.test","snippet":"s1"},`+
			`{"title":"T2","link":"https://r2.test","description":"s2"}]}`)
	}))
	client := localGuard(t).client
	provider := NewAnySearchProvider("https://api.test:"+portOf(srv.URL)+"/v1",
		func() string { return "SECRETKEY-9988" }, client)

	hits, err := provider.Search(context.Background(), SearchQuery{Query: "hello", MaxResults: 20})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(hits) != 2 || hits[1].URL != "https://r2.test" || hits[1].Snippet != "s2" {
		t.Fatalf("hit parsing wrong: %+v", hits)
	}
	mu.Lock()
	if got.path != "/v1/search" {
		mu.Unlock()
		t.Fatalf("path = %q, want /v1/search", got.path)
	}
	if got.auth != "Bearer SECRETKEY-9988" {
		mu.Unlock()
		t.Errorf("auth header = %q (R9: Authorization: Bearer <key>)", got.auth)
	}
	if q, _ := got.body["query"].(string); q != "hello" {
		mu.Unlock()
		t.Errorf("query field = %v", got.body["query"])
	}
	if n, _ := got.body["max_results"].(float64); n != 10 {
		mu.Unlock()
		t.Errorf("max_results must clamp to the registered 1..10 cap, got %v", got.body["max_results"])
	}
	mu.Unlock()

	// anonymous degraded chain (R9: no key still works, IP-rate-limited):
	provider2 := NewAnySearchProvider("https://api.test:"+portOf(srv.URL)+"/v1", func() string { return "" }, client)
	if _, err := provider2.Search(context.Background(), SearchQuery{Query: "x"}); err != nil {
		t.Fatalf("anonymous search: %v", err)
	}
	mu.Lock()
	auth2 := got.auth
	mu.Unlock()
	if auth2 != "" {
		t.Errorf("empty key must send NO Authorization header, got %q", auth2)
	}

	// R7 hygiene: the key must not land in the tool's rendered output.
	rs := newRecShared(t, nil)
	rs.w.provider = provider
	res := execCall(t, rs.searchTool(nil), map[string]any{"query": "hello"})
	if strings.Contains(res.Text, "SECRETKEY-9988") {
		t.Errorf("R7 violation: key echoed into tool output")
	}
}

func Test22AnySearchErrorShapes(t *testing.T) {
	srv := newCountingTLS(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/500" {
			w.WriteHeader(500)
			fmt.Fprint(w, "INTERNAL-LEAK-MARKER")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{not json`)
	}))
	client := localGuard(t).client
	p := NewAnySearchProvider("https://api.test:"+portOf(srv.URL)+"/v1/500", nil, client)
	_ = p
	// base with path suffix breaks the face; instead drive the failure via a
	// provider pointed at /v1 while the 500 lives at /v1/500: simulate by
	// swapping Base:
	p2 := NewAnySearchProvider("https://api.test:"+portOf(srv.URL)+"/v1", nil, client)
	_, err := p2.Search(context.Background(), SearchQuery{Query: "x"}) // hits /v1/search -> 404 body-less mux fallthrough serves 200 with "{not json"
	if err == nil || !strings.Contains(err.Error(), "JSON") {
		t.Errorf("garbage json must fail with wording, got %v", err)
	}
}

func Test22DefaultProviderWithoutBaseNeverDials(t *testing.T) {
	entries := BuiltinWebEntries(WebDeps{})
	for _, e := range entries {
		if e.Tool.Name() != "web.search" {
			continue
		}
		res, err := e.Tool.Execute(context.Background(), mustJSON(t, map[string]any{"query": "x"}), nil)
		if err != nil {
			t.Fatal(err)
		}
		if !res.IsError || !strings.Contains(res.Text, "baseUrl") {
			t.Errorf("half-wired composition must refuse explicitly: %s", res.Text)
		}
		return
	}
	t.Fatal("web.search not registered")
}
