package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 222 (the 丙 branch): "waiting for my child" is not "executing a tool",
// so it must not hold one of the bridge's D38d slots.
//
// What this file is the first ruler of: the child's tool surface running on the
// REAL bridge. subagent_197_test.go's harness drives the parent through
// h.bridge.Execute and hands the CHILD h.dir (a fake197Dir), so the one place
// where a waiting parent and a working child ever meet the same four slots was
// never executed by any test in this repository - ledger A420's ring 7, and the
// reason the shape stayed invisible through a green suite. Every leg below wires
// SubagentDeps.ParentTools to the bridge itself, which is what cmd/wisp/run.go:581
// does in production, so the mutual waiting is measured instead of inferred.
//
// No leg here judges a hang. Each one decides from a reading:
//   - Test222WaitingParentHoldsNoBridgeSlot samples len(bridge.sem) at the
//     instant every parent has announced "I am about to wait". The deterministic
//     marker is the task.spawn onUpdate delta, which the tool emits after the
//     point where it gives the slot up: 0 with the fix, 4 without it - and both
//     readings arrive in microseconds, no deadline involved.
//   - Test222SpawnConclusionArrivesThroughRealBridgeChildren reads which TEXT the
//     parent gets back: its own child's conclusion, or the timeout wording the
//     bridge writes when the parent's own per-tool budget ran out. Removing the
//     fix makes that a content failure.
//   - Test222CeilingStillCapsExecutedCallsWhileParentsWait is the reverse control
//     AC#2 demands: concurrent EXECUTIONS through the bridge stay at the frozen
//     four (PLAN.md D38d's reason - one turn must not swamp the machine) while the
//     parents wait. Each probe additionally records how many parents had already
//     returned at the moment it started, which is the causal half: before the fix
//     nothing can run until a parent's budget freed its slot, so the number it
//     writes down is nonzero.
//
// Ticket 235 AC#2 added one more ruler to all three legs, and it is the one that
// reads which side moved first: a child's arrival on the bridge must precede any
// parent's return (h.awaitChildOnBridge222). Causally it always does - a parent's
// conclusion comes from its child, and the child only reaches its conclusion
// after its tool call has come back through the bridge - so the reading costs
// microseconds on the fixed shape. Hand the child a tool surface that bypasses
// the bridge (the M3 shape, the child back on a fake197Dir) and nothing ever
// arrives while the parents still collect their conclusions one by one, which is
// exactly what that leg now reads red at 0.00s instead of leaning on the 30s
// guard rail below it.
//
// The per-tool budget these legs build the bridge with is a test-local knob
// (Options.DefaultTimeout; production reads cfg.Agent.PerToolTimeoutMS at
// cmd/wisp/run.go:562). It exists so the pre-fix reading arrives about three
// seconds after the spawn instead of half a minute later, with enough slack that
// a loaded machine cannot turn a post-fix run into a false red. No assertion here
// is "we waited and nothing came".

const (
	parent222   = "parent-task-222"
	probe222Tag = "probe.work222"

	// h222CeilingLiteral is AC#2's frozen number written as a literal, on purpose
	// and for the same reason subagent_197.go writes MaxConcurrentSubagents as 4:
	// two independent readings that can disagree are a measurement, one constant
	// compared against itself is not. TestToolConcurrencyCeilingIsFour compares
	// what it observes against MaxToolConcurrency and therefore cannot see that
	// constant being raised (A420 names that); these legs can.
	h222CeilingLiteral = 4

	// h222PreFixBudget is the C22 budget these legs hand the bridge: big enough
	// that real work here (a fake provider plus an in-memory probe) cannot exhaust
	// it by accident, small enough that the pre-fix reading is seconds.
	h222PreFixBudget = 3 * time.Second

	// h222SafetyBound is only a guard rail: if it ever fires, something in here is
	// parked forever (the first candidate being a slot handed back twice, which
	// blocks forever on an empty semaphore - the last check of
	// Test222CeilingStillCapsExecutedCallsWhileParentsWait is the other ruler for
	// that shape). It is never the assertion that decides.
	h222SafetyBound = 30 * time.Second
)

