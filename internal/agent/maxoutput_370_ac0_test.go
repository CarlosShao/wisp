package agent

// Ticket 370 AC#0(4): the positive control.
//
// The ticket's claim is that the state machine is fine and only the plumbing
// was missing. This probe isolates the state machine from the plumbing: it sets
// Config.MaxOutputTokens by hand (no config file, no llm.Resolver, no cmd/wisp),
// hands the loop a C5 provider that ends its one turn with llm.StopMaxTokens
// exactly as a real provider does when max_tokens runs out, and then reads the
// task_log row back out of a real SQLite store through the loop's own journal.
//
// Nothing here constructs a Result: res is what Loop.Run returned, and the row
// is what the loop wrote. If StatusTruncated were unreachable as a status, the
// first assertion would fail; if its row mapping were dead, the last one would.
//
// The reverse direction (ceiling omitted -> provider-own limit) is measured in
// cmd/wisp/task370_ac2_test.go's "no_ceiling_on_the_wire" control, where the
// same run ends "done" because nothing cut it.

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

const t370AC0Ceiling = 1

// truncatingProvider370 is the smallest honest C5 implementation for this
// probe: it records the ceiling it was handed, emits a partial text turn, and
// terminates with StopMaxTokens (the C6 reason a provider uses for "the output
// budget ran out"). It invents no Result - the loop assembles the turn.
type truncatingProvider370 struct {
	sawMaxOutput int
	calls        int
}

func (p *truncatingProvider370) Info() llm.ProviderInfo {
	return llm.ProviderInfo{
		Provider:         "probe-370",
		Model:            "probe-370",
		Protocol:         "probe",
		MaxContextWindow: 128000,
	}
}

func (p *truncatingProvider370) Stream(_ context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	p.calls++
	p.sawMaxOutput = req.MaxOutputTokens
	for _, ev := range []llm.StreamEvent{
		{Type: llm.EvTextDelta, Text: "half a sen"},
		{Type: llm.EvStop, Stop: llm.StopMaxTokens},
		{Type: llm.EvDone},
	} {
		if err := emit(ev); err != nil {
			return err
		}
	}
	return nil
}

func Test370AC0HandSetCeilingRowsTruncated(t *testing.T) {
	dir := sealableTempDir124(t)
	store, err := memory.Open(dir,
		memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()

	prov := &truncatingProvider370{}
	loop, err := New(Options{
		Provider: prov,
		Tools:    NewEchoProvider(),
		Sink:     &RecordSink{},
		Journal:  store,
		Registry: observe.NewRegistry(),
		Logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
		Config: Config{
			Model:         "probe-370",
			ContextWindow: 128000,
			// The hand-set ceiling: this is the value the missing hop used to
			// zero out. No resolver or composition root is involved here.
			MaxOutputTokens:             t370AC0Ceiling,
			ArtifactsDir:                filepath.Join(dir, "artifacts"),
			PerToolTimeout:              2 * time.Second,
			PassThroughUnclassifiedRisk: true,
			Currency:                    "CNY",
		},
	})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	if got := loop.asm.MaxOutputTokens; got != t370AC0Ceiling {
		t.Fatalf("Assembler.MaxOutputTokens = %d, want %d", got, t370AC0Ceiling)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res := loop.Run(ctx, "手设上限的一发")

	// Premise: the hand-set ceiling actually reached the C5 seam. Without this
	// the row below could be green for the wrong reason.
	if prov.calls != 1 {
		t.Fatalf("provider calls = %d, want 1", prov.calls)
	}
	if prov.sawMaxOutput != t370AC0Ceiling {
		t.Fatalf("the request carried MaxOutputTokens = %d, want %d: the hand-set ceiling never reached the seam",
			prov.sawMaxOutput, t370AC0Ceiling)
	}

	// The status the real loop produced from a real StopMaxTokens.
	if res.Status != StatusTruncated {
		t.Fatalf("Loop.Run status = %q, want %q: StopMaxTokens did not drive StatusTruncated",
			string(res.Status), string(StatusTruncated))
	}

	// The row the real journal wrote into real SQLite.
	tl, err := store.TaskLogByID(ctx, res.TaskID)
	if err != nil {
		t.Fatalf("task_log for %s: %v", res.TaskID, err)
	}
	t.Logf("positive control: request.MaxOutputTokens=%d; Loop.Run status=%q; task_log.state=%q",
		prov.sawMaxOutput, string(res.Status), tl.State)
	if tl.State != string(memory.TaskStateInterrupted) {
		t.Fatalf("task_log.state = %q, want %q for a task the provider cut at max_tokens",
			tl.State, string(memory.TaskStateInterrupted))
	}
	if !memory.ValidTaskState(tl.State) {
		t.Fatalf("task_log.state = %q is off the roster yet the row exists", tl.State)
	}
}

// Test370AC0ProviderOwnLimitTruncatesWithoutOurCeiling is AC#0(5): the other
// end of the range. The config sets no ceiling, so openaichat's
// `json:"max_tokens,omitempty"` drops the field from the wire, yet the
// recorded provider stream still says finish_reason "length" (the golden
// fixture's own bytes, parsed by the real adapter). The loop must still end
// StatusTruncated: "our ceiling is unreachable" is true when it is unset, but
// "the status is unreachable" is not.
func Test370AC0ProviderOwnLimitTruncatesWithoutOurCeiling(t *testing.T) {
	h := newHarness(t, "max-tokens-toolcall") // Config.MaxOutputTokens stays 0
	res := h.run("这一段会撞上供应商自己的上限")

	bodies := h.requestBodies()
	if len(bodies) == 0 {
		t.Fatal("no request went out at all")
	}
	if bytes.Contains(bodies[0], []byte(`"max_tokens"`)) {
		t.Fatalf("the request carries a max_tokens even though no ceiling was configured:\n%s",
			trunc(string(bodies[0]), 600))
	}
	t.Logf("no ceiling on the wire (body has no max_tokens); recorded finish_reason=length -> status=%q",
		string(res.Status))
	if res.Status != StatusTruncated {
		t.Fatalf("status = %q, want %q: a provider-own limit stop must still drive StatusTruncated",
			string(res.Status), string(StatusTruncated))
	}
}
