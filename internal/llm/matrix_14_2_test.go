package llm_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/llm/anthropic"
	"github.com/CarlosShao/wisp/internal/llm/openaichat"
	"github.com/CarlosShao/wisp/internal/llm/openairesponses"
	"github.com/CarlosShao/wisp/internal/observe"
)

// The 14.2 failure-semantics matrix, ONE TEST PER ROW, asserting the emitted
// user-visible event/state - not a summary claim. Each row names the plan row
// it covers so a reviewer can diff the table against the tests.

// scripted is a fake LlmProvider whose behavior a row test controls exactly.
type scripted struct {
	info     llm.ProviderInfo
	attempts *int
	// onCall returns the events of one attempt plus its terminal error.
	onCall func(attempt int) ([]llm.StreamEvent, error)
}

func (s *scripted) Info() llm.ProviderInfo { return s.info }

func (s *scripted) Stream(ctx context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	*s.attempts++
	events, err := s.onCall(*s.attempts)
	for _, ev := range events {
		if e := emit(ev); e != nil {
			return e
		}
	}
	return err
}

func failureEvents(class observe.ErrorClass, code string) []llm.StreamEvent {
	e := observe.New(class, "scripted failure")
	if code != "" {
		e.ProviderCode = code
	}
	return []llm.StreamEvent{
		{Type: llm.EvError, Err: e},
		{Type: llm.EvStop, Stop: llm.StopError},
		{Type: llm.EvDone},
	}
}

// fakeClock is the MONOTONIC time source the notice row drives (D42#9: never a
// wall-clock difference against a deadline - the clock moves only when the
// test says so, and every value comes from the same source).
type fakeClock struct{ t time.Time }

func (c *fakeClock) Now() time.Time { now := c.t; return now }
func (c *fakeClock) advance(d time.Duration) {
	c.t = c.t.Add(d)
}

// Row 1: 429/5xx/网络抖动 - exponential backoff <= 3 retries OUTSIDE the loop;
// the ball stays Thinking; past 5s the attempt number becomes visible.
func TestMatrixRow1RetryStaysThinkingAndNoticesAfterFiveSeconds(t *testing.T) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	attempts := 0
	var notices []llm.RetryNotice
	slept := []time.Duration{}

	inner := &scripted{
		info: llm.ProviderInfo{Protocol: "scripted"}, attempts: &attempts,
		onCall: func(n int) ([]llm.StreamEvent, error) {
			if n < 3 {
				return failureEvents(observe.ClassProvider, "500"),
					observe.New(observe.ClassProvider, "upstream 500")
			}
			return []llm.StreamEvent{
				{Type: llm.EvTextDelta, Text: "recovered"},
				{Type: llm.EvStop, Stop: llm.StopEndTurn},
				{Type: llm.EvDone},
			}, nil
		},
	}

	rp := llm.NewRetrying(inner, llm.RetryOptions{
		Max: 3, Base: 5 * time.Second, Now: clock.Now, NoticeAfter: 5 * time.Second,
		Sleep: func(ctx context.Context, d time.Duration) error {
			slept = append(slept, d)
			clock.advance(d) // the wait really happens, on the same monotonic source
			return nil
		},
		Notice: func(n llm.RetryNotice) { notices = append(notices, n) },
	})

	events, turn, err := adaptertest.Drain(t, rp, adaptertest.BaseRequest())
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3 (initial + 2 retries)", attempts)
	}
	if len(slept) != 2 || slept[0] != 5*time.Second || slept[1] != 10*time.Second {
		t.Errorf("backoff ladder = %v, want [5s 10s] (exponential)", slept)
	}
	// Under the 5s threshold no notice may fire: the first retry decision is
	// taken at t=0 (nothing spent yet), so the ball just stays Thinking; the
	// second decision at t=5s crosses the threshold and becomes visible.
	for _, n := range notices {
		if n.Elapsed < 5*time.Second {
			t.Errorf("notice fired before the 5s threshold: %+v", n)
		}
	}
	if len(notices) != 1 {
		t.Fatalf("notices = %+v, want exactly one (the >5s case)", notices)
	}
	n := notices[0]
	if n.Elapsed != 5*time.Second {
		t.Errorf("notice elapsed = %v, want exactly the 5s threshold", n.Elapsed)
	}
	if n.Label() != "重试中 (2/3)" {
		t.Errorf("label = %q, want 重试中 (2/3)", n.Label())
	}
	if n.BallStateName() != "Thinking" {
		t.Errorf("ball state during retry = %q, want Thinking (no transition)", n.BallStateName())
	}
	if n.Class != observe.ClassProvider {
		t.Errorf("notice class = %s", n.Class)
	}
	// A retry is NOT a user-visible failure: no Error event until the ladder
	// ends. This is the "保持 Thinking" half of the row.
	for _, ev := range events {
		if ev.Type == llm.EvError {
			t.Errorf("a retried failure leaked an Error event: %v", ev)
		}
	}
	if turn.Text != "recovered" || turn.Stop != llm.StopEndTurn {
		t.Errorf("turn = %+v", turn)
	}
}

