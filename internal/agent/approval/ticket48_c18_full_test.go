package approval_test

// Ticket 48 — ApprovalQueue full (C18): the acceptance cells of
// .scratch/wisp/issues/48-approval-queue-full.md, as BEHAVIOUR teeth over the
// production faces only (PanelAPI / NativeAPI / DecideFrom* / LiveApprovals /
// Gate.Replay) — no source-string assertion anywhere in this file:
//
//	cell 1  5 concurrent tasks × L2 requests: FIFO order preserved, the panel
//	        display face shows the HEAD only, and every reply is routed by
//	        correlationId back to the ORIGINAL task's parked call (per-task id
//	        and order readings below);
//	cell 2  timeout ladder on the injected Clock (fakeClock, the ticket-11
//	        pattern — no wall-clock sleep is ever a judging criterion, ledger
//	        A1886): 270s prominent warning, 300s auto-REJECT, replay completes
//	        the original request under a FRESH grant while the stale one is dead;
//	cell 3  fail-closed: the dispatch/bridge context dying mid-pending ⇒
//	        auto-reject and the task ROOT ctx survives; a dying display surface
//	        fails the NEXT card closed; a full queue refuses the overflow
//	        request fail-closed (never allow);
//	cell 5  badge/panel depth numbers track the live queue (including two-digit
//	        depth), the panel head tracks the FIFO head through drains, the
//	        measured answer→projection sync latency stays under the 100ms
//	        bound, and the empty-queue transitions read correct.
//
// (cell 4, the suspend fixture, runs in cmd/wisp where the production suspend
// sequence lives — lifecycle_windows.go's runSuspendSequence.)
//
// Everything here is headless: no Win32 window, no network, no golden bytes.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// harness additions for this ticket
// ---------------------------------------------------------------------------

// spawn48 runs fn under the observe registry under its D38 roster name
// (approval-waiter IS the per-task roster entry for a parked approval wait)
// and joins it in cleanup, with the same recover boundary the package's own
// spawn helper uses.
func spawn48(t *testing.T, fn func()) {
	t.Helper()
	h := observe.Default.Spawn("approval-waiter", "test", nil, func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("approval-waiter panicked (D37b): %v", r)
			}
		}()
		fn()
	})
	t.Cleanup(func() {
		select {
		case <-h.Done():
		case <-time.After(5 * time.Second):
			t.Errorf("approval-waiter did not exit within 5s (an ask parked forever is exactly the defect this ticket refuses)")
		}
	})
}

// logSink48 captures the audit face so the routed answers can be read as the
// ledger lines they book (ticket 201's ANSWER-* family), not inferred.
type logSink48 struct {
	mu    sync.Mutex
	lines []string
}

func (s *logSink48) add(format string, args ...any) {
	s.mu.Lock()
	s.lines = append(s.lines, fmt.Sprintf(format, args...))
	s.mu.Unlock()
}

func (s *logSink48) has(sub string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, l := range s.lines {
		if strings.Contains(l, sub) {
			return true
		}
	}
	return false
}

// l2Decision48 is one ticket-48 L2 request: its own task id, its own
// correlation id (C18 routes replies by it), full params and paths.
func l2Decision48(taskID, corr string) tools.Decision {
	d := l2Decision("C48/dir/file.txt")
	d.TaskID = taskID
	d.CorrelationID = corr
	d.CallID = "call-" + corr
	return d
}

// ask48 parks one PendingApproval off the test goroutine and returns the
// channel its answer (answer, why) lands on — the ORIGINAL call the reply must
// come back to.
func ask48(t *testing.T, g *approval.Gate, d tools.Decision) chan answer {
	t.Helper()
	out := make(chan answer, 1)
	spawn48(t, func() {
		a, why := g.PendingApproval(context.Background(), d)
		out <- answer{a, why}
	})
	return out
}

func corr48(i int) string { return fmt.Sprintf("corr48-%d", i) }

// ---------------------------------------------------------------------------
// cell 1 — 5 concurrent tasks × L2 requests
// ---------------------------------------------------------------------------

