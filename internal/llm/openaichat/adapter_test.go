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

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/observe"
)

// The protocol-neutral golden assertions that used to live in this file
// (tool-call assembly, max_tokens, mid-stream disconnect, usage aggregation,
// the 429/5xx/401 ladders, missing-usage/missing-finish tolerance and
// cancellation) were MOVED VERBATIM IN MEANING into the shared harness table:
// see harness_golden_test.go + internal/llm/adaptertest (they now run for all
// three adapters, not just this one). What stays here is openai-chat wire
// trivia with no equivalent in the other two protocols.

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

// ---------------------------------------------------------------------------
// Golden replay: max_tokens stop reason.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Golden replay: mid-stream disconnect keeps partial content, marks error.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Golden replay: usage aggregation across a multi-chunk stream.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Golden replay: backoff ladder over recorded 429s (retry-after honored).
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Golden replay: 5xx exhaustion -> provider error (fast backoff injection).
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Golden replay: 401 -> no retry, Unconfigured semantics.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Compat switches: missing usage / missing finish_reason / no [DONE] /
// fragments without index.
// ---------------------------------------------------------------------------

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
