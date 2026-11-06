package scheduler

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// The TaskScheduler (ticket 47, C20/D31/D38c). It owns what PLAN.md's D38(c)
// paragraph assigns to it: FIFO accept, concurrent run, one root ctx per task
// with every task goroutine spawned under it through the D38b registry, and
// completion that REQUIRES the join counter (root Pending) to drain - the
// RSS fall-back of D32 depends on that drain being observable.
//
// What this file deliberately does NOT own:
//
//   - the agent loop itself (internal/agent): the host hands Exec, the
//     scheduler only wraps it in the per-task root;
//   - the panel task list / ball Queued dot (internal/panel, ticket 145's
//     carrier): Rows() is the read port a display layer consumes;
//   - the composition wiring into cmd/wisp (the assembly ticket owns the
//     run.go half; the wave-45 disposition records that follow-up).
//
// D32 discipline: no timers are created here. The one deadline this package
// waits on (the PathLock deadlock ceiling) is a monotonic deadline check
// (observe.NewTimeout) inside the waiting task's own goroutine - stoppable,
// because the moment the lock frees or the task is cancelled the waiter
// leaves the loop. No time.AfterFunc anywhere.

// MaxTaskConcurrency is the scheduler's own FIFO run window. It is NOT a D38
// number (D38d's ceiling of 4 is the per-task tool concurrency cap and lives
// in internal/agent); this one is how many accepted tasks run at once before
// the rest wait in the queue. The default matches the ticket's acceptance
// shape (4 tasks at once); tests raise or shrink it explicitly.
const MaxTaskConcurrency = 4

// DefaultDeadlockTimeout is how long a task may sit blocked on a path hold
// before the deadlock guard gives up with an explicit error. The ticket names
// no number; 300s is the same "how long may a task visibly wait" scale D31
// already fixes for an approval card. The deadlock fixture shrinks it per
// task through PathLock.SetWaitCeiling.
const DefaultDeadlockTimeout = 300 * time.Second

// ErrClosed is what Submit answers once Close shut the door (D38e step 1:
// no new tasks). It is a sentinel, not a wrapped error, so a host can
// errors.Is it.
var ErrClosed = errors.New("scheduler: 门已关（D38e 步骤 1），不再接收新任务")

// ErrUnknownTask is what the row-addressing methods answer for an id this
// scheduler never accepted (or already reaped). Loud, never a silent no-op.
var ErrUnknownTask = errors.New("scheduler: 查无此任务")

// TaskRun is the handle Exec gets. It is the narrow face a task body needs:
// its own id, the task's own Root (spawn every task helper on THIS root so
// the D38b per-task names agent-task-<id> / tool-exec-<id> / approval-waiter
// join the task's join counter - the drain law of D38c is then observable:
// Scheduler.Wait cannot return while any helper is pending), the registry
// to spawn on, and the path lock when one is wired.
type TaskRun struct {
	ID    string
	Root  *observe.Root
	Reg   *observe.Registry
	Locks *PathLock
}

// TaskResult is what a task body reports. AppliedSteps is the D31 ledger: on
// a cancelled task it is the evidence of what already landed, and the
// scheduler files it on the row verbatim - never rendered as a clean undo.
type TaskResult struct {
	Status string // "completed" | "cancelled" | "failed" (host vocabulary)
	Text   string
	// AppliedSteps lists what already landed when the task stopped early.
	AppliedSteps []string
	// Err is the failure detail when Status == "failed".
	Err string
}

// TaskSpec is one accepted task.
type TaskSpec struct {
	Label string
	Exec  func(ctx context.Context, run *TaskRun) TaskResult
}

// TaskRow is one scheduler row as a reader sees it. State carries ONLY the
// D43 vocabulary (statemachine.State, the same discipline ticket 188 put on
// the roster's status dimension): Queued while waiting (slot or path),
// Thinking while running, Acting / AwaitingApproval when the host files them
// through MarkState, Settling when done, Error when failed. Note is free
// text - the 「在等 <task> 释放 <path>」 reason lives there.
type TaskRow struct {
	ID        string
	Label     string
	State     statemachine.State
	Note      string
	HeldPaths []string
}

