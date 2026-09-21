package agent

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// The single-task ReAct loop (SPEC-05 §2). Skeleton, in order:
//
//	input -> D11(1) control layer (local regex, no LLM)
//	      -> context assembly (D39 sections, cache-prefix order)
//	      -> C5 provider stream -> C6 event consumption
//	      -> tool calls (risk pass-through until ticket 21, D38d concurrency
//	         ceiling 4, C22 per-tool timeout, D15(3) result spill)
//	      -> results fed back -> next round
//
// Termination: no tool call plus finished text / cancellation / C22 brake /
// token budget / the 50-round last-resort floor. stopReason=max_tokens makes
// failToolCallsFromTruncatedMessage fail every tool call of that message
// (D21: truncated arguments are never executed).
//
// The loop owns no provider protocol code (llm), no tool implementations
// (tools), and no approval decisions (agent/approval, ticket 21).

// Status is the terminal state of one task.
type Status string

const (
	// StatusCompleted: the model finished with text and no open tool call.
	StatusCompleted Status = "completed"
	// StatusStuck: a C22 brake fired; a visible message was published.
	StatusStuck Status = "stuck"
	// StatusCancelled: the task ctx was cancelled.
	StatusCancelled Status = "cancelled"
	// StatusFailed: a classified provider/host failure.
	StatusFailed Status = "failed"
	// StatusTruncated: the response hit the output ceiling; nothing of that
	// message's tool calls was executed (D21).
	StatusTruncated Status = "truncated"
	// StatusControl: the utterance was a D11(1) control word; no LLM call was
	// made and no task_log row was opened.
	StatusControl Status = "control"
)

// Valid reports whether s is one of the loop's terminal states.
func (s Status) Valid() bool {
	switch s {
	case StatusCompleted, StatusStuck, StatusCancelled, StatusFailed,
		StatusTruncated, StatusControl:
		return true
	}
	return false
}

// ToolResultLog is one executed (or deliberately failed) call, kept for the
// result surface, the tests and the tool_call rows it mirrors.
type ToolResultLog struct {
	CallID     string
	Name       string
	Outcome    string // memory tool_call.outcome vocabulary
	ErrorClass string // D37 class or ""
	Text       string
	Spilled    bool
	Artifact   string
	Truncated  bool
	Rejected   bool // gate/pass-through refused before execution
}

// Result is the outcome of one Run.
type Result struct {
	TaskID string
	Status Status
	// Text is the final assistant text (the reply to present/route, D10).
	Text string
	// Message is the user-visible status line for non-clean endings (Stuck,
	// truncation, failure, control ack). Never empty when Status != completed.
	Message    string
	Rounds     int
	Usage      llm.Usage
	CostMicros int64
	Currency   string
	ToolCalls  int
	ToolLog    []ToolResultLog
	Brake      Brake
	// ReminderLevels records the C22 ladder rungs that fired (e.g. [3 5]).
	ReminderLevels []int
	// Compression reports the last history compression pass (zero if none).
	Compression CompressionReport
	// RootPending is the task root's undrained goroutine count at finish
	// (D38e: a completed task has drained to 0).
	RootPending int
	Err         *observe.Error
}

// Config carries the loop's policy knobs; zero values fall back to the scaled
// D15/D39/C22 defaults.
type Config struct {
	// Model and ContextWindow come from the resolved role endpoint. A zero
	// ContextWindow is read from provider Info().MaxContextWindow.
	Model         string
	ContextWindow int

	MaxRounds        int
	TokenBudget      int
	RepeatThresholds []int
	PerToolTimeout   time.Duration
	ToolConcurrency  int

	// ArtifactsDir is the memory store's artifacts directory (spill writes
	// tool-output-<id>.txt there; ticket 04's 500MB LRU job owns their life).
	ArtifactsDir string

	// PassThroughUnclassifiedRisk routes calls whose provider reports no risk
	// level as L0 pass-through. Ticket 21 replaces this with the C19 gate;
	// with the flag off an unclassified call fails closed instead.
	PassThroughUnclassifiedRisk bool

	Price    config.Price
	Currency string

	SteeringEnabled bool

	// Profiles supplies the (3) L1 section; the caller reads them from memory
	// (ticket 29 owns extraction, this ticket owns only the injection point).
	Profiles []string

	// SceneFunc renders the (4) section per turn (empty fields when unknown).
	SceneFunc func() Scene
}

// Options wire the loop to its seams.
type Options struct {
	// Provider is the C5 seam; a fallback chain wrapper (llm.ChainRunner) may
	// be passed here, so chain behaviour needs no loop change.
	Provider llm.LlmProvider
	// Tools is the C4 surface (echo provider in tests, host bridge at 20+).
	Tools ToolProvider
	// Sink receives user-visible events (nil = NopSink).
	Sink Sink
	// Journal is the 04-tables writer (nil = no persistence).
	Journal Journal
	// Summarizer backs history compression (nil = structural trim).
	Summarizer Summarizer
	// Control executes D11(1) verbs in the host (nil = the loop's own default,
	// which can cancel/stop its running task and reports the rest unhandled).
	Control ControlHandler
	// Registry spawns tool-execution goroutines (nil = observe.Default).
	Registry *observe.Registry
	// AdmitTask registers a task with the host's approval layer the moment the
	// loop knows its own task id, and returns the revocation the loop runs
	// when the task ends. D47 says only a task the TEXT loop registered may
	// reach a gate at all, and the task id is minted inside run(), so this is
	// the one place the registration can honestly happen (cmd/wisp passes
	// approval.Gate.AdmitTextTask here). Nil means the host has no approval
	// layer, which is also what keeps decideRisk refusing declared L1/L2
	// calls: an ungated host must never execute a write.
	AdmitTask func(taskID string) (revoke func())
	// Logger receives loop diagnostics (nil = slog.Default).
	Logger *slog.Logger
	Config
}

