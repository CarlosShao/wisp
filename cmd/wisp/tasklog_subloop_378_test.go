package main

// Ticket 378〔乙〕: a spawned child task opens ONE task_log row and writes NO
// tool_call rows.
//
// WHY THIS CASE SITS IN cmd/wisp. The decision this ticket makes is not in
// internal/agent - the loop already writes whatever Journal it is handed - but in
// THIS composition root's wiring: which Journal payload task.spawn's child
// inherits from rt.loopOpt. A test in internal/agent that hands a Loop a journal
// would prove the loop's half, which was never in doubt, and would stay green
// with the wiring reverted. So nothing below builds a child agent.Options by
// hand: the child is born from the real bridge, through the real task.spawn
// registered below, off an Options snapshot shaped exactly like the one
// execute() publishes. Every number is read back out of a real SQLite store
// through the real DAO.
//
// WHAT IT DISCRIMINATES. "Only the task row" and "the tool calls were recorded
// too" differ by exactly the shape [B]/[C] pin:
//
//	[A] child-rows delta      == number of spawned children  (one row per child)
//	[B] child tool_call rows  == 0                            (the inherited journal
//	                                                          half is inert)
//	[B'] the child's tool REALLY ran                          (so [B] is not "nothing
//	                                                          happened to be booked")
//	[C] parent-rows delta     == 0                            (adding the child row did
//	                                                          not disturb the root's
//	                                                          own rowing)
//	[D] parent tool_call rows == 0                            (this assembly books the
//	                                                          bridge's rows, not the
//	                                                          loop's: unchanged)
//
// A wiring that hands the child a *memory.Store (or any journal whose tool_call
// methods write) turns [B] red while [A] stays green - the "记录页刷屏" shape
// ticket 378's hard reminder #1 names, which [A] alone cannot see.

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/tools"
)

const (
	sub378RootID  = "root-378-sub-loop"
	sub378ChildID = "child-378-sub-loop"
	sub378Label   = "子环路取数"
	// sub378ChildTool is the tool the CHILD calls. It is a real, named call, so
	// the bridge (and any journal that wrote) would have booked it.
	sub378ChildTool = "echo"
	// sub378ChildEcho is what the child's tool returns, distinctive so the tool
	// call count proves execution rather than a skipped plan.
	sub378ChildEcho = "子环路工具真的跑过-378"
)

// sub378Provider is the C5 seam for the CHILD: the first response asks for one
// tool call, the second answers with text. It is the fake only because the
// golden fixtures this package's sibling uses live in internal/agent's
// testdata; the loop, the journal, the bridge and the store are the real ones.
type sub378Provider struct {
	calls int
}

func (p *sub378Provider) Stream(_ context.Context, _ *llm.Request, emit func(llm.StreamEvent) error) error {
	p.calls++
	if p.calls == 1 {
		args, _ := json.Marshal(map[string]string{"text": sub378ChildEcho})
		_ = emit(llm.StreamEvent{Type: llm.EvToolCallStart, ToolCallID: "call-378-child", ToolName: sub378ChildTool})
		_ = emit(llm.StreamEvent{Type: llm.EvToolCallArgsDelta, ToolCallID: "call-378-child", ArgsDelta: string(args)})
		_ = emit(llm.StreamEvent{Type: llm.EvToolCallEnd, ToolCallID: "call-378-child", ToolName: sub378ChildTool})
		_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopToolUse})
		_ = emit(llm.StreamEvent{Type: llm.EvDone})
		return nil
	}
	_ = emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: "子任务做完了。"})
	_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopEndTurn})
	_ = emit(llm.StreamEvent{Type: llm.EvDone})
	return nil
}

func (p *sub378Provider) Info() llm.ProviderInfo {
	return llm.ProviderInfo{Model: "fake-378", Provider: "fake-378", MaxContextWindow: 128000}
}

// sub378Dir is the C4 surface the child is handed (the parent's bridge, wrapped
// by task.spawn). It records executions so [B'] can prove the call really ran.
type sub378Dir struct {
	calls []string
}

func (d *sub378Dir) Tools(context.Context) ([]agent.ToolInfo, error) {
	return []agent.ToolInfo{{Name: sub378ChildTool, RiskLevel: "L0", Resident: true}}, nil
}

func (d *sub378Dir) Execute(_ context.Context, req agent.ToolRequest) (agent.ToolOutcome, error) {
	d.calls = append(d.calls, req.Name)
	var a struct {
		Text string `json:"text"`
	}
	_ = json.Unmarshal(req.Args, &a)
	return agent.ToolOutcome{Text: a.Text}, nil
}

