package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// web_search.go lands D34's web.search row. The implementation path is NOT
// open any more: ledger R4 (2026-20-20, docs/reports/pending-and-issues.md
// §"已解决" line 50) ruled 搜索 API with vendor anysearch (the user holds a
// key); the results-page scrape is demoted to backup/offline fallback. Call
// sites never learned the difference — everything sits behind SearchProvider
// exactly as the ticket required ("switching providers must not touch
// call-sites"), and the egress runs through the same guarded client as
// web.fetch: api.anysearch.com is a domain like any other for D30③/⑤.
//
// Key handling follows R7/R9: the key is NEVER read from config plaintext,
// argv or a dialog. Production feeds WebDeps.SearchKey a func that pulls
// from internal/secret (DPAPI, ticket 06) or the user's own environment
// reference; an empty key sends the request anonymously — the R9 entry says
// anysearch answers anonymous calls under an IP rate limit, so the degraded
// chain is real, not a promise.
//
// The scrape provider is deliberately NOT shipped as dead code: WebDeps.
// SearchProvider names whatever landed provider exists; selecting "scrape"
// returns an explicit tool-level error text so the model (and the ticket-25
// review) sees the gap instead of a silent fallback.

// SearchQuery is one provider call.
type SearchQuery struct {
	Query      string
	MaxResults int // 1..10; zero = provider default
	Zone       string
	Language   string
}

// SearchResult is one ranked hit.
type SearchResult struct {
	Title   string
	URL     string
	Snippet string
}

// SearchProvider is the R4 seam.
type SearchProvider interface {
	Name() string
	// EgressHost is the one domain the provider connects to (D30⑤ first
	// visit is asked of THIS domain, not of the result URLs — the built-in
	// tool itself only ever talks to the provider).
	EgressHost() string
	Search(ctx context.Context, q SearchQuery) ([]SearchResult, error)
}

// --- the anysearch-compatible provider --------------------------------------

// AnySearchProvider implements SearchProvider against the R9-registered
// interface face: POST {Base}/search with {"query","max_results","zone",
// "language"} and an OPTIONAL Bearer key.
type AnySearchProvider struct {
	Base string        // e.g. https://api.anysearch.com/v1
	Key  func() string // nil-safe; "" sends the call anonymous
	// Client is the guarded HTTP client (nil = a default guarded one built
	// at provider construction by newWebShared; tests inject their own).
	Client *http.Client
}

// NewAnySearchProvider builds the provider over a guarded client.
func NewAnySearchProvider(base string, key func() string, client *http.Client) *AnySearchProvider {
	if client == nil {
		client = newNetGuard(guardOptions{fetchTimeout: defaultFetchTimeout}).client
	}
	return &AnySearchProvider{Base: strings.TrimRight(base, "/"), Key: key, Client: client}
}

// Name implements SearchProvider.
func (p *AnySearchProvider) Name() string { return "anysearch" }

// EgressHost implements SearchProvider.
func (p *AnySearchProvider) EgressHost() string {
	u, err := url.Parse(p.Base)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

type anysearchResponse struct {
	Results []struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Link    string `json:"link"`
		Snippet string `json:"snippet"`
		Desc    string `json:"description"`
	} `json:"results"`
	Items []struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Snippet string `json:"snippet"`
	} `json:"items"`
	Error string `json:"error"`
}

// Search implements SearchProvider.
func (p *AnySearchProvider) Search(ctx context.Context, q SearchQuery) ([]SearchResult, error) {
	if strings.TrimSpace(p.Base) == "" {
		return nil, fmt.Errorf("web.search: anysearch provider 未配置 baseUrl")
	}
	body := map[string]any{"query": q.Query}
	if q.MaxResults > 0 {
		n := q.MaxResults
		if n > 10 {
			n = 10 // the registered API face caps max_results at 1..10
		}
		body["max_results"] = n
	} else {
		body["max_results"] = 10
	}
	if q.Zone != "" {
		body["zone"] = q.Zone
	}
	if q.Language != "" {
		body["language"] = q.Language
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.Base+"/search", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.Key != nil {
		if k := strings.TrimSpace(p.Key()); k != "" {
			// R7: the key value may exist only here, in memory, in the one
			// header. It is never logged, never echoed into an error text.
			req.Header.Set("Authorization", "Bearer "+k)
		}
	}
	c := p.Client
	if c == nil {
		c = http.DefaultClient
	}
	resp, err := c.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("web.search: 请求搜索服务失败：%v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("web.search: 读取响应失败：%v", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("web.search: 搜索服务返回 HTTP %d", resp.StatusCode)
	}
	var doc anysearchResponse
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, fmt.Errorf("web.search: 响应不是合法 JSON（HTTP %d）", resp.StatusCode)
	}
	out := make([]SearchResult, 0, len(doc.Results)+len(doc.Items))
	for _, r := range doc.Results {
		u := r.URL
		if u == "" {
			u = r.Link
		}
		sn := r.Snippet
		if sn == "" {
			sn = r.Desc
		}
		out = append(out, SearchResult{Title: r.Title, URL: u, Snippet: sn})
	}
	for _, r := range doc.Items {
		out = append(out, SearchResult{Title: r.Title, URL: r.URL, Snippet: r.Snippet})
	}
	return out, nil
}

