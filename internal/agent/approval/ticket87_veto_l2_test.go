package approval_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// 票 87: 卡片已经显示，否决却查无此项 ⇒ 人想现在拒也拒不掉
//
// Ticket 84 measured that the wait has a bound (the C18 300s auto-reject), so
// this is NOT a deadlock and NOT a hole on the allow side: the worst case today
// still refuses. What is broken is the human's control - the card advertises
// the veto channels (promptFor hands BOTH routes g.channels.Statuses()), and
// Gate.Veto looked only in the L1 window table and the started-call bookkeeping,
// so pressing one against an on-screen L2 card answered
// ErrUnknownCorrelation and the call sat out the whole deadline.
//
// The two halves below are both decidable without waiting 300s:
//
//	(i) the wait must END NOW - asserted against a 30s armed deadline, so a
//	    regression back to "wait it out" costs the mutation 30s of wall clock
//	    and shows up as an elapsed-time failure, not as a slow-but-green run;
//	(ii) it must still be a REFUSAL - asserted by the answer value, by the
//	    refusal wording, and by the fact that the item is gone from the queue
//	    so no grant can be spent on it afterwards.
//
// The default 300s / 30s numbers are untouched by this file, and
// TestDefaultDeadlineWallClockMeasurement (ticket 84) still owns the wall-clock
// reading of the full bound.
// ---------------------------------------------------------------------------

// vetoBudget is the test-side bound for "ended now". The armed deadline in
// these cases is 30s, so 2s is ~15x the work a settled reply needs and 1/15 of
// the wait it must NOT perform.
const vetoBudget = 2 * time.Second

// ticket87Gate builds a gate whose L2 deadline is long enough that "the veto
// worked" and "the veto was ignored and the deadline fired" cannot be confused.
func ticket87Gate(t *testing.T, ui approval.UI) *approval.Gate {
	t.Helper()
	g, _, _ := newGate(t, ui, approval.Options{
		Clock:           approval.SystemClock{},
		ApprovalTimeout: 30 * time.Second,
		WarningLead:     25 * time.Second,
	})
	g.AdmitTextTask(testTask)
	return g
}

// TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses is AC#2 on the
// exact-key shape: the veto carries the correlation id printed on the card, the
// card is displayed, and the gate answers 「no such pending item」.
func TestAVetoAgainstAnOnScreenL2CardEndsTheWaitNowAndStillRefuses(t *testing.T) {
	ui := newFakeUI()
	g := ticket87Gate(t, ui)

	d := l2Decision("C:/elsewhere/87-veto-on-screen.txt")
	d.CorrelationID = "corr-87-veto"
	res := runApproval(t, g, context.Background(), d)

	// Premise, read from the surfaces rather than assumed: the card IS on
	// screen and the gate DOES hold it pending. A red here means the shape this
	// ticket is about stopped being reachable, which is not a pass.
	p := ui.wait(t)
	if p.CorrelationID != "corr-87-veto" {
		t.Fatalf("card correlation_id = %q, want the id under test", p.CorrelationID)
	}
	if p.Level != "L2" {
		t.Fatalf("card level = %q, want L2 (this is the L2 route)", p.Level)
	}
	if n := g.Queue().Depth(); n != 1 {
		t.Fatalf("queue depth = %d, want the 1 pending item the card shows", n)
	}
	// The card's own promise: it lists the channel this test presses.
	if !channelAdvertised(p, approval.ChannelBall) {
		t.Fatalf("card does not advertise the ball veto channel: %+v", p.Channels)
	}

	start := time.Now()
	if err := g.Veto(approval.Veto{CorrelationID: p.CorrelationID, Channel: approval.ChannelBall}); err != nil {
		t.Fatalf("Veto on the displayed card: %v（票 87 的原始症状就是这个调用返回查无此项）", err)
	}

	got := mustAnswerWithin(t, res, vetoBudget, "veto against an on-screen L2 card")
	el := time.Since(start)

	// (ii) safety floor: the answer is a refusal, never an allow.
	if got.a != tools.AnswerReject {
		t.Errorf("answer=%v why=%q, want %v：否决只能把这次调用推向拒绝", got.a, got.why, tools.AnswerReject)
	}
	if got.a == tools.AnswerAllow {
		t.Errorf("SECURITY: a veto produced %v", tools.AnswerAllow)
	}
	if !strings.Contains(got.why, "拒") && !strings.Contains(got.why, "否决") {
		t.Errorf("why=%q，答复必须说明这是用户的否决", got.why)
	}
	// (i) the wait ended now, not at the armed 30s deadline.
	if el > vetoBudget {
		t.Errorf("ended after %v, want immediately：等待没有因这次否决而结束（ armed deadline %v）",
			el, g.Queue().Timeout())
	}
	if n := g.Queue().Depth(); n != 0 {
		t.Errorf("queue still holds %d items after the veto", n)
	}
	if len(ui.ofKind(approval.EventDismissed)) == 0 {
		t.Errorf("no dismissed event: the card was never taken off screen")
	}
	// A veto is not a grant: nothing may be answered on that id afterwards.
	if err := g.Native().Allow(context.Background(), p.CorrelationID, p.Grant); err == nil {
		t.Errorf("Allow on a vetoed correlation_id succeeded: 否决被当成了批准")
	} else if !errors.Is(err, approval.ErrNotPending) && !errors.Is(err, approval.ErrUnknownCorrelation) &&
		!errors.Is(err, approval.ErrBadGrant) {
		t.Errorf("Allow err=%v, want a routing error", err)
	}
}

// TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard is AC#2 on instance 3: the
// D31 / bridge bookkeeping key is the INCOMING correlation id or the task id
// (gate.go:451, tools/bridge.go:397), which is not the string the queue issues
// when it has re-stamped one. A card addressed by the host's own key must be
// refusable too - and only refusable.
func TestAVetoNamedByTheHostsOwnKeyStillRefusesThatCard(t *testing.T) {
	ui := newFakeUI()
	g := ticket87Gate(t, ui)

	d := l2Decision("C:/elsewhere/87-alias-key.txt")
	d.CorrelationID = "corr-87-incoming"
	d.TaskID = "task-87-host-key"
	g.AdmitTextTask("task-87-host-key")
	res := runApproval(t, g, context.Background(), d)
	p := ui.wait(t)
	if p.CorrelationID == "task-87-host-key" {
		t.Fatal("premise broke: the queue issued the task id as the card key, no alias to test")
	}

	start := time.Now()
	if err := g.Veto(approval.Veto{CorrelationID: "task-87-host-key", Channel: approval.ChannelEsc}); err != nil {
		t.Fatalf("Veto by the host's own key: %v", err)
	}
	got := mustAnswerWithin(t, res, vetoBudget, "veto addressed by the task id")
	if got.a != tools.AnswerReject {
		t.Errorf("answer=%v why=%q, want %v", got.a, got.why, tools.AnswerReject)
	}
	if el := time.Since(start); el > vetoBudget {
		t.Errorf("ended after %v, want immediately", el)
	}
}

// TestAnAmbiguousVetoNameIsNotGuessedOnTheRefusalSide is the floor under the
// lenient lookup: two live cards that a name cannot tell apart get NEITHER.
// Guessing here would let one keystroke refuse a card the user never looked at.
func TestAnAmbiguousVetoNameIsNotGuessedOnTheRefusalSide(t *testing.T) {
	ui := newFakeUI()
	g := ticket87Gate(t, ui)

	const sharedTask = "task-87-shared"
	g.AdmitTextTask(sharedTask)
	var cards []approval.Prompt
	for _, corr := range []string{"corr-87-twin-a", "corr-87-twin-b"} {
		d := l2Decision("C:/elsewhere/" + corr + ".txt")
		d.CorrelationID = corr
		d.TaskID = sharedTask
		runApproval(t, g, context.Background(), d)
		cards = append(cards, ui.wait(t))
	}
	if len(cards) != 2 || cards[0].CorrelationID == cards[1].CorrelationID {
		t.Fatalf("two distinct cards were not displayed: %+v", cards)
	}
	if n := g.Queue().Depth(); n != 2 {
		t.Fatalf("queue depth = %d, want 2 live cards", n)
	}

	err := g.Veto(approval.Veto{CorrelationID: sharedTask, Channel: approval.ChannelBall})
	if !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("ambiguous veto err=%v, want %v（两名同时有效时不许猜一个）",
			err, approval.ErrUnknownCorrelation)
	}
	if n := g.Queue().Depth(); n != 2 {
		t.Fatalf("an ambiguous veto settled %d card(s); nothing may be answered by it", 2-n)
	}

	// Both cards are still answerable by their own exact key, which is also the
	// proof that the ambiguous veto did not quietly burn one of them.
	for i, p := range cards {
		if err := g.Native().Reject(p.CorrelationID, "清理"); err != nil {
			t.Fatalf("Native().Reject(card %d): %v", i, err)
		}
	}
	waitFor(t, func() bool { return g.Queue().Depth() == 0 }, "cards never drained")
}

// TestAVetoOnAnUnknownNameLeavesTheCardAllowableByItsOwnGrant is the negative
// side of the same rule: the fallback may reach exactly the item the name
// addresses, and no other. A veto for somebody else's id must not touch the
// pending card, must not consume its grant, and must not answer it.
func TestAVetoOnAnUnknownNameLeavesTheCardAllowableByItsOwnGrant(t *testing.T) {
	ui := newFakeUI()
	g := ticket87Gate(t, ui)

	d := l2Decision("C:/elsewhere/87-untouchable.txt")
	d.CorrelationID = "corr-87-untouched"
	res := runApproval(t, g, context.Background(), d)
	p := ui.wait(t)

	err := g.Veto(approval.Veto{CorrelationID: "corr-87-someone-elses-card", Channel: approval.ChannelBall})
	if !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Errorf("err=%v, want %v", err, approval.ErrUnknownCorrelation)
	}
	if n := g.Queue().Depth(); n != 1 {
		t.Fatalf("queue depth = %d, want the card still pending", n)
	}
	select {
	case a := <-res:
		t.Fatalf("the card was answered by an unrelated veto: %+v", a)
	default:
	}
	// Still fully alive on its own route, both directions.
	if err := g.Native().Allow(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("Allow after an unrelated veto: %v", err)
	}
	got := mustAnswer(t, res)
	if got.a != tools.AnswerAllow {
		t.Errorf("answer=%v, want %v", got.a, tools.AnswerAllow)
	}
	g.Complete(p.CorrelationID)
}

// channelAdvertised reports whether one card renders one veto channel.
func channelAdvertised(p approval.Prompt, ch approval.Channel) bool {
	for _, c := range p.Channels {
		if c.Channel == ch {
			return true
		}
	}
	return false
}