// h222 is the ticket 222 harness: ONE real bridge carrying both sides of the
// chain - the parent's task.spawn call and the child's tool directory.
type h222 struct {
	roster *TaskRoster
	prov   *risk.Provenance
	bridge *Bridge

	// parents receives every parent's Result. The probes read its LENGTH while
	// they run, which is the causal (not time-based) half of AC#2: with the fix a
	// child's call executes while all parents are still waiting, so the number
	// recorded is 0; without the fix nothing executes until a parent's budget
	// freed its slot, so the number recorded is at least 1.
	parents chan Result

	// parked receives the task.spawn "已派生" delta, which the tool emits AFTER the
	// point where it gives the bridge slot back. Holding N of them therefore means
	// N parents stand at the wait with their slots already returned - no polling,
	// no wall clock.
	parked chan string

	probe *probe222

	started chan struct{}
	release chan struct{}
}

// probe222 is the tool the children call - and the only way to see the bridge's
// execution window from the inside. It counts the calls running at the same
// moment, samples len(bridge.sem) and the number of parents that had already
// returned, and parks on gate222 when a leg set one so that several of them are
// provably in flight at once.
type probe222 struct {
	h    *h222
	name string

	runs atomic.Int32

	// arrived is ticket 235 AC#2's arrival signal: one non-blocking token per
	// execution that really reached the bridge. It carries no verdict by itself -
	// h.awaitChildOnBridge222 below is what reads it - and its buffer (16) is
	// bigger than the most this file ever sends (the reverse-control leg's 6 host
	// calls plus MaxConcurrentSubagents children = 10), so the send can never
	// become the thing a leg waits on.
	arrived chan struct{}

	held    atomic.Int32
	maxSeen atomic.Int32

	in      chan probeRun222
	gate222 chan struct{}
}

type probeRun222 struct {
	mark        string
	slotsHeld   int
	returnedYet int
}

func (p *probe222) Name() string { return p.name }
func (p *probe222) Description() string {
	return "票 222 的探针：孩子在真桥上执行的那一枚工具"
}

func (p *probe222) Parameters() JSONSchema {
	return JSONSchema(`{"type":"object","properties":{"mark":{"type":"string"}},"additionalProperties":true}`)
}

func (p *probe222) Execute(ctx context.Context, params json.RawMessage, _ func(string)) (Result, error) {
	var a struct {
		Mark string `json:"mark"`
	}
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "探针参数解析失败：" + err.Error(), IsError: true}, nil
	}
	cur := p.held.Add(1)
	for {
		prev := p.maxSeen.Load()
		if cur <= prev || p.maxSeen.CompareAndSwap(prev, cur) {
			break
		}
	}
	p.runs.Add(1)
	select {
	case p.arrived <- struct{}{}:
	default:
	}
	p.in <- probeRun222{
		mark:        a.Mark,
		slotsHeld:   len(p.h.bridge.sem),
		returnedYet: len(p.h.parents),
	}
	if p.gate222 != nil {
		select {
		case <-p.gate222:
		case <-ctx.Done():
		}
	}
	p.held.Add(-1)
	return Result{Text: "探针 " + a.Mark + " 在桥上跑完了"}, nil
}

// child222 is the C5 seam (a frozen injection surface, AGENTS §1.3) for a child
// that does WORK: turn 1 asks for a tool through the real bridge, turn 2 answers.
// The two barriers make the interleaving deterministic without a wall clock - the
// child announces itself inside its first model call and parks there until the leg
// has every parent in place, so no child can slip onto the bridge early and let a
// pre-fix run pass by luck.
type child222 struct {
	id     string
	answer string

	mu   sync.Mutex
	turn int

	started chan<- struct{}
	release <-chan struct{}
}

func (c *child222) Stream(ctx context.Context, _ *llm.Request, emit func(llm.StreamEvent) error) error {
	c.mu.Lock()
	c.turn++
	n := c.turn
	c.mu.Unlock()

	if n == 1 {
		if c.started != nil {
			select {
			case c.started <- struct{}{}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		if c.release != nil {
			select {
			case <-c.release:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		call := "call-222-" + c.id
		_ = emit(llm.StreamEvent{Type: llm.EvToolCallStart, ToolCallID: call, ToolName: probe222Tag})
		_ = emit(llm.StreamEvent{
			Type: llm.EvToolCallArgsDelta, ToolCallID: call,
			ArgsDelta: `{"mark":"` + c.id + `"}`,
		})
		_ = emit(llm.StreamEvent{Type: llm.EvToolCallEnd, ToolCallID: call, ToolName: probe222Tag})
		_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopToolUse})
		_ = emit(llm.StreamEvent{Type: llm.EvDone})
		return nil
	}
	_ = emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: c.answer})
	_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopEndTurn})
	_ = emit(llm.StreamEvent{Type: llm.EvDone})
	return nil
}

