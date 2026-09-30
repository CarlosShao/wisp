//go:build windows && winlive

package main

// Ticket 246 AC#2 / AC#3 / AC#4, live tier: what only Win32 can answer.
//
// Three readings this file exists to produce, in this order and no other:
//
//	(a) BORROWED  - a real L1 card is up on a real orb, and a second process on
//	    the desktop can NO LONGER receive the injected bare Esc (the key is
//	    Wisp's for the length of the window);
//	(b) VETOED    - that same injection is the press the ball turns into a veto:
//	    the card answers AnswerVeto and the gate books ANSWER-VETO channel=esc;
//	(c) RETURNED  - once nothing is being waited on, the very next injected Esc
//	    reaches the second process again (keydown_esc=1), which is the half
//	    "we released it in our own bookkeeping" cannot fake.
//
// (a) and (c) are the same measurement on the same ruler with the opposite
// expectation, and the ruler's own positive control (-steal: a window that
// registers the key itself must stop seeing the keydown) runs first, so a blind
// rig cannot report the pass. That control was ticket 245-v1's E2, and the rig is
// in-repo now (cmd/wisp/testdata/esclistener) for the reason its AC#9 gives: an
// acceptance-side rig guards no regression.
//
// Nothing here skips. A desktop that will not host the ball, or that already has
// a foreign owner of bare Esc, is a FAIL that names which of the two it is -
// ticket 245 AC#8's lesson, that "someone else holds the key" and "we did not
// take it" are different reds and must never be smoothed into a skip.
//
// Run it alone, with no other Wisp or balldebug process on the desktop:
//
//	PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -tags winlive ./cmd/wisp \
//	  -run 'TestLive246' -v
//
// What this file does NOT claim: that anyone SAW a card. This binary links no
// WebView2 host, so the evidence is the orb's state, the ledger, the audit line
// and what the desktop did with the key.

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

// askResult246 is what one blocked confirmation reports when it finally answers.
type askResult246 struct {
	ans  tools.Answer
	why  string
	took time.Duration
}

