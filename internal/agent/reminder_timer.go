package agent

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// ---------------------------------------------------------------------------
// Ticket 31: the one-shot reminder timer (in-process, monotonic, notify-only).
//
// This is the process-side half the memory layer defers to by name
// (dao_reminders.go:13-16: "Lifecycle policy (who arms the timer, what fire
// does) belongs to internal/tools + internal/agent").
//
// Shape, and why it is this shape:
//
//   - ONE named+owned goroutine per armed reminder, started through
//     observe.Registry.Spawn — the only sanctioned `go` statement in the tree
//     (D38b / goroutine.go:22-25), which installs the recover boundary. There
//     is no scheduler and no resident ticker here: 票 31 rejects the
//     generalization, and D38b's frozen budget is 常驻 6 + 每任务 3, so a
//     permanent dispatcher goroutine would be a seventh resident nobody approved.
//   - SessionScope-INDEPENDENT: the waits hang off this type's own process root,
//     derived from context.Background() — not off a session.Scope and not off any
//     task root. A reminder fires minutes or hours after the task that set it
//     ended, and a Warm dispose (C31 / D43 #31) tears down ENGINES, not
//     reminders. That is the ticket's "survive Warm→Sleeping (process alive) but
//     NOT process exit", and Close is what process exit calls.
//   - MONOTONIC wait against a STORED WALL TARGET (D42#9): FireAt is wall clock
//     because 「下午三点」is a wall-clock notion (schema_reminders_v6.go:23-26
//     says the WAIT is computed on the monotonic clock by this caller). The
//     delta is derived from the wall reading ONCE, at arm time, and everything
//     after that is monotonic arithmetic on a value that carries the monotonic
//     reading. A system-clock jump while a reminder is pending therefore cannot
//     move the fire, because nothing re-reads the wall clock to decide
//     "is it time yet" — that re-subtraction is the shape AGENTS.md §1.2 and
//     d22scan ban, and 票 31 AC#4 is the case that keeps it banned here.
//
// Fire does exactly three things, in this order:
//
//  1. MarkReminderFired — the one-way armed→fired stamp, whose applied=false
//     answer IS the delivery gate: a timer that lost the race with a cancel
//     delivers nothing (the DAO's late-timer guard at dao_reminders.go:36-39,
//     first production reader of it). A cancel that landed before the wait was
//     released never reaches this stamp at all: run() re-checks the stop
//     channel after the release (ticket 180). The stamp covers the opposite
//     order - a cancel that arrives while fire is already in flight.
//  2. Notify — the ONLY thing a fire can do. 票 31's second hard exclusion says
//     a fire never executes a tool and never calls an LLM, so there is no
//     dispatch path in this file at all: just the C3 `notify` port.
//  3. one tool_call audit row (D35: "what did it actually do to my machine").
//     Booked post-hoc with outcome set, so memory's stampToolCallInsert lands
//     ended_at and the row ages normally in retention (ticket 368's rule).
//
// Boot: ReportMissed is the ONLY restart behaviour — CountMissedArmed, one
// count sentence, ⛔ no per-item replay (票 31 AC#2; PLAN.md:2610).
//
// ⛔ NOT implemented here, on purpose (票 31's four hard exclusions, mirrored in
// the tools layer's registry note): recurring/cron, action-on-fire, OS-level
// scheduling, restart replay.
// ---------------------------------------------------------------------------

// ReminderNotifier is the C3 `notify` port. Its method set is present.Notifier's
// verbatim (internal/present/dispatch.go:41), so whatever Notifier the
// composition root already built for tier 3 plugs in here unchanged. It is
// re-declared rather than imported because internal/present imports
// internal/agent — importing back would close a cycle, and the port is one
// method wide anyway.
type ReminderNotifier interface {
	Notify(title, body string) error
}

// ReminderRows is the reminders DAO subset this layer reads and stamps,
// satisfied by *memory.Store. An interface (not the Store) because this file
// owns two methods' behaviour and must not acquire the whole table's face.
type ReminderRows interface {
	MarkReminderFired(ctx context.Context, id string) (bool, error)
	CountMissedArmed(ctx context.Context, now int64) (int, error)
}

