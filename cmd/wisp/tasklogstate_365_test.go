package main

// Ticket 365: one status, one word - in the command and in the loop.
//
// The defect was two conversion ladders: cmd/wisp had no StatusStuck arm (so a
// braked task was rowed "error" while the screen said "stuck") and the loop's
// own ladder had no unknown arm (so an unnamed ending was rowed "done"). The
// fix is not "add the missing arms" - that would leave two tables to drift
// again - it is one ladder in agent.TaskLogState and a one-line delegation
// here. These tests hold both ends of that:
//
//	Test365TaskLogStateSpeaksOneWordOnBothPaths  the two paths, one word each
//	Test365RunRowsAStuckTaskAsStuck              a REAL run: screen and row
//	Test365RunKeepsItsWordForTheControls         REAL runs that must not move
//
// The run tests reach their status through the shipped composition root by
// configuring the model, not by manufacturing a Result: a context_window of 1
// shrinks the frozen 200k C22 token budget to 2 tokens, so the first response
// overdraws it and the loop ends through brakeStuck (StatusStuck). The run's
// row is read back through memory's own DAO.

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/secret"
)

// closingLineRe finds the CLI's own closing line, and statusWordRe pulls the
// status word out of it: "wisp run: 任务 <id> 结束（<word>，<rounds> 轮，...".
// That word is res.Status, i.e. what the operator is told the ending was.
var (
	closingLineRe = regexp.MustCompile(`(?m)^wisp run: 任务 \S+ 结束（.*$`)
	statusWordRe  = regexp.MustCompile(`任务 \S+ 结束（([^，]+)，`)
)

func Test365TaskLogStateSpeaksOneWordOnBothPaths(t *testing.T) {
	// Every Status the loop can produce, plus the two the ladder has to answer
	// without a constant: the zero value and a word from a future build.
	cases := []struct {
		status agent.Status
		want   string
	}{
		{agent.StatusCompleted, "done"},      // positive control: answered
		{agent.StatusControl, "done"},        // positive control: control word
		{agent.StatusCancelled, "cancelled"}, // positive control: host stopped it
		{agent.StatusFailed, "error"},        // positive control: no answer
		{agent.StatusStuck, "stuck"},         // the reported defect
		{agent.StatusTruncated, "interrupted"},
		{agent.Status(""), "error"},
		{agent.Status("a-status-this-build-does-not-know"), "error"},
	}
	for _, c := range cases {
		cmdWord := taskLogState(c.status)
		loopWord := string(agent.TaskLogState(c.status))
		if cmdWord != c.want {
			t.Errorf("cmd/wisp taskLogState(%q) = %q, want %q", string(c.status), cmdWord, c.want)
		}
		if loopWord != c.want {
			t.Errorf("agent.TaskLogState(%q) = %q, want %q", string(c.status), loopWord, c.want)
		}
		// AC#2's criterion in one line: not "some path changed", but "both
		// paths answer the same word for the same status".
		if cmdWord != loopWord {
			t.Errorf("the two paths disagree on %q: cmd/wisp writes %q, the loop writes %q",
				string(c.status), cmdWord, loopWord)
		}
		if !memory.ValidTaskState(cmdWord) {
			t.Errorf("taskLogState(%q) = %q is off ticket 196's roster: the row could not be written",
				string(c.status), cmdWord)
		}
	}
}

// probe365 is one `wisp run` over the real composition root with the model's
// window exposed. Everything else - config.toml, the DPAPI secret store, the
// SQLite store, the provider, the gate, the loop - is the production stack.
type probe365 struct {
	t    *testing.T
	srv  *adaptertest.Mockllm
	dir  string
	out  *bytes.Buffer
	errb *bytes.Buffer
}

func newProbe365(t *testing.T, contextWindow int) *probe365 {
	t.Helper()
	p := &probe365{t: t, srv: adaptertest.StartMockllm(t), out: &bytes.Buffer{}, errb: &bytes.Buffer{}}
	p.dir = t.TempDir()
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models.m1]
enabled = true
context_window = %d

