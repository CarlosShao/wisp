package observe

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

// Ticket 44 (C23 CostMeter) — the in-package instruments: the aggregate reads,
// the day/month buckets, the budget thresholds, and the price lookup. These
// use a fake store so the boundary cases (79/80/100%, a day with no rows, a
// negative booking) are reachable without a SQLite file and without waiting
// for a real midnight. The same three interfaces against the REAL store are
// pinned in internal/memory/costdao_44_test.go, and the loop's booking path in
// internal/agent/costbook_44_test.go.

// fakeStore is cost_daily + task_log in memory, reached through the same three
// interfaces the meter uses against memory.Store.
type fakeStore struct {
	days map[string]DayTotal
	// base is the task_log side: rows that ENDED at a unix instant with their
	// own C23 columns. Nothing writes it through the meter, so an aggregate
	// that disagrees with it is a real disagreement, not a shared bug.
	base []baseRow
	// bumps records the day labels PostTask wrote, in order.
	bumps []string
	// failBump makes the ledger fail, for the "booking error is reported" case.
	failBump error
}

type baseRow struct {
	endedAt int64
	total   DayTotal
}

func newFakeStore() *fakeStore {
	return &fakeStore{days: map[string]DayTotal{}}
}

func (f *fakeStore) BumpCostDay(_ context.Context, day string, in, out, micros int64) error {
	if f.failBump != nil {
		return f.failBump
	}
	if !validDay(day) {
		return fmt.Errorf("fake: bad day %q", day)
	}
	row := f.days[day]
	row.Day = day
	row.TokensIn += in
	row.TokensOut += out
	row.Micros += micros
	row.Tasks++
	f.days[day] = row
	f.bumps = append(f.bumps, day)
	// The base table gets the matching row only when the test planted one;
	// keeping them separate is what makes Reconcile mean something.
	return nil
}

func (f *fakeStore) SetCostDay(_ context.Context, day string, in, out, micros, tasks int64) error {
	if !validDay(day) {
		return fmt.Errorf("fake: bad day %q", day)
	}
	f.days[day] = DayTotal{Day: day, TokensIn: in, TokensOut: out, Micros: micros, Tasks: tasks}
	return nil
}

func (f *fakeStore) CostDaysBetween(_ context.Context, fromDay, toDay string) ([]DayTotal, error) {
	if !validDay(fromDay) || !validDay(toDay) {
		return nil, fmt.Errorf("fake: bad range %q..%q", fromDay, toDay)
	}
	if fromDay > toDay {
		return nil, fmt.Errorf("fake: inverted range")
	}
	var out []DayTotal
	for day, row := range f.days {
		if day >= fromDay && day <= toDay {
			out = append(out, row)
		}
	}
	return out, nil
}

func (f *fakeStore) SumTasksEndedBetween(_ context.Context, fromUnix, toUnix int64) (DayTotal, error) {
	if toUnix <= fromUnix {
		return DayTotal{}, fmt.Errorf("fake: empty range")
	}
	var sum DayTotal
	for _, r := range f.base {
		if r.endedAt >= fromUnix && r.endedAt < toUnix {
			// plantBase stores Tasks:1 per row, so Add carries the count. A
			// second ++ here would count every row twice, and a reconciliation
			// instrument that double-counts is worse than one that cannot read.
			sum = sum.Add(r.total)
		}
	}
	return sum, nil
}

// plantBase adds a finished task to the base table (the task_log mouth).
func (f *fakeStore) plantBase(endedAt time.Time, in, out, micros int64) {
	f.base = append(f.base, baseRow{
		endedAt: endedAt.Unix(),
		total:   DayTotal{TokensIn: in, TokensOut: out, Micros: micros, Tasks: 1},
	})
}

