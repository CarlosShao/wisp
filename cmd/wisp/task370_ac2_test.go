package main

// Ticket 370 AC#2, both halves, on ONE real run through the shipped
// composition root:
//
//	(1) task_log.state is "interrupted" for a task the provider cut at the
//	    output ceiling;
//	(2) the word printed on screen, recorded VERBATIM.
//
// AC#2's subject is not "a config field exists" - it is what an operator is
// told and what memory keeps when the answer hits max_tokens. Ticket 365 could
// not write this test: [llm.roles.chat].max_output_tokens stopped at the
// resolver, Request.MaxOutputTokens was always 0 and no real run could reach
// the ceiling, so 365 pinned the word at the ladder and recorded the gap in
// Test365TruncatedCannotBeReachedThroughThisCommand. Ticket 370 landed that
// hop (run.go hands the chat role's ceiling to agent.Config, the loop puts it
// on the Assembler, the Assembler puts it on every request as max_tokens).
// This file is the run 365 said would be written the day the hop existed:
// nothing here manufactures a Result - the ceiling is a config key, the
// truncation is the mock provider cutting the echo at max_tokens, and the row
// is read back through memory's own DAO.
//
// Test370TheInterruptedRowNeedsTheCeiling is the reverse criterion. It measures
// the same task with the wire hop removed, and judges it with literally the
// same AC#2(1) predicate: that predicate must be FALSE there. Without the
// control, a green (1) would only say "some run rows interrupted", not "this
// hop is what rows it interrupted".

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
)

const (
	// t370ContextWindow keeps the frozen C22 budget at its 200k reference
	// (budgets.go scales 200000 by context_window/128000), so the run cannot
	// brake on TOKENS first. The two endings are different rows on purpose -
	// stuck vs interrupted - and this test is about the second one.
	t370ContextWindow = 128000

	// t370Ceiling is the smallest ceiling the config surface accepts. mockllm
	// emulates a real chat API here (ticket 386): with max_tokens=1 it cuts the
	// echo at roughly four bytes and reports finish_reason "length", which the
	// C5 chat adapter normalizes to llm.StopMaxTokens - the same stop a real
	// provider sends when the output budget runs out mid-sentence.
	t370Ceiling = 1

	// t370CeilingNeverHit is the second control: the request DOES carry
	// max_tokens (the hop is present), but the ceiling is far above anything
	// the echo can produce, so the provider reports a normal stop. It separates
	// "the field reached the wire" from "the ceiling was actually hit".
	t370CeilingNeverHit = 5000

	// t370TaskText is the same request in every run below, so the only variable
	// between them is whether the ceiling reached the wire.
	t370TaskText = "把这份笔记总结成三行"
)

// t370WithChatCeiling adds the one key ticket 370 wired to the probe's
// config.toml: [llm.roles.chat].max_output_tokens. It appends a table to the
// bytes newProbe365 wrote rather than rebuilding that config here, so the probe
// keeps owning the rest of the stack (provider, secret store, fs allow-list)
// and a change there cannot silently drift out of this test.
//
// n == 0 writes nothing at all: that is the "no hop" control, and its config
// stays byte-identical to the one ticket 365 already runs.
func t370WithChatCeiling(t *testing.T, p *probe365, n int) {
	t.Helper()
	if n == 0 {
		return
	}
	path := filepath.Join(p.dir, configFileName)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the probe's config: %v", err)
	}
	body = append(body, []byte(fmt.Sprintf("\n[llm.roles.chat]\nmax_output_tokens = %d\n", n))...)
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatalf("writing the probe's config: %v", err)
	}
}

// t370Interrupted is AC#2(1) as a function, so the positive run and both
// controls are judged by literally the same statement instead of by three
// paraphrases of it.
func t370Interrupted(rowState string) bool {
	return rowState == string(memory.TaskStateInterrupted)
}