// Loop is a single-task agent loop holding one conversation history.
type Loop struct {
	opt      Options
	b        Budgets
	guard    GuardConfig
	asm      *Assembler
	sp       *Spiller
	comp     *Compressor
	reg      *observe.Registry
	provider llm.LlmProvider

	mu      sync.Mutex
	history []llm.Message
	steer   []llm.Message
	current *observe.Root
}

// New validates the wiring and freezes the scaled budgets for the endpoint's
// context window.
func New(opt Options) (*Loop, error) {
	if opt.Provider == nil {
		return nil, observe.New(observe.ClassConfig, "agent: no LlmProvider (C5 seam) wired")
	}
	if opt.Tools == nil {
		opt.Tools = NewEchoProvider()
	}
	if opt.Sink == nil {
		opt.Sink = NopSink{}
	}
	if opt.Registry == nil {
		opt.Registry = observe.Default
	}
	if opt.Logger == nil {
		opt.Logger = slog.Default()
	}
	info := opt.Provider.Info()
	window := opt.Config.ContextWindow
	if window <= 0 {
		window = info.MaxContextWindow
	}
	b := BudgetsFor(window)
	if opt.Config.MaxRounds <= 0 {
		opt.Config.MaxRounds = 50
	}
	if len(opt.Config.RepeatThresholds) == 0 {
		opt.Config.RepeatThresholds = []int{3, 5, 8}
	}
	if opt.Config.Currency == "" {
		opt.Config.Currency = "CNY"
	}
	if opt.Config.SceneFunc == nil {
		opt.Config.SceneFunc = func() Scene { return Scene{Now: observe.NowWallUTC()} }
	}
	model := opt.Config.Model
	if model == "" {
		model = info.Model
	}
	l := &Loop{
		opt:      opt,
		b:        b,
		reg:      opt.Registry,
		provider: opt.Provider,
		asm:      NewAssembler(model, b, info.Cache),
		sp:       NewSpiller(opt.Config.ArtifactsDir, b),
		comp:     NewCompressor(b, opt.Summarizer),
		guard: GuardConfig{
			RepeatThresholds: opt.Config.RepeatThresholds,
			MaxRounds:        opt.Config.MaxRounds,
			TokenBudget:      opt.Config.TokenBudget,
			PerToolTimeout:   opt.Config.PerToolTimeout,
			ToolConcurrency:  opt.Config.ToolConcurrency,
			Budgets:          b,
		},
	}
	return l, nil
}

// Budgets reports the scaled budget set the loop runs on.
func (l *Loop) Budgets() Budgets { return l.b }

// History returns a copy of the conversation so far.
func (l *Loop) History() []llm.Message {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]llm.Message, len(l.history))
	copy(out, l.history)
	return out
}

// Reset drops the conversation history (session end).
func (l *Loop) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.history = nil
	l.steer = nil
}

