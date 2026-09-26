package agent

// ACCEPTANCE PROBE 2 for ticket 153 (accept-r1), out-of-repo snapshot only.
// Closes the impl's admitted gap (its s7 item 2): attribution under REAL
// concurrency on a SHARED compressor - several goroutines, each tagging its own
// task id, all folding at once. If the id ever lived on the shared object, or
// if the ctx channel could cross, the per-id record counts would not add up.

import (
	"context"
	"sync"
	"testing"
)

func TestAccept153AttributionUnderRealConcurrency(t *testing.T) {
	recs := newTraceCapture()
	b := BudgetsFor(4096)
	c := NewCompressor(b, nil, WithLogger(recs.logger())) // ONE shared compressor
	hist := buildRoundHistory(6, 400)

	const n = 24
	ids := make([]string, n)
	for i := range ids {
		ids[i] = newTaskID()
	}
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx := withTraceTask(context.Background(), ids[i])
			if _, rep, err := c.Compress(ctx, hist); err != nil || !rep.Ran {
				t.Errorf("goroutine %d: err=%v ran=%v", i, err, rep.Ran)
			}
		}(i)
	}
	wg.Wait()

	hits := recs.with(traceMsg)
	if len(hits) != n {
		t.Fatalf("records = %d, want %d", len(hits), n)
	}
	seen := map[string]int{}
	for _, r := range hits {
		v, ok := r.attr["task"]
		if !ok {
			t.Fatalf("a concurrent record lost its task key: %v", r.attr)
		}
		s, isStr := v.(string)
		if !isStr {
			t.Fatalf("task attr is %T", v)
		}
		seen[s]++
	}
	for _, id := range ids {
		if seen[id] != 1 {
			t.Errorf("task %s appears %d times, want exactly 1 (cross-talk or self-overwrite)", id, seen[id])
		}
	}
	if len(seen) != n {
		t.Errorf("distinct task ids in records = %d, want %d", len(seen), n)
	}
	t.Logf("PROBE2-153 CONCURRENCY: %d goroutines / 1 shared compressor -> %d records, %d distinct ids, every id exactly once", n, len(hits), len(seen))
}
