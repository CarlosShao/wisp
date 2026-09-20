package llm_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/observe"
)

// AC#4: the text_chain fallback, proven against two LIVE mockllm instances.
//
// The trap this test is written against: a chain that never calls the primary
// still "passes" a naive fallback test. Every case here asserts the primary's
// own request counter, so a regression that skips the primary fails loudly.

// chainFixture stands up a primary (Anthropic dialect, to prove the failover
// is protocol-agnostic) and a fallback (chat dialect) on two live servers.
type chainFixture struct {
	primary, fallback *adaptertest.Mockllm
	runner            *llm.ChainRunner
	failovers         []llm.FailoverEvent
	names             []string
}

func newChain(t *testing.T) *chainFixture {
	t.Helper()
	cf := &chainFixture{
		primary:   adaptertest.StartMockllm(t),
		fallback:  adaptertest.StartMockllm(t),
		failovers: nil,
	}
	cfg := &config.Config{}
	cfg.LLM.Providers = map[string]config.Provider{
		"primary": {
			Protocol: "anthropic",
			BaseURL:  cf.primary.Base + "/v1",
			Models: map[string]config.ModelSpec{
				"mock-small": {Enabled: true, ContextWindow: 8192},
			},
		},
		"fallback": {
			Protocol: "openai-chat",
			BaseURL:  cf.fallback.Base + "/v1",
			Models: map[string]config.ModelSpec{
				"mock-small": {Enabled: true, ContextWindow: 8192},
			},
		},
	}
	cfg.LLM.TextChain = []string{"primary/mock-small", "fallback/mock-small"}

	res := llm.NewResolver(cfg, nil)
	providers, names, err := res.BuildChain(llm.ChainBuildOptions{
		RetryBase: 5 * time.Millisecond,
		RetryMax:  2, // 3 attempts: exactly fail_next(3) exhausts the primary
	})
	if err != nil {
		t.Fatalf("build chain: %v", err)
	}
	if len(providers) != 2 {
		t.Fatalf("chain elements = %d, want 2", len(providers))
	}
	cf.names = names
	cf.runner = &llm.ChainRunner{Elements: providers, Names: names,
		OnFailover: func(ev llm.FailoverEvent) {
			cf.failovers = append(cf.failovers, ev)
		}}
	return cf
}

func chainRequest() *llm.Request {
	return &llm.Request{Model: "mock-small", Messages: []llm.Message{{
		Role:    llm.RoleUser,
		Content: []llm.Content{llm.TextPart{Text: "chain probe"}}}}}
}

// TestFailoverAfterPrimaryExhaustion is the headline AC#4 case.
func TestFailoverAfterPrimaryExhaustion(t *testing.T) {
	cf := newChain(t)
	cf.primary.Reset(t)
	cf.fallback.Reset(t)
	// Three injected 5xx on the PRIMARY only, matching its retry ceiling.
	cf.primary.QueueFault(t, adaptertest.Fault{Status: 500, Times: 3})

	ctx := context.Background()

	collector := llm.NewTurnCollector()
	var events []llm.StreamEvent
	err := cf.runner.Stream(ctx, chainRequest(), func(ev llm.StreamEvent) error {
		events = append(events, ev)
		collector.Observe(ev)
		return nil
	})
	if err != nil {
		t.Fatalf("chain should complete on the fallback, got: %v", err)
	}
	turn := collector.Result()
	if turn.Stop != llm.StopEndTurn || turn.Text == "" {
		t.Fatalf("fallback turn = %+v", turn)
	}
	// The primary WAS actually hit: its own counter, not our bookkeeping.
	if n := cf.primary.RouteCount(t, "messages"); n != 3 {
		t.Errorf("primary saw %d requests, want 3 (initial + 2 retries before failover)", n)
	}
	if n := cf.fallback.RouteCount(t, "chat"); n != 1 {
		t.Errorf("fallback saw %d requests, want 1", n)
	}
	// The switch is never silent: exactly one failover event, correct names.
	if len(cf.failovers) != 1 {
		t.Fatalf("failover events = %d, want 1 (%+v)", len(cf.failovers), cf.failovers)
	}
	ev := cf.failovers[0]
	if ev.From != cf.names[0] || ev.To != cf.names[1] {
		t.Errorf("failover %s -> %s, want %s -> %s", ev.From, ev.To, cf.names[0], cf.names[1])
	}
	if ev.Cause == nil || ev.Cause.Class != observe.ClassProvider {
		t.Errorf("failover cause = %+v, want provider class", ev.Cause)
	}
	if ev.Attempt != 1 || ev.Total != 2 {
		t.Errorf("failover attempt/total = %d/%d, want 1/2", ev.Attempt, ev.Total)
	}
	// No terminal Error event reached the consumer (the fallback succeeded).
	for _, e := range events {
		if e.Type == llm.EvError {
			t.Errorf("a swallowed element failure leaked into the stream: %v", e)
		}
	}
}

