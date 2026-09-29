package tools

// Ticket 221 甲形, the behavioural legs (批准＝台账 A434; AC#1's resident ruler
// lives in task_cancel_221_test.go).
//
// What these legs refuse to repeat is what the ticket was filed for: take a
// promise in a description and call it a capability. So every leg below
// dispatches through the REAL bridge (real C19 routing, real roster, real
// agent.Loop, real memory.Store journal) with a REAL host-minted task id as the
// caller identity. The only fake is the C5 provider, a frozen injection seam
// (AGENTS §1.3); nothing here stands in for a real component.
//
// ⚠ Why these legs do NOT reuse subagent_197_test.go's harness bridge: that one
// runs with tools.NoGate, whose PendingWindow answers AnswerReject - so a
// task.cancel call (D34 freezes it at L1, PLAN.md:2564) would be refused by the
// GATE before Execute. A refusal from the gate is not a refusal from the
// authority check, and a "refused" leg that never reached Execute would be a
// false green for AC#3. So build221 takes its gate per leg, and
// Test221TaskCancelUnderNoGateStopsAtTheWindow pins that consequence on purpose.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// harness221 drives the 197 harness's roster / stream feed / taint engine
// through a bridge this file owns, so gate and sinks are chosen per leg.
type harness221 struct {
	h *sub197Harness

	mu    sync.Mutex
	audit []string

	mem *memory.Store // nil when the leg books no tool_call rows
}

func (x *harness221) logLine(format string, args ...any) {
	x.mu.Lock()
	defer x.mu.Unlock()
	x.audit = append(x.audit, fmt.Sprintf(format, args...))
}

func (x *harness221) auditText() string {
	x.mu.Lock()
	defer x.mu.Unlock()
	return strings.Join(x.audit, "\n")
}

// windowGate221 answers an L1 window the way SPEC-06 §2 says an unopposed window
// answers: it runs out, and running out MEANS EXECUTE (bridge.go's L1 arm).
func windowGate221() Gate {
	return GateFuncs{
		Window: func(context.Context, Decision) (Answer, string) { return AnswerTimeout, "" },
		Approval: func(context.Context, Decision) (Answer, string) {
			return AnswerReject, "本用例不带 L2 队列"
		},
	}
}

// build221 wires task.output + task.cancel + task.spawn over one roster through
// the real bridge. provider() runs once per spawn, so one held fake197Provider
// can back several children.
func build221(t *testing.T, gate Gate, provider func() llm.LlmProvider, withJournal bool) *harness221 {
	t.Helper()
	x := &harness221{h: newSub197Harness(t)}
	if gate == nil {
		gate = NoGate{}
	}
	var j agent.Journal
	if withJournal {
		x.mem = openStore(t, t.TempDir())
		mustStartTask(t, x.mem, parent197)
		j = x.mem
	}
	deps := SubagentDeps{
		Roster:      x.h.roster,
		Provenance:  x.h.prov,
		ParentTools: x.h.dir,
		Stream:      sub197StreamFeed{h: x.h},
		BaseOptions: func() (agent.Options, bool) {
			opt := agent.Options{Provider: provider(), Tools: x.h.dir}
			// AdmitTask must be wired or the spawner refuses the child outright
			// (197's 子代理永不自批 branch). This leg is about stopping, not
			// about approval, and the hook stays the host's single one queue.
			opt.AdmitTask = func(string) func() { return func() {} }
			return opt, true
		},
	}
	reg := NewRegistry()
	for _, e := range append(BuiltinTaskEntries(TaskDeps{Roster: x.h.roster}), BuiltinSubagentEntries(deps)...) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	x.h.bridge = New(Options{
		Registry: reg, Provenance: x.h.prov, Gate: gate, Journal: j,
		Logf: x.logLine,
	})
	return x
}

// cancel dispatches one task.cancel call AS callerID. The caller's identity is
// the one the bridge stamps into the context, which is why no parameter can name
// a different parent.
func (x *harness221) cancel(t *testing.T, callerID, targetID string) Result {
	t.Helper()
	out, err := x.h.bridge.Execute(t.Context(), agent.ToolRequest{
		TaskID: callerID, CorrelationID: callerID, CallID: "call-221-cancel-" + callerID + "-" + targetID,
		Name: "task.cancel",
		Args: json.RawMessage(`{"task_id":"` + targetID + `"}`),
	})
	if err != nil {
		t.Fatalf("task.cancel dispatch failed: %v", err)
	}
	return Result{Text: out.Text, IsError: out.IsError}
}

