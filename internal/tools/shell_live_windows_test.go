//go:build windows

package tools

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/proc"
)

// Ticket 163 live Windows arms: REAL cmd.exe carriers inside REAL
// kill-on-close jobs (proc.OpenJobScope), the production carrier shape.
// These are the readouts the ticket's ACs quote: cwd continuity (the AC#1
// gap, closed), the Wait-field exit code (a1's trap sample), the marker
// code, the timeout reset, and shell.exec's tree kill.

// procTreeFactory hands out one real kill-on-close job per call/session -
// the closure the composition root will pass. The compile-time assertion
// pins that proc.JobScope's method set IS SessionTree, so the production
// shape cannot drift away from the seam silently.
var _ SessionTree = (*proc.JobScope)(nil)

func procTreeFactory() func() (SessionTree, error) {
	return func() (SessionTree, error) { return proc.OpenJobScope() }
}

func liveShellDeps(t *testing.T) ShellDeps {
	t.Helper()
	roster := NewShellSessions()
	t.Cleanup(func() {
		for _, id := range []string{"t163-live", "t163-live2", "t163-live-exec"} {
			roster.TerminateForTask(id)
		}
	})
	return ShellDeps{
		Enabled:  true,
		Sessions: roster,
		Paths:    NewPathCanonicalizer(nil, nil),
		NewTree:  procTreeFactory(),
	}
}

func liveCall(t *testing.T, d ShellDeps, taskID string, raw string) Result {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ctx = withCancel(ctx, cancelHandle{taskID: taskID})
	res, err := (shellSession{d: d}).Execute(ctx, json.RawMessage(raw), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	return res
}

// Test163LiveSessionCwdCarriesOver - the AC#1 counter-read: two calls, ONE
// session. Call one establishes the session in a temp dir; call two asks
// `cd` and must see THAT dir - the exact measurement that made this ticket
// exist (two one-shot spawns printed different directories).
func Test163LiveSessionCwdCarriesOver(t *testing.T) {
	d := liveShellDeps(t)
	tmp := mustCanonical(t, t.TempDir())
	raw, _ := json.Marshal(map[string]any{"argv": []string{"echo", "up"}, "cwd": tmp})
	res := liveCall(t, d, "t163-live", string(raw))
	if res.IsError {
		t.Fatalf("session with cwd must open: %s", res.Text)
	}

	res = liveCall(t, d, "t163-live", `{"argv":["cd"]}`)
	if res.IsError {
		t.Fatalf("cd must run: %s", res.Text)
	}
	// Case-insensitive on purpose: cmd.exe prints the directory in its own
	// spelling; the fact under test is WHICH directory, not its case.
	lower := strings.ToLower(res.Text)
	if !strings.Contains(lower, strings.ToLower(tmp)) {
		t.Fatalf("second call lost the first call's cwd (AC#1 would ring):\ntext=%s\nwant cwd=%s", res.Text, tmp)
	}
	if strings.Contains(lower, strings.ToLower("退出码：1")) {
		t.Fatalf("cd must exit 0: %s", res.Text)
	}
	// Move the shell's working directory out of the temp dir BEFORE the test
	// ends: the live cmd.exe holds a handle on its cwd, and a kill that lands
	// while it still sits in the dir races t.TempDir's RemoveAll.
	if res := liveCall(t, d, "t163-live", `{"argv":["cd","/d","C:\\"]}`); res.IsError {
		t.Fatalf("cd out must run: %s", res.Text)
	}
}

// Test163LiveSessionMarkerExitCode - AC#2 forward on a real shell: a child
// cmd that exits 7 sets errorlevel, and the marker line carries exactly 7.
func Test163LiveSessionMarkerExitCode(t *testing.T) {
	d := liveShellDeps(t)
	res := liveCall(t, d, "t163-live", `{"argv":["cmd","/c","exit 7"]}`)
	if !strings.Contains(res.Text, "退出码：7") {
		t.Fatalf("marker exit code must be 7: %s", res.Text)
	}
	if !res.IsError {
		t.Fatalf("nonzero exit must read as error: %s", res.Text)
	}
}

// Test163LiveSessionShellExitWaitField - the a1 ruling's primary channel on
// a real shell: `exit 5` kills cmd.exe before the marker line runs, so the
// code comes from cmd's Wait state field.
func Test163LiveSessionShellExitWaitField(t *testing.T) {
	d := liveShellDeps(t)
	res := liveCall(t, d, "t163-live", `{"argv":["exit","5"]}`)
	if !strings.Contains(res.Text, "会话进程已退出") || !strings.Contains(res.Text, "退出码：5") {
		t.Fatalf("Wait-field exit code must be 5: %s", res.Text)
	}
}

// Test163LiveSessionTimeoutResets - AC#2/AC#3 on a real shell: a command
// that never finishes loses its budget; the session is collected (no exit
// code, no attributable output), and the next call runs in a fresh shell.
func Test163LiveSessionTimeoutResets(t *testing.T) {
	d := liveShellDeps(t)
	roster := d.Sessions
	raw := `{"argv":["ping","-n","30","127.0.0.1"]}`
	ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
	defer cancel()
	ctx = withCancel(ctx, cancelHandle{taskID: "t163-live2"})
	res, err := (shellSession{d: d}).Execute(ctx, json.RawMessage(raw), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "失去同步") {
		t.Fatalf("timeout must read as a desync: %s", res.Text)
	}
	if strings.Contains(res.Text, "退出码：") {
		t.Fatalf("timeout must not claim an exit code: %s", res.Text)
	}

	res2 := liveCall(t, d, "t163-live2", `{"argv":["echo","fresh"]}`)
	if res2.IsError || !strings.Contains(res2.Text, "fresh") {
		t.Fatalf("fresh session must answer: %s", res2.Text)
	}
	if roster.Count() != 1 {
		t.Fatalf("one live session row expected after the reset, count=%d", roster.Count())
	}
}

// fakeExecTree is a SessionTree that starts the command for real and kills
// it on Close - the kill path without a job object, so the timeout arm can
// prove tree.Close is what terminates the process.
type fakeExecTree struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	closed bool
}

