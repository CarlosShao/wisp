package approval

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Options configures a Gate. Every zero value lands on the contract default;
// nothing here can switch an enforcement off.
type Options struct {
	// UI is the injected presentation surface. nil fails closed: no prompt
	// means no answer means no execution.
	UI UI
	// Clock is the monotonic time source (D42#9). nil uses SystemClock.
	Clock Clock
	// Channels states which veto channels are really loaded. nil means
	// DefaultChannels (ball + Esc; panel is ticket 37, KWS is ticket 41).
	Channels *ChannelRegistry
	// Window is the L1 pre-execution block length. Clamped into
	// [MinL1Window, MaxL1Window]: SPEC-06 §2 says 2-3s and a value outside
	// that range is a misconfiguration, not a policy choice.
	Window time.Duration
	// ApprovalTimeout is the C18 deadline (default 300s) after which the
	// answer is auto-REJECT.
	ApprovalTimeout time.Duration
	// WarningLead is how long before that deadline the prominent warning
	// fires (default 30s).
	WarningLead time.Duration
	// MaxPending bounds the queue; overflow fails closed.
	MaxPending int
	// MaxTracked bounds the post-handoff bookkeeping used by D31's
	// applied-steps report. Oldest entries are evicted; eviction loses the
	// late-veto note for an already-finished call, never an approval.
	MaxTracked int
	// Logf is the audit sink. A forgery attempt is always logged.
	Logf func(format string, args ...any)
}

// Gate is ticket 21's approval mechanics behind the internal/tools/gate.go
// seam: the L1 pre-execution block window with its four veto channels, the L2
// route into the C18 queue with its native-only allow proof, D45-1 batch
// aggregation, and D31's applied-steps report shape.
//
// It implements tools.Gate. It holds no reference to the bridge and no
// reference to the loop, so the only way into it is the routing the bridge
// already performs.
type Gate struct {
	ui       UI
	clock    Clock
	channels *ChannelRegistry
	window   time.Duration
	track    int
	q        *Queue
	logf     func(string, ...any)

	mu       sync.Mutex
	windows  map[string]*window
	running  map[string]*runEntry
	order    []string
	admitted map[string]bool
	admitOrd []string
}

// window is one live L1 block window.
type window struct {
	ctx    context.Context
	corr   string
	tool   string
	vetoes chan Veto
}

// runEntry is the post-handoff record: the call was allowed (or timed out) and
// execution started, so a veto arriving now cannot un-write what already ran.
type runEntry struct {
	corr string
	tool string
	veto *Veto // set by a veto that arrived after handoff
}

// compile-time proof this is the seam internal/tools/gate.go describes.
var _ tools.Gate = (*Gate)(nil)

// New composes a Gate. A nil UI is legal and means "everything above L0 is
// refused", which is the same posture NoGate had - so wiring this type in is
// never a downgrade.
func New(o Options) *Gate {
	clock := o.Clock
	if clock == nil {
		clock = SystemClock{}
	}
	logf := o.Logf
	if logf == nil {
		logf = func(string, ...any) {}
	}
	ui := o.UI
	if ui == nil {
		ui = UIFuncs{}
	}
	ch := o.Channels
	if ch == nil {
		ch = DefaultChannels()
	}
	win := o.Window
	switch {
	case win <= 0:
		win = DefaultL1Window
	case win < MinL1Window:
		win = MinL1Window
	case win > MaxL1Window:
		win = MaxL1Window
	}
	tracked := o.MaxTracked
	if tracked <= 0 {
		tracked = 64
	}
	return &Gate{
		ui: ui, clock: clock, channels: ch, window: win, track: tracked,
		q:        NewQueue(o.ApprovalTimeout, o.WarningLead, o.MaxPending, logf),
		logf:     logf,
		windows:  map[string]*window{},
		running:  map[string]*runEntry{},
		admitted: map[string]bool{},
	}
}

// Queue exposes the approval queue for the host's own views (badge depth,
// replay listing). It hands out no authority beyond the two API faces.
func (g *Gate) Queue() *Queue { return g.q }

