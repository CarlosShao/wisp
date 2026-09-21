package approval_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// 票 97: 「别名买不到批准」今天靠习惯，不靠类型
//
// 票 87 的验收代理（docs/evidence/s1/87-adversarial-acceptance.md §11 R-1）留下
// 两条读数：队查找函数上那个方向 bool 的 true 侧零调用点，而它的注释描述的是一
// 条不存在的调用路径；并且全仓 *测试* 里钉"别名不能产生 allow"的断言是 0 条。
// 本文件把那 0 条补上（§6.2 攻击 1 的形状直接搬进来），并把
// R-2 的实际影响面（宽松解析坐在 Queue.reject 里 ⇒ 全部 5 条拒绝路线）写成一张
// 逐路线矩阵，而不是继续留在叙述里。
//
// 两条用例都不改判定结果，只钉方向：放行侧只认队列签发的精确键，拒绝侧可以
// 遍历这张卡的其它合法名字。
// ---------------------------------------------------------------------------

// ticket97Gate: an L2 deadline long enough that "refused now" and "timed out"
// cannot be confused, same instrument as ticket87Gate.
func ticket97Gate(t *testing.T, ui approval.UI) *approval.Gate {
	t.Helper()
	g, _, _ := newGate(t, ui, approval.Options{
		Clock:           approval.SystemClock{},
		ApprovalTimeout: 30 * time.Second,
		WarningLead:     25 * time.Second,
	})
	g.AdmitTextTask(testTask)
	return g
}

// aliasCard puts one L2 card on screen under an exact queue key that is NOT the
// alias, and returns the card plus the waiting call's answer channel. The alias
// is the task id the host keys its own bookkeeping on, i.e. a real
// qitem.names entry, so the refusal lookup does resolve it.
type aliasCard struct {
	p     approval.Prompt
	res   <-chan answer
	exact string
	alias string
}

func aliasCardOn(t *testing.T, ui *fakeUI, g *approval.Gate, path, exact, alias string) aliasCard {
	t.Helper()
	d := l2Decision(path)
	d.CorrelationID = exact
	d.TaskID = alias
	g.AdmitTextTask(alias)
	res := runApproval(t, g, context.Background(), d)
	p := ui.wait(t)
	if p.CorrelationID != exact {
		t.Fatalf("card key = %q, want the issued key %q", p.CorrelationID, exact)
	}
	if p.CorrelationID == alias {
		t.Fatalf("premise broke: the queue issued %q, so there is no alias to attack with", exact)
	}
	if p.Grant == "" {
		t.Fatal("premise broke: the card carries no live grant, an allow attack would be vacuous")
	}
	return aliasCard{p: p, res: res, exact: exact, alias: alias}
}

// TestAnAliasCanNeverBuyAnAllow is AC#2: the alias index is the reject-direction
// table, so no name from it may ever spend anything on the allow side. This is
// the shape the acceptance agent ran live in §6.2 attack 1 and measured as
// blocked; until this ticket nothing in the package said so.
func TestAnAliasCanNeverBuyAnAllow(t *testing.T) {
	ui := newFakeUI()
	g := ticket97Gate(t, ui)
	const exact, alias = "corr-97-exact-key", "task-97-host-key"
	card := aliasCardOn(t, ui, g, "C:/elsewhere/97-alias-never-buys-allow.txt", exact, alias)

	// Every allow route, addressed by the alias (and by shapes of it).
	names := []string{alias, alias + "  ", "  " + alias, alias + "-suffix", "no-such-card-97"}
	for _, name := range names {
		if err := g.Native().Allow(context.Background(), name, card.p.Grant); err == nil {
			t.Errorf("SECURITY: Native().Allow(%q, 卡片自己的活 grant) 返回 nil：别名买到了批准", name)
		} else if errors.Is(err, approval.ErrUnknownCorrelation) || errors.Is(err, approval.ErrNotPending) ||
			errors.Is(err, approval.ErrBadGrant) {
			// any of these is a refusal to authorize; the shape is asserted once below
		} else {
			t.Errorf("Allow(%q) err=%v, want a routing/grant error", name, err)
		}
		if err := g.DecideFromNative(context.Background(), approval.Request{
			CorrelationID: name, Allow: true, Grant: card.p.Grant, Source: "native",
		}); err == nil {
			t.Errorf("SECURITY: DecideFromNative(%q, allow=true, 真 grant) 返回 nil", name)
		}
		if err := g.DecideFromPanel(context.Background(), approval.Request{
			CorrelationID: name, Allow: true, Grant: card.p.Grant, Source: "native",
		}); err == nil {
			t.Errorf("SECURITY: DecideFromPanel(%q, allow=true, 真 grant) 返回 nil", name)
		}
	}

	// Nothing was consumed: the card is still the one live item and its call is
	// still waiting - an attack that quietly burned the entry is not a blocked
	// attack.
	if n := g.Queue().Depth(); n != 1 {
		t.Fatalf("queue depth = %d after %d allow attacks, want the card still pending", n, len(names))
	}
	select {
	case a := <-card.res:
		t.Fatalf("SECURITY: 别名攻击答掉了等待中的调用: %+v", a)
	default:
	}

	// The grant is still live, which is the proof that the refusals above were
	// about direction and not about a dead card.
	if err := g.Native().Allow(context.Background(), exact, card.p.Grant); err != nil {
		t.Fatalf("精确键的 Allow 失败: %v（前题破：卡片不可批准 ⇒ 本用例零信息）", err)
	}
	if got := mustAnswer(t, card.res); got.a != tools.AnswerAllow {
		t.Fatalf("answer = %v why=%q, want %v on the exact key", got.a, got.why, tools.AnswerAllow)
	}
	g.Complete(exact)
	waitFor(t, func() bool { return g.Queue().Depth() == 0 }, "the allowed card never left the queue")

	// Second card: the same alias buys exactly one thing - a refusal that lands.
	// Without this half the test would only prove the alias is meaningless
	// everywhere, which is not the direction the ticket is about.
	card2 := aliasCardOn(t, ui, g, "C:/elsewhere/97-alias-buys-refusal.txt", "corr-97-exact-key-2", alias+"-2")
	if err := g.Veto(approval.Veto{CorrelationID: card2.alias, Channel: approval.ChannelBall}); err != nil {
		t.Fatalf("Veto by the alias: %v（拒绝侧读不到这张表 ⇒ 下面的方向断言无从谈起）", err)
	}
	got := mustAnswerWithin(t, card2.res, vetoBudget, "别名驱动的拒绝")
	if got.a != tools.AnswerReject {
		t.Errorf("SECURITY: answer = %v why=%q, want %v：别名只准让拒绝更早落地",
			got.a, got.why, tools.AnswerReject)
	}
	g.Complete(card2.exact)
}

