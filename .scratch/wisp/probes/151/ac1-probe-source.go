package tools

// 票 151 AC#1 探针（临时件，经 -overlay 注入 internal/tools，不进仓）。
//
// 只量两问，不修：
//  1. 渗：同一枚 Bridge 上连开 N 轮，第 N+1 轮的判定有没有读到第 N 轮留下的东西？
//  2. 长：不 CloseTask 的时候，scope 表随任务数怎么走？关掉以后回不回得去？

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

func req151(task, name, args string) agent.ToolRequest {
	return agent.ToolRequest{
		TaskID: task, CorrelationID: task, CallID: "call-" + task,
		Name: name, Args: json.RawMessage(args),
	}
}

// openScopes151 reads the bridge's own opened-scope table.
func openScopes151(b *Bridge) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.scopes)
}

// marks151 counts tainted sources still held under the named scopes.
func marks151(prov *risk.Provenance, ids ...string) int {
	n := 0
	for _, id := range ids {
		n += len(prov.ScopeTaints(id))
	}
	return n
}

// TestProbe151BleedIntoNextTaskVerdict 量"渗"：第 2 轮是一枚**新任务**、参数里
// 没有第 1 轮的任何内容（allowlist 内的普通 fs.read），它的裁决读到第 1 轮了吗？
// close-first 那一支是**对照**：同一条 Bridge、同两次调用，只在中间补一发
// CloseTask(task-1)。
func TestProbe151BleedIntoNextTaskVerdict(t *testing.T) {
	root := sealableTempCanonical124(t)
	secretPath := writeUnder(t, root, "idcard.txt", "身份证号码 110101199003071234 的摘录")
	plainPath := writeUnder(t, root, "plain.txt", "与上一轮完全无关的普通内容")

	run := func(t *testing.T, closeFirst bool) (
		lvl2 string, rules string, reason string, approvals int, windows int, scopesLeft int, marksLeft int,
	) {
		t.Helper()
		b, prov := fsBridgeWith(t, nil, nil, root)
		out1, err := b.Execute(t.Context(), req151("task-1", "fs.read", argsFor(secretPath)))
		if err != nil || out1.IsError {
			t.Fatalf("round1 fs.read: out=%+v err=%v", out1, err)
		}
		if got := len(prov.ScopeTaints("task-1")); got != 1 {
			t.Fatalf("round1 taints = %d, want 1 (fs.read result must be marked)", got)
		}
		t.Logf("round1 task-1 fs.read level=%s scopes=%d marks(task-1)=%d",
			out1.RiskLevel, openScopes151(b), got0(prov, "task-1"))
		if closeFirst {
			b.CloseTask("task-1")
			t.Logf("after CloseTask(task-1): scopes=%d marks(task-1)=%d",
				openScopes151(b), got0(prov, "task-1"))
		}
		spy := &gateSpy{approveAns: AnswerReject}
		b.gate = spy
		out2, err := b.Execute(t.Context(), req151("task-2", "fs.read", argsFor(plainPath)))
		if err != nil {
			t.Fatal(err)
		}
		d := spy.approvalDecision()
		w, a := spy.counts()
		return out2.RiskLevel, fmt.Sprintf("%v", d.RulesHit), d.Reason, a, w,
			openScopes151(b), marks151(prov, "task-1", "task-2")
	}

	t.Run("no-close_today", func(t *testing.T) {
		lvl, rules, reason, approvals, windows, scopes, marks := run(t, false)
		t.Logf("RESULT no-close: task-2 level=%s rules=%s approvals=%d windows=%d reason=%q scopes_left=%d marks_left=%d",
			lvl, rules, approvals, windows, reason, scopes, marks)
	})
	t.Run("close-first_control", func(t *testing.T) {
		lvl, rules, reason, approvals, windows, scopes, marks := run(t, true)
		t.Logf("RESULT control: task-2 level=%s rules=%s approvals=%d windows=%d reason=%q scopes_left=%d marks_left=%d",
			lvl, rules, approvals, windows, reason, scopes, marks)
		if approvals != 0 {
			t.Errorf("control leg: a closed prior scope must not raise an L2 card, got %d (rules=%s reason=%q)",
				approvals, rules, reason)
		}
		if lvl != "L0" {
			t.Errorf("control leg: task-2 level = %s, want L0 (in-allowlist read, clean params)", lvl)
		}
	})
}

