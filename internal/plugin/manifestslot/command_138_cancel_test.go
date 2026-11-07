package manifestslot

// Ticket 138 AC#3 — the tooth on the TASK-CANCEL cascade.
//
// Test52TimeoutKillsTheProcessTree pins the BUDGET side of the kill: the
// deadline fires and the tree dies. This test pins the other half of the same
// watcher: the caller's ctx cancelling (what RunningTask.Cancel() does to the
// task root, internal/agent/loop.go) must ALSO reach the wrapped CLI's process
// tree. The cascade travels caller ctx -> withMonotonicBudget (a
// context.WithDeadline derived from the parent) -> the budget goroutine's
// childCtx.Done() arm -> ChildProcess.KillTree -> the per-child Windows Job
// Object -> TerminateJobObject.
//
// The mutation this test is the tooth against: detaching the budget from its
// parent (stopwatch.go, context.WithDeadline(context.Background(), ...)). The
// budget would still expire on its own (the timeout tooth stays green) but the
// caller's cancel would stop propagating - a stopped task would leave its
// hanging CLI running. This test goes red exactly there, and the evidence file
// docs/evidence/s1/138-cascade-ruling.md (§3) records the mutated line, the
// red line and the restored diff.

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/plugin"
	"github.com/CarlosShao/wisp/internal/tools"
)

func Test138TaskCancelKillsTheProcessTree(t *testing.T) {
	exe := installFixtureCLI(t)
	pidFile := filepath.Join(stableExeDir(t), "cancel-grandchild.pid")
	def := fixtureTool(exe, "hang")
	def.Extract = map[string]any{"kind": "regex", "pattern": `.+`}
	// The budget must NOT be the actor in this test: 120s of headroom means
	// only the caller's cancel can end the call while the test is running.
	def.TimeoutMs = plugin.MaxCommandTimeoutMS / 5 // 120000ms

	env := func(name string) (string, bool) {
		switch name {
		case "PATH":
			return os.Getenv("PATH"), true
		case "SystemRoot":
			return os.Getenv("SystemRoot"), true
		case ModuleNameEnv:
			return "1", true
		case fixtureHandshakeVar:
			return fixtureModePrefix + "hang:" + pidFile, true
		}
		return "", false
	}

	runtime, _ := runtimeFor(t, def, CommandRunOptions{EnvLookup: env})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	type outcome struct {
		res tools.Result
		err error
	}
	done := make(chan outcome, 1)
	go func() {
		res, err := runtime.Exec(ctx, "com.example.fixture", def, []byte(`{}`))
		done <- outcome{res, err}
	}()

	// The grandchild must be alive BEFORE the cancel: a tooth that never
	// bit a living tree proves nothing.
	var gcPID int
	deadline := time.Now().Add(10 * time.Second)
	for gcPID == 0 {
		if time.Now().After(deadline) {
			t.Fatalf("the fixture never reported a grandchild pid (%s)", pidFile)
		}
		gcPID = readPIDFile(t, pidFile)
		if gcPID == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	// THE action under test: the task's cancel, not a budget expiry.
	cancel()

	select {
	case o := <-done:
		if o.err != nil {
			t.Fatalf("Exec after the caller's cancel: %v", o.err)
		}
		_ = o.res
	case <-time.After(15 * time.Second):
		t.Fatalf("Exec did not return within 15s of the caller's cancel: the cancel cascade never reached the child (stopwatch budget detached from its parent?)")
	}

	// THE tooth: the grandchild must be gone. Killing only the direct child
	// (a plain Process.Kill) leaves exactly this pid running.
	alive := true
	for i := 0; i < 100; i++ {
		if !stillAlive(gcPID) {
			alive = false
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if alive {
		t.Fatalf("grandchild pid %d is still running after the task cancel: the cancel cascade did not reach the tree", gcPID)
	}
}