// Row 2: retries exhausted -> explicit Error, never a silent "the model did
// not answer". Asserted as: a terminal Error event carrying a classified
// D37 error, Stop{error}, and a non-empty mapped state.
func TestMatrixRow2ExhaustedIsExplicitErrorNotSilence(t *testing.T) {
	attempts := 0
	inner := &scripted{attempts: &attempts, onCall: func(int) ([]llm.StreamEvent, error) {
		evs := failureEvents(observe.ClassProvider, "502")
		evs[0].Err.Detail = "HTTP 502 upstream exploded"
		return evs, evs[0].Err
	}}
	rp := llm.NewRetrying(inner, llm.RetryOptions{
		Max: 3, Base: time.Millisecond,
		Sleep: func(context.Context, time.Duration) error { return nil },
	})

	events, turn, err := adaptertest.Drain(t, rp, adaptertest.BaseRequest())
	if err == nil {
		t.Fatal("exhausted ladder must surface an error")
	}
	if adaptertest.ClassOf(err) != observe.ClassProvider {
		t.Errorf("class = %s, want provider", adaptertest.ClassOf(err))
	}
	if attempts != 4 {
		t.Errorf("attempts = %d, want 4 (initial + 3)", attempts)
	}
	if turn.Err == nil || turn.Stop != llm.StopError {
		t.Fatalf("terminal pair missing: %+v", turn)
	}
	states := turn.Err.Class.MappedStates()
	if len(states) == 0 || states[0] != "Error" {
		t.Errorf("mapped states = %v, want Error first (panel shows the reason)", states)
	}
	if turn.Err.Class.MessageKey() == "" {
		t.Error("the failure must carry a copy key so the panel can render reason + advice")
	}
	// Exactly one Error event reached the consumer (no silent downgrade, and
	// no duplicate complaints per swallowed attempt).
	n := 0
	for _, ev := range events {
		if ev.Type == llm.EvError {
			n++
		}
	}
	if n != 1 {
		t.Errorf("Error events = %d, want 1", n)
	}
}

// Row 3: 401/403 -> never retried, straight to Unconfigured + guidance.
func TestMatrixRow3AuthIsUnconfiguredAndNeverRetried(t *testing.T) {
	for _, status := range []int{401, 403} {
		proc := adaptertest.StartMockllm(t)
		proc.Reset(t)
		proc.QueueFault(t, adaptertest.Fault{Status: status, Times: 5})
		rp := llm.NewRetrying(mustProvider(t, "openai-chat", proc.Base+"/v1"),
			llm.RetryOptions{
				Max: 3, Base: time.Millisecond,
				Sleep: func(context.Context, time.Duration) error { return nil },
			})

		_, turn, err := adaptertest.Drain(t, rp, adaptertest.BaseRequest())
		if err == nil {
			t.Fatalf("status %d: expected an auth failure", status)
		}
		if got := adaptertest.ClassOf(err); got != observe.ClassAuth {
			t.Errorf("status %d: class = %s, want auth", status, got)
		}
		if got := observe.ClassAuth.MappedStates(); len(got) == 0 || got[0] != "Unconfigured" {
			t.Errorf("auth mapped states = %v, want Unconfigured first (guide to the config panel)", got)
		}
		if n := proc.RouteCount(t, "chat"); n != 1 {
			t.Errorf("status %d: requests = %d, want 1 (retrying a bad key can get us banned)", status, n)
		}
		if turn.Err == nil {
			t.Error("auth failure must emit an Error event so the ball can move")
		}
		proc.Reset(t)
	}
}

