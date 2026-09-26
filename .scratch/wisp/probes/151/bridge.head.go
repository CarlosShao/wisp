package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
)

// MaxToolConcurrency is the D38d ceiling: at most four tool executions run at
// a time through one bridge. It is enforced HERE, in the choke point, not only
// in the agent loop - a host-internal caller (the panel, a future spill
// re-read, plugin dispatch) that reaches the bridge directly must hit the same
// ceiling, or the ceiling is a comment rather than a property.
const MaxToolConcurrency = 4

// DefaultToolTimeout is the C22 budget when neither the tool declaration nor
// the composition root names one.
const DefaultToolTimeout = 30 * time.Second

// Options configures a Bridge. Zero values are usable and pick the fail-closed
// default in every case; nothing here can switch enforcement off.
type Options struct {
	// Registry is the C4 directory. nil starts an empty one.
	Registry *Registry
	// Paths is the C26 canonicalizer plus the [fs] allowed_dirs allowlist.
	// nil means no path can be canonicalized, which R2 judges fail-closed as
	// out-of-scope L2 - it is not "no path check".
	Paths *PathCanonicalizer
	// Assessor overrides the composed C19 assessor (tests). When nil the
	// bridge builds one from Paths + the SPEC-06 §4.1 classifier +, per task,
	// the C25 detector bound to that task's scope.
	Assessor *risk.RiskAssessor
	// Provenance is the C25 engine: it feeds R4 and consumes the read results
	// that must be tainted. nil leaves R4 dormant.
	Provenance *risk.Provenance
	// Journal books the tool_call rows. nil disables booking (pure unit tests
	// of the decision logic).
	Journal agent.Journal
	// Gate receives the L1/L2 routing. nil is NoGate (fail closed).
	Gate Gate
	// DeclaredCaps is the HOST's view of what a tool declared - the seam
	// ticket 49 (session grants) and ticket 50 (manifest revocation) fill: a
	// plugin's capability face can shrink after registration, and the bridge
	// re-reads it per call so a stale registration can never keep executing.
	// nil falls back to each entry's Decl.Capabilities.
	DeclaredCaps func(toolName string) []Capability
	// Authorized is the host-level capability grant (config-fed, e.g. a
	// deployment built without the screen capability). Empty means "no filter
	// beyond what a tool declared".
	Authorized []Capability
	// MaxConcurrency clamps the in-bridge ceiling. Values above
	// MaxToolConcurrency are ignored: the D38d ceiling is not tunable upward.
	MaxConcurrency int
	// DefaultTimeout is the C22 budget for a tool without Decl.Timeout.
	DefaultTimeout time.Duration
	// OnUpdate receives a tool's streaming progress deltas (C1's onUpdate).
	// The bridge forwards them; what a host does with them (ball strip,
	// panel) is ticket 30/36's call.
	OnUpdate func(callID, toolName, delta string)
	// Logf is the security/audit sink. rules_hit and the reason are recorded
	// here because the frozen SPEC-02 tool_call DDL has no rules_hit column
	// (see book for the gap).
	Logf func(format string, args ...any)
	// OnDecision hands the structured verdict to the host: the confirmation
	// card renders RulesHit/Reason VERBATIM, so they must leave the bridge as
	// data, not as a string a UI would have to re-parse.
	OnDecision func(Decision)
	// Cancel is the D31 seam (see cancel.go): the veto poll a running tool
	// gets through the context, plus the applied-steps report the bridge adds
	// to that tool's answer. nil means no approval layer is wired, so nothing
	// can be vetoed mid-flight and no report can be built - which is the
	// honest state of a bridge whose gate is NoGate.
	Cancel CancelBus

	// Modes is the READ side of the user-facing permission mode (ticket 90,
	// ruling R20): nil means the strictest档 (ask_every_step) for every call,
	// which is the fail-closed direction - an unwired composition can never run
	// relaxed by accident. The bridge only ever reads it; switching lives in
	// internal/perm, where a switch costs a confirmation and an audit line.
	Modes ModeSource
	// Confirmations returns the B-tier single-file overrides on record (the
	// set risk.Gate's bOverrides argument was designed to receive). nil or
	// empty means "no file has ever been confirmed", so every B-tier path stays
	// "ask first". Populating it is ticket 21's confirmation flow; nothing in
	// this package can add to it.
	Confirmations func() map[string]bool
}

