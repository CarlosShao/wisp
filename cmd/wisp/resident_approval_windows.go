//go:build windows

package main

// The approval gate held by the resident process (ticket 246, 乙形：装配根注入).
//
// WHY THIS FILE EXISTS. Ticket 245 made "borrow Esc while a card is waiting,
// hand it back when the card is gone" the right shape, and ticket 228 moved the
// ball into the no-args resident process - but the two halves never met: the leg
// that runs tasks (`wisp run`) had the gate and no ball, and the leg that holds
// the ball had no gate at all. The measurement behind that statement is in
// .scratch/wisp/probes/246/a1/census.md §1 and it was re-run by the orchestrator
// (ledger A480 ②/A481 ③): TakeEscForCancel/ReleaseEscAfterSession had zero
// callers in any product process, and cmd/wisp/resident_ball_windows.go said so
// about itself in three places. So D43's veto row had no executor in the process
// the user actually launches.
//
// WHAT WAS CHOSEN, and by whom. Orchestrator ruling A481 picked form 乙: the
// assembly root (this package) builds the gate and hands it down, rather than the
// resident file assembling a second one. Two rulings from this repository are
// extended to this ticket by that choice, and the extension is named here rather
// than implied: ticket 238's cut-1 ("两枚正向依赖边一律不开、改注入") and ticket
// 197's "装配根 cmd/wisp 是唯一的接缝". This is an application of those rulings
// to a new case, not a rule that already covered it.
//
// WHAT THIS IS NOT.
//   - Not an agent loop, not a task pipeline, not a microphone. Nothing here
//     talks to a model. The cards this gate shows are host-initiated, the same
//     shape cmd/wisp/run.go's confirmModeSwitch and approval_always.go's
//     runWidening already use: a real admission (Gate.AdmitTextTask), a real
//     queue item or a real L1 window, a real single-use grant.
//   - Not a panel. This process links no WebView2 host (see the sentence at
//     approval_always.go:165), so "the card is on screen" is NOT a claim this
//     file makes. What it can and does claim: the gate opened a real pending
//     item, the orb's rendered state and Replies.WaitingState() moved to
//     Confirming, the cancel hot key was really taken from the desktop for the
//     length of the window, and the audit line says which of those happened.
//   - Not the other three veto channels. ChannelBall, ChannelKWS and
//     ChannelPanel stay unloaded in this process (ticket 246 AC#6); the gate's
//     own registry is what says so, and a veto arriving on an unloaded channel is
//     refused by that registry, not by a branch in this file.

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

// The transport names this host stamps on the audit line. Labels, not authority:
// approval/gate.go reads Request.Source nowhere, and no branch in this file
// either. They exist so a reader of the ledger can tell the resident leg's native
// surface from the console leg's.
const (
	residentNativeSource = "cmd-wisp-resident-ball-native"
	residentPanelSource  = "cmd-wisp-resident-panel-route"
)

// residentCardTool names a host-initiated card on the card and in the audit. It
// is not a registered C4 tool - no call route uses it - which is the same
// convention run.go:712 set with modeSwitchToolName.
const residentCardTool = "resident.confirmation"

// errResidentNoBall is what the injected UI returns when this process has no
// window to show a card on. The gate treats a UI error as "the host could not be
// reached" and fails closed, so a resident leg whose ball never came up refuses
// every confirmation instead of asking a question nobody can see.
var errResidentNoBall = errors.New("常驻进程没有悬浮球窗口，卡片无处呈现")

// residentApproval is this process's approval gate: the handle the assembly root
// builds and injects into the ball host, plus the task root the D38(e) cancel
// step owns.
type residentApproval struct {
	gate  *approval.Gate
	cards *approval.Replies
	ui    *ballCardUI

	mu   sync.Mutex
	root context.Context
	// cancel is the D38(e) step 3 action: it is what an L1 window and an L2
	// wait answer with when the process is leaving.
	cancel  context.CancelFunc
	live    int // raises between "asked" and "answered"
	closed  bool
	escLoad bool // the Esc channel was loaded (only ever with a real ball)
}

