//go:build windows

package main

// Ticket 48 — ApprovalQueue full (C18), the two cells whose production paths
// live in cmd/wisp (the write surface this leg was scoped to):
//
//	cell 4  the suspend fixture: pending approvals VOIDED + counted + the
//	        「错过 N 个待确认操作」 notice, and after resume NO stale approval is
//	        accepted — the stale grant is dead, a stale reply is refused, and a
//	        fresh card needs a fresh grant. The fixture is ticket 43's own
//	        (newSuspendFixture43: a REAL approval gate, REAL queue item, REAL
//	        runSuspendSequence on the REAL components) — this file adds only
//	        ticket 48's teeth on top of it.
//	cell 5  the panel half of the depth sync: the queue's live rows joined
//	        into the packet's task roster (the same pure hop
//	        panel_pump.nativeVerdictsFromGate + panel.TaskRosterSectionFrom the
//	        production pump runs), the measured answer→join sync against the
//	        100ms bound with the real numbers logged, head-only display
//	        depth, and the empty-queue migration (no row lit, head face
//	        closed, badge input 0).
//
// Everything here is headless: no Win32 window is created, nothing touches the
// owner's ball process.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// captureStdout48 runs fn with os.Stdout pointed at a pipe and returns what
// was printed (the runSuspendSequence notice is a fmt.Printf, not a log line,
// so the sentence the ticket requires is read off the operator's stdout).
func captureStdout48(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	prev := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = prev })
	done := make(chan string, 1)
	go func() {
		var sb strings.Builder
		buf := make([]byte, 4096)
		for {
			n, rerr := r.Read(buf)
			sb.Write(buf[:n])
			if rerr != nil {
				break
			}
		}
		done <- sb.String()
	}()
	fn()
	w.Close()
	os.Stdout = prev
	out := <-done
	r.Close()
	return out
}

