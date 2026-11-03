package tools

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// The D34 reminder family (ticket 31): one-shot, in-process, notify-only
// reminders — the deliberately narrow exception to the no-scheduler rule (D12).
//
// PLAN.md:2568 / SPEC-07 §3:68 freeze the row verbatim:
//
//	| `reminder.set/list/cancel` | 一次性提醒（**只发通知**） | L1 / L0 / L1 | `notify` | **S4** |
//
// so the档位 matrix here is set L1 / list L0 / cancel L1, and the family's
// declared capability face is the one token that row names: `notify`.
//
// ⛔ HARD EXCLUSIONS (票 31「NOT」清单 ＝ PLAN.md:2607-2611 ＋ SPEC-12 §5 的
// REJECTED 行「周期性调度器（cron）| D12 有意取舍；一次性 reminder 已移出为窄
// 例外（16.5.4）」). Each one is refused in THREE places on purpose — the
// parameter schema (the key does not exist and an unknown key is rejected),
// this comment, and the Description() the model reads:
//
//  1. ⛔ recurring / cron / interval reminders — that is the scheduler D12 keeps
//     REJECTED; there is no recurrence column in reminders either (schema v6).
//  2. ⛔ action-on-fire ("提醒我并顺手整理文件" → 只做前半句). fire is a fixed,
//     already-requested, L0 `notify` and NOTHING else: no tool execution, no
//     LLM call. There is no action/payload column to store one in.
//  3. ⛔ OS-level scheduling (registry / Task Scheduler) — that is L2 with
//     uninstall residue (PLAN.md:2611).
//  4. ⛔ missed-reminder replay across restarts — boot reports a COUNT only
//     ("你错过了 N 条提醒") and never re-delivers per item (PLAN.md:2610).
//     注意 这四条是 REJECTED，不是 DEFERRED：AGENTS.md §1.1 requires any
//     DEFERRED(D-xx) marker to match SPEC-12 §5 1:1, and none of them is registered there.
//
// What the timer layer (internal/agent) owes and this file does NOT do:
// the wait itself. reminder.set hands a stored row to the ReminderTimers seam
// and nothing here sleeps, spins a goroutine, or reads a clock: the monotonic
// discipline (D42#9) lives behind that seam, in one named+owned goroutine.
// ---------------------------------------------------------------------------

// ReminderRowStore is the reminders DAO subset the tools read and write. It is
// satisfied by *memory.Store directly; it is an interface here for the same
// reason grant.go's GrantSource is one: internal/tools has never held a database
// handle (grant.go:16-20 states that rule), and a narrow face is what keeps a
// reminder tool from reaching any table but this one.
type ReminderRowStore interface {
	InsertReminder(ctx context.Context, r memory.Reminder) error
	ReminderByID(ctx context.Context, id string) (memory.Reminder, error)
	ListReminders(ctx context.Context) ([]memory.Reminder, error)
	MarkReminderFired(ctx context.Context, id string) (bool, error)
	CancelReminder(ctx context.Context, id string) (bool, error)
}

// ReminderTimers is the process-side half of ticket 31: the layer that owns the
// one-shot wait for an armed row. *agent.ReminderTimers satisfies it.
//
// Arm never blocks (the wait runs on that package's named+owned goroutine), and
// its fire path is the notify-only path this family promises: Notifier.Notify
// + MarkReminderFired + one tool_call audit row, through the seams given to the
// timer at construction. Disarm is what reminder.cancel calls so the row and
// the wait stop together; a row that is only disarmed, or only cancelled in the
// table, is the half-cancel this ticket forbids.
type ReminderTimers interface {
	// Arm starts the one-shot wait for one already-stored row. A row whose
	// wall-clock target is already past is NOT armed: that is ticket 31's
	// boot rule (past-due rows are counted, never fired).
	Arm(r memory.Reminder) error
	// Disarm drops the wait for one id. False means nothing was waiting
	// (already fired, cancelled, or never armed in this process).
	Disarm(id string) bool
	// Pending reports how many waits this process is holding. reminder.list
	// reads it to say which rows are armed but NOT waiting (a row can be armed
	// in the table and un-waited here: armed before this process started).
	Pending() int
}