// task is the scheduler's internal record. Every field is guarded by
// Scheduler.mu; done is the only channel, closed exactly once.
type task struct {
	id    string
	label string
	spec  TaskSpec

	root *observe.Root
	h    *observe.Handle

	done   chan struct{}
	result TaskResult

	state statemachine.State
	note  string
}

// Scheduler is the C20 task scheduler. Create with New; one instance per
// process (the composition root owns the lifetime).
type Scheduler struct {
	reg      *observe.Registry
	locks    *PathLock
	maxRun   int
	deadlock time.Duration

	mu     sync.Mutex
	cond   *sync.Cond
	door   bool // true = open
	queue  []*task
	live   map[string]*task
	order  []string // FIFO accept order, for Rows()
	closed bool

	// parents is the authority join the tools' stop branch sits on: task id
	// -> the parent task id whose task.cancel may stop that row (the host
	// files it through SetParent when it submits a task on behalf of a
	// parent). A task with no entry is a root row, and a root has no
	// model-facing stop verb (the same rule the roster's ParentTaskID
	// empty-means-root reading follows).
	parents map[string]string
}

// Options parameterize New. Zero values fall back to the documented defaults.
type Options struct {
	// Registry spawns every task goroutine. nil = observe.Default (product
	// shape); tests pass their own so roster assertions stay isolated.
	Registry *observe.Registry
	// Locks is the C20 PathLock. nil disables path conflict detection
	// (a caller that never touches paths loses nothing).
	Locks *PathLock
	// MaxConcurrent caps how many accepted tasks run at once. 0 = the
	// documented default (MaxTaskConcurrency); negative falls back to 1.
	MaxConcurrent int
	// DeadlockTimeout is handed to Locks when the lock itself carries no
	// per-task ceiling; 0 leaves the lock's own default.
	DeadlockTimeout time.Duration
}

