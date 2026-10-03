package llm_test

// Ticket 261 leg p1 - INSTRUMENTATION ONLY (zero production-code change).
//
// AC#0 (a) asked: is there a REAL, RUNNING path by which a model whose
// catalog entry says enabled=false gets selected, dispatched to a provider,
// and priced? Reading, not inference - every claim below is produced by
// executing the same functions the shipped composition root executes
// (cmd/wisp/run.go:435-442: config.LoadFile -> llm.NewResolver ->
// ResolveChain/ResolveRole -> BuildEndpointProvider -> Stream), against the
// live mockllm HTTP server, and priced by the real billing function
// (internal/agent/cost.go:38 Cost.AddUsage, called by loop.go:435).
//
// The three-way ruler shape demanded by the ticket:
//   - ruler 1: enumeration (resolver.go:287-288 `for id := range p.Models`)
//     and selection (resolver.go:119 `spec, ok := p.Models[model]`) asked
//     about an enabled=false entry -> both admit it;
//   - ruler 2: the provider hop + the cost hop -> the request carrying that
//     model id reaches the server (mockllm's own capture, not our bookkeeping)
//     and the real price card on that same disabled entry produces a positive
//     micro-cost;
//   - ruler 3: flip the flag back to true -> the same rulers produce the
//     same readings (no inversion: the flag moves nothing), and the ghost
//     control (a model NOT in the catalog) is refused -> the rulers are live,
//     they react to catalog presence, only `enabled` is inert.
//
// These tests pin TODAY's behaviour. If AC#1 shape (a) lands (filter at the
// enumeration point), ruler 1/2/3 flip red exactly where the promise starts
// being kept - that is the change signal, not a flake.

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
)

// enabled261TOML builds the hand-written shape the ticket names: three model
// entries, one with no `enabled` key at all (the ticket-257 first-run
// guidance outcome), one explicitly false, one explicitly true, plus a real
// price card on the false/true pair. baseURL points at a live mockllm.
func enabled261TOML(baseURL string) string {
	return `schema_version = 2

[llm]
text_chain = ["mock261/t261-off"]

[llm.providers.mock261]
protocol = "openai-chat"
base_url = "` + baseURL + `"

[llm.providers.mock261.models.t261-nokey]
context_window = 4096

[llm.providers.mock261.models.t261-off]
context_window = 4096
enabled = false

[llm.providers.mock261.models.t261-off.price]
in = 2500000
out = 10000000

[llm.providers.mock261.models.t261-on]
context_window = 4096
enabled = true

[llm.providers.mock261.models.t261-on.price]
in = 2500000
out = 10000000
`
}

// loadEnabled261 writes the hand-written file and loads it through the real
// pipeline (readConfigFile inside LoadFile; nil resolver, keyless provider).
func loadEnabled261(t *testing.T, content string) (*config.Config, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config.toml: %v", err)
	}
	cfg, _, err := config.LoadFile(path, nil)
	if err != nil {
		t.Fatalf("LoadFile refused the hand-written file: %v", err)
	}
	return cfg, path
}

func enabled261Request(model, text string) *llm.Request {
	return &llm.Request{
		Model: model,
		Messages: []llm.Message{{
			Role:    llm.RoleUser,
			Content: []llm.Content{llm.TextPart{Text: text}},
		}},
	}
}

// HTTP isolation reuses the package's existing noProxyClient (matrix_14_2_test.go).
// Wire capture: mockllm only records request bodies for the messages/responses
// routes (server.go comment "the chat route is untouched"), so ruler 2/3 read
// the outgoing body through a recording RoundTripper that forwards to the
// live mockllm unchanged. The adapter, retry ladder, limiter and HTTP stack
// stay the shipped ones; the capture happens exactly at the wire seam.

type recordingRT struct {
	base http.RoundTripper
	mu   sync.Mutex
	body string
}

func (r *recordingRT) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		raw, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, err
		}
		r.mu.Lock()
		r.body = string(raw)
		r.mu.Unlock()
		req.Body = io.NopCloser(bytes.NewReader(raw))
	}
	return r.base.RoundTrip(req)
}

func (r *recordingRT) captured() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.body
}

func recordingClient() (*http.Client, *recordingRT) {
	rt := &recordingRT{base: &http.Transport{Proxy: nil}}
	return &http.Client{Transport: rt}, rt
}

