package agent

// 票 285 AC#4 第五形：per-call corr 的"每枚调用各不相同"这件事的尺（落点甲）。
//
// 这一枚补的是 callCorr 的两支各自都要"作差"的洞：
// 带 call id 的那一支（provider 给了 id）与回落到轮内位置的那一支
// （provider 没给 id）今天都只有逐枚形状尺（非空／前缀／不等于 task id），
// 没有把同一任务的两枚互相作差。回落那一支尤其：名册行不带 call id，
// 两枚调用塌成同一个键时既有的形状尺逐枚仍然全绿。
//
// 落点甲的理由：callCorr 是小写未导出，形状只能在 package agent 内测；
// 而夹具给不出"没有 call id 的两枚并发调用"（现量：全仓 SSE 夹具都带 id），
// 所以回落支的形状尺必须是包内直接调用，真 loop 那一发只覆盖带 id 支。

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// recordCorrProvider answers every call immediately and keeps the ToolRequest
// the loop actually dispatched, so the request side is read without a gate, a
// database or an approval card.
type recordCorrProvider struct {
	reqs []ToolRequest
}

func (p *recordCorrProvider) Tools(context.Context) ([]ToolInfo, error) {
	return []ToolInfo{{
		Name: "echo", Description: "Echo back the given text.", Resident: true,
		RiskLevel: RiskL0, Parameters: json.RawMessage(`{"type":"object"}`),
	}}, nil
}

func (p *recordCorrProvider) Execute(_ context.Context, req ToolRequest) (ToolOutcome, error) {
	p.reqs = append(p.reqs, req)
	return ToolOutcome{Text: "answered"}, nil
}

// Test285CallCorrIsDistinctPerCall is the differential ruler on both branches
// of callCorr: two calls of ONE task must never land on the same routing key,
// whichever branch mints them.
func Test285CallCorrIsDistinctPerCall(t *testing.T) {
	const task = "task-285-corr"

	withCallID := []struct {
		who  string
		corr string
	}{
		{"first call", callCorr(task, "call-285-a", 0)},
		{"second call", callCorr(task, "call-285-b", 1)},
	}
	idless := []struct {
		who  string
		corr string
	}{
		{"first id-less call", callCorr(task, "", 0)},
		{"second id-less call", callCorr(task, "", 1)},
		{"third id-less call", callCorr(task, "", 2)},
	}

	for _, branch := range [][]struct {
		who  string
		corr string
	}{withCallID, idless} {
		seen := map[string]string{}
		for _, b := range branch {
			if b.corr == "" {
				t.Fatalf("%s minted an empty corr: C18 routes the reply by it", b.who)
			}
			if b.corr == task {
				t.Errorf("%s corr = %q, want a per-call id beside the task id %q: it collapses onto the task", b.who, b.corr, task)
			}
			if !strings.HasPrefix(b.corr, task) {
				t.Errorf("%s corr = %q, want the task id %q as prefix", b.who, b.corr, task)
			}
			if prev, dup := seen[b.corr]; dup {
				t.Errorf("two calls of one task share routing key %q (%s and %s): the two asks collapse onto one card",
					b.corr, prev, b.who)
			}
			seen[b.corr] = b.who
		}
		if len(seen) != len(branch) {
			t.Errorf("routing keys = %d, want %d distinct ones (one per call of the task)", len(seen), len(branch))
		}
		// 作差那一发：两支各自的两枚互不相等，逐枚形状尺看不见这一条。
		if branch[0].corr == branch[1].corr {
			t.Errorf("%s and %s carry the same corr %q", branch[0].who, branch[1].who, branch[0].corr)
		}
	}

	if withCallID[0].corr == idless[0].corr {
		t.Errorf("a call-id-bearing corr %q equals an id-less one", withCallID[0].corr)
	}
}

// Test285LoopDispatchesOneCorrPerCall puts the same differential on the
// request side of the real loop: two calls of one turn are dispatched with two
// different correlation ids, and each id names exactly its own call.
func Test285LoopDispatchesOneCorrPerCall(t *testing.T) {
	p := &recordCorrProvider{}
	h := newHarness(t, "parallel-tools", withTools(p))
	res := h.run("两件事一起办")

	if res.Status != StatusCompleted {
		t.Fatalf("status = %s (%s), want completed", res.Status, res.Message)
	}
	if len(p.reqs) != 2 {
		t.Fatalf("dispatched requests = %d, want 2 (one per tool call): %+v", len(p.reqs), p.reqs)
	}
	first, second := p.reqs[0], p.reqs[1]

	// 夹具前置（票面 AC#4 那句）：任务身份与调用身份在这里必须是两枚值，
	// 否则两发必同绿不算读数。这一条本身也是断言，不是注释。
	if first.TaskID != res.TaskID || second.TaskID != res.TaskID {
		t.Fatalf("request task ids = %q/%q, want the run's task id %q", first.TaskID, second.TaskID, res.TaskID)
	}
	if first.CorrelationID == res.TaskID || second.CorrelationID == res.TaskID {
		t.Fatalf("corr = %q/%q collapses onto the task id %q: the fixture no longer separates task from call",
			first.CorrelationID, second.CorrelationID, res.TaskID)
	}
	if first.CorrelationID == second.CorrelationID {
		t.Fatalf("both calls of one task dispatched with corr %q: the two asks share one routing key",
			first.CorrelationID)
	}
	for _, req := range p.reqs {
		if req.CorrelationID == "" {
			t.Errorf("call %s carried an empty corr", req.CallID)
			continue
		}
		if !strings.HasPrefix(req.CorrelationID, res.TaskID) {
			t.Errorf("call %s corr = %q, want the task id %q as prefix", req.CallID, req.CorrelationID, res.TaskID)
		}
		// 各自可路由：corr 去掉任务前缀剩下的就是它自己那一枚调用的 id。
		if suffix := strings.TrimPrefix(req.CorrelationID, res.TaskID+"#"); suffix != req.CallID {
			t.Errorf("call %s corr = %q routes to %q, want its own call id", req.CallID, req.CorrelationID, suffix)
		}
	}
	if first.CallID == second.CallID {
		t.Fatalf("fixture dispatched the same call id twice: %+v/%+v", first, second)
	}
}