// Steer inserts a runtime instruction into the inner queue (Pi's two-layer
// steering, D21#4): it joins the conversation before the next model call; it
// is not a structural restart. It returns false when steering is disabled.
func (l *Loop) Steer(text string) bool {
	if !l.opt.Config.SteeringEnabled || strings.TrimSpace(text) == "" {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.steer = append(l.steer, userMessage(text))
	return true
}

// RunningTask is a task spawned through the registry under the D38 roster name
// agent-task-<id>.
type RunningTask struct {
	ID   string
	root *observe.Root
	h    *observe.Handle

	done   chan struct{}
	result Result
}

// Root exposes the task root (ctx + join counter).
func (t *RunningTask) Root() *observe.Root { return t.root }

// Cancel requests cancellation (root ctx, honoured loop -> stream -> tool).
func (t *RunningTask) Cancel() { t.root.Cancel() }

// Wait blocks until the task has fully finished and returns its result. It
// joins BOTH the result signal (t.done, closed inside the task goroutine) and
// the registry handle (t.h.done, which the registry closes only AFTER it has
// decremented the root's pending counter). Joining only t.done could let a
// caller read Pending()==1 immediately after Wait() returned - a task whose
// completion (D38e: the join counter drained) is therefore not yet observable.
func (t *RunningTask) Wait() Result {
	<-t.done
	if t.h != nil {
		<-t.h.Done()
	}
	return t.result
}

// Pending reports the root's undrained goroutines: 0 means the task's work is
// fully joined (D38e: completion means the WaitGroup drained).
func (t *RunningTask) Pending() int { return t.root.Pending() }

// RunAsync spawns the task on a roster-named goroutine and returns immediately.
func (l *Loop) RunAsync(ctx context.Context, input string) *RunningTask {
	id := newTaskID()
	t := &RunningTask{ID: id, root: observe.NewRootFrom(ctx, id), done: make(chan struct{})}
	t.h = l.reg.Spawn("agent-task-"+id, "agent", t.root, func(c context.Context) {
		defer close(t.done)
		t.result = l.run(c, id, input)
	})
	return t
}

// Run executes one task synchronously on the calling goroutine.
func (l *Loop) Run(ctx context.Context, input string) Result {
	return l.run(ctx, newTaskID(), input)
}

// ---------------------------------------------------------------------------
// core

func (l *Loop) run(ctx context.Context, taskID, input string) Result {
	// (1) D11(1): control words never reach the provider.
	if verb, ok := MatchControl(input); ok {
		return l.runControl(taskID, verb, input)
	}

	root := observe.NewRootFrom(ctx, taskID)
	defer root.Cancel()
	// Operate under the root's own cancelable context: this is what lets the
	// D11(1) control layer abort a running task by cancelling l.current (the
	// default handler reaches the loop only through this root, not the caller's
	// parent ctx). Cancellation set by the caller's parent propagates down here
	// as well, so both paths are observed by the loop's ctx.Err() checks.
	ctx = root.Ctx
	l.setCurrent(root)
	defer l.setCurrent(nil)

	j := newTaskJournal(l.opt.Journal, taskID, taskID) // C18: correlation == task id
	// D47: register the task with the host's approval layer before any tool
	// call can reach a gate. The revocation runs when the task ends, so an
	// admitted id cannot outlive the task that earned it.
	revokeAdmission := func() {}
	if l.opt.AdmitTask != nil {
		if r := l.opt.AdmitTask(taskID); r != nil {
			revokeAdmission = r
		}
	}
	defer revokeAdmission()
	guard := NewGuard(l.guard)
	cost := Cost{Currency: l.opt.Config.Currency}
	res := Result{TaskID: taskID, Status: StatusCompleted}

	l.append(userMessage(input))
	if err := j.startTask(ctx, input); err != nil {
		l.log().Warn("agent: task_log open failed", "err", err, "task", taskID)
	}

	for {
		if ctx.Err() != nil {
			return l.finish(ctx, root, j, guard, cost, &res, StatusCancelled, "任务已取消")
		}
		// Token brake first: an exhausted budget must not start another call.
		if guard.TokenBudget() > 0 && guard.TokensUsed() >= guard.TokenBudget() {
			return l.brakeStuck(ctx, root, j, guard, cost, &res, BrakeToken,
				fmt.Sprintf("本次任务已达 token 预算上限并停止：%s。可以说「继续」接着做，或把任务拆小。",
					guard.Consumption()))
		}
		if _, floorHit := guard.NextRound(); floorHit {
			return l.brakeStuck(ctx, root, j, guard, cost, &res, BrakeRounds,
				fmt.Sprintf("本次任务已达 %d 轮上限并停止：%s。", l.guard.MaxRounds, guard.Consumption()))
		}
		l.drainSteering()

		hist := l.History()
		if l.comp.Need(hist) {
			// DEFERRED(D28-1): the Warm-window hook owns this call; the
			// synchronous fallback here is what the S1 slice accepts.
			nh, rep, err := l.comp.Compress(ctx, hist)
			if err != nil {
				l.log().Warn("agent: history compression failed", "err", err)
			} else if rep.Ran {
				l.replaceHistory(nh)
				res.Compression = rep
			}
		}

		req, err := l.buildRequest(ctx)
		if err != nil {
			return l.failWith(ctx, root, j, guard, cost, &res,
				observe.Wrap(observe.ClassInternal, err, "agent: assemble context"))
		}

		turn, streamErr := l.stream(ctx, taskID, req)
		if turn.Err == nil && streamErr != nil {
			turn.Err = observe.Wrap(observe.ClassOfOrInternal(streamErr), streamErr, "agent: stream failed")
		}
		res.Rounds = guard.Rounds()

		if ctx.Err() != nil || turn.Stop == llm.StopCancelled {
			l.failOpenCalls(ctx, j, taskID, turn, OutcomeCancelled, string(observe.ClassCancelled), &res)
			if turn.Text != "" {
				l.append(assistantMessage(turn.Text))
			}
			return l.finish(ctx, root, j, guard, cost, &res, StatusCancelled, "任务已取消")
		}
		if turn.Err != nil {
			l.failOpenCalls(ctx, j, taskID, turn, OutcomeError, string(turn.Err.Class), &res)
			if turn.Text != "" {
				l.append(assistantMessage(turn.Text)) // partial kept, marked (SPEC-05 §3.4)
			}
			return l.failWith(ctx, root, j, guard, cost, &res, turn.Err)
		}
		if !turn.Usage.IsZero() {
			delta := cost.AddUsage(l.opt.Config.Price, turn.Usage)
			guard.AddCost(delta)
			if guard.AddUsage(turn.Usage) {
				l.publishUsage(&res, cost)
				return l.brakeStuck(ctx, root, j, guard, cost, &res, BrakeToken,
					fmt.Sprintf("本次任务已达 token 预算上限并停止：%s。", guard.Consumption()))
			}
			l.publishUsage(&res, cost)
		}

		// (D21) A length-truncated message never executes its tool calls.
		if turn.Stop == llm.StopMaxTokens {
			n := l.failToolCallsFromTruncatedMessage(ctx, j, taskID, turn, &res)
			msg := "响应在输出上限处被截断（未闭合的内容已标注）。"
			if n > 0 {
				msg = fmt.Sprintf(
					"响应在输出上限处被截断：本条消息的 %d 个工具调用已全部判失败，未执行任何截断参数。%s",
					n, guard.Consumption())
			}
			res.Text = turn.Text
			l.publish(Event{Kind: EvError, TaskID: taskID, Text: msg, Stop: turn.Stop})
			return l.finish(ctx, root, j, guard, cost, &res, StatusTruncated, msg)
		}

		if len(turn.ToolCalls) == 0 {
			// No tool call plus finished text = the loop's normal end.
			if turn.Text != "" {
				l.append(assistantMessage(turn.Text))
			}
			res.Text = turn.Text
			guard.ResetRepeats()
			return l.finish(ctx, root, j, guard, cost, &res, StatusCompleted, "")
		}

		// C22 gradient brake ladder on duplicate calls.
		rem, repeated := guard.ObserveTurn(turn.ToolCalls)
		if repeated {
			res.ReminderLevels = append(res.ReminderLevels, rem.Level)
			l.publish(Event{
				Kind: EvReminder, TaskID: taskID, Text: rem.Text,
				ToolName: rem.Tool, RepeatLevel: rem.Level,
			})
			if rem.Stuck {
				l.publish(Event{
					Kind: EvStuck, TaskID: taskID,
					Text: rem.Text, RepeatLevel: rem.Level,
				})
				res.Text = turn.Text
				return l.brakeStuck(ctx, root, j, guard, cost, &res, BrakeRepeat,
					fmt.Sprintf("%s 已连续调用 %d 次且参数相同，任务停在 Stuck 态。%s",
						rem.Tool, rem.Repeat, guard.Consumption()))
			}
		}

		// Record the assistant turn, then execute and feed results back.
		l.append(assistantTurnMessage(turn))
		l.executeCalls(ctx, root, j, taskID, guard, turn, &res)
		if rem.Text != "" {
			l.append(reminderMessage(rem.Text))
		}
		res.Text = turn.Text
		if ctx.Err() != nil {
			return l.finish(ctx, root, j, guard, cost, &res, StatusCancelled, "任务已取消")
		}
	}
}

// runControl handles a D11(1) utterance: zero provider round-trips, no
// task_log row (a control word is not a task).
func (l *Loop) runControl(taskID string, verb ControlVerb, utterance string) Result {
	l.publish(Event{Kind: EvControl, TaskID: taskID, Verb: verb, Text: utterance})
	h := l.opt.Control
	if h == nil {
		h = l.defaultControl
	}
	out := h(verb, utterance)
	res := Result{TaskID: taskID, Status: StatusControl}
	switch {
	case out.Text != "":
		res.Message = out.Text
	case !out.Handled:
		res.Message = "收到控制指令 " + string(verb) + "，但当前没有可执行对象。"
	}
	return res
}

// defaultControl covers what the loop itself can do: cancel/stop its running
// task. repeat/louder/confirm belong to the host (TTS, approval queue) and
// report unhandled unless a ControlHandler is wired.
func (l *Loop) defaultControl(verb ControlVerb, _ string) ControlOutcome {
	l.mu.Lock()
	root := l.current
	l.mu.Unlock()
	if root != nil {
		switch verb {
		case ControlStop, ControlCancel:
			root.Cancel()
			return ControlOutcome{Handled: true, Text: "已停止当前任务"}
		}
	}
	return ControlOutcome{Handled: false}
}

// stream runs one provider call and consumes the C6 events.
func (l *Loop) stream(ctx context.Context, taskID string, req *llm.Request) (llm.TurnResult, error) {
	col := llm.NewTurnCollector()
	err := l.provider.Stream(ctx, req, func(ev llm.StreamEvent) error {
		col.Observe(ev)
		switch ev.Type {
		case llm.EvTextDelta:
			l.publish(Event{Kind: EvTextDelta, TaskID: taskID, Text: ev.Text})
		case llm.EvReasoningDelta:
			l.publish(Event{Kind: EvReasoningDelta, TaskID: taskID, Text: ev.Text})
		}
		return nil
	})
	return col.Result(), err
}

// buildRequest assembles the D39 context for one turn.
func (l *Loop) buildRequest(ctx context.Context) (*llm.Request, error) {
	dir, err := l.opt.Tools.Tools(ctx)
	if err != nil {
		return nil, observe.Wrap(observe.ClassInternal, err, "agent: tool directory unavailable")
	}
	inj := Inject(latestUserText(l.History()), dir, l.b)
	in := PromptInput{
		Profiles:  l.opt.Config.Profiles,
		Scene:     l.opt.Config.SceneFunc(),
		Injection: inj,
	}
	return l.asm.Build(in, l.History()), nil
}

// toolPlan is one call's pre-flight decisions (row, gate outcome, skip).
type toolPlan struct {
	call     llm.ToolCall
	rowID    int64
	skip     bool
	rejected bool
	outcome  string
	class    string
	reason   string
	local    bool // answered by the loop (reserved list_tools path)
}

// executeCalls runs a turn's tool calls: risk policy first (pass-through until
// ticket 21), then D38d-capped concurrent execution with the C22 per-tool
// timeout, then D15(3) spill, then the results back into the history. Tool
// failures go to the model as class "tool" for self-correction (D37) instead
// of failing the task.
func (l *Loop) executeCalls(ctx context.Context, root *observe.Root, j *taskJournal,
	taskID string, guard *Guard, turn llm.TurnResult, res *Result,
) {
	dir, err := l.opt.Tools.Tools(ctx)
	if err != nil {
		l.log().Warn("agent: tool directory unavailable at dispatch", "err", err)
	}
	byName := make(map[string]ToolInfo, len(dir))
	for _, t := range dir {
		byName[t.Name] = t
	}
	listToolsDeclared := false
	if _, ok := byName[ToolListName]; ok {
		listToolsDeclared = ok
	}

	plans := make([]toolPlan, len(turn.ToolCalls))
	handles := []*observe.Handle{}
	sem := make(chan struct{}, guard.Concurrency())
	results := make([]ToolOutcome, len(turn.ToolCalls))
	execErr := make([]error, len(turn.ToolCalls))

	for i, c := range turn.ToolCalls {
		p := toolPlan{call: c}
		info := byName[c.Name]
		p.rowID = j.startCall(ctx, c.ID, c.Name, c.Args, riskColumn(info.RiskLevel))

		switch {
		case !c.Complete:
			// The seam never fabricates ToolCallEnd: an unclosed call is
			// failed, its partial arguments are never executed.
			p.skip, p.outcome, p.class = true, OutcomeTruncated, string(observe.ClassLoop)
			p.reason = "工具调用参数流未闭合，已判失败（未执行）"
		case c.Name == ToolListName && !listToolsDeclared:
			p.local = true // D15(2) fallback answered in-loop (L0, D34 table)
		default:
			ok, _, why := l.decideRisk(ctx, j, p.rowID, info)
			if !ok {
				p.skip, p.outcome, p.class, p.reason = true, OutcomeError,
					string(observe.ClassPermissionDenied), why
				p.reject()
			}
			// When ok came back for a declared L1/L2 call the loop booked no
			// decision: this row keeps decision='' (pending in the frozen
			// vocabulary) because the gate owns that column, and the host
			// bridge writes the authoritative row with the ASSESSED level.
		}
		plans[i] = p
	}

	for i, p := range plans {
		if p.skip {
			continue
		}
		p, i := p, i
		if p.local {
			results[i] = l.localListTools(ctx)
			continue
		}
		timeout := guard.PerToolTimeout()
		req := ToolRequest{
			TaskID: taskID, CorrelationID: taskID, CallID: p.call.ID,
			Name: p.call.Name, Args: p.call.Args, Timeout: timeout,
		}
		h := l.reg.Spawn("tool-exec-"+taskID, "agent", root, func(c context.Context) {
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i], execErr[i] = l.dispatch(c, j, p.rowID, req, timeout)
		})
		handles = append(handles, h)
	}
	for _, h := range handles {
		<-h.Done()
	}
	for _, h := range handles {
		if perr := h.Err(); perr != nil {
			// A recovered panic (D37b) is a tool failure, not a dead task.
			l.log().Warn("agent: tool goroutine panic recovered", "err", perr)
		}
	}

	// Terminal tool_call rows are forensics of how the task ended, so this
	// loop's writes run under a ctx detached from cancellation (see
	// terminalWriteCtx): a task cancelled mid-execution must still book its
	// in-flight calls instead of leaving them pending forever.
	wctx, cancelWrites := terminalWriteCtx(ctx)
	defer cancelWrites()
	for i, c := range turn.ToolCalls {
		p := plans[i]
		log := ToolResultLog{CallID: c.ID, Name: c.Name, Rejected: p.rejected}
		switch {
		case p.skip:
			log.Outcome, log.ErrorClass, log.Text = p.outcome, p.class, p.reason
		case execErr[i] != nil && ctx.Err() != nil:
			// The task was cancelled: an in-flight call is a cancelled call,
			// not an internal failure (the outcome vocabulary exists for it).
			log.Outcome = OutcomeCancelled
			log.ErrorClass = string(observe.ClassCancelled)
			log.Text = "任务已取消，调用被中止"
		case execErr[i] != nil:
			log.Outcome = OutcomeError
			log.ErrorClass = ErrorClassOfTurnError(execErr[i])
			log.Text = execErr[i].Error()
		default:
			out := results[i]
			log.Text = out.Text
			log.Truncated = out.Truncated
			if out.IsError {
				log.Outcome = OutcomeError
				log.ErrorClass = orClass(out.ErrorClass, observe.ClassTool)
			} else {
				log.Outcome = OutcomeSuccess
			}
		}

		// D15(3) spill of whatever goes back into the context.
		sp, err := l.sp.Prepare(c.ID, log.Text)
		if err != nil {
			l.log().Warn("agent: spill failed", "err", err, "call", c.ID)
		} else {
			log.Text = sp.Text
			log.Spilled = sp.Spilled
			log.Artifact = sp.Path
			log.Truncated = log.Truncated || sp.TruncatedRaw
		}

		j.finish(wctx, c.ID, p.rowID, log.Outcome, log.ErrorClass)
		res.ToolLog = append(res.ToolLog, log)
		res.ToolCalls++
		l.append(toolResultMessage(c.ID, log.Text, log.Outcome != OutcomeSuccess))
		l.publish(Event{
			Kind: EvToolEnd, TaskID: taskID, CallID: c.ID,
			ToolName: c.Name, Outcome: log.Outcome, Text: log.Text,
		})
	}
}

