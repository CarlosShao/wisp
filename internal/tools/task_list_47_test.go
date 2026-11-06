package tools

// Ticket 47 AC#6's tools-side legs: task.list's frozen declaration shape
// and task.cancel's scheduler branch authority. The behavioural round-trip
// over the REAL loop lives in internal/agent/scheduler (bridge_test.go);
// these legs pin the shape the D34 row froze, in the package that owns it.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// fake47ListView is a three-row SchedulerView double (rows in, nothing else).
type fake47ListView struct {
	rows []TaskListRow
}

func (f fake47ListView) Rows() []TaskListRow { return f.rows }

// fake47Canceler records one Cancel call.
type fake47Canceler struct {
	called []string
}

func (f *fake47Canceler) Cancel(id string) error {
	f.called = append(f.called, id)
	return nil
}

// Test47TaskListDeclMatchesTheFrozenRow pins the D34 shape: L0 (the frozen
// pair "task.list / task.cancel | L0 / L1"), empty capability face (the
// fourth column is a "—"), no path parameter (a list takes no path), resident.
func Test47TaskListDeclMatchesTheFrozenRow(t *testing.T) {
	entries := BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster()})
	var listEntry *Entry
	for i := range entries {
		if entries[i].Tool.Name() == "task.list" {
			listEntry = &entries[i]
		}
	}
	if listEntry == nil {
		t.Fatalf("BuiltinTaskEntries 没注册 task.list（现名册：%v）—— AC#6 的接线不完整", entries)
	}
	decl := listEntry.Decl
	if decl.Declared != risk.L0 {
		t.Errorf("task.list 声明 = %v, want L0（D34 那行 PLAN.md:2564 冻的是 L0 / L1，list 取 L0）", decl.Declared)
	}
	if len(decl.Capabilities) != 0 || len(decl.Needs) != 0 {
		t.Errorf("task.list 的能力声明 = %v / %v, want 两枚都空（那行的第四列是「—」，读进程内表不碰 C3）", decl.Capabilities, decl.Needs)
	}
	if len(decl.PathParams) != 0 {
		t.Errorf("task.list 的 PathParams = %v, want 空：列队任务不带路径参数", decl.PathParams)
	}
	if !decl.Resident {
		t.Error("task.list 应是常驻（D15(1)：每回合注入目录）")
	}
}

// Test47TaskListAnswersRowsAndFailsClosed exercises the three answer shapes
// through the REAL bridge: rows rendered, empty is a real answer, unwired
// fails closed loudly.
func Test47TaskListAnswersRowsAndFailsClosed(t *testing.T) {
	view := fake47ListView{rows: []TaskListRow{
		{
			ID: "row-1", Label: "甲", State: statemachine.StateQueued,
			Note: "在等 row-2 释放 C:/w/锁.txt", ParentTaskID: "root-1",
		},
		{
			ID: "row-2", Label: "乙", State: statemachine.StateActing,
			HeldPaths: []string{"C:/w/锁.txt"},
		},
	}}
	reg := NewRegistry()
	for _, e := range BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster(), Scheduler: SchedulerDeps{View: view}}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	bridge := New(Options{Registry: reg, Logf: func(string, ...any) {}})
	out, err := bridge.Execute(context.Background(), req47("task.list", "root-1", ""))
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if out.IsError {
		t.Fatalf("task.list 回了错误：%s", out.Text)
	}
	for _, want := range []string{"调度器任务 2 行", "row-1", "甲", "Queued", "在等 row-2 释放 C:/w/锁.txt", "row-2", "Acting", "占用路径"} {
		if !strings.Contains(out.Text, want) {
			t.Errorf("task.list 回执缺 %q：\n%s", want, out.Text)
		}
	}

	// Empty queue is a real answer, not a failure.
	emptyReg := NewRegistry()
	for _, e := range BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster(), Scheduler: SchedulerDeps{View: fake47ListView{}}}) {
		if err := emptyReg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	emptyBridge := New(Options{Registry: emptyReg, Logf: func(string, ...any) {}})
	out2, err := emptyBridge.Execute(context.Background(), req47("task.list", "root-1", ""))
	if err != nil {
		t.Fatalf("dispatch 2: %v", err)
	}
	if out2.IsError || !strings.Contains(out2.Text, "队列为空") {
		t.Errorf("空队列回了 %+v, want「队列为空」的真话", out2)
	}

	// Unwired fails closed.
	unwiredReg := NewRegistry()
	for _, e := range BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster()}) {
		if err := unwiredReg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	unwiredBridge := New(Options{Registry: unwiredReg, Logf: func(string, ...any) {}})
	out3, err := unwiredBridge.Execute(context.Background(), req47("task.list", "root-1", ""))
	if err != nil {
		t.Fatalf("dispatch 3: %v", err)
	}
	if !out3.IsError || !strings.Contains(out3.Text, "调度器未接线") {
		t.Errorf("未接线回了 %+v, want fail-closed", out3)
	}
}

