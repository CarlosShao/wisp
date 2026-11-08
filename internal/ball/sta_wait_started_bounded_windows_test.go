//go:build windows

package ball

// Ticket 408, the production half's tooth: waitStarted must be a BOUNDED
// wait. A never-signalled started channel (the WebView2 born-dead shape: the
// start goroutine panics before close(started) runs) has to come back as a
// NAMED error inside the budget - never as a silent permanent hang.
//
// How this tooth turns red: revert waitStarted to the bare <-s.started
// receive (the pre-408 shape) and the guard below fires - the first case
// never delivers and the test fails loudly at the guard instead of hanging
// the whole run. Either way the revert is red; the unbounded shape can no
// longer pass this file.

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestWaitStartedBoundedReturnsNamedError injects a started channel that is
// never closed and asserts waitStarted returns a named error within the
// (shortened) budget, guarded by its own slightly-longer timer so a revert
// to the unbounded receive fails the test loudly rather than hanging it.
func TestWaitStartedBoundedReturnsNamedError(t *testing.T) {
	old := bootWaitBudget
	bootWaitBudget = 200 * time.Millisecond
	t.Cleanup(func() { bootWaitBudget = old })

	s := newSTAThread(nil) // started never closes: the wedged-boot shape
	type outcome struct{ err error }
	done := make(chan outcome, 1)
	go func() { done <- outcome{s.waitStarted()} }()

	guard := time.NewTimer(5 * time.Second)
	defer guard.Stop()
	select {
	case o := <-done:
		if o.err == nil {
			t.Fatal("waitStarted returned nil on a never-signalled boot: the bounded wait is gone")
		}
		if !strings.Contains(o.err.Error(), "ui-sta thread did not report window creation within") {
			t.Fatalf("waitStarted returned a non-named error on a wedged boot: %v", o.err)
		}
	case <-guard.C:
		t.Fatalf("waitStarted did not come back within the %s budget + guard: the wait is unbounded again", bootWaitBudget)
	}
}

// TestWaitStartedStillPassesBootErrThrough pins the other half of the
// contract the bounded select must not break: a boot that DID signal (with an
// error) delivers that error verbatim, not the timeout error and not nil.
func TestWaitStartedStillPassesBootErrThrough(t *testing.T) {
	old := bootWaitBudget
	bootWaitBudget = 5 * time.Second
	t.Cleanup(func() { bootWaitBudget = old })

	s := newSTAThread(nil)
	s.mu.Lock()
	s.err = fmt.Errorf("ball: boot failed on purpose")
	s.mu.Unlock()
	close(s.started)

	if err := s.waitStarted(); err == nil || !strings.Contains(err.Error(), "boot failed on purpose") {
		t.Fatalf("waitStarted lost the boot error through the bounded select: %v", err)
	}
}
