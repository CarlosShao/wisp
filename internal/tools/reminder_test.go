package tools

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 31 acceptance for the tools layer: AC#1's set/cancel/list surface,
// AC#3's exclusion rejections, AC#5's caveat string. The timer's own behaviour
// (fire/late/clock-jump/missed-count) is covered in internal/agent, so nothing
// here sleeps.

// fakeTimers is a ReminderTimers seam that records calls without waiting.
type fakeTimers struct {
	mu      sync.Mutex
	armed   []memory.Reminder
	disarm  []string
	armErr  error
	pending int
}

func (f *fakeTimers) Arm(r memory.Reminder) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.armErr != nil {
		return f.armErr
	}
	f.armed = append(f.armed, r)
	f.pending++
	return nil
}

func (f *fakeTimers) Disarm(id string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, a := range f.armed {
		if a.ID == id {
			f.armed = append(f.armed[:i], f.armed[i+1:]...)
			f.pending--
			f.disarm = append(f.disarm, id)
			return true
		}
	}
	return false
}

func (f *fakeTimers) Pending() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.pending
}

func (f *fakeTimers) armedIDs() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	var ids []string
	for _, a := range f.armed {
		ids = append(ids, a.ID)
	}
	return ids
}

// reminderFixture wires the three tools against a REAL reminders table (schema
// v6) plus a recording timer seam and a deterministic id mint.
func reminderFixture(t *testing.T) (*memory.Store, *fakeTimers, ReminderDeps) {
	t.Helper()
	store := openStore(t, t.TempDir())
	tm := &fakeTimers{}
	n := 0
	deps := ReminderDeps{
		Store:  store,
		Timers: tm,
		MintID: func() (string, error) {
			n++
			return "rem_test" + string(rune('0'+n)), nil
		},
	}
	return store, tm, deps
}

// registerReminderFamily proves the entries pass C4 registration (the same
// contract face every family is pinned on).
func registerReminderFamily(t *testing.T, d ReminderDeps) *Registry {
	t.Helper()
	reg := NewRegistry()
	for _, e := range BuiltinReminderEntries(d) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	return reg
}

func callTool(t *testing.T, tool Tool, params string) Result {
	t.Helper()
	res, err := tool.Execute(context.Background(), json.RawMessage(params), nil)
	if err != nil {
		t.Fatalf("Execute(%s): host error %v", tool.Name(), err)
	}
	return res
}

// ---- AC#1: set / list / cancel ---------------------------------------------

func Test31SetWritesRowArmsWaitAndCarriesCaveat(t *testing.T) {
	store, tm, deps := reminderFixture(t)
	reg := registerReminderFamily(t, deps)
	setEntry, ok := reg.Lookup("reminder.set")
	if !ok {
		t.Fatal("reminder.set 没注册上")
	}

	target := time.Now().UTC().Add(2 * time.Hour).Unix()
	res := callTool(t, setEntry.Tool, `{"fire_at_wallclock_utc":`+itoa64(target)+`,"note":"下楼拿快递"}`)
	if res.IsError {
		t.Fatalf("set 报错了: %s", res.Text)
	}

	// AC#5: the honest caveat is required in the set confirmation (16.5.4).
	// ⚠ 断言用 PLAN.md:2612 的权威串，不用票面转述（台账 A1420 已裁"按 PLAN 权威"）。
	if !strings.Contains(res.Text, "进程没在跑就不会触发") {
		t.Fatalf("set 确认里没有局限句（PLAN.md:2612 逐字）:\n%s", res.Text)
	}
	if !strings.Contains(res.Text, reminderExclusions) {
		t.Fatalf("set 确认里没有拒绝清单句:\n%s", res.Text)
	}

	// The row landed, armed, with the given target.
	row, err := store.ReminderByID(context.Background(), "rem_test1")
	if err != nil {
		t.Fatalf("ReminderByID: %v", err)
	}
	if row.Status != memory.ReminderArmed {
		t.Errorf("status = %q, want armed", row.Status)
	}
	if row.FireAt != target {
		t.Errorf("fire_at = %d, want %d", row.FireAt, target)
	}
	if row.Note != "下楼拿快递" {
		t.Errorf("note = %q", row.Note)
	}

	// The process-side wait was armed for the SAME row.
	ids := tm.armedIDs()
	if len(ids) != 1 || ids[0] != "rem_test1" {
		t.Fatalf("计时器侧 armed = %v, want [rem_test1]", ids)
	}
}