// Row 4: quota exhausted -> Error that DISTINGUISHES account trouble from
// network trouble (the two must never share a class, a state or a copy key).
func TestMatrixRow4QuotaIsAccountNotNetwork(t *testing.T) {
	proc := adaptertest.StartMockllm(t)
	proc.Reset(t)
	proc.QueueFault(t, adaptertest.Fault{
		Status: 429, Times: 1,
		Body: `{"error":{"message":"no credits","type":"insufficient_quota","code":"insufficient_quota"}}`,
	})
	_, _, quotaErr := adaptertest.Drain(t,
		mustProvider(t, "openai-chat", proc.Base+"/v1"), adaptertest.BaseRequest())
	if adaptertest.ClassOf(quotaErr) != observe.ClassBudget {
		t.Fatalf("quota class = %s, want budget", adaptertest.ClassOf(quotaErr))
	}

	// The network counterpart: an unreachable host.
	dead := mustProvider(t, "openai-chat", "http://127.0.0.1:1/v1")
	_, _, netErr := adaptertest.Drain(t, dead, adaptertest.BaseRequest())
	if adaptertest.ClassOf(netErr) != observe.ClassNetwork {
		t.Fatalf("network class = %s, want network", adaptertest.ClassOf(netErr))
	}

	qe, ne := adaptertest.ObserveErr(quotaErr), adaptertest.ObserveErr(netErr)
	if qe.Class == ne.Class {
		t.Fatal("quota and network failures collapsed into one class")
	}
	if qe.Class.MessageKey() == ne.Class.MessageKey() {
		t.Error("the account-vs-network copy keys must differ (14.2 row 4)")
	}
	if qe.Retryable() {
		t.Error("an exhausted account is never retried")
	}
	if !ne.Retryable() {
		t.Error("a network failure stays retryable (the distinction must be actionable)")
	}
}

// Row 5: mid-stream drop -> partial content retained and marked, unclosed tool
// calls reported as open (the agent core turns each into a failed result).
func TestMatrixRow5MidStreamDropKeepsPartialAndFailsOpenCalls(t *testing.T) {
	for _, pair := range []struct {
		name     string
		fixture  string
		protocol string
	}{
		{"openai-chat", "chat-disconnect-tools", "openai-chat"},
		{"anthropic", "anthropic-disconnect-tools", "anthropic"},
		{"openai-responses", "responses-disconnect-tools", "openai-responses"},
	} {
		t.Run(pair.name, func(t *testing.T) {
			base, _ := adaptertest.ServeFixtureByName(t, pair.fixture)
			events, turn, err := adaptertest.Drain(t,
				mustProvider(t, pair.protocol, base), adaptertest.BaseRequest())
			if err == nil {
				t.Fatal("expected a mid-stream disconnect")
			}
			if adaptertest.CodeOf(err) != llm.CodeStreamDisconnect {
				t.Errorf("provider_code = %q, want stream_disconnected", adaptertest.CodeOf(err))
			}
			if !turn.Incomplete || turn.Text == "" {
				t.Errorf("partial content must survive and be marked: %+v", turn)
			}
			// The seam never invents ToolCallEnd: the truncated call is reported
			// with Complete=false and HasOpenToolCalls() is true, which is what
			// the agent core turns into a failed tool result (SPEC-05 sec 3.4,
			// D21 failToolCallsFromTruncatedMessage).
			if len(turn.ToolCalls) == 0 {
				t.Fatal("the fixture must deliver a tool call before the drop")
			}
			if turn.ToolCalls[0].Complete {
				t.Error("a truncated tool call must be reported incomplete so the loop fails it")
			}
			rebuild := llm.NewTurnCollector()
			for _, ev := range events {
				rebuild.Observe(ev)
			}
			if !rebuild.HasOpenToolCalls() {
				t.Error("collector must report the unclosed tool call as open")
			}
			// The 「响应中断」 marker is decidable from seam state alone:
			// incomplete payload + the stream_disconnected provider code. The
			// ball/panel own rendering the copy (SPEC-05 sec 3.4 row 5).
			oe := adaptertest.ObserveErr(err)
			if oe == nil || !turn.Incomplete || oe.ProviderCode != llm.CodeStreamDisconnect {
				t.Errorf("interrupted-marker state not exposed: turn=%+v err=%+v", turn, oe)
			}
		})
	}
}

// mustProvider builds a provider through the C5 registry (the adapters
// self-register from their own package init, imported below for exactly that).
func mustProvider(t *testing.T, protocol, base string) llm.LlmProvider {
	t.Helper()
	var o llm.EndpointOptions
	o.Provider = "mock"
	o.Model = "mock-small"
	o.Protocol = protocol
	o.BaseURL = base
	o.APIKey = "sk-test-not-a-real-key"
	o.ContextWindow = 8192
	o.HTTPClient = noProxyClient()
	p, err := llm.NewProvider(o)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// noProxyClient keeps machine proxy settings out of these tests.
func noProxyClient() *http.Client {
	return &http.Client{Transport: &http.Transport{Proxy: nil}}
}

// Adapter packages are imported only for their protocol registration.
var (
	_ = (*anthropic.Adapter)(nil)
	_ = (*openaichat.Adapter)(nil)
	_ = (*openairesponses.Adapter)(nil)
)