// newResidentApproval composes the gate the resident leg runs on. Every piece of
// it is the machinery `wisp run` already uses (run.go:523's approval.New with an
// injected UI and a bound Replies ledger); what differs is the UI implementation
// and the fact that this one outlives a single command line.
//
// The channel registry starts with NOTHING loaded. ChannelEsc is marked live
// only by bindBallHost, and only when Win32 actually gave this process a ball
// window - the shape 246-a1's D-5 warned about is a gate advertising a cancel key
// in a process that has no key to borrow.
func newResidentApproval() *residentApproval {
	root, cancel := context.WithCancel(context.Background())
	ra := &residentApproval{root: root, cancel: cancel}
	ra.ui = &ballCardUI{ra: ra}
	ra.gate = approval.New(approval.Options{
		UI:       ra.ui,
		Channels: approval.NewChannels(),
		Logf:     ra.residentAuditf,
	})
	ra.cards = approval.NewReplies()
	ra.cards.Attach(approval.HostBinding{
		Gate:         ra.gate,
		VetoChannel:  approval.ChannelEsc,
		NativeSource: residentNativeSource,
		PanelSource:  residentPanelSource,
	})
	return ra
}

// auditf is this gate's audit sink, in the same family the run leg writes: the
// gate's own lines (ANSWER-ALLOW / ANSWER-REJECT / ANSWER-VETO / ANSWER-EXPIRED)
// arrive here, and so do the two sentences only this file can say. stderr alone
// would be the pre-ticket-105 posture that R-105-1 called out, so the line goes
// through the installed sink as well.
func (ra *residentApproval) residentAuditf(format string, args ...any) {
	line := fmt.Sprintf(format, args...)
	slog.Info("audit: " + line)
	fmt.Printf("wisp: [audit] %s\n", line)
}

// bindBallHost hands the ball to the injected UI and, with it, loads the one
// veto channel this leg can honestly advertise. rb with no ball behind it is a
// real state (a machine with no desktop): the gate then keeps every channel
// unloaded and says why once.
func (ra *residentApproval) bindBallHost(rb *residentBall) bool {
	if rb == nil || rb.b == nil {
		ra.residentAuditf("resident-approval: 审批门已装配，但本进程没有悬浮球窗口，" +
			"四条否决通道全部保持未加载（Esc 无处可借，卡片无处可呈）")
		return false
	}
	ra.mu.Lock()
	ra.escLoad = true
	ra.mu.Unlock()
	ra.ui.attachBall(rb.b)
	ra.gate.Channels().SetLoaded(approval.ChannelEsc, true)
	ra.residentAuditf("resident-approval: 审批门已装配进常驻进程，取消通道 Esc 已加载（本票只落 Esc 一条通道；" +
		"单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）")
	return true
}

// vetoByEsc is the callback injected into the ball host: the resident leg's
// answer to "the user pressed the cancel key". It returns the sentence this
// process prints and books, so the ball host stays free of any knowledge of
// approvals - it only knows somebody handed it a function for that one gesture.
//
// It runs ON the ui-sta thread (internal/ball fires its events inline, and the
// contract there is "quick and non-blocking"), so it waits on nothing: the L1
// window's veto channel is buffered and Gate.Veto only ever sends.
func (ra *residentApproval) vetoByEsc() string {
	card, awaiting := ra.cards.AwaitingHuman()
	if !awaiting {
		// Said out loud rather than left silent: "the key was taken and there
		// was nothing to cancel" is the reading that tells a user the borrow and
		// the card are two different lifetimes.
		return "按下的取消键没有可否决的确认项：本进程此刻没有卡片在等人"
	}
	if err := ra.cards.Veto(card.CorrelationID); err != nil {
		return fmt.Sprintf("按 Esc 否决卡片 %s 被拒：%v", card.CorrelationID, err)
	}
	return fmt.Sprintf("按 Esc 否决了卡片 %s（%s / %s）：该调用未执行，答案已入审计",
		card.CorrelationID, card.Level, card.Tool)
}

// residentCard is one host-initiated confirmation: the whole of what this file
// can say about a card without a model in front of it.
type residentCard struct {
	TaskID string
	Tool   string
	Level  risk.Level
	Reason string
	Paths  []string
}

