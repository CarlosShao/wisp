package risk

import "fmt"

// BatchScaleThreshold is the frozen R7 edge (SPEC-06 §3 R7): a single call
// affecting >= 50 files escalates to L2, so that "every one of them is L1"
// cannot hide "delete 5000 at once".
const BatchScaleThreshold = 50

// ruleBatchScale implements R7: single-call batch scale >= 50 files -> L2.
// The authoritative count is Facts.BatchCount; when unset, len(Paths) is the
// count. Negative counts are invalid input and treated as unset.
func ruleBatchScale(ctx *assessCtx) *contribution {
	n := ctx.facts.BatchCount
	if n <= 0 {
		n = len(ctx.facts.Paths)
	}
	if n < BatchScaleThreshold {
		return nil
	}
	return &contribution{
		rules:  []RuleID{R7},
		level:  L2,
		reason: fmt.Sprintf("R7: 单次调用影响 %d 个文件（≥%d）", n, BatchScaleThreshold),
	}
}
