package tools

// 票 283 腿 283-r1：corr 落地后两处新命名盲区的仪器（补尺，非"登记不测"）。
//
// 盲区一（回退形无仪器）：票 242 把任务身份挪到 TaskID(ctx) 之后，
// cancel.go 的 TaskID 支若被种回 h.corr，既有用例全绿——它们的派发 fixture
// 一律 TaskID == CorrelationID（bridge.go 的 orDefault 回落也让缺 corr 的
// 构造塌成同一值），分不出两者。尺一 = 最小单元尺：一个 corr 与 taskID
// 刻意不同的 handle。
//
// 盲区二（身份改读无整链尺）：subagent_197.go 与 task.go 的任务身份改读
// TaskID(ctx) 前后，同一形状用例双双全绿——它们都经直接 bridge.Execute
// 派发、且 corr == taskID。真生产链是 agent.Loop 的 callCorr(taskID, callID, i)
// （loop.go 每枚调用一个 corr，形如 taskID + "#" + callID）→ bridge 把 corr
// 与 taskID 两枚分别写进 cancelHandle（bridge.go 的 withCancel 点）→ 工具读
// TaskID(ctx) 定位名册行。尺二 = 真 loop → 真 bridge → 真 task.spawn /
// task.cancel → 真名册 的整链：父环路先派生一枚子代理、再停它自己那一枚；
// 两枚调用在真 loop 下的 corr 都不等于 task id。种坏三处中的任一处，这条链
// 的指名断言变红（读数见 .scratch/wisp/probes/283/r1/10-positive-controls.md）。

import (
	"context"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// --- 尺一：handle 的 corr 与 taskID 不同时，TaskID 必须答 taskID --------

func Test283TaskIDAccessorIsTheTaskIDNotTheCorr(t *testing.T) {
	ctx := withCancel(context.Background(), cancelHandle{corr: "corr-283", taskID: "task-283"})
	if got := TaskID(ctx); got != "task-283" {
		t.Fatalf("TaskID = %q, want the task id: a handle whose corr differs from its task id must not answer with the corr", got)
	}
	if got := CorrelationID(ctx); got != "corr-283" {
		t.Fatalf("CorrelationID = %q, want corr-283", got)
	}
	// 回退形之二：没铸进 task id 的 handle 不许拿 corr 顶。两处身份读点的
	// fail-closed 支（spawn 的 parentID == ""、cancel 的 caller == ""）就建立在
	// 「空就是空」上；回落到 corr 会把「没有宿主给的任务 id」静默改写成一个
	// 身份声明，还会顺手绕过那两声拒绝。
	loose := withCancel(context.Background(), cancelHandle{corr: "corr-only-283"})
	if got := TaskID(loose); got != "" {
		t.Fatalf("TaskID = %q, want \"\": a handle with no task id stamped must answer empty, not fall back onto the corr", got)
	}
}

// --- 尺二：真 loop 的 per-call corr 整链 ---------------------------------

// parent283Provider is the C5 seam driving the PARENT loop through two tool
// turns: first a task.spawn, then a task.cancel aimed at the child id the
// spawned tool's own reply carries back to the model (subagent_197.go's
// answer() prints "子代理 <id>（<label>）已结束..."), so the second call uses
// the same id a model would read in production instead of a pinned constant.
type parent283Provider struct {
	mu       sync.Mutex
	calls    int
	cancelID string
}

func (p *parent283Provider) Stream(_ context.Context, req *llm.Request, emit func(llm.StreamEvent) error) error {
	p.mu.Lock()
	p.calls++
	n := p.calls
	p.mu.Unlock()
	switch n {
	case 1:
		return emit283ToolCall(emit, "call-283-spawn", "task.spawn",
			`{"description":"283 子代理","prompt":"跑完并给出结论"}`)
	case 2:
		id := childID283(req.Messages)
		p.mu.Lock()
		p.cancelID = id
		p.mu.Unlock()
		if id == "" {
			// 找不到 id 就走不完取消那一跳：这一程按"没测成"收口（断言先报
			// 这一条），而不是静默留一枚假绿。
			return emit283Text(emit, "283: 没有从工具回执里拿到子代理 id")
		}
		return emit283ToolCall(emit, "call-283-cancel", "task.cancel",
			`{"task_id":"`+id+`"}`)
	default:
		return emit283Text(emit, "283 收工")
	}
}

func (p *parent283Provider) Info() llm.ProviderInfo {
	return llm.ProviderInfo{Model: "fake-283", Provider: "parent283", MaxContextWindow: 8192}
}

func (p *parent283Provider) cancelTarget() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cancelID
}

func emit283ToolCall(emit func(llm.StreamEvent) error, id, name, args string) error {
	_ = emit(llm.StreamEvent{Type: llm.EvToolCallStart, ToolCallID: id, ToolName: name})
	_ = emit(llm.StreamEvent{Type: llm.EvToolCallArgsDelta, ToolCallID: id, ArgsDelta: args})
	_ = emit(llm.StreamEvent{Type: llm.EvToolCallEnd, ToolCallID: id, ToolName: name})
	_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopToolUse})
	_ = emit(llm.StreamEvent{Type: llm.EvDone})
	return nil
}

