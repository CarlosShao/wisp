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
// What changed since (ticket 290, form 甲-1). The two mute gestures - the global
// mute hot key and the tray's mute item - now have an executor in this process:
// they turn the capture gate this same process assembled (the residentAudio
// handle built by startResidentAudio in cmd/wisp/resident_audio_windows.go).
// That closes the half ticket 247 left open: the microphone leg was mounted,
// muted by default (which is correct and unchanged), and nothing in production
// could ever turn it. The gate's own muted flag is the truth source for "who is
// muted" (internal/audio/gate.go:20-21); this host keeps no copy of it - see
// muteGesture, whose sentence is read back off the gate after the write.
// What did NOT change here: no state machine event is dispatched (D43's
// Sleeping -> Muted edge does not exist, and inventing one is a human-approved
// contract change, not a wiring detail), and no second "am I muted" flag.
//
// Failure posture: a ball that cannot be created is loud and non-fatal, the
// same stance installLogSink takes above. A machine with no desktop (a service
// session, a CI runner) must still boot, still hold its Job Object and still
// leave through the D38(e) order; it just has to say it has no ball - and the
// gate it built then says out loud that it has nowhere to show a card.

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

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

// muteGestureFunc is the whole of what the ball host knows about the thing that
// happens when a mute gesture arrives: somebody handed this process a function
// that flips the capture gate and said back what the gate now reports. The
// sentence is the gate's, not this file's - the gate owns the muted flag
// (internal/audio/gate.go), and a second "who is muted" in this host is the
// shape ticket 246's ruling rejected for the approval side.
//
// executed=false means "there was no gate to turn in this process" (voice off,
// or the capture leg refused to assemble), and the caller says that instead of
// reporting a mute nobody asked for.
type muteGestureFunc func() (outcome string, executed bool)

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
	// hotkeyProvenance is ticket 258 AC#1's answer to "where did the four
	// bindings come from": config, defaults, or no host view at all. It is
	// named in the verdict and in the boot log record, which is what makes the
	// AC#1 fallback ("配置缺失 ⇒ DefaultHotkeys") a said sentence rather than a
	// silent swap. Callers who need it for an assertion read hotkeySourceName().
	hotkeyProvenance string
	// hotkeyBridgeArmed records whether the ticket-64 reloader bridge is running
	// over this ball (ticket 258, form A). It is a wiring fact, kept next to
	// cancelHosted for the same reason that field exists: the boot report and
	// the teardown must not claim more or less than was actually installed.
	hotkeyBridgeArmed bool
	// hotkeyBridgeStop cancels and joins the bridge loop. nil = no bridge was
	// installed (no config source, or the ball never came up). stop() runs it
	// BEFORE Close: see the join-order comment inside startResidentBall.
	hotkeyBridgeStop func()
	// hotkeyBridge is the installed reloader, kept for the test seam
	// (resident_hotkey_258_seam_windows.go) that drives one poll hop
	// synchronously instead of waiting on the real ticker. nil = no bridge.
	hotkeyBridge *ballHotkeyBridge258
	// muteGate is the injected gate-side executor for the two mute gestures
	// (ticket 290 AC#2, form 甲-1). It is attached AFTER the window exists
	// because the assembly root builds the ball before it builds the capture leg
	// that owns the gate (cmd/wisp/resident_windows.go: startResidentBall runs
	// first, startResidentAudio after it), so no closure handed to ball.New can
	// name the gate at construction time. nil is legal and is said out loud.
	//
	// muteMux is what makes that late attach visible to the ui-sta thread: the
	// hot key is live from the moment the window comes up, so a plain field
	// would be a race between the boot goroutine writing it and a key press
	// reading it. The attach also orders everything the boot wrote into the
	// residentAudio handle before any gesture may read it.
	muteMux  sync.Mutex
	muteGate muteGestureFunc
	// trayMuteRead and trayCheckPush are ticket 293 AC#2's projection pair for the
	// tray's "静音" checkmark (see attachTrayMuteProjection and mirrorTrayMute at
	// the bottom of this file); they hold no state - gate.Muted() stays the only
	// truth and the push is literally (*ball.Ball).SetTrayChecks.
	trayMuteRead  func() (muted, ownsGate bool)
	trayCheckPush func(muted, pausedWake bool)
}

// ballHotkeyBridge258 carries the reloader handle the host installed. The type
// exists so residentBall's field does not name internal/ball's bridge type
// directly - the host stays free to re-shape the poll without touching the
// struct's other readers.
type ballHotkeyBridge258 struct {
	check258 func()
}

