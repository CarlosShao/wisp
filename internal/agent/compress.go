package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/llm"
)

// D15(4) session-history compression (SPEC-05 §4.3): once the history exceeds
// the scaled threshold, the OLDEST rounds are compressed into one summary
// message while the last KeepRawRounds rounds stay verbatim; if the total is
// still over threshold the next-oldest rounds are folded too, until the budget
// holds or only KeepRawRounds rounds are left.
//
// Hard rule: tool_call ids and result references are NEVER dropped - the
// truncation-failure judgement (D21) and the C25 provenance chain depend on
// them. The compressed block therefore carries an id ledger that only ever
// grows.
//
// Placement: D15(4) says compression runs OUTSIDE the response path, in the
// Warm-window hook (ticket 28). Until that hook exists, Compress is called
// synchronously from the loop before the request is assembled (the spec allows
// the S1 fallback and requires it to be flagged).
// DEFERRED(D28-1): move the call into the Warm window hook; keep this
// synchronous path only as the fallback when no hook is registered.

// summaryMarker prefixes the compressed-history block so the compressor never
// re-folds its own output as if it were raw conversation.
const summaryMarker = "[已压缩的历史]"

// Summarizer compresses a slice of history into prose. The summarize role
// (llm.roles.summarize) is wired by whoever owns the role chain; with no
// Summarizer installed the loop uses the structural trim below, which keeps
// every tool_call id and result reference.
type Summarizer interface {
	SummarizeHistory(ctx context.Context, msgs []llm.Message) (string, error)
}

// CompressionReport describes one compression pass (assertion + log surface).
type CompressionReport struct {
	Ran            bool
	CompressedMsgs int
	KeptRawRounds  int
	PreservedIDs   []string
	// RawIDsInsideSummary lists tool_call ids that exist only inside the
	// compressed block (proof the ledger survived compression).
	RawIDsInsideSummary []string
	TokensBefore        int
	TokensAfter         int
}

// Compressor applies the D15(4) policy against a scaled budget set.
type Compressor struct {
	b   Budgets
	sum Summarizer
}

// NewCompressor builds a compressor; sum may be nil (structural trim).
func NewCompressor(b Budgets, sum Summarizer) *Compressor {
	return &Compressor{b: b, sum: sum}
}

// TotalTokens estimates the history cost.
func (c *Compressor) TotalTokens(hist []llm.Message) int {
	n := 0
	for _, m := range hist {
		n += ApproxTokensOf(m.Content)
	}
	return n
}

// Need reports whether the history is over the compression threshold.
func (c *Compressor) Need(hist []llm.Message) bool {
	return c.TotalTokens(hist) > c.b.HistoryCompressTokens
}

// round is one conversational round (a user turn plus what followed it).
type round struct {
	msgs      []llm.Message
	isSummary bool
}

// Compress folds the oldest rounds until the history fits the threshold or
// only KeepRawRounds raw rounds remain. It never mutates the input slices'
// content and always returns a freshly built history.
func (c *Compressor) Compress(ctx context.Context, hist []llm.Message) ([]llm.Message, CompressionReport, error) {
	rep := CompressionReport{TokensBefore: c.TotalTokens(hist)}
	rounds := groupRounds(hist)
	// ledger accumulates every tool_call id ever folded, so it can only grow.
	var ledger []string

	for c.totalOf(rounds) > c.b.HistoryCompressTokens {
		raw := rawRoundIndexes(rounds)
		if len(raw) <= c.b.KeepRawRounds {
			break // D15(4): stop at the last 3 raw rounds
		}
		fold := raw[:len(raw)-c.b.KeepRawRounds]

		var victims []llm.Message
		// A pre-existing summary block is re-folded in place so exactly one
		// compressed block ever exists and it cannot grow without bound.
		prefixSummary := !rounds[raw[0]].isSummary
		if idx := summaryRoundIndex(rounds); idx >= 0 {
			victims = append(victims, rounds[idx].msgs...)
			_ = prefixSummary
		}
		for _, i := range fold {
			victims = append(victims, rounds[i].msgs...)
		}
		ledger = append(ledger, toolCallIDs(victims)...)

		summary, err := c.summarize(ctx, victims)
		if err != nil {
			return hist, rep, err
		}
		block := llm.Message{
			Role:    llm.RoleUser,
			Content: []llm.Content{llm.TextPart{Text: summaryBlock(summary, ledger)}},
		}
		rounds = rebuildRounds(rounds, fold, block, summaryMarkerIn(victims))
		rep.Ran = true
		rep.CompressedMsgs += len(victims)
	}

	out := flattenRounds(rounds)
	rep.TokensAfter = c.TotalTokens(out)
	rep.KeptRawRounds = len(rawRoundIndexes(rounds))
	rep.PreservedIDs = dedupe(append(allToolCallIDs(out), ledger...))
	rep.RawIDsInsideSummary = dedupe(ledger)
	return out, rep, nil
}

func (c *Compressor) totalOf(rounds []round) int {
	n := 0
	for _, r := range rounds {
		n += c.TotalTokens(r.msgs)
	}
	return n
}

