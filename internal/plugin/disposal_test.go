package plugin

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

func testScope(t *testing.T, name string, opts ...Option) (*DisposalScope, *observe.Registry) {
	t.Helper()
	reg := observe.NewRegistry()
	opts = append(opts, WithRegistry(reg), WithMemoryReleaser(NopMemoryReleaser{}))
	s := NewDisposalScope(name, context.Background(), opts...)
	t.Cleanup(func() { s.Dispose() })
	return s, reg
}

// TestDisposeRunsInReverseOrder covers C11: reverse registration order.
func TestDisposeRunsInReverseOrder(t *testing.T) {
	s, reg := testScope(t, "task-reverse")
	var order []string
	s.DeferNamed("A", func() { order = append(order, "A") })
	s.DeferNamed("B", func() { order = append(order, "B") })
	s.DeferNamed("C", func() { order = append(order, "C") })

	res := s.Dispose()
	if len(order) != 3 || order[0] != "C" || order[1] != "B" || order[2] != "A" {
		t.Fatalf("execution order = %v, want [C B A]", order)
	}
	if len(res.Completed) != 3 {
		t.Fatalf("Completed = %v", res.Completed)
	}
	if res.Incomplete || res.TimedOut || len(res.Failed) != 0 {
		t.Fatalf("unexpected failure markers: %+v", res)
	}
	if got := reg.Count(); got != 0 {
		t.Fatalf("registry live = %d after dispose, want 0", got)
	}
}

// TestDisposeIsIdempotent covers C11: Dispose twice runs the actions once.
func TestDisposeIsIdempotent(t *testing.T) {
	s, _ := testScope(t, "task-idempotent")
	var calls sync.Map // name -> count
	record := func(n string) {
		v, _ := calls.LoadOrStore(n, new(int))
		atomicAdd(v.(*int))
	}
	s.DeferNamed("first", func() { record("first") })
	s.DeferNamed("second", func() { record("second") })

	first := s.Dispose()
	second := s.Dispose()

	for _, n := range []string{"first", "second"} {
		v, _ := calls.Load(n)
		if got := *v.(*int); got != 1 {
			t.Fatalf("fn %q ran %d times, want 1", n, got)
		}
	}
	if len(first.Completed) != 2 || len(second.Completed) != 2 {
		t.Fatalf("results: %+v / %+v", first, second)
	}
	// The second call returns the first result unchanged.
	if fmt.Sprint(first) != fmt.Sprint(second) {
		t.Fatalf("second Dispose result differs: %+v vs %+v", second, first)
	}
}

func atomicAdd(p *int) {
	mu.Lock()
	*p++
	mu.Unlock()
}

var mu sync.Mutex

// TestDisposePanicDoesNotBlockOthers covers C11: one fn panicking must not
// block the remaining steps, and the scope still completes.
func TestDisposePanicDoesNotBlockOthers(t *testing.T) {
	s, _ := testScope(t, "task-panic")
	var order []string
	s.DeferNamed("panics", func() { panic("destroy failed: injected") })
	s.DeferNamed("middle", func() { order = append(order, "middle") })
	s.DeferNamed("last", func() { order = append(order, "last") })

	res := s.Dispose()

	if len(order) != 2 || order[0] != "last" || order[1] != "middle" {
		t.Fatalf("steps after the panic did not run: order=%v", order)
	}
	if len(res.Completed) != 2 {
		t.Fatalf("Completed = %v, want the two healthy steps", res.Completed)
	}
	if len(res.Failed) != 1 || res.Failed[0].Name != "panics" || res.Failed[0].Kind != "panic" {
		t.Fatalf("Failed = %+v, want the panicking step recorded once", res.Failed)
	}
	if res.Failed[0].Err == nil {
		t.Fatal("panic step must carry its error")
	}
	// A panic is a failed step, not a timeout: disposal_incomplete stays false.
	if res.Incomplete || res.TimedOut {
		t.Fatalf("panic must not mark incomplete: %+v", res)
	}
	// The scope itself is healthy afterwards (idempotent path usable).
	if s.Incomplete() {
		t.Fatal("Incomplete() = true after panic-only failure")
	}
}