// Channels exposes the veto-channel registry, so the KWS loader (ticket 41)
// can mark its channel live once the model is actually resident.
func (g *Gate) Channels() *ChannelRegistry { return g.channels }

// Window returns the configured L1 block length.
func (g *Gate) Window() time.Duration { return g.window }

// ---------------------------------------------------------------------------
// D47: the only path that may reach a gate is the TEXT loop
// ---------------------------------------------------------------------------

// AdmitTextTask registers one TEXT-loop task (Path T, or a Path C conversation
// that handed its work back through C32) as allowed to reach the gates, and
// returns the revoke func the caller defers on the task's DisposalScope.
//
// This is how D47 is enforced as data rather than as a comment: the realtime
// brain has zero tool permissions, so nothing on that path ever holds the
// value that calls AdmitTextTask, and nothing that does not hold it can get a
// task registered. An approval request for an unregistered task is refused
// before the queue is touched, so a C-side request cannot even become a pending
// item. The admission decision is not taken from the request's own fields - a
// request claiming to be from the text loop says nothing about where it came
// from (rulings M-7/C-3 again).
func (g *Gate) AdmitTextTask(taskID string) (revoke func()) {
	if strings.TrimSpace(taskID) == "" {
		return func() {}
	}
	g.mu.Lock()
	g.admitted[taskID] = true
	g.admitOrd = append(g.admitOrd, taskID)
	for len(g.admitOrd) > maxAdmittedTasks {
		drop := g.admitOrd[0]
		g.admitOrd = g.admitOrd[1:]
		delete(g.admitted, drop)
	}
	g.mu.Unlock()
	var once sync.Once
	return func() {
		once.Do(func() {
			g.mu.Lock()
			delete(g.admitted, taskID)
			for i, t := range g.admitOrd {
				if t == taskID {
					g.admitOrd = append(g.admitOrd[:i], g.admitOrd[i+1:]...)
					break
				}
			}
			g.mu.Unlock()
		})
	}
}

// maxAdmittedTasks bounds the registry so a long-running host cannot grow it
// without limit if a revoke is ever missed. Eviction is fail-closed: a task
// evicted while still running gets refused, not waved through.
const maxAdmittedTasks = 256

// admitCheck is D47's guard on both routes.
func (g *Gate) admitCheck(taskID string) error {
	if strings.TrimSpace(taskID) == "" {
		return fmtw(ErrNotAdmitted, "审批请求未携带任务标识，已拒绝")
	}
	g.mu.Lock()
	ok := g.admitted[taskID]
	g.mu.Unlock()
	if !ok {
		return fmtw(ErrNotAdmitted,
			"实时语音（Path C）不具备工具权限；该任务未经文本循环登记（AdmitTextTask），已 fail-closed 拒绝")
	}
	return nil
}

// ---------------------------------------------------------------------------
// L1: the pre-execution block window (SPEC-06 §2, B1)
// ---------------------------------------------------------------------------

