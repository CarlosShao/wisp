package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

// Acceptance criterion 4 (D15(4) compression): history over the scaled
// threshold folds the OLDEST rounds into one summary while the last
// KeepRawRounds (3) stay verbatim. The hard rule - the one an implementation
// can silently get wrong - is that tool_call ids and their results are NEVER
// dropped: the C25 provenance chain and the truncation-failure judgement both
// key on them. So the assertions below track ids, not just token counts.

// buildRoundHistory makes `rounds` conversational rounds, each ~padBytes of
// user text plus an echo tool call and its result, with a unique call id
// "cNN". With padBytes chosen so the total clears HistoryCompressTokens.
func buildRoundHistory(rounds, padBytes int) []llm.Message {
	pad := strings.Repeat("x", padBytes)
	hist := make([]llm.Message, 0, rounds*3)
	for i := 0; i < rounds; i++ {
		id := fmt.Sprintf("c%02d", i)
		hist = append(hist,
			llm.Message{
				Role:    llm.RoleUser,
				Content: []llm.Content{llm.TextPart{Text: fmt.Sprintf("问题 %d %s", i, pad)}},
			},
			llm.Message{
				Role: llm.RoleAssistant,
				Content: []llm.Content{llm.ToolUsePart{
					ID: id, Name: "echo",
					Input: json.RawMessage(`{"text":"hi"}`),
				}},
			},
			llm.Message{
				Role: llm.RoleTool,
				Content: []llm.Content{llm.ToolResultPart{
					ID:      id,
					Content: []llm.Content{llm.TextPart{Text: "结果 " + pad}},
				}},
			},
		)
	}
	return hist
}

func TestCompressFoldsOldestKeepsLastThreePreservesIDs(t *testing.T) {
	b := BudgetsFor(128000) // HistoryCompressTokens = 12000
	if b.HistoryCompressTokens != 12000 {
		t.Fatalf("128k compress threshold = %d, want the frozen 12000", b.HistoryCompressTokens)
	}
	c := NewCompressor(b, nil)

	const rounds = 15
	hist := buildRoundHistory(rounds, 4000) // ~2000 tokens/round -> ~30000 total
	if !c.Need(hist) {
		t.Fatalf("Need=false for %d tokens, want true (> %d)", c.TotalTokens(hist), b.HistoryCompressTokens)
	}
	before := c.TotalTokens(hist)

	out, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if !rep.Ran {
		t.Fatal("Ran=false, want a compression pass")
	}
	if rep.TokensBefore != before {
		t.Errorf("report TokensBefore = %d, want %d", rep.TokensBefore, before)
	}
	if rep.TokensAfter >= rep.TokensBefore {
		t.Errorf("tokens did not shrink: %d -> %d", rep.TokensBefore, rep.TokensAfter)
	}
	if c.TotalTokens(out) > b.HistoryCompressTokens {
		t.Errorf("compressed history still %d tokens, over the %d threshold",
			c.TotalTokens(out), b.HistoryCompressTokens)
	}

	// D15(4): exactly the last KeepRawRounds rounds survive verbatim.
	if rep.KeptRawRounds != b.KeepRawRounds {
		t.Errorf("kept raw rounds = %d, want %d", rep.KeptRawRounds, b.KeepRawRounds)
	}
	if len(rawRoundIndexes(groupRounds(out))) != b.KeepRawRounds {
		t.Errorf("output holds %d raw rounds, want %d:\n%s",
			len(rawRoundIndexes(groupRounds(out))), b.KeepRawRounds, dumpHistory(out))
	}

	// The LAST three rounds' ids remain LIVE tool calls (not folded), i.e. their
	// tool_use + tool_result structure is intact.
	for i := rounds - b.KeepRawRounds; i < rounds; i++ {
		id := fmt.Sprintf("c%02d", i)
		if !hasLiveToolUse(out, id) {
			t.Errorf("round %d call %s was folded, but it is inside the last 3 (must stay raw)", i, id)
		}
		if !hasLiveToolResult(out, id) {
			t.Errorf("round %d result %s lost its tool_result part", i, id)
		}
	}

	// Every folded round's id is GONE as a live call but PRESENT in the ledger
	// block: the never-drop rule is satisfied via the id ledger, not by keeping
	// the prose.
	for i := 0; i < rounds-b.KeepRawRounds; i++ {
		id := fmt.Sprintf("c%02d", i)
		if hasLiveToolUse(out, id) {
			t.Errorf("folded call %s still live, want folded", id)
		}
		if !histTextHas(out, id) {
			t.Errorf("id %s vanished from the compressed history", id)
		}
	}

	// The report's PreservedIDs is the full id set (assertion surface used by
	// callers/logs): all 15 calls accounted for.
	if got := len(dedupe(rep.PreservedIDs)); got < rounds {
		t.Errorf("PreservedIDs = %d entries, want >= %d", got, rounds)
	}

	// Structural-trim summary is a single block; the ledger line is present.
	if !histTextHas(out, "[tool_call 记录（不得丢弃）]") {
		t.Errorf("compressed block lacks the id ledger marker:\n%s", dumpHistory(out))
	}
}

// A model-written summary that omits the call ids must not break the never-drop
// rule: the compressor appends the id ledger regardless of the prose.
type droppingSummarizer struct{}

func (droppingSummarizer) SummarizeHistory(context.Context, []llm.Message) (string, error) {
	return "（模型已省略具体调用标识的散文摘要）", nil
}

