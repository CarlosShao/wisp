package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// Ticket 163 unit tests: the session protocol driven through fake carriers
// (deterministic, no real cmd.exe here - the live Windows arms are in
// shell_live_windows_test.go). The fake simulates cmd.exe's half of the
// marker protocol: it answers every stdin line the way a shell would.

// fakeSessionShell is a SessionShell whose stdout is an io.Pipe the test's
// respond function writes into. It records every kill/wait/close so the
// teeth tests can read what actually happened to the carrier.
type fakeSessionShell struct {
	pipeW *io.PipeWriter
	pipeR *io.PipeReader

	mu       sync.Mutex
	closes   int
	kills    int
	waits    int
	closeErr error
	respond  func(f *fakeSessionShell, line string)

	closeOnce sync.Once
}

func newFakeSessionShell(respond func(f *fakeSessionShell, line string)) *fakeSessionShell {
	r, w := io.Pipe()
	return &fakeSessionShell{pipeW: w, pipeR: r, respond: respond}
}

func (f *fakeSessionShell) Stdin() io.Writer { return f }

// Write records each line and runs the responder for it. This runs on the
// tool's goroutine; the response bytes are consumed by the session's reader
// goroutine, so a pipe write here can only block while the session lives -
// after a kill the closed pipe errors instead of blocking.
func (f *fakeSessionShell) Write(p []byte) (int, error) {
	n := len(p)
	for _, line := range strings.Split(string(p), "\r\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		f.mu.Lock()
		respond := f.respond
		f.mu.Unlock()
		if respond != nil {
			respond(f, line)
		}
	}
	return n, nil
}

func (f *fakeSessionShell) Stdout() io.Reader { return f.pipeR }

func (f *fakeSessionShell) TakeStderr() string { return "" }

func (f *fakeSessionShell) Kill() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.kills++
	return nil
}

func (f *fakeSessionShell) Wait() (*os.ProcessState, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.waits++
	return nil, nil
}

// Close simulates the kill-on-close tree release: the pipe dies, the reader
// goroutine sees EOF, and the counters record the collection.
func (f *fakeSessionShell) Close() error {
	f.mu.Lock()
	f.closes++
	f.mu.Unlock()
	f.closeOnce.Do(func() { _ = f.pipeW.Close() })
	return nil
}

func (f *fakeSessionShell) closeCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.closes
}

// echoResponder simulates the shell's half of the marker protocol: plain
// lines answer "ok: <line>"; the marker echo line answers
// "<marker> <code>". This is the exact shape the live cmd.exe produces.
func echoResponder(code int) func(*fakeSessionShell, string) {
	return func(f *fakeSessionShell, line string) {
		if strings.HasPrefix(line, "echo w163") {
			fields := strings.Fields(line)
			if len(fields) >= 3 && fields[2] == "%errorlevel%" {
				_, _ = f.pipeW.Write([]byte(fmt.Sprintf("%s %d\r\n", fields[1], code)))
				return
			}
		}
		_, _ = f.pipeW.Write([]byte("ok: " + line + "\r\n"))
	}
}

// baitResponder is the "摘掉标记" world: output arrives, an exit-code-shaped
// line arrives, but the marker boundary never does.
func baitResponder(bait string) func(*fakeSessionShell, string) {
	return func(f *fakeSessionShell, line string) {
		if strings.HasPrefix(line, "echo w163") {
			return // no marker line: the stream never gets its boundary
		}
		_, _ = f.pipeW.Write([]byte(bait + "\r\n"))
	}
}

// endResponder simulates the shell exiting before the marker line runs
// (`exit N` kills cmd.exe before our echo executes).
func endResponder() func(*fakeSessionShell, string) {
	return func(f *fakeSessionShell, line string) {
		if strings.HasPrefix(line, "echo w163") {
			_ = f.pipeW.Close()
			return
		}
		_, _ = f.pipeW.Write([]byte("bye: " + line + "\r\n"))
	}
}

// fakeShellFactory hands out one fake per session creation, so a reset test
// can tell the first carrier from the second.
type fakeShellFactory struct {
	mu      sync.Mutex
	respond func(*fakeSessionShell, string)
	made    []*fakeSessionShell
}

func (ff *fakeShellFactory) new() (SessionShell, error) {
	ff.mu.Lock()
	defer ff.mu.Unlock()
	f := newFakeSessionShell(ff.respond)
	ff.made = append(ff.made, f)
	return f, nil
}

func shellDepsWith(ff *fakeShellFactory) ShellDeps {
	return ShellDeps{
		Enabled:         true,
		Sessions:        NewShellSessions(),
		Paths:           NewPathCanonicalizer(nil, nil),
		NewSessionShell: ff.new,
	}
}

func shellTaskCtx(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return withCancel(ctx, cancelHandle{taskID: "t163"}), cancel
}

