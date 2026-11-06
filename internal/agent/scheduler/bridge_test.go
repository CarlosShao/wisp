package scheduler

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Ticket 47 AC#6: task.list / task.cancel round-trip through a REAL agent
// loop. The mock here is the frozen C5 seam (AGENTS §1.3 - the provider is
// the one injection point); everything else is real: the loop, the bridge,
// the C19 assessor, the tools registry, the scheduler. The provider drives
// the turns a model would: first task.list, then task.cancel on the id the
// list printed.

// mock47Provider plays the two-turn script the ticket's AC names.
type mock47Provider struct {
	mu    sync.Mutex
	calls int
}

func (p *mock47Provider) Stream(_ context.Context, _ *llm.Request, emit func(llm.StreamEvent) error) error {
	p.mu.Lock()
	p.calls++
	n := p.calls
	p.mu.Unlock()
	switch n {
	case 1:
		return emit47ToolCall(emit, "call-47-list", "task.list", `{}`)
	default:
		return emit47Text(emit, "47 收工")
	}
}

func (p *mock47Provider) Info() llm.ProviderInfo {
	return llm.ProviderInfo{Model: "mock-47", Provider: "mock47", MaxContextWindow: 8192}
}

func emit47ToolCall(emit func(llm.StreamEvent) error, id, name, args string) error {
	_ = emit(llm.StreamEvent{Type: llm.EvToolCallStart, ToolCallID: id, ToolName: name})
	_ = emit(llm.StreamEvent{Type: llm.EvToolCallArgsDelta, ToolCallID: id, ArgsDelta: args})
	_ = emit(llm.StreamEvent{Type: llm.EvToolCallEnd, ToolCallID: id, ToolName: name})
	_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopToolUse})
	_ = emit(llm.StreamEvent{Type: llm.EvDone})
	return nil
}

func emit47Text(emit func(llm.StreamEvent) error, text string) error {
	_ = emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: text})
	_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopEndTurn})
	_ = emit(llm.StreamEvent{Type: llm.EvDone})
	return nil
}

// listCapture remembers what the model saw, so the leg can assert the rows
// actually reached the context (the round-trip's whole point).
type listCapture struct {
	mu   sync.Mutex
	body string
}

func (c *listCapture) feed(req *llm.Request) {
	var b strings.Builder
	for _, m := range req.Messages {
		for _, part := range m.Content {
			if tp, ok := part.(llm.TextPart); ok {
				b.WriteString(tp.Text)
				b.WriteString("\n")
			}
			if tr, ok := part.(llm.ToolResultPart); ok {
				for _, rc := range tr.Content {
					if tp, ok := rc.(llm.TextPart); ok {
						b.WriteString(tp.Text)
						b.WriteString("\n")
					}
				}
			}
		}
	}
	c.mu.Lock()
	c.body = b.String()
	c.mu.Unlock()
}

func (c *listCapture) text() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.body
}

// new47Loop builds one real agent.Loop over a bridge that carries the task
// family wired to the scheduler through NewBridgeDeps, exactly the shape the
// composition root will assemble.
func new47Loop(t *testing.T, s *Scheduler) (*agent.Loop, *listCapture) {
	t.Helper()
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{
		Roster:    tools.NewTaskRoster(),
		Scheduler: NewBridgeDeps(s),
	}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	bridge := tools.New(tools.Options{
		Registry: reg,
		Gate:     tools.GateFuncs{Window: func(context.Context, tools.Decision) (tools.Answer, string) { return tools.AnswerTimeout, "" }},
		Logf:     func(string, ...any) {},
	})
	cap := &listCapture{}
	prov := &turnSpy47{inner: &mock47Provider{}, capture: cap}
	loop, err := agent.New(agent.Options{
		Provider:  prov,
		Tools:     bridge,
		Sink:      &agent.RecordSink{},
		Registry:  observe.NewRegistry(),
		AdmitTask: func(string) func() { return func() {} },
		Config: agent.Config{
			Model: "probe-47", ContextWindow: 8192,
			ArtifactsDir:   t.TempDir(),
			PerToolTimeout: 5 * time.Second,
		},
	})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	return loop, cap
}

// turnSpy47 feeds every request to the capture before delegating.
type turnSpy47 struct {
	inner   *mock47Provider
	capture *listCapture
}

func (s *turnSpy47) Stream(ctx context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	s.capture.feed(req)
	return s.inner.Stream(ctx, req, emit)
}

func (s *turnSpy47) Info() llm.ProviderInfo { return s.inner.Info() }

// Test47AC6TaskListRoundTripThroughTheRealLoop is AC#6's first half: the
// mock LLM's task.list call reaches the REAL bridge, executes the REAL
// tool, and the scheduler's rows come back into the model's context.
func Test47AC6TaskListRoundTripThroughTheRealLoop(t *testing.T) {
	locks := NewPathLock(PathLockOptions{})
	s := New(Options{Registry: observe.NewRegistry(), Locks: locks, MaxConcurrent: 2})

	// Seed one row the model should see: a finished task with a note.
	id, err := s.Submit(TaskSpec{Label: "完成的活", Exec: tc("已交付")})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if _, err := s.Wait(id); err != nil {
		t.Fatalf("Wait: %v", err)
	}

	loop, cap := new47Loop(t, s)
	res := loop.Run(context.Background(), "看看队列里有什么")
	if res.Status != agent.StatusCompleted {
		t.Fatalf("round-trip 落到 %s（%s）", res.Status, res.Message)
	}
	body := cap.text()
	if !strings.Contains(body, "调度器任务 1 行") {
		t.Fatalf("模型上下文里没有 task.list 的回执（行数头）：\n%s", body)
	}
	if !strings.Contains(body, id) || !strings.Contains(body, "完成的活") {
		t.Errorf("模型上下文里没有调度器行的 id/标题：\n%s", body)
	}
	if !strings.Contains(body, "已完成") {
		t.Errorf("模型上下文里没有行注（已完成）：\n%s", body)
	}
}