// TestDisposeTimeoutMarksIncomplete covers C11: the 3s total budget (shortened
// for the test) marks timed-out and never-run steps as disposal_incomplete,
// visibly.
func TestDisposeTimeoutMarksIncomplete(t *testing.T) {
	s, reg := testScope(t, "task-timeout", WithTimeout(60*time.Millisecond))
	blocked := make(chan struct{})
	var ran []string
	// Reverse order: "blocks" is registered last, so it runs first and
	// exhausts the budget; "never-runs" must then be marked incomplete
	// without ever running.
	s.DeferNamed("never-runs", func() { ran = append(ran, "never-runs") })
	s.DeferNamed("blocks", func() { <-blocked }) // hangs until released below

	start := time.Now()
	res := s.Dispose()
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Fatalf("Dispose took %v, budget was 60ms", elapsed)
	}
	if !res.TimedOut || !res.Incomplete || !s.Incomplete() {
		t.Fatalf("timeout must mark disposal_incomplete: %+v", res)
	}
	if len(res.Completed) != 0 || len(ran) != 0 {
		t.Fatalf("no step may complete: completed=%v ran=%v", res.Completed, ran)
	}
	var kinds []string
	for _, f := range res.Failed {
		if f.Kind != "incomplete" {
			t.Fatalf("failed step kind = %q, want incomplete", f.Kind)
		}
		kinds = append(kinds, f.Name)
	}
	if len(kinds) != 2 {
		t.Fatalf("want both steps marked incomplete, got %v", kinds)
	}
	// The abandoned worker stays visible in the registry until it returns;
	// release it so the test leaves no goroutine behind.
	close(blocked)
	deadline := time.Now().Add(2 * time.Second)
	for reg.Count() != 0 && time.Now().Before(deadline) {
		time.Sleep(2 * time.Millisecond)
	}
	if got := reg.Count(); got != 0 {
		t.Fatalf("abandoned disposal worker still live: %d", got)
	}
}

// TestDisposeTailOrder audits the mandatory C11 tail order:
// cancel ctx -> wait goroutines -> deferred fns -> FreeOSMemory -> SLO sample.
func TestDisposeTailOrder(t *testing.T) {
	reg := observe.NewRegistry()
	audit := &tailAudit{entries: make(chan string, 16)}
	s := NewDisposalScope("task-tail", context.Background(),
		WithRegistry(reg),
		WithMemoryReleaser(audit),
		WithSLORecorder(audit),
	)
	t.Cleanup(func() { s.Dispose() })

	workerStarted := make(chan struct{})
	s.Go("disposal-worker", func(ctx context.Context) {
		close(workerStarted)
		// Exits only after the scope ctx is cancelled, i.e. during C11 step 2.
		<-ctx.Done()
		audit.append("worker-exited")
	})
	s.DeferNamed("fns", func() { audit.append("deferred-fns") })

	select {
	case <-workerStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("tracked goroutine never started")
	}

	res := s.Dispose()
	got := audit.drain(t)
	want := []string{"worker-exited", "deferred-fns", "free-os-memory", "slo-sample"}
	if len(got) != len(want) {
		t.Fatalf("tail order = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("tail order = %v, want %v", got, want)
		}
	}
	if res.Incomplete {
		t.Fatalf("clean dispose marked incomplete: %+v", res)
	}
	if !audit.recorded || audit.peak <= 0 || audit.tail <= 0 {
		t.Fatalf("SLO sample missing or non-positive: peak=%d tail=%d recorded=%v",
			audit.peak, audit.tail, audit.recorded)
	}
}

type tailAudit struct {
	mu       sync.Mutex
	entries  chan string
	order    []string
	peak     int64
	tail     int64
	recorded bool
}

func (a *tailAudit) append(e string) {
	a.mu.Lock()
	a.order = append(a.order, e)
	a.mu.Unlock()
}

func (a *tailAudit) FreeOSMemory() { a.append("free-os-memory") }

func (a *tailAudit) RecordScopeMemory(scope string, peakBytes, tailBytes int64) {
	a.mu.Lock()
	a.peak, a.tail, a.recorded = peakBytes, tailBytes, true
	a.mu.Unlock()
	a.append("slo-sample")
}

func (a *tailAudit) drain(t *testing.T) []string {
	t.Helper()
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]string(nil), a.order...)
}

// TestDeferAfterDisposeIsVisible covers the late-registration trap: it must
// not silently drop the action.
func TestDeferAfterDisposeIsVisible(t *testing.T) {
	s, _ := testScope(t, "task-late")
	res := s.Dispose()
	if res.Incomplete {
		t.Fatalf("clean dispose marked incomplete: %+v", res)
	}
	if ok := s.Defer(func() {}); ok {
		t.Fatal("Defer after Dispose must report false")
	}
	if !s.Incomplete() {
		t.Fatal("late registration must mark the scope incomplete")
	}
	again := s.Dispose()
	found := false
	for _, f := range again.Failed {
		if f.Kind == "late" {
			found = true
		}
	}
	if !found {
		t.Fatalf("late registration not visible in result: %+v", again)
	}
}

// TestGoGoroutineGetsScopeCtx verifies tracked goroutines see the scope ctx
// and are joined by Dispose.
func TestGoGoroutineGetsScopeCtx(t *testing.T) {
	s, reg := testScope(t, "task-go")
	started := make(chan struct{})
	s.Go("retention-job", func(ctx context.Context) {
		close(started)
		<-ctx.Done()
	})
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("scope goroutine never started")
	}
	if got := reg.Count(); got != 1 {
		t.Fatalf("registry live = %d, want 1", got)
	}
	res := s.Dispose()
	if res.Incomplete {
		t.Fatalf("dispose with joined goroutine marked incomplete: %+v", res)
	}
	if got := reg.Count(); got != 0 {
		t.Fatalf("registry live = %d after dispose, want 0", got)
	}
}
