package approval

// Readface wave 1 (2026-11-17), ticket 37 AC#5's Go half: the aggregate-card
// fields a PANEL card can carry.
//
// WHY A SECOND FUNCTION BESIDE Aggregate. Aggregate is the D45-1 CONTRACT
// surface: it decides whether many L1 operations collapse into ONE CHEAP
// CONFIRM, and it refuses L2 outright ("L2 is NEVER aggregated", batch.go's
// own header) and refuses tainted verdicts because a summary would hide the
// leak's source. CardSummary computes the DISPLAY fields only - total,
// affected dirs, the first preview rows, the deferred count and the summary
// sentence - and answers them for any decision with enough targets, including
// the R7-escalated L2 batch that Gate.PendingWindow already sends to the L2
// route (gate.go's belt-and-braces, unchanged here). Showing a panel card how
// many files one call touches aggregates nothing: the card's confirm is
// already one confirm for the whole call, the same one the native prompt
// raises. What it must never do is soften the verdict, so the taint rule
// carries over verbatim - a tainted card names its sources per item, and a
// summary would hide exactly that.

import (
	"fmt"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// CardSummary returns the panel-card display fields for one decision, or nil
// when the card must not show a summary at all:
//
//   - fewer than BatchMinOps targets: nothing to summarize;
//   - R4/C25 taint in play (SessionOverrideBlocked or an R4 hit): same rule
//     as Aggregate - the card must name the leak's source per item;
//   - R3 tier-A/B sensitive paths: same reasoning.
//
// Level does NOT gate this function: the R7-escalated batch arrives at the L2
// queue as one decision carrying 50+ paths, and its card still shows "what is
// in this batch". Aggregate keeps the confirm-aggregation contract untouched.
func CardSummary(d tools.Decision) *BatchView {
	if len(d.Paths) < BatchMinOps {
		return nil
	}
	if d.SessionOverrideBlocked || hasRule(d.RulesHit, risk.R4) || hasRule(d.RulesHit, risk.R3) {
		return nil
	}

	dirs := map[string]bool{}
	for _, p := range d.Paths {
		dirs[dirOf(p)] = true
	}
	uniq := make([]string, 0, len(dirs))
	for k := range dirs {
		uniq = append(uniq, k)
	}
	uniq = sortedStrings(uniq)

	first := len(d.Paths)
	if first > BatchPreviewCount {
		first = BatchPreviewCount
	}
	view := &BatchView{
		Total:        len(d.Paths),
		AffectedDirs: uniq,
		First:        append([]string(nil), d.Paths[:first]...),
		Deferred:     len(d.Paths) - first,
	}
	view.Summary = fmt.Sprintf(
		"本次调用包含 %d 项同类 %s 操作，涉及 %d 个目录：%s；前 %d 项：%s%s",
		view.Total, d.LevelString(), len(uniq), joinShort(uniq, 4), len(view.First),
		joinShort(view.First, len(view.First)),
		func() string {
			if view.Deferred > 0 {
				return fmt.Sprintf("（另有 %d 项可展开查看）", view.Deferred)
			}
			return ""
		}())
	return view
}
