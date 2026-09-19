package openaichat

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/observe"
)

// goldenDir points at the shared fixtures (the SAME files tools/mockllm
// replays - one format, two runners). Resolved from this file's location so
// the tests never depend on the test binary's working directory.
var goldenDir = func() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "testdata", "golden")
}()

// serveGolden loads a golden fixture and serves it sequentially.
func serveGolden(t *testing.T, name string, pace bool) (*httptestServer, *golden.Replayer) {
	t.Helper()
	rs, err := golden.LoadFile(filepath.Join(goldenDir, name+".sse"))
	if err != nil {
		t.Fatalf("load golden %s: %v", name, err)
	}
	rep := golden.NewReplayer(rs)
	rep.Pace = pace
	srv := rep.Server()
	t.Cleanup(srv.Close)
	return &httptestServer{srv}, rep
}

// httptestServer wraps *httptest.Server so tests can pass a BaseURL.
type httptestServer struct {
	*wrapper
}

// adapter over a golden-backed server.
func (s *httptestServer) adapter(t *testing.T, compatLoose, allowMissingUsage bool) *Adapter {
	t.Helper()
	return New(endpointOptionsFor(s.URL(), compatLoose, allowMissingUsage))
}

// drain runs the adapter and returns the collected events + final turn.
func drain(t *testing.T, p llm.LlmProvider, req *llm.Request) ([]llm.StreamEvent, llm.TurnResult, error) {
	t.Helper()
	var events []llm.StreamEvent
	collector := llm.NewTurnCollector()
	err := p.Stream(context.Background(), req, func(ev llm.StreamEvent) error {
		events = append(events, ev)
		collector.Observe(ev)
		return nil
	})
	return events, collector.Result(), err
}

func baseRequest() *llm.Request {
	return &llm.Request{Model: "mock-small", Messages: []llm.Message{{
		Role:    llm.RoleUser,
		Content: []llm.Content{llm.TextPart{Text: "hi"}},
	}}}
}

// ---------------------------------------------------------------------------
// Golden replay: tool-call assembly.
// ---------------------------------------------------------------------------

func TestGoldenToolCallAssembly(t *testing.T) {
	srv, _ := serveGolden(t, "tool-call", false)
	a := srv.adapter(t, false, false)

	events, turn, err := drain(t, a, baseRequest())
	if err != nil {
		t.Fatal(err)
	}
	if turn.Stop != llm.StopToolUse {
		t.Errorf("stop = %s, want tool_use", turn.Stop)
	}
	if turn.Text != "Let me check." {
		t.Errorf("text = %q", turn.Text)
	}
	if len(turn.ToolCalls) != 2 {
		t.Fatalf("tool calls = %d, want 2", len(turn.ToolCalls))
	}
	tc0 := turn.ToolCalls[0]
	if tc0.ID != "call_a1" || tc0.Name != "get_weather" || string(tc0.Args) != `{"city":"Zhuhai"}` || !tc0.Complete {
		t.Errorf("tool call 0 = %+v", tc0)
	}
	tc1 := turn.ToolCalls[1]
	if tc1.ID != "call_b2" || tc1.Name != "get_time" || string(tc1.Args) != "{}" || !tc1.Complete {
		t.Errorf("tool call 1 = %+v", tc1)
	}
	if turn.Usage != (llm.Usage{InputTokens: 21, OutputTokens: 9, CachedTokens: 8}) {
		t.Errorf("usage = %+v", turn.Usage)
	}
	if turn.Incomplete || turn.Err != nil {
		t.Errorf("turn = %+v", turn)
	}

	// Event-order contract: ArgsDelta only between Start and End of its id.
	open := map[string]bool{}
	for _, ev := range events {
		switch ev.Type {
		case llm.EvToolCallStart:
			open[ev.ToolCallID] = true
		case llm.EvToolCallArgsDelta:
			if !open[ev.ToolCallID] {
				t.Fatalf("ArgsDelta before Start for %s", ev.ToolCallID)
			}
		case llm.EvToolCallEnd:
			delete(open, ev.ToolCallID)
		}
	}
	if len(open) != 0 {
		t.Errorf("calls left open: %v", open)
	}
}

