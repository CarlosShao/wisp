// Package adaptertest is the ONE golden/fault harness every C5 adapter runs
// through (ticket 11, AC#1: "both adapters pass the same golden/fault test
// suite as OpenAI Chat (identical harness)").
//
// Why a non-test package: Go test files cannot be imported, so a shared test
// table that three adapter packages actually execute has to live in an
// ordinary package the adapter tests call (the httptest precedent). Nothing
// outside *_test.go imports it, and it imports no adapter (that would be a
// cycle): each adapter supplies a Unit describing how to build itself and
// which fixture carries each scenario.
//
// The contract this enforces:
//
//   - CaseTable() is the single shared table of scenarios. Its expectations
//     are adapter-agnostic (same text, same tool calls, same usage numbers,
//     same stop reason, same error class) because the per-adapter golden
//     fixtures were authored to carry the SAME semantics in each wire
//     dialect. An adapter cannot pass by running fewer cases:
//     TestScenarioCoverage (in each adapter package) asserts every scenario in
//     the table has a fixture and every fixture is exercised.
//   - Fixture naming lives with the Unit: protocol dialect bytes differ,
//     decoded C6 events must not.
//   - The fault suite (Mockllm) drives the real mockllm subprocess through
//     /__control injection so retries, 429 retry-after, 401 no-retry and
//     mid-stream truncation are proven against a live server for every
//     adapter, not only against the in-process replayer.
package adaptertest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/observe"
)

// ---------------------------------------------------------------------------
// Unit: what an adapter package supplies so the shared table can run it.
// ---------------------------------------------------------------------------

// Options are the compat switches the shared cases toggle.
type Options struct {
	Loose             bool
	AllowMissingUsage bool
}

// Unit describes one adapter under test.
type Unit struct {
	// Name is the config protocol enum (label in subtests + registry checks).
	Name string
	// Prefix is the golden-file prefix this unit's scenarios carry (fixture
	// names are Prefix+scenario, except where Fixture overrides).
	Prefix string
	// EndpointPath is appended to the server base URL by the adapter itself;
	// the harness only logs it. Kept for documentation of the seam shape.
	EndpointPath string
	// New builds the adapter over base (the harness server URL /v1-style root).
	New func(t *testing.T, base string, o Options) llm.LlmProvider
	// Fixture maps a scenario to the golden file base name carrying that
	// scenario in this adapter's wire dialect. Returning "" fails the
	// coverage test; there is no "not applicable" escape hatch, because every
	// scenario in CaseTable is expressible in all three dialects.
	Fixture func(sc ScenarioID) string
	// KeyForHeader is the auth header scheme the adapter must send
	// ("openai-bearer" | "anthropic-header"); the request-capture case reads
	// the body only, so this stays documentation-grade.
	AuthScheme string
}

// ---------------------------------------------------------------------------
// Scenario table
// ---------------------------------------------------------------------------

// ScenarioID names one shared case.
type ScenarioID string

const (
	ScToolCall        ScenarioID = "tool-call"
	ScMaxTokens       ScenarioID = "max-tokens"
	ScDisconnect      ScenarioID = "disconnect"
	ScDisconnectTools ScenarioID = "disconnect-tools"
	ScUsageMultichunk ScenarioID = "usage-multichunk"
	ScThinking        ScenarioID = "thinking"
	ScBackoff429      ScenarioID = "backoff-429"
	ScProvider500     ScenarioID = "provider-500"
	ScUnauthorized    ScenarioID = "unauthorized"
	ScMissingUsage    ScenarioID = "missing-usage"
	ScNoTerminal      ScenarioID = "no-terminal"
	ScMidStreamError  ScenarioID = "mid-stream-error"
	ScCancel          ScenarioID = "cancel"
)

// AllScenarios lists every shared scenario (coverage guard iterates it).
func AllScenarios() []ScenarioID {
	return []ScenarioID{
		ScToolCall, ScMaxTokens, ScDisconnect, ScDisconnectTools,
		ScUsageMultichunk, ScThinking, ScBackoff429, ScProvider500,
		ScUnauthorized, ScMissingUsage, ScNoTerminal, ScMidStreamError,
		ScCancel,
	}
}

// Mode selects how Run executes a case.
type Mode uint8

