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
// pipeline, no microphone, no approval gate and no panel host today, so this
// file starts no state machine and answers no card: every gesture the ball
// surfaces is recorded by name and said out loud as unhosted. Advancing the
// orb into Listening because a click arrived would be this process claiming a
// listen path it does not have - the same shape that made the run leg's empty
// runSpec.replyVeto a load-bearing choice rather than a missing assignment
// (cmd/wisp/run.go:141-147, ledger A472). Making those gestures real is
// ticket 228 AC#2/AC#3 and the approval-gate legs, not this one.
//
// Failure posture: a ball that cannot be created is loud and non-fatal, the
// same stance installLogSink takes above. A machine with no desktop (a service
// session, a CI runner) must still boot, still hold its Job Object and still
// leave through the D38(e) order; it just has to say it has no ball.

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// residentBall is this process's handle on the floating ball: the window and
// its owner's side of the teardown, plus the one sentence the boot report
// prints. It is nil-safe everywhere, so runResident cannot panic on a leg where
// the ball was never created.
type residentBall struct {
	b *ball.Ball
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
// How many keys this leg actually holds (ticket 245): THREE. summon / mute /
// panel register at boot; the fourth slot, cancel, is left standby by
// internal/ball because its production binding is a bare Esc and RegisterHotKey
// is desktop-wide. Confirming borrows that slot for the length of a card, and
// this leg has no card path (no approval gate, no state machine), so it never
// borrows it - which is the point of the ticket: the resident process must not
// take the user's Esc key while it is just sitting there.
//
// reg is the runtime's registry, not observe.Default, so a boot that overrode
// the registry (internal/proc.WithRegistry) books its UI thread where the rest
// of the process is booked.
//
// Errors are returned as a verdict, never as a failure to boot.
func startResidentBall(reg *observe.Registry) *residentBall {
	rb := &residentBall{}

	b, err := ball.New(ball.Options{
		// Sleeping is the only state this leg may claim: nothing here is
		// listening, thinking, acting or waiting on anyone.
		Initial:  statemachine.StateSleeping,
		Hotkeys:  ball.DefaultHotkeys(),
		Registry: reg,
		Events: ball.Events{
			OnClickBall:     func() { recordBallGesture("click") },
			OnSummonHotkey:  func() { recordBallGesture("summon-hotkey") },
			OnMuteHotkey:    func() { recordBallGesture("mute-hotkey") },
			OnCancelHotkey:  func() { recordBallGesture("cancel-hotkey") },
			OnPanelHotkey:   func() { recordBallGesture("panel-hotkey") },
			OnTrayPanel:     func() { recordBallGesture("tray-open-panel") },
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
		"gestures", "recorded only: this leg has no task pipeline, no microphone and no approval gate")
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
const ballGestureWhy = "this process has no task pipeline, no microphone, no approval gate and no panel host, so the gesture has no executor here"

// recordBallGesture books one ball gesture that arrived with nowhere to go.
//
// Callbacks run ON the ui-sta thread (Ball.fire calls them inline, and its
// contract is "quick and non-blocking"), so this does no Win32 and waits on
// nothing: one log record and one line.
func recordBallGesture(name string) {
	slog.Warn("ball gesture arrived with no executor", "gesture", name, "why", ballGestureWhy)
	fmt.Printf("wisp: ball %s: %s\n", name, ballGestureWhy)
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
