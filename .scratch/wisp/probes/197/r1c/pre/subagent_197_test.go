package tools

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// Ticket 197 leg A (实体层): a subagent really exists - a row in the roster with
// its own task id, parent link, kind and D43 state; it can be derived, listed and
// stopped on its own; and its conclusion comes back to the parent carrying the
// existing C25 source name.
//
// Injection stays on the frozen seams (AGENTS §1.3): the C5 LlmProvider is the
// fake - which is what lets a leg say "this child is still running" without a wall
// clock - and everything else is the real bridge, the real registry, the real
// risk.Provenance and the real agent.Loop. No mock stands in for a real component.

const (
	parent197 = "parent-task-197"
	label197  = "读日志"
	// secret197 is long enough for C25 to take a matchable fragment out of it.
	secret197   = "外来内容探针 197-subagent-4f2a：子代理从外面读回来的一段长文本"
	foreign197A = "甲号子代理的结论甲甲甲 197-a"
	foreign197B = "乙号子代理的结论乙乙乙 197-b"
)

// fake197Provider is the C5 seam with three knobs: the answer text, an optional
// "hold every call here until the test releases it" barrier, and an optional
// first-turn tool call for task.spawn (the depth leg).
type fake197Provider struct {
	name   string
	answer string

	mu     sync.Mutex
	calls  int
	roster *TaskRoster
	// firstCallSawSubagentRow is the DSH measurement: was a roster row already
	// published when this child made its FIRST model call?
	firstCallSawSubagentRow bool

	// started receives one token per in-flight model call; release, when closed,
	// lets them all finish. This is how the 9th-spawn refusal and the cancel
	// ordering are made deterministic without sleeping.
	started chan struct{}
	release chan struct{}

	askSpawnFirstTurn bool
}

func (p *fake197Provider) Stream(ctx context.Context, _ *llm.Request, emit func(llm.StreamEvent) error) error {
	p.mu.Lock()
	p.calls++
	n := p.calls
	if p.roster != nil && n == 1 {
		p.firstCallSawSubagentRow = len(p.roster.Descendants(parent197)) > 0
	}
	askSpawn := p.askSpawnFirstTurn && n == 1
	p.mu.Unlock()

	if p.started != nil {
		p.started <- struct{}{}
		if p.release != nil {
			select {
			case <-p.release:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	if askSpawn {
		_ = emit(llm.StreamEvent{Type: llm.EvToolCallStart, ToolCallID: "call-197-spawn", ToolName: "task.spawn"})
		_ = emit(llm.StreamEvent{Type: llm.EvToolCallArgsDelta, ToolCallID: "call-197-spawn",
			ArgsDelta: `{"description":"孙代理","prompt":"再派一枚"}`})
		_ = emit(llm.StreamEvent{Type: llm.EvToolCallEnd, ToolCallID: "call-197-spawn", ToolName: "task.spawn"})
		_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopToolUse})
		_ = emit(llm.StreamEvent{Type: llm.EvDone})
		return nil
	}
	text := p.answer
	if text == "" {
		text = "子代理完成了任务"
	}
	_ = emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: text})
	_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopEndTurn})
	_ = emit(llm.StreamEvent{Type: llm.EvDone})
	return nil
}

func (p *fake197Provider) Info() llm.ProviderInfo {
	return llm.ProviderInfo{Model: "fake-197", Provider: p.name, MaxContextWindow: 8192}
}

func (p *fake197Provider) callCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.calls
}

// fake197Dir is the parent's C4 directory as the child is handed it: it lists
// task.spawn, so the filtering leg has something to filter out.
type fake197Dir struct {
	mu       sync.Mutex
	executed []string
}

func (d *fake197Dir) Tools(_ context.Context) ([]agent.ToolInfo, error) {
	return []agent.ToolInfo{
		{Name: "task.spawn", RiskLevel: "L0"},
		{Name: "task.output", RiskLevel: "L0"},
	}, nil
}