// ReminderClock is this layer's only source of time. It follows the shape
// approval/approval.go:20-32 already froze (which cites internal/llm/ratelimit.go
// as the origin): production passes time.Now — whose readings carry the
// monotonic clock — and tests pass a fake they can advance, which is what makes
// 票 31 AC#4's +5h jump an assertion instead of a sleep.
//
// The two dials are the whole discipline in one interface:
//
//   - MonoNow anchors the deadline. It is the only reading any wait or lateness
//     decision is derived from, so a system-clock jump cannot move a fire.
//   - WallNow is for the ONE unavoidable translation (a stored wall target into
//     a duration, at arm time) and for the boot missed-count query, which is a
//     question about persisted records.
//
// ⚠ Named gap, stated instead of silently missing: 票 31's line "recompute delta
// on wake" needs an OS suspend/resume signal, and this layer has no seam for one
// (that is internal/proc / cmd/wisp territory, outside ticket 31's write face).
// The monotonic anchor is stored so the resume path that eventually lands has
// the deadline to recompute against; until then a suspended process simply
// fires late, and run() says so out loud (REMINDER-FIRE-LATE).
type ReminderClock interface {
	MonoNow() time.Time
	WallNow() time.Time
	After(d time.Duration) <-chan time.Time
}

// SystemReminderClock is the production ReminderClock.
type SystemReminderClock struct{}

// MonoNow implements ReminderClock: time.Now keeps the monotonic reading, so
// arithmetic on it is monotonic (observe/clock.go:9-16).
func (SystemReminderClock) MonoNow() time.Time { return time.Now() }

// WallNow implements ReminderClock.
func (SystemReminderClock) WallNow() time.Time { return observe.NowWallUTC() }

// After implements ReminderClock. Go timers run on the monotonic clock.
func (SystemReminderClock) After(d time.Duration) <-chan time.Time {
	if d < 0 {
		d = 0
	}
	return time.After(d)
}

// ReminderNotifyTitle is the title a fired reminder carries on the toast.
const ReminderNotifyTitle = "Wisp 提醒"

// ReminderAuditTool names the fire on its tool_call row. It is deliberately NOT
// a registered C1 tool name: the row records a host-side action (a fire), which
// no tool call made. Booking it as "reminder.set" would say the model called
// set at the fire instant, which is false.
const ReminderAuditTool = "reminder.fire"

// reminderTimerGoroutine is the D38b spawn name for one wait. It is NOT in the
// frozen roster (PLAN.md:2841-2846 / observe/goroutine.go:40-60), so
// ClassifyGoroutine answers CategoryUnknown and every arm logs the leak warning
// once. That is 31-r1's option 甲, and it is the branch this leg takes because
// the other one - a D38(b) roster line - is contract text no agent may edit
// (AGENTS.md §0「改契约＝人工批准」), and borrowing a roster name to silence the
// alarm is the shape ticket 229 was filed to end (its spawn_finish_watcher_229
// case convicts exactly that). The debt is booked in
// internal/tools/spawn_roster_census_229_test.go's documented-debt ledger with
// this reason, so the alarm stays meaningful for every OTHER name.
const reminderTimerGoroutine = "reminder-timer"

// ReminderTimerOptions configures a ReminderTimers.
type ReminderTimerOptions struct {
	// Rows is required: a fire that cannot stamp itself delivered is a
	// reminder that would later be delivered twice.
	Rows ReminderRows
	// Notify is required: it is the whole of what a fire is allowed to do.
	Notify ReminderNotifier
	// Journal books the fire's tool_call row. nil disables booking; it never
	// gates delivery (a dead audit sink must not eat the notification).
	Journal Journal
	// Clock is the time source. nil = SystemReminderClock.
	Clock ReminderClock
	// Registry is the D38b roster. nil = observe.Default.
	Registry *observe.Registry
	// Root is the owning root. nil = a fresh process root named
	// "reminder-timers", which Close cancels. It is deliberately not a task
	// root and not a session scope: see the header.
	Root *observe.Root
	// Logf receives this layer's diagnostic lines (nil = stay quiet).
	Logf func(format string, args ...any)
}

// reminderWait is one armed row's handle. stop is closed exactly once, by
// Disarm or by Close; the wait watches it because observe.Registry hands the
// spawned body the ROOT ctx (goroutine.go:287-292), so a per-row stop needs its
// own channel rather than a derived context.
type reminderWait struct {
	id   string
	note string
	// deadline is the monotonic anchor of this wait: MonoNow() PLUS the delta
	// derived at arm time. Nothing re-derives it from a wall reading, which is
	// what makes 票 31 AC#4's clock jump unable to move a scheduled fire; the
	// resume path only needs to recompute against this value to stay honest.
	deadline time.Time
	stop     chan struct{}
	once     sync.Once
}