func (p *toolPlan) reject() { p.rejected = true }

// dispatch executes one call under the C22 timeout (cooperative abort).
func (l *Loop) dispatch(ctx context.Context, j *taskJournal, rowID int64,
	req ToolRequest, timeout time.Duration,
) (ToolOutcome, error) {
	if timeout <= 0 {
		return l.opt.Tools.Execute(ctx, req)
	}
	tctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	out, err := l.opt.Tools.Execute(tctx, req)
	if errors.Is(tctx.Err(), context.DeadlineExceeded) {
		// Structured TOOL_TIMEOUT (C22): the row says timeout, the model is
		// told the call was aborted so it can shorten its request. The tool's
		// own error is dropped ON PURPOSE - a contract-honest host reports a
		// timed-out call as an error (tools.go:82-85), and returning it here
		// made executeCalls prefer it, so the model saw a bare "context
		// deadline exceeded" and the row was booked error_class="internal"
		// (a class the model cannot self-correct, D37) instead of "tool".
		// The outcome below expresses the failure completely.
		l.log().Warn("agent: tool timeout", "tool", req.Name, "task", req.TaskID,
			"timeout_ms", timeout.Milliseconds())
		return ToolOutcome{
			Text: fmt.Sprintf("工具 %s 超时（%dms），已协作式中止",
				req.Name, timeout.Milliseconds()), IsError: true,
			ErrorClass: string(observe.ClassTool),
		}, nil
	}
	return out, err
}

