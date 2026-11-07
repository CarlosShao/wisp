package agent

// 票 264 形甲的判据（咽喉一处：toolResultParts，A-账裁定）。
//
// 三问（AC#1 的两向 + 174 轴的保命钉）：
//
//	a  模型可见副本：一发带机器绝对路径的工具失败句进历史/上线后，
// 	   模型那边再也看不到盘符目录结构（<path> 顶替），但错误说明文字
// 	   （"： 系统找不到指定的文件。"）必须活下来——掩码不许吞正文。
// 	b  用户侧不掩：EvToolEnd 事件与 res.ToolLog 仍逐字带真路径，
// 	   用户仍然说得出是哪个文件（证据表 264-ac0 §5 的天然成立项）。
// 	c  174 指针豁免：task.output / spill 的回执「全文见 <绝对路径>」
// 	   必须逐字保留（那是模型回到全文的唯一路，掩掉＝只截不指，
// 	   PLAN.md:2564；cmd/wisp 的 174r8 钉靠它活着）。
//
// 尺的形状：e2e 走 harness 真 golden 流（C5 缝），断言打在历史、上线正文、
// 事件与 ToolLog 四面上；规则表与咽喉单元钉把掩码规则本身两面钉死
// （有路径必响、无路径逐字节不动——AC#2 的两向）。

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

// mask264RealPath 是这条尺自己的机器路径样本（固定串，不依赖 TempDir 拼写，
// 判据因此对盘符与用户名无关）。mask264FailText 是 fs.go:148 同形的失败句。
const (
	mask264RealPath = `C:\Users\operator\笔记\draft.txt`
	mask264FailText = "打开失败：open " + mask264RealPath + ": 系统找不到指定的文件。"
)

// pathMask264Provider 是 host 桥替身：一枚把机器绝对路径写进失败句的工具。
type pathMask264Provider struct {
	mu    sync.Mutex
	calls int
	text  string
}

func (p *pathMask264Provider) Tools(context.Context) ([]ToolInfo, error) {
	return []ToolInfo{{
		Name:        "echo",
		Description: "ticket 264 stand-in: 失败句里带机器绝对路径的工具",
		Resident:    true,
		RiskLevel:   RiskL0,
		Parameters:  json.RawMessage(`{"type":"object","properties":{}}`),
	}}, nil
}

func (p *pathMask264Provider) Execute(_ context.Context, _ ToolRequest) (ToolOutcome, error) {
	p.mu.Lock()
	p.calls++
	p.mu.Unlock()
	return ToolOutcome{Text: p.text, IsError: true}, nil
}

func (p *pathMask264Provider) count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.calls
}

// toolResultText264 收集历史里 RoleTool 消息的模型可见正文。
func toolResultText264(t *testing.T, hist []llm.Message) string {
	t.Helper()
	var out string
	for _, m := range hist {
		if m.Role != llmRoleTool {
			continue
		}
		for _, part := range m.Content {
			tr, ok := part.(llmToolResult)
			if !ok {
				continue
			}
			for _, inner := range tr.Content {
				if tp, ok := inner.(llmText); ok {
					out += tp.Text
				}
			}
		}
	}
	if out == "" {
		t.Fatal("history carries no tool result text: nothing to judge")
	}
	return out
}

// TestToolResultPathMaskedForModelButNotForUser264 是 (a)+(b) 的 e2e：
// 同一发失败句，模型副本无路径、用户副本逐字带路径。
func TestToolResultPathMaskedForModelButNotForUser264(t *testing.T) {
	p := &pathMask264Provider{text: mask264FailText}
	h := newHarness(t, "tool-then-text", withTools(p))

	res := h.run("读一下那份笔记")
	if res.Status != StatusCompleted {
		t.Fatalf("fixture did not run: status=%s msg=%s err=%v", res.Status, res.Message, res.Err)
	}
	if p.count() != 1 {
		t.Fatalf("tool dispatches = %d, want the one this fixture makes", p.count())
	}

	// (a) 模型可见副本：历史与上线正文都不得再带机器目录结构。
	modelText := toolResultText264(t, h.loop.History())
	if strings.Contains(modelText, mask264RealPath) || strings.Contains(modelText, `C:\Users`) {
		t.Errorf("model-visible copy carries the machine path: %q", modelText)
	}
	if !strings.Contains(modelText, maskPathPlaceholder) {
		t.Errorf("masked copy lost the placeholder: %q", modelText)
	}
	if !strings.Contains(modelText, "系统找不到指定的文件。") {
		t.Errorf("the mask ate the error prose: %q", modelText)
	}
	bodies := h.requestBodies()
	if len(bodies) < 2 {
		t.Fatalf("provider requests = %d, want 2 (tool turn + after-tool turn)", len(bodies))
	}
	body := string(bodies[1])
	if i := strings.Index(body, `C:\\Users`); i >= 0 {
		t.Errorf("the wire body carries the machine path at %d: …%s…", i,
			trunc(body[max(0, i-120):min(len(body), i+240)], 360))
	}
	if strings.Contains(body, "operator") {
		t.Errorf("the wire body carries a path fragment: %s", trunc(body, 800))
	}
	// The openai-chat adapter HTML-escapes < and > on the wire, so the
	// placeholder arrives as \u003cpath\u003e in the request body.
	if !strings.Contains(body, `\u003cpath\u003e`) {
		t.Errorf("the wire body lost the placeholder: %s", trunc(body, 800))
	}

	// (b) 用户侧：事件与 ToolLog 仍逐字说得出是哪个文件。
	ev := h.sink.LastOf(EvToolEnd)
	if ev == nil {
		t.Fatal("no EvToolEnd event published")
	}
	if !strings.Contains(ev.Text, mask264RealPath) {
		t.Errorf("user-facing EvToolEnd lost the real path: %q", ev.Text)
	}
	if len(res.ToolLog) != 1 {
		t.Fatalf("res.ToolLog = %d entries, want 1", len(res.ToolLog))
	}
	if !strings.Contains(res.ToolLog[0].Text, mask264RealPath) {
		t.Errorf("user-facing res.ToolLog lost the real path: %q", res.ToolLog[0].Text)
	}
}

