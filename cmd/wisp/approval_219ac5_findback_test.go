//go:build windows

package main

// Ticket 219 AC#5's last inch: a truncated reason's original must have a SECOND
// copy, findable through the mark's own pointer (票 147's same-family discipline
// for cuts). The 2026-11-15 receipt measured the hole precisely: dropped was
// only a count, the mark said 原文共 200 字, and the 40 cut runes had no home at
// all. This file pins the fix with three BEHAVIORAL teeth (A1880) - all three
// drive a real assembled run through the console reply pipe, not a source
// census:
//
//	Test219AC5TruncatedReasonIsFindableThroughTheMark
//	    200 runes in. The mark still stands with its dropped count (⛔ the
//	    pointer must never become an excuse to drop the mark), the mark now
//	    carries the pointer (答复理由档第 N 行), the `decisions` verb retrieves
//	    the verbatim 200-rune original, and the档's line N holds the same text
//	    with the same counts. Delete the retrieval plumbing (recordWithOriginal
//	    or decisions' 原文全文 line) => red on the verbatim retrieval; delete
//	    the archive append => red on the档 line assertions; delete the mark =>
//	    red on (1).
//	Test219AC5ShortReasonLeavesNoSecondCopy
//	    AC#5's 口径①: a reason within the ceiling grows NO second copy - the
//	    档 file must not even exist, and the review line carries no 原文全文
//	    field. Delete the "only when truly truncated" judgment in
//	    archivePointer (or in washReply) => the file appears => red.
//	Test219AC5EmptyReasonStaysDeterminate
//	    AC#5's 口径③: an empty reason keeps today's exact shape - the queue's
//	    default cause to the model, audit detail="", review （未填理由）, and
//	    no copy, no pointer, no档 line. Any "helpful" change to the empty
//	    semantics (inventing text, archiving empties) => red.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
)

