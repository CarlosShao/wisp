//go:build windows

package main

// Ticket 219 AC#1 (★ 双击启动那条) and AC#5 (理由的三条去向 / 空理由 / 超长理由).
//
// WHAT THIS FILE IS FOR. The 2026-11-15 收件 (commit 9c80263f2) measured two holes
// that were PRODUCTION gaps, not missing readings:
//
//	AC#1 ★  resident_windows.go:568-572 attached the tray's three answer items
//	        only inside `if src != nil && src.run != nil`. startResidentTaskSource
//	        returns nil before assembling anything when interactiveStdin() has no
//	        console, so a double-clicked wisp - the shape D2 and ticket 228 both
//	        call the process the user actually launches - had a tray whose three
//	        items booked one "unhosted gesture" line and nothing else.
//	AC#5 ③  a reason over 160 runes was cut with no mark anywhere (measured:
//	        200 in / 160 out, zero hits for any marker), which AC#5 forbids.
//	AC#5 ②  an empty reason had no case at all.
//	AC#5 ①  the reason's third destination (the operator reading it back) had no
//	        landing point: `tool_call`'s thirteen frozen columns carry no reason
//	        and nothing read `approval: REPLY` back to the person who typed it.
//
// WHAT EACH CASE WOULD GO RED ON, so none of them is a decorative census:
//
//	Test219AC1TrayAllowOnceHasAnExecutorWithoutAPipeline
//	    delete the else branch in resident_windows.go (or make
//	    answerSurfaceWithoutPipeline return nil) and the allowed write never
//	    lands: the sentence becomes the "no answer side" refusal and the bridge
//	    call it was blocking runs out its L2 deadline.
//	Test219AC1TrayRejectHasAnExecutorWithoutAPipeline
//	    same mutation, other polarity: the file the rejected call would have
//	    written exists.
//	Test219AC1TrayAlwaysWritesTheWideningWithoutAPipeline
//	    delete either the else branch or the widenSeat plumbing and config.toml
//	    never gains the [fs] line - the widening cannot even stage.
//	Test219AC1TheNoPipelineExecutorIsWiredInTheShippedBoot
//	    the WIRING census (the assembly root needs a real boot to run, so this
//	    one reads the shipped statement). It is named as a census, not offered as
//	    proof of behaviour: the three cases above are the proof.
//	Test219AC5LongReasonIsTruncatedVisibly / Test219AC5TheModelAndTheOperatorSee…
//	    put the silent cut back (return only the first 160 runes) and the model's
//	    tool result stops carrying the marker, so both go red on the text the
//	    MODEL and the OPERATOR receive - not on this package's source text.
//	Test219AC5EmptyReasonIsANamedState
//	    make an empty reason indistinguishable from "no box" (or have the wash
//	    invent a sentence) and the three readings below disagree.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/config"
)

// cardExec219 is one in-flight fs.write this file puts a real L2 card behind. It
// is deliberately the same shape approval_reply_201_test.go and the ticket-228
// trays cases use: the gate blocks inside bridge.Execute until somebody answers,
// so the answer has to arrive on another goroutine.
type cardExec219 struct {
	out  agent.ToolOutcome
	err  error
	done chan struct{}
}