// TestToolResultPartsMasksModelCopyOnly264 钉咽喉本身：toolResultParts 的
// 返回值已掩，调用方手里的字符串原封不动（用户侧上游不受影响）。
func TestToolResultPartsMasksModelCopyOnly264(t *testing.T) {
	parts := toolResultParts(&Loop{}, mask264FailText, nil)
	tp, ok := parts[0].(llm.TextPart)
	if !ok {
		t.Fatalf("parts[0] = %T, want llm.TextPart", parts[0])
	}
	if strings.Contains(tp.Text, mask264RealPath) {
		t.Errorf("the chokepoint let the machine path through: %q", tp.Text)
	}
	if !strings.Contains(tp.Text, maskPathPlaceholder) {
		t.Errorf("the chokepoint lost the placeholder: %q", tp.Text)
	}
	if mask264FailText != "打开失败：open "+mask264RealPath+": 系统找不到指定的文件。" {
		t.Errorf("the caller's own string was mutated: %q", mask264FailText)
	}
}

// TestToolResultPartsLeavesCleanText264 是 AC#2 的反向：一句不带任何绝对
// 路径的工具结果必须逐字节原样进历史——掩码不许顺手改写干净文本。
func TestToolResultPartsLeavesCleanText264(t *testing.T) {
	clean := "任务已取消，调用被中止"
	parts := toolResultParts(&Loop{}, clean, nil)
	if got := parts[0].(llm.TextPart).Text; got != clean {
		t.Errorf("clean text was rewritten: want %q got %q", clean, got)
	}
}

// TestMaskOutboundPathsRule264 是掩码规则的规则表：39 枚带路径装配点的真实
// 消息形逐形入表，另带不许误伤的形状（URL、时刻、裸盘符）。
func TestMaskOutboundPathsRule264(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{
			"patherror", mask264FailText,
			"打开失败：open " + maskPathPlaceholder + ": 系统找不到指定的文件。",
		},
		{
			"c26-expands", "路径无法解析（按 fail-closed 拒绝）：open D:\\a: D:\\a expands to D:\\a (不可达)",
			"路径无法解析（按 fail-closed 拒绝）：open " + maskPathPlaceholder +
				": " + maskPathPlaceholder + " expands to " + maskPathPlaceholder + " (不可达)",
		},
		{
			"rename-two-paths", "原子重命名失败：rename D:\\tmp\\x D:\\dst\\y: 系统找不到指定的文件。",
			"原子重命名失败：rename " + maskPathPlaceholder + ": 系统找不到指定的文件。",
		},
		{
			"space-in-middle-segment", "打开失败：open C:\\Users\\a b\\file.txt: 系统找不到指定的文件。",
			"打开失败：open " + maskPathPlaceholder + ": 系统找不到指定的文件。",
		},
		{
			"prose-after-path-kept", "路径 D:\\a\\b 不可读：详情如下",
			"路径 " + maskPathPlaceholder + " 不可读：详情如下",
		},
		{
			"fullwidth-colon-prefix", "拒绝（fail-closed）：D:\\secret\\file 不可读",
			"拒绝（fail-closed）：" + maskPathPlaceholder + " 不可读",
		},
		{
			"url-untouched", "打开 https://example.com/x?y=1 失败",
			"打开 https://example.com/x?y=1 失败",
		},
		{"clock-untouched", "重试于 12:30/45 之后", "重试于 12:30/45 之后"},
		{"bare-drive-untouched", "盘符 C: 已满", "盘符 C: 已满"},
		{"idempotent", "<path>: 已掩过一轮", "<path>: 已掩过一轮"},
	}
	for _, tc := range cases {
		if got := maskOutboundPaths(tc.in); got != tc.want {
			t.Errorf("%s: mask rule mismatch\n in  %q\nwant %q\n got %q", tc.name, tc.in, tc.want, got)
		}
	}
}

// TestMaskKeeps174TaskOutputPointer264 是 (c)：task.go:629 与 spill.go:189
// 同形的回执——注意句在前、指针在后——必须逐字节活着。掩掉它＝只截不指，
// 174r8 那两枚端到端钉（cmd/wisp）会先红。
func TestMaskKeeps174TaskOutputPointer264(t *testing.T) {
	artifact := `C:\Users\operator\AppData\Roaming\wisp\artifacts\tool-output-call-174r8.txt`
	notice := "；注意：这条路径现在读不到，它不在你被授权的目录范围内：" +
		"fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——" +
		"要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡"
	receipt := fmt.Sprintf(
		"%s\n[…输出已落文件：省略 3 字符，总长 20000 字节 / 约 5000 token%s，全文见 %s…]\n%s",
		"HEAD 段", notice, artifact, "TAIL 段")
	if got := maskOutboundPaths(receipt); got != receipt {
		t.Errorf("the 174 pointer was masked (只截不指):\nwant %q\n got %q", receipt, got)
	}
	if !strings.Contains(receipt, "全文见 "+artifact) {
		t.Errorf("the receipt lost its pointer: %q", receipt)
	}
}
