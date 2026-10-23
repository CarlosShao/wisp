package agent

// Ticket 365: the loop status -> task_log.state ladder, and the row it writes.
//
// The tooth that matters here is the DEFAULT arm. Every named status already
// answered with a word on ticket 196's roster before this ticket, so a table
// over the six constants would have been green either way; the ending nobody
// named - a status this build does not know - used to be rowed "done", i.e. an
// unnamed ending reported as a success. That is the case this file pins, and it
// is red against the pre-365 ladder.
//
// The second assertion is not about a word but about the column: every answer
// must be on memory's roster, because the DAO validates against it and the SQL
// CHECK (schema v4) declares the same set. A word off it cannot be written at
// all, so an off-roster answer here is a dead producer, not a wrong label.
//
// Test365JournalRows... is the path evidence, not the tooth: it drives the real
// loop against a real SQLite store through the real journal and reads the state
// column back, so "what this package writes into the row" is a measurement
// rather than a reading of the ladder.

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
)

func Test365TaskLogStateSpeaksThe196Roster(t *testing.T) {
	cases := []struct {
		status Status
		want   memory.TaskState
		why    string
	}{
		{StatusCompleted, memory.TaskStateDone, "the model answered and no call stayed open"},
		{StatusControl, memory.TaskStateDone, "a D11(1) control word ends with an answer"},
		{StatusStuck, memory.TaskStateStuck, "a C22 brake is exactly what memory calls stuck"},
		{StatusCancelled, memory.TaskStateCancelled, "the task ctx was cancelled"},
		{StatusTruncated, memory.TaskStateInterrupted, "memory spells the output ceiling as interrupted"},
		{StatusFailed, memory.TaskStateError, "a classified failure ended without an answer"},
		{Status(""), memory.TaskStateError, "the zero value is not a completion"},
		{Status("a-status-this-build-does-not-know"), memory.TaskStateError, "an unknown ending must fail closed"},
	}
	for _, c := range cases {
		got := TaskLogState(c.status)
		if got != c.want {
			t.Errorf("TaskLogState(%q) = %q, want %q (%s)", string(c.status), string(got), string(c.want), c.why)
		}
		if !memory.ValidTaskState(string(got)) {
			t.Errorf("TaskLogState(%q) = %q is off ticket 196's roster: the DAO and the SQL CHECK both refuse it",
				string(c.status), string(got))
		}
	}
}

// Test365TaskLogStateCoversEveryStatusTheLoopCanProduce keeps the table above
// honest as the enum grows: a status added to loop.go without an answer here
// lands in the default arm, and the default arm is "error". This test does not
// say which word is right; it says the ladder is defined for every status
// Status.Valid accepts, so nobody has to guess later whether the silence was
// deliberate.
func Test365TaskLogStateCoversEveryStatusTheLoopCanProduce(t *testing.T) {
	all := []Status{
		StatusCompleted, StatusStuck, StatusCancelled,
		StatusFailed, StatusTruncated, StatusControl,
	}
	for _, s := range all {
		if !s.Valid() {
			t.Fatalf("Status(%q).Valid() = false: the table below and loop.go's const block disagree", string(s))
		}
		if got := TaskLogState(s); string(got) == "" {
			t.Errorf("TaskLogState(%q) = empty state: an empty state is not a roster word", string(s))
		}
	}
}

// Test365JournalRowsTheRosterWordOnTheAgentPath reads this package's own write
// back out of a real task_log: the loop runs a golden fixture, the journal
// closes the row, and the state column is what this package put there.
//
// Three of the four statuses here are positive controls (an answer, a ceiling,
// a failure); the fourth is the ticket's defect - a C22 brake, which is how the
// loop reaches StatusStuck.
func Test365JournalRowsTheRosterWordOnTheAgentPath(t *testing.T) {
	cases := []struct {
		fixture string
		status  Status
		want    memory.TaskState
	}{
		{"text-reply", StatusCompleted, memory.TaskStateDone},
		{"repeat-echo", StatusStuck, memory.TaskStateStuck},
		{"max-tokens-toolcall", StatusTruncated, memory.TaskStateInterrupted},
		{"disconnect-mid-toolcall", StatusFailed, memory.TaskStateError},
	}
	for _, c := range cases {
		t.Run(c.fixture, func(t *testing.T) {
			dir := sealableTempDir124(t)
			store, err := memory.Open(dir,
				memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
			if err != nil {
				t.Fatalf("memory.Open: %v", err)
			}
			defer store.Close()

			h := newHarness(t, c.fixture, withJournal(store),
				withConfig(func(cfg *Config) { cfg.ArtifactsDir = filepath.Join(dir, "artifacts") }))
			res := h.run("把这件事做完")
			if res.Status != c.status {
				t.Fatalf("status = %q, want %q: this fixture cannot testify about %q",
					string(res.Status), string(c.status), string(c.want))
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			tl, err := store.TaskLogByID(ctx, res.TaskID)
			if err != nil {
				t.Fatalf("task_log for %s: %v", res.TaskID, err)
			}
			t.Logf("agent path wrote task_log.state = %q for status %q", tl.State, string(res.Status))
			if tl.State != string(c.want) {
				t.Errorf("task_log.state = %q for status %q, want %q",
					tl.State, string(res.Status), string(c.want))
			}
			if !memory.ValidTaskState(tl.State) {
				t.Errorf("task_log.state = %q is off the roster, yet the row exists: the write path and the roster disagree",
					tl.State)
			}
		})
	}
}