func (c *Compressor) summarize(ctx context.Context, msgs []llm.Message) (string, error) {
	if c.sum == nil {
		return structuralTrim(msgs), nil
	}
	s, err := c.sum.SummarizeHistory(ctx, msgs)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(s) == "" {
		return structuralTrim(msgs), nil
	}
	return s, nil
}

// summaryBlock renders the compressed-history message; the id ledger is always
// appended so ids survive even a model-written summary that omits them.
func summaryBlock(summary string, ledger []string) string {
	var b strings.Builder
	b.WriteString(summaryMarker)
	b.WriteString(" 以下轮次已被摘要压缩，仅供上下文参考：\n")
	b.WriteString(strings.TrimSpace(summary))
	if len(ledger) > 0 {
		b.WriteString("\n[tool_call 记录（不得丢弃）] ")
		b.WriteString(strings.Join(dedupe(ledger), "; "))
	}
	return b.String()
}

// structuralTrim is the no-model compression: every tool call id, tool name,
// argument size and result status survives; prose bodies are shortened.
func structuralTrim(msgs []llm.Message) string {
	var b strings.Builder
	for _, m := range msgs {
		for _, p := range m.Content {
			switch c := p.(type) {
			case llm.TextPart:
				if t := strings.TrimSpace(c.Text); t != "" {
					fmt.Fprintf(&b, "%s: %s\n", m.Role, firstRunes(t, 80))
				}
			case llm.ToolUsePart:
				fmt.Fprintf(&b, "tool_call %s name=%s args=%d字节\n", c.ID, c.Name, len(c.Input))
			case llm.ToolResultPart:
				status := "ok"
				if c.IsError {
					status = "error"
				}
				fmt.Fprintf(&b, "tool_result %s %s (%d token)\n", c.ID, status, ApproxTokensOf(c.Content))
			case llm.ImagePart:
				fmt.Fprintf(&b, "image mime=%s (字节引用不入历史)\n", c.MimeType)
			}
		}
	}
	return b.String()
}

func firstRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

// toolCallIDs lists tool_use ids plus tool_result references of a slice.
func toolCallIDs(msgs []llm.Message) []string {
	var out []string
	for _, m := range msgs {
		for _, p := range m.Content {
			switch c := p.(type) {
			case llm.ToolUsePart:
				out = append(out, c.ID+"("+c.Name+")")
			case llm.ToolResultPart:
				out = append(out, c.ID)
			}
		}
	}
	return out
}

func allToolCallIDs(hist []llm.Message) []string {
	var out []string
	for _, m := range hist {
		for _, p := range m.Content {
			if c, ok := p.(llm.ToolUsePart); ok {
				out = append(out, c.ID)
			}
		}
	}
	return out
}

// groupRounds splits a history into rounds. A round starts at a user-role
// message (tool results ride RoleTool messages, assistant turns follow their
// round) and holds everything up to the next user message.
func groupRounds(hist []llm.Message) []round {
	var rounds []round
	for _, m := range hist {
		startNew := len(rounds) == 0 || m.Role == llm.RoleUser
		if startNew {
			r := round{msgs: []llm.Message{m}}
			if isSummaryMessage(m) {
				r.isSummary = true
			}
			rounds = append(rounds, r)
			continue
		}
		rounds[len(rounds)-1].msgs = append(rounds[len(rounds)-1].msgs, m)
	}
	return rounds
}

func isSummaryMessage(m llm.Message) bool {
	if m.Role != llm.RoleUser {
		return false
	}
	for _, p := range m.Content {
		if t, ok := p.(llm.TextPart); ok && strings.HasPrefix(t.Text, summaryMarker) {
			return true
		}
	}
	return false
}

func rawRoundIndexes(rounds []round) []int {
	var out []int
	for i, r := range rounds {
		if !r.isSummary {
			out = append(out, i)
		}
	}
	return out
}

func summaryRoundIndex(rounds []round) int {
	for i, r := range rounds {
		if r.isSummary {
			return i
		}
	}
	return -1
}

func summaryMarkerIn(msgs []llm.Message) bool {
	for _, m := range msgs {
		if isSummaryMessage(m) {
			return true
		}
	}
	return false
}

// rebuildRounds drops the folded rounds (and any old summary round, which was
// re-folded into the new block) and inserts block at the front.
func rebuildRounds(rounds []round, fold []int, block llm.Message, hadSummary bool) []round {
	drop := map[int]bool{}
	if hadSummary {
		if i := summaryRoundIndex(rounds); i >= 0 {
			drop[i] = true
		}
	}
	for _, i := range fold {
		drop[i] = true
	}
	out := []round{{msgs: []llm.Message{block}, isSummary: true}}
	for i, r := range rounds {
		if drop[i] {
			continue
		}
		out = append(out, r)
	}
	return out
}

func flattenRounds(rounds []round) []llm.Message {
	var out []llm.Message
	for _, r := range rounds {
		out = append(out, r.msgs...)
	}
	return out
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