const (
	// ModeStream: replay once through a strict adapter, compare the turn.
	ModeStream Mode = iota
	// ModeRetry: replay through the retry ladder (Base 5ms), compare the turn
	// and the number of requests the server saw.
	ModeRetry
	// ModeCompat: strict must fail, loose and allow_missing_usage must pass.
	ModeCompat
	// ModeCancel: paced fixture cancelled after the first text delta.
	ModeCancel
)

// Case is one row of the shared table.
type Case struct {
	ID   ScenarioID
	Mode Mode
	Desc string

	// Request builds the normalized request (nil = BaseRequest()).
	Request func() *llm.Request

	// Expectations for success cases (and for the loose legs of ModeCompat).
	WantStop       llm.StopReason
	WantText       string
	WantReasoning  string
	WantToolCalls  []llm.ToolCall
	WantUsage      llm.Usage
	WantIncomplete bool

	// Expectations for failure cases.
	WantClass  observe.ErrorClass
	WantCode   string // provider_code of the emitted/returned error
	WantEvents []llm.StreamEventType
	// WantOpenToolCalls is set by the truncation scenarios: an unclosed tool
	// call is the CORRECT outcome there (the agent core fails it, the seam
	// never invents ToolCallEnd), so the open-call guard stands down.
	WantOpenToolCalls bool
}

// BaseRequest is the protocol-agnostic request the golden fixtures answer.
func BaseRequest() *llm.Request {
	return &llm.Request{Model: "mock-small", Messages: []llm.Message{{
		Role:    llm.RoleUser,
		Content: []llm.Content{llm.TextPart{Text: "hi"}},
	}}}
}

// The two tool calls every dialect's tool-call fixture must yield, identical
// down to the ids and argument bytes: that is what makes "the same table"
// more than a slogan.
func weatherAndTimeCalls() []llm.ToolCall {
	return []llm.ToolCall{
		{ID: "call_a1", Name: "get_weather",
			Args: json.RawMessage(`{"city":"Zhuhai"}`), Complete: true},
		{ID: "call_b2", Name: "get_time",
			Args: json.RawMessage(`{}`), Complete: true},
	}
}

// CaseTable is THE shared table. Every adapter runs every row.
func CaseTable() []Case {
	return []Case{
		{
			ID: ScToolCall, Mode: ModeStream,
			Desc:          "two fragmented tool calls + text + usage with cache reads",
			WantStop:      llm.StopToolUse,
			WantText:      "Let me check.",
			WantToolCalls: weatherAndTimeCalls(),
			WantUsage:     llm.Usage{InputTokens: 21, OutputTokens: 9, CachedTokens: 8},
		},
		{
			ID: ScMaxTokens, Mode: ModeStream,
			Desc:      "truncation must surface stop=max_tokens (SPEC-05 3.2 hard requirement)",
			WantStop:  llm.StopMaxTokens,
			WantText:  "This answer is cut mid-sentence because the token budget ran",
			WantUsage: llm.Usage{InputTokens: 10, OutputTokens: 6},
		},
		{
			ID: ScDisconnect, Mode: ModeStream,
			Desc:           "mid-stream drop keeps partial content, marks it incomplete",
			WantText:       "partial text kept",
			WantClass:      observe.ClassNetwork,
			WantCode:       llm.CodeStreamDisconnect,
			WantEvents:     []llm.StreamEventType{llm.EvTextDelta, llm.EvError, llm.EvStop, llm.EvDone},
			WantIncomplete: true,
		},
		{
			ID: ScDisconnectTools, Mode: ModeStream,
			Desc:           "drop with an unclosed tool call: partial text kept, call stays open (14.2 row 5)",
			WantText:       "before the tool call",
			WantClass:      observe.ClassNetwork,
			WantCode:       llm.CodeStreamDisconnect,
			WantEvents:     []llm.StreamEventType{llm.EvTextDelta, llm.EvToolCallStart, llm.EvToolCallArgsDelta, llm.EvError, llm.EvStop, llm.EvDone},
			WantIncomplete: true, WantOpenToolCalls: true,
		},
		{
			ID: ScUsageMultichunk, Mode: ModeStream,
			Desc:      "usage reported twice aggregates by max-merge",
			WantStop:  llm.StopEndTurn,
			WantText:  "Hello there",
			WantUsage: llm.Usage{InputTokens: 10, OutputTokens: 7},
		},
		{
			ID: ScThinking, Mode: ModeStream,
			Desc:          "extended thinking / reasoning items surface as ReasoningDelta (AC#3)",
			WantStop:      llm.StopEndTurn,
			WantText:      "four",
			WantReasoning: "Counting the dots.",
			WantUsage:     llm.Usage{InputTokens: 12, OutputTokens: 8, CachedTokens: 0},
		},
		{
			ID: ScBackoff429, Mode: ModeRetry,
			Desc:      "429 x2 with retry-after honored verbatim, then success",
			WantStop:  llm.StopEndTurn,
			WantText:  "ok after retries",
			WantUsage: llm.Usage{InputTokens: 5, OutputTokens: 3},
		},
		{
			ID: ScProvider500, Mode: ModeRetry,
			Desc:       "5xx exhausts the ladder into a provider-class terminal error",
			WantClass:  observe.ClassProvider,
			WantEvents: []llm.StreamEventType{llm.EvError, llm.EvStop, llm.EvDone},
		},
		{
			ID: ScUnauthorized, Mode: ModeRetry,
			Desc:       "401 is never retried and maps to Unconfigured semantics",
			WantClass:  observe.ClassAuth,
			WantEvents: []llm.StreamEventType{llm.EvError, llm.EvStop, llm.EvDone},
		},
		{
			ID: ScMissingUsage, Mode: ModeCompat,
			Desc:      "missing usage block: strict errors, compat switches tolerate",
			WantClass: observe.ClassProvider,
			WantStop:  llm.StopEndTurn,
			WantText:  "hi",
		},
		{
			ID: ScNoTerminal, Mode: ModeCompat,
			Desc:      "stream without its terminal marker: strict errors, loose completes",
			WantClass: observe.ClassProvider,
			WantStop:  llm.StopEndTurn,
			WantText:  "hi",
		},
		{
			ID: ScMidStreamError, Mode: ModeStream,
			Desc:       "an error payload inside a 200 stream ends it with Error{provider}",
			WantClass:  observe.ClassProvider,
			WantEvents: []llm.StreamEventType{llm.EvTextDelta, llm.EvError, llm.EvStop, llm.EvDone},
			WantText:   "up to the error", WantIncomplete: true,
		},
		{
			ID: ScCancel, Mode: ModeCancel,
			Desc: "cancellation is a clean end: Stop{cancelled}, no Error event",
		},
	}
}

