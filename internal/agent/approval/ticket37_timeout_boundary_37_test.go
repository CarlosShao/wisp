package approval_test

// Ticket 37 cell 2's boundary ruler: the C18 timeout boundary the card and the
// audit wording carry is READ from the contract constants, not restated by a
// second copy. Two readings are pinned apart:
//
//   - the LITERAL frozen boundary: the queue's default deadline is exactly
//     300s and the lead 30s (C18 「超时 300s 一律判拒绝 / 超时前 30s 醒目提示」).
//     A "logically equivalent" retune of the constant to 299 turns this file
//     red - that is the tooth: the boundary is the constant's own value and
//     any change to it is a contract change that must fail loudly here;
//   - the DERIVED wording: a card expiring under a DIFFERENT injected timeout
//     spells THAT value in its refusal, proving the wording is computed from
//     the queue's own deadline and not a literal "300" in the format string.
//
// The ladder itself (warn@270 / reject@300 / replay) lives in
// ticket48_c18_full_test.go cell 2; this file only owns where the numbers come
// from.

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/tools"
)

func Test37TimeoutBoundaryIsReadFromTheContractConstants(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui)
	g.AdmitTextTask(testTask)

	// The frozen boundary, read off the queue and pinned literally.
	if got := g.Queue().Timeout(); got != 300*time.Second {
		t.Errorf("Queue().Timeout() = %s, want the frozen 300s C18 deadline", got)
	}
	if got := g.Queue().Timeout(); got != approval.DefaultApprovalTimeout {
		t.Errorf("Queue().Timeout() = %s, want DefaultApprovalTimeout (the card must carry the constant, not a second copy)", got)
	}
	if got := g.Queue().WarningLead(); got != 30*time.Second {
		t.Errorf("Queue().WarningLead() = %s, want the frozen 30s lead", got)
	}
	if approval.DefaultApprovalWarning != 30*time.Second {
		t.Errorf("DefaultApprovalWarning = %s, want the frozen 30s lead", approval.DefaultApprovalWarning)
	}

	// The card carries the constant's value as its deadline.
	ansCh := ask48(t, g, l2Decision48(testTask, "corr37-boundary"))
	waitFor(t, func() bool { return g.Queue().Depth() == 1 }, "the boundary card never queued")
	p := ui.wait(t)
	if p.Deadline != approval.DefaultApprovalTimeout {
		t.Errorf("the displayed card's Deadline = %s, want DefaultApprovalTimeout (%s)", p.Deadline, approval.DefaultApprovalTimeout)
	}

	// The expiry wording spells the SAME constant's seconds.
	clk.Advance(approval.DefaultApprovalTimeout)
	a := mustAnswer(t, ansCh)
	want := fmt.Sprintf("审批超时（%d 秒未确认）", int(approval.DefaultApprovalTimeout.Seconds()))
	if !strings.Contains(a.why, want) {
		t.Errorf("the expiry why = %q, want the wording derived from the frozen constant (%q)", a.why, want)
	}
}

// Test37ExpiryWordingDerivesFromTheInjectedTimeout is the derivation half: a
// queue built with a DIFFERENT (injected) timeout spells THAT value - the
// wording is computed from the queue's own deadline, so a future "300" hard
// coded into the format string turns this red.
func Test37ExpiryWordingDerivesFromTheInjectedTimeout(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{ApprovalTimeout: 45 * time.Second})
	g.AdmitTextTask(testTask)

	ansCh := ask48(t, g, l2Decision48(testTask, "corr37-derived"))
	waitFor(t, func() bool { return g.Queue().Depth() == 1 }, "the derived-wording card never queued")
	ui.wait(t)
	clk.Advance(45 * time.Second)
	a := mustAnswer(t, ansCh)
	if a.a != tools.AnswerTimeout {
		t.Errorf("the 45s expiry resolved to %s, want the refusing timeout", a.a)
	}
	if !strings.Contains(a.why, "45 秒未确认") {
		t.Errorf("the expiry why = %q, want the injected 45s value spelled out (the wording must be computed, not a literal)", a.why)
	}
}