// TestTicket261P1EnumerationAndSelectionAdmitDisabled is ruler 1: the
// enumeration loop and the selection lookup both ignore `enabled`, and the
// hand-written missing-key shape lands on false at load and is still admitted.
func TestTicket261P1EnumerationAndSelectionAdmitDisabled(t *testing.T) {
	cfg, _ := loadEnabled261(t, enabled261TOML("http://127.0.0.1:9/v1"))

	// The load-side fact this ruler rests on (measured, not assumed).
	models := cfg.LLM.Providers["mock261"].Models
	if models["t261-off"].Enabled || models["t261-nokey"].Enabled {
		t.Fatalf("precondition broken: off/nokey entries loaded as enabled = true")
	}
	if !models["t261-on"].Enabled {
		t.Fatalf("precondition broken: the explicit true entry loaded as false")
	}

	res := llm.NewResolver(cfg, nil)

	// Enumeration hop (resolver.go:287-288): every id is listed, flag unseen.
	listed := res.DiscoveredModels("mock261")
	for _, want := range []string{"t261-nokey", "t261-off", "t261-on"} {
		if !strings.Contains(strings.Join(listed, ","), want) {
			t.Errorf("DiscoveredModels = %v, want it to contain %q (this is what the enumeration point does TODAY: no Enabled check)", listed, want)
		}
	}

	// Selection hop (resolver.go:119 via ResolveChain): text_chain already
	// names t261-off in the file; the resolver hands back its endpoint.
	eps, err := res.ResolveChain()
	if err != nil || len(eps) != 1 || eps[0].Model != "t261-off" {
		t.Fatalf("ResolveChain(disabled) = %+v, %v; want the disabled model resolved with no error", eps, err)
	}
	// The unset-role fallback takes the first chain element (run.go:436's
	// ResolveRole(RoleChat) shape) - same admission.
	if _, _, err := res.ResolveRole(llm.RoleChat); err != nil {
		t.Errorf("ResolveRole fell back to the chain's disabled element and errored: %v (want: accepted)", err)
	}
	// The nokey shape, driven exactly like cmd/wisp/providers.go:169 drives
	// the resolver for an operator-named target.
	res.TextChain = []string{"mock261/t261-nokey"}
	if _, err := res.ResolveChain(); err != nil {
		t.Errorf("ResolveChain(hand-written nokey) = %v, want nil (missing key == false == still selectable)", err)
	}
}

// TestTicket261P1DisabledModelReachesProviderAndCost is ruler 2: the full
// shipped dispatch - BuildEndpointProvider -> Stream against live mockllm -
// plus the real billing function on the disabled entry's own price card.
func TestTicket261P1DisabledModelReachesProviderAndCost(t *testing.T) {
	mock := adaptertest.StartMockllm(t)
	mock.Reset(t)
	cfg, _ := loadEnabled261(t, enabled261TOML(mock.Base+"/v1"))
	client, wire := recordingClient()

	res := llm.NewResolver(cfg, nil)
	eps, err := res.ResolveChain()
	if err != nil || len(eps) != 1 {
		t.Fatalf("ResolveChain: %v (%+v)", err, eps)
	}
	prov, err := llm.BuildEndpointProvider(eps[0], llm.ChainBuildOptions{
		HTTPClient: func(llm.Endpoint) *http.Client { return client },
	})
	if err != nil {
		t.Fatalf("BuildEndpointProvider: %v", err)
	}
	if got := prov.Info().Model; got != "t261-off" {
		t.Fatalf("provider Info().Model = %q, want t261-off", got)
	}

	collector := llm.NewTurnCollector()
	err = prov.Stream(context.Background(), enabled261Request("t261-off", "ticket261 dispatch"), func(ev llm.StreamEvent) error {
		collector.Observe(ev)
		return nil
	})
	if err != nil {
		t.Fatalf("Stream against the mock provider: %v", err)
	}
	turn := collector.Result()
	// The provider itself counted the hit - not our bookkeeping.
	if n := mock.RouteCount(t, "chat"); n != 1 {
		t.Errorf("mockllm chat route saw %d requests, want 1 (a disabled model still produced a real, billable completion call)", n)
	}
	if body := wire.captured(); !strings.Contains(body, `"model":"t261-off"`) {
		t.Errorf("the captured wire body does not carry the disabled model id:\n%s", body)
	}
	if turn.Usage.InputTokens == 0 || turn.Usage.OutputTokens == 0 {
		t.Fatalf("usage = %+v, want non-zero tokens (mockllm always answers with a usage chunk)", turn.Usage)
	}

	// The cost hop: the SHIPPED pricing function (internal/agent/cost.go,
	// the call loop.go:435 makes) against the disabled entry's own card.
	card := cfg.LLM.Providers["mock261"].Models["t261-off"].Price
	var cost agent.Cost
	delta := cost.AddUsage(card, turn.Usage)
	if delta <= 0 || cost.Micros <= 0 {
		t.Errorf("Cost.AddUsage(priceCardOf(disabled model), usage) = %d (total %d), want > 0: the billing hop has no Enabled check either", delta, cost.Micros)
	}
	t.Logf("TICKET261-COST-HOP READING: model=t261-off enabled=false usage=%+d in/%d out micros=%d currency=%q",
		turn.Usage.InputTokens, turn.Usage.OutputTokens, cost.Micros, cost.Currency)
}

