package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// C22 LoopGuard: the gradient brake (SPEC-05 §8, D31). Duplicate-call
// detection uses the frozen ladder [3,5,8]: each rung injects a graduated
// reminder into the conversation instead of hard-stopping, and the top rung
// drives the Stuck state with an explicit, user-visible message naming what
// repeated. Silent stopping is forbidden ("触发必须显式告知用户并进 Stuck 态，
// 不得静默停").
//
// The other brakes on the same object: per-tool timeoutMs (cooperative abort,
// structured TOOL_TIMEOUT), the per-task token budget (default 200k, scaled by
// context_window), and the 50-round last-resort floor. Any breach reports the
// consumption ("本次任务已用 X token / 约 ¥Y", C23 data).

// GuardConfig configures the brakes; zero values fall back to scaled defaults.
type GuardConfig struct {
	// RepeatThresholds is the ladder (default 3,5,8 from config/agent
	// .loop_guard.repeat_thresholds).
	RepeatThresholds []int
	// MaxRounds is the last-resort round floor (config agent.max_rounds, 50).
	MaxRounds int
	// TokenBudget overrides the scaled default when > 0 (config
	// agent.token_budget).
	TokenBudget int
	// PerToolTimeout is 0 = each tool's own default (config
	// agent.per_tool_timeout_ms).
	PerToolTimeout time.Duration
	// ToolConcurrency caps parallel executions (D38d ceiling 4).
	ToolConcurrency int
	// Budgets supplies the scaled defaults for token budget.
	Budgets Budgets
}

// Reminder is one graduated duplicate-call reminder (D43 row #19 side effect
// "agent.inject-reminder").
type Reminder struct {
	// Level is the ladder rung that fired (3, 5 or 8 by default).
	Level int
	// Repeat is how many consecutive identical calls were seen.
	Repeat int
	// Tool is the repeated tool name.
	Tool string
	// Text is the reminder injected as a user-role message.
	Text string
	// Stuck is true at the top rung: the loop must stop and say so.
	Stuck bool
}

// Brake describes which brake stopped the loop.
type Brake string

const (
	BrakeNone    Brake = ""
	BrakeRepeat  Brake = "repeat"
	BrakeToken   Brake = "token_budget"
	BrakeRounds  Brake = "rounds"
	BrakeTimeout Brake = "tool_timeout"
)

// Guard is the per-task C22 state machine. It is not safe for concurrent use;
// the loop owns it for the lifetime of one task.
type Guard struct {
	cfg       GuardConfig
	threshold []int

	lastSig   string
	repeat    int
	round     int
	tokensIn  int
	tokensOut int
	// CostMicros accumulates the C23 per-task money figure (1e-6 units of
	// Currency, computed from the model price card).
	CostMicros int64
	Currency   string
}

// NewGuard builds a guard, filling the frozen defaults.
func NewGuard(cfg GuardConfig) *Guard {
	th := make([]int, 0, len(cfg.RepeatThresholds))
	for _, v := range cfg.RepeatThresholds {
		if v > 0 {
			th = append(th, v)
		}
	}
	if len(th) == 0 {
		th = []int{3, 5, 8} // C22 frozen ladder
	}
	insertionSort(th)
	cfg.RepeatThresholds = th
	if cfg.MaxRounds <= 0 {
		cfg.MaxRounds = 50 // D15(4) last-resort round floor
	}
	if cfg.TokenBudget <= 0 {
		cfg.TokenBudget = cfg.Budgets.TokenBudget // 200k scaled by context_window
	}
	if cfg.ToolConcurrency <= 0 || cfg.ToolConcurrency > MaxToolConcurrency {
		cfg.ToolConcurrency = MaxToolConcurrency // D38d ceiling
	}
	return &Guard{cfg: cfg, threshold: th}
}

// Thresholds returns the active ladder (ascending).
func (g *Guard) Thresholds() []int {
	out := make([]int, len(g.threshold))
	copy(out, g.threshold)
	return out
}

// TopThreshold is the rung that triggers Stuck.
func (g *Guard) TopThreshold() int { return g.threshold[len(g.threshold)-1] }

// PerToolTimeout resolves the C22 per-tool timeout for one call: the config
// override wins, then the scaled default (0 = each tool's own default, which
// the host bridge supplies from ticket 20).
func (g *Guard) PerToolTimeout() time.Duration { return g.cfg.PerToolTimeout }

// Concurrency is the effective D38d tool-concurrency ceiling.
func (g *Guard) Concurrency() int { return g.cfg.ToolConcurrency }

// TokensUsed totals the task's consumed tokens (C23 data).
func (g *Guard) TokensUsed() int { return g.tokensIn + g.tokensOut }

// TokenBudget is the effective per-task budget after scaling/config.
func (g *Guard) TokenBudget() int { return g.cfg.TokenBudget }

