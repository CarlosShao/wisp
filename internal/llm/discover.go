package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Auto-discovery (2026-09-19 supplement): OpenAI-compatible providers expose
// GET /v1/models. Importing a discovered model creates a catalog entry with
// capabilities = unknown (all-false, SPEC-03 sec 3.1 "probe-first") and
// enabled = true. Discovery never claims capabilities the model was not
// probed for.

// DiscoverOptions is the discovery request input.
type DiscoverOptions struct {
	BaseURL string // e.g. https://api.openai.com/v1
	APIKey  string // may be empty for local providers
	// HTTPClient overrides the transport; nil = default proxy-aware client.
	HTTPClient *http.Client
	// Timeout bounds the whole call (monotonic via context).
	Timeout time.Duration
}

// DiscoveredModel is one entry of the provider's model list.
type DiscoveredModel struct {
	ID      string `json:"id"`
	OwnedBy string `json:"owned_by"`
}

// DiscoverModels fetches {BaseURL}/models and returns the model list sorted
// by id. Errors are classified (D37) at the seam like every other provider
// interaction.
func DiscoverModels(ctx context.Context, opts DiscoverOptions) ([]DiscoveredModel, error) {
	if opts.BaseURL == "" {
		return nil, observe.New(observe.ClassConfig, "llm: discovery needs a base_url")
	}
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}
	client := opts.HTTPClient
	if client == nil {
		client = NewDefaultHTTPClient("", "")
	}

	url := strings.TrimRight(opts.BaseURL, "/") + "/models"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, observe.Wrap(observe.ClassInternal, err, "llm: build discovery request")
	}
	req.Header.Set("Accept", "application/json")
	if opts.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+opts.APIKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, ClassifyTransportError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, NewHTTPError(HTTPErrorDetail{
			Status:  resp.StatusCode,
			Code:    upstreamErrorCode(body),
			Message: firstLine(body),
		}, ParseRetryAfter(resp.Header.Get("Retry-After"), time.Now()))
	}

	var parsed struct {
		Data []DiscoveredModel `json:"data"`
	}
	dec := json.NewDecoder(io.LimitReader(resp.Body, 4<<20))
	if err := dec.Decode(&parsed); err != nil {
		return nil, observe.Wrap(observe.ClassProvider, err, "llm: /models payload is not valid JSON")
	}
	out := make([]DiscoveredModel, 0, len(parsed.Data))
	for _, m := range parsed.Data {
		if strings.TrimSpace(m.ID) == "" {
			continue
		}
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// ImportDiscovered merges discovered models into a provider's catalog entry
// with capabilities unknown (all-false) and enabled = true. Existing entries
// are left untouched (user curation wins over discovery). The bool result is
// the number of ADDED entries.
func ImportDiscovered(providers map[string]config.Provider, provider string, models []DiscoveredModel) int {
	p, ok := providers[provider]
	if !ok {
		p = config.Provider{Billing: config.BillingPayPerToken}
	}
	if p.Models == nil {
		p.Models = map[string]config.ModelSpec{}
	}
	added := 0
	for _, m := range models {
		if _, exists := p.Models[m.ID]; exists {
			continue
		}
		// Capabilities all-false = unknown until probe (ticket 11) says
		// otherwise; probe results live in provider_health, not here.
		p.Models[m.ID] = config.ModelSpec{
			Display: m.ID,
			Enabled: true,
			Billing: config.BillingPayPerToken,
		}
		added++
	}
	providers[provider] = p
	return added
}

// firstLine returns the first non-empty line of b, bounded (log-safe).
func firstLine(b []byte) string {
	s := string(b)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

// upstreamErrorCode digs the "code" string out of an OpenAI-shaped error body
// without failing on other shapes.
func upstreamErrorCode(body []byte) string {
	var probe struct {
		Error struct {
			Code any    `json:"code"`
			Type string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &probe); err != nil {
		return ""
	}
	switch v := probe.Error.Code.(type) {
	case string:
		return v
	default:
		if probe.Error.Type != "" {
			return probe.Error.Type
		}
		return fmt.Sprintf("%v", v)
	}
}