// TestLive246ConfirmingCardBorrowsEscVetoesAndReturns is AC#2 and AC#3 in one
// real timeline: card up -> key taken -> key press vetoes the card -> key handed
// back -> another process gets it again.
func TestLive246ConfirmingCardBorrowsEscVetoesAndReturns(t *testing.T) {
	rig := buildEscListener246(t)

	// The ruler's positive control, before any Wisp code is in play.
	ctrl := observeEsc246(t, rig, "-steal")
	if ctrl.keydown != 0 || ctrl.hotkey != 1 {
		t.Fatalf("the observer rig is blind: with its OWN global bare-Esc registration it still "+
			"received keydown_esc=%d (want 0) and wm_hotkey=%d (want 1); no reading below means anything",
			ctrl.keydown, ctrl.hotkey)
	}
	// Baseline on this desktop with nothing of ours running.
	base := observeEsc246(t, rig, "-watch")
	if base.keydown != 1 {
		t.Fatalf("ticket 246 RED (desktop state, not our code): an injected bare Esc reached no "+
			"foreground window before Wisp started (keydown_esc=%d, wm_hotkey=%d): somebody else owns the key "+
			"on this desktop, so the borrow cannot be measured here", base.keydown, base.hotkey)
	}

	reg := observe.NewRegistry()
	ra := newResidentApproval()
	rb := startResidentBall(reg, ra.vetoByEsc)
	if rb.b == nil {
		t.Fatalf("ticket 246 RED: this leg could not create the ball window, so there is no key to borrow: %s",
			rb.verdict)
	}
	defer rb.stop()
	if !ra.bindBallHost(rb) {
		t.Fatal("the assembly root's gate did not bind the ball it was handed")
	}
	defer ra.detachBall()

	if !ra.gate.Channels().Loaded(channelEsc246()) {
		t.Fatal("Esc is not loaded with a real ball attached")
	}
	// Idle roster first: three keys, cancel standby (ticket 245's shape, now on
	// the leg that can actually reach the fourth).
	requireIdleCancelSlot246(t, rb)

	audit := captureAudit246(t)

	done := make(chan askResult246, 1)
	start := time.Now()
	go func() {
		a, w := ra.AskOnTaskRoot(residentCard{
			TaskID: "host:live246-esc",
			Tool:   residentCardTool,
			Level:  risk.L1,
			Reason: "真机验收：一张 L1 确认卡片挂起时，取消键归 Wisp",
		})
		done <- askResult246{ans: a, why: w, took: time.Since(start)}
	}()

	if !waitFor246(3*time.Second, func() bool {
		_, awaiting := ra.cards.AwaitingHuman()
		return awaiting
	}) {
		t.Fatal("no L1 card ever reached AwaitingHuman: the injected gate is not showing cards")
	}

	// AC#2's state位, from the seam that produces D43 names.
	st, ok := ra.cards.WaitingState()
	if !ok || st != statemachine.StateConfirming {
		t.Fatalf("WaitingState() = %q/%v with an L1 window open, want Confirming", st, ok)
	}
	if !rb.b.EscTakenOver() {
		t.Fatal("a card is up and the cancel key was not borrowed: D43's Confirming row has no veto route")
	}
	if rep := rb.b.HotkeyReport(); len(rep.Live()) != 4 {
		t.Fatalf("live hot key set during Confirming = %d (%+v), want 4", len(rep.Live()), rep.Bindings())
	}

	// (a) + (b): the observer's window takes the foreground and injects ONE Esc on
	// command, while the card is still counting down.
	borrowed := observeEsc246(t, rig, "-watch")
	if borrowed.keydown != 0 {
		t.Fatalf("(a) FAIL: the second process still received the injected Esc (keydown_esc=%d) "+
			"while a card was waiting - the key was not really taken", borrowed.keydown)
	}

	var res askResult246
	select {
	case res = <-done:
	case <-time.After(6 * time.Second):
		t.Fatal("the veto did not reach the blocked confirmation within 6s")
	}
	if res.ans != tools.AnswerVeto {
		t.Fatalf("(b) FAIL: the card answered %s (%q) after the injected Esc, want veto; window took %s",
			res.ans, res.why, res.took)
	}
	if !strings.Contains(res.why, "Esc") {
		t.Fatalf("(b) the veto reason does not name the channel: %q", res.why)
	}
	if line := audit.find("ANSWER-VETO"); line == "" {
		t.Fatalf("(b) FAIL: no ANSWER-VETO audit line after a veto that really cancelled a call; lines were:\n%s",
			audit.dump())
	} else if !strings.Contains(line, "channel=esc") {
		t.Fatalf("(b) the veto line does not name the esc channel: %q", line)
	}

	// (c): the borrow must be over, per our own bookkeeping, per Win32's roster,
	// and per the second process that now gets its key back.
	if rb.b.EscTakenOver() {
		t.Fatal("(c) FAIL: Esc still marked borrowed after the card closed")
	}
	requireIdleCancelSlot246(t, rb)
	returned := observeEsc246(t, rig, "-watch")
	if returned.keydown != 1 {
		t.Fatalf("(c) FAIL: after the card closed, another process received keydown_esc=%d, want 1 - "+
			"Wisp is still swallowing Esc (wm_hotkey=%d on the observer means the observer itself holds it)",
			returned.keydown, returned.hotkey)
	}
	if _, awaiting := ra.cards.AwaitingHuman(); awaiting {
		t.Fatal("a settled card is still waiting on a human")
	}
	t.Logf("AC#3 READING: borrow/veto/return in %s; veto=%s; observer keydown %d->%d->%d (steal control %d)",
		res.took, res.ans, base.keydown, borrowed.keydown, returned.keydown, ctrl.keydown)
}