// PendingWindow implements tools.Gate: the 2-3s pre-execution BLOCK (not an
// undo - B1). Four veto channels are published with their real availability;
// the window ending unopposed means EXECUTE (AnswerTimeout, which the bridge
// routes to allow), and a veto means the call is cancelled with the reason
// travelling back to the model as a failure so the loop continues.
func (g *Gate) PendingWindow(ctx context.Context, d tools.Decision) (tools.Answer, string) {
	if err := g.admitCheck(d.TaskID); err != nil {
		g.logf("approval: L1 refused by D47 guard task=%q tool=%s: %v", d.TaskID, d.Tool, err)
		return tools.AnswerReject, err.Error()
	}

	// R7 belt-and-braces: the assessor is supposed to have raised a >=50-file
	// batch to L2 already. If a caller hands this window a batch that big
	// while claiming L1, the batch is NOT aggregated and is escalated to the
	// L2 route instead of getting one cheap confirm.
	if n := len(d.Paths); n >= R7BatchFloor {
		g.logf("approval: R7 escalation - %d paths arrived on the L1 route (tool=%s); sending to L2 approval, unaggregated", n, d.Tool)
		if d.Level != risk.L2 {
			d.Level = risk.L2
			d.RulesHit = append(append([]risk.RuleID(nil), d.RulesHit...), risk.R7)
			d.Reason = orDefaultText(d.Reason,
				fmt.Sprintf("本次调用涉及 %d 个目标，达到 R7 批量规模阈值，已按 L2 处理", n))
		}
		a, why := g.PendingApproval(ctx, d)
		if a == tools.AnswerTimeout {
			// The L2 polarity is the opposite of L1's: an unanswered approval
			// is a REJECT. Returning the raw timeout here would let the
			// bridge's L1 branch read it as "execute".
			return tools.AnswerReject, "R7 升级后的 L2 审批未获答复：" + why
		}
		return a, why
	}

	bv := Aggregate(d)
	corr := orDefaultText(d.CorrelationID, d.TaskID)
	w := &window{ctx: ctx, corr: corr, tool: d.Tool, vetoes: make(chan Veto, len(allChannels)+1)}
	if !g.openWindow(w) {
		return tools.AnswerReject, "同一 correlation_id 已有确认窗口在跑，无法重复登记，已 fail-closed 拒绝"
	}
	defer g.closeWindow(w.corr)

	// Armed BEFORE the prompt is delivered: the countdown belongs to the
	// window, not to the UI's repaint, and arming after Prompt() would make
	// "the strip took 4s to draw" silently extend a 3s block.
	deadline := g.clock.After(g.window)

	p := g.promptFor(d, corr, g.window, 0, bv, "")
	if err := g.ui.Prompt(ctx, p); err != nil {
		// Host unreachable: fail closed. "Could not ask" is never "assumed ok".
		g.logf("approval: L1 prompt failed corr=%s tool=%s: %v", corr, d.Tool, err)
		return tools.AnswerReject, "确认界面不可达（" + err.Error() + "），已 fail-closed 拒绝执行"
	}

	for {
		select {
		case v := <-w.vetoes:
			// Re-checked here, at the point where the veto could actually
			// cancel: the registry, not the caller, decides availability.
			if err := g.channels.check(v.Channel); err != nil {
				g.logf("approval: veto on unloaded channel corr=%s: %v", corr, err)
				_ = g.ui.Update(ctx, Event{
					Kind: EventWarning, CorrelationID: corr, Remaining: g.window,
					Text: err.(*ChannelError).Msg,
				})
				continue
			}
			why := fmt.Sprintf("用户在 L1 确认窗口中通过「%s」否决了本次操作（取消窗口，非撤销已写出的内容）",
				channelNames[v.Channel])
			_ = g.ui.Update(ctx, Event{Kind: EventDismissed, CorrelationID: corr, Text: why})
			return tools.AnswerVeto, why

		case <-deadline:
			// Timeout MEANS EXECUTE on the L1 route (SPEC-06 §2). Execution has
			// now started, so cancellation is no longer atomic (D31): the
			// correlation stays tracked so a later veto yields an applied-steps
			// report instead of a false "nothing happened".
			g.markStarted(corr, d.Tool)
			_ = g.ui.Update(ctx, Event{
				Kind: EventStarted, CorrelationID: corr,
				Text: "确认窗口结束，未收到否决，开始执行",
			})
			return tools.AnswerTimeout, ""

		case <-ctx.Done():
			return tools.AnswerReject, "任务已取消，L1 确认窗口未放行，未执行"
		}
	}
}

func (g *Gate) openWindow(w *window) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, clash := g.windows[w.corr]; clash {
		return false
	}
	g.windows[w.corr] = w
	return true
}

func (g *Gate) closeWindow(corr string) {
	g.mu.Lock()
	delete(g.windows, corr)
	g.mu.Unlock()
}

