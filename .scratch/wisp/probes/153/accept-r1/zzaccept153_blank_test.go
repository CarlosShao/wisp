package agent

// ACCEPTANCE PROBE 3 (ticket 153 accept-r1), snapshot only: whitespace-only and
// empty ids must leave NO key at all (the "no fallback value" claim, measured
// rather than read off the source).

import (
	"context"
	"testing"
)

func TestAccept153EmptyAndBlankIdsLeaveNoKey(t *testing.T) {
	recs := newTraceCapture()
	c := NewCompressor(BudgetsFor(4096), nil, WithLogger(recs.logger()))
	hist := buildRoundHistory(6, 400)

	for _, in := range []string{"", "   ", "\t", " \t \n "} {
		if _, rep, err := c.Compress(withTraceTask(context.Background(), in), hist); err != nil || !rep.Ran {
			t.Fatalf("setup %q: err=%v ran=%v", in, err, rep.Ran)
		}
	}
	// and one plain untagged call, plus one real id as the control that the
	// ruler is not simply blind to the key
	id := newTaskID()
	if _, _, err := c.Compress(withTraceTask(context.Background(), id), hist); err != nil {
		t.Fatalf("control call: %v", err)
	}
	hits := recs.with(traceMsg)
	if len(hits) != 5 {
		t.Fatalf("records = %d, want 5", len(hits))
	}
	for i := 0; i < 4; i++ {
		if v, present := hits[i].attr["task"]; present {
			t.Errorf("record %d carries task=%q: a blank id must not reach the record", i, v)
		}
	}
	if got, present := hits[4].attr["task"]; !present || got != id {
		t.Errorf("positive control lost its key: present=%v got=%v want %q", present, got, id)
	}
	t.Logf("PROBE3-153: 4 blank-ish ids -> no task key at all; control id present as expected")
}
