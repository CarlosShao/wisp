package agent

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 31 acceptance for the timer layer. Two of these cases (AC#1's fire and
// AC#4's clock jump) are only possible because the clock is injected: a real
// 5-hour jump cannot be produced in a test, and sleeping for a real delta would
// be a test that hangs on a slow machine.

// fakeReminderClock is the injected ReminderClock. It carries two independent
// dials on purpose (票 31 AC#4's whole point):
//
//   - mono: the monotonic reading. Advance() moves it, and only this dial can
//     release a pending wait.
//   - wall: the wall reading. Jump() moves ONLY this one, which is exactly what
//     an NTP correction or a timezone change does to a live process.
//
// A wait that decided "is it time yet" by re-comparing a wall target to a wall
// now would fire (or never fire) on the Jump() alone - that is the mutation this
// clock exists to make red.
type fakeReminderClock struct {
	mu      sync.Mutex
	mono    time.Time
	wall    time.Time
	budgets []time.Duration // what each After() was asked to wait, in order
	timers  []*fakeTimer
}

// fakeTimer is one armed After(): the instant (on the MONOTONIC dial) it is due
// and the channel it releases.
type fakeTimer struct {
	due      time.Time
	ch       chan time.Time
	released bool
}

func newFakeReminderClock() *fakeReminderClock {
	base := time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)
	return &fakeReminderClock{mono: base, wall: base}
}

func (c *fakeReminderClock) MonoNow() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.mono
}

func (c *fakeReminderClock) WallNow() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.wall
}

func (c *fakeReminderClock) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	if d < 0 {
		d = 0
	}
	t := &fakeTimer{due: c.mono.Add(d), ch: make(chan time.Time, 1)}
	c.budgets = append(c.budgets, d)
	c.timers = append(c.timers, t)
	return t.ch
}

// Advance moves the MONOTONIC dial and releases ONLY the waits whose budget has
// run out. The wall dial is untouched, which is the point: nothing in this fake
// can release a wait because of anything the wall dial does.
func (c *fakeReminderClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.mono = c.mono.Add(d)
	for _, t := range c.timers {
		if t.released || t.due.After(c.mono) {
			continue
		}
		t.released = true
		t.ch <- c.mono
	}
	c.mu.Unlock()
}

// JumpWall moves ONLY the wall dial - a system-clock jump that must not move any
// scheduled fire.
func (c *fakeReminderClock) JumpWall(d time.Duration) {
	c.mu.Lock()
	c.wall = c.wall.Add(d)
	c.mu.Unlock()
}

func (c *fakeReminderClock) lastBudget() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.budgets) == 0 {
		return -1
	}
	return c.budgets[len(c.budgets)-1]
}

// fakeReminderRows is a ReminderRows face over an in-memory status map, with
// the DAO's one-way rule enforced: MarkReminderFired answers applied=false for
// a row that is no longer armed.
type fakeReminderRows struct {
	mu       sync.Mutex
	status   map[string]string
	missed   int
	missedAt int64
	marked   []string
	countRun int
}

func newFakeRows(status map[string]string) *fakeReminderRows {
	c := map[string]string{}
	for k, v := range status {
		c[k] = v
	}
	return &fakeReminderRows{status: c}
}

func (r *fakeReminderRows) MarkReminderFired(_ context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.status[id] != memory.ReminderArmed {
		return false, nil
	}
	r.status[id] = memory.ReminderFired
	r.marked = append(r.marked, id)
	return true, nil
}

func (r *fakeReminderRows) CountMissedArmed(_ context.Context, now int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.missedAt = now
	r.countRun++
	return r.missed, nil
}

func (r *fakeReminderRows) statusOf(id string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.status[id]
}

func (r *fakeReminderRows) markCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.marked)
}

func (r *fakeReminderRows) cancel(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.status[id] == memory.ReminderArmed {
		r.status[id] = memory.ReminderCancelled
	}
}

// fakeNotifier records every Notify.
type fakeNotifier struct {
	mu    sync.Mutex
	calls []string
	err   error
}

func (n *fakeNotifier) Notify(title, body string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.calls = append(n.calls, title+"\x00"+body)
	return n.err
}

func (n *fakeNotifier) count() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.calls)
}

func (n *fakeNotifier) last() string {
	n.mu.Lock()
	defer n.mu.Unlock()
	if len(n.calls) == 0 {
		return ""
	}
	return n.calls[len(n.calls)-1]
}

