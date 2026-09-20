package approval_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// The C18 queue and F2 layer 3: allow is only reachable from the native side,
// and "only" here means the untrusted path cannot mint the proof, not that it
// was asked politely to set a field honestly.

func TestL2NativeAllowExecutesThroughTheBridge(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui)
	g.AdmitTextTask(testTask)
	dir := t.TempDir()
	// Outside the allowlist: R2 raises the declared-L1 tool to L2, which is
	// the "declared level is a lower bound, never the conclusion" case.
	ft := newFakeTool("fs.write", tools.CapFSWrite, risk.L1, "paths")
	b := bridgeFor(t, g, ft, []string{dir})
	args := argsJSON(t, "C:/Windows/not-mine.txt")

	res := make(chan agent.ToolOutcome, 1)
	var gerr error
	spawn(t, func() {
		out, err := b.Execute(context.Background(), agentToolRequest("corr-l2", args))
		gerr = err
		res <- out
	})

	p := ui.wait(t)
	if p.Level != "L2" {
		t.Fatalf("card level=%q，越界写应为 L2", p.Level)
	}
	if p.Grant == "" || !strings.HasPrefix(p.Grant, "grant_") {
		t.Fatal("原生卡片必须拿到一枚签发令牌，否则用户根本没有可点击的允许")
	}
	if p.Deadline != approval.DefaultApprovalTimeout {
		t.Errorf("Deadline=%v，期望 300s", p.Deadline)
	}
	if p.Depth != 1 {
		t.Errorf("队头单显的深度角标应为 1，得 %d", p.Depth)
	}
	// C27: FULL params, not a summary, and the rules verbatim.
	raw, _ := json.Marshal(p.Params)
	if !strings.Contains(string(raw), "not-mine.txt") {
		t.Errorf("卡片必须携带完整参数：%s", raw)
	}
	if len(p.RulesHit) == 0 || !hasID(p.RulesHit, risk.R2) {
		t.Errorf("rules_hit 必须原样送达卡片，得 %v", p.RulesHit)
	}
	if p.Reason == "" {
		t.Error("判定理由必须原样送达卡片")
	}

	if err := g.Native().Allow(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("原生允许被拒：%v", err)
	}
	out := <-res
	clk.Advance(10 * time.Second)
	if gerr != nil {
		t.Fatalf("err=%v，审批通过不得成为宿主故障", gerr)
	}
	if out.IsError {
		t.Fatalf("批准后应执行：%+v", out)
	}
	if n := ft.count(); n != 1 {
		t.Errorf("工具执行 %d 次，期望 1", n)
	}
	if out.RiskLevel != "L2" {
		t.Errorf("risk_level=%q，期望 L2", out.RiskLevel)
	}
}

// The demanded forgery test. Three escalating attacks on the same endpoint:
// claim a source, carry a stolen grant, and use the real grant on the wrong
// route. All three are refused, and the item stays pending so nothing silently
// executed.
func TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui)
	g.AdmitTextTask(testTask)
	d := l2Decision("C:/data/a.txt")
	res := runApproval(t, g, context.Background(), d)
	p := ui.wait(t)
	if p.Grant == "" {
		t.Fatal("no grant minted")
	}

	cases := []struct {
		name string
		req  approval.Request
		want error
	}{
		{"面板声称自己是原生，且不带令牌", approval.Request{
			CorrelationID: p.CorrelationID, Allow: true, Source: "native"}, approval.ErrPanelAllow},
		{"面板声称原生并携带真实令牌", approval.Request{
			CorrelationID: p.CorrelationID, Allow: true, Source: "native", Grant: p.Grant}, approval.ErrPanelAllow},
		{"面板只带令牌不给答案（Allow=false 视为拒绝）", approval.Request{
			CorrelationID: p.CorrelationID, Allow: false, Source: "panel"}, nil},
	}
	// Case 3 answers the item, so run the two allow forgeries first.
	for _, tc := range cases[:2] {
		t.Run(tc.name, func(t *testing.T) {
			err := g.DecideFromPanel(context.Background(), tc.req)
			if !errors.Is(err, tc.want) {
				t.Fatalf("err=%v，期望 %v", err, tc.want)
			}
			select {
			case a := <-res:
				t.Fatalf("伪造的允许竟然结束了等待：%+v", a)
			default:
			}
		})
	}

	// A panel-side grant offer is treated as a compromise signal: the live
	// nonce is revoked, so even the real native card can no longer allow.
	if err := g.Native().Allow(context.Background(), p.CorrelationID, p.Grant); !errors.Is(err, approval.ErrBadGrant) {
		t.Fatalf("面板暴露过令牌后，原生侧仍可用它批准（err=%v）", err)
	}

	// The PanelAPI itself has no way to ask for an allow: there is no method
	// to lie to. Once the grant is burned, the honest outcome is rejection.
	if err := g.Panel().Reject(p.CorrelationID, "面板只能拒绝"); err != nil {
		t.Fatalf("Panel().Reject: %v", err)
	}
	a := mustAnswer(t, res)
	if a.a != tools.AnswerReject {
		t.Fatalf("answer=%q，期望 reject：%s", a.a, a.why)
	}
	clk.Advance(10 * time.Second)
}

