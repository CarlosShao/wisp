package agent

import (
	"github.com/CarlosShao/wisp/internal/llm"
)

// D15/D39 context budgets (SPEC-05 §4). The frozen reference numbers were
// tuned for a mainstream 128k-token model ("12000 / 4000 这两个数是按主流 128k
// 上下文模型留了充足余量定的"), and D15 makes scaling by the configured model's
// context_window a MUST, not an optimization: "所有阈值必须按所配模型的
// context_window 等比缩小，不得硬编码（否则小上下文模型直接 400）".
//
// Nothing in this file is a magic number at a call site: every consumer reads
// a Budgets field produced by BudgetsFor.

// ReferenceContextWindow is the model context the D15/D39 reference numbers
// were calibrated against. Scaling is a pure ratio against it.
const ReferenceContextWindow = 128000

// Reference budget values (SPEC-05 §4.1 table, §4.2 spill, §4.3 compression,
// §8 budgets) as frozen for a 128k-token model.
const (
	refSectionIdentity      = 150  // (1) identity & boundaries
	refSectionSafety        = 100  // (5) safety & confirmation rules
	refSectionStyle         = 50   // (6) output & speech style
	refSectionResidentTools = 1200 // (2) resident tool index (builtin full set)
	refSectionProfile       = 400  // (3) L1 user profile, <=20 rows
	refSectionToolIndex     = 300  // (2) BM25 top-K (K<=5) third-party index
	refSectionScene         = 100  // (4) time/focus/scene - MUST be last
	refPromptTotal          = 2300 // D39 corrected total (the old 2000 did not add up)

	refSpillTokens    = 4000 // D15(3) single tool result spill threshold
	refSpillHead      = 500  // head tokens kept in context
	refSpillTail      = 200  // tail tokens kept in context
	refRawOutputBytes = 1 << 20
	refHistoryTokens  = 12000 // D15(4) compression trigger
	refTokenBudget    = 200000
)

// KeepRawRounds is the D15(4) "keep the last 3 rounds verbatim" rule. Rounds
// are structural counts, not token budgets, so they do not scale.
const KeepRawRounds = 3

// MaxToolConcurrency is the D38d ceiling (tool execution never exceeds 4
// concurrent calls). Like KeepRawRounds it is a concurrency count, not a
// context threshold, so it does not scale by context_window.
const MaxToolConcurrency = 4

// Budgets is the scaled budget set for one provider/model pair.
type Budgets struct {
	// ContextWindow actually used for scaling (ReferenceContextWindow when
	// the model declares nothing; a declared window bigger than the
	// reference never raises thresholds above the frozen values).
	ContextWindow int
	// Scale is ContextWindow/ReferenceContextWindow clamped to (0,1].
	Scale float64

	SectionIdentity      int
	SectionSafety        int
	SectionStyle         int
	SectionResidentTools int
	SectionProfile       int
	SectionToolIndex     int
	SectionScene         int
	// PromptTotal caps the assembled prompt sections (D39: <=2300 at the
	// reference window).
	PromptTotal int

	SpillTokens       int
	SpillHeadTokens   int
	SpillTailTokens   int
	RawOutputCapBytes int

	HistoryCompressTokens int
	KeepRawRounds         int

	// TokenBudget is the per-task token brake (C22). It scales with the
	// window too: the frozen 200k is the 128k-model value.
	TokenBudget int
}

// BudgetsFor scales every D15/D39 threshold by ctxWindow. An unknown window
// (0) is treated as the reference window: the frozen numbers then apply
// verbatim, which is the documented behaviour of the 128k baseline.
func BudgetsFor(ctxWindow int) Budgets {
	if ctxWindow <= 0 || ctxWindow > ReferenceContextWindow {
		ctxWindow = ReferenceContextWindow
	}
	b := Budgets{
		ContextWindow: ctxWindow,
		Scale:         float64(ctxWindow) / float64(ReferenceContextWindow),
		KeepRawRounds: KeepRawRounds,
	}
	s := func(ref int) int { return scaleInt(ref, b.Scale) }
	b.SectionIdentity = s(refSectionIdentity)
	b.SectionSafety = s(refSectionSafety)
	b.SectionStyle = s(refSectionStyle)
	b.SectionResidentTools = s(refSectionResidentTools)
	b.SectionProfile = s(refSectionProfile)
	b.SectionToolIndex = s(refSectionToolIndex)
	b.SectionScene = s(refSectionScene)
	b.PromptTotal = s(refPromptTotal)
	b.SpillTokens = s(refSpillTokens)
	b.SpillHeadTokens = s(refSpillHead)
	b.SpillTailTokens = s(refSpillTail)
	b.RawOutputCapBytes = s(refRawOutputBytes)
	b.HistoryCompressTokens = s(refHistoryTokens)
	b.TokenBudget = s(refTokenBudget)
	return b
}

// scaleInt shrinks ref by factor, never below 1 and never above ref (scaling
// only ever tightens: the frozen numbers are ceilings, not floors).
func scaleInt(ref int, factor float64) int {
	if factor >= 1 {
		return ref
	}
	v := int(float64(ref)*factor + 0.5)
	if v < 1 {
		v = 1
	}
	return v
}

// ApproxTokens is the token heuristic used for every budget decision in this
// package. It deliberately matches llm.Request.EstimateTokens (utf-8 bytes /
// 4) so the loop's budget math and the seam's pre-flight estimate cannot
// disagree; exact accounting comes from Usage events (C23).
func ApproxTokens(s string) int { return len(s) / 4 }

// ApproxTokensOf estimates the token cost of content parts.
func ApproxTokensOf(parts []llm.Content) int {
	n := 0
	for _, p := range parts {
		switch c := p.(type) {
		case llm.TextPart:
			n += ApproxTokens(c.Text)
		case llm.ToolUsePart:
			n += ApproxTokens(string(c.Input))
		case llm.ToolResultPart:
			n += ApproxTokensOf(c.Content)
		case llm.ImagePart:
			n += ApproxTokens(c.Alt)
		}
	}
	return n
}

// ApproxTokensOfMessage estimates one message.
func ApproxTokensOfMessage(m llm.Message) int { return ApproxTokensOf(m.Content) }
