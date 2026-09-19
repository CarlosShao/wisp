package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// ---------------------------------------------------------------------------
// Provider catalog consumption (ticket 05 structures).
// ---------------------------------------------------------------------------

func testResolver() *Resolver {
	cfg := config.NewDefaults()
	cfg.LLM.Providers = map[string]config.Provider{
		"openai": {
			Protocol:  "openai-chat",
			BaseURL:   "https://api.openai.com/v1",
			APIKeyRef: "env:WISP_TEST_KEY",
			RPM:       100,
			TPM:       10000,
			Models: map[string]config.ModelSpec{
				"gpt-4o": {ContextWindow: 128000, Enabled: true},
			},
		},
		"deepseek": {
			Protocol: "openai-chat",
			BaseURL:  "https://api.deepseek.com/v1",
			Models: map[string]config.ModelSpec{
				"deepseek-chat": {Enabled: true},
			},
		},
		"anthropic": {
			Protocol: "anthropic", // adapter lands with ticket 11
			BaseURL:  "https://api.anthropic.com",
			Models: map[string]config.ModelSpec{
				"claude-sonnet": {Enabled: true},
			},
		},
	}
	cfg.LLM.TextChain = []string{"openai/gpt-4o", "deepseek/deepseek-chat"}
	cfg.LLM.Roles.Chat = config.Role{Provider: "openai", Model: "gpt-4o", Temperature: 0.5}
	return &Resolver{Providers: cfg.LLM.Providers, TextChain: cfg.LLM.TextChain, Roles: cfg.LLM.Roles}
}

type mapKeys map[string]string

func (m mapKeys) Resolve(ref string) (string, error) {
	if v, ok := m[ref]; ok {
		return v, nil
	}
	return "", observe.New(observe.ClassAuth, "unresolved ref "+ref)
}

func TestResolveChainConsumesCatalog(t *testing.T) {
	r := testResolver()
	r.Keys = mapKeys{"env:WISP_TEST_KEY": "sk-test"}
	eps, err := r.ResolveChain()
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 2 {
		t.Fatalf("chain = %d endpoints", len(eps))
	}
	e0 := eps[0]
	if e0.Provider != "openai" || e0.Model != "gpt-4o" || e0.Protocol != "openai-chat" ||
		e0.BaseURL != "https://api.openai.com/v1" || e0.APIKey != "sk-test" ||
		e0.RPM != 100 || e0.TPM != 10000 || e0.ContextWindow != 128000 {
		t.Errorf("endpoint 0 = %+v", e0)
	}
	// Keyless provider resolves without a key resolver entry.
	if eps[1].APIKey != "" {
		t.Errorf("deepseek unexpectedly has a key")
	}
	if strings.Contains(e0.String(), "sk-test") {
		t.Error("Endpoint.String() must never render the key")
	}
}

func TestResolveChainKeylessRefWithoutResolverIsAuth(t *testing.T) {
	r := testResolver()
	_, err := r.ResolveChain()
	if err == nil {
		t.Fatal("expected auth error for unresolvable key ref")
	}
	if classOf(err) != observe.ClassAuth {
		t.Errorf("class = %s, want auth (Unconfigured)", classOf(err))
	}
}

func TestResolveRoleAndFallback(t *testing.T) {
	r := testResolver()
	r.Keys = mapKeys{"env:WISP_TEST_KEY": "sk-test"}

	ep, settings, err := r.ResolveRole(RoleChat)
	if err != nil {
		t.Fatal(err)
	}
	if ep.Model != "gpt-4o" || settings.Temperature != 0.5 {
		t.Errorf("role chat = %+v %+v", ep, settings)
	}

	// Unset role falls back to the first chain element.
	ep, _, err = r.ResolveRole(RoleSummarize)
	if err != nil {
		t.Fatal(err)
	}
	if ep.Provider != "openai" || ep.Model != "gpt-4o" {
		t.Errorf("summarize fallback = %+v", ep)
	}

	// Unknown role is a config error.
	if _, _, err := r.ResolveRole(RoleName("mystery")); err == nil {
		t.Error("unknown role accepted")
	}
}

func TestBuildChainReportsUnimplementedProtocol(t *testing.T) {
	r := testResolver()
	r.TextChain = []string{"anthropic/claude-sonnet"}
	_, _, err := r.BuildChain(ChainBuildOptions{})
	if err == nil {
		t.Fatal("expected error for unimplemented protocol")
	}
	if classOf(err) != observe.ClassConfig || !strings.Contains(err.Error(), "ticket 11") {
		t.Errorf("err = %v, want config class naming ticket 11", err)
	}
}

// The preset table (ticket 05) must cover the 12 approved providers with a
// protocol and base URL the seam can consume.
func TestPresetTableCoversTwelveProviders(t *testing.T) {
	want := []string{"openai", "anthropic", "deepseek", "qwen", "zhipu", "moonshot",
		"siliconflow", "openrouter", "ollama", "minimax", "mimo", "stepfun"}
	names := config.PresetNames()
	if len(names) != len(want) {
		t.Fatalf("presets = %v, want %d entries", names, len(want))
	}
	for _, w := range want {
		p, ok := config.LookupPreset(w)
		if !ok {
			t.Errorf("preset %q missing", w)
			continue
		}
		if p.Protocol != config.ProtocolOpenAIChat && p.Protocol != config.ProtocolAnthropic {
			t.Errorf("preset %q protocol = %q", w, p.Protocol)
		}
		if p.BaseURL == "" {
			t.Errorf("preset %q has no base_url", w)
		}
	}
	if got, _ := config.LookupPreset("anthropic"); got.Protocol != config.ProtocolAnthropic {
		t.Errorf("anthropic preset protocol = %q", got.Protocol)
	}
}