func sessionCall(t *testing.T, d ShellDeps, argv []string, extra string) Result {
	t.Helper()
	raw := fmt.Sprintf(`{"argv":%s%s}`, shellArgvJSON(argv), extra)
	res, err := (shellSession{d: d}).Execute(mustTaskCtx(t), json.RawMessage(raw), nil)
	if err != nil {
		t.Fatalf("Execute: unexpected provider fault: %v", err)
	}
	return res
}

// shellArgvJSON marshals one argv vector; mustJSON in web_testhelp_test.go
// already owns the generic name, so this stays family-scoped.
func shellArgvJSON(argv []string) string {
	b, err := json.Marshal(argv)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func mustTaskCtx(t *testing.T) context.Context {
	t.Helper()
	return withCancel(context.Background(), cancelHandle{taskID: "t163"})
}

// Test163SessionMarkerCarriesExitCode - AC#2 forward: the exit code comes
// from the marker boundary, deterministically, and the marker line itself
// never reaches the model. A nonzero errorlevel reads as an error result.
func Test163SessionMarkerCarriesExitCode(t *testing.T) {
	ff := &fakeShellFactory{respond: echoResponder(7)}
	d := shellDepsWith(ff)
	res := sessionCall(t, d, []string{"echo", "hello"}, "")
	if !res.IsError {
		t.Fatalf("nonzero errorlevel must read as an error: %s", res.Text)
	}
	if !strings.Contains(res.Text, "ok: echo hello") {
		t.Fatalf("command output missing: %s", res.Text)
	}
	if !strings.Contains(res.Text, "退出码：7") {
		t.Fatalf("marker exit code missing: %s", res.Text)
	}
	if strings.Contains(res.Text, "w163") {
		t.Fatalf("marker leaked into the result: %s", res.Text)
	}
	// The zero arm: same protocol, clean exit reads as success.
	ff2 := &fakeShellFactory{respond: echoResponder(0)}
	res = sessionCall(t, shellDepsWith(ff2), []string{"echo", "hello"}, "")
	if res.IsError || !strings.Contains(res.Text, "退出码：0") {
		t.Fatalf("zero errorlevel must read as success: %s", res.Text)
	}
}

// Test163SessionMarkerIsNotAGuess - the a1 live-sample trap, session form:
// the bait output contains an exit-code-shaped line, but the marker boundary
// never arrives, so no code may be claimed (least of all the bait's).
func Test163SessionMarkerIsNotAGuess(t *testing.T) {
	ff := &fakeShellFactory{respond: baitResponder("exit 1")}
	d := shellDepsWith(ff)
	ctx, cancel := shellTaskCtx(300 * time.Millisecond)
	defer cancel()
	res, err := (shellSession{d: d}).Execute(ctx, json.RawMessage(`{"argv":["type","big.log"]}`), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError {
		t.Fatalf("a desynced stream must not read as success: %s", res.Text)
	}
	if !strings.Contains(res.Text, "失去同步") {
		t.Fatalf("timeout text must name the desync: %s", res.Text)
	}
	if strings.Contains(res.Text, "退出码：0") || strings.Contains(res.Text, "退出码：1") {
		t.Fatalf("no exit code may be claimed without the marker: %s", res.Text)
	}
	if ff.made[0].closeCount() == 0 {
		t.Fatalf("the desynced session must be collected (tree closed)")
	}
	// AC#2's answer, operationalized: after the reset the next call starts a
	// FRESH session - nothing of the dead one is reachable. The new carrier
	// gets a healthy responder: it is a new shell, not the desynced one.
	ff.mu.Lock()
	ff.respond = echoResponder(0)
	ff.mu.Unlock()
	res2 := sessionCall(t, d, []string{"echo", "fresh"}, "")
	if res2.IsError || !strings.Contains(res2.Text, "退出码：0") {
		t.Fatalf("fresh session must answer normally: %s", res2.Text)
	}
	if len(ff.made) != 2 {
		t.Fatalf("the reset must have built a new carrier, made=%d", len(ff.made))
	}
}

// Test163SessionShellExitUsesWaitField - the a1 ruling's primary channel,
// used exactly where a process boundary exists: the shell process exited,
// so the code comes from Wait, not from any marker.
func Test163SessionShellExitUsesWaitField(t *testing.T) {
	ff := &fakeShellFactory{respond: endResponder()}
	d := shellDepsWith(ff)
	res := sessionCall(t, d, []string{"exit"}, "")
	if !strings.Contains(res.Text, "会话进程已退出") {
		t.Fatalf("ended-session text missing: %s", res.Text)
	}
	if ff.made[0].closeCount() == 0 {
		t.Fatalf("ended session must be collected")
	}
}

// Test163SessionRosterTerminateKillsTree - AC#3 unit arm: the roster's
// collector closes the tree, and the collector is the only one.
func Test163SessionRosterTerminateKillsTree(t *testing.T) {
	ff := &fakeShellFactory{respond: echoResponder(0)}
	d := shellDepsWith(ff)
	roster := d.Sessions
	_ = sessionCall(t, d, []string{"echo", "one"}, "")
	if roster.Count() != 1 {
		t.Fatalf("session row expected, count=%d", roster.Count())
	}
	if got := roster.TerminateForTask("t163"); got != 1 {
		t.Fatalf("TerminateForTask=%d, want 1", got)
	}
	if roster.Count() != 0 {
		t.Fatalf("row must be gone after collection, count=%d", roster.Count())
	}
	if ff.made[0].closeCount() == 0 {
		t.Fatalf("collection must close the tree")
	}
	if got := roster.TerminateForTask("t163"); got != 0 {
		t.Fatalf("second collection=%d, want 0 (nothing left)", got)
	}
	if got := roster.TerminateForTask(""); got != 0 {
		t.Fatalf("empty id must collect nothing, got %d", got)
	}
	// A foreign task's id collects nothing (AC#3 scope arm, unit level).
	_ = sessionCall(t, d, []string{"echo", "two"}, "")
	if got := roster.TerminateForTask("t163-other"); got != 0 {
		t.Fatalf("foreign id=%d, want 0", got)
	}
	if roster.Count() != 1 {
		t.Fatalf("own session must survive a foreign collection")
	}
}

// Test163SessionRefusalsAreLoud - every unwired/misused shape refuses with
// a reason, never with an empty success.
func Test163SessionRefusalsAreLoud(t *testing.T) {
	// no task id: no scope to bind a session to.
	ctxNoTask := context.Background()
	res, err := (shellSession{d: shellDepsWith(&fakeShellFactory{respond: echoResponder(0)})}).
		Execute(ctxNoTask, json.RawMessage(`{"argv":["echo","x"]}`), nil)
	if err != nil || !res.IsError || !strings.Contains(res.Text, "任务 id") {
		t.Fatalf("no-task-id call must refuse loudly: %q %v", res.Text, err)
	}
	// nil roster: nothing can collect the session later.
	d := shellDepsWith(&fakeShellFactory{respond: echoResponder(0)})
	d.Sessions = nil
	res = sessionCall(t, d, []string{"echo", "x"}, "")
	if !res.IsError || !strings.Contains(res.Text, "会话名册未接线") {
		t.Fatalf("nil roster must refuse: %s", res.Text)
	}
	// no carrier factory and no tree factory: refuse.
	d2 := ShellDeps{Enabled: true, Sessions: NewShellSessions()}
	res = sessionCall(t, d2, []string{"echo", "x"}, "")
	if !res.IsError || !strings.Contains(res.Text, "进程树管控未接线") {
		t.Fatalf("unwired tree must refuse: %s", res.Text)
	}
	// empty argv.
	res = sessionCall(t, shellDepsWith(&fakeShellFactory{respond: echoResponder(0)}), nil, "")
	if !res.IsError || !strings.Contains(res.Text, "argv") {
		t.Fatalf("empty argv must refuse: %s", res.Text)
	}
}

// Test163SessionCwdOnlyAtCreation - cwd on a live session is a refusal,
// never a silent no-op.
func Test163SessionCwdOnlyAtCreation(t *testing.T) {
	ff := &fakeShellFactory{respond: echoResponder(0)}
	d := shellDepsWith(ff)
	res := sessionCall(t, d, []string{"echo", "up"}, `,"cwd":"C:\\tmp\\w163"`)
	if res.IsError {
		t.Fatalf("first call with cwd must create the session: %s", res.Text)
	}
	res = sessionCall(t, d, []string{"echo", "again"}, `,"cwd":"C:\\tmp\\other"`)
	if !res.IsError || !strings.Contains(res.Text, "只在会话建立的那一发生效") {
		t.Fatalf("cwd on a live session must refuse: %s", res.Text)
	}
}

// Test163JoinSessionArgv - the session line builder refuses what it cannot
// represent and quotes what it can.
func Test163JoinSessionArgv(t *testing.T) {
	line, err := joinSessionArgv([]string{"git", "status"})
	if err != nil || line != "git status" {
		t.Fatalf("simple join: %q %v", line, err)
	}
	line, err = joinSessionArgv([]string{`C:\Program Files\git\bin\git.exe`, "--version"})
	if err != nil || line != `"C:\Program Files\git\bin\git.exe" --version` {
		t.Fatalf("quoted join: %q %v", line, err)
	}
	for name, argv := range map[string][]string{
		"empty element":   {"git", ""},
		"space element":   {"git", "   "},
		"quote element":   {"echo", `he"llo`},
		"newline":         {"echo", "a\nb"},
		"carriage return": {"echo", "a\rb"},
		"nul":             {"echo", "a\x00b"},
	} {
		if _, err := joinSessionArgv(argv); err == nil {
			t.Fatalf("%s must be refused, got a line", name)
		}
	}
	long := make([]string, 0, 2)
	long = append(long, "echo", strings.Repeat("a", sessionLineLimit))
	if _, err := joinSessionArgv(long); err == nil {
		t.Fatalf("over-limit line must be refused")
	}
}