func (c *child222) Info() llm.ProviderInfo {
	return llm.ProviderInfo{Model: "fake-222", Provider: "child222", MaxContextWindow: 8192}
}

func newH222(t *testing.T) *h222 {
	t.Helper()
	h := &h222{
		roster:  NewTaskRoster(),
		prov:    risk.NewProvenance(risk.ProvOptions{NoProbe: true, SyncRoots: nil}),
		parents: make(chan Result, 2*MaxConcurrentSubagents),
		parked:  make(chan string, 2*MaxConcurrentSubagents),
		started: make(chan struct{}, 2*MaxConcurrentSubagents),
		release: make(chan struct{}),
	}
	h.roster.MarkRoot(parent222, "根任务 222")
	h.probe = &probe222{
		h:       h,
		name:    probe222Tag,
		in:      make(chan probeRun222, 16),
		arrived: make(chan struct{}, 16),
	}
	h.bridge = New(Options{
		Registry:       NewRegistry(),
		Provenance:     h.prov,
		Gate:           NoGate{},
		DefaultTimeout: h222PreFixBudget,
		Logf:           func(string, ...any) {},
		OnUpdate: func(_, _, delta string) {
			if strings.HasPrefix(delta, "task.spawn 已派生子代理") {
				h.parked <- delta
			}
		},
	})
	// Resident: true is D15(1) and it is load-bearing here, not decoration: a
	// non-resident entry is selected per turn by the loop's own retrieval, so a
	// child asking for this tool by name could be refused before the call ever
	// reached the bridge - which would make these legs measure the wrong thing.
	if err := h.bridge.Registry().Register(Entry{
		Tool: h.probe,
		Decl: Decl{Declared: risk.L0, Provider: KindBuiltin, Resident: true},
	}); err != nil {
		t.Fatalf("Register(%s): %v", probe222Tag, err)
	}
	return h
}

// wireSpawn registers task.spawn on the SAME bridge the children are then handed
// as their tool surface. That one line is the whole difference from
// subagent_197_test.go's harness, whose children get h.dir.
func (h *h222) wireSpawn(t *testing.T) {
	t.Helper()
	seq := &atomic.Int32{}
	deps := SubagentDeps{
		Roster:      h.roster,
		Provenance:  h.prov,
		ParentTools: h.bridge,
		BaseOptions: func() (agent.Options, bool) {
			n := int(seq.Add(1))
			return agent.Options{
				Provider: &child222{
					id:      fmt.Sprintf("c%d", n),
					answer:  fmt.Sprintf("子代理 %d 的结论正文 222", n),
					started: h.started,
					release: h.release,
				},
				Tools:     h.bridge,
				AdmitTask: func(string) func() { return func() {} },
			}, true
		},
	}
	for _, e := range BuiltinSubagentEntries(deps) {
		if err := h.bridge.Registry().Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
}

// launchParents issues n task.spawn calls through the bridge, one per goroutine.
// Results land in h.parents, whose LENGTH the probes read, so nothing here can
// block on an unread channel.
func (h *h222) launchParents(t *testing.T, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		go func(i int) { // test-only goroutine: d22scan scopes ban 1 to non-test files
			out, err := h.bridge.Execute(t.Context(), agent.ToolRequest{
				TaskID: parent222, CorrelationID: parent222, CallID: fmt.Sprintf("call-222-parent-%d", i),
				Name: "task.spawn",
				Args: json.RawMessage(fmt.Sprintf(`{"description":"甲%d","prompt":"让孩子在桥上干一件事"}`, i)),
			})
			if err != nil {
				h.parents <- Result{Text: err.Error(), IsError: true}
				return
			}
			h.parents <- Result{Text: out.Text, IsError: out.IsError}
		}(i)
	}
}