func validDay(day string) bool {
	if len(day) != 10 || day[4] != '-' || day[7] != '-' {
		return false
	}
	for i, r := range day {
		if i == 4 || i == 7 {
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// fixedZone is the calendar the tests bucket in: +08:00, the same offset the
// machine this product runs on uses, and spelled as a fixed zone so a rollover
// case cannot depend on the runner's tz database.
var fixedZone = time.FixedZone("t+8", 8*3600)

func dayAt(y int, mo time.Month, d, hh, mm int) time.Time {
	return time.Date(y, mo, d, hh, mm, 0, 0, fixedZone)
}

func newMeter(t *testing.T, f *fakeStore, at time.Time) *CostMeter {
	t.Helper()
	m, err := NewCostMeter(f,
		WithLedger(f), WithWriter(f), WithTasks(f),
		WithLocation(fixedZone), WithClock(func() time.Time { return at }))
	if err != nil {
		t.Fatalf("NewCostMeter: %v", err)
	}
	return m
}

func TestAC2TodayReadsOnlyItsOwnDay(t *testing.T) {
	f := newFakeStore()
	ctx := context.Background()
	f.days["2026-11-01"] = DayTotal{Day: "2026-11-01", TokensIn: 900, TokensOut: 400, Micros: 9000, Tasks: 3}
	m := newMeter(t, f, dayAt(2026, 11, 2, 9, 0))

	got, err := m.Today(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got.TokensIn != 0 || got.Micros != 0 || got.Tasks != 0 {
		t.Errorf("today = %+v, want the zero row: yesterday's 9000 micros must not leak into today", got)
	}
	if got.Day != "2026-11-02" {
		t.Errorf("today.Day = %q, want 2026-11-02 (the label is part of the answer)", got.Day)
	}
}

func TestAC2DayRolloverMovesTheBucket(t *testing.T) {
	f := newFakeStore()
	ctx := context.Background()
	before := dayAt(2026, 11, 2, 23, 59)
	m := newMeter(t, f, before)
	day, err := m.PostTask(ctx, 100, 20, 500)
	if err != nil {
		t.Fatal(err)
	}
	if day != "2026-11-02" {
		t.Fatalf("booked into %q, want 2026-11-02 at 23:59+08", day)
	}
	// The next minute is a new day: the same meter, only the clock moves.
	after := dayAt(2026, 11, 3, 0, 0)
	m = newMeter(t, f, after) // same store, a clock one minute later
	if f.days["2026-11-03"].Tasks != 0 {
		t.Fatalf("2026-11-03 already has %d tasks before anything ran", f.days["2026-11-03"].Tasks)
	}
	if _, err := m.PostTask(ctx, 7, 3, 11); err != nil {
		t.Fatal(err)
	}
	if f.days["2026-11-02"].Tasks != 1 || f.days["2026-11-03"].Tasks != 1 {
		t.Errorf("rollover split wrong: 11-02=%d 11-03=%d, want 1 and 1",
			f.days["2026-11-02"].Tasks, f.days["2026-11-03"].Tasks)
	}
	if f.days["2026-11-03"].Micros != 11 {
		t.Errorf("new day micros = %d, want 11 (yesterday's 500 stayed put)", f.days["2026-11-03"].Micros)
	}
	today, err := m.Today(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if today.Micros != 11 || today.Day != "2026-11-03" {
		t.Errorf("Today after rollover = %+v, want {2026-11-03, 11 micros}", today)
	}
}

func TestAC2MonthRolloverResetsTheMonthRange(t *testing.T) {
	f := newFakeStore()
	ctx := context.Background()
	f.days["2026-11-30"] = DayTotal{Day: "2026-11-30", TokensIn: 50, TokensOut: 10, Micros: 777, Tasks: 2}
	m := newMeter(t, f, dayAt(2026, 11, 30, 12, 0))

	nov, err := m.MonthToDate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if nov.Micros != 777 || nov.Tasks != 2 {
		t.Errorf("November = %+v, want the 2026-11-30 row (777 micros / 2 tasks)", nov)
	}
	if nov.Day != "2026-11-01..2026-11-30" {
		t.Errorf("November range label = %q, want 2026-11-01..2026-11-30", nov.Day)
	}
	// One day later the month is a new one and the old row is history.
	m = newMeter(t, f, dayAt(2026, 12, 1, 0, 30)) // same store, next month
	dec, err := m.MonthToDate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if dec.Micros != 0 || dec.Day != "2026-12-01..2026-12-01" {
		t.Errorf("December = %+v, want the empty month starting 2026-12-01", dec)
	}
	// And November's row is still exactly what it was: a rollover reads, it
	// never rewrites history.
	if f.days["2026-11-30"].Micros != 777 {
		t.Errorf("2026-11-30 was rewritten by the rollover: %d", f.days["2026-11-30"].Micros)
	}
}

func TestAC2ReconcileAgreesAndDisagrees(t *testing.T) {
	f := newFakeStore()
	ctx := context.Background()
	now := dayAt(2026, 11, 2, 10, 0)
	m := newMeter(t, f, now)

	// Two tasks booked through the meter, and the same two in the base table.
	for _, task := range []struct{ in, out, micros int64 }{{100, 20, 500}, {30, 5, 90}} {
		f.plantBase(now, task.in, task.out, task.micros)
		if _, err := m.PostTask(ctx, task.in, task.out, task.micros); err != nil {
			t.Fatal(err)
		}
	}
	rec, err := m.Reconcile(ctx, "2026-11-02")
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Agree {
		t.Fatalf("reconcile = %+v, want agreement (aggregate and base both saw the same two tasks)", rec)
	}
	if rec.Aggregate.Micros != 590 || rec.Base.Micros != 590 || rec.Aggregate.Tasks != 2 {
		t.Errorf("numbers moved during a matching reconcile: agg=%+v base=%+v", rec.Aggregate, rec.Base)
	}

	// A third row appears ONLY in the base table: this is the shape of a task
	// whose row was written by someone who never went through the meter (ticket
	// 44's reconciliation cell exists to catch exactly this).
	f.plantBase(now, 1, 1, 42)
	rec, err = m.Reconcile(ctx, "2026-11-02")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Agree {
		t.Fatalf("reconcile reported agreement after a base-only row: %+v", rec)
	}
	if rec.Aggregate.Micros != 590 || rec.Base.Micros != 632 || rec.Base.Tasks != 3 {
		t.Errorf("disagreement has the wrong numbers: agg=%+v base=%+v", rec.Aggregate, rec.Base)
	}

	// The repair path: refresh recomputes the day from the base table.
	fixed, err := m.RefreshDay(ctx, "2026-11-02")
	if err != nil {
		t.Fatal(err)
	}
	if fixed.Micros != 632 || fixed.Tasks != 3 {
		t.Errorf("refresh returned %+v, want the base numbers (632 micros / 3 tasks)", fixed)
	}
	if f.days["2026-11-02"].Micros != 632 || f.days["2026-11-02"].Tasks != 3 {
		t.Errorf("refresh did not store the recomputed row: %+v", f.days["2026-11-02"])
	}
	rec, err = m.Reconcile(ctx, "2026-11-02")
	if err != nil {
		t.Fatal(err)
	}
	if !rec.Agree {
		t.Errorf("after RefreshDay the day still disagrees: %+v", rec)
	}
}

func TestAC2ReconcileRefusesWithoutBothMouths(t *testing.T) {
	f := newFakeStore()
	m, err := NewCostMeter(f) // no WithTasks, no WithWriter
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.Reconcile(context.Background(), "2026-11-02"); err == nil {
		t.Fatal("Reconcile without a task reader must not answer, it must say it cannot check")
	}
	if _, err := m.RefreshDay(context.Background(), "2026-11-02"); err == nil {
		t.Fatal("RefreshDay without a writer must refuse: recomputing and then overwriting nothing is not a repair")
	}
}

func TestPostTaskRefusesNegativesAndNamesTheDay(t *testing.T) {
	f := newFakeStore()
	ctx := context.Background()
	m := newMeter(t, f, dayAt(2026, 11, 2, 10, 0))
	if _, err := m.PostTask(ctx, -1, 0, 0); err == nil {
		t.Error("negative tokens booked anyway")
	}
	if _, err := m.PostTask(ctx, 0, 0, -5); err == nil {
		t.Error("negative micros booked anyway")
	}
	if len(f.bumps) != 0 {
		t.Errorf("a refused booking still wrote %v to the ledger", f.bumps)
	}
	// A ledger that fails is reported, with the day it tried.
	f.failBump = fmt.Errorf("disk full")
	day, err := m.PostTask(ctx, 1, 1, 1)
	if err == nil || !strings.Contains(err.Error(), "disk full") {
		t.Fatalf("ledger error lost: day=%q err=%v", day, err)
	}
	if day != "2026-11-02" {
		t.Errorf("failed booking did not name the day: %q", day)
	}
	// No ledger at all is a refusal, not a silent no-op.
	bare, err := NewCostMeter(f, WithLocation(fixedZone))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bare.PostTask(ctx, 1, 1, 1); err == nil {
		t.Error("a read-only meter must refuse to book rather than report success")
	}
}

func TestAC3BudgetBoundaries79And80And100(t *testing.T) {
	// cap = 1000 micros, alert at 80%.
	params := BudgetParams{DailyMicros: 1000, AlertBasisPoints: 8000, PauseAtLimit: true}
	cases := []struct {
		spent     int64
		wantLevel string
		wantAllow bool
		wantBasis int
	}{
		{789, LevelOK, true, 78},
		{799, LevelOK, true, 79},
		{800, LevelWarn, true, 80},
		{999, LevelWarn, true, 99},
		{1000, LevelOver, false, 100},
		{1500, LevelOver, false, 150},
	}
	for _, c := range cases {
		v := params.Judge(DayTotal{Micros: c.spent}, DayTotal{})
		if v.Level != c.wantLevel {
			t.Errorf("spent %d: level %q, want %q", c.spent, v.Level, c.wantLevel)
		}
		if v.AllowNewTask != c.wantAllow {
			t.Errorf("spent %d: AllowNewTask %v, want %v", c.spent, v.AllowNewTask, c.wantAllow)
		}
		if Percent(v.Binding.Basis) != c.wantBasis {
			t.Errorf("spent %d: percent %d, want %d", c.spent, Percent(v.Binding.Basis), c.wantBasis)
		}
		if c.wantLevel != LevelOK && v.Notice == "" {
			t.Errorf("spent %d: %s verdict with no user-visible notice", c.spent, c.wantLevel)
		}
		if c.wantLevel == LevelOK && v.Notice != "" {
			t.Errorf("spent %d: ok verdict still produced a notice %q", c.spent, v.Notice)
		}
	}
}

func TestAC3OverBudgetWarnModeNeverBlocksAndNamesItself(t *testing.T) {
	warn := BudgetParams{DailyMicros: 1000, AlertBasisPoints: 8000, PauseAtLimit: false}
	v := warn.Judge(DayTotal{Micros: 2000}, DayTotal{})
	if v.Level != LevelOver {
		t.Fatalf("level = %q, want over at 200%%", v.Level)
	}
	if !v.AllowNewTask {
		t.Error("over_budget=warn must still allow new tasks (that is the whole alternative)")
	}
	if !strings.Contains(v.Notice, "warn") && !strings.Contains(v.Notice, "不拦") {
		t.Errorf("warn-mode notice does not say the mode: %q", v.Notice)
	}
	pause := BudgetParams{DailyMicros: 1000, AlertBasisPoints: 8000, PauseAtLimit: true}
	pv := pause.Judge(DayTotal{Micros: 2000}, DayTotal{})
	if pv.AllowNewTask {
		t.Error("over_budget=pause at 200% must block new tasks")
	}
	if !strings.Contains(pv.Notice, "暂停新任务") {
		t.Errorf("pause notice must say 暂停新任务: %q", pv.Notice)
	}
	if !strings.Contains(pv.Notice, "运行中的任务不受影响") {
		t.Errorf("pause notice must say running tasks are unaffected: %q", pv.Notice)
	}
	if !strings.Contains(pv.Notice, "明天") {
		t.Errorf("pause notice must say when it recovers (resumable next day): %q", pv.Notice)
	}
}

func TestAC3BindingCapIsTheTighterOneAndBothAreListed(t *testing.T) {
	params := BudgetParams{DailyMicros: 1000, MonthlyMicros: 10000, AlertBasisPoints: 8000, PauseAtLimit: true}
	// Day at 90% (basis 9000), month at 50% (basis 5000): the day binds.
	v := params.Judge(DayTotal{Micros: 900}, DayTotal{Micros: 5000})
	if v.Binding.Kind != "daily" {
		t.Errorf("binding = %q, want daily (90%% beats 50%%)", v.Binding.Kind)
	}
	if v.Level != LevelWarn || !v.AllowNewTask {
		t.Errorf("verdict = %+v, want warn + allowed", v)
	}
	if len(v.Caps) != 2 {
		t.Fatalf("Caps has %d rows, want both budgets listed so a panel can show two lines", len(v.Caps))
	}
	// Now the month is the one past 100%: it binds, and the recovery wording
	// changes to the month's.
	v = params.Judge(DayTotal{Micros: 100}, DayTotal{Micros: 20000})
	if v.Binding.Kind != "monthly" {
		t.Errorf("binding = %q, want monthly", v.Binding.Kind)
	}
	if v.Level != LevelOver || v.AllowNewTask {
		t.Errorf("verdict = %+v, want over + blocked", v)
	}
	if !strings.Contains(v.Notice, "下月") {
		t.Errorf("monthly over notice must name the month recovery: %q", v.Notice)
	}
}

func TestAC3UnsetBudgetsAreUncappedNotBlocking(t *testing.T) {
	v := BudgetParams{}.Judge(DayTotal{Micros: math.MaxInt64 / 4}, DayTotal{Micros: math.MaxInt64 / 4})
	if v.Level != LevelOK || !v.AllowNewTask {
		t.Errorf("no caps configured produced %+v: an unset budget must not silence the product", v)
	}
	if len(v.Caps) != 0 {
		t.Errorf("unset caps still produced %d cap rows", len(v.Caps))
	}
	// A negative stored spend (a bug upstream) judges at 0 basis, never as a
	// huge percentage.
	v = BudgetParams{DailyMicros: 1000, AlertBasisPoints: 8000}.Judge(DayTotal{Micros: -5}, DayTotal{})
	if v.Level != LevelOK {
		t.Errorf("negative spend judged %q, want ok (never a false alarm, and the writer refuses negatives anyway)", v.Level)
	}
}

func TestBudgetFromRejectsBadConfig(t *testing.T) {
	if _, err := BudgetFrom(-1, 0, 0.8, "pause"); err == nil {
		t.Error("a negative daily budget loaded")
	}
	if _, err := BudgetFrom(0, 0, 1.5, "pause"); err == nil {
		t.Error("alert_threshold 1.5 loaded (validate.go range is [0,1])")
	}
	if _, err := BudgetFrom(0, 0, math.NaN(), "pause"); err == nil {
		t.Error("NaN alert_threshold loaded")
	}
	if _, err := BudgetFrom(0, 0, 0.8, "stop"); err == nil {
		t.Error("over_budget=stop loaded (the enum is pause|warn)")
	}
	p, err := BudgetFrom(1_000_000, 20_000_000, 0.8, "pause")
	if err != nil {
		t.Fatal(err)
	}
	if p.AlertBasisPoints != 8000 || !p.PauseAtLimit || p.DailyMicros != 1_000_000 {
		t.Errorf("BudgetFrom = %+v, want the [cost] leaves carried across with 0.8 -> 8000 bp", p)
	}
	w, err := BudgetFrom(0, 0, 0.8, "warn")
	if err != nil {
		t.Fatal(err)
	}
	if w.PauseAtLimit {
		t.Error("over_budget=warn produced PauseAtLimit=true")
	}
}

func TestBasisPointsSaturateInsteadOfWrapping(t *testing.T) {
	// spent*10000 overflows int64 here; a wrapped value would read as a small
	// percentage and let a task start past a budget that is obviously gone.
	b := basisPoints(math.MaxInt64/2, 1)
	if b != math.MaxInt {
		t.Errorf("basisPoints wrapped: %d, want saturation at MaxInt", b)
	}
	v := BudgetParams{DailyMicros: 1, AlertBasisPoints: 8000, PauseAtLimit: true}.
		Judge(DayTotal{Micros: math.MaxInt64 / 2}, DayTotal{})
	if v.AllowNewTask {
		t.Error("an overflow-sized spend was allowed past a 1-micro cap")
	}
	if Percent(v.Binding.Basis) < 100 {
		t.Errorf("display percent = %d, want a huge number rather than a wrapped small one", Percent(v.Binding.Basis))
	}
}

func TestDayBoundsUseTheInjectedCalendar(t *testing.T) {
	f := newFakeStore()
	// 2026-11-02 in +08:00 runs from 178... to +24h; check the bounds land the
	// base row on the right side of the boundary.
	m := newMeter(t, f, dayAt(2026, 11, 2, 12, 0))
	start, end, err := m.dayBounds("2026-11-02")
	if err != nil {
		t.Fatal(err)
	}
	if !start.Equal(dayAt(2026, 11, 2, 0, 0)) || !end.Equal(dayAt(2026, 11, 3, 0, 0)) {
		t.Errorf("bounds = %v..%v, want local midnight to local midnight", start, end)
	}
	if end.Sub(start) != 24*time.Hour {
		t.Errorf("day length = %v, want 24h in a fixed zone", end.Sub(start))
	}
	f.plantBase(dayAt(2026, 11, 2, 0, 0), 10, 1, 100)                // first second of the day: in
	f.plantBase(dayAt(2026, 11, 3, 0, 0).Add(-time.Second), 1, 1, 1) // last second: in
	f.plantBase(dayAt(2026, 11, 3, 0, 0), 1, 1, 1)                   // first second of the NEXT day: out
	// The DayReader mouth is bypassed here; the base sum alone proves the bounds.
	sum, err := m.Reconcile(context.Background(), "2026-11-02")
	if err != nil {
		t.Fatal(err)
	}
	if sum.Base.Micros != 101 || sum.Base.Tasks != 2 {
		t.Errorf("day bounds put the wrong rows inside: %+v, want 101 micros over 2 tasks", sum.Base)
	}
}

func TestMeterRefusesWithoutReader(t *testing.T) {
	if _, err := NewCostMeter(nil); err == nil {
		t.Fatal("a meter with nothing to read must not be built")
	}
}

func TestDayKeyIsTheShapeMemoryValidates(t *testing.T) {
	got := DayKey(dayAt(2026, 1, 9, 3, 4))
	if got != "2026-01-09" {
		t.Fatalf("DayKey = %q, want zero-padded 2026-01-09", got)
	}
	if !validDay(got) {
		t.Errorf("DayKey output %q does not satisfy memory's validDay shape", got)
	}
}