func Test31ListShowsPendingOnly(t *testing.T) {
	store, _, deps := reminderFixture(t)
	ctx := context.Background()
	future := time.Now().UTC().Add(time.Hour).Unix()

	// armed / fired / cancelled: only the armed row may show.
	if err := store.InsertReminder(ctx, memory.Reminder{ID: "rem_p1", FireAt: future, Note: "待触发", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("seed armed: %v", err)
	}
	if err := store.InsertReminder(ctx, memory.Reminder{ID: "rem_p2", FireAt: future + 1, Note: "已触发", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("seed to-fire: %v", err)
	}
	if _, err := store.MarkReminderFired(ctx, "rem_p2"); err != nil {
		t.Fatalf("mark fired: %v", err)
	}
	if err := store.InsertReminder(ctx, memory.Reminder{ID: "rem_p3", FireAt: future + 2, Note: "已取消", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("seed to-cancel: %v", err)
	}
	if _, err := store.CancelReminder(ctx, "rem_p3"); err != nil {
		t.Fatalf("cancel seed: %v", err)
	}

	reg := NewRegistry()
	for _, e := range BuiltinReminderEntries(deps) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	listEntry, _ := reg.Lookup("reminder.list")
	res := callTool(t, listEntry.Tool, `{}`)
	if res.IsError {
		t.Fatalf("list 报错了: %s", res.Text)
	}
	if !strings.Contains(res.Text, "rem_p1") || !strings.Contains(res.Text, "待触发") {
		t.Fatalf("list 没列出 armed 行:\n%s", res.Text)
	}
	if strings.Contains(res.Text, "rem_p2") {
		t.Fatalf("list 把 fired 行当待触发展示了:\n%s", res.Text)
	}
	if strings.Contains(res.Text, "rem_p3") {
		t.Fatalf("list 把 cancelled 行当待触发展示了:\n%s", res.Text)
	}
}

func Test31ListEmptyIsQuiet(t *testing.T) {
	_, _, deps := reminderFixture(t)
	deps.Timers = nil // no waits at all
	e := reminderList{d: deps}
	res := callTool(t, e, `{}`)
	if res.IsError {
		t.Fatalf("空表不该报错: %s", res.Text)
	}
	if !strings.Contains(res.Text, "没有待触发的提醒") {
		t.Fatalf("空表文案: %q", res.Text)
	}
}

func Test31CancelRemovesWaitAndStampsRow(t *testing.T) {
	store, tm, deps := reminderFixture(t)
	ctx := context.Background()
	future := time.Now().UTC().Add(time.Hour).Unix()
	if err := store.InsertReminder(ctx, memory.Reminder{ID: "rem_c1", FireAt: future, Note: "取消我", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := tm.Arm(memory.Reminder{ID: "rem_c1", FireAt: future, Note: "取消我", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("pre-arm: %v", err)
	}

	e := reminderCancel{d: deps}
	res := callTool(t, e, `{"id":"rem_c1"}`)
	if res.IsError {
		t.Fatalf("cancel 报错了: %s", res.Text)
	}
	if tm.Pending() != 0 {
		t.Fatalf("计时器侧还在等 %d 条，cancel 没停掉等待", tm.Pending())
	}
	row, err := store.ReminderByID(ctx, "rem_c1")
	if err != nil {
		t.Fatalf("ReminderByID: %v", err)
	}
	if row.Status != memory.ReminderCancelled {
		t.Fatalf("status = %q, want cancelled（等待与行必须一起停）", row.Status)
	}
}

func Test31CancelOfFiredRowIsRefused(t *testing.T) {
	store, _, deps := reminderFixture(t)
	ctx := context.Background()
	future := time.Now().UTC().Add(time.Hour).Unix()
	if err := store.InsertReminder(ctx, memory.Reminder{ID: "rem_c2", FireAt: future, Note: "已触发", Status: memory.ReminderArmed}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, err := store.MarkReminderFired(ctx, "rem_c2"); err != nil {
		t.Fatalf("mark fired: %v", err)
	}

	res := callTool(t, reminderCancel{d: deps}, `{"id":"rem_c2"}`)
	if !res.IsError {
		t.Fatalf("取消 fired 行居然成功: %s", res.Text)
	}
	if !strings.Contains(res.Text, memory.ReminderFired) {
		t.Errorf("拒绝句该点名那一行的真实状态: %s", res.Text)
	}
	if !strings.Contains(res.Text, "单向终态") {
		t.Errorf("拒绝句该说清单向终态: %s", res.Text)
	}
}

func Test31CancelUnknownIDIsLoud(t *testing.T) {
	_, _, deps := reminderFixture(t)
	res := callTool(t, reminderCancel{d: deps}, `{"id":"rem_nope"}`)
	if !res.IsError {
		t.Fatal("不存在的 id 竟被当成取消成功")
	}
	// "查不到" and "取消不了" are different facts (票 164 的 fail-closed 形状).
	if !strings.Contains(res.Text, "查不到") {
		t.Fatalf("未知 id 的拒绝句: %s", res.Text)
	}
}

// ---- AC#3: exclusions rejected by schema ------------------------------------

func Test31SchemaRejectsRecurrence(t *testing.T) {
	_, _, deps := reminderFixture(t)
	deps.MintID = func() (string, error) { return "rem_x", nil }
	e := reminderSet{d: deps}

	for _, body := range []string{
		`{"fire_at_wallclock_utc":` + itoa64(time.Now().UTC().Add(time.Hour).Unix()) + `,"note":"x","recurrence":"daily"}`,
		`{"fire_at_wallclock_utc":` + itoa64(time.Now().UTC().Add(time.Hour).Unix()) + `,"note":"x","cron":"0 9 * * *"}`,
		`{"fire_at_wallclock_utc":` + itoa64(time.Now().UTC().Add(time.Hour).Unix()) + `,"note":"x","every":"1h"}`,
	} {
		res := callTool(t, e, body)
		if !res.IsError {
			t.Fatalf("周期性参数 %s 竞被接受: %s", body, res.Text)
		}
		if !strings.Contains(res.Text, "周期性") && !strings.Contains(res.Text, "周期") {
			t.Errorf("拒绝句没点名周期性: %s", res.Text)
		}
	}

	// The schema itself must not even carry the concept.
	for _, banned := range []string{"recurrence", "cron", "repeat", "interval", "every"} {
		if strings.Contains(string(reminderSetSchema), `"`+banned+`"`) {
			t.Errorf("reminder.set 的 schema 里出现了 %q：那等于给了一个不存在的把手", banned)
		}
	}
	if !strings.Contains(string(reminderSetSchema), `"additionalProperties":false`) {
		t.Error("reminder.set schema 缺 additionalProperties:false：未知键会被静默吞掉")
	}
}

func Test31SchemaRejectsActionOnFire(t *testing.T) {
	_, _, deps := reminderFixture(t)
	deps.MintID = func() (string, error) { return "rem_x", nil }
	e := reminderSet{d: deps}

	// 「提醒我并顺手整理文件」的第二半句必须被拒。
	for _, body := range []string{
		`{"fire_at_wallclock_utc":` + itoa64(time.Now().UTC().Add(time.Hour).Unix()) + `,"note":"x","action":"整理文件"}`,
		`{"fire_at_wallclock_utc":` + itoa64(time.Now().UTC().Add(time.Hour).Unix()) + `,"note":"x","tool":"fs.write"}`,
		`{"fire_at_wallclock_utc":` + itoa64(time.Now().UTC().Add(time.Hour).Unix()) + `,"note":"x","payload":"run ls"}`,
	} {
		res := callTool(t, e, body)
		if !res.IsError {
			t.Fatalf("到点执行参数 %s 竞被接受: %s", body, res.Text)
		}
		if !strings.Contains(res.Text, "只做前半句") && !strings.Contains(res.Text, "不执行任何工具") {
			t.Errorf("拒绝句没点名「只做前半句」: %s", res.Text)
		}
	}

	for _, banned := range []string{"action", "payload", "tool", "command"} {
		if strings.Contains(string(reminderSetSchema), `"`+banned+`"`) {
			t.Errorf("reminder.set 的 schema 里出现了 %q：fire 只发通知，不该有执行把手", banned)
		}
	}
}

// Test31ExclusionsAreStatedAtTheRegistryFace: the four refusals must be
// discoverable from what the model reads, not only from a Go comment.
func Test31ExclusionsAreStatedAtTheRegistryFace(t *testing.T) {
	_, _, deps := reminderFixture(t)
	for _, e := range BuiltinReminderEntries(deps) {
		desc := e.Tool.Description()
		if !strings.Contains(desc, "一次性") {
			t.Errorf("%s 的 description 没写清一次性: %s", e.Tool.Name(), desc)
		}
		if !strings.Contains(desc, reminderCaveat) && e.Tool.Name() == "reminder.set" {
			t.Errorf("reminder.set 的 description 没带 16.5.4 的局限句: %s", desc)
		}
	}
}

// ---- AC#5: the caveat string -------------------------------------------------

// Test31SetConfirmationCarriesCaveatVerbatim is AC#5. The string is PLAN.md:2612
// verbatim, NOT the ticket's paraphrase (ledger A1420 ruled 按 PLAN 权威).
func Test31SetConfirmationCarriesCaveatVerbatim(t *testing.T) {
	if reminderCaveat != "进程没在跑就不会触发" {
		t.Fatalf("caveat 串漂了: %q（权威串＝PLAN.md:2612 逐字「进程没在跑就不会触发」）", reminderCaveat)
	}
	store, _, deps := reminderFixture(t)
	ctx := context.Background()
	target := time.Now().UTC().Add(90 * time.Minute).Unix()

	e := reminderSet{d: deps}
	res := callTool(t, e, `{"fire_at_wallclock_utc":`+itoa64(target)+`,"note":"三点开会"}`)
	if res.IsError {
		t.Fatalf("set 失败: %s", res.Text)
	}
	if !strings.Contains(res.Text, reminderCaveat) {
		t.Fatalf("确认文案里没有局限句:\n%s", res.Text)
	}
	// And it is a real row, not just a promise.
	if _, err := store.ReminderByID(ctx, "rem_test1"); err != nil {
		t.Fatalf("确认说了却没写行: %v", err)
	}
}

// Test31ArmedButNotWaitedIsSaidApart: a row that cannot be armed (timer seam
// down) must NOT be reported as a set reminder. 16.5.4: 不得静默失败.
func Test31ArmedButNotWaitedIsSaidApart(t *testing.T) {
	_, _, deps := reminderFixture(t)
	deps.MintID = func() (string, error) { return "rem_nowait", nil }
	e := reminderSet{d: deps}
	res := callTool(t, e, `{"fire_at_wallclock_utc":`+itoa64(time.Now().UTC().Add(time.Hour).Unix())+`,"note":"x"}`)
	if res.IsError {
		t.Fatalf("set 失败: %s", res.Text)
	}
	if !strings.Contains(res.Text, reminderCaveat) {
		t.Errorf("仍要带局限句: %s", res.Text)
	}

	// Now a timer seam whose Arm fails.
	tm := &fakeTimers{armErr: context.Canceled}
	deps2 := deps
	deps2.Timers = tm
	deps2.MintID = func() (string, error) { return "rem_armfail", nil }
	res2 := callTool(t, reminderSet{d: deps2}, `{"fire_at_wallclock_utc":`+itoa64(time.Now().UTC().Add(time.Hour).Unix())+`,"note":"y"}`)
	if !res2.IsError {
		t.Fatalf("Arm 失败却报成功: %s", res2.Text)
	}
	if !strings.Contains(res2.Text, "没有挂上计时器") {
		t.Errorf("该说清写进了表但没挂等待: %s", res2.Text)
	}
}

// ---- unwired seams are fail-closed -------------------------------------------

func Test31ReminderToolsRejectMissingStore(t *testing.T) {
	res := callTool(t, reminderSet{d: ReminderDeps{}}, `{"fire_at_wallclock_utc":1,"note":"x"}`)
	if !res.IsError || !strings.Contains(res.Text, "未接线") {
		t.Fatalf("没接表却像设成功: %s", res.Text)
	}
	res = callTool(t, reminderList{d: ReminderDeps{}}, `{}`)
	if !res.IsError || !strings.Contains(res.Text, "未接线") {
		t.Fatalf("没接表却像列成功: %s", res.Text)
	}
	res = callTool(t, reminderCancel{d: ReminderDeps{}}, `{"id":"rem_x"}`)
	if !res.IsError || !strings.Contains(res.Text, "未接线") {
		t.Fatalf("没接表却像取消成功: %s", res.Text)
	}
}

func Test31SetValidatesArgs(t *testing.T) {
	_, _, deps := reminderFixture(t)
	e := reminderSet{d: deps}

	res := callTool(t, e, `{"fire_at_wallclock_utc":0,"note":"x"}`)
	if !res.IsError || !strings.Contains(res.Text, "fire_at_wallclock_utc") {
		t.Errorf("零时刻该被拒并点名: %s", res.Text)
	}
	res = callTool(t, e, `{"fire_at_wallclock_utc":`+itoa64(time.Now().UTC().Add(time.Hour).Unix())+`,"note":"   "}`)
	if !res.IsError || !strings.Contains(res.Text, "note") {
		t.Errorf("空正文该被拒并点名: %s", res.Text)
	}
}

// ---- the C1/C4 declaration face (D34 matrix L1/L0/L1, notify) ----------------

func Test31ReminderDeclarationsMatchD34Row(t *testing.T) {
	cases := []struct {
		name  string
		decl  Decl
		level risk.Level
		caps  []Capability
	}{
		{"reminder.set", ReminderSetDecl(), risk.L1, []Capability{CapNotify}},
		{"reminder.list", ReminderListDecl(), risk.L0, nil},
		{"reminder.cancel", ReminderCancelDecl(), risk.L1, []Capability{CapNotify}},
	}
	for _, c := range cases {
		if c.decl.Declared != c.level {
			t.Errorf("%s 声明档位 = %v, want %v（SPEC-07 §3:68 冻结 L1/L0/L1）", c.name, c.decl.Declared, c.level)
		}
		have := map[Capability]bool{}
		for _, cp := range c.decl.Capabilities {
			have[cp] = true
		}
		for _, cp := range c.caps {
			if !have[cp] {
				t.Errorf("%s 没声明 %s（D34 那一行的能力列只有 notify）", c.name, cp)
			}
		}
		for _, n := range c.decl.Needs {
			if !have[n] {
				t.Errorf("%s needs 未声明能力 %s（C3 不对称执法：注册期就该拒）", c.name, n)
			}
		}
	}
}

func Test31ReminderFamilyRegistersAndNamesAreC1(t *testing.T) {
	_, _, deps := reminderFixture(t)
	reg := registerReminderFamily(t, deps)
	want := []string{"reminder.cancel", "reminder.list", "reminder.set"}
	var got []string
	for _, e := range reg.List() {
		got = append(got, e.Tool.Name())
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("注册名册 = %v, want %v（三枚，按名排序）", got, want)
	}
}

func itoa64(n int64) string { return strconv.FormatInt(n, 10) }