func (w *reminderWait) halt() { w.once.Do(func() { close(w.stop) }) }

// ReminderTimers is the process's one-shot reminder wait table. Safe for
// concurrent use, and it satisfies internal/tools' ReminderTimers seam
// (Arm / Disarm / Pending) by construction: tools must not import this package
// for it (agent is the lower layer - internal/tools/task.go:13 already imports
// agent), so the face is an interface on the tools side.
type ReminderTimers struct {
	rows    ReminderRows
	notify  ReminderNotifier
	j       Journal
	clock   ReminderClock
	reg     *observe.Registry
	root    *observe.Root
	ownRoot bool
	logf    func(format string, args ...any)

	mu     sync.Mutex
	wait   map[string]*reminderWait
	seq    int64
	closed bool
}

// NewReminderTimers composes the timer. A missing Rows or Notify is a
// construction error rather than an inert timer: a reminder that cannot stamp
// itself, or cannot notify, is not a reminder.
func NewReminderTimers(o ReminderTimerOptions) (*ReminderTimers, error) {
	if o.Rows == nil {
		return nil, fmt.Errorf("agent: NewReminderTimers 需要 Rows（标不出 fired 就不知道自己送达过）")
	}
	if o.Notify == nil {
		return nil, fmt.Errorf("agent: NewReminderTimers 需要 Notify（发不出通知就没有可做的事）")
	}
	clock := o.Clock
	if clock == nil {
		clock = SystemReminderClock{}
	}
	reg := o.Registry
	if reg == nil {
		reg = observe.Default
	}
	root := o.Root
	ownRoot := false
	if root == nil {
		root = observe.NewRoot("reminder-timers")
		ownRoot = true
	}
	return &ReminderTimers{
		rows: o.Rows, notify: o.Notify, j: o.Journal, clock: clock,
		reg: reg, root: root, ownRoot: ownRoot, logf: o.Logf,
		wait: map[string]*reminderWait{},
	}, nil
}

// Root exposes the owning root so a composition root can place this timer in
// its shutdown order (D38e: joins before exit).
func (t *ReminderTimers) Root() *observe.Root { return t.root }

// Arm starts the one-shot wait for one already-stored row. It is
// reminder.set's process side, and it never blocks.
//
// A row whose wall target is already past is refused, not fired: boot COUNTS
// past-due rows and never delivers them (票 31 AC#2 / PLAN.md:2610). Getting a
// future reminder means a future FireAt, and that is the caller's decision to
// make with a fresh one - not this layer's to invent.
func (t *ReminderTimers) Arm(r memory.Reminder) error {
	if r.ID == "" {
		return fmt.Errorf("agent: Arm 需要提醒 id")
	}
	if r.Status != memory.ReminderArmed {
		return fmt.Errorf("agent: Arm 只接 armed 行（当前 %q）", r.Status)
	}

	wallNow := t.clock.WallNow()
	delta := time.Unix(r.FireAt, 0).UTC().Sub(wallNow)
	if delta <= 0 {
		return fmt.Errorf("agent: 提醒 %s 的目标时刻已经过去（还差 %s），不挂等待：到点前进程没在跑就不会触发，重启也不补发",
			r.ID, delta.Round(time.Second))
	}

	w := &reminderWait{
		stop: make(chan struct{}),
		id:   r.ID,
		note: r.Note,
	}
	// The anchor and the budget come from the SAME monotonic reading, so the
	// wait and the lateness check below cannot disagree. `delta` is the one
	// place a wall reading enters timing math: the stored target is a wall
	// instant (「下午三点」) and must be translated to a duration once, at arm
	// time. After this line no wall value is read to decide "is it time yet"
	// (D42#9; AGENTS.md §1.2 bans the wall-delta shape outright).
	monoNow := t.clock.MonoNow()
	w.deadline = monoNow.Add(delta)

	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return fmt.Errorf("agent: 计时器已 Close，这条提醒挂不上去")
	}
	if _, dup := t.wait[r.ID]; dup {
		t.mu.Unlock()
		return fmt.Errorf("agent: 提醒 %s 已经挂着了（同一条不许挂两枚，否则会通知两次）", r.ID)
	}
	// Register the wait BEFORE spawning, so a Disarm that arrives immediately
	// after Arm returns always finds the handle it means to stop.
	t.wait[r.ID] = w
	t.mu.Unlock()

	// The timer is armed on the caller's goroutine (not inside the spawned one):
	// the ticket's "进程没在跑就不会触发" is only honest if Arm's return means the
	// wait really exists, and a test that advances a fake clock cannot race a
	// goroutine that has not registered yet.
	ready := t.clock.After(delta)
	t.reg.Spawn(reminderTimerGoroutine, "agent", t.root, func(ctx context.Context) {
		t.run(ctx, w, ready)
	})

	t.emit("agent: REMINDER-ARM id=%s wait=%s", r.ID, delta.Round(time.Second))
	return nil
}