// ---------------------------------------------------------------------------
// Auto-discovery (/v1/models).
// ---------------------------------------------------------------------------

func TestDiscoverModelsAndImportUnknownCapabilities(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			http.NotFound(w, r)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer sk-disco" {
			t.Errorf("Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[
			{"id":"m2","object":"model","owned_by":"org"},
			{"id":"m1","object":"model","owned_by":"org"},
			{"id":"","object":"model","owned_by":"org"}]}`))
	}))
	defer srv.Close()

	models, err := DiscoverModels(context.Background(), DiscoverOptions{
		BaseURL: srv.URL + "/v1", APIKey: "sk-disco",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(models) != 2 || models[0].ID != "m1" || models[1].ID != "m2" {
		t.Errorf("models = %+v (blank ids dropped, sorted)", models)
	}

	// Import: capabilities unknown (all-false), enabled, existing kept.
	providers := map[string]config.Provider{
		"p": {Models: map[string]config.ModelSpec{
			"m1": {Display: "curated m1", Enabled: true},
		}},
	}
	added := ImportDiscovered(providers, "p", models)
	if added != 1 {
		t.Errorf("added = %d, want 1 (m1 already curated)", added)
	}
	m2 := providers["p"].Models["m2"]
	if !m2.Enabled || m2.Capabilities != (config.Capabilities{}) {
		t.Errorf("imported m2 = %+v, want enabled with unknown capabilities", m2)
	}
	if providers["p"].Models["m1"].Display != "curated m1" {
		t.Error("import clobbered a curated entry")
	}
}

func TestDiscoverModelsClassifiesErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key","code":"invalid_api_key"}}`))
	}))
	defer srv.Close()

	_, err := DiscoverModels(context.Background(), DiscoverOptions{BaseURL: srv.URL + "/v1"})
	if err == nil {
		t.Fatal("expected auth error")
	}
	if classOf(err) != observe.ClassAuth {
		t.Errorf("class = %s, want auth", classOf(err))
	}
}

// ---------------------------------------------------------------------------
// Probe primitive contract.
// ---------------------------------------------------------------------------

func TestProbeContractFourCapabilities(t *testing.T) {
	caps := AllProbeCapabilities()
	if len(caps) != 4 {
		t.Fatalf("capabilities = %v, want fc/vision/thinking/audio", caps)
	}
	for _, cap := range caps {
		cse, ok := ProbeCaseFor(cap)
		if !ok {
			t.Errorf("no case for %s", cap)
			continue
		}
		if cse.Validate == nil {
			t.Errorf("%s: response check undefined", cap)
		}
		req, err := cse.BuildRequest()
		if cap == ProbeAudio {
			// DEFERRED: C7 has no audio part yet; the case must say so.
			if req != nil || cse.NotImplementable == "" {
				t.Errorf("%s: expected a nil request WITH an explicit not-implementable note", cap)
			}
			continue
		}
		if err != nil || req == nil {
			t.Errorf("%s: build request = %v, %v", cap, req, err)
			continue
		}
		if req.Model == "" || len(req.Messages) == 0 {
			t.Errorf("%s: minimal request incomplete: %+v", cap, req)
		}
		// fc and vision carry their capability's payload.
		switch cap {
		case ProbeFC:
			if len(req.Tools) != 1 || req.ToolChoice == nil {
				t.Errorf("fc probe request = %+v", req)
			}
		case ProbeVision:
			hasImage := false
			for _, p := range req.Messages[0].Content {
				if _, ok := p.(ImagePart); ok {
					hasImage = true
				}
			}
			if !hasImage {
				t.Error("vision probe request carries no ImagePart")
			}
		}
	}
}

func TestRunProbeNeverFakesAPass(t *testing.T) {
	// Audio must report not-implementable honestly.
	res := RunProbe(context.Background(), &scriptedProvider{}, ProbeAudio)
	if res.OK || res.Detail == "" {
		t.Errorf("audio probe = %+v, want OK=false with the deferral reason", res)
	}

	// A provider that answers text only fails the fc check.
	res = RunProbe(context.Background(), &scriptedProvider{}, ProbeFC)
	if res.OK {
		t.Error("fc probe passed against a provider that never calls tools")
	}

	// A provider that calls the tool passes.
	res = RunProbe(context.Background(), &fcProvider{}, ProbeFC)
	if !res.OK {
		t.Errorf("fc probe against a tool-calling provider failed: %s", res.Detail)
	}
}

type fcProvider struct{}

func (fcProvider) Stream(ctx context.Context, req *Request, emit func(StreamEvent) error) error {
	emit(StreamEvent{Type: EvToolCallStart, ToolCallID: "c1", ToolName: "echo"})
	emit(StreamEvent{Type: EvToolCallArgsDelta, ToolCallID: "c1", ArgsDelta: `{"text":"ping"}`})
	emit(StreamEvent{Type: EvToolCallEnd, ToolCallID: "c1"})
	emit(StreamEvent{Type: EvStop, Stop: StopToolUse})
	emit(StreamEvent{Type: EvDone})
	return nil
}

func (fcProvider) Info() ProviderInfo { return ProviderInfo{} }

func TestProbeRequestJSONShapes(t *testing.T) {
	// The tool schema must be valid JSON the adapters can pass through.
	cse, _ := ProbeCaseFor(ProbeFC)
	req, _ := cse.BuildRequest()
	var schema map[string]any
	if err := json.Unmarshal(req.Tools[0].Parameters, &schema); err != nil {
		t.Fatalf("fc tool parameters not JSON: %v", err)
	}
}
