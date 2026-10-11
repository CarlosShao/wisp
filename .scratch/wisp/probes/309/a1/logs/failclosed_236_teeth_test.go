package tools

// Ticket 236, AC#1 + AC#1b: the five fail-closed refusals that only exist as a
// line of code.
//
// What this file is for. Five branches refuse a call because a wiring is
// missing rather than because a judgement came out "no":
//
//	task.go:699-701            taskCancel, t.d.Roster == nil
//	task.go:707-709            taskCancel, caller == ""
//	subagent_197.go:253-255    task.spawn, t.d.Roster == nil
//	subagent_197.go:256-258    task.spawn, t.d.BaseOptions == nil || t.d.ParentTools == nil
//	subagent_197.go:259-263    task.spawn, parentID == ""
//
// Before this file, deleting any one of them changed NOTHING observable to the
// test suite (ticket 236's "现量" row 1, measured there by 221-v1 as m4/m5).
// The reason a leg that only asserts IsError == true has zero teeth is that
// every one of these branches has a NEIGHBOUR that also answers IsError, so
// fall-through keeps the call refused and a substring/IsError ruler still
// passes:
//
//	task cancel, roster branch gone -> TaskRoster.Look is nil-safe (task.go:269),
//	    so the call falls to "查不到这个任务" (task.go:711-715).
//	task cancel, caller branch gone  -> caller "" is neither the target nor its
//	    parent, so it falls to "不是调用者的孩子" (task.go:730-734).
//	task spawn, roster branch gone   -> TryAcquireSubagentSlot on nil answers
//	    ok=false (task.go:364-366), so it falls to "同时在跑的子代理已达上限".
//	task spawn, ParentTools half gone -> the child really derives with an empty
//	    tool directory, so the reply becomes a success line.
//	task spawn, BaseOptions half gone -> :278 CALLS t.d.BaseOptions two lines
//	    after the guard, so removing it is not a fall-through at all but a nil-func
//	    panic inside Execute. That is why the spawn legs dispatch through
//	    spawnGuarded: the panic is reported as ONE named red instead of killing the
//	    test binary and washing every reading that came after it (measured, 交件 §3
//	    m-1d). Recovering here can only turn "crash" into "red"; it never turns a
//	    red into a pass, and the text it hands back starts with "panic".
//	task spawn, parentID branch gone  -> the child really derives and the reply
//	    becomes a success line plus the 盖戳 notice (subagent_197.go:471-472).
//
// So every leg below pins the COMPLETE literal, byte for byte, by equality and
// not by Contains: the sentence is the only thing that identifies WHICH branch
// refused, which is the whole claim 票 221's security conclusion rests on
// ("a refusal is only evidence if you can tell which door it came from").
//
// ---------------------------------------------------------------------------
// Two carriers, chosen per branch level, and the difference is not cosmetic
// ---------------------------------------------------------------------------
//
// task.cancel is frozen at L1 (task.go:662, D34 row PLAN.md:2564), so a call
// only reaches Execute if the gate lets the L1 window run out. NoGate's
// PendingWindow answers AnswerReject ("L1 确认窗口尚未接入"), so under NoGate
// the refusal a leg observes comes from the GATE and never from the authority
// check - a false green, warned about verbatim at task_cancel_221_legs_test.go:13-19.
// These legs therefore use windowGate221(), whose PendingWindow answers
// AnswerTimeout, and bridge.go's L1 arm treats an unanswered window as EXECUTE
// (bridge.go:443-447, SPEC-06 §2).
//
// task.spawn is frozen at L0 (subagent_197.go:222), and L0 routes straight
// through (bridge.go:425-427) without consulting any gate - so the existing
// sub197 harness's NoGate is safe HERE, and only here. Anyone reading this file
// must not generalise it: an L0 spawn under NoGate proves nothing about an L1
// call under NoGate.
//
// ---------------------------------------------------------------------------
// Instrument note, pinned here because it is what makes these legs measurable
// ---------------------------------------------------------------------------
//
// All five legs are BEHAVIOURAL: each dispatches one real call through the real
// bridge and reads the reply back. None of them opens a source file. That is
// deliberate, because 读盘型 rulers are structurally blind to the instrument
// this ticket measures with: `go test -overlay` replaces the bytes the
// COMPILER sees, while os.ReadFile (e.g. task_cancel_221_legs_test.go:224)
// hits the physical disk, so an overlay mutation of a lexical claim is
// invisible to the ruler that claims it. A lexical leg here could not be shown
// to have teeth at all. These five can.

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
)