// run is one reminder's whole life: wait, then fire. It never loops and never
// re-reads the wall clock to decide whether it is time.
func (t *ReminderTimers) run(ctx context.Context, w *reminderWait, ready <-chan time.Time) {
	defer t.retire(w.id)

	select {
	case <-ctx.Done():
		t.emit("agent: REMINDER-DISARMED id=%s - 根停了，没有送达任何东西", w.id)
		return
	case <-w.stop:
		t.emit("agent: REMINDER-DISARMED id=%s - 没有送达任何东西", w.id)
		return
	case <-ready:
		// A cancel that landed BEFORE this release was observable still wins.
		// Once stop is closed and ready has a value, Go's select picks
		// uniformly among the ready cases, so without the re-check below the
		// same user cancel sometimes delivered a toast (measured red: 3 of 4
		// package runs, ticket 180 orphan-red repair). The stop channel is this
		// wait's authority; the fired-stamp guard in fire stays the second belt
		// for the order this check cannot see - a cancel that arrives after
		// fire has already claimed the row.
		select {
		case <-w.stop:
			t.emit("agent: REMINDER-DISARMED id=%s - 没有送达任何东西", w.id)
			return
		default:
		}
	}

	// Monotonic lateness, for the resume case: a wait released long after its
	// anchor means the process was suspended or starved, not that the clock
	// moved. It is a log line, never a re-derivation of the fire time.
	if late := t.clock.MonoNow().Sub(w.deadline); late > time.Minute {
		t.emit("agent: REMINDER-FIRE-LATE id=%s late=%s（进程被挂起过？按单调时钟补发这一条，不重算目标）",
			w.id, late.Round(time.Second))
	}
	t.fire(ctx, w)
}

// fire performs the delivery. See the header for the order and why each step is
// where it is.
func (t *ReminderTimers) fire(ctx context.Context, w *reminderWait) {
	applied, err := t.rows.MarkReminderFired(ctx, w.id)
	if err != nil {
		// Unknown row state: deliver nothing. A notification this layer cannot
		// record as delivered is one the user could later get twice.
		t.book(ctx, w.id, "error", "internal")
		t.emit("agent: REMINDER-FIRE-FAILED id=%s mark=error - 不送达任何东西: %v", w.id, err)
		return
	}
	if !applied {
		// The DAO's late-timer guard spoke: the row is no longer armed (a cancel
		// won the race, or another path already fired it). Deliver nothing and
		// book nothing - the event did not happen, and a row for it would be a
		// record of something the user never got.
		t.emit("agent: REMINDER-FIRE-SKIPPED id=%s - 那一行已经不是 armed（多半刚被取消），不投递", w.id)
		return
	}
	if nerr := t.notify.Notify(ReminderNotifyTitle, w.note); nerr != nil {
		// The row stays fired: the reminder DID reach its moment, and 票 31
		// forbids using a fired row as a re-arming handle (models.go:134-136).
		// The failure is booked and logged loudly rather than retried.
		t.book(ctx, w.id, "error", "internal")
		t.emit("agent: REMINDER-NOTIFY-FAILED id=%s err=%v - 通知没发出去，行已标 fired（不重发）", w.id, nerr)
		return
	}
	t.book(ctx, w.id, "success", "")
	t.emit("agent: REMINDER-FIRED id=%s - 已通知并落 fired", w.id)
}

