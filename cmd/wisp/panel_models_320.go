package main

// Ticket 320 AC#4's fetch leg: the page asks the host what one provider's own
// model endpoint reports, and this type is the executor behind that hop. It is
// the Go half of a two-sided change - the page's `config.models.list` literal
// and the marker parser landed first (frontend/src/lib/panel.ts,
// config-view.ts), and the roster ruler in this tree fails until both sides
// carry the name, which is why this file and bridge.go's MethodModelsList move
// together.
//
// WHY IT LIVES IN cmd/wisp AND NOT internal/panel: the panel package declares
// sockets, it does not make network calls (its own rule). The three pieces this
// hop needs - the config manager, the secret store and the llm client factory -
// are all assembled here already, so the fetch is a leaf of the same tree the
// `wisp providers discover` CLI leg uses.
//
// WHAT IT DOES NOT DO: it never returns the credential (the key goes into a
// request header inside internal/llm and nowhere else), it never writes to disk,
// and it does not relay the upstream error body verbatim - the sentence the page
// shows is this file's own, built from the classified error, because the reply is
// page-visible and an upstream body is not this host's to quote (ticket 320
// ruling 4, ledger A972).

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/secret"
)

// modelsFetchBudget320 bounds the ONE synchronous call this hop makes. The hop
// runs on the WebView2 callback thread, so a long budget would freeze the panel
// window; eight seconds is short enough that a stuck provider costs the operator
// a visible pause rather than a hung page, and the page shows a loading state for
// the duration.
const modelsFetchBudget320 = 8 * time.Second

// modelsListErrorCap320 is how much of a classified error sentence this hop
// relays. The errors internal/llm returns are already classified (D37), but their
// text can carry the first line of an upstream body; capping keeps that to one
// bounded line and the page-visible reply stays this host's sentence.
const modelsListErrorCap320 = 200

// modelsListHandler320 implements panel.ModelsListRequestHandler against the same
// pair the settings legs use: the live config manager for base_url/api_key_ref
// and the secret store for the resolved key.
type modelsListHandler320 struct {
	mgr     *config.Manager
	secrets *secret.Store
}

// modelsListPayload320 is the JSON half of the receipt, on the marker line. The
// field names are the page's contract (config-view.ts parses `models[].id` and
// ignores what it does not know); internal/llm's own `owned_by` spelling is
// deliberately NOT reused here so the wire shape has exactly one source.
type modelsListPayload320 struct {
	Provider string               `json:"provider"`
	Models   []modelsListEntry320 `json:"models"`
}

type modelsListEntry320 struct {
	ID      string `json:"id"`
	OwnedBy string `json:"ownedBy,omitempty"`
}

// HandleModelsListRequest is the whole hop: one provider name in, the receipt out.
func (h *modelsListHandler320) HandleModelsListRequest(ctx context.Context, req panel.ComposerRequest) (string, error) {
	provider := strings.TrimSpace(req.ConfigProvider)
	if provider == "" {
		return "", fmt.Errorf("请求里没有指名供应商（configProvider 为空），不知道该拉哪一行")
	}
	if h == nil || h.mgr == nil {
		return "", fmt.Errorf("配置不可用，拉取这一跳没有底座")
	}
	cfg := h.mgr.Config()
	if cfg == nil {
		return "", fmt.Errorf("配置还没读进来，先保存一次再拉")
	}
	p, ok := cfg.LLM.Providers[provider]
	if !ok {
		return "", fmt.Errorf("配置里没有供应商 %q（先「添加供应商」建行，再拉它的模型列表）", provider)
	}
	if strings.TrimSpace(p.BaseURL) == "" {
		return "", fmt.Errorf("供应商 %q 还没设置基址（base_url）；先填「基址」并保存，再拉模型列表", provider)
	}
	key := ""
	if strings.TrimSpace(p.APIKeyRef) != "" {
		k, err := h.secrets.Resolve(p.APIKeyRef)
		if err != nil || k == "" {
			// The Unconfigured path, same posture as the CLI leg: a ref that does
			// not resolve is never answered with an unauthenticated request.
			return "", fmt.Errorf("供应商 %q 的密钥引用解析不了（未配置）；不会用未认证的请求去问它的模型列表", provider)
		}
		key = k
	}

	ctx, cancel := context.WithTimeout(ctx, modelsFetchBudget320)
	defer cancel()
	found, err := llm.DiscoverModels(ctx, llm.DiscoverOptions{
		BaseURL:    p.BaseURL,
		APIKey:     key,
		HTTPClient: llm.NewDefaultHTTPClient(cfg.Net.Proxy.Mode, cfg.Net.Proxy.URL),
	})
	if err != nil {
		return "", fmt.Errorf("拉取失败（%s）：%s", provider, boundedFirstLine320(err.Error(), modelsListErrorCap320))
	}

	payload := modelsListPayload320{Provider: provider, Models: make([]modelsListEntry320, 0, len(found))}
	for _, m := range found {
		payload.Models = append(payload.Models, modelsListEntry320{ID: m.ID, OwnedBy: m.OwnedBy})
	}
	line, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("拉到 %d 个模型但清单编码失败（%v）", len(found), err)
	}
	return fmt.Sprintf("已从 %s 拉到 %d 个模型；勾选后逐枚「添加」即可（保存后重启生效）。\n%s%s",
		provider, len(found), panel.ModelsListPayloadMarker, line), nil
}

// boundedFirstLine320 keeps one line and caps it, so an upstream body cannot ride
// into a page-visible sentence in full.
func boundedFirstLine320(s string, cap int) string {
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	if len(s) > cap {
		s = s[:cap] + "…"
	}
	return s
}
