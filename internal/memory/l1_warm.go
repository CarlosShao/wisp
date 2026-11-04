package memory

// Ticket 29 AC#1: the host seam between a live session's Warm window and the
// L1 extraction chain (l1_extract.go).
//
// THE PROBLEM THIS FILE SOLVES. The session scope calls WarmWindow once
// each time a session lands in Warm (internal/session/scope.go), and the Warm
// window is the ONLY idle time the host has inside a live session. Ticket 29's
// L1 extraction has to run there - AFTER the response, never on the response
// path - so the two halves are literally two different goroutines:
//
//	Record(conv)  the RESPONSE PATH. The host hands over the finished task's
//	              conversation right after the loop returned. It copies turns
//	              into the pending slot under a mutex and returns. No C5 call,
//	              no database write, no waiting: there is nothing on this line
//	              that can delay a token. (The latency test asserts exactly
//	              this, and it asserts it as a COUNT rather than a stopwatch:
//	              the provider never sees a request while the response path
//	              runs.)
//	Fire()        THE WARM HOOK BODY. The host assigns this function to
//	              session.WarmWindow. It takes the pending batch, hands it to
//	              the "memory-extract" roster goroutine, and returns - the hook
//	              never waits for the work (the session scope's own contract:
//	              "the scope never runs the work itself and never waits for it,
//	              D38b: the work has an owner or it does not run").
//
// WHY THE PENDING SLOT IS ONE TASK, NOT A QUEUE. The resident leg runs one task
// at a time (the task entry refuses a second submission while one is in
// flight), so a queue would only collect work the next Warm window could not
// reach. A newer Record that arrives before the window fires replaces the older
// one and is counted as Replaced, so the loss is a reading and not a silent
// drop.
//
// FAILURE POSTURE. Ticket 29: "extraction failures -> silent skip + debug log
// (never disturb user)". Every failure path here logs at Debug and books a
// counter; nothing is retried, nothing surfaces to the user, and the batch that
// failed left the profile untouched (the store commit is atomic).

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// The goroutine name comes from the frozen D38b roster
// (observe.TemporaryNames already lists memory-extract); the owner is this
// package. Using the roster name - rather than an invented one - is what keeps
// RosterReport able to classify the spawn.
const (
	l1ExtractGoroutine = "memory-extract"
	l1ExtractOwner     = "memory"
)

// l1ExtractTimeout bounds one extraction call. It is a timer deadline on the
// spawned goroutine's context, never a wall-clock subtraction (D42#9).
const l1ExtractTimeout = 30 * time.Second

// L1WarmStats is the reading the host's boot/stop line and the latency test
// quote. Every counter answers a question an operator would otherwise have to
// guess: was anything recorded, did it run, did it fail, was anything replaced.
type L1WarmStats struct {
	Recorded  int // Record calls that stored a task conversation
	Replaced  int // a recorded conversation overwritten before the window fired
	Started   int // batches handed to the memory-extract goroutine
	Succeeded int // batches that extracted and committed
	Failed    int // batches that failed (silent skip, debug log)
	Dropped   int // Fire calls with nothing pending (the ordinary no-op case)
}

// L1WarmWindow is the host seam. Construct it once per process and give the
// SAME value to the two call sites: Record from the response path, Fire as
// session.WarmWindow.
type L1WarmWindow struct {
	pipeline *L1Pipeline
	logger   *slog.Logger
	reg      *observe.Registry
	timeout  time.Duration

	// onDone is an observation seam (tests, and a host that wants the outcome);
	// it runs on the extraction goroutine, after the stats are booked.
	onDone func(ExtractionOutcome, error)

	mu       sync.Mutex
	pending  *GoldenConversation
	inflight bool
	stats    L1WarmStats
	lastErr  string
}

// L1WarmOption adjusts the warm window.
type L1WarmOption func(*L1WarmWindow)

// WithL1WarmLogger sets the debug sink. nil keeps slog.Default().
func WithL1WarmLogger(l *slog.Logger) L1WarmOption {
	return func(w *L1WarmWindow) {
		if l != nil {
			w.logger = l
		}
	}
}

// WithL1WarmRegistry sets the goroutine registry. nil keeps the store's own
// registry (the same one db-writer is spawned through), which is what the host
// already assembled.
func WithL1WarmRegistry(r *observe.Registry) L1WarmOption {
	return func(w *L1WarmWindow) {
		if r != nil {
			w.reg = r
		}
	}
}

// WithL1WarmObserve installs the post-batch observation seam.
func WithL1WarmObserve(fn func(ExtractionOutcome, error)) L1WarmOption {
	return func(w *L1WarmWindow) {
		if fn != nil {
			w.onDone = fn
		}
	}
}

// NewL1WarmWindow composes the seam around an assembled pipeline.
func NewL1WarmWindow(p *L1Pipeline, opts ...L1WarmOption) (*L1WarmWindow, error) {
	if p == nil {
		return nil, errors.New("memory: L1 warm window needs a pipeline")
	}
	w := &L1WarmWindow{
		pipeline: p,
		logger:   slog.Default(),
		reg:      p.store.reg,
		timeout:  l1ExtractTimeout,
	}
	if w.reg == nil {
		w.reg = observe.Default
	}
	if w.logger == nil {
		w.logger = slog.Default()
	}
	for _, o := range opts {
		o(w)
	}
	return w, nil
}