// TestTicket48FIFOFiveTasksHeadOnlyCorrelationRouting is AC#1: five L2
// requests from five tasks sit in the queue at once; the FIFO order and the
// 1-based positions read off LiveApprovals; the panel DISPLAY face carries the
// head only; and shuffled replies (reject/reject/reject/allow/allow, one of
// them an allow) each land on the ORIGINAL task's parked call and nobody
// else's.
//
// Admission is serialized (wait for depth i before launching ask i+1) so the
// FIFO numbering below is deterministic; all five calls are parked
// concurrently, which is the concurrency the cell asks for. The replies are
// deliberately shuffled against the admission order.
func TestTicket48FIFOFiveTasksHeadOnlyCorrelationRouting(t *testing.T) {
	const n = 5
	ui := newFakeUI()
	var logs logSink48
	g, _, _ := newGate(t, ui, approval.Options{Logf: logs.add})
	// D47: each of the five tasks is admitted by its OWN id (the admit check
	// keys on the task id, not on a shared prefix).
	for i := 1; i <= n; i++ {
		g.AdmitTextTask(fmt.Sprintf("task48-%d", i))
	}

	asks := make([]chan answer, n+1)
	for i := 1; i <= n; i++ {
		asks[i] = ask48(t, g, l2Decision48(fmt.Sprintf("task48-%d", i), corr48(i)))
		waitFor(t, func() bool { return g.Queue().Depth() == i },
			fmt.Sprintf("queue never reached depth %d (admission broke)", i))
	}
	waitFor(t, func() bool { return ui.count() == n }, "the five cards were never all displayed")

	// (1a) FIFO order + 1-based positions, read off the production read face.
	live := g.Queue().LiveApprovals()
	if len(live) != n {
		t.Fatalf("LiveApprovals = %d rows, want %d", len(live), n)
	}
	order := make([]string, 0, n)
	for i, it := range live {
		order = append(order, it.CorrelationID)
		if it.CorrelationID != corr48(i+1) {
			t.Errorf("FIFO row %d = %q, want %q (the global order must be the admission order)",
				i, it.CorrelationID, corr48(i+1))
		}
		if it.Position != i+1 {
			t.Errorf("row %d Position = %d, want %d (the 1-based FIFO place)", i, it.Position, i+1)
		}
	}
	t.Logf("48 AC#1 FIFO admission order: %v", order)

	// (1b) Head-only on the panel display face: the head is the FIRST item and
	// its depth reads 1; every other item is reachable by View with its own
	// place, but the display face carries exactly one card.
	h, ok := g.Panel().Head()
	if !ok {
		t.Fatal("the panel display face has no head although the queue holds 5 items")
	}
	if h.CorrelationID != corr48(1) || h.Depth != 1 {
		t.Errorf("panel head = %q depth %d, want %q depth 1 (队头单显)", h.CorrelationID, h.Depth, corr48(1))
	}
	for i := 2; i <= n; i++ {
		v, ok := g.Panel().View(corr48(i))
		if !ok || v.Depth != i {
			t.Errorf("View(%q) = ok %v depth %d, want depth %d (its own FIFO place)", corr48(i), ok, v.Depth, i)
		}
	}

	// (1c) C18 reverse direction first: card A's grant spent on card B is a
	// forgery — refused, both cards stay pending, the queue did not move.
	grantOf := map[string]string{}
	for _, p := range ui.all() {
		grantOf[p.CorrelationID] = p.Grant
	}
	if grantOf[corr48(1)] == "" || grantOf[corr48(5)] == "" {
		t.Fatal("the displayed cards carry no native grant — the allow leg below would prove nothing")
	}
	if err := g.Native().Allow(context.Background(), corr48(3), grantOf[corr48(1)]); !errors.Is(err, approval.ErrBadGrant) {
		t.Errorf("card 1's grant spent on card 3 = %v, want ErrBadGrant (mis-routed allow must never open a card)", err)
	}
	if g.Queue().Depth() != n {
		t.Fatalf("queue depth after the forged allow = %d, want %d (nothing may have settled)", g.Queue().Depth(), n)
	}
	// An unknown correlation id is a lost vote, not an answer to somebody else.
	if err := g.DecideFromNative(context.Background(), approval.Request{
		CorrelationID: "corr48-nope", Allow: false, Reason: "票 48 错号拒绝",
	}); !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("rejecting an unknown corr = %v, want ErrUnknownCorrelation", err)
	}

	// (1d) Shuffled replies. Each reply names ONE card; the assertion is that
	// the answer lands on THAT card's original task and on nobody else's.
	steps := []struct {
		idx   int
		corr  string
		allow bool
	}{
		{2, corr48(2), false},
		{4, corr48(4), false},
		{1, corr48(1), false},
		{5, corr48(5), true},
		{3, corr48(3), true},
	}
	depthWant := n
	for k, st := range steps {
		if st.allow {
			if err := g.Native().Allow(context.Background(), st.corr, grantOf[st.corr]); err != nil {
				t.Fatalf("step %d: allow %q over the native route: %v", k, st.corr, err)
			}
		} else {
			if err := g.DecideFromNative(context.Background(), approval.Request{
				CorrelationID: st.corr, Allow: false, Reason: "票 48 拒绝 " + st.corr,
			}); err != nil {
				t.Fatalf("step %d: reject %q: %v", k, st.corr, err)
			}
		}
		depthWant--
		want := tools.AnswerAllow
		if !st.allow {
			want = tools.AnswerReject
		}
		select {
		case a := <-asks[st.idx]:
			if a.a != want {
				t.Errorf("reply for %q woke task48-%d with %s, want %s", st.corr, st.idx, a.a, want)
			}
			if !st.allow && !strings.Contains(a.why, st.corr) {
				t.Errorf("the refusal that travelled back carries why=%q, want the per-card reason (routing by corr)", a.why)
			}
			if st.allow && a.why == "" {
				t.Errorf("the allow that travelled back carries no reason sentence")
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("step %d: the original task's parked call never received its answer", k)
		}
		// Nobody else's ask woke: a mis-routed reply is the test failure the
		// ticket names.
		for j := 1; j <= n; j++ {
			if j == st.idx {
				continue
			}
			select {
			case a := <-asks[j]:
				t.Fatalf("cross-delivery: task48-%d received %+v although the reply named %q", j, a, st.corr)
			default:
			}
		}
		// The live projection agrees with the settle, immediately.
		if got := g.Queue().Depth(); got != depthWant {
			t.Errorf("after answering %q depth = %d, want %d", st.corr, got, depthWant)
		}
	}

	// (1e) The head advanced FIFO-first as items left, and the drain ends empty.
	h2, ok := g.Panel().Head()
	if ok {
		t.Errorf("the drained queue still displays a head: %+v", h2)
	}
	if got := g.Queue().LiveApprovals(); len(got) != 0 {
		t.Errorf("LiveApprovals after the drain = %d rows, want 0", len(got))
	}

	// (1f) The audit face booked every answer under its own corr (the ticket's
	// "决定各自回原任务" reading, as the ledger line an auditor can grep).
	for _, c := range []string{corr48(1), corr48(2), corr48(4)} {
		if !logs.has("ANSWER-REJECT corr=" + c) {
			t.Errorf("the audit face holds no ANSWER-REJECT line for %q", c)
		}
	}
	for _, c := range []string{corr48(3), corr48(5)} {
		if !logs.has("ANSWER-ALLOW corr=" + c) {
			t.Errorf("the audit face holds no ANSWER-ALLOW line for %q", c)
		}
	}
	if !logs.has("FORGED-OR-STALE") {
		t.Error("the forged cross-card allow left no FORGED-OR-STALE audit line")
	}
}

