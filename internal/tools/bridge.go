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
	// Grants is the READ side of the D45 scoped session grant (ticket 224), the
	// same shape as Modes and for the same reason: nil means "no grant has ever
	// been answered", which is the fail-closed direction, and the bridge only
	// reads it. Populating it is the native answer route's job
	// (approval.Gate's allow-session path over internal/session's ledger), and
	// nothing in this package can add to it.
	//
	// What it can and cannot cover is a contract line, not an implementation
	// choice: SPEC-06 §8.3's first bullet ("L2 永不进入任何持久授权（含会话级）")
	// and PLAN.md:2148 D45-3 mean this seam is consulted for an unsilenced L1
	// verdict and for nothing else. See Execute.
	Grants GrantSource
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
	// grants is ticket 224's third read-only injection, same shape and the same
	// rule: read per call, no setter, and nil means no session has ever
	// answered anything.
	grants GrantSource

	// classifier is shared: it is stateless over the frozen blacklist.
	classifier risk.SensitiveClassifier

	mu     sync.Mutex
	seqs   map[string]int64       // task id -> tool_call.seq
	scopes map[string]*risk.Scope // task id -> the handle that opened it (ticket 160: the ONLY thing that can close it)
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
		grants:   o.Grants,
		confirmations: func() map[string]bool {
			if o.Confirmations == nil {
				return nil // no file was ever confirmed: B tier stays "ask first"
			}
			return o.Confirmations()
		},
		classifier: NewSensitiveClassifier(),
		seqs:       map[string]int64{},
		scopes:     map[string]*risk.Scope{},
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
// hazard in agent.Loop.decideRisk (internal/agent/loop.go:783-806, measured at
// ticket 179's anchor; the pointer this comment carried, 719-738, had drifted),
// and why
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
// The returned error is reserved for host-internal faults, because the loop
// treats non-nil as "the provider broke" (internal/agent/loop.go:692-694) and
// books it through ErrorClassOfTurnError (internal/agent/guard.go:271): the
// class recorded is the one the error itself carries, and an error carrying
// no class falls back to D37 class "internal" (observe.ClassInternal). The
// hyphenated pair internal-provider is not one of D37's 17 classes - a
// provider fault is class provider (observe.ClassProvider) and reaches the
// row only if the error is wrapped with that class. Either way it is a class
// the model cannot self-correct against.
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

	// --- ticket 224: a live D45 session grant may stand in for ONE question ----
	//
	// The gate on this branch is the contract, not a preference. SPEC-06 §8.3's
	// first bullet ("L2 永不进入任何持久授权（含会话级）", restated at
	// PLAN.md:2148 D45-3) says an L2 verdict is never covered by a session grant,
	// and Deny is beyond everything. What is left is an L1 the mode did NOT
	// already silence - the one question a stored native click can legitimately
	// answer, because the click happened (approval's allow-session route spent a
	// live nonce on a card naming this tool and this path) and the answer was
	// scoped to this session. A mode that already silenced the call needs no
	// help from here, and must not be reported as if it had been granted.
	var grantID int64
	if !sil.Silenced && sil.Level == risk.L1 {
		grantID = b.sessionGrantID(ctx, req.Name, rawPaths)
		if grantID != 0 {
			b.log("tools: GRANT-USE tool=%s paths=%d grant_id=%d assessed=%s mode=%s "+
				"(an authorization answered by a native click earlier in THIS session, "+
				"not a click on this call)",
				req.Name, len(rawPaths), grantID, verdict.Level, mode)
		}
	}
	b.emit(dec)

	// --- routing: L0 pass, L1 window, L2 approval, Deny never reaches a gate -
	ok2, why := b.route(ctx, &dec, sil, grantID)
	if !ok2 {
		class := string(observe.ClassUserRejected)
		if dec.Level == risk.Deny {
			class = string(observe.ClassPermissionDenied)
		}
		return b.reject(ctx, req, dec, why, class, dispositionOf(dec))
	}

	// --- execution under the D38d ceiling and the C22 deadline --------------
	return b.run(ctx, req, entry, dec, grantID)
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
//
// grantID follows the same rule for the same reason (ticket 224): whether a
// stored session authorization answers this call is decided by the caller and
// handed in, never re-derived here, and it is only ever consulted on the L1
// branch. L2 and Deny have no path to it, which is SPEC-06 §8.3's first bullet
// expressed as a switch statement rather than as a check this function might
// forget.
func (b *Bridge) route(ctx context.Context, dec *Decision, sil risk.Silenced, grantID int64) (bool, string) {
	switch sil.Level {
	case risk.L0:
		dec.DecisionColumn = agent.DecisionAllow
		return true, ""

	case risk.Deny:
		dec.DecisionColumn = agent.DecisionReject
		return false, "该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：" + dec.Reason

	case risk.L1:
		if grantID != 0 {
			// Booked as allow_session_grant, not allow: the frozen vocabulary
			// already separates the two (agent/journal.go:30-36, and the same
			// string is in the tool_call.decision comment at memory/schema.go:76)
			// precisely so a reader can tell "a click answered this card" from
			// "an earlier click, scoped to this session, covered it".
			dec.DecisionColumn = agent.DecisionAllowGrant
			return true, ""
		}
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
//
// grantID is the D45 row that covered this call (0 when nothing did). It travels
// as a parameter down to book for one reason: tool_call.grant_id is a frozen
// column (SPEC-02 §3 / memory/schema.go:83), and "每次使用被授权通道都写
// tool_call 日志（可取证）" (SPEC-06 §8.3, D45-2) is only true if the row names
// which grant was spent. Decision cannot carry it - a grant field on that
// struct is exactly what internal/tools/ticket90_test.go:444 forbids.
func (b *Bridge) run(ctx context.Context, req agent.ToolRequest, entry Entry,
	dec Decision, grantID int64,
) (agent.ToolOutcome, error) {
	timeout := b.timeoutFor(entry)
	dec.Timeout = timeout

	// D38d ceiling. A cancelled task must not queue behind a full semaphore.
	slot := &inFlightSlot{}
	select {
	case b.sem <- struct{}{}:
		slot.take(b.sem)
		defer slot.giveBack()
	case <-ctx.Done():
		return b.close(ctx, req, dec, grantID, agent.ToolOutcome{
			Text:       "任务已取消，调用未执行",
			IsError:    true,
			ErrorClass: string(observe.ClassCancelled),
		}, OutcomeKindCancelled)
	}

	ectx, cancel := b.execContext(ctx, timeout)
	defer cancel()
	ectx = withInFlightSlot(ectx, slot)
	// D31: the running tool's only window onto the approval layer is its own
	// context. The handle carries the correlation id the veto is keyed on (and
	// the task id, the tool's task-level identity since the corr went per-call
	// in ticket 242), and Complete releases the gate's post-handoff record
	// when the call is over.
	ectx = withCancel(ectx, cancelHandle{
		corr: orDefault(req.CorrelationID, req.TaskID), taskID: req.TaskID, bus: b.cancel,
	})
	// Ticket 177 shape A: a fresh per-call box for the ONE path this call's
	// tool might itself write into its result text (task.output's re-read
	// pointer is the only declaration today). Per-call by construction, so an
	// exemption can never be keyed onto another call's output.
	ectx, hostPaths := withHostPathBox(ectx)
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
		return b.close(ctx, req, dec, grantID, agent.ToolOutcome{
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
		return b.close(ctx, req, dec, grantID, agent.ToolOutcome{
			Text:       "工具内部故障：" + err.Error(),
			IsError:    true,
			ErrorClass: string(observe.ClassTool),
		}, OutcomeKindError)
	}
	if errors.Is(ectx.Err(), context.DeadlineExceeded) && ctx.Err() == nil {
		b.log("tools: %s exceeded its %dms budget", req.Name, timeout.Milliseconds())
		return b.close(ctx, req, dec, grantID, agent.ToolOutcome{
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
		b.mark(dec, res, hostPaths.get())
	}
	return b.close(ctx, req, dec, grantID, out, kind)
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

// mark records C25 provenance for one result. hostPath is this call's shape-A
// declaration from the box below ("" when the tool wrote no host-minted path):
// risk only ever excludes that exact span inside THIS mark, and the marking
// gate itself (IsSensitiveSource) is unchanged. task.output joined that roster
// with ticket 175-r2, which is the first day a non-empty hostPath is reachable
// here at all: 177 shipped the carrier, nothing shipped the name. The
// bridge-level pair that keeps both halves honest - a foreign text naming the
// same artifact still hits R4, the host's own pointer stays re-readable - lives
// in internal/tools/ticket175r2_stamp_live_test.go.
func (b *Bridge) mark(dec Decision, res Result, hostPath string) {
	if b.prov == nil || dec.TaskID == "" || !risk.IsSensitiveSource(dec.Tool) {
		return
	}
	origin := res.Origin
	if origin == "" && len(dec.Paths) > 0 {
		origin = dec.Paths[0]
	}
	b.OpenTask(dec.TaskID)
	if !b.prov.MarkWithHostPath(dec.TaskID, dec.Tool, origin, res.Text, hostPath) {
		b.log("tools: %s result produced no matchable taint fragment (origin=%q)",
			dec.Tool, origin)
	}
}

// hostPathBox is the ticket 177 shape-A carrier between a tool that just wrote
// a concrete path into its own result text and the bridge that marks that text
// ("包内未导出的通道" - the vehicle 177-m1 measured as zero-red: no new field
// on the C1 Result, no new exported method anywhere, no per-scope roster).
// One box per call, hung on the call's context; the tool sets it, the bridge
// reads it once at marking time, and it dies with the call.
type hostPathBox struct {
	mu   sync.Mutex
	path string
}

type hostPathBoxKey struct{}

func withHostPathBox(ctx context.Context) (context.Context, *hostPathBox) {
	box := &hostPathBox{}
	return context.WithValue(ctx, hostPathBoxKey{}, box), box
}

// hostPathBoxFromCtx lets a tool declare its own host-minted path. A nil
// result (context built by a test or another caller, not by Bridge.execute)
// is not an error - the declaration is simply absent, and marking proceeds
// with nothing excluded.
func hostPathBoxFromCtx(ctx context.Context) *hostPathBox {
	box, _ := ctx.Value(hostPathBoxKey{}).(*hostPathBox)
	return box
}

func (h *hostPathBox) set(p string) {
	if h == nil || p == "" {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.path = p
}

func (h *hostPathBox) get() string {
	if h == nil {
		return ""
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.path
}

// inFlightSlot is ticket 222's carrier for the ONE D38d ceiling token a running
// call holds, hung on that call's context the way ticket 177 hangs its per-call
// box: no new exported name, no new field on any contract type, nothing outside
// this package can reach for it.
//
// Why it exists: run takes the token before entry.Tool.Execute and gave it back
// only when the tool returned. For every tool that is right - holding while
// executing is what "at most four tool calls at once" means. For task.spawn it
// was wrong: the call spends its whole life waiting for a child that needs one of
// those SAME four tokens to do its work, so every parent collected its own
// per-tool deadline while its child queued behind the slot the parent was holding
// (ticket 222's chain, rings 1-6; ledger A420). The frozen reason behind the
// ceiling (PLAN.md D38d) is one turn must not swamp the machine, and a call that
// is waiting runs nothing.
//
// The handoff is one-way: a call that gave its token back does not take another
// one before returning. Re-acquiring would put the parent straight back into the
// queue it just left, ending the wait with a subtler copy of the same block.
// Everything such a call still does after the wait is bookkeeping (roster row,
// C25 stamp, journal row), not capability execution, so nothing runs outside the
// ceiling: the invariant moves to the one D38d's wording actually states, at most
// four EXECUTIONS at once.
//
// It is deliberately not a public seam. Handing back is only honest for a call
// that does nothing but wait from that point on, and today exactly one tool
// qualifies. It is deliberately not a SubagentDeps field either: 197's
// Test197SubagentHasNoSelfApprovalOutlet enumerates that struct's field set as
// the "is this a second approval outlet?" check, and a slot handle would trip it.
type inFlightSlot struct {
	once sync.Once
	sem  chan struct{}
}

// take records the token run has just acquired. Only run calls it, on the
// goroutine that owns the slot, before that call's context is handed to the tool.
func (s *inFlightSlot) take(sem chan struct{}) { s.sem = sem }

// giveBack returns the token at most once and reports whether THIS call is the
// one that gave it up. The once is load-bearing twice over: run defers it and the
// tool may already have called it, and draining an empty semaphore would block
// the returning goroutine forever - the hang this ticket refuses to leave behind.
func (s *inFlightSlot) giveBack() bool {
	gave := false
	s.once.Do(func() {
		if s.sem != nil {
			<-s.sem
			s.sem = nil
			gave = true
		}
	})
	return gave
}

type inFlightSlotKey struct{}

func withInFlightSlot(ctx context.Context, s *inFlightSlot) context.Context {
	return context.WithValue(ctx, inFlightSlotKey{}, s)
}

// giveBackWhileWaiting is the whole seam a waiting tool needs: it hands the
// bridge token back for as long as the call only waits, and reports whether this
// call ever held one. A context not built by Bridge.run (a tool invoked
// directly, e.g. by a test or another host) holds no token and gets false, which
// is a statement about the caller, never an error - the tool waits exactly as it
// did before. It never blocks.
func giveBackWhileWaiting(ctx context.Context) bool {
	s, ok := ctx.Value(inFlightSlotKey{}).(*inFlightSlot)
	if !ok || s == nil {
		return false
	}
	return s.giveBack()
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
// DEFERRED(C25-loop-wiring) item (1), landed as a shape by ticket 160:
// risk.OpenScope now hands back the *Scope that owns the close, and the bridge
// holds on to it. Idempotent; Execute opens lazily so a caller that forgot
// cannot get an untainted read.
//
// What changed for the reader who was warned by the old text: the close side no
// longer lives in another package as a bare function anyone can call with any
// id. It lives in THIS map, as the handle this line stored. The residual the
// old text named is still true in one respect - Execute/close are separate
// calls, and a caller that opens a scope and never reaches CloseTask still
// leaks; G5 in .scratch/wisp/probes/154/gate-clauses.sh is the doorbell for
// that, not the compiler (measured, 160-handle-r1.md §3).
func (b *Bridge) OpenTask(taskID string) {
	if b.prov == nil || taskID == "" {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, open := b.scopes[taskID]; open {
		return
	}
	b.scopes[taskID] = b.prov.OpenScope(taskID)
}

// CloseTask closes one task's taint scope. The composition root defers this on
// the task's own boundary: cmd/wisp wires it into agent.Options.AdmitTask's
// revoke, which the TEXT loop already defers (internal/agent/loop.go:366), so a
// task that ends takes its taint with it (ticket 151). A caller that dispatches
// on the bridge outside that boundary - a host-internal call with its own task
// id - still has no owner here, and stays fail-closed: a task that never closes
// leaks only its own indexed sources, which is the fail-closed direction, since
// dropping marks would open R4.
//
// WHICH LEGS STILL HAVE NO OWNER HERE - ticket 154, read this before assuming
// ticket 151 finished the job. 151 put the close on ONE boundary: cmd/wisp's
// admitTask hook, which the TEXT loop defers. That covers only ids the loop
// mints itself (Run/RunAsync -> newTaskID, internal/agent/loop.go:332/:321).
// So ask one question about your own leg: does the TaskID you dispatch with
// come from that boundary? If it comes from anywhere else - a panel or ball
// host inventing its own id, a retry wrapper, a scheduled task carrying one id
// across turns - nothing in this tree closes it, and no test will tell you so.
// Today no production code dispatches on the bridge except loop.go, so being
// such a caller means being the first one; the gates meant to ring when that
// happens are clauses of .scratch/wisp/probes/154/gate-clauses.sh - G1/G1b for
// this bridge leg, and G5 for a C25 taint scope opened with no paired close in
// the same file, i.e. the leg that never crosses this bridge at all (G5 is
// ticket 158's addition; §2 of 158-gate-scope-blind-spot-r1.md is its write-up,
// §2.3 of 154-host-id-never-closed-r1.md the write-up of G1-G4). Each is one
// `git grep` away, not a note in someone's head.
//
// Being first is also an API decision, not just a missing defer: Loop has no
// exported method that takes a caller-supplied task id, so a host cannot route
// its id through the boundary that owns the close. That choice is Q-56 and this
// file does not answer it. When the shape does become reachable, closing your
// leg is only half of what ticket 154 owes - the reading for two tasks with
// scopes open at once is still zero, and it has to be measured in the same
// change (same evidence file, §2.5).
//
// The audit line runs on EVERY call, closed scope or not, because "who took
// this task's taint off the table, and was there anything on it" is exactly the
// question the C25 ledger owes a long-running host. A scope that was never
// opened (a task that read no sensitive source) closes as dropped=0.
func (b *Bridge) CloseTask(taskID string) {
	if b.prov == nil || taskID == "" {
		return
	}
	dropped := len(b.prov.ScopeTaints(taskID))
	b.mu.Lock()
	scope := b.scopes[taskID]
	open := scope != nil
	delete(b.scopes, taskID)
	left := len(b.scopes)
	b.mu.Unlock()
	var closeErr error
	if open {
		// The handle this bridge was issued is what closes the scope - there is
		// no by-id route left to call, and this one refuses a caller that does
		// not hold the handle. The refusal is a return value, so it is readable
		// here and on the audit line below, not only in a log the engine prints.
		closeErr = scope.Close()
	}
	b.log("tools: C25 scope closed task=%s was_open=%v dropped=%d open_scopes=%d close_err=%v",
		taskID, open, dropped, left, closeErr)
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
	// grantID is 0 by construction: a rejected call was never covered by a
	// session grant, and Execute reaches reject before the grant check on the
	// capability/args branches and after it on the routing branch, where route
	// only refuses when it did NOT use the grant.
	b.book(ctx, req, dec, 0, kind)
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
	grantID int64, out agent.ToolOutcome, kind OutcomeKind,
) (agent.ToolOutcome, error) {
	b.book(ctx, req, dec, grantID, kind)
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
//
// grant_id is NOT that missing column: tool_call has carried it since SPEC-02
// §3 froze the DDL (memory/schema.go:83, "关联 approval_grant（D45）"), and until
// ticket 224 nothing in production could fill it because nothing wrote an
// approval_grant row. It travels as a parameter for the same reason rules_hit
// cannot be a Decision field - the field-name scan at
// internal/tools/ticket90_test.go:444 forbids a grant-shaped member on Decision,
// and that ban is correct: the verdict a caller can set is the bug this
// repository keeps re-registering.
func (b *Bridge) book(ctx context.Context, req agent.ToolRequest, dec Decision,
	grantID int64, kind OutcomeKind,
) {
	b.log("tools: call kind=%s task=%s corr=%s tool=%s risk=%s decision=%s "+
		"outcome=%s rules_hit=%v in_allowlist_scope=%v grant_id=%d reason=%q",
		kind, req.TaskID, orDefault(req.CorrelationID, req.TaskID), req.Name,
		dec.LevelString(), dec.DecisionColumn, kind.outcomeColumn(),
		dec.RulesHit, b.inScope(dec), grantID, dec.Reason)

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
	if grantID != 0 {
		g := grantID
		tc.GrantID = &g
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
