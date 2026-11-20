package tools

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Ticket 191 cap-pin (audit follow-up, docs/evidence/s1/191-audit-r1.md
// section 2): sessionPendingCap's EXISTENCE was already pinned by
// Test191StreamChunkCapIsTheBufferCap, but its NUMBER was not - that test
// derives both its flood size and its bound from the constant itself, so
// growing the constant 64MiB -> 128MiB stayed green (the audit's mutation
// M2a, rc=0). A cap whose value can silently double is not pinned. Two
// teeth here close that gap:
//
//   - Test191SessionPendingCapNumberPinnedTo64MiB compares the constant to
//     the literal 64<<20 and names the consumer list on failure.
//     WHY 64MiB: the bound answers "one command streams faster than it
//     finishes" and stays in the same budget family as C22 (PLAN.md:1378 -
//     per-tool timeoutMs, per-task token budget): a runaway command costs
//     at most one 64MiB buffer before the session resets and refuses, one
//     order of magnitude under the task-level budget the same call may
//     already spend. It is also the number the panel side quotes
//     (internal/panel/terminal_carrier_191.go:55 calls this cap the tool
//     side's own 64MiB transport bound). Changing it means changing its
//     consumers (shell_session.go's sink guard, the search-loop cap, the
//     refusal TEXT that prints the number) plus that panel wording plus the
//     ticket rationale - all together, never alone.
//
//   - Test191FloodRefusedPastLiteral68157440Bytes derives NOTHING from the
//     constant: it floods 67108864 + 1048576 = 68157440 literal bytes (the
//     audit's own 64MiB + 1MiB dose, spelled as literals) and demands the
//     cap-path refusal, a refusal text naming 67108864, and delivery
//     staying inside 67108864. When the constant grows, every one of those
//     literals stays put and the test goes red.

// Test191SessionPendingCapNumberPinnedTo64MiB - the number nail: the cap is
// exactly 64<<20 bytes, not "whatever the constant says". The failure text
// carries the consumer list because the number is budget-sourced (C22), and
// a lone change here desynchronizes the refusal text the model reads and
// the panel-side wording that quotes this cap.
func Test191SessionPendingCapNumberPinnedTo64MiB(t *testing.T) {
	if sessionPendingCap != 64<<20 {
		t.Fatalf("sessionPendingCap moved off 64MiB: got %d bytes, want 67108864 (64<<20). "+
			"The number is C22-budget-sourced; if it moved, move its consumers with it - "+
			"shell_session.go's sink guard and search-loop cap and the refusal text that prints "+
			"the number, internal/panel/terminal_carrier_191.go:55's wording, the ticket rationale - "+
			"then update this nail and docs/evidence/s1/191-cap-pin-r1.md",
			sessionPendingCap)
	}
}

// Test191FloodRefusedPastLiteral68157440Bytes - the derived-from-literal
// arm: the dose is written as 67108864 + 1048576, never as
// sessionPendingCap + (1 << 20), so growing the constant cannot drag this
// test's expectations along. The flood never writes the marker line, so the
// only way out is the cap path: the call is REFUSED (cap refusal, not the
// deadline), the refusal text names the literal 67108864 bound, and the
// streaming consumer is TRUNCATED inside 67108864 bytes.
func Test191FloodRefusedPastLiteral68157440Bytes(t *testing.T) {
	total := 67108864 + 1048576 // literal bytes: 64MiB dose + 1MiB overshoot, no constant involved
	ff := &fakeShellFactory{respond: floodResponder(total)}
	d := shellDepsWith(ff)
	col := &streamCollector{}
	ctx, cancel := shellTaskCtx(60 * time.Second)
	defer cancel()
	rawJSON := fmt.Sprintf(`{"argv":%s}`, shellArgvJSON([]string{"type", "huge.log"}))
	res, err := (shellSession{d: d}).Execute(ctx, json.RawMessage(rawJSON), col.sink)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "单发上限") {
		t.Fatalf("68157440 literal bytes must take the cap refusal (the constant was 67108864): %s", res.Text)
	}
	if !strings.Contains(res.Text, "67108864") {
		t.Fatalf("the refusal text must name the literal 64MiB bound 67108864: %s", res.Text)
	}
	if col.bytes() > 67108864 {
		t.Fatalf("streamed chunks must stay inside the literal 67108864 cap: delivered=%d cap=67108864",
			col.bytes())
	}
	if strings.Contains(res.Text, "失去同步") {
		t.Fatalf("the cap path, not the deadline, must end a 68157440-byte flood: %s", res.Text)
	}
}