// await222Tokens collects n values from ch. The guard rail is the only way it
// comes back early, and it says what a fired rail would mean.
func await222Tokens[T any](t *testing.T, name string, ch <-chan T, n int) []T {
	t.Helper()
	out := make([]T, 0, n)
	timer := time.NewTimer(h222SafetyBound)
	defer timer.Stop()
	for len(out) < n {
		select {
		case v := <-ch:
			out = append(out, v)
		case <-timer.C:
			t.Fatalf("%s：只等到 %d/%d 枚，护栏到点（这一发真挂住了，不是读数不同）", name, len(out), n)
		}
	}
	return out
}

// parkParents puts n parents at the wait and proves the children are alive but
// held: every child is inside its first model call and every parent has emitted
// its spawn delta, which the tool emits AFTER it gave the bridge slot back.
func (h *h222) parkParents(t *testing.T, n int) {
	t.Helper()
	h.launchParents(t, n)
	await222Tokens(t, "孩子进入第一枚模型调用", h.started, n)
	if got := len(h.parked); got != n {
		t.Fatalf("只有 %d/%d 枚父任务到达等待点（父任务没派生成功，等待读数就无从谈起）", got, n)
	}
	if got := len(h.parents); got != 0 {
		t.Fatalf("已经有 %d 枚父任务返回了，这一发的等待读数就是假的", got)
	}
}

// awaitChildOnBridge222 is ticket 235 AC#2's pre-read: it reads WHICH side moved
// first before the leg goes on to collect n records. On the fixed shape a child's
// arrival on the bridge always precedes any parent's return - the parent's
// conclusion comes from its child, and the child only reaches that conclusion
// after its tool call has come back through the bridge - so this returns in
// microseconds and asserts nothing but the ordering. Hand the child a tool
// surface that bypasses the bridge (the M3 shape) and nothing arrives while the
// parents keep collecting conclusions, which is the reading this now fails on
// instead of the 30s guard rail of the await behind it.
//
// The parent result it may pick up is handed straight back, so the leg's own
// await over h.parents still sees all n of them. h.parents is buffered at
// 2*MaxConcurrentSubagents and never holds more than n at once.
func (h *h222) awaitChildOnBridge222(t *testing.T) {
	t.Helper()
	timer := time.NewTimer(h222SafetyBound)
	defer timer.Stop()
	select {
	case <-h.probe.arrived:
	case res := <-h.parents:
		h.parents <- res
		t.Fatalf("桥上的到达信号一枚都没有，而第一枚父任务已经拿到了结论（真桥上累计执行 %d 次，"+
			"反控那一发里这些全是宿主调用，孩子一次都没上桥）：%q ——"+
			"孩子的工具面绕开了真桥＝A420 的第 7 环，票 235 AC#2 这一发就是为它加的",
			h.probe.runs.Load(), res.Text)
	case <-timer.C:
		t.Fatalf("前置读数两头都没到（桥上到达 0 枚，父任务已返回 %d/%d 枚），护栏到点："+
			"这一发真挂住了，不是读数不同", len(h.parents), MaxConcurrentSubagents)
	}
}

// drainBridgeArrivals empties the arrival signal, so that every token
// awaitChildOnBridge222 reads afterwards can only have come from a child. The
// reverse-control leg calls it after wg.Wait(): a host call sends its arrival
// before it returns, so by then all of them are already buffered here.
func (h *h222) drainBridgeArrivals() {
	for {
		select {
		case <-h.probe.arrived:
		default:
			return
		}
	}
}