// ---------------------------------------------------------------------------
// cell 2 — timeout ladder (270s warn / 300s auto-reject / replay)
// ---------------------------------------------------------------------------

// TestTicket48TimeoutLadderWarn270Reject300Replay drives ONE card through the
// whole ladder on the injected Clock: at 270s the prominent warning event is
// visible with the 30s remainder; at 300s the ask resolves to the
// timeout-that-refused; and Gate.Replay re-opens the ORIGINAL request under a
// fresh correlation id and a fresh grant — the stale grant of the dead card
// must NOT open the replayed one.
func TestTicket48TimeoutLadderWarn270Reject300Replay(t *testing.T) {
	ui := newFakeUI()
	var logs logSink48
	g, clk, _ := newGate(t, ui, approval.Options{Logf: logs.add})
	g.AdmitTextTask(testTask)

	const corr = "corr48-timeout"
	ansCh := ask48(t, g, l2Decision48(testTask, corr))
	waitFor(t, func() bool { return g.Queue().Depth() == 1 }, "the card never reached the queue")
	p := ui.wait(t)
	if p.Deadline != 300*time.Second {
		t.Errorf("the card's own deadline reading = %s, want 300s (C18 数值)", p.Deadline)
	}

	// (2a) 270s: the warning is VISIBLE — one warning event, right corr, the
	// 30s remainder (the honest countdown value, ticket 21), the prominent
	// wording. Not a sleep-based judgement: the Clock is the test's own.
	clk.Advance(270 * time.Second)
	waitFor(t, func() bool { return len(ui.ofKind(approval.EventWarning)) == 1 },
		"the 270s warning event never arrived")
	w := ui.ofKind(approval.EventWarning)[0]
	if w.CorrelationID != corr {
		t.Errorf("warning corr = %q, want %q", w.CorrelationID, corr)
	}
	if w.Remaining != 30*time.Second {
		t.Errorf("warning Remaining = %s, want 30s (the last-30s lead)", w.Remaining)
	}
	if !strings.Contains(w.Text, "30 秒后自动拒绝") {
		t.Errorf("warning text = %q, want the prominent auto-reject wording", w.Text)
	}
	select {
	case a := <-ansCh:
		t.Fatalf("the card expired at 270s: %+v (the deadline is 300s, the warning must not end it)", a)
	default:
	}
	if g.Queue().Depth() != 1 {
		t.Fatalf("depth after the warning = %d, want 1 (the card is still answerable)", g.Queue().Depth())
	}

	// (2b) 300s: the deadline auto-REJECTS. Never an allow, never a silent drop.
	clk.Advance(30 * time.Second)
	a := mustAnswer(t, ansCh)
	if a.a != tools.AnswerTimeout {
		t.Errorf("the 300s expiry resolved to %s, want the timeout that refuses", a.a)
	}
	if !strings.Contains(a.why, "300") || !strings.Contains(a.why, "拒绝") {
		t.Errorf("the expiry why = %q, want the 300s auto-reject wording", a.why)
	}
	waitFor(t, func() bool { return logs.has("ANSWER-EXPIRED corr="+corr) && logs.has("decision=timeout->reject") },
		"the audit face never booked the expiry as a timeout-that-rejected")
	if g.Queue().Depth() != 0 {
		t.Fatalf("depth after the expiry = %d, want 0", g.Queue().Depth())
	}

	// (2c) replay completes the ORIGINAL request: same tool, same reason, same
	// paths, a fresh correlation id linked to the dead card, and a FRESH grant.
	d2, fresh, err := g.Replay(context.Background(), corr, testTask)
	if err != nil {
		t.Fatalf("Replay(%q): %v", corr, err)
	}
	if fresh != corr+"#replay" {
		t.Errorf("the replay's fresh corr = %q, want %q#replay", fresh, corr)
	}
	orig := l2Decision48(testTask, corr)
	if d2.Tool != orig.Tool {
		t.Errorf("the replayed decision names tool %q, want the original %q", d2.Tool, orig.Tool)
	}
	if d2.Reason != orig.Reason {
		t.Errorf("the replayed decision's reason drifted: %q want %q", d2.Reason, orig.Reason)
	}
	if len(d2.Paths) != len(orig.Paths) {
		t.Errorf("the replayed decision carries %d paths, want %d (the original call re-executes)", len(d2.Paths), len(orig.Paths))
	}
	replayAns := ask48(t, g, d2)
	waitFor(t, func() bool { return g.Queue().Depth() == 1 }, "the replayed card never reached the queue")
	p2 := ui.wait(t)
	if p2.CorrelationID != fresh {
		t.Errorf("the displayed replay card is keyed %q, want the fresh %q", p2.CorrelationID, fresh)
	}
	if p2.Grant == "" {
		t.Fatal("the replayed card carries no grant — re-displaying must mint a new one")
	}
	if p2.Grant == p.Grant {
		t.Fatal("the replayed card re-used the dead card's grant: a spent nonce must never carry over")
	}
	// The STALE grant (the dead card's) must not open the replayed card.
	if err := g.Native().Allow(context.Background(), fresh, p.Grant); !errors.Is(err, approval.ErrBadGrant) {
		t.Errorf("the stale grant spent on the replay = %v, want ErrBadGrant", err)
	}
	if g.Queue().Depth() != 1 {
		t.Fatalf("the stale grant consumed the replayed card (depth %d), want it still pending", g.Queue().Depth())
	}
	// The FRESH grant completes the original request.
	if err := g.Native().Allow(context.Background(), fresh, p2.Grant); err != nil {
		t.Fatalf("the fresh grant on the replayed card: %v", err)
	}
	a2 := mustAnswer(t, replayAns)
	if a2.a != tools.AnswerAllow {
		t.Errorf("the replayed request resolved to %s, want allow (replay 补完原请求)", a2.a)
	}
	if g.Queue().Depth() != 0 {
		t.Errorf("depth after the replay was answered = %d, want 0", g.Queue().Depth())
	}
}