// spawnHeld derives one child and holds it inside its first model call, so the
// leg controls the child's liveness with channels instead of a clock (the shape
// subagent_197_test.go's AC#5 leg uses).
func (x *harness221) spawnHeld(t *testing.T, label string, started chan struct{}) (childID string, parked chan Result) {
	t.Helper()
	done := make(chan Result, 1)
	go func() {
		res, err := x.h.spawnResult(t.Context(), label, "长任务")
		if err != nil {
			res = Result{Text: "dispatch error: " + err.Error(), IsError: true}
		}
		done <- res
	}()
	select {
	case <-started:
	case <-t.Context().Done():
		t.Fatalf("子代理 %s 没有开始它的模型调用", label)
	}
	// By now the row exists twice over, and neither publisher needs this leg to
	// wait: the spawn goroutine filed it right after RunAsync, and the child's own
	// admission hook files the identical row before its first model call, which is
	// the call that just sent started.
	for _, r := range x.h.childRows() {
		if r.Out.Label == label {
			return r.TaskID, done
		}
	}
	t.Fatalf("名册里没有 %s 那一行：%+v", label, x.h.childRows())
	return "", done
}

// heldProvider is the C5 seam parked inside its first call.
func heldProvider(started chan struct{}, release chan struct{}) func() llm.LlmProvider {
	p := &fake197Provider{name: "child", answer: foreign197A, started: started, release: release}
	return func() llm.LlmProvider { return p }
}

// ---------------------------------------------------------------------------
// AC#1 / AC#2 - the registration itself, at the level D34 froze
// ---------------------------------------------------------------------------

func Test221TaskCancelIsRegisteredAtItsFrozenLevel(t *testing.T) {
	entries := BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster()})
	names := make([]string, 0, len(entries))
	var cancelEntry *Entry
	for i, e := range entries {
		names = append(names, e.Tool.Name())
		if e.Tool.Name() == "task.cancel" {
			cancelEntry = &entries[i]
		}
	}
	if cancelEntry == nil {
		t.Fatalf("BuiltinTaskEntries 没注册 task.cancel（现名册：%v）—— 说明书又在许诺一枚不存在的工具", names)
	}
	if len(names) != 2 {
		t.Errorf("task 家族注册了 %d 枚：%v（task.list 仍须是 DEFERRED，本票不许顺手注册它）", len(names), names)
	}
	decl := cancelEntry.Decl
	if decl.Declared != risk.L1 {
		t.Errorf("task.cancel 声明 = %v, want L1（D34 那行 PLAN.md:2564 冻的是 L0 / L1，cancel 取 L1）", decl.Declared)
	}
	if len(decl.Capabilities) != 0 || len(decl.Needs) != 0 {
		t.Errorf("task.cancel 的能力声明 = %v / %v, want 两枚都空（那行的第四列是「—」）", decl.Capabilities, decl.Needs)
	}
	if len(decl.PathParams) != 0 {
		t.Errorf("task.cancel 的 PathParams = %v, want 空：唯一的模型输入是名册 id，不是路径", decl.PathParams)
	}

	// The registration is not decorative: a call that reaches Execute does reach
	// TaskRoster.Cancel. Here the child has already joined, so the ONLY honest
	// answer is "没有在跑" - not a success that stopped nothing.
	x := build221(t, windowGate221(), func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, false)
	res := x.h.spawn(t, t.Context(), label197, "读完就好")
	if res.IsError {
		t.Fatalf("派生失败：%s", res.Text)
	}
	row := x.h.childRow(t)
	if st := row.Out.State; st != statemachine.StateSettling {
		t.Fatalf("跑完的孩子 state = %q, want Settling", string(st))
	}
	late := x.cancel(t, parent197, row.TaskID)
	if !late.IsError || !strings.Contains(late.Text, "没有在跑") {
		t.Errorf("停一枚已经结束的孩子回了 %q, want 报「没有在跑」", late.Text)
	}
}

