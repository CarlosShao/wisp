package approval

import (
	"fmt"
	"strings"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// The two thresholds D45-1 and R7 fix in place. They are numbers from the
// contract, not tunables: changing one is a contract change.
const (
	// BatchMinOps is D45-1's "N >= 3 同质 L1 操作 -> 合并为一个确认".
	BatchMinOps = 3
	// R7BatchFloor is C19 R7's "单次调用影响 >= 50 个文件 -> L2".
	R7BatchFloor = 50
)

// Aggregate decides whether one call's L1 operations collapse into a single
// confirm, and returns the summary the card shows: total + affected dirs +
// first 5, with the remainder expandable (D45-1).
//
// It returns nil - "show every item, ask the slow way" - whenever aggregation
// would weaken a gate:
//
//   - the verdict is not L1: L2 is NEVER aggregated;
//   - fewer than BatchMinOps targets: nothing to gain;
//   - R7BatchFloor or more targets: that is an L2 batch (see
//     Gate.PendingWindow's escalation), so no cheap collapsed confirm;
//   - R4/C25 taint is in play (SessionOverrideBlocked or an R4 hit): the card
//     must name the source of the leaked content per item, which a summary by
//     definition hides;
//   - R3 tier-A/B sensitive paths: same reasoning.
func Aggregate(d tools.Decision) *BatchView {
	if d.Level != risk.L1 || len(d.Paths) < BatchMinOps || len(d.Paths) >= R7BatchFloor {
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
		"本次调用包含 %d 项同类 %s 操作（合并为一次确认，D45-1），涉及 %d 个目录：%s；前 %d 项：%s%s",
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

// hasRule scans the frozen rule list for one hit.
func hasRule(rules []risk.RuleID, want risk.RuleID) bool {
	for _, r := range rules {
		if r == want {
			return true
		}
	}
	return false
}

// dirOf takes the parent segment of an already-canonical path. It deliberately
// does NOT call filepath.Dir: the path is canonicalized by C26 already (the
// bridge's Decision.Paths), and filepath.* outside internal/risk's resolver is
// the D22 ban-2 pattern this repo's static scan fails on. Splitting a string
// that no longer needs resolving cannot re-open a normalization hole.
func dirOf(p string) string {
	t := strings.TrimRight(p, `/\`)
	if i := strings.LastIndexAny(t, `/\`); i > 0 {
		return t[:i]
	}
	if t == "" {
		return "."
	}
	return t
}

// joinShort renders at most n items, comma separated, with an ellipsis.
func joinShort(items []string, n int) string {
	if len(items) == 0 {
		return "（无）"
	}
	if n <= 0 || len(items) <= n {
		return strings.Join(items, "、")
	}
	return strings.Join(items[:n], "、") + fmt.Sprintf("…（共 %d 项）", len(items))
}