// ---------------------------------------------------------------------------
// Golden replay: max_tokens stop reason.
// ---------------------------------------------------------------------------

func TestGoldenMaxTokens(t *testing.T) {
	srv, _ := serveGolden(t, "max-tokens", false)
	_, turn, err := drain(t, srv.adapter(t, false, false), baseRequest())
	if err != nil {
		t.Fatal(err)
	}
	if turn.Stop != llm.StopMaxTokens {
		t.Fatalf("stop = %s, want max_tokens (hard requirement of SPEC-05 3.2)", turn.Stop)
	}
	if turn.Text != "This answer is cut mid-sentence because the token budget ran" {
		t.Errorf("partial text = %q", turn.Text)
	}
	if turn.Usage.OutputTokens != 6 {
		t.Errorf("usage = %+v", turn.Usage)
	}
}

// ---------------------------------------------------------------------------
// Golden replay: mid-stream disconnect keeps partial content, marks error.
// ---------------------------------------------------------------------------

func TestGoldenMidStreamDisconnect(t *testing.T) {
	srv, _ := serveGolden(t, "disconnect", false)
	events, turn, err := drain(t, srv.adapter(t, false, false), baseRequest())
	if err == nil {
		t.Fatal("expected a stream_disconnected error")
	}
	if classOf(err) != observe.ClassNetwork {
		t.Errorf("class = %s, want network", classOf(err))
	}
	var e *observe.Error
	if !errorsAs(err, &e) || e.ProviderCode != llm.CodeStreamDisconnect {
		t.Errorf("provider code = %v, want stream_disconnected", e)
	}
	if turn.Text != "partial text kept" {
		t.Errorf("partial content lost: %q", turn.Text)
	}
	if !turn.Incomplete {
		t.Error("disconnect turn must be marked incomplete")
	}
	// Terminal triple present in the stream.
	last := events[len(events)-1]
	if last.Type != llm.EvDone {
		t.Errorf("last event = %s, want done", last.Type)
	}
}

// ---------------------------------------------------------------------------
// Golden replay: usage aggregation across a multi-chunk stream.
// ---------------------------------------------------------------------------

func TestGoldenUsageAggregation(t *testing.T) {
	srv, _ := serveGolden(t, "usage-multichunk", false)
	_, turn, err := drain(t, srv.adapter(t, false, false), baseRequest())
	if err != nil {
		t.Fatal(err)
	}
	if turn.Usage != (llm.Usage{InputTokens: 10, OutputTokens: 7, CachedTokens: 0}) {
		t.Errorf("aggregated usage = %+v, want {10 7 0} (max-merge of 10/4 and 10/7)", turn.Usage)
	}
}

// ---------------------------------------------------------------------------
// Golden replay: backoff ladder over recorded 429s (retry-after honored).
// ---------------------------------------------------------------------------

func TestGoldenBackoffLadder429(t *testing.T) {
	srv, rep := serveGolden(t, "backoff-429", false)
	inner := srv.adapter(t, false, false)
	rp := llm.NewRetrying(inner, llm.RetryOptions{Max: 3, Base: 5 * time.Millisecond})

	start := time.Now()
	_, turn, err := drain(t, rp, baseRequest())
	elapsed := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	if turn.Stop != llm.StopEndTurn || turn.Text != "ok after retries" {
		t.Errorf("turn = %+v", turn)
	}
	// 3 recorded responses consumed: 429, 429, 200.
	if n := len(rep.Requests); n != 3 {
		t.Errorf("requests = %d, want 3", n)
	}
	// retry-after:1s x2 honored verbatim (well above the 5ms backoff base).
	if elapsed < 1900*time.Millisecond {
		t.Errorf("elapsed = %v, want >= 2s (retry-after honored)", elapsed)
	}
}

// ---------------------------------------------------------------------------
// Golden replay: 5xx exhaustion -> provider error (fast backoff injection).
// ---------------------------------------------------------------------------

