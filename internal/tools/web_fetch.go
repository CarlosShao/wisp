package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
)

// web_fetch.go implements D34's web.fetch row over the webnet.go guard.
// The three contract faces this file owns:
//
//   - EGRESS: everything connectable goes through w.guard.fetch — URL shape
//     (http/https only, 2048 cap), encoded/private literals rejected before
//     DNS, DNS answers rejected before connect, redirects re-guarded per hop.
//     There is no other route to the network in this file: w.client exists
//     ONLY inside netGuard's Transport, which only netGuard's dialContext
//     builds conns from.
//   - CONVERSION: the result text is Markdown or fenced JSON — htmlToMarkdown
//     (webmd.go) cannot emit live markup, which is the F2 "never rendered as
//     HTML" rule held by construction.
//   - TAINT: the bridge marks every successful web.fetch result under C25
//     (risk.IsSensitiveSource covers "web.fetch"); this file's only taint job
//     is the honest Origin (the FINAL URL after redirects), so the
//     confirmation card can name which page the leaked fragment rode.
//
// The body is never logged by anything here, and the Result text goes to the
// model context only (D37: a model-actionable failure is an IsError result,
// never a Go error out of the tool).

type webFetch struct{ w *webShared }

type webFetchArgs struct {
	URL      string `json:"url"`
	MaxBytes int    `json:"max_bytes"`
}

var webFetchSchema = JSONSchema(`{"type":"object","properties":{` +
	`"url":{"type":"string","description":"网页 URL（仅 http/https；私有/回环网段在建立连接之前即拒绝）"},` +
	`"max_bytes":{"type":"integer","minimum":1,"description":"本次最多读取多少字节（默认 1 MiB，超出截断）"}}` +
	`,"required":["url"],"additionalProperties":false}`)

func (webFetch) Name() string { return "web.fetch" }

func (webFetch) Description() string {
	return "抓取一个网页或 JSON 并转为 Markdown 文本（非 http(s) 协议与私有网段在连接前拒绝；结果按 C25 打污染标记，绝不以 HTML 呈现）"
}

func (webFetch) Parameters() JSONSchema { return webFetchSchema }

// Execute implements Tool.
func (t webFetch) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a webFetchArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	raw := strings.TrimSpace(a.URL)
	if raw == "" {
		return Result{Text: "缺少 url 参数", IsError: true}, nil
	}
	u, err := normalizeGuardedURL(raw)
	if err != nil {
		// Shape/scheme/encoded-private refusal BEFORE DNS, and the same
		// verdicts R5 would book: this refusal does not wait for the gate.
		return Result{Text: "web.fetch: " + err.Error(), IsError: true}, nil
	}
	if host := u.Hostname(); host != "" {
		if err := t.w.firstVisit.ensure(ctx, "web.fetch", host); err != nil {
			return Result{Text: "web.fetch: " + err.Error(), IsError: true}, nil
		}
	}
	cctx, cancel := context.WithTimeout(ctx, t.w.fetchTimeout)
	defer cancel()
	maxBytes := t.w.maxResponseBytes
	if a.MaxBytes > 0 && a.MaxBytes < maxBytes {
		maxBytes = a.MaxBytes
	}
	doc, err := t.w.guard.fetch(cctx, u.String(), maxBytes)
	if err != nil {
		if ctx.Err() != nil {
			return Result{Text: "web.fetch：任务取消或超时，请求未返回", IsError: true}, nil
		}
		return Result{Text: "web.fetch: " + err.Error(), IsError: true}, nil
	}
	if doc.Status < 200 || doc.Status >= 400 {
		// The status is protocol metadata, not body: naming it does not
		// touch the never-log-body rule.
		return Result{Text: fmt.Sprintf("web.fetch: HTTP %d（%s），未返回内容（正文不回显）",
			doc.Status, http.StatusText(doc.Status)), IsError: true, Origin: doc.FinalURL}, nil
	}
	ctype := strings.ToLower(doc.ContentType)
	body := ensureUTF8(doc.Body)
	var text string
	switch {
	case strings.Contains(ctype, "html"):
		text = htmlToMarkdown([]byte(body))
	case strings.Contains(ctype, "json"):
		if pretty := prettyJSON(doc.Body); pretty != "" {
			text = pretty
		} else {
			text = foldInline(ensureUTF8(doc.Body))
		}
	case strings.Contains(ctype, "text/") || ctype == "":
		text = foldInline(body)
	default:
		return Result{Text: fmt.Sprintf(
			"web.fetch: 不支持的内容类型 %q（本工具只取网页/JSON；文档请用 doc.read）",
			primaryType(ctype)), IsError: true, Origin: doc.FinalURL}, nil
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("web.fetch 取得 %d 字节并转换为 Markdown", len(doc.Body)))
	}
	trailer := ""
	if doc.Truncated {
		trailer = fmt.Sprintf("\n\n[truncated=true：响应超过 %d 字节硬上限，已截断。]", maxBytes)
	}
	return Result{
		Text:      text + trailer,
		Origin:    doc.FinalURL,
		Truncated: doc.Truncated,
	}, nil
}

// primaryType trims a content-type to its type/subtype token for messages.
func primaryType(ctype string) string {
	if i := strings.IndexAny(ctype, ";,"); i >= 0 {
		return strings.TrimSpace(ctype[:i])
	}
	return strings.TrimSpace(ctype)
}

// WebFetchDecl is the host-side declaration for web.fetch (D34 row: L0,
// capability net, private ranges refused unconditionally). Facts feeds R5 the
// model-supplied target so the bridge judges it like every other egress.
func WebFetchDecl(d WebDeps) Decl {
	return Decl{
		Capabilities: []Capability{CapNet},
		Needs:        []Capability{CapNet},
		Declared:     risk.L0,
		Resident:     true,
		Provider:     KindBuiltin,
		Timeout:      d.fetchBudget(),
		Facts:        netURLFacts(d.NetAllowlist),
	}
}

// netURLFacts builds the Decl.Facts hook that hands a "url"-parameter call
// to R5 (protocol / private host / allowlist / length). Host decoding here is
// R5's own face (literal + encoded spellings); the socket-time checks stay in
// netGuard — two layers, same frozen list.
func netURLFacts(allowlist []string) func(ctx context.Context, params map[string]any, paths []string) risk.Facts {
	return func(_ context.Context, params map[string]any, _ []string) risk.Facts {
		raw, _ := params["url"].(string)
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return risk.Facts{}
		}
		u, err := url.Parse(raw)
		if err != nil {
			// an unparseable url is a tool-level error at Execute; R5 still
			// gets the raw length and an unknown protocol, which is its
			// fail-closed side already.
			return risk.Facts{Network: &risk.NetTarget{URL: raw}}
		}
		return risk.Facts{Network: &risk.NetTarget{
			Protocol:  strings.ToLower(u.Scheme),
			Host:      u.Hostname(),
			URL:       raw,
			Allowlist: allowlist,
		}}
	}
}
