package memory

// Ticket 29 AC#1 tests: the real L1 extraction chain and the Warm-window seam.
//
// The three instruments are the AC's three claims:
//
//   - TestTicket29_AC1_WarmWindowLatencyIsolation: extraction runs
//     POST-RESPONSE, and the response path pays nothing for it. The mechanism
//     is what is asserted, not a stopwatch: Record (the response-path half) is
//     checked to make ZERO C5 calls, Fire (the Warm-hook half) is checked to
//     return while the provider is still blocked, and the first-token reading
//     is the supporting number - with a deliberately slow provider, the
//     extraction cannot possibly have contributed to it.
//   - TestTicket29_AC1_OneCallPerTaskBatchCommit: one task = one C5 call (so at
//     most one prompt-cache invalidation), one store transaction, and a broken
//     reply writes NOTHING (no half-applied profile).
//   - TestTicket29_AC1_ParserMatchesGoldenContract: the real chain's reply
//     grammar is the golden contract's - for all 10 pinned conversations, the
//     pinned expected outcome, fed back through the extractor as the model's
//     JSON answer, reproduces the fixture exactly through CompareExtraction.
//     (Whether a LIVE model produces those answers is AC#2's golden-SSE work,
//     the named residual; this test pins the shape the future chain replays
//     into.)

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// l1FakeProvider is the C5 seam injection (SPEC-05's golden-SSE harness is the
// same seam for the recorded-bytes leg). It answers with scripted assistant
// text and can hold the call open so a test can prove the caller did not wait.
type l1FakeProvider struct {
	mu        sync.Mutex
	reply     string
	calls     int
	reqs      []llm.Request
	streamErr error

	// entered receives one value whenever Stream starts.
	entered chan struct{}
	// release, when non-nil, blocks the stream after the request is recorded:
	// the "slow C5 call" the latency test needs to be able to catch.
	release chan struct{}
}

func newL1FakeProvider(reply string) *l1FakeProvider {
	return &l1FakeProvider{reply: reply, entered: make(chan struct{}, 16)}
}

func (p *l1FakeProvider) Info() llm.ProviderInfo {
	return llm.ProviderInfo{Protocol: "openai-chat", Provider: "fake", Model: "l1-test-model"}
}

func (p *l1FakeProvider) Stream(ctx context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	p.mu.Lock()
	p.calls++
	p.reqs = append(p.reqs, *req)
	p.mu.Unlock()

	select {
	case p.entered <- struct{}{}:
	default:
	}
	p.mu.Lock()
	release := p.release
	streamErr := p.streamErr
	p.mu.Unlock()

	if release != nil {
		select {
		case <-release:
		case <-ctx.Done():
			// C6: cancellation emits Stop{cancelled} then Done, and Stream
			// returns nil (this is the contract the adapters honor).
			if err := emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopCancelled}); err != nil {
				return err
			}
			return emit(llm.StreamEvent{Type: llm.EvDone})
		}
	}
	if streamErr != nil {
		oe := observe.Wrap(observe.ClassProvider, streamErr, "l1 fake provider")
		if err := emit(llm.StreamEvent{Type: llm.EvError, Err: oe}); err != nil {
			return err
		}
		if err := emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopError}); err != nil {
			return err
		}
		_ = emit(llm.StreamEvent{Type: llm.EvDone})
		return oe
	}
	if err := emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: p.reply}); err != nil {
		return err
	}
	if err := emit(llm.StreamEvent{Type: llm.EvUsage, Usage: llm.Usage{InputTokens: 120, OutputTokens: 24}}); err != nil {
		return err
	}
	if err := emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopEndTurn}); err != nil {
		return err
	}
	return emit(llm.StreamEvent{Type: llm.EvDone})
}

func (p *l1FakeProvider) callCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.calls
}

func (p *l1FakeProvider) lastRequest() (llm.Request, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.reqs) == 0 {
		return llm.Request{}, false
	}
	return p.reqs[len(p.reqs)-1], true
}

// l1FirstToken is the response path in miniature: one emitted text delta, which
// is the exact moment "first-token latency" is measured at in the loop.
func l1FirstToken() (time.Duration, bool) {
	col := llm.NewTurnCollector()
	start := time.Now()
	col.Observe(llm.StreamEvent{Type: llm.EvTextDelta, Text: "ok"})
	d := time.Since(start)
	return d, col.SawPayload()
}