// fakeJournal records the fire's tool_call rows.
type fakeJournal struct {
	mu   sync.Mutex
	rows []memory.ToolCall
}

func (j *fakeJournal) StartTaskLog(context.Context, memory.TaskLog) error { return nil }
func (j *fakeJournal) FinishTaskLog(context.Context, string, string, string, int64, int64, int64, string) error {
	return nil
}

func (j *fakeJournal) DecideToolCall(context.Context, int64, string, *int64) error {
	return nil
}

func (j *fakeJournal) FinishToolCall(context.Context, int64, string, string) error {
	return nil
}

func (j *fakeJournal) InsertToolCall(_ context.Context, tc memory.ToolCall) (int64, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.rows = append(j.rows, tc)
	return int64(len(j.rows)), nil
}

func (j *fakeJournal) all() []memory.ToolCall {
	j.mu.Lock()
	defer j.mu.Unlock()
	return append([]memory.ToolCall(nil), j.rows...)
}

// ---- AC#1: set arms a wait that fires with a notify ------------------------

// Test31ArmFiresWithNotify is the ticket's "set → fires at T+delta with notify
// (fixture clock)". The delta is what the fake clock was ASKED to wait, so the
// assertion is about the arithmetic, not about sleeping.
func Test31ArmFiresWithNotify(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{"rem_a": memory.ReminderArmed})
	ntf := &fakeNotifier{}
	reg := observe.NewRegistry()
	root := observe.NewRoot("test-reminder-timers")

	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Clock: clock, Registry: reg, Root: root,
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	target := clock.WallNow().Add(20 * time.Minute).Unix()
	if err := tm.Arm(memory.Reminder{ID: "rem_a", FireAt: target, Note: "开会", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("Arm: %v", err)
	}

	if got := clock.lastBudget(); got != 20*time.Minute {
		t.Fatalf("单调等待 = %s, want 20m（FireAt 墙钟目标减一次墙钟读数）", got)
	}
	if n := reg.CountByName(reminderTimerGoroutine); n != 1 {
		t.Fatalf("名册里的 reminder 协程 = %d, want 1（named+owned 的一半是 named）", n)
	}
	if tm.Pending() != 1 {
		t.Fatalf("Pending = %d, want 1", tm.Pending())
	}

	clock.Advance(20 * time.Minute)
	waitFor(t, 2*time.Second, func() bool { return ntf.count() == 1 }, "通知没发出去")

	if got := ntf.last(); got != ReminderNotifyTitle+"\x00"+"开会" {
		t.Fatalf("通知内容 = %q, want title+note", got)
	}
	if got := rows.statusOf("rem_a"); got != memory.ReminderFired {
		t.Fatalf("行状态 = %q, want fired", got)
	}
	if n := rows.markCount(); n != 1 {
		t.Fatalf("MarkReminderFired 调用 = %d, want 1", n)
	}
	if tm.Pending() != 0 {
		t.Fatalf("fire 之后 Pending = %d, want 0（协程自己退休）", tm.Pending())
	}
}

// Test31FireBooksOneToolCallRow is 票 31's "fire records a tool_call row for
// auditability": one row, L0, outcome success, keyed to the reminder id.
func Test31FireBooksOneToolCallRow(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{"rem_b": memory.ReminderArmed})
	j := &fakeJournal{}
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: &fakeNotifier{}, Journal: j, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-book"),
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	target := clock.WallNow().Add(5 * time.Minute).Unix()
	if err := tm.Arm(memory.Reminder{ID: "rem_b", FireAt: target, Note: "喝水", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("Arm: %v", err)
	}
	clock.Advance(5 * time.Minute)
	waitFor(t, 2*time.Second, func() bool { return len(j.all()) == 1 }, "审计行没落")

	tc := j.all()[0]
	if tc.Tool != ReminderAuditTool {
		t.Errorf("tool = %q, want %q", tc.Tool, ReminderAuditTool)
	}
	if tc.RiskLevel != memory.RiskL0 {
		t.Errorf("risk_level = %q, want L0（notify 是 L0）", tc.RiskLevel)
	}
	if tc.Outcome != "success" {
		t.Errorf("outcome = %q, want success", tc.Outcome)
	}
	if tc.TaskID != reminderAuditTaskID("rem_b") {
		t.Errorf("task_id = %q, want 按提醒 id 归名（fire 不属于任何任务）", tc.TaskID)
	}
	if tc.CorrelationID != "rem_b" {
		t.Errorf("correlation_id = %q, want rem_b", tc.CorrelationID)
	}
}

