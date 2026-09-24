package observe

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestNoopTaskReturnsToBaseline is the SPEC-01 §7 registry assertion: after
// boot plus one no-op task, the goroutine count returns to the resident
// baseline with a tolerance of 1.
//
// Ticket 136 AC#11 (was an ever-green-by-luck read): Registry.Spawn registers
// the key synchronously, but the two legs whose bodies were empty returned
// immediately, so their deferred unregistration could land BEFORE the mid-task
// reads. Measured on the pre-fix tree: 33 of 500 replays of that read window
// saw fewer than 3 per-task legs (tool-exec-noop missing in 31 of the 33,
// approval-waiter in 9, agent-task-noop in 0). The missing edge is now waited
// for with a judgment (each leg reports entering its goroutine and stays
// inside it until the count has been read) under a monotonic timeout bound;
// no sleep is used to paper over the window and no threshold moved.
func TestNoopTaskReturnsToBaseline(t *testing.T) {
	reg := NewRegistry()
	before := runtime.NumGoroutine()

	root := NewRoot("task-noop")
	taskRan := make(chan struct{})
	entered := make(chan string, 3)
	release := make(chan struct{})
	var stopOnce sync.Once
	stop := func() { stopOnce.Do(func() { close(release) }) }
	defer stop() // failure paths must not leave the 3 legs holding a goroutine

	// hold waits until the test releases it, so a leg that has announced
	// itself cannot retire during the counting window.
	hold := func(name string) func(context.Context) {
		return func(ctx context.Context) {
			entered <- name
			<-release
		}
	}
	reg.Spawn("agent-task-noop", "test", root, func(ctx context.Context) {
		entered <- "agent-task-noop"
		time.Sleep(5 * time.Millisecond)
		close(taskRan)
		<-release
	})
	reg.Spawn("tool-exec-noop", "test", root, hold("tool-exec-noop"))
	reg.Spawn("approval-waiter", "test", root, hold("approval-waiter"))

	// Poll to the condition instead of racing it: all 3 legs have entered and
	// the registry counts them. Bounded; on timeout the reads below report
	// exactly what was seen rather than a bare count.
	tm := NewTimeout(2 * time.Second)
	seen := make([]string, 0, 3)
	for (len(seen) < 3 || reg.Count() != 3) && !tm.Expired() {
		select {
		case n := <-entered:
			seen = append(seen, n)
		case <-time.After(2 * time.Millisecond):
		}
	}
	if len(seen) != 3 {
		t.Fatalf("mid-task legs entered = %v after %v, want all 3", seen, tm.Budget())
	}

	// While the task runs, the registry must see its 3 per-task goroutines.
	if got := reg.Count(); got != 3 {
		t.Fatalf("live count mid-task = %d, want 3 (live=%v)", got, reg.Snapshot())
	}
	rep := reg.RosterReport()
	if rep.PerTask != 3 {
		t.Fatalf("PerTask mid-task = %d, want 3 (live=%v)", rep.PerTask, reg.Snapshot())
	}
	select {
	case <-taskRan:
	case <-time.After(2 * time.Second):
		t.Fatal("no-op task fn never ran")
	}

	// Task completion = all derived goroutines drained (D38c). The held legs
	// are released first; the join below is what proves they drained.
	stop()
	if p := root.Wait(3 * time.Second); p != 0 {
		t.Fatalf("root still has %d pending after task completion", p)
	}

	// Count returns to baseline (+1 tolerance for scheduler lag).
	deadline := time.Now().Add(3 * time.Second)
	for runtime.NumGoroutine() > before+1 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if after := runtime.NumGoroutine(); after > before+1 {
		t.Fatalf("goroutine count did not return to baseline: before=%d after=%d", before, after)
	}
	if got := reg.Count(); got != 0 {
		t.Fatalf("registry live count = %d after task, want 0", got)
	}
	rep = reg.RosterReport()
	if len(rep.Unknown) != 0 {
		t.Fatalf("unexpected roster entries: %v", rep.Unknown)
	}
	if rep.ResidentOverBaseline {
		t.Fatalf("resident baseline exceeded: %+v", rep)
	}
}