// decideRisk resolves what the LOOP may decide on its own, and books the
// decision column. It is not the risk verdict: ToolInfo.RiskLevel carries the
// tool's DECLARED level, i.e. the R1 lower bound only, and C19 owns the
// conclusion. Two shapes follow from that:
//
//   - a declared L1/L2 call is passed through UNLESS the host registered its
//     tasks with an approval layer (Options.AdmitTask). Where a bridge and a
//     gate are wired, the routing below the loop does the asking, and the loop
//     books no decision, so it never writes "allow" over a call the user is
//     about to veto (ticket 12's assembly, ruling A13). Where none is wired the
//     old refusal stands: an unapproved write must not run just because a host
//     forgot its gate.
//   - an unclassified call still obeys PassThroughUnclassifiedRisk, which is a
//     host-level policy switch, not a verdict.
//
// The second result reports "booked nothing" (gate-owned routing) so the caller
// can drop its own pending row instead of leaving one open forever.

func (l *Loop) decideRisk(ctx context.Context, j *taskJournal, rowID int64, info ToolInfo) (ok, unbooked bool, why string) {
	risk := info.RiskLevel
	switch risk {
	case RiskL1, RiskL2:
		if l.opt.AdmitTask == nil {
			j.decide(ctx, rowID, DecisionReject)
			return false, false, "该操作属于 " + risk + " 级，审批通道尚未接入（ticket 21），已拒绝执行"
		}
		return true, true, ""
	default:
		if !l.opt.Config.PassThroughUnclassifiedRisk {
			j.decide(ctx, rowID, DecisionReject)
			return false, false, "风险未分级且直通开关关闭，已拒绝执行"
		}
		j.decide(ctx, rowID, DecisionAllow)
		return true, false, ""
	}
}