func (d *fake197Dir) Execute(_ context.Context, req agent.ToolRequest) (agent.ToolOutcome, error) {
	d.mu.Lock()
	d.executed = append(d.executed, req.Name)
	d.mu.Unlock()
	return agent.ToolOutcome{Text: "工具 " + req.Name + " 的假执行结果"}, nil
}

func (d *fake197Dir) calls() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]string(nil), d.executed...)
}

// sub197Harness wires the real bridge with the task family plus task.spawn, and
// keeps the roster, the stream feed, the taint engine and the admission spy.
type sub197Harness struct {
	roster   *TaskRoster
	bridge   *Bridge
	deps     *SubagentDeps
	prov     *risk.Provenance
	dir      *fake197Dir
	streamMu sync.Mutex
	streams  map[string][]string
	muAdmit  sync.Mutex
	admits   []string
}

func newSub197Harness(t *testing.T) *sub197Harness {
	t.Helper()
	h := &sub197Harness{
		roster:  NewTaskRoster(),
		prov:    risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil}),
		dir:     &fake197Dir{},
		streams: map[string][]string{},
	}
	h.roster.MarkRoot(parent197, "根任务")
	return h
}

func (h *sub197Harness) build(t *testing.T, childProvider func() llm.LlmProvider, admitWired bool) {
	t.Helper()
	h.buildWith(t, childProvider, admitWired, nil)
}

// buildWith is build with one mutation of the deps, which is how the
// "provenance engine not wired" fail-closed branch gets exercised against the
// real tool instead of a stand-in.
func (h *sub197Harness) buildWith(t *testing.T, childProvider func() llm.LlmProvider,
	admitWired bool, mutate func(*SubagentDeps)) {
	t.Helper()
	reg := NewRegistry()
	deps := SubagentDeps{
		Roster:      h.roster,
		Provenance:  h.prov,
		ParentTools: h.dir,
		Stream: func(key, text string) {
			h.streamMu.Lock()
			h.streams[key] = append(h.streams[key], text)
			h.streamMu.Unlock()
		},
		BaseOptions: func() (agent.Options, bool) {
			opt := agent.Options{Provider: childProvider(), Tools: h.dir}
			if admitWired {
				opt.AdmitTask = func(taskID string) func() {
					h.muAdmit.Lock()
					h.admits = append(h.admits, taskID)
					h.muAdmit.Unlock()
					return func() {}
				}
			}
			return opt, true
		},
	}
	if mutate != nil {
		mutate(&deps)
	}
	h.deps = &deps
	for _, e := range append(BuiltinTaskEntries(TaskDeps{Roster: h.roster}), BuiltinSubagentEntries(deps)...) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	h.bridge = New(Options{Registry: reg, Provenance: h.prov, Gate: NoGate{},
		Logf: func(string, ...any) {}})
}

// spawnResult dispatches one task.spawn call through the real bridge, so the
// parent task id reaches the tool the way production delivers it. No t.Fatalf
// inside, so a leg that expects a refusal can read the text.
func (h *sub197Harness) spawnResult(ctx context.Context, label, prompt string) (Result, error) {
	out, err := h.bridge.Execute(ctx, agent.ToolRequest{
		TaskID: parent197, CorrelationID: parent197, CallID: "call-197-spawn-" + label,
		Name: "task.spawn",
		Args: json.RawMessage(`{"description":"` + label + `","prompt":"` + prompt + `"}`),
	})
	if err != nil {
		return Result{}, err
	}
	return Result{Text: out.Text, IsError: out.IsError}, nil
}

func (h *sub197Harness) spawn(t *testing.T, ctx context.Context, label, prompt string) Result {
	t.Helper()
	res, err := h.spawnResult(ctx, label, prompt)
	if err != nil {
		t.Fatalf("task.spawn: %v", err)
	}
	return res
}

func (h *sub197Harness) childRows() []TaskRecord { return h.roster.Descendants(parent197) }

func (h *sub197Harness) childRow(t *testing.T) TaskRecord {
	t.Helper()
	rows := h.childRows()
	if len(rows) != 1 {
		t.Fatalf("名册里子代理行数 = %d, want 1（%+v）", len(rows), rows)
	}
	return rows[0]
}