[fs]
allowed_dirs = [%q]
`, p.srv.Base+"/v1", contextWindow, filepath.ToSlash(p.dir))
	if err := os.WriteFile(filepath.Join(p.dir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	store, err := secret.NewStore(p.dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}
	return p
}

// run drives the same entry point main() calls. taskCtx nil is what `wisp run`
// means: the process is the task's lifetime.
func (p *probe365) run(task string) int {
	p.t.Helper()
	return runTextTask(runSpec{
		argv:    []string{task},
		stdout:  p.out,
		stderr:  p.errb,
		dataDir: p.dir,
		notify:  func(string, string) error { return nil },
	})
}

// surfaces reads back the two claims about how this task ended: the word the
// command printed in its closing line, and the state memory holds in its row.
func (p *probe365) surfaces() (screenWord, rowState string) {
	p.t.Helper()
	line := closingLineRe.FindString(p.out.String())
	m := statusWordRe.FindStringSubmatch(p.out.String())
	if line == "" || m == nil {
		p.t.Fatalf("the run printed no closing status line:\nstdout:\n%s\nstderr:\n%s",
			p.out.String(), p.errb.String())
	}
	p.t.Logf("屏幕逐字: %s", line)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	store, err := memory.Open(p.dir)
	if err != nil {
		p.t.Fatalf("reopening the run's store: %v", err)
	}
	defer func() { _ = store.Close() }()
	id := taskIDRe.FindStringSubmatch(p.out.String())
	if id == nil {
		p.t.Fatalf("the status line carries no task id:\n%s", p.out.String())
	}
	tl, err := store.TaskLogByID(ctx, id[1])
	if err != nil {
		p.t.Fatalf("task_log for %s: %v", id[1], err)
	}
	p.t.Logf("库逐字: task_log.state = %q", tl.State)
	return m[1], tl.State
}

// Test365RunRowsAStuckTaskAsStuck is the defect itself, measured on one real
// run: the CLI used to print "（stuck，" and row state="error". Both surfaces
// must now name the same ending, and that ending is the roster word.
//
// A context_window of 1 scales the frozen 200k C22 token budget down to 2
// tokens, so the first response overdraws it and the loop ends through
// brakeStuck - StatusStuck, the status the ticket was filed about. Nothing is
// simulated: the brake fires, the loop finishes, this command writes the row.
func Test365RunRowsAStuckTaskAsStuck(t *testing.T) {
	p := newProbe365(t, 1)
	code := p.run("把目录再列一遍")
	screenWord, rowState := p.surfaces()

	if screenWord != string(agent.StatusStuck) {
		t.Fatalf("this run did not end stuck (screen word %q, exit %d), so it cannot testify about stuck:\nstderr:\n%s",
			screenWord, code, p.errb.String())
	}
	if code != 3 {
		t.Errorf("exit = %d, want 3 for a loop brake", code)
	}
	if rowState != string(memory.TaskStateStuck) {
		t.Errorf("task_log.state = %q for a stuck task, want %q: the screen said %q and the row says %q",
			rowState, string(memory.TaskStateStuck), screenWord, rowState)
	}
	// The ticket's own sentence: the screen and the row must not be two claims
	// about the same ending. This holds here because "stuck" is a word both
	// vocabularies spell the same way.
	if screenWord != rowState {
		t.Errorf("screen and row disagree about the same task: screen %q, row %q", screenWord, rowState)
	}
}

// Test365RunKeepsItsWordForTheControls pins the endings this ticket must NOT
// move: the answer, the dead provider, the cancelled task. All three are real
// runs through the same composition root; the two failures reach their status
// through different mechanisms (a provider fault ladder, a cancelled parent
// context), so neither is the other restated.
//
// The screen prints the loop's status word and the row prints the roster's -
// "completed"/"done", "failed"/"error", "cancelled"/"cancelled" - because
// ticket 196 keeps both vocabularies on purpose. What must not happen is the
// row changing to some third word.
func Test365RunKeepsItsWordForTheControls(t *testing.T) {
	t.Run("completed", func(t *testing.T) {
		p := newProbe365(t, 128000)
		code := p.run("总结一下 这份笔记讲了什么")
		screenWord, rowState := p.surfaces()
		if code != 0 {
			t.Fatalf("exit = %d, want 0\nstderr:\n%s", code, p.errb.String())
		}
		if screenWord != string(agent.StatusCompleted) {
			t.Errorf("screen word = %q, want %q", screenWord, string(agent.StatusCompleted))
		}
		if rowState != string(memory.TaskStateDone) {
			t.Errorf("task_log.state = %q, want %q", rowState, string(memory.TaskStateDone))
		}
	})
	t.Run("failed", func(t *testing.T) {
		p := newProbe365(t, 128000)
		p.srv.QueueFault(t, adaptertest.Fault{Status: 500, Times: 3})
		code := p.run("总结一下 这份笔记")
		screenWord, rowState := p.surfaces()
		if code != 1 {
			t.Errorf("exit = %d, want 1 for a classified provider failure", code)
		}
		if screenWord != string(agent.StatusFailed) {
			t.Errorf("screen word = %q, want %q", screenWord, string(agent.StatusFailed))
		}
		if rowState != string(memory.TaskStateError) {
			t.Errorf("task_log.state = %q, want %q", rowState, string(memory.TaskStateError))
		}
	})
	t.Run("cancelled", func(t *testing.T) {
		p := newProbe365(t, 128000)
		code := p.runWithCancelledParent("总结一下 这份笔记")
		screenWord, rowState := p.surfaces()
		if screenWord != string(agent.StatusCancelled) {
			t.Fatalf("this run did not end cancelled (screen word %q, exit %d):\nstderr:\n%s",
				screenWord, code, p.errb.String())
		}
		if rowState != string(memory.TaskStateCancelled) {
			t.Errorf("task_log.state = %q, want %q", rowState, string(memory.TaskStateCancelled))
		}
		if screenWord != rowState {
			t.Errorf("screen and row disagree about the same task: screen %q, row %q", screenWord, rowState)
		}
	})
}

// runWithCancelledParent hands the run a parent context that is already
// cancelled - the console's Ctrl-C, at the moment it is pressed.
func (p *probe365) runWithCancelledParent(task string) int {
	p.t.Helper()
	parent, cancel := context.WithCancel(context.Background())
	cancel()
	return runTextTask(runSpec{
		argv:    []string{task},
		stdout:  p.out,
		stderr:  p.errb,
		dataDir: p.dir,
		notify:  func(string, string) error { return nil },
		taskCtx: parent,
	})
}

// Test365TruncatedCannotBeReachedThroughThisCommand records a limit this ticket
// ran into rather than hid: `wisp run` never sends max_tokens, so no real run
// can end at the output ceiling. [llm.roles.chat].max_output_tokens is resolved
// into llm.Endpoint, but agent.Config has no field carrying it to the
// assembler, so Request.MaxOutputTokens is always 0 and a provider has no
// ceiling to hit. The word for that status is therefore pinned at the ladder
// (both tables above) and in the agent-path row test, not by a CLI run.
//
// This test does not assert the gap is acceptable; it asserts the gap is
// VISIBLE, so the day someone wires the field this test fails and the real run
// gets written.
func Test365TruncatedCannotBeReachedThroughThisCommand(t *testing.T) {
	if got := agent.TaskLogState(agent.StatusTruncated); string(got) != "interrupted" {
		t.Fatalf("TaskLogState(truncated) = %q, want interrupted: the ladder this note rests on changed", string(got))
	}
	if !strings.Contains(string(agent.StatusTruncated), "truncated") {
		t.Fatalf("StatusTruncated's spelling changed to %q: re-measure the roster split", string(agent.StatusTruncated))
	}
}
