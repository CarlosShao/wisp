package approval_test

// Ticket 375's rulers for the approval layer's half of the fix.
//
// WHAT EACH ONE IS FOR. The seam 375 adds (approval.AnswerSource) does not exist
// before the change, so these cases cannot compile against the old tree: the red
// that proves they are load-bearing is taken at the ASSEMBLY level instead
// (cmd/wisp/approval_unanswerable_375_test.go's Test375UnattendedConsoleRefuses-
// TheL1Write and Test375ResidentGateRefusesTheWindowWithNoBall, and the
// 375-r1 evidence file's READ-F / READ-S readings, all of which use no new API
// and go red on the old tree by running). What these cases add is the sharp,
// fast end of the same three claims: the L1 guard, the L2 boundary that must NOT
// move, and the mode read that has to agree with the guard.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// logSink375 captures the gate's audit lines (the gate logs from its own
// goroutine, so every read takes the lock).
type logSink375 struct {
	mu    sync.Mutex
	lines []string
}

func (l *logSink375) f(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, fmt.Sprintf(format, args...))
}

func (l *logSink375) count(needle string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := 0
	for _, line := range l.lines {
		if strings.Contains(line, needle) {
			n++
		}
	}
	return n
}

// gate375 builds a gate directly (not through fakes_test.go's newGate) so this
// file can pass Options.Answer without changing a shared helper.
func gate375(t *testing.T, ui approval.UI, clk approval.Clock, answer approval.AnswerSource, sink *logSink375) *approval.Gate {
	t.Helper()
	if clk == nil {
		clk = newFakeClock()
	}
	g := approval.New(approval.Options{
		UI:       ui,
		Clock:    clk,
		Channels: approval.NewChannels(),
		Answer:   answer,
		Logf:     sink.f,
	})
	t.Cleanup(func() {
		if n := g.Queue().Depth(); n != 0 {
			t.Errorf("gate exited with %d pending approval items", n)
		}
	})
	return g
}

func w375Decision(level risk.Level) tools.Decision {
	return tools.Decision{
		Tool: "fs.write", Provider: tools.KindBuiltin,
		Params:   map[string]any{"path": "/work/w375.txt"},
		Args:     json.RawMessage(`{"path":"/work/w375.txt"}`),
		Level:    level,
		RulesHit: []risk.RuleID{risk.R1},
		Reason:   "R1: 工具声明为下界（L1）",
		Paths:    []string{"/work/w375.txt"},
		TaskID:   testTask,
	}
}

// Test375L1RefusesAWindowNobodyCouldSee is the core claim: no answer source
// means no window, and NO window means the call is refused - not "waited and
// then executed", which is what the L1 route did before this ticket.
func Test375L1RefusesAWindowNobodyCouldSee(t *testing.T) {
	ui := newFakeUI()
	sink := &logSink375{}
	g := gate375(t, ui, nil, func() bool { return false }, sink)
	revoke := g.AdmitTextTask(testTask)
	defer revoke()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := mustAnswer(t, runWindow(t, g, ctx, w375Decision(risk.L1)))

	if a.a != tools.AnswerReject {
		t.Fatalf("answer = %q, want %q: a window nobody could see must be refused, and its expiry "+
			"must never be read as consent (375 READ-S)", a.a, tools.AnswerReject)
	}
	if n := ui.count(); n != 0 {
		t.Errorf("the surface was asked to show %d card(s); a card nobody can answer is a lie on screen", n)
	}
	if !strings.Contains(a.why, "无人可答") {
		t.Errorf("the refusal does not name the cause: %q", a.why)
	}
	if sink.count("L1-UNANSWERABLE") != 1 {
		t.Errorf("L1-UNANSWERABLE lines = %d, want exactly 1: %v", sink.count("L1-UNANSWERABLE"), sink.lines)
	}
	if sink.count("ANSWER-EXPIRED") != 0 {
		t.Errorf("ANSWER-EXPIRED was booked for a window that must never have opened: %v", sink.lines)
	}
}

// Test375L1AnswerableKeepsTimeoutMeansExecute is the positive control at this
// layer: with an answer source declared, the L1 route is byte-for-byte the
// pre-375 route (card shown, expiry -> execute).
func Test375L1AnswerableKeepsTimeoutMeansExecute(t *testing.T) {
	ui := newFakeUI()
	clk := newFakeClock()
	sink := &logSink375{}
	g := gate375(t, ui, clk, func() bool { return true }, sink)
	revoke := g.AdmitTextTask(testTask)
	defer revoke()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := runWindow(t, g, ctx, w375Decision(risk.L1))
	p := ui.wait(t)
	if p.Level != "L1" {
		t.Fatalf("card level = %q, want L1", p.Level)
	}
	clk.Advance(g.Window() + time.Millisecond)
	a := mustAnswer(t, ch)
	if a.a != tools.AnswerTimeout {
		t.Fatalf("answer = %q, want %q (an unopposed window still MEANS EXECUTE)", a.a, tools.AnswerTimeout)
	}
	if sink.count("ANSWER-EXPIRED") != 1 || sink.count("L1-UNANSWERABLE") != 0 {
		t.Errorf("the routed window's audit lines moved: %v", sink.lines)
	}
}

// Test375UndeclaredHostIsUntouched pins nil != false: a host that declared
// nothing (every host and fixture built before this ticket) keeps the old
// behaviour. This is the case that stops the fix from becoming a purge.
func Test375UndeclaredHostIsUntouched(t *testing.T) {
	ui := newFakeUI()
	clk := newFakeClock()
	sink := &logSink375{}
	g := gate375(t, ui, clk, nil, sink)
	revoke := g.AdmitTextTask(testTask)
	defer revoke()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := runWindow(t, g, ctx, w375Decision(risk.L1))
	ui.wait(t)
	clk.Advance(g.Window() + time.Millisecond)
	if a := mustAnswer(t, ch); a.a != tools.AnswerTimeout {
		t.Fatalf("answer = %q, want %q: an undeclared host is not a host that said no", a.a, tools.AnswerTimeout)
	}
}

