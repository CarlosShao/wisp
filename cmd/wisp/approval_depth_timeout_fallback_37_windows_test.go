//go:build windows

package main

// Ticket 37 (.scratch/wisp/issues/37-approval-ui-l2.md AC#1/AC#2/AC#3), the
// three cells whose production paths live in cmd/wisp. The panel-page half of
// this ticket already landed (readface legs, render:approval37); this file is
// the queue/badge/native-card half:
//
//	cell 1  two SEQUENTIAL L2 requests from ONE task (the pre-47 shape, driven
//	        with the ticket-43/48 raise-and-wait-for-depth queue test mode):
//	        the orb badge input reads 1 then 2 (approvalDepthBadge over the
//	        card the gate stamped), the panel face shows the HEAD only, and
//	        answers route by correlationId - a reply that names the SECOND
//	        card settles the second card's ask, never the first arrival's.
//	        Plus the ticket-87 R-3 acceptance case this ticket's wiring note
//	        names: two cards coexist, cancelling one never hits the other.
//	cell 2  the C18 ladder on the injected Clock (approval.Clock seam, no
//	        wall-clock comparison is ever a judging criterion, ledger A1886):
//	        warning event at 270s with the 30s remainder, auto-REJECT at 300s,
//	        Gate.Replay completes the original request under a fresh corr and
//	        a fresh grant, and the badge inputs go to zero (queue depth 0,
//	        ledger empty, ballCardUI.Update's real settle path).
//	cell 3  the panel runtime MASKED (ticket 33 AC#5's own fixture:
//	        withMaskedRuntime -> named refusal -> panel.unavailable ->
//	        Unavailable() ON): an L2 card is still approvable through the
//	        NATIVE card path (Replies.Allow, the ball/tray verb), while the
//	        panel route keeps refusing allows (ErrPanelAllow) and burns the
//	        offered nonce - C27's "L2 capability never disappears" with the
//	        F2 no-panel-allow invariant intact.
//
// Everything here is headless: no Win32 window is created, nothing touches the
// owner's ball process. The Win32 posting itself (SetBadge) stays pinned by
// TestTicket21S2BadgeWiringCensus; the tests below pin the VALUES that
// posting carries.

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// the fixture: the ticket-43 approval assembly, audio/voice legs removed
// ---------------------------------------------------------------------------

// recordingUI37 is the injected card face: it records every prompt the gate
// displays (so the DEPTH stamp is a reading, not an inference) and forwards
// into the same reply ledger the resident boots with.
type recordingUI37 struct {
	ra *residentApproval

	mu      sync.Mutex
	prompts []approval.Prompt
	events  []approval.Event
}

func (u *recordingUI37) Prompt(_ context.Context, p approval.Prompt) error {
	u.mu.Lock()
	u.prompts = append(u.prompts, p)
	u.mu.Unlock()
	u.ra.cards.Record(p)
	return nil
}

func (u *recordingUI37) Update(_ context.Context, e approval.Event) error {
	u.mu.Lock()
	u.events = append(u.events, e)
	u.mu.Unlock()
	if e.Kind == approval.EventDismissed || e.Kind == approval.EventStarted {
		u.ra.cards.Forget(e.CorrelationID)
	}
	return nil
}

func (u *recordingUI37) all() []approval.Prompt {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]approval.Prompt(nil), u.prompts...)
}

func (u *recordingUI37) promptOf(corr string) (approval.Prompt, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, p := range u.prompts {
		if p.CorrelationID == corr {
			return p, true
		}
	}
	return approval.Prompt{}, false
}

func (u *recordingUI37) ofKind(k approval.EventKind) []approval.Event {
	u.mu.Lock()
	defer u.mu.Unlock()
	var out []approval.Event
	for _, e := range u.events {
		if e.Kind == k {
			out = append(out, e)
		}
	}
	return out
}

// fixture37 is the minimal residentApproval assembly: the REAL gate, the REAL
// reply ledger with its HostBinding (the verb table the native card holds),
// and the recording UI. The audio and voice legs of newSuspendFixture43 are
// not needed by any cell below and are deliberately absent.
type fixture37 struct {
	ra   *residentApproval
	ui   *recordingUI37
	logs *logSink37
}

