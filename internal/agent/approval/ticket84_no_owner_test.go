package approval_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// 票 84: 「等一个永远不会来的决定」——无对应待审批项时到底等多久
//
// The claim under test is NOT "PendingApproval looks up an existing pending
// item and parks when it finds none" (it registers its own item via q.push, so
// there is always exactly one owner waiting on it). The claim is that the wait
// has an UPPER BOUND, and that every route which really can hit "no pending
// item with that id" fails fast with a named error instead of parking.
//
// These are the wall-clock readings for that, on whatever OS runs them; the
// two-sided numbers live in docs/evidence/s1/84-ac1-bounded-wait.md.
// ---------------------------------------------------------------------------

// mustAnswerWithin waits for one gate answer and fails with a TIMESTAMPED
// reading, so a red run says how long the gate actually held the call rather
// than just "it did not answer".
func mustAnswerWithin(t *testing.T, ch <-chan answer, budget time.Duration, what string) answer {
	t.Helper()
	start := time.Now()
	select {
	case a := <-ch:
		t.Logf("%s: returned after %v (start=%s, end=%s)",
			what, time.Since(start), start.Format(time.TimeOnly), time.Now().Format(time.TimeOnly))
		return a
	case <-time.After(budget):
		t.Fatalf("%s: STILL BLOCKED after %v (start=%s, now=%s) - 闸门没有上界",
			what, budget, start.Format(time.TimeOnly), time.Now().Format(time.TimeOnly))
		return answer{}
	}
}

// TestAnUnansweredL2CardResolvesToRejectAtTheQueueDeadline is the boundedness
// reading at the shape the production route has: nobody answers, so the C18
// deadline must resolve it to REJECT.
//
// The deadline length here comes from Options.ApprovalTimeout, which is the
// same knob cmd/wisp/run.go:259 feeds from [risk].confirm_timeout_sec - it is
// a supported production configuration face, not a test-only back window. The
// default (300s) itself is asserted in
// TestDefaultQueueDeadlineIsFiniteAtThreeHundredSeconds, and measured on the
// real wall clock in TestDefaultDeadlineWallClockMeasurement.
func TestAnUnansweredL2CardResolvesToRejectAtTheQueueDeadline(t *testing.T) {
	ui := newFakeUI()
	g, _, _ := newGate(t, ui, approval.Options{
		Clock:           approval.SystemClock{},
		ApprovalTimeout: 250 * time.Millisecond,
		WarningLead:     100 * time.Millisecond,
	})
	g.AdmitTextTask(testTask)

	// An EMPTY gate: nothing was queued before this call, and the incoming
	// correlation id matches no live item - the exact shape the 75 mutation
	// evidence described as 「correlation_id 无对应待审批项」.
	if n := g.Queue().Depth(); n != 0 {
		t.Fatalf("前置条件破了：门里有 %d 个待审批项，应为空门", n)
	}
	d := l2Decision("C:/elsewhere/never-approved-by-anyone.txt")
	d.CorrelationID = "corr-with-no-pending-item"

	start := time.Now()
	res := runApproval(t, g, context.Background(), d)
	got := mustAnswerWithin(t, res, 5*time.Second, "unanswered L2 (250ms deadline)")

	if got.a != tools.AnswerTimeout {
		t.Errorf("answer=%v, want %v：无人答复的 L2 必须按超时判拒绝", got.a, tools.AnswerTimeout)
	}
	if !strings.Contains(got.why, "超时") {
		t.Errorf("why=%q，拒绝理由必须说明是超时自动拒绝", got.why)
	}
	if el := time.Since(start); el < 250*time.Millisecond {
		t.Errorf("returned in %v, under the armed 250ms deadline：上界没有真正生效", el)
	}
	if el := time.Since(start); el > 2*time.Second {
		t.Errorf("returned in %v: the 250ms deadline fired far too late, 上界形同虚设", el)
	}
	if p := ui.all(); len(p) != 1 {
		t.Errorf("displayed %d cards, want 1", len(p))
	}
}