// ask raises a REAL confirmation through the gate and blocks until the route
// answers it - the same blocking shape a tool call takes through
// internal/tools' bridge, minus the bridge. L1 opens the pre-execution window
// (D43's Confirming, 2-3s, veto-only); L2 enters the C18 queue (D43's
// AwaitingApproval, auto-reject at the deadline).
//
// D47 is respected the way run.go respects it for a host card: the task is
// admitted for the length of the confirmation and revoked on the way out, so an
// unregistered task cannot reach this gate at all.
func (ra *residentApproval) askConfirmation(ctx context.Context, c residentCard) (tools.Answer, string) {
	if c.TaskID == "" || c.Reason == "" {
		return tools.AnswerReject, "宿主发起的确认必须自报任务标识与理由，已 fail-closed 拒绝"
	}
	ra.mu.Lock()
	if ra.closed {
		ra.mu.Unlock()
		return tools.AnswerReject, "审批门已随退出序列封闭，本进程不再签发任何确认"
	}
	ra.live++
	ra.mu.Unlock()
	defer func() {
		ra.mu.Lock()
		ra.live--
		ra.mu.Unlock()
		// Every return path of the two routes below ends the wait, including the
		// one where the task context was cancelled - and THAT one sends the UI no
		// dismissal event at all (gate.go returns straight away on ctx.Done). A
		// card that left without settling here would leave the orb in Confirming
		// and the desktop without its Esc key, which is precisely the residue
		// ticket 245 was filed to remove. settle() is idempotent and re-checks
		// the ledger, so a concurrent card still waiting is never disturbed.
		ra.ui.settleOrb()
	}()

	revoke := ra.gate.AdmitTextTask(c.TaskID)
	defer revoke()

	d := tools.Decision{
		TaskID: c.TaskID,
		Tool:   c.Tool,
		Level:  c.Level,
		Reason: c.Reason,
		Paths:  c.Paths,
	}
	switch c.Level {
	case risk.L1:
		return ra.gate.PendingWindow(ctx, d)
	case risk.L2:
		return ra.gate.PendingApproval(ctx, d)
	default:
		// L0 asks nobody, which is not what a caller of this function means.
		// Refusing an unlevelled request is cheaper than inventing a route for it.
		return tools.AnswerReject, "宿主发起的确认必须指名 L1 或 L2 路线，已 fail-closed 拒绝"
	}
}

// AskOnTaskRoot raises one card against this process's own task root - the
// context D38(e) step 3 cancels on the way out. It is the seam a future task
// source (voice chain, ticket 228's config wiring) will call; today its only
// callers are this package's own cases, and this file says so in the boot report
// rather than letting the existence of the seam read as a running pipeline.
//
// There is deliberately no caller-supplied context: a card that could be detached
// from the root would be a card the exit sequence cannot reach, which is the
// exact shape AC#4 exists to close.
func (ra *residentApproval) AskOnTaskRoot(c residentCard) (tools.Answer, string) {
	taskCtx, cancelTask := context.WithCancel(ra.root)
	defer cancelTask()
	return ra.askConfirmation(taskCtx, c)
}