func TestGoldenProvider500Exhausted(t *testing.T) {
	srv, rep := serveGolden(t, "provider-500", false)
	inner := srv.adapter(t, false, false)
	rp := llm.NewRetrying(inner, llm.RetryOptions{Max: 3, Base: 5 * time.Millisecond})

	_, turn, err := drain(t, rp, baseRequest())
	if err == nil {
		t.Fatal("expected provider error after the ladder")
	}
	if classOf(err) != observe.ClassProvider {
		t.Errorf("class = %s, want provider", classOf(err))
	}
	if len(rep.Requests) != 4 {
		t.Errorf("requests = %d, want 4 (initial + 3 retries)", len(rep.Requests))
	}
	if turn.Err == nil || turn.Stop != llm.StopError {
		t.Errorf("terminal events wrong: stop=%s err=%v", turn.Stop, turn.Err)
	}
}

// ---------------------------------------------------------------------------
// Golden replay: 401 -> no retry, Unconfigured semantics.
// ---------------------------------------------------------------------------

func TestGoldenUnauthorizedNoRetry(t *testing.T) {
	srv, rep := serveGolden(t, "unauthorized", false)
	inner := srv.adapter(t, false, false)
	rp := llm.NewRetrying(inner, llm.RetryOptions{Max: 3, Base: 5 * time.Millisecond})

	_, turn, err := drain(t, rp, baseRequest())
	if err == nil {
		t.Fatal("expected auth error")
	}
	if classOf(err) != observe.ClassAuth {
		t.Errorf("class = %s, want auth (Unconfigured semantics)", classOf(err))
	}
	if len(rep.Requests) != 1 {
		t.Errorf("requests = %d, want 1 (401 never retried)", len(rep.Requests))
	}
	if turn.Incomplete {
		t.Error("a pre-payload auth failure is not an incomplete stream")
	}
}

// ---------------------------------------------------------------------------
// Compat switches: missing usage / missing finish_reason / no [DONE] /
// fragments without index.
// ---------------------------------------------------------------------------

func TestGoldenMissingUsageStrictErrorsLooseTolerates(t *testing.T) {
	newSrv := func() *httptestServer {
		srv, _ := serveGolden(t, "missing-usage", false)
		return srv
	}

	_, _, err := drain(t, newSrv().adapter(t, false, false), baseRequest())
	if err == nil || classOf(err) != observe.ClassProvider {
		t.Errorf("strict mode: err = %v, want provider", err)
	}

	_, turn, err := drain(t, newSrv().adapter(t, true, false), baseRequest())
	if err != nil {
		t.Fatalf("loose mode: %v", err)
	}
	if turn.Stop != llm.StopEndTurn || turn.Text != "hi" {
		t.Errorf("loose turn = %+v", turn)
	}
	if turn.Usage != (llm.Usage{}) {
		t.Errorf("loose turn usage = %+v, want zero (no usage event)", turn.Usage)
	}

	// allow_missing_usage (without loose) also tolerates.
	if _, _, err := drain(t, newSrv().adapter(t, false, true), baseRequest()); err != nil {
		t.Errorf("allow_missing_usage mode: %v", err)
	}
}

func TestGoldenMissingFinishReasonStrictErrorsLooseTolerates(t *testing.T) {
	newSrv := func() *httptestServer {
		srv, _ := serveGolden(t, "missing-finish", false)
		return srv
	}
	if _, _, err := drain(t, newSrv().adapter(t, false, false), baseRequest()); err == nil {
		t.Error("strict mode accepted a stream without finish_reason")
	}
	_, turn, err := drain(t, newSrv().adapter(t, true, false), baseRequest())
	if err != nil {
		t.Fatal(err)
	}
	if turn.Stop != llm.StopEndTurn || turn.Usage.OutputTokens != 1 {
		t.Errorf("loose turn = %+v (usage must still arrive)", turn)
	}
}