// Test375L2IgnoresTheAnswerSource is this ticket's boundary 2: the L2 route's
// polarity is 到点即拒 and is already correct, so the guard must not touch it. A
// no-answer-source L2 card is still raised and still rides its clock; only the
// L1 route gained a reason to refuse early.
func Test375L2IgnoresTheAnswerSource(t *testing.T) {
	ui := newFakeUI()
	clk := newFakeClock()
	sink := &logSink375{}
	g := gate375(t, ui, clk, func() bool { return false }, sink)
	revoke := g.AdmitTextTask(testTask)
	defer revoke()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := runApproval(t, g, ctx, w375Decision(risk.L2))
	if p := ui.wait(t); p.Level != "L2" {
		t.Fatalf("card level = %q, want L2 (the L2 route still raises its card)", p.Level)
	}
	clk.Advance(g.Queue().Timeout() + time.Millisecond)
	a := mustAnswer(t, ch)
	// The gate's L2 route hands the raw expiry back (AnswerTimeout) and the
	// bridge's L2 branch is what turns it into "decision=timeout->reject" - that
	// split predates this ticket and is exactly what "到点即拒" is made of, so
	// this case pins the value the route produces rather than a rewritten one.
	if a.a != tools.AnswerTimeout {
		t.Fatalf("answer = %q, want %q (the L2 route's own expiry answer)", a.a, tools.AnswerTimeout)
	}
	if !strings.Contains(a.why, "拒绝") {
		t.Errorf("the L2 expiry sentence no longer says the call is refused: %q", a.why)
	}
	if strings.Contains(a.why, "无人可答") {
		t.Errorf("the L2 refusal was rewritten by ticket 375's L1 guard: %q", a.why)
	}
	if sink.count("L1-UNANSWERABLE") != 0 {
		t.Errorf("the L1 guard fired on the L2 route: %v", sink.lines)
	}
}

// modeSrc375 is a tools.ModeSource that answers one fixed mode.
type modeSrc375 struct{ m risk.Mode }

func (s modeSrc375) PermissionMode() risk.Mode { return s.m }

// Test375UnsilencedModeOnlyDropsAModeNobodyCouldUse is the other half: the mode
// the bridge reads may not silence a question this process cannot ask. It must
// move ONLY in that case, and it must book one line rather than one per call.
func Test375UnsilencedModeOnlyDropsAModeNobodyCouldUse(t *testing.T) {
	cases := []struct {
		name       string
		answer     approval.AnswerSource
		configured risk.Mode
		want       risk.Mode
		wantLines  int
	}{
		{"undeclared host keeps auto_approve", nil, risk.ModeAutoApprove, risk.ModeAutoApprove, 0},
		{"answerable host keeps auto_approve", func() bool { return true }, risk.ModeAutoApprove, risk.ModeAutoApprove, 0},
		{"unanswerable host drops auto_approve", func() bool { return false }, risk.ModeAutoApprove, risk.ModeAskEveryStep, 1},
		{"unanswerable host drops ask_high_risk", func() bool { return false }, risk.ModeAskHighRisk, risk.ModeAskEveryStep, 1},
		{"unanswerable host keeps the strictest", func() bool { return false }, risk.ModeAskEveryStep, risk.ModeAskEveryStep, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sink := &logSink375{}
			src := approval.UnsilencedMode(modeSrc375{tc.configured}, tc.answer, sink.f)
			if got := src.PermissionMode(); got != tc.want {
				t.Fatalf("mode = %s, want %s", got, tc.want)
			}
			// Twice: the read happens once per tool call, and the line must not
			// be per call.
			if got := src.PermissionMode(); got != tc.want {
				t.Fatalf("second read = %s, want %s", got, tc.want)
			}
			if n := sink.count("MODE-UNANSWERABLE"); n != tc.wantLines {
				t.Fatalf("MODE-UNANSWERABLE lines = %d, want %d: %v", n, tc.wantLines, sink.lines)
			}
			if tc.wantLines == 1 && !strings.Contains(sink.lines[0], "configured=auto_approve") &&
				!strings.Contains(sink.lines[0], "configured=ask_high_risk") {
				t.Errorf("the downgrade line does not name the mode it took away: %v", sink.lines)
			}
		})
	}
}

// Test375UnsilencedModeReArmsAfterItClears: once the host can answer again, a
// later downgrade is a NEW event and is booked again - otherwise the ledger
// would show one downgrade for an episode that happened twice.
func Test375UnsilencedModeReArmsAfterItClears(t *testing.T) {
	sink := &logSink375{}
	can := false
	src := approval.UnsilencedMode(modeSrc375{risk.ModeAutoApprove}, func() bool { return can }, sink.f)
	if got := src.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Fatalf("first read = %s, want ask_every_step", got)
	}
	can = true
	if got := src.PermissionMode(); got != risk.ModeAutoApprove {
		t.Fatalf("answerable read = %s, want auto_approve", got)
	}
	can = false
	if got := src.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Fatalf("second downgrade = %s, want ask_every_step", got)
	}
	if n := sink.count("MODE-UNANSWERABLE"); n != 2 {
		t.Fatalf("MODE-UNANSWERABLE lines = %d, want 2 (one per episode): %v", n, sink.lines)
	}
}