// ---- cancel removes the wait ------------------------------------------------

// Test31DisarmStopsTheWait is 票 31's "cancel removes". Nothing is delivered,
// and the row is NOT stamped fired (reminder.cancel owns that write, and the
// timer must never be the thing that says fired for a cancelled reminder).
func Test31DisarmStopsTheWait(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{"rem_c": memory.ReminderArmed})
	ntf := &fakeNotifier{}
	j := &fakeJournal{}
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Journal: j, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-disarm"),
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	target := clock.WallNow().Add(10 * time.Minute).Unix()
	if err := tm.Arm(memory.Reminder{ID: "rem_c", FireAt: target, Note: "别发", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("Arm: %v", err)
	}
	if !tm.Disarm("rem_c") {
		t.Fatal("Disarm = false, want true（刚挂上的等待必须能停）")
	}
	if tm.Pending() != 0 {
		t.Fatalf("Pending = %d, want 0", tm.Pending())
	}
	if tm.Disarm("rem_c") {
		t.Fatal("第二次 Disarm 仍报 true：同一枚停了两次")
	}

	clock.Advance(10 * time.Minute)
	// Give any in-flight fire the chance to land, then prove it did not.
	time.Sleep(50 * time.Millisecond)
	if ntf.count() != 0 {
		t.Fatalf("取消后仍发了 %d 条通知", ntf.count())
	}
	if got := rows.statusOf("rem_c"); got != memory.ReminderArmed {
		t.Fatalf("行状态 = %q, want 仍是 armed（取消那一写是 tools 侧的）", got)
	}
	if n := len(j.all()); n != 0 {
		t.Fatalf("审计行 = %d 枚, want 0（没发生的事不留记录）", n)
	}
}

// Test31LateTimerDeliversNothing is the DAO's late-timer guard as a production
// caller: a fire that arrives after a cancel must see applied=false from
// MarkReminderFired and deliver nothing.
func Test31LateTimerDeliversNothing(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{"rem_d": memory.ReminderArmed})
	ntf := &fakeNotifier{}
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-late"),
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	target := clock.WallNow().Add(1 * time.Minute).Unix()
	if err := tm.Arm(memory.Reminder{ID: "rem_d", FireAt: target, Note: "该被取消", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("Arm: %v", err)
	}
	// Cancel the ROW without stopping the wait: this is the race the guard is
	// written for, and the wait is not Disarmed so it really does reach fire.
	rows.cancel("rem_d")

	clock.Advance(1 * time.Minute)
	waitFor(t, 2*time.Second, func() bool { return tm.Pending() == 0 }, "协程没退休")
	time.Sleep(50 * time.Millisecond)

	if ntf.count() != 0 {
		t.Fatalf("已取消的行仍通知了 %d 次（DAO 的迟到守卫没接住）", ntf.count())
	}
	if rows.markCount() != 0 {
		t.Fatalf("MarkReminderFired 生效次数 = %d, want 0", rows.markCount())
	}
}

// ---- AC#4: clock jump -------------------------------------------------------

// Test31ClockJumpKeepsMonotonicFire is 票 31 AC#4: the system clock jumps +5h
// while a reminder is pending → it still fires relative to the monotonic clock.
//
// The mutation this pins is the banned one: recomputing "FireAt - wallNow" when
// deciding whether it is time. Under that shape the +5h jump alone would move
// the fire, and this case goes red without advancing the monotonic dial at all.
func Test31ClockJumpKeepsMonotonicFire(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{"rem_e": memory.ReminderArmed})
	ntf := &fakeNotifier{}
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-jump"),
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	// Arming is deliberately BEFORE the jump, so the jump is mid-pending.
	target := clock.WallNow().Add(30 * time.Minute).Unix()
	if err := tm.Arm(memory.Reminder{ID: "rem_e", FireAt: target, Note: "跳钟也要提醒", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("Arm: %v", err)
	}
	if got := clock.lastBudget(); got != 30*time.Minute {
		t.Fatalf("单调等待 = %s, want 30m", got)
	}

	// The +5h wall jump: nothing may react to it.
	clock.JumpWall(5 * time.Hour)

	time.Sleep(50 * time.Millisecond)
	if ntf.count() != 0 {
		t.Fatal("墙钟跳 +5h 就把提醒提前了：等待用的是墙钟差（D42#9 禁的形状）")
	}
	if tm.Pending() != 1 {
		t.Fatalf("跳钟后 Pending = %d, want 1（等待还在，没被动过）", tm.Pending())
	}
	if n := rows.markCount(); n != 0 {
		t.Fatalf("跳钟后行被标 fired %d 次, want 0", n)
	}

	// Only the monotonic dial may release it, and it releases at T+30m.
	clock.Advance(30 * time.Minute)
	waitFor(t, 2*time.Second, func() bool { return ntf.count() == 1 }, "按单调时钟到点了却没触发")
	if got := rows.statusOf("rem_e"); got != memory.ReminderFired {
		t.Fatalf("行状态 = %q, want fired", got)
	}
}

