//go:build windows

package main

// The floating ball, hosted by the leg that stays alive.
//
// Why this file exists (ticket 228 AC#1). Before it, the only production
// importer of internal/ball was cmd/balldebug - a debug harness. The reading
// was one ruler, run twice by two different agents and agreeing:
//
//	grep -rln "wisp/internal/ball" --include=*.go . | grep -v .scratch
//	=> ./cmd/balldebug/main.go
//
// so the ball, the tray and the four global hotkeys existed only in a process
// nobody launches to get work done, while the two real legs (`wisp run`, and
// this no-args resident process) showed the user console text. D2 (PLAN.md:74,
// :83-88, :103) draws the layered window, the tray, the hotkeys and the Job
// Object holder as ONE resident main process; orchestrator ruling 228 AC#0
// picked that leg ("走甲"), and this file is that ruling's first installment.
//
// What this installment deliberately is NOT. The resident leg has no task
// pipeline and no microphone, so this file starts no state machine of its own and
// answers no card on the user's behalf: every gesture the ball surfaces except one
// is recorded by name and said out loud as unhosted.
// Advancing the orb into Listening because a click arrived would be this process
// claiming a listen path it does not have - the same shape that made the run
// leg's empty runSpec.replyVeto a load-bearing choice rather than a missing
// assignment (cmd/wisp/run.go:141-147, ledger A472).
//
// What changed since (ticket 246, the sentence above is the ticket-228 wording
// and it is kept because it is still most of the truth): the assembly root now
// injects an approval gate into this leg (cmd/wisp/resident_approval_windows.go),
// so ONE of D43's four veto channels - the Esc cancel hot key - has an executor
// here. It is handed in as a function value by startResidentBall's parameter, so
// this file still knows nothing about approvals, cards or channels: that
// knowledge lives in the assembly root, which is the shape ticket 246's ruling
// (ledger A481, form 乙) asked for. The ball-click veto, the KWS veto word and
// the panel's reject button remain unhosted here, each with its own ticket.
//
// Failure posture: a ball that cannot be created is loud and non-fatal, the
// same stance installLogSink takes above. A machine with no desktop (a service
// session, a CI runner) must still boot, still hold its Job Object and still
// leave through the D38(e) order; it just has to say it has no ball - and the
// gate it built then says out loud that it has nowhere to show a card.

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// escVetoFunc is the whole of what the ball host knows about the thing that
// happens when the cancel key is pressed: somebody handed this process a function
// that turns that press into a sentence. nil means no gate was injected (the
// shape every leg had before ticket 246, and the shape a host that cannot build
// one falls back to), and the gesture is then recorded as unhosted.
type escVetoFunc func() string

// ballHostHook is how the assembly root hands the ball host one capability it does
// not know how to build. It follows the shape ticket 246's ruling (ledger A481) set
// for the cancel executor: the resident ball file receives a function value and
// stays ignorant of what is behind it, so this package keeps one dependency edge
// per thing instead of importing the panel host from here.
type ballHostHook func(*panelHostHooks)

// panelHostHooks holds the panel-side executors, all optional. A nil showPanel is
// the pre-ticket-33 shape: the gesture is recorded and said out loud as unhosted.
type panelHostHooks struct {
	showPanel func(via string) bool
}

// withPanelHost hands the ball host the thing that opens the panel window
// (ticket 33). It is a hook and not a parameter so the three existing
// two-argument call sites of startResidentBall keep compiling untouched.
func withPanelHost(showPanel func(via string) bool) ballHostHook {
	return func(h *panelHostHooks) { h.showPanel = showPanel }
}