// markStarted moves a correlation from "asking" to "running" for D31 tracking.
// The set is bounded by g.track; an entry holding a late veto is kept so the
// applied-steps report cannot lose the fact that a veto arrived.
func (g *Gate) markStarted(corr, tool string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, dup := g.running[corr]; !dup {
		g.order = append(g.order, corr)
	}
	g.running[corr] = &runEntry{corr: corr, tool: tool}
	for len(g.order) > g.track {
		drop := g.order[0]
		g.order = g.order[1:]
		if e, ok := g.running[drop]; ok && e.veto == nil {
			delete(g.running, drop)
		}
	}
}

// Complete drops one call's post-handoff record; a host that uses the Bus
// defers this when the tool returns.
func (g *Gate) Complete(corr string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.running, corr)
	for i, c := range g.order {
		if c == corr {
			g.order = append(g.order[:i], g.order[i+1:]...)
			break
		}
	}
}

// ---------------------------------------------------------------------------
// veto entry (ball click / Esc / panel / KWS all call this)
// ---------------------------------------------------------------------------

// ErrAlreadyStarted means a veto arrived after execution began. It is not a
// fault: the call keeps running to completion and reports its applied steps
// (D31). The veto is recorded so the running tool can see it through Bus.
var ErrAlreadyStarted = errors.New("approval: 执行已开始，取消不是原子的（D31），将报告已生效步骤")

// Veto cancels the L1 window named by v.CorrelationID.
//
// It returns ErrChannelUnavailable (wrapping the exact 「语音取消不可用」 style
// wording) for a channel the host never loaded, ErrUnknownCorrelation when the
// id matches no live window (C18 routing: a click cannot land on someone
// else's request), and ErrAlreadyStarted when the call already went live.
func (g *Gate) Veto(v Veto) error {
	err := g.channels.check(v.Channel)
	if strings.TrimSpace(v.CorrelationID) == "" {
		return ErrUnknownCorrelation
	}
	g.mu.Lock()
	w := g.windows[v.CorrelationID]
	g.mu.Unlock()
	if err != nil {
		// The user just tried to cancel by a channel that cannot fire. The
		// answer goes back to the SAME live window (B1: 「语音取消不可用」 must
		// be said, never left implied), and the window keeps counting down.
		if w != nil {
			ce := &ChannelError{}
			if errors.As(err, &ce) {
				_ = g.ui.Update(w.ctx, Event{
					Kind: EventWarning, CorrelationID: w.corr, Remaining: g.window,
					Text: ce.Msg,
				})
			}
		}
		return err
	}
	g.mu.Lock()
	r := g.running[v.CorrelationID]
	if r != nil {
		cp := v
		r.veto = &cp
	}
	g.mu.Unlock()
	if w != nil {
		select {
		case w.vetoes <- v:
			return nil
		default:
			return errors.New("取消请求积压，上一条否决尚未生效")
		}
	}
	if r != nil {
		return ErrAlreadyStarted
	}
	return ErrUnknownCorrelation
}

// LateVeto reports whether a started call was vetoed after it began (the
// running tool's poll; also available through Bus).
func (g *Gate) LateVeto(corr string) (Veto, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	r := g.running[corr]
	if r == nil || r.veto == nil {
		return Veto{}, false
	}
	return *r.veto, true
}

// ---------------------------------------------------------------------------
// L2: the C18 queue, native-source-only allow
// ---------------------------------------------------------------------------