// Test370TruncatedRunRowsAnInterruptedTask is AC#2(1): a real `wisp run` whose
// provider stops at the output ceiling leaves task_log.state = "interrupted".
func Test370TruncatedRunRowsAnInterruptedTask(t *testing.T) {
	p := newProbe365(t, t370ContextWindow)
	t370WithChatCeiling(t, p, t370Ceiling)

	code := p.run(t370TaskText)
	screenWord, rowState := p.surfaces()

	// AC#2(2), verbatim. Today the screen says "truncated" while the row says
	// "interrupted" - ticket 196 keeps the loop's vocabulary and memory's
	// roster on purpose (memory spells the ceiling ending "interrupted", and
	// "truncated" is not on the roster). Recorded as the measured pair, not
	// asserted to be one word.
	t.Logf("AC#2(2) 屏幕逐字 = %q；库 task_log.state 逐字 = %q", screenWord, rowState)

	// The premise of the measurement, checked before the claim: if this run did
	// not end at the ceiling, its row cannot testify about what the ceiling
	// rows, and a row-shaped assertion would be reading the wrong ending.
	if screenWord != string(agent.StatusTruncated) {
		t.Fatalf("this run did not end at the output ceiling (screen word %q, exit %d), so its row "+
			"is not evidence about truncation\nstdout:\n%s\nstderr:\n%s",
			screenWord, code, p.out.String(), p.errb.String())
	}
	if code != 3 {
		t.Errorf("exit = %d, want 3: runExitCode gives truncated the brake code it gives stuck", code)
	}

	// AC#2(1).
	if !t370Interrupted(rowState) {
		t.Errorf("task_log.state = %q for a task the provider cut at max_tokens, want %q: "+
			"the screen said %q and the row says %q",
			rowState, string(memory.TaskStateInterrupted), screenWord, rowState)
	}
}

// Test370TheInterruptedRowNeedsTheCeiling is the reverse criterion, in two
// directions: the AC#2(1) predicate must be FALSE in each of them.
//
// The runs are real `wisp run` invocations of the same task over the same
// composition root. Neither asserts that "done" is desirable - it asserts that
// the interrupted row above is caused by this hop, which is what makes that
// assertion able to go red.
func Test370TheInterruptedRowNeedsTheCeiling(t *testing.T) {
	t.Run("no_ceiling_on_the_wire", func(t *testing.T) {
		// Ticket 365's own config, untouched: max_output_tokens is 0, so no
		// request carries max_tokens and the provider cannot be cut at it.
		p := newProbe365(t, t370ContextWindow)
		code := p.run(t370TaskText)
		screenWord, rowState := p.surfaces()
		t.Logf("摘掉 wire 上那一跳：屏幕逐字 = %q；库 task_log.state 逐字 = %q（exit %d）",
			screenWord, rowState, code)

		if t370Interrupted(rowState) {
			t.Fatalf("AC#2(1) is TRUE without the hop: task_log.state = %q for a run whose request "+
				"carried no max_tokens, so the interrupted row above is not caused by the ceiling",
				rowState)
		}
		if rowState != string(memory.TaskStateDone) {
			t.Errorf("task_log.state = %q for an uncut answer, want %q", rowState, string(memory.TaskStateDone))
		}
		if code != 0 {
			t.Errorf("exit = %d, want 0 for an uncut answer\nstderr:\n%s", code, p.errb.String())
		}
	})

	t.Run("ceiling_never_hit", func(t *testing.T) {
		// The hop is present but the ceiling is out of reach, so the provider
		// reports a normal stop. This is the "摘掉 mock 的 StopMaxTokens" half:
		// the same config key, a different ending.
		p := newProbe365(t, t370ContextWindow)
		t370WithChatCeiling(t, p, t370CeilingNeverHit)
		code := p.run(t370TaskText)
		screenWord, rowState := p.surfaces()
		t.Logf("wire 上有 max_tokens 但撞不到：屏幕逐字 = %q；库 task_log.state 逐字 = %q（exit %d）",
			screenWord, rowState, code)

		if t370Interrupted(rowState) {
			t.Fatalf("AC#2(1) is TRUE for a run whose ceiling was never reached: task_log.state = %q, "+
				"so the positive run's row is not evidence that hitting the ceiling rows interrupted",
				rowState)
		}
		// The wire really did carry the ceiling: the request body is the product
		// of the hop under test, and reading it back is what separates "the key
		// was ignored" from "the ceiling was simply not hit".
		body, _ := p.srv.LastRequest(t, "chat")
		if !strings.Contains(body, `"max_tokens"`) {
			t.Fatalf("the request body carries no max_tokens field, so this control never exercised "+
				"the hop it is reversing: %s", body)
		}
		if rowState != string(memory.TaskStateDone) {
			t.Errorf("task_log.state = %q for an answer that never hit the ceiling, want %q",
				rowState, string(memory.TaskStateDone))
		}
		if code != 0 {
			t.Errorf("exit = %d, want 0\nstderr:\n%s", code, p.errb.String())
		}
	})
}
