package plugin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// DisposalScope — contract C11 (D39 deepened; SPEC-01 §6).
//
// Semantics: Defer registers destroy actions; Dispose runs them
//   - idempotently (second call returns the first result),
//   - in reverse registration order,
//   - each fn with an independent recover (one failure must not block the rest),
//   - under a total budget of DisposalTimeout (3s); items that exceed the
//     budget - or never run because the budget was exhausted - are marked
//     disposal_incomplete and logged (never silent).
//
// Mandatory tail order (fixed): cancel the derived ctx -> wait the tracked
// goroutines -> run the deferred fns (native session release happens here)
// -> release Go memory to the OS -> record peak/tail memory to the SLO
// sampler. Steps 4 and 5 are injectable (MemoryReleaser / SLORecorder); the
// interfaces are frozen by this ticket, the real SLO sampler lands with
// ticket 08.
//
// Hierarchy: session scope (C31) > task scope > tool scope > plugin scope;
// each level derives its ctx from its parent.

// DisposalTimeout is the frozen C11 total budget for one Dispose.
const DisposalTimeout = 3 * time.Second

// MemoryReleaser releases Go runtime memory back to the OS (C11 step 4).
// Without it the D32 "RSS back within limits in 10s" target is unreachable
// because the Go runtime does not return memory eagerly.
type MemoryReleaser interface {
	FreeOSMemory()
}

// DefaultMemoryReleaser is the C11 step-4 implementation (debug.FreeOSMemory).
type DefaultMemoryReleaser struct{}

// FreeOSMemory implements MemoryReleaser.
func (DefaultMemoryReleaser) FreeOSMemory() { debug.FreeOSMemory() }

// NopMemoryReleaser is the injectable empty implementation (tests, boot
// sequences where the real release is wired elsewhere).
type NopMemoryReleaser struct{}

// FreeOSMemory implements MemoryReleaser as a no-op.
func (NopMemoryReleaser) FreeOSMemory() {}

// SLORecorder receives per-scope peak/tail memory samples (C11 step 5). The
// real sampler and its ledger are ticket 08; until then NopSLORecorder is
// injected by default and the interface is the contract.
type SLORecorder interface {
	RecordScopeMemory(scope string, peakBytes, tailBytes int64)
}

// NopSLORecorder is the injectable empty implementation.
type NopSLORecorder struct{}

// RecordScopeMemory implements SLORecorder as a no-op.
func (NopSLORecorder) RecordScopeMemory(string, int64, int64) {}

// DisposeError describes one failed disposal step.
type DisposeError struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"` // "panic" | "incomplete" | "late"
	Err   error  `json:"-"`
	Stack string `json:"stack,omitempty"` // set for Kind "panic"
}

func (e DisposeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("disposal step %q (%s): %v", e.Name, e.Kind, e.Err)
	}
	return fmt.Sprintf("disposal step %q (%s)", e.Name, e.Kind)
}

// DisposeResult is what one Dispose pass produced. It is returned again
// (unchanged) on subsequent idempotent calls.
type DisposeResult struct {
	Scope      string         `json:"scope"`
	Completed  []string       `json:"completed"`
	Failed     []DisposeError `json:"failed,omitempty"`
	Incomplete bool           `json:"disposal_incomplete"` // any step timed out, never ran, or registered too late
	TimedOut   bool           `json:"timed_out"`
}

// Option configures a DisposalScope.
type Option func(*DisposalScope)

// WithTimeout overrides the total Dispose budget (frozen default: 3s). For
// tests and fast paths only; production scopes keep DisposalTimeout.
func WithTimeout(d time.Duration) Option {
	return func(s *DisposalScope) { s.timeout = d }
}

// WithRegistry spawns disposal workers through a specific registry instead
// of observe.Default (tests, isolated subsystems).
func WithRegistry(r *observe.Registry) Option {
	return func(s *DisposalScope) { s.reg = r }
}

// WithMemoryReleaser replaces the C11 step-4 implementation.
func WithMemoryReleaser(mr MemoryReleaser) Option {
	return func(s *DisposalScope) { s.releaser = mr }
}