// TestLive246ExitRefusesAHangingL2Card is AC#4's one shot on a real window: a
// card that is genuinely waiting (the C18 route, 300s, so no stop watch is
// involved) when the D38(e) sequence runs. The step must be executed, not
// recorded skipped, and the card must leave as a refusal booked in the same
// ANSWER-REJECT family a native 「拒绝」 writes.
func TestLive246ExitRefusesAHangingL2Card(t *testing.T) {
	reg := observe.NewRegistry()
	rt, err := proc.Boot(buildinfo.EnvTest, proc.WithRegistry(reg))
	if err != nil {
		t.Fatalf("proc.Boot(test): %v", err)
	}
	ra := newResidentApproval()
	rb := startResidentBall(reg, ra.vetoByEsc)
	if rb.b == nil {
		t.Fatalf("no ball window, so no card can hang: %s", rb.verdict)
	}
	defer rb.stop()
	if !ra.bindBallHost(rb) {
		t.Fatal("gate did not bind the ball")
	}
	defer ra.detachBall()

	if err := rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots); err != nil {
		t.Fatalf("RegisterShutdownHook(step 3): %v", err)
	}
	audit := captureAudit246(t)

	done := make(chan askResult246, 1)
	go func() {
		start := time.Now()
		a, w := ra.AskOnTaskRoot(residentCard{
			TaskID: "host:live246-exit",
			Tool:   residentCardTool,
			Level:  risk.L2,
			Reason: "真机验收：卡片挂着时收到退出信号",
		})
		done <- askResult246{ans: a, why: w, took: time.Since(start)}
	}()
	if !waitFor246(3*time.Second, func() bool {
		card, awaiting := ra.cards.AwaitingHuman()
		return awaiting && card.Level == "L2"
	}) {
		t.Fatal("no L2 card ever hung")
	}
	st, ok := ra.cards.WaitingState()
	if !ok || st != statemachine.StateAwaitingApproval {
		t.Fatalf("WaitingState() = %q/%v with an L2 card pending, want AwaitingApproval", st, ok)
	}
	// An L2 card must NOT borrow the desktop key: that is ticket 245's whole
	// point, and the exit path cannot quietly become a second reason to hold it.
	if rb.b.EscTakenOver() {
		t.Fatal("an L2 card borrowed the cancel key: a 300s deadline is not a 2-3s window")
	}

	records := rt.Shutdown(false)
	if len(records) != 10 {
		t.Fatalf("audit trail = %d records, want the frozen 10", len(records))
	}
	step3 := records[proc.StepCancelTasks-1]
	if step3.Skipped {
		t.Fatalf("AC#4 FAIL: step 3 recorded SKIPPED while a card was hanging: %+v", step3)
	}
	if step3.Err != nil {
		t.Fatalf("AC#4: step 3 reported an error: %+v", step3)
	}

	var res askResult246
	select {
	case res = <-done:
	case <-time.After(8 * time.Second):
		t.Fatal("AC#4 FAIL: the hanging card was never resolved by the exit sequence")
	}
	if res.ans != tools.AnswerReject {
		t.Fatalf("AC#4 FAIL: the card answered %s (%q), want reject: leaving must never be an allow", res.ans, res.why)
	}
	if line := audit.find("ANSWER-REJECT"); line == "" {
		t.Fatalf("AC#4 FAIL: no ANSWER-REJECT audit line for the refused card:\n%s", audit.dump())
	}
	if _, awaiting := ra.cards.AwaitingHuman(); awaiting {
		t.Fatal("a card survived the exit sequence still waiting on a human")
	}
	if rb.b.EscTakenOver() {
		t.Fatal("the cancel key is borrowed after the sequence ran")
	}
	requireIdleCancelSlot246(t, rb)
	// AC#4's other half: the hotkey report must not grow an "unregistered" complaint
	// because a card came and went - standby is not a problem line, and a false
	// problem is the same species of lie as a false skip.
	if probs := rb.b.HotkeyReport().Problems(); len(probs) != 0 {
		t.Fatalf("the exit path left %d problem lines behind: %v", len(probs), probs)
	}
	t.Logf("AC#4 READING: step3 executed in %s, card answered %s, ANSWER-REJECT booked; no-skipped=%v",
		step3.Elapsed, res.ans, !step3.Skipped)
}