// TestDefaultQueueDeadlineIsFiniteAtThreeHundredSeconds pins the answer to
// "does the DEFAULT configuration have an upper bound" without sleeping, and
// pins the two shapes that could silently mean "forever".
func TestDefaultQueueDeadlineIsFiniteAtThreeHundredSeconds(t *testing.T) {
	for _, tc := range []struct {
		name    string
		timeout time.Duration
	}{
		{"unset (zero)", 0},
		{"negative", -time.Second},
		{"explicit default", approval.DefaultApprovalTimeout},
	} {
		g := approval.New(approval.Options{ApprovalTimeout: tc.timeout})
		if got := g.Queue().Timeout(); got != approval.DefaultApprovalTimeout {
			t.Errorf("%s: Queue().Timeout()=%v, want the finite %v default",
				tc.name, got, approval.DefaultApprovalTimeout)
		}
	}
	if approval.DefaultApprovalTimeout <= 0 {
		t.Fatal("DefaultApprovalTimeout must be a finite positive bound")
	}
}

// TestEveryAnswerRouteForAnUnknownCorrelationFailsFast is the other half of
// AC#1: the routes that genuinely face "no pending item carries that id" must
// return a NAMED error immediately, because that is where a human's click
// lands when the queue has already moved on.
func TestEveryAnswerRouteForAnUnknownCorrelationFailsFast(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name string
		call func(t *testing.T, g *approval.Gate) error
	}{
		{"Native.Allow", func(_ *testing.T, g *approval.Gate) error {
			return g.Native().Allow(ctx, "corr-never-queued", "grant_00")
		}},
		{"Native.Reject", func(_ *testing.T, g *approval.Gate) error {
			return g.Native().Reject("corr-never-queued", "拒绝一个不存在的东西")
		}},
		{"Panel.Reject", func(_ *testing.T, g *approval.Gate) error {
			return g.Panel().Reject("corr-never-queued", "面板拒绝")
		}},
		{"DecideFromNative.reject", func(_ *testing.T, g *approval.Gate) error {
			return g.DecideFromNative(ctx, approval.Request{CorrelationID: "corr-never-queued", Reason: "r"})
		}},
		{"DecideFromNative.allow", func(_ *testing.T, g *approval.Gate) error {
			return g.DecideFromNative(ctx, approval.Request{CorrelationID: "corr-never-queued", Allow: true, Grant: "grant_00"})
		}},
		{"DecideFromPanel.reject", func(_ *testing.T, g *approval.Gate) error {
			return g.DecideFromPanel(ctx, approval.Request{CorrelationID: "corr-never-queued", Reason: "r"})
		}},
		{"Veto.unknown", func(_ *testing.T, g *approval.Gate) error {
			return g.Veto(approval.Veto{CorrelationID: "corr-never-queued", Channel: approval.ChannelBall})
		}},
		{"Veto.empty", func(_ *testing.T, g *approval.Gate) error {
			return g.Veto(approval.Veto{Channel: approval.ChannelEsc})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g, _, _ := newGate(t, newFakeUI(), approval.Options{Clock: approval.SystemClock{}})
			start := time.Now()
			err := tc.call(t, g)
			el := time.Since(start)
			if !errors.Is(err, approval.ErrUnknownCorrelation) {
				t.Fatalf("err=%v, want ErrUnknownCorrelation（具名错误，不 panic、不静默重试）", err)
			}
			// 250ms is 20x what an in-memory map lookup takes on either host;
			// a red here means a reply route grew a wait, which is the shape
			// this ticket was opened about.
			if el > 250*time.Millisecond {
				t.Errorf("%s took %v (start=%s, end=%s)：无对应待审批项必须立即失败",
					tc.name, el, start.Format(time.TimeOnly), time.Now().Format(time.TimeOnly))
			}
			t.Logf("%s: %v in %v", tc.name, err, el)
		})
	}
}