// logSink37 captures the gate's own audit face (approval.Options.Logf) - the
// channel the PANEL-ALLOW-REJECTED and FORGED-OR-STALE lines travel on, which
// is NOT the slog stream the panel emits on.
type logSink37 struct {
	mu    sync.Mutex
	lines []string
}

func (s *logSink37) add(format string, args ...any) {
	s.mu.Lock()
	s.lines = append(s.lines, fmt.Sprintf(format, args...))
	s.mu.Unlock()
}

func (s *logSink37) has(sub string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, l := range s.lines {
		if strings.Contains(l, sub) {
			return true
		}
	}
	return false
}

func newFixture37(t *testing.T, clk approval.Clock) *fixture37 {
	t.Helper()
	root, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	f := &fixture37{ra: &residentApproval{root: root, cancel: cancel}, logs: &logSink37{}}
	f.ra.ui = &ballCardUI{ra: f.ra}
	f.ra.grants = &residentGrantHolder{}
	f.ui = &recordingUI37{ra: f.ra}
	o := approval.Options{
		UI:       f.ui,
		Channels: approval.NewChannels(),
		Logf:     f.logs.add,
	}
	if clk != nil {
		o.Clock = clk
	}
	f.ra.gate = approval.New(o)
	f.ra.cards = approval.NewReplies()
	f.ra.cards.Attach(approval.HostBinding{Gate: f.ra.gate})
	return f
}

// answer37 is one parked ask's outcome, the pair PendingApproval returns.
type answer37 struct {
	a   tools.Answer
	why string
}

// ask37 parks ONE L2 ask off the test goroutine (observe.Default.Spawn under
// the D38 approval-waiter name, recover owned by the test) and returns the
// channel its answer lands on. The task is admitted here exactly the way
// execute() does - the D47 admit guard is a production precondition.
func ask37(t *testing.T, f *fixture37, taskID, corr string) chan answer37 {
	t.Helper()
	f.ra.gate.AdmitTextTask(taskID)
	out := make(chan answer37, 1)
	h := observe.Default.Spawn("approval-waiter", "test", nil, func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("approval-waiter panicked (D37b): %v", r)
			}
		}()
		a, why := f.ra.gate.PendingApproval(ctx, decision37(taskID, corr))
		out <- answer37{a, why}
	})
	t.Cleanup(func() {
		select {
		case <-h.Done():
		case <-time.After(5 * time.Second):
			t.Errorf("approval-waiter did not exit within 5s (a parked ask that never settles is the defect this ticket refuses)")
		}
	})
	return out
}

func decision37(taskID, corr string) tools.Decision {
	return tools.Decision{
		TaskID: taskID, Tool: "fs.write", Level: risk.L2,
		Reason:        "票 37 深度/超时/降级用例的待确认操作",
		CorrelationID: corr,
	}
}

// waitFor37 polls a condition the gate's own goroutine satisfies
// asynchronously, with a bounded real-time budget (the fixture-48 pattern).
func waitFor37(t *testing.T, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal(msg)
}

// grantOf37 reads one displayed card's native grant off the ledger - the value
// the native card's allow button would spend.
func grantOf37(t *testing.T, f *fixture37, corr string) string {
	t.Helper()
	card, ok := f.ra.cards.Look(corr)
	if !ok {
		t.Fatalf("the ledger never recorded card %q - the native reading below would prove nothing", corr)
	}
	if card.Grant == "" {
		t.Fatalf("card %q carries no native grant", corr)
	}
	return card.Grant
}

// ---------------------------------------------------------------------------
// cell 1 - depth badge 1->2, head-only, answers route by correlationId
// ---------------------------------------------------------------------------