// Test219AC5TruncatedReasonIsFindableThroughTheMark is tooth ①: the mark's
// pointer retrieves the verbatim original - through the operator's own read
// entry (`decisions`) and through the durable档 line it names.
func Test219AC5TruncatedReasonIsFindableThroughTheMark(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	pr, pw := newPipe219(t)
	h.reply = pr
	target := filepath.Join(h.outside, "219-findback.txt")
	const corr = "t219-findback-corr"
	long := strings.Repeat("长", 200)
	var (
		out    agent.ToolOutcome
		outErr error
		done   = make(chan struct{})
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t219-findback")
		defer revoke()
		card := startWriteCard219(rt, "t219-findback", corr, target)
		if !waitForCard187(rt.liveCards, corr, 5*time.Second) {
			t.Errorf("the L2 card never reached the native ledger")
			return
		}
		if _, err := fmt.Fprintln(pw, "no "+corr+" "+long); err != nil {
			t.Errorf("writing the reply: %v", err)
			return
		}
		<-card.done
		out, outErr = card.out, card.err
		if _, err := fmt.Fprintln(pw, "decisions"); err != nil {
			t.Errorf("asking for the review list: %v", err)
			return
		}
		if !waitForConsole219(h, "原文全文=", 10*time.Second) {
			t.Errorf("the review list never printed the archived original:\n%s", h.out.String())
		}
		close(done)
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if outErr != nil {
		t.Fatalf("bridge.Execute: %v", outErr)
	}
	// (1) The mark still stands, with both numbers: the pointer is ADDITIVE,
	// never a replacement for the visible cut (AC#5's hard constraint).
	if !strings.Contains(out.Text, replyTruncationMark) {
		t.Fatalf("the model-visible refusal lost the truncation mark: %q", out.Text)
	}
	for _, want := range []string{"40 字", "原文共 200 字"} {
		if !strings.Contains(out.Text, want) {
			t.Errorf("the mark does not name %q: %q", want, out.Text)
		}
	}
	// (2) The pointer rides the same mark, and names a real档 line: a fresh
	// data root means this run's truncation is the档's line 1.
	if !strings.Contains(out.Text, "原文另存：") || !strings.Contains(out.Text, "第 1 行") {
		t.Errorf("the mark carries no resolvable pointer: %q", out.Text)
	}
	// (3) The retrieval entry AC#5 names: `decisions` prints the verbatim
	// original. 200 consecutive 长 appear ONLY here - the detail line carries
	// the cut 160, so this assertion is exactly the retrieval tooth.
	console := h.out.String()
	if !strings.Contains(console, long) {
		t.Error("the decisions read entry does not carry the verbatim 200-rune original: "+
			"the retrieval was deleted or broken")
	}
	if !strings.Contains(console, "原文全文=") {
		t.Error("the review line grew no 原文全文 field: the entry lost its second copy")
	}
	// (4) The durable second copy: the档 holds exactly one line - the only
	// truncation this run performed - and it is the original verbatim, with the
	// counts the mark claims.
	lines, err := readReplyReasonArchive(h.dir)
	if err != nil {
		t.Fatalf("reading the reason档: %v", err)
	}
	if len(lines) != 1 {
		t.Fatalf("the档 has %d lines, want exactly 1 (口径①: only a real truncation writes)", len(lines))
	}
	l := lines[0]
	if l.Seq != 1 {
		t.Errorf("档 line seq = %d, want 1 (the pointer in the mark says 第 1 行)", l.Seq)
	}
	if l.Corr != corr {
		t.Errorf("档 line corr = %q, want %q", l.Corr, corr)
	}
	if l.Dropped != 40 || l.OriginalRunes != 200 {
		t.Errorf("档 line counts dropped=%d runes=%d, want 40/200 (must match the mark's numbers)", l.Dropped, l.OriginalRunes)
	}
	if l.Original != long {
		t.Error("the档 line is not the verbatim original: the second copy lied about being full")
	}
}

// Test219AC5ShortReasonLeavesNoSecondCopy is tooth ②: a reason within the
// ceiling must not produce the second copy anywhere - no档 line, no 原文全文
// field. Deleting the "only when truly truncated" judgment makes this red.
func Test219AC5ShortReasonLeavesNoSecondCopy(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	pr, pw := newPipe219(t)
	h.reply = pr
	target := filepath.Join(h.outside, "219-no-copy.txt")
	const corr = "t219-noshort-corr"
	short := "理由够短，不需要另存一份"
	var (
		out    agent.ToolOutcome
		outErr error
		done   = make(chan struct{})
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t219-noshort")
		defer revoke()
		card := startWriteCard219(rt, "t219-noshort", corr, target)
		if !waitForCard187(rt.liveCards, corr, 5*time.Second) {
			t.Errorf("the L2 card never reached the native ledger")
			return
		}
		if _, err := fmt.Fprintln(pw, "no "+corr+" "+short); err != nil {
			t.Errorf("writing the reply: %v", err)
			return
		}
		<-card.done
		out, outErr = card.out, card.err
		if _, err := fmt.Fprintln(pw, "decisions"); err != nil {
			t.Errorf("asking for the review list: %v", err)
			return
		}
		if !waitForConsole219(h, "理由/明细=", 10*time.Second) {
			t.Errorf("the review list never reached the console:\n%s", h.out.String())
		}
		close(done)
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if outErr != nil {
		t.Fatalf("bridge.Execute: %v", outErr)
	}
	// The reply itself worked, so the absences below are about the copy, not a
	// broken path.
	if !strings.Contains(out.Text, short) {
		t.Fatalf("the model lost the within-ceiling reason: %q", out.Text)
	}
	if strings.Contains(out.Text, replyTruncationMark) || strings.Contains(out.Text, "原文另存") {
		t.Errorf("a within-ceiling reason was marked or pointed as truncated: %q", out.Text)
	}
	// 口径①: the second copy does not exist - not an empty file, no file at all.
	if _, err := os.Stat(replyReasonArchivePath(h.dir)); !os.IsNotExist(err) {
		t.Errorf("a reason within the ceiling grew a档 copy (stat err = %v): "+
			"the only-on-truncation judgment was deleted", err)
	}
	if console := h.out.String(); strings.Contains(console, "原文全文=") {
		t.Errorf("the review list grew an 原文全文 field for a within-ceiling reason:\n%s", console)
	}
}

// Test219AC5EmptyReasonStaysDeterminate is tooth ③: an empty reason keeps
// today's determinate shape AND grows no copy. Test219AC5EmptyReasonIsANamedState
// already pins the three readings; this tooth adds the archive side so a
// "helpful" change to the empty semantics cannot slip past silently.
func Test219AC5EmptyReasonStaysDeterminate(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	pr, pw := newPipe219(t)
	h.reply = pr
	target := filepath.Join(h.outside, "219-empty-nocopy.txt")
	const corr = "t219-empty-arch-corr"
	var (
		out    agent.ToolOutcome
		outErr error
		done   = make(chan struct{})
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t219-empty-arch")
		defer revoke()
		card := startWriteCard219(rt, "t219-empty-arch", corr, target)
		if !waitForCard187(rt.liveCards, corr, 5*time.Second) {
			t.Errorf("the L2 card never reached the native ledger")
			return
		}
		// A bare `no <id>`: nothing typed into the reason box.
		if _, err := fmt.Fprintln(pw, "no "+corr); err != nil {
			t.Errorf("writing the reply: %v", err)
			return
		}
		<-card.done
		out, outErr = card.out, card.err
		if _, err := fmt.Fprintln(pw, "decisions"); err != nil {
			t.Errorf("asking for the review list: %v", err)
			return
		}
		if !waitForConsole219(h, "（未填理由）", 10*time.Second) {
			t.Errorf("the review list never named the empty state:\n%s", h.out.String())
		}
		close(done)
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if outErr != nil {
		t.Fatalf("bridge.Execute: %v", outErr)
	}
	// (a) The model: the queue's own default cause, not an invented or emptied
	// reason, and no operator-facing marker leaked into it.
	if !strings.Contains(out.Text, "拒绝") {
		t.Errorf("the model-visible refusal lost its cause: %q", out.Text)
	}
	if strings.Contains(out.Text, "（未填理由）") || strings.Contains(out.Text, replyTruncationMark) {
		t.Errorf("an operator-facing marker leaked into the model's text: %q", out.Text)
	}
	// (b) The audit: the stored detail is empty - the truth is "nothing typed".
	if audit := h.err.String(); !strings.Contains(audit, `route="native/reject" detail=""`) {
		t.Errorf("the audit line did not record the empty reason as empty:\n%s", audit)
	}
	// (c) 口径③'s archive side: an empty reason grows NO copy and NO pointer.
	// The header's constant statement about the档 is legitimate; what must be
	// absent is the per-entry copy field and the档 line itself.
	if _, err := os.Stat(replyReasonArchivePath(h.dir)); !os.IsNotExist(err) {
		t.Errorf("an empty reason grew a档 copy (stat err = %v): the empty semantics changed", err)
	}
	if console := h.out.String(); strings.Contains(console, "原文全文=") {
		t.Errorf("an empty reason produced a copy field:\n%s", console)
	}
}