// waitFor48 polls a condition the gate's own goroutine satisfies
// asynchronously, with a bounded real-time budget.
func waitFor48(t *testing.T, cond func() bool, msg string) {
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

// ---------------------------------------------------------------------------
// cell 4 — suspend fixture: void + notice; resume refuses stale approvals
// ---------------------------------------------------------------------------

// Test48SuspendFixtureVoidsPendingAndResumeRefusesStale drives ticket 48's AC#4
// through the REAL suspend sequence (lifecycle_windows.go's runSuspendSequence,
// the production body the power door calls):
//
//  1. two pending L2 cards, FIFO positions 1/2, each parked on its own ask;
//  2. the suspend sequence voids BOTH through the queue's refusal funnel and
//     prints the 「错过 N 个待确认操作」 notice with the exact count;
//  3. the queue drains, the ledger forgets the cards, and the parked asks wake
//     with the reject + suspend wording;
//  4. after resume the STALE approval is not accepted: the card's captured
//     grant cannot allow anything, a reply routed for the stale corr is
//     refused, and a fresh post-resume card needs (and gets) a fresh grant that
//     differs from the dead one.
func Test48SuspendFixtureVoidsPendingAndResumeRefusesStale(t *testing.T) {
	f := newSuspendFixture43(t)

	// Raise one card at a time and WAIT for its admission before raising the
	// next: raiseL2 parks the ask off-goroutine, so without the gate the two
	// pushes could land in either order and the FIFO reading below would
	// measure the fixture, not the queue.
	f.raiseL2(t, "corr-48-a")
	waitFor48(t, func() bool { return f.ra.gate.Queue().Depth() == 1 }, "card A never queued")
	f.raiseL2(t, "corr-48-b")
	waitFor48(t, func() bool { return f.ra.gate.Queue().Depth() == 2 }, "card B never queued")
	cards := waitPending43(t, f.ra, 2)
	if cards[0].Level != "L2" || cards[1].Level != "L2" {
		t.Fatalf("pending levels = %q/%q, want two L2 cards", cards[0].Level, cards[1].Level)
	}
	staleGrant := cards[0].Grant
	if staleGrant == "" {
		t.Fatal("the first pending card carries no grant — the stale-approval leg below would prove nothing")
	}

	// FIFO reads before the suspend: head-first rows, 1-based positions.
	live := f.ra.gate.Queue().LiveApprovals()
	if len(live) != 2 {
		t.Fatalf("live rows before the suspend = %d, want 2", len(live))
	}
	if live[0].CorrelationID != "corr-48-a" || live[1].CorrelationID != "corr-48-b" {
		t.Errorf("FIFO order before the suspend = [%q %q], want [corr-48-a corr-48-b]",
			live[0].CorrelationID, live[1].CorrelationID)
	}
	if live[0].Position != 1 || live[1].Position != 2 {
		t.Errorf("positions before the suspend = %d/%d, want 1/2", live[0].Position, live[1].Position)
	}

	out := captureStdout48(t, func() {
		voided := runSuspendSequence(f.ra, f.raudio, f.rv, "48-suspend")
		if voided != 2 {
			t.Errorf("voided = %d, want 2 (both pending approvals)", voided)
		}
	})
	if !strings.Contains(out, "错过 2 个待确认操作") {
		t.Errorf("the suspend notice did not name the count: %q", strings.TrimSpace(out))
	}

	// The queue drained and the ledger forgot both cards.
	if d := f.ra.gate.Queue().Depth(); d != 0 {
		t.Errorf("queue depth after the suspend = %d, want 0", d)
	}
	if left := f.ra.cards.Pending(); len(left) != 0 {
		t.Fatalf("tracked cards after the suspend = %d, want 0", len(left))
	}
	if _, ok := f.ra.cards.Look("corr-48-a"); ok {
		t.Error("the stale card is still tracked after the suspend — the acceptance face must refuse it, not remember it")
	}

	// Both parked asks woke with the reject + suspend wording.
	for i := 0; i < 2; i++ {
		select {
		case ans := <-f.askRc:
			if ans != tools.AnswerReject {
				t.Errorf("parked ask %d answered %s, want reject", i, ans)
			}
		case <-time.After(10 * time.Second):
			t.Fatalf("parked ask %d never came back from the suspend void", i)
		}
		why := <-f.askWhy
		if !strings.Contains(why, "挂起") {
			t.Errorf("parked ask %d reason = %q, want the suspend wording", i, why)
		}
	}

	// After resume, the STALE approval is not accepted:
	// (a) the captured grant cannot allow anything any more;
	if err := f.ra.gate.Native().Allow(context.Background(), "corr-48-a", staleGrant); !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("a stale allow with the dead card's grant = %v, want ErrUnknownCorrelation", err)
	}
	// (b) a stale reply routed for the voided corr is refused, never settled;
	if err := f.ra.gate.Native().Reject("corr-48-a", "恢复后的陈旧拒绝"); !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("a stale reject after resume = %v, want ErrUnknownCorrelation", err)
	}
	// (c) a fresh post-resume card needs a FRESH grant — the dead nonce did not
	// carry over — and the fresh grant answers its own card normally.
	f.raiseL2(t, "corr-48-post")
	fresh := waitPending43(t, f.ra, 1)
	if len(fresh) == 0 {
		t.Fatal("the post-resume card never displayed")
	}
	if fresh[0].Grant == "" {
		t.Fatal("the post-resume card carries no grant (re-display must mint a new nonce)")
	}
	if fresh[0].Grant == staleGrant {
		t.Fatal("the post-resume card re-used the dead card's grant — a stale approval would be spendable")
	}
	if err := f.ra.gate.Native().Allow(context.Background(), "corr-48-post", fresh[0].Grant); err != nil {
		t.Fatalf("the fresh grant on the post-resume card: %v", err)
	}
	select {
	case ans := <-f.askRc:
		if ans != tools.AnswerAllow {
			t.Errorf("the post-resume card answered %s, want allow (a fresh approval still works)", ans)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("the post-resume ask never came back after its fresh allow")
	}
	if d := f.ra.gate.Queue().Depth(); d != 0 {
		t.Errorf("queue depth after the post-resume allow = %d, want 0", d)
	}
	t.Logf("48 AC#4 readings: voided=2, notice 「错过 2 个待确认操作」 printed, stale grant refused (ErrUnknownCorrelation), post-resume fresh grant answered allow")
}

// ---------------------------------------------------------------------------
// cell 5 — the panel half: roster join depth sync, head-only, empty migration
// ---------------------------------------------------------------------------

