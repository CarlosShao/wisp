package tools

// 票 165（plan.present）用例件。尺都在这里，生产件在 plan_present.go。
//
// 五格各有一条：
//
//	AC#1 Test165UnmodifiedRosterHasNoPlanTool    未修码名册里零枚 plan/propose
//	AC#2 Test165PlanPresentAsksThroughTheExistingL2Card
//	     （不新增状态/通道：走的还是 Gate 的两个现成回调、D43 的四行）
//	AC#3 Test165PlanApprovalRingsAgainForAnUnlistedWrite（清单外的写必须响）
//	     Test165ScopeClauseNamesTheBatch（"批过的范围"那一味按批大小写死）
//	AC#4 Test165PlanPresentAddsNoAllowExit（只加发起请求的出口）
//
// 两发变异（-overlay 注入、树零改）见 probes/165：把 Declared 从 L2 降成 L1
// 必须叫 AC#2 那条红；摘掉 planPresentScopeClause 必须叫 AC#3 的条款那条红。

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// plan165Gate records which of the two existing callbacks was asked. It answers
// with one fixed answer so a test can isolate "was the human asked" from "what
// did they say".
type plan165Gate struct {
	answer    Answer
	windows   []Decision
	approvals []Decision
}

func (g *plan165Gate) PendingWindow(_ context.Context, d Decision) (Answer, string) {
	g.windows = append(g.windows, d)
	return g.answer, "165 记录门：L1 静默窗口被问过一次"
}

func (g *plan165Gate) PendingApproval(_ context.Context, d Decision) (Answer, string) {
	g.approvals = append(g.approvals, d)
	return g.answer, "165 记录门：L2 卡被问过一次"
}

// plan165Ledger is the read side of a D45 session grant. Rows exist only when
// the TEST puts them there, which is the honest shape of "a native click minted
// a row earlier in this session": nothing internal/tools can reach writes here,
// so "the plan approval minted nothing" is observable as this staying empty.
type plan165Ledger struct {
	rows  []string
	asked int
}

func (l *plan165Ledger) Covering(_ context.Context, tool string, _ []string) (int64, bool) {
	l.asked++
	for i, r := range l.rows {
		if r == tool {
			return int64(i + 1), true
		}
	}
	return 0, false
}

// plan165Bridge composes the machine the way cmd/wisp/run.go does (fs family)
// plus the plan family this ticket lands, with the gate and the grant ledger
// under test control.
func plan165Bridge(t *testing.T, allowed string, g Gate, grants GrantSource) (*Bridge, *[]Decision) {
	t.Helper()
	paths := NewPathCanonicalizer([]string{allowed}, nil)
	reg := NewRegistry()
	roster := BuiltinFSEntries(FSDeps{Paths: paths, DeleteEnabled: true})
	roster = append(roster, BuiltinPlanEntries()...)
	for _, e := range roster {
		if err := reg.Register(e); err != nil {
			t.Fatalf("register %s: %v", e.Tool.Name(), err)
		}
	}
	var decs []Decision
	b := New(Options{
		Registry:   reg,
		Paths:      paths,
		Gate:       g,
		Grants:     grants,
		Logf:       func(string, ...any) {},
		OnDecision: func(d Decision) { decs = append(decs, d) },
	})
	return b, &decs
}

// plan165Args builds one plan.present argument object.
func plan165Args(t *testing.T, steps ...string) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(map[string]any{"summary": "165 探针", "steps": steps})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// plan165WriteArgs builds one fs.write argument object for a path under dir.
func plan165WriteArgs(t *testing.T, dir, name string) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(map[string]any{
		"path":    filepath.ToSlash(filepath.Join(dir, name)),
		"content": "165\n",
	})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

// Test165UnmodifiedRosterHasNoPlanTool is AC#1's reading in machine form: the
// roster that exists WITHOUT this ticket (fs + task + subagent builders) holds
// zero plan/propose-shaped tools, which is why the ticket stands. It reads the
// three existing builders only, so adding the plan family cannot make it pass
// by itself.
func Test165UnmodifiedRosterHasNoPlanTool(t *testing.T) {
	paths := NewPathCanonicalizer([]string{t.TempDir()}, nil)
	roster := BuiltinFSEntries(FSDeps{Paths: paths, DeleteEnabled: true})
	roster = append(roster, BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster(), Paths: paths})...)
	roster = append(roster, BuiltinSubagentEntries(SubagentDeps{})...)
	if len(roster) == 0 {
		t.Fatal("未修码名册是空的：这条读数就无从谈起了")
	}
	var names []string
	for _, e := range roster {
		names = append(names, e.Tool.Name())
		low := strings.ToLower(e.Tool.Name())
		if strings.Contains(low, "plan") || strings.Contains(low, "propose") {
			t.Errorf("未修码名册里出现了 plan 形的工具 %s：这张票的前提读数要重做", e.Tool.Name())
		}
	}
	t.Logf("未修码名册 %d 枚、零枚 plan/propose 形：%v", len(names), names)
}