func l1Median(d []time.Duration) time.Duration {
	if len(d) == 0 {
		return 0
	}
	cp := append([]time.Duration(nil), d...)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })
	return cp[len(cp)/2]
}

// TestTicket29_AC1_WarmWindowLatencyIsolation is the AC's latency-isolation
// instrument. It asserts the MECHANISM that makes extraction post-response:
//
//	Record (response path) -> zero C5 calls
//	Fire (Warm hook body)  -> returns before the C5 call can finish
//
// and it carries the first-token number as the supporting reading, taken with
// an extraction that is deliberately still blocked while the turn is streamed.
func TestTicket29_AC1_WarmWindowLatencyIsolation(t *testing.T) {
	store := openTestStore(t)
	const reply = `{"upserts":[{"slot":"pref.tone","content":"concise"},` +
		`{"slot":"habit.editor","content":"vscode"}],"forgets":[]}`
	prov := newL1FakeProvider(reply)
	prov.release = make(chan struct{}) // hold the extraction open until we free it

	ex, err := NewL1Extractor(prov, "l1-test-model")
	if err != nil {
		t.Fatalf("NewL1Extractor: %v", err)
	}
	pipe, err := NewL1Pipeline(store, ex)
	if err != nil {
		t.Fatalf("NewL1Pipeline: %v", err)
	}
	done := make(chan error, 1)
	warm, err := NewL1WarmWindow(pipe, WithL1WarmObserve(func(_ ExtractionOutcome, err error) {
		done <- err
	}))
	if err != nil {
		t.Fatalf("NewL1WarmWindow: %v", err)
	}

	conv := GoldenConversation{
		ID:    "task-ac1-latency",
		Title: "latency isolation",
		Turns: []GoldenTurn{
			{Role: "user", Text: "以后回答简短一点，说重点"},
			{Role: "assistant", Text: "好，我尽量简洁。"},
			{Role: "user", Text: "我平时用 vscode 写代码"},
		},
	}

	// Baseline: the response path with NO extraction pipeline in the process.
	const turns = 40
	base := make([]time.Duration, 0, turns)
	for i := 0; i < turns; i++ {
		d, saw := l1FirstToken()
		if !saw {
			t.Fatal("baseline first token was not observed: the miniature response path is broken")
		}
		base = append(base, d)
	}

	// Wired: the same response path, but each finished task is handed to the
	// Warm window AFTER its response (Record is post-response by construction).
	wired := make([]time.Duration, 0, turns)
	for i := 0; i < turns; i++ {
		d, saw := l1FirstToken()
		if !saw {
			t.Fatal("wired first token was not observed: the miniature response path is broken")
		}
		warm.Record(conv)
		wired = append(wired, d)
	}

	// MECHANISM, part 1: nothing on the response path reached C5. This is the
	// assertion a stopwatch cannot make for us: even a fast extraction could
	// hide inside a noisy timing, but a call count cannot.
	if got := prov.callCount(); got != 0 {
		t.Fatalf("the response path made %d C5 call(s): extraction is inline, not post-response", got)
	}
	st := warm.Stats()
	if st.Recorded != turns || !warm.Pending() {
		t.Fatalf("warm window recorded %d conversations (pending=%v), want %d recorded and one pending",
			st.Recorded, warm.Pending(), turns)
	}
	if st.Started != 0 {
		t.Fatalf("warm window already started %d batch(es) before any Warm entry", st.Started)
	}

	// First-token reading: the wired path may not regress beyond measurement
	// noise. The hard proof is the call count above; this is the number the AC
	// asks to keep unchanged.
	if baseMed, wiredMed := l1Median(base), l1Median(wired); wiredMed > baseMed+5*time.Millisecond {
		t.Fatalf("first-token median regressed with extraction wired: baseline %v, wired %v", baseMed, wiredMed)
	}

	// MECHANISM, part 2: Fire is the Warm hook body, and it must RETURN without
	// waiting for the C5 call (the hook runs on the voice owner goroutine).
	fireStart := time.Now()
	warm.Fire()
	fireDur := time.Since(fireStart)
	if fireDur > 20*time.Millisecond {
		t.Fatalf("Fire blocked for %v: the Warm hook waited for the extraction call", fireDur)
	}
	select {
	case <-prov.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("Fire did not start the extraction goroutine")
	}
	if got := prov.callCount(); got != 1 {
		t.Fatalf("C5 calls after Fire = %d, want exactly 1 (one call per task)", got)
	}
	if st := warm.Stats(); st.Succeeded != 0 {
		t.Fatalf("batch completed while the provider was still blocked: %+v", st)
	}
	if warm.Pending() {
		t.Fatal("Fire left the batch pending: the warm window would extract it twice")
	}

	// Release the C5 call; the batch must land, once, as one transaction.
	close(prov.release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("extraction failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("extraction did not finish after the provider was released")
	}
	if st := warm.Stats(); st.Succeeded != 1 || st.Failed != 0 || st.Started != 1 {
		t.Fatalf("warm window stats after one batch = %+v, want Started=1 Succeeded=1", st)
	}
	ctx := context.Background()
	for slot, want := range map[string]string{"pref.tone": "concise", "habit.editor": "vscode"} {
		p, err := store.ProfileBySlot(ctx, slot)
		if err != nil {
			t.Fatalf("ProfileBySlot(%s): %v", slot, err)
		}
		if p.Content != want || p.Source != string(ProfileExtracted) {
			t.Errorf("slot %s = %q/%q, want %q/extracted", slot, p.Content, p.Source, want)
		}
	}

	// The request itself is the conversation the task held, not a shadow.
	req, ok := prov.lastRequest()
	if !ok {
		t.Fatal("no request recorded")
	}
	if len(req.Messages) != 1 || len(req.System) != 1 {
		t.Fatalf("extraction request shape = %d system / %d messages, want 1/1", len(req.System), len(req.Messages))
	}
	body := ""
	if tp, ok := req.Messages[0].Content[0].(llm.TextPart); ok {
		body = tp.Text
	}
	for _, want := range []string{"以后回答简短一点", "我平时用 vscode 写代码", "user: "} {
		if !strings.Contains(body, want) {
			t.Errorf("extraction request body does not carry %q", want)
		}
	}
	if strings.Contains(body, conv.ID) || strings.Contains(body, conv.Title) {
		t.Error("extraction request leaks the fixture metadata (id/title) into the prompt")
	}
}

// TestTicket29_AC1_OneCallPerTaskBatchCommit covers the batching claim: ONE C5
// call per task, ONE store transaction for its whole outcome, and a reply that
// breaks the contract writes nothing at all.
func TestTicket29_AC1_OneCallPerTaskBatchCommit(t *testing.T) {
	store := openTestStore(t)
	ctx := context.Background()

	// Three user facts in one task: one call, three slots, one commit.
	prov := newL1FakeProvider(`{"upserts":[
		{"slot":"pref.language","content":"zh"},
		{"slot":"pref.tone","content":"direct"},
		{"slot":"habit.work_hours","content":"weekday 09:00-18:00"}],"forgets":[]}`)
	ex, err := NewL1Extractor(prov, "l1-test-model")
	if err != nil {
		t.Fatalf("NewL1Extractor: %v", err)
	}
	pipe, err := NewL1Pipeline(store, ex)
	if err != nil {
		t.Fatalf("NewL1Pipeline: %v", err)
	}
	conv := GoldenConversation{ID: "task-ac1-batch", Turns: []GoldenTurn{
		{Role: "user", Text: "以后都用中文"},
		{Role: "user", Text: "直接说结论，别绕"},
		{Role: "user", Text: "我工作日早上9点上班，下午6点下班"},
		{Role: "assistant", Text: "记下了。"},
	}}
	if _, err := pipe.RunTask(ctx, conv); err != nil {
		t.Fatalf("RunTask: %v", err)
	}
	if got := prov.callCount(); got != 1 {
		t.Fatalf("C5 calls for one task = %d, want 1 (batch-commit: one prompt-cache invalidation per task, not per fact)", got)
	}
	rows, err := store.ListProfiles(ctx)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("committed %d rows, want 3: %+v", len(rows), rows)
	}

	// A reply with one slot outside the SPEC-02 namespace fails the WHOLE
	// batch: nothing from it may land.
	bad := newL1FakeProvider(`{"upserts":[{"slot":"pref.tone","content":"casual"},` +
		`{"slot":"user.name","content":"x"}],"forgets":[]}`)
	badEx, _ := NewL1Extractor(bad, "l1-test-model")
	badPipe, _ := NewL1Pipeline(store, badEx)
	if _, err := badPipe.RunTask(ctx, conv); err == nil {
		t.Fatal("a slot outside pref./habit./fact. was accepted; the batch must be refused whole")
	}
	if _, err := store.ProfileBySlot(ctx, "pref.tone"); err != nil {
		t.Fatalf("the refused batch wrote pref.tone anyway: %v", err)
	}
	p, err := store.ProfileBySlot(ctx, "pref.tone")
	if err != nil || p.Content != "direct" {
		t.Fatalf("the refused batch overwrote pref.tone: %+v err=%v", p, err)
	}

	// A provider failure is a silent skip: no rows change, an error is returned
	// for the caller's debug log.
	failing := newL1FakeProvider("")
	failing.streamErr = errors.New("provider down")
	failEx, _ := NewL1Extractor(failing, "l1-test-model")
	failPipe, _ := NewL1Pipeline(store, failEx)
	if _, err := failPipe.RunTask(ctx, conv); err == nil {
		t.Fatal("a failing provider was reported as success")
	}
	after, err := store.ListProfiles(ctx)
	if err != nil {
		t.Fatalf("ListProfiles: %v", err)
	}
	if len(after) != 3 {
		t.Fatalf("a failed extraction changed the profile: %d rows, want 3", len(after))
	}

	// Net semantics inside one outcome: a forget cancels a same-slot upsert.
	net := newL1FakeProvider(`{"upserts":[{"slot":"fact.family","content":"mother birthday 05-06"}],` +
		`"forgets":["fact.family"]}`)
	netEx, _ := NewL1Extractor(net, "l1-test-model")
	// A slot stored by an EARLIER task: the forget has a real row to remove.
	seedID, err := store.UpsertProfile(ctx, "fact.family", "father birthday 01-02", ProfileManual)
	if err != nil {
		t.Fatalf("seed fact.family: %v", err)
	}
	netPipe, _ := NewL1Pipeline(store, netEx)
	out, err := netPipe.RunTask(ctx, conv)
	if err != nil {
		t.Fatalf("RunTask(net): %v", err)
	}
	if len(out.Upserts) != 0 || len(out.Forgets) != 1 {
		t.Fatalf("net outcome = %+v, want no upsert and one forget", out)
	}
	if _, err := store.ProfileBySlot(ctx, "fact.family"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("forget did not remove the slot (seeded id %d): %v", seedID, err)
	}
}