func emit283Text(emit func(llm.StreamEvent) error, text string) error {
	_ = emit(llm.StreamEvent{Type: llm.EvTextDelta, Text: text})
	_ = emit(llm.StreamEvent{Type: llm.EvStop, Stop: llm.StopEndTurn})
	_ = emit(llm.StreamEvent{Type: llm.EvDone})
	return nil
}

// childID283 lifts the child id out of the conversation the loop hands the
// provider: the spawn reply reads "子代理 <id>（<label>）已结束", so the id is
// the run between that marker and the label's opening parenthesis.
func childID283(msgs []llm.Message) string {
	re := regexp.MustCompile(`子代理 (\S+?)（`)
	for _, m := range msgs {
		for _, c := range m.Content {
			switch v := c.(type) {
			case llm.TextPart:
				if sub := re.FindStringSubmatch(v.Text); sub != nil {
					return sub[1]
				}
			case llm.ToolResultPart:
				for _, rc := range v.Content {
					if tp, ok := rc.(llm.TextPart); ok {
						if sub := re.FindStringSubmatch(tp.Text); sub != nil {
							return sub[1]
						}
					}
				}
			}
		}
	}
	return ""
}

func Test283IdentityChainThroughTheRealLoop(t *testing.T) {
	x := build221(t, windowGate221(), func() llm.LlmProvider {
		return &fake197Provider{name: "child283", answer: "283 号子代理的结论"}
	}, true)

	p := &parent283Provider{}
	loop, err := agent.New(agent.Options{
		Provider:  p,
		Tools:     x.h.bridge,
		Sink:      &agent.RecordSink{},
		Registry:  observe.NewRegistry(),
		AdmitTask: func(string) func() { return func() {} },
		Config: agent.Config{
			Model: "probe-283", ContextWindow: 8192,
			ArtifactsDir:                filepath.Join(t.TempDir(), "artifacts"),
			PerToolTimeout:              10 * time.Second,
			PassThroughUnclassifiedRisk: false,
		},
	})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}

	res := loop.Run(t.Context(), "先派一枚子代理，然后把它停掉")
	if res.Status != agent.StatusCompleted {
		t.Fatalf("status = %s (%s)", res.Status, res.Message)
	}
	if res.ToolCalls != 2 {
		t.Fatalf("ToolCalls = %d, want 2 (task.spawn 后 task.cancel)", res.ToolCalls)
	}
	for _, l := range res.ToolLog {
		if l.Rejected {
			t.Fatalf("真 loop 拒了一枚它该放行的调用: %+v", l)
		}
	}
	childID := p.cancelTarget()
	if childID == "" {
		t.Fatalf("这程没有测成：provider 没能从工具回执里拿到子代理 id")
	}

	// 孩子那一行：父链接必须是宿主给的那枚 task id（真 loop 铸的
	// res.TaskID），不是这次调用的 per-call corr。种坏 subagent_197.go 的
	// 身份读点、或 cancel.go 的 TaskID 访问器，这条先红。
	rec, ok := x.h.roster.Look(childID)
	if !ok {
		t.Fatalf("名册里查不到子代理 %s", childID)
	}
	if rec.ParentTaskID != res.TaskID {
		t.Errorf("子代理行的 ParentTaskID = %q, want 宿主 task id %q（task.spawn 读的是 TaskID(ctx)，不是 per-call corr）",
			rec.ParentTaskID, res.TaskID)
	}
	if rec.Kind != TaskKindSubagent {
		t.Errorf("Kind = %q, want %q", rec.Kind, TaskKindSubagent)
	}

	// 取消那一跳：调用者是同一个 task id，目标是它的孩子。孩子此时已跑完，
	// 所以权威判定通过、停不到句柄的形状是「没能停掉」；「拒绝停止」只该在
	// 调用者身份读错时出现（task.go 的身份读点读 corr 会走到那儿）。
	cancelRes := res.ToolLog[1]
	if !strings.Contains(cancelRes.Text, "没能停掉") {
		t.Errorf("task.cancel 的回执没有走到权威判定通过那一支: %q", cancelRes.Text)
	}
	if strings.Contains(cancelRes.Text, "拒绝停止") {
		t.Errorf("父任务停自己的孩子被身份判定拒了（调用者读到的不是 task id）: %q", cancelRes.Text)
	}

	// 真 loop 的 per-call corr 确实到达 bridge 的记账层：每枚调用一行，
	// corr 非空、以 task id 为前缀、且不等于 task id 本身（票 242 的铸形）。
	rows, err := x.mem.ListToolCallsByTask(t.Context(), res.TaskID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("tool_call 行数 = %d, want 2 (一枚调用一行): %+v", len(rows), rows)
	}
	for _, r := range rows {
		if r.CorrelationID == "" || r.CorrelationID == res.TaskID || !strings.HasPrefix(r.CorrelationID, res.TaskID) {
			t.Errorf("tool %s 的 correlation_id = %q, want 非空 per-call id（前缀 = task id %q，且不等于它）",
				r.Tool, r.CorrelationID, res.TaskID)
		}
	}
}
