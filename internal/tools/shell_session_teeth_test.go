package tools

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
)

// Ticket 163 bridge-level teeth tests: AC#3 (the task-end boundary collects
// the session) and AC#4 (a session never becomes an authorization - every
// command is judged by the full C19 + gate path, separately).

// countingGate records how many times the bridge ASKED, and answers what
// the test told it to. The count is the AC#4 readout: two session commands
// => two gate visits, no matter that the session already existed.
type countingGate struct {
	mu        sync.Mutex
	approvals int
	allow     bool
}

func (g *countingGate) PendingWindow(context.Context, Decision) (Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.approvals++
	if g.allow {
		return AnswerAllow, ""
	}
	return AnswerReject, "test window refuses"
}

func (g *countingGate) PendingApproval(context.Context, Decision) (Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.approvals++
	if g.allow {
		return AnswerAllow, ""
	}
	return AnswerReject, "test approval refuses"
}

func (g *countingGate) visits() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.approvals
}

func (g *countingGate) setAllow(v bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.allow = v
}

// shellTestBridge composes the bridge the way ticket 163's AC#3 wiring will:
// the shell family registered, the session roster shared between the family
// and the bridge's task-end janitor.
func shellTestBridge(t *testing.T, ff *fakeShellFactory, g Gate, withJanitor bool) (*Bridge, *ShellSessions) {
	t.Helper()
	roster := NewShellSessions()
	reg := NewRegistry()
	for _, e := range BuiltinShellEntries(ShellDeps{
		Enabled:         true,
		Sessions:        roster,
		Allowlist:       []string{"echo"},
		NewSessionShell: ff.new,
	}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	opts := Options{
		Registry: reg,
		Paths:    NewPathCanonicalizer(nil, nil),
		Gate:     g,
		Logf:     func(string, ...any) {},
	}
	if withJanitor {
		opts.TaskCleanup = roster.TerminateForTask
	}
	return New(opts), roster
}

func shellBridgeCall(t *testing.T, b *Bridge, taskID string, argv ...string) agent.ToolOutcome {
	t.Helper()
	raw, err := json.Marshal(map[string]any{"argv": argv})
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		Name: "shell.session", TaskID: taskID, CallID: taskID + "-call", Args: raw,
	})
	if err != nil {
		t.Fatalf("Execute provider fault: %v", err)
	}
	return out
}

// Test163AC4SessionNeverCarriesAuthority - AC#4: the FIRST command got an
// approval, but the session it created authorizes NOTHING. The second
// command is judged again (visits==2), a refused command leaves the session
// alive but runs nothing, and the next approved command works. No grant
// path is consulted on the L2 route at all - the verdict that reaches the
// gate is L2 every time.
func Test163AC4SessionNeverCarriesAuthority(t *testing.T) {
	g := &countingGate{allow: true}
	ff := &fakeShellFactory{respond: echoResponder(0)}
	b, _ := shellTestBridge(t, ff, g, true)

	out := shellBridgeCall(t, b, "t163-ac4", "echo", "first")
	if out.IsError {
		t.Fatalf("first command must run: %s", out.Text)
	}
	if visits := g.visits(); visits != 1 {
		t.Fatalf("first command must visit the gate once, visits=%d", visits)
	}

	out = shellBridgeCall(t, b, "t163-ac4", "echo", "second")
	if out.IsError {
		t.Fatalf("second command must run: %s", out.Text)
	}
	if visits := g.visits(); visits != 2 {
		t.Fatalf("second command must be judged SEPARATELY (AC#4), visits=%d", visits)
	}

	// The refusal arm: the session exists, the approval is denied, nothing
	// runs - and the session is not destroyed by the refusal either. The
	// refusal text is the gate's own reason, carried through the bridge.
	g.setAllow(false)
	out = shellBridgeCall(t, b, "t163-ac4", "echo", "refused")
	if !out.IsError || !strings.Contains(out.Text, "test approval refuses") {
		t.Fatalf("denied command must refuse with the gate's reason: %s", out.Text)
	}
	if visits := g.visits(); visits != 3 {
		t.Fatalf("a denied command still goes through the gate, visits=%d", visits)
	}

	// And the session keeps working for approved commands: refusal was an
	// authorization answer, not a session verdict.
	g.setAllow(true)
	out = shellBridgeCall(t, b, "t163-ac4", "echo", "after-refusal")
	if out.IsError {
		t.Fatalf("session must survive a refused command: %s", out.Text)
	}
}