// TestTicket29_AC1_ParserMatchesGoldenContract pins the real chain's reply
// grammar to the reviewed fixture: for every pinned conversation, the pinned
// expected outcome - rendered as the model's JSON answer, wrapped in the prose
// and fences real models add - must come back out of Extract and pass the same
// CompareExtraction the stub passes. The AC#2 leg replays recorded SSE bytes
// into this same extractor; this test is what makes that replay's target the
// fixture and not whatever the model happened to answer.
func TestTicket29_AC1_ParserMatchesGoldenContract(t *testing.T) {
	set, err := LoadGoldenExtractSet()
	if err != nil {
		t.Fatalf("LoadGoldenExtractSet: %v", err)
	}
	ctx := context.Background()
	for i := range set.Cases {
		c := set.Cases[i]
		t.Run(c.ID, func(t *testing.T) {
			payload, err := json.Marshal(c.Expect)
			if err != nil {
				t.Fatalf("marshal expectation: %v", err)
			}
			// The tolerant wrapper: prose first, a fence, then the object. The
			// extractor must still read exactly the pinned outcome.
			reply := "Here is the extraction result:\n```json\n" + string(payload) + "\n```\nDone."
			prov := newL1FakeProvider(reply)
			ex, err := NewL1Extractor(prov, "l1-test-model")
			if err != nil {
				t.Fatalf("NewL1Extractor: %v", err)
			}
			got, err := ex.Extract(ctx, GoldenConversation{ID: c.ID, Title: c.Title, Turns: c.Turns})
			if err != nil {
				t.Fatalf("Extract: %v", err)
			}
			if err := CompareExtraction(got, c.Expect); err != nil {
				t.Fatalf("the real chain's grammar disagrees with the pinned contract: %v", err)
			}
			if prov.callCount() != 1 {
				t.Fatalf("C5 calls = %d, want 1 per task", prov.callCount())
			}
		})
	}
}