// Test221DeferredMarkerForCancelLiftedButListStillMarked is AC#4 in the rewritten
// form (票面 ⑤: this ticket does NOT claim a two-way 1:1 with SPEC-12 §5, whose
// instrument is ticket 225's and which carries no task.* row at all). What it
// claims is its own marker: (a) task.cancel's DEFERRED line is lifted, (b) lifted
// in the same landing as the wiring - which this file proves by asserting the
// wiring and the marker in one leg - and task.list's line is untouched.
func Test221DeferredMarkerForCancelLiftedButListStillMarked(t *testing.T) {
	src, err := os.ReadFile("D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/221/v1/mutations/m14-cancel-marked-again.go")
	if err != nil {
		t.Fatalf("读不到 task.go（这把尺的落点是包内源文件）：%v", err)
	}
	var cancelMarked, listMarked int
	for _, line := range strings.Split(string(src), "\n") {
		if !strings.Contains(line, "DEFERRED") {
			continue
		}
		if strings.Contains(line, "task.cancel") {
			cancelMarked++
		}
		if strings.Contains(line, "task.list") {
			listMarked++
		}
	}
	if cancelMarked != 0 {
		t.Errorf("task.go 里还有 %d 行把 task.cancel 标成 DEFERRED：它已经注册并有实现，"+
			"「在册＋无实现＋无人认领」正是票 164 AC#1 要杀的那一形", cancelMarked)
	}
	if listMarked == 0 {
		t.Errorf("task.list 的 DEFERRED 标记不见了：本票不许摘它（那半支仍然没有实现）")
	}
}

// Test221SpawnDescriptionPromisesOnlyWhatIsTrue is 乙形: the three promise shapes
// the ticket names (task.spawn 的说明文字那句无条件许诺、它后面那个括号、父任务不等了的
// 回执那半句) must not survive, and the two literals other nails own must.
func Test221SpawnDescriptionPromisesOnlyWhatIsTrue(t *testing.T) {
	desc := (subagentSpawn{}).Description()
	for _, gone := range []string{
		"可以用 task.cancel 单独停它", // :196 的裸许诺
		"（要停它得单独停）",            // :197 那个括号尾巴
	} {
		if strings.Contains(desc, gone) {
			t.Errorf("说明文字还留着 %q：那一句在改前是不真的话（票 221 标题）", gone)
		}
	}
	for _, keep := range []string{
		"停掉父任务不会级联",   // subagent_197_test.go:797 钉的字面量，不是装饰
		"task.cancel", // 甲形之后这句话才允许出现
		fmt.Sprintf("上限 %d 枚", MaxConcurrentSubagents),
		fmt.Sprintf("深度 %d", MaxSubagentDepth),
	} {
		if !strings.Contains(desc, keep) {
			t.Errorf("说明文字丢了 %q：%q", keep, desc)
		}
	}
	// The bounded form of the promise is what has to be there in its place: the
	// model must be able to tell from this text alone that it cannot stop a
	// sibling or itself.
	for _, must := range []string{"只有派生它的那枚父任务", "停兄弟", "停自己都一律被拒"} {
		if !strings.Contains(desc, must) {
			t.Errorf("删掉裸许诺之后，边界那一半也没写出来（缺 %q）：%q", must, desc)
		}
	}
}

// ---------------------------------------------------------------------------
// AC#2 - the stopped child's row and stream reach a terminal state, and the
//        trail says who stopped whom (A434 item 4: 审计可查, no new on-screen field)
// ---------------------------------------------------------------------------

