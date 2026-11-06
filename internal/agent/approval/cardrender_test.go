package approval_test

// Ticket 21 segment 2: the native card renderer. Every case below names the
// contract line it pins (SPEC-06 §7 item content, §8.1 batch shape, the
// ticket's "FULL params (not summary)" clause) rather than a golden blob, so a
// refactor that drops a disclosure cannot pass by rewriting a fixture.

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
)

func cardChannels() []approval.ChannelStatus {
	return approval.NewChannels(approval.ChannelBall, approval.ChannelEsc).Statuses()
}

// l2Card is the subject the ticket's AC#2 names: an L2 fs.write with full
// params and two rules, plus the depth the queue stamped.
func l2Card() approval.Prompt {
	return approval.Prompt{
		CorrelationID: "corr-21s2",
		TaskID:        "task-21s2",
		Tool:          "fs.write",
		Level:         "L2",
		Args:          json.RawMessage(`{"content":"x","path":"C:/work/a.txt"}`),
		Params:        map[string]any{"content": "x", "path": "C:/work/a.txt"},
		RulesHit:      []risk.RuleID{risk.R2, risk.R8},
		Reason:        "覆盖已有文件",
		Paths:         []string{"C:/work/a.txt"},
		Deadline:      300 * time.Second,
		Channels:      cardChannels(),
		Depth:         2,
		Grant:         "grant_live",
	}
}

func joinCard(lines []string) string { return strings.Join(lines, "\n") }

// TestTicket21S2CardCarriesRulesTaskAndDepth walks SPEC-06 §7's item-content
// list onto the rendered card: rules verbatim, the owning task, the queue
// depth («深度角标», both faces), and the C18 window the route really waits.
func TestTicket21S2CardCarriesRulesTaskAndDepth(t *testing.T) {
	text := joinCard(approval.RenderCard(l2Card()))
	for _, want := range []string{
		"[确认 L2 fs.write] 覆盖已有文件",
		"  规则 R2",
		"  规则 R8",
		"  所属任务：task-21s2",
		"  队列深度：2",
		"  影响路径：C:/work/a.txt",
		"  窗口 300.0s",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("card missing %q\nfull card:\n%s", want, text)
		}
	}
	if iR2, iR8 := strings.Index(text, "规则 R2"), strings.Index(text, "规则 R8"); iR2 > iR8 {
		t.Errorf("rules_hit order changed (R2 at %d, R8 at %d); the verdict list arrives ordered", iR2, iR8)
	}
}

// TestTicket21S2CardPrintsFullParamsNotSummary is AC#2's first clause. The
// check is per-value, not per-line-count: a rendering that folded the object
// into "2 个参数" would satisfy a length check and still fail the ticket.
func TestTicket21S2CardPrintsFullParamsNotSummary(t *testing.T) {
	p := l2Card()
	p.Args = json.RawMessage(`{"content":"a very long body","path":"C:/work/long/name.txt","mode":"overwrite"}`)
	text := joinCard(approval.RenderCard(p))
	for _, want := range []string{
		"完整参数（非摘要）：",
		`"content": "a very long body"`,
		`"path": "C:/work/long/name.txt"`,
		`"mode": "overwrite"`,
	} {
		if !strings.Contains(text, want) {
			t.Errorf("card missing %q\nfull card:\n%s", want, text)
		}
	}
}

// TestTicket21S2CardFallsBackToParamsMap covers the other producer shape: a
// Decision whose Args is empty still carries the parsed Params, and the card
// must not blank out for it. json.Marshal orders map keys, so the render is
// deterministic (two runs byte-equal).
func TestTicket21S2CardFallsBackToParamsMap(t *testing.T) {
	p := l2Card()
	p.Args = nil
	first := joinCard(approval.RenderCard(p))
	second := joinCard(approval.RenderCard(p))
	if first != second {
		t.Errorf("params render is not deterministic:\n%s\n---\n%s", first, second)
	}
	for _, want := range []string{`"content": "x"`, `"path": "C:/work/a.txt"`} {
		if !strings.Contains(first, want) {
			t.Errorf("card missing %q\nfull card:\n%s", want, first)
		}
	}
}

// TestTicket21S2CardSaysOutLoudWhenNoParamsArrives is the 宁缺毋造 boundary:
// with neither field carried, the card names its own gap instead of printing
// an empty block a reader could misread as "the call had no arguments".
func TestTicket21S2CardSaysOutLoudWhenNoParamsArrives(t *testing.T) {
	p := l2Card()
	p.Args = nil
	p.Params = nil
	text := joinCard(approval.RenderCard(p))
	if !strings.Contains(text, "（本次卡片未随附参数对象）") {
		t.Errorf("silent empty params block; want the explicit gap sentence\n%s", text)
	}
}