// Test163AC3TaskEndCollectsSessionTeeth - AC#3: the task-end boundary is
// the collector. Positive arm: CloseTask kills the session's tree. Scope
// arm: another task's boundary cannot reach this session. Negative control
// (the mutation): WITHOUT the janitor wiring the session survives CloseTask,
// which is exactly the defect the wired test above rings on.
func Test163AC3TaskEndCollectsSessionTeeth(t *testing.T) {
	g := &countingGate{allow: true}
	ff := &fakeShellFactory{respond: echoResponder(0)}
	b, roster := shellTestBridge(t, ff, g, true)

	out := shellBridgeCall(t, b, "t163-ac3", "echo", "up")
	if out.IsError {
		t.Fatalf("session must open: %s", out.Text)
	}
	if roster.Count() != 1 {
		t.Fatalf("session row expected, count=%d", roster.Count())
	}

	// The task ends: the session goes with it.
	b.CloseTask("t163-ac3")
	if roster.Count() != 0 {
		t.Fatalf("session must be collected at the task boundary, count=%d", roster.Count())
	}
	if ff.made[0].closeCount() == 0 {
		t.Fatalf("collection must close the carrier's tree")
	}

	// Scope arm: two tasks, one boundary each, no cascade.
	outA := shellBridgeCall(t, b, "t163-ac3-a", "echo", "a")
	outB := shellBridgeCall(t, b, "t163-ac3-b", "echo", "b")
	if outA.IsError || outB.IsError {
		t.Fatalf("both sessions must open: %q %q", outA.Text, outB.Text)
	}
	b.CloseTask("t163-ac3-a")
	if roster.Count() != 1 {
		t.Fatalf("only task A's session must be collected, count=%d", roster.Count())
	}
	outB2 := shellBridgeCall(t, b, "t163-ac3-b", "echo", "b-alive")
	if outB2.IsError || !strings.Contains(outB2.Text, "退出码：0") {
		t.Fatalf("task B's session must survive A's boundary: %s", outB2.Text)
	}
	b.CloseTask("t163-ac3-b")
	if roster.Count() != 0 {
		t.Fatalf("all sessions collected at the end, count=%d", roster.Count())
	}
}

// Test163AC3UnwiredJanitorIsTheDefect - the negative control that gives the
// test above its teeth: with the janitor unwired, CloseTask leaves the
// session alive. Read together with Test163AC3TaskEndCollectsSessionTeeth:
// that test fails (session survives) exactly when this wiring is removed,
// which is the mutation AC#3 demands proof of.
func Test163AC3UnwiredJanitorIsTheDefect(t *testing.T) {
	g := &countingGate{allow: true}
	ff := &fakeShellFactory{respond: echoResponder(0)}
	b, roster := shellTestBridge(t, ff, g, false)

	if out := shellBridgeCall(t, b, "t163-mut", "echo", "up"); out.IsError {
		t.Fatalf("session must open: %s", out.Text)
	}
	b.CloseTask("t163-mut")
	if roster.Count() != 1 {
		t.Fatalf("unwired janitor leaves the session alive (that IS the defect this pair detects)")
	}
	if ff.made[0].closeCount() != 0 {
		t.Fatalf("nothing may kill the session outside the wired boundary")
	}
}

// Test163AC4GrantsNeverCoverTheSessionRoute - the D45 grant source is
// consulted only on L1; the session's L2 verdict never sees it. With a
// grant source that answers EVERYTHING, the gate is still asked.
func Test163AC4GrantsNeverCoverTheSessionRoute(t *testing.T) {
	g := &countingGate{allow: true}
	ff := &fakeShellFactory{respond: echoResponder(0)}
	roster := NewShellSessions()
	reg := NewRegistry()
	for _, e := range BuiltinShellEntries(ShellDeps{
		Enabled: true, Sessions: roster, NewSessionShell: ff.new,
	}) {
		if err := reg.Register(e); err != nil {
			t.Fatal(err)
		}
	}
	b := New(Options{
		Registry: reg,
		Paths:    NewPathCanonicalizer(nil, nil),
		Gate:     g,
		Grants:   promiscuousGrantSource{},
		Logf:     func(string, ...any) {},
	})
	out := shellBridgeCall(t, b, "t163-grant", "echo", "x")
	if out.IsError {
		t.Fatalf("allowed command must run: %s", out.Text)
	}
	if g.visits() != 1 {
		t.Fatalf("the L2 verdict must reach the gate despite a grant source, visits=%d", g.visits())
	}
}

// promiscuousGrantSource answers every lookup - the widest possible grant
// world, under which the session's L2 route must STILL ask the gate.
type promiscuousGrantSource struct{}

func (promiscuousGrantSource) Covering(context.Context, string, []string) (int64, bool) {
	return 99, true
}

var _ GrantSource = promiscuousGrantSource{}