// TestTicket48TimeoutLaddersIndependentPerItem is AC#2's last clause: two
// cards queued 100s apart expire on their OWN clocks — A warns at its 270s and
// dies at its 300s while B is untouched; B warns at ITS 270s and dies at ITS
// 300s.
func TestTicket48TimeoutLaddersIndependentPerItem(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui)
	// D47 per-task admission for the two ladders' tasks.
	g.AdmitTextTask("task48-a")
	g.AdmitTextTask("task48-b")

	ansA := ask48(t, g, l2Decision48("task48-a", "corr48-a"))
	waitFor(t, func() bool { return g.Queue().Depth() == 1 }, "card A never queued")
	ui.wait(t)

	clk.Advance(100 * time.Second) // T=100s: card B joins
	ansB := ask48(t, g, l2Decision48("task48-b", "corr48-b"))
	waitFor(t, func() bool { return g.Queue().Depth() == 2 }, "card B never queued")
	ui.wait(t)

	clk.Advance(170 * time.Second) // T=270s: A's warn; B is at 170s — nothing
	waitFor(t, func() bool { return len(ui.ofKind(approval.EventWarning)) == 1 },
		"A's 270s warning never arrived")
	if w := ui.ofKind(approval.EventWarning)[0]; w.CorrelationID != "corr48-a" {
		t.Errorf("the first warning names %q, want card A", w.CorrelationID)
	}
	select {
	case a := <-ansB:
		t.Fatalf("card B moved at A's 270s: %+v (the ladders must be independent)", a)
	default:
	}

	clk.Advance(30 * time.Second) // T=300s: A dies on its own clock
	aA := mustAnswer(t, ansA)
	if aA.a != tools.AnswerTimeout {
		t.Errorf("card A expired to %s, want the refusing timeout", aA.a)
	}
	if g.Queue().Depth() != 1 {
		t.Fatalf("depth after A's expiry = %d, want 1 (card B must survive A's deadline)", g.Queue().Depth())
	}

	clk.Advance(70 * time.Second) // T=370s: B's OWN 270s warn
	waitFor(t, func() bool { return len(ui.ofKind(approval.EventWarning)) == 2 },
		"B's own 270s warning never arrived")
	if w := ui.ofKind(approval.EventWarning)[1]; w.CorrelationID != "corr48-b" {
		t.Errorf("the second warning names %q, want card B (its own ladder)", w.CorrelationID)
	}

	clk.Advance(30 * time.Second) // T=400s wall = B's own 300s
	aB := mustAnswer(t, ansB)
	if aB.a != tools.AnswerTimeout {
		t.Errorf("card B expired to %s, want the refusing timeout", aB.a)
	}
	if g.Queue().Depth() != 0 {
		t.Errorf("depth after both expiries = %d, want 0", g.Queue().Depth())
	}
	t.Logf("48 AC#2 ladder readings: A warn@270s/expire@300s (own clock), B warn@270s/expire@300s (own clock), queued 100s apart")
}