// --- the tool ----------------------------------------------------------------

// unlandedProvider answers for a configured provider name nothing shipped
// (the scrape fallback): the gap is visible, never silently swallowed.
type unlandedProvider struct{ name string }

func (p unlandedProvider) Name() string       { return p.name }
func (p unlandedProvider) EgressHost() string { return "" }
func (p unlandedProvider) Search(context.Context, SearchQuery) ([]SearchResult, error) {
	return nil, fmt.Errorf("web.search: provider %q 尚未落地（票 22 只交付 anysearch 兼容实现；抓结果页方案为 R4 裁定的备用路径，未在本票构建）", p.name)
}

type webSearch struct{ w *webShared }

type webSearchArgs struct {
	Query      string `json:"query"`
	MaxResults int    `json:"max_results"`
	Zone       string `json:"zone"`
	Language   string `json:"language"`
}

var webSearchSchema = JSONSchema(`{"type":"object","properties":{` +
	`"query":{"type":"string","description":"搜索词"},` +
	`"max_results":{"type":"integer","minimum":1,"maximum":10,"description":"最多返回几条（默认 10）"},` +
	`"zone":{"type":"string","enum":["cn","intl"],"description":"搜索区域"},` +
	`"language":{"type":"string","enum":["zh-CN","en"],"description":"结果语言"}}` +
	`,"required":["query"],"additionalProperties":false}`)

func (webSearch) Name() string           { return "web.search" }
func (webSearch) Parameters() JSONSchema { return webSearchSchema }
func (webSearch) Description() string {
	return "搜索网页（经配置的搜索服务；返回标题/链接/摘要列表）。搜索词若含此前读取的外部内容，调用会被升为 L2 强确认（D30①）"
}

// Execute implements Tool.
func (t webSearch) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a webSearchArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if strings.TrimSpace(a.Query) == "" {
		return Result{Text: "缺少 query 参数", IsError: true}, nil
	}
	// D30③ shape rule re-applied tool-side as defense in depth: R5 books the
	// verdict on ANY url-shaped fact, and the query itself can smuggle one.
	if len(a.Query) > risk.MaxURLLength {
		return Result{Text: fmt.Sprintf("query 超过 %d 字符上限（D30③）", risk.MaxURLLength), IsError: true}, nil
	}
	p := t.w.provider
	if p == nil {
		return Result{Text: "web.search: 未配置搜索 provider（WebDeps.SearchProvider 为空，票 22 交付形状为 anysearch 兼容实现）", IsError: true}, nil
	}
	if host := p.EgressHost(); host != "" {
		if err := t.w.firstVisit.ensure(ctx, "web.search", host); err != nil {
			return Result{Text: "web.search: " + err.Error(), IsError: true}, nil
		}
	}
	cctx, cancel := context.WithTimeout(ctx, t.w.searchTimeout)
	defer cancel()
	hits, err := p.Search(cctx, SearchQuery(a))
	if err != nil {
		return Result{Text: "web.search: " + err.Error(), IsError: true}, nil
	}
	if len(hits) == 0 {
		return Result{Text: "web.search: 无结果"}, nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "web.search 结果（provider=%s，%d 条，query 已按 D30③ 长度上限受控）：\n\n", p.Name(), len(hits))
	for i, h := range hits {
		title := strings.TrimSpace(h.Title)
		if title == "" {
			title = h.URL
		}
		fmt.Fprintf(&b, "%d. [%s](%s)\n", i+1, mdEscape(title), h.URL)
		if s := strings.TrimSpace(h.Snippet); s != "" {
			fmt.Fprintf(&b, "   %s\n", foldInline(s))
		}
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("web.search 返回 %d 条结果", len(hits)))
	}
	// Origin names the egress (the search RESULTS themselves are provider
	// data going into the context; SPEC-06 §5 marks web.fetch bodies — the
	// search row is an index of links, and each followed link becomes a
	// marked web.fetch read the moment the model opens it).
	return Result{Text: b.String(), Origin: p.Name() + ":" + p.EgressHost()}, nil
}

// WebSearchDecl is the host-side declaration for web.search (D34 row: L0;
// capability net; the query→L2 escalation is C25/R4, computed by the bridge,
// never declared here).
func WebSearchDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapNet},
		Needs:        []Capability{CapNet},
		Declared:     risk.L0,
		Resident:     true,
		Provider:     KindBuiltin,
		Timeout:      searchBudget,
	}
}

var searchBudget = 20 * time.Second