func (h *sub197Harness) mustLook(t *testing.T, id string) TaskRecord {
	t.Helper()
	out, ok := h.roster.Look(id)
	if !ok {
		t.Fatalf("查不到任务 %s", id)
	}
	return TaskRecord{TaskID: id, Out: out}
}

func (h *sub197Harness) streamText(key string) string {
	h.streamMu.Lock()
	defer h.streamMu.Unlock()
	return strings.Join(h.streams[key], " | ")
}

// AC#1 - the row exists with three real identity values, and the conclusion is
// filed back onto the SAME row. Field existence alone would not do: ticket 181's
// AC#7 lesson is that a permanently empty field has to be judgable red, so every
// value below is compared against what the test actually passed in.
func Test197SpawnPublishesIdentityRow(t *testing.T) {
	h := newSub197Harness(t)
	prov := &fake197Provider{name: "child", answer: foreign197A, roster: h.roster}
	h.build(t, func() llm.LlmProvider { return prov }, true)

	res := h.spawn(t, t.Context(), label197, "把日志读完并总结")
	if res.IsError {
		t.Fatalf("派生失败：%s", res.Text)
	}
	row := h.childRow(t)
	if row.Out.Kind != TaskKindSubagent {
		t.Errorf("Kind = %q, want %q", row.Out.Kind, TaskKindSubagent)
	}
	if row.Out.ParentTaskID != parent197 {
		t.Errorf("ParentTaskID = %q, want %q", row.Out.ParentTaskID, parent197)
	}
	if row.Out.Label != label197 {
		t.Errorf("Label = %q, want %q（description 与 prompt 两枚字段不能混用）", row.Out.Label, label197)
	}
	if row.Out.Text != foreign197A {
		t.Errorf("结论没回到同一条记录上：Text = %q, want %q", row.Out.Text, foreign197A)
	}
	if st, why := row.Out.StateAnswer(); st == "" {
		t.Errorf("State 维不可读：%s（state=%q）", why, string(row.Out.State))
	}
	if st := row.Out.State; st != statemachine.StateSettling {
		t.Errorf("State = %q, want %q（跑完的孩子落在 D43 的 Settling 上）", string(st), string(statemachine.StateSettling))
	}
	if !statemachine.Valid(row.Out.State) {
		t.Errorf("State %q 不在 D43 的 20 枚名字上", string(row.Out.State))
	}
	// The stream key shape (§0): subagent:<taskID>, id passed through verbatim.
	key := SubagentStreamKey(row.TaskID)
	if key != "subagent:"+row.TaskID {
		t.Errorf("流键 = %q, want subagent:%s", key, row.TaskID)
	}
	if got := h.streamText(key); !strings.Contains(got, "已派生") || !strings.Contains(got, "已结束") {
		t.Errorf("这条子代理的流没有两段生命周期：%q", got)
	}
	if !strings.Contains(res.Text, foreign197A) {
		t.Errorf("结论没有回到父任务：%q", res.Text)
	}
	// The root row's identity is on file too, which is what makes the tree walkable.
	if root := h.mustLook(t, parent197); root.Out.Kind != TaskKindRoot {
		t.Errorf("根行 kind = %q, want %q", root.Out.Kind, TaskKindRoot)
	}
	// 跑完的孩子不再占池位。
	if got := h.roster.InFlightSubagents(); got != 0 {
		t.Errorf("跑完后池没清空：%d, want 0", got)
	}
}