// AC#1 - the production shape gets its first ruler: n parents wait, each child
// runs a real tool call ON THE SAME BRIDGE, and every parent must get its own
// child's conclusion back rather than the timeout wording the bridge writes when
// the parent gave up waiting.
//
// Removing the fix turns this red on CONTENT (IsError plus the wrong text) and on
// the probe's own record of how many parents had already returned when it was
// finally allowed to run - not on a hang.
func Test222SpawnConclusionArrivesThroughRealBridgeChildren(t *testing.T) {
	h := newH222(t)
	h.wireSpawn(t)
	h.probe.gate222 = make(chan struct{})

	n := MaxConcurrentSubagents
	h.parkParents(t, n)
	close(h.release)
	// Ticket 235 AC#2: the pre-read in front of the await below. This leg was the
	// only one of the three that ever noticed a child whose tool surface bypasses
	// the bridge, and it noticed only after the 30s guard rail fired; the arrival
	// ordering decides that in microseconds, so the rail is a rail again and not
	// the assertion.
	h.awaitChildOnBridge222(t)

	// The gate keeps every child inside its bridge call until all n arrived, so
	// "the child ran while its parent was still waiting" is a fact this leg
	// establishes rather than a race it hopes to win.
	runs := await222Tokens(t, "孩子在桥上跑工具调用", h.probe.in, n)
	if got := h.probe.runs.Load(); got != int32(n) {
		t.Errorf("孩子在真桥上的工具调用跑了 %d 次, want %d（孩子的工具面没接到桥上＝A420 的第 7 环）", got, n)
	}
	for _, r := range runs {
		if r.returnedYet != 0 {
			t.Errorf("孩子 %s 的工具调用起跑时已有 %d 枚父任务返回，want 0：它是等父任务的 per-tool 预算到点才挤上桥的",
				r.mark, r.returnedYet)
		}
	}
	// The n-th token proves all n children were inside the bridge at the same
	// instant, i.e. the waiting parents had really given the slots back.
	if last := runs[n-1]; last.slotsHeld != h222CeilingLiteral {
		t.Errorf("第 %d 枚孩子的工具调用起跑时桥位占用 %d 枚, want %d（池内 %d 枚孩子该同时跑在 %d 枚许可上）",
			n, last.slotsHeld, h222CeilingLiteral, n, h222CeilingLiteral)
	}

	close(h.probe.gate222)
	results := await222Tokens(t, "父任务拿到结论", h.parents, n)
	for i, res := range results {
		if res.IsError {
			t.Errorf("父任务 %d 收到的是错误而不是结论：%q", i, res.Text)
			continue
		}
		if strings.Contains(res.Text, "不等了") {
			t.Errorf("父任务 %d 又落回票 222 那条文本：%q", i, res.Text)
		}
	}
	// Every child's conclusion is back in exactly one parent's text: the answers
	// are numbered, so a crossed wire or a lost child cannot pass for free.
	for k := 1; k <= n; k++ {
		want := fmt.Sprintf("子代理 %d 的结论正文 222", k)
		got := 0
		for _, res := range results {
			if strings.Contains(res.Text, want) {
				got++
			}
		}
		if got != 1 {
			t.Errorf("结论 %q 出现在 %d 枚父任务回复里, want 1（%v）", want, got, results)
		}
	}
	if got := h.roster.InFlightSubagents(); got != 0 {
		t.Errorf("全部跑完池位没回收：%d, want 0", got)
	}
	if got := h.probe.maxSeen.Load(); got > h222CeilingLiteral {
		t.Errorf("同时在桥上跑的工具调用数 %d 超过了冻结的 %d（不占许可地等，不等于把天花板说破）",
			got, h222CeilingLiteral)
	}
	if got := len(h.bridge.sem); got != 0 {
		t.Errorf("全部结束后桥位没回到空：%d, want 0", got)
	}
}

// AC#2 - the shape that caused the whole ticket gets pinned: while every parent
// is waiting for its child, the bridge holds ZERO slots on their behalf.
//
// The reading is len(bridge.sem), sampled the instant the leg holds all n
// "已派生" deltas. With the fix the number is 0, without it the number is n, and
// neither answer depends on a deadline: this leg decides in microseconds.
func Test222WaitingParentHoldsNoBridgeSlot(t *testing.T) {
	h := newH222(t)
	h.wireSpawn(t)

	n := MaxConcurrentSubagents
	h.parkParents(t, n)

	held := len(h.bridge.sem)
	if held != 0 {
		t.Fatalf("等待孩子的父任务占着 %d 枚桥位，want 0：全桥 %d 枚，父任务一边干等一边占位＝票 222 现量链的形状，"+
			"孩子只能排在 sem 后面，直到父任务自己的 per-tool 超时松绑", held, cap(h.bridge.sem))
	}
	if got := cap(h.bridge.sem); got != h222CeilingLiteral {
		t.Errorf("桥的执行许可总数 = %d, want %d（D38d 那枚数字不许被这枚票改动）", got, h222CeilingLiteral)
	}
	if got := h.roster.InFlightSubagents(); got != n {
		t.Errorf("池内占位 = %d, want %d（不占桥的许可不等于把孩子从池里摘掉）", got, n)
	}
	if got := len(h.roster.RunningSubagentIDs()); got != n {
		t.Errorf("名册里在跑的行数 = %d, want %d", got, n)
	}

	close(h.release)
	// Ticket 235 AC#2: this leg used to be blind to a child whose tool surface
	// bypasses the bridge - the bypassing child still hands its parent a
	// conclusion, so every assertion above reads the same either way. The arrival
	// ordering is what tells the two shapes apart, and it reads in microseconds.
	h.awaitChildOnBridge222(t)
	results := await222Tokens(t, "父任务拿到结论", h.parents, n)
	for i, res := range results {
		if res.IsError {
			t.Errorf("放开之后父任务 %d 还是收到了错误：%q", i, res.Text)
		}
	}
}