// Bridge is the host bridge: the single choke point every capability flows
// through (SPEC-01 §3 / SPEC-07 §2). It implements agent.ToolProvider, the
// consumer-side seam the agent loop already calls (loop.go:563 and
// loop.go:696-700), so there is no second dispatch path to forget to secure.
//
// Every Execute runs, in this order:
//
//	C3 capability check      undeclared -> hard reject, tool never runs
//	C19 risk.Assess          over C26-canonical paths + C25 taint (R1..R9)
//	routing                  L0 pass / L1 PendingWindow / L2 PendingApproval /
//	                         Deny refuse without asking
//	D38d ceiling + C22 deadline, then Tool.Execute
//	C25 provenance Mark      on every sensitive-source result
//	tool_call row            assessed risk_level + decision + outcome +
//	                         correlation_id
//
// Safe for concurrent use by multiple tasks.
type Bridge struct {
	reg      *Registry
	paths    *PathCanonicalizer
	assess   *risk.RiskAssessor
	injected bool // assess came from Options.Assessor: do not rebuild per call
	prov     *risk.Provenance
	j        agent.Journal
	gate     Gate
	declared func(string) []Capability
	authz    capSet
	sem      chan struct{}
	defTm    time.Duration
	onUpdate func(callID, tool, delta string)
	logf     func(string, ...any)
	onDec    func(Decision)
	cancel   CancelBus
	// modes and confirmations are ticket 90's two read-only injections: the
	// permission mode source, and the B-tier override set risk.Gate reads. Both
	// are read per call, and neither has a setter on the bridge.
	modes         ModeSource
	confirmations func() map[string]bool

	// classifier is shared: it is stateless over the frozen blacklist.
	classifier risk.SensitiveClassifier

	mu     sync.Mutex
	seqs   map[string]int64 // task id -> tool_call.seq
	scopes map[string]bool  // task id -> C25 scope opened
}

// New composes a Bridge. It never returns a partially enforcing bridge: a nil
// Registry, a nil Gate and an out-of-range ceiling are filled/clamped here.
func New(o Options) *Bridge {
	if o.Registry == nil {
		o.Registry = NewRegistry()
	}
	ceiling := o.MaxConcurrency
	if ceiling <= 0 || ceiling > MaxToolConcurrency {
		ceiling = MaxToolConcurrency
	}
	defTm := o.DefaultTimeout
	if defTm <= 0 {
		defTm = DefaultToolTimeout
	}
	gate := o.Gate
	if gate == nil {
		gate = NoGate{}
	}
	b := &Bridge{
		reg:      o.Registry,
		paths:    o.Paths,
		prov:     o.Provenance,
		j:        o.Journal,
		gate:     gate,
		declared: o.DeclaredCaps,
		sem:      make(chan struct{}, ceiling),
		defTm:    defTm,
		onUpdate: o.OnUpdate,
		logf:     o.Logf,
		onDec:    o.OnDecision,
		cancel:   o.Cancel,
		modes:    o.Modes,
		confirmations: func() map[string]bool {
			if o.Confirmations == nil {
				return nil // no file was ever confirmed: B tier stays "ask first"
			}
			return o.Confirmations()
		},
		classifier: NewSensitiveClassifier(),
		seqs:       map[string]int64{},
		scopes:     map[string]bool{},
	}
	if len(o.Authorized) > 0 {
		b.authz = newCapSet(o.Authorized...)
	}
	if o.Assessor != nil {
		b.assess, b.injected = o.Assessor, true
	} else {
		b.assess = risk.NewRiskAssessor().
			WithCanonicalizer(o.Paths).
			WithSensitiveClassifier(b.classifier)
	}
	return b
}

// Registry exposes the directory so a composition root can register tools.
func (b *Bridge) Registry() *Registry { return b.reg }

// Paths exposes the C26 seam the composition root shares with the tools.
func (b *Bridge) Paths() *PathCanonicalizer { return b.paths }

// Tools implements agent.ToolProvider: the whole directory every turn
// (D15(2): the loop selects locally, no extra round-trip).
//
// RiskLevel carries the DECLARED level, which is only the R1 lower bound; the
// verdict is computed per call in Execute. A host that gates on this field
// before calling Execute is gating on a lower bound - which is exactly the
// hazard in agent.Loop.decideRisk (internal/agent/loop.go:719-738), and why
// ticket 21 replaces that function with this bridge's verdict.
func (b *Bridge) Tools(_ context.Context) ([]agent.ToolInfo, error) {
	entries := b.reg.List()
	out := make([]agent.ToolInfo, 0, len(entries))
	for _, e := range entries {
		out = append(out, agent.ToolInfo{
			Name:        e.Tool.Name(),
			Description: e.Tool.Description(),
			Parameters:  e.Tool.Parameters(),
			Resident:    e.Decl.Resident,
			RiskLevel:   levelString(e.Decl.Declared),
		})
	}
	return out, nil
}

