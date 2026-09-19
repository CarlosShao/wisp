package observe

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"time"
)

// Named-goroutine registry (D38b, SPEC-01 §4.2, §14.10).
//
// Every goroutine in Wisp must be spawned through Registry.Spawn: it carries a
// name (matched against the frozen D38 roster), an owner, and an optional task
// Root. The spawn site installs the mandatory recover boundary: a recovered
// panic becomes a structured "internal" error record (D37b) with the goroutine
// name and stack, and cancels the owning task's root ctx. Panics do not cross
// goroutines; the process survives.
//
// This file contains the ONLY sanctioned `go` statement of the codebase;
// production code elsewhere must not start goroutines directly (D22 ban,
// AST scan lands with ticket 08; a grep-level regression test runs meanwhile).

// GoroutineCategory is the D38b roster category of a goroutine name.
type GoroutineCategory string

const (
	CategoryResident  GoroutineCategory = "resident"  // baseline 6, D38b
	CategoryOnDemand  GoroutineCategory = "on_demand" // kws-infer, Armed state only
	CategoryPerTask   GoroutineCategory = "per_task"  // 3 per active task
	CategoryTemporary GoroutineCategory = "temporary" // inside a DisposalScope
	CategoryUnknown   GoroutineCategory = "unknown"   // NOT in the roster: leak symptom
)

// ResidentBaseline is the frozen resident goroutine budget (D38b: ui-sta,
// audio-capture, hotkey-listener, db-writer, watchdog, log-flusher).
const ResidentBaseline = 6

// ResidentNames are the 6 resident roster names (D38b), in table order.
var ResidentNames = []string{
	"ui-sta", "audio-capture", "hotkey-listener", "db-writer", "watchdog", "log-flusher",
}

// OnDemandNames are state-gated goroutines (D38b).
var OnDemandNames = []string{"kws-infer"}

// PerTaskPrefixes are the per-task roster entries (D38b): agent-task-<id>
// holds the task root ctx; tool-exec-<id> and approval-waiter derive from it.
var PerTaskPrefixes = []string{"agent-task-", "tool-exec-", "approval-waiter"}

// TemporaryNames are DisposalScope-scoped goroutines (D38b). disposal-worker
// is the C11 per-fn disposal runner (ticket 03); it is bookkeeping, not a
// product goroutine, and lives inside whichever scope is disposing.
var TemporaryNames = []string{
	"asr-infer", "tts-infer", "panel-host", "retention-job", "model-downloader",
	"memory-extract", "disposal-worker",
}

// ClassifyGoroutine maps a goroutine name to its roster category. Unknown
// names are leak symptoms and surface in RosterReport.
func ClassifyGoroutine(name string) GoroutineCategory {
	if slices.Contains(ResidentNames, name) {
		return CategoryResident
	}
	if slices.Contains(OnDemandNames, name) {
		return CategoryOnDemand
	}
	for _, p := range PerTaskPrefixes {
		// approval-waiter is the roster name itself (no per-instance id);
		// agent-task-/tool-exec- carry the task id as suffix.
		if strings.HasPrefix(name, p) {
			return CategoryPerTask
		}
	}
	if slices.Contains(TemporaryNames, name) {
		return CategoryTemporary
	}
	return CategoryUnknown
}

// Root is a task (or process) root: one cancelable context plus a join
// counter. All goroutines of a task derive from the same Root (D38c); task
// completion means Root.Wait() drains to zero, which is what keeps RSS
// fall-back (D32) reachable.
type Root struct {
	ID     string
	Ctx    context.Context
	Cancel context.CancelFunc

	mu      sync.Mutex
	pending int
}

// NewRoot creates a process-level root (detached parent).
func NewRoot(id string) *Root {
	return NewRootFrom(context.Background(), id)
}

// NewRootFrom derives a root from a parent ctx (e.g. a session root).
func NewRootFrom(parent context.Context, id string) *Root {
	ctx, cancel := context.WithCancel(parent)
	return &Root{ID: id, Ctx: ctx, Cancel: cancel}
}

func (r *Root) addPending(n int) {
	if r == nil {
		return
	}
	r.mu.Lock()
	r.pending += n
	r.mu.Unlock()
}