// localListTools answers list_tools from the loop when the provider does not
// expose it (D15(2) fallback: an in-loop call, never an extra round-trip).
func (l *Loop) localListTools(ctx context.Context) ToolOutcome {
	dir, err := l.opt.Tools.Tools(ctx)
	if err != nil {
		return ToolOutcome{
			Text: "工具目录不可用：" + err.Error(), IsError: true,
			ErrorClass: string(observe.ClassInternal),
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "共 %d 个工具：\n", len(dir))
	for _, t := range dir {
		kind := "第三方"
		if t.Resident {
			kind = "内置"
		}
		fmt.Fprintf(&b, "- %s [%s/%s] %s\n", t.Name, kind, riskLabel(t.RiskLevel), t.Description)
	}
	return ToolOutcome{Text: b.String()}
}

// failToolCallsFromTruncatedMessage is the D21 must-steal rule: when the
// provider stops a message with max_tokens, EVERY tool call collected in that
// message is failed and none of its arguments is executed. It is asserted at
// the golden level (testdata/golden/max-tokens-toolcall.sse). The return value
// is how many calls were failed.
func (l *Loop) failToolCallsFromTruncatedMessage(ctx context.Context, j *taskJournal,
	taskID string, turn llm.TurnResult, res *Result,
) int {
	if len(turn.ToolCalls) == 0 {
		return 0
	}
	// The assistant turn is recorded first so the failed results attach to the
	// call ids they belong to (tool_call ids must stay in the chain).
	l.append(assistantTurnMessage(turn))
	for _, c := range turn.ToolCalls {
		row := j.startCall(ctx, c.ID, c.Name, c.Args, memoryRiskL0)
		j.decide(ctx, row, DecisionReject)
		j.finish(ctx, c.ID, row, OutcomeTruncated, string(observe.ClassLoop))
		text := "该消息因达到输出上限被截断，其工具调用已全部判失败（截断参数未执行）"
		l.append(toolResultMessage(c.ID, text, true))
		res.ToolLog = append(res.ToolLog, ToolResultLog{
			CallID: c.ID, Name: c.Name, Outcome: OutcomeTruncated,
			ErrorClass: string(observe.ClassLoop), Text: text,
		})
		res.ToolCalls++
		l.publish(Event{
			Kind: EvToolEnd, TaskID: taskID, CallID: c.ID, ToolName: c.Name,
			Outcome: OutcomeTruncated, Text: text,
		})
	}
	return len(turn.ToolCalls)
}

// failOpenCalls fails every call the seam left open after a failure or a
// cancellation (SPEC-05 §3.4: partial content is kept and marked; open tool
// calls are all judged failed, never executed with partial arguments).
func (l *Loop) failOpenCalls(ctx context.Context, j *taskJournal, taskID string,
	turn llm.TurnResult, outcome, class string, res *Result,
) {
	// These rows describe how the task ENDED (cancellation included), so they
	// are written under a ctx detached from that cancellation.
	wctx, cancelWrites := terminalWriteCtx(ctx)
	defer cancelWrites()
	for _, c := range turn.ToolCalls {
		if c.Complete {
			continue
		}
		row := j.startCall(wctx, c.ID, c.Name, c.Args, memoryRiskL0)
		// The gate column is part of the row contract and this path made the
		// same judgement the max_tokens path does (never execute an unclosed
		// call), so it must book the decision too - otherwise one table ends up
		// with two write disciplines and the row reads as "pending, forever".
		j.decide(wctx, row, DecisionReject)
		j.finish(wctx, c.ID, row, outcome, class)
		text := "调用未闭合（响应中断或已取消），未执行"
		res.ToolLog = append(res.ToolLog, ToolResultLog{
			CallID: c.ID, Name: c.Name,
			Outcome: outcome, ErrorClass: class, Text: text,
		})
		res.ToolCalls++
		l.publish(Event{
			Kind: EvToolEnd, TaskID: taskID, CallID: c.ID, ToolName: c.Name,
			Outcome: outcome, Text: text,
		})
	}
}

// brakeStuck ends the task in Stuck with the visible message C22 demands.
func (l *Loop) brakeStuck(ctx context.Context, root *observe.Root, j *taskJournal,
	guard *Guard, cost Cost, res *Result, brake Brake, msg string,
) Result {
	res.Brake = brake
	l.publish(Event{
		Kind: EvStuck, TaskID: res.TaskID, Text: msg,
		TokensIn: guard.tokensIn, TokensOut: guard.tokensOut,
	})
	return l.finish(ctx, root, j, guard, cost, res, StatusStuck, msg)
}

func (l *Loop) failWith(ctx context.Context, root *observe.Root, j *taskJournal,
	guard *Guard, cost Cost, res *Result, err *observe.Error,
) Result {
	res.Err = err
	res.Message = err.Error()
	l.publish(Event{Kind: EvError, TaskID: res.TaskID, Text: err.Error(), Err: err})
	return l.finish(ctx, root, j, guard, cost, res, StatusFailed, res.Message)
}

// finish writes the task_log row, drains the task root (D38e: completion means
// every spawned goroutine joined) and returns the result.
func (l *Loop) finish(ctx context.Context, root *observe.Root, j *taskJournal,
	guard *Guard, cost Cost, res *Result, status Status, msg string,
) Result {
	res.Status = status
	if msg != "" && res.Message == "" {
		res.Message = msg
	}
	res.Rounds = guard.Rounds()
	res.Usage = cost.Usage
	res.CostMicros = cost.Micros
	res.Currency = cost.Currency
	if root != nil {
		res.RootPending = root.Wait(3 * time.Second)
	}
	errClass := ""
	if res.Err != nil {
		errClass = string(res.Err.Class)
	}
	summary := res.Text
	if summary == "" {
		summary = msg
	}
	// The terminal row is the audit record of HOW the task ended, so it must
	// survive the very cancellation that ended it: detaching from ctx keeps
	// memory's writer from abandoning the write (an already-dead ctx makes
	// BeginTx fail outright). The 2s bound is a context deadline, not a
	// wall-clock difference (D42#9).
	wctx, cancelWrite := terminalWriteCtx(ctx)
	defer cancelWrite()
	if err := j.finishTask(wctx, taskLogState(status), firstRunes(summary, 400),
		int64(cost.Usage.InputTokens), int64(cost.Usage.OutputTokens),
		cost.Micros, errClass); err != nil {
		l.log().Warn("agent: task_log finish failed", "err", err, "task", res.TaskID)
	}
	l.publish(Event{
		Kind: EvDone, TaskID: res.TaskID, Text: summary,
		TokensIn: cost.Usage.InputTokens, TokensOut: cost.Usage.OutputTokens,
	})
	return *res
}

// terminalWriteTimeout bounds a detached forensics write. It is a context
// deadline (D42#9 forbids wall-clock-difference timeouts).
const terminalWriteTimeout = 2 * time.Second

// terminalWriteCtx detaches a forensics write from the task's own cancellation
// and bounds it. Cancellation is one of the four terminal states the contract
// names, and it requires task_log/tool_call rows with error_class/decision
// just like the clean paths do - so a cancelled task must still be able to
// persist its ending. Without this the write is passed to memory's writer with
// an already-cancelled ctx, BeginTx fails, and the row is stranded at
// state="running"/ended_at=NULL.
func terminalWriteCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), terminalWriteTimeout)
}

