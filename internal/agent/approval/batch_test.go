package approval_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// D45-1 batch aggregation and R7's floor under it.

func TestBatchAggregatesHomogeneousL1Ops(t *testing.T) {
	d := l1Decision("C:/data/x/a.txt", "C:/data/x/b.txt", "C:/data/y/c.txt",
		"C:/data/y/d.txt", "C:/data/y/e.txt", "C:/data/y/f.txt")
	bv := approval.Aggregate(d)
	if bv == nil {
		t.Fatal("6 项同质 L1 操作必须合并为一次确认（D45-1）")
	}
	if bv.Total != 6 {
		t.Errorf("Total=%d，期望 6", bv.Total)
	}
	if len(bv.First) != approval.BatchPreviewCount {
		t.Errorf("明细只显示前 %d 条，得 %d", approval.BatchPreviewCount, len(bv.First))
	}
	if bv.Deferred != 1 {
		t.Errorf("其余 %d 项应标为可展开，得 %d", 1, bv.Deferred)
	}
	if len(bv.AffectedDirs) != 2 || bv.AffectedDirs[0] != "C:/data/x" {
		t.Errorf("影响目录：%v", bv.AffectedDirs)
	}
	for _, want := range []string{"6 项", "2 个目录", "另有 1 项"} {
		if !strings.Contains(bv.Summary, want) {
			t.Errorf("摘要要含 %q：%q", want, bv.Summary)
		}
	}

	// One call, one confirm - not six.
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)
	res := runWindow(t, g, context.Background(), d)
	p := ui.wait(t)
	if p.Batch == nil || p.Batch.Total != 6 {
		t.Fatalf("L1 提示必须携带聚合摘要，得 %+v", p.Batch)
	}
	if p.Params["paths"] == nil {
		t.Error("聚合的是确认，不是参数：完整参数仍要随卡送达")
	}
	clk.Advance(3 * time.Second)
	if a := mustAnswer(t, res); a.a != tools.AnswerTimeout {
		t.Fatalf("answer=%q，期望 timeout：%s", a.a, a.why)
	}
	if n := ui.count(); n != 1 {
		t.Errorf("6 项操作只该弹一次确认，得 %d 次", n)
	}
}

func TestBatchRefusalsBelowThresholdAndOnSensitiveVerdicts(t *testing.T) {
	if bv := approval.Aggregate(l1Decision("C:/a.txt", "C:/b.txt")); bv != nil {
		t.Errorf("2 项不足 D45-1 的 >=3 门槛，不该聚合：%+v", bv)
	}
	if bv := approval.Aggregate(l2Decision("C:/a.txt", "C:/b.txt", "C:/c.txt")); bv != nil {
		t.Errorf("L2 永不聚合：%+v", bv)
	}
	tainted := l1Decision("C:/a.txt", "C:/b.txt", "C:/c.txt")
	tainted.RulesHit = []risk.RuleID{risk.R1, risk.R4}
	tainted.SessionOverrideBlocked = true
	if bv := approval.Aggregate(tainted); bv != nil {
		t.Errorf("C25/R4 污染命中必须逐项确认，摘要会藏掉来源：%+v", bv)
	}
	blocked := l1Decision("C:/a.txt", "C:/b.txt", "C:/c.txt")
	blocked.RulesHit = []risk.RuleID{risk.R1, risk.R3}
	if bv := approval.Aggregate(blocked); bv != nil {
		t.Errorf("R3 敏感路径不得被聚合一并放行：%+v", bv)
	}
}

// R7: a batch of >=50 targets is an L2 batch, and it is left unaggregated.
// The polarity flip is the point - an unanswered L2 must never execute.
func TestR7SizedBatchEscalatesToL2Unaggregated(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{
		Window: 3 * time.Second, ApprovalTimeout: 300 * time.Second, WarningLead: 30 * time.Second,
	})
	g.AdmitTextTask(testTask)
	res := runWindow(t, g, context.Background(), l1Decision(manyPaths(60)...))
	p := ui.wait(t)

	if p.Level != "L2" {
		t.Fatalf("card level=%q，R7 规模批次必须按 L2 走", p.Level)
	}
	if p.Batch != nil {
		t.Error("R7 批次不得聚合")
	}
	if !hasID(p.RulesHit, risk.R7) {
		t.Errorf("rules_hit 要带上 R7，卡片才说得出为什么升级：%v", p.RulesHit)
	}
	if p.Deadline != 300*time.Second {
		t.Errorf("L2 卡片要挂 C18 的 300s 截止，得 %v", p.Deadline)
	}
	// An allow really does let it run (it is still the same decision object).
	if err := g.Native().Allow(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("原生允许失败：%v", err)
	}
	if a := mustAnswer(t, res); a.a != tools.AnswerAllow {
		t.Fatalf("answer=%q，期望 allow：%s", a.a, a.why)
	}
	clk.Advance(10 * time.Second)
}

// C19 owns R7; the gate's own escalation is the backstop for a verdict that
// arrives as L1 with 60 targets on it.
func TestR7SizedBatchThatNobodyAnsweredIsRejectedNotRun(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{
		Window: 3 * time.Second, ApprovalTimeout: 300 * time.Second, WarningLead: 30 * time.Second,
	})
	g.AdmitTextTask(testTask)
	res := runWindow(t, g, context.Background(), l1Decision(manyPaths(60)...))
	ui.wait(t)
	clk.Advance(300 * time.Second)

	a := mustAnswer(t, res)
	if a.a == tools.AnswerTimeout || a.a == tools.AnswerAllow {
		t.Fatalf("升级后的批次超时竟然放行执行：%+v", a)
	}
	if a.a != tools.AnswerReject {
		t.Fatalf("answer=%q，期望 reject：%s", a.a, a.why)
	}
	clk.Advance(10 * time.Second)
}

// The whole thing at the choke point: ten L1 targets in ONE tool call is one
// confirm, and vetoing that confirm stops every one of them.
func TestTenOpsInOneToolCallGetOneConfirm(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)
	dir := t.TempDir()
	ft := newFakeTool("fs.write", tools.CapFSWrite, risk.L1, "paths")
	b := bridgeFor(t, g, ft, []string{dir})

	paths := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		paths = append(paths, filepath.Join(dir, "f"+itoa(i)+".txt"))
	}
	raw, err := json.Marshal(map[string]any{"paths": pathsAny(paths)})
	if err != nil {
		t.Fatal(err)
	}
	res := make(chan agent.ToolOutcome, 1)
	spawn(t, func() {
		out, _ := b.Execute(context.Background(), agentToolRequest("corr-batch", string(raw)))
		res <- out
	})

	p := ui.wait(t)
	if p.Level != "L1" {
		t.Fatalf("level=%q：授权目录内的声明 L1 写应留在 L1", p.Level)
	}
	if p.Batch == nil || p.Batch.Total != 10 {
		t.Fatalf("10 项目标要合并成一次确认，得 %+v", p.Batch)
	}
	if err := g.Veto(approval.Veto{CorrelationID: "corr-batch", Channel: approval.ChannelEsc}); err != nil {
		t.Fatalf("Veto: %v", err)
	}
	out := <-res
	clk.Advance(10 * time.Second)
	if runs := ft.count(); !out.IsError || runs != 0 {
		t.Fatalf("否决一次要取消整批：out=%+v 执行次数=%d", out, runs)
	}
	if n := ui.count(); n != 1 {
		t.Errorf("一次调用只该弹一次确认，得 %d 次", n)
	}
}
