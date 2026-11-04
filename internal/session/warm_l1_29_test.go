package session

// Ticket 29 AC#1, the cross-package half of the latency-isolation instrument.
//
// internal/memory's own test proves the two halves it owns: Record (the
// response-path half) makes ZERO C5 calls, and Fire returns while the provider
// is still blocked. What it cannot prove from inside internal/memory is the
// MECHANISM that runs Fire at all: that a live session's Warm entry
// (scope.go's EnterWarm) is what invokes WarmWindow, which is the exact
// assignment cmd/wisp/resident_voice_engines_windows.go makes
// ("WarmWindow: r.warmWindow", whose body is the memory warm window's Fire).
//
// That is the assertion this file adds: with a provider that cannot answer until
// the test releases it, the response path finishes, EnterWarm returns, and the
// extraction only then completes. A response path that ran the extraction inline
// would have to block here - and this test would hang on the release it never
// gets.

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
)

// warmL1Provider is the C5 seam injection for this test: it records the request,
// then blocks until the test releases it (or the context expires).
type warmL1Provider struct {
	mu      sync.Mutex
	reply   string
	calls   int
	entered chan struct{}
	release chan struct{}
}

func (p *warmL1Provider) Info() llm.ProviderInfo {
	return llm.ProviderInfo{Protocol: "openai-chat", Provider: "fake", Model: "warm-l1-test"}
}

func (p *warmL1Provider) Stream(ctx context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	p.mu.Lock()
	p.calls++
	reply := p.reply
	p.mu.Unlock()
	select {
	case p.entered <- struct{}{}:
	default:
	}
	select {
	case <-p.release:
	case <-ctx.Done():
		if err := emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopCancelled}); err != nil {
			return err
		}
		return emit(llm.StreamEvent{Type: llm.EvDone})
	}
	if err := emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: reply}); err != nil {
		return err
	}
	if err := emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopEndTurn}); err != nil {
		return err
	}
	return emit(llm.StreamEvent{Type: llm.EvDone})
}

func (p *warmL1Provider) callCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.calls
}