// requestPanelOpen answers one of the two panel gestures. Without an executor it
// records the gesture the way every other unhosted one is recorded; with one it
// asks, and says out loud when the answer was "the panel thread will not take the
// request" - a hot key that fires into a full or finished queue is exactly the
// silence tickets 245 and 246 were filed about.
func (h *panelHostHooks) requestPanelOpen(via string) func() {
	if h.showPanel == nil {
		return func() { recordBallGesture(via) }
	}
	return func() {
		if !h.showPanel(via) {
			const why = "the panel thread took no request (it is starting up, its queue is full, or it already exited)"
			slog.Warn("panel gesture could not be handed over", "gesture", via, "why", why)
			fmt.Printf("wisp: ball %s: %s\n", via, why)
		}
	}
}

// residentBall is this process's handle on the floating ball: the window and
// its owner's side of the teardown, plus the one sentence the boot report
// prints. It is nil-safe everywhere, so runResident cannot panic on a leg where
// the ball was never created.
type residentBall struct {
	b *ball.Ball
	// cancelHosted records whether the assembly root actually handed this host a
	// cancel executor (ticket 246 AC#3). It is a fact about the wiring, taken in
	// the same statement that installs the event, and the gate reads it before it
	// advertises the Esc channel: a window whose cancel key fires into nothing
	// must not be reported as a leg that can be vetoed by Esc.
	cancelHosted bool
	// verdict is the operator-facing fact about the ball in THIS process,
	// built from what Win32 actually returned rather than from a hope. It is
	// written in the same statement that stores (or does not store) the ball,
	// which is what keeps the boot line from claiming a window the process
	// does not hold: rb.b != nil is that fact's other half.
	verdict string
}

// startResidentBall creates the floating ball window, the tray icon and the
// ball's global hotkeys on the ui-sta thread internal/ball owns (that thread is
// the ball's own Registry.Spawn with the frozen D38b resident name "ui-sta",
// so this function starts no goroutine of its own - D38a's "one STA/UI thread,
// ball and panel share it" law is satisfied by not adding a second one).
//
// How many keys this leg actually holds (ticket 245): THREE at idle. summon /
// mute / panel register at boot; the fourth slot, cancel, is left standby by
// internal/ball because its production binding is a bare Esc and RegisterHotKey
// is desktop-wide.
//
// Since ticket 246 the Confirming borrow is reachable FROM this leg: the
// assembly root injects an approval gate (resident_approval_windows.go), a real
// card puts the orb in Confirming, and that card's Prompt is what calls
// Ball.TakeEscForCancel - so the cancel slot is held for the 2-3 seconds a
// window is open and handed back when it closes. Before that injection the
// borrow existed only in cmd/balldebug, which is the reverse gap ticket 246 was
// filed for. What has NOT changed: a resident leg with no card in flight holds
// three keys, not four, and an L2 queue card does not borrow the key at all.
//
// What changed since (ticket 33, this file's second injection). The panel hot key
// and the tray's open-panel item stop here too: they are handed to the assembly
// root's panel host as a function value (withPanelHost), which posts them to the
// resident panel's OWN STA thread (cmd/wisp/panel_resident_windows.go, orchestrator
// ruling P1). What this file still does not know is anything about WebView2, the
// window, or what a shown panel means - and the two gestures still do not carry an
// approval: opening the panel grants nothing, and the panel side can never answer
// "allow" (internal/agent/approval/ui.go: PanelAPI has no Allow method).
//
// reg is the runtime's registry, not observe.Default, so a boot that overrode
// the registry (internal/proc.WithRegistry) books its UI thread where the rest
// of the process is booked.
//
// onCancelEsc is the injected cancel executor (escVetoFunc). nil is legal and
// means what it always meant: the gesture is recorded as unhosted.
//
// Errors are returned as a verdict, never as a failure to boot.
func startResidentBall(reg *observe.Registry, onCancelEsc escVetoFunc, hooks ...ballHostHook) *residentBall {
	rb := &residentBall{cancelHosted: onCancelEsc != nil}
	hs := panelHostHooks{}
	for _, apply := range hooks {
		apply(&hs)
	}

	b, err := ball.New(ball.Options{
		// Sleeping is the only state this leg may claim on its own: nothing here
		// is listening, thinking, acting or waiting on anyone. A card that IS
		// waiting moves the orb through the injected gate's UI, which returns it
		// here the moment nobody is being waited on.
		Initial:  statemachine.StateSleeping,
		Hotkeys:  ball.DefaultHotkeys(),
		Registry: reg,
		Events: ball.Events{
			OnClickBall:     func() { recordBallGesture("click") },
			OnSummonHotkey:  func() { recordBallGesture("summon-hotkey") },
			OnMuteHotkey:    func() { recordBallGesture("mute-hotkey") },
			OnCancelHotkey:  func() { recordCancelHotkey(onCancelEsc) },
			OnPanelHotkey:   hs.requestPanelOpen("panel-hotkey"),
			OnTrayPanel:     hs.requestPanelOpen("tray-open-panel"),
			OnTrayMute:      func() { recordBallGesture("tray-mute") },
			OnTrayPauseWake: func() { recordBallGesture("tray-pause-wake") },
			OnTrayExit:      recordTrayExit,
			OnDragEnd:       func() { recordBallDragEnd() },
		},
	})
	if err != nil {
		rb.verdict = fmt.Sprintf("this process has NO floating ball window: %v", err)
		slog.Error("ball: the resident leg could not create the floating ball window", "err", err)
		fmt.Printf("wisp: %s\n", rb.verdict)
		return rb
	}

	rb.b = b
	rep := b.HotkeyReport()

	// Every [hotkey] outcome is said, never swallowed. A ball that came up with
	// none of its keys registered looked exactly like a working demo until
	// ticket 64 A1b named it, so the problems print before the summary does.
	// Since ticket 245 the idle set is three of the four slots: cancel is
	// standby, which is NOT a problem line (Problems() skips it on purpose), so
	// a clean boot here prints nothing and still holds no Esc.
	for _, line := range rep.Problems() {
		slog.Error("ball: " + line)
		fmt.Printf("wisp: %s\n", line)
	}

	rb.verdict = fmt.Sprintf("the floating ball window is up in this process (tray icon added, hotkeys live %d/4: %s)",
		len(rep.Live()), hotkeySummary(b))
	slog.Info("ball: the resident leg created the floating ball window",
		"hotkeys_live", len(rep.Live()), "hotkeys", hotkeySummary(b),
		"gestures", "recorded only except the cancel key: this leg has no task pipeline and no microphone, "+
			"and D43's four veto channels are reduced here to the one the assembly root injected (ticket 246)")
	return rb
}