// PendingApproval implements tools.Gate: the L2 route. The request is queued
// with a correlation id, displayed once through the injected UI together with
// the single-use grant only the native side receives, and then waits under the
// C18 deadline. An unanswered 300s always resolves to REJECT; a UI that cannot
// be reached resolves to REJECT immediately; the task's own context is only
// ever observed, never cancelled by this layer.
func (g *Gate) PendingApproval(ctx context.Context, d tools.Decision) (tools.Answer, string) {
	if err := g.admitCheck(d.TaskID); err != nil {
		g.logf("approval: L2 refused by D47 guard task=%q tool=%s: %v", d.TaskID, d.Tool, err)
		return tools.AnswerReject, err.Error()
	}
	it, err := g.q.push(d)
	if err != nil {
		// Queue broken/full -> fail closed (ticket 21).
		return tools.AnswerReject, "审批队列不可用：" + err.Error()
	}
	nonce, err := g.q.grantNonce(it)
	if err != nil {
		return tools.AnswerReject, "审批令牌无法签发：" + err.Error()
	}
	corr := it.Corr
	d.CorrelationID = corr

	// Armed before display, for the same reason as the L1 countdown: the C18
	// deadline starts when the request is queued.
	deadline := g.clock.After(g.q.Timeout())
	var warn <-chan time.Time
	if lead := g.q.Timeout() - g.q.WarningLead(); lead > 0 {
		warn = g.clock.After(lead)
	}

	// L2 is NEVER aggregated (D45-1): the batch argument is dropped here on
	// purpose, so no future refactor can hand an L2 card a collapsed summary.
	p := g.promptFor(d, corr, 0, g.q.Timeout(), nil, nonce)
	p.Depth = g.q.Depth()
	if err := g.ui.Prompt(ctx, p); err != nil {
		g.q.abandon(it, "审批界面不可达（"+err.Error()+"），已 fail-closed 拒绝执行")
		g.logf("approval: L2 prompt failed corr=%s: %v", corr, err)
		return tools.AnswerReject, "审批界面不可达，已 fail-closed 拒绝执行"
	}

	for {
		select {
		case a := <-it.answer:
			_ = g.ui.Update(ctx, Event{Kind: EventDismissed, CorrelationID: corr, Text: a.why})
			return a.a, a.why

		case <-warn:
			// C18: 「超时前 30s 醒目提示」. Prominent, and it does not shorten
			// the deadline.
			w := Event{
				Kind: EventWarning, CorrelationID: corr,
				Remaining: g.q.WarningLead(),
				Text: fmt.Sprintf("审批将在 %d 秒后自动拒绝，请尽快确认",
					int(g.q.WarningLead().Seconds())),
			}
			if err := g.ui.Update(ctx, w); err != nil {
				g.logf("approval: warning delivery failed corr=%s: %v", corr, err)
			}
			warn = nil // fire once

		case <-deadline:
			a := g.q.expire(it)
			_ = g.ui.Update(ctx, Event{Kind: EventDismissed, CorrelationID: corr, Text: a.why})
			return a.a, a.why

		case <-ctx.Done():
			// The waiting call went away. Refuse, and do NOT touch the task
			// root context (C18: 拒绝后任务 root ctx 不取消).
			a := g.q.abandon(it, "任务上下文已结束，审批请求已作废并按拒绝处理")
			return a.a, a.why
		}
	}
}

// ---------------------------------------------------------------------------
// prompt construction (shared by both routes)
// ---------------------------------------------------------------------------

// promptFor assembles the card content. Full params, verbatim rules_hit and
// reason, honest channel availability, and the grant only where one exists.
//
// grant is the single argument that carries native authority, and there is
// exactly one call site that passes a non-empty value: PendingApproval, with
// the nonce minted seconds earlier in the same function. The L1 route passes
// "" because a countdown window has no allow action at all - its only user
// verb is veto, and refusing never needs proof.
func (g *Gate) promptFor(d tools.Decision, corr string, win, deadline time.Duration, bv *BatchView, grant string) Prompt {
	return Prompt{
		CorrelationID: corr,
		TaskID:        d.TaskID,
		Tool:          d.Tool,
		Level:         d.LevelString(),
		Params:        d.Params,
		Args:          d.Args,
		RulesHit:      append([]risk.RuleID(nil), d.RulesHit...),
		Reason:        d.Reason,
		Paths:         append([]string(nil), d.Paths...),
		Batch:         bv,
		Window:        win,
		Deadline:      deadline,
		Channels:      g.channels.Statuses(),
		Grant:         grant,
	}
}

// ---------------------------------------------------------------------------
// API faces
// ---------------------------------------------------------------------------

// Native returns the surface the ball click handler, the native card buttons
// and the global hotkey hold.
func (g *Gate) Native() NativeAPI { return nativeAPI{g} }