// AC#2's reverse control - the frozen reason behind the ceiling still holds with
// waiting parents in the picture: the bridge never runs more than
// h222CeilingLiteral tool calls at once, and the extra host-side calls are proven
// to have run while EVERY parent was still waiting.
func Test222CeilingStillCapsExecutedCallsWhileParentsWait(t *testing.T) {
	h := newH222(t)
	h.wireSpawn(t)
	h.probe.gate222 = make(chan struct{})

	n := MaxConcurrentSubagents
	h.parkParents(t, n)
	if held := len(h.bridge.sem); held != 0 {
		t.Fatalf("等待孩子的父任务占着 %d 枚桥位，want 0（反控的前半：先把许可还不回来）", held)
	}

	oversubscribed := h222CeilingLiteral + 2
	var wg sync.WaitGroup
	for i := 0; i < oversubscribed; i++ {
		wg.Add(1)
		go func(i int) { // test-only goroutine: d22scan scopes ban 1 to non-test files
			defer wg.Done()
			if _, err := h.bridge.Execute(t.Context(), agent.ToolRequest{
				TaskID: parent222, CorrelationID: parent222, CallID: fmt.Sprintf("call-222-host-%d", i),
				Name: probe222Tag,
				Args: json.RawMessage(fmt.Sprintf(`{"mark":"host%d"}`, i)),
			}); err != nil {
				t.Errorf("宿主探针 %d: %v", i, err)
			}
		}(i)
	}

	running := await222Tokens(t, "宿主侧探针起跑", h.probe.in, h222CeilingLiteral)
	if got := h.probe.maxSeen.Load(); got != int32(h222CeilingLiteral) {
		t.Errorf("同时执行的工具调用数 = %d, want %d（等待中的父任务不该算进来，多出来的宿主调用要挤得进来）",
			got, h222CeilingLiteral)
	}
	for _, r := range running {
		if r.returnedYet != 0 {
			t.Errorf("探针 %s 起跑时已有 %d 枚父任务返回，want 0：修前正是这样，谁都只能在父任务预算到点之后才上桥",
				r.mark, r.returnedYet)
		}
		if r.slotsHeld > h222CeilingLiteral {
			t.Errorf("探针 %s 起跑时桥位占用 %d 枚，超过冻结的 %d", r.mark, r.slotsHeld, h222CeilingLiteral)
		}
	}

	close(h.probe.gate222)
	wg.Wait()
	if got := h.probe.runs.Load(); got != int32(oversubscribed) {
		t.Errorf("宿主探针跑了 %d 次, want %d（排队的调用挤不进来＝天花板被说破的反面）", got, oversubscribed)
	}

	// The host probes above filled the arrival signal, and ticket 235's pre-read
	// is about the children only. wg.Wait() has returned, so every host send is
	// already buffered here and no child has been released yet: emptying the
	// channel now means whatever awaitChildOnBridge222 reads next came from a
	// child.
	h.drainBridgeArrivals()
	close(h.release)
	h.awaitChildOnBridge222(t)
	results := await222Tokens(t, "父任务拿到结论", h.parents, n)
	for i, res := range results {
		if res.IsError {
			t.Errorf("父任务 %d 在被反控压过之后收到了错误：%q", i, res.Text)
		}
	}
	if got := len(h.bridge.sem); got != 0 {
		t.Errorf("全部结束后桥位没回到空：%d, want 0（同一枚许可被还不回去的那一枚扣了两次，就是这一发先红）", got)
	}
}