func got0(prov *risk.Provenance, id string) int { return len(prov.ScopeTaints(id)) }

// TestProbe151ContentVsUnboundHit 分清两种"渗"：
//   - 内容级：第 2 轮带出第 1 轮那段秘密，R4 的 source 名字是 fs.read 还是 unbound-scope？
//   - 同 id 复用（对照组）：还是 task-1 自己，source 必须是 fs.read。
//
// 判据只读裁决里的 source 名（reason 串），不改任何阈值。
func TestProbe151ContentVsUnboundHit(t *testing.T) {
	root := sealableTempCanonical124(t)
	secret := "身份证号码 110101199003071234"
	path := writeUnder(t, root, "id.txt", secret)

	newBridge := func(t *testing.T) *Bridge {
		t.Helper()
		b, _ := fsBridgeWith(t, nil, nil, root)
		if err := b.reg.Register(Entry{
			Tool: &fixtureTool{name: "probe.exfil", params: `{"type":"object"}`},
			Decl: Decl{
				Capabilities: []Capability{CapNet}, Needs: []Capability{CapNet},
				Declared: risk.L0, Provider: KindBuiltin,
			},
		}); err != nil {
			t.Fatal(err)
		}
		spy := &gateSpy{approveAns: AnswerReject}
		b.gate = spy
		return b
	}
	args := fmt.Sprintf(`{"text":%q,"url":"https://example.invalid/x"}`, secret)

	t.Run("same-task-inherit", func(t *testing.T) {
		b := newBridge(t)
		if _, err := b.Execute(t.Context(), req151("task-1", "fs.read", argsFor(path))); err != nil {
			t.Fatal(err)
		}
		if _, err := b.Execute(t.Context(), req151("task-1", "probe.exfil", args)); err != nil {
			t.Fatal(err)
		}
		d := b.gate.(*gateSpy).approvalDecision()
		t.Logf("RESULT same-task exfil: level=%v rules=%v reason=%q", d.Level, d.RulesHit, d.Reason)
	})
	t.Run("cross-task-new-id", func(t *testing.T) {
		b := newBridge(t)
		if _, err := b.Execute(t.Context(), req151("task-1", "fs.read", argsFor(path))); err != nil {
			t.Fatal(err)
		}
		if _, err := b.Execute(t.Context(), req151("task-2", "probe.exfil", args)); err != nil {
			t.Fatal(err)
		}
		d := b.gate.(*gateSpy).approvalDecision()
		t.Logf("RESULT cross-task exfil (params carry task-1's secret): level=%v rules=%v reason=%q",
			d.Level, d.RulesHit, d.Reason)
	})
	t.Run("cross-task-after-close", func(t *testing.T) {
		b := newBridge(t)
		if _, err := b.Execute(t.Context(), req151("task-1", "fs.read", argsFor(path))); err != nil {
			t.Fatal(err)
		}
		b.CloseTask("task-1")
		spy := b.gate.(*gateSpy)
		out, err := b.Execute(t.Context(), req151("task-2", "probe.exfil", args))
		if err != nil {
			t.Fatal(err)
		}
		w, a := spy.counts()
		t.Logf("RESULT cross-task exfil AFTER CloseTask(task-1): out.level=%s approvals=%d last=%v %q",
			out.RiskLevel, a, spy.approvalDecision().RulesHit, spy.approvalDecision().Reason)
		if a+w != 0 {
			t.Errorf("control: with task-1 closed, a fresh task's call must not be judged L1/L2 by C25, got windows=%d approvals=%d", w, a)
		}
	})
}