// TestDoubleFailureKeepsTaskContextAndIsResumable covers D40#4's second half:
// Error(provider) + task ctx preserved + retry from the breakpoint.
func TestDoubleFailureKeepsTaskContextAndIsResumable(t *testing.T) {
	cf := newChain(t)
	cf.primary.Reset(t)
	cf.fallback.Reset(t)
	cf.primary.QueueFault(t, adaptertest.Fault{Status: 500, Times: 3})
	cf.fallback.QueueFault(t, adaptertest.Fault{Status: 503, Times: 3})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collector := llm.NewTurnCollector()
	req := chainRequest()
	err := cf.runner.Stream(ctx, req, func(ev llm.StreamEvent) error {
		collector.Observe(ev)
		return nil
	})
	if err == nil {
		t.Fatal("expected a provider error after the whole chain failed")
	}
	if got := adaptertest.ClassOf(err); got != observe.ClassProvider {
		t.Errorf("class = %s, want provider (Error(provider) state)", got)
	}
	if collector.Result().Stop != llm.StopError || collector.Result().Err == nil {
		t.Errorf("terminal events = %+v, want Error + Stop{error}", collector.Result())
	}
	// Both elements really ran.
	if n := cf.primary.RouteCount(t, "messages"); n != 3 {
		t.Errorf("primary requests = %d, want 3", n)
	}
	if n := cf.fallback.RouteCount(t, "chat"); n != 3 {
		t.Errorf("fallback requests = %d, want 3", n)
	}
	// Exactly one switch is announced (the last element has no successor, so
	// the terminal Error it emitted is what the user sees - chain.go
	// documents that a failover event is only ever a real switch).
	if len(cf.failovers) != 1 || cf.failovers[0].To != cf.names[1] {
		t.Fatalf("failovers = %+v, want one switch to %s", cf.failovers, cf.names[1])
	}
	// D40#4: the TASK context is still usable - the failure cancelled nothing.
	if ctx.Err() != nil {
		t.Fatalf("task context was cancelled by the provider failure: %v", ctx.Err())
	}
	// Resumable: with the fallback healthy again, the SAME request object and
	// the SAME ctx complete the turn.
	cf.fallback.Reset(t)
	collector2 := llm.NewTurnCollector()
	if err := cf.runner.Stream(ctx, req, func(ev llm.StreamEvent) error {
		collector2.Observe(ev)
		return nil
	}); err != nil {
		t.Fatalf("resume after the transient failure: %v", err)
	}
	if collector2.Result().Text == "" || collector2.Result().Stop != llm.StopEndTurn {
		t.Errorf("resumed turn = %+v", collector2.Result())
	}
}

// TestNoFailoverAfterDeliveredPayload pins the "never append a second answer
// to a half-delivered one" rule (chain.go).
func TestNoFailoverAfterDeliveredPayload(t *testing.T) {
	cf := newChain(t)
	cf.primary.Reset(t)
	cf.fallback.Reset(t)
	// Primary streams real content, then the connection is cut mid-body.
	cf.primary.Control(t, "/__control/truncate", `{"chunks":6}`)

	collector := llm.NewTurnCollector()
	err := cf.runner.Stream(context.Background(), chainRequest(), func(ev llm.StreamEvent) error {
		collector.Observe(ev)
		return nil
	})
	if err == nil {
		t.Fatal("expected the truncated stream to fail")
	}
	if got := adaptertest.ClassOf(err); got != observe.ClassNetwork {
		t.Errorf("class = %s, want network", got)
	}
	if len(cf.failovers) != 0 {
		t.Errorf("failover after delivered payload must not happen: %+v", cf.failovers)
	}
	if n := cf.fallback.RouteCount(t, "chat"); n != 0 {
		t.Errorf("fallback saw %d requests after a half-delivered primary turn, want 0", n)
	}
	turn := collector.Result()
	if !turn.Incomplete {
		t.Error("a truncated delivered turn must stay marked incomplete")
	}
	if turn.Text == "" {
		t.Error("partial content must be retained, not discarded")
	}
}

// TestFailoverChainFromConfigOrder: the chain order IS the config order, and a
// healthy primary never touches the fallback.
func TestFailoverChainOrderHealthyPrimary(t *testing.T) {
	cf := newChain(t)
	cf.primary.Reset(t)
	cf.fallback.Reset(t)
	collector := llm.NewTurnCollector()
	if err := cf.runner.Stream(context.Background(), chainRequest(), func(ev llm.StreamEvent) error {
		collector.Observe(ev)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(collector.Result().Text, "echo:") {
		t.Errorf("primary answer = %q", collector.Result().Text)
	}
	if n := cf.primary.RouteCount(t, "messages"); n != 1 {
		t.Errorf("primary requests = %d, want 1", n)
	}
	if n := cf.fallback.RouteCount(t, "chat"); n != 0 {
		t.Errorf("fallback requests = %d, want 0 (no failover when the primary is healthy)", n)
	}
	if len(cf.failovers) != 0 {
		t.Errorf("unexpected failovers: %+v", cf.failovers)
	}
}

// TestResolverRejectsUnknownProtocol keeps a mis-typed chain element loud
// (ClassConfig naming the catalog entry) instead of silently shorter.
func TestResolverRejectsUnknownChainElement(t *testing.T) {
	cfg := &config.Config{}
	cfg.LLM.Providers = map[string]config.Provider{"p": {
		Protocol: "openai-chat", BaseURL: "http://127.0.0.1:1/v1",
		Models: map[string]config.ModelSpec{"m": {Enabled: true}},
	}}
	cfg.LLM.TextChain = []string{"p/m", "ghost/model"}
	if _, _, err := llm.NewResolver(cfg, nil).BuildChain(llm.ChainBuildOptions{}); err == nil {
		t.Fatal("unknown provider must fail the chain build")
	} else if adaptertest.ClassOf(err) != observe.ClassConfig {
		t.Errorf("class = %v, want config", adaptertest.ClassOf(err))
	}
}

var _ = json.Marshal