// Panel returns the surface the panel host (ticket 37) may be handed. Note
// what is NOT on it: no Allow.
func (g *Gate) Panel() PanelAPI { return panelAPI{g.q} }

type nativeAPI struct{ g *Gate }

func (n nativeAPI) Allow(_ context.Context, corr, grant string) error {
	return n.g.q.allow(corr, grant)
}

func (n nativeAPI) Reject(corr, reason string) error { return n.g.q.reject(corr, reason) }

type panelAPI struct{ q *Queue }

func (p panelAPI) Reject(corr, reason string) error { return p.q.reject(corr, reason) }
func (p panelAPI) Head() (PanelItem, bool)          { return p.q.head() }
func (p panelAPI) View(corr string) (PanelItem, bool) {
	return p.q.view(corr)
}

// Request is one decision as it arrives from a transport that carries both
// kinds of answer (a future PanelBridge, a test, a host router).
//
// Source is ADVISORY ONLY. It is set by the caller, which is precisely why no
// branch in this package reads it: an authority check keyed on a caller-owned
// field is the M-7/C-3 bug wearing a new hat. It exists so the audit line can
// record what a caller claimed next to what the route proved.
type Request struct {
	CorrelationID string
	Allow         bool
	Grant         string
	Reason        string
	Source        string
}

// DecideFromNative is the router for the native transport.
func (g *Gate) DecideFromNative(ctx context.Context, r Request) error {
	g.logf("approval: native route corr=%s allow=%v claimed_source=%q", r.CorrelationID, r.Allow, r.Source)
	if !r.Allow {
		return g.q.reject(r.CorrelationID, r.Reason)
	}
	return g.q.allow(r.CorrelationID, r.Grant)
}

// DecideFromPanel is the server-side API the panel reaches (SPEC-06 §9 layer
// 3). It rejects an allow on the ROUTE alone: a grant presented here is
// refused even if it is real and live, and a Source claiming "native" changes
// nothing because the field is never read for authority.
func (g *Gate) DecideFromPanel(ctx context.Context, r Request) error {
	if r.Allow {
		g.logf("approval: PANEL-ALLOW-REJECTED corr=%s claimed_source=%q grant_offered=%v",
			r.CorrelationID, r.Source, r.Grant != "")
		if r.Grant != "" {
			// Burn it: an offer to allow from an untrusted route is evidence of
			// a compromised or confused panel; the nonce must not survive to be
			// tried again somewhere else.
			_ = g.q.allow(r.CorrelationID, r.Grant+"!burned")
		}
		return fmtw(ErrPanelAllow, "面板来源的「允许」被服务端 API 直接拒绝（F2 第三层），可改为拒绝或查看完整参数")
	}
	return g.q.reject(r.CorrelationID, r.Reason)
}

// Replay re-displays a refused or expired L2 request under a fresh
// correlation id (C18 一键重放). It requires the task to be admitted again
// (D47 still applies to a replay), and it is NOT an answer: the new item must
// still be approved with a newly minted grant.
func (g *Gate) Replay(ctx context.Context, corr, taskID string) (tools.Decision, string, error) {
	if err := g.admitCheck(taskID); err != nil {
		return tools.Decision{}, "", err
	}
	d, fresh, err := g.q.replay(corr)
	if err != nil {
		return tools.Decision{}, "", err
	}
	d.TaskID = taskID
	g.logf("approval: replay %s -> %s tool=%s", corr, fresh, d.Tool)
	_ = ctx
	return d, fresh, nil
}

// ---------------------------------------------------------------------------
// small shared helpers
// ---------------------------------------------------------------------------

// wrappedError lets a sentinel carry a user-visible sentence.
type wrappedError struct {
	msg string
	err error
}

func (w *wrappedError) Error() string { return w.msg }
func (w *wrappedError) Unwrap() error { return w.err }

func fmtw(err error, format string, args ...any) error {
	return &wrappedError{msg: sprintf(format, args...), err: err}
}

func sprintf(format string, args ...any) string { return fmt.Sprintf(format, args...) }

func orDefaultText(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}