// statusLine is the fact this leg prints about its own ball. It is one string
// built in startResidentBall, which is what keeps the printed sentence and the
// created window from drifting apart.
func (rb *residentBall) statusLine() string {
	if rb == nil {
		return "this process has NO floating ball window (the ball host never ran)"
	}
	return rb.verdict
}

// stop tears the ball down: Ball.Close runs the animation stop, the hotkey
// unregistration, the tray removal, the renderer release and the window destroy
// on the STA thread, then posts WM_QUIT and joins the ui-sta handle
// (internal/ball/ball_windows.go:904). It is the work D38(e) step 2 names -
// "hotkey + wake-word listening stops" (internal/proc/shutdown.go:18) - and
// runResident registers it as the defer that runs BEFORE rt.Shutdown walks the
// frozen sequence, so the hot keys are released ahead of the Job Object close
// in step 9 and the ball's own log lines still reach the installed sink.
//
// Close is idempotent and joins its thread, so after stop returns this process
// holds no ball-side window, USER/GDI/D2D object or registered hotkey.
func (rb *residentBall) stop() {
	if rb == nil || rb.b == nil {
		return
	}
	rb.b.Close()
	rb.b = nil
	slog.Info("ball: resident leg destroyed the ball window, its tray icon and its hotkeys")
}