// AC#1b (the DSH leg) - the row is published before the child's FIRST model call,
// not after it: an unlisted child is exactly the "实时列表会在调用方最需要它的
// 几十秒里保持不可见" failure the dispatch names.
func Test197RowExistsBeforeFirstChildModelCall(t *testing.T) {
	h := newSub197Harness(t)
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	prov := &fake197Provider{name: "child", answer: foreign197A, roster: h.roster,
		started: started, release: release}
	h.build(t, func() llm.LlmProvider { return prov }, true)

	done := make(chan Result, 1)
	go func() {
		res, err := h.spawnResult(context.Background(), label197, "跑一段慢任务")
		if err != nil {
			t.Errorf("spawn: %v", err)
		}
		done <- res
	}()
	<-started // the child is inside its first model call right now
	rows := h.childRows()
	if len(rows) != 1 {
		t.Fatalf("孩子正在跑，名册里却有 %d 行（跑起来才登记＝界面在说谎）", len(rows))
	}
	if st := rows[0].Out.State; st != statemachine.StateThinking {
		t.Errorf("在跑的行 State = %q, want %q", string(st), string(statemachine.StateThinking))
	}
	close(release)
	<-done
	if !prov.firstCallSawSubagentRow {
		t.Error("孩子的第一次模型调用发生时名册里还没有它（先起跑后登记）")
	}
}

// spawnDirect runs the real tool with the same per-call context the bridge builds
// (corr = the parent task id), bypassing the bridge's in-bridge tool-concurrency
// ceiling (D38d). That ceiling is lower than this ticket's 8-slot pool, so holding
// eight children at once through eight bridge calls would queue them against the
// C22 30s budget instead of measuring the cap. The ceiling itself is reported in
// the evidence file; every other leg here still goes through the real bridge.
func (h *sub197Harness) spawnDirect(ctx context.Context, label, prompt string) (Result, error) {
	tctx := withCancel(ctx, cancelHandle{corr: parent197})
	return subagentSpawn{d: *h.deps}.Execute(tctx,
		json.RawMessage(`{"description":"`+label+`","prompt":"`+prompt+`"}`), nil)
}

// AC#2 - the pool is 8, shared by the root and all its descendants, and the 9th
// is a hard refusal with a readable reason (never a silent queue). Reverse
// control: 8 all finish.
func Test197SubagentPoolCapsAtEight(t *testing.T) {
	h := newSub197Harness(t)
	started := make(chan struct{}, MaxConcurrentSubagents+1)
	release := make(chan struct{})
	h.build(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A,
			started: started, release: release}
	}, true)

	results := make(chan Result, MaxConcurrentSubagents)
	errs := make(chan error, MaxConcurrentSubagents)
	for i := 0; i < MaxConcurrentSubagents; i++ {
		go func(i int) {
			res, err := h.spawnDirect(context.Background(), string(rune('A'+i)), "并发任务")
			if err != nil {
				errs <- err
				return
			}
			results <- res
		}(i)
	}
	for i := 0; i < MaxConcurrentSubagents; i++ {
		<-started
	}
	if got := h.roster.InFlightSubagents(); got != MaxConcurrentSubagents {
		t.Fatalf("池内占位 = %d, want %d", got, MaxConcurrentSubagents)
	}

	ninth, err := h.spawnDirect(t.Context(), "第九枚", "再来一枚")
	if err != nil {
		t.Fatalf("第 9 枚调用: %v", err)
	}
	if !ninth.IsError {
		t.Fatalf("第 9 枚没有被硬拒：%s", ninth.Text)
	}
	for _, want := range []string{"8", "拒绝派生"} {
		if !strings.Contains(ninth.Text, want) {
			t.Errorf("超限理由不可读：%q 里没有 %q", ninth.Text, want)
		}
	}
	if got := h.roster.InFlightSubagents(); got != MaxConcurrentSubagents {
		t.Errorf("被拒的那一枚占了池位：in-flight = %d, want %d（硬拒不是排队）", got, MaxConcurrentSubagents)
	}
	if rows := len(h.childRows()); rows != MaxConcurrentSubagents {
		t.Errorf("被拒的那一枚在名册里留了行：%d, want %d", rows, MaxConcurrentSubagents)
	}

	close(release)
	for i := 0; i < MaxConcurrentSubagents; i++ {
		select {
		case err := <-errs:
			t.Fatalf("并发派生出错：%v", err)
		case r := <-results:
			if r.IsError {
				t.Errorf("8 枚以内的一枚失败了（反控）：%s", r.Text)
			}
		}
	}
	if rows := h.childRows(); len(rows) != MaxConcurrentSubagents {
		t.Errorf("名册行数 = %d, want %d", len(rows), MaxConcurrentSubagents)
	}
	if got := h.roster.InFlightSubagents(); got != 0 {
		t.Errorf("全跑完后可位没回收：%d, want 0", got)
	}
}

