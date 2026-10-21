package agent

import (
	"context"
	"fmt"
	"log/slog"
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
//
// Trace placement (ticket 139 AC#2): a pass that actually folded history is
// booked as ONE slog.Info record ("agent: history compressed") emitted right
// here, on the success edge, so it carries the counts and never the text.
// Field set (ticket 153 AC#2 appends the last key): tokens_before,
// tokens_after, threshold, msgs_before, msgs_after, compressed_msgs,
// kept_raw_rounds, history_changed, and "task" - the owning task id, present
// only when the caller handed this pass one (see withTraceTask below).
// Three things follow from putting it in the compressor rather than at a call
// site: the record exists even for a caller that forgets to log (the
// Warm-window hook DEFERRED(D28-1) installs will get it for free); it pairs
// with the caller's failure Warn ("agent: history compression failed") so both
// outcomes of one call answer to the search prefix "agent: history compress";
// and it cannot be emitted for a pass that folded nothing, because the fold is
// decided inside this function (rep.Ran), not inferred afterwards.
//
// Privacy: [privacy] keep_transcript forbids putting history text on this
// line, so every attribute is a count or a flag.
// CompressionState is the outcome of one compression pass (ticket 358 AC#1,
// the shape the orchestrator ruled: a tri-state enum, NOT a companion bool).
//
// Why an enum and not a second bool on Ran: "did not trigger" and "triggered
// but could not fold" are different answers to different questions, and two
// booleans can express the illegal pair Ran=true && Skipped=true. A single
// enum makes exactly one state readable at a time. The vocabulary is the one
// mainstream compaction uses.
type CompressionState string

const (
	// CompressionSkipped: this pass was not an attempt. The trigger never said
	// the history was over its budget, so the fold loop was never entered.
	CompressionSkipped CompressionState = "skipped"
	// CompressionAttempted: the pass WAS an attempt - the trigger said the
	// history is over budget, or the fold loop was entered - and it folded
	// nothing. This is the state ticket 358 exists for: "该压但压不动" used to
	// be indistinguishable from "没触发" on every surface Wisp has.
	//
	// Two shapes reach it, and the report's own numbers tell them apart:
	//   - the fold loop was entered and broke at the KeepRawRounds floor
	//     (TokensBefore > Threshold, KeptRawRounds == KeepRawRounds);
	//   - the provider-reported input tokens cleared the threshold while this
	//     pass's own estimate says the history is under the fold budget
	//     (TokensBefore <= Threshold): the trigger fired, the fold had
	//     nothing of its own to take.
	CompressionAttempted CompressionState = "attempted"
	// CompressionCompacted: at least one round was folded into the summary
	// block. Ran is the compatibility mirror of this state.
	CompressionCompacted CompressionState = "compacted"
)

// Valid reports whether s is one of the three states.
func (s CompressionState) Valid() bool {
	switch s {
	case CompressionSkipped, CompressionAttempted, CompressionCompacted:
		return true
	}
	return false
}

// String makes the state readable in a record or a failure message.
func (s CompressionState) String() string {
	if s.Valid() {
		return string(s)
	}
	return "unknown(" + string(s) + ")"
}

// TriggerSource names the numerator a trigger decision was made on (ticket 358
// AC#5). The two are never mixed silently: the state report carries the source
// and so does every trace record, because "we compressed because the provider
// said 40k" and "we compressed because our own estimate said 40k" are
// different facts about one run.
type TriggerSource string

const (
	// TriggerEstimate: the local approximation (ApproxTokensOf). This is what
	// Need() has always used and what a first turn has to use - nothing has
	// been observed yet.
	TriggerEstimate TriggerSource = "estimate"
	// TriggerProviderUsage: the input-token count the provider itself reported
	// for the previous response (llm.Usage.InputTokens). Mainstream harnesses
	// trigger compaction on this number, and it is the only one that knows
	// what the request actually cost on the wire.
	TriggerProviderUsage TriggerSource = "provider_usage"
)

// Valid reports whether s names a numerator.
func (s TriggerSource) Valid() bool {
	return s == TriggerEstimate || s == TriggerProviderUsage
}

// TriggerReading is one trigger decision WITH its basis (ticket 358 AC#5).
// Tokens is the numerator that was compared, Threshold the frozen budget it was
// compared against; nothing else decides the trigger.
type TriggerReading struct {
	Need      bool
	Tokens    int
	Threshold int
	Source    TriggerSource
	// Observed is the provider-reported input-token count handed in (0 when
	// the caller had none: first turn, or a provider that reports nothing).
	Observed int
}

// skip-reason values carried by the "attempted" trace record. They name WHY an
// attempt folded nothing, so a reader does not have to infer it.
const (
	// skipReasonKeepRawFloor: the fold loop ran and stopped at the D15(4) floor
	// - the last KeepRawRounds rounds must stay verbatim, so there was nothing
	// left to fold.
	skipReasonKeepRawFloor = "keep_raw_rounds_floor"
	// skipReasonUnderFoldBudget: the trigger fired (on the provider's number)
	// while this pass's own estimate puts the history at or under the fold
	// budget, so the fold loop never took a round.
	skipReasonUnderFoldBudget = "under_fold_budget"
)

type CompressionReport struct {
	// Ran is the compatibility mirror of State == CompressionCompacted. State
	// is authoritative: every report Compress returns satisfies
	// Ran == (State == CompressionCompacted) (pinned by
	// TestCompressionStateRanIsTheMirrorOfCompacted).
	Ran bool
	// State is the three-state outcome (ticket 358 AC#1). Every report Compress
	// returns carries one of the three valid states. The empty string is NOT a
	// state: it is what a CompressionReport nobody filled in holds, i.e. "no
	// pass ran at all" (the loop's Result.Compression for a run that never
	// reached the compressor). Keeping that outside the enum is deliberate - an
	// unfilled struct must not read as "we looked and decided not to".
	State CompressionState
	// TriggerSource and TriggerTokens are the basis this pass was triggered on
	// (ticket 358 AC#5): never defaulted silently, always named.
	TriggerSource  TriggerSource
	TriggerTokens  int
	CompressedMsgs int
	KeptRawRounds  int
	PreservedIDs   []string
	// RawIDsInsideSummary lists tool_call ids that exist only inside the
	// compressed block (proof the ledger survived compression).
	RawIDsInsideSummary []string
	TokensBefore        int
	TokensAfter         int
	// HistoryChanged reports whether this pass actually moved the history:
	// it is the same expression the trace records as "history_changed"
	// (TokensAfter != TokensBefore || the message count changed), NOT a copy
	// of Ran. The two can disagree - a pass can fold rounds whose summary is
	// the same size, leaving Ran true and the bytes identical - and a caller
	// that asked "did the history move" must be able to hear that.
	//
	// It is a read-only reading of a completed pass: Compress sets it, nothing
	// in this package consults it. It exists so the answer stops living only
	// on the log line (ticket 355 AC#2).
	HistoryChanged bool
}

// Compressor applies the D15(4) policy against a scaled budget set.
type Compressor struct {
	b   Budgets
	sum Summarizer
	lg  *slog.Logger
}

// CompressorOpt tweaks a Compressor at construction.
type CompressorOpt func(*Compressor)

// WithLogger sends the compression trace to lg instead of the process default
// (nil lg = slog.Default(), which is where the host's persistent JSONL sink
// lives). The two-argument call sites keep compiling: the option list is open
// by design, so a caller with no logger changes nothing.
func WithLogger(lg *slog.Logger) CompressorOpt {
	return func(c *Compressor) { c.lg = lg }
}

// NewCompressor builds a compressor; sum may be nil (structural trim).
func NewCompressor(b Budgets, sum Summarizer, opts ...CompressorOpt) *Compressor {
	c := &Compressor{b: b, sum: sum}
	for _, o := range opts {
		o(c)
	}
	return c
}

// log is the compressor's trace sink: the injected one if a caller provided it,
// else the process default. Never nil, so call sites need no guard.
func (c *Compressor) log() *slog.Logger {
	if c.lg != nil {
		return c.lg
	}
	return slog.Default()
}

// TotalTokens estimates the history cost.
func (c *Compressor) TotalTokens(hist []llm.Message) int {
	n := 0
	for _, m := range hist {
		n += ApproxTokensOf(m.Content)
	}
	return n
}

// Need reports whether the history is over the compression threshold, measured
// with the local estimate. Its threshold semantics are the frozen D15 ones and
// are NOT touched by ticket 358; it is the estimate-only reading, kept because
// callers that have no provider measurement (and the tests that pin the budget)
// still need it.
//
// The loop does not call this any more: it calls Trigger, which prefers the
// provider-reported input tokens and only falls back to exactly this estimate
// when the provider reported nothing (ticket 358 AC#5).
func (c *Compressor) Need(hist []llm.Message) bool {
	return c.TotalTokens(hist) > c.b.HistoryCompressTokens
}

// Trigger decides whether to compress and records which numerator decided it
// (ticket 358 AC#5).
//
// observed is the input-token count the provider reported for the PREVIOUS
// response (llm.Usage.InputTokens), or 0 when there is none - the first turn of
// a task, or a provider that reports no usage. The two numerators are never
// mixed silently: whichever one produced the decision is named in the reading
// and travels into the compression report and its trace.
//
// The threshold is the same frozen scaled HistoryCompressTokens in both
// branches: ticket 358 changes the numerator, never the threshold.
func (c *Compressor) Trigger(hist []llm.Message, observed int) TriggerReading {
	th := c.b.HistoryCompressTokens
	if observed > 0 {
		return TriggerReading{
			Need:      observed > th,
			Tokens:    observed,
			Threshold: th,
			Source:    TriggerProviderUsage,
			Observed:  observed,
		}
	}
	est := c.TotalTokens(hist)
	return TriggerReading{
		Need:      est > th,
		Tokens:    est,
		Threshold: th,
		Source:    TriggerEstimate,
		Observed:  observed,
	}
}

// ---------------------------------------------------------------------------
// trace attribution (ticket 153 AC#2)
//
// The record booked on the success edge cannot see a task id on its own: the
// Compressor is built ONCE per Loop in agent.New (loop.go:239), the id is
// minted later and only inside Loop.run (loop.go:322 / :333), and neither the
// receiver nor Compress's parameter list carries it. Nor is it on the context
// the pass already receives - observe.NewRootFrom derives a plain cancelable
// child and keeps the id on the Root struct (goroutine.go:103-106), and there
// is no context lookup for it anywhere in the package.
//
// withTraceTask is the smallest channel that closes that gap without moving the
// record out of the compressor (139 AC#2's placement, which this ticket does
// not reopen) and without adding mutable per-task state to a shared object.
//
// What it does NOT buy, stated plainly: a caller that does not tag its context
// gets a record with no "task" key at all. Absence is the designed outcome -
// the compressor never fills in "" or "unknown", because a placeholder reads as
// attributed and turns "cannot attribute" into "attributed to nobody". Pinned by
// TestCompressionTraceNeverInventsATaskID. Today the only non-test caller is
// loop.go's synchronous fallback, which tags every call; the Warm-window hook
// this call is supposed to move to (the D28-1 deferral flagged at the top of
// this file) has to tag too.

// traceTaskKey is the unexported context key the owning task id rides on.
type traceTaskKey struct{}

// withTraceTask tags ctx with the task a compression pass belongs to. An empty
// (or whitespace-only) id is not a task, so it is dropped rather than carried.
func withTraceTask(ctx context.Context, taskID string) context.Context {
	if strings.TrimSpace(taskID) == "" {
		return ctx
	}
	return context.WithValue(ctx, traceTaskKey{}, taskID)
}

// traceTaskID is the read side: "" when the caller tagged nothing.
func traceTaskID(ctx context.Context) string {
	s, _ := ctx.Value(traceTaskKey{}).(string)
	return s
}

// round is one conversational round (a user turn plus what followed it).
type round struct {
	msgs      []llm.Message
	isSummary bool
}

// CompressOption tweaks a single Compress call. The option list is open by
// design, so existing callers compile unchanged.
type CompressOption func(*compressCall)

// compressCall carries the per-call inputs that are not part of the frozen
// signature.
type compressCall struct {
	// trigger is the reading that made the loop enter the compressor, or the
	// zero TriggerReading when the caller did not compute one. It is what turns
	// a pass that folded nothing from "we do not know why" into "attempted,
	// and here is the numerator that said so".
	trigger TriggerReading
	// hasTrigger distinguishes "the caller computed a reading" from "the zero
	// reading", which matter differently: with no reading the pass reports its
	// own estimate as the basis.
	hasTrigger bool
}

// WithTrigger hands Compress the trigger reading the caller already computed
// (ticket 358 AC#5). Compress does not re-decide the trigger - the loop owns
// that decision - but a pass that folds nothing must still be able to say
// whether it was an attempt, and on which numerator.
func WithTrigger(r TriggerReading) CompressOption {
	return func(cc *compressCall) {
		cc.trigger = r
		cc.hasTrigger = true
	}
}

// Compress folds the oldest rounds until the history fits the threshold or
// only KeepRawRounds raw rounds remain. It never mutates the input slices'
// content and always returns a freshly built history.
//
// Ticket 358 AC#1: the returned report carries a three-state outcome. A pass
// that was triggered and folded nothing is CompressionAttempted, and it leaves
// a trace of its own ("agent: history compress skipped") - the state that used
// to be invisible everywhere.
func (c *Compressor) Compress(ctx context.Context, hist []llm.Message, opts ...CompressOption) ([]llm.Message, CompressionReport, error) {
	var call compressCall
	for _, o := range opts {
		o(&call)
	}

	rep := CompressionReport{
		TokensBefore:  c.TotalTokens(hist),
		TriggerSource: TriggerEstimate,
		TriggerTokens: c.TotalTokens(hist),
	}
	if call.hasTrigger {
		rep.TriggerSource = call.trigger.Source
		rep.TriggerTokens = call.trigger.Tokens
	}
	rounds := groupRounds(hist)
	// entered is the fact the tri-state keys on: did this pass actually try to
	// fold? Two things can make it true - the caller's trigger reading, or this
	// pass's own estimate clearing the fold budget (a bare Compress call with
	// no reading is still an attempt when the history is over).
	entered := call.hasTrigger && call.trigger.Need
	if c.totalOf(rounds) > c.b.HistoryCompressTokens {
		entered = true
	}
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
		if idx := summaryRoundIndex(rounds); idx >= 0 {
			victims = append(victims, rounds[idx].msgs...)
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

	// Ticket 358 AC#1: exactly one of the three states, decided here and
	// nowhere else. Ran stays the mirror of the compacted state so every
	// existing reader keeps its meaning.
	switch {
	case rep.Ran:
		rep.State = CompressionCompacted
	case entered:
		rep.State = CompressionAttempted
	default:
		rep.State = CompressionSkipped
	}
	rep.Ran = rep.State == CompressionCompacted

	if rep.Ran {
		// The one trace a successful pass leaves (AC#2): counts only, never the
		// folded text. "history_changed" is recomputed from what this call
		// actually returned instead of echoing rep.Ran, so it can contradict
		// the fold flag and a reader is never told the history moved when the
		// bytes coming back are the same size.
		//
		// Ticket 355 AC#2: the same value is stored on the report instead of
		// living only on this record. It is computed ONCE here - the expression
		// is unchanged - and both readers take it from rep, so the log line and
		// the caller's reading cannot drift apart.
		rep.HistoryChanged = rep.TokensAfter != rep.TokensBefore || len(out) != len(hist)
		attrs := []any{
			"state", rep.State.String(),
			"tokens_before", rep.TokensBefore,
			"tokens_after", rep.TokensAfter,
			"threshold", c.b.HistoryCompressTokens,
			"msgs_before", len(hist),
			"msgs_after", len(out),
			"compressed_msgs", rep.CompressedMsgs,
			"kept_raw_rounds", rep.KeptRawRounds,
			"history_changed", rep.HistoryChanged,
		}
		// Ticket 358 AC#5: which numerator triggered this pass. Named on every
		// record - a reader must never have to guess whether "token" meant the
		// provider's number or our own estimate.
		attrs = append(attrs,
			"trigger_source", string(rep.TriggerSource),
			"trigger_tokens", rep.TriggerTokens)
		// Ticket 153 AC#2: whose pass this was. The key is "task", the name the
		// loop already uses for the same value (loop.go's task_log Warns), and
		// the value is the id the Result and every Event of that run carry, so
		// one grep joins the line to the task. It is a correlation id, not
		// history content, so [privacy] keep_transcript does not reach it.
		if task := traceTaskID(ctx); task != "" {
			attrs = append(attrs, "task", task)
		}
		c.log().Info("agent: history compressed", attrs...)
	} else if rep.State == CompressionAttempted {
		// Ticket 358 AC#1: the shape this whole ticket is about. The history
		// was over its budget on the numerator that counts, and the pass folded
		// nothing. Without this record the run is indistinguishable from one
		// whose history was never over budget - which is exactly the hole the
		// ticket names. Same search prefix as its two siblings ("agent:
		// history compress"), so one grep answers all three outcomes.
		reason := skipReasonUnderFoldBudget
		if rep.TokensBefore > c.b.HistoryCompressTokens {
			reason = skipReasonKeepRawFloor
		}
		attrs := []any{
			"state", string(CompressionAttempted),
			"reason", reason,
			"tokens_before", rep.TokensBefore,
			"tokens_after", rep.TokensAfter,
			"threshold", c.b.HistoryCompressTokens,
			"msgs_before", len(hist),
			"msgs_after", len(out),
			"compressed_msgs", 0,
			"kept_raw_rounds", rep.KeptRawRounds,
			"history_changed", false,
			"trigger_source", string(rep.TriggerSource),
			"trigger_tokens", rep.TriggerTokens,
		}
		if task := traceTaskID(ctx); task != "" {
			attrs = append(attrs, "task", task)
		}
		c.log().Info("agent: history compress skipped", attrs...)
	}
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
