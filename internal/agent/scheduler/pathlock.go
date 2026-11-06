package scheduler

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// PathLock is C20 (ticket 47): path-level conflict detection across tasks.
//
// Before a path-touching tool call runs, its task registers the call's
// normalized target paths here (Acquire). A path already held by ANOTHER
// task is a conflict: the later task QUEUES with the visible reason
// 「在等 <task> 释放 <path>」 and unblocks when the holder releases (D43
// edges #20 -> #39, both rows already frozen in the transition table). Two
// tasks mutually waiting is a deadlock: the guard gives up with an explicit
// error after the ceiling, never a silent hang.
//
// Normalization is the caller's contract: Acquire receives what C26 already
// canonicalized (the bridge extracts paths through Decl.PathParams and C26
// before this lock ever sees them), and folds pure spelling once more
// defensively (backslashes, case on Windows) so a caller that skipped C26
// still cannot alias two spellings of one file. No filepath call happens
// here - this lock decides nothing about the filesystem, it only makes the
// strings the canonicalizer produced comparable (D22 ban #2 stays out).
//
// The lock starts NO goroutine. Release re-launches the scheduler's queue on
// the releasing goroutine, so a release is just a function return and there
// is no timer to stop (D32). The deadlock ceiling is enforced by the waiting
// task's own goroutine with a monotonic observe.NewTimeout check - the wait
// is stoppable (lock frees / ctx cancels) and arms no timer.
type PathLock struct {
	mu      sync.Mutex
	owners  map[string]string   // canonical path -> task id holding it
	heldBy  map[string][]string // task id -> canonical paths it holds
	waiting map[string]waitRec  // task id -> what it waits for
	ceil    map[string]time.Duration

	s *Scheduler // nil until bindScheduler; the lock is usable standalone

	// launchQueued is the #39 re-launch seam (set by bindScheduler).
	launchQueued func()
}

// waitRec is one task's blocked-on record.
type waitRec struct {
	path  string
	owner string
	since time.Time
}

// PathLockOptions parameterize NewPathLock.
type PathLockOptions struct {
	// DeadlockTimeout is the default wait ceiling for tasks without a
	// per-task one. 0 = DefaultDeadlockTimeout (300s).
	DeadlockTimeout time.Duration
}

// ErrDeadlock is what a timed-out wait answers. It is a sentinel so a host
// (and the deadlock fixture) can errors.Is it.
var ErrDeadlock = errors.New("pathlock: 等待路径释放超时（死锁防护）")

// ErrWaitCancelled is what an Acquire answers when the task's own ctx died
// while it waited. A cancelled task must not grab the lock it was waiting
// for on the way out.
var ErrWaitCancelled = errors.New("pathlock: 等待期间任务被取消")

// NewPathLock builds an empty lock. It starts no goroutine and owns no
// timer; the deadlock ceiling is checked by the waiting task itself.
func NewPathLock(o PathLockOptions) *PathLock {
	d := o.DeadlockTimeout
	if d <= 0 {
		d = DefaultDeadlockTimeout
	}
	return &PathLock{
		owners:  map[string]string{},
		heldBy:  map[string][]string{},
		waiting: map[string]waitRec{},
		ceil:    map[string]time.Duration{},
	}
}

// bindScheduler wires the #39 re-launch seam and a per-scheduler default
// ceiling (Options.DeadlockTimeout). New calls it.
func (p *PathLock) bindScheduler(s *Scheduler, deadlock time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.s = s
	p.launchQueued = s.launchQueued
	if deadlock > 0 {
		p.ceil[""] = deadlock // the "" key is the scheduler-level default
	}
}

// normalizeCanonical folds pure spelling: backslashes to forward slashes,
// and case-folded on Windows (goosIsWindows is set per-GOOS below). This is
// comparability only; no filesystem decision is made (D22 ban #2).
func normalizeCanonical(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	if goosIsWindows {
		p = strings.ToLower(p)
	}
	return p
}

// Acquire blocks until all paths are free of FOREIGN holds, then takes them.
//
// Call shape: the task calls Acquire on its own goroutine right before the
// path-touching tool executes (the composition adapter wires it beside the
// bridge). The ctx is honored - a cancelled task unblocks with
// ErrWaitCancelled instead of holding the queue. Re-asking for a path the
// task already holds is fine (idempotent; a multi-step task re-touches its
// working file). Returns ErrDeadlock when the wait outlives the ceiling.
func (p *PathLock) Acquire(ctx context.Context, taskID string, paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	norm := normalizeAll(paths)
	if len(norm) == 0 {
		return nil
	}

	p.mu.Lock()
	foreign, holder := p.foreignHeld(taskID, norm)
	if !foreign {
		p.take(taskID, norm)
		p.mu.Unlock()
		return nil
	}
	// Record what we wait for BEFORE unlocking: this is the record the
	// scheduler's Rows() reads for the queued note. The deadlock ceiling
	// also starts HERE (monotonic, per Acquire), not before the poll loop.
	first := p.firstForeignPath(taskID, norm)
	p.waiting[taskID] = waitRec{path: first, owner: holder, since: time.Now()}
	p.mu.Unlock()

	tm := observe.NewTimeout(p.waitCeiling(taskID))
	for {
		p.mu.Lock()
		if foreign, _ := p.foreignHeld(taskID, norm); !foreign {
			delete(p.waiting, taskID)
			p.take(taskID, norm)
			p.mu.Unlock()
			return nil
		}
		// Refresh the visible reason: the holder may have changed while we
		// waited (A released, B took the path).
		first := p.firstForeignPath(taskID, norm)
		w := p.waiting[taskID]
		w.path, w.owner = first, p.owners[first]
		p.waiting[taskID] = w
		p.mu.Unlock()

		select {
		case <-ctx.Done():
			p.mu.Lock()
			w := p.waiting[taskID]
			delete(p.waiting, taskID)
			p.mu.Unlock()
			return fmt.Errorf("pathlock: 任务在等 %q（持有者 %s）释放 %s 时被取消: %w",
				taskID, w.owner, w.path, ErrWaitCancelled)
		case <-time.After(pollStep):
		}

		if tm.Expired() {
			p.mu.Lock()
			w := p.waiting[taskID]
			delete(p.waiting, taskID)
			p.mu.Unlock()
			return fmt.Errorf("pathlock: 任务 %q 等待 %q（持有者 %s）释放 %s 超过上限 %v: %w",
				taskID, taskID, w.owner, w.path, p.waitCeiling(taskID), ErrDeadlock)
		}
	}
}