// TestEveryRefusalRouteOnAnUnknownEntryStillRefuses is AC#3: 票 87 §11 R-2 的
// 如实登记。宽松解析不在某条路线上，而在全部 5 条拒绝路线共用的 `Queue.reject`
// 漏斗里，所以它的影响面是这张表，不是票 87 叙述的"Veto 查不到时"。
// 每条路线断两件事：
//
//	(i) 查不到条目 ⇒ 诚实报错，且队列里那张活卡一格都没被动过（绝不是放行）；
//	(ii) 同一路线用卡片自己的别名点名 ⇒ 落地的是拒绝。
func TestEveryRefusalRouteOnAnUnknownEntryStillRefuses(t *testing.T) {
	type route struct {
		name string
		call func(t *testing.T, g *approval.Gate, corr string) error
	}
	routes := []route{
		{"veto", func(_ *testing.T, g *approval.Gate, corr string) error {
			return g.Veto(approval.Veto{CorrelationID: corr, Channel: approval.ChannelEsc})
		}},
		{"native_reject", func(_ *testing.T, g *approval.Gate, corr string) error {
			return g.Native().Reject(corr, "票 97 矩阵：原生拒绝路线")
		}},
		{"panel_reject", func(_ *testing.T, g *approval.Gate, corr string) error {
			return g.Panel().Reject(corr, "票 97 矩阵：面板拒绝路线")
		}},
		{"decide_from_native", func(_ *testing.T, g *approval.Gate, corr string) error {
			return g.DecideFromNative(context.Background(), approval.Request{
				CorrelationID: corr, Allow: false, Reason: "票 97 矩阵：原生路由", Source: "native",
			})
		}},
		{"decide_from_panel", func(_ *testing.T, g *approval.Gate, corr string) error {
			return g.DecideFromPanel(context.Background(), approval.Request{
				CorrelationID: corr, Allow: false, Reason: "票 97 矩阵：面板路由", Source: "panel",
			})
		}},
	}
	if len(routes) != 5 {
		t.Fatalf("R-2 登记的是 5 条拒绝路线，本矩阵只有 %d 条", len(routes))
	}

	for _, rt := range routes {
		rt := rt
		t.Run(rt.name, func(t *testing.T) {
			ui := newFakeUI()
			g := ticket97Gate(t, ui)
			const exact = "corr-97-matrix-exact"
			alias := "task-97-matrix-" + rt.name
			card := aliasCardOn(t, ui, g, "C:/elsewhere/97-matrix-"+rt.name+".txt", exact, alias)

			// (i) 查不到条目：诚实的 ErrUnknownCorrelation，绝不是一次成功。
			unknown := "corr-97-no-such-entry-" + rt.name
			err := rt.call(t, g, unknown)
			if err == nil {
				t.Errorf("SECURITY: 路线 %s 对查不到的条目返回 nil：拒绝路线把一次没有发生的拒绝当成了成功", rt.name)
			} else if !errors.Is(err, approval.ErrUnknownCorrelation) {
				t.Errorf("路线 %s unknown-entry err = %v, want %v", rt.name, err, approval.ErrUnknownCorrelation)
			}
			if n := g.Queue().Depth(); n != 1 {
				t.Fatalf("路线 %s 之后队列深度 = %d, want 1：查不到条目却动了那张活卡", rt.name, 2-n)
			}
			select {
			case a := <-card.res:
				t.Errorf("SECURITY: 路线 %s 用查不到的条目答掉了等待中的调用: %+v", rt.name, a)
			default:
			}
			if _, ok := g.Panel().View(exact); !ok {
				t.Fatalf("卡片在 %s 之后从面板视图消失，队列深度却是 1", rt.name)
			}

			// (ii) 同一路线的别名点名：影响面确实覆盖这 5 条，而落地方向是拒绝。
			if err := rt.call(t, g, alias); err != nil {
				t.Errorf("路线 %s 用卡片自己的别名点名: %v, want nil（R-2：宽松解析对全部 5 条路线生效）", rt.name, err)
			}
			got := mustAnswerWithin(t, card.res, vetoBudget, "路线 "+rt.name+" 的别名拒绝")
			if got.a != tools.AnswerReject {
				t.Errorf("SECURITY: 路线 %s 的答案 = %v why=%q, want %v：拒绝路线绝不产出放行",
					rt.name, got.a, got.why, tools.AnswerReject)
			}
			g.Complete(exact)
		})
	}
}