// The five literals, verbatim as the production code writes them. Each one is
// the whole Text field of exactly one branch, so pinning it pins the branch.
const (
	lit236CancelRosterUnwired = "任务名册未接线（fail-closed：拒绝停掉任何任务——名册才是唯一的停法）"

	lit236CancelNoHostTaskID = "这条调用没有宿主给的任务 id（fail-closed：不知道是谁要停，只能拒）"

	lit236SpawnRosterUnwired = "任务名册未接线（fail-closed：拒绝派生子代理，派生了也没有地方登记它）"

	lit236SpawnNoAssembly = "宿主没有给出派生用的装配（BaseOptions/ParentTools 未接线，" +
		"fail-closed：不起第二套运行时）"

	lit236SpawnNoHostTaskID = "这条调用没有宿主给的任务 id（fail-closed：不知道父任务是谁，" +
		"就没法登记父子关系，也没法把结论盖戳进父任务的 C25 作用域）"
)

// AC#1, task.go:707-709. The caller's identity is the host's, not an argument
// (bridge.go:269-271 only falls back CorrelationID -> TaskID, so passing both
// empty is how production delivers "nobody told us who this is"), so the target
// is a REAL derived child rather than an invented id: the row must exist for the
// leg to prove it was the caller check that spoke, not the lookup.
func Test236R2TaskCancelRefusesWhenHostGaveNoCallerID(t *testing.T) {
	x := build221(t, windowGate221(), func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, false)
	res := x.h.spawn(t, t.Context(), label197, "读完就好")
	if res.IsError {
		t.Fatalf("派生失败，本枚用例没有对象可停：%s", res.Text)
	}
	target := x.h.childRow(t).TaskID

	// Same carrier 票 236 §三 names: x.cancel with an empty caller id.
	callerUnknown := x.cancel(t, "", target)
	if !callerUnknown.IsError {
		t.Fatalf("caller 为空时 task.cancel 回了非拒绝：%q", callerUnknown.Text)
	}
	if callerUnknown.Text != lit236CancelNoHostTaskID {
		t.Errorf("task.go:707 那一支的完整字面量没被原样说出（got %q, want %q）——"+
			"若 got 是「不是调用者的孩子」那一支，说明 caller == \"\" 的拒绝分支已经不在了",
			callerUnknown.Text, lit236CancelNoHostTaskID)
	}
}

