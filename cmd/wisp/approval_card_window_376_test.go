package main

// approval_card_window_376_test.go — ticket 376 AC#2's teeth.
//
// WHAT WAS BROKEN. The console L2 card printed 「窗口 0.0s」 while the gate really
// waited out the C18 ceiling. Two sources were in play and only one was read:
// the printer read Prompt.Window, which the L2 route fills with a literal 0
// (internal/agent/approval/gate.go:534), while the wait the gate actually arms is
// q.Timeout() - the same value the L2 route also stamps onto Prompt.Deadline, a
// field that had NO production reader before this ticket. The fix (form 乙′) is
// one hop in the RENDERER: cmd/wisp/run.go's card printer selects the source by
// level (L1 block length, L2 C18 ceiling). The producer and both field meanings
// are deliberately untouched - replies.go:83's 「zero on an L2 card」 stays true.
//
// WHY THESE CASES ARE RED BEFORE THE FIX AND GREEN AFTER. Nothing in the tree
// pinned the card's own 「窗口」 line: the cross-surface card regex at
// panel_pump_test.go:443-444 stops after the rule lines, both existing
// Prompt.Window assertions are L1, both Prompt.Deadline assertions are L2, and
// the L2 Window==0 had no assertion at all. So the L2 case below reddens on the
// pre-fix printer (0.0s against the gate's own 300.0s) and passes after it.
// The L1 case is the counter-example: form 乙 (print Deadline for every card)
// would print 「0.0s」 on an L1 card, because the L1 route passes Deadline=0.
// That case is green before and after - its job is to keep it that way.
//
// THE WAITS HERE ARE FAKE ON PURPOSE. The clock is the approval package's own
// Clock seam: each After registers a channel the test fires by hand, so a 300s
// ceiling costs microseconds and the assertions are about which SOURCE the
// printer read, never about wall time. The real-wait reading (L2 prints 31.0s
// and really waits ~31s; L1's configured 1s clamps to a 2s block and prints
// 2.0s) is measured with a real clock in docs/evidence/s1/376-r1-impl.md.

import (
	"bytes"
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// cardText376 is the console sink the gate prints into, plus the one signal this
// file needs: cardLine is closed once the whole 「窗口」 line has been written.
// Waiting on the printer instead of sleeping is what makes the hand-fired clock
// safe - PendingApproval / PendingWindow arm the clock BEFORE calling Prompt, so
// the line existing proves the After channel exists.
type cardText376 struct {
	mu       sync.Mutex
	buf      bytes.Buffer
	cardLine chan struct{}
}

func newCardText376() *cardText376 {
	return &cardText376{cardLine: make(chan struct{})}
}

func (w *cardText376) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n, err := w.buf.Write(p)
	if strings.Contains(w.buf.String(), "窗口 ") {
		select {
		case <-w.cardLine:
		default:
			close(w.cardLine)
		}
	}
	return n, err
}

func (w *cardText376) text() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

func (w *cardText376) waitCard(t *testing.T) {
	t.Helper()
	select {
	case <-w.cardLine:
	case <-time.After(5 * time.Second):
		t.Fatalf("no card was printed within 5s; output so far:\n%s", w.text())
	}
}

// windowLine376 is the reader-shaped ruler: the number a user actually sees on
// the card, not a variable the test could have typed itself.
var windowLine376 = regexp.MustCompile(`窗口 ([0-9]+\.[0-9])s`)

func printedWindow376(t *testing.T, w *cardText376) (float64, string) {
	t.Helper()
	out := w.text()
	m := windowLine376.FindStringSubmatch(out)
	if m == nil {
		t.Fatalf("the card carries no 「窗口 N.Ns」 line at all:\n%s", out)
	}
	sec, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		t.Fatalf("the card's window number %q is not a number: %v", m[1], err)
	}
	return sec, m[0]
}

// clock376 is approval.Clock with a hand-fired trigger: After records the
// duration and returns a buffered channel the test sends on. Now is never used
// for a timeout decision (approval.Clock:26-29 forbids that).
type clock376 struct {
	mu    sync.Mutex
	armed map[time.Duration]chan time.Time
}