// ballGestureWhy is the single reason every unhosted gesture gets: said to the
// console, booked in the log file, and identical between the two so neither can
// drift into a claim the other does not make.
//
// It names the two gestures that are NOT unhosted as exceptions, on purpose. Since
// ticket 246 the resident process holds an approval gate (the cancel key), and
// since ticket 33 it holds a panel host (the panel hot key and the tray's open-panel
// item, which reach the dedicated panel thread through panelHostHooks). The
// pre-246 wording ("no approval gate", "no panel host") would otherwise be a
// sentence every remaining gesture tells in order to cover the ones that stopped
// being one.
const ballGestureWhy = "this process has no task pipeline and no microphone, so the gesture has no executor here; " +
	"the gestures that DO have one are the cancel key the assembly root wired (ticket 246) and the two panel " +
	"gestures that reach the resident panel thread (ticket 33)"

// recordBallGesture books one ball gesture that arrived with nowhere to go.
//
// Callbacks run ON the ui-sta thread (Ball.fire calls them inline, and its
// contract is "quick and non-blocking"), so this does no Win32 and waits on
// nothing: one log record and one line.
func recordBallGesture(name string) {
	slog.Warn("ball gesture arrived with no executor", "gesture", name, "why", ballGestureWhy)
	fmt.Printf("wisp: ball %s: %s\n", name, ballGestureWhy)
}

// recordCancelHotkey is the ONE gesture in this leg with an executor: the
// assembly root's injected cancel function decides what a cancel-key press
// means, and this file only says what came back.
//
// A nil executor is the pre-ticket-246 shape and is reported as such rather than
// swallowed - a cancel key that fires into nothing is exactly the defect class
// ticket 245 and 246 were filed about, so it gets a line of its own in both the
// console and the ledger.
func recordCancelHotkey(onCancelEsc escVetoFunc) {
	if onCancelEsc == nil {
		slog.Warn("cancel hotkey fired with no executor", "why",
			"the assembly root injected no approval gate into this leg, so the borrow cannot be spent")
		fmt.Printf("wisp: ball cancel-hotkey: no approval gate was injected into this process, so the press decided nothing\n")
		return
	}
	sentence := onCancelEsc()
	slog.Info("ball: the cancel key was handled by the injected approval gate", "outcome", sentence)
	fmt.Printf("wisp: %s\n", sentence)
}

// recordBallDragEnd says the one thing the drag path owes the user: the orb
// moved, and this build does not remember where. Options.Store is nil, so
// internal/ball persists nothing - which has to be said, not left for the user
// to discover on the next boot.
func recordBallDragEnd() {
	const why = "this leg passes no ball.Options.Store, so the position is not remembered across boots"
	slog.Info("ball: drag ended", "persisted", false, "why", why)
	fmt.Printf("wisp: ball moved: %s\n", why)
}

// recordTrayExit is the tray's 退出 item. The request is recorded; the stop is
// not performed, because the decision to leave this process lives inside
// internal/proc's event loop, which today takes input from a console signal
// only. Giving this item an executor needs either a stop-request hook in
// internal/proc or a seventh live goroutine that is not in the frozen D38b
// resident roster (observe.ResidentNames) - both are new contracts, and ticket
// 228 AC#1 does not authorize either, so it is reported instead of assumed.
func recordTrayExit() {
	const why = "this process leaves when its event loop returns, which today only a console signal (Ctrl+C) asks for; the tray Exit item has no stop path attached to it"
	slog.Warn("tray exit requested", "outcome", "ignored", "why", why)
	fmt.Printf("wisp: tray exit requested: %s\n", why)
}

// hotkeySummary renders the four bindings as one log-friendly line: the
// configured spelling and what Win32 made of it, per slot.
func hotkeySummary(b *ball.Ball) string {
	lines := make([]string, 0, 4)
	for _, bd := range b.HotkeyReport().Bindings() {
		binding := bd.Binding
		if binding == "" {
			binding = "(unset)"
		}
		lines = append(lines, fmt.Sprintf("%s=%s %s", bd.Name, binding, bd.Status))
	}
	return strings.Join(lines, ", ")
}