// Record hands one finished task's conversation to the Warm window. It runs on
// the response path and does no extraction work: the LLM call happens only in
// Fire/Drain, which is what makes the first token of the NEXT turn - and the
// whole response path of this one - independent of the extraction's latency.
//
// A conversation with no user turn carries no extraction evidence (only
// assistant turns are never facts about the user), so it is not recorded at
// all: that is not a drop, it is "there was nothing to extract".
func (w *L1WarmWindow) Record(conv GoldenConversation) {
	if w == nil {
		return
	}
	turns := make([]GoldenTurn, 0, len(conv.Turns))
	for _, t := range conv.Turns {
		role := strings.ToLower(strings.TrimSpace(t.Role))
		text := strings.TrimSpace(t.Text)
		if role != "user" && role != "assistant" {
			continue
		}
		if text == "" {
			continue
		}
		turns = append(turns, GoldenTurn{Role: role, Text: text})
	}
	hasUser := false
	for _, t := range turns {
		if t.Role == "user" {
			hasUser = true
			break
		}
	}
	if !hasUser {
		return
	}
	cp := GoldenConversation{ID: conv.ID, Title: conv.Title, Turns: turns}

	w.mu.Lock()
	if w.pending != nil {
		w.stats.Replaced++
	}
	w.pending = &cp
	w.stats.Recorded++
	w.mu.Unlock()
}

// Fire is the session Warm hook body: take the pending batch, spawn the
// "memory-extract" roster goroutine to run it, return. It is a no-op when
// nothing is pending (the ordinary case: the warm entry that follows a turn in
// which no task finished) and when a batch is already in flight.
func (w *L1WarmWindow) Fire() {
	if w == nil {
		return
	}
	w.mu.Lock()
	if w.pending == nil {
		w.stats.Dropped++
		w.mu.Unlock()
		return
	}
	if w.inflight {
		// An extraction is already running; its own window owns the pending
		// slot. Leave this one to the next warm entry rather than racing two
		// batches for one prompt-cache budget.
		w.mu.Unlock()
		return
	}
	w.inflight = true
	w.stats.Started++
	reg := w.reg
	w.mu.Unlock()

	reg.Spawn(l1ExtractGoroutine, l1ExtractOwner, nil, func(ctx context.Context) {
		defer func() {
			w.mu.Lock()
			w.inflight = false
			w.mu.Unlock()
		}()
		cctx, cancel := context.WithTimeout(ctx, w.timeout)
		defer cancel()
		_, _, _ = w.Drain(cctx)
	})
}

// Drain runs the pending batch synchronously, if there is one. It is the body
// the Warm goroutine runs, and the deterministic entry the tests drive so they
// do not have to race the roster. ran reports whether there was work.
func (w *L1WarmWindow) Drain(ctx context.Context) (out ExtractionOutcome, ran bool, err error) {
	if w == nil {
		return ExtractionOutcome{}, false, nil
	}
	w.mu.Lock()
	conv := w.pending
	w.pending = nil
	w.mu.Unlock()
	if conv == nil {
		return ExtractionOutcome{}, false, nil
	}
	out, err = w.pipeline.RunTask(ctx, *conv)
	w.mu.Lock()
	if err != nil {
		w.stats.Failed++
		w.lastErr = err.Error()
	} else {
		w.stats.Succeeded++
	}
	done := w.onDone
	w.mu.Unlock()
	if done != nil {
		done(out, err)
	}
	if err != nil {
		w.logger.Debug("memory: L1 extraction failed (silent skip, profile untouched)",
			"component", "memory",
			"task", conv.ID,
			"error_class", string(classOf(err)),
			"err", err.Error())
		return out, true, err
	}
	w.logger.Debug("memory: L1 extraction committed",
		"component", "memory", "task", conv.ID,
		"upserts", len(out.Upserts), "forgets", len(out.Forgets))
	return out, true, nil
}

// Pending reports whether a batch is waiting for the next Warm window.
func (w *L1WarmWindow) Pending() bool {
	if w == nil {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.pending != nil
}

// Stats returns a snapshot of the counters.
func (w *L1WarmWindow) Stats() L1WarmStats {
	if w == nil {
		return L1WarmStats{}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.stats
}

// LastError is the most recent extraction failure, for diagnostics. Empty while
// no batch has failed.
func (w *L1WarmWindow) LastError() string {
	if w == nil {
		return ""
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lastErr
}

// classOf keeps the debug line's error_class field classified without importing
// the classifier's rules: an observe error carries its own class, anything else
// is internal (the same fallback the seam's error mapping uses).
func classOf(err error) observe.ErrorClass {
	var oe *observe.Error
	if errors.As(err, &oe) {
		return oe.Class
	}
	return observe.ClassInternal
}