// Test31ArmsPastTargetRefusesNoFire is AC#2's other half at this layer: a row
// whose wall target already passed is NOT armed, so no restart can make it fire.
func Test31ArmsPastTargetRefusesNoFire(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{"rem_old": memory.ReminderArmed})
	ntf := &fakeNotifier{}
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-past"),
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	past := clock.WallNow().Add(-1 * time.Hour).Unix()
	err = tm.Arm(memory.Reminder{ID: "rem_old", FireAt: past, Note: "已经过了", Status: memory.ReminderArmed})
	if err == nil {
		t.Fatal("过期行被 Arm 接了：重启补发的门就是这个形状打开的")
	}
	if tm.Pending() != 0 {
		t.Fatalf("Pending = %d, want 0", tm.Pending())
	}
	if ntf.count() != 0 {
		t.Fatalf("过期行还发了 %d 条通知", ntf.count())
	}
}

// ---- AC#2: boot reports a count, never a replay -----------------------------

// Test31BootReportsMissedCountOnly is 票 31 AC#2's restart half: boot gets one
// count sentence and NOTHING is delivered per item.
func Test31BootReportsMissedCountOnly(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{})
	rows.missed = 3
	ntf := &fakeNotifier{}
	j := &fakeJournal{}
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Journal: j, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-boot"),
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	rep, err := tm.ReportMissed(context.Background())
	if err != nil {
		t.Fatalf("ReportMissed: %v", err)
	}
	if rep.Count != 3 {
		t.Fatalf("Count = %d, want 3", rep.Count)
	}
	want := "你错过了 3 条提醒"
	if len(rep.Sentence) < len(want) || rep.Sentence[:len(want)] != want {
		t.Fatalf("Sentence = %q, want 以 %q 开头（票 31 逐字句形）", rep.Sentence, want)
	}
	if ntf.count() != 0 {
		t.Fatalf("boot 补发了 %d 条通知：跨重启只计数，不逐条弹（PLAN.md:2610）", ntf.count())
	}
	if n := len(j.all()); n != 0 {
		t.Fatalf("boot 落了 %d 枚 tool_call 行, want 0（没有 fire 就没有审计）", n)
	}
	// The query went out on the wall clock (it is a persisted-record question),
	// and exactly once.
	if rows.countRun != 1 {
		t.Fatalf("CountMissedArmed 跑了 %d 次, want 1", rows.countRun)
	}
	if rows.missedAt != clock.WallNow().Unix() {
		t.Fatalf("CountMissedArmed 的时刻 = %d, want 当前墙钟 %d", rows.missedAt, clock.WallNow().Unix())
	}
}

// Test31BootQuietWhenNothingMissed: zero missed → no sentence, no toast. A boot
// that apologises for nothing is noise, and noise trains the user to ignore the
// one message that matters.
func Test31BootQuietWhenNothingMissed(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{})
	ntf := &fakeNotifier{}
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-quiet"),
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	rep, err := tm.ReportMissed(context.Background())
	if err != nil {
		t.Fatalf("ReportMissed: %v", err)
	}
	if rep.Count != 0 || rep.Sentence != "" {
		t.Fatalf("零错过却报了 %d / %q", rep.Count, rep.Sentence)
	}
	if ntf.count() != 0 {
		t.Fatalf("零错过仍发了 %d 条", ntf.count())
	}
}

// ---- notify failure is loud, never silent ----------------------------------