// Test165PlanPresentAsksThroughTheExistingL2Card is AC#2: the plan档 reuses the
// frozen state machine instead of adding one. Three readings, one per claim:
//
//	the level the gate is asked at is L2 (D43 #17's AwaitingApproval branch);
//	it is the CARD that is asked, never the silent L1 window;
//	Gate still has exactly its two frozen callbacks - no third channel.
func Test165PlanPresentAsksThroughTheExistingL2Card(t *testing.T) {
	dir := t.TempDir()
	g := &plan165Gate{answer: AnswerAllow}
	b, decs := plan165Bridge(t, dir, g, nil)

	if n := reflect.TypeOf((*Gate)(nil)).Elem().NumMethod(); n != 2 {
		t.Fatalf("Gate 有 %d 个方法：本票复用现成的 PendingWindow/PendingApproval 两个回调，"+
			"多出来的第三个就是一张新审批通道", n)
	}
	entry, ok := b.reg.Lookup("plan.present")
	if !ok {
		t.Fatal("plan.present 没在名册里")
	}
	if entry.Decl.Declared != risk.L2 {
		t.Errorf("Declared=%v，冻结文本（PLAN.md:2567 / SPEC-07 §3:73）写的是 L2", entry.Decl.Declared)
	}
	if len(entry.Decl.Capabilities) != 0 || len(entry.Decl.Needs) != 0 || len(entry.Decl.PathParams) != 0 {
		t.Errorf("清单工具的能力/路径面不是空的：caps=%v needs=%v paths=%v",
			entry.Decl.Capabilities, entry.Decl.Needs, entry.Decl.PathParams)
	}

	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t165", CallID: "c165", Name: "plan.present",
		Args: plan165Args(t, "改 A 文件", "改 B 文件"),
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out.IsError {
		t.Fatalf("清单被拒了：%s", out.Text)
	}
	if len(g.windows) != 0 {
		t.Errorf("清单走了 L1 静默窗口 %d 次：这一档要的是卡（3 秒无人应答不是点头）", len(g.windows))
	}
	if len(g.approvals) != 1 {
		t.Fatalf("L2 卡被问了 %d 次，要 1 次", len(g.approvals))
	}
	if got := g.approvals[0].Level; got != risk.L2 {
		t.Errorf("卡上的等级是 %v，要 L2（D43 #17 的 AwaitingApproval 那一支）", got)
	}
	if len(*decs) != 1 {
		t.Fatalf("OnDecision 收到 %d 条判定，要 1 条", len(*decs))
	}
	// 这一发的 L2 是从工具自己的声明（R1）来的：档位照冻结行读，没有新规则、
	// 没有新模式，也没有被判成静默（静默＝没有卡，没有卡就没有点头）。
	dec := (*decs)[0]
	if !slices.Contains(dec.RulesHit, risk.R1) {
		t.Errorf("这一发的 L2 不是从工具声明（R1）来的：rules=%v", dec.RulesHit)
	}
	if dec.ModeSilenced {
		t.Errorf("清单这一发被判成了静默（mode=%s kept=%s）：没有卡就没有点头", dec.Mode, dec.ModeKept)
	}
	if !strings.Contains(out.Text, "1. 改 A 文件") || !strings.Contains(out.Text, "2. 改 B 文件") {
		t.Errorf("清单文本没有把两步列出来：%q", out.Text)
	}
}