// Test37TwoSequentialL2BadgeOneToTwoHeadOnlyRoutedByCorrelation is AC#1 in the
// pre-47 shape: ONE text task raises two L2 requests in sequence (the second
// while the first is still pending, simulated with the ticket-43/48 queue test
// mode), and the readings are
//
//	(1a) the orb badge input goes 1 -> 2: the first card the gate stamped
//	     carries Depth 1, the second Depth 2, and approvalDepthBadge (the
//	     producer of ball.SetBadge) maps them verbatim;
//	(1b) the panel face shows the HEAD only (corr A at depth 1) while card B
//	     stays reachable at its own FIFO place;
//	(1c) one explicit key reconciliation, the ticket-87 R-3 acceptance item:
//	     the corr the reply must carry == the queue's issued key == the corr
//	     on the displayed card == the corr the panel verdict carries;
//	(1d) an answer that names the SECOND card settles the SECOND card's ask
//	     and nobody else's - a reply routed by arrival order would wake the
//	     first card here, and the negative check below is exactly that red;
//	(1e) a grant spent across cards is a forgery (ErrBadGrant), nothing
//	     settles, and after both cards are answered the badge inputs read 0.
func Test37TwoSequentialL2BadgeOneToTwoHeadOnlyRoutedByCorrelation(t *testing.T) {
	f := newFixture37(t, nil)
	const task = "task-37" // pre-47: one task, two sequential L2 asks

	askA := ask37(t, f, task, "corr-37-a")
	waitFor37(t, func() bool { return f.ra.gate.Queue().Depth() == 1 }, "card A never queued")
	pA, ok := f.ui.promptOf("corr-37-a")
	if !ok {
		t.Fatal("card A was never displayed")
	}
	if pA.Depth != 1 {
		t.Errorf("card A's stamped depth = %d, want 1 (first pending card)", pA.Depth)
	}
	if n := approvalDepthBadge(pA); n != 1 {
		t.Errorf("approvalDepthBadge(card A) = %d, want 1 (the badge the orb posts after the first ask)", n)
	}

	askB := ask37(t, f, task, "corr-37-b")
	waitFor37(t, func() bool { return f.ra.gate.Queue().Depth() == 2 }, "card B never queued (the two asks collapsed)")
	pB, ok := f.ui.promptOf("corr-37-b")
	if !ok {
		t.Fatal("card B was never displayed")
	}
	if pB.Depth != 2 {
		t.Errorf("card B's stamped depth = %d, want 2 (two sends are TWO pending cards, not one)", pB.Depth)
	}
	if n := approvalDepthBadge(pB); n != 2 {
		t.Errorf("approvalDepthBadge(card B) = %d, want 2 (the badge goes 1 -> 2 on the second ask)", n)
	}

	// (1b) head-only: the panel face carries the FIRST arrival; card B is not
	// the displayed head, it sits at its own FIFO place.
	h, ok := f.ra.gate.Panel().Head()
	if !ok || h.CorrelationID != "corr-37-a" || h.Depth != 1 {
		t.Errorf("panel head = ok %v %+v, want corr-37-a at depth 1 (队头单显)", ok, h)
	}
	if v, ok := f.ra.gate.Panel().View("corr-37-b"); !ok || v.Depth != 2 {
		t.Errorf("View(corr-37-b) = ok %v %+v, want depth 2 (its own FIFO place)", ok, v)
	}

	// (1c) the one explicit key reconciliation: queue-issued key == displayed
	// card corr == panel verdict corr. Today they agree by construction; the
	// pairing below is the reading that makes a future divergence loud here
	// instead of at a user's click (ticket-87 R-3).
	live := f.ra.gate.Queue().LiveApprovals()
	if len(live) != 2 || live[0].CorrelationID != "corr-37-a" || live[1].CorrelationID != "corr-37-b" {
		t.Fatalf("issued keys = %+v, want [corr-37-a corr-37-b] in FIFO order", live)
	}
	for _, it := range live {
		card, ok := f.ui.promptOf(it.CorrelationID)
		if !ok || card.CorrelationID != it.CorrelationID {
			t.Errorf("the displayed card corr %v != the issued key %q (the C17 inbound key and the queue key diverged)", ok, it.CorrelationID)
		}
		if v, ok := f.ra.gate.Panel().View(it.CorrelationID); !ok || v.CorrelationID != it.CorrelationID {
			t.Errorf("the panel verdict corr for %q diverged (ok %v)", it.CorrelationID, ok)
		}
	}

	// (1d) the anti-arrival-order attack, done BEFORE anything settles: card
	// A's grant spent on card B must be a forgery, and a reply that names B
	// must wake B's ask while A's stays parked.
	grantA := grantOf37(t, f, "corr-37-a")
	if err := f.ra.gate.Native().Allow(context.Background(), "corr-37-b", grantA); !errors.Is(err, approval.ErrBadGrant) {
		t.Errorf("card A's grant spent on card B = %v, want ErrBadGrant (a seat is not an identity)", err)
	}
	if got := f.ra.gate.Queue().Depth(); got != 2 {
		t.Fatalf("depth after the forged allow = %d, want 2 (nothing may have settled)", got)
	}
	if err := f.ra.cards.Reject(context.Background(), "corr-37-b", "票 37 先答第二发"); err != nil {
		t.Fatalf("reject by corr-37-b: %v", err)
	}
	select {
	case a := <-askB:
		if a.a != tools.AnswerReject {
			t.Errorf("corr-37-b's ask woke with %s, want reject", a.a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the corr-named card's ask never woke - the reply did not route by correlationId")
	}
	select {
	case a := <-askA:
		t.Fatalf("cross-delivery: the FIRST card's ask woke (%+v) although the reply named corr-37-b - that is routing by arrival order", a)
	default:
	}
	if got := f.ra.gate.Queue().Depth(); got != 1 {
		t.Fatalf("depth after answering B = %d, want 1", got)
	}
	if h, ok := f.ra.gate.Panel().Head(); !ok || h.CorrelationID != "corr-37-a" {
		t.Errorf("after answering B the head = ok %v %+v, want card A still (FIFO head did not move)", ok, h)
	}
	// The badge re-sync input after one card left: the queue depth the orb
	// posts on the dismissal edge (syncBadge) is 1, not the stale 2.
	if got := f.ra.gate.Queue().Depth(); got != 1 {
		t.Errorf("the badge re-sync input = %d, want 1", got)
	}

	// (1e) the first card is answered through the native card verb, and the
	// drain ends with zeroed badge inputs.
	if err := f.ra.cards.Allow(context.Background(), "corr-37-a"); err != nil {
		t.Fatalf("native allow of corr-37-a: %v", err)
	}
	select {
	case a := <-askA:
		if a.a != tools.AnswerAllow {
			t.Errorf("corr-37-a's ask woke with %s, want allow", a.a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the first card's ask never woke after its native allow")
	}
	if got := f.ra.gate.Queue().Depth(); got != 0 {
		t.Errorf("badge input after the drain = %d, want 0", got)
	}
	if _, awaiting := f.ra.cards.AwaitingHuman(); awaiting {
		t.Errorf("the ledger still reports a card awaiting after both answered - the orb would keep the badge up")
	}
	t.Logf("37 AC#1 readings: badge 1->2 (stamped depths 1,2), head-only held, replies routed by corr (B first, A untouched), forged cross-card grant ErrBadGrant, drain to 0")
}

// Test37TwoCardsCoexistingCancelOneNeverHitsTheOther is ticket 37's progress
// note acceptance item (the ticket-87 R-3 reconciliation, 2026-09-21 entry):
// with TWO cards coexisting, cancelling one must never hit the other. Two
// collision shapes are pinned:
//
//	(a) cancel by an EXACT issued key settles only that card - the sibling's
//	    ask stays parked and its grant stays spendable;
//	(b) cancel by a name BOTH cards answer to (their shared task id - the
//	    C17 inbound key the host used to stamp both, the exact collision the
//	    87 report's attack 2 exploited) refuses as ambiguous for BOTH: the
//	    reply means neither card and neither ask wakes. The same name becomes
//	    spendable for the single survivor afterwards - the refusal-direction
//	    leniency is unambiguous-only, and that boundary is pinned here too.
func Test37TwoCardsCoexistingCancelOneNeverHitsTheOther(t *testing.T) {
	f := newFixture37(t, nil)
	const task = "task-37-cx" // the shared inbound name both cards answer to

	askA := ask37(t, f, task, "corr-37-cx-a")
	waitFor37(t, func() bool { return f.ra.gate.Queue().Depth() == 1 }, "card A never queued")
	askB := ask37(t, f, task, "corr-37-cx-b")
	waitFor37(t, func() bool { return f.ra.gate.Queue().Depth() == 2 }, "card B never queued")
	grantA := grantOf37(t, f, "corr-37-cx-a")

	// (b) FIRST, while BOTH cards are live: the shared name names TWO cards
	// -> ambiguous -> refused for both, never resolved by arrival order.
	if err := f.ra.gate.Native().Reject(task, "票 37 共享名取消"); !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("a cancel by the shared inbound name with both cards live = %v, want ErrUnknownCorrelation (two cards answer that name; it must not pick one)", err)
	}
	select {
	case a := <-askA:
		t.Fatalf("the ambiguous cancel settled card A (%+v) - exactly the cancel-one-hits-the-other defect", a)
	default:
	}
	select {
	case a := <-askB:
		t.Fatalf("the ambiguous cancel settled card B (%+v) - exactly the cancel-one-hits-the-other defect", a)
	default:
	}
	if got := f.ra.gate.Queue().Depth(); got != 2 {
		t.Fatalf("depth after the ambiguous cancel = %d, want 2 (nothing settled)", got)
	}

	// (a) cancel B by its exact issued key: only B settles.
	if err := f.ra.gate.Native().Reject("corr-37-cx-b", "票 37 按键取消 B"); err != nil {
		t.Fatalf("cancel by B's issued key: %v", err)
	}
	select {
	case a := <-askB:
		if a.a != tools.AnswerReject {
			t.Errorf("card B answered %s, want reject", a.a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("card B's ask never woke after its own key was cancelled")
	}
	select {
	case a := <-askA:
		t.Fatalf("cancel hit the WRONG card: card A woke (%+v) although the key named B", a)
	default:
	}
	if got := f.ra.gate.Queue().Depth(); got != 1 {
		t.Fatalf("depth after cancelling B = %d, want 1 (A must survive the cancel of B)", got)
	}

	// The boundary's other side: with B gone, the same name answers to ONE
	// live card, and the refusal-direction leniency lands it - the human's
	// 「no」 is a lost vote only while it is genuinely ambiguous.
	if err := f.ra.gate.Native().Reject(task, "票 37 单活卡共享名取消"); err != nil {
		t.Fatalf("the shared name after B left = %v, want it to land on the single survivor", err)
	}
	select {
	case a := <-askA:
		if a.a != tools.AnswerReject {
			t.Errorf("the survivor answered %s, want reject", a.a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the survivor's ask never woke after the unambiguous cancel")
	}
	// And A's grant never spent: the native allow with the still-held nonce
	// on a settled corr is refused, the ledger forgot it, the badge reads 0.
	if err := f.ra.gate.Native().Allow(context.Background(), "corr-37-cx-a", grantA); !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("the settled card's grant still spent: %v, want ErrUnknownCorrelation", err)
	}
	if got := f.ra.gate.Queue().Depth(); got != 0 {
		t.Errorf("badge input after both cancels = %d, want 0", got)
	}
	t.Logf("37 R-3 readings: exact-key cancel hit only its card, shared-name cancel refused ErrUnknownCorrelation while ambiguous, leniency landed once unambiguous, badge 0")
}

// ---------------------------------------------------------------------------
// cell 2 - 270s warning, 300s auto-reject, replay, badge clears
// ---------------------------------------------------------------------------

// fakeClock37 is the approval.Clock seam's test double (the ticket-11 pattern,
// ledger A1886): one time source for stamps and timers, advanced by the test.
// No wall-clock difference is ever a judging criterion here.
type fakeClock37 struct {
	mu     sync.Mutex
	now    time.Time
	timers []fakeTimer37
}

type fakeTimer37 struct {
	at    time.Time
	c     chan time.Time
	fired bool
}

func newFakeClock37() *fakeClock37 { return &fakeClock37{now: time.Unix(0, 0)} }

func (c *fakeClock37) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock37) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := &fakeTimer37{at: c.now.Add(d), c: make(chan time.Time, 1)}
	if d <= 0 {
		t.c <- c.now
		t.fired = true
		return t.c
	}
	c.timers = append(c.timers, *t)
	return t.c
}

// Advance moves the single time source and fires what is due, oldest first,
// so a warn-then-deadline sequence is observed in order.
func (c *fakeClock37) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	var fire []fakeTimer37
	rest := c.timers[:0]
	for _, tm := range c.timers {
		if !tm.at.After(c.now) {
			fire = append(fire, tm)
			continue
		}
		rest = append(rest, tm)
	}
	c.timers = rest
	for i := 1; i < len(fire); i++ {
		for j := i; j > 0 && fire[j].at.Before(fire[j-1].at); j-- {
			fire[j], fire[j-1] = fire[j-1], fire[j]
		}
	}
	for _, tm := range fire {
		if !tm.fired {
			tm.fired = true
			tm.c <- tm.at
		}
	}
	c.mu.Unlock()
}

// Test37TimeoutWarnsAt270AutoRejectsAt300ReplayWorksBadgeClears is AC#2 on the
// injected Clock at the REAL contract boundary (the gate is built with zero
// timeout options, so 300s/30s come from the queue's frozen defaults):
//
//	(2a) before 270s no warning; at 270s ONE warning event with the 30s
//	     remainder and the prominent wording - and the card is still
//	     answerable (the badge input is still 1);
//	(2b) at 300s the ask resolves to the timeout-that-refused with the
//	     auto-reject wording, the queue drains, and the badge inputs read 0
//	     (queue depth 0, ledger empty, ballCardUI.Update's real settle path
//	     runs with nothing left awaiting);
//	(2c) Gate.Replay re-opens the ORIGINAL request under a fresh corr with a
//	     FRESH grant - the dead card's grant cannot open the replay - and a
//	     native allow on the fresh grant completes the original request.
func Test37TimeoutWarnsAt270AutoRejectsAt300ReplayWorksBadgeClears(t *testing.T) {
	clk := newFakeClock37()
	f := newFixture37(t, clk)
	const task = "task-37-timeout"
	const corr = "corr-37-t"

	ans := ask37(t, f, task, corr)
	waitFor37(t, func() bool { return f.ra.gate.Queue().Depth() == 1 }, "the card never queued")
	p, ok := f.ui.promptOf(corr)
	if !ok {
		t.Fatal("the card was never displayed")
	}
	if p.Deadline != 300*time.Second {
		t.Errorf("the card's deadline reading = %s, want the frozen 300s C18 boundary (read off the card, not assumed)", p.Deadline)
	}
	// The dead card's grant, captured while it is still spendable - (2c) pins
	// that a spent nonce never carries over onto the replay.
	deadGrant := grantOf37(t, f, corr)

	// (2a) 270s: the warning is a reading, the deadline never moved.
	clk.Advance(270 * time.Second)
	waitFor37(t, func() bool { return len(f.ui.ofKind(approval.EventWarning)) == 1 }, "the 270s warning event never arrived")
	w := f.ui.ofKind(approval.EventWarning)[0]
	if w.CorrelationID != corr {
		t.Errorf("the warning names %q, want %q", w.CorrelationID, corr)
	}
	if w.Remaining != 30*time.Second {
		t.Errorf("the warning's remainder = %s, want the 30s lead", w.Remaining)
	}
	if !strings.Contains(w.Text, "30 秒后自动拒绝") {
		t.Errorf("the warning text = %q, want the prominent auto-reject wording", w.Text)
	}
	select {
	case a := <-ans:
		t.Fatalf("the card ended at 270s: %+v (a warning must not end a 300s card)", a)
	default:
	}
	if got := f.ra.gate.Queue().Depth(); got != 1 {
		t.Fatalf("badge input during the warning = %d, want 1 (the card is still answerable)", got)
	}

	// (2b) 300s: the auto-reject, the drain, the badge zero.
	clk.Advance(30 * time.Second)
	select {
	case a := <-ans:
		if a.a != tools.AnswerTimeout {
			t.Errorf("the 300s expiry resolved to %s, want the refusing timeout", a.a)
		}
		if !strings.Contains(a.why, "300") || !strings.Contains(a.why, "拒绝") {
			t.Errorf("the expiry why = %q, want the 300s auto-reject wording", a.why)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the card never expired at its 300s deadline")
	}
	if got := f.ra.gate.Queue().Depth(); got != 0 {
		t.Fatalf("badge input after the expiry = %d, want 0", got)
	}
	if left := f.ra.cards.Pending(); len(left) != 0 {
		t.Fatalf("the ledger still holds %d cards after the expiry, want 0", len(left))
	}
	if _, awaiting := f.ra.cards.AwaitingHuman(); awaiting {
		t.Fatal("the ledger still reports a human awaited after the expiry")
	}
	// The REAL settle path the orb runs on the dismissal event (settleOrb +
	// syncBadge): headless it is safe by construction (no ball, nothing
	// awaited), and running it proves the Update branch the badge zeroing
	// rides stays executable on this shape.
	if err := f.ra.ui.Update(context.Background(), approval.Event{
		Kind: approval.EventDismissed, CorrelationID: corr, Text: "300s auto-reject",
	}); err != nil {
		t.Fatalf("ballCardUI.Update on the dismissal: %v", err)
	}

	// (2c) replay completes the original request under a fresh corr and a
	// fresh grant; the dead grant must not open the replayed card.
	d2, fresh, err := f.ra.gate.Replay(context.Background(), corr, task)
	if err != nil {
		t.Fatalf("Replay(%q): %v", corr, err)
	}
	if fresh != corr+"#replay" {
		t.Errorf("the replay's fresh corr = %q, want %q#replay", fresh, corr)
	}
	orig := decision37(task, corr)
	if d2.Tool != orig.Tool || d2.Reason != orig.Reason {
		t.Errorf("the replayed request drifted: tool %q reason %q, want %q/%q", d2.Tool, d2.Reason, orig.Tool, orig.Reason)
	}
	replayAns := ask37(t, f, task, fresh)
	waitFor37(t, func() bool { return f.ra.gate.Queue().Depth() == 1 }, "the replayed card never queued")
	p2, ok := f.ui.promptOf(fresh)
	if !ok {
		t.Fatal("the replayed card was never displayed")
	}
	if p2.Grant == "" {
		t.Fatal("the replayed card carries no grant - re-displaying must mint a new one")
	}
	if p2.Grant == deadGrant {
		t.Fatal("the replayed card re-used the dead card's grant: a spent nonce must never carry over")
	}
	if err := f.ra.gate.Native().Allow(context.Background(), fresh, p2.Grant); err != nil {
		t.Fatalf("the fresh grant on the replayed card: %v", err)
	}
	select {
	case a := <-replayAns:
		if a.a != tools.AnswerAllow {
			t.Errorf("the replayed request resolved to %s, want allow (replay 补完原请求)", a.a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the replayed ask never came back after its fresh allow")
	}
	if got := f.ra.gate.Queue().Depth(); got != 0 {
		t.Errorf("badge input after the replay was answered = %d, want 0", got)
	}
	t.Logf("37 AC#2 readings: warn@270s (remainder 30s), auto-reject@300s (timeout->refuse), badge inputs 1->0, replay fresh corr %q with a fresh grant answered allow", fresh)
}

// ---------------------------------------------------------------------------
// cell 3 - panel runtime masked (the 33 fixture) => native card still approves
// ---------------------------------------------------------------------------

// Test37MaskedPanelL2StillApprovableViaNativeCard33fixture is AC#3 with the
// ticket-33 AC#5 fixture itself: the WebView2 runtime probe seam is masked,
// the panel refuses BY NAME (never fatally, the window is never created), the
// panel.unavailable event is emitted and Unavailable() reads ON - and in that
// masked world an L2 card is STILL approvable through the native card path
// (Replies.Allow, the ball/tray verb) while the panel route keeps refusing
// allows (F2, ErrPanelAllow) and burning the offered nonce. C27: the L2
// capability never disappears just because the panel cannot run; and the
// fallback is NOT a panel-side allow opening.
func Test37MaskedPanelL2StillApprovableViaNativeCard33fixture(t *testing.T) {
	withMaskedRuntime(t)
	logs := withCapturedSlog(t)
	mgr := NewPanelManager(&panel.ComposerDispatch{}, nil, filepath.Join(t.TempDir(), "panel-webview2"))

	// The mask lands: named refusal, flag on, no window.
	err := mgr.Show(context.Background())
	if !errors.Is(err, errPanelRuntimeMissing) {
		t.Fatalf("the masked runtime must refuse by name, got: %v", err)
	}
	if on, reason := mgr.Unavailable(); !on || reason == "" {
		t.Fatalf("the 33 fallback flag must read (true, reason) in the masked world, got (%t, %q)", on, reason)
	}
	if !strings.Contains(logs.String(), "panel.unavailable") {
		t.Errorf("the panel.unavailable event was not emitted, log buffer: %q", logs.String())
	}

	// In that masked world the approval flow runs on the resident components:
	// two L2 cards queue, and the native card verb answers BOTH, head first.
	f := newFixture37(t, nil)
	const task = "task-37-fallback"
	askA := ask37(t, f, task, "corr-37-fa")
	waitFor37(t, func() bool { return f.ra.gate.Queue().Depth() == 1 }, "fallback card A never queued")
	askB := ask37(t, f, task, "corr-37-fb")
	waitFor37(t, func() bool { return f.ra.gate.Queue().Depth() == 2 }, "fallback card B never queued")
	// The badge still reads the depth in the masked world - the orb corner
	// mark does not depend on the panel either.
	if pB, ok := f.ui.promptOf("corr-37-fb"); !ok || pB.Depth != 2 {
		t.Errorf("card B's stamped depth in the masked world = ok %v %d, want 2 (the badge never depended on the panel)", ok, pB.Depth)
	}

	// The native card approves the head, then the next card: the two asks
	// wake with their OWN answers.
	if err := f.ra.cards.Allow(context.Background(), "corr-37-fa"); err != nil {
		t.Fatalf("native allow of the head card in the masked world: %v", err)
	}
	select {
	case a := <-askA:
		if a.a != tools.AnswerAllow {
			t.Errorf("the head card answered %s, want allow (native approval works with the panel masked)", a.a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the head card's ask never woke after its native allow")
	}
	if err := f.ra.cards.Allow(context.Background(), "corr-37-fb"); err != nil {
		t.Fatalf("native allow of the second card: %v", err)
	}
	select {
	case a := <-askB:
		if a.a != tools.AnswerAllow {
			t.Errorf("the second card answered %s, want allow", a.a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the second card's ask never woke after its native allow")
	}
	if got := f.ra.gate.Queue().Depth(); got != 0 {
		t.Errorf("badge input after both native allows = %d, want 0", got)
	}

	// The fallback is NOT a panel-side allow opening: a THIRD card's real,
	// live grant offered over the panel route is refused ON THE ROUTE, the
	// offer burns the nonce (a token that reached the untrusted route must
	// not be spendable anywhere), and the card can then only be refused.
	askC := ask37(t, f, task, "corr-37-fc")
	waitFor37(t, func() bool { return f.ra.gate.Queue().Depth() == 1 }, "fallback card C never queued")
	grantC := grantOf37(t, f, "corr-37-fc")
	errP := f.ra.gate.DecideFromPanel(context.Background(), approval.Request{
		CorrelationID: "corr-37-fc", Allow: true, Grant: grantC, Source: "panel",
	})
	if !errors.Is(errP, approval.ErrPanelAllow) {
		t.Errorf("a panel-sourced allow = %v, want ErrPanelAllow (the fallback never opens an allow on the panel route)", errP)
	}
	if !f.logs.has("PANEL-ALLOW-REJECTED") {
		t.Errorf("the refused panel allow left no PANEL-ALLOW-REJECTED audit line (the gate's own audit face)")
	}
	if err := f.ra.gate.Native().Allow(context.Background(), "corr-37-fc", grantC); !errors.Is(err, approval.ErrBadGrant) {
		t.Errorf("the burned nonce still spent on the native route: %v, want ErrBadGrant", err)
	}
	if got := f.ra.gate.Queue().Depth(); got != 1 {
		t.Fatalf("depth after the panel-allow attack = %d, want 1 (the card survived, refused)", got)
	}
	if err := f.ra.gate.DecideFromPanel(context.Background(), approval.Request{
		CorrelationID: "corr-37-fc", Allow: false, Reason: "票 37 面板遮蔽世界的拒绝方向",
	}); err != nil {
		t.Fatalf("the panel refusal direction: %v", err)
	}
	select {
	case a := <-askC:
		if a.a != tools.AnswerReject {
			t.Errorf("card C answered %s, want reject", a.a)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("card C's ask never woke after the panel refusal")
	}
	t.Logf("37 AC#3 readings: masked runtime refused by name (flag on, panel.unavailable), two L2 cards approved via the native card verb, panel allow refused ErrPanelAllow + nonce burned, refusal direction open")
}
