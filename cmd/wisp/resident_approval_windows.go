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
	"path/filepath"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
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

	// riskWindow / riskTimeout / riskProvenance are ticket 256's [risk] receipt:
	// what this gate was ACTUALLY built with, and where the numbers came from.
	// They are the construction-time answer only - see the limitation spelled
	// out on newResidentApprovalWithConfig.
	riskWindow     time.Duration
	riskTimeout    time.Duration
	riskProvenance string

	mu   sync.Mutex
	root context.Context
	// cancel is the D38(e) step 3 action: it is what an L1 window and an L2
	// wait answer with when the process is leaving.
	cancel  context.CancelFunc
	live    int // raises between "asked" and "answered"
	closed  bool
	escLoad bool // the Esc channel was loaded (only ever with a real ball)

	// cancelKeyRead is the test seam of the two seeded-wording judges this leg
	// owes (ticket 260-r3, ledger A595): it stands in for the ball's report line
	// so a case can name a NON-default cancel key without opening a real window.
	// A real window is a winlive rig and this leg is forbidden to run one (see
	// .scratch/wisp/probes/260/r3/wording.md §6), so the ball-side half of the
	// chain is asserted by reading, not by execution - and that limit is stated
	// there rather than hidden here. nil is the production reader, which is what
	// the shipped process runs: this field has no production writer.
	cancelKeyRead func() string
}

// ---------------------------------------------------------------- the cancel key
//
// residentCancelKeySpelling is the ONE source for "which key does this process
// call the cancel key", and cancelKeySpelling is the only way this file reaches
// it. The value is read off the ball's own hotkey report: the cancel line is the
// line TakeEscForCancel rewrites from the exact cancelBorrow it hands to
// RegisterHotKey (internal/ball/hotkey_windows.go's cancelBorrowedLineFor), and
// the idle pass fills it with the configured binding (cancelIdleLine). Printing
// from that line is what ticket 260 AC#1 condition 2 ("说到做到") asks for;
// composing a combination in this package is the shape the ticket was filed to
// remove, so nothing here parses a binding, spells a modifier, or re-decides
// which value wins.
//
// The single fallback is the ball's own default constant, not a new literal: it
// covers the two shapes where the report has no cancel binding to show - no ball
// window at all, and an empty [hotkey] cancel, which is exactly the input
// resolveCancelBorrow maps to that same default. So the printed name and the
// borrowed key still agree in both.
//
// It is safe to call from the cancel callback (vetoByEsc runs on the ball's UI
// thread): HotkeyReport goes through Ball.uiRun, which runs inline when the
// caller already IS that thread (internal/ball/ball_windows.go:741-745) instead
// of posting and waiting on itself.
func residentCancelKeySpelling(b *ball.Ball) string {
	if b != nil {
		for _, bd := range b.HotkeyReport().Bindings() {
			if bd.Name == "cancel" && bd.Binding != "" {
				return bd.Binding
			}
		}
	}
	return ball.DefaultHotkeys().Cancel
}

// cancelKeySpelling answers for this gate: the injected reader when a case seeded
// one, otherwise the live ball the UI is holding.
func (ra *residentApproval) cancelKeySpelling() string {
	if ra.cancelKeyRead != nil {
		return ra.cancelKeyRead()
	}
	return residentCancelKeySpelling(ra.ui.currentBall())
}

// newResidentApproval composes the gate the resident leg runs on WITHOUT a host
// config view. It is the pre-ticket-256 shape and it is kept for the 246 test
// roster (orchestrator ruling 256 §8.2, form 甲ⓐ: the signature does not gain a
// parameter, because that would put 13 existing test call sites in compile-red).
//
// What it now owes, and says: the gate it builds runs on the compiled approval
// constants, and that has to be NAMEABLE rather than inferred from silence. This
// is the same debt ticket 258 took out on the hotkey chain
// (resident_ball_windows.go's hotkeyProvenanceNone for a nil host view), and the
// value is that a reader of the log can tell "the file said 300" from "nobody
// asked the file" - both land on 300s, and only one of them is a config reading.
func newResidentApproval() *residentApproval {
	return newResidentApprovalWithConfig("")
}