// ReminderDeps is what the reminder.* tools need from the host. Either seam
// being nil is answered fail-closed in Execute, never by pretending the
// reminder was set.
type ReminderDeps struct {
	// Store is the reminders table. nil = no table wired.
	Store ReminderRowStore
	// Timers is the process-side one-shot wait. nil = nothing can fire.
	Timers ReminderTimers
	// MintID mints a reminder id. nil = the crypto/rand default below; a test
	// injects a deterministic mint so an id-shaped assertion is not a
	// randomness assertion.
	MintID func() (string, error)
}

// reminderIDPrefix is the row's id shape: 'rem_' + crypto/rand hex, the same
// mint shape C31 uses for session ids (memory.Reminder's own doc comment).
const reminderIDPrefix = "rem_"

// mintReminderID is the default ReminderDeps.MintID.
func mintReminderID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("tools: mint reminder id: %w", err)
	}
	return reminderIDPrefix + hex.EncodeToString(b[:]), nil
}

func (d ReminderDeps) mint() (string, error) {
	if d.MintID != nil {
		return d.MintID()
	}
	return mintReminderID()
}

// reminderExclusions is the refusal sentence every one of these tools carries:
// the four hard exclusions of 票 31, in the order the ticket lists them, plus
// the honest limit 16.5.4 requires the user to be told at set-time. It is one
// constant so set/cancel/list and the schema rejections cannot say four
// different things about the same boundary.
const reminderExclusions = "只做一次性提醒：不做周期/cron、不做到点执行任何动作" +
	"（「提醒我并顺手整理文件」只做前半句）、不做 OS 计划任务、不做重启补发。" +
	"局限：进程没在跑就不会触发。"

// reminderCaveat is the set-time disclosure sentence 16.5.4 requires, kept as
// its own constant because ticket 31's AC#5 asserts THIS string in the set
// confirmation and nothing else: PLAN.md:2612-2613 is verbatim
// 「**进程没在跑就不会触发**……必须在设置提醒时告知用户，不得静默失败」.
// 注意 措辞以 PLAN.md 权威行为准：票 31 票面第 21 行是转述（「进程未运行则不会触发」），
// 台账 A1420 已就地订正为「按 PLAN 权威」，本腿不拿转述串当判据。
const reminderCaveat = "进程没在跑就不会触发"

// rejectedReminderKeys maps every parameter key name that would ask for one of
// the four hard exclusions to the sentence that refuses it. The schema already
// rejects unknown keys (additionalProperties:false); this table is the SECOND
// belt and the one that can name WHICH boundary was asked for, so a model that
// guessed a plausible-looking field gets an answer it can act on (D37: a
// failure the model can self-correct goes back as IsError text, not a Go error).
var rejectedReminderKeys = map[string]string{
	"recurrence": "拒绝周期性提醒：一次性提醒是唯一落点（SPEC-12 §5 REJECTED：周期性调度器/cron）。",
	"cron":       "拒绝 cron：一次性提醒是唯一落点（SPEC-12 §5 REJECTED：周期性调度器/cron）。",
	"repeat":     "拒绝重复提醒：一次性提醒是唯一落点（SPEC-12 §5 REJECTED：周期性调度器/cron）。",
	"interval":   "拒绝间隔/重复提醒：一次性提醒是唯一落点（SPEC-12 §5 REJECTED：周期性调度器/cron）。",
	"every":      "拒绝「每 N 分钟」这类周期：一次性提醒是唯一落点（SPEC-12 §5 REJECTED）。",
	"action":     "拒绝到点执行动作：fire 只发通知，不执行任何工具、不调 LLM（「提醒我并顺手整理文件」只做前半句）。",
	"payload":    "拒绝到点执行动作：fire 只发通知，不执行任何工具、不调 LLM（PLAN.md:2609）。",
	"tool":       "拒绝到点调用工具：fire 只发通知（PLAN.md:2609 的窄例外只覆盖 notify）。",
	"task":       "拒绝写 OS 计划任务：本工具不起常驻调度、不碰注册表/任务计划（PLAN.md:2611）。",
	"schedule":   "拒绝常驻调度：只有进程内一次性定时器（PLAN.md:2611、SPEC-12 §5 REJECTED）。",
	"replay":     "拒绝重启补发：错过的提醒只在 boot 报条数，不逐条弹（PLAN.md:2610）。",
}