// Test48PanelRosterJoinDepthSyncUnder100ms runs the production join the packet
// pump runs (queue → NativeVerdict.CardView → panel.TaskRosterSectionFrom)
// against a live three-card queue:
//
//  1. per-task blocked readings for the three tasks, each card lighting its
//     OWN task row (the per-task id readings the ticket asks for);
//  2. the measured answer→join sync after each settle, asserted under the
//     100ms bound with the numbers logged (可判事实, no wall-clock guess);
//  3. head-only display depth at depth 3, and the empty-queue migration: no
//     row lit, head face closed, badge input (queue depth) 0.
func Test48PanelRosterJoinDepthSyncUnder100ms(t *testing.T) {
	ui := approval.UIFuncs{
		PromptFn: func(context.Context, approval.Prompt) error { return nil },
		UpdateFn: func(context.Context, approval.Event) error { return nil },
	}
	var logs []string
	g := approval.New(approval.Options{
		UI: ui, Logf: func(f string, a ...any) { logs = append(logs, fmt.Sprintf(f, a...)) },
	})
	const n = 3
	for i := 1; i <= n; i++ {
		g.AdmitTextTask(fmt.Sprintf("task48-p%d", i))
	}
	type parked struct {
		done chan approval48Answer
	}
	asks := make([]parked, n+1)
	// The correlations the loop actually mints are "<taskID>#<callID>"
	// (agent/loop.go's callCorr) - the shape the roster join's third arm is
	// built to recognise. The fixture uses that production shape, not bare ids
	// a real host never sends.
	corrOf := func(i int) string { return fmt.Sprintf("task48-p%d#call-48-%d", i, i) }
	for i := 1; i <= n; i++ {
		i := i
		asks[i].done = make(chan approval48Answer, 1)
		h := observeSpawn48(fmt.Sprintf("task48-p%d", i), corrOf(i), g)
		observe.Default.Spawn("approval-waiter", "test", nil, func(ctx context.Context) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("approval-waiter panicked (D37b): %v", r)
				}
			}()
			a, why := h.gate.PendingApproval(ctx, h.decision)
			asks[i].done <- approval48Answer{a, why}
		})
		waitFor48(t, func() bool { return g.Queue().Depth() == i },
			fmt.Sprintf("queue never reached depth %d", i))
	}

	// (5.1) the join lights each card's OWN task row, per-task readings.
	live := g.Queue().LiveApprovals()
	if len(live) != n {
		t.Fatalf("live rows = %d, want %d", len(live), n)
	}
	for i, it := range live {
		if it.CorrelationID != corrOf(i+1) || it.Position != i+1 {
			t.Errorf("FIFO row %d = %q pos %d, want %q pos %d",
				i, it.CorrelationID, it.Position, corrOf(i+1), i+1)
		}
	}
	state := panel.TaskRosterState{
		Rows: []panel.TaskRow{
			{TaskID: "task48-p1", Label: "任务一", Kind: "root", Status: "Thinking", StatusKnown: true, StreamKey: "task48-p1"},
			{TaskID: "task48-p2", Label: "任务二", Kind: "root", Status: "Thinking", StatusKnown: true, StreamKey: "task48-p2"},
			{TaskID: "task48-p3", Label: "任务三", Kind: "root", Status: "Thinking", StatusKnown: true, StreamKey: "task48-p3"},
		},
		InFlightSlots: n, PoolCap: 4,
	}
	verdict := func(corr string) panel.ApprovalCardView {
		return panel.NativeVerdict{
			CorrelationID: corr, Tool: "fs.write", Level: risk.L2, Reason: "票 48 面板 join 用例",
		}.CardView()
	}
	views := make([]panel.ApprovalCardView, 0, n)
	for _, it := range live {
		views = append(views, verdict(it.CorrelationID))
	}
	sect := panel.TaskRosterSectionFrom(state, views)
	blocked := 0
	for _, row := range sect.Rows {
		if row.BlockedOnApproval {
			blocked++
		}
	}
	if blocked != n {
		t.Errorf("%d rows blocked after the join, want %d (one card per task, all lit)", blocked, n)
	}
	// Head-only display depth at depth 3: the panel card face carries the head
	// with depth 1, not the queue count.
	if h, ok := g.Panel().Head(); !ok || h.CorrelationID != corrOf(1) || h.Depth != 1 {
		t.Errorf("panel head = ok %v %+v, want %q at depth 1 (队头单显)", ok, h, corrOf(1))
	}

	// (5.2) the answer→join sync, as a JUDGEABLE FACT (⛔ a wall-clock guess is
	// not a judging criterion): after each settle the packet join must read
	// the settled state with ZERO extra polls — the first read after the
	// answer returns already shows the row unlit and the head advanced. That
	// is the synchronous bound the <100ms requirement asks for; the wall
	// reading is logged beside it as the number, sub-tick runs print 0s.
	for i := 1; i <= n; i++ {
		corr := corrOf(i)
		start := time.Now()
		if err := g.DecideFromNative(context.Background(), approval.Request{
			CorrelationID: corr, Allow: false, Reason: "票 48 排水 " + corr,
		}); err != nil {
			t.Fatalf("reject %q: %v", corr, err)
		}
		select {
		case a := <-asks[i].done:
			if a.a != tools.AnswerReject {
				t.Errorf("%q answered %s, want reject", corr, a.a)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%q's parked ask never woke", corr)
		}
		// First read after the answer: the row must already be unlit (extra
		// polls would be the measured lag; zero is the synchronous bound).
		polls := 0
		for {
			liveNow := g.Queue().LiveApprovals()
			viewsNow := make([]panel.ApprovalCardView, 0, len(liveNow))
			for _, it := range liveNow {
				viewsNow = append(viewsNow, verdict(it.CorrelationID))
			}
			sectNow := panel.TaskRosterSectionFrom(state, viewsNow)
			lit := false
			for _, row := range sectNow.Rows {
				if row.TaskID == fmt.Sprintf("task48-p%d", i) && row.BlockedOnApproval {
					lit = true
				}
			}
			if !lit {
				break
			}
			polls++
			if polls > 1 {
				t.Fatalf("the panel join lagged the answer for %q by %d extra polls (want 0: the join must read synchronously)", corr, polls)
			}
			time.Sleep(100 * time.Millisecond)
		}
		// The head advanced FIFO-first at the same first read (i>1 already
		// drained; the head is the NEXT item until the last one).
		if i < n {
			if h, ok := g.Panel().Head(); !ok || h.CorrelationID != corrOf(i+1) {
				t.Errorf("after settling %q the head = ok=%v %+v, want the next FIFO item %q", corr, ok, h, corrOf(i+1))
			}
		}
		elapsed := time.Since(start)
		if elapsed > 100*time.Millisecond {
			t.Errorf("answer→panel-join sync for %q took %v, over the 100ms bound", corr, elapsed)
		}
		t.Logf("48 AC#5 answer→panel-join sync %q: extra polls=%d (want 0), wall reading %v (sub-tick prints 0s), bound 100ms", corr, polls, elapsed)
	}
	// input (the queue depth ballCardUI posts) 0.
	if d := g.Queue().Depth(); d != 0 {
		t.Fatalf("queue depth after the drain = %d, want 0", d)
	}
	sectEmpty := panel.TaskRosterSectionFrom(state, nil)
	for _, row := range sectEmpty.Rows {
		if row.BlockedOnApproval {
			t.Errorf("row %q still lit on the empty queue — the empty migration is wrong", row.TaskID)
		}
	}
	if h, ok := g.Panel().Head(); ok {
		t.Errorf("the empty queue still displays a head: %+v", h)
	}
	if live := g.Queue().LiveApprovals(); len(live) != 0 {
		t.Errorf("LiveApprovals on the empty queue = %d rows, want 0", len(live))
	}
	rejected := 0
	for _, l := range logs {
		if strings.Contains(l, "ANSWER-REJECT") {
			rejected++
		}
	}
	if rejected != n {
		t.Errorf("%d ANSWER-REJECT lines booked, want %d (one per drained card)", rejected, n)
	}
}

// joinAsk48 carries what one parked panel-join ask needs (its gate and its own
// decision); the goroutine that runs it goes through observe.Default.Spawn
// under the D38 approval-waiter roster name.
type joinAsk48 struct {
	gate     *approval.Gate
	decision tools.Decision
}

// approval48Answer is one parked ask's outcome (answer, reason), the same pair
// PendingApproval returns.
type approval48Answer struct {
	a   tools.Answer
	why string
}

func observeSpawn48(taskID, corr string, g *approval.Gate) *joinAsk48 {
	return &joinAsk48{
		gate: g,
		decision: tools.Decision{
			TaskID: taskID, Tool: "fs.write", Provider: tools.KindBuiltin,
			Level: risk.L2, Reason: "票 48 面板 join 用例",
			CorrelationID: corr, CallID: "call-" + corr,
		},
	}
}
