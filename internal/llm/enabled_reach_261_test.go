package llm_test

// Ticket 261 - the Enabled promise, both ends of the story.
//
// Leg p1 (2026-10-03) landed these as rulers against the DEFECT: an
// enabled=false catalog entry was enumerated, selected, dispatched and
// priced, with zero production readers of ModelSpec.Enabled. Leg r1 (same
// day) landed the fix in internal/llm/resolver.go - the enumeration point
// (DiscoveredModels) skips Enabled==false entries and the selection point
// (resolveEndpoint) refuses them by name - and REWROTE the two rulers that
// described the old behaviour:
//
//   - TestTicket261P1EnumerationAndSelectionAdmitDisabled asserted the
//     disabled entry WAS listed and WAS resolvable; it is now
//     TestTicket261P1EnumerationSkipsAndSelectionRefusesDisabled and asserts
//     the opposite (the promise being kept). The p1 assertions are quoted in
//     probes/261/r1/impl.md section 1.
//   - TestTicket261P1FlagInversionChangesNoReading asserted the flag flips
//     NOTHING (proving the key was inert); it is now
//     TestTicket261P1FlagInversionReversesEveryReading and asserts every
//     reading inverts.
//   - TestTicket261P1DisabledModelReachesProviderAndCost was ruler 2 of the
//     defect reading (disabled entry reaches the provider and bills 90
//     micros); it survives as the inverted ruler
//     TestTicket261P1GateBlocksProviderAndCost: gate on -> the provider sees
//     zero requests and the cost hop reads zero. Same ruler, both colours.
//
// The companion config-side rulers (internal/config/enabled_261_test.go:
// missing key decodes false, settings write materializes false) are
// UNCHANGED - they pin the load/write side, which leg r1 does not touch.

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
	"github.com/CarlosShao/wisp/internal/observe"
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

// TestTicket261P1EnumerationSkipsAndSelectionRefusesDisabled (rewritten by
// leg r1; was TestTicket261P1EnumerationAndSelectionAdmitDisabled) pins the
// promise being KEPT: the enumeration loop skips Enabled==false entries, and
// the selection lookup refuses them by name while still admitting every
// enabled shape - explicit true AND the hand-written missing-key shape that
// the ticket-257 first-run guidance produces (which decodes to... wait, it
// decodes to false: see internal/config/enabled_261_test.go. That shape is
// therefore expected to be REFUSED too - the missing-key entry below is the
// control that it is, exactly like an explicit false).
func TestTicket261P1EnumerationSkipsAndSelectionRefusesDisabled(t *testing.T) {
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

	// Enumeration hop: the disabled and missing-key entries are removed from
	// discovery; the enabled one is listed. (p1's assertion wanted all three
	// listed - that described the defect, quote kept in probes/261/r1/impl.md.)
	listed := res.DiscoveredModels("mock261")
	joined := strings.Join(listed, ",")
	if !strings.Contains(joined, "t261-on") {
		t.Errorf("DiscoveredModels = %v, want the enabled entry %q still listed", listed, "t261-on")
	}
	for _, banned := range []string{"t261-off", "t261-nokey"} {
		if strings.Contains(joined, banned) {
			t.Errorf("DiscoveredModels = %v, want the disabled entry %q removed from discovery (ModelSpec.Enabled promise)", listed, banned)
		}
	}

	// Selection hop: text_chain already names t261-off; the resolver must
	// refuse it with a DISABLED-shaped error - not the unknown-model error,
	// so an operator can tell "no such entry" from "the user turned it off".
	eps, err := res.ResolveChain()
	if err == nil {
		t.Fatalf("ResolveChain(disabled) = %+v, want a named refusal", eps)
	}
	if !strings.Contains(err.Error(), "disabled") {
		t.Errorf("ResolveChain(disabled) error = %v, want it to name the model as disabled", err)
	}
	// Leg r2 (A571 fix B): the refusal carries the missing-key guidance as well
	// (wording addition, description updated with it).
	for _, guidance := range []string{"缺 enabled 键的条目视为关闭", "enabled = true"} {
		if !strings.Contains(err.Error(), guidance) {
			t.Errorf("ResolveChain(disabled) error = %v, want the missing-key guidance %q", err, guidance)
		}
	}
	if strings.Contains(err.Error(), "unknown model") {
		t.Errorf("ResolveChain(disabled) error = %v, must NOT reuse the unknown-model wording", err)
	}
	if class := adaptertest.ClassOf(err); class != observe.ClassConfig {
		t.Errorf("ResolveChain(disabled) class = %v, want config", class)
	}
	// The unset-role fallback walks the same resolveEndpoint line (run.go:436's
	// ResolveRole(RoleChat) shape) - the refusal must hold there too.
	if _, _, err := res.ResolveRole(llm.RoleChat); err == nil || !strings.Contains(err.Error(), "disabled") {
		t.Errorf("ResolveRole fallback to the chain's disabled element = %v, want the same disabled refusal", err)
	}
	// The hand-written missing-key shape is refused identically: it decodes
	// to false (defaults.go never enters maps), so "absent key" cannot be
	// used to sneak an entry past the gate.
	res.TextChain = []string{"mock261/t261-nokey"}
	if _, err := res.ResolveChain(); err == nil || !strings.Contains(err.Error(), "disabled") {
		t.Errorf("ResolveChain(hand-written nokey) = %v, want the disabled refusal (missing key decodes to false)", err)
	} else if !strings.Contains(err.Error(), "enabled = true") {
		// Leg r2 (A571 fix B): the nokey shape is exactly the shape the refusal's
		// guidance answers, so this control must carry it verbatim.
		t.Errorf("ResolveChain(hand-written nokey) error = %v, want the guidance sentence (enabled = true) on the nokey control too", err)
	}
}