// Execute implements agent.ToolProvider - the choke point.
//
// The returned error is reserved for host-internal faults (D37 class
// internal/provider), because the loop treats non-nil as "the provider broke"
// (loop.go:654-657) and books a class the model cannot self-correct against.
//
// A REJECT IS NOT A FAULT. SPEC-07 §2 states the C3 rule as "未声明即拒绝调用
// （不是报错，是拒绝）": an undeclared capability is refused, not reported. So
// every rejection - capability, Deny, gate veto, timeout - comes back as an
// IsError outcome carrying the right D37 class with err == nil: the task
// continues, the model learns the call was refused, and the row records the
// refusal. Returning a Go error instead would surface a policy decision as
// "Wisp crashed", which is the failure mode the constraint exists to prevent.
func (b *Bridge) Execute(ctx context.Context, req agent.ToolRequest) (agent.ToolOutcome, error) {
	if req.CorrelationID == "" {
		req.CorrelationID = req.TaskID
	}
	entry, ok := b.reg.Lookup(req.Name)
	if !ok {
		// Unknown tool: class tool (self-correctable, D37), never internal.
		return b.reject(ctx, req, Decision{Tool: req.Name},
			"未知工具 "+req.Name+"，可用工具见 list_tools",
			string(observe.ClassTool), OutcomeKindNotFound)
	}

	// --- C3: capability enforcement, ahead of any risk work or execution ----
	dec, why := b.checkCaps(req, entry)
	if why != "" {
		return b.reject(ctx, req, dec, why, string(observe.ClassPermissionDenied), OutcomeKindCapability)
	}

	// --- parameter decode ---------------------------------------------------
	params, err := decodeArgs(req.Args)
	if err != nil {
		return b.reject(ctx, req, dec, "参数不是合法的 JSON 对象："+err.Error(),
			string(observe.ClassTool), OutcomeKindBadArgs)
	}
	dec.Params = params

	// --- C19: the one risk verdict (R1..R9 over C26 paths + C25 taint) ------
	rawPaths := pathArgs(params, entry.Decl.PathParams)
	dec.Paths = b.displayPaths(rawPaths)
	verdict := b.assessorFor(req.TaskID).Assess(req.Name, params,
		b.factsFor(ctx, entry, params, rawPaths))
	dec.Level = verdict.Level
	dec.RulesHit = verdict.RulesHit
	dec.Reason = verdict.Reason
	dec.SessionOverrideBlocked = verdict.SessionOverrideBlocked

	// --- ticket 90: the permission mode screens that verdict -----------------
	// Read once per call (AC#1: an explicit value, never a package global), then
	// handed to route as a PARAMETER. dec.Level keeps the assessed level, because
	// that is what the tool_call row's risk_level column means; the mode's effect
	// travels as Mode/ModeSilenced/ModeKept so "why was nobody asked" is readable
	// from the record.
	mode := b.permissionMode()
	sil := mode.Screen(verdict)
	dec.Mode = mode
	dec.ModeSilenced = sil.Silenced
	dec.ModeKept = sil.Kept
	if len(rawPaths) > 0 && verdict.Level >= risk.L1 {
		dec.Blacklist = b.readBlacklist(rawPaths)
	}
	switch {
	case sil.Silenced:
		b.log("tools: MODE-SILENCE mode=%s assessed=%s effective=%s tool=%s rules=%v "+
			"(a silenced question, not an allow: no user click happened)",
			mode, verdict.Level, sil.Level, req.Name, verdict.RulesHit)
	case sil.Kept != "":
		b.log("tools: MODE-REDLINE mode=%s still asks tool=%s rules=%v because=%s",
			mode, req.Name, verdict.RulesHit, sil.Kept)
	}
	b.emit(dec)

	// --- routing: L0 pass, L1 window, L2 approval, Deny never reaches a gate -
	ok2, why := b.route(ctx, &dec, sil)
	if !ok2 {
		class := string(observe.ClassUserRejected)
		if dec.Level == risk.Deny {
			class = string(observe.ClassPermissionDenied)
		}
		return b.reject(ctx, req, dec, why, class, dispositionOf(dec))
	}

	// --- execution under the D38d ceiling and the C22 deadline --------------
	return b.run(ctx, req, entry, dec)
}

// ---------------------------------------------------------------------------
// C3
// ---------------------------------------------------------------------------

// checkCaps applies the C3 rule: every capability a tool NEEDS must be in the
// set the host believes it DECLARED, and that set must sit inside the
// host-level authorization. A miss is a rejection verdict, never an error.
func (b *Bridge) checkCaps(req agent.ToolRequest, entry Entry) (Decision, string) {
	dec := Decision{
		Tool: req.Name, Provider: entry.Decl.Provider, Args: req.Args,
		Capabilities: dedupe(entry.Decl.Needs), CorrelationID: req.CorrelationID,
		TaskID: req.TaskID, CallID: req.CallID,
		Level:   entry.Decl.Declared,
		Timeout: b.timeoutFor(entry),
	}
	need := entry.Decl.Needs
	declared := b.capsOf(req.Name, entry)
	if miss := declared.missing(need); len(miss) > 0 {
		return dec, fmt.Sprintf("工具 %s 需要能力 %s，但该能力从未声明（C3：未声明即拒绝调用），已拒绝执行",
			req.Name, joinCaps(miss))
	}
	if b.authz != nil {
		if miss := b.authz.missing(need); len(miss) > 0 {
			return dec, fmt.Sprintf("工具 %s 需要的能力 %s 未获本机授权（C3），已拒绝执行",
				req.Name, joinCaps(miss))
		}
	}
	return dec, ""
}