func TestGoldenNoDoneSentinel(t *testing.T) {
	// A full stream that ends without [DONE]: strict -> network error,
	// loose (finish_reason seen) -> tolerated as complete.
	src := "data: {\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"\"},\"finish_reason\":null}]}\n\n" +
		"data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"done-ish\"},\"finish_reason\":null}]}\n\n" +
		"data: {\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n" +
		"data: {\"choices\":[],\"usage\":{\"prompt_tokens\":3,\"completion_tokens\":2,\"total_tokens\":5}}\n\n"
	rs, err := golden.Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	mkSrv := func() *httptestServer {
		rep := golden.NewReplayer(rs)
		srv := rep.Server()
		t.Cleanup(srv.Close)
		return &httptestServer{srv}
	}

	s1 := mkSrv()
	if _, _, err := drain(t, s1.adapter(t, false, false), baseRequest()); err == nil ||
		classOf(err) != observe.ClassNetwork {
		t.Errorf("strict: err = %v, want network (missing [DONE])", err)
	}

	s2 := mkSrv()
	_, turn, err := drain(t, s2.adapter(t, true, false), baseRequest())
	if err != nil {
		t.Fatalf("loose: %v", err)
	}
	if turn.Stop != llm.StopEndTurn || turn.Text != "done-ish" {
		t.Errorf("loose turn = %+v", turn)
	}
}

func TestGoldenMalformedToolIndexStrictVsLoose(t *testing.T) {
	src := "data: {\"choices\":[{\"index\":0,\"delta\":{\"tool_calls\":[{\"id\":\"c1\",\"type\":\"function\",\"function\":{\"name\":\"f\",\"arguments\":\"{}\"}}]},\"finish_reason\":null}]}\n\n" +
		"data: {\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"tool_calls\"}]}\n\n" +
		"data: {\"choices\":[],\"usage\":{\"prompt_tokens\":3,\"completion_tokens\":2,\"total_tokens\":5}}\n\n" +
		"data: [DONE]\n"
	rs, err := golden.Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	mkSrv := func() *httptestServer {
		rep := golden.NewReplayer(rs)
		srv := rep.Server()
		t.Cleanup(srv.Close)
		return &httptestServer{srv}
	}

	if _, _, err := drain(t, mkSrv().adapter(t, false, false), baseRequest()); err == nil {
		t.Error("strict mode accepted tool_calls without index")
	}
	_, turn, err := drain(t, mkSrv().adapter(t, true, false), baseRequest())
	if err != nil {
		t.Fatal(err)
	}
	if len(turn.ToolCalls) != 1 || turn.ToolCalls[0].Name != "f" {
		t.Errorf("loose tool calls = %+v", turn.ToolCalls)
	}
}

// ---------------------------------------------------------------------------
// Cancellation (paced golden fixture).
// ---------------------------------------------------------------------------

func TestGoldenCancellationMidStream(t *testing.T) {
	srv, _ := serveGolden(t, "long-text", true) // 40ms per chunk
	a := srv.adapter(t, false, false)

	ctx, cancel := context.WithCancel(context.Background())
	var events []llm.StreamEvent
	var stop llm.StreamEvent
	afterFirst := false
	err := a.Stream(ctx, baseRequest(), func(ev llm.StreamEvent) error {
		events = append(events, ev)
		if ev.Type == llm.EvTextDelta && !afterFirst {
			afterFirst = true
			cancel() // cancel after the first text delta
		}
		if ev.Type == llm.EvStop || ev.Type == llm.EvDone {
			stop = ev
		}
		return nil
	})
	if err != nil {
		t.Fatalf("cancellation is a clean end, got: %v", err)
	}
	if stop.Type != llm.EvDone {
		t.Fatalf("last event = %s, want done", stop.Type)
	}
	// A Stop{cancelled} must precede Done.
	sawCancelled := false
	for _, ev := range events {
		if ev.Type == llm.EvStop && ev.Stop == llm.StopCancelled {
			sawCancelled = true
		}
	}
	if !sawCancelled {
		t.Errorf("no Stop{cancelled} in %d events", len(events))
	}
	if len(events) >= 32 {
		t.Errorf("stream ran to completion despite cancellation (%d events)", len(events))
	}
}

// ---------------------------------------------------------------------------
// Proxy: HTTP(S)_PROXY honored (plain http target through a recording proxy).
// ---------------------------------------------------------------------------