// rejectReminderKeys answers the exclusion refusal for a raw parameter object,
// or "" when every key is one this tool actually takes. It reads the raw JSON
// keys rather than a decoded struct precisely because a struct cannot see the
// key a caller invented: the rejection has to name the boundary that was asked
// for, or the model would read "unknown field" as "try another spelling".
func rejectReminderKeys(params json.RawMessage) string {
	if len(params) == 0 {
		return ""
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(params, &raw); err != nil {
		return "" // the caller's own decode error is the tool's next answer
	}
	var hits []string
	for k := range raw {
		if why, bad := rejectedReminderKeys[k]; bad {
			hits = append(hits, why)
		}
	}
	if len(hits) == 0 {
		return ""
	}
	// Stable order so the refusal text is comparable across runs.
	seen := map[string]bool{}
	var out []string
	for _, k := range sortedRejectedHits(hits) {
		if !seen[k] {
			seen[k] = true
			out = append(out, k)
		}
	}
	return strings.Join(out, " ")
}

func sortedRejectedHits(hits []string) []string {
	out := append([]string(nil), hits...)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// reminderArgErrors names the two argument shapes that are not a capability
// refusal but still cannot be armed: a missing/blank note (the row requires it)
// and a target that is not a positive unix-seconds instant.
func reminderArgErrors(fireAt int64, note string) string {
	switch {
	case note == "":
		return "缺少 note：提醒正文是通知里唯一给用户看的内容（reminders.note NOT NULL）"
	case fireAt <= 0:
		return "缺少 fire_at_wallclock_utc：提醒目标必须是正的 unix 秒时刻"
	}
	return ""
}

// ---------------------------------------------------------------------------
// reminder.set (L1, PLAN.md:2568 / SPEC-07 §3:68)
// ---------------------------------------------------------------------------

type reminderSet struct{ d ReminderDeps }

type reminderSetArgs struct {
	FireAtWallclockUTC int64  `json:"fire_at_wallclock_utc"`
	Note               string `json:"note"`
}

// reminderSetSchema has exactly two properties and additionalProperties:false.
// ⛔ There is deliberately NO recurrence / cron / action / payload / tool /
// schedule / replay property: those are 票 31's four hard exclusions, and the
// schema is where the ticket's "rejected by schema" AC is decided (AC#3).
var reminderSetSchema = JSONSchema(`{"type":"object","properties":{` +
	`"fire_at_wallclock_utc":{"type":"integer","minimum":1,"description":"到点时刻（墙钟 unix 秒，UTC）。等待按单调时钟算，系统时钟跳变不移动它"},` +
	`"note":{"type":"string","description":"提醒正文：到点原样推给用户的一句话，不是指令、也不会被当成要做的事"}` +
	`},"required":["fire_at_wallclock_utc","note"],"additionalProperties":false}`)

func (reminderSet) Name() string { return "reminder.set" }

// Description carries the honest limit at the face the model reads (16.5.4),
// and the exclusion list, because the ticket requires the refusal to be
// discoverable without a round-trip: a model that only ever reads this line
// must not be able to believe a recurring or action-carrying reminder exists.
func (reminderSet) Description() string {
	return "设一条一次性提醒：到点只发一条系统通知。" + reminderExclusions +
		" 只接受 fire_at_wallclock_utc 与 note 两个参数，多余的键（周期/cron/动作/补发）一律拒。"
}

func (reminderSet) Parameters() JSONSchema { return reminderSetSchema }

// Execute implements Tool: store the row, then hand it to the timer seam.
//
// Order matters and is the ticket's AC#2 shape: the row is written FIRST, so a
// timer that cannot be armed (nil seam, a past-due target) still leaves an
// honest armed row for the next boot's missed count, and the reply says which
// of the two happened rather than claiming a fire that cannot come.
func (t reminderSet) Execute(ctx context.Context, params json.RawMessage, onUpdate func(delta string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if why := rejectReminderKeys(params); why != "" {
		return Result{Text: why + " " + reminderExclusions, IsError: true}, nil
	}
	var a reminderSetArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	note := strings.TrimSpace(a.Note)
	if why := reminderArgErrors(a.FireAtWallclockUTC, note); why != "" {
		return Result{Text: why, IsError: true}, nil
	}
	if t.d.Store == nil {
		return Result{Text: "提醒表未接线（fail-closed：没有 reminders 表就不假装设过提醒）", IsError: true}, nil
	}

	id, err := t.d.mint()
	if err != nil {
		return Result{Text: "铸造提醒 id 失败：" + err.Error(), IsError: true}, nil
	}
	r := memory.Reminder{
		ID:     id,
		FireAt: a.FireAtWallclockUTC,
		Note:   note,
		Status: memory.ReminderArmed,
	}
	if err := t.d.Store.InsertReminder(ctx, r); nil != err {
		return Result{Text: "写入提醒失败：" + err.Error(), IsError: true}, nil
	}

	waited := false
	if t.d.Timers != nil {
		// A refusal here is not a failed set: the row exists, and the reply
		// says the wait is not held. The common cause is a target already past.
		if err := t.d.Timers.Arm(r); err != nil {
			return Result{Text: fmt.Sprintf(
				"提醒 %s 已写进表，但本次没有挂上计时器：%v。这条提醒在本进程内不会触发；"+
					"若它到点前进程还在跑，下一次启动只会把它计入「你错过了 N 条提醒」，不补发。"+
					"局限：%s", id, err, reminderCaveat), IsError: true}, nil
		}
		waited = true
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("reminder.set 已设提醒 %s", id))
	}

	text := fmt.Sprintf("已设一次性提醒 %s：到点只发一条通知，正文是「%s」。", id, note)
	if !waited {
		// The timer seam was never wired: say it, do not imply a fire.
		text += "注意 计时器未接线，这条提醒只是写进了表，本进程内不会触发。"
	}
	// 16.5.4's honesty rule: the limit is spoken at set-time, every time, and
	// AC#5 asserts this exact string is present in this output.
	text += " 局限：" + reminderCaveat + "（不做 OS 级调度，重启后也只报错过的条数、不逐条补发）。"
	text += " " + reminderExclusions
	return Result{Text: text}, nil
}

// ReminderSetDecl declares the set row: L1 with the notify face the frozen
// D34 row names.
func ReminderSetDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapNotify},
		Needs:        []Capability{CapNotify},
		Declared:     risk.L1,
		PathParams:   nil,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// ---------------------------------------------------------------------------
// reminder.list (L0)
// ---------------------------------------------------------------------------

type reminderList struct{ d ReminderDeps }

// reminderListSchema takes nothing. There is no "also show fired" flag because
// the row's vocabulary is a roster, not a filter, and a list that could render
// history as a pending reminder is the double-promise this tool must not make.
var reminderListSchema = JSONSchema(`{"type":"object","properties":{},"additionalProperties":false}`)

func (reminderList) Name() string { return "reminder.list" }

func (reminderList) Description() string {
	return "列出还没到点的提醒（armed）。已触发/已取消的历史不显示——提醒不是任务队列。" + reminderExclusions
}

func (reminderList) Parameters() JSONSchema { return reminderListSchema }

// Execute implements Tool.
func (t reminderList) Execute(ctx context.Context, params json.RawMessage, onUpdate func(delta string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if t.d.Store == nil {
		return Result{Text: "提醒表未接线（fail-closed：说不出有哪些提醒）", IsError: true}, nil
	}
	rows, err := t.d.Store.ListReminders(ctx)
	if err != nil {
		return Result{Text: "读提醒表失败：" + err.Error(), IsError: true}, nil
	}

	// The table is ordered by fire-at already (ListReminders' ORDER BY), so the
	// display order is inherited, not re-sorted here.
	var armed []memory.Reminder
	for _, r := range rows {
		if r.Status == memory.ReminderArmed {
			armed = append(armed, r)
		}
	}
	if len(armed) == 0 {
		return Result{Text: "没有待触发的提醒。"}, nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "待触发提醒 %d 条", len(armed))
	if t.d.Timers != nil {
		// Honest split: rows can be armed in the table and NOT waited on here
		// (armed by a previous process). Saying "都会触发" about those would be
		// the silent-failure shape 16.5.4 forbids.
		fmt.Fprintf(&b, "（本进程挂着 %d 条）", t.d.Timers.Pending())
	}
	b.WriteString("：\n")
	for _, r := range armed {
		fmt.Fprintf(&b, "  %s · %s · %s\n", r.ID, wallClockLabel(r.FireAt), r.Note)
	}
	b.WriteString("提醒只发通知，不执行任何动作。")
	return Result{Text: strings.TrimRight(b.String(), "\n")}, nil
}

// wallClockLabel renders a stored wall-clock target for a user-visible line.
// It is display only: no interval is ever computed from a value produced here
// (observe/clock.go's rule — wall readings are for records and display).
func wallClockLabel(unixSec int64) string {
	return time.Unix(unixSec, 0).UTC().Format("2006-01-02 15:04:05 UTC")
}

// ReminderListDecl declares the list row: L0 and no C3 face (it reads one table
// the host gave it, the same shape task.output declares at task.go:586).
func ReminderListDecl() Decl {
	return Decl{
		Capabilities: nil,
		Needs:        nil,
		Declared:     risk.L0,
		PathParams:   nil,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// ---------------------------------------------------------------------------
// reminder.cancel (L1)
// ---------------------------------------------------------------------------

type reminderCancel struct{ d ReminderDeps }

type reminderCancelArgs struct {
	ID string `json:"id"`
}

var reminderCancelSchema = JSONSchema(`{"type":"object","properties":{` +
	`"id":{"type":"string","description":"要取消的提醒 id（reminder.set/list 给出的那个 rem_…，不是路径）"}` +
	`},"required":["id"],"additionalProperties":false}`)

func (reminderCancel) Name() string { return "reminder.cancel" }

func (reminderCancel) Description() string {
	return "取消一条还没到点的提醒：停掉进程内的等待，并把那一行标成 cancelled（单向，取消不了的只有已触发/已取消的）。" +
		reminderExclusions
}

func (reminderCancel) Parameters() JSONSchema { return reminderCancelSchema }

// Execute implements Tool: stop the wait first, then stamp the row.
//
// That order is the one-way-transition guard: the DAO refuses to stamp a row
// that is no longer armed, and stopping the wait first means a timer cannot
// deliver a notification for a reminder the user just cancelled. The DAO's own
// late-timer guard (MarkReminderFired affects 0 rows once cancelled) is the
// second belt, so a fire already in flight cannot re-stamp it either.
func (t reminderCancel) Execute(ctx context.Context, params json.RawMessage, onUpdate func(delta string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a reminderCancelArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	id := strings.TrimSpace(a.ID)
	if id == "" {
		return Result{Text: "缺少 id 参数", IsError: true}, nil
	}
	if t.d.Store == nil {
		return Result{Text: "提醒表未接线（fail-closed：不动任何提醒）", IsError: true}, nil
	}

	// Read the row first: cancel of a row that does not exist, and cancel of a
	// row that already fired, are different facts and must be said apart.
	cur, err := t.d.Store.ReminderByID(ctx, id)
	if err != nil {
		return Result{Text: fmt.Sprintf("查不到提醒 %s：%v", id, err), IsError: true}, nil
	}
	if cur.Status != memory.ReminderArmed {
		return Result{Text: fmt.Sprintf(
			"提醒 %s 已经是 %s，取消不了（单向终态）。已触发的那条是历史，不是可以重挂的把手。",
			id, cur.Status), IsError: true}, nil
	}

	dropped := false
	if t.d.Timers != nil {
		dropped = t.d.Timers.Disarm(id)
	}
	applied, err := t.d.Store.CancelReminder(ctx, id)
	if err != nil {
		return Result{Text: "取消提醒失败：" + err.Error(), IsError: true}, nil
	}
	if !applied {
		// Lost the race with the timer's fire path: the notification may already
		// be out. Say so instead of reporting a clean cancel.
		return Result{Text: fmt.Sprintf(
			"没能取消 %s：它在这一瞬间已经不是 armed 状态（多半正好触发了）。计时器侧已停手=%v。",
			id, dropped), IsError: true}, nil
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("reminder.cancel 已取消 %s", id))
	}
	return Result{Text: fmt.Sprintf(
		"已取消提醒 %s（正文「%s」）：等待已停、那一行标成 cancelled，之后不会再触发。", id, cur.Note)}, nil
}

// ReminderCancelDecl declares the cancel row: L1 per the frozen D34 pair.
func ReminderCancelDecl() Decl {
	return Decl{
		Capabilities: []Capability{CapNotify},
		Needs:        []Capability{CapNotify},
		Declared:     risk.L1,
		PathParams:   nil,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// BuiltinReminderEntries returns the reminder family: set / list / cancel at
// the L1 / L0 / L1 matrix SPEC-07 §3:68 and PLAN.md:2568 freeze.
//
// The registry note (Description on every entry) names the four hard exclusions
// — recurring/cron, action-on-fire, OS scheduling, restart replay — so nothing
// downstream has to read this file to learn what the family refuses to do.
func BuiltinReminderEntries(d ReminderDeps) []Entry {
	return []Entry{
		{Tool: reminderSet{d: d}, Decl: ReminderSetDecl()},
		{Tool: reminderList{d: d}, Decl: ReminderListDecl()},
		{Tool: reminderCancel{d: d}, Decl: ReminderCancelDecl()},
	}
}