// capsOf resolves the host's view of a tool's declared capability face.
func (b *Bridge) capsOf(name string, entry Entry) capSet {
	if b.declared != nil {
		return newCapSet(b.declared(name)...)
	}
	return newCapSet(entry.Decl.Capabilities...)
}

// ---------------------------------------------------------------------------
// routing
// ---------------------------------------------------------------------------

// route maps the EFFECTIVE level (the assessed verdict as screened by the
// permission mode, ticket 90) onto the gate branch that owns it (SPEC-06 §2)
// and books the decision column. It returns whether the call may proceed plus
// the user-visible reason when it may not.
//
// sil is a parameter, not a lookup: this is the one place a level becomes "ask
// the user" or "do not ask", so the mode that decided it has to be visible in
// the signature (AC#1). An implementation that reached for a global here would
// put the permission switch inside the enforcement layer.
func (b *Bridge) route(ctx context.Context, dec *Decision, sil risk.Silenced) (bool, string) {
	switch sil.Level {
	case risk.L0:
		dec.DecisionColumn = agent.DecisionAllow
		return true, ""

	case risk.Deny:
		dec.DecisionColumn = agent.DecisionReject
		return false, "该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：" + dec.Reason

	case risk.L1:
		a, why := b.gate.PendingWindow(ctx, *dec)
		switch a {
		case AnswerAllow, AnswerTimeout:
			// The L1 window running out unopposed MEANS EXECUTE (SPEC-06 §2).
			dec.DecisionColumn = agent.DecisionAllow
			return true, ""
		case AnswerVeto:
			dec.DecisionColumn = agent.DecisionReject
			return false, orDefault(why, "用户在 L1 确认窗口中否决了本次操作")
		default:
			dec.DecisionColumn = agent.DecisionReject
			return false, orDefault(why, "L1 确认窗口未放行，已拒绝执行")
		}

	case risk.L2:
		a, why := b.gate.PendingApproval(ctx, *dec)
		switch a {
		case AnswerAllow:
			dec.DecisionColumn = agent.DecisionAllow
			return true, ""
		case AnswerTimeout:
			// The queue never waits forever: an unanswered L2 auto-REJECTS
			// (ticket 21), deliberately the opposite polarity from L1.
			dec.DecisionColumn = agent.DecisionTimeout
			return false, orDefault(why, "审批超时未确认，已自动拒绝")
		default:
			dec.DecisionColumn = agent.DecisionReject
			return false, orDefault(why, "L2 审批未通过，已拒绝执行")
		}

	default:
		// An impossible level means a broken judge: fail closed (R9's posture).
		dec.DecisionColumn = agent.DecisionReject
		return false, "风险判定返回了不可能的等级，已 fail-closed 拒绝"
	}
}

// dispositionOf names a refused call's disposition for the audit line.
func dispositionOf(dec Decision) OutcomeKind {
	if dec.Level == risk.Deny {
		return OutcomeKindDenied
	}
	if dec.DecisionColumn == agent.DecisionTimeout {
		return OutcomeKindApprovalTimeout
	}
	return OutcomeKindRefused
}