// newResidentApprovalWithConfig composes the same gate, fed with the host's
// [risk] pair. It is what the shipped resident process calls
// (resident_windows.go, the one production site).
//
// WHAT THIS CLOSES (ticket 256 AC#1, the [risk] half only). Before this leg the
// gate was built from three Options fields - UI / Channels / Logf - so
// confirm_timeout_sec and l1_window_sec changed nothing in this process: it ran
// the compiled 300s / 3s no matter what config.toml said. Two of the ten fields
// are now passed, and the two are the only two the assembly root can obtain at
// this moment.
//
// WHAT THIS DELIBERATELY DOES NOT CLOSE, NAMED SO IT IS NOT MISTAKEN FOR DONE:
// Options.Grants stays unset here. The session ledger's only production
// construction site is cmd/wisp/run.go inside assembleRuntime, which this file's
// gate is built roughly 129 lines BEFORE (256-a2 census §0/§2), and g.grants has
// exactly one writer in the repository (gate.go's New) with no late-binding
// entry. Making that half move means either a second minted ledger or a new
// holder type - both new seams, not a field to fill in. That half is filed as
// ticket 255's sibling (pending ticket 265) and is NOT this leg's work.
//
// THE LIMITATION THAT MUST TRAVEL WITH THIS FUNCTION (census §3, pinned by
// TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly): approval.Gate
// copies the window into g.window and the queue copies the deadline into
// q.timeout inside New/NewQueue, and nothing in the repository re-applies either
// afterwards - the only read faces are Gate.Window() and Queue.Timeout(). So
// what this closes is "the gate THIS launch of the resident process is built
// with follows [risk]". It is not "edit config.toml and the running process
// follows you", and no assertion here may be written as if it were. [risk] is a
// locked tier (internal/config/tiers.go) with no reload or restart hook, which
// is the same direction the facts point.
//
// WHY A PER-USE LoadFile AND NOT THE MANAGER: exactly the reason
// resident_windows.go's hotkey closures give (its :168-171 comment, ticket 258) -
// the task pipeline that owns rt.mgr is assembled later, by
// startResidentTaskSource, so at this point in the boot nobody holds a Manager.
// Importing internal/config here is the shape those two closures, models.go:184
// and panel_resident_windows.go:201 already have; it is not a new package-level
// dependency edge of the kind ticket 238's cut-1 refuses, because cmd/wisp
// already imports internal/config in production code.
func newResidentApprovalWithConfig(dataDir string) *residentApproval {
	root, cancel := context.WithCancel(context.Background())
	ra := &residentApproval{root: root, cancel: cancel}
	ra.ui = &ballCardUI{ra: ra}
	window, timeout, provenance := residentRiskGateValues(dataDir)
	ra.riskWindow, ra.riskTimeout, ra.riskProvenance = window, timeout, provenance
	ra.gate = approval.New(approval.Options{
		UI:              ra.ui,
		Channels:        approval.NewChannels(),
		Window:          window,
		ApprovalTimeout: timeout,
		Logf:            ra.residentAuditf,
	})
	// The receipt is read OFF THE GATE, not off the arithmetic above: what is
	// logged is what the object actually holds after New's own clamping, so the
	// sentence cannot disagree with the value it names.
	slog.Info("resident gate: [risk] tier taken at construction",
		"provenance", provenance,
		"config_path", residentRiskConfigPathForLog(dataDir),
		"window_sec_read", window.Seconds(), "confirm_timeout_sec_read", timeout.Seconds(),
		"gate_window", ra.gate.Window().String(), "gate_queue_timeout", ra.gate.Queue().Timeout().String(),
		"scope", "construction time only: approval.Gate and Queue copy these in New and nothing re-applies them (ticket 256)")
	ra.cards = approval.NewReplies()
	ra.cards.Attach(approval.HostBinding{
		Gate:         ra.gate,
		VetoChannel:  approval.ChannelEsc,
		NativeSource: residentNativeSource,
		PanelSource:  residentPanelSource,
	})
	return ra
}

