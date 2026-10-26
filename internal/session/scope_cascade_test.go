package session

// Ticket 28's teardown readings: the D43 #31 cascade's RSS return (AC5) and the
// 100-cycle refcount/goroutine leak proof (AC6).
//
// These two are the only cases in this package that touch the real runtime
// (heap accounting, real time.Timer, real goroutine count). The policy cases in
// scope_test.go run on a fake clock; these deliberately do not, because "the
// RSS comes back" and "nothing is left running" are claims about the process,
// not about the state machine.

import (
	"runtime"
	"runtime/debug"
	"testing"
	"time"
)

// rssSampleStep / rssSampleBudget are the sampler's grid. AC5 asks for a return
// WITHIN 10 s; a bounded sample COUNT (200 x 50 ms = 10 s) expresses that
// budget without any code deciding "has it been 10 s yet" from a wall-clock
// delta - the loop simply runs out.
const (
	rssSampleStep   = 50 * time.Millisecond
	rssSampleBudget = 200 // rssSampleStep * budget = 10 s
)

// heapEngine is an engine whose whole body is a payload: it exists so the
// cascade's "models unload" step frees something a sampler can see, and so the
// FreeOSMemory step has something to return to the OS.
type heapEngine struct {
	kind    EngineKind
	payload []byte
	closed  bool
	onClose func()
}

func (e *heapEngine) Close() error {
	if e.onClose != nil {
		e.onClose()
	}
	e.payload = nil
	e.closed = true
	return nil
}

// TestTicket28DisposeReturnsTheHeapWithinTenSeconds is AC5's second half: after
// the cascade the process's own accounting comes back down (the sampler), and it
// does so inside the 10 s budget. The order half of AC5 is asserted in
// scope_test.go (assertCascadeOrder) because it needs the seam log, not the heap.
func TestTicket28DisposeReturnsTheHeapWithinTenSeconds(t *testing.T) {
	runtime.GC()
	debug.FreeOSMemory()
	base := heapInUse()

	held := make([]byte, 48<<20)
	for i := range held {
		held[i] = byte(i)
	}
	eng := &heapEngine{kind: EngineASR, payload: held}
	held = nil

	s, err := NewScope(ScopeOptions{
		Loader: func(kind EngineKind) (Engine, error) { return eng, nil },
		Hooks: Hooks{
			FreeOSMemory: debug.FreeOSMemory,
			Logf:         func(string, ...any) {},
		},
	})
	if err != nil {
		t.Fatalf("NewScope: %v", err)
	}
	if _, err := s.Acquire(EngineASR); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	during := heapInUse()
	if during < base+(24<<20) {
		t.Fatalf("the 48 MiB payload is not visible to the sampler (base=%d during=%d); the reading would prove nothing", base, during)
	}

	start := time.Now()
	if err := s.Dispose("menu.end"); err != nil {
		t.Fatalf("dispose: %v", err)
	}
	if !eng.closed {
		t.Fatal("dispose did not close the engine")
	}

	samples := 0
	for ; samples < rssSampleBudget; samples++ {
		runtime.GC()
		if heapInUse() <= base+(8<<20) {
			break
		}
		time.Sleep(rssSampleStep)
	}
	elapsed := time.Since(start)
	if samples >= rssSampleBudget {
		t.Fatalf("heap did not return within %d samples (%s): base=%d now=%d",
			rssSampleBudget, elapsed, base, heapInUse())
	}
	if elapsed > 10*time.Second {
		t.Fatalf("heap returned after %s, want at most 10s (AC5)", elapsed)
	}
	t.Logf("AC5 sampler: base=%d during=%d after=%d samples=%d elapsed=%s",
		base, during, heapInUse(), samples, elapsed)
}

// heapInUse is the process-side number the sampler reads. HeapAlloc rather than
// Sys, because Sys keeps what the runtime RESERVED and would report a return
// that did not happen.
func heapInUse() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.HeapAlloc
}

