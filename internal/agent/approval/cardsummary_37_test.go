package approval_test

// Readface wave 1 (2026-11-17), ticket 37 AC#5's Go teeth: CardSummary is the
// panel card's aggregate-FIELD face. It is behavior over real decisions, and
// its negative half pins the contract line that must NOT move: Aggregate
// still refuses L2 and taint (D45-1 / C25 / R3), while the card fields ride
// the same decision.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
)

// TestCardSummaryTenOpsCarriesTheAggregateFields is the AC#5 scenario: ten
// targets in one call produce the total / dirs / first-5 / deferred fields a
// panel card renders.
func TestCardSummaryTenOpsCarriesTheAggregateFields(t *testing.T) {
	d := l1Decision("C:/data/x/a.txt", "C:/data/x/b.txt", "C:/data/x/c.txt",
		"C:/data/y/d.txt", "C:/data/y/e.txt", "C:/data/y/f.txt",
		"C:/data/y/g.txt", "C:/data/y/h.txt", "C:/data/y/i.txt",
		"C:/data/y/j.txt")
	s := approval.CardSummary(d)
	if s == nil {
		t.Fatal("10 项目标必须给出聚合卡字段（AC#5），得 nil")
	}
	if s.Total != 10 {
		t.Errorf("Total=%d，期望 10", s.Total)
	}
	if len(s.First) != approval.BatchPreviewCount {
		t.Errorf("First=%d 条，期望前 %d", len(s.First), approval.BatchPreviewCount)
	}
	if s.Deferred != 5 {
		t.Errorf("Deferred=%d，期望 5", s.Deferred)
	}
	if len(s.AffectedDirs) != 2 {
		t.Errorf("AffectedDirs=%v，期望 2 个目录", s.AffectedDirs)
	}
	if s.Summary == "" {
		t.Error("Summary 是空的：卡片没有可读的聚合句")
	}
}

// TestCardSummarySurvivesR7Escalation: the escalated batch arrives at the L2
// queue as ONE decision carrying 50+ paths, and its CARD still says what is
// in the batch - while Aggregate (the confirm-collapse contract) stays nil.
func TestCardSummarySurvivesR7Escalation(t *testing.T) {
	l1Big := l1Decision(manyPaths(60)...)
	escalated := l1Big
	escalated.Level = risk.L2
	escalated.RulesHit = append(append([]risk.RuleID(nil), escalated.RulesHit...), risk.R7)

	if s := approval.CardSummary(escalated); s == nil || s.Total != 60 {
		t.Errorf("升级后的 L2 批次卡字段 = %+v，期望 Total=60（卡要说清这一批是什么）", s)
	}
	if bv := approval.Aggregate(escalated); bv != nil {
		t.Errorf("Aggregate 对 L2 返回了 %+v：确认合并的契约不许动（R7 批次不得聚合成一次便宜确认）", bv)
	}
}

// TestCardSummaryRefusesTaintAndSmallBatches: the taint rule carries over
// verbatim from Aggregate - a summary would hide the leak's source.
func TestCardSummaryRefusesTaintAndSmallBatches(t *testing.T) {
	if s := approval.CardSummary(l1Decision("C:/a.txt", "C:/b.txt")); s != nil {
		t.Errorf("2 项不足 %d 门槛，卡字段也不该出：%+v", approval.BatchMinOps, s)
	}
	tainted := l1Decision("C:/a.txt", "C:/b.txt", "C:/c.txt")
	tainted.SessionOverrideBlocked = true
	if s := approval.CardSummary(tainted); s != nil {
		t.Errorf("R4 污染卡字段也不该出（摘要会藏泄露源）：%+v", s)
	}
	r4 := l1Decision("C:/a.txt", "C:/b.txt", "C:/c.txt")
	r4.RulesHit = []risk.RuleID{risk.R4}
	if s := approval.CardSummary(r4); s != nil {
		t.Errorf("R4 命中卡字段也不该出：%+v", s)
	}
	r3 := l1Decision("C:/a.txt", "C:/b.txt", "C:/c.txt")
	r3.RulesHit = []risk.RuleID{risk.R3}
	if s := approval.CardSummary(r3); s != nil {
		t.Errorf("R3 敏感路径卡字段也不该出：%+v", s)
	}
}
