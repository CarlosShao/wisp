package observe

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Per-model quota gates and per-provider plan-credit bookkeeping (ticket 44
// AC#5, SPEC-03 §3.1 "配额三层" as quoted by the ticket at :29-37).
//
// The three layers, kept distinct because they expire differently:
//
//   - per-model caps: quota_daily_micro / quota_monthly_micro (config,
//     internal/config/schema.go:372-373 — the shape is SPEC-defined, this code
//     does not invent it);
//   - per-provider plan credit: plan_credit_total_micro for billing=plan
//     (schema.go:397). The total is user-entered, the consumption is C23
//     bookkeeping, and the remaining is marked 估算 ("estimated") because no
//     public balance API is assumed (票 44:30-32);
//   - the global [cost] budgets: that is BudgetParams/costbudget.go, a
//     different gate checked elsewhere; plan vs pay-per-token consumption is
//     tracked separately per model here, exactly as AC#5 requires.
//
// Anti-negative is a PRE-DISPATCH check (票 44:33): JudgeModel answers "may
// this model be dispatched at all", and SelectChain walks the ordered
// text_chain so an over-quota model comes back in the skipped set — the chain
// stepping itself belongs to ticket 11 (票 44:39, "C23 owns the quota signal").
//
// Money stays int64 micro-units end to end; the one float in the neighborhood
// (the alert ratio) is converted once by BudgetFrom, and this file reuses its
// basisPoints helper rather than growing a second percentage calculator.
//
// Persistence note (named, not hidden): D35's cost_daily/task_log have no
// model dimension and their DDL is contract-frozen, so this ledger is
// session-scoped. The permanent per-model row is the owner-approved schema
// step (44-v1-verdict §3.3.1/§7.1); until then the disk-side quota signal is
// provider_health.quota_state (SPEC-02 §3, schema v2), which ticket 11's
// probe writer owns.

// Billing kinds the quota gate is aware of (SPEC-03 §3.1 billing enum).
const (
	BillingPayPerToken = "pay-per-token"
	BillingPlan        = "plan"
)

// PlanCreditEstimatedMarker labels every plan-credit remaining figure: the
// total is manually entered and consumption is bookkeeping, never a balance
// read (票 44:31 "remaining marked estimated").
const PlanCreditEstimatedMarker = "估算"

// QuotaModel is one model's quota card, adapted from
// models.<id>{quota_daily_micro, quota_monthly_micro, billing}. A zero cap is
// "no cap" (uncapped), never "block everything" — the same reading as
// BudgetParams' zero caps.
type QuotaModel struct {
	Provider      string
	Model         string
	Billing       string // BillingPayPerToken | BillingPlan
	DailyMicros   int64
	MonthlyMicros int64
}

// Key is the "provider/model" spelling the config text_chain elements use
// (SPEC-03 §3.1), so chain walking needs no second naming scheme.
func (m QuotaModel) Key() string { return m.Provider + "/" + m.Model }

// QuotaProvider is one provider's plan-credit line: the user-entered total.
// Zero means "not entered" and gates nothing — an unset total must not
// silence the product on a user's first day (same rule as the global budgets).
type QuotaProvider struct {
	Provider             string
	PlanCreditTotalMicro int64
}

// QuotaSpend is the consumption side of a judgment, bucketed to the caller's
// calendar: model day/month totals keyed by QuotaModel.Key(), and the
// provider-level plan-billing consumption keyed by provider name. Plan spend
// and pay-per-token spend are separate maps by construction: a plan credit
// draw must never darken a pay-per-token model's number, or AC#5's
// "tracked separately" fails silently.
type QuotaSpend struct {
	ModelDay     map[string]int64
	ModelMonth   map[string]int64
	ProviderPlan map[string]int64
}

// QuotaVerdict is one model's pre-dispatch reading.
type QuotaVerdict struct {
	Model string
	Level string // LevelOK | LevelWarn | LevelOver
	// Binding is "daily" | "monthly" | "plan-credit" | "" — which line the
	// user has to be told about. Like CapVerdict, "今天的钱" and "这个月的钱"
	// and "套餐 credit" expire on different clocks.
	Binding string
	// RemainingMicros is floored at 0: a negative remaining would promise a
	// dispatch into debt, which the pre-dispatch check exists to prevent.
	RemainingMicros int64
	// Estimated is true exactly for plan-credit remaining (manual total minus
	// bookkeeping, never a provider balance read).
	Estimated bool
	// AllowDispatch is the anti-negative gate: false only at 100% of a set
	// cap. Warn and over-with-no-total both still dispatch.
	AllowDispatch bool
	Notice        string
}

// QuotaParams carries the warn line. AlertBasisPoints is the SAME
// configuration point as [cost] alert_threshold (80% = 8000): SPEC-03 §3.1
// prescribes "80% warn, 100% hard-stop" for the quota layers without a second
// knob, so there is no invented per-model alert field.
type QuotaParams struct {
	AlertBasisPoints int
}

// NewQuotaParamsFrom adapts a BudgetParams' alert line so the global budget
// and the quota gates cannot drift two different 80%s apart.
func NewQuotaParamsFrom(b BudgetParams) QuotaParams {
	return QuotaParams{AlertBasisPoints: b.AlertBasisPoints}
}

// JudgeModel evaluates one model against its caps and (for billing=plan) its
// provider's credit line. It reads no store and takes no clock — the caller
// supplies the bucketed numbers, exactly like BudgetParams.Judge, which keeps
// the 80/100 boundaries testable without waiting for a rollover.
func (p QuotaParams) JudgeModel(m QuotaModel, prov *QuotaProvider, s QuotaSpend) QuotaVerdict {
	v := QuotaVerdict{Model: m.Key(), Level: LevelOK, AllowDispatch: true}
	remaining := int64(-1) // -1 = no cap seen at all
	basis := -1

	consider := func(kind string, spent, cap int64, estimated bool) {
		if cap <= 0 {
			return
		}
		b := basisPoints(spent, cap)
		rem := cap - spent
		if rem < 0 {
			rem = 0
		}
		// The tighter line wins, and a tie goes to the more specific layer
		// first-seen; order below is daily, monthly, plan-credit.
		if basis < 0 || b > basis {
			basis, v.Binding, v.Estimated = b, kind, estimated
			remaining = rem
		}
		switch {
		case spent >= cap:
			v.Level = LevelOver
		case b >= p.AlertBasisPoints:
			if v.Level != LevelOver {
				v.Level = LevelWarn
			}
		}
	}

	consider("daily", s.ModelDay[m.Key()], m.DailyMicros, false)
	consider("monthly", s.ModelMonth[m.Key()], m.MonthlyMicros, false)
	if m.Billing == BillingPlan && prov != nil {
		consider("plan-credit", s.ProviderPlan[prov.Provider], prov.PlanCreditTotalMicro, true)
	}

	if remaining >= 0 {
		v.RemainingMicros = remaining
	}
	if v.Level == LevelOver {
		v.AllowDispatch = false
		v.Notice = p.quotaStopNotice(m, v)
	} else if v.Level == LevelWarn {
		v.Notice = p.quotaWarnNotice(m, v, basis)
	}
	return v
}

// QuotaSkip records one chain entry the quota gate refused — the signal half
// of "over-quota model is SKIPPED via the text_chain" (票 44:34-35). Ticket 11
// owns the stepping; this is what it steps on.
type QuotaSkip struct {
	Model  string // "provider/model"
	Reason string // user-visible: why this entry was not dispatched
}

// SelectChain walks an ordered "provider/model" chain and returns the first
// entry that may still be dispatched plus the skipped ones in chain order.
// An entry not present in cards is treated as no-card (uncapped) and allowed:
// inventing a quota for a model nobody configured one would block real work.
func (p QuotaParams) SelectChain(chain []string, cards map[string]QuotaModel, providers map[string]QuotaProvider, s QuotaSpend) (chosen string, skips []QuotaSkip, verdicts map[string]QuotaVerdict) {
	verdicts = make(map[string]QuotaVerdict, len(chain))
	for _, key := range chain {
		m, ok := cards[key]
		if !ok {
			var split []string
			if i := strings.IndexByte(key, '/'); i >= 0 {
				split = []string{key[:i], key[i+1:]}
			} else {
				split = []string{key, ""}
			}
			m = QuotaModel{Provider: split[0], Model: split[1]}
		}
		var prov *QuotaProvider
		if q, ok := providers[m.Provider]; ok {
			prov = &q
		}
		v := p.JudgeModel(m, prov, s)
		verdicts[key] = v
		if !v.AllowDispatch {
			skips = append(skips, QuotaSkip{Model: key, Reason: v.Notice})
			continue
		}
		if chosen == "" {
			chosen = key
		}
	}
	return chosen, skips, verdicts
}

// QuotaLedger is the session-scoped bucket bookkeeping (see the persistence
// note at the top of the file). Day and month buckets are stamped with the
// same DayKey shape memory validates; on a bucket change the old numbers are
// dropped, because a new day/month resets per-model caps by definition —
// exactly the "resumable next day" behavior 票 44:22 asks of budgets.
//
// Concurrency (ticket 44 race leg, docs/evidence/s1/44-quota-ledger-race.md):
// ONE ledger serves the whole run. QuotaGate.Decide snapshots it before every
// dispatch (internal/agent/quota_44r4.go) while Book posts spend after each
// turn, and a root agent plus its subagents do both at the same time. So
// Spend and Post are concurrent entry points, and rollLocked is a WRITER of
// the bucket stamps plus a clear() of the maps whenever the calendar day or
// month moved — the torn stamp read CI's Test197Subagent… hit inside roll is
// exactly the access this mutex now serializes.
//
// mu guards EVERY field below, the two stamps included: deciding "did the
// bucket change" against a stamp another goroutine is rewriting is the same
// race as touching the maps, and the pair must be seen together or a reader
// can judge a fresh day against yesterday's buckets. The critical sections are
// pure in-memory map work (no I/O, no clock wait); Spend already copies all
// three maps into its QuotaSpend result, so the lock covers the work the
// reader would do anyway, and it is taken once per DISPATCH, not per model —
// a dispatch is a network round trip, so this cannot become the throughput
// ceiling. The alternative shape, an immutable snapshot replaced by an atomic
// pointer, would buy a lock-free Spend by moving that same copy onto every Post
// (three fresh maps per booking) and turning "the day rolled" into a
// compare-and-swap over a whole state struct: more machinery for no measured
// gain on a path where two goroutines hand each other the ball once per turn.
type QuotaLedger struct {
	mu       sync.Mutex       // guards every field below
	now      func() time.Time // injected clock; written once by NewQuotaLedger
	dayKey   string
	monthKey string
	modelDay map[string]int64
	modelMon map[string]int64
	planMon  map[string]int64
}

// NewQuotaLedger takes the same injected clock the cost meter's calendar uses
// (fixture clocks, AC#2 style; wall clock only at the composition root).
func NewQuotaLedger(now func() time.Time) *QuotaLedger {
	if now == nil {
		now = time.Now
	}
	return &QuotaLedger{
		now:      now,
		modelDay: map[string]int64{},
		modelMon: map[string]int64{},
		planMon:  map[string]int64{},
	}
}

// rollLocked advances the bucket stamps; a zero or changed day empties the day
// maps, a changed month empties the month maps. Called on every read and
// write so a long idle session cannot judge against yesterday's numbers.
// The caller holds l.mu: roll takes it, Post and Spend take it, so no reader
// ever sees a stamp its own map state was not rolled with.
func (l *QuotaLedger) rollLocked() {
	t := l.now()
	d := DayKey(t)
	m := d[:len("YYYY-MM")]
	if m != l.monthKey {
		l.monthKey = m
		clear(l.modelMon)
		clear(l.planMon)
	}
	if d != l.dayKey {
		l.dayKey = d
		clear(l.modelDay)
	}
}

// roll is rollLocked with the lock taken, for the one writer whose call site
// lives outside this file: SeedFrom (internal/observe/costmodel_44b.go) opens
// with l.roll() and then books rows into the maps itself. That seed is
// documented as a boot-time call made before any task books
// (QuotaGate.SeedBaseline), so it is not on the concurrent path this leg
// closes; its own map writes are a NAMED open item in
// docs/evidence/s1/44-quota-ledger-race.md, because protecting them needs one
// line in costmodel_44b.go, which is outside this leg's write face.
func (l *QuotaLedger) roll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rollLocked()
}

