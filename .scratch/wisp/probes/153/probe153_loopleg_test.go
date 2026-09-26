package agent

import (
	"testing"
)

// Ticket 153 AC#3 probe (NOT a deliverable; snapshot only).
//
// Sentence (2) of the load-bearing pair asks whether removing this ticket's
// pieces changes anything an outside reader can see. The tests elsewhere read
// an in-memory capture, so this one puts the record on the PROCESS DEFAULT
// logger (harness withLogger(nil) -> Options.Logger nil -> slog.Default()) on
// the real Loop.Run leg, and prints the id the run minted, so the line a log
// file would carry is visible verbatim in `go test -v` output.

func TestProbe153LoopLegLineOnDefaultLogger(t *testing.T) {
	h := newHarness(t, "text-reply",
		withConfig(func(c *Config) { c.ContextWindow = 4096 }),
		withLogger(nil)) // nil = the process default handler, i.e. the outside world
	for _, m := range buildRoundHistory(4, 400) {
		h.loop.append(m)
	}
	res := h.run("收个尾")
	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}
	t.Logf("PROBE-LOOPLEG this run's Result.TaskID = %q (the line above must carry exactly this id)", res.TaskID)
}