// pollStep is the Acquire slow path's wake-up interval. Same scale as
// observe.Root.Wait's internal poll; short enough that a release unblocks
// the waiter within a scheduling beat, coarse enough to stay invisible in
// CPU terms for a 300s wait.
const pollStep = 2 * time.Millisecond

// foreignHeld reports whether any of paths is held by a task other than
// taskID, and names the first such holder. Caller holds p.mu.
func (p *PathLock) foreignHeld(taskID string, paths []string) (bool, string) {
	for _, c := range paths {
		if o, ok := p.owners[c]; ok && o != taskID {
			return true, o
		}
	}
	return false, ""
}

// firstForeignPath names the first path of paths held by a foreign task.
// Caller holds p.mu; at least one foreign hold must exist (the caller just
// measured one).
func (p *PathLock) firstForeignPath(taskID string, paths []string) string {
	for _, c := range paths {
		if o, ok := p.owners[c]; ok && o != taskID {
			return c
		}
	}
	return ""
}

// take records the holds. Caller holds p.mu. Idempotent per path: a path
// already owned (by this task, or raced by the fast path) keeps its owner.
func (p *PathLock) take(taskID string, paths []string) {
	for _, c := range paths {
		if _, ok := p.owners[c]; ok {
			continue
		}
		p.owners[c] = taskID
		p.heldBy[taskID] = append(p.heldBy[taskID], c)
	}
}

// waitCeiling resolves the deadlock ceiling for one task: its per-task cap
// (SetWaitCeiling), else the scheduler-level default, else the package
// default.
func (p *PathLock) waitCeiling(taskID string) time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	if d, ok := p.ceil[taskID]; ok && d > 0 {
		return d
	}
	if d, ok := p.ceil[""]; ok && d > 0 {
		return d
	}
	return DefaultDeadlockTimeout
}

// SetWaitCeiling caps one task's wait (the deadlock fixture shrinks it so a
// test does not sit 300 seconds). Production never calls it.
func (p *PathLock) SetWaitCeiling(taskID string, d time.Duration) {
	p.mu.Lock()
	p.ceil[taskID] = d
	p.mu.Unlock()
}

// Release drops taskID's holds on paths and re-launches the queue (the #39
// half: a queued task becomes runnable the moment the path frees). Unknown
// paths are ignored; safe to call twice.
func (p *PathLock) Release(taskID string, paths ...string) {
	rel := normalizeAll(paths)
	if len(rel) == 0 {
		return
	}
	p.mu.Lock()
	freed := false
	for _, c := range rel {
		if p.owners[c] == taskID {
			delete(p.owners, c)
			p.heldBy[taskID] = slices.DeleteFunc(p.heldBy[taskID], func(x string) bool { return x == c })
			freed = true
		}
	}
	p.mu.Unlock()
	if freed {
		p.relaunch()
	}
}

// ReleaseAll drops every hold taskID has (and any wait record). The
// scheduler calls it on task completion so a task that died holding a path
// cannot wedge the queue behind it.
func (p *PathLock) ReleaseAll(taskID string) {
	p.mu.Lock()
	paths := slices.Clone(p.heldBy[taskID])
	delete(p.heldBy, taskID)
	for _, c := range paths {
		if p.owners[c] == taskID {
			delete(p.owners, c)
		}
	}
	_, wasWaiting := p.waiting[taskID]
	delete(p.waiting, taskID)
	delete(p.ceil, taskID)
	p.mu.Unlock()
	if len(paths) > 0 || wasWaiting {
		p.relaunch()
	}
}

// relaunch fires the #39 re-launch exactly once per state change.
func (p *PathLock) relaunch() {
	p.mu.Lock()
	fn := p.launchQueued
	p.mu.Unlock()
	if fn != nil {
		fn()
	}
}

// Waiting reports what taskID is blocked on: the holder's task id, the
// canonical path, and ok. The scheduler's Rows() consumes it for the
// 「在等 <task> 释放 <path>」 note.
func (p *PathLock) Waiting(taskID string) (holder, path string, ok bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	w, ok := p.waiting[taskID]
	if !ok {
		return "", "", false
	}
	return w.owner, w.path, true
}

// Held returns the canonical paths taskID currently holds, sorted.
func (p *PathLock) Held(taskID string) []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return slices.Sorted(slices.Values(p.heldBy[taskID]))
}

// HolderOf reports which task holds the path (ok=false when free).
func (p *PathLock) HolderOf(path string) (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	o, ok := p.owners[normalizeCanonical(path)]
	return o, ok
}

// normalizeAll canonicalizes and dedupes, dropping empties.
func normalizeAll(paths []string) []string {
	norm := make([]string, 0, len(paths))
	for _, raw := range paths {
		c := normalizeCanonical(raw)
		if c == "" || slices.Contains(norm, c) {
			continue
		}
		norm = append(norm, c)
	}
	return norm
}
