//go:build windows

package main

// Ticket 246 AC#1 / AC#2 / AC#4, default tier: the shapes that do NOT need a
// desktop, and the two postures that must hold precisely when one is missing.
//
// The real-machine half of this family - a card that is really on the orb, the
// cancel key really borrowed, and a second process really receiving it again - is
// resident_approval_live_246_windows_test.go behind -tags winlive, for the same
// reason ticket 245's hot-key judgements are: only Win32 can answer them.
//
// What this file refuses to assert: any claim about a card being VISIBLE. The
// panel this binary could show it on does not exist in this tree
// (approval_always.go:165, "this binary links no WebView2 host"), so the readings
// here are log records, ledger state, and what the gate answered - never a
// screenshot.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

// TestAC246CancelGestureUsesTheInjectedExecutor is AC#1's plumbing, measured on
// the ball host without a window: the cancel gesture goes to the function the
// assembly root handed in, exactly once, and what comes back is what this process
// says. The nil case is the pre-ticket-246 shape and has to be REPORTED, not
// silently accepted - a cancel key that fires into nothing is the defect class
// this ticket exists to close.
func TestAC246CancelGestureUsesTheInjectedExecutor(t *testing.T) {
	calls := 0
	executor := func() string {
		calls++
		return "injected executor ran"
	}
	recordCancelHotkey(executor)
	if calls != 1 {
		t.Fatalf("the injected executor ran %d times for one cancel gesture, want 1", calls)
	}

	// nil: no gate was injected. Must not panic, must not claim anything happened.
	recordCancelHotkey(nil)
	if calls != 1 {
		t.Fatalf("the nil case reached the real executor (calls = %d)", calls)
	}
}

// TestAC246EscChannelStaysUnloadedWithoutABallWindow is 246-a1's D-5 guard: a
// gate that advertises a cancel key in a process with no key to borrow is the
// loaded-but-no-ball inconsistency the census named. bindBallHost is the only
// thing that may flip that flag, and only on a real window.
func TestAC246EscChannelStaysUnloadedWithoutABallWindow(t *testing.T) {
	ra := newResidentApproval()
	if ra.gate.Channels().Loaded(approval.ChannelEsc) {
		t.Fatal("a freshly built resident gate claims the Esc channel before any ball exists")
	}
	// The exact struct startResidentBall returns when ball.New failed: verdict set,
	// no window behind it.
	bound := ra.bindBallHost(&residentBall{verdict: "this process has NO floating ball window: injected failure"})
	if bound {
		t.Fatal("bindBallHost reported a ball this residentBall does not have")
	}
	if ra.gate.Channels().Loaded(approval.ChannelEsc) {
		t.Fatal("Esc got loaded with no ball window")
	}
	// And the registry, not a branch in cmd/wisp, is what refuses the veto.
	err := ra.gate.Veto(approval.Veto{CorrelationID: "corr-not-real", Channel: approval.ChannelEsc})
	if !errors.Is(err, approval.ErrChannelUnavailable) {
		t.Fatalf("Veto on an unloaded channel err = %v, want ErrChannelUnavailable", err)
	}
	for _, ch := range []approval.Channel{approval.ChannelBall, approval.ChannelKWS, approval.ChannelPanel} {
		if ra.gate.Channels().Loaded(ch) {
			t.Fatalf("channel %s loaded in the resident leg: ticket 246 AC#6 says only Esc is wired", ch)
		}
	}
}