func Test221ParentStopsItsOwnChildRowAndStreamSettle(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	x := build221(t, windowGate221(), heldProvider(started, release), true)
	t.Cleanup(func() { close(release) })

	childID, parked := x.spawnHeld(t, label197, started)
	key := SubagentStreamKey(childID)

	// A sibling row that no cancel of mine may touch, and a watcher on it: the
	// "只动那一行" claim is measured, not asserted.
	x.h.roster.PublishSubagent("sibling-221", parent197, "隔壁那枚")
	sibWatch := x.h.roster.WatchRow("sibling-221")
	watch := x.h.roster.WatchRow(childID)

	res := x.cancel(t, parent197, childID)
	if res.IsError {
		t.Fatalf("父任务停自己的孩子被拒了：%s", res.Text)
	}
	for _, want := range []string{childID, parent197, "只动了这一行", "停掉父任务不会级联", "这次停止由 " + parent197 + " 发起"} {
		if !strings.Contains(res.Text, want) {
			t.Errorf("回执里没有 %q（这句话是写给模型看的）：%q", want, res.Text)
		}
	}

	// The row settles through its own write stream - no sleep, no clock.
	var final TaskOutput
	for final.State == "" || final.State == statemachine.State(subagentStateRunning) {
		select {
		case o, open := <-watch:
			if !open {
				t.Fatal("行监视通道被关掉，那一行没有回报终态")
			}
			final = o
		case <-t.Context().Done():
			t.Fatalf("停掉的子代理 %s 那一行没有落到终态：%+v", childID, final)
		}
	}
	if final.State != statemachine.State(subagentStateStopped) {
		t.Errorf("那一行的 state = %q, want %q（被宿主停掉那一支借的是 D43 的 Muted，本票不新造态名）",
			string(final.State), string(statemachine.StateMuted))
	}
	if !statemachine.Valid(final.State) {
		t.Errorf("state %q 不在 D43 的 20 枚名字上", string(final.State))
	}
	rec := x.h.mustLook(t, childID)
	if rec.Out.ParentTaskID != parent197 || rec.Out.Kind != TaskKindSubagent || rec.Out.Label != label197 {
		t.Errorf("那一行消失了或被改写了：parent=%q kind=%q label=%q（A434 item 3：停了以后行不许消失，三家别家都留行）",
			rec.Out.ParentTaskID, rec.Out.Kind, rec.Out.Label)
	}
	if !x.h.streamClosed(key) {
		t.Errorf("子代理 %s 的流键 %q 没有收到终态：那一页会一直显示在流", childID, key)
	}

	// The parked task.spawn call answers with that terminal state too - the
	// parent's model reads it, it is not left holding a spawn that never ends.
	select {
	case spawnRes := <-parked:
		if spawnRes.IsError {
			t.Errorf("父任务的 task.spawn 回了错误：%s", spawnRes.Text)
		}
		if !strings.Contains(spawnRes.Text, string(statemachine.StateMuted)) {
			t.Errorf("task.spawn 的结论里没有那句终态名（%s）：%q", statemachine.StateMuted, spawnRes.Text)
		}
	case <-t.Context().Done():
		t.Error("父任务那次 task.spawn 在孩子被停之后没有收口")
	}

	// Nothing else moved.
	select {
	case moved := <-sibWatch:
		t.Errorf("停这一枚把隔壁也写了：%+v", moved)
	default:
	}
	if sib := x.h.mustLook(t, "sibling-221"); sib.Out.State != statemachine.StateThinking {
		t.Errorf("隔壁那枚被动了：state = %q", string(sib.Out.State))
	}

	// 被谁停的 = 现成审计 sink 里查得到（A434 item 4：这一维只进审计，不加上屏字段）。
	// Two records, both already produced by the bridge for every call:
	//   - the audit line, naming the CALLER's task id and the tool;
	//   - the persisted tool_call row, whose ArgsJSON names the TARGET.
	audit := x.auditText()
	if !strings.Contains(audit, "tool=task.cancel") || !strings.Contains(audit, "task="+parent197) {
		t.Errorf("审计日志里没有这一停（缺 tool=task.cancel 或 task=%s）：\n%s", parent197, audit)
	}
	rows, err := x.mem.ListToolCallsByTask(t.Context(), parent197)
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	var booked bool
	for _, r := range rows {
		if r.Tool == "task.cancel" && strings.Contains(r.ArgsJSON, childID) {
			if r.TaskID != parent197 {
				t.Errorf("tool_call 行的调用者是 %q, want %q", r.TaskID, parent197)
			}
			booked = true
		}
	}
	if !booked {
		t.Errorf("tool_call 账上没有一条「谁停了谁」的记录（task.cancel＋%s）：%+v", childID, rows)
	}
}

// ---------------------------------------------------------------------------
// AC#3 - the two positive controls, plus the non-cascade nail
// ---------------------------------------------------------------------------

func Test221SubagentCannotStopSiblingOrItself(t *testing.T) {
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	x := build221(t, windowGate221(), heldProvider(started, release), false)
	t.Cleanup(func() { close(release) })

	aID, parkA := x.spawnHeld(t, "甲号孩子", started)
	bID, parkB := x.spawnHeld(t, "乙号孩子", started)

	// 正控①：子代理停兄弟。
	sib := x.cancel(t, aID, bID)
	if !sib.IsError {
		t.Errorf("子代理停兄弟被放行了（AC#3 的正控①）：%s", sib.Text)
	}
	for _, want := range []string{bID, aID, parent197, "只有父任务能停自己的孩子"} {
		if !strings.Contains(sib.Text, want) {
			t.Errorf("停兄弟的拒绝理由里没有 %q（要指名是谁家的孩子）：%q", want, sib.Text)
		}
	}
	// 正控②：子代理停自己。
	self := x.cancel(t, bID, bID)
	if !self.IsError {
		t.Errorf("子代理停自己被放行了（AC#3 的正控②）：%s", self.Text)
	}
	if !strings.Contains(self.Text, "不许停自己") {
		t.Errorf("停自己的拒绝理由没说不许自停（要指名是「不许停自己」那一支，不是父归属那一支）：%q", self.Text)
	}

	// Both refusals must have touched nothing: each row still runs, each handle is
	// still attached, and the parent can still stop them for real. That last call
	// is also this leg's cleanup.
	for _, pair := range []struct{ id, label string }{{aID, "甲号孩子"}, {bID, "乙号孩子"}} {
		rec := x.h.mustLook(t, pair.id)
		if rec.Out.State != statemachine.StateThinking {
			t.Errorf("%s 被拒绝的调用改写了：%q", pair.label, string(rec.Out.State))
		}
	}
	stopA := x.cancel(t, parent197, aID)
	if stopA.IsError {
		t.Errorf("父任务停甲号孩子被拒（拒绝类正控不许把有权的那一支也一起堵死）：%s", stopA.Text)
	}
	stopB := x.cancel(t, parent197, bID)
	if stopB.IsError {
		t.Errorf("父任务停乙号孩子被拒：%s", stopB.Text)
	}
	for _, park := range []chan Result{parkA, parkB} {
		select {
		case <-park:
		case <-t.Context().Done():
			t.Error("被停的孩子没有让父任务那次 task.spawn 收口")
		}
	}
}