// TestTicket21S2CardKeepsUnparseableArgsVerbatim covers the hand-built Prompt
// shape (a test or a future host that carries bytes the bridge never validated):
// the field is printed as it arrived, never dropped because formatting failed.
func TestTicket21S2CardKeepsUnparseableArgsVerbatim(t *testing.T) {
	p := l2Card()
	p.Args = json.RawMessage(`{"broken":`)
	text := joinCard(approval.RenderCard(p))
	if !strings.Contains(text, `{"broken":`) {
		t.Errorf("unparseable args were dropped instead of printed verbatim\n%s", text)
	}
}

// TestTicket21S2CardBatchViewIsD45Shape pins §8.1's four pieces: total,
// affected-dir summary, the first 5 detail lines, and the expandable
// remainder. The view is built by Aggregate so the pair is exercised as the
// production route composes it.
func TestTicket21S2CardBatchViewIsD45Shape(t *testing.T) {
	d := l1Decision(manyPaths(10)...)
	bv := approval.Aggregate(d)
	if bv == nil {
		t.Fatal("premise broke: 10 homogeneous L1 paths must aggregate (D45-1)")
	}
	p := l2Card()
	p.Level = "L1"
	p.Window = 3 * time.Second
	p.Deadline = 0
	p.Depth = 0
	p.Batch = bv
	text := joinCard(approval.RenderCard(p))
	for _, want := range []string{
		"批量确认：共 10 项",
		"涉及目录：",
		"明细（前 5 条）：",
		"其余 5 条一并由本次确认覆盖",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("batch card missing %q\nfull card:\n%s", want, text)
		}
	}
	// The preview must show the first 5 named and not the 6th.
	for _, item := range bv.First {
		if !strings.Contains(text, "    - "+item) {
			t.Errorf("preview line missing for %q", item)
		}
	}
	if strings.Contains(text, "    - "+d.Paths[5]) {
		t.Errorf("the 6th path leaked into the preview; D45-1 says 前 5 条")
	}
}

// TestTicket21S2NoBatchMeansNoBatchLines keeps the nil case silent: an
// unaggregated card must not hint that a batch was collapsed.
func TestTicket21S2NoBatchMeansNoBatchLines(t *testing.T) {
	text := joinCard(approval.RenderCard(l2Card()))
	if strings.Contains(text, "批量确认") {
		t.Errorf("nil BatchView printed batch lines:\n%s", text)
	}
}

// TestTicket21S2ChannelsSayUnavailableNotAvailable is B1's rendering half:
// an unloaded channel line must carry its 「不可用」 sentence, and a loaded one
// must not. The registry, not this file, decides which is which.
func TestTicket21S2ChannelsSayUnavailableNotAvailable(t *testing.T) {
	text := joinCard(approval.RenderCard(l2Card())) // DefaultChannels pair loaded
	// The loaded line names the channel; an unloaded one carries the whole
	// availability sentence the registry produced (B1: the 「不可用」 is part of
	// Text, not something the card may drop or imply).
	for _, want := range []string{"单击悬浮球（已加载）", "说取消词：语音取消不可用（未加载）"} {
		if !strings.Contains(text, want) {
			t.Errorf("channel line %q missing\nfull card:\n%s", want, text)
		}
	}
	if !strings.Contains(text, "语音取消不可用") {
		t.Errorf("the unloaded KWS line lost its 「语音取消不可用」 sentence (B1)\n%s", text)
	}
}

// TestTicket21S2DepthLineIsL2Only keeps the depth badge off the L1 route:
// an L1 window is not a queue item and has no depth to advertise.
func TestTicket21S2DepthLineIsL2Only(t *testing.T) {
	p := l2Card()
	p.Level = "L1"
	p.Window = 3 * time.Second
	p.Deadline = 0
	p.Depth = 3
	if text := joinCard(approval.RenderCard(p)); strings.Contains(text, "队列深度") {
		t.Errorf("L1 window card advertised a queue depth it cannot have:\n%s", text)
	}
}

// TestTicket21S2HeaderShapeStableForParsers is a compatibility ruler, not a
// design rule: cmd/wisp's cross-surface tests parse 「[确认 <level> <tool>]」
// and the contiguous 「  规则 」 block. The renderer must keep that shape or
// those parsers have to move with it - and this test says so at the producer.
func TestTicket21S2HeaderShapeStableForParsers(t *testing.T) {
	lines := approval.RenderCard(l2Card())
	if lines[0] != "[确认 L2 fs.write] 覆盖已有文件" {
		t.Errorf("header changed shape: %q", lines[0])
	}
	for _, l := range lines {
		if strings.HasPrefix(l, "规则") {
			t.Errorf("rule line lost its two-space indent: %q", l)
		}
	}
}

// TestTicket21S2RendererCarriesNoAuthority is the structural guard the
// renderer's header comment promises: RenderCard is pure over Prompt - it reads
// no registry, no clock, no queue, and this case proves it cannot mutate the
// Prompt it was handed either.
func TestTicket21S2RendererCarriesNoAuthority(t *testing.T) {
	p := l2Card()
	before := p
	_ = approval.RenderCard(p)
	if p.Grant != before.Grant || p.Depth != before.Depth || p.Level != before.Level {
		t.Error("RenderCard mutated the Prompt it was handed")
	}
}