// run executes the tool under the ceiling and the per-tool deadline.
func (b *Bridge) run(ctx context.Context, req agent.ToolRequest, entry Entry,
	dec Decision,
) (agent.ToolOutcome, error) {
	timeout := b.timeoutFor(entry)
	dec.Timeout = timeout

	// D38d ceiling. A cancelled task must not queue behind a full semaphore.
	select {
	case b.sem <- struct{}{}:
		defer func() { <-b.sem }()
	case <-ctx.Done():
		return b.close(ctx, req, dec, agent.ToolOutcome{
			Text:       "任务已取消，调用未执行",
			IsError:    true,
			ErrorClass: string(observe.ClassCancelled),
		}, OutcomeKindCancelled)
	}

	ectx, cancel := b.execContext(ctx, timeout)
	defer cancel()
	// D31: the running tool's only window onto the approval layer is its own
	// context. The handle carries the correlation id the veto is keyed on, and
	// Complete releases the gate's post-handoff record when the call is over.
	ectx = withCancel(ectx, cancelHandle{corr: orDefault(req.CorrelationID, req.TaskID), bus: b.cancel})
	if b.cancel != nil {
		defer b.cancel.Complete(orDefault(req.CorrelationID, req.TaskID))
	}

	onUpdate := func(delta string) {
		if b.onUpdate != nil && delta != "" {
			b.onUpdate(req.CallID, req.Name, delta)
		}
	}

	res, err := entry.Tool.Execute(ectx, req.Args, onUpdate)
	// The deadline is checked BEFORE the fault, because a well-behaved tool
	// aborts on ctx and returns ctx.Err(): that IS the timeout, and reporting
	// it as an internal fault would both lose the C22 structured wording and
	// book a class the model cannot self-correct against.
	if errors.Is(ectx.Err(), context.DeadlineExceeded) && ctx.Err() == nil {
		b.log("tools: %s exceeded its %dms budget", req.Name, timeout.Milliseconds())
		return b.close(ctx, req, dec, agent.ToolOutcome{
			Text: fmt.Sprintf("工具 %s 超时（%dms），已协作式中止",
				req.Name, timeout.Milliseconds()),
			IsError:    true,
			ErrorClass: string(observe.ClassTool),
		}, OutcomeKindTimeout)
	}
	if err != nil {
		// A fault INSIDE a tool. The bridge itself did not break, so this is
		// still a tool-class outcome for the model to reason about (D37) and
		// not a provider fault.
		b.log("tools: %s execute fault: %v", req.Name, err)
		return b.close(ctx, req, dec, agent.ToolOutcome{
			Text:       "工具内部故障：" + err.Error(),
			IsError:    true,
			ErrorClass: string(observe.ClassTool),
		}, OutcomeKindError)
	}
	if errors.Is(ectx.Err(), context.DeadlineExceeded) && ctx.Err() == nil {
		b.log("tools: %s exceeded its %dms budget", req.Name, timeout.Milliseconds())
		return b.close(ctx, req, dec, agent.ToolOutcome{
			Text: fmt.Sprintf("工具 %s 超时（%dms），已协作式中止",
				req.Name, timeout.Milliseconds()),
			IsError:    true,
			ErrorClass: string(observe.ClassTool),
		}, OutcomeKindTimeout)
	}

	out := agent.ToolOutcome{
		Text:         res.Text,
		IsError:      res.IsError,
		RiskLevel:    dec.LevelString(),
		Truncated:    res.Truncated,
		AppliedSteps: append([]string(nil), res.AppliedSteps...),
	}
	kind := OutcomeKindSuccess
	if res.IsError {
		out.ErrorClass = string(observe.ClassTool)
		kind = OutcomeKindError
	}

	// D31: a call that stopped because the user vetoed it is a cancellation,
	// not a tool fault, and it must carry the applied-steps report. The report
	// type belongs to the approval layer (this package deliberately keeps no
	// second one); the bridge only appends its user-visible text so the model
	// and the transcript both see what actually landed.
	if txt := b.cancelText(dec, res); txt != "" {
		out.Text = orDefault(out.Text, "调用已中止") + "\n\n" + txt
		if kind == OutcomeKindError && b.cancel.Vetoed(orDefault(dec.CorrelationID, dec.TaskID)) {
			kind = OutcomeKindCancelled
			out.ErrorClass = string(observe.ClassCancelled)
		}
	}

	// C25: a sensitive-source result is tainted from the moment it can reach
	// the context. Marking runs on the CONTENT and only for successful
	// results (an error text carries no user data).
	if !res.IsError {
		b.mark(dec, res)
	}
	return b.close(ctx, req, dec, out, kind)
}

// execContext puts the C22 per-tool deadline on the caller's ctx. The budget
// is a context deadline, never a wall-clock subtraction (D42#9, d22scan ban 4).
func (b *Bridge) execContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}

// timeoutFor resolves one call's C22 budget.
func (b *Bridge) timeoutFor(entry Entry) time.Duration {
	if entry.Decl.Timeout > 0 {
		return entry.Decl.Timeout
	}
	return b.defTm
}

// mark records C25 provenance for one result.
func (b *Bridge) mark(dec Decision, res Result) {
	if b.prov == nil || dec.TaskID == "" || !risk.IsSensitiveSource(dec.Tool) {
		return
	}
	origin := res.Origin
	if origin == "" && len(dec.Paths) > 0 {
		origin = dec.Paths[0]
	}
	b.OpenTask(dec.TaskID)
	if !b.prov.Mark(dec.TaskID, dec.Tool, origin, res.Text) {
		b.log("tools: %s result produced no matchable taint fragment (origin=%q)",
			dec.Tool, origin)
	}
}

// cancelText asks the D31 bus for one call's applied-steps report, and returns
// "" when there is nothing to report: no bus wired, no step applied and no
// veto recorded. A tool that applied steps WITHOUT saying so is why the bus,
// not the bridge, owns the wording - see approval.CancellationReport.
func (b *Bridge) cancelText(dec Decision, res Result) string {
	if b.cancel == nil {
		return ""
	}
	corr := orDefault(dec.CorrelationID, dec.TaskID)
	if len(res.AppliedSteps) == 0 && !b.cancel.Vetoed(corr) {
		return ""
	}
	rep := b.cancel.Report(dec, res)
	if rep == nil {
		return ""
	}
	txt := strings.TrimSpace(rep.String())
	if txt == "" {
		return ""
	}
	b.log("tools: d31 report corr=%s tool=%s applied=%d %s",
		corr, dec.Tool, len(res.AppliedSteps), txt)
	return txt
}