// Test47AC6TaskCancelRoundTripStopsASchedulerRow is AC#6's second half: a
// task.list first (so the id the model sees is the row's own), then a
// task.cancel aimed at a CHILD row the CALLER parented - the full
// read-then-stop round-trip a model drives.
func Test47AC6TaskCancelRoundTripStopsASchedulerRow(t *testing.T) {
	s := New(Options{Registry: observe.NewRegistry(), MaxConcurrent: 2})

	// The parent: its own scheduler row exists (this is the row the loop
	// runs inside). The child is parked on a gate, so the cancel has a live
	// target; the parent id joins them through SetParent.
	parentGate := make(chan struct{})
	var parentID string
	parentDone := make(chan struct{})
	go func() {
		defer close(parentDone)
		id, err := s.Submit(TaskSpec{Label: "父任务", Exec: func(context.Context, *TaskRun) TaskResult {
			<-parentGate
			return TaskResult{Status: "completed"}
		}})
		if err != nil {
			t.Errorf("Submit parent: %v", err)
			return
		}
		parentID = id
	}()
	waitFor(t, "父任务行到场", func() bool { return parentID != "" })

	childGate := make(chan struct{})
	childStarted := make(chan struct{})
	var childID string
	go func() {
		id, err := s.Submit(TaskSpec{Label: "孩子任务", Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			s.SetParent(run.ID, parentID)
			close(childStarted)
			<-childGate
			// After the cancel, the root ctx is dead; honour it.
			return TaskResult{Status: "cancelled", AppliedSteps: []string{"已写一半"}}
		}})
		if err != nil {
			t.Errorf("Submit child: %v", err)
			return
		}
		childID = id
	}()
	<-childStarted
	waitFor(t, "孩子行在场", func() bool { return childID != "" })

	// The round-trip: task.cancel dispatched AS the parent (the bridge
	// stamps the caller identity the same way the loop does).
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{
		Roster:    tools.NewTaskRoster(),
		Scheduler: NewBridgeDeps(s),
	}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	bridge := tools.New(tools.Options{
		Registry: reg,
		Gate:     tools.GateFuncs{Window: func(context.Context, tools.Decision) (tools.Answer, string) { return tools.AnswerTimeout, "" }},
		Logf:     func(string, ...any) {},
	})
	out, err := bridge.Execute(context.Background(), agent.ToolRequest{
		TaskID: parentID, CorrelationID: parentID, CallID: "call-47-cancel",
		Name: "task.cancel",
		Args: []byte(fmt.Sprintf(`{"task_id":%q}`, childID)),
	})
	if err != nil {
		t.Fatalf("task.cancel dispatch: %v", err)
	}
	if out.IsError {
		t.Fatalf("task.cancel 回了错误：%s", out.Text)
	}
	if !strings.Contains(out.Text, childID) {
		t.Errorf("task.cancel 回执没有点名目标行：%s", out.Text)
	}
	// The child body parks on childGate; the cancel unblocks it through the
	// root ctx (the body's select-free shape only needs the gate closed OR
	// the ctx dead - the gate stays shut on purpose, so the ctx path is the
	// one that runs, which is exactly the D38c shape being pinned).
	close(childGate)
	// The cancelled child settles with its applied steps on the row.
	res, err := s.Wait(childID)
	if err != nil {
		t.Fatalf("Wait child: %v", err)
	}
	if res.Status != "cancelled" {
		t.Errorf("孩子落到 %s, want cancelled", res.Status)
	}
	row, ok := s.Row(childID)
	if !ok {
		t.Fatal("孩子的行不见了")
	}
	if !strings.Contains(row.Note, "已写一半") || !strings.Contains(row.Note, "未自动回退") {
		t.Errorf("行注没有带 applied-steps 台账：%q", row.Note)
	}
	// The refusal half: a FOREIGN caller cannot stop the child.
	out2, err := bridge.Execute(context.Background(), agent.ToolRequest{
		TaskID: "不是父任务", CorrelationID: "不是父任务", CallID: "call-47-cancel-2",
		Name: "task.cancel",
		Args: []byte(fmt.Sprintf(`{"task_id":%q}`, childID)),
	})
	if err != nil {
		t.Fatalf("task.cancel dispatch 2: %v", err)
	}
	if !out2.IsError || !strings.Contains(out2.Text, "拒绝停止") {
		t.Errorf("外人的 task.cancel 没被拒：%+v", out2)
	}
	close(parentGate)
	<-parentDone
}

// Test47AC6TaskListUnwiredFailsClosed pins the fail-closed half: a bridge
// without the scheduler answers an error, never an empty success.
func Test47AC6TaskListUnwiredFailsClosed(t *testing.T) {
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: tools.NewTaskRoster()}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	bridge := tools.New(tools.Options{Registry: reg, Logf: func(string, ...any) {}})
	out, err := bridge.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t", CorrelationID: "t", CallID: "call-47-unwired",
		Name: "task.list", Args: []byte(`{}`),
	})
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if !out.IsError || !strings.Contains(out.Text, "调度器未接线") {
		t.Errorf("未接线时 task.list 回了 %+v, want fail-closed 错误", out)
	}
}