// hotkeyProvenanceConfig / Defaults / None are the three shapes the AC#1
// sentence can take. The names are words on purpose: they print in the verdict
// ("hotkeys from ..."), and the ticket-258 tests read them off the running
// process rather than inferring them from the absence of an error.
const (
	hotkeyProvenanceConfig   = "config"
	hotkeyProvenanceDefaults = "defaults"
	hotkeyProvenanceNone     = "no host config view"
)

// bridgeJoinBudget258 bounds the bridge join inside stop(). The loop's own tick
// is one second, so a stop request is answered on the next tick at the latest;
// the budget exists so a wedged join can never hang the shutdown (D38(e)) and
// gets named instead - the same posture detachBall takes for the approval side.
const bridgeJoinBudget258 = 3 * time.Second

// residentBallHotkeyChain258 is ticket 258 AC#1's construction-time chain: it
// maps the host's [hotkey] view through ball.ApplyHotkeyDefaults and names the
// provenance of the result.
//
// THE RULE 258-r2 FIXED (AC#1 wording; adjudicated by 258-v1 §3 and the
// orchestrator's 翻勾 note): the tier word describes WHAT THE LIVE SET IS, not
// merely whether a file was readable. A readable file without a [hotkey]
// section arrives as three empty slots plus the schema's cancel = "Esc" tag
// (258-a1 census Q1a), and ApplyHotkeyDefaults merges exactly
// ball.DefaultHotkeys() out of it - so it prints "defaults". Claiming
// "config" there named a file that contributed no binding at all.
//
// The three shapes, and the sentence each one owes:
//   - a non-nil hotCfg whose merged set differs from the compiled defaults: at
//     least one slot is the file's, provenance "config"; the audit line in
//     cmd/wisp/resident_windows.go names the per-field picture.
//   - a non-nil hotCfg whose merged set IS the compiled defaults - the file is
//     missing or unreadable (the view answers nothing), or the file exists
//     without a [hotkey] section (only the schema tag arrives): the result is
//     literally ball.DefaultHotkeys(), provenance "defaults". AC#1's 不许静默
//     换成另一套 is a VALUE statement: none of these shapes returns a fourth
//     set; only the word could ever have lied.
//   - nil hotCfg (no host view at all): DefaultHotkeys, provenance "no host
//     config view" - the pre-258 shape, still sayable because a test may start
//     the host without a config source.
func residentBallHotkeyChain258(hotCfg func() ball.HotkeyConfig) (ball.HotkeyConfig, string) {
	if hotCfg == nil {
		return ball.DefaultHotkeys(), hotkeyProvenanceNone
	}
	cfg := ball.ApplyHotkeyDefaults(hotCfg())
	if cfg == ball.DefaultHotkeys() {
		// The file (or its absence) moved no slot: the live set is the compiled
		// table and the verdict owes the word "defaults" (ticket 258-r2).
		return cfg, hotkeyProvenanceDefaults
	}
	return cfg, hotkeyProvenanceConfig
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
// hotCfg is ticket 258's construction-time half of "the [hotkey] section is the
// authority": the chain (residentBallHotkeyChain258 below) turns the assembly
// root's config view into the four bindings ball.New registers. The verdict and
// the slog record both name which of the three provenance shapes produced the
// set, so "the bindings came from config.toml" and "the bindings fell back to
// the compiled defaults" stay two sentences a reader can tell apart (AC#1:
// 不许静默换成另一套).
//
// hotReload is ticket 258's hot-tier half: a function that returns the CURRENT
// effective [hotkey] view (nil = the host has none, and no bridge is installed).
// startResidentBall installs the ticket-64 reloader bridge over it (the same
// machine cmd/balldebug/main.go:237-253 runs), so editing config.toml rebinds
// the live keys inside the same process - hot tier's promise, finally wired to
// the leg that hosts the ball.
//
// Errors are returned as a verdict, never as a failure to boot.
func startResidentBall(reg *observe.Registry, onCancelEsc escVetoFunc, hotCfg func() ball.HotkeyConfig, hotReload func() ball.HotkeyConfig, hooks ...ballHostHook) *residentBall {
	rb := &residentBall{cancelHosted: onCancelEsc != nil}
	hs := panelHostHooks{}
	for _, apply := range hooks {
		apply(&hs)
	}

	// Ticket 258 AC#1: the four bindings come from the config chain, not from a
	// compiled literal. provenance is said in the verdict and in the log record
	// below, so the AC#1 fallback sentence ("配置缺失 ⇒ DefaultHotkeys") is a
	// sentence the process actually says, not an inference.
	cfg, provenance := residentBallHotkeyChain258(hotCfg)
	rb.hotkeyProvenance = provenance

	b, err := ball.New(ball.Options{
		// Sleeping is the only state this leg may claim on its own: nothing here
		// is listening, thinking, acting or waiting on anyone. A card that IS
		// waiting moves the orb through the injected gate's UI, which returns it
		// here the moment nobody is being waited on.
		Initial:  statemachine.StateSleeping,
		Hotkeys:  cfg,
		Registry: reg,
		Events: ball.Events{
			OnClickBall:     func() { recordBallGesture("click") },
			OnSummonHotkey:  func() { recordBallGesture("summon-hotkey") },
			OnMuteHotkey:    func() { rb.muteGesture("mute-hotkey") },
			OnCancelHotkey:  func() { recordCancelHotkey(onCancelEsc) },
			OnPanelHotkey:   hs.requestPanelOpen("panel-hotkey"),
			OnTrayPanel:     hs.requestPanelOpen("tray-open-panel"),
			OnTrayMute:      func() { rb.muteGesture("tray-mute") },
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

	// Ticket 258 (form A, ledger A538): the ticket-64 reloader machine, the one
	// balldebug ran alone until now, wired into the leg that hosts the ball. The
	// poll drives config.Manager.CheckAndReload through Refresh and rebinds on a
	// real diff (hotkey_reload.go diffs first; an unchanged file never touches
	// Win32). The loop goroutine is spawned through the injected registry and
	// joined HERE, before rb.stop()'s Close destroys the window: after WM_QUIT
	// nothing answers a posted rebind, and a Check that arrives at that moment
	// would hang its caller forever (sta_windows.go:240-247 drops the task but
	// uiRun's <-done never fires). Join order is the guard the census asked for
	// (258-a2 §6(b)); the rebind-drop-in-flight-borrow residual stays ticket
	// 245's, not this file's.
	if hotReload != nil && rb.b != nil {
		bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), hotReload)
		bridgeRoot := observe.NewRoot("ball-hotkey-bridge")
		bridgeHandle := reg.Spawn("hotkey-listener", "ball", bridgeRoot, func(ctx context.Context) {
			bridge.Run(ctx, time.Second)
		})
		rb.hotkeyBridgeStop = func() {
			bridgeRoot.Cancel()
			select {
			case <-bridgeHandle.Done():
			case <-time.After(bridgeJoinBudget258):
			}
		}
		rb.hotkeyBridge = &ballHotkeyBridge258{check258: func() { bridge.Check() }}
		// Said once per boot, so the operator can tell "the bridge is armed" from
		// "the ball came up" - two facts, one sentence each, neither swallowed.
		rb.hotkeyBridgeArmed = true
		slog.Info("ball: hotkey reload bridge armed (polls config, rebinds on change)",
			"poll", "1s", "provenance", provenance)
		fmt.Printf("wisp: ball hotkey reload bridge armed (provenance=%s); a hand edit of [hotkey] rebinds the live keys without a restart\n", provenance)
	}

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

	rb.verdict = fmt.Sprintf("the floating ball window is up in this process (tray icon added, hotkeys from %s, hotkeys live %d/4: %s)",
		provenance, len(rep.Live()), hotkeySummary(b))
	slog.Info("ball: the resident leg created the floating ball window",
		"hotkeys_provenance", provenance,
		"hotkeys_live", len(rep.Live()), "hotkeys", hotkeySummary(b),
		"gestures", "recorded only except the ones the assembly root hands an executor for: the cancel key (ticket 246) "+
			"and the panel gestures (ticket 33); the two mute gestures turn this process's capture gate once that leg is "+
			"assembled (ticket 290 - this line prints before the attach, because the ball is built first), and D43's four "+
			"veto channels stay reduced here to the one the assembly root injected")
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
	if rb == nil {
		return
	}
	// The bridge joins BEFORE the window dies (ticket 258): a rebind posted to
	// a WM_QUIT'ed STA thread is a task that never runs and a caller that never
	// wakes (sta_windows.go:240-247). D38(e) step 2's "hotkey listening stops"
	// covers this loop - it is the thing that would re-register keys behind the
	// teardown's back.
	if rb.hotkeyBridgeStop != nil {
		rb.hotkeyBridgeStop()
		rb.hotkeyBridgeStop = nil
	}
	if rb.b == nil {
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
// RE-DERIVED BY TICKET 290, because the sentence it replaced was false twice
// over. The old wording opened with "this process has no task pipeline and no
// microphone", and ticket 247 mounted a real capture leg in this very process
// (cmd/wisp/resident_audio_windows.go), while ticket 290 handed the mute
// gestures their executor. A reason a gesture gives in order to explain itself
// has to be a reason that gesture still has.
//
// What the gestures this file still records as unhosted actually lack, named:
//   - the ball click and the summon key: no executor is handed to THIS host for
//     them. The leg does assemble a task pipeline
//     (cmd/wisp/resident_task_source_windows.go), but it reads the console, and
//     nothing here turns a click on the orb into a task; that hop is not this
//     ticket's to invent.
//   - the tray's pause-wake item: the wake-word listener it would pause lives in
//     internal/speech, which is unwritten, and wake_word.enabled ships false.
//   - the drag end and the tray Exit item state their own reasons
//     (recordBallDragEnd, recordTrayExit) instead of this one.
//
// The gestures that DO have an executor are named conditionally on purpose: each
// one depends on something the assembly root had to inject, so a flat "these
// always work" would become the lie the moment a host starts without it.
const ballGestureWhy = "no executor was handed to this ball host for this gesture; the gestures with one are the " +
	"cancel key when the assembly root injected an approval gate (ticket 246), the panel gestures when it " +
	"injected a panel host (ticket 33) and the two mute gestures when this process assembled a capture leg " +
	"(ticket 290) - what is left is recorded by name, never invented"

// attachMuteGate hands this host the executor that turns a mute gesture into a
// flip of this process's capture gate (ticket 290 AC#2). The assembly root calls
// it once, right after the capture leg exists, and it is the only writer.
func (rb *residentBall) attachMuteGate(fn muteGestureFunc) {
	if rb == nil {
		return
	}
	rb.muteMux.Lock()
	defer rb.muteMux.Unlock()
	rb.muteGate = fn
}

// currentMuteGate reads the attached executor back under the same lock, so the
// ui-sta thread can never see a half-written handoff.
func (rb *residentBall) currentMuteGate() muteGestureFunc {
	if rb == nil {
		return nil
	}
	rb.muteMux.Lock()
	defer rb.muteMux.Unlock()
	return rb.muteGate
}

// attachTrayMuteProjection hands this host the two halves of ticket 293 AC#2's
// tray checkmark projection (form 甲, orchestrator ruling 2026-10-09 20:2x):
// `read` is the capture leg's own answer about the gate it owns, `push` is the
// ball window's display setter. The assembly root calls it once, immediately
// after attachMuteGate, and it is the only writer.
//
// Both halves are optional on purpose: with no ball window the host never
// attaches anything (the guard is at the call site, where taking the method
// value rb.b.SetTrayChecks would itself panic on a nil window), and with no
// capture leg there is no gate to read. Either gap leaves mirrorTrayMute saying
// "nothing was projected" instead of writing a checkmark out of a hope.
func (rb *residentBall) attachTrayMuteProjection(read func() (bool, bool), push func(bool, bool)) {
	if rb == nil {
		return
	}
	rb.muteMux.Lock()
	defer rb.muteMux.Unlock()
	rb.trayMuteRead = read
	rb.trayCheckPush = push
}

// mirrorTrayMute projects this process's gate onto the tray's "静音" checkmark:
// it reads gate.Muted() back off the capture leg - never off a flag this host
// keeps, never off what the gesture hoped for - and hands that single bool to
// the ball's tray display surface. It returns true only when a projection really
// was written.
//
// Threading: callers are the boot goroutine and muteGesture, which runs ON the
// ui-sta thread. Ball.SetTrayChecks only queues onto that thread's own task
// (internal/ball/sta_windows.go's PostTask is documented safe from any goroutine
// and does not wait), so this neither blocks the message pump nor needs a thread
// of its own: no new goroutine (D38b's resident roster stays at its booked six).
//
// The push's SECOND argument is a literal false, and that asymmetry is the point
// rather than an oversight (ticket 293's ruling, 禁区②):
//   - the mute half has a truth source to read (the gate), so leaving it undriven
//     would show "not muted" next to a microphone the user just opened - a lie
//     about privacy, which is what this whole hop exists to stop;
//   - the pause-wake half has NOTHING to read yet: the wake-word listener it
//     would pause lives in internal/speech, which is still only a doc.go, the
//     shipped [voice].wake_word.enabled default is false, and the tray's
//     pause-wake item is still an unhosted gesture booked by name
//     (recordBallGesture, Events.OnTrayPauseWake above). Unchecked is therefore
//     the true reading today, not a suppressed one; inventing a flag to make the
//     two arguments look symmetric would create exactly the second authority
//     this ticket forbids.
func (rb *residentBall) mirrorTrayMute() bool {
	if rb == nil {
		return false
	}
	rb.muteMux.Lock()
	read, push := rb.trayMuteRead, rb.trayCheckPush
	rb.muteMux.Unlock()
	if read == nil || push == nil {
		return false
	}
	muted, ownsGate := read()
	if !ownsGate {
		// This process owns no gate (voice off, or a leg that never assembled),
		// so there is no mute state to show and nothing is written.
		return false
	}
	push(muted, false)
	return true
}

// muteGesture is the one shape both mute gestures take: the global mute hot key
// (internal/ball's hkMute branch) and the tray's mute menu item both end up here,
// and both get the same sentence the gate answered with.
//
// It asks for the OPPOSITE of what the gate currently reports and reads the
// result back, so the printed state is the gate's, never a hope and never a
// mirrored flag (internal/audio/gate.go:20-21 makes the gate's muted the truth
// source; ticket 246's ruling is the precedent for not keeping a second copy).
//
// Since ticket 293 the gate's flag is also PROJECTED outward after a successful
// flip (mirrorTrayMute, at the bottom of this function) so the tray's "静音"
// checkmark says what the gate now is. That projection is one-way and never read
// back here: this sentence, like the outcome it prints, is still taken off the
// gate alone.
//
// Threading, because it is not free: callbacks run ON the ui-sta thread and
// Ball.fire's contract is "quick and non-blocking". Turning the gate on is
// synchronous inside internal/audio - HalfDuplexGate.SetMuted starts the inner
// source, WASAPIMicrophone.Start waits for the pinned capture thread's single
// device open, and that open is one enumeration plus one activation, not a retry
// loop. So the key can hold the message pump for the length of a device open,
// which is the same bounded COM work the boot goroutine already does in
// assembleCapture, and the alternative (a goroutine per gesture) is forbidden by
// D38b's resident roster. Named here rather than left for the next reader.
//
// It returns the sentence it said. ball.Events wants func(), and the two closures
// that call this discard the result on purpose; the return value exists so a
// caller that hands a key press to this host can read what the gate answered
// instead of scraping the console.
func (rb *residentBall) muteGesture(name string) string {
	fn := rb.currentMuteGate()
	if fn == nil {
		// Only reachable inside the boot window: the window and its hot key exist
		// before the capture leg does. Say which window that is, do not pretend
		// the key did something.
		const why = "no capture leg had been assembled when this key arrived, so there was no gate to turn; " +
			"the assembly root attaches the mute key to the gate at the end of boot (ticket 290)"
		slog.Warn("mute gesture arrived before its executor was attached", "gesture", name, "why", why)
		fmt.Printf("wisp: ball %s: %s\n", name, why)
		return why
	}
	outcome, executed := fn()
	if !executed {
		// The capture leg answered "this process owns no gate" and named why
		// (voice off, or the leg never assembled). The gesture is not a mute.
		slog.Warn("mute gesture found no gate to turn", "gesture", name, "why", outcome)
		fmt.Printf("wisp: ball %s: %s\n", name, outcome)
		return outcome
	}
	// Ticket 293 AC#2: the gate changed hands on the line above, so the tray's
	// "静音" checkmark is re-projected from the gate itself before this returns -
	// both start points a user has (the tray item and the global hot key) come
	// through this one exit, which is why the mirror lives here and nowhere else.
	rb.mirrorTrayMute()
	slog.Info("mute gesture turned this process's capture gate", "gesture", name, "outcome", outcome)
	fmt.Printf("wisp: ball %s: %s\n", name, outcome)
	return outcome
}

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

// hotkeySourceName reports where this host's four bindings came from (ticket
// 258 AC#1): "config", "defaults" (the named AC#1 fallback) or "no host config
// view". Empty on a host whose ball never came up - the verdict already said
// why, and a no-ball host has no bindings to attribute.
func (rb *residentBall) hotkeySourceName() string {
	if rb == nil {
		return ""
	}
	return rb.hotkeyProvenance
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
