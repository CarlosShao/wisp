package agent

// ACCEPTANCE PROBE 4 (ticket 153 accept-r1), snapshot only: a context-carried
// tag is INHERITED by derived contexts. Measure what that means for attribution:
// a fold that happens on a child ctx of a tagged run reports the PARENT's id.

import (
	"context"
	"testing"
)

func TestAccept153TagIsInheritedByDerivedContexts(t *testing.T) {
	recs := newTraceCapture()
	c := NewCompressor(BudgetsFor(4096), nil, WithLogger(recs.logger()))
	hist := buildRoundHistory(6, 400)

	id := newTaskID()
	parent := withTraceTask(context.Background(), id)
	child, cancel := context.WithCancel(parent)
	defer cancel()

	if _, rep, err := c.Compress(child, hist); err != nil || !rep.Ran {
		t.Fatalf("child-ctx Compress: err=%v ran=%v", err, rep.Ran)
	}
	got, present := recs.with(traceMsg)[0].attr["task"]
	if !present || got != id {
		t.Fatalf("derived ctx: task=%v present=%v, want %q (if this fails the ruler is wrong)", got, present, id)
	}
	t.Logf("PROBE4-153 INHERIT: fold on a ctx DERIVED from a tagged run still reports the parent task id %q", id)

	// and an unrelated background ctx (same process, no lineage) stays empty
	if _, _, err := c.Compress(context.Background(), hist); err != nil {
		t.Fatalf("untagged control: %v", err)
	}
	second := recs.with(traceMsg)[1]
	if _, present := second.attr["task"]; present {
		t.Errorf("untagged control got a task key: %v", second.attr)
	}
}