// Post books one dispatch's spend against the model's caps and, when the
// model bills as a plan, against its provider's credit line. Negative micros
// are refused rather than credited: a ledger that accepts negatives is one
// subtraction from fabricating remaining.
func (l *QuotaLedger) Post(m QuotaModel, micros int64) error {
	if micros < 0 {
		return fmt.Errorf("observe: quota ledger refuses negative spend (%d micros for %s)", micros, m.Key())
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rollLocked()
	if micros == 0 {
		return nil
	}
	l.modelDay[m.Key()] += micros
	l.modelMon[m.Key()] += micros
	if m.Billing == BillingPlan {
		l.planMon[m.Provider] += micros
	}
	return nil
}

// Spend snapshots the current buckets for JudgeModel/SelectChain. The copy is
// made under l.mu, which is what makes the returned QuotaSpend a consistent
// picture instead of three maps sampled while another task booked.
func (l *QuotaLedger) Spend() QuotaSpend {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rollLocked()
	cp := func(src map[string]int64) map[string]int64 {
		dst := make(map[string]int64, len(src))
		for k, v := range src {
			dst[k] = v
		}
		return dst
	}
	return QuotaSpend{
		ModelDay:     cp(l.modelDay),
		ModelMonth:   cp(l.modelMon),
		ProviderPlan: cp(l.planMon),
	}
}

// quotaStopNotice is the 100% line: dispatch is refused for THIS model and
// the chain moves on — SPEC-03 §3.1 "100% 硬停并自动走兜底链，降级事件
// 悬浮球+面板可见，不得静默". The message names which cap died and when it
// comes back, because the three layers recover on different clocks.
func (p QuotaParams) quotaStopNotice(m QuotaModel, v QuotaVerdict) string {
	word, recovered := quotaBindingWords(v.Binding)
	est := ""
	if v.Estimated {
		est = "（" + PlanCreditEstimatedMarker + "：手录总额减记账，非余额接口）"
	}
	return fmt.Sprintf(
		"配额已用尽：%s %s的额度到达上限%s，剩余 0。本次派发跳过该模型，按兜底链继续；%s后恢复。",
		m.Key(), word, est, recovered)
}

// quotaWarnNotice is the 80% line for one model's binding cap.
func (p QuotaParams) quotaWarnNotice(m QuotaModel, v QuotaVerdict, basis int) string {
	word, _ := quotaBindingWords(v.Binding)
	est := ""
	if v.Estimated {
		est = "（" + PlanCreditEstimatedMarker + "）"
	}
	return fmt.Sprintf(
		"配额预警：%s %s已用到 %d%%%s，剩余 %d（1e-6 货币单位）。到 100%% 会跳过该模型并走兜底链。",
		m.Key(), word, basis/basisPerPercent, est, v.RemainingMicros)
}

func quotaBindingWords(binding string) (word, recovered string) {
	switch binding {
	case "daily":
		return "当日配额", "明天 00:00（本地日历换日）"
	case "monthly":
		return "当月配额", "下月 1 日 00:00（本地日历换月）"
	case "plan-credit":
		return "套餐 credit", "录入新的 plan_credit_total_micro 或下月重置"
	default:
		return binding, "配额重置"
	}
}
