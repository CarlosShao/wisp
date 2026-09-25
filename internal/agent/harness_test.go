package agent

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	_ "github.com/CarlosShao/wisp/internal/llm/openaichat" // protocol registration
	"github.com/CarlosShao/wisp/internal/observe"
)

// goldenDir points at this package's golden fixtures: the same format and
// replayer as ticket 09, i.e. recorded SSE bytes served by a real HTTP server
// and parsed by the genuine openai-chat adapter. SPEC-05 §9 requires the loop
// tests to hit the C5 seam with real byte streams rather than function mocks.
var goldenDir = func() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(thisFile), "testdata", "golden")
}()

// harness is a Loop wired to a golden-backed C5 provider.
type harness struct {
	t      *testing.T
	rep    *golden.Replayer
	server *httptest.Server

	mu     sync.Mutex
	bodies [][]byte // raw request bodies, in arrival order

	reg   *observe.Registry
	loop  *Loop
	sink  *RecordSink
	tools ToolProvider
	prov  llm.LlmProvider
}

// harnessOpt tweaks the wiring before the Loop is built.
type harnessOpt func(*harness, *Options)

// withConfig sets loop config fields.
func withConfig(f func(*Config)) harnessOpt {
	return func(_ *harness, o *Options) { f(&o.Config) }
}

// withTools replaces the ToolProvider.
func withTools(tp ToolProvider) harnessOpt {
	return func(h *harness, o *Options) { h.tools = tp; o.Tools = tp }
}

// withJournal attaches a Journal implementation.
func withJournal(j Journal) harnessOpt {
	return func(_ *harness, o *Options) { o.Journal = j }
}

// withSummarizer attaches a history summarizer.
func withSummarizer(s Summarizer) harnessOpt {
	return func(_ *harness, o *Options) { o.Summarizer = s }
}

// withControl installs a control handler.
func withControl(f ControlHandler) harnessOpt {
	return func(_ *harness, o *Options) { o.Control = f }
}

// withRegistry installs a specific goroutine registry.
func withRegistry(r *observe.Registry) harnessOpt {
	return func(h *harness, o *Options) { h.reg = r; o.Registry = r }
}

// withLogger overrides the harness logger. The default a harness builds is a
// discard handler (so a normal run stays quiet); a test that asserts on a
// record the loop produces installs its own handler here, which is also what
// proves the loop hands that logger down to the components it wires.
func withLogger(lg *slog.Logger) harnessOpt {
	return func(_ *harness, o *Options) { o.Logger = lg }
}

// newHarness loads a golden fixture and serves its responses sequentially.
func newHarness(t *testing.T, fixture string, opts ...harnessOpt) *harness {
	t.Helper()
	rs, err := golden.LoadFile(filepath.Join(goldenDir, fixture+".sse"))
	if err != nil {
		t.Fatalf("load golden %s: %v", fixture, err)
	}
	rep := golden.NewReplayer(rs)

	h := &harness{t: t, rep: rep, sink: &RecordSink{}, reg: observe.NewRegistry()}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(body))
		h.mu.Lock()
		h.bodies = append(h.bodies, body)
		h.mu.Unlock()
		rep.ServeHTTP(w, r)
	}))
	t.Cleanup(srv.Close)
	h.server = srv

	o := Options{}
	o.Config = Config{
		Model: "mock-small", ContextWindow: 128000,
		ArtifactsDir:                t.TempDir(),
		PassThroughUnclassifiedRisk: true,
		SteeringEnabled:             true,
		Currency:                    "CNY",
		PerToolTimeout:              2 * time.Second,
	}
	o.Sink = h.sink
	o.Registry = h.reg
	o.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	for _, opt := range opts {
		opt(h, &o)
	}
	if o.Tools == nil {
		ep := NewEchoProvider()
		o.Tools = ep
		h.tools = ep
	}

	ep := llm.EndpointOptions{
		Provider:      "mock",
		Model:         o.Config.Model,
		Protocol:      "openai-chat",
		BaseURL:       srv.URL,
		ContextWindow: o.Config.ContextWindow,
		HTTPClient:    &http.Client{Transport: &http.Transport{Proxy: nil}},
	}
	prov, err := llm.NewProvider(ep)
	if err != nil {
		t.Fatalf("build provider: %v", err)
	}
	o.Provider = prov
	h.prov = prov

	loop, err := New(o)
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	h.loop = loop
	return h
}

// requests is how many HTTP requests the provider made - the "mockllm request
// count" the control-layer test asserts on. It counts the bodies the handler
// captured under h.mu rather than reading the replayer's internal slice: that
// slice is written on the httptest handler goroutine, and a cancelled task can
// abandon a stream before its handler has finished, which left an unsynchronised
// read (a data race under -race) if we touched rep.Requests directly.
func (h *harness) requests() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.bodies)
}

// requestBodies returns the recorded request bodies in arrival order.
func (h *harness) requestBodies() [][]byte {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([][]byte, len(h.bodies))
	copy(out, h.bodies)
	return out
}

// echo returns the EchoProvider under test.
func (h *harness) echo() *EchoProvider {
	ep, ok := h.tools.(*EchoProvider)
	if !ok {
		h.t.Fatalf("harness has %T, not *EchoProvider", h.tools)
	}
	return ep
}

// run executes one task synchronously on the test goroutine.
func (h *harness) run(text string) Result {
	h.t.Helper()
	return h.loop.Run(context.Background(), text)
}

// waitForFirstRequest blocks until the provider has been called at least once.
func (h *harness) waitForFirstRequest(within time.Duration) bool {
	deadline := observe.NewTimeout(within)
	for !deadline.Expired() {
		if len(h.requestBodies()) > 0 {
			return true
		}
		time.Sleep(2 * time.Millisecond)
	}
	return false
}

// ---------------------------------------------------------------------------
// shared shorthands for the C7 types (test readability only)

type (
	llmMessage    = llm.Message
	llmText       = llm.TextPart
	llmToolResult = llm.ToolResultPart
)

const llmRoleTool = llm.RoleTool

// toolCalls builds a turn's assembled calls (guard unit tests).
func toolCalls(id, name, args string) []llm.ToolCall {
	return []llm.ToolCall{{ID: id, Name: name, Args: []byte(args), Complete: true}}
}

// ---------------------------------------------------------------------------
// shared assertions

func mustContain(t *testing.T, label, hay, needle string) {
	t.Helper()
	if !stringsContains(hay, needle) {
		t.Fatalf("%s: missing %q in %q", label, needle, trunc(hay, 600))
	}
}

func stringsContains(hay, needle string) bool {
	return bytes.Contains([]byte(hay), []byte(needle))
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