// The three provenance words this leg's [risk] reading can answer with. Named
// constants (not inline strings) for the same reason ticket 258 named its
// hotkeyProvenance* trio: the tests read them off the constructed object, and a
// prose claim that a value "came from config" is exactly the thing that must not
// be writable twice with two different meanings.
const (
	// riskProvenanceRead: config.toml was read, and these two numbers are its
	// [risk] view. Note what the word does NOT claim: a readable file with no
	// [risk] section still answers through the schema tags (confirm_timeout_sec
	// 300 / l1_window_sec 2), so "config" names the SOURCE OF THE READ, not a
	// promise that a human typed a number. That boundary is stated here rather
	// than left to the reader because ticket 258's adjudication (258-r2, "the
	// tier word describes WHAT THE LIVE SET IS") is the nearest precedent and a
	// value-for-value copy of it would have been wrong: the schema's 2s and the
	// gate's compiled 3s differ, so a file that says nothing still moves this
	// gate's window, and calling that "defaults" would be the lie in the other
	// direction.
	riskProvenanceRead = "config"
	// riskProvenanceUnreadable: a dataDir was handed over but config.toml would
	// not read (missing or broken). The gate falls back to the compiled approval
	// constants, and the word "defaults" is what says so out loud.
	riskProvenanceUnreadable = "defaults (config.toml unreadable)"
	// riskProvenanceNoView: no host config view exists at all (dataDir empty -
	// the shape newResidentApproval keeps for the 246 roster). Same compiled
	// constants, different reason, and the reason is the part a reader needs.
	riskProvenanceNoView = "defaults (no host config view)"
)

// residentRiskGateValues reads the two [risk] numbers the resident gate is built
// with, and names which of the three shapes the reading is. Zero durations are
// returned for the two fallback shapes on purpose: approval.New's documented
// zero-value fallback (Options.Window -> DefaultL1Window, and NewQueue's
// timeout <= 0 -> DefaultApprovalTimeout) then becomes the mechanism, so this
// file cannot drift its own private copy of those constants.
func residentRiskGateValues(dataDir string) (window, timeout time.Duration, provenance string) {
	if dataDir == "" {
		return 0, 0, riskProvenanceNoView
	}
	cfgPath := filepath.Join(dataDir, configFileName)
	c, _, err := config.LoadFile(cfgPath, nil)
	if err != nil || c == nil {
		slog.Warn("resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants",
			"path", cfgPath, "err", err,
			"fallback", "DefaultApprovalTimeout=300s / DefaultL1Window=3s")
		return 0, 0, riskProvenanceUnreadable
	}
	return time.Duration(c.Risk.L1WindowSec) * time.Second,
		time.Duration(c.Risk.ConfirmTimeoutSec) * time.Second,
		riskProvenanceRead
}

// residentRiskConfigPathForLog keeps the log honest about WHERE it looked: with
// no host view there is no path, and printing an empty path would read like a
// failed read rather than the absence of one.
func residentRiskConfigPathForLog(dataDir string) string {
	if dataDir == "" {
		return "(no host config view)"
	}
	return filepath.Join(dataDir, configFileName)
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
		// 档一 (A595 §1 boundary ④): this branch has NO key to name - no window
		// means nothing can be borrowed - so the sentence says the slot is empty
		// instead of claiming a key. It used to read 「Esc 无处可借」.
		ra.residentAuditf("resident-approval: 审批门已装配，但本进程没有悬浮球窗口，" +
			"四条否决通道全部保持未加载（取消键无处可借，卡片无处可呈）")
		return false
	}
	if !rb.cancelHosted {
		// A window with no executor behind the cancel key: the borrow would be
		// advertised and the press would be booked as "no executor". Both halves
		// of that are lies the channel registry is not allowed to tell.
		// 档一 for the same reason: the channel is not loaded, so no key is named.
		ra.residentAuditf("resident-approval: 本进程有悬浮球窗口，但装配根没有注入取消执行者，" +
			"取消键通道保持未加载（advertise 一枚按不动的键＝B1 禁止的形状）")
		return false
	}
	ra.mu.Lock()
	ra.escLoad = true
	ra.mu.Unlock()
	ra.ui.attachBall(rb.b)
	ra.gate.Channels().SetLoaded(approval.ChannelEsc, true)
	// 档二: the channel IS loaded here, so the sentence names a key - and it names
	// the one this ball will actually borrow (its report's cancel line, filled from
	// [hotkey] cancel). With the default 档 the rendered line is byte-for-byte the
	// sentence this file printed before ticket 260-r3, which is the zero-drift half
	// of the acceptance pair.
	key := ra.cancelKeySpelling()
	ra.residentAuditf("resident-approval: 审批门已装配进常驻进程，取消通道 %s 已加载（本票只落 %s 一条通道；"+
		"单击球 / KWS 否决词 / 面板拒绝三条仍按各自归口未接入）", key, key)
	return true
}