// AC#1, task.go:699-701. Carrier is the one 票 236 §三 points at:
// mustRegisterTaskEntries(TaskDeps{Roster: nil}) registers the whole
// BuiltinTaskEntries family, task.cancel included, so the call reaches the real
// tool with the real unwired deps. The sibling shape at
// task_output_leg_test.go:126-138 nails task.output's roster branch and is not
// touched here; the roster sentence's PREFIX (任务名册未接线) is shared by five
// production sites, so a Contains("未接线") ruler cannot tell them apart.
func Test236R2TaskCancelRefusesWhenRosterIsUnwired(t *testing.T) {
	bare := New(Options{
		Registry: mustRegisterTaskEntries(t, TaskDeps{Roster: nil}),
		Gate:     windowGate221(),
		Logf:     func(string, ...any) {},
	})
	// TaskID non-empty so the caller check is satisfied and the roster branch is
	// the one that speaks, in this branch's own order (it sits ahead of it).
	out, err := bare.Execute(t.Context(), agent.ToolRequest{
		TaskID: "caller-236r2", CorrelationID: "caller-236r2", CallID: "call-236r2-cancel-no-roster",
		Name: "task.cancel",
		Args: json.RawMessage(`{"task_id":"target-236r2"}`),
	})
	if err != nil {
		t.Fatalf("task.cancel dispatch failed: %v", err)
	}
	if !out.IsError {
		t.Fatalf("名册没接线时 task.cancel 回了非拒绝：%q", out.Text)
	}
	if out.Text != lit236CancelRosterUnwired {
		t.Errorf("task.go:699 那一支的完整字面量没被原样说出（got %q, want %q）——"+
			"若 got 是「查不到这个任务」那一支，说明 Roster == nil 的拒绝分支已经不在了",
			out.Text, lit236CancelRosterUnwired)
	}
}

// AC#1b, subagent_197.go:253-255. The mutate hook of buildWith
// (subagent_197_test.go:175-197) is the carrier already in the repo - the same
// hook 票 197 used to exercise the "provenance engine not wired" branch at
// :719-725 - so this leg adds no harness of its own.
func Test236R2TaskSpawnRefusesWhenRosterIsUnwired(t *testing.T) {
	h := newSub197Harness(t)
	h.buildWith(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true, func(d *SubagentDeps) { d.Roster = nil })

	res := spawnGuarded(t, h, "名册没接线", "读一段")
	if !res.IsError {
		t.Fatalf("名册没接线时 task.spawn 回了非拒绝：%q", res.Text)
	}
	if res.Text != lit236SpawnRosterUnwired {
		t.Errorf("subagent_197.go:253 那一支的完整字面量没被原样说出（got %q, want %q）——"+
			"若 got 是「同时在跑的子代理已达上限」，说明 Roster == nil 的拒绝分支已经不在了",
			res.Text, lit236SpawnRosterUnwired)
	}
}

// AC#1b, subagent_197.go:256-258, the ParentTools half of the ||.
func Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired(t *testing.T) {
	h := newSub197Harness(t)
	h.buildWith(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true, func(d *SubagentDeps) { d.ParentTools = nil })

	res := spawnGuarded(t, h, "父工具面没接线", "读一段")
	if !res.IsError {
		t.Errorf("ParentTools 没接线时 task.spawn 回了非拒绝：%q（guard 不在时子代理照样派生，"+
			"只是拿到一个空目录）", res.Text)
	}
	if res.Text != lit236SpawnNoAssembly {
		t.Errorf("subagent_197.go:256 那一支的完整字面量没被原样说出（got %q, want %q）——"+
			"这一支是 || 的 ParentTools 半边", res.Text, lit236SpawnNoAssembly)
	}
}

// AC#1b, subagent_197.go:256-258, the BaseOptions half of the ||. Same sentence,
// other half - and the half where "guard gone" is a CRASH, not a fall-through,
// because :278 calls t.d.BaseOptions() unconditionally two lines later. Measured
// (see 交件 §3 m-1d): deleting the whole guard with BaseOptions nil panics inside
// Execute, which kills the test BINARY and washes every reading after it. So the
// dispatch below recovers, and this leg reports the crash as its own red instead
// of taking the rest of the package down with it.
func Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired(t *testing.T) {
	h := newSub197Harness(t)
	h.buildWith(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true, func(d *SubagentDeps) { d.BaseOptions = nil })

	res := spawnGuarded(t, h, "装配快照没接线", "读一段")
	if !res.IsError {
		t.Errorf("BaseOptions 没接线时 task.spawn 回了非拒绝：%q", res.Text)
	}
	if res.Text != lit236SpawnNoAssembly {
		t.Errorf("subagent_197.go:256 那一支的完整字面量没被原样说出（got %q, want %q）——"+
			"这一支是 || 的 BaseOptions 半边；got 以「panic」开头就是那道 guard 已经不在了",
			res.Text, lit236SpawnNoAssembly)
	}
}