func TestHTTProxyEnvHonored(t *testing.T) {
	var viaProxy atomic.Int32
	proxy := newTestServerHandler(t, func(w http.ResponseWriter, r *http.Request) {
		viaProxy.Add(1)
		// The proxy answers as the LLM (plain http forward behavior).
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"via-proxy\"},\"finish_reason\":null}]}\n\n" +
			"data: {\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n" +
			"data: {\"choices\":[],\"usage\":{\"prompt_tokens\":1,\"completion_tokens\":1,\"total_tokens\":2}}\n\n" +
			"data: [DONE]\n\n"))
	})
	defer proxy.Close()

	t.Setenv("HTTP_PROXY", proxy.URL())
	t.Setenv("HTTPS_PROXY", "")
	t.Setenv("NO_PROXY", "")

	opts := endpointOptionsWithClient(proxy.URL(), false, false, envProxyClient())
	a := New(opts)
	_, turn, err := drain(t, a, baseRequest())
	if err != nil {
		t.Fatal(err)
	}
	if turn.Text != "via-proxy" {
		t.Errorf("text = %q", turn.Text)
	}
	if viaProxy.Load() != 1 {
		t.Errorf("proxy saw %d requests, want 1 (HTTP_PROXY must be honored)", viaProxy.Load())
	}
}

// ---------------------------------------------------------------------------
// TLS: untrusted certificate is a DISTINCT code from unreachable (D42#4).
// ---------------------------------------------------------------------------

func TestTLSUntrustedCertDistinctFromUnreachable(t *testing.T) {
	// 1. Self-signed server + default roots -> tls_untrusted_cert.
	tlsSrv := newTLSServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	})
	defer tlsSrv.Close()

	a := New(endpointOptionsFor(tlsSrv.URL(), false, false))
	_, _, err := drain(t, a, baseRequest())
	if err == nil {
		t.Fatal("expected a TLS error")
	}
	var oe *observe.Error
	if !errorsAs(err, &oe) {
		t.Fatalf("unclassified TLS error: %v", err)
	}
	if oe.Class != observe.ClassNetwork || oe.ProviderCode != llm.CodeTLSUntrusted {
		t.Errorf("TLS error = %+v, want network/tls_untrusted_cert", oe)
	}

	// 2. Closed port -> connect_failed.
	dead := New(endpointOptionsFor("http://127.0.0.1:1/v1", false, false))
	_, _, err = drain(t, dead, baseRequest())
	if err == nil {
		t.Fatal("expected a connection error")
	}
	var oe2 *observe.Error
	if !errorsAs(err, &oe2) || oe2.ProviderCode != llm.CodeConnectFailed {
		t.Errorf("unreachable = %+v, want network/connect_failed", oe2)
	}

	// The two causes must classify DIFFERENTLY.
	if oe.ProviderCode == oe2.ProviderCode {
		t.Error("untrusted cert and unreachable must not collapse to one code (D42#4)")
	}
}

// ---------------------------------------------------------------------------
// Info(): cache capability bits for the seam.
// ---------------------------------------------------------------------------

func TestInfoCacheCapabilityBits(t *testing.T) {
	a := New(endpointOptionsFor("http://127.0.0.1:1/v1", false, false))
	info := a.Info()
	if info.Protocol != "openai-chat" || info.Model != "mock-small" {
		t.Errorf("info = %+v", info)
	}
	// OpenAI-compatible caching is implicit-prefix only: no explicit
	// breakpoints (that bit belongs to the anthropic adapter, ticket 11).
	if info.Cache.Mode != llm.CacheImplicit || info.Cache.Breakpoints {
		t.Errorf("cache bits = %+v, want implicit/no-breakpoints", info.Cache)
	}
}

// ensure the golden fixtures on disk match what the tests expect (guards
// against fixture drift between the two runners).
func TestGoldenFixturesExist(t *testing.T) {
	for _, name := range []string{"tool-call", "max-tokens", "disconnect",
		"usage-multichunk", "backoff-429", "provider-500", "unauthorized",
		"missing-usage", "missing-finish", "long-text"} {
		if _, err := os.Stat(filepath.Join(goldenDir, name+".sse")); err != nil {
			t.Errorf("fixture %s missing: %v", name, err)
		}
	}
	if !strings.Contains(goldenDir, "testdata") {
		t.Error("golden dir must stay under testdata")
	}
}