func (c *clock376) Now() time.Time { return time.Unix(0, 0).UTC() }

func (c *clock376) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.armed == nil {
		c.armed = map[time.Duration]chan time.Time{}
	}
	ch := make(chan time.Time, 1)
	c.armed[d] = ch
	return ch
}

func (c *clock376) fire(t *testing.T, d time.Duration) {
	t.Helper()
	c.mu.Lock()
	ch, ok := c.armed[d]
	have := make([]string, 0, len(c.armed))
	for k := range c.armed {
		have = append(have, k.String())
	}
	c.mu.Unlock()
	if !ok {
		t.Fatalf("the gate armed nothing for %v (armed: %v): it is not waiting on the value this case claims", d, have)
	}
	ch <- time.Unix(0, 0).UTC()
}

// card376 hands each Prompt to the real console printer under test and keeps a
// copy, so a case can state the producer's shape next to the printed number.
type card376 struct {
	inner *consoleApprovalUI
	mu    sync.Mutex
	last  approval.Prompt
}

func (c *card376) Prompt(ctx context.Context, p approval.Prompt) error {
	c.mu.Lock()
	c.last = p
	c.mu.Unlock()
	return c.inner.Prompt(ctx, p)
}

func (c *card376) Update(ctx context.Context, e approval.Event) error {
	return c.inner.Update(ctx, e)
}

func (c *card376) prompt(t *testing.T) approval.Prompt {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.last
}

// l2Card376 is the L2 subject the ticket names: overwriting a file that already
// exists. Level is risk.L2, i.e. the verdict a real assessor produces for an
// out-of-allowlist overwrite.
func l2Card376() tools.Decision {
	return tools.Decision{
		Tool:     "fs.write",
		Provider: tools.KindBuiltin,
		Params:   map[string]any{"paths": []any{"C:/work/existing.txt"}},
		Args:     json.RawMessage(`{"paths":["C:/work/existing.txt"]}`),
		Level:    risk.L2,
		RulesHit: []risk.RuleID{risk.R2},
		Reason:   "覆盖已有文件越出 [fs] 授权目录，需 L2 强确认",

		Paths:         []string{"C:/work/existing.txt"},
		CorrelationID: "376-r1-l2",
		TaskID:        "task-376-r1-l2",
		CallID:        "call-376-r1-l2",
	}
}

// l1Card376 is the low-risk notice the ticket's second shot names: one path, so
// Aggregate stays out of the way and the card is the plain L1 block.
func l1Card376() tools.Decision {
	return tools.Decision{
		Tool:     "fs.read",
		Provider: tools.KindBuiltin,
		Params:   map[string]any{"paths": []any{"C:/data/a.txt"}},
		Args:     json.RawMessage(`{"paths":["C:/data/a.txt"]}`),
		Level:    risk.L1,
		RulesHit: []risk.RuleID{risk.R1},
		Reason:   "读被约束目录内的文件，执行前 2 秒可取消",

		Paths:         []string{"C:/data/a.txt"},
		CorrelationID: "376-r1-l1",
		TaskID:        "task-376-r1-l1",
		CallID:        "call-376-r1-l1",
	}
}