// TestTicket261P1GateBlocksProviderAndCost (rewritten by leg r1; was
// TestTicket261P1DisabledModelReachesProviderAndCost) is the inverted ruler 2:
// with the gate in, a disabled entry that text_chain names must never reach
// the provider and never produce a cost reading. The enabled control proves
// the stack itself still works on the same fixture.
func TestTicket261P1GateBlocksProviderAndCost(t *testing.T) {
	mock := adaptertest.StartMockllm(t)
	mock.Reset(t)
	cfg, _ := loadEnabled261(t, enabled261TOML(mock.Base+"/v1"))

	res := llm.NewResolver(cfg, nil)
	_, err := res.ResolveChain()
	if err == nil {
		t.Fatal("ResolveChain resolved a disabled model with the gate in")
	}
	if adaptertest.ClassOf(err) != observe.ClassConfig {
		t.Errorf("class = %v, want config", adaptertest.ClassOf(err))
	}

	// The enabled control: the same stack, same fixture, reaches the provider
	// and bills - proving the zero-count below is the gate, not a dead fixture.
	res.TextChain = []string{"mock261/t261-on"}
	eps, err := res.ResolveChain()
	if err != nil || len(eps) != 1 {
		t.Fatalf("ResolveChain(enabled control): %v (%+v)", err, eps)
	}
	prov, err := llm.BuildEndpointProvider(eps[0], llm.ChainBuildOptions{
		HTTPClient: func(llm.Endpoint) *http.Client { return nil },
	})
	if err != nil {
		t.Fatalf("BuildEndpointProvider(control): %v", err)
	}
	collector := llm.NewTurnCollector()
	if err := prov.Stream(context.Background(), enabled261Request("t261-on", "ticket261 gate control"), func(ev llm.StreamEvent) error {
		collector.Observe(ev)
		return nil
	}); err != nil {
		t.Fatalf("Stream(control): %v", err)
	}
	if n := mock.RouteCount(t, "chat"); n != 1 {
		t.Errorf("mockllm chat route saw %d requests for the ENABLED control, want 1", n)
	}
	turn := collector.Result()
	if turn.Usage.InputTokens == 0 || turn.Usage.OutputTokens == 0 {
		t.Fatalf("control usage = %+v, want non-zero tokens", turn.Usage)
	}
	card := cfg.LLM.Providers["mock261"].Models["t261-on"].Price
	var cost agent.Cost
	if delta := cost.AddUsage(card, turn.Usage); delta <= 0 || cost.Micros <= 0 {
		t.Errorf("Cost.AddUsage(priceCardOf(enabled control), usage) = %d (total %d), want > 0", delta, cost.Micros)
	}

	// The disabled side: the provider hop must see ZERO requests (the gate is
	// upstream of dispatch) and the cost hop must read ZERO - the price card
	// is still on the entry (the catalog entry survives), but no request can
	// ever be priced against it.
	if n := mock.RouteCount(t, "chat"); n != 1 {
		t.Errorf("mockllm chat route saw %d requests total, want exactly 1 (the enabled control only): the disabled entry reached the provider", n)
	}
	wireBody := ""
	_ = wireBody // wire capture lives in the inversion test; here the route
	// count alone separates reach/no-reach because only the control resolves.
	offCard := cfg.LLM.Providers["mock261"].Models["t261-off"].Price
	var offCost agent.Cost
	if delta := offCost.AddUsage(offCard, llm.Usage{}); delta != 0 || offCost.Micros != 0 {
		t.Errorf("the disabled entry produced a cost reading %d/%d with zero usage, want 0/0", delta, offCost.Micros)
	}
}