// Pending returns the number of goroutines spawned under this root that have
// not finished yet.
func (r *Root) Pending() int {
	if r == nil {
		return 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.pending
}

// Wait waits up to timeout (monotonic; D42#9) for all spawned goroutines of
// this root to finish and returns how many are still pending. timeout <= 0
// returns the current pending count without waiting. The 3s cap used by the
// D38(e) shutdown sequence lives in proc, not here.
func (r *Root) Wait(timeout time.Duration) int {
	if r == nil {
		return 0
	}
	if timeout <= 0 {
		return r.Pending()
	}
	tm := NewTimeout(timeout)
	for {
		if r.Pending() == 0 {
			return 0
		}
		if tm.Expired() {
			return r.Pending()
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// PanicEvent is the structured record emitted when a spawned goroutine
// recovers a panic (D37b: class internal, name + stack, cancel owning root).
type PanicEvent struct {
	Goroutine string     `json:"goroutine"`
	Owner     string     `json:"owner"`
	RootID    string     `json:"root_id"`
	Class     ErrorClass `json:"error_class"` // always ClassInternal
	Recovered string     `json:"recovered"`
	Stack     string     `json:"stack"`
	At        time.Time  `json:"at"` // wall clock, for the persisted record only
}

// PanicSink receives panic records. The default writes a structured slog
// error; ticket 08 replaces it with the JSONL pipeline.
type PanicSink func(PanicEvent)

func defaultPanicSink(ev PanicEvent) {
	slog.Error("goroutine panic recovered",
		"goroutine", ev.Goroutine,
		"owner", ev.Owner,
		"root", ev.RootID,
		"error_class", string(ev.Class),
		"recovered", ev.Recovered,
		"stack", ev.Stack,
		"at", WallTimestampUTC(ev.At),
	)
}

// GoroutineInfo is one live roster entry (aggregated by name+owner).
type GoroutineInfo struct {
	Name     string
	Owner    string
	Category GoroutineCategory
	Count    int
}

// RosterReport compares the live set against the D38b roster.
type RosterReport struct {
	Total                int
	Resident             int
	OnDemand             int
	PerTask              int
	Temporary            int
	Unknown              []string // live names outside the roster
	ResidentOverBaseline bool     // resident count > 6: leak symptom (D38b)
}

// Handle joins a spawned goroutine: Done closes after the fn returns or
// panics; Err returns nil or the panic converted to an internal Error.
type Handle struct {
	done chan struct{}
	errc chan error
}

// Done closes when the goroutine has finished (including recover handling).
func (h *Handle) Done() <-chan struct{} { return h.done }

// Err blocks until the goroutine finishes and returns nil, or the recovered
// panic as a classified internal Error.
func (h *Handle) Err() error { return <-h.errc }

// Registry tracks every goroutine spawned through it.
type Registry struct {
	mu           sync.Mutex
	live         map[gkey]int
	totalStarted uint64
	panicCount   uint64

	sinkMu sync.RWMutex
	sink   PanicSink
}

type gkey struct {
	name  string
	owner string
}

// Default is the process registry; Spawn helpers on it are used by all
// product code so the watchdog sees one picture.
var Default = NewRegistry()

// NewRegistry creates an independent registry (tests, isolated subsystems).
func NewRegistry() *Registry {
	return &Registry{
		live: map[gkey]int{},
		sink: defaultPanicSink,
	}
}

// SetPanicSink replaces the panic record sink (ticket 08 JSONL pipeline).
func (r *Registry) SetPanicSink(s PanicSink) {
	r.sinkMu.Lock()
	defer r.sinkMu.Unlock()
	if s == nil {
		s = defaultPanicSink
	}
	r.sink = s
}

// Spawn starts fn as a named goroutine.
//
//   - name MUST follow the D38 roster (unknown names are recorded as leak
//     symptoms; see ClassifyGoroutine).
//   - owner is the subsystem that started it (e.g. "agent", "test").
//   - root, when non-nil, contributes its ctx to fn, joins the root's pending
//     counter, and gets cancelled if fn panics (D37b).
//   - fn receives the root ctx (or a detached background ctx when root is
//     nil) and must honor cancellation.
//
// The recover boundary is installed here; callers must not recover their own
// panics in a way that bypasses this boundary.
func (r *Registry) Spawn(name, owner string, root *Root, fn func(ctx context.Context)) *Handle {
	h := &Handle{done: make(chan struct{}), errc: make(chan error, 1)}
	if fn == nil {
		h.errc <- New(ClassInternal, "Spawn called with nil fn")
		close(h.done)
		return h
	}
	cat := ClassifyGoroutine(name)
	if cat == CategoryUnknown {
		slog.Warn("goroutine outside the D38 roster (leak symptom)",
			"goroutine", name, "owner", owner)
	}

	r.mu.Lock()
	r.live[gkey{name, owner}]++
	r.totalStarted++
	r.mu.Unlock()
	root.addPending(1)

	go r.run(h, name, owner, root, cat, fn)
	return h
}

// run is the single sanctioned goroutine body of the codebase.
func (r *Registry) run(h *Handle, name, owner string, root *Root, cat GoroutineCategory, fn func(context.Context)) {
	var ctx context.Context
	if root != nil {
		ctx = root.Ctx
	} else {
		ctx = context.Background()
	}
	_ = cat // roster classification is recorded at Spawn/Snapshot time

	defer func() {
		var err error
		if rec := recover(); rec != nil {
			ev := PanicEvent{
				Goroutine: name,
				Owner:     owner,
				Class:     ClassInternal,
				Recovered: fmt.Sprint(rec),
				Stack:     string(debug.Stack()),
				At:        NowWallUTC(),
			}
			if root != nil {
				ev.RootID = root.ID
			}
			r.mu.Lock()
			r.panicCount++
			r.mu.Unlock()
			r.sinkMu.RLock()
			sink := r.sink
			r.sinkMu.RUnlock()
			if sink != nil {
				sink(ev)
			}
			err = &Error{
				Class:  ClassInternal,
				Detail: fmt.Sprintf("panic in goroutine %q (owner %s): %s", name, owner, ev.Recovered),
			}
			// Cancel the owning task root ctx (D37b): a panicked worker must
			// not leave its task half-alive.
			if root != nil && root.Cancel != nil {
				root.Cancel()
			}
		}
		h.errc <- err
		if root != nil {
			root.addPending(-1)
		}
		r.mu.Lock()
		r.live[gkey{name, owner}]--
		if r.live[gkey{name, owner}] <= 0 {
			delete(r.live, gkey{name, owner})
		}
		r.mu.Unlock()
		close(h.done)
	}()

	fn(ctx)
}

// Count returns the number of live goroutines spawned through this registry.
func (r *Registry) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, c := range r.live {
		n += c
	}
	return n
}

// CountByName returns the number of live goroutines with the given name.
func (r *Registry) CountByName(name string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for k, c := range r.live {
		if k.name == name {
			n += c
		}
	}
	return n
}

// TotalStarted returns the cumulative number of spawns (for diagnostics).
func (r *Registry) TotalStarted() uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.totalStarted
}

// PanicCount returns the number of panics recovered since creation.
func (r *Registry) PanicCount() uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.panicCount
}

// Snapshot lists the live goroutines aggregated by (name, owner), sorted for
// stable diagnostics output.
func (r *Registry) Snapshot() []GoroutineInfo {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]GoroutineInfo, 0, len(r.live))
	for k, c := range r.live {
		out = append(out, GoroutineInfo{
			Name:     k.name,
			Owner:    k.owner,
			Category: ClassifyGoroutine(k.name),
			Count:    c,
		})
	}
	slices.SortFunc(out, func(a, b GoroutineInfo) int {
		if a.Name != b.Name {
			return strings.Compare(a.Name, b.Name)
		}
		return strings.Compare(a.Owner, b.Owner)
	})
	return out
}

// RosterReport verifies the live set against the frozen D38b roster.
func (r *Registry) RosterReport() RosterReport {
	snap := r.Snapshot()
	rep := RosterReport{Total: r.Count()}
	for _, gi := range snap {
		switch gi.Category {
		case CategoryResident:
			rep.Resident += gi.Count
		case CategoryOnDemand:
			rep.OnDemand += gi.Count
		case CategoryPerTask:
			rep.PerTask += gi.Count
		case CategoryTemporary:
			rep.Temporary += gi.Count
		default:
			rep.Unknown = append(rep.Unknown, gi.Name)
		}
	}
	rep.ResidentOverBaseline = rep.Resident > ResidentBaseline
	slices.Sort(rep.Unknown)
	return rep
}