// TestLive246ExitAbandonsAHangingL1Window covers the other route: an L1 window
// has no reject verb at all (SPEC-06 §2 B1 - a countdown can be opposed, not
// answered), so what the exit sequence owes it is a booked abandonment plus a
// root cancel that the window answers as a refusal.
func TestLive246ExitAbandonsAHangingL1Window(t *testing.T) {
	reg := observe.NewRegistry()
	rt, err := proc.Boot(buildinfo.EnvTest, proc.WithRegistry(reg))
	if err != nil {
		t.Fatalf("proc.Boot(test): %v", err)
	}
	ra := newResidentApproval()
	rb := startResidentBall(reg, ra.vetoByEsc)
	if rb.b == nil {
		t.Fatalf("no ball window: %s", rb.verdict)
	}
	defer rb.stop()
	ra.bindBallHost(rb)
	defer ra.detachBall()
	if err := rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots); err != nil {
		t.Fatalf("RegisterShutdownHook(step 3): %v", err)
	}
	audit := captureAudit246(t)

	done := make(chan askResult246, 1)
	go func() {
		start := time.Now()
		a, w := ra.AskOnTaskRoot(residentCard{
			TaskID: "host:live246-exit-l1", Tool: residentCardTool, Level: risk.L1,
			Reason: "真机验收：L1 窗口挂着时收到退出信号",
		})
		done <- askResult246{ans: a, why: w, took: time.Since(start)}
	}()
	if !waitFor246(2*time.Second, func() bool { return rb.b.EscTakenOver() }) {
		t.Fatal("the L1 window never borrowed the key, so there is nothing to abandon")
	}
	records := rt.Shutdown(false)
	if rec := records[proc.StepCancelTasks-1]; rec.Skipped {
		t.Fatalf("AC#4 FAIL: step 3 skipped with an L1 window open: %+v", rec)
	}
	var res askResult246
	select {
	case res = <-done:
	case <-time.After(8 * time.Second):
		t.Fatal("the L1 window never answered")
	}
	if res.ans != tools.AnswerReject {
		t.Fatalf("an abandoned L1 window answered %s (%q), want reject", res.ans, res.why)
	}
	if line := audit.find("RESIDENT-WINDOW-ABANDONED"); line == "" {
		t.Fatalf("the abandoned L1 window left no audit line (the gate writes none for this branch):\n%s",
			audit.dump())
	}
	if rb.b.EscTakenOver() {
		t.Fatal("the key is still borrowed after the window was abandoned")
	}
	requireIdleCancelSlot246(t, rb)
}

// ------------------------------------------------------------------- helpers

// channelEsc246 is the one veto channel ticket 246 loads. Named through the
// approval package's own constant, never a string literal, so the case below
// cannot pass on a typo.
func channelEsc246() approval.Channel { return approval.ChannelEsc }

func requireIdleCancelSlot246(t *testing.T, rb *residentBall) {
	t.Helper()
	rep := rb.b.HotkeyReport()
	if len(rep.Live()) != 3 {
		t.Fatalf("idle live hot key set = %d (%+v), want 3 (summon/mute/panel): the cancel key was not "+
			"handed back, or never handed at all", len(rep.Live()), rep.Bindings())
	}
	if rb.b.EscTakenOver() {
		t.Fatal("EscTakenOver still true while the roster says idle")
	}
}

// captureAudit246 points the process default logger at an in-memory buffer for
// the duration of one case and restores it afterwards. It is the only way to read
// the gate's own lines from an in-process assembly.
func captureAudit246(t *testing.T) *auditCapture246 {
	t.Helper()
	c := &auditCapture246{}
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(c, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return c
}

type auditCapture246 struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (c *auditCapture246) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buf.Write(p)
}

func (c *auditCapture246) dump() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buf.String()
}

func (c *auditCapture246) find(needle string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, line := range strings.Split(c.buf.String(), "\n") {
		if strings.Contains(line, needle) {
			return line
		}
	}
	return ""
}