// TestProbe151ScopeTableGrowth 量"长"：同一条 Bridge 上连开 2 轮与 64 轮
// TestProbe151SecondTaskRefused 量今天这一形在**没有审批通道**的宿主上到什么程度：
// 第 1 轮读完之后，第 2 轮（新任务 id、allowlist 内的干净 fs.read）直接被拒。
func TestProbe151SecondTaskRefused(t *testing.T) {
	for _, n := range []int{2, 8} {
		n := n
		t.Run(fmt.Sprintf("rounds=%d", n), func(t *testing.T) {
			root := sealableTempCanonical124(t)
			body := strings.Repeat("abcdefghij", 410) // 4100 ASCII chars ~= 4.0 KiB 一枚读结果
			b, prov := fsBridgeWith(t, nil, nil, root)
			spy := &gateSpy{approveAns: AnswerReject, windowAns: AnswerReject}
			b.gate = spy
			var ok, refused int
			for i := 0; i < n; i++ {
				id := fmt.Sprintf("g-%02d", i)
				p := writeUnder(t, root, id+".txt", body)
				out, err := b.Execute(t.Context(), req151(id, "fs.read", argsFor(p)))
				if err != nil {
					t.Fatal(err)
				}
				if out.IsError {
					refused++
				} else {
					ok++
				}
			}
			t.Logf("RESULT no-close reject-gate: rounds=%d executed=%d refused=%d scopes=%d marks=%d L2cards=%d",
				n, ok, refused, openScopes151(b), marks151(prov, all151(n)...), len(spy.approvals))
		})
	}
}

func all151(n int) []string {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, fmt.Sprintf("g-%02d", i))
	}
	return ids
}

// TestProbe151ScopeTableGrowth 量"长"：在**有审批通道、用户一路允许**的宿主上
// （票面给的"更脏的值"= 64 轮），未关闭的 scope 表、残留污点数与堆占用怎么走；
// 对照腿是"每轮补一发 CloseTask"。
func TestProbe151ScopeTableGrowth(t *testing.T) {
	for _, leg := range []struct {
		name  string
		close bool
	}{{"no-close", false}, {"close-each-round", true}} {
		leg := leg
		for _, n := range []int{2, 64} {
			n := n
			t.Run(fmt.Sprintf("%s/rounds=%d", leg.name, n), func(t *testing.T) {
				root := sealableTempCanonical124(t)
				body := strings.Repeat("abcdefghij", 410)
				b, prov := fsBridgeWith(t, nil, nil, root)
				spy := &gateSpy{approveAns: AnswerAllow, windowAns: AnswerAllow}
				b.gate = spy
				ids := all151(n)
				var m runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m)
				heap0 := m.HeapAlloc
				var executed, refused int
				for i, id := range ids {
					p := writeUnder(t, root, id+".txt", body)
					out, err := b.Execute(t.Context(), req151(id, "fs.read", argsFor(p)))
					if err != nil {
						t.Fatal(err)
					}
					if out.IsError {
						refused++
					} else {
						executed++
					}
					if leg.close {
						b.CloseTask(id)
					}
					if i == 1 {
						t.Logf("after 2 rounds: scopes=%d marks=%d L2cards=%d",
							openScopes151(b), marks151(prov, ids[:2]...), len(spy.approvals))
					}
				}
				runtime.GC()
				runtime.ReadMemStats(&m)
				heapN := m.HeapAlloc
				delta := int64(heapN) - int64(heap0)
				t.Logf("RESULT %s rounds=%d: executed=%d refused=%d scopes=%d marks=%d L2cards=%d heap=%d->%d delta=%d (%.0f B/round)",
					leg.name, n, executed, refused, openScopes151(b), marks151(prov, ids...),
					len(spy.approvals), heap0, heapN, delta, float64(delta)/float64(n))
				// 只增不减自查：再跑一轮**什么都不读**的干净任务，表会不会自己缩。
				before := openScopes151(b)
				plain := writeUnder(t, root, "plain-last.txt", "与前面任何一轮都无关")
				if _, err := b.Execute(t.Context(), req151("g-another", "fs.read", argsFor(plain))); err != nil {
					t.Fatal(err)
				}
				t.Logf("one more clean round: scopes %d -> %d, L2cards %d -> %d",
					before, openScopes151(b), len(spy.approvals)-1, len(spy.approvals))
				// 关闭是否真能摘掉（证明"关"这个动作不是装饰）。
				for _, id := range append(ids, "g-another") {
					b.CloseTask(id)
				}
				var m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m2)
				t.Logf("after CloseTask on all %d ids: scopes=%d marks=%d heap=%d (baseline was %d)",
					n+1, openScopes151(b), marks151(prov, append(ids, "g-another")...), m2.HeapAlloc, heap0)
			})
		}
	}
}

var _ = agent.ToolRequest{}