// AC#3 - depth 1, structurally: task.spawn is not in the child's directory at
// all, a child that asks for it by name is refused with the reason, and the
// child's own run derives no second-generation row.
func Test197ChildCannotDeriveSubagent(t *testing.T) {
	h := newSub197Harness(t)
	child := &fake197Provider{name: "child", answer: "我试了派一枚，失败了",
		askSpawnFirstTurn: true}
	h.build(t, func() llm.LlmProvider { return child }, true)

	// The directory the child is offered is the parent's minus task.spawn.
	childDir := newSubagentToolProvider(h.dir)
	list, err := childDir.Tools(t.Context())
	if err != nil {
		t.Fatalf("child Tools: %v", err)
	}
	for _, info := range list {
		if info.Name == "task.spawn" {
			t.Errorf("孩子的工具目录里还有 %s（深度 1 要靠结构保证，不靠运行期计数）", info.Name)
		}
	}
	if len(list) != 1 {
		t.Errorf("孩子的目录少了不该少的工具：%+v", list)
	}
	// And a child that asks for it by name anyway gets a readable refusal.
	out, err := childDir.Execute(t.Context(), agent.ToolRequest{
		TaskID: "child-197", Name: "task.spawn",
		Args: json.RawMessage(`{"description":"孙代理","prompt":"再派一枚"}`),
	})
	if err != nil {
		t.Fatalf("child Execute: %v", err)
	}
	if !out.IsError || !strings.Contains(out.Text, "深度") {
		t.Errorf("孩子派孙没有被拒并说明原因：%+v", out)
	}
	for _, name := range h.dir.calls() {
		if name == "task.spawn" {
			t.Errorf("被拒的那一枚还是打到了父目录上：%v", h.dir.calls())
		}
	}

	// The full-loop leg: the child really asks for task.spawn, and no grandchild
	// row ever appears.
	if res := h.spawn(t, t.Context(), label197, "让孩子试着再派一枚"); res.IsError {
		t.Fatalf("派生失败：%s", res.Text)
	}
	rows := h.childRows()
	if len(rows) != 1 {
		t.Fatalf("名册里出现第二代（深度上限 1 被破）：%+v", rows)
	}
	if grand := h.roster.Descendants(rows[0].TaskID); len(grand) != 0 {
		t.Errorf("孙代理行数 = %d, want 0：%+v", len(grand), grand)
	}
	if child.callCount() < 2 {
		t.Errorf("孩子那发 task.spawn 请求没走满两轮（call=%d），这一发正控没跑到", child.callCount())
	}
}

// AC#4 - the conclusion carries the EXISTING C25 source name, and nothing outside
// the register. Deleting the stamp step has to judge red: that is what the ScopeTaints
// and CheckText legs below measure.
func Test197ConclusionCarriesTaskOutputTaint(t *testing.T) {
	h := newSub197Harness(t)
	h.build(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: secret197 + " 后续"}
	}, true)

	res := h.spawn(t, t.Context(), label197, "读一份外来的东西")
	if res.IsError {
		t.Fatalf("派生失败：%s", res.Text)
	}
	marks := h.prov.ScopeTaints(parent197)
	var found bool
	for _, m := range marks {
		if m.Tool == risk.SrcTaskOutput {
			found = true
		}
		if strings.HasPrefix(m.Tool, "subagent.") {
			t.Errorf("名册外的源名被发明了：%q（§0 禁用）", m.Tool)
		}
	}
	if !found {
		t.Fatalf("父任务作用域里没有 %q 源名的戳（子代理结论走了没人认识的通道）：%+v",
			risk.SrcTaskOutput, marks)
	}
	hit, ok := h.prov.CheckText(parent197, risk.ChClipboard, secret197)
	if !ok || hit.Fragment == "" {
		t.Errorf("子代理结论没有被 R4 认出来（戳没盖实）：%+v %v", hit, ok)
	}
	// 盖戳失败的分支也要说出口，不许静默。
	h2 := newSub197Harness(t)
	h2.buildWith(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: secret197 + " 第二段"}
	}, true, func(d *SubagentDeps) { d.Provenance = nil })
	if r := h2.spawn(t, t.Context(), "没有溯源", "读一段"); !strings.Contains(r.Text, "没接线") {
		t.Errorf("溯源引擎没接线时回复里没有那句注意：%q", r.Text)
	}
}