// TestTicket28HundredWakeSleepCyclesLeaveNothing is AC6: 100 wake -> warm -> wake
// -> dispose cycles, on the REAL clock (so the timers are real time.Timer values
// that must be stopped), ending with zero refs, zero hot engines, zero timers and
// no goroutine growth.
func TestTicket28HundredWakeSleepCyclesLeaveNothing(t *testing.T) {
	log := &eventLog{}
	stub := newEngineStub(log)
	base := settledGoroutines()

	for i := 0; i < 100; i++ {
		s, err := NewScope(ScopeOptions{Loader: stub.load, Hooks: Hooks{
			Logf:           func(string, ...any) {},
			MicClose:       func() {},
			PanelHide:      func() {},
			FreeOSMemory:   func() {},
			GrantsExpire:   func() {},
			PrivacyConfirm: func() bool { return true },
		}})
		if err != nil {
			t.Fatalf("cycle %d: NewScope: %v", i, err)
		}
		for _, k := range []EngineKind{EngineVAD, EngineASR} {
			if _, err := s.Acquire(k); err != nil {
				t.Fatalf("cycle %d: acquire %s: %v", i, k, err)
			}
			if err := s.Release(k); err != nil {
				t.Fatalf("cycle %d: release %s: %v", i, k, err)
			}
		}
		if err := s.EnterWarm("turn.end"); err != nil {
			t.Fatalf("cycle %d: enter warm: %v", i, err)
		}
		if err := s.Resume(); err != nil {
			t.Fatalf("cycle %d: resume: %v", i, err)
		}
		for _, k := range []EngineKind{EngineVAD, EngineASR} {
			if _, err := s.Acquire(k); err != nil {
				t.Fatalf("cycle %d: warm acquire %s: %v", i, k, err)
			}
			if err := s.Release(k); err != nil {
				t.Fatalf("cycle %d: warm release %s: %v", i, k, err)
			}
		}
		if err := s.Dispose("cycle.end"); err != nil {
			t.Fatalf("cycle %d: dispose: %v", i, err)
		}
		st := s.Stats()
		if st.Refs != 0 || st.EnginesHot != 0 {
			t.Fatalf("cycle %d: residual state refs=%d hot=%d, want 0/0 after dispose", i, st.Refs, st.EnginesHot)
		}
		if st.LoadsAfterWarm != 0 {
			t.Fatalf("cycle %d: LoadsAfterWarm=%d, want 0", i, st.LoadsAfterWarm)
		}
		if st.WarmResumes != 1 {
			t.Fatalf("cycle %d: WarmResumes=%d, want 1", i, st.WarmResumes)
		}
		if st.Disposals != 1 {
			t.Fatalf("cycle %d: Disposals=%d, want 1", i, st.Disposals)
		}
		if got, want := stub.totalLoads(), 2*(i+1); got != want {
			t.Fatalf("cycle %d: loader calls=%d, want %d (two cold loads per cycle, none on the wake)", i, got, want)
		}
	}

	// Every handle the 100 cycles loaded was closed exactly once: 100 closes per
	// kind, not "the last handle looks fine".
	for _, kind := range []EngineKind{EngineVAD, EngineASR} {
		if got, want := stub.closesOf(kind), 100; got != want {
			t.Fatalf("closes[%s] = %d, want %d (every load must be closed exactly once)", kind, got, want)
		}
	}

	after := settledGoroutines()
	if after > base {
		t.Fatalf("goroutines: %d after 100 cycles, %d before, want no growth", after, base)
	}
	t.Logf("AC6: 100 cycles, loader calls=%d, goroutines base=%d after=%d", stub.totalLoads(), base, after)
}

// settledGoroutines is the smallest goroutine count over a bounded number of
// scheduling points. The runtime's own GC workers are goroutines too, so a
// single reading can be one or two high for reasons nobody controls; the minimum
// over a fixed sample count is the number a leak would move, and it needs no
// wall clock to say "wait a bit".
func settledGoroutines() int {
	runtime.GC()
	min := runtime.NumGoroutine()
	for i := 0; i < 200; i++ {
		runtime.Gosched()
		if n := runtime.NumGoroutine(); n < min {
			min = n
		}
	}
	return min
}