func (l *Loop) publishUsage(res *Result, cost Cost) {
	res.Usage = cost.Usage
	res.CostMicros = cost.Micros
	res.Currency = cost.Currency
}

// taskLogState maps a loop status onto the task_log.state vocabulary (a
// free-text column; 'running' and 'interrupted' are memory's own values).
func taskLogState(s Status) string {
	switch s {
	case StatusCompleted:
		return "done"
	case StatusStuck:
		return "stuck"
	case StatusCancelled:
		return "cancelled"
	case StatusTruncated:
		return "interrupted"
	case StatusFailed:
		return "error"
	default:
		return "done"
	}
}

// ---------------------------------------------------------------------------
// helpers

func (l *Loop) setCurrent(r *observe.Root) {
	l.mu.Lock()
	l.current = r
	l.mu.Unlock()
}

func (l *Loop) append(m llm.Message) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.history = append(l.history, m)
}

func (l *Loop) replaceHistory(hist []llm.Message) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.history = hist
}

func (l *Loop) drainSteering() {
	l.mu.Lock()
	st := l.steer
	l.steer = nil
	l.mu.Unlock()
	for _, m := range st {
		l.append(m)
	}
}

func (l *Loop) publish(e Event) {
	if l.opt.Sink != nil {
		l.opt.Sink.Publish(e)
	}
}