// Test376R1L2CardPrintsQueueTimeout is the case that reddens on the pre-fix
// printer. Its expected value is the gate's OWN C18 number (Queue().Timeout()),
// never a literal written here: that is what makes it 「卡上印的 == 真等的」
// instead of 「卡上不再是 0.0」.
func Test376R1L2CardPrintsQueueTimeout(t *testing.T) {
	out := newCardText376()
	ui := &card376{inner: &consoleApprovalUI{out: out}}
	clk := &clock376{}
	g := approval.New(approval.Options{UI: ui, Clock: clk})
	revoke := g.AdmitTextTask("task-376-r1-l2")
	defer revoke()

	done := make(chan struct{})
	var got tools.Answer
	var why string
	go func() {
		defer close(done)
		got, why = g.PendingApproval(context.Background(), l2Card376())
	}()

	out.waitCard(t)
	p := ui.prompt(t)
	if p.Window != 0 {
		t.Fatalf("premise broke: the L2 card's Window is %v, want 0 - the L2 route passes a literal 0 "+
			"(approval/gate.go:534), and ticket 376 chose form 乙′ precisely so that producer stays put. "+
			"If that changed, replies.go:83's 「zero on an L2 card」 has to change with it", p.Window)
	}
	timeout := g.Queue().Timeout()
	if p.Deadline != timeout {
		t.Fatalf("premise broke: the L2 card's Deadline is %v but Queue().Timeout() is %v - the card's "+
			"only honest source is the gate's own C18 value", p.Deadline, timeout)
	}
	clk.fire(t, timeout)
	<-done

	printed, line := printedWindow376(t, out)
	if want := timeout.Seconds(); printed != want {
		t.Errorf("the L2 card printed %q, but the gate waits Queue().Timeout() = %.1fs (%v).\n"+
			"Ticket 376: the card's seconds must EQUAL the queue's timeout, so a user deciding whether "+
			"to wait is reading the same number the process is waiting on.\nfull output:\n%s",
			line, want, timeout, out.text())
	}
	if got != tools.AnswerTimeout {
		t.Errorf("an unanswered L2 card answered %q (%s), want timeout: the polarity is the bridge's, "+
			"and C18's timeout->reject is what makes 300s safe to print", got, why)
	}
	if n := g.Queue().Depth(); n != 0 {
		t.Errorf("%d approvals still pending after the deadline fired; the queue leaked a card", n)
	}
}

// Test376R1L1CardKeepsItsBlockLength is the counter-example that form 乙 (print
// Deadline for every card) would break: an L1 Prompt carries Deadline=0, so a
// blind source swap prints 「窗口 0.0s」 on the L1 card. The expected value is the
// gate's own clamped block length (Gate.Window()), and the configured 1s is
// below the contract floor on purpose - that is the value a real
// l1_window_sec=1 config arrives with.
func Test376R1L1CardKeepsItsBlockLength(t *testing.T) {
	out := newCardText376()
	ui := &card376{inner: &consoleApprovalUI{out: out}}
	clk := &clock376{}
	g := approval.New(approval.Options{UI: ui, Clock: clk, Window: 1 * time.Second})
	revoke := g.AdmitTextTask("task-376-r1-l1")
	defer revoke()

	block := g.Window()
	if block != approval.MinL1Window {
		t.Fatalf("premise broke: a configured 1s window clamped to %v, want the %v floor "+
			"(approval/gate.go:143-151)", block, approval.MinL1Window)
	}

	done := make(chan struct{})
	var got tools.Answer
	var why string
	go func() {
		defer close(done)
		got, why = g.PendingWindow(context.Background(), l1Card376())
	}()

	out.waitCard(t)
	p := ui.prompt(t)
	if p.Level != "L1" {
		t.Fatalf("the card's level is %q, want L1: this case exists to hold the L1 half of the split", p.Level)
	}
	if p.Window != block {
		t.Fatalf("premise broke: the L1 card's Window is %v, want the gate's block length %v", p.Window, block)
	}
	if p.Deadline != 0 {
		t.Fatalf("premise broke: the L1 card's Deadline is %v, want 0 - this is exactly why printing "+
			"Deadline for every level (form 乙) would put 「0.0s」 on an L1 card", p.Deadline)
	}
	clk.fire(t, block)
	<-done

	printed, line := printedWindow376(t, out)
	if want := block.Seconds(); printed != want {
		t.Errorf("the L1 card printed %q, but its real block is Gate.Window() = %.1fs (%v).\n"+
			"Ticket 376 AC#2 forbids trading the L2 defect for an L1 one.\nfull output:\n%s",
			line, want, block, out.text())
	}
	if printed == 0 {
		t.Errorf("the L1 card printed a zero-second window: an L1 block nobody can read is not a window")
	}
	if got != tools.AnswerTimeout {
		t.Errorf("an unopposed L1 window answered %q (%s), want timeout: for L1 that MEANS execute (B1)", got, why)
	}
	if n := g.Queue().Depth(); n != 0 {
		t.Errorf("%d approvals still pending after the L1 block ended; the L1 route must not queue", n)
	}
}