// WithSLORecorder replaces the C11 step-5 sampler.
func WithSLORecorder(sr SLORecorder) Option {
	return func(s *DisposalScope) { s.slo = sr }
}

type deferredFn struct {
	name string
	fn   func()
}

// DisposalScope is the C11 deterministic-disposal primitive. Zero value is
// not usable; use NewDisposalScope.
type DisposalScope struct {
	name   string
	ctx    context.Context
	cancel context.CancelFunc

	mu         sync.Mutex
	fns        []deferredFn
	steps      int
	disposed   bool
	result     DisposeResult
	incomplete bool

	wgMu      sync.Mutex
	wgPending int

	peak int64 // sampled heap bytes, high-water mark

	timeout  time.Duration
	reg      *observe.Registry
	releaser MemoryReleaser
	slo      SLORecorder
}

// NewDisposalScope creates a scope derived from parent (nil = detached). The
// scope name labels logs and SLO samples (e.g. "session", "task-<id>",
// "tool", "plugin:<id>").
func NewDisposalScope(name string, parent context.Context, opts ...Option) *DisposalScope {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	s := &DisposalScope{
		name:     name,
		ctx:      ctx,
		cancel:   cancel,
		timeout:  DisposalTimeout,
		reg:      observe.Default,
		releaser: DefaultMemoryReleaser{},
		slo:      NopSLORecorder{},
	}
	s.peak = heapAllocBytes()
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Ctx returns the scope's derived context; it is cancelled when Dispose
// starts (C11 step 1) so tracked goroutines can wind down.
func (s *DisposalScope) Ctx() context.Context { return s.ctx }

// Name returns the scope label.
func (s *DisposalScope) Name() string { return s.name }

// Defer registers a destroy action (auto-named step-<n>). Returns false if
// the scope was already disposed: the action will NOT run and the scope is
// marked disposal_incomplete (visible, never silent).
func (s *DisposalScope) Defer(fn func()) bool {
	return s.DeferNamed("", fn)
}

// DeferNamed registers a named destroy action. See Defer.
func (s *DisposalScope) DeferNamed(name string, fn func()) bool {
	if fn == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.disposed {
		n := name
		if n == "" {
			n = fmt.Sprintf("step-%d", s.steps+1)
		}
		s.incomplete = true
		s.result.Incomplete = true
		s.result.Failed = append(s.result.Failed, DisposeError{
			Name: n, Kind: "late",
			Err: errors.New("registered after Dispose; the action did not run"),
		})
		slog.Error("disposal_incomplete: Defer called after Dispose",
			"scope", s.name, "step", n)
		return false
	}
	s.steps++
	if name == "" {
		name = fmt.Sprintf("step-%d", s.steps)
	}
	s.fns = append(s.fns, deferredFn{name: name, fn: fn})
	s.peak = max64(s.peak, heapAllocBytes())
	return true
}

// Go spawns a tracked goroutine through the registry: it runs under the
// scope ctx and is joined by the C11 step-2 wait during Dispose. The name
// must be a D38 roster name (temporary category in practice). Panics inside
// fn are handled by the registry boundary (internal error + log), not by the
// scope; the scope's own per-fn recover contract covers Defer actions.
func (s *DisposalScope) Go(name string, fn func(ctx context.Context)) *observe.Handle {
	s.wgMu.Lock()
	s.wgPending++
	s.wgMu.Unlock()
	// Throwaway root: the registry needs a cancelable root for its panic
	// bookkeeping; the scope ctx is the actual cancellation channel.
	root := observe.NewRootFrom(s.ctx, s.name+"/go")
	return s.reg.Spawn(name, "scope:"+s.name, root, func(ctx context.Context) {
		defer func() {
			s.wgMu.Lock()
			s.wgPending--
			s.wgMu.Unlock()
		}()
		fn(ctx)
	})
}

// Incomplete reports whether disposal did not fully succeed (timeouts, steps
// that never ran, late registrations, panics are reported via result too).
func (s *DisposalScope) Incomplete() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.incomplete
}