// Test47TaskCancelSchedulerBranchAuthority pins the branch's authority shape:
// parent stops its child; foreign and self are refused; the stop reaches the
// wired Control exactly once.
func Test47TaskCancelSchedulerBranchAuthority(t *testing.T) {
	view := fake47ListView{rows: []TaskListRow{
		{ID: "sched-parent", Label: "父", State: statemachine.StateActing},
		{ID: "sched-child", Label: "孩子", State: statemachine.StateActing, ParentTaskID: "sched-parent"},
	}}
	canceler := &fake47Canceler{}
	reg := NewRegistry()
	for _, e := range BuiltinTaskEntries(TaskDeps{
		Roster:    NewTaskRoster(),
		Scheduler: SchedulerDeps{View: view, Control: canceler},
	}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	bridge := New(Options{
		Registry: reg,
		Gate:     GateFuncs{Window: func(context.Context, Decision) (Answer, string) { return AnswerTimeout, "" }},
		Logf:     func(string, ...any) {},
	})

	// The parent stops its own child: reaches Control, answers success.
	out, err := bridge.Execute(context.Background(), req47("task.cancel", "sched-parent", "sched-child"))
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	if out.IsError || !strings.Contains(out.Text, "sched-child") {
		t.Errorf("父停孩子回了 %+v", out)
	}
	if len(canceler.called) != 1 || canceler.called[0] != "sched-child" {
		t.Errorf("Control 收到的停止 = %v, want 恰 [sched-child]", canceler.called)
	}

	// A foreign caller is refused, and Control is NOT reached again.
	out2, err := bridge.Execute(context.Background(), req47("task.cancel", "别人", "sched-child"))
	if err != nil {
		t.Fatalf("dispatch 2: %v", err)
	}
	if !out2.IsError || !strings.Contains(out2.Text, "拒绝停止") {
		t.Errorf("外人停孩子回了 %+v, want 拒绝", out2)
	}
	// Self-stop is refused.
	out3, err := bridge.Execute(context.Background(), req47("task.cancel", "sched-parent", "sched-parent"))
	if err != nil {
		t.Fatalf("dispatch 3: %v", err)
	}
	if !out3.IsError || !strings.Contains(out3.Text, "拒绝停止") {
		t.Errorf("自停回了 %+v, want 拒绝", out3)
	}
	if len(canceler.called) != 1 {
		t.Errorf("被拒的调用不该到达 Control：现 %v", canceler.called)
	}

	// A root row (no parent) is refused for anyone: nobody is its parent.
	rootView := fake47ListView{rows: []TaskListRow{
		{ID: "sched-root", Label: "根", State: statemachine.StateActing},
	}}
	rootReg := NewRegistry()
	for _, e := range BuiltinTaskEntries(TaskDeps{
		Roster:    NewTaskRoster(),
		Scheduler: SchedulerDeps{View: rootView, Control: &fake47Canceler{}},
	}) {
		if err := rootReg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	rootBridge := New(Options{
		Registry: rootReg,
		Gate:     GateFuncs{Window: func(context.Context, Decision) (Answer, string) { return AnswerTimeout, "" }},
		Logf:     func(string, ...any) {},
	})
	out4, err := rootBridge.Execute(context.Background(), req47("task.cancel", "sched-root", "sched-root"))
	if err != nil {
		t.Fatalf("dispatch 4: %v", err)
	}
	if !out4.IsError || !strings.Contains(out4.Text, "拒绝停止") {
		t.Errorf("根行自停回了 %+v, want 拒绝", out4)
	}
}

// req47 builds one tool request with the caller identity stamped (the way
// the bridge receives it from the loop) and, for task.cancel, the target id.
func req47(name, callerID, targetID string) agent.ToolRequest {
	args := `{}`
	if name == "task.cancel" {
		args = fmt.Sprintf(`{"task_id":%q}`, targetID)
	}
	return agent.ToolRequest{
		TaskID: callerID, CorrelationID: callerID,
		CallID: "call-47-" + name + "-" + callerID,
		Name:   name, Args: json.RawMessage(args),
	}
}