// ---------------------------------------------------------------------------
// C25 scope + C19 composition helpers
// ---------------------------------------------------------------------------

// factsFor assembles one call's C19 input. Declared and Paths belong to the
// BRIDGE - a tool must not be able to lower its own floor or hide a target from
// the judge. Everything else (R8's irreversibility classes, R7's batch count)
// comes from the host-side Decl.Facts hook, because only the tool can tell
// whether the target already exists or whether a move leaves its volume: the
// frozen rules are fed facts, not opinions. A hook that panics contributes one
// unknown irreversibility class, which R8 fail-closes to L2 rather than letting
// a broken judge read as "nothing to worry about".
func (b *Bridge) factsFor(ctx context.Context, entry Entry, params map[string]any,
	rawPaths []string,
) risk.Facts {
	var f risk.Facts
	if hook := entry.Decl.Facts; hook != nil {
		f = safeFacts(ctx, hook, params, rawPaths)
	}
	f.Declared = entry.Decl.Declared
	f.Paths = rawPaths
	return f
}

func safeFacts(ctx context.Context, hook func(context.Context, map[string]any, []string) risk.Facts,
	params map[string]any, rawPaths []string,
) (f risk.Facts) {
	defer func() {
		if rec := recover(); rec != nil {
			f = risk.Facts{Irreversible: []string{fmt.Sprintf("宿主事实钩子 panic: %v", rec)}}
		}
	}()
	return hook(ctx, params, rawPaths)
}

// OpenTask opens the C25 taint scope for one task - ticket 19's
// DEFERRED(C25-loop-wiring) item (1). Idempotent; Execute opens lazily so a
// caller that forgot cannot get an untainted read.
func (b *Bridge) OpenTask(taskID string) {
	if b.prov == nil || taskID == "" {
		return
	}
	b.mu.Lock()
	open := b.scopes[taskID]
	b.scopes[taskID] = true
	b.mu.Unlock()
	if !open {
		b.prov.OpenScope(taskID)
	}
}

// CloseTask closes one task's taint scope. The composition root defers this on
// the task's DisposalScope (tickets 12/28). A task that never closes leaks only
// its own indexed sources, which is the fail-closed direction: dropping marks
// would open R4.
func (b *Bridge) CloseTask(taskID string) {
	if b.prov == nil || taskID == "" {
		return
	}
	b.mu.Lock()
	open := b.scopes[taskID]
	delete(b.scopes, taskID)
	b.mu.Unlock()
	if open {
		b.prov.CloseScope(taskID)
	}
}

// assessorFor returns the assessor whose R4 input is bound to this task's
// scope. The frozen C19 TaintDetector seam takes no tool name
// (internal/risk/assessor.go:165-171), so the binding is per scope, hence per
// task; the canonicalizer and classifier are shared and immutable.
func (b *Bridge) assessorFor(taskID string) *risk.RiskAssessor {
	if b.injected || b.prov == nil || taskID == "" {
		return b.assess
	}
	return risk.NewRiskAssessor().
		WithCanonicalizer(b.paths).
		WithSensitiveClassifier(b.classifier).
		WithTaintDetector(b.prov.Detector(taskID))
}

// displayPaths resolves each raw path for the audit line and the card. The
// ASSessor is fed the RAW values instead, because R2/R3 own canonicalization
// through this same seam; handing a judge pre-resolved paths would let a second
// fold decide what the gate sees.
func (b *Bridge) displayPaths(raw []string) []string {
	if b.paths == nil || len(raw) == 0 {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		c, err := b.paths.Canonicalize(p)
		if err != nil {
			out = append(out, p+" (无法规范化: "+err.Error()+")")
			continue
		}
		out = append(out, c)
	}
	return out
}

// ---------------------------------------------------------------------------
// forensics: the tool_call row
// ---------------------------------------------------------------------------

// OutcomeKind classifies how a call left the bridge, for the audit line and
// the tests. It is NOT a storage enum: the columns stay in agent's frozen
// vocabulary (agent/journal.go:30-46).
type OutcomeKind string

// The bridge's dispositions.
const (
	OutcomeKindSuccess         OutcomeKind = "success"
	OutcomeKindError           OutcomeKind = "error"
	OutcomeKindCancelled       OutcomeKind = "cancelled"
	OutcomeKindTimeout         OutcomeKind = "tool-timeout"
	OutcomeKindRefused         OutcomeKind = "refused" // gate said no
	OutcomeKindApprovalTimeout OutcomeKind = "approval-timeout"
	OutcomeKindDenied          OutcomeKind = "denied" // verdict said no (Deny)
	OutcomeKindCapability      OutcomeKind = "capability"
	OutcomeKindNotFound        OutcomeKind = "not-found"
	OutcomeKindBadArgs         OutcomeKind = "bad-args"
)