// The grant is single-use and bound to one pending item; a caller-controlled
// Source string buys nothing in either direction.
func TestGrantIsSingleUseAndBoundToItsItem(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui)
	g.AdmitTextTask(testTask)

	d := l2Decision("C:/data/a.txt")
	d.CorrelationID = "corr-A"
	resA := runApproval(t, g, context.Background(), d)
	pA := ui.wait(t)

	// A valid grant for item B: unknown correlation id.
	if err := g.Native().Allow(context.Background(), "corr-B", pA.Grant); !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Fatalf("err=%v，期望 ErrUnknownCorrelation", err)
	}
	// A guessed grant for item A - the value is 256 bits from crypto/rand, so
	// this is the only attack an outsider has, and it is refused.
	if err := g.Native().Allow(context.Background(), pA.CorrelationID, "grant_00"); !errors.Is(err, approval.ErrBadGrant) {
		t.Fatalf("err=%v，期望 ErrBadGrant", err)
	}
	// The route is what carries authority, not the claim: a native decision
	// that describes itself as a panel still needs (and gets) the grant.
	if err := g.DecideFromNative(context.Background(), approval.Request{
		CorrelationID: pA.CorrelationID, Allow: true, Grant: pA.Grant, Source: "panel"}); err != nil {
		t.Fatalf("DecideFromNative with a real grant: %v", err)
	}
	if a := mustAnswer(t, resA); a.a != tools.AnswerAllow {
		t.Fatalf("answer=%q，期望 allow：%s", a.a, a.why)
	}
	// Replaying the spent nonce after the fact must fail even though the
	// attacker has a perfectly valid-looking string.
	if err := g.Native().Allow(context.Background(), pA.CorrelationID, pA.Grant); err == nil {
		t.Fatal("已消费的令牌竟然二次生效")
	}
	clk.Advance(10 * time.Second)
}

// C18: 超时 300s 一律判拒绝（never an infinite wait），超时前 30s 醒目提示,
// and a reject must not cancel the task's root context.
func TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{
		ApprovalTimeout: 300 * time.Second, WarningLead: 30 * time.Second})
	revoke := g.AdmitTextTask(testTask)
	defer revoke()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dir := t.TempDir()
	ft := newFakeTool("fs.write", tools.CapFSWrite, risk.L1, "paths")
	b := bridgeFor(t, g, ft, []string{dir})
	res := make(chan agent.ToolOutcome, 1)
	var gerr error
	spawn(t, func() {
		out, err := b.Execute(ctx, agentToolRequest("corr-slow",
			argsJSON(t, "C:/Windows/not-mine.txt")))
		gerr = err
		res <- out
	})
	ui.wait(t)

	clk.Advance(270*time.Second - time.Millisecond)
	if got := ui.ofKind(approval.EventWarning); len(got) != 0 {
		t.Fatalf("270s 之前不该弹出醒目提示：%+v", got)
	}
	clk.Advance(time.Millisecond)
	var got []approval.Event
	waitFor(t, func() bool {
		got = ui.ofKind(approval.EventWarning)
		return len(got) == 1
	}, "超时前 30s 必须有一次醒目提示")
	if got[0].Remaining != 30*time.Second || !strings.Contains(got[0].Text, "30 秒") {
		t.Errorf("提示要说清还剩 30s：%+v", got[0])
	}
	select {
	case a := <-res:
		t.Fatalf("提示阶段就结束了审批等待：%+v", a)
	default:
	}

	clk.Advance(30 * time.Second)
	out := <-res
	if gerr != nil {
		t.Fatalf("err=%v，自动拒绝不是宿主故障", gerr)
	}
	if !out.IsError || out.ErrorClass != string(observe.ClassUserRejected) {
		t.Fatalf("自动拒绝必须以可自纠的失败回给模型：%+v", out)
	}
	if n := ft.count(); n != 0 {
		t.Errorf("超时后绝不可执行，工具跑了 %d 次", n)
	}
	if ctx.Err() != nil {
		t.Fatalf("C18 明写拒绝后任务 root ctx 不取消，得 %v", ctx.Err())
	}

	// The task really does continue: a later approval on the same live context
	// is answered normally.
	d := l2Decision("C:/data/b.txt")
	d.CorrelationID = "corr-next"
	next := runApproval(t, g, ctx, d)
	p := ui.wait(t)
	if len(ui.ofKind(approval.EventWarning)) != 1 {
		t.Error("前一次超时的一次性提示不得泄漏到新条目")
	}
	if err := g.Native().Allow(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("后续审批被拒：%v", err)
	}
	if a := mustAnswer(t, next); a.a != tools.AnswerAllow {
		t.Fatalf("answer=%q，期望 allow：%s", a.a, a.why)
	}
	if ctx.Err() != nil {
		t.Errorf("任务上下文被审批层取消了：%v", ctx.Err())
	}
	clk.Advance(10 * time.Second)
}

