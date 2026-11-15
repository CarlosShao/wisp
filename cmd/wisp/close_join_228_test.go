package main

// Ticket 228 AC#4: D38(c) (PLAN.md:2845, "任务完成必须等所有派生 goroutine 退出")
// says a run is not finished while goroutines it started are still alive. Until
// this slice agentRuntime.close did not meet that: it cancelled the two roots and
// returned, so `wisp run` declared itself done with the answer listener and the
// config-reload tick still on the books ("the cancel is not a join" was the
// comment's own wording).
//
// These cases pin the two halves of the fix WITHOUT a desktop and without a
// whole run: one owns a real runtime, spawns the two goroutines through the same
// observe.Registry.Spawn the production assembly uses, and watches close().
//
//   - TestAC4CloseWaitsForBothGoroutinesItSpawned: with one line of work still in
//     flight, close() must come back only after both goroutines finished. Delete
//     the wait and this case reds on the pending counters - it does not "look
//     steadier", it names which goroutine was left behind.
//   - TestAC4CloseNamesTheParkedListenerInsteadOfWaitingOnIt: the listener's
//     normal resting place is a BLOCKING read, which a cancel cannot wake. close()
//     must not spend its budget on a wake-up that cannot arrive; it must RECORD
//     that fact under the goroutine's name instead. Delete that sentence and this
//     case reds.
//
// Both are `windows`-free: close() lives in run.go (no build tag), and everything
// here is the registry plus timers.

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// ac4Runtime228 is a runtime whose close() can run with no store, no scheduler
// and no retention job - each of those is nil-guarded in close(), which is what
// lets this case own exactly the two goroutines under test. stderr is a buffer so
// the audit line close() books can be read back (auditf writes to
// agentRuntime.stderr, which assembleRuntime fills from the run spec).
func ac4Runtime228(t *testing.T) (*agentRuntime, *bytes.Buffer) {
	t.Helper()
	audit := &bytes.Buffer{}
	return &agentRuntime{spec: runSpec{stderr: audit}, stderr: audit}, audit
}

// doneClosed answers whether the goroutine behind h has finished.
func doneClosed(h *observe.Handle) bool {
	if h == nil {
		return true
	}
	select {
	case <-h.Done():
		return true
	default:
		return false
	}
}

func TestAC4CloseWaitsForBothGoroutinesItSpawned(t *testing.T) {
	rt, audit := ac4Runtime228(t)
	rt.replyRoot = observe.NewRoot("approval-reply")
	rt.reloadRoot = observe.NewRoot("config-reload")

	// The listener is NOT parked: it is finishing one accepted reply line, which
	// is the work the budget exists to protect (a ledger spend + an audit append).
	rt.reply = &replySurface{}
	rt.reply.noteParked(false)
	rt.replyHandle = observe.Default.Spawn("approval-waiter", "test-ac4", rt.replyRoot, func(ctx context.Context) {
		<-ctx.Done()
		time.Sleep(120 * time.Millisecond)
	})
	// The tick leaves on its cancelled root, which is the ordinary case.
	rt.reloadHandle = observe.Default.Spawn("watchdog", "test-ac4", rt.reloadRoot, func(ctx context.Context) {
		<-ctx.Done()
	})

	rt.close()

	if !doneClosed(rt.replyHandle) {
		t.Errorf("AC#4 RED: close() returned while the answer listener (approval-waiter) was still finishing "+
			"its line. D38(c) says the goroutine has to be waited for; deleting the join in close() is exactly "+
			"the shape this case exists to catch. audit line(s):\n%s", audit.String())
	}
	if !doneClosed(rt.reloadHandle) {
		t.Errorf("AC#4 RED: close() returned while the config-reload tick (watchdog) was still alive.\n%s", audit.String())
	}
	if p := rt.replyRoot.Pending(); p != 0 {
		t.Errorf("AC#4 RED: the reply root still counts %d pending goroutine(s) after close()", p)
	}
	if p := rt.reloadRoot.Pending(); p != 0 {
		t.Errorf("AC#4 RED: the reload root still counts %d pending goroutine(s) after close()", p)
	}
	if strings.Contains(audit.String(), "joined=false") {
		t.Errorf("AC#4 RED: close() booked a failed join for work it did wait for:\n%s", audit.String())
	}
}

func TestAC4CloseNamesTheParkedListenerInsteadOfWaitingOnIt(t *testing.T) {
	rt, audit := ac4Runtime228(t)
	rt.replyRoot = observe.NewRoot("approval-reply")

	// The production resting place: a blocking read no cancel can reach. The
	// goroutine deliberately ignores its context, exactly as a scanner parked in
	// Read on the operator's console does.
	rt.reply = &replySurface{}
	rt.reply.noteParked(true)
	block := make(chan struct{})
	rt.replyHandle = observe.Default.Spawn("approval-waiter", "test-ac4", rt.replyRoot, func(ctx context.Context) {
		<-block
	})
	t.Cleanup(func() {
		close(block)
		select {
		case <-rt.replyHandle.Done():
		case <-time.After(5 * time.Second):
			t.Errorf("AC#4: the parked listener did not leave after its read was released")
		}
	})

	start := time.Now()
	rt.close()
	elapsed := time.Since(start)

	if elapsed >= joinOnCloseBudget228 {
		t.Errorf("AC#4 RED: close() spent %s (the whole join budget is %s) on a listener that told it it was "+
			"parked in its read; a cancel cannot wake that read, so the wait buys nothing and the operator "+
			"pays for it on every run.", elapsed, joinOnCloseBudget228)
	}
	said := audit.String()
	if !strings.Contains(said, "GOROUTINE-JOIN goroutine=approval-waiter joined=false reason=parked-in-read") {
		t.Errorf("AC#4 RED: close() left the answer listener unjoined WITHOUT recording it. That sentence is the "+
			"difference between the pre-ticket assumption (\"the goroutine dies with the process\", stated in a "+
			"comment with nothing measuring it) and a reading taken at the moment it happens. audit line(s):\n%s",
			said)
	}
}