// New builds an empty, open scheduler. It starts no goroutine: launches
// happen on Submit, on task completion, and on path release, all through
// Registry.Spawn on the calling goroutine, so this package adds zero new
// goroutine shapes to the D38b roster beyond the per-task trio it is named
// for.
func New(o Options) *Scheduler {
	if o.Registry == nil {
		o.Registry = observe.Default
	}
	maxRun := o.MaxConcurrent
	if maxRun == 0 {
		maxRun = MaxTaskConcurrency
	}
	if maxRun < 1 {
		maxRun = 1
	}
	s := &Scheduler{
		reg:      o.Registry,
		locks:    o.Locks,
		maxRun:   maxRun,
		deadlock: o.DeadlockTimeout,
		door:     true,
		live:     map[string]*task{},
		parents:  map[string]string{},
	}
	if s.locks != nil {
		s.locks.bindScheduler(s, s.deadlock)
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

// Submit accepts one task (FIFO) and returns its id. If a run slot is free
// the task launches before Submit returns (the Spawn itself is async); the
// row is visible as Queued from the moment of acceptance either way. After
// Close, Submit answers ErrClosed and accepts nothing.
func (s *Scheduler) Submit(spec TaskSpec) (string, error) {
	if spec.Exec == nil {
		return "", errors.New("scheduler: Exec 为空，拒绝收下无体的任务")
	}
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return "", ErrClosed
	}
	t := &task{
		id:    newTaskID(),
		label: spec.Label,
		spec:  spec,
		done:  make(chan struct{}),
		state: statemachine.StateQueued,
	}
	s.live[t.id] = t
	s.order = append(s.order, t.id)
	s.queue = append(s.queue, t)
	s.mu.Unlock()

	s.launchQueued()
	return t.id, nil
}

// runningCount counts tasks that hold a run slot: not queued, and NOT yet
// terminal (a Settling/Error row has already freed its slot - that is the
// whole point of launchQueued running on the finishing task's goroutine).
// Caller holds s.mu.
func (s *Scheduler) runningCount() int {
	n := 0
	for _, id := range s.order {
		t := s.live[id]
		if t == nil {
			continue
		}
		switch t.state {
		case statemachine.StateQueued, statemachine.StateSettling, statemachine.StateError:
			// no slot held (queued: never started; terminal: slot freed)
		default:
			n++
		}
	}
	return n
}

// launchQueued starts every queued task that fits the run window, in FIFO
// order. It runs on the caller's goroutine (Submit, a finishing task, a path
// release) and never blocks: each launch is one Registry.Spawn.
func (s *Scheduler) launchQueued() {
	for {
		s.mu.Lock()
		if len(s.queue) == 0 || s.runningCount() >= s.maxRun {
			s.mu.Unlock()
			return
		}
		t := s.queue[0]
		s.queue = s.queue[1:]
		if s.live[t.id] != t { // cancelled while queued
			s.mu.Unlock()
			continue
		}
		t.root = observe.NewRootFrom(context.Background(), t.id)
		t.h = s.reg.Spawn("agent-task-"+t.id, "scheduler", t.root, func(c context.Context) {
			defer close(t.done)
			t.result = s.runTask(c, t)
		})
		t.state = statemachine.StateThinking
		s.mu.Unlock()
	}
}

// runTask wraps one Exec: it hands the task its run handle (the task root
// travels ON it, so helpers join the join counter), releases the task's path
// holds on the way out, and settles the row.
//
// The drain law (D38c) is enforced by the JOIN side, not here: the registry
// closes the spawn handle only AFTER it decremented the root's pending
// counter, and Scheduler.Wait joins that handle - so a Wait() returner
// never observes a half-drained root, and an Exec that "returns" while its
// helper still runs keeps the task observably incomplete. (Waiting for the
// root from INSIDE this goroutine would deadlock by construction: the task
// goroutine is itself one pending count.)
func (s *Scheduler) runTask(ctx context.Context, t *task) TaskResult {
	run := &TaskRun{ID: t.id, Root: t.root, Reg: s.reg, Locks: s.locks}
	res := t.spec.Exec(ctx, run)

	// Release every path this task still holds, so a task that dies mid-run
	// cannot wedge the queue behind it (the #39 re-launch is the same call).
	if s.locks != nil {
		s.locks.ReleaseAll(t.id)
	}
	s.mu.Lock()
	t.state = terminalState(res)
	t.note = terminalNote(res)
	s.mu.Unlock()

	// Free the slot: launch the next queued task from this (dying)
	// goroutine, so FIFO keeps moving without any pump goroutine.
	s.launchQueued()
	return res
}

// taskDrainWait is the per-task completion drain bound. proc.TaskWaitTimeout
// (3s) is the process-wide D38e number; a task that has not drained its own
// helpers in that window is abandoned with a recorded note.
const taskDrainWait = 3 * time.Second

// terminalState maps a TaskResult onto the D43 vocabulary. Failed lands on
// Error, everything else on Settling (D43 has no per-task Cancelled name;
// the note says what happened).
func terminalState(res TaskResult) statemachine.State {
	if res.Status == "failed" {
		return statemachine.StateError
	}
	return statemachine.StateSettling
}

// terminalNote builds the row note a reader sees after the task ends: the
// status, and on a stopped task the applied steps verbatim (D31: never a
// clean-undo wording over steps that landed).
func terminalNote(res TaskResult) string {
	var b strings.Builder
	switch res.Status {
	case "completed":
		b.WriteString("已完成")
	case "cancelled":
		b.WriteString("已取消（任务根 ctx）")
	case "failed":
		b.WriteString("失败：" + res.Err)
	default:
		b.WriteString("结束（宿主上报状态 " + res.Status + "）")
	}
	if res.Text != "" {
		b.WriteString("；")
		b.WriteString(res.Text)
	}
	if len(res.AppliedSteps) > 0 {
		b.WriteString("；已生效步骤：")
		for i, step := range res.AppliedSteps {
			if i > 0 {
				b.WriteString("、")
			}
			b.WriteString(step)
		}
		b.WriteString("（未自动回退）")
	}
	return b.String()
}

// Wait blocks until the task has fully finished AND its root's join counter
// has drained (D38c: completion means the WaitGroup drained - the RSS
// fall-back of D32 depends on that being observable). The drain bound is
// taskDrainWait; a root still pending beyond it is recorded on the row as
// abandoned (the D37c shape), never silently forgotten. An unknown id
// answers ErrUnknownTask.
func (s *Scheduler) Wait(id string) (TaskResult, error) {
	s.mu.Lock()
	t := s.live[id]
	s.mu.Unlock()
	if t == nil {
		return TaskResult{}, ErrUnknownTask
	}
	<-t.done
	if t.h != nil {
		<-t.h.Done()
	}
	// The drain: Pending counts every goroutine still joined on the task
	// root. Zero means the WaitGroup drained; anything left after the bound
	// is abandoned and recorded.
	if p := t.root.Wait(taskDrainWait); p > 0 {
		s.mu.Lock()
		t.note = strings.TrimSpace(t.note + "；任务根未排干（还有 " +
			fmt.Sprint(p) + " 个派生协程未退出，已放弃等待并记录）")
		s.mu.Unlock()
	}
	return t.result, nil
}

// Done returns a channel that closes when the task has fully finished (same
// join discipline as Wait). A nil channel comes back for an unknown id, so a
// reader selecting on it blocks instead of spinning.
func (s *Scheduler) Done(id string) <-chan struct{} {
	s.mu.Lock()
	t := s.live[id]
	s.mu.Unlock()
	if t == nil {
		return nil
	}
	return t.done
}

// Cancel stops one task: a queued one is dequeued and settled without ever
// running; a running one gets its root ctx cancelled (D38c: 取消任务 = cancel
// root ctx). The row settles through the normal path and carries whatever
// applied steps the exec reported.
func (s *Scheduler) Cancel(id string) error {
	s.mu.Lock()
	t := s.live[id]
	if t == nil {
		s.mu.Unlock()
		return ErrUnknownTask
	}
	if t.state == statemachine.StateQueued {
		for i, q := range s.queue {
			if q == t {
				s.queue = slices.Delete(s.queue, i, i+1)
				break
			}
		}
		t.state = statemachine.StateSettling
		t.note = "已取消（还在排队，未开始执行）"
		close(t.done)
		s.mu.Unlock()
		return nil
	}
	s.mu.Unlock()
	if t.root != nil && t.root.Cancel != nil {
		t.root.Cancel()
	}
	return nil
}

// Rows returns a snapshot of every row in FIFO accept order. The wait note
// of a path-blocked task comes from the lock itself, so the 「在等 <task>
// 释放 <path>」 reason a reader needs is always the live one.
func (s *Scheduler) Rows() []TaskRow {
	s.mu.Lock()
	out := make([]TaskRow, 0, len(s.order))
	for _, id := range s.order {
		t := s.live[id]
		if t == nil {
			continue
		}
		row := TaskRow{ID: t.id, Label: t.label, State: t.state, Note: t.note}
		if s.locks != nil {
			if owner, path, ok := s.locks.Waiting(id); ok {
				row.State = statemachine.StateQueued
				row.Note = waitingNote(owner, path)
			}
			row.HeldPaths = s.locks.Held(id)
		}
		out = append(out, row)
	}
	s.mu.Unlock()
	return out
}

// Row returns one row snapshot. ok is false for an unknown id.
func (s *Scheduler) Row(id string) (TaskRow, bool) {
	for _, r := range s.Rows() {
		if r.ID == id {
			return r, true
		}
	}
	return TaskRow{}, false
}

// MarkState lets the host file the per-task D43 dimension it owns: Acting
// while tools run, AwaitingApproval while an L2 card is up (Thinking is the
// scheduler's own run-state and is accepted idempotently so a host adapter
// can mark unconditionally). Everything else is refused - the scheduler
// alone owns Queued/Settling/Error, and a host writing over those would be
// forging the #20→#39 transition trace.
func (s *Scheduler) MarkState(id string, st statemachine.State) error {
	switch st {
	case statemachine.StateActing, statemachine.StateAwaitingApproval,
		statemachine.StateThinking:
		// host-filable dimensions
	default:
		return fmt.Errorf("scheduler: 宿主不许登记状态 %s（只许 Acting / AwaitingApproval / Thinking）", string(st))
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t := s.live[id]
	if t == nil {
		return ErrUnknownTask
	}
	if t.state == statemachine.StateSettling || t.state == statemachine.StateError {
		return fmt.Errorf("scheduler: 任务 %s 已结束，状态定在 %s 不再接受改写", id, string(t.state))
	}
	t.state = st
	return nil
}

// Close shuts the door (D38e step 1): accepted tasks keep running, new
// Submits are refused. It is idempotent.
func (s *Scheduler) Close() {
	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()
}

// SetParent files the authority join for one task: parentID is the task id
// whose task.cancel may stop id through the tools. An empty parentID clears
// the entry (root row). Unknown tasks answer ErrUnknownTask.
func (s *Scheduler) SetParent(id, parentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.live[id] == nil {
		return ErrUnknownTask
	}
	if parentID == "" {
		delete(s.parents, id)
	} else {
		s.parents[id] = parentID
	}
	return nil
}

// ParentOf returns the filed parent of a task ("" = root or unknown).
func (s *Scheduler) ParentOf(id string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.parents[id]
}

// Closed reports whether the door is shut.
func (s *Scheduler) Closed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

// CancelAll implements the D38e step-3 shape over this scheduler's tasks:
// the door closes, queued tasks are settled, every running root ctx is
// cancelled, and the caller waits up to timeout for all of them to drain.
// It returns the ids that had NOT drained when the wait ran out (abandoned
// and recorded, D37c). The ctx a caller passes bounds the whole step;
// deadline exceeded answers context.DeadlineExceeded alongside the list.
func (s *Scheduler) CancelAll(ctx context.Context, timeout time.Duration) (abandoned []string, err error) {
	s.Close()

	s.mu.Lock()
	var running []chan struct{}
	var runningIDs []string
	var runningTasks []*task
	for _, id := range s.order {
		t := s.live[id]
		if t == nil {
			continue
		}
		switch t.state {
		case statemachine.StateQueued:
			for i, q := range s.queue {
				if q == t {
					s.queue = slices.Delete(s.queue, i, i+1)
					break
				}
			}
			select {
			case <-t.done:
			default:
				t.state = statemachine.StateSettling
				t.note = "已取消（关停排队清场）"
				close(t.done)
			}
		default:
			if t.root != nil && t.root.Cancel != nil {
				t.root.Cancel()
			}
			running = append(running, t.done)
			runningIDs = append(runningIDs, t.id)
			runningTasks = append(runningTasks, t)
		}
	}
	s.mu.Unlock()

	if timeout <= 0 {
		timeout = taskDrainWait
	}
	tm := observe.NewTimeout(timeout)
	for len(running) > 0 {
		stillDones := running[:0]
		stillIDs := runningIDs[:0]
		still := runningTasks[:0]
		for i, done := range running {
			select {
			case <-done:
			default:
				stillDones = append(stillDones, done)
				stillIDs = append(stillIDs, runningIDs[i])
				still = append(still, runningTasks[i])
			}
		}
		running, runningIDs, runningTasks = stillDones, stillIDs, still
		if len(running) == 0 {
			break
		}
		if tm.Expired() {
			for _, t := range runningTasks {
				abandoned = append(abandoned, t.id)
				s.mu.Lock()
				if t.state != statemachine.StateSettling && t.state != statemachine.StateError {
					t.note = strings.TrimSpace(t.note + "；关停放弃等待（根未排干，D37c）")
				}
				s.mu.Unlock()
			}
			break
		}
		if err := ctx.Err(); err != nil {
			for _, t := range runningTasks {
				abandoned = append(abandoned, t.id)
			}
			return abandoned, err
		}
		time.Sleep(2 * time.Millisecond)
	}
	return abandoned, nil
}

// waitingNote is the ONE conflict sentence the ticket freezes: the queued
// task's visible reason. It reads 「在等 <task> 释放 <path>」 with the holder's
// task id and the canonical path.
func waitingNote(owner, path string) string {
	return "在等 " + owner + " 释放 " + path
}

// newTaskID returns a UUIDv4-shaped task id (the same shape class the agent
// loop's task_log.id uses), built from crypto/rand so no new dependency
// appears, with the same never-panic clock fallback as internal/agent.
func newTaskID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// Entropy source unavailable: fall back to a clock-derived unique id.
		n := time.Now().UnixNano()
		for i := range b {
			b[i] = byte(n >> (8 * uint(i%8)))
		}
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