// Dispose runs the C11 sequence. Idempotent: later calls return the first
// result without re-running anything.
func (s *DisposalScope) Dispose() DisposeResult {
	s.mu.Lock()
	if s.disposed {
		res := s.result
		s.mu.Unlock()
		return res
	}
	s.disposed = true
	fns := make([]deferredFn, len(s.fns))
	copy(fns, s.fns)
	s.fns = nil
	s.mu.Unlock()

	res := DisposeResult{Scope: s.name}
	total := observe.NewTimeout(s.timeout)

	// C11 step 1: cancel the derived ctx.
	s.cancel()

	// C11 step 2: wait tracked goroutines within the shared budget.
	if p := s.waitPending(total); p > 0 {
		res.TimedOut = true
		res.Incomplete = true
		res.Failed = append(res.Failed, DisposeError{
			Name: "wait-goroutines", Kind: "incomplete",
			Err: fmt.Errorf("%d goroutine(s) still running after disposal budget", p),
		})
	}

	// C11 step 3: deferred fns, reverse order, per-fn recover (via the
	// registry boundary), shared total budget.
	for i := len(fns) - 1; i >= 0; i-- {
		dfn := fns[i]
		if total.Expired() {
			res.TimedOut = true
			res.Incomplete = true
			res.Failed = append(res.Failed, DisposeError{
				Name: dfn.name, Kind: "incomplete",
				Err: errors.New("disposal budget exhausted before this step ran"),
			})
			continue
		}
		remaining := total.Remaining()
		h := s.reg.Spawn("disposal-worker", "scope:"+s.name, nil, func(context.Context) {
			dfn.fn()
		})
		select {
		case <-h.Done():
			if err := h.Err(); err != nil {
				// The registry boundary recovered a panic in this fn; the
				// remaining steps still run (C11: one failure must not
				// block the rest).
				res.Failed = append(res.Failed, DisposeError{
					Name: dfn.name, Kind: "panic", Err: err,
					Stack: panicStack(err),
				})
			} else {
				res.Completed = append(res.Completed, dfn.name)
			}
		case <-time.After(remaining):
			// The worker goroutine cannot be killed; it stays visible in the
			// registry until it returns on its own.
			res.TimedOut = true
			res.Incomplete = true
			res.Failed = append(res.Failed, DisposeError{
				Name: dfn.name, Kind: "incomplete",
				Err: fmt.Errorf("timed out after %v (total budget %v)", remaining, s.timeout),
			})
		}
	}

	// C11 step 4: release Go memory to the OS.
	s.releaser.FreeOSMemory()

	// C11 step 5: record peak/tail memory to the SLO sampler.
	tail := heapAllocBytes()
	s.slo.RecordScopeMemory(s.name, s.peakBytes(), tail)

	for _, f := range res.Failed {
		slog.Error("disposal step failed", "scope", s.name, "step", f.Name, "kind", f.Kind, "err", f.Err)
	}
	slog.Info("scope disposed",
		"scope", s.name,
		"completed", len(res.Completed),
		"failed", len(res.Failed),
		"disposal_incomplete", res.Incomplete,
		"timed_out", res.TimedOut,
	)

	s.mu.Lock()
	s.incomplete = s.incomplete || res.Incomplete
	s.result = res
	s.mu.Unlock()
	return res
}

func (s *DisposalScope) waitPending(total observe.Timeout) int {
	for {
		s.wgMu.Lock()
		p := s.wgPending
		s.wgMu.Unlock()
		if p == 0 {
			return 0
		}
		if total.Expired() {
			return p
		}
		time.Sleep(2 * time.Millisecond)
	}
}

func (s *DisposalScope) peakBytes() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.peak
}

func heapAllocBytes() int64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return int64(ms.HeapAlloc)
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// panicStack digs the stack out of the registry's internal error detail when
// available; the full stack already went to the panic sink.
func panicStack(err error) string {
	var oe *observe.Error
	if errors.As(err, &oe) {
		return oe.Detail
	}
	return ""
}