func startWriteCard219(rt *agentRuntime, taskID, corr, target string) *cardExec219 {
	c := &cardExec219{done: make(chan struct{})}
	args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"written because the tray answered without a console"}`,
		filepath.ToSlash(target)))
	go func() {
		defer close(c.done)
		c.out, c.err = rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: taskID, CorrelationID: corr, CallID: "a1", Name: "fs.write", Args: args,
		})
	}()
	return c
}

// trayExec219 builds the executor the NO-PIPELINE way: a residentApproval over
// the run's own gate and ledger (the same two objects a boot's ra holds), handed
// to the production constructor answerSurfaceWithoutPipeline rather than to the
// console's newReplySurface. The point of the case is that this constructor's
// output is a real executor - if it were nil, the tray items would fall back to
// the unhosted-gesture record and nothing would be answered.
func trayExec219(rt *agentRuntime) trayAnswerFunc {
	ra := &residentApproval{
		cards: rt.liveCards.h,
		gate:  rt.gate,
		root:  context.Background(),
	}
	return ra.trayAnswerExecutor(ra.answerSurfaceWithoutPipeline(nil, nil))
}

// Test219AC1TrayAllowOnceHasAnExecutorWithoutAPipeline is AC#1's ★ finding cut
// the other way: with NO pipeline-shaped answer surface at all, the tray's
// allow-once item still spends the card's single-use native grant and the call
// it was blocking really executes.
func Test219AC1TrayAllowOnceHasAnExecutorWithoutAPipeline(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	target := filepath.Join(h.outside, "219-tray-allow-without-pipeline.txt")
	var seamErr error

	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t219-nopipe-allow")
		defer revoke()
		exec := trayExec219(rt)
		if exec == nil {
			seamErr = fmt.Errorf("the no-pipeline constructor produced no executor, so the tray items have nothing behind them")
			return
		}
		card := startWriteCard219(rt, "t219-nopipe-allow", "t219-nopipe-allow-corr", target)
		if !waitForCard187(rt.liveCards, "t219-nopipe-allow-corr", 5*time.Second) {
			seamErr = fmt.Errorf("the L2 card never reached the host ledger, so no tray item could answer it")
			return
		}
		sentence := exec(trayAnswerAllowOnce)
		if !strings.Contains(sentence, "已允许") || !strings.Contains(sentence, "t219-nopipe-allow-corr") {
			seamErr = fmt.Errorf("tray allow-once sentence = %q, want the surface's own allow sentence naming the card "+
				"(the no-executor branch answers with the named refusal instead)", sentence)
			return
		}
		<-card.done
		if card.err != nil {
			seamErr = fmt.Errorf("bridge.Execute: %w", card.err)
			return
		}
		if card.out.IsError {
			seamErr = fmt.Errorf("a call the tray allowed without a pipeline must execute, got: %s", card.out.Text)
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if seamErr != nil {
		t.Fatal(seamErr)
	}
	body, err := os.ReadFile(target)
	if err != nil || !strings.Contains(string(body), "written because the tray answered without a console") {
		t.Fatalf("the no-pipeline tray allow never reached the write: %v %q", err, body)
	}
	if audit := h.err.String(); !strings.Contains(audit, "ANSWER-ALLOW corr=t219-nopipe-allow-corr") {
		t.Errorf("the answer did not travel the native route (no ANSWER-ALLOW line); audit:\n%s", audit)
	}
}

// Test219AC1TrayRejectHasAnExecutorWithoutAPipeline is the same hole, other
// polarity: with no pipeline the tray's 拒绝 item really refuses, and the file
// the call would have written does not exist afterwards.
func Test219AC1TrayRejectHasAnExecutorWithoutAPipeline(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	target := filepath.Join(h.outside, "219-tray-reject-without-pipeline.txt")
	var seamErr error

	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t219-nopipe-reject")
		defer revoke()
		exec := trayExec219(rt)
		card := startWriteCard219(rt, "t219-nopipe-reject", "t219-nopipe-reject-corr", target)
		if !waitForCard187(rt.liveCards, "t219-nopipe-reject-corr", 5*time.Second) {
			seamErr = fmt.Errorf("the L2 card never reached the host ledger")
			return
		}
		sentence := exec(trayAnswerReject)
		if !strings.Contains(sentence, "已拒绝") {
			seamErr = fmt.Errorf("tray reject sentence = %q, want the surface's own refusal sentence", sentence)
			return
		}
		<-card.done
		if card.err != nil {
			seamErr = fmt.Errorf("bridge.Execute: %w", card.err)
			return
		}
		if !card.out.IsError {
			seamErr = fmt.Errorf("a call the tray refused must not execute, got: %s", card.out.Text)
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if seamErr != nil {
		t.Fatal(seamErr)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("the refused call wrote its file anyway (stat err = %v)", err)
	}
	// The tray's own reason sentence reaches the record: it is a reason the
	// operator never typed, and it says so. The line is the GATE's (ANSWER-REJECT,
	// internal/agent/approval/queue.go:528), booked through this run's audit sink -
	// what proves the decision travelled the native route the tray executor names.
	if audit := h.err.String(); !strings.Contains(audit, "ANSWER-REJECT corr=t219-nopipe-reject-corr") {
		t.Errorf("the no-pipeline reject never reached the native route; audit:\n%s", audit)
	}
}

// Test219AC1TrayAlwaysWritesTheWideningWithoutAPipeline is AC#1's 「真写长期允许」
// in the shape that has no runtime at all. The 长期允许 item must reach
// replySurface.always (which stages, then raises its OWN second L2 card - the
// PLAN.md:1645-1646 re-confirmation), and approving that second card must really
// write the [fs] line into config.toml.
func Test219AC1TrayAlwaysWritesTheWideningWithoutAPipeline(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	target := filepath.Join(h.outside, "219-tray-always-without-pipeline.txt")
	var durable string

	var seamErr error
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t219-nopipe-always")
		defer revoke()
		// The 长期 branch's write door, in the no-pipeline shape: the boot's own
		// config Manager. Same file the run loaded, a second Manager over it,
		// exactly what cfgMgr65 is relative to the pipeline's rt.mgr.
		mgr, err := config.NewManager(filepath.Join(h.dir, configFileName), nil)
		if err != nil {
			seamErr = fmt.Errorf("config.NewManager: %w", err)
			return
		}
		ra := &residentApproval{cards: rt.liveCards.h, gate: rt.gate, root: context.Background()}
		exec := ra.trayAnswerExecutor(ra.answerSurfaceWithoutPipeline(mgr, nil))
		card := startWriteCard219(rt, "t219-nopipe-always", "t219-nopipe-always-corr", target)
		if !waitForCard187(rt.liveCards, "t219-nopipe-always-corr", 5*time.Second) {
			seamErr = fmt.Errorf("the L2 card never reached the host ledger")
			return
		}
		sentence := exec(trayAnswerAllowAlways)
		if !strings.Contains(sentence, "要存的规则") {
			seamErr = fmt.Errorf("the 长期允许 item did not stage a rule, sentence = %q", sentence)
			return
		}
		// The second card: always() raises it off the ui thread, so it arrives a
		// moment later on the same ledger.
		widenCorr, ok := waitForWideningCard219(rt.liveCards, 10*time.Second)
		if !ok {
			seamErr = fmt.Errorf("no second L2 card was raised: the widening was staged but never re-confirmed")
			return
		}
		if err := rt.liveCards.h.Allow(context.Background(), widenCorr); err != nil {
			seamErr = fmt.Errorf("approving the widening card: %w", err)
			return
		}
		if _, ok := waitForDurable219(h.dir, "WIDEN-APPLIED", 15*time.Second); !ok {
			text, _ := waitForDurable219(h.dir, "WIDEN-APPLIED", 0)
			durable = text
			seamErr = fmt.Errorf("the second card was approved but no WIDEN-APPLIED line reached the durable log; "+
				"the log said:\n%s", durable)
			return
		}
		// The original card is still waiting (always() does not answer it), so the
		// run could not finish without this.
		if err := rt.liveCards.h.Allow(context.Background(), "t219-nopipe-always-corr"); err != nil {
			seamErr = fmt.Errorf("answering the original card after the widening: %w", err)
			return
		}
		<-card.done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if seamErr != nil {
		t.Fatal(seamErr)
	}
	// The reading that matters: config.toml really carries the line.
	body, err := os.ReadFile(filepath.Join(h.dir, configFileName))
	if err != nil {
		t.Fatalf("reading config.toml: %v", err)
	}
	want := filepath.ToSlash(h.outside)
	if !strings.Contains(string(body), want) {
		t.Fatalf("the no-pipeline 长期允许 never reached [fs] allowed_dirs (want %q); config.toml:\n%s", want, body)
	}
	// And the durable record says which line it wrote, not merely that it wrote.
	if text, ok := waitForDurable219(h.dir, "WIDEN-APPLIED", time.Second); !ok ||
		!strings.Contains(text, want) || !strings.Contains(text, "key=allowed_dirs") {
		t.Errorf("the durable WIDEN-APPLIED line does not name the line and the key it wrote (want %q):\n%s",
			want, text)
	}
}

// Test219AC1TheNoPipelineExecutorIsWiredInTheShippedBoot is the WIRING census.
// runResident needs a real desktop (ball window, WebView2 host, event loop), so
// this half of AC#1 is read off the shipped statement instead of executed - and
// it is named as a census rather than offered as proof: the three cases above are
// what would go red if the behaviour broke. What this adds is the one claim they
// cannot make: the boot's `if src != nil && src.run != nil` has a branch for the
// other shape.
func Test219AC1TheNoPipelineExecutorIsWiredInTheShippedBoot(t *testing.T) {
	body, err := os.ReadFile("resident_windows.go")
	if err != nil {
		t.Fatalf("read resident_windows.go: %v", err)
	}
	src := string(body)
	// The console shape keeps its own line (ticket 228's pin, unchanged).
	if !strings.Contains(src, "rb.attachTrayAnswers(ra.trayAnswerExecutor(src.answers))") {
		t.Error("the console shape lost its tray executor line")
	}
	// The no-console shape has one too, and it is the production constructor.
	for _, want := range []string{
		"rb.attachTrayAnswers(ra.trayAnswerExecutor(ra.answerSurfaceWithoutPipeline(cfgMgr65, residentAnswerMode(rp))))",
		"residentAnswerMode(rp)",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the shipped resident boot is missing %q: the double-click shape would fall back to the "+
				"unhosted-gesture record, which is ticket 219 AC#1's measured defect", want)
		}
	}
}

// waitForWideningCard219 polls the ledger for the second card always() raises.
// The correlation id is minted by the gate, so the card is found by the tool name
// the widening uses (allowWidenTool), which is exactly how an operator would know
// it: the card says what it is widening.
func waitForWideningCard219(live *nativeCards, limit time.Duration) (string, bool) {
	deadline := time.Now().Add(limit)
	for time.Now().Before(deadline) {
		for _, c := range live.pending() {
			if c.Tool == allowWidenTool {
				return c.CorrelationID, true
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	return "", false
}

// durableLog219 concatenates every wisp-*.jsonl file this run wrote under
// <dataDir>/logs. That is the copy that outlives the process, and the one the
// receipt's "audit line" means; the console copy of the same sentence goes to
// stdout, which a test binary cannot read back.
func durableLog219(dir string) string {
	logDir := filepath.Join(dir, "logs")
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return ""
	}
	var b strings.Builder
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "wisp-") || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(logDir, e.Name()))
		if err != nil {
			continue
		}
		b.Write(data)
	}
	return b.String()
}

// waitForDurable219 polls that log for one needle. limit 0 means "read once".
func waitForDurable219(dir, needle string, limit time.Duration) (string, bool) {
	deadline := time.Now().Add(limit)
	for {
		text := durableLog219(dir)
		if strings.Contains(text, needle) {
			return text, true
		}
		if limit <= 0 || !time.Now().Before(deadline) {
			return text, false
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// Test219AC5LongReasonIsTruncatedVisibly is AC#5 ③ at the boundary: the ceiling
// stays, and what the ceiling removes is COUNTED and STATED. Before this leg a
// 200-rune reason came out as 160 runes with nothing anywhere saying so.
func Test219AC5LongReasonIsTruncatedVisibly(t *testing.T) {
	long := strings.Repeat("长", 200)
	text, dropped := washReplyText(long)
	if dropped != 40 {
		t.Errorf("dropped = %d, want 40 (200 typed - 160 ceiling)", dropped)
	}
	if !strings.Contains(text, replyTruncationMark) {
		t.Fatalf("the cut reason carries no truncation mark: %q", text)
	}
	for _, want := range []string{"40 字", "原文共 200 字"} {
		if !strings.Contains(text, want) {
			t.Errorf("the mark does not name %q: %q", want, text)
		}
	}
	if strings.Contains(text, long) {
		t.Error("the washed text still carries all 200 runes, so nothing was actually bounded")
	}
	if n := len([]rune(text)); n <= replyReasonMax {
		t.Errorf("the mark was charged against the ceiling (total runes %d <= %d): a mark that can itself be "+
			"cut off is the same silence one layer deeper", n, replyReasonMax)
	}

	// Exactly at the ceiling: no mark, nothing dropped. A wash that marked a
	// 160-rune reason would make every long-but-legal reason look cut.
	exact := strings.Repeat("短", replyReasonMax)
	textAt, droppedAt := washReplyText(exact)
	if droppedAt != 0 || strings.Contains(textAt, replyTruncationMark) {
		t.Errorf("a %d-rune reason was reported as truncated (dropped=%d): %q", replyReasonMax, droppedAt, textAt)
	}
	if textAt != exact {
		t.Errorf("a legal reason was altered by the wash: %q", textAt)
	}

	// And the note a caller adds names the same two numbers.
	note := replyReason{text: text, dropped: dropped}.truncationNote()
	if !strings.Contains(note, "40") || !strings.Contains(note, "160") {
		t.Errorf("truncationNote = %q, want both the dropped count and the ceiling", note)
	}
	if (replyReason{text: "短"}).truncationNote() != "" {
		t.Error("a reason within the ceiling produced a truncation note")
	}
}

// Test219AC5TheModelAndTheOperatorSeeTheTruncation is the end-to-end half: the
// two readers AC#5 names - the model (through the tool result) and the operator
// (through the console sentence and the review list) - both hear that the reason
// was cut. The audit line is the third reader and is asserted too.
func Test219AC5TheModelAndTheOperatorSeeTheTruncation(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	pr, pw := newPipe219(t)
	h.reply = pr
	target := filepath.Join(h.outside, "219-long-reason.txt")
	long := strings.Repeat("长", 200)
	var (
		out    agent.ToolOutcome
		outErr error
		done   = make(chan struct{})
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t219-long")
		defer revoke()
		card := startWriteCard219(rt, "t219-long", "t219-long-corr", target)
		if !waitForCard187(rt.liveCards, "t219-long-corr", 5*time.Second) {
			t.Errorf("the L2 card never reached the native ledger")
			return
		}
		if _, err := fmt.Fprintln(pw, "no t219-long-corr "+long); err != nil {
			t.Errorf("writing the reply: %v", err)
			return
		}
		<-card.done
		out, outErr = card.out, card.err
		if _, err := fmt.Fprintln(pw, "decisions"); err != nil {
			t.Errorf("asking for the review list: %v", err)
			return
		}
		// The listener prints from its own goroutine: a case that returned here
		// would race the process exit and read an empty console.
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
	// (1) The model's copy: the reason it reads says it was cut, and names how.
	if !strings.Contains(out.Text, replyTruncationMark) {
		t.Fatalf("the model-visible refusal carries no truncation mark: %q", out.Text)
	}
	if !strings.Contains(out.Text, "40 字") || !strings.Contains(out.Text, "原文共 200 字") {
		t.Errorf("the model-visible text does not name what was dropped: %q", out.Text)
	}
	if strings.Contains(out.Text, long) {
		t.Error("the model received all 200 runes, so the ceiling did not apply")
	}
	// (2) The operator's copy, on the console: the answer sentence says it too.
	if console := h.out.String(); !strings.Contains(console, "已截断") || !strings.Contains(console, "40 字") {
		t.Errorf("the console never told the operator the reason was cut:\n%s", console)
	}
	// (3) The audit copy.
	if audit := h.err.String(); !strings.Contains(audit, replyTruncationMark) {
		t.Errorf("the audit REPLY line does not carry the mark:\n%s", audit)
	}
	// (4) AC#5 ①'s landing point: the reason is readable back, with the mark.
	if console := h.out.String(); !strings.Contains(console, "t219-long-corr") || !strings.Contains(console, "理由/明细=") {
		t.Errorf("the review list did not show the answered card and its reason:\n%s", console)
	}
}

// Test219AC5EmptyReasonIsANamedState is AC#5 ②: an empty reason is a STATE with
// three agreeing readings - the model is told the default refusal cause (not a
// fabricated reason), the console says 未填理由, and the review list says the
// same word instead of printing an empty string.
func Test219AC5EmptyReasonIsANamedState(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	pr, pw := newPipe219(t)
	h.reply = pr
	target := filepath.Join(h.outside, "219-empty-reason.txt")
	var (
		out    agent.ToolOutcome
		outErr error
		done   = make(chan struct{})
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t219-empty")
		defer revoke()
		card := startWriteCard219(rt, "t219-empty", "t219-empty-corr", target)
		if !waitForCard187(rt.liveCards, "t219-empty-corr", 5*time.Second) {
			t.Errorf("the L2 card never reached the native ledger")
			return
		}
		// A bare `no <id>`: the operator typed nothing in the reason box.
		if _, err := fmt.Fprintln(pw, "no t219-empty-corr"); err != nil {
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
	// (a) The model: the queue's own empty-reason sentence (internal/agent/approval/
	// queue.go:515), NOT an invented reason and not a blank where the cause should
	// be - and NOT the operator-facing marker, which is this surface's copy and not
	// something to pass off as the operator's words.
	if strings.TrimSpace(out.Text) == "" {
		t.Errorf("the model-visible refusal is empty: an empty reason must still produce a determinate cause")
	}
	if !strings.Contains(out.Text, "拒绝") {
		t.Errorf("the model-visible refusal lost its cause: %q", out.Text)
	}
	if strings.Contains(out.Text, "（未填理由）") || strings.Contains(out.Text, replyTruncationMark) {
		t.Errorf("the operator-facing marker leaked into the model's text: %q", out.Text)
	}
	// (b) The operator: the surface names the empty box instead of pretending a
	// reason arrived.
	if console := h.out.String(); !strings.Contains(console, "未填理由") {
		t.Errorf("the console did not name the empty reason:\n%s", console)
	}
	// (c) The audit: the stored detail is empty - the truth is "nothing typed",
	// not a composed sentence.
	if audit := h.err.String(); !strings.Contains(audit, `route="native/reject" detail=""`) {
		t.Errorf("the audit line did not record the empty reason as empty:\n%s", audit)
	}
	// (d) The review list: the empty state is displayed as a state.
	if console := h.out.String(); !strings.Contains(console, "（未填理由）") {
		t.Errorf("the review list printed an empty reason instead of naming the state:\n%s", console)
	}
}

// waitForConsole219 waits for one sentence to reach the run's console. The reply
// listener prints from its own goroutine, so a case that wrote a verb and then
// returned would race the process exit.
func waitForConsole219(h *replyHost, needle string, limit time.Duration) bool {
	deadline := time.Now().Add(limit)
	for time.Now().Before(deadline) {
		if strings.Contains(h.out.String(), needle) {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// newPipe219 is the console's own reply stream: the listener reads it, the case
// writes lines into it, exactly as approval_reply_201_test.go's cases do.
func newPipe219(t *testing.T) (*os.File, *os.File) {
	t.Helper()
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })
	return pr, pw
}
