package tools

// 票 158 AC#1 —— "过桥的敏感读会打开它自己的 C25 污点 scope" 这一味在本包必须有断言。
//
// 来路：票 154 非实现者验收表 D-6 记的正控 P-2 —— 摘掉 bridge.go 里 mark() 的那一发
// b.OpenTask(dec.TaskID)，本包 115/79/0/0 零枚红；同一发变异换到 cmd/wisp 才红
// （cmd/wisp/task_scope_close_151_test.go）。也就是说这件事唯一的守卫长在**另一个包**里，
// 下一个动 mark() 的程在本包改坏了是听不到声音的。
//
// 三枚判据各自守什么：
//  1. 敏感读跑完后，桥自己的开账本 b.scopes 必须记着这一枚 task id —— 这是被复现成
//     读数的那枚变异（摘掉 OpenTask）唯一直接改变的包内状态。
//  2. 收尾那一条 CloseTask 审计行必须报 was_open=true / dropped=1 —— 判据 1 的**后果**，
//     它证明这本账真的连到了"谁把这枚 task 的污点从表上摘下去"那一句上，而不是一个
//     只有测试会读的字段。
//  3. 反向对照：非敏感源（fs.list）不得开账。这一枚挡住一种"看起来也满足了判据 1"的
//     退化改法——在 Execute 里无条件给每个 task 开 scope。
//
// 为什么不去断言 prov 那一侧：risk.Provenance.Mark 走 p.scopes[id]=append(...)，一枚
// 从没 OpenScope 的 scope 被 Mark 之后同样会在表上出现并带着污点（provenance.go:394，
// 注释里明写"marking a closed/never-opened scope still records the taint"）。所以
// "Mark 之后 ScopeTaints 非空"这把尺**量不出 OpenTask 在不在**——本包 fs_test.go:35 正是
// 这一族，它对这枚变异是瞎的。真正瞎得少的只有桥自己的那本开账本。

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
)

// reqTask158 is one host-shaped request carrying its own task id. The shared
// req() helper hardcodes "task-1", and this test needs to tell two tasks apart.
func reqTask158(task, name, args string) agent.ToolRequest {
	return agent.ToolRequest{
		TaskID: task, CorrelationID: task, CallID: "call-" + task,
		Name: name, Args: json.RawMessage(args),
	}
}

// logSink158 captures the bridge's audit lines. In-package field assignment is
// the shape this package already uses for its own seams (fs_test.go:72).
type logSink158 struct {
	mu    sync.Mutex
	lines []string
}

func (s *logSink158) write(format string, args ...any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lines = append(s.lines, fmt.Sprintf(format, args...))
}

func (s *logSink158) containing(sub string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []string
	for _, l := range s.lines {
		if strings.Contains(l, sub) {
			out = append(out, l)
		}
	}
	return out
}

// scopeOpened158 reads the bridge's own open-scope ledger under its lock.
// Since ticket 160 the ledger holds the *risk.Scope handles (the only things
// that can close those scopes) rather than a bool, so "opened" is "a handle is
// registered" - the same predicate, no weaker assertion.
func scopeOpened158(b *Bridge, taskID string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.scopes[taskID] != nil
}

func TestSensitiveReadAcrossTheBridgeOpensItsTaskScope(t *testing.T) {
	root := sealableTempCanonical124(t)
	secret := "本机凭据摘录：" + strings.Repeat("abcdefgh", 4)
	path := writeUnder(t, root, "note.txt", secret)

	b, _ := fsBridgeWith(t, nil, nil, root)
	sink := &logSink158{}
	b.logf = sink.write

	// 判据 3（反向对照，先在一颗干净的污点表上跑）：非敏感源不得开账。
	// 这一发必须在敏感读**之前**：表上只要挂着别枚 task 的污点，一枚没开账的 task id
	// 就会被 risk 层的 unbound-scope fail-closed 顶到 L2（实测量到，见证据件 §1.5），
	// 那样这一腿量的就变成"审批通道"而不是"开账"了。
	const plainTask = "task-158-plain"
	out2, err2 := b.Execute(t.Context(), reqTask158(plainTask, "fs.list", argsFor(root)))
	if err2 != nil || out2.IsError {
		t.Fatalf("对照腿 fs.list 没有跑成：out=%+v err=%v", out2, err2)
	}
	if scopeOpened158(b, plainTask) {
		t.Errorf("非敏感源 fs.list 也开了 scope，task=%q："+
			"开账必须挂在敏感源上，不能挂在每一次 Execute 上", plainTask)
	}

	// 前置：这一发读必须真的跑成功，否则 mark() 根本不会被走到
	// （Execute 里 !res.IsError 那道门），红就会是假红。
	const taintedTask = "task-158-taint"
	out, err := b.Execute(t.Context(), reqTask158(taintedTask, "fs.read", argsFor(path)))
	if err != nil || out.IsError {
		t.Fatalf("敏感读没有跑成，本条判据失去前提：out=%+v err=%v", out, err)
	}
	if out.RiskLevel != "L0" {
		t.Fatalf("敏感读 judged %s, want L0（allowlist 内的 fs.read）", out.RiskLevel)
	}

	// 判据 1：桥的开账本记着这一枚 task。
	if !scopeOpened158(b, taintedTask) {
		t.Errorf("一枚过了桥的敏感读没有打开它自己的 C25 污点 scope："+
			"bridge.scopes 里没有 %q。这一味（mark 里的 b.OpenTask）是本包对 C25 开侧唯一的"+
			"断言点，它一旦被改掉，敏感读的污点就没有主人了", taintedTask)
	}

	// 判据 2：收尾那一行报的是"确实开着、摘掉了 1 枚"。
	b.CloseTask(taintedTask)
	want := "task=" + taintedTask + " was_open=true"
	got := sink.containing("C25 scope closed")
	var hit []string
	for _, l := range got {
		if strings.Contains(l, want) {
			hit = append(hit, l)
		}
	}
	if len(hit) == 0 {
		t.Fatalf("收尾审计行没有承认这枚 scope 是开过的：找不到含 %q 的行；"+
			"C25 scope closed 全部读数=%q\n"+
			"（was_open=false 就是开侧那一味丢了的形状：CloseTask 会照常打点、"+
			"照常把 risk 层的污点留在表上）", want, got)
	}
	if !strings.Contains(strings.Join(hit, "\n"), "dropped=1") {
		t.Errorf("收尾审计行没有带着 dropped=1（本轮只读了 1 枚敏感源），行=%q", hit)
	}
}