// TestPanicInWorkerSurvivesAndCancelsRoot covers D37b/§14.10: a panicking
// worker is converted to an internal error with stack, the owning root ctx is
// cancelled, siblings exit, and the process survives.
func TestPanicInWorkerSurvivesAndCancelsRoot(t *testing.T) {
	reg := NewRegistry()
	events := make(chan PanicEvent, 8)
	reg.SetPanicSink(func(ev PanicEvent) { events <- ev })
	defer reg.SetPanicSink(nil)

	panicsBefore := reg.PanicCount()
	root := NewRoot("task-panic")
	siblingDone := make(chan struct{})
	reg.Spawn("tool-exec-sibling", "test", root, func(ctx context.Context) {
		<-ctx.Done()
		close(siblingDone)
	})
	boom := reg.Spawn("tool-exec-boom", "test", root, func(ctx context.Context) {
		panic("boom: injected test panic")
	})

	// Root ctx cancelled by the recover boundary.
	select {
	case <-root.Ctx.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("root ctx was not cancelled after worker panic")
	}

	// Sibling exits through the cancelled root ctx; process stays alive.
	select {
	case <-siblingDone:
	case <-time.After(2 * time.Second):
		t.Fatal("sibling goroutine did not exit after root cancellation")
	}

	// The handle reports the panic as a classified internal error.
	err := boom.Err()
	if got := ClassOfOrInternal(err); got != ClassInternal {
		t.Fatalf("panic error class = %q, want internal", got)
	}
	var oe *Error
	if !errors.As(err, &oe) || !errors.As(err, &oe) {
		t.Fatalf("panic error not an observe *Error: %v", err)
	}

	// Structured record with name, root and stack reached the sink.
	select {
	case ev := <-events:
		if ev.Goroutine != "tool-exec-boom" {
			t.Errorf("event goroutine = %q", ev.Goroutine)
		}
		if ev.RootID != "task-panic" {
			t.Errorf("event root = %q", ev.RootID)
		}
		if ev.Class != ClassInternal {
			t.Errorf("event class = %q, want internal", ev.Class)
		}
		if len(ev.Stack) == 0 {
			t.Error("event stack is empty")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no panic event reached the sink")
	}

	if got := reg.PanicCount(); got != panicsBefore+1 {
		t.Fatalf("PanicCount = %d, want %d", got, panicsBefore+1)
	}
	if p := root.Wait(2 * time.Second); p != 0 {
		t.Fatalf("root pending after panic = %d, want 0", p)
	}
	if got := reg.Count(); got != 0 {
		t.Fatalf("registry live count = %d after panic, want 0", got)
	}
}

// TestUnknownNameIsLeakSymptom verifies roster verification flags goroutines
// outside the D38b roster (what the watchdog will alert on).
func TestUnknownNameIsLeakSymptom(t *testing.T) {
	reg := NewRegistry()
	release := make(chan struct{})
	mystery := reg.Spawn("mystery-worker", "test", nil, func(ctx context.Context) {
		<-release
	})

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		rep := reg.RosterReport()
		if len(rep.Unknown) == 1 && rep.Unknown[0] == "mystery-worker" && rep.Temporary == 0 {
			break
		}
		time.Sleep(2 * time.Millisecond)
	}
	rep := reg.RosterReport()
	if len(rep.Unknown) != 1 || rep.Unknown[0] != "mystery-worker" {
		t.Fatalf("RosterReport.Unknown = %v, want [mystery-worker]", rep.Unknown)
	}

	close(release)
	select {
	case <-mystery.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("mystery goroutine did not exit after release")
	}
	if got := reg.Count(); got != 0 {
		t.Fatalf("registry live count = %d after release, want 0", got)
	}
}

// TestRootWaitTimeout verifies the bounded join (monotonic timeout).
func TestRootWaitTimeout(t *testing.T) {
	reg := NewRegistry()
	root := NewRoot("task-blocked")
	release := make(chan struct{})
	reg.Spawn("tool-exec-blocked", "test", root, func(ctx context.Context) {
		<-release
	})

	start := time.Now()
	if p := root.Wait(60 * time.Millisecond); p != 1 {
		t.Fatalf("Wait(60ms) pending = %d, want 1", p)
	}
	if el := time.Since(start); el < 50*time.Millisecond {
		t.Fatalf("Wait returned too early: %v", el)
	}
	close(release)
	if p := root.Wait(2 * time.Second); p != 0 {
		t.Fatalf("Wait after release pending = %d, want 0", p)
	}
}

// TestClassifyGoroutine pins the roster classification.
func TestClassifyGoroutine(t *testing.T) {
	cases := map[string]GoroutineCategory{
		"ui-sta":           CategoryResident,
		"audio-capture":    CategoryResident,
		"hotkey-listener":  CategoryResident,
		"db-writer":        CategoryResident,
		"watchdog":         CategoryResident,
		"log-flusher":      CategoryResident,
		"kws-infer":        CategoryOnDemand,
		"agent-task-7f3a":  CategoryPerTask,
		"tool-exec-42":     CategoryPerTask,
		"approval-waiter":  CategoryPerTask,
		"asr-infer":        CategoryTemporary,
		"tts-infer":        CategoryTemporary,
		"panel-host":       CategoryTemporary,
		"retention-job":    CategoryTemporary,
		"model-downloader": CategoryTemporary,
		"memory-extract":   CategoryTemporary,
		"disposal-worker":  CategoryTemporary,
		"mystery-worker":   CategoryUnknown,
		"":                 CategoryUnknown,
	}
	for name, want := range cases {
		if got := ClassifyGoroutine(name); got != want {
			t.Errorf("ClassifyGoroutine(%q) = %q, want %q", name, got, want)
		}
	}
	if ResidentBaseline != len(ResidentNames) {
		t.Fatalf("ResidentBaseline = %d but ResidentNames has %d entries", ResidentBaseline, len(ResidentNames))
	}
}