// Test31NotifyFailureIsBookedAndNotRetried: the row stays fired (the reminder
// DID reach its moment) and the toast failure is booked + logged. D42#8's rule
// is that a toast failure never silences the rest of the path; here "the rest"
// is the record, and it is still written.
func Test31NotifyFailureIsBookedAndNotRetried(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{"rem_f": memory.ReminderArmed})
	ntf := &fakeNotifier{err: context.DeadlineExceeded}
	j := &fakeJournal{}
	var logged []string
	var lm sync.Mutex
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Journal: j, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-ntffail"),
		Logf: func(f string, a ...any) {
			lm.Lock()
			defer lm.Unlock()
			logged = append(logged, f)
		},
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	target := clock.WallNow().Add(2 * time.Minute).Unix()
	if err := tm.Arm(memory.Reminder{ID: "rem_f", FireAt: target, Note: "发不出的", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("Arm: %v", err)
	}
	clock.Advance(2 * time.Minute)
	waitFor(t, 2*time.Second, func() bool { return len(j.all()) == 1 }, "失败没落审计")

	if got := j.all()[0].Outcome; got != "error" {
		t.Fatalf("outcome = %q, want error", got)
	}
	if got := rows.statusOf("rem_f"); got != memory.ReminderFired {
		t.Fatalf("行状态 = %q, want fired（不重发、不借道重挂）", got)
	}
	lm.Lock()
	joined := ""
	for _, l := range logged {
		joined += l
	}
	lm.Unlock()
	if want := "REMINDER-NOTIFY-FAILED"; !containsLog(joined, want) {
		t.Fatalf("日志里没有 %q：%s", want, joined)
	}
	if ntf.count() != 1 {
		t.Fatalf("通知发了 %d 次, want 1（fire 不重试）", ntf.count())
	}
}

func containsLog(s, want string) bool {
	for i := 0; i+len(want) <= len(s); i++ {
		if s[i:i+len(want)] == want {
			return true
		}
	}
	return false
}

// ---- lifecycle --------------------------------------------------------------

// Test31CloseStopsEveryWait is the process-exit half of "reminders do NOT
// survive process exit": Close halts every wait and joins the goroutines.
func Test31CloseStopsEveryWait(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{
		"rem_g1": memory.ReminderArmed,
		"rem_g2": memory.ReminderArmed,
	})
	ntf := &fakeNotifier{}
	reg := observe.NewRegistry()
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Clock: clock, Registry: reg,
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}

	for _, id := range []string{"rem_g1", "rem_g2"} {
		target := clock.WallNow().Add(1 * time.Hour).Unix()
		if err := tm.Arm(memory.Reminder{ID: id, FireAt: target, Note: "退出前", Status: memory.ReminderArmed}); err != nil {
			t.Fatalf("Arm %s: %v", id, err)
		}
	}
	if tm.Pending() != 2 {
		t.Fatalf("Pending = %d, want 2", tm.Pending())
	}
	if n := reg.CountByName(reminderTimerGoroutine); n != 2 {
		t.Fatalf("在飞协程 = %d, want 2", n)
	}

	// Close waits on a monotonic budget (D42#9): a hang here is a bug, not slow
	// hardware, because nothing is scheduled to elapse.
	if pending := tm.Close(2 * time.Second); pending != 0 {
		t.Fatalf("Close 后仍有 %d 枚协程没收（D38c：任务/生命周期结束必须等派生协程退出）", pending)
	}
	if tm.Pending() != 0 {
		t.Fatalf("Close 后 Pending = %d, want 0", tm.Pending())
	}
	if n := reg.CountByName(reminderTimerGoroutine); n != 0 {
		t.Fatalf("Close 后在飞协程 = %d, want 0", n)
	}

	// And the table refuses new work afterwards rather than leaking it.
	if err := tm.Arm(memory.Reminder{
		ID:     "rem_g3",
		FireAt: clock.WallNow().Add(time.Hour).Unix(), Note: "Close 后", Status: memory.ReminderArmed,
	}); err == nil {
		t.Fatal("Close 之后 Arm 仍成功：这枚提醒永远不会触发")
	}
	if ntf.count() != 0 {
		t.Fatalf("Close 仍触发了 %d 条通知", ntf.count())
	}
}