// Test378ChildLoopWritesOneTaskRowAndNoToolCallRows is the ticket's one directed
// case. See the header for what each assertion discriminates.
func Test378ChildLoopWritesOneTaskRowAndNoToolCallRows(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	store, err := memory.Open(t.TempDir(),
		memory.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	defer store.Close()

	dir := &sub378Dir{}
	roster := tools.NewTaskRoster()
	roster.MarkRoot(sub378RootID, "根任务")

	// The one journal payload this composition root hands over (run.go's
	// taskLifecycleJournal): StartTaskLog writes, everything else is inert.
	// task.spawn's child copies the snapshot verbatim (subagent_197.go:320,
	// `opt := base`), which is the wiring under test.
	// One provider per child loop. A single shared counter would make the two
	// children share "call #1 / call #2", so the second child would answer with
	// text and never call its tool - which would leave [B'] reading 1 instead of
	// the per-child shape this case is about.
	base := agent.Options{
		Provider: nil,
		Tools:    dir,
		Journal:  &taskLifecycleJournal{store: store},
		AdmitTask: func(string) func() {
			return func() {}
		},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Sink:   agent.NopSink{},
		Config: agent.Config{
			Model:                       "fake-378",
			ContextWindow:               128000,
			PassThroughUnclassifiedRisk: true,
			Currency:                    "CNY",
			PerToolTimeout:              5 * time.Second,
		},
	}
	baseOptions := func() (agent.Options, bool) {
		o := base
		o.Provider = &sub378Provider{}
		return o, true
	}
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinSubagentEntries(tools.SubagentDeps{
		Roster:      roster,
		ParentTools: dir,
		BaseOptions: baseOptions,
	}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	bridge := tools.New(tools.Options{
		Registry: reg, Gate: tools.NoGate{},
		Logf: func(string, ...any) {},
	})

	tasksBefore := listTaskRows378(t, ctx, store)
	callsBefore := len(listToolCallRows378(t, ctx, store))

	const spawns = 2
	for i := 0; i < spawns; i++ {
		out, err := bridge.Execute(ctx, agent.ToolRequest{
			TaskID: sub378RootID, CorrelationID: sub378RootID,
			CallID: sub378ChildID + "-spawn",
			Name:   "task.spawn",
			Args: json.RawMessage(`{"description":"` + sub378Label +
				`","prompt":"取一次数，然后收尾"}`),
		})
		if err != nil {
			t.Fatalf("task.spawn #%d: %v", i+1, err)
		}
		if out.IsError {
			t.Fatalf("task.spawn #%d refused: %s", i+1, out.Text)
		}
	}

	tasksAfter := listTaskRows378(t, ctx, store)
	callsAfter := listToolCallRows378(t, ctx, store)

	// [A] exactly one task_log row per spawned child.
	if got := len(tasksAfter) - len(tasksBefore); got != spawns {
		t.Errorf("[A] task_log 新增行 = %d，want %d（每个子任务一行）：%+v",
			got, spawns, relIDs378(tasksAfter, tasksBefore))
	}
	// [C] the root's own rowing did not move: this ticket added children, not a
	// second root row, and the root is not among the rows this spawn produced.
	rootRowsAfter := 0
	for _, tl := range tasksAfter {
		if tl.ID == sub378RootID {
			rootRowsAfter++
		}
	}
	rootRowsBefore := 0
	for _, tl := range tasksBefore {
		if tl.ID == sub378RootID {
			rootRowsBefore++
		}
	}
	if rootRowsAfter != rootRowsBefore {
		t.Errorf("[C] 根任务在 task_log 的行数 %d -> %d，want 不变（本票只加子任务行）",
			rootRowsBefore, rootRowsAfter)
	}

	// The children are the roster's descendants: the row count above is about
	// the real child tasks, not about anything else that happened to write.
	kids := roster.Descendants(sub378RootID)
	if len(kids) != spawns {
		t.Fatalf("名册里子代理行数 = %d，want %d（%+v），后面的读数就没有对准子任务",
			len(kids), spawns, kids)
	}

	// [B] the decisive half: no tool_call row for any child, even though the
	// child called sub378ChildTool exactly once per turn shape.
	for _, k := range kids {
		rows, err := store.ListToolCallsByTask(ctx, k.TaskID)
		if err != nil {
			t.Fatalf("ListToolCallsByTask(%s): %v", k.TaskID, err)
		}
		if len(rows) != 0 {
			var names []string
			for _, r := range rows {
				names = append(names, r.Tool)
			}
			t.Errorf("[B] 子任务 %s 写了 %d 行 tool_call（%v），want 0："+
				"继承的 journal 只该开任务行，不该把子环路的每次工具调用都记成取证行",
				k.TaskID, len(rows), names)
		}
	}
	// [D] and the assembly as a whole booked exactly its own shape (zero here:
	// this test wires the bridge without a journal, so the loop's inert half is
	// the only journal in play). This is the "工具调用数 -> 行数" relation staying
	// put: adding the child row changed no tool_call count.
	if got := len(callsAfter) - callsBefore; got != 0 {
		t.Errorf("[D] tool_call 总行数 +%d，want 0：本票只加任务行，不许改工具调用数→行数的关系", got)
	}

	// [B'] the negative half is not vacuous: the child's tool really executed.
	if len(dir.calls) != spawns {
		t.Errorf("[B'] 子环路实际执行工具 %d 次，want %d：子工具没跑起来的话 [B] 的 0 行读不出任何东西",
			len(dir.calls), spawns)
	}
	for _, k := range kids {
		row, ok := roster.Look(k.TaskID)
		if !ok {
			t.Errorf("名册查不到子任务 %s", k.TaskID)
			continue
		}
		if row.Text != "子任务做完了。" {
			t.Errorf("子任务 %s 的结论 = %q，want %q", k.TaskID, row.Text, "子任务做完了。")
		}
	}
}

// The three helpers below are the store read surface this case needs, kept
// deliberately thin so the assertions above speak in row counts.

func listTaskRows378(t *testing.T, ctx context.Context, store *memory.Store) []memory.TaskLog {
	t.Helper()
	rows, err := store.ListTaskLogs(ctx)
	if err != nil {
		t.Fatalf("ListTaskLogs: %v", err)
	}
	return rows
}

func listToolCallRows378(t *testing.T, ctx context.Context, store *memory.Store) []memory.ToolCall {
	t.Helper()
	rows, err := store.ListToolCalls(ctx)
	if err != nil {
		t.Fatalf("ListToolCalls: %v", err)
	}
	return rows
}

func relIDs378(after, before []memory.TaskLog) []string {
	seen := make(map[string]bool, len(before))
	for _, tl := range before {
		seen[tl.ID] = true
	}
	var out []string
	for _, tl := range after {
		if !seen[tl.ID] {
			out = append(out, tl.ID+"/"+tl.State)
		}
	}
	return out
}
