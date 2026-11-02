package observe

import (
	"fmt"
	"math"
)

// Money budgets (C23 [cost], ticket 44 AC#3).
//
// The unit is micro-currency everywhere. The comparison is integer: a spend is
// judged against a cap in BASIS POINTS (1/100 of a percent), so "80%" is 8000
// and no float ever multiplies a money amount. alert_threshold arrives as a
// ratio because config's [cost] field is a float64 in the frozen D36 schema
// (internal/config/schema.go:561) — converting that ratio once, here, is the
// only place a float is touched, and the money numbers never see one.
//
// Two distinct layers, on purpose (44-costmeter-c23.md:24, quoted below): the
// single-task token budget (default 200k, anti-loop, internal/agent's Guard)
// and THIS money budget. They are enforced independently and neither reads the
// other's number.

// Basis points in one whole percent.
const basisPerPercent = 100

// Levels a budget BudgetVerdict can be at.
const (
	LevelOK   = "ok"
	LevelWarn = "warn"
	LevelOver = "over"
)

// BudgetParams is the meter's view of the [cost] section.
type BudgetParams struct {
	// DailyMicros / MonthlyMicros are caps in micro-currency units; 0 = no cap
	// (a 0 cap is "uncapped", NOT "block everything": an unset budget must not
	// silence the product on a user's first day).
	DailyMicros   int64
	MonthlyMicros int64
	// AlertBasisPoints is the warn line (8000 = 80%).
	AlertBasisPoints int
	// PauseAtLimit is [cost] over_budget: true = pause new tasks at 100%,
	// false = keep warning (the "warn" alternative). It never touches a
	// running task either way.
	PauseAtLimit bool
}

// BudgetFrom adapts the [cost] fields. observe cannot import config (config
// imports observe), so the host calls this with the section's four leaves. An
// alert ratio outside [0,1] is a config bug and is refused rather than clamped:
// clamping would turn a typo into a silently different policy.
func BudgetFrom(dailyMicros, monthlyMicros int64, alertRatio float64, overBudget string) (BudgetParams, error) {
	if dailyMicros < 0 || monthlyMicros < 0 {
		return BudgetParams{}, fmt.Errorf("observe: [cost] budgets must be >= 0 (got daily=%d monthly=%d)", dailyMicros, monthlyMicros)
	}
	if math.IsNaN(alertRatio) || math.IsInf(alertRatio, 0) || alertRatio < 0 || alertRatio > 1 {
		return BudgetParams{}, fmt.Errorf("observe: [cost] alert_threshold %v out of range [0, 1]", alertRatio)
	}
	switch overBudget {
	case "pause", "warn":
	default:
		return BudgetParams{}, fmt.Errorf("observe: [cost] over_budget %q must be pause or warn", overBudget)
	}
	return BudgetParams{
		DailyMicros:      dailyMicros,
		MonthlyMicros:    monthlyMicros,
		AlertBasisPoints: int(math.Round(alertRatio * 100 * basisPerPercent)),
		PauseAtLimit:     overBudget == "pause",
	}, nil
}

// CapVerdict is one cap's reading (daily or monthly), kept separate because the
// user has to be told WHICH budget ran out: "今天的钱用完了" and "这个月的钱用完了"
// are different facts with different recovery days.
type CapVerdict struct {
	Kind      string // "daily" | "monthly"
	Spent     int64
	Cap       int64
	Basis     int // floor(spend/cap * 10000), saturating at MaxInt for a huge spend
	Level     string
	Recovered string // when the limit lifts, in the user's calendar terms
}

// BudgetVerdict is the answer the meter gives before a new task starts.
type BudgetVerdict struct {
	// AllowNewTask is the gate: false only when the limit is reached AND
	// over_budget=pause. In warn mode it stays true, with a Notice.
	AllowNewTask bool
	Level        string
	// Binding is the tighter of the two caps (the one worth showing).
	Binding CapVerdict
	// Caps lists every cap that is set, so a panel can show both lines instead
	// of only the one that bites.
	Caps []CapVerdict
	// Notice is user-visible text; empty at LevelOK. A rejection message is
	// explicit by requirement (AC#3: "rejected with explicit message").
	Notice string
}