// outcomeColumn maps a disposition onto tool_call.outcome (agent vocabulary).
func (k OutcomeKind) outcomeColumn() string {
	switch k {
	case OutcomeKindSuccess:
		return agent.OutcomeSuccess
	case OutcomeKindCancelled:
		return agent.OutcomeCancelled
	default:
		return agent.OutcomeError
	}
}

// errorClass picks the D37 class the row carries.
func (k OutcomeKind) errorClass() string {
	switch k {
	case OutcomeKindSuccess:
		return ""
	case OutcomeKindCapability, OutcomeKindDenied:
		return string(observe.ClassPermissionDenied)
	case OutcomeKindRefused, OutcomeKindApprovalTimeout:
		return string(observe.ClassUserRejected)
	case OutcomeKindCancelled:
		return string(observe.ClassCancelled)
	default:
		return string(observe.ClassTool)
	}
}

// reject books a refused call and returns what the model sees. Signature note:
// (outcome, nil) - see Execute's comment for why a rejection must never travel
// as a Go error.
func (b *Bridge) reject(ctx context.Context, req agent.ToolRequest, dec Decision,
	why, class string, kind OutcomeKind,
) (agent.ToolOutcome, error) {
	if dec.Reason == "" {
		dec.Reason = why
	}
	if dec.DecisionColumn == "" {
		dec.DecisionColumn = agent.DecisionReject
	}
	b.book(ctx, req, dec, kind)
	b.emit(dec)
	return agent.ToolOutcome{
		Text:       why,
		IsError:    true,
		RiskLevel:  dec.LevelString(),
		ErrorClass: class,
	}, nil
}

// close books a call that reached execution and returns its outcome.
func (b *Bridge) close(ctx context.Context, req agent.ToolRequest, dec Decision,
	out agent.ToolOutcome, kind OutcomeKind,
) (agent.ToolOutcome, error) {
	b.book(ctx, req, dec, kind)
	if out.RiskLevel == "" {
		out.RiskLevel = dec.LevelString()
	}
	if out.IsError && out.ErrorClass == "" {
		out.ErrorClass = kind.errorClass()
	}
	return out, nil
}

// emit publishes the structured verdict once per call.
func (b *Bridge) emit(d Decision) {
	if b.onDec != nil {
		b.onDec(d)
	}
}

// book writes the complete tool_call row for one finished call.
//
// SINGLE WRITER: the agent loop books its own pre-gate row today
// (loop.go:585 and journal.go:70-93). When a Bridge is composed under a Loop,
// the Loop must run with a nil Journal so the bridge owns the tool_call rows:
// the row below carries the ASSESSED level, the gate decision and the outcome,
// so a loop-side pre-gate row - which can only carry the declared lower bound -
// would be strictly less informative, not a duplicate worth keeping. Ticket 21
// rewrites loop.decideRisk into this call path; ticket 12 owns the composition.
//
// rules_hit: the frozen SPEC-02 §3 DDL for tool_call has no rules_hit column,
// and DDL text is contract-level (it is copied byte-for-byte from the spec), so
// the rule list travels in Decision (OnDecision -> the ticket 21/37 card, which
// renders it verbatim) and in the audit line below; the row keeps the columns it
// owns. Adding the column is a contract change, not this ticket's to make.
func (b *Bridge) book(ctx context.Context, req agent.ToolRequest, dec Decision, kind OutcomeKind) {
	b.log("tools: call kind=%s task=%s corr=%s tool=%s risk=%s decision=%s "+
		"outcome=%s rules_hit=%v in_allowlist_scope=%v reason=%q",
		kind, req.TaskID, orDefault(req.CorrelationID, req.TaskID), req.Name,
		dec.LevelString(), dec.DecisionColumn, kind.outcomeColumn(),
		dec.RulesHit, b.inScope(dec), dec.Reason)

	// Ticket 105 AC#2: this is the production reader of ticket 102's rewrite
	// account. Until here the only consumers of Roots()/RewrittenRoots()/
	// UnusableRoots() were tests, so "the operator can see which root C26 moved"
	// was a comment. It is now a record in the same sink as MODE-READ and the
	// line above (Options.Logf, wired to agentRuntime.auditf by cmd/wisp).
	// It goes out for EVERY call judged on a path, not only when the account has
	// something to complain about: empty lists are the evidence that the book was
	// read and came back clean, which is precisely what a reader cannot recover
	// from a missing line. The account is ticket 102's - no second book is kept
	// here, and the joined form is printed unescaped because the entries already
	// name the config spelling, the tree it moved to and the construct.
	if len(dec.Paths) > 0 {
		roots, rewritten, unusable := b.pathAccount()
		b.log("tools: PATH-ACCOUNT task=%s tool=%s roots=%d rewritten=[%s] unusable=[%s]",
			req.TaskID, req.Name, len(roots),
			strings.Join(rewritten, "; "), strings.Join(unusable, "; "))
	}

	if b.j == nil || req.TaskID == "" {
		return
	}
	decision := dec.DecisionColumn
	if decision == "" {
		decision = agent.DecisionAllow
		if kind != OutcomeKindSuccess {
			decision = agent.DecisionReject
		}
	}
	outcome := kind.outcomeColumn()
	errClass := kind.errorClass()
	tc := memory.ToolCall{
		TaskID:        req.TaskID,
		Seq:           b.nextSeq(req.TaskID),
		Tool:          req.Name,
		ArgsJSON:      sanitizeArgs(req.Args),
		RiskLevel:     dec.LevelString(),
		Decision:      decision,
		Outcome:       outcome,
		ErrorClass:    errClass,
		CorrelationID: orDefault(req.CorrelationID, req.TaskID),
	}
	if _, err := b.j.InsertToolCall(ctx, tc); err != nil {
		b.log("tools: tool_call insert failed task=%s tool=%s: %v", req.TaskID, req.Name, err)
	}
}