// buildEscListener246 compiles the in-repo rig. It lives under testdata/ on
// purpose: the go tool's package walk and this repository's d22scan walk both
// skip that directory, so the case that needs it names it explicitly.
func buildEscListener246(t *testing.T) string {
	t.Helper()
	if runtime.GOARCH != "amd64" {
		t.Fatalf("GOARCH = %s, want amd64", runtime.GOARCH)
	}
	goExe := "go"
	if p, err := exec.LookPath("go"); err == nil {
		goExe = p
	} else if rooted := filepath.Join(runtime.GOROOT(), "bin", "go.exe"); fileExists(rooted) {
		goExe = rooted
	}
	exe := filepath.Join(t.TempDir(), "esclistener.exe")
	cmd := exec.Command(goExe, "build", "-o", exe, "./cmd/wisp/testdata/esclistener")
	cmd.Dir = filepath.Join("..", "..")
	var log bytes.Buffer
	cmd.Stdout, cmd.Stderr = &log, &log
	if err := cmd.Run(); err != nil {
		t.Fatalf("go build ./cmd/wisp/testdata/esclistener: %v\n%s", err, log.String())
	}
	if !fileExists(exe) {
		t.Fatalf("build produced no %s", exe)
	}
	return exe
}

// observerReading246 is one run of the second process.
type observerReading246 struct {
	keydown int
	syskey  int
	hotkey  int
	other   int
	fgMine  bool
	steal   bool
	raw     string
}

// observeEsc246 starts the rig, waits for its window to say READY (a window that
// never owned the foreground is reported, not papered over), triggers the single
// Esc injection, and reads back the READING line.
func observeEsc246(t *testing.T, exe string, mode ...string) observerReading246 {
	t.Helper()
	cmd := exec.Command(exe, mode...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start observer: %v", err)
	}
	t.Cleanup(func() {
		_ = stdin.Close()
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	})

	lines := make(chan string, 32)
	go func() {
		sc := bufio.NewScanner(stdout)
		sc.Buffer(make([]byte, 0, 4096), 1<<16)
		for sc.Scan() {
			lines <- sc.Text()
		}
		close(lines)
	}()

	var got observerReading246
	ready := false
	timeout := time.After(20 * time.Second)
	for {
		select {
		case line, ok := <-lines:
			if !ok {
				t.Fatalf("observer exited before printing READING (mode %v, stderr: %s)", mode, stderr.String())
			}
			got.raw += line + "\n"
			switch {
			case strings.HasPrefix(line, "READY"):
				ready = true
				got.fgMine = strings.Contains(line, "foreground_is_mine=true")
				got.steal = strings.Contains(line, "steal=true")
				if _, err := fmt.Fprintln(stdin, "go"); err != nil {
					t.Fatalf("trigger write: %v", err)
				}
			case strings.HasPrefix(line, "READING"):
				if !ready {
					t.Fatalf("observer printed a reading without ever being ready: %s", got.raw)
				}
				if !got.fgMine {
					t.Fatalf("the observer window never owned the foreground, so its delivery count proves "+
						"nothing about Wisp: %s", got.raw)
				}
				got.keydown = mustField246(t, line, "keydown_esc=")
				got.syskey = mustField246(t, line, "syskeydown_esc=")
				got.hotkey = mustField246(t, line, "wm_hotkey=")
				got.other = mustField246(t, line, "other_keys=")
				t.Logf("observer %v: %s", mode, strings.TrimSpace(got.raw))
				return got
			}
		case <-timeout:
			t.Fatalf("observer produced no READING within 20s (mode %v, so far: %s, stderr: %s)",
				mode, got.raw, stderr.String())
		}
	}
}

func mustField246(t *testing.T, line, key string) int {
	t.Helper()
	i := strings.Index(line, key)
	if i < 0 {
		t.Fatalf("field %q missing from %q", key, line)
	}
	rest := line[i+len(key):]
	j := strings.IndexAny(rest, " \t\n")
	if j >= 0 {
		rest = rest[:j]
	}
	n, err := strconv.Atoi(rest)
	if err != nil {
		t.Fatalf("field %q = %q is not a number in %q", key, rest, line)
	}
	return n
}

func waitFor246(budget time.Duration, cond func() bool) bool {
	deadline := time.Now().Add(budget)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}