// book writes the fire's tool_call row (票 31: "fire records a tool_call row for
// auditability"). Booking never gates delivery.
func (t *ReminderTimers) book(ctx context.Context, id, outcome, errClass string) {
	if t.j == nil {
		return
	}
	t.mu.Lock()
	t.seq++
	seq := t.seq
	t.mu.Unlock()

	tc := memory.ToolCall{
		TaskID:        reminderAuditTaskID(id),
		Seq:           seq,
		Tool:          ReminderAuditTool,
		ArgsJSON:      fmt.Sprintf(`{"id":%q}`, id),
		RiskLevel:     memory.RiskL0, // notify is L0 (PLAN.md:2568, SPEC-07 §3:66)
		Decision:      DecisionAllow, // not gated: it IS the already-granted notify capability
		Outcome:       outcome,
		ErrorClass:    errClass,
		CorrelationID: id,
	}
	if _, err := t.j.InsertToolCall(ctx, tc); err != nil {
		t.emit("agent: REMINDER-AUDIT-FAILED id=%s err=%v", id, err)
	}
}

// reminderAuditTaskID keys the fire's tool_call row. tool_call.task_id is
// required by the DAO (memory/models.go:187) and a fire belongs to no task - it
// happens after the task that set it ended - so the row hangs off a namespaced
// reminder id. The "reminder:" prefix is what lets a reader that groups
// tool_call by task_id tell a host-side fire from a task's own calls.
func reminderAuditTaskID(id string) string {
	if id == "" {
		return "reminder:fire"
	}
	return "reminder:" + id
}

// retire drops one wait from the table (the wait's own goroutine exit path).
func (t *ReminderTimers) retire(id string) {
	t.mu.Lock()
	delete(t.wait, id)
	t.mu.Unlock()
}

// Disarm stops waiting for one reminder. reminder.cancel calls it BEFORE it
// stamps the row cancelled, so a wait cannot deliver a toast for a reminder the
// user just took away; the fired-stamp's own guard is the second belt.
func (t *ReminderTimers) Disarm(id string) bool {
	t.mu.Lock()
	w := t.wait[id]
	delete(t.wait, id)
	t.mu.Unlock()
	if w == nil {
		return false
	}
	w.halt()
	t.emit("agent: REMINDER-DISARM id=%s", id)
	return true
}

// Pending reports how many waits this process holds. reminder.list reads it to
// keep "armed in the table" and "actually waiting here" apart.
func (t *ReminderTimers) Pending() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.wait)
}

// PendingIDs names the waiting ids, sorted, for a host that has to show or
// reconcile the table (and for tests that need a deterministic view).
func (t *ReminderTimers) PendingIDs() []string {
	t.mu.Lock()
	var ids []string
	for id := range t.wait {
		ids = append(ids, id)
	}
	t.mu.Unlock()
	sort.Strings(ids)
	return ids
}

// MissedReport is boot's answer: how many armed rows had already come due, and
// the single sentence the user is allowed to get about them.
type MissedReport struct {
	Count int
	// Sentence is "" when Count == 0: nothing to say, so nothing is said.
	Sentence string
}

// ReportMissed is 票 31's restart rule and the ONLY thing boot may do with rows
// that came due while the process was down: count them, say the count, deliver
// none. A per-item replay would need this layer to fire rows whose targets are
// past - which Arm refuses by construction - and it is what PLAN.md:2610 rejects.
//
// It only reads: the rows stay armed in the table exactly as CountMissedArmed
// found them (dao_reminders.go:152-156 records that choice on purpose).
func (t *ReminderTimers) ReportMissed(ctx context.Context) (MissedReport, error) {
	now := t.clock.WallNow().Unix()
	n, err := t.rows.CountMissedArmed(ctx, now)
	if err != nil {
		return MissedReport{}, err
	}
	out := MissedReport{Count: n}
	if n > 0 {
		out.Sentence = fmt.Sprintf("你错过了 %d 条提醒（进程没在跑就不会触发，重启后也不逐条补发）", n)
	}
	return out, nil
}

// Close stops every wait and joins them under a monotonic budget (D42#9),
// returning how many are still pending. It is what process exit does, and
// reminders do NOT survive it - the limit reminder.set must state up front.
func (t *ReminderTimers) Close(timeout time.Duration) int {
	t.mu.Lock()
	var waits []*reminderWait
	for id, w := range t.wait {
		waits = append(waits, w)
		delete(t.wait, id)
	}
	t.closed = true
	t.mu.Unlock()

	for _, w := range waits {
		w.halt()
	}
	if t.ownRoot && t.root.Cancel != nil {
		// Only the root this type minted: a root a composition root handed in
		// stays that owner's to cancel.
		t.root.Cancel()
	}
	return t.root.Wait(timeout)
}

func (t *ReminderTimers) emit(format string, args ...any) {
	if t.logf == nil {
		return
	}
	t.logf(format, args...)
}
