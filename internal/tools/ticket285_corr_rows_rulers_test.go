package tools

// 票 285 AC#4 第五形：per-call corr 的"每枚调用各不相同"这件事的尺（落点乙）。
//
// 续的是票 283 那条真链（真 agent.Loop → 真 bridge → 真工具 → 真 tool_call
// 名册行）。票 283 的尺二逐枚钉了非空／不等于 task id／以 task id 为前缀，
// 没有把两枚互相作差：把 callCorr 带 call id 那一支种成同一枚常量后缀，
// 那一枚尺逐枚看仍然全绿，名册里两枚调用已经塌成一枚键。
// 这一枚补的就是那一发差集，断言面落在持久化的名册行上（不是内存里的请求）。
//
// 夹具前置（票面 AC#4 那句）：TaskID 与 CorrelationID 必须是两枚值，否则两发
// 必同绿不算读数——这一条在本用例里也是断言，不是注释。

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Test285RosterRowsCarryDistinctCorrPerCall drives one task through the real
// loop with two tool calls and requires the two tool_call rows to carry two
// different correlation ids, each routing back to its own call.
func Test285RosterRowsCarryDistinctCorrPerCall(t *testing.T) {
	x := build221(t, windowGate221(), func() llm.LlmProvider {
		return &fake197Provider{name: "child285", answer: "285 号子代理的结论"}
	}, true)

	p := &parent283Provider{}
	loop, err := agent.New(agent.Options{
		Provider:  p,
		Tools:     x.h.bridge,
		Sink:      &agent.RecordSink{},
		Registry:  observe.NewRegistry(),
		AdmitTask: func(string) func() { return func() {} },
		Config: agent.Config{
			Model: "probe-285", ContextWindow: 8192,
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

	rows, err := x.mem.ListToolCallsByTask(t.Context(), res.TaskID)
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("tool_call 行数 = %d, want 2 (一枚调用一行): %+v", len(rows), rows)
	}

	// 夹具前置：两枚行各自的 corr 都不得等于 task id，也不得为空。
	for _, r := range rows {
		if r.CorrelationID == "" {
			t.Fatalf("tool %s 的行 corr 为空，名册无从按 corr 路由", r.Tool)
		}
		if r.CorrelationID == res.TaskID {
			t.Fatalf("tool %s 的行 corr = %q 塌在 task id 上：这一发的作差不算读数", r.Tool, r.CorrelationID)
		}
		if !strings.HasPrefix(r.CorrelationID, res.TaskID) {
			t.Errorf("tool %s 的行 corr = %q, want 以 task id %q 为前缀", r.Tool, r.CorrelationID, res.TaskID)
		}
	}

	// 作差那一发（本用例的存在理由）：同一任务的两枚调用，corr 互不相等。
	if rows[0].CorrelationID == rows[1].CorrelationID {
		t.Errorf("同一任务的两枚调用共用一枚 corr %q（%s 与 %s）：两枚 ask 塌成同一个路由键",
			rows[0].CorrelationID, rows[0].Tool, rows[1].Tool)
	}

	// 各自可路由：两枚 corr 在路由表里各占一格，且各自归给自己那一枚调用。
	keyed := map[string]string{}
	for _, r := range rows {
		if prev, dup := keyed[r.CorrelationID]; dup {
			t.Errorf("corr %q 同时是 %s 与 %s 的路由键：名册里两枚调用不可分", r.CorrelationID, prev, r.Tool)
		}
		keyed[r.CorrelationID] = r.Tool
	}
	if len(keyed) != 2 {
		t.Errorf("路由键 = %d 枚，want 2（一枚调用一键）: %+v", len(keyed), keyed)
	}
	for corr, tool := range keyed {
		switch tool {
		case "task.spawn", "task.cancel":
			// 路由键带着它自己那一枚调用的标识：spawn 的键里必须有 spawn，
			// cancel 的键里必须有 cancel，两枚不可互换、也不可同名。
			want := "spawn"
			if tool == "task.cancel" {
				want = "cancel"
			}
			if !strings.Contains(corr, want) {
				t.Errorf("%s 的行 corr = %q, want 后缀带着 %q（各自可路由，不是同一本键）", tool, corr, want)
			}
		default:
			t.Errorf("意外的调用名 %q in row corr %q", tool, corr)
		}
	}
}