func TestCompressPreservesIDsEvenWhenSummarizerDropsThem(t *testing.T) {
	b := BudgetsFor(128000)
	c := NewCompressor(b, droppingSummarizer{})
	const rounds = 15
	hist := buildRoundHistory(rounds, 4000)

	out, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if !rep.Ran {
		t.Fatal("expected a compression pass")
	}
	// Folded ids survive via the ledger even though the prose summary dropped them.
	for i := 0; i < rounds-b.KeepRawRounds; i++ {
		id := fmt.Sprintf("c%02d", i)
		if !histTextHas(out, id) {
			t.Errorf("id %s lost: model summary dropped it and the ledger failed to carry it", id)
		}
	}
	// The last three rounds stay live regardless.
	for i := rounds - b.KeepRawRounds; i < rounds; i++ {
		if !hasLiveToolUse(out, fmt.Sprintf("c%02d", i)) {
			t.Errorf("kept-round id c%02d must remain a live tool call", i)
		}
	}
}

// Compression only runs past the threshold, and that threshold scales with the
// context window: a history a 128k model keeps intact forces a 4k model to
// fold - proving the trigger is config-derived, not a bare 12000 literal.
func TestCompressTriggerScalesWithWindow(t *testing.T) {
	big, tiny := BudgetsFor(128000), BudgetsFor(4096)
	if big.HistoryCompressTokens != 12000 {
		t.Fatalf("128k threshold = %d, want 12000", big.HistoryCompressTokens)
	}
	want := (4096 * 12000) / 128000 // 384
	if tiny.HistoryCompressTokens != want {
		t.Fatalf("4096 threshold = %d, want proportional %d", tiny.HistoryCompressTokens, want)
	}

	// A history that sits between the two thresholds: > tiny (384), < big
	// (12000). Five rounds so the tiny pass has oldest rounds to fold while
	// still keeping three verbatim.
	pad := strings.Repeat("y", 400) // ~100 tokens/round -> ~500 total
	const rounds = 5
	var hist []llm.Message
	for i := 0; i < rounds; i++ {
		id := fmt.Sprintf("s%02d", i)
		hist = append(hist,
			llm.Message{
				Role:    llm.RoleUser,
				Content: []llm.Content{llm.TextPart{Text: fmt.Sprintf("问 %d %s", i, pad)}},
			},
			llm.Message{
				Role: llm.RoleAssistant,
				Content: []llm.Content{llm.ToolUsePart{
					ID: id, Name: "echo",
					Input: json.RawMessage(`{}`),
				}},
			},
		)
	}
	n := NewCompressor(big, nil).TotalTokens(hist)
	if n <= tiny.HistoryCompressTokens || n >= big.HistoryCompressTokens {
		t.Fatalf("setup: %d tokens not between tiny %d and big %d", n, tiny.HistoryCompressTokens, big.HistoryCompressTokens)
	}
	if NewCompressor(big, nil).Need(hist) {
		t.Error("128k model must NOT compress this history")
	}
	if !NewCompressor(tiny, nil).Need(hist) {
		t.Fatal("4k model must compress this history (scaled threshold crossed)")
	}
	out, rep, err := NewCompressor(tiny, nil).Compress(context.Background(), hist)
	if err != nil || !rep.Ran {
		t.Fatalf("tiny-window compression did not run: ran=%v err=%v", rep.Ran, err)
	}
	// ids still preserved through the scaled-window pass.
	for i := 0; i < rounds; i++ {
		if !histTextHas(out, fmt.Sprintf("s%02d", i)) {
			t.Errorf("scaled-window compression dropped id s%02d", i)
		}
	}
}

// Under threshold, compression is a no-op (Ran=false) and never runs.
func TestCompressNoopUnderThreshold(t *testing.T) {
	c := NewCompressor(BudgetsFor(128000), nil)
	hist := buildRoundHistory(2, 100) // tiny
	if c.Need(hist) {
		t.Fatal("Need=true for a tiny history")
	}
	out, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatal(err)
	}
	if rep.Ran {
		t.Error("Ran=true, want a no-op")
	}
	if c.TotalTokens(out) != c.TotalTokens(hist) {
		t.Error("no-op pass altered the history")
	}
}

// ---------------------------------------------------------------------------

func hasLiveToolUse(hist []llm.Message, id string) bool {
	for _, m := range hist {
		for _, p := range m.Content {
			if c, ok := p.(llm.ToolUsePart); ok && c.ID == id {
				return true
			}
		}
	}
	return false
}

func hasLiveToolResult(hist []llm.Message, id string) bool {
	for _, m := range hist {
		for _, p := range m.Content {
			if c, ok := p.(llm.ToolResultPart); ok && c.ID == id {
				return true
			}
		}
	}
	return false
}

// histTextHas reports whether `needle` appears anywhere in the history: a live
// text part, or a tool_use / tool_result id.
func histTextHas(hist []llm.Message, needle string) bool {
	for _, m := range hist {
		for _, p := range m.Content {
			switch c := p.(type) {
			case llm.TextPart:
				if strings.Contains(c.Text, needle) {
					return true
				}
			case llm.ToolUsePart:
				if c.ID == needle {
					return true
				}
			case llm.ToolResultPart:
				if c.ID == needle {
					return true
				}
			}
		}
	}
	return false
}
