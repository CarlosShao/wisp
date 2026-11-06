package approval

// Native card rendering (ticket 21 segment 2).
//
// WHY THIS FILE EXISTS. Segment 1 proved the card's DATA arrives complete:
// Prompt.Params / Prompt.Args / Prompt.RulesHit / Prompt.Batch are asserted
// through the bridge (docs/evidence/s1/21-segment1-adversarial-acceptance.md
// carries that reading, and the ticket's progress log names the card itself as
// the segment-2 half). But no code in this tree turned those fields into the
// card a human reads: cmd/wisp's console printer printed level/tool/reason,
// rules, channels, paths and the window - and stopped. The full parameter
// object, which SPEC-06 §7 names as queue-item content ("每项含
// correlationId/所属任务/工具名/完整参数/风险级/请求理由") and the ticket
// states as C27 fallback-grade ("FULL params (not summary)"), never reached
// any surface, and neither did the D45-1 batch view that Aggregate builds and
// gate.go hands to every Prompt.
//
// ONE COPY. Both native surfaces - cmd/wisp/run.go's consoleApprovalUI and
// cmd/wisp/resident_approval_windows.go's ballCardUI - render through this one
// function, for the reason approval_reply.go states for the answer side: a
// second copy of the same content logic is how one surface grows a detail the
// other does not show, and on a security card a missing detail is a missing
// disclosure. The panel-facing projection is NOT here and cannot use this file:
// PanelItem carries no grant (ui.go), and the panel's card half belongs to
// ticket 37.
//
// WHAT THIS IS NOT: an authority. RenderCard is a pure function over Prompt -
// it reads no registry, mints no grant, settles nothing. Every decision the
// lines describe was already made by the Gate before this function was
// reachable.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// RenderCard turns one confirmation into the ordered text lines of a native
// card, exactly as SPEC-06 §7 and §8.1 describe it: identity and reason first,
// the verbatim rules that raised the level, the FULL parameters (not a
// summary), the D45-1 batch view when one exists, the affected paths, the
// veto channels with their honest availability, and the time the route waits.
//
// The output is lines without trailing newlines; the host decides its own
// prefix and blank-line shape. The header keeps cmd/wisp's historical
// 「[确认 <level> <tool>] <reason>」 spelling and every indented line keeps
// the two-space shape the console's cross-surface tests parse.
func RenderCard(p Prompt) []string {
	lines := []string{fmt.Sprintf("[确认 %s %s] %s", p.Level, p.Tool, p.Reason)}
	for _, r := range p.RulesHit {
		lines = append(lines, fmt.Sprintf("  规则 %s", r))
	}
	if strings.TrimSpace(p.TaskID) != "" {
		lines = append(lines, fmt.Sprintf("  所属任务：%s", p.TaskID))
	}
	if p.Level == "L2" && p.Depth > 0 {
		// SPEC-06 §7's 「队头单显 + 深度角标（悬浮球必显，不能只藏面板）」:
		// the orb carries the same number as a badge (cmd/wisp's
		// approvalDepthBadge), and the text card states it too.
		lines = append(lines, fmt.Sprintf("  队列深度：%d", p.Depth))
	}
	lines = append(lines, renderBatch(p.Batch)...)
	lines = append(lines, renderParams(p)...)
	if len(p.Paths) > 0 {
		lines = append(lines, fmt.Sprintf("  影响路径：%s", strings.Join(p.Paths, ", ")))
	}
	for _, c := range p.Channels {
		lines = append(lines, fmt.Sprintf("  取消方式：%s（%s）", c.Text, channelAvailability(c)))
	}
	shown := p.Window
	if p.Level == "L2" { // an L2 card's Window is 0 by construction; its wait is the Deadline
		shown = p.Deadline
	}
	lines = append(lines, fmt.Sprintf("  窗口 %.1fs", shown.Seconds()))
	return lines
}

// renderBatch prints the D45-1 view (§8.1: total + affected-dirs summary +
// first 5 + expandable remainder). A nil view is the ordinary case and prints
// nothing; the data says so upstream - Aggregate answers nil for every
// non-L1, sub-3, R7-sized, R4-tainted or R3-sensitive call, and promptFor
// passes that nil through unchanged.
func renderBatch(bv *BatchView) []string {
	if bv == nil {
		return nil
	}
	out := []string{
		fmt.Sprintf("  批量确认：共 %d 项同类操作，已合并为这一次确认（D45-1；否决＝全部否决）", bv.Total),
		fmt.Sprintf("  涉及目录：%s", joinDirs(bv.AffectedDirs)),
	}
	if len(bv.First) > 0 {
		out = append(out, fmt.Sprintf("  明细（前 %d 条）：", len(bv.First)))
		for _, item := range bv.First {
			out = append(out, "    - "+item)
		}
	}
	if bv.Deferred > 0 {
		out = append(out, fmt.Sprintf("  其余 %d 条一并由本次确认覆盖，完整列表见下方完整参数", bv.Deferred))
	}
	return out
}

// renderParams prints the FULL parameter object, never a summary.
//
// Priority is Args first, Params second. Decision.Args is the model's own
// argument JSON relayed verbatim (internal/tools/bridge.go), so printing it is
// printing what the tool will actually receive - the strongest form of "full
// params" this layer can offer. Params is the parsed view the assessor read;
// when it is all there is, it is re-serialized with Go's deterministic map-key
// ordering. A card whose Prompt carries neither says so out loud instead of
// showing an empty block a reader could mistake for "no parameters existed".
//
// No truncation, no elision, no redaction here: a secret-bearing value would
// have to be washed at the producer (the queue item is defined to carry the
// full params), and inventing a local filter would silently break the one
// property this block exists to keep.
func renderParams(p Prompt) []string {
	if raw := bytes.TrimSpace(p.Args); len(raw) > 0 {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, raw, "    ", "  "); err != nil {
			// Un-parseable (a hand-built Prompt in a test, not a bridge
			// product - bridge validates before routing). Verbatim beats
			// "the card dropped the field it could not format".
			return []string{"  完整参数（非摘要）：", "    " + string(raw)}
		}
		out := []string{"  完整参数（非摘要）："}
		for _, line := range strings.Split(pretty.String(), "\n") {
			out = append(out, "    "+line)
		}
		return out
	}
	if len(p.Params) > 0 {
		raw, err := json.MarshalIndent(p.Params, "    ", "  ")
		if err != nil {
			return []string{"  完整参数（非摘要）：（参数无法序列化，原始对象未随卡片送达）"}
		}
		out := []string{"  完整参数（非摘要）："}
		for _, line := range strings.Split(string(raw), "\n") {
			out = append(out, "    "+line)
		}
		return out
	}
	return []string{"  完整参数（非摘要）：（本次卡片未随附参数对象）"}
}

// joinDirs renders the affected-directory set the way the batch summary words
// it: enumerated, never just counted, because "which dirs" is the sentence
// D45-1 asks the user to weigh against one confirm.
func joinDirs(dirs []string) string {
	if len(dirs) == 0 {
		return "（无）"
	}
	return strings.Join(dirs, "、")
}

// channelAvailability is the loaded/unloaded suffix of a veto-channel line.
// The words are cmd/wisp's own historical pair (channelState before this file
// existed); moving them here keeps both native surfaces saying the same thing.
func channelAvailability(c ChannelStatus) string {
	if c.Loaded {
		return "已加载"
	}
	return "未加载"
}