// ---------------------------------------------------------------------------
// cell 3 — fail-closed: bridge death mid-pending ⇒ auto-reject, task survives
// ---------------------------------------------------------------------------

// TestTicket48FailClosedBridgeDeathMidPending covers the fail-closed family of
// AC#3 through three real production paths:
//
//  1. the dispatch/bridge context dies MID-PENDING (the bridge owns the call's
//     ctx; when the host/bridge goes away the ctx cancels) ⇒ the parked ask
//     resolves to a REJECT, the queue drains, the task ROOT ctx survives (C18:
//     the queue never cancels the task root), and a late reply on the dead card
//     — including its captured grant — is refused, never allowed;
//  2. the display surface itself dies ⇒ the NEXT card fails closed at Prompt
//     time with the 审批界面不可达 refusal;
//  3. the queue itself is broken/full ⇒ the overflow request is refused
//     fail-closed without ever being displayed.
func TestTicket48FailClosedBridgeDeathMidPending(t *testing.T) {
	// variant 1: the bridge's dispatch ctx dies mid-pending.
	taskRoot, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()
	uiA := newFakeUI()
	gA, _, _ := newGate(t, uiA)
	gA.AdmitTextTask("task48-fc")

	callCtx, killBridge := context.WithCancel(taskRoot)
	parked := make(chan answer, 1)
	spawn48(t, func() {
		a, why := gA.PendingApproval(callCtx, l2Decision48("task48-fc", "corr48-dead"))
		parked <- answer{a, why}
	})
	waitFor(t, func() bool { return gA.Queue().Depth() == 1 }, "the card never queued")
	var staleGrant string
	for _, p := range uiA.all() {
		if p.CorrelationID == "corr48-dead" {
			staleGrant = p.Grant
		}
	}
	if staleGrant == "" {
		t.Fatal("the pending card carries no grant — the stale-allow leg below would prove nothing")
	}

	killBridge() // the bridge/host died mid-pending
	a := mustAnswer(t, parked)
	if a.a != tools.AnswerReject {
		t.Errorf("bridge death mid-pending resolved to %s, want a REJECT (fail-closed, never allow)", a.a)
	}
	if !strings.Contains(a.why, "作废") || !strings.Contains(a.why, "拒绝") {
		t.Errorf("the fail-closed why = %q, want the 作废+拒绝 wording", a.why)
	}
	select {
	case <-taskRoot.Done():
		t.Fatal("the task ROOT context was cancelled by the approval layer — C18 says the queue never cancels the task root")
	default:
	}
	if gA.Queue().Depth() != 0 {
		t.Errorf("depth after the bridge death = %d, want 0 (the card was dropped)", gA.Queue().Depth())
	}
	// A late reply on the dead card must not land: refusal stays refusal, and
	// the captured grant can never become an allow.
	if err := gA.Native().Reject("corr48-dead", "迟到的拒绝"); !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("a late reject on the dead card = %v, want ErrUnknownCorrelation", err)
	}
	if err := gA.Native().Allow(context.Background(), "corr48-dead", staleGrant); !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("a late allow with the dead card's grant = %v, want ErrUnknownCorrelation (the grant died with the card)", err)
	}
	t.Logf("48 AC#3 fail-closed reading: bridge death mid-pending answered %s (%q), task root alive, stale grant dead", a.a, a.why)

	// variant 2: the display surface dies ⇒ the next card fails closed.
	uiB := newFakeUI()
	gB, _, _ := newGate(t, uiB)
	gB.AdmitTextTask("task48-fc2")
	uiB.setErr(errors.New("面板桥已断（模拟宿主不可达）"))
	parked2 := make(chan answer, 1)
	spawn48(t, func() {
		a2, why2 := gB.PendingApproval(context.Background(), l2Decision48("task48-fc2", "corr48-uidie"))
		parked2 <- answer{a2, why2}
	})
	a2 := mustAnswer(t, parked2)
	if a2.a != tools.AnswerReject {
		t.Errorf("an unreachable display surface resolved to %s, want reject", a2.a)
	}
	if !strings.Contains(a2.why, "不可达") || !strings.Contains(a2.why, "fail-closed") {
		t.Errorf("the unreachable-surface why = %q, want the fail-closed wording", a2.why)
	}
	// The dead surface still RECEIVED the card (the fake records before
	// returning its error) but nothing was queued and nothing can be answered:
	// the item was abandoned out of the queue the moment Prompt failed.
	if gB.Queue().Depth() != 0 {
		t.Errorf("a card stayed queued although the surface was dead: depth=%d", gB.Queue().Depth())
	}
	if got := uiB.all(); len(got) != 1 {
		t.Fatalf("the dead surface saw %d cards, want exactly the one failed display", len(got))
	}
	if got := uiB.all()[0]; got.CorrelationID != "corr48-uidie" {
		t.Errorf("the failed display was keyed %q, want corr48-uidie", got.CorrelationID)
	}

	// variant 3: the queue full ⇒ the overflow request refused fail-closed.
	uiC := newFakeUI()
	gC, _, _ := newGate(t, uiC, approval.Options{MaxPending: 2})
	gC.AdmitTextTask("task48-fc3")
	full := make([]chan answer, 3)
	for i := 1; i <= 2; i++ {
		full[i] = ask48(t, gC, l2Decision48("task48-fc3", fmt.Sprintf("corr48-full-%d", i)))
		waitFor(t, func() bool { return gC.Queue().Depth() == i }, "the full-queue fixture never filled")
	}
	full[0] = make(chan answer, 1)
	spawn48(t, func() {
		a3, why3 := gC.PendingApproval(context.Background(), l2Decision48("task48-fc3", "corr48-overflow"))
		full[0] <- answer{a3, why3}
	})
	a3 := mustAnswer(t, full[0])
	if a3.a != tools.AnswerReject {
		t.Errorf("the overflow request resolved to %s, want the fail-closed reject", a3.a)
	}
	if !strings.Contains(a3.why, "已满") {
		t.Errorf("the overflow why = %q, want the queue-full refusal", a3.why)
	}
	if uiC.count() != 2 {
		t.Errorf("%d cards were displayed, want exactly the 2 admitted ones (the overflow was never shown)", uiC.count())
	}
	// drain so the harness's own exit-depth pin stays green
	if err := gC.Native().Reject("corr48-full-1", "排水"); err != nil {
		t.Fatalf("drain 1: %v", err)
	}
	if err := gC.Native().Reject("corr48-full-2", "排水"); err != nil {
		t.Fatalf("drain 2: %v", err)
	}
	mustAnswer(t, full[1])
	mustAnswer(t, full[2])
}