// pathAccount reads ticket 102's book: the roots in effect and the two things
// the C26 expansion step left behind on the way there (roots that moved onto
// another tree, roots that were dropped and authorize nothing). It exists so
// the audit record above names the same three lists the canonicalizer keeps,
// with no copy of its own to fall out of sync.
func (b *Bridge) pathAccount() (roots, rewritten, unusable []string) {
	if b.paths == nil {
		return nil, nil, nil
	}
	return b.paths.Roots(), b.paths.RewrittenRoots(), b.paths.UnusableRoots()
}

// inScope reports whether every judged path landed inside [fs] allowed_dirs -
// the audit line's cheap answer to "did the gate wave this through or did it
// actually belong here".
func (b *Bridge) inScope(dec Decision) bool {
	if b.paths == nil || len(dec.Paths) == 0 {
		return true
	}
	for _, p := range dec.Paths {
		if strings.Contains(p, "无法规范化") || !b.paths.InAllowlist(p) {
			return false
		}
	}
	return true
}

// nextSeq allocates one task's monotonic tool_call.seq.
func (b *Bridge) nextSeq(taskID string) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.seqs[taskID]++
	return b.seqs[taskID]
}

// sanitizeArgs keeps the forensics column useful without letting it carry an
// arbitrary payload (SPEC-02 §5.1 long-string truncation).
func sanitizeArgs(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if s == "" {
		return "{}"
	}
	const cap = 2048
	if len(s) <= cap {
		return s
	}
	cut := cap
	for cut > 0 && !utf8Start(s[cut]) {
		cut-- // never split a rune
	}
	return s[:cut] + "...[truncated]"
}

// utf8Start reports whether b begins a UTF-8 code point.
func utf8Start(b byte) bool { return b&0xC0 != 0x80 }

func (b *Bridge) log(format string, args ...any) {
	if b.logf == nil {
		return
	}
	b.logf(format, args...)
}

// ---------------------------------------------------------------------------
// small shared helpers
// ---------------------------------------------------------------------------

// decodeArgs decodes a call's arguments into the map shape C19 and C25 judge.
func decodeArgs(raw json.RawMessage) (map[string]any, error) {
	if len(strings.TrimSpace(string(raw))) == 0 {
		return map[string]any{}, nil
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("arguments are JSON null")
	}
	return m, nil
}

// pathArgs pulls the declared path parameters out of a call's arguments, in
// sorted key order so the judge's input is stable across dict iterations.
func pathArgs(params map[string]any, keys []string) []string {
	if len(keys) == 0 || len(params) == 0 {
		return nil
	}
	var out []string
	seen := map[string]bool{}
	for _, k := range keys {
		v, ok := params[k]
		if !ok {
			continue
		}
		for _, s := range stringValues(v) {
			if s == "" || seen[s] {
				continue
			}
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

// stringValues flattens a parameter value into strings (a scalar, or every
// string in a list). Non-string members are dropped: they are not paths.
func stringValues(v any) []string {
	switch t := v.(type) {
	case string:
		return []string{strings.TrimSpace(t)}
	case []any:
		var out []string
		for _, e := range t {
			if s, ok := e.(string); ok {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	default:
		return nil
	}
}

func levelString(l risk.Level) string {
	switch l {
	case risk.L1:
		return memory.RiskL1
	case risk.L2, risk.Deny:
		return memory.RiskL2
	default:
		return memory.RiskL0
	}
}

func orDefault(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}

// compile-time proof that the bridge is the seam the loop already calls.
var (
	_ agent.ToolProvider       = (*Bridge)(nil)
	_ risk.PathCanonicalizer   = (*PathCanonicalizer)(nil)
	_ risk.SensitiveClassifier = sensitiveClassifier{}
)