func (l *Loop) log() *slog.Logger {
	if l.opt.Logger != nil {
		return l.opt.Logger
	}
	return slog.Default()
}

func userMessage(text string) llm.Message {
	return llm.Message{Role: llm.RoleUser, Content: []llm.Content{llm.TextPart{Text: text}}}
}

func assistantMessage(text string) llm.Message {
	return llm.Message{Role: llm.RoleAssistant, Content: []llm.Content{llm.TextPart{Text: text}}}
}

func assistantTurnMessage(turn llm.TurnResult) llm.Message {
	var parts []llm.Content
	if turn.Text != "" {
		parts = append(parts, llm.TextPart{Text: turn.Text})
	}
	for _, c := range turn.ToolCalls {
		args := c.Args
		if len(args) == 0 {
			args = json.RawMessage(`{}`)
		}
		parts = append(parts, llm.ToolUsePart{ID: c.ID, Name: c.Name, Input: args})
	}
	if len(parts) == 0 {
		parts = []llm.Content{llm.TextPart{Text: ""}}
	}
	return llm.Message{Role: llm.RoleAssistant, Content: parts}
}

// toolResultMessage builds the one-ToolResultPart message the C5 seam requires
// per call id (openai-chat: a tool message carries exactly one result).
func toolResultMessage(id, text string, isErr bool) llm.Message {
	return llm.Message{Role: llm.RoleTool, Content: []llm.Content{llm.ToolResultPart{
		ID: id, Content: []llm.Content{llm.TextPart{Text: text}}, IsError: isErr,
	}}}
}

func reminderMessage(text string) llm.Message {
	return userMessage("[系统提醒] " + text)
}

// latestUserText is the retrieval key for the (2) BM25 tool index: the newest
// real user utterance in the history.
func latestUserText(hist []llm.Message) string {
	for i := len(hist) - 1; i >= 0; i-- {
		m := hist[i]
		if m.Role != llm.RoleUser {
			continue
		}
		for _, p := range m.Content {
			if t, ok := p.(llm.TextPart); ok && strings.TrimSpace(t.Text) != "" {
				return t.Text
			}
		}
	}
	return ""
}

func riskLabel(r string) string {
	if r == RiskUnclassified {
		return "未分级"
	}
	return r
}

// memoryRiskL0 is the NOT NULL placeholder the pre-gate phase must write;
// ticket 21 replaces it with the assessed level.
const memoryRiskL0 = "L0"

func orClass(s string, def observe.ErrorClass) string {
	if s == "" {
		return string(def)
	}
	if err := observe.ValidateErrorClass(s); err != nil {
		return string(def)
	}
	return s
}

// newTaskID returns a UUIDv4-shaped task id (the task_log.id column is
// documented as a UUID), built from crypto/rand so no new dependency appears.
func newTaskID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// Entropy source unavailable: fall back to a clock-derived unique id.
		n := time.Now().UnixNano()
		for i := range b {
			b[i] = byte(n >> (8 * uint(i%8)))
		}
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b[:])
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
}