// TestAC246CardWithNoWindowFailsClosedThroughTheRealGate is AC#2's honest
// negative: no fake UI is standing in here. The gate is the production one, the
// task really is admitted, and the presentation surface is the one the assembly
// root injected - which says "I have nowhere to show this", so the card is
// REFUSED. A confirmation nobody can be shown must never run the call.
func TestAC246CardWithNoWindowFailsClosedThroughTheRealGate(t *testing.T) {
	ra := newResidentApproval()
	ra.bindBallHost(&residentBall{})

	ans, why := ra.askConfirmation(context.Background(), residentCard{
		TaskID: "host:246-no-window",
		Tool:   residentCardTool,
		Level:  risk.L1,
		Reason: "测试用宿主确认：本进程没有窗口",
	})
	if ans != tools.AnswerReject {
		t.Fatalf("answer = %s (%s), want reject: a card with nowhere to appear cannot be an allow", ans, why)
	}
	if !strings.Contains(why, "确认界面不可达") {
		t.Fatalf("reason = %q, want the gate's own fail-closed sentence", why)
	}
	if _, awaiting := ra.cards.AwaitingHuman(); awaiting {
		t.Fatal("a card that was never displayed is still waiting on a human")
	}
	if n := ra.liveConfirmations(); n != 0 {
		t.Fatalf("liveAsks() = %d after the ask returned, want 0", n)
	}
	if n := ra.ui.displayedCards(); n != 0 {
		t.Fatalf("the UI displayed %d cards with no ball attached", n)
	}
	// The orb's state name for this level, and the fact no state was claimed here.
	if s := stateForCardLevel("L1"); s != statemachine.StateConfirming {
		t.Fatalf("L1 maps to %q, want Confirming (D43 row 17)", s)
	}
	if s := stateForCardLevel("L2"); s != statemachine.StateAwaitingApproval {
		t.Fatalf("L2 maps to %q, want AwaitingApproval (D43 row 17)", s)
	}
}

// TestAC246CancelStepHookRunsOnTheRealShutdownSequence is AC#4's headless half:
// the step is registered through the seam, executed inside the frozen sequence,
// and the gate closes so an exiting process cannot open a NEW question on the way
// out. The one-shot with a card actually hanging needs a window, and lives in the
// winlive file.
func TestAC246CancelStepHookRunsOnTheRealShutdownSequence(t *testing.T) {
	rt, err := proc.Boot(buildinfo.EnvTest, proc.WithRegistry(observe.NewRegistry()))
	if err != nil {
		t.Fatalf("proc.Boot(test): %v", err)
	}
	ra := newResidentApproval()
	if got := rt.RegisteredShutdownSteps(); len(got) != 0 {
		t.Fatalf("a fresh runtime already owns steps %v", got)
	}
	if err := rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots); err != nil {
		t.Fatalf("RegisterShutdownHook(step 3): %v", err)
	}

	records := rt.Shutdown(false)
	if len(records) != 10 {
		t.Fatalf("audit trail = %d records, want the frozen 10", len(records))
	}
	step3 := records[proc.StepCancelTasks-1]
	if step3.Skipped {
		t.Fatalf("step 3 recorded skipped with a hook registered: %+v", step3)
	}
	if step3.Err != nil {
		t.Fatalf("step 3 reported an error with nothing hanging: %+v", step3)
	}
	if step3.Name != "cancel-task-roots" {
		t.Fatalf("step 3 name = %q, want cancel-task-roots (the frozen label)", step3.Name)
	}

	// After the sequence has walked, this gate must refuse to ask anything.
	ans, why := ra.askConfirmation(context.Background(), residentCard{
		TaskID: "host:246-after-shutdown", Tool: residentCardTool, Level: risk.L1,
		Reason: "退出序列之后还想签一张卡",
	})
	if ans != tools.AnswerReject || !strings.Contains(why, "退出序列") {
		t.Fatalf("answer after shutdown = %s (%q), want reject naming the exit sequence", ans, why)
	}
}