// Test31RemindersSurviveRootHandover is the SessionScope-independence claim as
// a mechanism: waits hang off the timer's OWN root, so a caller that hands its
// root in keeps it alive across a task ending, and only Close stops them.
func Test31RemindersSurviveRootHandover(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{"rem_h": memory.ReminderArmed})
	ntf := &fakeNotifier{}
	root := observe.NewRoot("handover")
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: ntf, Clock: clock,
		Registry: observe.NewRegistry(), Root: root,
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	target := clock.WallNow().Add(45 * time.Minute).Unix()
	if err := tm.Arm(memory.Reminder{ID: "rem_h", FireAt: target, Note: "活过任务结束", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("Arm: %v", err)
	}
	// A composition root that does NOT own this root must not be able to stop
	// the wait by cancelling something of its own: the root handed in is shared,
	// so the assertion is that the wait is still here.
	if tm.Pending() != 1 {
		t.Fatalf("Pending = %d, want 1", tm.Pending())
	}
	if ids := tm.PendingIDs(); len(ids) != 1 || ids[0] != "rem_h" {
		t.Fatalf("PendingIDs = %v, want [rem_h]", ids)
	}
	clock.Advance(45 * time.Minute)
	waitFor(t, 2*time.Second, func() bool { return ntf.count() == 1 }, "跨过去仍没触发")
	if ntf.last() != ReminderNotifyTitle+"\x00"+"活过任务结束" {
		t.Fatalf("通知 = %q", ntf.last())
	}
}

// Test31ArmRejectsNonArmedAndDuplicate is the row-state half of the guard: a
// fired/cancelled row is history, never a re-arming handle (memory/models.go
// says exactly that), and one row cannot be waited on twice.
func Test31ArmRejectsNonArmedAndDuplicate(t *testing.T) {
	clock := newFakeReminderClock()
	rows := newFakeRows(map[string]string{"rem_i": memory.ReminderFired})
	tm, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows, Notify: &fakeNotifier{}, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-rearm"),
	})
	if err != nil {
		t.Fatalf("NewReminderTimers: %v", err)
	}
	defer func() { _ = tm.Close(2 * time.Second) }()

	target := clock.WallNow().Add(time.Hour).Unix()
	if err := tm.Arm(memory.Reminder{ID: "rem_i", FireAt: target, Note: "x", Status: memory.ReminderFired}); err == nil {
		t.Fatal("fired 行被重新挂上等待：单向终态被绕开了")
	}
	rows2 := newFakeRows(map[string]string{"rem_j": memory.ReminderArmed})
	tm2, err := NewReminderTimers(ReminderTimerOptions{
		Rows: rows2, Notify: &fakeNotifier{}, Clock: clock,
		Registry: observe.NewRegistry(), Root: observe.NewRoot("t31-dup"),
	})
	if err != nil {
		t.Fatalf("NewReminderTimers 2: %v", err)
	}
	defer func() { _ = tm2.Close(2 * time.Second) }()
	armed := memory.Reminder{ID: "rem_j", FireAt: target, Note: "y", Status: memory.ReminderArmed}
	if err := tm2.Arm(armed); err != nil {
		t.Fatalf("第一次 Arm: %v", err)
	}
	if err := tm2.Arm(armed); err == nil {
		t.Fatal("同一行挂了两枚等待：会通知两次")
	}
	if tm2.Pending() != 1 {
		t.Fatalf("Pending = %d, want 1", tm2.Pending())
	}
}

// Test31NewRequiresRowsAndNotify: both required faces are loud at construction,
// so a half-wired timer cannot exist and quietly do nothing.
func Test31NewRequiresRowsAndNotify(t *testing.T) {
	if _, err := NewReminderTimers(ReminderTimerOptions{Notify: &fakeNotifier{}}); err == nil {
		t.Fatal("没有 Rows 也建起来了")
	}
	if _, err := NewReminderTimers(ReminderTimerOptions{Rows: newFakeRows(nil)}); err == nil {
		t.Fatal("没有 Notify 也建起来了：fire 无事可做")
	}
}

// waitFor polls a condition under a monotonic bound. It is this file's own tiny
// helper (the agent package's tests have no such thing) and it never converts
// an unmet condition into a pass.
func waitFor(t *testing.T, bound time.Duration, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(bound) // carries the monotonic reading
	for !cond() {
		if remaining := time.Until(deadline); remaining <= 0 {
			t.Fatalf("超过 %s 仍未满足：%s", bound, msg)
		}
		time.Sleep(2 * time.Millisecond)
	}
}