// AC#5 - cancellation: stopping the parent does NOT cascade; the child stays
// findable and stoppable; and stopping one row touches only that row.
func Test197CancelIsPerRowAndNeverCascades(t *testing.T) {
	h := newSub197Harness(t)
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	h.build(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197B, started: started, release: release}
	}, true)

	// A sibling row that nobody's cancel may touch.
	h.roster.PublishSubagent("sibling-197", parent197, "隔壁那枚")
	sibWatch := h.roster.WatchRow("sibling-197")

	parentCtx, cancelParent := context.WithCancel(t.Context())
	done := make(chan Result, 1)
	go func() {
		res, err := h.spawnResult(parentCtx, label197, "长任务")
		if err != nil {
			t.Errorf("spawn: %v", err)
		}
		done <- res
	}()
	<-started

	cancelParent()
	res := <-done
	if !res.IsError || !strings.Contains(res.Text, "没有被级联取消") {
		t.Fatalf("父取消没被如实说出口：%q", res.Text)
	}
	childID := ""
	for _, r := range h.childRows() {
		if r.Out.Label == label197 {
			childID = r.TaskID
		}
	}
	if childID == "" {
		t.Fatalf("父任务已经不听了，孩子却从名册里消失了：%+v", h.childRows())
	}
	watch := h.roster.WatchRow(childID)
	if st := h.mustLook(t, childID).Out.State; st != statemachine.StateThinking {
		t.Errorf("父取消后孩子的 State = %q, want %q（级联会正好体现在这一维上）", string(st), string(statemachine.StateThinking))
	}

	ok, why := h.roster.Cancel(childID)
	if !ok {
		t.Fatalf("停这一枚子代理没成功：%s", why)
	}
	// The child's own loop observes the cancel; its row is finalised by the finish
	// watcher, and "only that row" means exactly one row moved.
	final := waitUntilSettled(t, watch)
	if final.State == statemachine.StateThinking {
		t.Errorf("停掉的孩子还在跑：%+v", final)
	}
	if st, why := final.StateAnswer(); st == "" {
		t.Errorf("停掉的孩子那一维不可读：%s", why)
	}
	if got := h.roster.InFlightSubagents(); got != 0 {
		t.Errorf("孩子结束后可位没回收：%d, want 0", got)
	}
	select {
	case moved := <-sibWatch:
		t.Errorf("停子代理把隔壁也写了：%+v", moved)
	default:
	}
	if sib := h.mustLook(t, "sibling-197"); sib.Out.State != statemachine.StateThinking {
		t.Errorf("隔壁那枚被动了：State = %q", string(sib.Out.State))
	}
	// 停完再停是「没有在跑」，不是假装停成功。
	if ok2, why2 := h.roster.Cancel(childID); ok2 || why2 == "" {
		t.Errorf("重复取消报了成功：%v %q", ok2, why2)
	}
	if ok3, why3 := h.roster.Cancel("nope-197"); ok3 || !strings.Contains(why3, "查不到") {
		t.Errorf("停一枚不存在的孩子没报「查不到」：%v %q", ok3, why3)
	}
	close(release)
	// 取消语义要写进给模型读的那段话（§0 / DSH tool-subagent-control 的理由）。
	if desc := (subagentSpawn{}).Description(); !strings.Contains(desc, "停掉父任务不会级联") {
		t.Errorf("给模型读的说明文本里没有「父取消不级联」这句：%q", desc)
	}
}