// TestTicket261P1FlagInversionChangesNoReading is ruler 3, both directions on
// the same ruler: the SAME entry measured disabled and re-measured enabled
// must produce the same reachability (proving `enabled` is inert, not that the
// ruler is dead), and a model absent from the catalog must be refused (the
// live control).
func TestTicket261P1FlagInversionChangesNoReading(t *testing.T) {
	mock := adaptertest.StartMockllm(t)

	// pass 1: as loaded from the hand-written file - enabled=false.
	mock.Reset(t)
	cfgOff, path := loadEnabled261(t, enabled261TOML(mock.Base+"/v1"))
	off := reachReading(t, cfgOff, "t261-off", mock)

	// flip: the same entry, the same file, only the flag line changed.
	flipped := strings.Replace(enabled261TOML(mock.Base+"/v1"), "enabled = false", "enabled = true", 1)
	if flipped == enabled261TOML(mock.Base+"/v1") {
		t.Fatal("the flip did not flip anything")
	}
	if err := os.WriteFile(path, []byte(flipped), 0o600); err != nil {
		t.Fatalf("rewrite flipped file: %v", err)
	}
	cfgOn, _, err := config.LoadFile(path, nil)
	if err != nil {
		t.Fatalf("reload flipped: %v", err)
	}
	if !cfgOn.LLM.Providers["mock261"].Models["t261-off"].Enabled {
		t.Fatal("flip did not reach the decoded state")
	}
	mock.Reset(t)
	on := reachReading(t, cfgOn, "t261-off", mock)

	// The same ruler must invert under a working Enabled gate. It does not:
	// every hop reports the disabled model exactly as it reports the enabled
	// one. The reading below is the evidence for AC#0 (a) = YES.
	if !off.listed || !on.listed {
		t.Errorf("enumeration asymmetry: off listed=%v on listed=%v", off.listed, on.listed)
	}
	if !off.resolved || !on.resolved {
		t.Errorf("selection asymmetry: off resolved=%v on resolved=%v", off.resolved, on.resolved)
	}
	if !off.reachedProvider || !on.reachedProvider {
		t.Errorf("provider-hop asymmetry: off=%v on=%v", off.reachedProvider, on.reachedProvider)
	}
	if off.micros <= 0 || on.micros <= 0 {
		t.Errorf("cost-hop asymmetry: off=%d on=%d (both must be positive: the ruler says the disabled entry is billable TODAY)",
			off.micros, on.micros)
	}

	// Live control: the ruler DOES discriminate when the catalog speaks -
	// a model that is absent is refused at the same selection line.
	res := llm.NewResolver(cfgOn, nil)
	res.TextChain = []string{"mock261/t261-ghost"}
	if _, err := res.ResolveChain(); err == nil {
		t.Errorf("ResolveChain(ghost) succeeded: the selection ruler is not reacting to catalog presence either - suspect ruler")
	}
}

// reach is one pass of every ruler for one model id.
type reach struct {
	listed          bool
	resolved        bool
	reachedProvider bool
	micros          int64
}

func reachReading(t *testing.T, cfg *config.Config, model string, mock *adaptertest.Mockllm) reach {
	t.Helper()
	res := llm.NewResolver(cfg, nil)
	var r reach
	for _, id := range res.DiscoveredModels("mock261") {
		if id == model {
			r.listed = true
		}
	}
	res.TextChain = []string{"mock261/" + model}
	eps, err := res.ResolveChain()
	if err != nil || len(eps) != 1 {
		return r
	}
	r.resolved = true
	client, wire := recordingClient()
	prov, err := llm.BuildEndpointProvider(eps[0], llm.ChainBuildOptions{
		HTTPClient: func(llm.Endpoint) *http.Client { return client },
	})
	if err != nil {
		t.Fatalf("BuildEndpointProvider(%s): %v", model, err)
	}
	collector := llm.NewTurnCollector()
	if err := prov.Stream(context.Background(), enabled261Request(model, "ticket261 inversion"), func(ev llm.StreamEvent) error {
		collector.Observe(ev)
		return nil
	}); err != nil {
		t.Fatalf("Stream(%s): %v", model, err)
	}
	turn := collector.Result()
	r.reachedProvider = mock.RouteCount(t, "chat") == 1
	if body := wire.captured(); !strings.Contains(body, `"model":"`+model+`"`) {
		t.Errorf("wire body for %s lacks the model id:\n%s", model, body)
	}
	card := cfg.LLM.Providers["mock261"].Models[model].Price
	var cost agent.Cost
	cost.AddUsage(card, turn.Usage)
	r.micros = cost.Micros
	return r
}