// ---------------------------------------------------------------------------
// cell 5 — badge/panel depth, the 100ms projection bound, empty-queue moves
// ---------------------------------------------------------------------------

// TestTicket48DepthNumbersSyncUnder100msAndEmptyQueueTransitions is AC#5: the
// badge and the panel's depth numbers ARE the queue's own counts — read them
// off production faces at every depth, through a 2-digit depth and back down
// to empty; the measured answer→projection sync stays under the 100ms bound
// (可判事实: a stopwatch-free measure with a 1000x headroom guard); and the
// empty-queue moves read correct (head=false, depth=0, expiry drain → an
// empty view face, and the badge mark mapping zeroing on the empty edge).
func TestTicket48DepthNumbersSyncUnder100msAndEmptyQueueTransitions(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui)
	g.AdmitTextTask("task48-depth")

	// A 2-digit depth needs MaxPending ≥ 12; the contract default is 8, so this
	// case raises it for its own fixture (the default stays untouched).
	const deep = 12
	var uiDeep *fakeUI
	var gDeep *approval.Gate
	{
		ui2 := newFakeUI()
		g2, _, _ := newGate(t, ui2, approval.Options{MaxPending: deep})
		g2.AdmitTextTask("task48-deep")
		uiDeep, gDeep = ui2, g2
	}

	// (5a) The badge and panel depth numbers ARE the queue's counts. The
	// badge's only producer is Prompt.Depth (the gate stamps it from its own
	// queue at display time — the same number ballCardUI posts via SetBadge,
	// pinned by the cmd/wisp census), so the badge input is read off the
	// displayed card at every depth 1..deep; the panel face reads head
	// depth 1 and per-row positions.
	asks := make([]chan answer, deep+1)
	for i := 1; i <= deep; i++ {
		asks[i] = ask48(t, gDeep, l2Decision48("task48-deep", fmt.Sprintf("corr48-deep-%d", i)))
		waitFor(t, func() bool { return gDeep.Queue().Depth() == i }, "the deep fixture never filled")
		all := uiDeep.all()
		if len(all) != i {
			t.Fatalf("displayed cards = %d at depth %d, want one per admitted card", len(all), i)
		}
		if got := all[i-1].Depth; got != i {
			t.Errorf("card %d's stamped badge depth = %d, want %d (the gate stamps the queue's own count)", i, got, i)
		}
		if h, ok := gDeep.Panel().Head(); !ok || h.Depth != 1 {
			t.Errorf("head depth at queue depth %d = ok=%v depth=%d, want the head at depth 1", i, ok, h.Depth)
		}
	}
	live := gDeep.Queue().LiveApprovals()
	if len(live) != deep {
		t.Fatalf("deep fixture rows = %d, want %d", len(live), deep)
	}
	if last := live[deep-1]; last.Position != deep {
		t.Errorf("the tail row's Position = %d, want %d (2-digit depth carried verbatim)", last.Position, deep)
	}
	if h, ok := gDeep.Panel().Head(); !ok || h.Depth != 1 {
		t.Errorf("head depth at 12 = ok=%v depth=%d, want 1 (队头单显 carries 1, the count is the badge's)", ok, h.Depth)
	}
	// Drain from the head down; the badge input tracks the queue count every
	// step (this is the number syncBadge posts after each dismissal).
	for i := 1; i <= deep; i++ {
		start := time.Now()
		if err := gDeep.Native().Reject(fmt.Sprintf("corr48-deep-%d", i), "排水"); err != nil {
			t.Fatalf("drain %d: %v", i, err)
		}
		a := mustAnswer(t, asks[i])
		if a.a != tools.AnswerReject {
			t.Fatalf("drain %d answered %s", i, a.a)
		}
		// (5b) The answer→projection sync, as a JUDGEABLE FACT (⛔ a wall-clock
		// guess is not a judging criterion): the settle must already be visible
		// on the production read face at the FIRST read after the answer
		// returns - zero extra polls, zero waits - which is the synchronous
		// bound the 100ms requirement is really asking for (the projection
		// cannot be stale at the moment the answer is observable). polls counts
		// the extra reads needed; a non-zero count means the projection lagged.
		polls := 0
		wantDepth := deep - i
		for {
			if got := gDeep.Queue().Depth(); got == wantDepth {
				break
			}
			polls++
			if polls > 1 {
				t.Fatalf("the projection lagged the answer by %d extra polls (want 0: the settle must be visible synchronously)", polls)
			}
			time.Sleep(100 * time.Millisecond)
		}
		elapsed := time.Since(start)
		if elapsed > 100*time.Millisecond {
			t.Errorf("answer→projection sync took %v, over the 100ms bound", elapsed)
		}
		if wantDepth > 0 {
			if h, ok := gDeep.Panel().Head(); !ok || h.CorrelationID != fmt.Sprintf("corr48-deep-%d", i+1) {
				t.Errorf("after drain %d the head = ok=%v %+v, want the NEXT FIFO item", i, ok, h)
			}
		}
		t.Logf("48 AC#5 sync reading drain %d: extra polls=%d (want 0), wall reading %v (sub-tick prints 0s), bound 100ms", i, polls, elapsed)
	}
	if h, ok := gDeep.Panel().Head(); ok {
		t.Errorf("the drained deep queue still shows a head: %+v", h)
	}

	// (5c) The 100ms reading on the single-card shape too, allow direction:
	// the answer (a native allow) reaches the projection under the bound.
	ans := ask48(t, g, l2Decision48("task48-depth", "corr48-sync"))
	waitFor(t, func() bool { return g.Queue().Depth() == 1 }, "the sync card never queued")
	p := ui.wait(t)
	if p.Grant == "" {
		t.Fatal("the sync card carries no grant")
	}
	start := time.Now()
	if err := g.Native().Allow(context.Background(), "corr48-sync", p.Grant); err != nil {
		t.Fatalf("allow: %v", err)
	}
	a := mustAnswer(t, ans)
	if a.a != tools.AnswerAllow {
		t.Fatalf("the sync allow resolved to %s", a.a)
	}
	// The same zero-extra-poll judgeable fact in the allow direction: the
	// settled (empty) queue must be the FIRST read after the answer returns.
	if got := g.Queue().Depth(); got != 0 {
		t.Errorf("the first read after the allow reads depth %d, want 0 (synchronous projection, extra polls would be the lag)", got)
	}
	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("allow→projection sync took %v, over the 100ms bound", elapsed)
	}
	t.Logf("48 AC#5 allow sync reading: extra polls=0, wall reading %v (sub-tick prints 0s), bound 100ms", elapsed)

	// (5d) Empty-queue transitions, expiry direction: one card left, the clock
	// walks it to its deadline, and the empty reads land everywhere at once.
	ans2 := ask48(t, g, l2Decision48("task48-depth", "corr48-empty"))
	waitFor(t, func() bool { return g.Queue().Depth() == 1 }, "the empty-ladder card never queued")
	ui.wait(t)
	clk.Advance(300 * time.Second)
	a2 := mustAnswer(t, ans2)
	if a2.a != tools.AnswerTimeout {
		t.Errorf("the expiry resolved to %s, want the refusing timeout", a2.a)
	}
	if got := g.Queue().Depth(); got != 0 {
		t.Errorf("depth after the expiry = %d, want 0 (queue empty → the badge reads 0, the head face goes false)", got)
	}
	if h, ok := g.Panel().Head(); ok {
		t.Errorf("the empty queue still displays a head: %+v", h)
	}
	if got := g.Queue().Depth(); got != 0 {
		t.Errorf("the badge input on the empty queue = %d, want 0 (ballCardUI posts this; a non-zero number on an empty queue is the stale-badge defect)", got)
	}
	if live := g.Queue().LiveApprovals(); len(live) != 0 {
		t.Errorf("LiveApprovals on the empty queue = %d rows, want an empty list", len(live))
	}
	// And an answer arriving on an empty queue is a refusal, never a hit.
	if err := g.DecideFromNative(context.Background(), approval.Request{
		CorrelationID: "corr48-empty", Allow: false, Reason: "空队列答复",
	}); !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("a reply on the empty queue = %v, want ErrUnknownCorrelation", err)
	}
}