func (t *fakeExecTree) StartInJob(cmd *exec.Cmd) (*os.Process, error) {
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	t.mu.Lock()
	t.cmd = cmd
	t.mu.Unlock()
	return cmd.Process, nil
}

func (t *fakeExecTree) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
	if t.cmd != nil && t.cmd.Process != nil {
		_ = t.cmd.Process.Kill()
	}
	return nil
}

func (t *fakeExecTree) wasClosed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.closed
}

// Test163LiveShellExecExitCodeFromWaitField - shell.exec's AC#2 arm on a
// real process: code from Wait; and a1's trap sample (a process printing
// "exit 1" while exiting 0) must read 0.
func Test163LiveShellExecExitCodeFromWaitField(t *testing.T) {
	tree := &fakeExecTree{}
	d := ShellDeps{Enabled: true, NewTree: func() (SessionTree, error) { return tree, nil }}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ctx = withCancel(ctx, cancelHandle{taskID: "t163-live-exec"})

	res, err := (shellExec{d: d}).Execute(ctx, json.RawMessage(`{"argv":["cmd","/c","exit 7"]}`), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Text, "退出码：7") || !res.IsError {
		t.Fatalf("exit 7 must read as 7/error: %s", res.Text)
	}

	tree2 := &fakeExecTree{}
	d2 := ShellDeps{Enabled: true, NewTree: func() (SessionTree, error) { return tree2, nil }}
	res, err = (shellExec{d: d2}).Execute(ctx, json.RawMessage(`{"argv":["cmd","/c","echo exit 1"]}`), nil)
	if err != nil {
		t.Fatal(err)
	}
	// The a1 trap: the TEXT says "exit 1"; the verdict must say 0 - the code
	// comes from the Wait state field, never from the output text.
	if !strings.Contains(res.Text, "exit 1") {
		t.Fatalf("trap sample output missing: %s", res.Text)
	}
	if strings.Contains(res.Text, "退出码：1") || res.IsError {
		t.Fatalf("the field must say 0 even though the text says exit 1: %s", res.Text)
	}
	if !strings.Contains(res.Text, "退出码：0") {
		t.Fatalf("exit code 0 must be reported: %s", res.Text)
	}
}

// Test163LiveShellExecTimeoutKillsTree - the C22 deadline arm: no exit code
// is claimed, the tree is closed (kill), and nothing is left running.
func Test163LiveShellExecTimeoutKillsTree(t *testing.T) {
	tree := &fakeExecTree{}
	d := ShellDeps{Enabled: true, NewTree: func() (SessionTree, error) { return tree, nil }}
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()
	ctx = withCancel(ctx, cancelHandle{taskID: "t163-live-exec"})

	res, err := (shellExec{d: d}).Execute(ctx, json.RawMessage(`{"argv":["cmd","/c","ping -n 30 127.0.0.1 > NUL"]}`), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !res.IsError || !strings.Contains(res.Text, "C22") {
		t.Fatalf("timeout must be reported: %s", res.Text)
	}
	if strings.Contains(res.Text, "退出码：") {
		t.Fatalf("a killed process must not yield an exit code: %s", res.Text)
	}
	if !tree.wasClosed() {
		t.Fatalf("the tree must be closed on the timeout path")
	}
}