// waitUntilSettled joins the roster's own write stream: it returns the first
// state that is no longer "running", and fails the test if the row never moves
// before the test's own context ends. No sleep, no wall-clock comparison.
func waitUntilSettled(t *testing.T, watch <-chan TaskOutput) TaskOutput {
	t.Helper()
	for {
		select {
		case o, ok := <-watch:
			if !ok {
				t.Fatal("行监视通道被关掉")
			}
			if o.State != statemachine.StateThinking {
				return o
			}
		case <-t.Context().Done():
			t.Fatal("等名册那一行收尾没等到")
		}
	}
}

// subagentDepsFieldNames lists the seams the spawn path can be handed at all. The
// count is the guard: a subagent must not acquire a sixth knob that turns out to
// be a private approval outlet (ticket 197 §0: 子代理永不自批, and AGENTS §1.2
// bans any L2 allow whose source is the panel side).
func subagentDepsFieldNames() []string {
	typ := reflect.TypeOf(SubagentDeps{})
	out := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		out = append(out, typ.Field(i).Name)
	}
	return out
}

// AC#6 - a subagent never approves itself: with no host admission hook the
// spawner refuses to derive at all, and when one is wired the child asks the
// parent's queue through that SAME hook (one queue, not a private one).
func Test197SubagentHasNoSelfApprovalOutlet(t *testing.T) {
	h := newSub197Harness(t)
	h.build(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, false)

	res := h.spawn(t, t.Context(), label197, "没有审批层就想派一枚")
	if !res.IsError || !strings.Contains(res.Text, "永不自批") {
		t.Fatalf("宿主没有准入钩子时派生成功了（子代理就有了自己的出口）：%q", res.Text)
	}
	if got := len(h.childRows()); got != 0 {
		t.Errorf("被拒的派生还是在名册里留了行：%d", got)
	}
	if got := h.roster.InFlightSubagents(); got != 0 {
		t.Errorf("被拒的派生占了池位：%d", got)
	}

	h2 := newSub197Harness(t)
	h2.build(t, func() llm.LlmProvider {
		return &fake197Provider{name: "child", answer: foreign197A}
	}, true)
	if res := h2.spawn(t, t.Context(), label197, "走父队列"); res.IsError {
		t.Fatalf("派生失败：%s", res.Text)
	}
	childID := h2.childRow(t).TaskID
	h2.muAdmit.Lock()
	joined := strings.Join(h2.admits, ",")
	h2.muAdmit.Unlock()
	if !strings.Contains(joined, childID) {
		t.Errorf("孩子的请求没有走宿主那一枚准入钩子（＝父队列）：%q", joined)
	}
	// 派生路径上没有第二枚队列/第二份 gate 可接：SubagentDeps 的字段只有装配、
	// 目录、溯源、流四枚，没有一处能塞进"允许出口"。
	if got := subagentDepsFieldNames(); len(got) != 5 {
		t.Errorf("SubagentDeps 字段集合变了：%v（多出来的每一枚都要问一句是不是新的允许出口）", got)
	}
}

// The stream key shape ticket 197 §0 froze, pinned at the source that writes it.
func Test197StreamKeyShapeIsLiteral(t *testing.T) {
	for _, id := range []string{"aB-1", "00000000-0000-4000-8000-000000000000", "带中文的id"} {
		if got := SubagentStreamKey(id); got != "subagent:"+id {
			t.Errorf("SubagentStreamKey(%q) = %q, want subagent:%q（id 逐字透传、前缀全小写、冒号分隔）", id, got, id)
		}
	}
	if SubagentStreamKeyPrefix != "subagent:" {
		t.Errorf("前缀 = %q, want subagent:", SubagentStreamKeyPrefix)
	}
	if MaxConcurrentSubagents != 8 || MaxSubagentDepth != 1 {
		t.Errorf("池/深度常量漂了：%d / %d, want 8 / 1", MaxConcurrentSubagents, MaxSubagentDepth)
	}
}