// vetoRejectedLine and vetoDoneLine are vetoByEsc's two answers, lifted out
// verbatim so the wording pair this ticket owes (ledger A595 §2) has an executable
// ruler. The reason is measurement, not design: the success sentence needs a live
// L1 window on the gate, which only a real ball window can raise, and a real
// window is a winlive rig this leg is forbidden to run (see
// .scratch/wisp/probes/260/r3/wording.md §6). Nothing here changes behaviour -
// vetoByEsc still returns exactly these two strings with exactly these arguments,
// and the key it names still comes from the single source above.
func vetoRejectedLine(key, corr string, err error) string {
	return fmt.Sprintf("按 %s 否决卡片 %s 被拒：%v", key, corr, err)
}

func vetoDoneLine(key string, card approval.ReplyCard) string {
	return fmt.Sprintf("按 %s 否决了卡片 %s（%s / %s）：该调用未执行，答案已入审计",
		key, card.CorrelationID, card.Level, card.Tool)
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
		// 档二: the key named here is the one this card's borrow used, read off the
		// ball's report (vetoByEsc runs on the UI thread, where that read is inline
		// - see residentCancelKeySpelling's note). Default 档 renders the old
		// sentence unchanged.
		return vetoRejectedLine(ra.cancelKeySpelling(), card.CorrelationID, err)
	}
	return vetoDoneLine(ra.cancelKeySpelling(), card)
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
	ra.mu.Lock()
	loaded := ra.escLoad
	ra.mu.Unlock()
	if !loaded {
		return "审批门未装配（本进程没有可承载卡片的悬浮球窗口，或装配根没有注入取消执行者）"
	}
	return fmt.Sprintf("审批门已装配进本进程（取消通道：%s 已加载；等待中的确认项：%d）",
		ra.cancelKeySpelling(), len(ra.cards.Pending()))
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
		// 档三: the borrow WAS attempted and Win32 refused it, so there is a key to
		// name - the one this card tried to borrow, which the ball left on its
		// report's cancel line (cancelFailedLineFor carries exactly the accelerator
		// it handed to RegisterHotKey). The hard-coded 「Esc」 here used to promise a
		// key that may never have been attempted, so the name now comes from that
		// same report line, inside the sentence, like the line below it.
		//
		// ONE WORD HAD TO LEAVE: the old message read 「裸 Esc」, and 「裸」 (bare,
		// no modifiers) is a claim about the key - true of the shipped default and
		// false of any combination the user moves [hotkey] cancel to. Keeping it
		// would have planted a second sentence of exactly the class this ticket was
		// filed to remove, so the default 档 renders this line one word shorter than
		// before. That single-word drift is the whole deviation from AC#4 ① and it
		// is written up in .scratch/wisp/probes/260/r3/wording.md §5 rather than
		// passed off as zero drift. The other shape (key name in a slog attribute,
		// message kept constant) drifts the sentence much further; this call is the
		// tree's first slog message built with fmt.Sprintf - 尺＝grep -rn
		// "slog\.[A-Za-z]*(fmt.Sprintf" --include=*.go internal cmd ⇒ 0 枚先例 - and
		// it is behaviour-neutral, which is what buys the shorter drift.
		key := u.ra.cancelKeySpelling()
		slog.Warn(fmt.Sprintf("approval: L1 窗口挂起期间取消键 %s 未借到，本张卡片无法用 %s 否决", key, key),
			"corr", p.CorrelationID, "why", "取消键位被占用或注册被拒，见热键报告")
		fmt.Printf("wisp: 卡片 %s 的取消键未借到（桌面已有占位者），按 %s 不会否决它\n", p.CorrelationID, key)
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
