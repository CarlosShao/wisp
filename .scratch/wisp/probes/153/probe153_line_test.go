package agent

import (
	"context"
	"testing"
)

// Ticket 153 AC#2/AC#3 probe (NOT a deliverable: it exists only in the
// out-of-repo snapshot). It runs the compressor with NO injected logger, so the
// record goes out slog.Default() and shows up verbatim in `go test -v` output.
// Purpose: one externally visible reading of the line, tagged and untagged side
// by side - the assertions elsewhere only ever look at the in-memory capture.

func TestProbe153PrintsTheRealLine(t *testing.T) {
	b := BudgetsFor(4096)
	c := NewCompressor(b, nil) // default logger on purpose
	hist := buildRoundHistory(6, 400)

	if _, rep, err := c.Compress(context.Background(), hist); err != nil || !rep.Ran {
		t.Fatalf("setup: untagged fold did not run (err=%v ran=%v)", err, rep.Ran)
	}
	t.Log("PROBE-A above: untagged call (no task in ctx)")

	tagged := withTraceTask(context.Background(), newTaskID())
	if _, rep, err := c.Compress(tagged, hist); err != nil || !rep.Ran {
		t.Fatalf("setup: tagged fold did not run (err=%v ran=%v)", err, rep.Ran)
	}
	t.Log("PROBE-B above: tagged call (id from newTaskID)")
}