// C18 一键重放: a refused request can be re-displayed, but re-displaying is
// not answering - the old nonce is worthless on the new item.
func TestReplayRedisplaysUnderAFreshGrant(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui)
	g.AdmitTextTask(testTask)
	d := l2Decision("C:/data/a.txt")
	res := runApproval(t, g, context.Background(), d)
	p := ui.wait(t)
	if err := g.Panel().Reject(p.CorrelationID, "先看一眼再说"); err != nil {
		t.Fatalf("Reject: %v", err)
	}
	if a := mustAnswer(t, res); a.a != tools.AnswerReject {
		t.Fatalf("answer=%q，期望 reject", a.a)
	}

	fresh, corr, err := g.Replay(context.Background(), p.CorrelationID, testTask)
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if fresh.Tool != d.Tool || string(fresh.Args) != string(d.Args) {
		t.Fatalf("重放必须带回原请求：%+v", fresh)
	}
	if corr == p.CorrelationID {
		t.Fatal("重放项必须换 correlation_id，否则答复会路由到已作废的旧卡片")
	}
	// The dead nonce cannot authorize the new item.
	if err := g.Native().Allow(context.Background(), corr, p.Grant); err == nil {
		t.Fatal("旧令牌竟然能批准重放后的新条目")
	}

	res2 := runApproval(t, g, context.Background(), tools.Decision{
		Tool: fresh.Tool, Provider: fresh.Provider, Params: fresh.Params, Args: fresh.Args,
		Level: risk.L2, RulesHit: fresh.RulesHit, Reason: fresh.Reason, Paths: fresh.Paths,
		CorrelationID: corr, TaskID: testTask, CallID: "call-replay",
	})
	p2 := ui.wait(t)
	if p2.Grant == p.Grant {
		t.Fatal("重放项复用了同一枚令牌（应重新签发）")
	}
	if err := g.Native().Allow(context.Background(), p2.CorrelationID, p2.Grant); err != nil {
		t.Fatalf("重放后原生允许失败：%v", err)
	}
	if a := mustAnswer(t, res2); a.a != tools.AnswerAllow {
		t.Fatalf("answer=%q，期望 allow：%s", a.a, a.why)
	}
	clk.Advance(10 * time.Second)
}

// Fail-closed in both directions a queue can break: the host cannot be
// reached, and the queue is full.
func TestHostUnreachableAndFullQueueFailClosed(t *testing.T) {
	t.Run("宿主不可达", func(t *testing.T) {
		ui := newFakeUI()
		ui.setErr(errors.New("WebView2 创建失败"))
		g, _, _ := newGate(t, ui)
		g.AdmitTextTask(testTask)
		a, why := g.PendingApproval(context.Background(), l2Decision("C:/data/a.txt"))
		if a != tools.AnswerReject || !strings.Contains(why, "fail-closed") {
			t.Fatalf("answer=%q why=%q，宿主不可达必须 fail-closed", a, why)
		}
		if n := g.Queue().Depth(); n != 0 {
			t.Errorf("失败的提示不得留下待审批项，深度=%d", n)
		}
	})

	t.Run("队列已满", func(t *testing.T) {
		ui := newFakeUI()
		g, clk, _ := newGate(t, ui, approval.Options{MaxPending: 1})
		g.AdmitTextTask(testTask)
		res := runApproval(t, g, context.Background(), l2Decision("C:/data/a.txt"))
		ui.wait(t)

		a, why := g.PendingApproval(context.Background(), l2Decision("C:/data/b.txt"))
		if a != tools.AnswerReject {
			t.Fatalf("answer=%q，队列满必须拒绝第二项：%s", a, why)
		}
		if !strings.Contains(why, "已满") {
			t.Errorf("文案要说清队列已满：%q", why)
		}
		clk.Advance(300 * time.Second)
		if got := mustAnswer(t, res); got.a != tools.AnswerTimeout {
			t.Fatalf("第一项应超时自动拒绝，得 %+v", got)
		}
	})
}

func hasID(rules []risk.RuleID, want risk.RuleID) bool {
	for _, r := range rules {
		if r == want {
			return true
		}
	}
	return false
}