// TestTicket29WarmHookRunsL1Extraction is the AC#1 mechanism test: the session
// Warm entry, and nothing else, starts the L1 extraction, and it starts it off
// the response path.
func TestTicket29WarmHookRunsL1Extraction(t *testing.T) {
	store, err := memory.Open(t.TempDir())
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer func() { _ = store.Close() }()

	prov := &warmL1Provider{
		reply:   `{"upserts":[{"slot":"pref.tone","content":"concise"},{"slot":"habit.editor","content":"vscode"}],"forgets":[]}`,
		entered: make(chan struct{}, 1),
		release: make(chan struct{}),
	}
	ex, err := memory.NewL1Extractor(prov, "warm-l1-test")
	if err != nil {
		t.Fatalf("NewL1Extractor: %v", err)
	}
	pipe, err := memory.NewL1Pipeline(store, ex)
	if err != nil {
		t.Fatalf("NewL1Pipeline: %v", err)
	}
	done := make(chan error, 1)
	warm, err := memory.NewL1WarmWindow(pipe, memory.WithL1WarmObserve(func(_ memory.ExtractionOutcome, err error) {
		done <- err
	}))
	if err != nil {
		t.Fatalf("NewL1WarmWindow: %v", err)
	}

	log := &eventLog{}
	scope, err := NewScope(ScopeOptions{
		Loader: newEngineStub(log).load,
		Clock:  newFakeClock(),
		Hooks: Hooks{
			MicClose: func() {},
			// THIS is the assignment cmd/wisp makes at assembly
			// (resident_voice_engines_windows.go: "WarmWindow: r.warmWindow",
			// whose body is warm.Fire). A session with no such hook is the
			// pre-AC#1 shape the ticket's 29-v1 ruling found in the host.
			WarmWindow: warm.Fire,
			Logf:       func(string, ...any) {},
		},
	})
	if err != nil {
		t.Fatalf("NewScope: %v", err)
	}
	defer func() { _ = scope.Dispose("test-end") }()

	// THE RESPONSE PATH: the task finished, so its conversation is handed over.
	// Nothing here may reach C5 - that is what "never on the response path" is.
	responseStart := time.Now()
	warm.Record(memory.GoldenConversation{
		ID: "task-ac1-hook",
		Turns: []memory.GoldenTurn{
			{Role: "user", Text: "以后回答简短一点"},
			{Role: "assistant", Text: "好。"},
			{Role: "user", Text: "我用 vscode"},
		},
	})
	responseDur := time.Since(responseStart)
	if n := prov.callCount(); n != 0 {
		t.Fatalf("the response path made %d C5 call(s); extraction must not run inline", n)
	}
	if responseDur > 5*time.Millisecond {
		t.Fatalf("Record took %v on the response path; it must only copy the turns", responseDur)
	}
	if !warm.Pending() {
		t.Fatal("the finished task's conversation was not recorded for the Warm window")
	}

	// THE WARM ENTRY: this is the only thing that may start the extraction.
	warmStart := time.Now()
	if err := scope.EnterWarm("announce-end"); err != nil {
		t.Fatalf("EnterWarm: %v", err)
	}
	warmDur := time.Since(warmStart)
	if warmDur > 20*time.Millisecond {
		t.Fatalf("EnterWarm blocked for %v: the Warm hook waited for the extraction call", warmDur)
	}
	select {
	case <-prov.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("the session Warm hook never started the extraction: WarmWindow is not wired to the pipeline")
	}
	if n := prov.callCount(); n != 1 {
		t.Fatalf("C5 calls after one Warm entry = %d, want 1 (one call per task)", n)
	}
	if st := warm.Stats(); st.Succeeded != 0 {
		t.Fatalf("batch completed while the provider was still blocked: %+v", st)
	}
	// The stage landed BEFORE the work started - the window is open when the
	// extraction runs, which is what the AC means by "in the Warm window".
	if got := scope.Stage(); got != StageWarm {
		t.Fatalf("stage during the extraction = %s, want Warm", got)
	}
	ctx := context.Background()
	if _, err := store.ProfileBySlot(ctx, "pref.tone"); !errors.Is(err, memory.ErrNotFound) {
		t.Fatalf("a profile row landed before the provider answered: %v", err)
	}

	// Release the C5 call: the whole batch commits now, and only now.
	close(prov.release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("extraction failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("extraction did not finish after the provider was released")
	}
	for slot, want := range map[string]string{"pref.tone": "concise", "habit.editor": "vscode"} {
		p, err := store.ProfileBySlot(ctx, slot)
		if err != nil {
			t.Fatalf("ProfileBySlot(%s) after the Warm window: %v", slot, err)
		}
		if p.Content != want {
			t.Errorf("slot %s = %q, want %q (the Warm window committed a different value)", slot, p.Content, want)
		}
	}
	if st := warm.Stats(); st.Succeeded != 1 || st.Failed != 0 || st.Started != 1 {
		t.Fatalf("warm stats after one batch = %+v, want Started=1 Succeeded=1", st)
	}
	if prov.callCount() != 1 {
		t.Fatalf("C5 calls = %d after the batch, want exactly 1", prov.callCount())
	}

	// A second Warm entry with nothing pending is the ordinary no-op: no new
	// call, no failure (the hook fires on every Warm landing, not just after a
	// task).
	if err := scope.Resume(); err != nil {
		t.Fatalf("Resume: %v", err)
	}
	if err := scope.EnterWarm("turn-end"); err != nil {
		t.Fatalf("EnterWarm(second): %v", err)
	}
	if prov.callCount() != 1 {
		t.Fatalf("an empty Warm window made another C5 call: %d", prov.callCount())
	}
	if st := warm.Stats(); st.Dropped == 0 {
		t.Fatalf("the empty Warm entry was not counted as a drop: %+v", st)
	}
}

// compile-time proof that the test provider is a real C5 seam value.
var _ llm.LlmProvider = (*warmL1Provider)(nil)