// AssertCoverage is the anti-shortcut guard every adapter package calls: the
// unit must have a real, parseable fixture for EVERY scenario in the shared
// table (no skipped rows), and every fixture carrying this unit's prefix must
// be claimed by a scenario (an orphan fixture is fake coverage).
func AssertCoverage(t *testing.T, u Unit) {
	t.Helper()
	table := CaseTable()
	if len(table) != len(AllScenarios()) {
		t.Fatalf("CaseTable has %d rows but AllScenarios lists %d", len(table), len(AllScenarios()))
	}
	claimed := map[string]bool{}
	for _, c := range table {
		name := u.Fixture(c.ID)
		if name == "" {
			t.Fatalf("scenario %s has no fixture for unit %s", c.ID, u.Name)
		}
		claimed[name] = true
		rs, err := golden.LoadFile(filepath.Join(GoldenDir(), name+".sse"))
		if err != nil {
			t.Fatalf("scenario %s: fixture %s: %v", c.ID, name, err)
		}
		if len(rs) == 0 {
			t.Fatalf("scenario %s: fixture %s parses to zero responses", c.ID, name)
		}
	}
	entries, err := os.ReadDir(GoldenDir())
	if err != nil {
		t.Fatal(err)
	}
	prefix := u.Prefix
	if prefix == "" {
		prefix = u.Name
	}
	for _, e := range entries {
		n := e.Name()
		if !strings.HasSuffix(n, ".sse") {
			continue
		}
		base := strings.TrimSuffix(n, ".sse")
		if !strings.HasPrefix(base, prefix) {
			continue
		}
		if !claimed[base] {
			t.Errorf("orphan fixture %s: no scenario of the shared table maps to it", base)
		}
	}
}

// ---------------------------------------------------------------------------
// Fixture plumbing
// ---------------------------------------------------------------------------

// GoldenDir is the shared fixture directory: the SAME files tools/mockllm
// replays over HTTP (one format, two runners).
func GoldenDir() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "testdata", "golden")
}