// Test221ParentCancellationStillDoesNotCascade is the other half of AC#3: the
// shape at subagent_197.go's context.WithoutCancel call must not be broken by
// this ticket, and the reply the model reads on that path must stay true after
// 乙形's deletion.
func Test221ParentCancellationStillDoesNotCascade(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	x := build221(t, windowGate221(), heldProvider(started, release), false)
	t.Cleanup(func() { close(release) })

	parentCtx, cancelParent := context.WithCancel(t.Context())
	done := make(chan Result, 1)
	go func() {
		res, err := x.h.spawnResult(parentCtx, label197, "长任务")
		if err != nil {
			res = Result{Text: "dispatch error: " + err.Error(), IsError: true}
		}
		done <- res
	}()
	select {
	case <-started:
	case <-t.Context().Done():
		t.Fatal("子代理没有开始")
	}
	childID := ""
	for _, r := range x.h.childRows() {
		if r.Out.Label == label197 {
			childID = r.TaskID
		}
	}
	if childID == "" {
		t.Fatalf("名册里没有孩子那一行：%+v", x.h.childRows())
	}

	cancelParent()
	var res Result
	select {
	case res = <-done:
	case <-t.Context().Done():
		t.Fatal("父任务取消后 task.spawn 没有收口")
	}
	if !res.IsError || !strings.Contains(res.Text, "没有被级联取消") {
		t.Fatalf("父取消没被如实说出口：%q", res.Text)
	}
	if strings.Contains(res.Text, "可以单独停它") {
		t.Errorf("那条回执还留着裸许诺 %q（乙形第三处删句）：%q", "可以单独停它", res.Text)
	}
	if !strings.Contains(res.Text, "还能用 task.cancel 单独停它") {
		t.Errorf("回执没写成当下为真的那半句（谁停得动它）：%q", res.Text)
	}
	rec := x.h.mustLook(t, childID)
	if rec.Out.State != statemachine.StateThinking {
		t.Errorf("父取消把孩子也停了：state = %q, want %q（级联会正好体现在这一维上）",
			string(rec.Out.State), string(statemachine.StateThinking))
	}
	// The row is still stoppable by its own parent - the promise this ticket made
	// true, now measured on the very path that used to only claim it.
	stop := x.cancel(t, parent197, childID)
	if stop.IsError {
		t.Errorf("父任务事后停它被拒：%s", stop.Text)
	}
}

// Test221TaskCancelUnderNoGateStopsAtTheWindow pins the consequence this leg
// found during its nail check instead of designing around it: with the D34 level
// (L1) and a gate that cannot run a window - tools.NoGate, i.e. 票 21 未接入,
// which several existing harnesses use - the call is refused BEFORE Execute.
// Anyone measuring task.cancel through a NoGate bridge is measuring the gate.
func Test221TaskCancelUnderNoGateStopsAtTheWindow(t *testing.T) {
	x := build221(t, NoGate{}, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, false)
	x.h.roster.PublishSubagent("kid-221", parent197, "只登记不起环的孩子")
	res := x.cancel(t, parent197, "kid-221")
	if !res.IsError {
		t.Fatalf("NoGate 下 task.cancel 竟然执行了（那意味着 L1 窗口被绕过）：%s", res.Text)
	}
	if !strings.Contains(res.Text, "L1 确认窗口") {
		t.Errorf("拒绝理由没说是窗口拦的：%q", res.Text)
	}
	if strings.Contains(res.Text, "只有父任务能停自己的孩子") {
		t.Errorf("这条拒绝被当成了权限判定，其实它根本没进到权限那一步：%q", res.Text)
	}
	if rec := x.h.mustLook(t, "kid-221"); rec.Out.State != statemachine.StateThinking {
		t.Errorf("被窗口拦下却动了那一行：%q", string(rec.Out.State))
	}
}