// Judge evaluates the budgets against the day and month-to-date aggregates.
// It reads no store and takes no clock: the caller supplies the numbers it
// already has, which is what makes the 80%/100% boundaries testable without
// waiting for a rollover.
func (p BudgetParams) Judge(day, month DayTotal) BudgetVerdict {
	v := BudgetVerdict{AllowNewTask: true, Level: LevelOK}
	if c, ok := p.judgeOne("daily", day.Micros, p.DailyMicros, "明天 00:00（本地日历换日）"); ok {
		v.Caps = append(v.Caps, c)
		v = merge(v, c)
	}
	if c, ok := p.judgeOne("monthly", month.Micros, p.MonthlyMicros, "下月 1 日 00:00（本地日历换月）"); ok {
		v.Caps = append(v.Caps, c)
		v = merge(v, c)
	}
	switch v.Level {
	case LevelOver:
		if p.PauseAtLimit {
			v.AllowNewTask = false
			v.Notice = p.pauseNotice(v.Binding)
		} else {
			v.Notice = p.warnNotice(v.Binding, true)
		}
	case LevelWarn:
		v.Notice = p.warnNotice(v.Binding, false)
	}
	return v
}

// merge keeps the tighter cap as the binding one (and the worst level). The
// first cap set always becomes the binding one: without a baseline, "0 basis"
// would silently win over a cap that is actually at 5%.
func merge(v BudgetVerdict, c CapVerdict) BudgetVerdict {
	if len(v.Caps) <= 1 || c.Basis >= v.Binding.Basis {
		v.Binding = c
	}
	if levelRank(c.Level) > levelRank(v.Level) {
		v.Level = c.Level
	}
	return v
}

func levelRank(l string) int {
	switch l {
	case LevelOver:
		return 2
	case LevelWarn:
		return 1
	default:
		return 0
	}
}

// judgeOne returns ok=false when the cap is unset (no budget to judge).
func (p BudgetParams) judgeOne(kind string, spent, cap int64, recovered string) (CapVerdict, bool) {
	if cap <= 0 {
		return CapVerdict{}, false
	}
	c := CapVerdict{Kind: kind, Spent: spent, Cap: cap, Basis: basisPoints(spent, cap), Recovered: recovered}
	switch {
	case spent >= cap:
		c.Level = LevelOver
	case c.Basis >= p.AlertBasisPoints:
		c.Level = LevelWarn
	default:
		c.Level = LevelOK
	}
	return c, true
}

// basisPoints is floor(spent*10000/cap) with saturation instead of wraparound:
// an overflowed division would read as a tiny percentage and let a task start
// past a budget that is obviously gone.
func basisPoints(spent, cap int64) int {
	if cap <= 0 {
		return 0
	}
	if spent <= 0 {
		return 0
	}
	if spent > math.MaxInt64/10000 {
		return math.MaxInt
	}
	b := spent * 10000 / cap
	if b > math.MaxInt {
		return math.MaxInt
	}
	return int(b)
}

func (p BudgetParams) pauseNotice(c CapVerdict) string {
	return fmt.Sprintf(
		"费用预算已用尽：%s已用 %d / 上限 %d（1e-6 货币单位，%d%%）。为避免费用继续增长，Wisp 已暂停新任务；运行中的任务不受影响，会跑完它自己。%s后自动恢复，也可以调高 [cost] 的预算或把 over_budget 设为 warn。",
		capWord(c.Kind), c.Spent, c.Cap, c.Basis/100, c.Recovered)
}

func (p BudgetParams) warnNotice(c CapVerdict, at100 bool) string {
	if at100 {
		return fmt.Sprintf(
			"费用预算已达上限：%s已用 %d / 上限 %d（1e-6 货币单位，%d%%）。over_budget=warn，所以新任务照常开始；运行中的任务不受影响。%s后恢复。",
			capWord(c.Kind), c.Spent, c.Cap, c.Basis/100, c.Recovered)
	}
	return fmt.Sprintf(
		"费用已用到预算的 %d%%：%s当前 %d / 上限 %d（1e-6 货币单位）。到 100%% 会%s。",
		c.Basis/100, capWord(c.Kind), c.Spent, c.Cap,
		p.overBudgetWord())
}

// overBudgetWord states what 100% will do, in the mode the user configured —
// the warn line is worth reading before the limit bites, and it has to say the
// right thing in both modes.
func (p BudgetParams) overBudgetWord() string {
	if p.PauseAtLimit {
		return "暂停新任务"
	}
	return "只告警、不拦任务"
}

func capWord(kind string) string {
	switch kind {
	case "daily":
		return "今天（本地日历）"
	case "monthly":
		return "本月（本地日历）"
	default:
		return kind
	}
}

// Percent renders a basis-point value as an integer percent for display
// (12345 -> "123%"). Floor, never round-up: telling a user "100%" when the
// spend is 99.5% would announce a limit that has not been reached.
func Percent(basis int) int { return basis / basisPerPercent }