// RepoRoot resolves the module root from the fixture directory
// (<root>/internal/llm/testdata/golden).
func RepoRoot() string {
	return filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(GoldenDir()))))
}

// Serve starts the in-process replayer for one unit's fixture of a scenario.
func Serve(t *testing.T, u Unit, sc ScenarioID, pace bool) (base string, rep *golden.Replayer) {
	t.Helper()
	name := u.Fixture(sc)
	if name == "" {
		t.Fatalf("adaptertest: unit %s has no fixture for scenario %s", u.Name, sc)
	}
	rs, err := golden.LoadFile(filepath.Join(GoldenDir(), name+".sse"))
	if err != nil {
		t.Fatalf("load golden %s: %v", name, err)
	}
	rep = golden.NewReplayer(rs)
	rep.Pace = pace
	srv := rep.Server()
	t.Cleanup(srv.Close)
	return srv.URL, rep
}

// ServeBody starts the in-process replayer over inline golden source.
func ServeBody(t *testing.T, src string) (base string, rep *golden.Replayer) {
	t.Helper()
	rs, err := golden.Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	rep = golden.NewReplayer(rs)
	srv := rep.Server()
	t.Cleanup(srv.Close)
	return srv.URL, rep
}

// ---------------------------------------------------------------------------
// Drain helpers
// ---------------------------------------------------------------------------