// TestAnUnreachablePromptSurfaceIsRefusedImmediatelyNotWaitedOut is the
// conservative-direction guarantee: when the gate can tell up front that
// nobody will ever answer, it must refuse now instead of parking until the
// deadline. This is also AC#4(ii)'s anchor test - turning this refusal into a
// silent pass makes it red.
func TestAnUnreachablePromptSurfaceIsRefusedImmediatelyNotWaitedOut(t *testing.T) {
	ui := newFakeUI()
	ui.setErr(errors.New("native card window gone"))
	g, _, _ := newGate(t, ui, approval.Options{
		Clock:           approval.SystemClock{},
		ApprovalTimeout: 30 * time.Second,
	})
	g.AdmitTextTask(testTask)

	start := time.Now()
	res := runApproval(t, g, context.Background(), l2Decision("C:/elsewhere/unreachable-card.txt"))
	got := mustAnswerWithin(t, res, 2*time.Second, "prompt surface unreachable")

	if got.a != tools.AnswerReject {
		t.Errorf("answer=%v why=%q, want REJECT：界面不可达时只能拒绝", got.a, got.why)
	}
	if !strings.Contains(got.why, "fail-closed") {
		t.Errorf("why=%q，拒绝必须说明是 fail-closed", got.why)
	}
	if el := time.Since(start); el > time.Second {
		t.Errorf("refused only after %v, want an immediate refusal", el)
	}
	if n := g.Queue().Depth(); n != 0 {
		t.Errorf("abandoned item left %d pending entries behind", n)
	}
}

// TestDefaultDeadlineWallClockMeasurement is AC#1's literal reading: the
// DEFAULT gate (300s, real clock, nothing answering) measured on the wall
// clock, with timestamps, on whichever host runs it. It exists so "有上界"
// is a measured statement rather than an inference from the code.
//
// Skipped unless WISP_84_MEASURE=1, because it costs 300s of wall clock. The
// evidence runs are recorded in docs/evidence/s1/84-ac1-bounded-wait.md
// (Windows host + docker golang:1.27 on a git-archive snapshot, outside the
// repo), each with the real exit code.
func TestDefaultDeadlineWallClockMeasurement(t *testing.T) {
	if os.Getenv("WISP_84_MEASURE") == "" {
		t.Skip("有意慢：300s 墙钟计量，只在 WISP_84_MEASURE=1 时跑（票 84 AC#1；见 docs/evidence/s1/84-ac1-bounded-wait.md）")
	}
	g := approval.New(approval.Options{})
	g.AdmitTextTask(testTask)
	if got := g.Queue().Timeout(); got != approval.DefaultApprovalTimeout {
		t.Fatalf("measurement premise broke: Timeout()=%v, want %v", got, approval.DefaultApprovalTimeout)
	}

	d := l2Decision("C:/elsewhere/measured-on-the-wall-clock.txt")
	d.CorrelationID = "corr-84-measure"
	start := time.Now()
	done := make(chan answer, 1)
	go func() {
		a, why := g.PendingApproval(context.Background(), d)
		done <- answer{a, why}
	}()

	const observationCap = 400 * time.Second
	var got answer
	select {
	case got = <-done:
	case <-time.After(observationCap):
		t.Fatalf("(c) UNBOUNDED: no return within %v (start=%s, now=%s) - 等待没有上界",
			observationCap, start.Format(time.RFC3339), time.Now().Format(time.RFC3339))
	}
	el := time.Since(start)
	t.Logf("(b) BOUNDED: returned after %v at %s (start=%s), answer=%v why=%q, configured deadline=%v",
		el, time.Now().Format(time.RFC3339), start.Format(time.RFC3339), got.a, got.why, g.Queue().Timeout())
	if got.a != tools.AnswerTimeout {
		t.Errorf("answer=%v, want %v（300s 未确认一律判拒绝，C18）", got.a, tools.AnswerTimeout)
	}
	if el < approval.DefaultApprovalTimeout {
		t.Errorf("returned in %v, before the %v deadline", el, approval.DefaultApprovalTimeout)
	}
	if el > 2*approval.DefaultApprovalTimeout {
		t.Errorf("returned in %v, more than 2x the %v bound", el, approval.DefaultApprovalTimeout)
	}
}