// TestTicket261P1FlagInversionReversesEveryReading (rewritten by leg r1; was
// TestTicket261P1FlagInversionChangesNoReading) is the two-direction ruler on
// one entry: SAME file, only the flag line changed. disabled -> not listed /
// refused / unreachable / zero cost; flipped true -> all four reverse. A ghost
// model absent from the catalog keeps the unknown-model wording, distinct
// from the disabled refusal.
func TestTicket261P1FlagInversionReversesEveryReading(t *testing.T) {
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

	// Every reading must invert (p1 asserted the opposite: nothing inverted).
	if off.listed {
		t.Errorf("enumeration: the disabled entry was listed, want not listed")
	}
	if !on.listed {
		t.Errorf("enumeration: the enabled entry is missing from discovery, want listed")
	}
	if off.resolved {
		t.Errorf("selection: the disabled entry resolved, want refused")
	}
	if !on.resolved {
		t.Errorf("selection: the enabled entry failed to resolve")
	}
	if off.reachedProvider {
		t.Errorf("provider-hop: the disabled entry reached the provider, want zero requests")
	}
	if !on.reachedProvider {
		t.Errorf("provider-hop: the enabled entry never reached the provider")
	}
	if off.micros != 0 {
		t.Errorf("cost-hop: the disabled entry produced %d micros, want 0", off.micros)
	}
	if on.micros <= 0 {
		t.Errorf("cost-hop: the enabled entry produced %d micros, want > 0", on.micros)
	}

	// Live control, wording split: a model ABSENT from the catalog is refused
	// with the unknown-model wording - visibly different from the disabled
	// refusal an operator sees for an entry that exists but is switched off.
	res := llm.NewResolver(cfgOn, nil)
	res.TextChain = []string{"mock261/t261-ghost"}
	_, ghostErr := res.ResolveChain()
	if ghostErr == nil {
		t.Fatal("ResolveChain(ghost) succeeded")
	}
	if !strings.Contains(ghostErr.Error(), "unknown model") {
		t.Errorf("ghost error = %v, want the unknown-model wording", ghostErr)
	}
	if strings.Contains(ghostErr.Error(), "disabled") {
		t.Errorf("sanity: ghost error %v unexpectedly mentions disabled", ghostErr)
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
	r.reachedProvider = mock.RouteCount(t, "chat") >= 1
	if body := wire.captured(); !strings.Contains(body, `"model":"`+model+`"`) {
		t.Errorf("wire body for %s lacks the model id:\n%s", model, body)
	}
	card := cfg.LLM.Providers["mock261"].Models[model].Price
	var cost agent.Cost
	cost.AddUsage(card, turn.Usage)
	r.micros = cost.Micros
	return r
}