// Rounds is the completed round count.
func (g *Guard) Rounds() int { return g.round }

// NextRound advances the round counter and reports the 50-round floor.
func (g *Guard) NextRound() (round int, floorHit bool) {
	g.round++
	return g.round, g.round > g.cfg.MaxRounds
}

// AddUsage folds one Usage event into the task totals and reports whether the
// token budget is now exhausted.
func (g *Guard) AddUsage(u llm.Usage) bool {
	g.tokensIn += u.InputTokens
	g.tokensOut += u.OutputTokens
	return g.TokenBudget() > 0 && g.TokensUsed() >= g.TokenBudget()
}

// AddCost accumulates a C23 cost delta in 1e-6 currency units.
func (g *Guard) AddCost(micros int64) {
	if micros > 0 {
		g.CostMicros += micros
	}
}

// ObserveTurn is the duplicate-call ladder. It must be called once per model
// turn that produced tool calls, with the calls in wire order. The rule is the
// D43 row #19 wording: the SAME tool with the SAME arguments called
// consecutively ("已连续调用 N 次，参数相同").
func (g *Guard) ObserveTurn(calls []llm.ToolCall) (Reminder, bool) {
	sig := turnSignature(calls)
	if sig == "" {
		g.lastSig, g.repeat = "", 0
		return Reminder{}, false
	}
	if sig == g.lastSig {
		g.repeat++
	} else {
		g.lastSig = sig
		g.repeat = 1
	}
	level := 0
	for _, th := range g.threshold {
		if g.repeat == th {
			level = th
			break
		}
	}
	if level == 0 {
		return Reminder{}, false
	}
	name := calls[0].Name
	r := Reminder{
		Level:  level,
		Repeat: g.repeat,
		Tool:   name,
		Text: fmt.Sprintf(
			"提醒：你已经连续 %d 次以完全相同的参数调用 %s。请换一个做法，或直接把当前结果告诉用户；继续重复同样的调用不会改变结果。",
			g.repeat, name),
		Stuck: level >= g.TopThreshold(),
	}
	return r, true
}

// ResetRepeats clears the duplicate counter (called when the model changed its
// behaviour, e.g. after a tool result arrived).
func (g *Guard) ResetRepeats() { g.lastSig, g.repeat = "", 0 }

// Consumption renders the C22 "what did this cost me" line. The money part is
// only shown when a price card produced a figure: currency conversion is not
// this module's invention (see the ticket report's open question).
func (g *Guard) Consumption() string {
	if g.CostMicros > 0 {
		return fmt.Sprintf("本次任务已用 %d token / 约 %s（估算 %s）",
			g.TokensUsed(), formatMicros(g.CostMicros), g.currency())
	}
	return fmt.Sprintf("本次任务已用 %d token", g.TokensUsed())
}

func (g *Guard) currency() string {
	if g.Currency == "" {
		return "micro"
	}
	return g.Currency
}

func formatMicros(micros int64) string {
	whole := micros / 1_000_000
	frac := micros % 1_000_000
	return fmt.Sprintf("%d.%06d", whole, frac)
}

// turnSignature is the duplicate key of a turn: the ordered (name, args)
// pairs. JSON arguments are canonicalized with encoding/json (decoded, then
// re-encoded, which sorts object keys and drops insignificant whitespace), so
// equivalent objects written with a different key order or spacing count as
// repeats. Arguments that are not a single JSON value fall back to the trimmed
// raw string.
func turnSignature(calls []llm.ToolCall) string {
	if len(calls) == 0 {
		return ""
	}
	var b strings.Builder
	for _, c := range calls {
		fmt.Fprintf(&b, "%s|%s;", c.Name, canonicalArgs(c.Args))
	}
	return b.String()
}

// canonicalArgs is the canonicalization turnSignature promises. Decoding with
// json.Number keeps number literals verbatim, so the rewrite cannot lose
// precision on a big id; re-encoding is what sorts the keys.
func canonicalArgs(raw json.RawMessage) string {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return ""
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return string(raw) // not JSON: exact-text comparison is all there is
	}
	if dec.More() {
		return string(raw) // trailing junk after the first value: stay literal
	}
	enc, err := json.Marshal(v)
	if err != nil {
		return string(raw)
	}
	return string(enc)
}

// ErrorClassOfTurnError maps a tool/host failure onto the D37 class used for
// the task_log and tool_call rows.
func ErrorClassOfTurnError(err error) string {
	c, ok := observe.ClassOf(err)
	if !ok {
		return string(observe.ClassInternal)
	}
	return string(c)
}

// insertionSort keeps the ladder ascending without pulling in "sort" for a
// three-element slice.
func insertionSort(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j-1] > a[j]; j-- {
			a[j-1], a[j] = a[j], a[j-1]
		}
	}
}