// spawnGuarded dispatches one task.spawn and turns a panic into an IsError
// result whose text starts with "panic:", so a guard that has been deleted reads
// as ONE named red rather than as an aborted run. It recovers only here, only in
// test code, and only to keep the rest of the package measurable; nothing it
// returns is treated as a pass.
func spawnGuarded(t *testing.T, h *sub197Harness, label, prompt string) (out Result) {
	t.Helper()
	defer func() {
		if rec := recover(); rec != nil {
			out = Result{Text: fmt.Sprintf("panic（guard 摘掉后不是拒绝而是崩）: %v", rec), IsError: true}
		}
	}()
	return h.spawn(t, t.Context(), label, prompt)
}

// AC#1b, subagent_197.go:259-263. The real bridge, with both id fields empty
// (subagent_197_test.go:217 hardcodes parent197, so it cannot carry this leg;
// 票 236 §三 says the same about the cancel side - the fix is to hand the
// request an empty host id, never to invent one inside the tool).
func Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID(t *testing.T) {
	h := newSub197Harness(t)
	h.build(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true)

	res := spawnWithoutHostTaskID(t, h)
	if !res.IsError {
		// Errorf, not Fatalf: the roster check below must still run, because the
		// row this branch prevents is the part of the claim that survives a rewrite
		// of the wording.
		t.Errorf("宿主没给任务 id 时 task.spawn 回了非拒绝：%q", res.Text)
	}
	if res.Text != lit236SpawnNoHostTaskID {
		t.Errorf("subagent_197.go:259 那一支的完整字面量没被原样说出（got %q, want %q）——"+
			"若 got 是「子代理 … 已结束」，说明 parentID == \"\" 的拒绝分支已经不在了",
			res.Text, lit236SpawnNoHostTaskID)
	}
	// The consequence that branch exists to prevent: a row nobody is the parent
	// of. Red under the same mutation, and it is the half a sentence-only leg
	// would miss if somebody rewrote the wording but kept the hole. TaskRoster
	// deliberately has no full-table read port (Descendants walks ParentTaskID
	// only and a parentless row is by design nobody's descendant, task.go:302-304),
	// so the orphan signature is the DIFFERENCE between the two existing ports:
	// a row in the table that is not anybody's descendant.
	if n := h.roster.Count(); n != 1 {
		t.Errorf("名册行数 = %d, want 1（只有 harness 那枚根任务）——"+
			"宿主没给父任务 id 时不许派生、更不许登记：%v", n, h.roster.Descendants(parent197))
	}
}

// spawnWithoutHostTaskID dispatches one task.spawn exactly as production does
// for a call whose host stamped no task id at all: both id fields empty, because
// bridge.go:269-271's only fall-back is CorrelationID -> TaskID. Package-unexported,
// so no new name leaves the package. Guarded the same way spawnGuarded is, so a
// mutant that turns this dispatch into a crash still reports one red.
func spawnWithoutHostTaskID(t *testing.T, h *sub197Harness) (out Result) {
	t.Helper()
	defer func() {
		if rec := recover(); rec != nil {
			out = Result{Text: fmt.Sprintf("panic（guard 摘掉后不是拒绝而是崩）: %v", rec), IsError: true}
		}
	}()
	res, err := h.bridge.Execute(t.Context(), agent.ToolRequest{
		TaskID: "", CorrelationID: "", CallID: "call-236r2-spawn-no-parent",
		Name: "task.spawn",
		Args: json.RawMessage(`{"description":"没有父任务","prompt":"读一段就好"}`),
	})
	if err != nil {
		t.Fatalf("task.spawn dispatch failed: %v", err)
	}
	return Result{Text: res.Text, IsError: res.IsError}
}