// Test165PlanApprovalRingsAgainForAnUnlistedWrite is AC#3①: approving the plan
// is not a blanket cover. After the plan call went through its card, a write
// that the plan never listed must still raise its own question, and the plan
// approval must not have left a grant row behind.
func Test165PlanApprovalRingsAgainForAnUnlistedWrite(t *testing.T) {
	dir := sealableTempDir124(t)
	g := &plan165Gate{answer: AnswerAllow}
	led := &plan165Ledger{}
	b, decs := plan165Bridge(t, dir, g, led)

	if _, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t165", CallID: "c165a", Name: "plan.present",
		Args: plan165Args(t, "只改 in-plan.txt"),
	}); err != nil {
		t.Fatalf("plan.present: %v", err)
	}

	// 清单之外的一发写：名字都没在清单里出现过。
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t165", CallID: "c165b", Name: "fs.write",
		Args: plan165WriteArgs(t, dir, "not-in-plan.txt"),
	})
	if err != nil {
		t.Fatalf("fs.write: %v", err)
	}
	if out.IsError {
		t.Fatalf("清单外的写被拒了（要的是「响」＝被问，不是被拒）：%s", out.Text)
	}
	if len(g.approvals) != 1 {
		t.Fatalf("L2 卡被问了 %d 次，要 1 次（只有清单那一发）", len(g.approvals))
	}
	if len(g.windows) != 1 {
		t.Errorf("方案批过之后，清单外的写没有响（L1 窗口被问 %d 次，要 1 次）——"+
			"批准变成了对后续写的放行", len(g.windows))
	}
	if len(led.rows) != 0 {
		t.Errorf("批准一份清单之后，会话授权里多出了 %v：批准只批那一发，不许铸出授权行", led.rows)
	}
	last := (*decs)[len(*decs)-1]
	if last.DecisionColumn == agent.DecisionAllowGrant {
		t.Errorf("清单外的写靠一笔会话授权被静默放行（decision=%s）：这正是「方案被批准」变成「后续写都被批准」的那条路",
			last.DecisionColumn)
	}
}

// Test165ScopeClauseNamesTheBatch is AC#3②'s ruler over the production text:
// the sentence that carries the approved scope must be IN what the human reads
// and must name this batch's size. Remove it (or make it generic) and there is
// no longer a reading that says an out-of-scope write is out of scope.
func Test165ScopeClauseNamesTheBatch(t *testing.T) {
	res, err := planPresent{}.Execute(context.Background(), plan165Args(t, "甲", "乙", "丙"), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if res.IsError {
		t.Fatalf("三步骤清单被拒了：%s", res.Text)
	}
	for _, want := range []string{"只限上面这 3 步", "清单之外"} {
		if !strings.Contains(res.Text, want) {
			t.Errorf("方案文本里读不到批准范围的边界（缺 %q）：人类点头时看不到「批到哪儿为止」：%q",
				want, res.Text)
		}
	}

	// 空清单不是一份方案：没有这一步，"清单之外"就无从谈起。
	empty, err := planPresent{}.Execute(context.Background(), json.RawMessage(`{"steps":[]}`), nil)
	if err != nil {
		t.Fatalf("Execute(empty): %v", err)
	}
	if !empty.IsError {
		t.Errorf("空清单被当成了一份方案：%q", empty.Text)
	}
}

// Test165PlanPresentAddsNoAllowExit is AC#4: this ticket adds a way to ASK,
// never a way to ALLOW. Four readings:
//
//	no field on the declaration can carry an approval;
//	no parameter name the model can set means "approved"/"scope";
//	the grant read seam is still read-only (one method, no writer);
//	the tool reports that it applied nothing.
func Test165PlanPresentAddsNoAllowExit(t *testing.T) {
	entry := BuiltinPlanEntries()[0]

	typ := reflect.TypeOf(entry.Decl)
	for i := 0; i < typ.NumField(); i++ {
		low := strings.ToLower(typ.Field(i).Name)
		if strings.Contains(low, "allow") || strings.Contains(low, "approve") || strings.Contains(low, "grant") {
			t.Errorf("Decl 多了字段 %s：声明面能装下「已批准」就等于加了一个给出允许的出口", typ.Field(i).Name)
		}
	}

	var schema struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(entry.Tool.Parameters(), &schema); err != nil {
		t.Fatalf("参数 schema 读不出来：%v", err)
	}
	if len(schema.Properties) == 0 {
		t.Fatal("参数 schema 解析出零个属性：阳性对照失败，扫描无从谈起")
	}
	for k := range schema.Properties {
		low := strings.ToLower(k)
		for _, bad := range []string{"allow", "approve", "grant", "permission", "scope", "auto"} {
			if strings.Contains(low, bad) {
				t.Errorf("参数 %q 让模型自己说「这批已获准」：本票只许加发起请求的出口", k)
			}
		}
	}

	gs := reflect.TypeOf((*GrantSource)(nil)).Elem()
	if gs.NumMethod() != 1 || gs.Method(0).Name != "Covering" {
		t.Errorf("GrantSource 的方法面变了（%d 个，第一个 %s）：授权读取面必须保持只读",
			gs.NumMethod(), gs.Method(0).Name)
	}

	res, err := planPresent{}.Execute(context.Background(), plan165Args(t, "甲"), nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if res.AppliedSteps != nil || res.Truncated {
		t.Errorf("清单工具报告自己动了东西：AppliedSteps=%v Truncated=%v", res.AppliedSteps, res.Truncated)
	}
}