// cancelTasks IS D38(e) step 3 for this process ("all task root ctxs cancelled
// -> wait <= 3s"). It is registered through proc.Runtime.RegisterShutdownHook,
// so the audit trail says this step ran instead of recording it skipped while the
// cancellation happens somewhere outside the sequence - the fork 246-a1 §2③
// names as the cost of not having a registration seam.
//
// The order inside the step is deliberate:
//  1. every card still awaiting a human is refused through the queue's own
//     refusal funnel, which is the same route a native 「拒绝」 takes and writes
//     the same ANSWER-REJECT line - an unanswered question is never allowed by
//     leaving;
//  2. the task root is cancelled, which is what an L1 window (not a queue item)
//     answers with, and each such abandonment is booked by this file because the
//     gate writes no line for that branch;
//  3. the in-flight asks are waited for, bounded by the ctx this hook was handed
//     (shutdown.go:83 - exceeding it abandons and records, it does not hang).
func (ra *residentApproval) cancelTaskRoots(ctx context.Context) error {
	ra.mu.Lock()
	ra.closed = true
	ra.mu.Unlock()

	refused, abandoned, failed := 0, 0, 0
	for _, card := range ra.cards.Pending() {
		switch card.Level {
		case "L2":
			if err := ra.cards.Reject(ctx, card.CorrelationID,
				"常驻进程退出：未获批准，按拒绝处理（未执行）"); err != nil {
				failed++
				ra.residentAuditf("resident-approval: 退出序列拒绝对待卡片 %s 失败：%v", card.CorrelationID, err)
				continue
			}
			refused++
		default:
			// An L1 window has no reject verb at all (SPEC-06 §2 B1: a countdown
			// can be opposed, not answered), so the honest record for it is the
			// one this file writes before the root cancel below resolves it.
			abandoned++
			ra.residentAuditf("approval: RESIDENT-WINDOW-ABANDONED corr=%s tool=%s decision=reject "+
				"reason=%q channel=none", card.CorrelationID, card.Tool,
				"常驻进程退出，L1 确认窗口未放行，未执行")
			ra.cards.Forget(card.CorrelationID)
		}
	}

	ra.cancel()

	deadline := ctx.Done()
	for {
		ra.mu.Lock()
		n := ra.live
		ra.mu.Unlock()
		if n == 0 {
			break
		}
		select {
		case <-deadline:
			ra.residentAuditf("resident-approval: 退出序列第 3 步到点仍有 %d 张卡片在等，已弃等并记录（D38e）", n)
			return fmt.Errorf("proc: 退出第 3 步到点，仍有 %d 个等待未收口: %w", n, ctx.Err())
		case <-time.After(10 * time.Millisecond):
		}
	}
	ra.residentAuditf("resident-approval: 退出第 3 步完成：拒绝待批卡片 %d 张、作废 L1 窗口 %d 张、路由失败 %d 张；"+
		"任务根已取消，无等待残留", refused, abandoned, failed)
	return nil
}

// liveConfirmations reports how many confirmations are between "asked" and "answered".
func (ra *residentApproval) liveConfirmations() int {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	return ra.live
}

// detachBall is called before the ball host tears its window down. After this,
// the UI posts nothing to the ball's thread (a stopped STA thread would never
// answer), and the veto channel is unloaded again so a late arrival cannot be
// booked as a cancel on a key this process no longer holds.
func (ra *residentApproval) detachBall() {
	ra.ui.releaseBall()
	ra.mu.Lock()
	was := ra.escLoad
	ra.escLoad = false
	ra.mu.Unlock()
	if was {
		ra.gate.Channels().SetLoaded(approval.ChannelEsc, false)
	}
}

// statusLine is the sentence the boot report prints about this gate. It names
// what exists (a gate, one loaded channel) and what does not (a task source),
// because "imported the approval package" and "runs confirmations" are different
// claims and only the first is true here.
func (ra *residentApproval) residentStatusLine() string {
	if !ra.ui.hasBall() {
		return "审批门未装配（本进程没有可承载卡片的悬浮球窗口）"
	}
	return fmt.Sprintf("审批门已装配进本进程（取消通道：Esc 已加载；等待中的确认项：%d）",
		len(ra.cards.Pending()))
}

// --------------------------------------------------------------------- the UI

// ballCardUI is the injected presentation surface for the resident leg: the orb's
// state and the Esc borrow are its output, and it decides nothing. It mirrors
// cmd/wisp/run.go's consoleApprovalUI in the two things a native surface owes the
// answer side: it books the card's single-use grant into the ledger (only the UI
// is handed that value), and it forgets it the moment the card stops being
// answerable.
type ballCardUI struct {
	mu sync.Mutex
	// b is nil until bindBallHost runs and again after detach; every use takes a
	// local copy under the lock so no method holds the mutex across a Win32 post.
	b     *ball.Ball
	ra    *residentApproval
	shown int
}

func (u *ballCardUI) attachBall(b *ball.Ball) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.b = b
}

func (u *ballCardUI) releaseBall() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.b = nil
}

func (u *ballCardUI) hasBall() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.b != nil
}

func (u *ballCardUI) currentBall() *ball.Ball {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.b
}

func (u *ballCardUI) displayedCards() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.shown
}