// TestAC246ShippedResidentProcessOwnsItsCancelStep is AC#1 and AC#4 on the
// command line the owner actually double clicks: no arguments, in a real process,
// read out of its own console and its own on-disk log.
//
// It exists because the in-process cases above could be satisfied by a gate that
// nothing assembles. Two facts have to survive the boundary of the shipped
// binary: the boot report names the step this process registered as its own, and
// the D38(e) trail the process writes on the way out contains the sentence only
// the hook can write. Deleting the registration in resident_windows.go takes both
// readings away, and this case goes red with them.
//
// No branch of the desktop state is skipped here: a leg with no ball prints the
// un-assembled gate sentence, and the cancel step is still registered, because
// what step 3 cancels is the task root, not a window.
func TestAC246ShippedResidentProcessOwnsItsCancelStep(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	leg := bootResidentLeg(t, exe, dataDir)
	sinkDir := logSinkDir(dataDir)

	said := pollUntil127(400, func() bool {
		out := leg.stdout.String()
		posture := strings.Contains(out, ballUpClaim) || strings.Contains(out, ballAbsentClaim)
		return posture && strings.Contains(out, cancelStepRosterClaim)
	})
	out := leg.stdout.String()
	if !said {
		leg.stop()
		t.Fatalf("AC#1/#4 RED: the shipped resident process never printed both its ball posture and its "+
			"registered D38(e) steps.\n%s\n"+
			"Expected one of %q / %q plus %q: the pair is what this process knows about its own window and "+
			"the step it registered a hook for. A leg that stopped wiring says neither.",
			leg.console(), ballUpClaim, ballAbsentClaim, cancelStepRosterClaim)
	}
	if !strings.Contains(out, cancelStepRosterClaim) {
		leg.stop()
		t.Fatalf("AC#4 RED: the boot report does not name step 3 as owned by this process (%q):\n%s",
			cancelStepRosterClaim, leg.console())
	}
	// The coherence check, which is what makes the injection itself observable.
	// Which half is required depends on the desktop, and that is the point:
	// bindBallHost loads the Esc channel only when Win32 gave this process a window
	// AND the assembly root handed it a cancel executor, so a leg that lost the
	// injection has to say 审批门未装配 while still owning step 3 - and a leg that
	// advertises the channel with no executor behind the key fails here.
	switch {
	case strings.Contains(out, ballUpClaim):
		if !strings.Contains(out, gateAssembledClaim) {
			t.Errorf("AC#1 RED: the ball window is up but the gate is not assembled: %s"+
				" A window whose cancel key fires into nothing must not advertise the channel; if this is the"+
				" residue of deleting the injection in resident_windows.go, that is exactly the lie.",
				leg.console())
		}
	case strings.Contains(out, ballAbsentClaim):
		if !strings.Contains(out, gateAbsentClaim) {
			t.Errorf("AC#1 RED: no ball window, yet the gate claims to be assembled: %s", leg.console())
		}
	default:
		t.Fatalf("AC#1 RED: the process named neither ball posture: %s", leg.console())
	}

	if err := leg.breakToLoop(); err != nil {
		leg.stop()
		t.Fatalf("GenerateConsoleCtrlEvent on the leg's own group: %v\n%s", leg.pid(), err)
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("the leg took the exit request and never left through its own shutdown path.\n%s", leg.console())
	}
	if leg.waitErr != nil {
		t.Errorf("AC#4 RED: the shipped process exited with %v; a cancel step that errors on an empty set of "+
			"tasks is a teardown defect, not a cosmetic one.\n%s", leg.waitErr, leg.console())
	}

	recs := readResidentSink(t, sinkDir)
	hookLine := indexOfMsgContaining246(recs, cancelStepBookedMsg)
	if hookLine < 0 {
		t.Fatalf("AC#4 RED: the on-disk log of the shipped process holds no %q record (of %d records: %v).\n"+
			"That sentence is written by the registered hook, inside the frozen sequence: no record means the "+
			"step ran outside it or not at all, which is the fork ticket 246 closed.",
			cancelStepBookedMsg, len(recs), msgsOf127(recs))
	}
	// And the trail still says ten steps with nothing failing, so the registration
	// did not cost the sequence anything. "shutdown step skipped (module not present)"
	// is the honest record for the six steps this ticket does not own and is NOT a
	// failure; only a step that ran and broke, or blew its deadline, is.
	for _, bad := range []string{"shutdown step failed", "abandoned wait"} {
		if idx := indexOfMsgContaining246(recs, bad); idx >= 0 {
			t.Errorf("AC#4 RED: the shipped process booked %q: %q", bad, msgsOf127(recs)[idx])
		}
	}
	if !strings.Contains(leg.stdout.String(), "10 steps, 0 failed") {
		t.Errorf("AC#4 RED: the console does not report a clean 10-step exit: %q",
			tailContaining246(leg.stdout.String(), "shutdown order"))
	}
}

const (
	// gateAssembledClaim / gateAbsentClaim are the two forms the boot report can
	// take, copied as literals on purpose: rewording resident_approval_windows.go
	// has to be a deliberate act here too.
	gateAssembledClaim = "审批门已装配进本进程"
	gateAbsentClaim    = "审批门未装配"
	// cancelStepRosterClaim is the half that names the registered step.
	cancelStepRosterClaim = "3:cancel-task-roots"
	// cancelStepBookedMsg is the hook's own completion sentence in the log file.
	cancelStepBookedMsg = "resident-approval: 退出第 3 步完成"
)

func indexOfMsgContaining246(recs []sinkInstallRecord, needle string) int {
	for i, r := range recs {
		if strings.Contains(r.Msg, needle) {
			return i
		}
	}
	return -1
}