// Drain streams a request and returns the events, the assembled turn and the
// error.
func Drain(t *testing.T, p llm.LlmProvider, req *llm.Request) ([]llm.StreamEvent, llm.TurnResult, error) {
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

// ClassOf extracts the D37 class of an error chain.
func ClassOf(err error) observe.ErrorClass {
	c, _ := observe.ClassOf(err)
	return c
}

// CodeOf extracts the provider_code of an *observe.Error in the chain ("" when
// there is none).
func CodeOf(err error) string {
	var e *observe.Error
	if as(err, &e) {
		return e.ProviderCode
	}
	return ""
}

func as(err error, target **observe.Error) bool {
	for err != nil {
		if e, ok := err.(*observe.Error); ok {
			*target = e
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

func eqToolCalls(got, want []llm.ToolCall) string {
	if len(got) != len(want) {
		return fmt.Sprintf("tool calls = %d, want %d (%+v)", len(got), len(want), got)
	}
	for i := range got {
		if got[i].ID != want[i].ID || got[i].Name != want[i].Name ||
			string(got[i].Args) != string(want[i].Args) || got[i].Complete != want[i].Complete {
			return fmt.Sprintf("tool call %d = %+v, want %+v", i, got[i], want[i])
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Run: the shared golden suite
// ---------------------------------------------------------------------------

// Run executes one case of the shared table against one adapter.
func Run(t *testing.T, u Unit, c Case) {
	t.Helper()
	switch c.Mode {
	case ModeStream:
		runStream(t, u, c)
	case ModeRetry:
		runRetry(t, u, c)
	case ModeCompat:
		runCompat(t, u, c)
	case ModeCancel:
		runCancel(t, u, c)
	default:
		t.Fatalf("adaptertest: case %s has unknown mode %d", c.ID, c.Mode)
	}
}

// RunAll runs the whole shared table for one unit.
func RunAll(t *testing.T, u Unit) {
	t.Helper()
	for _, c := range CaseTable() {
		c := c
		t.Run(string(c.ID), func(t *testing.T) {
			Run(t, u, c)
		})
	}
}

func requestFor(c Case) *llm.Request {
	if c.Request != nil {
		return c.Request()
	}
	return BaseRequest()
}

func runStream(t *testing.T, u Unit, c Case) {
	t.Helper()
	base, _ := Serve(t, u, c.ID, false)
	p := u.New(t, base, Options{})
	events, turn, err := Drain(t, p, requestFor(c))
	checkTerminalContract(t, c, events, turn, err)
	checkTurn(t, c, turn)
}

func runRetry(t *testing.T, u Unit, c Case) {
	t.Helper()
	base, rep := Serve(t, u, c.ID, false)
	inner := u.New(t, base, Options{})
	rp := llm.NewRetrying(inner, RetryOptionsOf())

	start := time.Now()
	events, turn, err := Drain(t, rp, requestFor(c))
	elapsed := time.Since(start)

	checkTerminalContract(t, c, events, turn, err)
	if c.WantClass == "" {
		checkTurn(t, c, turn)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	// Request counts prove the ladder actually ran (and that 401 never did).
	switch c.ID {
	case ScBackoff429:
		if n := len(rep.Requests); n != 3 {
			t.Errorf("requests = %d, want 3 (429,429,200)", n)
		}
		// Measured threshold: the two retry-after:1 hints sum to 2s; 1.9s
		// leaves jitter margin while proving they were honored verbatim (a
		// non-honoring backoff would finish in ~10ms with the 5ms base).
		if elapsed < 1900*time.Millisecond {
			t.Errorf("elapsed = %v, want >= 1.9s (retry-after honored)", elapsed)
		}
	case ScProvider500:
		if n := len(rep.Requests); n != 4 {
			t.Errorf("requests = %d, want 4 (initial + 3 retries)", n)
		}
	case ScUnauthorized:
		if n := len(rep.Requests); n != 1 {
			t.Errorf("requests = %d, want 1 (401 never retried)", n)
		}
	}
}

func runCompat(t *testing.T, u Unit, c Case) {
	t.Helper()

	// Strict: must fail with the expected class.
	strictBase, _ := Serve(t, u, c.ID, false)
	_, turn, err := Drain(t, u.New(t, strictBase, Options{}), requestFor(c))
	if err == nil || ClassOf(err) != c.WantClass {
		t.Errorf("strict: err = %v, want class %s", err, c.WantClass)
	}
	if turn.Err == nil {
		t.Error("strict: failure stream carries no Error event payload")
	}

	// Loose: must succeed with the shared expectations.
	looseBase, _ := Serve(t, u, c.ID, false)
	_, lturn, lerr := Drain(t, u.New(t, looseBase, Options{Loose: true}), requestFor(c))
	if lerr != nil {
		t.Fatalf("loose: %v", lerr)
	}
	if lturn.Stop != c.WantStop || lturn.Text != c.WantText {
		t.Errorf("loose turn = %s %q, want %s %q", lturn.Stop, lturn.Text, c.WantStop, c.WantText)
	}

	// allow_missing_usage alone also tolerates a missing usage block.
	if c.ID == ScMissingUsage {
		amuBase, _ := Serve(t, u, c.ID, false)
		if _, _, err := Drain(t, u.New(t, amuBase, Options{AllowMissingUsage: true}), requestFor(c)); err != nil {
			t.Errorf("allow_missing_usage: %v", err)
		}
	}
}

func runCancel(t *testing.T, u Unit, c Case) {
	t.Helper()
	base, _ := Serve(t, u, c.ID, true) // paced fixture
	p := u.New(t, base, Options{})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	var events []llm.StreamEvent
	var last llm.StreamEvent
	afterFirst := false
	err := p.Stream(ctx, requestFor(c), func(ev llm.StreamEvent) error {
		events = append(events, ev)
		if ev.Type == llm.EvTextDelta && !afterFirst {
			afterFirst = true
			cancel()
		}
		if ev.Type == llm.EvStop || ev.Type == llm.EvDone {
			last = ev
		}
		return nil
	})
	if err != nil {
		t.Fatalf("cancellation is a clean end, got: %v", err)
	}
	if last.Type != llm.EvDone {
		t.Fatalf("last event = %s, want done", last.Type)
	}
	sawCancelled, sawError := false, false
	for _, ev := range events {
		switch {
		case ev.Type == llm.EvStop && ev.Stop == llm.StopCancelled:
			sawCancelled = true
		case ev.Type == llm.EvError:
			sawError = true
		}
	}
	if !sawCancelled {
		t.Errorf("no Stop{cancelled} in %d events", len(events))
	}
	if sawError {
		t.Errorf("cancellation stream contains an Error event: %v", events)
	}
	if len(events) >= 32 {
		t.Errorf("stream ran to completion despite cancellation (%d events)", len(events))
	}
}

// checkTurn asserts the shared success expectations. Failure cases only carry
// the partial-content and incomplete markers here (their terminal triple is
// asserted by checkTerminalContract), and a failure case that declares no text
// must not have emitted any.
func checkTurn(t *testing.T, c Case, turn llm.TurnResult) {
	t.Helper()
	if c.WantClass != "" {
		if c.WantText != "" && turn.Text != c.WantText {
			t.Errorf("partial text = %q, want %q (partial content must be retained)",
				turn.Text, c.WantText)
		}
		if c.WantText == "" && turn.Text != "" {
			t.Errorf("unexpected text on a pre-payload failure: %q", turn.Text)
		}
		if turn.Incomplete != c.WantIncomplete {
			t.Errorf("incomplete = %v, want %v", turn.Incomplete, c.WantIncomplete)
		}
		return
	}
	if c.WantStop != "" && turn.Stop != c.WantStop {
		t.Errorf("stop = %s, want %s", turn.Stop, c.WantStop)
	}
	if turn.Text != c.WantText {
		t.Errorf("text = %q, want %q", turn.Text, c.WantText)
	}
	if c.WantReasoning != "" {
		if turn.Reasoning == "" {
			t.Errorf("no ReasoningDelta assembled (reasoning=%q events below)", turn.Reasoning)
		}
		if turn.Reasoning != c.WantReasoning {
			t.Errorf("reasoning = %q, want %q", turn.Reasoning, c.WantReasoning)
		}
	}
	if msg := eqToolCalls(turn.ToolCalls, c.WantToolCalls); msg != "" {
		t.Error(msg)
	}
	if turn.Usage != c.WantUsage {
		t.Errorf("usage = %+v, want %+v", turn.Usage, c.WantUsage)
	}
	if turn.Incomplete != c.WantIncomplete {
		t.Errorf("incomplete = %v, want %v", turn.Incomplete, c.WantIncomplete)
	}
	if c.WantClass == "" && turn.Err != nil {
		t.Errorf("unexpected Error event: %v", turn.Err)
	}
}

// checkTerminalContract asserts the event-ordering contract plus failure
// expectations for every mode.
func checkTerminalContract(t *testing.T, c Case, events []llm.StreamEvent, turn llm.TurnResult, err error) {
	t.Helper()
	if len(events) == 0 {
		t.Fatal("stream emitted no events at all")
	}
	if last := events[len(events)-1]; last.Type != llm.EvDone {
		t.Errorf("last event = %s, want done", last.Type)
	}
	done := 0
	for _, ev := range events {
		if ev.Type == llm.EvDone {
			done++
		}
	}
	if done != 1 {
		t.Errorf("Done emitted %d times, want exactly 1", done)
	}
	// Exactly one terminal pair.
	stops := 0
	for _, ev := range events {
		if ev.Type == llm.EvStop {
			stops++
		}
	}
	if stops != 1 {
		t.Errorf("Stop emitted %d times, want exactly 1", stops)
	}

	if c.WantClass == "" {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if turn.Stop == llm.StopError {
			t.Error("success turn ended with Stop{error}")
		}
	} else {
		if err == nil {
			t.Fatal("expected an error")
		}
		if ClassOf(err) != c.WantClass {
			t.Errorf("class = %s, want %s", ClassOf(err), c.WantClass)
		}
		if c.WantCode != "" && CodeOf(err) != c.WantCode {
			t.Errorf("provider_code = %q, want %q", CodeOf(err), c.WantCode)
		}
		// Failure streams must carry the terminal triple in order.
		if turn.Stop != llm.StopError {
			t.Errorf("terminal stop = %s, want error", turn.Stop)
		}
		if turn.Err == nil {
			t.Error("failure stream emitted no Error event")
		}
	}
	if len(c.WantEvents) > 0 {
		var got []llm.StreamEventType
		for _, ev := range events {
			got = append(got, ev.Type)
		}
		if len(got) < len(c.WantEvents) {
			t.Fatalf("event kinds %v, want to contain %v", got, c.WantEvents)
		}
		// The wanted sequence must appear in order (not necessarily
		// contiguously: dialects may interleave more deltas).
		i := 0
		for _, ev := range got {
			if ev == c.WantEvents[i] {
				i++
				if i == len(c.WantEvents) {
					break
				}
			}
		}
		if i != len(c.WantEvents) {
			t.Errorf("event kinds %v do not contain ordered %v", got, c.WantEvents[:i])
		}
	}
	// Tool-call ordering contract on every adapter: ArgsDelta only between
	// Start and End of its id.
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
	if len(open) != 0 && !c.WantOpenToolCalls {
		t.Errorf("calls left open: %v", open)
	}
}

// RetryOptionsOf is the ladder the shared retry cases use (Base 5ms so the
// only measured delays come from Retry-After hints).
func RetryOptionsOf() llm.RetryOptions {
	return llm.RetryOptions{Max: 3, Base: 5 * time.Millisecond}
}