// Prompt puts one confirmation on the orb: the D43 state the card itself implies
// (Confirming for an L1 window, AwaitingApproval for an L2 one - the two names
// Replies.WaitingState produces, never a fifth), and for the L1 window only, the
// borrow of the cancel key that ticket 245 scoped to exactly that window.
//
// An L2 card must NOT borrow Esc: the C18 deadline is 300s, and holding a
// desktop-wide key for that long would be strictly worse than the standby bug
// ticket 245 removed. An L2 card is refused by the exit path (see cancelTasks)
// or by its own deadline; the veto machinery still routes an Esc naming it into
// the queue's refusal funnel (approval/gate.go's ticket-87 branch), which is the
// gate's behaviour, not this file's.
func (u *ballCardUI) Prompt(_ context.Context, p approval.Prompt) error {
	b := u.currentBall()
	if b == nil {
		slog.Error("approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝",
			"corr", p.CorrelationID, "tool", p.Tool, "level", p.Level)
		return errResidentNoBall
	}
	u.mu.Lock()
	u.shown++
	u.mu.Unlock()

	u.ra.cards.Record(p)
	if p.Level == "L1" {
		b.TakeEscForCancel()
	}
	b.SetState(stateForCardLevel(p.Level))
	slog.Info("approval: 常驻进程显示一张确认卡片",
		"corr", p.CorrelationID, "level", p.Level, "tool", p.Tool,
		"orb_state", string(stateForCardLevel(p.Level)),
		"esc_borrowed", p.Level == "L1",
		"channels", channelRosterText(p.Channels))
	fmt.Printf("wisp: 卡片挂起：%s %s（编号 %s）\n", p.Level, p.Tool, p.CorrelationID)

	// A card whose borrow Win32 refused is a card the user cannot cancel. The
	// report is the authority on that (internal/ball books the refusal there),
	// so the sentence is said now rather than discovered at the deadline.
	if p.Level == "L1" && !b.EscTakenOver() {
		slog.Warn("approval: L1 窗口挂起期间裸 Esc 未借到，本张卡片无法用 Esc 否决",
			"corr", p.CorrelationID, "why", "取消键位被占用或注册被拒，见热键报告")
		fmt.Printf("wisp: 卡片 %s 的取消键未借到（桌面已有占位者），按 Esc 不会否决它\n", p.CorrelationID)
	}
	return nil
}

// Update prints the transient states and ends answerability the way the console
// surface does: dismissed and started both mean the card can no longer be
// answered, so the ledger must not keep a spendable grant for it.
func (u *ballCardUI) Update(_ context.Context, e approval.Event) error {
	switch e.Kind {
	case approval.EventDismissed, approval.EventStarted:
		u.ra.cards.Forget(e.CorrelationID)
		u.settleOrb()
	case approval.EventWarning:
		slog.Warn("approval: 常驻进程卡片进入醒目提示", "corr", e.CorrelationID, "text", e.Text)
	}
	return nil
}

// settle returns the orb and the desktop to the state this leg may claim once no
// card is being waited on: Esc handed back (ticket 245's "会话结束必须归还") and
// the rendered state back to Sleeping - the only state a process with no
// listening, thinking or acting path may show.
//
// It re-reads the ledger before touching anything, so with two cards live the
// first dismissal does not cancel the second one's borrow.
func (u *ballCardUI) settleOrb() {
	if _, awaiting := u.ra.cards.AwaitingHuman(); awaiting {
		return
	}
	b := u.currentBall()
	if b == nil {
		return
	}
	b.ReleaseEscAfterSession()
	b.SetState(statemachine.StateSleeping)
}

// stateForCardLevel is the D43 name for "a human is being waited on at this
// level". It is the same pairing replies.go's WaitingState produces from the same
// two routes (StateConfirming for an L1 window, StateAwaitingApproval for a queue
// card) - this file does not pick a fifth name, and D43's row for a veto's
// TARGET (Acting) is deliberately not claimed here: this leg has no loop to go
// back to, and pretending it does is the lie ticket 228's header warns about.
func stateForCardLevel(level string) statemachine.State {
	if level == "L1" {
		return statemachine.StateConfirming
	}
	return statemachine.StateAwaitingApproval
}

// channelRosterText renders the four channels' availability as one log line, so
// the record shows B1's honest limit rather than implying all four can fire.
func channelRosterText(statuses []approval.ChannelStatus) string {
	out := ""
	for _, s := range statuses {
		if out != "" {
			out += " | "
		}
		out += string(s.Channel) + "="
		if s.Loaded {
			out += "loaded"
		} else {
			out += "unloaded"
		}
	}
	return out
}