func tailContaining246(all, needle string) string {
	lines := strings.Split(all, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], needle) {
			return lines[i]
		}
	}
	return "(no line containing " + needle + ")"
}

// TestAC246ChannelNeedsBothWindowAndExecutor is the coherence guard on its own
// two legs. The assembly may advertise the Esc cancel channel only when BOTH
// facts hold: Win32 gave this process a ball window to borrow the key for, and
// the assembly root handed that window a cancel executor. Each half is refused
// here from the side this package can build:
//
//   - no executor: the host is started exactly the way resident_windows.go starts
//     it, with a real startResidentBall call, so on a desktop rb.b is a live
//     window and the refusal can only come from the missing executor;
//   - no window: the failed-creation shape (246-a1's D-5).
//
// On a machine with no desktop both legs collapse to the second reading, which is
// still a real assertion - it is the loaded flag that must stay false either way.
func TestAC246ChannelNeedsBothWindowAndExecutor(t *testing.T) {
	// leg 1: window present or not, but definitely no injected executor.
	ra := newResidentApproval()
	rb := startResidentBall(observe.NewRegistry(), nil, nil, nil)
	if rb.cancelHosted {
		t.Fatal("cancelHosted true for a host started with a nil executor")
	}
	if bound := ra.bindBallHost(rb); bound {
		t.Fatalf("bindBallHost loaded the cancel channel with no executor behind it (ball up = %v)", rb.b != nil)
	}
	if ra.gate.Channels().Loaded(approval.ChannelEsc) {
		t.Fatalf("Esc advertised with no cancel executor injected; the ball window behind this host: %v", rb.b != nil)
	}
	if line := ra.residentStatusLine(); !strings.Contains(line, "审批门未装配") {
		t.Fatalf("statusLine = %q while nothing was injected", line)
	}
	rb.stop()

	// leg 2: the other half, taken from the same function the production host uses.
	ra2 := newResidentApproval()
	if bound := ra2.bindBallHost(&residentBall{cancelHosted: true}); bound {
		t.Fatal("bindBallHost claimed a window that the failed-creation shape does not have")
	}
	if ra2.gate.Channels().Loaded(approval.ChannelEsc) {
		t.Fatal("Esc advertised with no ball window")
	}
}

// TestAC246VetoSentenceWithNoCard pins that the injected cancel path cannot
// invent a decision: with nothing on screen, a press says so and touches nothing.
func TestAC246VetoSentenceWithNoCard(t *testing.T) {
	ra := newResidentApproval()
	said := ra.vetoByEsc()
	if !strings.Contains(said, "没有可否决的确认项") {
		t.Fatalf("cancel press with no card said %q, want the nothing-was-waiting sentence", said)
	}
	if _, awaiting := ra.cards.AwaitingHuman(); awaiting {
		t.Fatal("AwaitingHuman true with no card ever raised")
	}
	if st, ok := ra.cards.WaitingState(); ok {
		t.Fatalf("WaitingState() = %q with an empty ledger", st)
	}
}

// TestAC246StatusLineSaysWhatTheLegDoesNot is the sentence-decides-nothing guard:
// the boot report may not read as a running task pipeline, and with no window it
// may not read as a leg that can show a card either.
func TestAC246StatusLineSaysWhatTheLegDoesNot(t *testing.T) {
	ra := newResidentApproval()
	line := ra.residentStatusLine()
	if !strings.Contains(line, "审批门未装配") {
		t.Fatalf("statusLine before any ball = %q, want the un-assembled sentence", line)
	}
	if bound := ra.bindBallHost(&residentBall{}); bound {
		t.Fatal("bindBallHost claimed a window that is not there")
	}
	if line := ra.residentStatusLine(); !strings.Contains(line, "审批门未装配") {
		t.Fatalf("statusLine after binding an empty host = %q, want the un-assembled sentence", line)
	}
	// The other half: the sentence this leg prints must never be a promise that a
	// card is on screen somewhere. That claim belongs to a panel host, and this
	// process links none (approval_always.go:165).
	for _, forbidden := range []string{"看得见", "面板已就绪", "已显示卡片"} {
		if strings.Contains(ra.residentStatusLine(), forbidden) {
			t.Fatalf("statusLine says %q, which this leg cannot mean: %q", forbidden, line)
		}
	}
}
